import { store } from '../store';
import {
  setConnectionStatus,
  setConnectionId,
  incrementReconnectAttempts,
  resetReconnectAttempts,
  setConnectionError,
  updateLatency,
  setHealthStatus,
  cleanup,
  ConnectionStatus,
  setPollingStatus,
  updatePollTime,
  setPollInterval,
  incrementErrorCount,
  resetErrorCount,
  setLastSuccessfulPoll,
} from '../store/slices/connectionSlice';
import {
  addMessage,
  addOptimisticMessage,
  updateMessageStatus,
  removeFailedMessage,
  setTyping,
  setError,
  ChatMessage,
  markMessageAsDelivered,
  addMessages,
  queueOfflineMessage,
  clearOfflineQueue,
} from '../store/slices/chatSlice';
import { ConnectionManager, connectionManager, PollingConnectionState } from './ConnectionManager';

// Types for polling chat service
export interface PollingChatMessage {
  id: string;
  type: 'user' | 'assistant' | 'system';
  content: string;
  timestamp: string;
  session_id: string;
  status?: 'sending' | 'sent' | 'delivered' | 'failed';
  metadata?: Record<string, any>;
}

export interface SendMessageRequest {
  content: string;
  session_id: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

export interface SendMessageResponse {
  success: boolean;
  message_id: string;
  content?: string;    // AI response content
  type?: string;       // Message type (assistant)
  error?: string;
}

export interface GetMessagesResponse {
  success: boolean;
  messages: PollingChatMessage[];
  has_more: boolean;
}

export interface ChatRequest {
  message: string;
  session_id?: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

// Enhanced error handling types
export interface RetryConfig {
  maxRetries: number;
  baseDelay: number;
  maxDelay: number;
  backoffMultiplier: number;
  jitterFactor: number;
}

export interface QueuedMessage {
  id: string;
  request: SendMessageRequest;
  timestamp: Date;
  retryCount: number;
  nextRetryAt: Date;
}

export interface NetworkError extends Error {
  code?: string;
  status?: number;
  isNetworkError: boolean;
  isRetryable: boolean;
}

export interface ErrorHandler {
  handleNetworkError(error: NetworkError): void;
  handleServerError(error: NetworkError): void;
  handleTimeoutError(error: NetworkError): void;
  shouldRetry(error: NetworkError, attempt: number): boolean;
}

// Polling strategy configuration
interface PollingStrategy {
  baseInterval: number;        // 3 seconds default
  maxInterval: number;         // 30 seconds max
  backoffMultiplier: number;   // 1.5x on errors
  activeInterval: number;      // 2 seconds when actively chatting
  inactiveInterval: number;    // 10 seconds when idle
}

// Re-export types from ConnectionManager
export type { PollingConnectionState } from './ConnectionManager';

class PollingChatService {
  private baseUrl: string;
  private authToken: string | null = null;
  private pollingIntervalId: NodeJS.Timeout | null = null;
  private connectionManager: ConnectionManager;
  private pollingStrategy: PollingStrategy;
  private messageQueue: QueuedMessage[] = [];
  private lastMessageId: string | null = null;
  private lastActivity: Date = new Date();
  private isUserActive: boolean = true;
  private retryTimeouts: Map<string, NodeJS.Timeout> = new Map();
  private sentMessageIds: Set<string> = new Set(); // Duplicate prevention
  private retryConfig: RetryConfig;
  private errorHandler: ErrorHandler;
  private isOnline: boolean = navigator.onLine;
  
  // Smart polling activity tracking
  private activityTracker = {
    lastTyping: new Date(0),
    lastMessageSent: new Date(0),
    lastPageFocus: new Date(),
    lastMouseMove: new Date(0),
    lastKeyPress: new Date(0),
    isTyping: false,
    typingTimeout: null as NodeJS.Timeout | null,
    focusTimeout: null as NodeJS.Timeout | null,
  };

  // Message caching for optimization
  private messageCache = new Map<string, {
    messages: PollingChatMessage[];
    timestamp: Date;
    etag?: string;
    lastModified?: string;
  }>();
  private cacheExpiryTime = 5 * 60 * 1000; // 5 minutes

  constructor() {
    console.log('[PollingChat] Service instance created', {
      timestamp: new Date().toISOString()
    });

    // Initialize configuration
    this.baseUrl = process.env.REACT_APP_API_URL || 'http://localhost:8061';
    this.authToken = localStorage.getItem('adminToken');

    // Initialize polling strategy
    this.pollingStrategy = {
      baseInterval: 3000,      // 3 seconds
      maxInterval: 30000,      // 30 seconds
      backoffMultiplier: 1.5,
      activeInterval: 2000,    // 2 seconds when active
      inactiveInterval: 10000, // 10 seconds when idle
    };

    // Initialize retry configuration
    this.retryConfig = {
      maxRetries: 4,
      baseDelay: 1000,         // 1 second
      maxDelay: 8000,          // 8 seconds max
      backoffMultiplier: 2,    // Double each time (1s, 2s, 4s, 8s)
      jitterFactor: 0.1,       // ±10% jitter
    };

    // Use singleton connection manager
    this.connectionManager = connectionManager;

    // Initialize error handler
    this.errorHandler = this.createErrorHandler();

    // Bind methods to preserve context
    this.startPolling = this.startPolling.bind(this);
    this.stopPolling = this.stopPolling.bind(this);
    this.sendMessage = this.sendMessage.bind(this);
    this.getMessages = this.getMessages.bind(this);

    // Track user activity for smart polling
    this.setupActivityTracking();
    this.setupAdvancedActivityTracking();

    // Setup online/offline detection
    this.setupNetworkDetection();

    // Process message queue periodically
    this.setupMessageQueueProcessor();

    // Handle page unload
    window.addEventListener('beforeunload', () => {
      console.log('[PollingChat] Page unloading, cleaning up');
      this.stopPolling();
    });
  }

  /**
   * Start polling for messages with configurable intervals
   */
  public startPolling(): void {
    if (this.connectionManager.isPollingActive()) {
      console.log('[PollingChat] Polling already active');
      return;
    }

    console.log('[PollingChat] Starting polling with smart intervals');
    this.connectionManager.startPolling();
    
    // Generate connection ID
    const connectionId = `poll-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
    store.dispatch(setConnectionId(connectionId));
    
    // Update Redux state for polling
    store.dispatch(setPollingStatus(true));
    store.dispatch(setConnectionStatus('polling'));

    // Initialize last message ID from current messages
    const currentMessages = store.getState().chat.messages;
    if (currentMessages.length > 0) {
      this.lastMessageId = currentMessages[currentMessages.length - 1].id;
      console.log('[PollingChat] Initialized with last message ID:', this.lastMessageId);
    }

    // Start the polling loop
    this.createPollingLoop();
  }

  /**
   * Stop polling for messages
   */
  public stopPolling(): void {
    console.log('[PollingChat] Stopping polling');
    
    this.connectionManager.stopPolling();
    
    // Clear polling interval
    if (this.pollingIntervalId) {
      clearTimeout(this.pollingIntervalId);
      this.pollingIntervalId = null;
    }

    // Clear retry timeouts
    this.retryTimeouts.forEach(timeout => clearTimeout(timeout));
    this.retryTimeouts.clear();

    // Update Redux state
    store.dispatch(setPollingStatus(false));
    store.dispatch(cleanup());
  }

  /**
   * Send a chat message with optimistic updates and retry logic
   */
  public async sendMessage(request: ChatRequest): Promise<string> {
    const messageId = `msg-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
    
    // Check for duplicate message
    if (this.isDuplicateMessage(messageId)) {
      console.warn('[PollingChat] Duplicate message detected, skipping:', messageId);
      return messageId;
    }

    // Create optimistic message
    const optimisticMessage: ChatMessage = {
      id: messageId,
      type: 'user',
      content: request.message,
      timestamp: new Date().toISOString(),
      session_id: request.session_id || '',
      status: 'sending',
    };

    // Add optimistic message to store
    store.dispatch(addOptimisticMessage(optimisticMessage));

    // Update last activity for smart polling
    this.trackMessageSendingActivity();

    // Prepare request payload
    const sendRequest: SendMessageRequest = {
      content: request.message,
      session_id: request.session_id || '',
      client_name: request.client_name,
      context: request.context,
      quick_action: request.quick_action,
    };

    // If offline, queue the message
    if (!this.isOnline) {
      console.log('[PollingChat] Offline, queuing message:', messageId);
      const queuedMessage: QueuedMessage = {
        id: messageId,
        request: sendRequest,
        timestamp: new Date(),
        retryCount: 0,
        nextRetryAt: new Date(),
      };
      this.messageQueue.push(queuedMessage);
      
      // Use new Redux action for offline queuing
      store.dispatch(queueOfflineMessage(optimisticMessage));
      store.dispatch(setError('Offline. Message queued for sending when connection is restored.'));
      return messageId;
    }

    try {
      // Attempt to send message with retry logic
      const response = await this.sendMessageWithRetry(sendRequest, messageId);
      
      if (response.success) {
        // Mark as sent to prevent duplicates
        this.markMessageAsSent(messageId);
        
        // Clear message cache since new message was sent
        this.clearMessageCache(request.session_id);
        
        // Update message status to sent and then delivered
        store.dispatch(updateMessageStatus({ id: messageId, status: 'sent' }));
        // Mark as delivered since HTTP is immediate
        setTimeout(() => {
          store.dispatch(markMessageAsDelivered(messageId));
        }, 100);
        
        // Process AI response if present
        if (response.content && response.type) {
          const aiMessage: ChatMessage = {
            id: `ai-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`,
            type: response.type as 'assistant',
            content: response.content,
            timestamp: new Date().toISOString(),
            session_id: request.session_id || '',
            status: 'delivered',
          };
          store.dispatch(addMessage(aiMessage));
          console.log('[PollingChat] AI response added to chat:', aiMessage.id);
        }
        
        console.log('[PollingChat] Message sent successfully:', messageId);
      } else {
        // Queue for retry if retryable error
        const queuedMessage: QueuedMessage = {
          id: messageId,
          request: sendRequest,
          timestamp: new Date(),
          retryCount: 0,
          nextRetryAt: new Date(Date.now() + this.calculateRetryDelay(1)),
        };
        this.messageQueue.push(queuedMessage);
        console.error('[PollingChat] Message send failed, queued for retry:', response.error);
      }
    } catch (error) {
      console.error('[PollingChat] Message send error:', error);
      
      // Create network error and check if retryable
      const networkError = this.createNetworkError(error);
      
      if (this.errorHandler.shouldRetry(networkError, 1)) {
        // Queue for retry
        const queuedMessage: QueuedMessage = {
          id: messageId,
          request: sendRequest,
          timestamp: new Date(),
          retryCount: 0,
          nextRetryAt: new Date(Date.now() + this.calculateRetryDelay(1)),
        };
        this.messageQueue.push(queuedMessage);
        console.log('[PollingChat] Message queued for retry due to network error');
      } else {
        // Not retryable, mark as failed
        store.dispatch(updateMessageStatus({ id: messageId, status: 'failed' }));
        this.errorHandler.handleNetworkError(networkError);
      }
    }

    return messageId;
  }

  /**
   * Get messages from the server using efficient polling with last message ID/timestamp
   * Implements client-side message caching to avoid duplicate requests
   */
  public async getMessages(sessionId: string, lastMessageId?: string): Promise<PollingChatMessage[]> {
    try {
      // Check cache first if no specific lastMessageId (full session request)
      if (!lastMessageId) {
        const cachedMessages = this.getCachedMessages(sessionId);
        if (cachedMessages) {
          return cachedMessages;
        }
      }

      const queryParams = new URLSearchParams();
      queryParams.append('session_id', sessionId);
      
      // Always get all messages for now (simpler approach)
      queryParams.append('limit', '100');

      // Get cache metadata for conditional requests
      const cacheMetadata = !lastMessageId ? this.getCacheMetadata(sessionId) : null;

      const response = await this.makeRequest<GetMessagesResponse>(
        `/api/v1/admin/chat/messages?${queryParams.toString()}`,
        { method: 'GET' },
        cacheMetadata || undefined
      );

      // Handle 304 Not Modified response
      if ((response as any).notModified) {
        const cachedMessages = this.getCachedMessages(sessionId);
        if (cachedMessages) {
          console.log('[PollingChat] Using cached messages due to 304 Not Modified');
          return cachedMessages;
        }
      }

      if (response.success) {
        const messages = response.messages || [];
        console.log(`[PollingChat] Retrieved ${messages.length} messages for session ${sessionId}`);
        
        // Cache messages if this is a full session request (not incremental)
        if (!lastMessageId && messages.length > 0) {
          const cacheMetadata = (response as any)._cacheMetadata;
          this.cacheMessages(
            sessionId, 
            messages, 
            cacheMetadata?.etag, 
            cacheMetadata?.lastModified
          );
        }
        
        return messages;
      } else {
        console.error('[PollingChat] Failed to get messages:', response);
        return [];
      }
    } catch (error) {
      console.error('[PollingChat] Error getting messages:', error);
      throw error; // Re-throw to trigger error handling in polling loop
    }
  }

  /**
   * Get messages using timestamp-based polling (alternative approach)
   */
  public async getMessagesSince(sessionId: string, timestamp?: string): Promise<PollingChatMessage[]> {
    try {
      const queryParams = new URLSearchParams();
      queryParams.append('session_id', sessionId);
      
      if (timestamp) {
        queryParams.append('since_timestamp', timestamp);
      } else {
        // Get messages from last 5 minutes if no timestamp
        const fiveMinutesAgo = new Date(Date.now() - 5 * 60 * 1000).toISOString();
        queryParams.append('since_timestamp', fiveMinutesAgo);
      }

      const response = await this.makeRequest<GetMessagesResponse>(
        `/api/v1/admin/chat/messages?${queryParams.toString()}`,
        { method: 'GET' }
      );

      if (response.success) {
        return response.messages || [];
      } else {
        console.error('[PollingChat] Failed to get messages by timestamp');
        return [];
      }
    } catch (error) {
      console.error('[PollingChat] Error getting messages by timestamp:', error);
      throw error;
    }
  }

  /**
   * Set polling interval
   */
  public setPollingInterval(interval: number): void {
    this.pollingStrategy.baseInterval = interval;
  }

  /**
   * Get current connection status
   */
  public getConnectionStatus(): ConnectionStatus {
    return store.getState().connection.status;
  }

  /**
   * Get detailed connection status information
   */
  public getConnectionStatusInfo() {
    return this.connectionManager.getStatusInfo();
  }

  /**
   * Get user-friendly status message
   */
  public getStatusMessage(): string {
    return this.connectionManager.getStatusMessage();
  }

  /**
   * Check if service is healthy
   */
  public isHealthy(): boolean {
    return this.connectionManager.isHealthy();
  }

  /**
   * Force reconnection attempt
   */
  public forceReconnect(): void {
    this.connectionManager.forceReconnect();
    
    // Restart polling if it was active
    if (this.connectionManager.isPollingActive()) {
      this.createPollingLoop();
    }
  }

  /**
   * Subscribe to connection status changes
   */
  public onStatusChange(callback: (status: any) => void): () => void {
    return this.connectionManager.onStatusChange(callback);
  }

  // Private methods

  /**
   * Create error handler with comprehensive error detection and handling
   */
  private createErrorHandler(): ErrorHandler {
    return {
      handleNetworkError: (error: NetworkError) => {
        console.error('[PollingChat] Network error:', error);
        // Connection manager will handle state updates via recordPollingError
        store.dispatch(setConnectionStatus('reconnecting'));
        store.dispatch(setConnectionError(`Network error: ${error.message}`));
      },

      handleServerError: (error: NetworkError) => {
        console.error('[PollingChat] Server error:', error);
        if (error.status === 401) {
          // Authentication error - stop polling and require re-login
          this.stopPolling();
          store.dispatch(setConnectionStatus('failed'));
          store.dispatch(setConnectionError('Authentication required. Please log in to use chat.'));
        } else if (error.status && error.status >= 500) {
          // Server error - retry with backoff
          store.dispatch(setConnectionStatus('reconnecting'));
          store.dispatch(setConnectionError(`Server error: ${error.message}`));
        }
      },

      handleTimeoutError: (error: NetworkError) => {
        console.error('[PollingChat] Timeout error:', error);
        store.dispatch(setConnectionStatus('reconnecting'));
        store.dispatch(setConnectionError('Request timeout. Retrying...'));
      },

      shouldRetry: (error: NetworkError, attempt: number): boolean => {
        // Don't retry authentication errors
        if (error.status === 401 || error.status === 403) {
          return false;
        }

        // Don't retry client errors (4xx except 401/403)
        if (error.status && error.status >= 400 && error.status < 500) {
          return false;
        }

        // Retry network errors and server errors
        return attempt < this.retryConfig.maxRetries && error.isRetryable;
      }
    };
  }

  /**
   * Setup network detection for offline/online state management
   */
  private setupNetworkDetection(): void {
    // Listen for online/offline events
    window.addEventListener('online', () => {
      console.log('[PollingChat] Network came online');
      this.isOnline = true;
      this.connectionManager.setOnline();
      
      // Resume polling if it was active
      if (this.connectionManager.isPollingActive()) {
        this.processMessageQueue(); // Process any queued messages
      }
    });

    window.addEventListener('offline', () => {
      console.log('[PollingChat] Network went offline');
      this.isOnline = false;
      this.connectionManager.setOffline();
    });
  }

  /**
   * Setup message queue processor for offline scenarios
   */
  private setupMessageQueueProcessor(): void {
    // Process message queue every 5 seconds
    setInterval(() => {
      if (this.isOnline && this.messageQueue.length > 0) {
        this.processMessageQueue();
      }
    }, 5000);
  }

  /**
   * Process queued messages when network is restored
   */
  private async processMessageQueue(): Promise<void> {
    if (!this.isOnline || this.messageQueue.length === 0) {
      return;
    }

    console.log(`[PollingChat] Processing ${this.messageQueue.length} queued messages`);

    // Process messages in order
    const messagesToProcess = [...this.messageQueue];
    this.messageQueue = [];

    for (const queuedMessage of messagesToProcess) {
      try {
        // Check if enough time has passed for retry
        if (new Date() < queuedMessage.nextRetryAt) {
          // Re-queue for later
          this.messageQueue.push(queuedMessage);
          continue;
        }

        // Attempt to send the message
        const response = await this.sendMessageWithRetry(
          queuedMessage.request, 
          queuedMessage.id, 
          queuedMessage.retryCount + 1
        );

        if (response.success) {
          // Update message status to sent and then delivered
          store.dispatch(updateMessageStatus({ id: queuedMessage.id, status: 'sent' }));
          setTimeout(() => {
            store.dispatch(markMessageAsDelivered(queuedMessage.id));
          }, 100);
          console.log('[PollingChat] Queued message sent successfully:', queuedMessage.id);
        } else {
          // Re-queue if retries remaining
          if (queuedMessage.retryCount < this.retryConfig.maxRetries) {
            const delay = this.calculateRetryDelay(queuedMessage.retryCount + 1);
            queuedMessage.retryCount++;
            queuedMessage.nextRetryAt = new Date(Date.now() + delay);
            this.messageQueue.push(queuedMessage);
          } else {
            // Max retries exceeded
            store.dispatch(updateMessageStatus({ id: queuedMessage.id, status: 'failed' }));
            console.error('[PollingChat] Queued message failed after max retries:', queuedMessage.id);
          }
        }
      } catch (error) {
        console.error('[PollingChat] Error processing queued message:', error);
        
        // Re-queue if retries remaining
        if (queuedMessage.retryCount < this.retryConfig.maxRetries) {
          const delay = this.calculateRetryDelay(queuedMessage.retryCount + 1);
          queuedMessage.retryCount++;
          queuedMessage.nextRetryAt = new Date(Date.now() + delay);
          this.messageQueue.push(queuedMessage);
        } else {
          // Max retries exceeded
          store.dispatch(updateMessageStatus({ id: queuedMessage.id, status: 'failed' }));
        }
      }
    }
  }

  /**
   * Calculate retry delay with exponential backoff and jitter
   */
  private calculateRetryDelay(attempt: number): number {
    const baseDelay = Math.min(
      this.retryConfig.baseDelay * Math.pow(this.retryConfig.backoffMultiplier, attempt - 1),
      this.retryConfig.maxDelay
    );

    // Add jitter to prevent thundering herd
    const jitter = baseDelay * this.retryConfig.jitterFactor * (Math.random() - 0.5);
    return Math.max(1000, baseDelay + jitter);
  }

  /**
   * Create network error from fetch error
   */
  private createNetworkError(error: any, status?: number): NetworkError {
    const networkError = new Error(error.message || 'Network error') as NetworkError;
    networkError.code = error.code;
    networkError.status = status;
    networkError.isNetworkError = true;
    
    // Determine if error is retryable
    networkError.isRetryable = this.isRetryableError(error, status);
    
    return networkError;
  }

  /**
   * Determine if an error is retryable
   */
  private isRetryableError(error: any, status?: number): boolean {
    // Network errors are retryable
    if (error.name === 'TypeError' && error.message.includes('fetch')) {
      return true;
    }

    // Timeout errors are retryable
    if (error.name === 'AbortError' || error.message.includes('timeout')) {
      return true;
    }

    // Server errors (5xx) are retryable
    if (status && status >= 500) {
      return true;
    }

    // Rate limiting (429) is retryable
    if (status === 429) {
      return true;
    }

    // Client errors (4xx except 401/403) are not retryable
    if (status && status >= 400 && status < 500 && status !== 401 && status !== 403) {
      return false;
    }

    // Default to retryable for unknown errors
    return true;
  }

  /**
   * Add duplicate message prevention using client-side tracking
   */
  private isDuplicateMessage(messageId: string): boolean {
    return this.sentMessageIds.has(messageId);
  }

  /**
   * Mark message as sent to prevent duplicates
   */
  private markMessageAsSent(messageId: string): void {
    this.sentMessageIds.add(messageId);
    
    // Clean up old message IDs (keep last 1000)
    if (this.sentMessageIds.size > 1000) {
      const idsArray = Array.from(this.sentMessageIds);
      const toRemove = idsArray.slice(0, idsArray.length - 1000);
      toRemove.forEach(id => this.sentMessageIds.delete(id));
    }
  }

  /**
   * Schedule the next polling attempt with proper interval management
   */
  private scheduleNextPoll(): void {
    if (!this.connectionManager.isPollingActive()) {
      console.log('[PollingChat] Polling inactive, not scheduling next poll');
      return;
    }

    // Clear any existing timeout
    if (this.pollingIntervalId) {
      clearTimeout(this.pollingIntervalId);
    }

    const interval = this.calculatePollingInterval();
    
    // Check if we should pause polling (page not visible and user away)
    if (!document.hidden || interval < this.pollingStrategy.maxInterval) {
      console.log(`[PollingChat] Scheduling next poll in ${interval}ms`);
      
      this.pollingIntervalId = setTimeout(async () => {
        if (this.connectionManager.isPollingActive()) {
          await this.pollForMessages();
          this.scheduleNextPoll(); // Schedule next poll recursively
        }
      }, interval);
    } else {
      // Page is hidden and using max interval, schedule a longer check
      console.log('[PollingChat] Page hidden, scheduling visibility check instead of poll');
      this.pollingIntervalId = setTimeout(() => {
        if (this.connectionManager.isPollingActive()) {
          this.scheduleNextPoll(); // Re-evaluate scheduling
        }
      }, 5000); // Check every 5 seconds if we should resume polling
    }
  }

  /**
   * Create polling loop that fetches new messages every 3-5 seconds with smart intervals
   */
  private createPollingLoop(): void {
    console.log('[PollingChat] Creating polling loop');
    
    // Initial poll
    this.pollForMessages().then(() => {
      // Start the recurring polling schedule
      this.scheduleNextPoll();
    }).catch(error => {
      console.error('[PollingChat] Initial poll failed:', error);
      this.handlePollingError(error);
      // Still schedule next poll to retry
      this.scheduleNextPoll();
    });
  }

  /**
   * Calculate the appropriate polling interval based on activity and errors
   * Implements smart polling intervals (faster when active, slower when idle)
   */
  private calculatePollingInterval(): number {
    // If there are errors, use exponential backoff with appropriate backoff strategies
    const errorCount = this.connectionManager.getErrorCount();
    if (errorCount > 0) {
      const backoffInterval = Math.min(
        this.pollingStrategy.baseInterval * Math.pow(this.pollingStrategy.backoffMultiplier, errorCount),
        this.pollingStrategy.maxInterval
      );
      console.log(`[PollingChat] Using backoff interval: ${backoffInterval}ms (errors: ${errorCount})`);
      return backoffInterval;
    }

    const now = Date.now();
    const isPageVisible = !document.hidden;
    
    // Enhanced activity detection using multiple signals
    const timeSinceLastActivity = now - this.lastActivity.getTime();
    const timeSinceLastTyping = now - this.activityTracker.lastTyping.getTime();
    const timeSinceLastMessageSent = now - this.activityTracker.lastMessageSent.getTime();
    const timeSinceLastMouseMove = now - this.activityTracker.lastMouseMove.getTime();
    const timeSinceLastKeyPress = now - this.activityTracker.lastKeyPress.getTime();
    const timeSincePageFocus = now - this.activityTracker.lastPageFocus.getTime();

    // Check if user is actively chatting (recent message activity)
    const chatState = store.getState().chat;
    const recentMessages = chatState.messages.filter(msg => 
      now - new Date(msg.timestamp).getTime() < 60000 // Messages in last minute
    );
    const hasRecentChatActivity = recentMessages.length > 0;

    // Determine activity level based on multiple signals
    const isCurrentlyTyping = this.activityTracker.isTyping;
    const isRecentlyTyping = timeSinceLastTyping < 5000; // 5 seconds
    const isRecentlyActive = timeSinceLastActivity < 30000; // 30 seconds
    const hasRecentMessageActivity = timeSinceLastMessageSent < 30000; // 30 seconds
    const hasRecentMouseActivity = timeSinceLastMouseMove < 60000; // 1 minute
    const hasRecentKeyActivity = timeSinceLastKeyPress < 30000; // 30 seconds
    const isPageRecentlyFocused = timeSincePageFocus < 60000; // 1 minute

    let interval: number;
    let reason: string;

    // Stop polling when page is not visible or user is away
    if (!isPageVisible) {
      interval = this.pollingStrategy.maxInterval;
      reason = 'Page not visible, using max interval to conserve resources';
    }
    // Use faster polling (2s) when user is actively chatting
    else if (isCurrentlyTyping || isRecentlyTyping) {
      interval = this.pollingStrategy.activeInterval;
      reason = 'User is typing, using fast polling for real-time experience';
    }
    // Use fast polling when user just sent a message (expecting response)
    else if (hasRecentMessageActivity) {
      interval = this.pollingStrategy.activeInterval;
      reason = 'Recent message sent, using fast polling for response';
    }
    // Use fast polling when there's recent chat activity and user is active
    else if (hasRecentChatActivity && (isRecentlyActive || hasRecentKeyActivity)) {
      interval = this.pollingStrategy.activeInterval;
      reason = 'Active conversation detected, using fast polling';
    }
    // Use base interval when user is generally active
    else if (isRecentlyActive && (hasRecentMouseActivity || hasRecentKeyActivity || isPageRecentlyFocused)) {
      interval = this.pollingStrategy.baseInterval;
      reason = 'User active, using base polling interval';
    }
    // Switch to slower polling (10s) when user is idle
    else {
      interval = this.pollingStrategy.inactiveInterval;
      reason = 'User idle, using slow polling to conserve resources';
    }

    console.log(`[PollingChat] Smart polling: ${interval}ms - ${reason}`, {
      isPageVisible,
      isCurrentlyTyping,
      isRecentlyTyping,
      hasRecentMessageActivity,
      hasRecentChatActivity,
      isRecentlyActive,
      hasRecentMouseActivity,
      hasRecentKeyActivity,
      timeSinceLastActivity: Math.round(timeSinceLastActivity / 1000) + 's',
      timeSinceLastTyping: Math.round(timeSinceLastTyping / 1000) + 's',
      timeSinceLastMessageSent: Math.round(timeSinceLastMessageSent / 1000) + 's',
    });

    return interval;
  }

  /**
   * Poll for new messages with efficient polling using last message ID/timestamp
   */
  private async pollForMessages(): Promise<void> {
    if (!this.connectionManager.isPollingActive()) {
      return;
    }

    try {
      const currentSession = store.getState().chat.currentSession;
      if (!currentSession) {
        // No active session, skip polling but don't treat as error
        console.log('[PollingChat] No active session, skipping poll');
        return;
      }

      const startTime = Date.now();
      
      // Get all messages for the session
      const messages = await this.getMessages(currentSession.id);
      const latency = Date.now() - startTime;

      // Record successful poll with connection manager
      this.connectionManager.recordSuccessfulPoll(latency);
      
      // Update Redux state for successful poll
      store.dispatch(updatePollTime());
      store.dispatch(setLastSuccessfulPoll(new Date().toISOString()));
      store.dispatch(updateLatency(latency));

      // Process new messages if any
      if (messages.length > 0) {
        console.log(`[PollingChat] Received ${messages.length} new messages`);
        
        // Use batch update for efficiency
        store.dispatch(addMessages(messages as ChatMessage[]));
        
        // Update last message ID for efficient polling
        if (messages.length > 0) {
          this.lastMessageId = messages[messages.length - 1].id;
        }
        
        // Show typing indicator for assistant messages
        const hasAssistantMessage = messages.some(msg => msg.type === 'assistant');
        if (hasAssistantMessage) {
          store.dispatch(setTyping(true));
          setTimeout(() => store.dispatch(setTyping(false)), 1000);
        }
        
        // Update activity when new messages arrive
        this.updateActivity();
      }

    } catch (error) {
      // Only log polling errors occasionally to reduce spam
      const errorCount = this.connectionManager.getErrorCount();
      if (errorCount <= 2 || errorCount % 10 === 0) {
        console.warn('[PollingChat] Polling error:', error);
      }
      this.handlePollingError(error);
    }
  }

  /**
   * Handle polling errors with appropriate backoff strategies
   */
  private handlePollingError(error: any): void {
    // Create network error for proper error handling
    const networkError = this.createNetworkError(error);
    const errorMessage = networkError.message || 'Polling connection error';
    
    // Record error with connection manager
    this.connectionManager.recordPollingError(errorMessage);
    
    // Update Redux state for error
    store.dispatch(incrementErrorCount());

    // Handle different types of errors using error handler
    if (networkError.status === 401 || networkError.status === 403) {
      // Authentication error - stop polling and require re-login
      console.warn('[PollingChat] Authentication error, stopping polling');
      this.errorHandler.handleServerError(networkError);
      this.stopPolling();
      return;
    }

    // Handle network errors (reduce logging frequency)
    const errorCount = this.connectionManager.getErrorCount();
    if (errorCount <= 3) { // Only log first few errors
      if (networkError.isNetworkError) {
        this.errorHandler.handleNetworkError(networkError);
      } else if (networkError.status && networkError.status >= 500) {
        this.errorHandler.handleServerError(networkError);
      } else if (error.name === 'AbortError' || error.message.includes('timeout')) {
        this.errorHandler.handleTimeoutError(networkError);
      }
    }

    // If too many consecutive errors, pause polling temporarily
    if (errorCount >= 5) {
      console.warn('[PollingChat] Too many consecutive polling errors, pausing');
      
      // Retry after a longer delay (30 seconds) with automatic reconnection when network is restored
      setTimeout(() => {
        if (this.connectionManager.isPollingActive() && this.isOnline) {
          console.log('[PollingChat] Resuming polling after error recovery period');
          this.connectionManager.resetErrorCount();
          store.dispatch(resetErrorCount());
        }
      }, 30000); // 30 second recovery delay
    }

    // Show progressive error messages based on error count
    if (errorCount >= 3 && errorCount < 5) {
      store.dispatch(setError('Connection unstable. Retrying with longer intervals.'));
    } else if (errorCount >= 2) {
      store.dispatch(setError('Connection issues detected. Retrying...'));
    }
  }

  /**
   * Send message with retry logic and exponential backoff (1s, 2s, 4s, 8s intervals)
   */
  private async sendMessageWithRetry(request: SendMessageRequest, messageId: string, attempt: number = 1): Promise<SendMessageResponse> {
    try {
      const response = await this.makeRequest<SendMessageResponse>(
        '/api/v1/admin/chat/messages',
        {
          method: 'POST',
          body: JSON.stringify(request),
        }
      );

      return response;
    } catch (error) {
      console.error(`[PollingChat] Send attempt ${attempt} failed:`, error);

      // Create network error for proper error handling
      const networkError = this.createNetworkError(error);
      
      // Check if we should retry
      if (!this.errorHandler.shouldRetry(networkError, attempt)) {
        // Handle the error appropriately
        if (networkError.status === 401) {
          this.errorHandler.handleServerError(networkError);
        } else if (networkError.isNetworkError) {
          this.errorHandler.handleNetworkError(networkError);
        } else {
          this.errorHandler.handleServerError(networkError);
        }
        throw error;
      }

      if (attempt >= this.retryConfig.maxRetries) {
        throw error;
      }

      // Calculate delay with exponential backoff and jitter
      const delay = this.calculateRetryDelay(attempt);
      
      console.log(`[PollingChat] Retrying send in ${delay}ms (attempt ${attempt + 1}/${this.retryConfig.maxRetries})`);
      
      return new Promise((resolve, reject) => {
        const timeoutId = setTimeout(async () => {
          this.retryTimeouts.delete(messageId);
          try {
            const result = await this.sendMessageWithRetry(request, messageId, attempt + 1);
            resolve(result);
          } catch (retryError) {
            reject(retryError);
          }
        }, delay);

        this.retryTimeouts.set(messageId, timeoutId);
      });
    }
  }

  /**
   * Make HTTP request with authentication, timeout, and enhanced error handling
   * Supports conditional requests and compression
   */
  private async makeRequest<T>(endpoint: string, options: RequestInit = {}, conditionalHeaders?: { etag?: string; lastModified?: string }): Promise<T> {
    const url = `${this.baseUrl}${endpoint}`;
    
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      'Accept-Encoding': 'gzip, deflate, br', // Request compression
      ...options.headers,
    };
    
    // Add conditional request headers for caching
    if (conditionalHeaders) {
      if (conditionalHeaders.etag) {
        (headers as Record<string, string>)['If-None-Match'] = conditionalHeaders.etag;
      }
      if (conditionalHeaders.lastModified) {
        (headers as Record<string, string>)['If-Modified-Since'] = conditionalHeaders.lastModified;
      }
    }
    
    // Add authorization header if token exists
    if (this.authToken) {
      (headers as Record<string, string>)['Authorization'] = `Bearer ${this.authToken}`;
    }
    
    // Create abort controller for timeout
    const controller = new AbortController();
    const timeoutId = setTimeout(() => controller.abort(), 10000); // 10 second timeout
    
    const config: RequestInit = {
      headers,
      signal: controller.signal,
      ...options,
    };

    try {
      const response = await fetch(url, config);
      clearTimeout(timeoutId);
      
      // Handle 304 Not Modified for conditional requests
      if (response.status === 304) {
        console.log('[PollingChat] Received 304 Not Modified, using cached data');
        // Return a special response indicating cache should be used
        return { notModified: true } as T;
      }
      
      if (!response.ok) {
        // If unauthorized, clear token
        if (response.status === 401) {
          this.authToken = null;
          localStorage.removeItem('adminToken');
          store.dispatch(setConnectionError('Authentication required. Please log in to use chat.'));
        }
        
        const errorData = await response.json().catch(() => ({}));
        const error = new Error(errorData.error || `HTTP error! status: ${response.status}`) as NetworkError;
        error.status = response.status;
        error.isNetworkError = false;
        error.isRetryable = this.isRetryableError(error, response.status);
        throw error;
      }

      const data = await response.json();
      
      // Extract caching headers for future conditional requests
      const etag = response.headers.get('ETag');
      const lastModified = response.headers.get('Last-Modified');
      
      // Add caching metadata to response if available
      if (etag || lastModified) {
        (data as any)._cacheMetadata = { etag, lastModified };
      }
      
      return data;
    } catch (error: any) {
      clearTimeout(timeoutId);
      
      // Handle abort/timeout errors
      if (error.name === 'AbortError') {
        const timeoutError = new Error('Request timeout') as NetworkError;
        timeoutError.isNetworkError = true;
        timeoutError.isRetryable = true;
        throw timeoutError;
      }
      
      // Handle network errors
      if (error.name === 'TypeError' && error.message.includes('fetch')) {
        const networkError = new Error('Network error: Unable to connect to server') as NetworkError;
        networkError.isNetworkError = true;
        networkError.isRetryable = true;
        throw networkError;
      }
      
      // Re-throw other errors
      throw error;
    }
  }

  /**
   * Update last activity timestamp
   */
  private updateActivity(): void {
    this.lastActivity = new Date();
    this.isUserActive = true;
  }

  /**
   * Track message sending activity for smart polling
   */
  private trackMessageSendingActivity(): void {
    this.activityTracker.lastMessageSent = new Date();
    this.updateActivity();
    console.log('[PollingChat] Message sending activity tracked');
  }

  /**
   * Get cached messages if available and not expired
   */
  private getCachedMessages(sessionId: string): PollingChatMessage[] | null {
    const cacheKey = `session-${sessionId}`;
    const cached = this.messageCache.get(cacheKey);
    
    if (!cached) {
      return null;
    }
    
    // Check if cache is expired
    const now = Date.now();
    const cacheAge = now - cached.timestamp.getTime();
    
    if (cacheAge > this.cacheExpiryTime) {
      console.log('[PollingChat] Message cache expired, removing');
      this.messageCache.delete(cacheKey);
      return null;
    }
    
    console.log(`[PollingChat] Using cached messages (${cached.messages.length} messages, age: ${Math.round(cacheAge / 1000)}s)`);
    return cached.messages;
  }

  /**
   * Cache messages with metadata for conditional requests
   */
  private cacheMessages(sessionId: string, messages: PollingChatMessage[], etag?: string, lastModified?: string): void {
    const cacheKey = `session-${sessionId}`;
    
    this.messageCache.set(cacheKey, {
      messages: [...messages], // Create a copy to avoid mutations
      timestamp: new Date(),
      etag,
      lastModified,
    });
    
    console.log(`[PollingChat] Cached ${messages.length} messages for session ${sessionId}`);
    
    // Clean up old cache entries (keep last 10 sessions)
    if (this.messageCache.size > 10) {
      const entries = Array.from(this.messageCache.entries());
      const oldestEntries = entries
        .sort((a, b) => a[1].timestamp.getTime() - b[1].timestamp.getTime())
        .slice(0, this.messageCache.size - 10);
      
      oldestEntries.forEach(([key]) => {
        this.messageCache.delete(key);
      });
    }
  }

  /**
   * Get cache metadata for conditional requests
   */
  private getCacheMetadata(sessionId: string): { etag?: string; lastModified?: string } | null {
    const cacheKey = `session-${sessionId}`;
    const cached = this.messageCache.get(cacheKey);
    
    if (!cached) {
      return null;
    }
    
    return {
      etag: cached.etag,
      lastModified: cached.lastModified,
    };
  }

  /**
   * Clear message cache for a session
   */
  private clearMessageCache(sessionId?: string): void {
    if (sessionId) {
      const cacheKey = `session-${sessionId}`;
      this.messageCache.delete(cacheKey);
      console.log(`[PollingChat] Cleared cache for session ${sessionId}`);
    } else {
      this.messageCache.clear();
      console.log('[PollingChat] Cleared all message cache');
    }
  }

  /**
   * Setup activity tracking for smart polling
   */
  private setupActivityTracking(): void {
    // Track user interactions
    const activityEvents = ['mousedown', 'mousemove', 'keypress', 'scroll', 'touchstart'];
    
    const handleActivity = () => {
      this.updateActivity();
    };

    activityEvents.forEach(event => {
      document.addEventListener(event, handleActivity, { passive: true });
    });

    // Track page visibility
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        this.isUserActive = false;
      } else {
        this.updateActivity();
      }
    });

    // Set user as inactive after 5 minutes of no activity
    setInterval(() => {
      const timeSinceLastActivity = Date.now() - this.lastActivity.getTime();
      if (timeSinceLastActivity > 300000) { // 5 minutes
        this.isUserActive = false;
      }
    }, 60000); // Check every minute
  }

  /**
   * Setup advanced activity tracking for smart polling intervals
   * Detects user activity (typing, sending messages, page focus)
   */
  private setupAdvancedActivityTracking(): void {
    // Track typing activity in chat input fields
    const trackTyping = (event: Event) => {
      const target = event.target as HTMLElement;
      if (target && (target.tagName === 'INPUT' || target.tagName === 'TEXTAREA')) {
        this.activityTracker.lastTyping = new Date();
        this.activityTracker.lastKeyPress = new Date();
        this.activityTracker.isTyping = true;
        
        // Clear existing typing timeout
        if (this.activityTracker.typingTimeout) {
          clearTimeout(this.activityTracker.typingTimeout);
        }
        
        // Set typing to false after 3 seconds of no typing
        this.activityTracker.typingTimeout = setTimeout(() => {
          this.activityTracker.isTyping = false;
        }, 3000);
        
        this.updateActivity();
      }
    };

    // Track mouse movement for activity detection
    const trackMouseActivity = () => {
      this.activityTracker.lastMouseMove = new Date();
      this.updateActivity();
    };

    // Track page focus/blur for smart polling
    const trackPageFocus = () => {
      this.activityTracker.lastPageFocus = new Date();
      this.updateActivity();
      
      // Clear focus timeout
      if (this.activityTracker.focusTimeout) {
        clearTimeout(this.activityTracker.focusTimeout);
      }
    };

    const trackPageBlur = () => {
      // Set a timeout to reduce polling when page loses focus
      this.activityTracker.focusTimeout = setTimeout(() => {
        // Page has been unfocused for 30 seconds, reduce activity
        this.isUserActive = false;
      }, 30000);
    };

    // Add event listeners for advanced activity tracking
    document.addEventListener('keydown', trackTyping, { passive: true });
    document.addEventListener('input', trackTyping, { passive: true });
    document.addEventListener('mousemove', trackMouseActivity, { passive: true });
    window.addEventListener('focus', trackPageFocus);
    window.addEventListener('blur', trackPageBlur);

    // Track page visibility changes for smart polling
    document.addEventListener('visibilitychange', () => {
      if (document.hidden) {
        console.log('[PollingChat] Page hidden, reducing polling frequency');
        this.isUserActive = false;
      } else {
        console.log('[PollingChat] Page visible, resuming normal polling');
        this.activityTracker.lastPageFocus = new Date();
        this.updateActivity();
      }
    });

    // Cleanup on page unload
    window.addEventListener('beforeunload', () => {
      if (this.activityTracker.typingTimeout) {
        clearTimeout(this.activityTracker.typingTimeout);
      }
      if (this.activityTracker.focusTimeout) {
        clearTimeout(this.activityTracker.focusTimeout);
      }
    });
  }
}

// Create singleton instance
export const pollingChatService = new PollingChatService();

// Export for use in components
export default pollingChatService;
import { store } from '../store';
import {
  setConnectionStatus,
  setWebSocket,
  setConnectionId,
  incrementReconnectAttempts,
  resetReconnectAttempts,
  setConnectionError,
  updateLatency,
  setPingTime,
  setPongTime,
  setHealthStatus,
  cleanup,
  ConnectionStatus,
} from '../store/slices/connectionSlice';
import {
  addMessage,
  addOptimisticMessage,
  updateMessageStatus,
  removeFailedMessage,
  setTyping,
  setError,
  ChatMessage,
} from '../store/slices/chatSlice';
import { connectionDiagnostics } from './connectionDiagnostics';

export interface WebSocketMessage {
  type: 'message' | 'typing' | 'status' | 'error' | 'ping' | 'pong' | 'heartbeat';
  session_id?: string;
  message_id?: string;
  content?: string;
  metadata?: Record<string, any>;
  timestamp?: string;
  error?: string;
}

export interface ChatRequest {
  message: string;
  session_id?: string;
  client_name?: string;
  context?: string;
  quick_action?: string;
}

export interface ChatResponse {
  message: ChatMessage;
  session_id: string;
  success: boolean;
  error?: string;
}

class WebSocketService {
  private ws: WebSocket | null = null;
  private reconnectTimeoutId: NodeJS.Timeout | null = null;
  private pingIntervalId: NodeJS.Timeout | null = null;
  private healthCheckIntervalId: NodeJS.Timeout | null = null;
  private messageQueue: WebSocketMessage[] = [];
  private isReconnecting = false;
  private lastPingTime = 0;
  private lastDiagnosticsRun = 0;
  private readonly pingInterval = 30000; // 30 seconds
  private readonly healthCheckInterval = 5000; // 5 seconds
  private readonly connectionTimeout = 10000; // 10 seconds
  private connectionPromise: Promise<void> | null = null;
  private isConnecting = false;

  constructor() {
    console.log('[WebSocket] Service instance created', {
      timestamp: new Date().toISOString(),
      stack: new Error().stack?.split('\n').slice(1, 3).join('\n')
    });
    
    // Bind methods to preserve context
    this.connect = this.connect.bind(this);
    this.disconnect = this.disconnect.bind(this);
    this.send = this.send.bind(this);
    this.handleOpen = this.handleOpen.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.handleError = this.handleError.bind(this);

    // Handle page unload to prevent connection errors
    window.addEventListener('beforeunload', () => {
      console.log('[WebSocket] Page unloading, cleaning up connection');
      this.disconnect();
    });
  }

  /**
   * Establish WebSocket connection
   */
  public connect(): Promise<void> {
    // Return existing connection promise if already connecting
    if (this.connectionPromise) {
      console.log('[WebSocket] Returning existing connection promise');
      return this.connectionPromise;
    }

    // If already connected, resolve immediately
    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('[WebSocket] Already connected, resolving immediately');
      return Promise.resolve();
    }

    // Create new connection promise
    this.connectionPromise = new Promise((resolve, reject) => {
      try {
        console.log('[WebSocket] Starting new connection attempt', {
          currentState: this.ws?.readyState,
          isConnecting: this.isConnecting,
          hasExistingWs: !!this.ws
        });

        // Prevent multiple simultaneous connection attempts
        if (this.isConnecting) {
          console.log('[WebSocket] Connection already in progress, waiting...');
          return;
        }

        this.isConnecting = true;

        // Clean up existing connection
        this.cleanup();

        // Update connection status
        store.dispatch(setConnectionStatus('connecting'));

        // Create WebSocket URL with environment-based configuration and authentication
        const wsUrl = this.getWebSocketUrl();
        console.log('[WebSocket] Connecting to:', wsUrl.replace(/token=[^&]+/, 'token=***'));

        // Create WebSocket connection
        this.ws = new WebSocket(wsUrl);
        store.dispatch(setWebSocket(this.ws));

        // Set up event listeners
        this.ws.onopen = () => {
          console.log('[WebSocket] onopen event fired');
          this.isConnecting = false;
          this.connectionPromise = null;
          this.handleOpen();
          resolve();
        };
        this.ws.onmessage = this.handleMessage;
        this.ws.onclose = (event) => {
          console.error('[WebSocket] onclose event fired - DETAILED DEBUG', { 
            code: event.code, 
            reason: event.reason,
            wasClean: event.wasClean,
            timeStamp: event.timeStamp,
            readyState: this.ws?.readyState,
            stack: new Error().stack?.split('\n').slice(1, 8).join('\n')
          });
          this.isConnecting = false;
          this.connectionPromise = null;
          
          // Handle authentication failures (code 1008 or 1006 with specific reasons)
          if (event.code === 1008 || (event.code === 1006 && event.reason?.includes('auth'))) {
            const authError = 'WebSocket authentication failed. Please log in again.';
            store.dispatch(setConnectionError(authError));
            // Clear invalid token
            localStorage.removeItem('adminToken');
            reject(new Error(authError));
          } else {
            this.handleClose(event);
            // If the connection closes before it was fully established, reject the promise
            if (event.code !== 1000) {
              reject(new Error(`WebSocket connection closed unexpectedly: ${event.code} ${event.reason || 'No reason provided'}`));
            }
          }
        };
        this.ws.onerror = (error) => {
          console.log('[WebSocket] onerror event fired', error);
          this.isConnecting = false;
          this.connectionPromise = null;
          this.handleError(error);
          reject(error);
        };

        // Set connection timeout
        setTimeout(() => {
          if (this.ws?.readyState === WebSocket.CONNECTING) {
            this.isConnecting = false;
            this.connectionPromise = null;
            this.ws.close();
            reject(new Error('Connection timeout'));
          }
        }, this.connectionTimeout);

      } catch (error) {
        this.isConnecting = false;
        this.connectionPromise = null;
        
        const errorMessage = error instanceof Error ? error.message : 'Failed to create WebSocket connection';
        console.error('WebSocket connection error:', errorMessage);
        
        // Handle authentication token errors specifically
        if (errorMessage.includes('No admin token found')) {
          store.dispatch(setConnectionError('Authentication required. Please log in to use chat.'));
        } else {
          store.dispatch(setConnectionError(errorMessage));
        }
        
        store.dispatch(setConnectionStatus('disconnected'));
        reject(error);
      }
    });

    return this.connectionPromise;
  }

  /**
   * Disconnect WebSocket
   */
  public disconnect(): void {
    console.log('[WebSocket] Disconnect requested');
    this.isReconnecting = false;
    this.isConnecting = false;
    this.connectionPromise = null;
    this.cleanup();
    store.dispatch(setConnectionStatus('disconnected'));
  }

  /**
   * Send message through WebSocket with queuing support
   */
  public send(message: WebSocketMessage): boolean {
    if (this.ws?.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(JSON.stringify(message));
        return true;
      } catch (error) {
        console.error('Failed to send WebSocket message:', error);
        this.queueMessage(message);
        return false;
      }
    } else {
      // Queue message for later sending
      this.queueMessage(message);
      
      // Attempt to reconnect if not already doing so
      if (!this.isReconnecting) {
        this.reconnect();
      }
      
      return false;
    }
  }

  /**
   * Send chat message with optimistic updates
   */
  public sendChatMessage(request: ChatRequest): string {
    const messageId = `msg-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    
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

    // Send message
    const wsMessage: WebSocketMessage = {
      type: 'message',
      message_id: messageId,
      content: request.message,
      session_id: request.session_id,
      metadata: {
        client_name: request.client_name,
        context: request.context,
        quick_action: request.quick_action,
      },
      timestamp: new Date().toISOString(),
    };

    const sent = this.send(wsMessage);
    
    if (!sent) {
      // Update message status to failed if couldn't send
      store.dispatch(updateMessageStatus({ id: messageId, status: 'failed' }));
    } else {
      // Update message status to sent
      store.dispatch(updateMessageStatus({ id: messageId, status: 'sent' }));
    }

    return messageId;
  }

  /**
   * Get current connection status
   */
  public getConnectionStatus(): ConnectionStatus {
    return store.getState().connection.status;
  }

  /**
   * Get connection health status
   */
  public isHealthy(): boolean {
    return store.getState().connection.isHealthy;
  }

  /**
   * Force reconnection
   */
  public forceReconnect(): void {
    this.disconnect();
    setTimeout(() => this.connect(), 1000);
  }

  /**
   * Run connection diagnostics
   */
  public async runDiagnostics(): Promise<void> {
    console.log('🔍 Running WebSocket connection diagnostics...');
    const report = await connectionDiagnostics.runDiagnostics();
    connectionDiagnostics.printReport(report);
    return;
  }

  // Private methods

  /**
   * Get WebSocket URL with environment-based configuration and authentication token
   */
  private getWebSocketUrl(): string {
    // Get authentication token
    const token = localStorage.getItem('adminToken');
    if (!token) {
      throw new Error('No admin token found for WebSocket authentication');
    }

    let baseUrl: string;

    // Check for environment variable first
    const envWsUrl = process.env.REACT_APP_WS_URL;
    console.log('[WebSocket] Environment WS URL:', envWsUrl);
    
    if (envWsUrl) {
      baseUrl = envWsUrl;
    } else {
      // Check for API URL environment variable and construct WebSocket URL
      const apiUrl = process.env.REACT_APP_API_URL;
      console.log('[WebSocket] Environment API URL:', apiUrl);
      
      if (apiUrl) {
        try {
          const url = new URL(apiUrl);
          const protocol = url.protocol === 'https:' ? 'wss:' : 'ws:';
          baseUrl = `${protocol}//${url.host}/api/v1/admin/chat/ws`;
          console.log('[WebSocket] Constructed URL from API URL:', baseUrl);
        } catch (error) {
          console.warn('Invalid REACT_APP_API_URL, falling back to default:', error);
          // Fallback to current host-based URL for backward compatibility
          const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
          baseUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;
          console.log('[WebSocket] Fallback URL (from window.location):', baseUrl);
        }
      } else {
        // Fallback to current host-based URL for backward compatibility
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        baseUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;
        console.log('[WebSocket] Default fallback URL:', baseUrl);
      }
    }

    // Add authentication token as query parameter
    const separator = baseUrl.includes('?') ? '&' : '?';
    const finalUrl = `${baseUrl}${separator}token=${encodeURIComponent(token)}`;
    console.log('[WebSocket] Final WebSocket URL (token masked):', finalUrl.replace(/token=[^&]+/, 'token=***'));
    
    return finalUrl;
  }

  private handleOpen(): void {
    console.log('[WebSocket] CONNECTION OPENED SUCCESSFULLY', {
      time: new Date().toISOString(),
      readyState: this.ws?.readyState,
      isReconnecting: this.isReconnecting,
      queuedMessages: this.messageQueue.length,
      connectionStatus: store.getState().connection.status,
    });
    
    // Implement connection stability check to prevent immediate disconnections
    // Wait 2 seconds to ensure the connection is stable before proceeding
    setTimeout(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        console.log('[WebSocket] Connection stability check passed, proceeding with setup');
        
        // Update connection status
        store.dispatch(setConnectionStatus('connected'));
        store.dispatch(resetReconnectAttempts());
        store.dispatch(setConnectionError(null));
        
        // Generate connection ID
        const connectionId = `conn-${Date.now()}-${Math.random().toString(36).substring(2, 9)}`;
        store.dispatch(setConnectionId(connectionId));
        
        // Send queued messages
        this.sendQueuedMessages();
        
        // Send initial heartbeat to keep connection alive
        setTimeout(() => {
          if (this.ws?.readyState === WebSocket.OPEN) {
            console.log('[WebSocket] Sending initial heartbeat');
            this.send({ type: 'heartbeat', timestamp: new Date().toISOString() });
          }
        }, 1000);
        
        // Start health monitoring after connection is stable
        setTimeout(() => {
          if (this.ws?.readyState === WebSocket.OPEN) {
            console.log('[WebSocket] Starting health monitoring after connection stabilized');
            this.startHealthMonitoring();
          }
        }, 3000); // Wait additional 3 seconds before starting health monitoring
        
        this.isReconnecting = false;
      } else {
        console.error('[WebSocket] Connection stability check failed - connection closed during setup');
        // Connection was closed during the stability check
        this.handleClose(new CloseEvent('close', { code: 1005, reason: 'Connection closed during stability check' }));
      }
    }, 2000); // 2-second stability check
  }

  private handleMessage(event: MessageEvent): void {
    try {
      const data: WebSocketMessage | ChatResponse = JSON.parse(event.data);
      
      if ('type' in data) {
        // Handle WebSocket protocol messages
        this.handleWebSocketMessage(data);
      } else {
        // Handle chat response messages
        this.handleChatResponse(data);
      }
    } catch (error) {
      console.error('Failed to parse WebSocket message:', error);
      store.dispatch(setError('Failed to parse server message'));
    }
  }

  private handleWebSocketMessage(message: WebSocketMessage): void {
    switch (message.type) {
      case 'ping':
        // Respond to ping with pong
        this.send({ type: 'pong', timestamp: new Date().toISOString() });
        break;
        
      case 'pong':
        // Calculate latency
        if (this.lastPingTime) {
          const latency = Date.now() - this.lastPingTime;
          store.dispatch(updateLatency(latency));
        }
        break;
        
      case 'heartbeat':
        // Handle heartbeat response - calculate latency and update health status
        if (this.lastPingTime) {
          const latency = Date.now() - this.lastPingTime;
          store.dispatch(updateLatency(latency));
        }
        store.dispatch(setHealthStatus(true));
        store.dispatch(setPongTime(Date.now()));
        break;
        
      case 'typing':
        // Handle typing indicator
        store.dispatch(setTyping(true));
        setTimeout(() => store.dispatch(setTyping(false)), 3000);
        break;
        
      case 'status':
        // Handle status updates
        if (message.content === 'healthy') {
          store.dispatch(setHealthStatus(true));
        }
        break;
        
      case 'error':
        // Handle server errors
        store.dispatch(setError(message.error || 'Server error'));
        break;
        
      default:
        console.warn('Unknown WebSocket message type:', message.type);
    }
  }

  private handleChatResponse(response: ChatResponse): void {
    if (response.success) {
      // Add message to store
      store.dispatch(addMessage(response.message));
      
      // Update message status if it was optimistic
      if (response.message.id) {
        store.dispatch(updateMessageStatus({ 
          id: response.message.id, 
          status: 'delivered' 
        }));
      }
    } else {
      // Handle chat error
      const error = response.error || 'Chat error occurred';
      store.dispatch(setError(error));
      
      // If there's a message ID, mark it as failed
      if (response.message?.id) {
        store.dispatch(updateMessageStatus({ 
          id: response.message.id, 
          status: 'failed' 
        }));
      }
    }
  }

  private handleClose(event: CloseEvent): void {
    console.error('[WebSocket] CONNECTION CLOSED', {
      time: new Date().toISOString(),
      code: event.code,
      reason: event.reason,
      wasClean: event.wasClean,
      readyState: this.ws?.readyState,
      isReconnecting: this.isReconnecting,
      connectionStatus: store.getState().connection.status,
      reconnectAttempts: store.getState().connection.reconnectAttempts,
    });
    
    // Clean up intervals
    this.stopHealthMonitoring();
    
    // Update connection status
    store.dispatch(setConnectionStatus('disconnected'));
    store.dispatch(setWebSocket(null));
    store.dispatch(setHealthStatus(false));
    
    // Attempt to reconnect if not intentionally closed
    // Code 1005 means no status received - often indicates client-side close
    if (event.code !== 1000 && !this.isReconnecting) {
      console.log('[WebSocket] Attempting to reconnect due to unexpected close');
      this.reconnect();
    }
  }

  private handleError(error: Event): void {
    console.error('[WebSocket] ERROR', {
      time: new Date().toISOString(),
      error,
      readyState: this.ws?.readyState,
      isReconnecting: this.isReconnecting,
      connectionStatus: store.getState().connection.status,
      reconnectAttempts: store.getState().connection.reconnectAttempts,
    });
    const errorMessage = 'WebSocket connection error';
    store.dispatch(setConnectionError(errorMessage));
    store.dispatch(setHealthStatus(false));
    // Run diagnostics when WebSocket errors occur
    this.runDiagnosticsOnError();
  }

  private reconnect(): void {
    if (this.isReconnecting) {
      return;
    }

    this.isReconnecting = true;
    const state = store.getState().connection;
    
    // Check if we've exceeded max attempts
    if (state.reconnectAttempts >= state.maxReconnectAttempts) {
      console.error('Max reconnection attempts reached');
      store.dispatch(setConnectionStatus('failed'));
      store.dispatch(setConnectionError('Failed to reconnect after maximum attempts'));
      this.isReconnecting = false;
      return;
    }

    // Update status and increment attempts
    store.dispatch(setConnectionStatus('reconnecting'));
    store.dispatch(incrementReconnectAttempts());
    
    // Get updated delay after incrementing attempts
    const updatedState = store.getState().connection;
    
    console.log(`Reconnecting in ${updatedState.reconnectDelay}ms (attempt ${updatedState.reconnectAttempts})`);
    
    // Schedule reconnection
    this.reconnectTimeoutId = setTimeout(async () => {
      try {
        await this.connect();
      } catch (error) {
        console.error('Reconnection failed:', error);
        this.isReconnecting = false;
        // Will trigger another reconnection attempt via handleClose
      }
    }, updatedState.reconnectDelay);
  }

  private queueMessage(message: WebSocketMessage): void {
    // Limit queue size to prevent memory issues
    if (this.messageQueue.length >= 100) {
      this.messageQueue.shift(); // Remove oldest message
    }
    
    this.messageQueue.push(message);
  }

  private sendQueuedMessages(): void {
    while (this.messageQueue.length > 0 && this.ws?.readyState === WebSocket.OPEN) {
      const message = this.messageQueue.shift();
      if (message) {
        try {
          this.ws.send(JSON.stringify(message));
        } catch (error) {
          console.error('Failed to send queued message:', error);
          // Put message back at the front of the queue
          this.messageQueue.unshift(message);
          break;
        }
      }
    }
  }

  private startHealthMonitoring(): void {
    console.log('[WebSocket] Starting health monitoring');
    
    // Start ping interval - send heartbeat every 30 seconds
    this.pingIntervalId = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        console.log('[WebSocket] Sending heartbeat');
        this.lastPingTime = Date.now();
        store.dispatch(setPingTime(this.lastPingTime));
        this.send({ type: 'heartbeat', timestamp: new Date().toISOString() });
      }
    }, this.pingInterval);

    // Start health check interval - check every 60 seconds (less aggressive)
    this.healthCheckIntervalId = setInterval(() => {
      const state = store.getState().connection;
      
      // Only check for stale connections after we've sent at least one ping
      // and give more time before considering it stale (3x ping interval instead of 2x)
      if (state.lastPingTime && Date.now() - state.lastPingTime > this.pingInterval * 3) {
        console.warn('[WebSocket] Connection appears stale, forcing reconnection', {
          lastPingTime: new Date(state.lastPingTime).toISOString(),
          timeSinceLastPing: Date.now() - state.lastPingTime,
          threshold: this.pingInterval * 3
        });
        store.dispatch(setHealthStatus(false));
        this.forceReconnect();
      }
    }, 60000); // Check every 60 seconds instead of 5 seconds
  }

  private stopHealthMonitoring(): void {
    if (this.pingIntervalId) {
      clearInterval(this.pingIntervalId);
      this.pingIntervalId = null;
    }
    
    if (this.healthCheckIntervalId) {
      clearInterval(this.healthCheckIntervalId);
      this.healthCheckIntervalId = null;
    }
  }

  private cleanup(): void {
    console.log('[WebSocket] Cleanup called', {
      hasWebSocket: !!this.ws,
      readyState: this.ws?.readyState,
      isConnecting: this.isConnecting,
      stack: new Error().stack?.split('\n').slice(1, 6).join('\n')
    });
    
    // Clear timeouts and intervals
    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId);
      this.reconnectTimeoutId = null;
    }
    
    this.stopHealthMonitoring();
    
    // Close WebSocket connection only if it's not in the process of connecting
    if (this.ws && this.ws.readyState !== WebSocket.CONNECTING) {
      try {
        console.log('[WebSocket] Closing WebSocket connection from cleanup', {
          readyState: this.ws.readyState,
          wasConnecting: this.isConnecting
        });
        this.ws.close(1000, 'Client disconnect');
      } catch (error) {
        console.warn('Error closing WebSocket:', error);
      }
      this.ws = null;
    } else if (this.ws && this.ws.readyState === WebSocket.CONNECTING) {
      console.log('[WebSocket] WebSocket is connecting, deferring cleanup');
      // Don't close a connecting WebSocket, let it complete or fail naturally
    }
    
    // Clear message queue
    this.messageQueue = [];
    
    // Update store
    store.dispatch(cleanup());
  }

  /**
   * Run diagnostics when WebSocket errors occur
   */
  private async runDiagnosticsOnError(): Promise<void> {
    // Debounce diagnostics to avoid running too frequently
    if (this.lastDiagnosticsRun && Date.now() - this.lastDiagnosticsRun < 30000) {
      return; // Don't run diagnostics more than once every 30 seconds
    }

    this.lastDiagnosticsRun = Date.now();
    
    try {
      console.log('🔍 WebSocket error detected, running diagnostics...');
      const report = await connectionDiagnostics.runDiagnostics();
      connectionDiagnostics.printReport(report);
      
      // Update connection error with more specific information
      if (report.recommendations.length > 0) {
        const primaryRecommendation = report.recommendations[0];
        store.dispatch(setConnectionError(primaryRecommendation));
      }
    } catch (error) {
      console.error('Failed to run diagnostics:', error);
    }
  }
}

// Create singleton instance
export const websocketService = new WebSocketService();

// Export for use in components
export default websocketService;
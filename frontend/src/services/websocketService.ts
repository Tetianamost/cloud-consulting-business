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

export interface WebSocketMessage {
  type: 'message' | 'typing' | 'status' | 'error' | 'ping' | 'pong';
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
  private readonly pingInterval = 30000; // 30 seconds
  private readonly healthCheckInterval = 5000; // 5 seconds
  private readonly connectionTimeout = 10000; // 10 seconds

  constructor() {
    // Bind methods to preserve context
    this.connect = this.connect.bind(this);
    this.disconnect = this.disconnect.bind(this);
    this.send = this.send.bind(this);
    this.handleOpen = this.handleOpen.bind(this);
    this.handleMessage = this.handleMessage.bind(this);
    this.handleClose = this.handleClose.bind(this);
    this.handleError = this.handleError.bind(this);
  }

  /**
   * Establish WebSocket connection
   */
  public connect(): Promise<void> {
    return new Promise((resolve, reject) => {
      try {
        const token = localStorage.getItem('adminToken');
        if (!token) {
          const error = 'No admin token found';
          store.dispatch(setConnectionError(error));
          reject(new Error(error));
          return;
        }

        // Clean up existing connection
        this.cleanup();

        // Update connection status
        store.dispatch(setConnectionStatus('connecting'));

        // Create WebSocket URL
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;

        // Create WebSocket connection
        this.ws = new WebSocket(wsUrl);
        store.dispatch(setWebSocket(this.ws));

        // Set up event listeners
        this.ws.onopen = () => {
          this.handleOpen();
          resolve();
        };
        this.ws.onmessage = this.handleMessage;
        this.ws.onclose = this.handleClose;
        this.ws.onerror = (error) => {
          this.handleError(error);
          reject(error);
        };

        // Set connection timeout
        setTimeout(() => {
          if (this.ws?.readyState === WebSocket.CONNECTING) {
            this.ws.close();
            reject(new Error('Connection timeout'));
          }
        }, this.connectionTimeout);

      } catch (error) {
        const errorMessage = error instanceof Error ? error.message : 'Failed to create WebSocket connection';
        store.dispatch(setConnectionError(errorMessage));
        reject(error);
      }
    });
  }

  /**
   * Disconnect WebSocket
   */
  public disconnect(): void {
    this.isReconnecting = false;
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

  // Private methods

  private handleOpen(): void {
    console.log('WebSocket connected');
    
    // Update connection status
    store.dispatch(setConnectionStatus('connected'));
    store.dispatch(resetReconnectAttempts());
    store.dispatch(setConnectionError(null));
    
    // Generate connection ID
    const connectionId = `conn-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
    store.dispatch(setConnectionId(connectionId));
    
    // Start health monitoring
    this.startHealthMonitoring();
    
    // Send queued messages
    this.sendQueuedMessages();
    
    this.isReconnecting = false;
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
    console.log('WebSocket disconnected:', event.code, event.reason);
    
    // Clean up intervals
    this.stopHealthMonitoring();
    
    // Update connection status
    store.dispatch(setConnectionStatus('disconnected'));
    store.dispatch(setWebSocket(null));
    store.dispatch(setHealthStatus(false));
    
    // Attempt to reconnect if not intentionally closed
    if (event.code !== 1000 && !this.isReconnecting) {
      this.reconnect();
    }
  }

  private handleError(error: Event): void {
    console.error('WebSocket error:', error);
    
    const errorMessage = 'WebSocket connection error';
    store.dispatch(setConnectionError(errorMessage));
    store.dispatch(setHealthStatus(false));
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
    // Start ping interval
    this.pingIntervalId = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        this.lastPingTime = Date.now();
        store.dispatch(setPingTime(this.lastPingTime));
        this.send({ type: 'ping', timestamp: new Date().toISOString() });
      }
    }, this.pingInterval);

    // Start health check interval
    this.healthCheckIntervalId = setInterval(() => {
      const state = store.getState().connection;
      
      // Check if connection is stale (no pong received)
      if (state.lastPingTime && Date.now() - state.lastPingTime > this.pingInterval * 2) {
        console.warn('Connection appears stale, forcing reconnection');
        store.dispatch(setHealthStatus(false));
        this.forceReconnect();
      }
    }, this.healthCheckInterval);
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
    // Clear timeouts and intervals
    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId);
      this.reconnectTimeoutId = null;
    }
    
    this.stopHealthMonitoring();
    
    // Close WebSocket connection
    if (this.ws) {
      try {
        this.ws.close(1000, 'Client disconnect');
      } catch (error) {
        console.warn('Error closing WebSocket:', error);
      }
      this.ws = null;
    }
    
    // Clear message queue
    this.messageQueue = [];
    
    // Update store
    store.dispatch(cleanup());
  }
}

// Create singleton instance
export const websocketService = new WebSocketService();

// Export for use in components
export default websocketService;
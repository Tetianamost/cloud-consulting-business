/**
 * Simple WebSocket Service
 * 
 * This is a simplified WebSocket service designed to work around React lifecycle issues
 * that cause the main WebSocket service to disconnect immediately after connecting.
 * 
 * Key differences from the main service:
 * - Minimal React integration to avoid lifecycle conflicts
 * - Simplified connection management
 * - Focus on connection stability over advanced features
 */

export interface SimpleWebSocketMessage {
  type: 'message' | 'heartbeat' | 'ping' | 'pong';
  content?: string;
  timestamp?: string;
}

class SimpleWebSocketService {
  private ws: WebSocket | null = null;
  private reconnectTimeoutId: NodeJS.Timeout | null = null;
  private heartbeatIntervalId: NodeJS.Timeout | null = null;
  private isConnecting = false;
  private isDestroyed = false;
  private connectionAttempts = 0;
  private maxReconnectAttempts = 5;
  private reconnectDelay = 1000;

  constructor() {
    console.log('[SimpleWebSocket] Service created');
    
    // Handle page unload
    window.addEventListener('beforeunload', () => {
      this.destroy();
    });
  }

  /**
   * Connect to WebSocket
   */
  public async connect(): Promise<void> {
    if (this.isDestroyed) {
      console.log('[SimpleWebSocket] Service is destroyed, cannot connect');
      return;
    }

    if (this.isConnecting) {
      console.log('[SimpleWebSocket] Already connecting, skipping');
      return;
    }

    if (this.ws?.readyState === WebSocket.OPEN) {
      console.log('[SimpleWebSocket] Already connected');
      return;
    }

    return new Promise((resolve, reject) => {
      try {
        this.isConnecting = true;
        this.cleanup();

        const wsUrl = this.getWebSocketUrl();
        console.log('[SimpleWebSocket] Connecting to:', wsUrl.replace(/token=[^&]+/, 'token=***'));

        this.ws = new WebSocket(wsUrl);

        const connectionTimeout = setTimeout(() => {
          if (this.ws?.readyState === WebSocket.CONNECTING) {
            console.error('[SimpleWebSocket] Connection timeout');
            this.ws.close();
            reject(new Error('Connection timeout'));
          }
        }, 10000);

        this.ws.onopen = () => {
          clearTimeout(connectionTimeout);
          this.isConnecting = false;
          this.connectionAttempts = 0;
          console.log('[SimpleWebSocket] Connected successfully');
          
          // Start heartbeat to keep connection alive
          this.startHeartbeat();
          
          resolve();
        };

        this.ws.onmessage = (event) => {
          try {
            const data = JSON.parse(event.data);
            console.log('[SimpleWebSocket] Received message:', data);
            this.handleMessage(data);
          } catch (error) {
            console.error('[SimpleWebSocket] Failed to parse message:', error);
          }
        };

        this.ws.onclose = (event) => {
          clearTimeout(connectionTimeout);
          this.isConnecting = false;
          console.log('[SimpleWebSocket] Connection closed:', {
            code: event.code,
            reason: event.reason,
            wasClean: event.wasClean
          });

          this.stopHeartbeat();

          // Attempt reconnection if not intentionally closed and not destroyed
          if (event.code !== 1000 && !this.isDestroyed && this.connectionAttempts < this.maxReconnectAttempts) {
            this.scheduleReconnect();
          } else if (this.connectionAttempts >= this.maxReconnectAttempts) {
            console.error('[SimpleWebSocket] Max reconnection attempts reached');
            reject(new Error('Max reconnection attempts reached'));
          }
        };

        this.ws.onerror = (error) => {
          clearTimeout(connectionTimeout);
          this.isConnecting = false;
          console.error('[SimpleWebSocket] Connection error:', error);
          reject(error);
        };

      } catch (error) {
        this.isConnecting = false;
        console.error('[SimpleWebSocket] Failed to create connection:', error);
        reject(error);
      }
    });
  }

  /**
   * Send message
   */
  public send(message: SimpleWebSocketMessage): boolean {
    if (this.ws?.readyState === WebSocket.OPEN) {
      try {
        this.ws.send(JSON.stringify(message));
        return true;
      } catch (error) {
        console.error('[SimpleWebSocket] Failed to send message:', error);
        return false;
      }
    } else {
      console.warn('[SimpleWebSocket] Cannot send message, not connected');
      return false;
    }
  }

  /**
   * Disconnect
   */
  public disconnect(): void {
    console.log('[SimpleWebSocket] Disconnecting');
    this.cleanup();
  }

  /**
   * Destroy service
   */
  public destroy(): void {
    console.log('[SimpleWebSocket] Destroying service');
    this.isDestroyed = true;
    this.cleanup();
  }

  /**
   * Get connection status
   */
  public isConnected(): boolean {
    return this.ws?.readyState === WebSocket.OPEN;
  }

  /**
   * Get WebSocket URL
   */
  private getWebSocketUrl(): string {
    const token = localStorage.getItem('adminToken');
    if (!token) {
      throw new Error('No admin token found');
    }

    // Use environment variable or fallback to localhost:8061
    const wsUrl = process.env.REACT_APP_WS_URL || 'ws://localhost:8061/api/v1/admin/chat/ws';
    return `${wsUrl}?token=${encodeURIComponent(token)}`;
  }

  /**
   * Handle incoming messages
   */
  private handleMessage(message: SimpleWebSocketMessage): void {
    switch (message.type) {
      case 'ping':
        // Respond to ping with pong
        this.send({ type: 'pong', timestamp: new Date().toISOString() });
        break;
      
      case 'pong':
        console.log('[SimpleWebSocket] Received pong');
        break;
      
      case 'message':
        console.log('[SimpleWebSocket] Received chat message:', message.content);
        break;
      
      default:
        console.log('[SimpleWebSocket] Unknown message type:', message.type);
    }
  }

  /**
   * Start heartbeat to keep connection alive
   */
  private startHeartbeat(): void {
    this.stopHeartbeat(); // Clear any existing heartbeat
    
    this.heartbeatIntervalId = setInterval(() => {
      if (this.ws?.readyState === WebSocket.OPEN) {
        console.log('[SimpleWebSocket] Sending heartbeat');
        this.send({ type: 'heartbeat', timestamp: new Date().toISOString() });
      }
    }, 30000); // Send heartbeat every 30 seconds
  }

  /**
   * Stop heartbeat
   */
  private stopHeartbeat(): void {
    if (this.heartbeatIntervalId) {
      clearInterval(this.heartbeatIntervalId);
      this.heartbeatIntervalId = null;
    }
  }

  /**
   * Schedule reconnection
   */
  private scheduleReconnect(): void {
    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId);
    }

    this.connectionAttempts++;
    const delay = this.reconnectDelay * Math.pow(2, this.connectionAttempts - 1); // Exponential backoff

    console.log(`[SimpleWebSocket] Scheduling reconnection in ${delay}ms (attempt ${this.connectionAttempts})`);

    this.reconnectTimeoutId = setTimeout(async () => {
      try {
        await this.connect();
      } catch (error) {
        console.error('[SimpleWebSocket] Reconnection failed:', error);
      }
    }, delay);
  }

  /**
   * Cleanup resources
   */
  private cleanup(): void {
    if (this.reconnectTimeoutId) {
      clearTimeout(this.reconnectTimeoutId);
      this.reconnectTimeoutId = null;
    }

    this.stopHeartbeat();

    if (this.ws) {
      // Remove event listeners to prevent them from firing during cleanup
      this.ws.onopen = null;
      this.ws.onmessage = null;
      this.ws.onclose = null;
      this.ws.onerror = null;

      if (this.ws.readyState === WebSocket.OPEN || this.ws.readyState === WebSocket.CONNECTING) {
        this.ws.close(1000, 'Client disconnect');
      }
      
      this.ws = null;
    }
  }
}

// Create singleton instance
export const simpleWebSocketService = new SimpleWebSocketService();

export default simpleWebSocketService;
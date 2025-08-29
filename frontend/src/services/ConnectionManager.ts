import { store } from '../store';
import {
  setConnectionStatus,
  setConnectionError,
  setHealthStatus,
  incrementReconnectAttempts,
  resetReconnectAttempts,
  updateLatency,
  ConnectionStatus,
} from '../store/slices/connectionSlice';

// Connection state for polling
export type PollingConnectionState = 'connected' | 'polling' | 'error' | 'offline';

export interface ConnectionManagerState {
  state: PollingConnectionState;
  lastSuccessfulPoll: Date | null;
  errorCount: number;
  retryAttempts: number;
  isPollingActive: boolean;
  lastErrorTime: Date | null;
  consecutiveSuccesses: number;
}

export interface ConnectionStatusInfo {
  state: PollingConnectionState;
  isHealthy: boolean;
  errorCount: number;
  lastError: string | null;
  uptime: number;
  lastSuccessfulConnection: Date | null;
}

/**
 * ConnectionManager class to track polling state and manage connection status
 */
export class ConnectionManager {
  private state: ConnectionManagerState;
  private startTime: Date;
  private statusUpdateCallbacks: Array<(status: ConnectionStatusInfo) => void> = [];

  constructor() {
    this.startTime = new Date();
    this.state = {
      state: 'offline',
      lastSuccessfulPoll: null,
      errorCount: 0,
      retryAttempts: 0,
      isPollingActive: false,
      lastErrorTime: null,
      consecutiveSuccesses: 0,
    };
  }

  /**
   * Start polling and update connection state
   */
  public startPolling(): void {
    console.log('[ConnectionManager] Starting polling');
    this.state.isPollingActive = true;
    this.state.state = 'polling';
    this.state.errorCount = 0;
    this.state.retryAttempts = 0;
    this.state.consecutiveSuccesses = 0;
    
    this.updateReduxState('connecting');
    this.notifyStatusChange();
  }

  /**
   * Stop polling and update connection state
   */
  public stopPolling(): void {
    console.log('[ConnectionManager] Stopping polling');
    this.state.isPollingActive = false;
    this.state.state = 'offline';
    
    this.updateReduxState('disconnected');
    this.notifyStatusChange();
  }

  /**
   * Record successful poll
   */
  public recordSuccessfulPoll(latency: number): void {
    const now = new Date();
    this.state.lastSuccessfulPoll = now;
    this.state.errorCount = 0;
    this.state.retryAttempts = 0;
    this.state.consecutiveSuccesses++;
    this.state.lastErrorTime = null;
    
    // Update state based on consecutive successes
    if (this.state.consecutiveSuccesses >= 3) {
      this.state.state = 'connected';
    } else {
      this.state.state = 'polling';
    }
    
    // Update Redux state
    store.dispatch(updateLatency(latency));
    
    const currentStatus = store.getState().connection.status;
    if (currentStatus !== 'connected' && this.state.state === 'connected') {
      console.log('[ConnectionManager] Connection established');
      this.updateReduxState('connected');
      store.dispatch(resetReconnectAttempts());
      store.dispatch(setConnectionError(null));
      store.dispatch(setHealthStatus(true));
    }
    
    this.notifyStatusChange();
  }

  /**
   * Record polling error
   */
  public recordPollingError(error: string): void {
    this.state.errorCount++;
    this.state.retryAttempts++;
    this.state.lastErrorTime = new Date();
    this.state.consecutiveSuccesses = 0;
    this.state.state = 'error';
    
    console.error(`[ConnectionManager] Polling error (count: ${this.state.errorCount}):`, error);
    
    // Update Redux state based on error severity
    if (this.state.errorCount === 1) {
      this.updateReduxState('reconnecting');
      store.dispatch(setConnectionError('Connection issue detected. Retrying...'));
    } else if (this.state.errorCount >= 5) {
      this.updateReduxState('failed');
      store.dispatch(setConnectionError('Connection failed. Will retry automatically.'));
    } else {
      this.updateReduxState('reconnecting');
      store.dispatch(setConnectionError(`Connection unstable (${this.state.errorCount} errors). Retrying...`));
    }
    
    store.dispatch(incrementReconnectAttempts());
    store.dispatch(setHealthStatus(false));
    
    this.notifyStatusChange();
  }

  /**
   * Handle network offline state
   */
  public setOffline(): void {
    console.log('[ConnectionManager] Network offline');
    this.state.state = 'offline';
    this.state.lastErrorTime = new Date();
    
    this.updateReduxState('reconnecting');
    store.dispatch(setConnectionError('Network offline. Will reconnect automatically.'));
    store.dispatch(setHealthStatus(false));
    
    this.notifyStatusChange();
  }

  /**
   * Handle network online state
   */
  public setOnline(): void {
    console.log('[ConnectionManager] Network online');
    if (this.state.isPollingActive) {
      this.state.state = 'polling';
      this.state.errorCount = Math.max(0, this.state.errorCount - 1); // Reduce error count
      
      this.updateReduxState('reconnecting');
      store.dispatch(setConnectionError(null));
    }
    
    this.notifyStatusChange();
  }

  /**
   * Get current connection status information
   */
  public getStatusInfo(): ConnectionStatusInfo {
    const uptime = Date.now() - this.startTime.getTime();
    
    return {
      state: this.state.state,
      isHealthy: this.isHealthy(),
      errorCount: this.state.errorCount,
      lastError: store.getState().connection.error,
      uptime,
      lastSuccessfulConnection: this.state.lastSuccessfulPoll,
    };
  }

  /**
   * Check if connection is healthy
   */
  public isHealthy(): boolean {
    // Consider healthy if:
    // 1. Connected state, or
    // 2. Polling state with recent successful poll (within last 30 seconds), or
    // 3. Low error count (< 3)
    if (this.state.state === 'connected') {
      return true;
    }
    
    if (this.state.state === 'polling' && this.state.lastSuccessfulPoll) {
      const timeSinceLastSuccess = Date.now() - this.state.lastSuccessfulPoll.getTime();
      return timeSinceLastSuccess < 30000; // 30 seconds
    }
    
    return this.state.errorCount < 3;
  }

  /**
   * Get appropriate status message for users
   */
  public getStatusMessage(): string {
    switch (this.state.state) {
      case 'connected':
        return 'Connected';
      case 'polling':
        if (this.state.consecutiveSuccesses > 0) {
          return 'Connected (polling)';
        }
        return 'Connecting...';
      case 'error':
        if (this.state.errorCount >= 5) {
          return 'Connection failed - retrying automatically';
        } else if (this.state.errorCount >= 3) {
          return 'Connection unstable - retrying';
        }
        return 'Connection issues - retrying';
      case 'offline':
        return 'Offline - will reconnect when network is available';
      default:
        return 'Disconnected';
    }
  }

  /**
   * Get connection state
   */
  public getState(): PollingConnectionState {
    return this.state.state;
  }

  /**
   * Check if polling is active
   */
  public isPollingActive(): boolean {
    return this.state.isPollingActive;
  }

  /**
   * Get error count
   */
  public getErrorCount(): number {
    return this.state.errorCount;
  }

  /**
   * Reset error count (used when connection is restored)
   */
  public resetErrorCount(): void {
    this.state.errorCount = 0;
    this.state.retryAttempts = 0;
    this.state.lastErrorTime = null;
    this.notifyStatusChange();
  }

  /**
   * Subscribe to status changes
   */
  public onStatusChange(callback: (status: ConnectionStatusInfo) => void): () => void {
    this.statusUpdateCallbacks.push(callback);
    
    // Return unsubscribe function
    return () => {
      const index = this.statusUpdateCallbacks.indexOf(callback);
      if (index > -1) {
        this.statusUpdateCallbacks.splice(index, 1);
      }
    };
  }

  /**
   * Force reconnection attempt
   */
  public forceReconnect(): void {
    console.log('[ConnectionManager] Forcing reconnection');
    this.state.errorCount = 0;
    this.state.retryAttempts = 0;
    this.state.consecutiveSuccesses = 0;
    this.state.state = 'polling';
    
    this.updateReduxState('reconnecting');
    store.dispatch(setConnectionError(null));
    
    this.notifyStatusChange();
  }

  // Private methods

  /**
   * Update Redux connection state
   */
  private updateReduxState(status: ConnectionStatus): void {
    store.dispatch(setConnectionStatus(status));
  }

  /**
   * Notify all status change callbacks
   */
  private notifyStatusChange(): void {
    const statusInfo = this.getStatusInfo();
    this.statusUpdateCallbacks.forEach(callback => {
      try {
        callback(statusInfo);
      } catch (error) {
        console.error('[ConnectionManager] Error in status change callback:', error);
      }
    });
  }
}

// Export singleton instance
export const connectionManager = new ConnectionManager();
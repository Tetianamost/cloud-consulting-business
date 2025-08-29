import { store } from '../store';
import { setConnectionStatus, setConnectionError } from '../store/slices/connectionSlice';
import websocketService from './websocketService';
import { pollingChatService } from './pollingChatService';

// Chat mode types
export type ChatMode = 'websocket' | 'polling' | 'auto';

// Chat configuration interface
export interface ChatConfig {
  mode: ChatMode;
  enable_websocket_fallback: boolean;
  websocket_timeout: number;
  polling_interval: number;
  max_reconnect_attempts: number;
  fallback_delay: number;
}

// Fallback state interface
interface FallbackState {
  isInFallback: boolean;
  fallbackReason: string | null;
  websocketFailureCount: number;
  lastWebSocketAttempt: Date | null;
  fallbackStartTime: Date | null;
}

// Performance metrics interface
interface PerformanceMetrics {
  websocket: {
    connectionTime: number;
    messageLatency: number;
    disconnectionCount: number;
    lastDisconnection: Date | null;
  };
  polling: {
    averageResponseTime: number;
    successRate: number;
    errorCount: number;
    lastError: Date | null;
  };
}

class ChatModeManager {
  private currentMode: ChatMode = 'auto';
  private config: ChatConfig | null = null;
  private fallbackState: FallbackState = {
    isInFallback: false,
    fallbackReason: null,
    websocketFailureCount: 0,
    lastWebSocketAttempt: null,
    fallbackStartTime: null,
  };
  private performanceMetrics: PerformanceMetrics = {
    websocket: {
      connectionTime: 0,
      messageLatency: 0,
      disconnectionCount: 0,
      lastDisconnection: null,
    },
    polling: {
      averageResponseTime: 0,
      successRate: 100,
      errorCount: 0,
      lastError: null,
    },
  };
  private activeService: 'websocket' | 'polling' | null = null;
  private fallbackTimeoutId: NodeJS.Timeout | null = null;
  private statusChangeCallbacks: Array<() => void> = [];

  constructor() {
    console.log('[ChatModeManager] Initializing chat mode manager');
    this.loadConfiguration();
    this.setupEventListeners();
  }

  /**
   * Load chat configuration from backend
   */
  private async loadConfiguration(): Promise<void> {
    try {
      const response = await fetch('/api/v1/admin/chat/config', {
        headers: {
          'Authorization': `Bearer ${localStorage.getItem('adminToken')}`,
        },
      });

      if (response.ok) {
        const data = await response.json();
        this.config = data.config;
        this.currentMode = this.config?.mode || 'auto';
        console.log('[ChatModeManager] Configuration loaded:', this.config);
      } else {
        console.warn('[ChatModeManager] Failed to load configuration, using defaults');
        this.useDefaultConfiguration();
      }
    } catch (error) {
      console.error('[ChatModeManager] Error loading configuration:', error);
      this.useDefaultConfiguration();
    }
  }

  /**
   * Use default configuration if backend config is unavailable
   */
  private useDefaultConfiguration(): void {
    this.config = {
      mode: 'polling', // Default to polling since WebSocket has issues
      enable_websocket_fallback: true,
      websocket_timeout: 10,
      polling_interval: 3000,
      max_reconnect_attempts: 3,
      fallback_delay: 5000,
    };
    this.currentMode = 'polling';
    console.log('[ChatModeManager] Using default configuration - FORCING POLLING MODE');
  }

  /**
   * Setup event listeners for connection monitoring
   */
  private setupEventListeners(): void {
    // Listen for polling service status changes
    try {
      const unsubscribePolling = pollingChatService.onStatusChange((status: any) => {
        this.handlePollingStatusChange(status);
      });

      // Store unsubscribe functions for cleanup
      if (unsubscribePolling && typeof unsubscribePolling === 'function') {
        this.statusChangeCallbacks.push(unsubscribePolling);
      }
    } catch (error) {
      console.warn('[ChatModeManager] Failed to setup polling status listener:', error);
    }
  }

  /**
   * Initialize chat service based on current mode
   */
  public async initializeChatService(): Promise<void> {
    if (!this.config) {
      await this.loadConfiguration();
    }

    // FORCE POLLING MODE FOR NOW
    this.currentMode = 'polling';
    console.log(`[ChatModeManager] FORCING POLLING MODE - Initializing chat service with mode: ${this.currentMode}`);

    // Force polling mode
    await this.initializePollingService();
  }

  /**
   * Initialize WebSocket service
   */
  private async initializeWebSocketService(): Promise<void> {
    try {
      console.log('[ChatModeManager] Initializing WebSocket service');
      this.activeService = 'websocket';
      
      const startTime = Date.now();
      await websocketService.connect();
      
      // Record connection time
      this.performanceMetrics.websocket.connectionTime = Date.now() - startTime;
      
      // Reset failure count on successful connection
      this.fallbackState.websocketFailureCount = 0;
      
      console.log('[ChatModeManager] WebSocket service initialized successfully');
    } catch (error) {
      console.error('[ChatModeManager] WebSocket initialization failed:', error);
      this.handleWebSocketFailure(error as Error);
    }
  }

  /**
   * Initialize polling service
   */
  private async initializePollingService(): Promise<void> {
    try {
      console.log('[ChatModeManager] Initializing polling service');
      this.activeService = 'polling';
      
      pollingChatService.startPolling();
      
      console.log('[ChatModeManager] Polling service initialized successfully');
    } catch (error) {
      console.error('[ChatModeManager] Polling initialization failed:', error);
      store.dispatch(setConnectionError(`Failed to initialize polling service: ${error}`));
    }
  }

  /**
   * Initialize auto mode (try WebSocket first, fallback to polling)
   */
  private async initializeAutoMode(): Promise<void> {
    console.log('[ChatModeManager] Initializing auto mode');
    
    // If already in fallback mode, use polling
    if (this.fallbackState.isInFallback) {
      await this.initializePollingService();
      return;
    }

    // Try WebSocket first
    try {
      await this.initializeWebSocketService();
    } catch (error) {
      console.log('[ChatModeManager] WebSocket failed in auto mode, falling back to polling');
      await this.triggerFallback('WebSocket connection failed in auto mode');
    }
  }

  /**
   * Handle WebSocket status changes
   */
  private handleWebSocketStatusChange(status: any): void {
    console.log('[ChatModeManager] WebSocket status change:', status);

    if (status.status === 'disconnected' || status.status === 'failed') {
      this.performanceMetrics.websocket.disconnectionCount++;
      this.performanceMetrics.websocket.lastDisconnection = new Date();

      // Check if we should trigger fallback
      if (this.shouldTriggerFallback()) {
        this.triggerFallback('WebSocket connection repeatedly failed');
      }
    }
  }

  /**
   * Handle polling status changes
   */
  private handlePollingStatusChange(status: any): void {
    console.log('[ChatModeManager] Polling status change:', status);

    if (status.error) {
      this.performanceMetrics.polling.errorCount++;
      this.performanceMetrics.polling.lastError = new Date();
    }
  }

  /**
   * Handle WebSocket failure
   */
  private handleWebSocketFailure(error: Error): void {
    this.fallbackState.websocketFailureCount++;
    this.fallbackState.lastWebSocketAttempt = new Date();

    console.log(`[ChatModeManager] WebSocket failure #${this.fallbackState.websocketFailureCount}:`, error);

    // Check if we should trigger fallback
    if (this.shouldTriggerFallback()) {
      this.triggerFallback(`WebSocket failed ${this.fallbackState.websocketFailureCount} times`);
    }
  }

  /**
   * Determine if fallback should be triggered
   */
  private shouldTriggerFallback(): boolean {
    if (!this.config?.enable_websocket_fallback) {
      return false;
    }

    // Don't fallback if already in fallback mode
    if (this.fallbackState.isInFallback) {
      return false;
    }

    // Don't fallback if mode is explicitly set to websocket
    if (this.currentMode === 'websocket') {
      return false;
    }

    // Trigger fallback if we've exceeded max reconnect attempts
    return this.fallbackState.websocketFailureCount >= (this.config?.max_reconnect_attempts || 3);
  }

  /**
   * Trigger fallback from WebSocket to polling
   */
  private async triggerFallback(reason: string): Promise<void> {
    if (this.fallbackState.isInFallback) {
      console.log('[ChatModeManager] Already in fallback mode, ignoring trigger');
      return;
    }

    console.log(`[ChatModeManager] Triggering fallback to polling: ${reason}`);

    // Update fallback state
    this.fallbackState.isInFallback = true;
    this.fallbackState.fallbackReason = reason;
    this.fallbackState.fallbackStartTime = new Date();

    // Disconnect WebSocket service
    websocketService.disconnect();

    // Show user notification
    store.dispatch(setConnectionStatus('reconnecting'));
    store.dispatch(setConnectionError(`Switching to polling mode: ${reason}`));

    // Wait for fallback delay before switching
    if (this.fallbackTimeoutId) {
      clearTimeout(this.fallbackTimeoutId);
    }

    this.fallbackTimeoutId = setTimeout(async () => {
      try {
        await this.initializePollingService();
        
        // Show success notification
        store.dispatch(setConnectionError('Successfully switched to polling mode'));
        
        // Clear error after a few seconds
        setTimeout(() => {
          store.dispatch(setConnectionError(null));
        }, 3000);
        
      } catch (error) {
        console.error('[ChatModeManager] Fallback to polling failed:', error);
        store.dispatch(setConnectionError('Failed to switch to polling mode'));
      }
    }, this.config?.fallback_delay || 5000);
  }

  /**
   * Manually switch chat mode
   */
  public async switchMode(mode: ChatMode): Promise<void> {
    console.log(`[ChatModeManager] Manually switching to mode: ${mode}`);

    // Stop current service
    this.stopCurrentService();

    // Update current mode
    this.currentMode = mode;

    // Reset fallback state if switching manually
    this.resetFallbackState();

    // Initialize new service
    await this.initializeChatService();
  }

  /**
   * Stop current active service
   */
  private stopCurrentService(): void {
    if (this.activeService === 'websocket') {
      websocketService.disconnect();
    } else if (this.activeService === 'polling') {
      pollingChatService.stopPolling();
    }
    this.activeService = null;
  }

  /**
   * Reset fallback state
   */
  private resetFallbackState(): void {
    this.fallbackState = {
      isInFallback: false,
      fallbackReason: null,
      websocketFailureCount: 0,
      lastWebSocketAttempt: null,
      fallbackStartTime: null,
    };

    if (this.fallbackTimeoutId) {
      clearTimeout(this.fallbackTimeoutId);
      this.fallbackTimeoutId = null;
    }
  }

  /**
   * Get current chat mode
   */
  public getCurrentMode(): ChatMode {
    return this.currentMode;
  }

  /**
   * Get active service
   */
  public getActiveService(): 'websocket' | 'polling' | null {
    return this.activeService;
  }

  /**
   * Get fallback state
   */
  public getFallbackState(): FallbackState {
    return { ...this.fallbackState };
  }

  /**
   * Get performance metrics
   */
  public getPerformanceMetrics(): PerformanceMetrics {
    return { ...this.performanceMetrics };
  }

  /**
   * Get configuration
   */
  public getConfiguration(): ChatConfig | null {
    return this.config ? { ...this.config } : null;
  }

  /**
   * Update configuration
   */
  public async updateConfiguration(updates: Partial<ChatConfig>): Promise<void> {
    try {
      const response = await fetch('/api/v1/admin/chat/config', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('adminToken')}`,
        },
        body: JSON.stringify(updates),
      });

      if (response.ok) {
        const data = await response.json();
        this.config = data.config;
        console.log('[ChatModeManager] Configuration updated:', this.config);
      } else {
        throw new Error('Failed to update configuration');
      }
    } catch (error) {
      console.error('[ChatModeManager] Error updating configuration:', error);
      throw error;
    }
  }

  /**
   * Get status message for UI
   */
  public getStatusMessage(): string {
    if (this.fallbackState.isInFallback) {
      return `Using polling mode (${this.fallbackState.fallbackReason})`;
    }

    switch (this.activeService) {
      case 'websocket':
        return 'Connected via WebSocket';
      case 'polling':
        return 'Connected via polling';
      default:
        return 'Connecting...';
    }
  }

  /**
   * Check if service is healthy
   */
  public isHealthy(): boolean {
    if (this.activeService === 'websocket') {
      return websocketService.isHealthy();
    } else if (this.activeService === 'polling') {
      return pollingChatService.isHealthy();
    }
    return false;
  }

  /**
   * Force reconnection
   */
  public forceReconnect(): void {
    console.log('[ChatModeManager] Force reconnect requested');
    
    if (this.activeService === 'websocket') {
      websocketService.forceReconnect();
    } else if (this.activeService === 'polling') {
      pollingChatService.forceReconnect();
    }
  }

  /**
   * Cleanup resources
   */
  public cleanup(): void {
    console.log('[ChatModeManager] Cleaning up resources');
    
    this.stopCurrentService();
    this.resetFallbackState();
    
    // Cleanup event listeners
    this.statusChangeCallbacks.forEach(callback => {
      if (typeof callback === 'function') {
        callback();
      }
    });
    this.statusChangeCallbacks = [];
  }
}

// Create singleton instance
export const chatModeManager = new ChatModeManager();

export default chatModeManager;
import { store } from '../store';
import { setConnectionStatus, setConnectionError } from '../store/slices/connectionSlice';
import { pollingChatService } from './pollingChatService';

// Chat mode types - simplified to only polling
export type ChatMode = 'polling';

// Chat configuration interface - simplified for polling only
export interface ChatConfig {
  mode: ChatMode;
  polling_interval: number;
  max_reconnect_attempts: number;
}

// Performance metrics interface - polling only
interface PerformanceMetrics {
  polling: {
    averageResponseTime: number;
    successRate: number;
    errorCount: number;
    lastError: Date | null;
  };
}

class ChatModeManager {
  private currentMode: ChatMode = 'polling';
  private config: ChatConfig | null = null;
  private performanceMetrics: PerformanceMetrics = {
    polling: {
      averageResponseTime: 0,
      successRate: 100,
      errorCount: 0,
      lastError: null,
    },
  };
  private activeService: 'polling' | null = null;
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
        this.currentMode = this.config?.mode || 'polling';
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
      mode: 'polling',
      polling_interval: 3000,
      max_reconnect_attempts: 3,
    };
    this.currentMode = 'polling';
    console.log('[ChatModeManager] Using default polling configuration');
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
   * Initialize chat service - polling only
   */
  public async initializeChatService(): Promise<void> {
    if (!this.config) {
      await this.loadConfiguration();
    }

    this.currentMode = 'polling';
    console.log(`[ChatModeManager] Initializing polling chat service`);

    await this.initializePollingService();
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
   * Stop current active service
   */
  private stopCurrentService(): void {
    if (this.activeService === 'polling') {
      pollingChatService.stopPolling();
    }
    this.activeService = null;
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
  public getActiveService(): 'polling' | null {
    return this.activeService;
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
    switch (this.activeService) {
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
    if (this.activeService === 'polling') {
      return pollingChatService.isHealthy();
    }
    return false;
  }

  /**
   * Force reconnection
   */
  public forceReconnect(): void {
    console.log('[ChatModeManager] Force reconnect requested');
    
    if (this.activeService === 'polling') {
      pollingChatService.forceReconnect();
    }
  }

  /**
   * Cleanup resources
   */
  public cleanup(): void {
    console.log('[ChatModeManager] Cleaning up resources');
    
    this.stopCurrentService();
    
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
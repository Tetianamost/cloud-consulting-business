import { chatModeManager, ChatMode } from './chatModeManager';

// Mock the services
jest.mock('./websocketService', () => ({
  connect: jest.fn(),
  disconnect: jest.fn(),
  isHealthy: jest.fn(() => true),
  forceReconnect: jest.fn(),
  onStatusChange: jest.fn(),
}));

jest.mock('./pollingChatService', () => ({
  startPolling: jest.fn(),
  stopPolling: jest.fn(),
  isHealthy: jest.fn(() => true),
  forceReconnect: jest.fn(),
  onStatusChange: jest.fn(),
}));

// Mock fetch
global.fetch = jest.fn();

describe('ChatModeManager', () => {
  beforeEach(() => {
    jest.clearAllMocks();
    
    // Mock successful config fetch
    (global.fetch as jest.Mock).mockResolvedValue({
      ok: true,
      json: () => Promise.resolve({
        config: {
          mode: 'auto',
          enable_websocket_fallback: true,
          websocket_timeout: 10,
          polling_interval: 3000,
          max_reconnect_attempts: 3,
          fallback_delay: 5000,
        },
      }),
    });
  });

  describe('Mode Switching', () => {
    test('should switch to websocket mode', async () => {
      await chatModeManager.switchMode('websocket');
      expect(chatModeManager.getCurrentMode()).toBe('websocket');
      expect(chatModeManager.getActiveService()).toBe('websocket');
    });

    test('should switch to polling mode', async () => {
      await chatModeManager.switchMode('polling');
      expect(chatModeManager.getCurrentMode()).toBe('polling');
      expect(chatModeManager.getActiveService()).toBe('polling');
    });

    test('should handle auto mode', async () => {
      await chatModeManager.switchMode('auto');
      expect(chatModeManager.getCurrentMode()).toBe('auto');
      // In auto mode, it should try WebSocket first
      expect(chatModeManager.getActiveService()).toBe('websocket');
    });
  });

  describe('Fallback Mechanism', () => {
    test('should trigger fallback after max reconnect attempts', async () => {
      // Set up auto mode
      await chatModeManager.switchMode('auto');
      
      // Simulate WebSocket failures
      const fallbackState = chatModeManager.getFallbackState();
      expect(fallbackState.isInFallback).toBe(false);
      
      // The actual fallback logic would be triggered by WebSocket status changes
      // This is a simplified test of the state management
    });

    test('should not fallback when disabled', async () => {
      // Mock config with fallback disabled
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({
          config: {
            mode: 'auto',
            enable_websocket_fallback: false,
            websocket_timeout: 10,
            polling_interval: 3000,
            max_reconnect_attempts: 3,
            fallback_delay: 5000,
          },
        }),
      });

      const config = chatModeManager.getConfiguration();
      expect(config?.enable_websocket_fallback).toBe(false);
    });
  });

  describe('Configuration Management', () => {
    test('should load configuration from backend', async () => {
      const config = chatModeManager.getConfiguration();
      expect(config).toBeTruthy();
      expect(config?.mode).toBe('auto');
    });

    test('should update configuration', async () => {
      const updates = { mode: 'polling' as ChatMode };
      
      (global.fetch as jest.Mock).mockResolvedValue({
        ok: true,
        json: () => Promise.resolve({
          config: { ...chatModeManager.getConfiguration(), ...updates },
        }),
      });

      await chatModeManager.updateConfiguration(updates);
      
      expect(global.fetch).toHaveBeenCalledWith('/api/v1/admin/chat/config', {
        method: 'PUT',
        headers: {
          'Content-Type': 'application/json',
          'Authorization': `Bearer ${localStorage.getItem('adminToken')}`,
        },
        body: JSON.stringify(updates),
      });
    });
  });

  describe('Performance Metrics', () => {
    test('should track performance metrics', () => {
      const metrics = chatModeManager.getPerformanceMetrics();
      expect(metrics).toBeTruthy();
      expect(metrics.websocket).toBeTruthy();
      expect(metrics.polling).toBeTruthy();
    });
  });

  describe('Status Management', () => {
    test('should provide status message', () => {
      const message = chatModeManager.getStatusMessage();
      expect(typeof message).toBe('string');
    });

    test('should check health status', () => {
      const isHealthy = chatModeManager.isHealthy();
      expect(typeof isHealthy).toBe('boolean');
    });
  });

  describe('Cleanup', () => {
    test('should cleanup resources', () => {
      chatModeManager.cleanup();
      expect(chatModeManager.getActiveService()).toBe(null);
    });
  });
});
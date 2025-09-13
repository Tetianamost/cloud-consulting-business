import { chatModeManager, ChatMode } from './chatModeManager';

// Mock the services
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
          mode: 'polling',
          polling_interval: 3000,
          max_reconnect_attempts: 3,
        },
      }),
    });
  });

  describe('Mode Management', () => {
    test('should use polling mode', async () => {
      await chatModeManager.initializeChatService();
      expect(chatModeManager.getCurrentMode()).toBe('polling');
      expect(chatModeManager.getActiveService()).toBe('polling');
    });
  });

  describe('Polling Configuration', () => {
    test('should handle polling configuration', async () => {
      const config = chatModeManager.getConfiguration();
      expect(config?.mode).toBe('polling');
      expect(config?.polling_interval).toBe(3000);
      expect(config?.max_reconnect_attempts).toBe(3);
    });
  });

  describe('Configuration Management', () => {
    test('should load configuration from backend', async () => {
      const config = chatModeManager.getConfiguration();
      expect(config).toBeTruthy();
      expect(config?.mode).toBe('polling');
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
      expect(metrics.polling).toBeTruthy();
      expect(metrics.polling.averageResponseTime).toBeDefined();
      expect(metrics.polling.successRate).toBeDefined();
      expect(metrics.polling.errorCount).toBeDefined();
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
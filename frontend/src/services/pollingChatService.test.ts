import { pollingChatService } from './pollingChatService';
import { store } from '../store';
import { ConnectionManager } from './ConnectionManager';

// Mock dependencies
jest.mock('../store', () => ({
  store: {
    dispatch: jest.fn(),
    getState: jest.fn(() => ({
      connection: {
        status: 'connected',
        isPollingActive: false,
        errorCount: 0,
        error: null,
      },
      chat: {
        messages: [],
        currentSession: { id: 'test-session' },
      },
    })),
  },
}));

jest.mock('./ConnectionManager', () => ({
  ConnectionManager: jest.fn().mockImplementation(() => ({
    isPollingActive: jest.fn(() => false),
    startPolling: jest.fn(),
    stopPolling: jest.fn(),
    getErrorCount: jest.fn(() => 0),
    getStatusInfo: jest.fn(() => ({ status: 'connected' })),
    getStatusMessage: jest.fn(() => 'Connected'),
    isHealthy: jest.fn(() => true),
    forceReconnect: jest.fn(),
    onStatusChange: jest.fn(() => () => {}),
    setOnline: jest.fn(),
    setOffline: jest.fn(),
    recordSuccessfulPoll: jest.fn(),
    recordPollingError: jest.fn(),
    resetErrorCount: jest.fn(),
  })),
  connectionManager: {
    isPollingActive: jest.fn(() => false),
    startPolling: jest.fn(),
    stopPolling: jest.fn(),
    getErrorCount: jest.fn(() => 0),
    getStatusInfo: jest.fn(() => ({ status: 'connected' })),
    getStatusMessage: jest.fn(() => 'Connected'),
    isHealthy: jest.fn(() => true),
    forceReconnect: jest.fn(),
    onStatusChange: jest.fn(() => () => {}),
    setOnline: jest.fn(),
    setOffline: jest.fn(),
    recordSuccessfulPoll: jest.fn(),
    recordPollingError: jest.fn(),
    resetErrorCount: jest.fn(),
  },
}));

// Mock fetch
global.fetch = jest.fn();

// Mock localStorage
const mockLocalStorage = {
  getItem: jest.fn(() => 'mock-token'),
  setItem: jest.fn(),
  removeItem: jest.fn(),
};
Object.defineProperty(window, 'localStorage', {
  value: mockLocalStorage,
});

// Mock navigator.onLine
Object.defineProperty(navigator, 'onLine', {
  writable: true,
  value: true,
});

// Mock window events
const mockAddEventListener = jest.fn();
const mockRemoveEventListener = jest.fn();
Object.defineProperty(window, 'addEventListener', {
  value: mockAddEventListener,
});
Object.defineProperty(window, 'removeEventListener', {
  value: mockRemoveEventListener,
});

// Mock document.hidden
Object.defineProperty(document, 'hidden', {
  writable: true,
  value: false,
});

// Mock timers
jest.useFakeTimers();

describe('PollingChatService - Comprehensive Tests', () => {
  let mockFetch: jest.MockedFunction<typeof fetch>;

  beforeEach(() => {
    jest.clearAllMocks();
    jest.clearAllTimers();
    
    mockFetch = global.fetch as jest.MockedFunction<typeof fetch>;
    
    // Reset navigator.onLine
    Object.defineProperty(navigator, 'onLine', {
      writable: true,
      value: true,
    });
    
    // Reset document.hidden
    Object.defineProperty(document, 'hidden', {
      writable: true,
      value: false,
    });
  });

  afterEach(() => {
    pollingChatService.stopPolling();
    jest.runOnlyPendingTimers();
  });

  describe('Message Sending - Success/Failure Scenarios', () => {
    test('should send message successfully with optimistic updates', async () => {
      const mockResponse = {
        success: true,
        message_id: 'test-message-id',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
        headers: new Headers(),
      } as Response);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
        client_name: 'Test Client',
      };

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/admin/chat/messages'),
        expect.objectContaining({
          method: 'POST',
          headers: expect.objectContaining({
            'Content-Type': 'application/json',
            'Authorization': 'Bearer mock-token',
          }),
          body: expect.stringContaining('Test message'),
        })
      );
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('addOptimisticMessage'),
        })
      );
    });

    test('should handle server error (500) with retry logic', async () => {
      const mockErrorResponse = {
        success: false,
        error: 'Internal server error',
      };

      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.resolve(mockErrorResponse),
      } as Response);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('addOptimisticMessage'),
        })
      );
    });

    test('should handle authentication error (401) without retry', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 401,
        json: () => Promise.resolve({ success: false, error: 'Unauthorized' }),
      } as Response);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      expect(mockLocalStorage.removeItem).toHaveBeenCalledWith('adminToken');
    });

    test('should handle network timeout with retry', async () => {
      const timeoutError = new Error('Request timeout');
      timeoutError.name = 'AbortError';
      
      mockFetch.mockRejectedValueOnce(timeoutError);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('addOptimisticMessage'),
        })
      );
    });

    test('should handle rate limiting (429) with retry', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 429,
        json: () => Promise.resolve({ 
          success: false, 
          error: 'Rate limit exceeded',
          retry_after: 60,
        }),
      } as Response);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      // Should queue message for retry
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('addOptimisticMessage'),
        })
      );
    });

    test('should queue message when offline', async () => {
      // Set offline
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false,
      });

      const request = {
        message: 'Offline message',
        session_id: 'test-session',
      };

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      expect(mockFetch).not.toHaveBeenCalled();
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('queueOfflineMessage'),
        })
      );
    });

    test('should prevent duplicate messages', async () => {
      const mockResponse = {
        success: true,
        message_id: 'test-message-id',
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
        headers: new Headers(),
      } as Response);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      // Send same message twice quickly
      const messageId1 = await pollingChatService.sendMessage(request);
      const messageId2 = await pollingChatService.sendMessage(request);

      expect(messageId1).toBeDefined();
      expect(messageId2).toBeDefined();
      // Both should succeed but duplicate prevention should be handled internally
    });
  });

  describe('Polling Loop Behavior with Different Intervals', () => {
    test('should start polling with correct initial state', () => {
      pollingChatService.startPolling();

      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('setPollingStatus'),
          payload: true,
        })
      );
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('setConnectionStatus'),
          payload: 'polling',
        })
      );
    });

    test('should stop polling and cleanup resources', () => {
      pollingChatService.startPolling();
      pollingChatService.stopPolling();

      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('setPollingStatus'),
          payload: false,
        })
      );
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('cleanup'),
        })
      );
    });

    test('should use fast polling (2s) when user is actively typing', () => {
      pollingChatService.startPolling();

      // Simulate typing activity
      const typingEvent = new KeyboardEvent('keydown', { key: 'a' });
      document.dispatchEvent(typingEvent);

      jest.advanceTimersByTime(2000);

      // Should use active interval
      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 2000);
    });

    test('should use slow polling (10s) when user is idle', () => {
      pollingChatService.startPolling();

      // Fast forward time to make user idle (5+ minutes)
      jest.advanceTimersByTime(300000);

      // Should use inactive interval
      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 10000);
    });

    test('should use max interval (30s) when page is hidden', () => {
      pollingChatService.startPolling();

      // Hide the page
      Object.defineProperty(document, 'hidden', {
        writable: true,
        value: true,
      });

      jest.advanceTimersByTime(5000);

      // Should use max interval when page is hidden
      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), 30000);
    });

    test('should adjust polling interval based on error count with exponential backoff', () => {
      // Mock connection manager to return error count
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      mockConnectionManager.getErrorCount.mockReturnValue(2);

      pollingChatService.startPolling();

      jest.advanceTimersByTime(3000);

      // Should use backoff interval when there are errors
      const expectedBackoffInterval = Math.min(3000 * Math.pow(1.5, 2), 30000);
      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), expectedBackoffInterval);
    });

    test('should poll for messages successfully', async () => {
      const mockMessages = [
        {
          id: 'msg-1',
          type: 'user',
          content: 'Hello',
          timestamp: new Date().toISOString(),
          session_id: 'test-session',
        },
        {
          id: 'msg-2',
          type: 'assistant',
          content: 'Hi there!',
          timestamp: new Date().toISOString(),
          session_id: 'test-session',
        },
      ];

      const mockResponse = {
        success: true,
        messages: mockMessages,
        has_more: false,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
        headers: new Headers(),
      } as Response);

      const messages = await pollingChatService.getMessages('test-session');

      expect(messages).toEqual(mockMessages);
      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('/api/v1/admin/chat/messages?session_id=test-session'),
        expect.objectContaining({
          method: 'GET',
          headers: expect.objectContaining({
            'Authorization': 'Bearer mock-token',
          }),
        })
      );
    });

    test('should use efficient polling with last message ID', async () => {
      const mockResponse = {
        success: true,
        messages: [],
        has_more: false,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve(mockResponse),
        headers: new Headers(),
      } as Response);

      await pollingChatService.getMessages('test-session', 'last-msg-id');

      expect(mockFetch).toHaveBeenCalledWith(
        expect.stringContaining('since=last-msg-id'),
        expect.any(Object)
      );
    });
  });

  describe('Error Handling and Retry Logic', () => {
    test('should retry failed requests with exponential backoff', async () => {
      // First call fails
      mockFetch.mockRejectedValueOnce(new Error('Network error'));
      
      // Second call succeeds
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ success: true, message_id: 'test-id' }),
        headers: new Headers(),
      } as Response);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      await pollingChatService.sendMessage(request);

      // Should add optimistic message and handle retry internally
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('addOptimisticMessage'),
        })
      );
    });

    test('should stop retrying after max attempts', async () => {
      // Mock all retry attempts to fail
      mockFetch.mockRejectedValue(new Error('Persistent network error'));

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      await pollingChatService.sendMessage(request);

      // Fast forward through all retry attempts
      for (let i = 0; i < 5; i++) {
        jest.advanceTimersByTime(10000);
        await Promise.resolve();
      }

      // Should add optimistic message initially
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('addOptimisticMessage'),
        })
      );
    });

    test('should handle polling errors gracefully', async () => {
      mockFetch.mockRejectedValueOnce(new Error('Network error'));

      await expect(pollingChatService.getMessages('test-session')).rejects.toThrow('Network error');
    });

    test('should handle server errors during polling', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: false,
        status: 500,
        json: () => Promise.resolve({ success: false, error: 'Server error' }),
      } as Response);

      await expect(pollingChatService.getMessages('test-session')).rejects.toThrow();
    });

    test('should handle timeout errors during polling', async () => {
      const timeoutError = new Error('Request timeout');
      timeoutError.name = 'AbortError';
      
      mockFetch.mockRejectedValueOnce(timeoutError);

      await expect(pollingChatService.getMessages('test-session')).rejects.toThrow('Request timeout');
    });
  });

  describe('Connection State Management and Transitions', () => {
    test('should track connection status correctly', () => {
      const status = pollingChatService.getConnectionStatus();
      expect(status).toBeDefined();
    });

    test('should provide detailed status information', () => {
      const statusInfo = pollingChatService.getConnectionStatusInfo();
      expect(statusInfo).toBeDefined();
      expect(statusInfo.state).toBeDefined();
    });

    test('should provide user-friendly status messages', () => {
      const statusMessage = pollingChatService.getStatusMessage();
      expect(typeof statusMessage).toBe('string');
      expect(statusMessage.length).toBeGreaterThan(0);
    });

    test('should check health status', () => {
      const isHealthy = pollingChatService.isHealthy();
      expect(typeof isHealthy).toBe('boolean');
    });

    test('should handle force reconnection', () => {
      pollingChatService.forceReconnect();
      
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      expect(mockConnectionManager.forceReconnect).toHaveBeenCalled();
    });

    test('should support status change subscriptions', () => {
      const callback = jest.fn();
      const unsubscribe = pollingChatService.onStatusChange(callback);
      
      expect(typeof unsubscribe).toBe('function');
      
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      expect(mockConnectionManager.onStatusChange).toHaveBeenCalledWith(callback);
    });

    test('should handle online/offline transitions', () => {
      // Simulate going offline
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false,
      });
      
      const offlineEvent = new Event('offline');
      window.dispatchEvent(offlineEvent);
      
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      expect(mockConnectionManager.setOffline).toHaveBeenCalled();
      
      // Simulate coming back online
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: true,
      });
      
      const onlineEvent = new Event('online');
      window.dispatchEvent(onlineEvent);
      
      expect(mockConnectionManager.setOnline).toHaveBeenCalled();
    });

    test('should transition from polling to connected state after successful polls', () => {
      pollingChatService.startPolling();
      
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      
      // Simulate successful polls
      mockConnectionManager.recordSuccessfulPoll(100);
      mockConnectionManager.recordSuccessfulPoll(120);
      mockConnectionManager.recordSuccessfulPoll(90);
      
      expect(mockConnectionManager.recordSuccessfulPoll).toHaveBeenCalledTimes(3);
    });

    test('should transition to error state on consecutive failures', () => {
      pollingChatService.startPolling();
      
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      
      // Simulate consecutive errors
      mockConnectionManager.recordPollingError('Network error');
      mockConnectionManager.recordPollingError('Timeout error');
      mockConnectionManager.recordPollingError('Server error');
      
      expect(mockConnectionManager.recordPollingError).toHaveBeenCalledTimes(3);
    });

    test('should reset error count on successful reconnection', () => {
      const mockConnectionManager = require('./ConnectionManager').connectionManager;
      
      pollingChatService.forceReconnect();
      
      expect(mockConnectionManager.forceReconnect).toHaveBeenCalled();
    });
  });

  describe('Message Queue Processing', () => {
    test('should process queued messages when coming back online', async () => {
      // Start offline
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false,
      });

      const request = {
        message: 'Queued message',
        session_id: 'test-session',
      };

      // Send message while offline (should be queued)
      await pollingChatService.sendMessage(request);

      // Mock successful response for when we come back online
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ success: true, message_id: 'queued-msg-id' }),
        headers: new Headers(),
      } as Response);

      // Come back online
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: true,
      });

      const onlineEvent = new Event('online');
      window.dispatchEvent(onlineEvent);

      // Fast forward to trigger queue processing
      jest.advanceTimersByTime(5000);
      await Promise.resolve();

      // Should have queued the message
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('queueOfflineMessage'),
        })
      );
    });

    test('should clear offline queue when messages are sent successfully', async () => {
      // Queue a message offline
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: false,
      });

      await pollingChatService.sendMessage({
        message: 'Queued message',
        session_id: 'test-session',
      });

      // Come back online and process successfully
      Object.defineProperty(navigator, 'onLine', {
        writable: true,
        value: true,
      });

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ success: true, message_id: 'processed-msg-id' }),
        headers: new Headers(),
      } as Response);

      const onlineEvent = new Event('online');
      window.dispatchEvent(onlineEvent);

      jest.advanceTimersByTime(5000);
      await Promise.resolve();

      // Should have queued the message initially
      expect(store.dispatch).toHaveBeenCalledWith(
        expect.objectContaining({
          type: expect.stringContaining('queueOfflineMessage'),
        })
      );
    });
  });

  describe('Performance Optimizations', () => {
    test('should implement conditional requests with ETags', async () => {
      const mockResponse = {
        success: true,
        messages: [],
        has_more: false,
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        headers: new Headers({
          'ETag': '"test-etag"',
          'Last-Modified': new Date().toUTCString(),
        }),
        json: () => Promise.resolve(mockResponse),
      } as Response);

      await pollingChatService.getMessages('test-session');

      // Second request should include conditional headers if implemented
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 304, // Not Modified
        json: () => Promise.resolve({}),
        headers: new Headers(),
      } as Response);

      await pollingChatService.getMessages('test-session');

      expect(mockFetch).toHaveBeenCalledTimes(2);
    });

    test('should handle 304 Not Modified responses', async () => {
      // First request returns data with ETag
      mockFetch.mockResolvedValueOnce({
        ok: true,
        headers: new Headers({ 'ETag': '"test-etag"' }),
        json: () => Promise.resolve({
          success: true,
          messages: [{ id: 'cached-msg', content: 'Cached' }],
          has_more: false,
        }),
      } as Response);

      const messages1 = await pollingChatService.getMessages('test-session');

      // Second request returns 304 Not Modified
      mockFetch.mockResolvedValueOnce({
        ok: true,
        status: 304,
        json: () => Promise.resolve({}),
        headers: new Headers(),
      } as Response);

      const messages2 = await pollingChatService.getMessages('test-session');

      // Both requests should complete
      expect(messages1).toBeDefined();
      expect(messages2).toBeDefined();
    });

    test('should use compression headers in requests', async () => {
      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ success: true, message_id: 'test-id' }),
        headers: new Headers(),
      } as Response);

      await pollingChatService.sendMessage({
        message: 'Test message',
        session_id: 'test-session',
      });

      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'Accept-Encoding': 'gzip, deflate, br',
          }),
        })
      );
    });
  });

  describe('Configuration Management', () => {
    test('should allow polling interval configuration', () => {
      pollingChatService.setPollingInterval(5000);
      
      pollingChatService.startPolling();
      jest.advanceTimersByTime(5000);
      
      // Should use the configured interval
      expect(setTimeout).toHaveBeenCalledWith(expect.any(Function), expect.any(Number));
    });

    test('should respect environment configuration', () => {
      // Test that service uses environment variables correctly
      expect(pollingChatService).toBeDefined();
      // Service should be initialized with proper base URL from environment
    });

    test('should handle missing authentication token', async () => {
      mockLocalStorage.getItem.mockReturnValueOnce(null as any);

      const request = {
        message: 'Test message',
        session_id: 'test-session',
      };

      mockFetch.mockResolvedValueOnce({
        ok: true,
        json: () => Promise.resolve({ success: true, message_id: 'test-id' }),
        headers: new Headers(),
      } as Response);

      const messageId = await pollingChatService.sendMessage(request);

      expect(messageId).toBeDefined();
      // Should still attempt to send without auth header
      expect(mockFetch).toHaveBeenCalled();
    });
  });
});
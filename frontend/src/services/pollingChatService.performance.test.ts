import { pollingChatService } from './pollingChatService';
import { store } from '../store';

// Mock dependencies for performance testing
jest.mock('../store', () => ({
  store: {
    dispatch: jest.fn(),
    getState: jest.fn(() => ({
      connection: {
        status: 'connected',
        isPollingActive: false,
        errorCount: 0,
      },
      chat: {
        messages: [],
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
    recordPollingSuccess: jest.fn(),
    recordPollingError: jest.fn(),
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
    recordPollingSuccess: jest.fn(),
    recordPollingError: jest.fn(),
  },
}));

// Mock fetch for performance testing
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
Object.defineProperty(window, 'addEventListener', {
  value: mockAddEventListener,
});

// Mock document.hidden
Object.defineProperty(document, 'hidden', {
  writable: true,
  value: false,
});

// Performance test utilities
interface PerformanceMetrics {
  averageResponseTime: number;
  minResponseTime: number;
  maxResponseTime: number;
  totalRequests: number;
  successfulRequests: number;
  failedRequests: number;
  requestsPerSecond: number;
  memoryUsage?: number;
}

class PerformanceTestRunner {
  private metrics: PerformanceMetrics = {
    averageResponseTime: 0,
    minResponseTime: Infinity,
    maxResponseTime: 0,
    totalRequests: 0,
    successfulRequests: 0,
    failedRequests: 0,
    requestsPerSecond: 0,
  };

  private responseTimes: number[] = [];
  private startTime: number = 0;

  start(): void {
    this.startTime = performance.now();
    this.metrics = {
      averageResponseTime: 0,
      minResponseTime: Infinity,
      maxResponseTime: 0,
      totalRequests: 0,
      successfulRequests: 0,
      failedRequests: 0,
      requestsPerSecond: 0,
    };
    this.responseTimes = [];
  }

  recordRequest(responseTime: number, success: boolean): void {
    this.metrics.totalRequests++;
    
    if (success) {
      this.metrics.successfulRequests++;
      this.responseTimes.push(responseTime);
      this.metrics.minResponseTime = Math.min(this.metrics.minResponseTime, responseTime);
      this.metrics.maxResponseTime = Math.max(this.metrics.maxResponseTime, responseTime);
    } else {
      this.metrics.failedRequests++;
    }
  }

  finish(): PerformanceMetrics {
    const endTime = performance.now();
    const totalTime = endTime - this.startTime;
    
    if (this.responseTimes.length > 0) {
      this.metrics.averageResponseTime = this.responseTimes.reduce((a, b) => a + b, 0) / this.responseTimes.length;
    }
    
    this.metrics.requestsPerSecond = (this.metrics.totalRequests / totalTime) * 1000;
    
    // Get memory usage if available
    if ((performance as any).memory) {
      this.metrics.memoryUsage = (performance as any).memory.usedJSHeapSize;
    }
    
    return { ...this.metrics };
  }
}

describe('PollingChatService Performance Tests', () => {
  let mockFetch: jest.MockedFunction<typeof fetch>;
  let performanceRunner: PerformanceTestRunner;

  beforeEach(() => {
    jest.clearAllMocks();
    
    mockFetch = global.fetch as jest.MockedFunction<typeof fetch>;
    performanceRunner = new PerformanceTestRunner();
  });

  afterEach(() => {
    pollingChatService.stopPolling();
  });

  describe('Message Sending Performance', () => {
    test('should handle high-frequency message sending', async () => {
      const messageCount = 100;
      const mockResponse = {
        success: true,
        message_id: 'test-message-id',
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      performanceRunner.start();

      const promises = [];
      for (let i = 0; i < messageCount; i++) {
        const startTime = performance.now();
        
        const promise = pollingChatService.sendMessage({
          message: `Performance test message ${i}`,
          session_id: 'performance-test-session',
        }).then(() => {
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, true);
        }).catch(() => {
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, false);
        });
        
        promises.push(promise);
      }

      await Promise.all(promises);
      const metrics = performanceRunner.finish();

      // Performance assertions
      expect(metrics.totalRequests).toBe(messageCount);
      expect(metrics.successfulRequests).toBe(messageCount);
      expect(metrics.failedRequests).toBe(0);
      expect(metrics.averageResponseTime).toBeLessThan(100); // Less than 100ms average
      expect(metrics.maxResponseTime).toBeLessThan(500); // Less than 500ms max
      expect(metrics.requestsPerSecond).toBeGreaterThan(10); // At least 10 RPS

      console.log('Message Sending Performance Metrics:', metrics);
    });

    test('should maintain performance with large message payloads', async () => {
      const largeMessage = 'x'.repeat(5000); // 5KB message
      const messageCount = 20;
      
      const mockResponse = {
        success: true,
        message_id: 'large-message-id',
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      performanceRunner.start();

      for (let i = 0; i < messageCount; i++) {
        const startTime = performance.now();
        
        try {
          await pollingChatService.sendMessage({
            message: largeMessage,
            session_id: 'large-message-test-session',
          });
          
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, true);
        } catch (error) {
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, false);
        }
      }

      const metrics = performanceRunner.finish();

      // Performance assertions for large messages
      expect(metrics.totalRequests).toBe(messageCount);
      expect(metrics.successfulRequests).toBe(messageCount);
      expect(metrics.averageResponseTime).toBeLessThan(200); // Less than 200ms for large messages
      expect(metrics.maxResponseTime).toBeLessThan(1000); // Less than 1s max

      console.log('Large Message Performance Metrics:', metrics);
    });
  });

  describe('Message Polling Performance', () => {
    test('should efficiently poll for messages', async () => {
      const pollCount = 50;
      const mockMessages = Array.from({ length: 10 }, (_, i) => ({
        id: `msg-${i}`,
        type: 'user',
        content: `Message ${i}`,
        timestamp: new Date().toISOString(),
        session_id: 'test-session',
      }));

      const mockResponse = {
        success: true,
        messages: mockMessages,
        has_more: false,
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      performanceRunner.start();

      for (let i = 0; i < pollCount; i++) {
        const startTime = performance.now();
        
        try {
          await pollingChatService.getMessages('test-session');
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, true);
        } catch (error) {
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, false);
        }
      }

      const metrics = performanceRunner.finish();

      // Performance assertions for polling
      expect(metrics.totalRequests).toBe(pollCount);
      expect(metrics.successfulRequests).toBe(pollCount);
      expect(metrics.averageResponseTime).toBeLessThan(50); // Less than 50ms average for polling
      expect(metrics.maxResponseTime).toBeLessThan(200); // Less than 200ms max
      expect(metrics.requestsPerSecond).toBeGreaterThan(20); // At least 20 RPS for polling

      console.log('Message Polling Performance Metrics:', metrics);
    });

    test('should handle empty polling responses efficiently', async () => {
      const pollCount = 100;
      const mockResponse = {
        success: true,
        messages: [],
        has_more: false,
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      performanceRunner.start();

      for (let i = 0; i < pollCount; i++) {
        const startTime = performance.now();
        
        try {
          await pollingChatService.getMessages('empty-session');
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, true);
        } catch (error) {
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, false);
        }
      }

      const metrics = performanceRunner.finish();

      // Empty responses should be even faster
      expect(metrics.totalRequests).toBe(pollCount);
      expect(metrics.successfulRequests).toBe(pollCount);
      expect(metrics.averageResponseTime).toBeLessThan(30); // Less than 30ms for empty responses
      expect(metrics.requestsPerSecond).toBeGreaterThan(30); // At least 30 RPS

      console.log('Empty Polling Performance Metrics:', metrics);
    });
  });

  describe('Concurrent Operations Performance', () => {
    test('should handle concurrent message sending and polling', async () => {
      const concurrentOperations = 50;
      const mockSendResponse = {
        success: true,
        message_id: 'concurrent-message-id',
      };
      const mockGetResponse = {
        success: true,
        messages: [],
        has_more: false,
      };

      mockFetch.mockImplementation((url) => {
        if (url.toString().includes('POST')) {
          return Promise.resolve({
            ok: true,
            json: () => Promise.resolve(mockSendResponse),
          } as Response);
        } else {
          return Promise.resolve({
            ok: true,
            json: () => Promise.resolve(mockGetResponse),
          } as Response);
        }
      });

      performanceRunner.start();

      const operations = [];
      
      // Mix of send and get operations
      for (let i = 0; i < concurrentOperations; i++) {
        const startTime = performance.now();
        
        let operation;
        if (i % 2 === 0) {
          // Send message
          operation = pollingChatService.sendMessage({
            message: `Concurrent message ${i}`,
            session_id: 'concurrent-test-session',
          });
        } else {
          // Get messages
          operation = pollingChatService.getMessages('concurrent-test-session');
        }
        
        operations.push(
          operation.then(() => {
            const endTime = performance.now();
            performanceRunner.recordRequest(endTime - startTime, true);
          }).catch(() => {
            const endTime = performance.now();
            performanceRunner.recordRequest(endTime - startTime, false);
          })
        );
      }

      await Promise.all(operations);
      const metrics = performanceRunner.finish();

      // Performance assertions for concurrent operations
      expect(metrics.totalRequests).toBe(concurrentOperations);
      expect(metrics.successfulRequests).toBe(concurrentOperations);
      expect(metrics.averageResponseTime).toBeLessThan(150); // Less than 150ms average
      expect(metrics.requestsPerSecond).toBeGreaterThan(5); // At least 5 RPS under concurrent load

      console.log('Concurrent Operations Performance Metrics:', metrics);
    });
  });

  describe('Memory Usage Performance', () => {
    test('should not leak memory during extended polling', async () => {
      const initialMemory = (performance as any).memory?.usedJSHeapSize || 0;
      const pollCount = 200;
      
      const mockResponse = {
        success: true,
        messages: [
          {
            id: 'memory-test-msg',
            type: 'user',
            content: 'Memory test message',
            timestamp: new Date().toISOString(),
            session_id: 'memory-test-session',
          },
        ],
        has_more: false,
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      // Perform many polling operations
      for (let i = 0; i < pollCount; i++) {
        await pollingChatService.getMessages('memory-test-session');
        
        // Force garbage collection periodically (if available)
        if (i % 50 === 0 && global.gc) {
          global.gc();
        }
      }

      const finalMemory = (performance as any).memory?.usedJSHeapSize || 0;
      const memoryIncrease = finalMemory - initialMemory;
      
      // Memory increase should be reasonable (less than 10MB)
      expect(memoryIncrease).toBeLessThan(10 * 1024 * 1024);

      console.log('Memory Usage:', {
        initial: initialMemory,
        final: finalMemory,
        increase: memoryIncrease,
        increasePerOperation: memoryIncrease / pollCount,
      });
    });

    test('should clean up resources properly', async () => {
      const initialMemory = (performance as any).memory?.usedJSHeapSize || 0;
      
      // Start and stop polling multiple times to test cleanup
      for (let i = 0; i < 10; i++) {
        pollingChatService.startPolling();
        pollingChatService.stopPolling();
      }

      // Force garbage collection if available
      if (global.gc) {
        global.gc();
      }

      const finalMemory = (performance as any).memory?.usedJSHeapSize || 0;
      const memoryIncrease = finalMemory - initialMemory;
      
      // Memory increase should be minimal after cleanup
      expect(memoryIncrease).toBeLessThan(5 * 1024 * 1024); // Less than 5MB

      console.log('Resource Cleanup Memory Usage:', {
        initial: initialMemory,
        final: finalMemory,
        increase: memoryIncrease,
      });
    });
  });

  describe('Network Efficiency Performance', () => {
    test('should minimize redundant requests with caching', async () => {
      const sessionId = 'cache-test-session';
      const mockMessages = [
        {
          id: 'cached-msg-1',
          type: 'user',
          content: 'Cached message',
          timestamp: new Date().toISOString(),
          session_id: sessionId,
        },
      ];

      const mockResponse = {
        success: true,
        messages: mockMessages,
        has_more: false,
      };

      // First request returns data with ETag
      mockFetch.mockResolvedValueOnce({
        ok: true,
        headers: new Headers({
          'ETag': '"test-etag"',
          'Last-Modified': new Date().toUTCString(),
        }),
        json: () => Promise.resolve(mockResponse),
      } as Response);

      // Subsequent requests return 304 Not Modified
      mockFetch.mockResolvedValue({
        ok: true,
        status: 304,
        json: () => Promise.resolve({}),
      } as Response);

      performanceRunner.start();

      // First request (should fetch data)
      const startTime1 = performance.now();
      await pollingChatService.getMessages(sessionId);
      const endTime1 = performance.now();
      performanceRunner.recordRequest(endTime1 - startTime1, true);

      // Subsequent requests (should use cache/304 responses)
      for (let i = 0; i < 10; i++) {
        const startTime = performance.now();
        await pollingChatService.getMessages(sessionId);
        const endTime = performance.now();
        performanceRunner.recordRequest(endTime - startTime, true);
      }

      const metrics = performanceRunner.finish();

      // Cached requests should be faster
      expect(metrics.totalRequests).toBe(11);
      expect(metrics.successfulRequests).toBe(11);
      expect(metrics.averageResponseTime).toBeLessThan(50); // Should be fast due to caching

      // Verify that conditional requests were made
      expect(mockFetch).toHaveBeenCalledWith(
        expect.any(String),
        expect.objectContaining({
          headers: expect.objectContaining({
            'If-None-Match': '"test-etag"',
          }),
        })
      );

      console.log('Network Caching Performance Metrics:', metrics);
    });

    test('should handle network latency gracefully', async () => {
      const latencyMs = 100; // Simulate 100ms network latency
      const requestCount = 20;

      const mockResponse = {
        success: true,
        messages: [],
        has_more: false,
      };

      mockFetch.mockImplementation(() => 
        new Promise(resolve => {
          setTimeout(() => {
            resolve({
              ok: true,
              json: () => Promise.resolve(mockResponse),
            } as Response);
          }, latencyMs);
        })
      );

      performanceRunner.start();

      for (let i = 0; i < requestCount; i++) {
        const startTime = performance.now();
        
        try {
          await pollingChatService.getMessages('latency-test-session');
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, true);
        } catch (error) {
          const endTime = performance.now();
          performanceRunner.recordRequest(endTime - startTime, false);
        }
      }

      const metrics = performanceRunner.finish();

      // Should handle latency gracefully
      expect(metrics.totalRequests).toBe(requestCount);
      expect(metrics.successfulRequests).toBe(requestCount);
      expect(metrics.averageResponseTime).toBeGreaterThan(latencyMs); // Should account for latency
      expect(metrics.averageResponseTime).toBeLessThan(latencyMs + 50); // But not much overhead

      console.log('Network Latency Performance Metrics:', metrics);
    });
  });

  describe('Polling Interval Optimization', () => {
    test('should optimize polling intervals based on activity', async () => {
      const mockResponse = {
        success: true,
        messages: [],
        has_more: false,
      };

      mockFetch.mockResolvedValue({
        ok: true,
        json: () => Promise.resolve(mockResponse),
      } as Response);

      // Start polling
      pollingChatService.startPolling();

      // Measure initial polling frequency (should be base interval)
      const initialRequestCount = mockFetch.mock.calls.length;
      
      // Wait for a few polling cycles
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      const afterWaitRequestCount = mockFetch.mock.calls.length;
      const requestsInPeriod = afterWaitRequestCount - initialRequestCount;

      // Should have made some requests but not too many (smart interval)
      expect(requestsInPeriod).toBeGreaterThan(0);
      expect(requestsInPeriod).toBeLessThan(10); // Not more than 10 requests per second

      console.log('Polling Interval Optimization:', {
        requestsInPeriod,
        estimatedInterval: 1000 / requestsInPeriod,
      });
    });
  });
});

// Benchmark utility for comparing different implementations
export class PollingPerformanceBenchmark {
  static async compareImplementations(
    implementations: { name: string; service: typeof pollingChatService }[],
    testConfig: {
      messageCount: number;
      pollCount: number;
      sessionId: string;
    }
  ): Promise<{ [key: string]: PerformanceMetrics }> {
    const results: { [key: string]: PerformanceMetrics } = {};

    for (const impl of implementations) {
      const runner = new PerformanceTestRunner();
      runner.start();

      // Test message sending
      for (let i = 0; i < testConfig.messageCount; i++) {
        const startTime = performance.now();
        try {
          await impl.service.sendMessage({
            message: `Benchmark message ${i}`,
            session_id: testConfig.sessionId,
          });
          const endTime = performance.now();
          runner.recordRequest(endTime - startTime, true);
        } catch (error) {
          const endTime = performance.now();
          runner.recordRequest(endTime - startTime, false);
        }
      }

      // Test message polling
      for (let i = 0; i < testConfig.pollCount; i++) {
        const startTime = performance.now();
        try {
          await impl.service.getMessages(testConfig.sessionId);
          const endTime = performance.now();
          runner.recordRequest(endTime - startTime, true);
        } catch (error) {
          const endTime = performance.now();
          runner.recordRequest(endTime - startTime, false);
        }
      }

      results[impl.name] = runner.finish();
    }

    return results;
  }
}
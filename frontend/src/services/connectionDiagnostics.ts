import { websocketService } from './websocketService';

export interface HealthStatus {
  status: 'healthy' | 'degraded' | 'unhealthy';
  responseTime: number;
  timestamp: Date;
  details: string[];
}

export interface EndpointStatus {
  available: boolean;
  responseTime: number;
  error?: string;
  timestamp: Date;
}

export interface NetworkStatus {
  online: boolean;
  connectionType?: string;
  effectiveType?: string;
  downlink?: number;
  rtt?: number;
  timestamp: Date;
}

export interface ConfigStatus {
  valid: boolean;
  issues: string[];
  configuration: {
    wsUrl: string;
    apiUrl: string;
    protocol: string;
    host: string;
    token: boolean;
  };
}

export interface ConnectionAttempt {
  timestamp: Date;
  url: string;
  result: 'success' | 'failed' | 'timeout';
  errorCode?: number;
  errorMessage?: string;
  duration: number;
}

export interface DiagnosticReport {
  timestamp: Date;
  connectionAttempts: ConnectionAttempt[];
  healthStatus: HealthStatus;
  endpointStatus: EndpointStatus;
  networkStatus: NetworkStatus;
  configStatus: ConfigStatus;
  recommendations: string[];
}

class ConnectionDiagnostics {
  private connectionHistory: ConnectionAttempt[] = [];
  private readonly maxHistorySize = 50;

  /**
   * Check backend health status
   */
  public async checkBackendHealth(): Promise<HealthStatus> {
    const startTime = Date.now();
    const details: string[] = [];
    
    try {
      const response = await fetch('/api/v1/health', {
        method: 'GET',
        headers: {
          'Content-Type': 'application/json',
        },
        signal: AbortSignal.timeout(5000), // 5 second timeout
      });

      const responseTime = Date.now() - startTime;
      
      if (response.ok) {
        const data = await response.json();
        details.push(`Backend responded in ${responseTime}ms`);
        details.push(`Status: ${data.status || 'unknown'}`);
        
        return {
          status: responseTime < 1000 ? 'healthy' : 'degraded',
          responseTime,
          timestamp: new Date(),
          details,
        };
      } else {
        details.push(`HTTP ${response.status}: ${response.statusText}`);
        return {
          status: 'unhealthy',
          responseTime,
          timestamp: new Date(),
          details,
        };
      }
    } catch (error) {
      const responseTime = Date.now() - startTime;
      const errorMessage = error instanceof Error ? error.message : 'Unknown error';
      details.push(`Connection failed: ${errorMessage}`);
      
      return {
        status: 'unhealthy',
        responseTime,
        timestamp: new Date(),
        details,
      };
    }
  }

  /**
   * Validate WebSocket endpoint availability
   */
  public async validateWebSocketEndpoint(): Promise<EndpointStatus> {
    const startTime = Date.now();
    
    return new Promise((resolve) => {
      try {
        const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
        const wsUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;
        
        const testWs = new WebSocket(wsUrl);
        
        const timeout = setTimeout(() => {
          testWs.close();
          resolve({
            available: false,
            responseTime: Date.now() - startTime,
            error: 'Connection timeout',
            timestamp: new Date(),
          });
        }, 5000);

        testWs.onopen = () => {
          clearTimeout(timeout);
          testWs.close();
          resolve({
            available: true,
            responseTime: Date.now() - startTime,
            timestamp: new Date(),
          });
        };

        testWs.onerror = (error) => {
          clearTimeout(timeout);
          resolve({
            available: false,
            responseTime: Date.now() - startTime,
            error: 'WebSocket connection error',
            timestamp: new Date(),
          });
        };

        testWs.onclose = (event) => {
          clearTimeout(timeout);
          if (event.code !== 1000) {
            resolve({
              available: false,
              responseTime: Date.now() - startTime,
              error: `WebSocket closed with code ${event.code}: ${event.reason}`,
              timestamp: new Date(),
            });
          }
        };
      } catch (error) {
        resolve({
          available: false,
          responseTime: Date.now() - startTime,
          error: error instanceof Error ? error.message : 'Unknown error',
          timestamp: new Date(),
        });
      }
    });
  }

  /**
   * Test network connectivity
   */
  public async testNetworkConnectivity(): Promise<NetworkStatus> {
    const networkStatus: NetworkStatus = {
      online: navigator.onLine,
      timestamp: new Date(),
    };

    // Get connection information if available
    if ('connection' in navigator) {
      const connection = (navigator as any).connection;
      if (connection) {
        networkStatus.connectionType = connection.type;
        networkStatus.effectiveType = connection.effectiveType;
        networkStatus.downlink = connection.downlink;
        networkStatus.rtt = connection.rtt;
      }
    }

    return networkStatus;
  }

  /**
   * Validate frontend WebSocket configuration
   */
  public validateConfiguration(): ConfigStatus {
    const issues: string[] = [];
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;
    const apiUrl = `/api/v1`;
    const token = !!localStorage.getItem('adminToken');

    // Check token
    if (!token) {
      issues.push('No admin token found in localStorage');
    }

    // Check protocol consistency
    if (window.location.protocol === 'https:' && protocol !== 'wss:') {
      issues.push('HTTPS page should use WSS protocol');
    }

    // Check URL format
    try {
      new URL(wsUrl.replace('ws:', 'http:').replace('wss:', 'https:'));
    } catch (error) {
      issues.push('Invalid WebSocket URL format');
    }

    // Check if running on localhost
    if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
      if (!window.location.port) {
        issues.push('Localhost connection without port may cause issues');
      }
    }

    return {
      valid: issues.length === 0,
      issues,
      configuration: {
        wsUrl,
        apiUrl,
        protocol,
        host: window.location.host,
        token,
      },
    };
  }

  /**
   * Record a connection attempt
   */
  public recordConnectionAttempt(attempt: ConnectionAttempt): void {
    this.connectionHistory.unshift(attempt);
    
    // Limit history size
    if (this.connectionHistory.length > this.maxHistorySize) {
      this.connectionHistory = this.connectionHistory.slice(0, this.maxHistorySize);
    }
  }

  /**
   * Get connection attempt history
   */
  public getConnectionHistory(): ConnectionAttempt[] {
    return [...this.connectionHistory];
  }

  /**
   * Generate comprehensive diagnostic report
   */
  public async generateDiagnosticReport(): Promise<DiagnosticReport> {
    const [healthStatus, endpointStatus, networkStatus] = await Promise.all([
      this.checkBackendHealth(),
      this.validateWebSocketEndpoint(),
      this.testNetworkConnectivity(),
    ]);

    const configStatus = this.validateConfiguration();
    const recommendations = this.generateRecommendations(
      healthStatus,
      endpointStatus,
      networkStatus,
      configStatus
    );

    return {
      timestamp: new Date(),
      connectionAttempts: this.getConnectionHistory(),
      healthStatus,
      endpointStatus,
      networkStatus,
      configStatus,
      recommendations,
    };
  }

  /**
   * Generate recommendations based on diagnostic results
   */
  private generateRecommendations(
    health: HealthStatus,
    endpoint: EndpointStatus,
    network: NetworkStatus,
    config: ConfigStatus
  ): string[] {
    const recommendations: string[] = [];

    // Health-based recommendations
    if (health.status === 'unhealthy') {
      recommendations.push('Backend server appears to be down or unreachable');
      recommendations.push('Check if the backend service is running');
      recommendations.push('Verify network connectivity to the server');
    } else if (health.status === 'degraded') {
      recommendations.push('Backend server is responding slowly');
      recommendations.push('Check server load and performance');
    }

    // Endpoint-based recommendations
    if (!endpoint.available) {
      recommendations.push('WebSocket endpoint is not available');
      if (endpoint.error?.includes('timeout')) {
        recommendations.push('Connection timeout - check firewall settings');
      }
      if (endpoint.error?.includes('1006')) {
        recommendations.push('WebSocket closed abnormally - check server logs');
      }
    }

    // Network-based recommendations
    if (!network.online) {
      recommendations.push('Device appears to be offline');
      recommendations.push('Check internet connection');
    } else if (network.effectiveType === 'slow-2g' || network.effectiveType === '2g') {
      recommendations.push('Slow network connection detected');
      recommendations.push('WebSocket may have difficulty maintaining connection');
    }

    // Configuration-based recommendations
    config.issues.forEach(issue => {
      recommendations.push(`Configuration issue: ${issue}`);
    });

    if (config.issues.includes('No admin token found')) {
      recommendations.push('Please log in again to refresh authentication');
    }

    // General recommendations
    if (recommendations.length === 0) {
      recommendations.push('All diagnostics passed - connection should work normally');
    } else {
      recommendations.push('Try refreshing the page after addressing the issues above');
      recommendations.push('If problems persist, contact system administrator');
    }

    return recommendations;
  }

  /**
   * Test WebSocket connection with detailed logging
   */
  public async testWebSocketConnection(): Promise<ConnectionAttempt> {
    const startTime = Date.now();
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/v1/admin/chat/ws`;

    return new Promise((resolve) => {
      try {
        const testWs = new WebSocket(wsUrl);
        
        const timeout = setTimeout(() => {
          testWs.close();
          const attempt: ConnectionAttempt = {
            timestamp: new Date(),
            url: wsUrl,
            result: 'timeout',
            duration: Date.now() - startTime,
            errorMessage: 'Connection timeout after 10 seconds',
          };
          this.recordConnectionAttempt(attempt);
          resolve(attempt);
        }, 10000);

        testWs.onopen = () => {
          clearTimeout(timeout);
          testWs.close();
          const attempt: ConnectionAttempt = {
            timestamp: new Date(),
            url: wsUrl,
            result: 'success',
            duration: Date.now() - startTime,
          };
          this.recordConnectionAttempt(attempt);
          resolve(attempt);
        };

        testWs.onerror = () => {
          clearTimeout(timeout);
          const attempt: ConnectionAttempt = {
            timestamp: new Date(),
            url: wsUrl,
            result: 'failed',
            duration: Date.now() - startTime,
            errorMessage: 'WebSocket connection error',
          };
          this.recordConnectionAttempt(attempt);
          resolve(attempt);
        };

        testWs.onclose = (event) => {
          clearTimeout(timeout);
          if (event.code !== 1000) {
            const attempt: ConnectionAttempt = {
              timestamp: new Date(),
              url: wsUrl,
              result: 'failed',
              duration: Date.now() - startTime,
              errorCode: event.code,
              errorMessage: `WebSocket closed with code ${event.code}: ${event.reason || 'No reason provided'}`,
            };
            this.recordConnectionAttempt(attempt);
            resolve(attempt);
          }
        };
      } catch (error) {
        const attempt: ConnectionAttempt = {
          timestamp: new Date(),
          url: wsUrl,
          result: 'failed',
          duration: Date.now() - startTime,
          errorMessage: error instanceof Error ? error.message : 'Unknown error',
        };
        this.recordConnectionAttempt(attempt);
        resolve(attempt);
      }
    });
  }

  /**
   * Clear connection history
   */
  public clearHistory(): void {
    this.connectionHistory = [];
  }
}

// Create singleton instance
export const connectionDiagnostics = new ConnectionDiagnostics();

export default connectionDiagnostics;
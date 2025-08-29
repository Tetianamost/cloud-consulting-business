/**
 * Connection Diagnostics Service
 * Provides comprehensive diagnostics for WebSocket connectivity issues
 */

export interface DiagnosticResult {
  test: string;
  status: 'pass' | 'fail' | 'warning';
  message: string;
  details?: any;
  timestamp: Date;
}

export interface DiagnosticReport {
  timestamp: Date;
  overallStatus: 'healthy' | 'degraded' | 'unhealthy';
  results: DiagnosticResult[];
  recommendations: string[];
}

export class ConnectionDiagnosticsService {
  private apiUrl: string;
  private wsUrl: string;

  constructor() {
    this.apiUrl = process.env.REACT_APP_API_URL || 'http://localhost:8061';
    this.wsUrl = process.env.REACT_APP_WS_URL || 'ws://localhost:8061/api/v1/admin/chat/ws';
  }

  /**
   * Run comprehensive connection diagnostics
   */
  async runDiagnostics(): Promise<DiagnosticReport> {
    const results: DiagnosticResult[] = [];
    const recommendations: string[] = [];

    console.log('🔍 Starting WebSocket connection diagnostics...');

    // Test 1: Check if backend server is reachable
    results.push(await this.testBackendHealth());

    // Test 2: Check authentication token
    results.push(await this.testAuthenticationToken());

    // Test 3: Check WebSocket endpoint configuration
    results.push(await this.testWebSocketConfiguration());

    // Test 4: Test basic HTTP connectivity
    results.push(await this.testHttpConnectivity());

    // Test 5: Test CORS configuration
    results.push(await this.testCorsConfiguration());

    // Test 6: Test WebSocket endpoint availability
    results.push(await this.testWebSocketEndpoint());

    // Analyze results and generate recommendations
    const failedTests = results.filter(r => r.status === 'fail');
    const warningTests = results.filter(r => r.status === 'warning');

    if (failedTests.length === 0 && warningTests.length === 0) {
      recommendations.push('All diagnostics passed. WebSocket should be working correctly.');
    } else {
      if (failedTests.some(t => t.test === 'Backend Health Check')) {
        recommendations.push('🚨 CRITICAL: Backend server is not running or not accessible. Start the backend server on port 8061.');
        recommendations.push('Run: cd backend && go run cmd/server/main.go');
      }
      
      if (failedTests.some(t => t.test === 'Authentication Token')) {
        recommendations.push('🔐 Authentication required: Please log in to get a valid admin token.');
      }
      
      if (failedTests.some(t => t.test === 'WebSocket Endpoint')) {
        const wsTest = failedTests.find(t => t.test === 'WebSocket Endpoint');
        if (wsTest?.details?.diagnosis === 'immediate_close_issue') {
          recommendations.push('🔌 WebSocket connects but closes immediately (code 1005). This is likely caused by React StrictMode or client-side connection management issues.');
          recommendations.push('Try: 1) Disable React StrictMode in development, 2) Check for multiple connection attempts, 3) Review WebSocket service cleanup logic.');
        } else {
          recommendations.push('🔌 WebSocket endpoint is not accessible. Check if the backend WebSocket handler is properly configured.');
        }
      }
      
      if (failedTests.some(t => t.test === 'CORS Configuration')) {
        recommendations.push('🌐 CORS issue detected. Ensure backend CORS allows origin: http://localhost:3007');
      }
    }

    const overallStatus = failedTests.length > 0 ? 'unhealthy' : 
                         warningTests.length > 0 ? 'degraded' : 'healthy';

    return {
      timestamp: new Date(),
      overallStatus,
      results,
      recommendations
    };
  }

  /**
   * Test if backend server is running and healthy
   */
  private async testBackendHealth(): Promise<DiagnosticResult> {
    try {
      console.log('🏥 Testing backend health...');
      const response = await fetch(`${this.apiUrl}/health`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });

      if (response.ok) {
        const data = await response.json();
        return {
          test: 'Backend Health Check',
          status: 'pass',
          message: `Backend server is running (${data.status})`,
          details: data,
          timestamp: new Date()
        };
      } else {
        return {
          test: 'Backend Health Check',
          status: 'fail',
          message: `Backend server returned ${response.status}: ${response.statusText}`,
          details: { status: response.status, statusText: response.statusText },
          timestamp: new Date()
        };
      }
    } catch (error) {
      return {
        test: 'Backend Health Check',
        status: 'fail',
        message: `Cannot reach backend server: ${error instanceof Error ? error.message : 'Unknown error'}`,
        details: { error: error instanceof Error ? error.message : error },
        timestamp: new Date()
      };
    }
  }

  /**
   * Test authentication token availability
   */
  private async testAuthenticationToken(): Promise<DiagnosticResult> {
    try {
      console.log('🔐 Testing authentication token...');
      const token = localStorage.getItem('adminToken');
      
      if (!token) {
        return {
          test: 'Authentication Token',
          status: 'fail',
          message: 'No admin token found in localStorage',
          details: { tokenExists: false },
          timestamp: new Date()
        };
      }

      // Try to decode JWT to check if it's valid format
      try {
        const payload = JSON.parse(atob(token.split('.')[1]));
        const now = Math.floor(Date.now() / 1000);
        
        if (payload.exp && payload.exp < now) {
          return {
            test: 'Authentication Token',
            status: 'fail',
            message: 'Admin token has expired',
            details: { tokenExists: true, expired: true, exp: payload.exp, now },
            timestamp: new Date()
          };
        }

        return {
          test: 'Authentication Token',
          status: 'pass',
          message: 'Valid admin token found',
          details: { tokenExists: true, expired: false, exp: payload.exp },
          timestamp: new Date()
        };
      } catch (decodeError) {
        return {
          test: 'Authentication Token',
          status: 'warning',
          message: 'Admin token exists but format is invalid',
          details: { tokenExists: true, validFormat: false },
          timestamp: new Date()
        };
      }
    } catch (error) {
      return {
        test: 'Authentication Token',
        status: 'fail',
        message: `Error checking authentication token: ${error instanceof Error ? error.message : 'Unknown error'}`,
        details: { error: error instanceof Error ? error.message : error },
        timestamp: new Date()
      };
    }
  }

  /**
   * Test WebSocket configuration
   */
  private async testWebSocketConfiguration(): Promise<DiagnosticResult> {
    try {
      console.log('⚙️ Testing WebSocket configuration...');
      const issues: string[] = [];
      
      // Check environment variables
      const apiUrl = process.env.REACT_APP_API_URL;
      const wsUrl = process.env.REACT_APP_WS_URL;
      
      if (!apiUrl) {
        issues.push('REACT_APP_API_URL environment variable not set');
      }
      
      if (!wsUrl) {
        issues.push('REACT_APP_WS_URL environment variable not set');
      }
      
      // Check URL formats
      try {
        new URL(this.apiUrl);
      } catch {
        issues.push(`Invalid API URL format: ${this.apiUrl}`);
      }
      
      // Check WebSocket URL format
      if (!this.wsUrl.startsWith('ws://') && !this.wsUrl.startsWith('wss://')) {
        issues.push(`Invalid WebSocket URL format: ${this.wsUrl}`);
      }

      if (issues.length === 0) {
        return {
          test: 'WebSocket Configuration',
          status: 'pass',
          message: 'WebSocket configuration is valid',
          details: { apiUrl: this.apiUrl, wsUrl: this.wsUrl },
          timestamp: new Date()
        };
      } else {
        return {
          test: 'WebSocket Configuration',
          status: 'fail',
          message: `Configuration issues found: ${issues.join(', ')}`,
          details: { issues, apiUrl: this.apiUrl, wsUrl: this.wsUrl },
          timestamp: new Date()
        };
      }
    } catch (error) {
      return {
        test: 'WebSocket Configuration',
        status: 'fail',
        message: `Error checking WebSocket configuration: ${error instanceof Error ? error.message : 'Unknown error'}`,
        details: { error: error instanceof Error ? error.message : error },
        timestamp: new Date()
      };
    }
  }

  /**
   * Test basic HTTP connectivity to backend
   */
  private async testHttpConnectivity(): Promise<DiagnosticResult> {
    try {
      console.log('🌐 Testing HTTP connectivity...');
      const startTime = Date.now();
      
      const response = await fetch(`${this.apiUrl}/api/v1/config/services`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
        },
      });

      const endTime = Date.now();
      const responseTime = endTime - startTime;

      if (response.ok) {
        return {
          test: 'HTTP Connectivity',
          status: 'pass',
          message: `HTTP connectivity working (${responseTime}ms)`,
          details: { responseTime, status: response.status },
          timestamp: new Date()
        };
      } else {
        return {
          test: 'HTTP Connectivity',
          status: 'fail',
          message: `HTTP request failed: ${response.status} ${response.statusText}`,
          details: { status: response.status, statusText: response.statusText, responseTime },
          timestamp: new Date()
        };
      }
    } catch (error) {
      return {
        test: 'HTTP Connectivity',
        status: 'fail',
        message: `HTTP connectivity failed: ${error instanceof Error ? error.message : 'Unknown error'}`,
        details: { error: error instanceof Error ? error.message : error },
        timestamp: new Date()
      };
    }
  }

  /**
   * Test CORS configuration
   */
  private async testCorsConfiguration(): Promise<DiagnosticResult> {
    try {
      console.log('🔒 Testing CORS configuration...');
      
      const response = await fetch(`${this.apiUrl}/health`, {
        method: 'GET',
        headers: {
          'Accept': 'application/json',
          'Origin': 'http://localhost:3007', // Explicitly set origin
        },
      });

      const corsHeaders = {
        'access-control-allow-origin': response.headers.get('Access-Control-Allow-Origin'),
        'access-control-allow-methods': response.headers.get('Access-Control-Allow-Methods'),
        'access-control-allow-headers': response.headers.get('Access-Control-Allow-Headers'),
        'access-control-allow-credentials': response.headers.get('Access-Control-Allow-Credentials'),
      };

      const allowedOrigin = corsHeaders['access-control-allow-origin'];
      
      // Check if CORS is properly configured
      if (allowedOrigin === 'http://localhost:3007' || allowedOrigin === '*') {
        return {
          test: 'CORS Configuration',
          status: 'pass',
          message: 'CORS is properly configured for frontend origin',
          details: { corsHeaders, allowedOrigin },
          timestamp: new Date()
        };
      } else {
        return {
          test: 'CORS Configuration',
          status: 'fail',
          message: `CORS not configured for frontend origin. Allowed: ${allowedOrigin}`,
          details: { corsHeaders, expectedOrigin: 'http://localhost:3007', allowedOrigin },
          timestamp: new Date()
        };
      }
    } catch (error) {
      return {
        test: 'CORS Configuration',
        status: 'warning',
        message: `Could not test CORS: ${error instanceof Error ? error.message : 'Unknown error'}`,
        details: { error: error instanceof Error ? error.message : error },
        timestamp: new Date()
      };
    }
  }

  /**
   * Test WebSocket endpoint availability
   */
  private async testWebSocketEndpoint(): Promise<DiagnosticResult> {
    return new Promise((resolve) => {
      console.log('🔌 Testing WebSocket endpoint...');
      
      const token = localStorage.getItem('adminToken');
      if (!token) {
        resolve({
          test: 'WebSocket Endpoint',
          status: 'fail',
          message: 'Cannot test WebSocket endpoint without authentication token',
          details: { reason: 'no_token' },
          timestamp: new Date()
        });
        return;
      }

      const wsUrl = `${this.wsUrl}?token=${encodeURIComponent(token)}`;
      const ws = new WebSocket(wsUrl);
      
      const timeout = setTimeout(() => {
        ws.close();
        resolve({
          test: 'WebSocket Endpoint',
          status: 'fail',
          message: 'WebSocket connection timeout (10 seconds)',
          details: { reason: 'timeout', url: wsUrl.replace(/token=[^&]+/, 'token=***') },
          timestamp: new Date()
        });
      }, 10000);

      ws.onopen = () => {
        console.log('🔌 WebSocket opened, waiting to check stability...');
        // Don't close immediately - wait to see if connection stays stable
        setTimeout(() => {
          if (ws.readyState === WebSocket.OPEN) {
            clearTimeout(timeout);
            ws.close();
            resolve({
              test: 'WebSocket Endpoint',
              status: 'pass',
              message: 'WebSocket endpoint is accessible and maintains stable connections',
              details: { url: wsUrl.replace(/token=[^&]+/, 'token=***') },
              timestamp: new Date()
            });
          }
        }, 2000); // Wait 2 seconds to check if connection stays open
      };

      ws.onerror = (error) => {
        clearTimeout(timeout);
        resolve({
          test: 'WebSocket Endpoint',
          status: 'fail',
          message: 'WebSocket connection failed',
          details: { error, url: wsUrl.replace(/token=[^&]+/, 'token=***') },
          timestamp: new Date()
        });
      };

      ws.onclose = (event) => {
        clearTimeout(timeout);
        if (event.code !== 1000) { // 1000 is normal closure
          let message = `WebSocket closed unexpectedly: ${event.code}`;
          if (event.code === 1005) {
            message += ' (no status) - Connection closes immediately after opening. This indicates a client-side connection management issue.';
          }
          if (event.reason) {
            message += ` ${event.reason}`;
          }
          
          resolve({
            test: 'WebSocket Endpoint',
            status: 'fail',
            message,
            details: { 
              code: event.code, 
              reason: event.reason || 'No reason provided', 
              url: wsUrl.replace(/token=[^&]+/, 'token=***'),
              diagnosis: event.code === 1005 ? 'immediate_close_issue' : 'unexpected_close'
            },
            timestamp: new Date()
          });
        }
      };
    });
  }

  /**
   * Print diagnostic report to console
   */
  printReport(report: DiagnosticReport): void {
    console.log('\n🔍 WebSocket Connection Diagnostic Report');
    console.log('==========================================');
    console.log(`Timestamp: ${report.timestamp.toISOString()}`);
    console.log(`Overall Status: ${report.overallStatus.toUpperCase()}`);
    console.log('\nTest Results:');
    
    report.results.forEach((result, index) => {
      const icon = result.status === 'pass' ? '✅' : result.status === 'warning' ? '⚠️' : '❌';
      console.log(`${index + 1}. ${icon} ${result.test}: ${result.message}`);
      if (result.details) {
        console.log(`   Details:`, result.details);
      }
    });

    if (report.recommendations.length > 0) {
      console.log('\nRecommendations:');
      report.recommendations.forEach((rec, index) => {
        console.log(`${index + 1}. ${rec}`);
      });
    }
    
    console.log('\n==========================================\n');
  }
}

// Export singleton instance
export const connectionDiagnostics = new ConnectionDiagnosticsService();
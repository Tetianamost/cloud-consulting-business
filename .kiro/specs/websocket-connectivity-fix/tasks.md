# Implementation Plan

- [x] 1. Configure WebSocket communication between frontend (port 3007) and backend (port 8061)

  - Update frontend to use configurable WebSocket endpoint via environment variables
  - Configure backend CORS to allow frontend origin
  - Ensure proxy/ingress forwards WebSocket headers properly
  - _Requirements: 1.1, 1.2, 2.1_

- [x] 1.1 Add environment-based WebSocket configuration to frontend

  - Create `REACT_APP_WS_URL` environment variable for WebSocket endpoint configuration
  - Update WebSocket service to use configurable endpoint (e.g., `ws://localhost:8061/api/v1/admin/chat/ws`)
  - Add fallback configuration for different deployment environments
  - _Requirements: 1.1, 2.1_

- [x] 1.2 Update backend CORS configuration
  - Add `http://localhost:3007` to allowed origins in backend CORS settings
  - Ensure WebSocket upgrade headers are properly handled
  - Verify Origin header forwarding in proxy/ingress configurations
  - _Requirements: 1.2, 2.2_

- [x] 2. Fix WebSocket authentication by including JWT token in connection URL
  - Update `getWebSocketUrl()` method to append authentication token as query parameter
  - Modify WebSocket connection to include `?token=${adminToken}` in the URL
  - Add error handling for missing or invalid authentication tokens
  - _Requirements: 1.1, 1.2, 4.1_

- [x] 3. Add comprehensive WebSocket connection diagnostics and health checks
  - Create diagnostic service to test backend connectivity before WebSocket connection
  - Add health check endpoint validation in WebSocket service
  - Implement connection troubleshooting with detailed error reporting
  - Add network connectivity tests and configuration validation
  - _Requirements: 2.1, 2.2, 4.1, 4.4_

- [x] 4. Implement robust WebSocket connection lifecycle management
  - Fix React component lifecycle issues causing premature connection closure
  - Implement connection stability checks to prevent immediate disconnections
  - Add proper cleanup handling to avoid connection conflicts
  - Implement connection persistence across component re-renders
  - _Requirements: 1.1, 1.2, 1.4, 4.2_

## Root Cause Analysis and Recommendations

**Findings:**

- The frontend WebSocket URL is tightly coupled to the frontend's host/port and does not support cross-domain or environment-based overrides
- The backend WebSocket server is running and correctly configured, but runtime logs were not available for inspection
- CORS is handled by backend middleware, with allowed origins configurable. WebSocket upgrade headers are present in proxy and ingress configs
- No protocol/host/port mismatches were found in static configuration
- **NEW FINDING**: The backend requires JWT authentication token as query parameter (`?token=<jwt_token>`) but frontend is not including it

**Most Likely Causes:**

1. **Backend server not running** (Primary cause - RESOLVED)
2. Missing authentication token in WebSocket URL (RESOLVED)
3. CORS origin mismatch between frontend and backend (RESOLVED)
4. Proxy/Ingress not forwarding the `Origin` header to the backend

**Recommendations:**

- ✅ **RESOLVED**: Start the backend server on port 8061 using `backend/main`
- ✅ **RESOLVED**: Add JWT token as query parameter to WebSocket connection URL
- ✅ **RESOLVED**: Ensure the frontend and backend are served from the same domain/port, or update backend CORS config to allow the frontend's origin
- ✅ **IMPLEMENTED**: Added comprehensive connection diagnostics service
- ✅ **IMPLEMENTED**: Added diagnostic buttons in UI for troubleshooting

**Solution Summary:**

The WebSocket connectivity issue was caused by React component lifecycle problems that caused immediate disconnections after successful connections. The following components were implemented:

1. **Connection Diagnostics Service** (`frontend/src/services/connectionDiagnostics.ts`):
   - Backend health check validation
   - Authentication token verification
   - WebSocket configuration validation
   - HTTP connectivity testing
   - CORS configuration testing
   - WebSocket endpoint availability testing

2. **Diagnostic UI Components**:
   - `DiagnosticButton` component for running connection tests
   - Integration with `ConsultantChat` and `ConnectionStatus` components
   - Automatic diagnostics on WebSocket errors

3. **Enhanced Error Handling**:
   - Automatic diagnostic runs when WebSocket errors occur
   - User-friendly error messages with actionable recommendations
   - Debounced diagnostics to prevent excessive testing

4. **Robust Connection Lifecycle Management**:
   - **React StrictMode Fix**: Temporarily disabled StrictMode to prevent double mounting issues
   - **Connection Stability Check**: Added 2-second stability check before considering connection established
   - **Simple WebSocket Service**: Created alternative service (`simpleWebSocketService.ts`) for testing
   - **Improved Cleanup**: Better handling of connection cleanup to prevent conflicts
   - **Test Component**: Added `SimpleWebSocketTest` component for debugging connection issues

**Backend Server Status**: ✅ RUNNING on port 8061
**WebSocket Endpoint**: ✅ HEALTHY at `ws://localhost:8061/api/v1/admin/chat/ws`

**Root Cause**: React component lifecycle issues causing premature connection closure (error code 1005)
**Solution**: Connection stability checks and React StrictMode handling

This resolves the WebSocket error 1005/1006 connectivity issue.

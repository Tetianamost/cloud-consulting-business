# Implementation Plan

- [x] 1. Create backend REST API endpoints for polling-based chat

  - Implement POST /api/v1/admin/chat/messages endpoint for sending messages
  - Implement GET /api/v1/admin/chat/messages endpoint for retrieving messages
  - Add proper authentication and authorization middleware
  - Include message validation and error handling
  - _Requirements: 1.2, 2.1, 3.2_

- [x] 1.1 Implement message sending endpoint

  - Create SendMessageRequest and SendMessageResponse structs
  - Add message validation (content length, session ID format)
  - Store messages in database with unique IDs and timestamps
  - Return immediate confirmation with message ID
  - _Requirements: 1.2, 4.1_

- [x] 1.2 Implement message retrieval endpoint

  - Create GetMessagesResponse struct with pagination support
  - Add query parameters for session_id and since (last message ID/timestamp)
  - Implement efficient database queries to fetch new messages
  - Include message ordering and limit controls
  - _Requirements: 1.4, 2.2_

- [x] 2. Create frontend polling chat service

  - Implement PollingChatService class to replace WebSocket service
  - Add HTTP client methods for sending and receiving messages
  - Create polling loop with configurable intervals
  - Implement connection state management
  - _Requirements: 1.1, 2.1, 4.2_

- [x] 2.1 Implement message sending functionality

  - Create sendMessage method with HTTP POST to backend
  - Add optimistic message updates (show immediately as "sending")
  - Implement retry logic with exponential backoff for failed sends
  - Update message status based on server response
  - _Requirements: 1.2, 4.1, 4.3_

- [x] 2.2 Implement message polling functionality

  - Create polling loop that fetches new messages every 3-5 seconds
  - Add smart polling intervals (faster when active, slower when idle)
  - Implement efficient polling using last message ID/timestamp
  - Handle polling errors with appropriate backoff strategies
  - _Requirements: 1.4, 2.2, 3.1_

- [ ] 3. Add comprehensive error handling and retry logic

  - Implement network error detection and handling
  - Add exponential backoff for failed requests
  - Create offline/online state detection
  - Add user-friendly error messages and status indicators
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 3.1 Implement retry mechanisms

  - Add retry logic for failed message sends (1s, 2s, 4s, 8s intervals)
  - Implement polling retry with backoff when server is unavailable
  - Create message queue for offline scenarios
  - Add duplicate message prevention using client-side tracking
  - _Requirements: 3.1, 3.3, 4.3_

- [x] 3.2 Add connection status management

  - Create ConnectionManager class to track polling state
  - Implement connection status indicators (connected, polling, error, offline)
  - Add automatic reconnection when network is restored
  - Show appropriate status messages to users
  - _Requirements: 3.2, 3.4_

- [x] 4. Integrate polling chat service with existing UI components

  - Update ChatPage component to use polling service instead of WebSocket
  - Modify ConsultantChat component to work with HTTP-based messaging
  - Update ConnectionStatus component to show polling-specific status
  - Ensure seamless integration with existing Redux store
  - _Requirements: 5.1, 5.2_

- [x] 4.1 Update chat components

  - Modify ChatPage to initialize polling service instead of WebSocket
  - Update message sending logic to use HTTP POST requests
  - Change message receiving to use polling instead of WebSocket events
  - Maintain existing UI behavior and appearance
  - _Requirements: 5.1, 5.3_

- [x] 4.2 Update Redux integration

  - Modify chat slice actions to work with polling service
  - Update connection slice to handle polling-specific states
  - Ensure message state management works with HTTP-based flow
  - Preserve existing chat history and session management
  - _Requirements: 5.2, 5.3_

- [x] 5. Add performance optimizations and smart polling

  - Implement adaptive polling intervals based on user activity
  - Add message caching to reduce redundant server requests
  - Optimize polling efficiency with conditional requests
  - Add performance monitoring and metrics
  - _Requirements: 2.2, 2.4_

- [x] 5.1 Implement smart polling intervals

  - Detect user activity (typing, sending messages, page focus)
  - Use faster polling (2s) when user is actively chatting
  - Switch to slower polling (10s) when user is idle
  - Stop polling when page is not visible or user is away
  - _Requirements: 2.2, 2.4_

- [x] 5.2 Add message caching and optimization

  - Implement client-side message caching to avoid duplicate requests
  - Add conditional requests using ETags or timestamps
  - Optimize database queries on backend for message retrieval
  - Add response compression for large message payloads
  - _Requirements: 2.2, 2.3_

- [x] 6. Create feature flag system for WebSocket/Polling toggle

  - Add configuration option to switch between WebSocket and polling
  - Implement automatic fallback from WebSocket to polling on connection failures
  - Create admin interface to toggle chat implementation
  - Add monitoring to compare performance of both approaches
  - _Requirements: 1.1, 3.1_

- [x] 6.1 Implement fallback mechanism

  - Detect WebSocket connection failures automatically
  - Switch to polling mode when WebSocket fails repeatedly
  - Add user notification when fallback occurs
  - Allow manual switching between modes via admin interface
  - _Requirements: 3.1, 3.4_

- [ ] 7. Add comprehensive testing for polling chat system

  - Create unit tests for polling service functionality
  - Add integration tests for backend API endpoints
  - Implement end-to-end tests for complete chat flow
  - Add performance tests for polling efficiency
  - _Requirements: 1.1, 2.1, 3.1_

- [ ] 7.1 Create polling service tests

  - Test message sending with various success/failure scenarios
  - Test polling loop behavior with different intervals
  - Test error handling and retry logic
  - Test connection state management and transitions
  - _Requirements: 1.2, 2.1, 3.1_

- [x] 7.2 Create backend API tests

  - Test message sending endpoint with validation and error cases
  - Test message retrieval endpoint with pagination and filtering
  - Test concurrent access and race conditions
  - Test rate limiting and performance under load
  - _Requirements: 1.2, 1.4, 2.1_

- [ ] 8. Remove WebSocket implementation and cleanup codebase

  - Remove WebSocket service and related frontend code
  - Remove WebSocket handlers and backend infrastructure
  - Clean up WebSocket-related configuration and dependencies
  - Update documentation to reflect polling-based approach
  - _Requirements: Repository cleanup and maintenance_

- [ ] 8.1 Remove frontend WebSocket code

  - Delete `frontend/src/services/websocketService.ts`
  - Delete `frontend/src/services/simpleWebSocketService.ts`
  - Delete `frontend/src/services/connectionDiagnostics.ts`
  - Delete `frontend/src/components/admin/SimpleWebSocketTest.tsx`
  - Delete `frontend/src/components/admin/DiagnosticButton.tsx`
  - Remove WebSocket test pages (`frontend/public/test-websocket.html`, `frontend/src/test-websocket.html`)
  - _Requirements: Frontend cleanup_

- [ ] 8.2 Remove backend WebSocket code

  - Remove WebSocket handler from `backend/internal/handlers/chat_handler.go`
  - Remove WebSocket route registration from `backend/internal/server/server.go`
  - Remove WebSocket-related imports and dependencies
  - Clean up WebSocket authentication middleware if no longer needed
  - _Requirements: Backend cleanup_

- [ ] 8.3 Update Redux store and components

  - Remove WebSocket-specific actions from `frontend/src/store/slices/connectionSlice.ts`
  - Clean up connection state management to only handle polling states
  - Remove WebSocket imports from React components
  - Update `IntegratedAdminDashboard.tsx` to remove WebSocket initialization
  - _Requirements: State management cleanup_

- [ ] 8.4 Remove WebSocket configuration and infrastructure

  - Remove WebSocket-related environment variables from `.env.example`
  - Clean up nginx WebSocket configuration (`nginx/nginx-websocket-lb.conf`)
  - Remove WebSocket-related Docker and Kubernetes configurations
  - Update monitoring configurations to remove WebSocket metrics
  - _Requirements: Infrastructure cleanup_

- [ ] 8.5 Clean up WebSocket tests and documentation

  - Remove WebSocket integration tests (`backend/test_websocket_integration.go`)
  - Remove WebSocket-related test files and specs
  - Delete the `websocket-connectivity-fix` spec directory
  - Update API documentation to reflect polling endpoints only
  - Update troubleshooting guides to remove WebSocket sections
  - _Requirements: Documentation and testing cleanup_

- [ ] 8.6 Remove WebSocket dependencies and imports
  - Remove WebSocket-related npm packages from `package.json` if any
  - Remove Go WebSocket dependencies from `go.mod` if no longer needed
  - Clean up unused imports across all files
  - Run linting to ensure no dead code remains
  - _Requirements: Dependency cleanup_

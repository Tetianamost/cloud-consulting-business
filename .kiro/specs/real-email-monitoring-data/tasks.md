# Implementation Plan

- [x] 1. Create database schema and migration for email events tracking

  - Create email_events table with proper indexes for performance
  - Add migration script to backend/scripts/ directory
  - Include rollback migration for safe deployment
  - _Requirements: 4.1, 4.2, 4.3_

- [x] 2. Implement email event data models and interfaces

  - Define EmailEvent, EmailEventType, and EmailEventStatus types in domain models
  - Create EmailEventRecorder interface in interfaces directory
  - Create EmailMetricsService interface with comprehensive methods
  - _Requirements: 2.1, 2.2, 2.3, 3.1_

- [x] 3. Create email event repository with database operations

  - Implement EmailEventRepository interface with CRUD operations
  - Add Create, Update, GetByInquiryID, GetByMessageID methods
  - Implement GetMetrics method for aggregated statistics calculation
  - Add proper error handling and logging for database operations
  - _Requirements: 4.2, 4.3, 3.2_

- [x] 4. Implement email event recorder service

  - Create emailEventRecorder service implementing EmailEventRecorder interface
  - Add RecordEmailSent method for capturing email send events
  - Add UpdateEmailStatus method for delivery status updates
  - Implement non-blocking event recording to prevent email delivery failures
  - _Requirements: 2.1, 2.2, 2.3, 7.1_

- [x] 5. Create email metrics service for real-time calculations

  - Implement EmailMetricsService interface with metrics calculation logic
  - Add GetEmailMetrics method with time range filtering
  - Add GetEmailStatusByInquiry method for individual inquiry status
  - Implement efficient database queries with proper aggregation
  - _Requirements: 3.1, 3.2, 3.3_

- [x] 6. Enhance existing email service with event recording integration

  - Modify SendReportEmail method to record consultant notification events
  - Modify SendCustomerConfirmation method to record customer confirmation events
  - Modify SendInquiryNotification method to record inquiry notification events
  - Ensure email delivery continues even if event recording fails
  - _Requirements: 2.1, 2.2, 2.3, 7.1_

- [x] 7. Update SES service to capture message IDs and delivery status

  - Modify SES service to return message IDs from AWS SES responses
  - Add delivery status tracking capabilities to SES service
  - Implement SES webhook handling for delivery confirmations (optional)
  - Add error categorization for bounces and spam complaints
  - _Requirements: 5.1, 5.2, 5.3, 5.4_

- [x] 8. Enhance admin handler with real email metrics endpoints

  - Update GetSystemMetrics method to use real email metrics instead of mock calculations
  - Update GetEmailStatus method to return actual email event data
  - Add new endpoint for detailed email event history if needed
  - Remove mock data generation and replace with real data or proper error responses
  - _Requirements: 3.1, 3.2, 3.3, 6.4_

- [x] 9. Update frontend V0DataAdapter to prioritize real data over mock data

  - Modify safeAdaptEmailMetrics method to handle real API responses properly
  - Remove generateMockEmailMetrics fallback when real data is available
  - Update error handling to show appropriate messages instead of mock data
  - Ensure V0EmailDeliveryDashboardConnected uses real data when available
  - _Requirements: 6.1, 6.2, 6.3, 7.2, 7.3_

- [x] 10. Add comprehensive error handling and fallback mechanisms

  - Implement proper error responses in API endpoints when data is unavailable
  - Add logging for email event recording failures without blocking email delivery
  - Update frontend to display meaningful error messages instead of mock data warnings
  - Add health checks for email event recording system
  - _Requirements: 7.1, 7.2, 7.3, 7.4_

- [x] 11. Create database migration and update server initialization

  - Add email events table migration to backend/scripts/
  - Update server.go to initialize email event recorder and metrics services
  - Ensure proper dependency injection for new services
  - Add configuration options for email event tracking if needed
  - _Requirements: 4.1, 4.2_

- [x] 12. Write comprehensive tests for email event tracking system

  - Create unit tests for email event repository operations
  - Create unit tests for email metrics service calculations
  - Create integration tests for email service with event recording
  - Create tests for enhanced admin handler endpoints with real data
  - Add tests for frontend data adapter with real API responses
  - _Requirements: All requirements - testing coverage_

- [x] 13. Update email service factory to include event recording

  - Modify NewEmailServiceWithSES to include email event recorder dependency
  - Update email service constructor to accept event recorder parameter
  - Ensure proper initialization of event recording in production and test environments
  - Add configuration validation for email event tracking
  - _Requirements: 2.1, 2.2, 2.3_

- [x] 14. Add monitoring and observability for email event system

  - Add structured logging for email event recording operations
  - Add metrics collection for email event recording success/failure rates
  - Add health check endpoint for email event tracking system
  - Implement alerting for high email event recording failure rates
  - _Requirements: 7.1, 7.4_

- [x] 15. Performance optimization and indexing
  - Verify database indexes are properly created for email events queries
  - Optimize email metrics calculation queries for large datasets
  - Add caching for frequently accessed email metrics if needed
  - Implement email event retention policies to manage database growth
  - _Requirements: 4.3, 3.3_

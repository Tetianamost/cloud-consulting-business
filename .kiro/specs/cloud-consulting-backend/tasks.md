# Implementation Plan - Minimal Working Backend

This implementation plan focuses on creating a minimal, working Go backend that can accept and store inquiries. Advanced features like validation, testing, hooks, and external integrations will be added in future iterations.

- [x] 1. Set up basic project structure
  - Create Go module with minimal directory structure (cmd, internal)
  - Set up basic configuration management with environment variables
  - Create simple .env file for local development
  - _Requirements: 5.3, 7.3_

- [x] 2. Create minimal HTTP server
  - Set up Gin web framework with basic middleware (CORS, logging, recovery)
  - Implement basic server startup and graceful shutdown
  - Add simple configuration loading without validation
  - _Requirements: 5.1, 5.2_

- [x] 3. Add health check endpoint
  - Implement GET /health endpoint that returns server status
  - Include basic server information (service name, version, timestamp)
  - Test endpoint responds correctly
  - _Requirements: 4.3, 6.1_

- [x] 4. Create basic data models
  - Define simple Inquiry struct with essential fields (name, email, company, phone, services, message)
  - Create CreateInquiryRequest struct for API input
  - Add basic service type constants (assessment, migration, optimization, architecture_review)
  - _Requirements: 1.1, 1.4_

- [x] 5. Implement in-memory storage
  - Create simple in-memory storage for inquiries using Go maps
  - Implement basic CRUD operations (Create, Read by ID, List all)
  - Add thread-safe access with mutex
  - _Requirements: 1.1, 4.2_

- [x] 6. Build inquiry endpoints
  - Implement POST /api/v1/inquiries endpoint with basic validation
  - Add GET /api/v1/inquiries/{id} endpoint to retrieve single inquiry
  - Create GET /api/v1/inquiries endpoint to list all inquiries
  - Include basic error handling and JSON responses
  - _Requirements: 1.1, 1.3, 5.3_

- [x] 7. Add service configuration endpoint
  - Implement GET /api/v1/config/services endpoint
  - Return available service types with descriptions
  - Format response for frontend consumption
  - _Requirements: 1.2_

- [x] 8. Create basic documentation
  - Add README.md with setup and usage instructions
  - Document available API endpoints with examples
  - Include curl commands for testing endpoints
  - _Requirements: 8.1, 8.2_

- [x] 9. Test the complete flow
  - Start the server locally
  - Test health check endpoint
  - Create sample inquiry via POST endpoint
  - Retrieve inquiry via GET endpoint
  - Verify all endpoints work correctly
  - _Requirements: 5.1, 5.2_

- [x] 10. Add Amazon Bedrock configuration
  - Add Bedrock configuration struct with environment variable support
  - Include AWS_BEARER_TOKEN_BEDROCK, region, model ID, and timeout settings
  - Update config loading to include Bedrock settings with defaults
  - _Requirements: 2.6_

- [x] 11. Implement Bedrock service interface
  - Create BedrockService interface with GenerateText method
  - Implement HTTP client for Bedrock API calls with proper authentication
  - Add request/response structs for Bedrock API communication
  - Include timeout handling and basic error handling
  - _Requirements: 2.2, 2.3_

- [x] 12. Create report generator component
  - Implement ReportGenerator that uses BedrockService
  - Build structured prompts based on inquiry service type and content
  - Handle Bedrock responses and format as plain text reports
  - Add graceful error handling when Bedrock calls fail
  - _Requirements: 2.1, 2.4, 2.5_

- [x] 13. Update data models for reports
  - Add Report struct with ID, InquiryID, content, status, and timestamps
  - Modify Inquiry struct to include optional Reports field
  - Update in-memory storage to handle inquiry-report relationships
  - _Requirements: 2.4_

- [x] 14. Integrate report generation into inquiry creation
  - Modify inquiry creation endpoint to trigger report generation
  - Call report generator after successfully storing inquiry
  - Store generated report with inquiry in memory storage
  - Ensure inquiry creation succeeds even if report generation fails
  - _Requirements: 2.1, 2.5_

- [x] 15. Add report retrieval endpoint
  - Implement GET /api/v1/inquiries/{id}/report endpoint
  - Return generated report content for specific inquiry
  - Handle cases where no report exists yet
  - _Requirements: 2.4_

- [ ] 16. Update API documentation
  - Document new Bedrock integration in README.md
  - Add environment variable setup instructions for AWS_BEARER_TOKEN_BEDROCK
  - Include examples of inquiry creation with report generation
  - Document new report retrieval endpoint
  - _Requirements: 8.1, 8.2_

- [ ] 17. Test Bedrock integration end-to-end
  - Set up Bedrock API key in environment
  - Create test inquiry and verify report generation
  - Test error scenarios (invalid API key, network issues)
  - Verify inquiry creation works when Bedrock fails
  - Test report retrieval endpoint
  - _Requirements: 2.1, 2.2, 2.3, 2.4, 2.5_

- [ ] 18. Add AWS SES configuration
  - Add SES configuration struct with environment variable support
  - Include AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, region, and sender email settings
  - Update config loading to include SES settings with validation
  - _Requirements: 3.3, 3.6_

- [ ] 19. Implement SES service interface
  - Create SESService interface with SendEmail method
  - Implement AWS SDK v2 client for SES API calls with proper authentication
  - Add request/response structs for SES API communication
  - Include timeout handling and basic error handling
  - _Requirements: 3.3_

- [ ] 20. Create email service component
  - Implement EmailService that uses SESService
  - Build HTML and text email templates for report notifications
  - Format inquiry and report data into professional email content
  - Add graceful error handling when SES calls fail
  - _Requirements: 3.1, 3.2, 3.5_

- [ ] 21. Integrate email notifications into inquiry flow
  - Modify inquiry creation endpoint to trigger email notifications
  - Send email after successfully generating report
  - Ensure inquiry creation succeeds even if email sending fails
  - Log email sending attempts and failures
  - _Requirements: 3.1, 3.2, 3.4_

- [ ] 22. Update environment configuration
  - Add SES environment variables to .env.example
  - Document required AWS SES setup and verification steps
  - Include sender email verification requirements
  - _Requirements: 3.6_

- [ ] 23. Test email integration end-to-end
  - Set up AWS SES credentials and verify sender email
  - Create test inquiry and verify email delivery
  - Test error scenarios (invalid credentials, unverified sender)
  - Verify inquiry creation works when email fails
  - Test email content formatting and delivery
  - _Requirements: 3.1, 3.2, 3.3, 3.4, 3.5_

## Future Iterations (Not in this implementation)

The following features will be added in subsequent iterations:
- Database integration (PostgreSQL)
- Advanced validation and error handling
- AI report generation with LLM integration
- Email and Slack notifications
- Agent hooks system
- Comprehensive testing suite
- External service integrations (AWS S3, SQS)
- Authentication and authorization
- Rate limiting and security features
- Monitoring and metrics
- Docker containerization
- Production deployment setup
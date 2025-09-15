# Requirements Document

## Introduction

This specification addresses the critical issue where report downloads are failing with 404 errors due to a URL path mismatch between the frontend and backend. Users are unable to download generated reports in both PDF and HTML formats, which significantly impacts the usability of the AI consultant reporting system.

## Requirements

### Requirement 1

**User Story:** As an admin user, I want to download generated reports in PDF format, so that I can save and share consultant reports offline.

#### Acceptance Criteria

1. WHEN I click the "Download PDF" button on a report THEN the system SHALL initiate a PDF download without errors
2. WHEN the PDF download completes THEN the system SHALL provide a properly named PDF file with the report content
3. IF no report exists for the inquiry THEN the system SHALL display an appropriate error message
4. WHEN the download fails due to server errors THEN the system SHALL display a user-friendly error message

### Requirement 2

**User Story:** As an admin user, I want to download generated reports in HTML format, so that I can view and customize report content in web browsers.

#### Acceptance Criteria

1. WHEN I click the "Download HTML" button on a report THEN the system SHALL initiate an HTML download without errors
2. WHEN the HTML download completes THEN the system SHALL provide a properly named HTML file with formatted report content
3. WHEN the HTML file is opened THEN it SHALL display the report with proper styling and formatting
4. IF the HTML generation fails THEN the system SHALL return an appropriate error response

### Requirement 3

**User Story:** As a system administrator, I want consistent API endpoints between frontend and backend, so that all report download functionality works reliably.

#### Acceptance Criteria

1. WHEN the frontend makes a download request THEN the backend SHALL have a matching route to handle the request
2. WHEN API endpoints are defined THEN they SHALL follow consistent naming conventions across the application
3. WHEN route parameters are used THEN they SHALL be consistently named between frontend and backend
4. IF route mismatches exist THEN the system SHALL be updated to use consistent endpoint patterns

### Requirement 4

**User Story:** As a developer, I want proper error handling for download failures, so that users receive clear feedback when downloads fail.

#### Acceptance Criteria

1. WHEN a download request fails with 404 THEN the system SHALL log the specific endpoint being called
2. WHEN authentication fails during download THEN the system SHALL redirect to login or show auth error
3. WHEN server errors occur during download THEN the system SHALL display technical details to admin users
4. WHEN network errors occur THEN the system SHALL provide retry options to users

### Requirement 5

**User Story:** As an admin user, I want download functionality to work from all report interfaces, so that I can access reports consistently throughout the application.

#### Acceptance Criteria

1. WHEN I access reports from the AI Reports page THEN download buttons SHALL work correctly
2. WHEN I access reports from the Inquiry List THEN download options SHALL function properly  
3. WHEN I view reports in modal dialogs THEN download actions SHALL complete successfully
4. WHEN I use dropdown menus for downloads THEN all format options SHALL be available and functional
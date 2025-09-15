# Design Document

## Overview

This design addresses the critical 404 error occurring during report downloads by fixing the URL path mismatch between frontend and backend systems. The solution ensures consistent API endpoints, proper error handling, and reliable download functionality across all report interfaces.

## Architecture

### Current Issue Analysis

**Problem:** Frontend requests `/api/v1/admin/reports/{inquiryId}/download/{format}` but backend serves `/api/v1/admin/inquiries/{inquiryId}/download/{format}`

**Root Cause:** Inconsistent naming convention between "reports" and "inquiries" in the API endpoint paths

**Impact:** All report downloads fail with 404 errors, breaking core functionality

### Solution Approach

**Option 1: Update Backend Routes (Recommended)**
- Change backend route from `/inquiries/{id}/download/{format}` to `/reports/{inquiryId}/download/{format}`
- Maintains frontend consistency and follows RESTful resource naming
- Aligns with existing `/reports` endpoints in the admin API

**Option 2: Update Frontend URLs**
- Change frontend to use `/inquiries/{id}/download/{format}`
- Less preferred as it breaks semantic consistency with other report endpoints

## Components and Interfaces

### Backend Route Updates

#### Current Route Structure
```go
// Current (incorrect)
admin.GET("/inquiries/:inquiryId/download/:format", s.adminHandler.DownloadReport)

// Should be
admin.GET("/reports/:inquiryId/download/:format", s.adminHandler.DownloadReport)
```

#### Handler Interface
```go
type AdminHandler interface {
    DownloadReport(c *gin.Context) // Existing - no changes needed
}
```

### Frontend API Service

#### Current Implementation (Working)
```typescript
async downloadReport(inquiryId: string, format: 'pdf' | 'html'): Promise<Blob> {
    const url = `${this.baseUrl}/api/v1/admin/reports/${inquiryId}/download/${format}`;
    // ... rest of implementation
}
```

### Error Handling Enhancement

#### Backend Error Responses
```go
type DownloadErrorResponse struct {
    Success bool   `json:"success"`
    Error   string `json:"error"`
    Code    string `json:"code"`
}
```

#### Frontend Error Handling
```typescript
interface DownloadError {
    message: string;
    code?: string;
    retry?: boolean;
}
```

## Data Models

### Report Download Request
```go
type DownloadRequest struct {
    InquiryID string `uri:"inquiryId" binding:"required"`
    Format    string `uri:"format" binding:"required,oneof=pdf html"`
}
```

### Download Response Headers
```go
type DownloadHeaders struct {
    ContentType        string // "application/pdf" or "text/html"
    ContentDisposition string // "attachment; filename=\"report.pdf\""
    ContentLength      string // File size in bytes
}
```

## Error Handling

### Backend Error Scenarios

1. **Invalid Format**
   - Status: 400 Bad Request
   - Response: `{"success": false, "error": "Invalid format", "code": "INVALID_FORMAT"}`

2. **Inquiry Not Found**
   - Status: 404 Not Found  
   - Response: `{"success": false, "error": "Inquiry not found", "code": "INQUIRY_NOT_FOUND"}`

3. **No Reports Available**
   - Status: 404 Not Found
   - Response: `{"success": false, "error": "No reports found", "code": "NO_REPORTS"}`

4. **Report Generation Failed**
   - Status: 500 Internal Server Error
   - Response: `{"success": false, "error": "Report generation failed", "code": "GENERATION_ERROR"}`

### Frontend Error Handling

1. **Network Errors**
   - Display: "Network error. Please check your connection and try again."
   - Action: Provide retry button

2. **Authentication Errors (401)**
   - Display: "Session expired. Please log in again."
   - Action: Redirect to login or refresh token

3. **Not Found Errors (404)**
   - Display: "Report not found or not yet generated."
   - Action: Suggest generating report first

4. **Server Errors (500)**
   - Display: "Server error occurred. Please try again later."
   - Action: Log error details for debugging

## Testing Strategy

### Unit Tests

#### Backend Route Testing
```go
func TestDownloadReportRoute(t *testing.T) {
    // Test correct route registration
    // Test parameter extraction
    // Test format validation
}
```

#### Frontend API Testing
```typescript
describe('downloadReport', () => {
    it('should call correct endpoint URL', () => {
        // Verify URL construction
        // Test different formats
        // Test error handling
    });
});
```

### Integration Tests

#### End-to-End Download Flow
1. Generate a report for test inquiry
2. Attempt PDF download via frontend
3. Verify file download and content
4. Attempt HTML download via frontend
5. Verify file download and formatting

#### Error Scenario Testing
1. Test download with non-existent inquiry ID
2. Test download with invalid format parameter
3. Test download without authentication
4. Test download with server errors

### Manual Testing Checklist

1. **AI Reports Page**
   - [ ] PDF download button works
   - [ ] HTML download button works
   - [ ] Error messages display correctly

2. **Inquiry List Page**
   - [ ] Dropdown PDF download works
   - [ ] Dropdown HTML download works
   - [ ] Loading states display properly

3. **Report Preview Modal**
   - [ ] Download buttons function correctly
   - [ ] Modal remains open during download
   - [ ] Success feedback provided

## Implementation Plan

### Phase 1: Backend Route Fix
1. Update route registration in `server.go`
2. Verify handler parameter extraction still works
3. Test route with curl/Postman

### Phase 2: Error Handling Enhancement
1. Add structured error responses
2. Improve error logging with context
3. Add error code constants

### Phase 3: Frontend Error Handling
1. Enhance error message display
2. Add retry functionality for failed downloads
3. Improve loading states during downloads

### Phase 4: Testing and Validation
1. Run integration tests
2. Perform manual testing across all interfaces
3. Validate error scenarios
4. Performance testing for large reports

## Security Considerations

### Authentication
- Maintain existing JWT token validation
- Ensure download endpoints require admin authentication
- Validate user permissions for inquiry access

### Data Protection
- Sanitize filenames to prevent path traversal
- Validate inquiry ownership/access rights
- Log download activities for audit trail

### Rate Limiting
- Consider implementing download rate limiting
- Prevent abuse of report generation resources
- Monitor download patterns for anomalies

## Performance Considerations

### Caching Strategy
- Cache generated PDF/HTML content temporarily
- Implement cache invalidation when reports update
- Consider CDN for static report assets

### Large File Handling
- Stream large PDF files instead of loading in memory
- Implement progress indicators for large downloads
- Add timeout handling for slow generations

### Concurrent Downloads
- Handle multiple simultaneous download requests
- Implement queue system for resource-intensive generations
- Monitor server resources during peak usage
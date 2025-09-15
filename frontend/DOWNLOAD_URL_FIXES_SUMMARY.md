# Download URL Fixes Summary

## Issue
The frontend was making download requests to incorrect URLs, causing 404 errors:
- **Incorrect URL**: `/api/reports/{inquiryId}/download?format={format}`
- **Correct URL**: `/api/v1/admin/reports/{inquiryId}/download/{format}`

## Changes Made

### 1. Updated reportService.ts
**File**: `frontend/src/services/reportService.ts`

**Before**:
```typescript
const response = await fetch(`/api/reports/${inquiryId}/download?format=${format}`, {
  method: 'GET',
  credentials: 'include',
});
```

**After**:
```typescript
return await apiService.downloadReport(inquiryId, format);
```

**Impact**: Now uses the centralized API service which has the correct URL and proper error handling.

### 2. Verified API Service (Already Correct)
**File**: `frontend/src/services/api.ts`

The API service already had the correct URL:
```typescript
const url = `${this.baseUrl}/api/v1/admin/reports/${inquiryId}/download/${format}`;
```

## Components Using Download Functionality

All the following components are now using the correct download URLs:

### ✅ Components Using apiService.downloadReport() (Correct)
1. **inquiry-list.tsx** - Uses `apiService.downloadReport(inquiryId, format)`
2. **V0InquiryList.tsx** - Uses `apiService.downloadReport(inquiryId, format)`

### ✅ Components Using reportService.downloadReport() (Fixed)
1. **AIReportsPage.tsx** - Uses `downloadReport(report.inquiry_id, format)` from reportService
2. **report-preview-modal.tsx** - Uses `onDownload` prop which comes from AIReportsPage

### ✅ Components with Simulated Downloads (No Changes Needed)
1. **metrics-dashboard.tsx** - Simulates downloads, doesn't make real API calls
2. **V0InquiryAnalysisSection.tsx** - Simulates downloads with fake blobs
3. **ChatPage.tsx** - Downloads chat exports, not reports

## URL Structure Comparison

| Component | Old URL Pattern | New URL Pattern | Status |
|-----------|----------------|-----------------|---------|
| reportService | `/api/reports/{id}/download?format={format}` | Uses apiService (correct) | ✅ Fixed |
| apiService | N/A | `/api/v1/admin/reports/{id}/download/{format}` | ✅ Already Correct |
| inquiry-list | N/A | Uses apiService (correct) | ✅ Already Correct |
| V0InquiryList | N/A | Uses apiService (correct) | ✅ Already Correct |
| AIReportsPage | N/A | Uses reportService (now fixed) | ✅ Fixed |

## Backend Route (Already Fixed in Task 1)
The backend route is correctly registered as:
```go
admin.GET("/reports/:inquiryId/download/:format", s.adminHandler.DownloadReport)
```

This creates the full path: `/api/v1/admin/reports/{inquiryId}/download/{format}`

## Testing
After these changes, all download functionality should work correctly:
- PDF downloads from all report interfaces
- HTML downloads from all report interfaces
- Proper error handling with structured error responses
- Contextual logging for debugging

## Error Handling Improvements
The API service includes proper error handling:
- Automatic token refresh on 401 errors
- Structured error messages from backend
- Proper HTTP status code handling
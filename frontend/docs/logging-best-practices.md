# Frontend Logging Best Practices

## Overview

This document outlines the logging best practices for the Cloud Consulting Platform frontend, focusing on production-ready logging strategies and debugging approaches.

## Logging Principles

### 1. Environment-Aware Logging
- **Development**: Verbose logging for debugging
- **Production**: Error-only logging for security and performance
- **Testing**: Minimal logging to avoid test noise

### 2. Log Levels
- **Error**: Critical issues that need immediate attention
- **Warn**: Potential issues that should be monitored
- **Info**: General information about application flow
- **Debug**: Detailed information for debugging (development only)

### 3. Security Considerations
- Never log sensitive information (passwords, tokens, personal data)
- Avoid exposing internal system details in production
- Use generic error messages for user-facing errors

## Current Implementation

### Authentication Context
```typescript
// Good: Clean error logging
const login = async (username: string, password: string): Promise<boolean> => {
  try {
    const response = await apiService.login(username, password);
    
    if (response.success && response.token) {
      localStorage.setItem('adminToken', response.token);
      apiService.setAuthToken(response.token);
      setIsAuthenticated(true);
      return true;
    }
    return false;
  } catch (error) {
    console.error('Login failed:', error);
    return false;
  }
};
```

### What We Avoid
```typescript
// Avoid: Verbose debug logging in production
console.log('AuthContext: Attempting login...');
console.log('AuthContext: Login response:', response);
console.log('AuthContext: Login successful, setting token and auth state');
```

## Recommended Patterns

### 1. Error Logging
```typescript
// Good: Structured error logging
try {
  await riskyOperation();
} catch (error) {
  console.error('Operation failed:', {
    operation: 'riskyOperation',
    error: error.message,
    timestamp: new Date().toISOString()
  });
}
```

### 2. Conditional Debug Logging
```typescript
// Good: Environment-based debug logging
const DEBUG = process.env.NODE_ENV === 'development';

if (DEBUG) {
  console.log('Debug info:', debugData);
}
```

### 3. User-Friendly Error Messages
```typescript
// Good: Generic user messages, detailed logging
try {
  await apiCall();
} catch (error) {
  console.error('API call failed:', error);
  setUserError('Something went wrong. Please try again.');
}
```

## Tools and Libraries

### Development Tools
- **React DevTools**: Component debugging
- **Redux DevTools**: State debugging
- **Browser DevTools**: Network and performance debugging

### Future Considerations
- **Sentry**: Error tracking and monitoring
- **LogRocket**: Session replay and debugging
- **Winston**: Structured logging library

## Implementation Guidelines

### Do's
- ✅ Log errors with context
- ✅ Use structured logging when possible
- ✅ Include timestamps for debugging
- ✅ Log user actions for analytics
- ✅ Use appropriate log levels

### Don'ts
- ❌ Log sensitive information
- ❌ Use verbose logging in production
- ❌ Log every function call
- ❌ Expose internal errors to users
- ❌ Leave debug logs in production code

## Code Review Checklist

### Before Merging
- [ ] No debug console.log statements in production code
- [ ] Error logging includes sufficient context
- [ ] No sensitive information in logs
- [ ] User-facing errors are generic and helpful
- [ ] Development-only logging is properly conditional

### Testing
- [ ] Console output is clean in production build
- [ ] Error scenarios are properly logged
- [ ] No information leakage through logs
- [ ] Performance impact of logging is minimal

## Migration Strategy

### Phase 1: Cleanup (Current)
- Remove debug logging from production code
- Standardize error logging patterns
- Implement basic security measures

### Phase 2: Enhancement
- Implement structured logging
- Add environment-based logging configuration
- Integrate error tracking service

### Phase 3: Advanced
- Add performance monitoring
- Implement user analytics
- Create logging dashboard

## Examples

### Authentication Logging
```typescript
// Current implementation
const login = async (username: string, password: string): Promise<boolean> => {
  try {
    const response = await apiService.login(username, password);
    
    if (response.success && response.token) {
      localStorage.setItem('adminToken', response.token);
      apiService.setAuthToken(response.token);
      setIsAuthenticated(true);
      return true;
    }
    return false;
  } catch (error) {
    console.error('Login failed:', error);
    return false;
  }
};
```

### Chat System Logging
```typescript
// Good: Error logging with context
const sendMessage = async (message: string) => {
  try {
    await chatService.sendMessage(message);
  } catch (error) {
    console.error('Failed to send message:', {
      error: error.message,
      messageLength: message.length,
      sessionId: currentSession?.id
    });
    throw new Error('Failed to send message. Please try again.');
  }
};
```

### API Service Logging
```typescript
// Good: Request/response logging for debugging
private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
  try {
    const response = await fetch(`${this.baseUrl}${endpoint}`, options);
    
    if (!response.ok) {
      console.error('API request failed:', {
        endpoint,
        status: response.status,
        statusText: response.statusText
      });
      throw new Error(`API request failed: ${response.statusText}`);
    }
    
    return await response.json();
  } catch (error) {
    console.error('Network error:', { endpoint, error: error.message });
    throw error;
  }
}
```

This logging strategy ensures clean, secure, and maintainable code while providing adequate debugging information for development and error tracking in production.
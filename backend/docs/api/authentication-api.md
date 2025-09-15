# Authentication API Documentation

## Overview

The Cloud Consulting Platform uses JWT (JSON Web Token) based authentication for admin access. This document describes the authentication endpoints, token management, and security considerations.

## Authentication Flow

```mermaid
sequenceDiagram
    participant Client as Admin Client
    participant API as Backend API
    participant JWT as JWT Service
    
    Client->>API: POST /api/v1/auth/login
    API->>JWT: Validate credentials
    JWT->>API: Generate JWT token
    API->>Client: Return token + success
    
    Client->>API: Subsequent requests with Bearer token
    API->>JWT: Validate token
    JWT->>API: Token valid
    API->>Client: Protected resource
```

## Endpoints

### POST /api/v1/auth/login

Authenticates admin users and returns a JWT token for subsequent requests.

#### Request

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "admin",
  "password": "cloudadmin"
}
```

#### Request Body

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| username | string | Yes | Admin username |
| password | string | Yes | Admin password |

#### Response

**Success (200 OK):**
```json
{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**Error (401 Unauthorized):**
```json
{
  "success": false,
  "error": "Invalid username or password"
}
```

**Error (400 Bad Request):**
```json
{
  "success": false,
  "error": "Invalid request: username is required"
}
```

#### JWT Token Structure

The JWT token contains the following claims:

```json
{
  "user_id": "admin",
  "username": "admin",
  "email": "",
  "roles": ["admin"],
  "permissions": [],
  "token_type": "access",
  "iat": 1640995200,
  "exp": 1641081600,
  "nbf": 1640995200,
  "iss": "auth-service",
  "sub": "admin"
}
```

| Claim | Description |
|-------|-------------|
| user_id | Unique user identifier |
| username | Admin username |
| email | User email (currently empty) |
| roles | User roles array |
| permissions | User permissions array |
| token_type | Token type (always "access") |
| iat | Issued at timestamp |
| exp | Expiration timestamp (24 hours) |
| nbf | Not before timestamp |
| iss | Token issuer |
| sub | Token subject |

## Authentication Middleware

All protected endpoints require the `Authorization` header with a valid JWT token.

### Header Format

```http
Authorization: Bearer <jwt_token>
```

### Example Protected Request

```http
GET /api/v1/admin/inquiries
Authorization: Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...
```

### Middleware Validation

The authentication middleware performs the following validations:

1. **Header Presence**: Checks for `Authorization` header
2. **Format Validation**: Ensures "Bearer " prefix
3. **Token Parsing**: Parses JWT token structure
4. **Signature Verification**: Validates token signature with secret key
5. **Expiration Check**: Ensures token hasn't expired
6. **Claims Extraction**: Extracts user information for request context

## Security Considerations

### Token Security

- **Secret Key**: JWT tokens are signed with a secure secret key
- **Expiration**: Tokens expire after 24 hours
- **HTTPS Only**: All authentication should use HTTPS in production
- **Secure Storage**: Frontend stores tokens in localStorage with proper cleanup

### Password Security

- **Demo Credentials**: Current implementation uses hardcoded credentials for demo
- **Production**: Should implement proper password hashing (bcrypt) and user database
- **Rate Limiting**: Consider implementing login attempt rate limiting

### Session Management

- **Token Refresh**: Currently no refresh token mechanism (24-hour expiration)
- **Logout**: Frontend clears token from localStorage
- **Multiple Sessions**: Each login generates a new token

## Frontend Integration

### Login Flow

```typescript
// Login function in AuthContext
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

### Token Management

```typescript
// API service token handling
class ApiService {
  private authToken: string | null = null;

  setAuthToken(token: string | null) {
    this.authToken = token;
    if (token) {
      localStorage.setItem('adminToken', token);
    } else {
      localStorage.removeItem('adminToken');
    }
  }

  private async request<T>(endpoint: string, options: RequestInit = {}): Promise<T> {
    const headers: HeadersInit = {
      'Content-Type': 'application/json',
      ...options.headers,
    };
    
    if (this.authToken && !endpoint.includes('/auth/login')) {
      headers['Authorization'] = `Bearer ${this.authToken}`;
    }
    
    // ... rest of request logic
  }
}
```

### Automatic Logout

The system automatically logs out users when:
- Token expires (401 response)
- Invalid token detected
- Manual logout action

## Error Handling

### Common Error Responses

| Status Code | Error | Description |
|-------------|-------|-------------|
| 400 | Invalid request | Missing or invalid request body |
| 401 | Invalid credentials | Wrong username/password |
| 401 | Invalid token | Expired or malformed JWT token |
| 401 | Missing authorization | No Authorization header |
| 500 | Token generation failed | Server error during token creation |

### Frontend Error Handling

```typescript
// Error handling in login component
const handleSubmit = async (e: React.FormEvent) => {
  e.preventDefault();
  
  try {
    const success = await login(username, password);
    
    if (success) {
      navigate('/admin/dashboard');
    } else {
      setError('Invalid username or password');
    }
  } catch (err: any) {
    setError(err.message || 'Login failed. Please try again.');
  }
};
```

### Code Quality Improvements

The authentication system has been updated to improve code quality:

- **Clean Logging**: Removed debug console.log statements from production code
- **Error-Only Logging**: Maintains error logging for debugging while removing verbose debug output
- **Production Ready**: Cleaner console output in production environments 

## Testing

### Manual Testing

```bash
# Login request
curl -X POST http://localhost:8061/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","password":"cloudadmin"}'

# Protected request
curl -X GET http://localhost:8061/api/v1/admin/inquiries \
  -H "Authorization: Bearer <token>"
```

### Automated Testing

```go
func TestAuthHandler_Login(t *testing.T) {
    handler := NewAuthHandler("test-secret")
    
    tests := []struct {
        name           string
        requestBody    string
        expectedStatus int
        expectedSuccess bool
    }{
        {
            name: "valid credentials",
            requestBody: `{"username":"admin","password":"cloudadmin"}`,
            expectedStatus: http.StatusOK,
            expectedSuccess: true,
        },
        {
            name: "invalid credentials",
            requestBody: `{"username":"admin","password":"wrong"}`,
            expectedStatus: http.StatusUnauthorized,
            expectedSuccess: false,
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Test implementation
        })
    }
}
```

## Migration Notes

### From Demo to Production

1. **Replace hardcoded credentials** with proper user database
2. **Implement password hashing** using bcrypt or similar
3. **Add user management** endpoints for admin user creation
4. **Implement refresh tokens** for better security
5. **Add rate limiting** for login attempts
6. **Enable HTTPS** for all authentication endpoints
7. **Add audit logging** for authentication events

### Database Schema (Future)

```sql
CREATE TABLE admin_users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    roles TEXT[] DEFAULT ARRAY['admin'],
    is_active BOOLEAN DEFAULT true,
    last_login TIMESTAMP,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE auth_sessions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id UUID REFERENCES admin_users(id),
    token_hash VARCHAR(255) NOT NULL,
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);
```

This authentication system provides a secure foundation for the admin dashboard while maintaining simplicity for demo purposes.
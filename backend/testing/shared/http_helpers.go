// Package shared provides HTTP testing utilities
package shared

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// HTTPTestHelper provides utilities for HTTP testing
type HTTPTestHelper struct {
	Router *gin.Engine
	Server *httptest.Server
}

// NewHTTPTestHelper creates a new HTTP test helper
func NewHTTPTestHelper() *HTTPTestHelper {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	return &HTTPTestHelper{
		Router: router,
	}
}

// StartServer starts a test HTTP server
func (h *HTTPTestHelper) StartServer() {
	h.Server = httptest.NewServer(h.Router)
}

// StopServer stops the test HTTP server
func (h *HTTPTestHelper) StopServer() {
	if h.Server != nil {
		h.Server.Close()
	}
}

// GetBaseURL returns the base URL of the test server
func (h *HTTPTestHelper) GetBaseURL() string {
	if h.Server != nil {
		return h.Server.URL
	}
	return ""
}

// HTTPTestRequest represents an HTTP test request
type HTTPTestRequest struct {
	Method      string
	Path        string
	Body        interface{}
	Headers     map[string]string
	QueryParams map[string]string
}

// HTTPTestResponse represents an HTTP test response
type HTTPTestResponse struct {
	StatusCode int
	Body       string
	Headers    http.Header
}

// MakeRequest makes an HTTP request and returns the response
func (h *HTTPTestHelper) MakeRequest(t *testing.T, req *HTTPTestRequest) *HTTPTestResponse {
	t.Helper()

	var bodyReader io.Reader
	if req.Body != nil {
		bodyBytes, err := json.Marshal(req.Body)
		require.NoError(t, err)
		bodyReader = bytes.NewReader(bodyBytes)
	}

	httpReq := httptest.NewRequest(req.Method, req.Path, bodyReader)

	// Set headers
	if req.Headers != nil {
		for key, value := range req.Headers {
			httpReq.Header.Set(key, value)
		}
	}

	// Set query parameters
	if req.QueryParams != nil {
		q := httpReq.URL.Query()
		for key, value := range req.QueryParams {
			q.Add(key, value)
		}
		httpReq.URL.RawQuery = q.Encode()
	}

	// Set default content type for JSON requests
	if req.Body != nil && httpReq.Header.Get("Content-Type") == "" {
		httpReq.Header.Set("Content-Type", "application/json")
	}

	w := httptest.NewRecorder()
	h.Router.ServeHTTP(w, httpReq)

	return &HTTPTestResponse{
		StatusCode: w.Code,
		Body:       w.Body.String(),
		Headers:    w.Header(),
	}
}

// MakeJSONRequest makes a JSON HTTP request
func (h *HTTPTestHelper) MakeJSONRequest(t *testing.T, method, path string, body interface{}) *HTTPTestResponse {
	t.Helper()

	return h.MakeRequest(t, &HTTPTestRequest{
		Method: method,
		Path:   path,
		Body:   body,
		Headers: map[string]string{
			"Content-Type": "application/json",
		},
	})
}

// AssertJSONResponse asserts that the response is valid JSON and matches expected status
func (h *HTTPTestHelper) AssertJSONResponse(t *testing.T, resp *HTTPTestResponse, expectedStatus int, expectedBody interface{}) {
	t.Helper()

	assert.Equal(t, expectedStatus, resp.StatusCode)
	assert.Equal(t, "application/json; charset=utf-8", resp.Headers.Get("Content-Type"))

	if expectedBody != nil {
		expectedJSON, err := json.Marshal(expectedBody)
		require.NoError(t, err)

		assert.JSONEq(t, string(expectedJSON), resp.Body)
	}
}

// AssertErrorResponse asserts that the response contains an error message
func (h *HTTPTestHelper) AssertErrorResponse(t *testing.T, resp *HTTPTestResponse, expectedStatus int, expectedError string) {
	t.Helper()

	assert.Equal(t, expectedStatus, resp.StatusCode)

	var errorResp map[string]interface{}
	err := json.Unmarshal([]byte(resp.Body), &errorResp)
	require.NoError(t, err)

	assert.Contains(t, errorResp["error"], expectedError)
}

// MockHTTPServer provides a mock HTTP server for external service testing
type MockHTTPServer struct {
	Server   *httptest.Server
	Handlers map[string]http.HandlerFunc
}

// NewMockHTTPServer creates a new mock HTTP server
func NewMockHTTPServer() *MockHTTPServer {
	handlers := make(map[string]http.HandlerFunc)

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mockServer := &MockHTTPServer{
		Server:   server,
		Handlers: handlers,
	}

	// Set up a catch-all handler that routes to registered handlers
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		key := fmt.Sprintf("%s %s", r.Method, r.URL.Path)
		if handler, exists := handlers[key]; exists {
			handler(w, r)
		} else {
			http.NotFound(w, r)
		}
	})

	return mockServer
}

// RegisterHandler registers a handler for a specific method and path
func (m *MockHTTPServer) RegisterHandler(method, path string, handler http.HandlerFunc) {
	key := fmt.Sprintf("%s %s", method, path)
	m.Handlers[key] = handler
}

// RegisterJSONHandler registers a handler that returns JSON response
func (m *MockHTTPServer) RegisterJSONHandler(method, path string, statusCode int, response interface{}) {
	m.RegisterHandler(method, path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(statusCode)

		if response != nil {
			json.NewEncoder(w).Encode(response)
		}
	})
}

// RegisterErrorHandler registers a handler that returns an error response
func (m *MockHTTPServer) RegisterErrorHandler(method, path string, statusCode int, errorMessage string) {
	m.RegisterJSONHandler(method, path, statusCode, map[string]string{
		"error": errorMessage,
	})
}

// Close closes the mock HTTP server
func (m *MockHTTPServer) Close() {
	m.Server.Close()
}

// GetURL returns the base URL of the mock server
func (m *MockHTTPServer) GetURL() string {
	return m.Server.URL
}

// HTTPTestCase represents a test case for HTTP endpoints
type HTTPTestCase struct {
	Name           string
	Request        *HTTPTestRequest
	ExpectedStatus int
	ExpectedBody   interface{}
	ExpectedError  string
	SetupMocks     func()
	Assertions     func(t *testing.T, resp *HTTPTestResponse)
}

// RunHTTPTestCases runs a series of HTTP test cases
func (h *HTTPTestHelper) RunHTTPTestCases(t *testing.T, testCases []HTTPTestCase) {
	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			// Setup mocks if provided
			if tc.SetupMocks != nil {
				tc.SetupMocks()
			}

			// Make request
			resp := h.MakeRequest(t, tc.Request)

			// Assert status code
			assert.Equal(t, tc.ExpectedStatus, resp.StatusCode)

			// Assert expected body if provided
			if tc.ExpectedBody != nil {
				h.AssertJSONResponse(t, resp, tc.ExpectedStatus, tc.ExpectedBody)
			}

			// Assert expected error if provided
			if tc.ExpectedError != "" {
				h.AssertErrorResponse(t, resp, tc.ExpectedStatus, tc.ExpectedError)
			}

			// Run custom assertions if provided
			if tc.Assertions != nil {
				tc.Assertions(t, resp)
			}
		})
	}
}

// AuthHelper provides utilities for authentication testing
type AuthHelper struct {
	JWTSecret string
}

// NewAuthHelper creates a new auth helper
func NewAuthHelper(jwtSecret string) *AuthHelper {
	return &AuthHelper{
		JWTSecret: jwtSecret,
	}
}

// CreateTestJWT creates a test JWT token
func (a *AuthHelper) CreateTestJWT(userID, role string) (string, error) {
	// This would implement JWT creation logic
	// For now, return a mock token
	return fmt.Sprintf("test-jwt-token-%s-%s", userID, role), nil
}

// AddAuthHeader adds authentication header to request
func (a *AuthHelper) AddAuthHeader(req *HTTPTestRequest, token string) {
	if req.Headers == nil {
		req.Headers = make(map[string]string)
	}
	req.Headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
}

// WebSocketTestHelper provides utilities for WebSocket testing
type WebSocketTestHelper struct {
	// This will be expanded when WebSocket tests are moved
}

// NewWebSocketTestHelper creates a new WebSocket test helper
func NewWebSocketTestHelper() *WebSocketTestHelper {
	return &WebSocketTestHelper{}
}

// RequestBuilder provides a fluent interface for building HTTP requests
type RequestBuilder struct {
	request *HTTPTestRequest
}

// NewRequestBuilder creates a new request builder
func NewRequestBuilder() *RequestBuilder {
	return &RequestBuilder{
		request: &HTTPTestRequest{
			Headers:     make(map[string]string),
			QueryParams: make(map[string]string),
		},
	}
}

// Method sets the HTTP method
func (rb *RequestBuilder) Method(method string) *RequestBuilder {
	rb.request.Method = method
	return rb
}

// Path sets the request path
func (rb *RequestBuilder) Path(path string) *RequestBuilder {
	rb.request.Path = path
	return rb
}

// Body sets the request body
func (rb *RequestBuilder) Body(body interface{}) *RequestBuilder {
	rb.request.Body = body
	return rb
}

// Header adds a header to the request
func (rb *RequestBuilder) Header(key, value string) *RequestBuilder {
	rb.request.Headers[key] = value
	return rb
}

// QueryParam adds a query parameter to the request
func (rb *RequestBuilder) QueryParam(key, value string) *RequestBuilder {
	rb.request.QueryParams[key] = value
	return rb
}

// JSON sets the content type to JSON
func (rb *RequestBuilder) JSON() *RequestBuilder {
	rb.request.Headers["Content-Type"] = "application/json"
	return rb
}

// Auth adds authorization header
func (rb *RequestBuilder) Auth(token string) *RequestBuilder {
	rb.request.Headers["Authorization"] = fmt.Sprintf("Bearer %s", token)
	return rb
}

// Build returns the built request
func (rb *RequestBuilder) Build() *HTTPTestRequest {
	return rb.request
}

// Common HTTP test methods

// GET creates a GET request builder
func GET(path string) *RequestBuilder {
	return NewRequestBuilder().Method("GET").Path(path)
}

// POST creates a POST request builder
func POST(path string) *RequestBuilder {
	return NewRequestBuilder().Method("POST").Path(path).JSON()
}

// PUT creates a PUT request builder
func PUT(path string) *RequestBuilder {
	return NewRequestBuilder().Method("PUT").Path(path).JSON()
}

// DELETE creates a DELETE request builder
func DELETE(path string) *RequestBuilder {
	return NewRequestBuilder().Method("DELETE").Path(path)
}

// PATCH creates a PATCH request builder
func PATCH(path string) *RequestBuilder {
	return NewRequestBuilder().Method("PATCH").Path(path).JSON()
}

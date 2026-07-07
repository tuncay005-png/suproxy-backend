package testutil

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// HTTPTestContext provides utilities for HTTP testing
type HTTPTestContext struct {
	Router   *gin.Engine
	Recorder *httptest.ResponseRecorder
	t        *testing.T
}

// NewHTTPTestContext creates a new HTTP test context
func NewHTTPTestContext(t *testing.T) *HTTPTestContext {
	gin.SetMode(gin.TestMode)

	return &HTTPTestContext{
		Router:   gin.New(),
		Recorder: httptest.NewRecorder(),
		t:        t,
	}
}

// MakeRequest makes an HTTP request for testing
func (ctx *HTTPTestContext) MakeRequest(method, path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	ctx.t.Helper()

	var bodyReader io.Reader
	if body != nil {
		jsonBody, err := json.Marshal(body)
		require.NoError(ctx.t, err, "Failed to marshal request body")
		bodyReader = bytes.NewBuffer(jsonBody)
	}

	req := httptest.NewRequest(method, path, bodyReader)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add custom headers
	for key, value := range headers {
		req.Header.Set(key, value)
	}

	ctx.Recorder = httptest.NewRecorder()
	ctx.Router.ServeHTTP(ctx.Recorder, req)

	return ctx.Recorder
}

// GET makes a GET request
func (ctx *HTTPTestContext) GET(path string, headers map[string]string) *httptest.ResponseRecorder {
	return ctx.MakeRequest(http.MethodGet, path, nil, headers)
}

// POST makes a POST request
func (ctx *HTTPTestContext) POST(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return ctx.MakeRequest(http.MethodPost, path, body, headers)
}

// PUT makes a PUT request
func (ctx *HTTPTestContext) PUT(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return ctx.MakeRequest(http.MethodPut, path, body, headers)
}

// PATCH makes a PATCH request
func (ctx *HTTPTestContext) PATCH(path string, body interface{}, headers map[string]string) *httptest.ResponseRecorder {
	return ctx.MakeRequest(http.MethodPatch, path, body, headers)
}

// DELETE makes a DELETE request
func (ctx *HTTPTestContext) DELETE(path string, headers map[string]string) *httptest.ResponseRecorder {
	return ctx.MakeRequest(http.MethodDelete, path, nil, headers)
}

// AssertStatusCode asserts the response status code
func (ctx *HTTPTestContext) AssertStatusCode(expected int) {
	ctx.t.Helper()
	require.Equal(ctx.t, expected, ctx.Recorder.Code, "Status code mismatch")
}

// AssertJSONResponse asserts and decodes JSON response
func (ctx *HTTPTestContext) AssertJSONResponse(expected int, target interface{}) {
	ctx.t.Helper()

	ctx.AssertStatusCode(expected)

	err := json.Unmarshal(ctx.Recorder.Body.Bytes(), target)
	require.NoError(ctx.t, err, "Failed to unmarshal response body")
}

// GetResponseBody returns the response body as string
func (ctx *HTTPTestContext) GetResponseBody() string {
	return ctx.Recorder.Body.String()
}

// GetResponseJSON decodes response body as JSON
func (ctx *HTTPTestContext) GetResponseJSON(target interface{}) {
	ctx.t.Helper()

	err := json.Unmarshal(ctx.Recorder.Body.Bytes(), target)
	require.NoError(ctx.t, err, "Failed to unmarshal response body")
}

// AuthHeader creates an authorization header
func AuthHeader(token string) map[string]string {
	return map[string]string{
		"Authorization": "Bearer " + token,
	}
}

// JSONHeaders returns JSON content-type headers
func JSONHeaders() map[string]string {
	return map[string]string{
		"Content-Type": "application/json",
	}
}

// MergeHeaders merges multiple header maps
func MergeHeaders(headers ...map[string]string) map[string]string {
	result := make(map[string]string)
	for _, h := range headers {
		for k, v := range h {
			result[k] = v
		}
	}
	return result
}


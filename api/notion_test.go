package api

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"testing"
)

// MockHttpClient is a mock implementation of HttpClientInterface for testing
type MockHttpClient struct {
	Response *http.Response // The mock response to return
	Err      error          // The mock error to return
}

// Do is the mock of the HttpClientInterface Do method
func (m *MockHttpClient) Do(req *http.Request) (*http.Response, error) {
	return m.Response, m.Err
}

// Example test for a successful API call
func TestCallAPI_Success(t *testing.T) {
	mockRespBody := `{"results":[{"id":"exampleID","paragraph":{"rich_text":[{"type":"text","text":{"content":"Example content"}}]}}]}`
	mockClient := &MockHttpClient{
		Response: &http.Response{
			StatusCode: 200,
			Body:       io.NopCloser(bytes.NewBufferString(mockRespBody)),
			Header:     make(http.Header),
		},
		Err: nil,
	}
	apiClient := NewNotionApiClient(mockClient)

	err := apiClient.CallAPI("validID", "validToken")
	if err != nil {
		t.Errorf("CallAPI() error = %v, wantErr %v", err, nil)
	}
	// Additional assertions can be added here to check the behavior
}

func TestCallAPI_ErrorResponse(t *testing.T) {
	// Setup
	mockRespBody := `{"object":"error","status":400,"code":"invalid_request","message":"Invalid Request"}`
	mockClient := &MockHttpClient{
		Response: &http.Response{
			StatusCode: 400,
			Body:       io.NopCloser(bytes.NewBufferString(mockRespBody)),
		},
		Err: nil,
	}
	apiClient := NewNotionApiClient(mockClient)

	// Execute
	err := apiClient.CallAPI("invalidID", "validToken")

	// Assert
	if err == nil {
		t.Errorf("Expected an error, got nil")
	}
	if err != nil && err.Error() != "API Error: invalid_request - Invalid Request" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestCallAPI_HttpError(t *testing.T) {
	// Setup
	mockClient := &MockHttpClient{
		Response: nil,
		Err:      errors.New("network error"),
	}
	apiClient := NewNotionApiClient(mockClient)

	// Execute
	err := apiClient.CallAPI("anyID", "validToken")

	// Assert
	if err == nil || err.Error() != "error sending request: network error" {
		t.Errorf("Expected network error, got %v", err)
	}
}

package api

import (
	"bytes"
	"io/ioutil"
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
			Body:       ioutil.NopCloser(bytes.NewBufferString(mockRespBody)),
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

// Example test for an API error response
func TestCallAPI_ApiError(t *testing.T) {
	// Similar to TestCallAPI_Success, but simulate an API error response
}

// Example test for a network error
func TestCallAPI_NetworkError(t *testing.T) {
	// Similar to TestCallAPI_Success, but simulate a network error using the Err field of MockHttpClient
}

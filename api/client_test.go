package api

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"testing"
)

type MockHTTPClient struct {
	MockDo func(req *http.Request) (*http.Response, error)
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	return m.MockDo(req)
}

func TestGetNotionChildBlocks(t *testing.T) {
	testCases := []struct {
		name           string
		mockResponse   string
		mockStatusCode int
		mockErr        error
		expectErr      bool
	}{
		{
			name:           "Successful Fetch",
			mockResponse:   `{"object":"list","results":[{"object":"block","type":"heading_1","heading_1":{"rich_text":[{"type":"text","text":{"content":"I am a heading one"}}]}}]}`,
			mockStatusCode: http.StatusOK,
			expectErr:      false,
		},
		{
			name:           "API Error",
			mockResponse:   `{"object":"error","status":400,"code":"validation_error","message":"Invalid block ID"}`,
			mockStatusCode: http.StatusBadRequest,
			expectErr:      true,
		},
		{
			name:      "HTTP Client Error",
			mockErr:   fmt.Errorf("network error"),
			expectErr: true,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			responseBody := io.NopCloser(bytes.NewReader([]byte(tc.mockResponse)))
			mockClient := &MockHTTPClient{
				MockDo: func(req *http.Request) (*http.Response, error) {
					return &http.Response{
						StatusCode: tc.mockStatusCode,
						Body:       responseBody,
					}, tc.mockErr
				},
			}
			client := NewNotionApiClient(mockClient)
			blockID, bearerToken := "test-block-id", "test-bearer-token"
			_, err := client.GetNotionChildBlocks(blockID, bearerToken)

			// Check if an error was expected
			if tc.expectErr && err == nil {
				t.Errorf("Expected an error but did not get one")
			} else if !tc.expectErr && err != nil {
				t.Errorf("Did not expect an error but got one: %v", err)
			}
		})
	}
}

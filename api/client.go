package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// NotionApiClient struct holds any dependencies for your API client, e.g., the HTTP client.
type NotionApiClient struct {
	Client HttpClientInterface
}

// NewNotionApiClient creates a new API client with the provided HTTP client.
func NewNotionApiClient(client HttpClientInterface) *NotionApiClient {
	return &NotionApiClient{
		Client: client,
	}
}

// GetNotionBlocks performs the actual API call to retrieve the blocks and processes the response.
func (api *NotionApiClient) GetNotionBlocks(blockID, bearerToken string) (*ResultsWrapper, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", blockID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	//@TODO REMOVE THIS
	req.Header.Add("Authorization", "Bearer secret_hVDPuHdW5ec7WzM2WicFHNCT7dWy8F5mOE9MMIY2PjK")
	// req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Notion-Version", "2022-06-28")

	resp, err := api.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		var apiError APIErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return nil, fmt.Errorf("error parsing API error response: %w", err)
		}
		return nil, fmt.Errorf("API Error: %s - %s", apiError.Code, apiError.Message)
	}

	// Correctly read the response body here
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading response body: %w", err)
	}

	var results ResultsWrapper
	if err := json.Unmarshal(body, &results); err != nil {
		return nil, fmt.Errorf("error parsing JSON: %w", err)
	}

	return &results, nil
}

// HttpClientInterface defines the interface for the HTTP client, allowing for easy mocking/testing.
type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

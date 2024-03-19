package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

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

// NotionAPI defines the interface for interacting with the Notion API.
type NotionAPI interface {
	GetNotionBlockTitle(blockID, bearerToken string) (string, error)
	GetNotionChildBlocks(blockID, bearerToken string) (*ResultsWrapper, error)
}

var _ NotionAPI = (*NotionApiClient)(nil)

func FetchBlockTitle(apiClient NotionAPI, pageID, bearerToken string) (string, error) {
	return apiClient.GetNotionBlockTitle(pageID, bearerToken)
}

func FetchChildBlocks(apiClient NotionAPI, blockID, bearerToken string) (*ResultsWrapper, error) {
	return apiClient.GetNotionChildBlocks(blockID, bearerToken)
}

// GetNotionBlockTitle makes an API request to Notion to get the title of a block by its ID.
func (api *NotionApiClient) GetNotionBlockTitle(blockID, bearerToken string) (string, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s", blockID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Notion-Version", "2022-06-28")

	resp, err := api.Client.Do(req)
	if err != nil {
		return "", fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return "", fmt.Errorf("API request failed with status code: %d", resp.StatusCode)
	}

	var blockTitleResponse BlockTitleResponse
	if err := json.NewDecoder(resp.Body).Decode(&blockTitleResponse); err != nil {
		return "", fmt.Errorf("error parsing JSON: %w", err)
	}

	return blockTitleResponse.ChildPage.Title, nil
}

func (api *NotionApiClient) GetNotionChildBlocks(blockID, bearerToken string) (*ResultsWrapper, error) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", blockID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Add("Authorization", "Bearer "+bearerToken)
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

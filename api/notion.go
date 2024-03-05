package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

// APIErrorResponse represents JSON error responses from the API
type APIErrorResponse struct {
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// ResultsWrapper is the structure of your successful response
type ResultsWrapper struct {
	Results []Block `json:"results"`
}

type Block struct {
	ID        string     `json:"id"`
	Type      string     `json:"type"`
	Heading1  *Heading   `json:"heading_1,omitempty"`
	Heading2  *Heading   `json:"heading_2,omitempty"`
	Heading3  *Heading   `json:"heading_3,omitempty"`
	Paragraph *Paragraph `json:"paragraph,omitempty"`
}

// Heading represents a generic heading, which can be used for both heading_1, heading_2, heading_3 etc.
type Heading struct {
	RichText []RichText `json:"rich_text"`
}
type Paragraph struct {
	RichText []RichText `json:"rich_text"`
}

type RichText struct {
	Type string `json:"type"`
	Text struct {
		Content string  `json:"content"`
		Link    *string `json:"link,omitempty"`
	} `json:"text"`
	Annotations struct {
		Bold          bool   `json:"bold"`
		Italic        bool   `json:"italic"`
		Strikethrough bool   `json:"strikethrough"`
		Underline     bool   `json:"underline"`
		Code          bool   `json:"code"`
		Color         string `json:"color"`
	} `json:"annotations"`
	PlainText string  `json:"plain_text"`
	Href      *string `json:"href,omitempty"`
}

// RichTextProvider interface for blocks that contain Rich Text
type RichTextProvider interface {
	GetRichText() []RichText
}

// Implement GetRichText for Heading
func (h *Heading) GetRichText() []RichText {
	return h.RichText
}

// Implement GetRichText for Paragraph
func (p *Paragraph) GetRichText() []RichText {
	return p.RichText
}

// HttpClientInterface defines the interface for the HTTP client
type HttpClientInterface interface {
	Do(req *http.Request) (*http.Response, error)
}

// NotionApiClient struct will hold any dependencies for your API client, e.g., the HTTP client
type NotionApiClient struct {
	Client HttpClientInterface
}

// NewNotionApiClient creates a new API client with the provided HTTP client
func NewNotionApiClient(client HttpClientInterface) *NotionApiClient {
	return &NotionApiClient{
		Client: client,
	}
}

// CallAPI performs the actual API call and processes the response
func (api *NotionApiClient) CallAPI(customID, bearerToken string) error {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", customID)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// @todo we'll need the user to provide their own bearer token
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Notion-Version", "2022-06-28")

	resp, err := api.Client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		var apiError APIErrorResponse
		if err := json.NewDecoder(resp.Body).Decode(&apiError); err != nil {
			return fmt.Errorf("error parsing API error response: %w", err)
		}
		return fmt.Errorf("API Error: %s - %s", apiError.Code, apiError.Message)
	} else if resp.StatusCode != 200 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}

	var results ResultsWrapper
	if err := json.Unmarshal(body, &results); err != nil {
		return fmt.Errorf("error parsing JSON: %w", err)
	}

	// Open a new file for writing, create if not exists, truncate if exists
	outputPath := "test.md"
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating markdown file: %w", err)
	}
	defer file.Close()

	for _, block := range results.Results {
		fmt.Printf("Block ID: %s\n", block.ID)
		fmt.Printf("Block Type: %s\n", block.Type)

		var provider RichTextProvider

		switch block.Type {
		case "heading_1":
			provider = block.Heading1
		case "heading_2":
			provider = block.Heading2
		case "heading_3":
			provider = block.Heading3
		case "paragraph":
			provider = block.Paragraph
		}

		if provider != nil {
			for _, rt := range provider.GetRichText() {
				// @TODO: Depending on the type of Rich Text, I need to format the content differently
				_, err := file.WriteString(rt.Text.Content + "\n")
				if err != nil {
					return fmt.Errorf("error writing to markdown file: %w", err)
				}
				fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
			}
		}
	}

	return nil // Return nil if everything was successful
}

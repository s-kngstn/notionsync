package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/s-kngstn/notionsync/api"
	"github.com/s-kngstn/notionsync/format"
)

// WriteBlocksToMarkdown writes the contents of fetched blocks to a Markdown file.
func main() {
	var uuid string
	var err error

	for {
		userInput := RealUserInput{}
		token := Prompt(userInput, "Please enter the Notion API bearer token: ")
		PersistToken(token)

		url := Prompt(userInput, "Please enter the Notion page URL: ")
		if url == "" {
			fmt.Println("URL is required, please try again.")
			continue
		}

		uuid, err = FetchDataBlockString(url)
		if err != nil {
			fmt.Printf("Error: %v. Please try again.\n", err)
			continue // If an error occurs (e.g., no UUID found), prompt for the URL again
		}

		fmt.Printf("Extracted UUID: %s\n", uuid)
		break // Exit the loop if a valid UUID is found
	}

	// Initialize the API client with http.Client
	client := &http.Client{}
	apiClient := api.NewNotionApiClient(client)
	// @todo have user provide their own bearer token
	// Set the bearer token
	// bearerToken := "secret_hVDPuHdW5ec7WzM2WicFHNCT7dWy8F5mOE9MMIY2PjK"
	bearerToken := os.Getenv("NOTION_BEARER_TOKEN")

	// Call the API with the extracted UUID
	results, err := apiClient.GetNotionBlocks(uuid, bearerToken)
	if err != nil {
		fmt.Println("Error calling API:", err)
		return
	}

	// Now process and write the results to a Markdown file
	outputPath := "output/test.md"
	if err := format.WriteBlocksToMarkdown(results, outputPath); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
	}
}

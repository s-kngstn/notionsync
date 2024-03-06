package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/s-kngstn/notionsync/api"
	"github.com/s-kngstn/notionsync/format"
)

func main() {
	var uuid, url string
	var err error

	userInput := RealUserInput{}
	token := Prompt(userInput, "Please enter the Notion API bearer token: ")
	PersistToken(token)

	for {
		url = Prompt(userInput, "Please enter the Notion page URL: ")
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

	// Extract the page name from the URL and use it as the filename
	pageName, err := ExtractNameFromURL(url)
	if err != nil {
		fmt.Println("Error extracting page name from URL:", err)
		return
	}
	outputPath := fmt.Sprintf("output/%s.md", pageName)

	client := &http.Client{}
	apiClient := api.NewNotionApiClient(client)

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
	// outputPath := "output/test.md"
	if err := format.WriteBlocksToMarkdown(results, outputPath); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
	}
}

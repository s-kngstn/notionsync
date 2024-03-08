package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/s-kngstn/notionsync/api"
	"github.com/s-kngstn/notionsync/format"
	"github.com/s-kngstn/notionsync/pkg/cli"
	"github.com/s-kngstn/notionsync/pkg/fetch"
	"github.com/s-kngstn/notionsync/pkg/token"
)

func main() {
	userInput := cli.RealUserInput{}
	tokenValue := cli.Prompt(userInput, "Please enter the Notion API bearer token: ")
	token.PersistToken(tokenValue)
	uuid, url := cli.PromptForURL(userInput)

	// Extract the page name from the URL and use it as the filename
	pageName, err := fetch.ExtractNameFromURL(url)
	if err != nil {
		fmt.Println("Error extracting page name from URL:", err)
		return
	}

	client := &http.Client{}
	apiClient := api.NewNotionApiClient(client)
	bearerToken := os.Getenv("NOTION_BEARER_TOKEN")

	results, err := apiClient.GetNotionBlocks(uuid, bearerToken)
	if err != nil {
		fmt.Println("Error calling API:", err)
		return
	}
	processedBlocks := make(map[string]string) // Initialize the map

	outputPath := fmt.Sprintf("output/%s.md", pageName)
	format.ProcessBlocks(uuid, results, outputPath, pageName, apiClient, bearerToken, processedBlocks)
	// Once processing is complete, print the map to view the processed blocks
	// This is just for debugging purposes @TODO: Remove this
	for blockID, filePath := range processedBlocks {
		fmt.Printf("BlockID: %s, FilePath: %s\n", blockID, filePath)
	}
}

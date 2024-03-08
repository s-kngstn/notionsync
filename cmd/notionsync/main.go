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

	results, err := apiClient.GetNotionChildBlocks(uuid, bearerToken)
	if err != nil {
		fmt.Println("Error calling API:", err)
		return
	}
	// Before we process the blocks we need to check if the type is a link_to_page:
	// - If it is we need to get the page_id from that block, and then check if that page has been processed Before
	// by checking the processedBlocks map. If it has been processed before we can just link to the file that was created
	// for that page.
	// - If it has not been processed before we need to process that page and then link to the file that was created

	// We want to keep track of the processed blocks
	processedBlocks := make(map[string]string)

	outputPath := fmt.Sprintf("output/%s.md", pageName)
	format.ProcessBlocks(uuid, results, outputPath, pageName, apiClient, bearerToken, processedBlocks)

	// This is just for debugging purposes @TODO: Remove this
	// for blockID, filePath := range processedBlocks {
	// 	fmt.Printf("BlockID: %s, FilePath: %s\n", blockID, filePath)
	// }
}

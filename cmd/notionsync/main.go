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

	outputPath := fmt.Sprintf("output/%s.md", pageName)
	format.ProcessBlocks(results, outputPath, pageName, apiClient, bearerToken)
}

// func processBlocks(results *api.ResultsWrapper, outputPath string, pageName string, apiClient *api.NotionApiClient, bearerToken string) {
// 	// We need to handle if we have a child page block has children and if so, we need to call the API again
// 	// The ID of the child page is the UUID of the block
// 	for _, block := range results.Results {
// 		if block.Type == "child_page" && block.HasChildren {
// 			// Call the API with the ID of the block
// 			childResults, err := apiClient.GetNotionBlocks(block.ID, bearerToken)
// 			if err != nil {
// 				fmt.Println("Error calling API:", err)
// 				return
// 			}
// 			// The output path will be the same as the parent page, but with the title of the child page name
// 			childPageName := format.ToKebabCase(block.ChildPage.Title)
// 			childOutputPath := fmt.Sprintf("output/%s.md", childPageName)
// 			// Write the child blocks to a Markdown file with the title of the child pageName
// 			if err := format.WriteBlocksToMarkdown(childResults, childOutputPath, childPageName); err != nil {
// 				fmt.Println("Error writing blocks to Markdown:", err)
// 			}
// 		}
// 	}

// 	// Now process and write the results to a Markdown file
// 	if err := format.WriteBlocksToMarkdown(results, outputPath, pageName); err != nil {
// 		fmt.Println("Error writing blocks to Markdown:", err)
// 	}
// }

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

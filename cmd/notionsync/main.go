package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"
	"github.com/s-kngstn/notionsync/api"
	"github.com/s-kngstn/notionsync/format"
	"github.com/s-kngstn/notionsync/pkg/cli"
	"github.com/s-kngstn/notionsync/pkg/fetch"
	"github.com/s-kngstn/notionsync/pkg/utils"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file or no .env file found:", err)
	}

	appEnv := os.Getenv("APP_ENV")

	tokenFlag := flag.String("token", "", "Notion API bearer token")
	filePath := flag.String("file", "", "Path to the file containing URLs")
	flag.Parse()

	if appEnv == "development" {
		if err := godotenv.Load(); err != nil {
			fmt.Println("No .env file found")
		}
	}

	bearerToken := os.Getenv("NOTION_API_KEY")

	// If in production and token not provided through flag, prompt for it
	if appEnv != "development" && *tokenFlag == "" {
		// Initialize RealUserInput with os.Stdin
		if bearerToken == "" {
			inputReader := bufio.NewReader(os.Stdin)
			userInput := cli.NewRealUserInput(inputReader)
			bearerToken = cli.Prompt(userInput, "Please enter the Notion API bearer token: ")
		}
	} else if *tokenFlag != "" {
		// If token is provided through flag, use it
		bearerToken = *tokenFlag
	}

	var urls []string

	if *filePath == "" {
		// No file path provided, prompt for a single URL
		inputReader := bufio.NewReader(os.Stdin)
		userInput := cli.NewRealUserInput(inputReader)
		url := cli.Prompt(userInput, "Please enter the Notion page URL: ")
		urls = append(urls, url)
	} else {
		// File path provided, read URLs from the file
		urls, err = utils.ReadURLs(*filePath)
		if err != nil {
			fmt.Printf("Failed to read URLs from file: %v\n", err)
			return
		}
	}

	client := &http.Client{}
	apiClient := api.NewNotionApiClient(client)

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish
	// This mutal exclusion is for safely updating the processedBlocks map
	var mu sync.Mutex

	processedBlocks := make(map[string]map[string]string)

	for _, url := range urls {
		wg.Add(1)
		go func(url string) {
			defer wg.Done()
			processURL(url, apiClient, bearerToken, &mu, processedBlocks)
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()
}

// processURL handles the processing of a single URL
func processURL(url string, apiClient api.NotionAPI, bearerToken string, mu *sync.Mutex, processedBlocks map[string]map[string]string) {
	blockIDFetcher := fetch.DefaultBlockIDFetcher{}
	uuid, err := blockIDFetcher.GetBlockID(url)
	if err != nil {
		fmt.Printf("Error extracting UUID from URL %s: %v\n", url, err)
		return
	}
	pageName, err := fetch.ExtractNameFromURL(url)
	if err != nil {
		fmt.Printf("Error extracting page name from URL %s: %v\n", url, err)
		return
	}

	results, err := api.FetchChildBlocks(apiClient, uuid, bearerToken)
	if err != nil {
		fmt.Printf("Error calling API for URL %s: %v\n", url, err)
		return
	}

	outputPath := fmt.Sprintf("output/%s.md", pageName)

	// Initialize the inner map if it doesn't exist
	mu.Lock()
	if processedBlocks[uuid] == nil {
		processedBlocks[uuid] = make(map[string]string)
	}
	processedBlocks[uuid]["outputPath"] = outputPath
	// You can add more entries to processedBlocks[uuid] as needed
	mu.Unlock()

	format.ProcessBlocks(uuid, results, outputPath, pageName, apiClient, bearerToken, processedBlocks[uuid])
}

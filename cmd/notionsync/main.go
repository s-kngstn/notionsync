package main

import (
	"bufio"
	"flag"
	"fmt"
	"net/http"
	"os"
	"sync"

	"github.com/s-kngstn/notionsync/api"
	"github.com/s-kngstn/notionsync/format"
	"github.com/s-kngstn/notionsync/pkg/cli"
	"github.com/s-kngstn/notionsync/pkg/fetch"
	"github.com/s-kngstn/notionsync/pkg/utils"
)

func main() {
	var err error
	bearerToken := os.Getenv("NOTION_API_KEY")
	tokenFlag := flag.String("token", "", "Notion API bearer token")
	filePath := flag.String("file", "", "Path to the file containing URLs to process")
	outputDir := flag.String("dir", "notion-notes", "Directory to save markdown files in")
	flag.Parse()

	// Ensure output directory exists
	if _, err := os.Stat(*outputDir); os.IsNotExist(err) {
		os.Mkdir(*outputDir, 0755)
	}

	if bearerToken == "" && *tokenFlag == "" {
		// Initialize RealUserInput with os.Stdin
		fmt.Println("No Notion API Token found in env[`NOTION_API_KEY`] or flag provided")
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
			processURL(url, apiClient, bearerToken, &mu, processedBlocks, *outputDir)
		}(url)
	}

	// Wait for all goroutines to finish
	wg.Wait()
	fmt.Println("All URLs processed")
}

// processURL handles the processing of a single URL
func processURL(url string, apiClient api.NotionAPI, bearerToken string, mu *sync.Mutex, processedBlocks map[string]map[string]string, outputDir string) {
	// Checking if the URL is a notion page
	urlChecker := fetch.DefaultURLChecker{}
	urlIsValid, err := urlChecker.CheckURL(url)
	if err != nil {
		fmt.Printf("%v\n", err)
	}

	if !urlIsValid {
		return
	}
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
		// @todo if there is a error here we need to do a graceful shutdown and have the user try again
		return
	}

	outputPath := fmt.Sprintf("%s/%s.md", outputDir, pageName)
	// Initialize the inner map if it doesn't exist
	mu.Lock()
	if processedBlocks[uuid] == nil {
		processedBlocks[uuid] = make(map[string]string)
	}
	processedBlocks[uuid]["outputPath"] = outputPath
	// You can add more entries to processedBlocks[uuid] as needed
	mu.Unlock()

	format.ProcessBlocks(uuid, results, outputPath, pageName, apiClient, bearerToken, processedBlocks[uuid], outputDir)
}

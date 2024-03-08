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
	"github.com/s-kngstn/notionsync/pkg/token"
)

func readURLs(filePath string) ([]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var urls []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		urls = append(urls, scanner.Text())
	}
	return urls, scanner.Err()
}

func main() {
	defaultFilePath := "PLACEHOLDER" // Indicates no file path was provided
	filePath := flag.String("file", defaultFilePath, "Path to the file containing URLs")
	flag.Parse()

	userInput := cli.RealUserInput{}
	tokenValue := cli.Prompt(userInput, "Please enter the Notion API bearer token: ")
	token.PersistToken(tokenValue)

	var urls []string
	var err error

	if *filePath == defaultFilePath {
		// No file path provided, prompt for a single URL
		userInput := cli.RealUserInput{}
		url := cli.Prompt(userInput, "Please enter the Notion page URL: ")
		urls = append(urls, url)
	} else {
		// File path provided, read URLs from the file
		urls, err = readURLs(*filePath)
		if err != nil {
			fmt.Printf("Failed to read URLs from file: %v\n", err)
			return
		}
	}

	client := &http.Client{}
	apiClient := api.NewNotionApiClient(client)
	bearerToken := os.Getenv("NOTION_BEARER_TOKEN")

	var wg sync.WaitGroup // WaitGroup to wait for all goroutines to finish
	// Mutex for safely updating the processedBlocks map
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
	uuid, err := fetch.GetBlockID(url)
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

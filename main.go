package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/s-kngstn/notionsync/api" // This is the adjusted import path
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	var uuid string
	var err error

	for {
		userInput := RealUserInput{}
		token := PromptForToken(userInput)
		PersistToken(token)

		fmt.Print("Please enter the URL: ")
		url, _ := reader.ReadString('\n')
		url = strings.TrimSpace(url)

		// Check if the URL was provided
		if url == "" {
			fmt.Println("URL is required, please try again.")
			continue // Skip the rest of the loop and prompt again
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
	apiClient := api.NewNotionApiClient(client) // Assuming this is the constructor function
	// @todo have user provide their own bearer token
	// Set the bearer token
	// bearerToken := "secret_hVDPuHdW5ec7WzM2WicFHNCT7dWy8F5mOE9MMIY2PjK"
	bearerToken := os.Getenv("NOTION_BEARER_TOKEN")

	// Call the API with the extracted UUID
	err = apiClient.CallAPI(uuid, bearerToken)
	if err != nil {
		fmt.Printf("API call error: %v\n", err)
	} else {
		fmt.Println("API call was successful.")
	}
}

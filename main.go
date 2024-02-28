package main

import (
	"flag"
	"fmt"
)

func main() {
	// Define a command-line flag
	var url string
	flag.StringVar(&url, "url", "", "URL to extract the UUID from")
	flag.Parse()

	// Check if the URL was provided
	if url == "" {
		fmt.Println("URL is required")
		return
	}

	uuid, err := FetchDataBlockString(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Extracted UUID: %s\n", uuid)
	}
}

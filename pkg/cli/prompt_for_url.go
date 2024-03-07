package cli

import (
	"fmt"

	"github.com/s-kngstn/notionsync/pkg/fetch"
)

func PromptForURL(userInput UserInput) (string, string) {
	for {
		url := Prompt(userInput, "Please enter the Notion page URL: ")
		if url != "" {
			uuid, err := fetch.GetBlockID(url)
			if err != nil {
				fmt.Printf("Error: %v. Please try again.\n", err)
				continue
			}
			fmt.Printf("Preparing to sync..\n")
			return uuid, url
		}
		fmt.Println("URL is required, please try again.")
	}
}

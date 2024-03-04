package main

import (
	"fmt"
	"os"
)

// PromptForToken asks the user for the Notion API bearer token and returns it.
func PromptForToken(ui UserInput) string {
	return ui.ReadString("Please enter the Notion API bearer token: ")
}

// PersistToken stores the token in an environment variable.
func PersistToken(token string) {
	os.Setenv("NOTION_BEARER_TOKEN", token)
	fmt.Println("Notion API bearer token has been saved.")
}

/**
* Environment Variables do not persist after the program exits.
* But lets keep this logic here incase we want to use a different method to store the token later.
 */
func RemoveToken() {
	os.Unsetenv("NOTION_BEARER_TOKEN")
	fmt.Println("Notion API bearer token removed.")
}

// PromptToRemoveToken asks the user if they would like to remove the token from the environment.
func PromptToRemoveToken(ui UserInput) bool {
	response := ui.ReadString("Would you like to remove the Notion API bearer token? (y/N): ")
	return response == "y"
}

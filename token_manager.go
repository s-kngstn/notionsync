package main

import (
	"fmt"
	"os"
)

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

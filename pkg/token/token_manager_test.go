package token

import (
	"os"
	"testing"
)

func TestPersistToken(t *testing.T) {
	token := "test_token"
	PersistToken(token)
	if os.Getenv("NOTION_BEARER_TOKEN") != token {
		t.Errorf("PersistToken() failed to set the environment variable")
	}
	// Cleanup
	os.Unsetenv("NOTION_BEARER_TOKEN")
}

func TestRemoveToken(t *testing.T) {
	// Setup: Set the NOTION_BEARER_TOKEN environment variable
	testToken := "test_token"
	os.Setenv("NOTION_BEARER_TOKEN", testToken)

	// Act: Call RemoveToken to unset the NOTION_BEARER_TOKEN
	RemoveToken()

	// Assert: Check if NOTION_BEARER_TOKEN is unset
	if os.Getenv("NOTION_BEARER_TOKEN") != "" {
		t.Errorf("RemoveToken() failed to unset the NOTION_BEARER_TOKEN environment variable")
	}
}

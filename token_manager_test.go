// token_manager_test.go
package main

import (
	"os"
	"testing"
)

func TestPromptForToken(t *testing.T) {
	expectedToken := "test_token"
	ui := NewMockUserInput(map[string]string{
		"Please enter the Notion API bearer token: ": expectedToken,
	})

	token := PromptForToken(ui)
	if token != expectedToken {
		t.Errorf("PromptForToken() = %v, want %v", token, expectedToken)
	}
}

func TestPersistToken(t *testing.T) {
	token := "test_token"
	PersistToken(token)
	if os.Getenv("NOTION_BEARER_TOKEN") != token {
		t.Errorf("PersistToken() failed to set the environment variable")
	}
	// Cleanup
	os.Unsetenv("NOTION_BEARER_TOKEN")
}

func TestPromptToRemoveToken(t *testing.T) {
	ui := NewMockUserInput(map[string]string{
		"Would you like to remove the Notion API bearer token? (y/N): ": "y",
	})

	if !PromptToRemoveToken(ui) {
		t.Errorf("PromptToRemoveToken() did not return true for 'y' response")
	}
}

// Testing RemoveToken is tricky because it affects the environment.
// If your test suite relies on environment variables, altering them could lead to unpredictable test outcomes.
// Consider isolating these tests or using a dedicated test environment.

// token_manager_test.go
package main

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

package utils

import (
	"os"
	"testing"
)

// TestReadURLs tests the readURLs function for various scenarios.
func TestReadURLs(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		wantURLs  []string
		wantError bool
	}{
		{
			name:      "single URL",
			content:   "https://example.com",
			wantURLs:  []string{"https://example.com"},
			wantError: false,
		},
		{
			name:      "multiple URLs",
			content:   "https://example.com\nhttps://example.org",
			wantURLs:  []string{"https://example.com", "https://example.org"},
			wantError: false,
		},
		{
			name:      "empty file",
			content:   "",
			wantURLs:  nil,
			wantError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a temporary file
			tmpFile, err := os.CreateTemp("", "urls")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name()) // Clean up

			// Write the content to the temp file
			if _, err := tmpFile.WriteString(tt.content); err != nil {
				t.Fatalf("Failed to write to temp file: %v", err)
			}
			if err := tmpFile.Close(); err != nil {
				t.Fatalf("Failed to close temp file: %v", err)
			}

			// Test readURLs
			gotURLs, err := ReadURLs(tmpFile.Name())
			if (err != nil) != tt.wantError {
				t.Errorf("readURLs() error = %v, wantErr %v", err, tt.wantError)
				return
			}
			if !stringsEqual(gotURLs, tt.wantURLs) {
				t.Errorf("readURLs() gotURLs = %v, want %v", gotURLs, tt.wantURLs)
			}
		})
	}
}

// stringsEqual checks if two slices of strings are equal.
func stringsEqual(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

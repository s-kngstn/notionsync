package main

import (
	"fmt"
	"net/url"
	"regexp"
)

// Assumes this function does something with the UUID extracted from the URL.
// The implementation would depend on what you want to do with the UUID.
func FetchDataBlockString(inputURL string) (string, error) {
	// Check if the input string is a valid URL
	parsedURL, err := url.Parse(inputURL)
	if err != nil || parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL")
	}

	// Proceed with UUID extraction from the valid URL
	extracted, uuid := extractUUIDFromURL(inputURL)
	if !extracted {
		return "", fmt.Errorf("no UUID found in URL")
	}

	// Placeholder for further processing with the UUID
	return uuid, nil
}

func extractUUIDFromURL(url string) (bool, string) {
	// Regex pattern that matches both standard and compact forms of UUIDs
	uuidPattern := `([0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[1-5][0-9a-fA-F]{3}-[89abAB][0-9a-fA-F]{3}-[0-9a-fA-F]{12})|([0-9a-fA-F]{32})`
	regex := regexp.MustCompile(uuidPattern)

	// Search for UUID in the URL
	matches := regex.FindStringSubmatch(url)
	if len(matches) > 0 {
		// Return the first match found
		for _, match := range matches {
			if match != "" {
				return true, match
			}
		}
	}

	// No UUID found
	return false, ""
}

package main

import (
	"fmt"
	"net/url"
	"regexp"
)

func FetchDataBlockString(inputURL string) (string, error) {
	// Check if the input string is a valid URL
	parsedURL, err := url.Parse(inputURL)

	if err != nil {
		return "", fmt.Errorf("error parsing URL: %w", err)
	}

	if parsedURL.Scheme == "" || parsedURL.Host == "" {
		return "", fmt.Errorf("invalid URL: missing scheme or host")
	}

	extracted, uuid := extractUUID(inputURL)
	if !extracted {
		return "", fmt.Errorf("no UUID found in URL")
	}

	return uuid, nil
}

func extractUUID(url string) (bool, string) {
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

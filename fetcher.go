package main

import (
	"fmt"
	"strings"
)

func FetchDataBlockString(url string) (string, error) {
	// Split the URL by "/"
	parts := strings.Split(url, "/")
	// The last part of the URL should be "Title-UUID"
	lastPart := parts[len(parts)-1]
	// Now, split the last part by "-" to isolate the UUID
	uuidParts := strings.Split(lastPart, "-")
	// The UUID is the last element after the split
	uuid := uuidParts[len(uuidParts)-1]

	// Check if UUID is not empty
	if uuid == "" {
		return "", fmt.Errorf("no UUID found in URL")
	}

	return uuid, nil
}

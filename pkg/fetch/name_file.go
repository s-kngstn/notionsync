package fetch

import (
	"fmt"
	"path"
	"strings"
)

// Extracts the desired substring from the URL.
func ExtractNameFromURL(url string) (string, error) {
	// Extract the last segment of the URL path
	lastSegment := path.Base(url)
	// Find the last dash before the UUID
	dashIndex := strings.LastIndex(lastSegment, "-")
	if dashIndex == -1 {
		return "", fmt.Errorf("invalid URL format")
	}
	// Extract and return the substring before the last dash, converted to lowercase
	return strings.ToLower(lastSegment[:dashIndex]), nil
}

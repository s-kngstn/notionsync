package fetch

import (
	"errors"
	"fmt"
	"net/url"
	"regexp"
)

type URLFetcher interface {
	CheckURL(inputURL string) (bool, error)
}
type BlockIDFetcher interface {
	GetBlockID(inputURL string) (string, error)
}

type DefaultBlockIDFetcher struct{}
type DefaultURLChecker struct{}

func (f DefaultURLChecker) CheckURL(inputURL string) (bool, error) {
	parsedURL, err := url.Parse(inputURL)
	if err != nil {
		// Return false and the parsing error if URL parsing fails.
		return false, err
	}

	// Define the correct domain. This is just an example; adjust it as needed.
	const correctDomain = "www.notion.so"
	if parsedURL.Hostname() != correctDomain {
		// Return false and an error if the domain does not match.
		return false, errors.New("URL domain needs to be from notion.so, skipping...")
	}

	// If we reach this point, the domain matches.
	return true, nil
}

func (f DefaultBlockIDFetcher) GetBlockID(inputURL string) (string, error) {
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

/**
* How should we handle images?
* - [ ] Re: Images - if we save locally, we need to handle the file type & size
* - [ ] Re: Images - maybe have it as a command line flag to save locally or not
 */

// This function could be used to download images from Notion pages.
// DownloadImage downloads an image from the given URL and saves it to the specified file path.
// func DownloadImage(imageUrl, filePath string) error {
// 	resp, err := http.Get(imageUrl)
// 	if err != nil {
// 		return err
// 	}
// 	defer resp.Body.Close()

// 	if resp.StatusCode != http.StatusOK {
// 		return err
// 	}

// 	file, err := os.Create(filePath)
// 	if err != nil {
// 		return err
// 	}
// 	defer file.Close()

// 	_, err = io.Copy(file, resp.Body)
// 	if err != nil {
// 		return err
// 	}

// 	return nil
// }

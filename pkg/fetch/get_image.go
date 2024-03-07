package fetch

import (
	"io"
	"net/http"
	"os"
)

/**
* How should we handle images?
* - [ ] Re: Images - if we save locally, we need to handle the file type & size
* - [ ] Re: Images - maybe have it as a command line flag to save locally or not
 */

// This function could be used to download images from Notion pages.
// DownloadImage downloads an image from the given URL and saves it to the specified file path.
func DownloadImage(imageUrl, filePath string) error {
	resp, err := http.Get(imageUrl)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return err
	}

	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return err
	}

	return nil
}

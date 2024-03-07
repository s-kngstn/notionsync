package format

import (
	"fmt"

	"github.com/s-kngstn/notionsync/api"
)

func ProcessBlocks(results *api.ResultsWrapper, outputPath string, pageName string, apiClient *api.NotionApiClient, bearerToken string) {
	// We need to handle if we have a child page block has children and if so, we need to call the API again
	// The ID of the child page is the UUID of the block
	for _, block := range results.Results {
		if block.Type == "child_page" && block.HasChildren {
			// Call the API with the ID of the block
			childResults, err := apiClient.GetNotionBlocks(block.ID, bearerToken)
			if err != nil {
				fmt.Println("Error calling API:", err)
				return
			}
			// The output path will be the same as the parent page, but with the title of the child page name
			childPageName := ToKebabCase(block.ChildPage.Title)
			childOutputPath := fmt.Sprintf("output/%s.md", childPageName)
			// Write the child blocks to a Markdown file with the title of the child pageName
			if err := WriteBlocksToMarkdown(childResults, childOutputPath, childPageName); err != nil {
				fmt.Println("Error writing blocks to Markdown:", err)
			}
		}
	}

	// Now process and write the results to a Markdown file
	if err := WriteBlocksToMarkdown(results, outputPath, pageName); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
	}
}

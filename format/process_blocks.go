package format

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/s-kngstn/notionsync/api"
)

func ProcessBlocks(uuid string, results *api.ResultsWrapper, outputPath string, pageName string, apiClient *api.NotionApiClient, bearerToken string, processedBlocks map[string]string) {
	// We need to handle if we have a child page block has children and if so, we need to call the API again
	// The ID of the child page is the UUID of the block
	for _, block := range results.Results {
		// Assume every block has a unique ID that can be used as a key in the map
		blockID := block.ID

		if _, processed := processedBlocks[blockID]; processed {
			// If the block has already been processed, skip it or handle as needed
			continue
		}

		if block.Type == "child_page" && block.HasChildren {
			// Call the API with the ID of the block
			childResults, err := apiClient.GetNotionBlocks(block.ID, bearerToken)
			if err != nil {
				fmt.Println("Error calling API:", err)
				return
			}
			// The output path will be the same as the parent page, but with the title of the child page name
			childPageName := strcase.ToKebab(block.ChildPage.Title)
			childOutputPath := fmt.Sprintf("output/%s.md", childPageName)

			// Before processing child blocks, mark this block as processed to avoid infinite recursion
			processedBlocks[blockID] = childOutputPath

			// Write the child blocks to a Markdown file with the title of the child pageName
			if err := WriteBlocksToMarkdown(childResults, childOutputPath, childPageName); err != nil {
				fmt.Println("Error writing blocks to Markdown:", err)
			}
		}
	}

	// Mark the parent block as processed using the original UUID and the output path
	processedBlocks[uuid] = outputPath
	// Now process and write the results to a Markdown file
	if err := WriteBlocksToMarkdown(results, outputPath, pageName); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
	}
}

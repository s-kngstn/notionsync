package format

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/s-kngstn/notionsync/api"
)

func ProcessBlocks(uuid string, results *api.ResultsWrapper, outputPath, pageName string, apiClient *api.NotionApiClient, bearerToken string, processedBlocks map[string]string) {
	for _, block := range results.Results {
		if _, processed := processedBlocks[block.ID]; processed {
			// If the block has already been processed, skip it or handle as needed
			continue
		}

		if block.Type == "child_page" && block.HasChildren {
			processChildBlocks(&block, apiClient, bearerToken, processedBlocks)
		} else if block.Type == "link_to_page" {
			print("LINK TO PAGE")
			// Handle link_to_page blocks similarly if they require recursive processing
			if _, processed := processedBlocks[block.LinkToPage.PageID]; processed {
				fmt.Println("LTP PROCCESSED BEFORE")
				// I need to link to the file that was created for that page
			} else {
				fmt.Println("LTP NOT YET PROCCESSED")
				// I need to processChildBlocks here
			}
		}

		// After processing children, mark the current block as processed
		processedBlocks[block.ID] = outputPath
	}

	// After all blocks are processed, write the current page's blocks to markdown
	if err := WriteBlocksToMarkdown(results, outputPath, pageName); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
	}
}

func processChildBlocks(parentBlock *api.Block, apiClient *api.NotionApiClient, bearerToken string, processedBlocks map[string]string) {
	childResults, err := apiClient.GetNotionChildBlocks(parentBlock.ID, bearerToken)
	if err != nil {
		fmt.Println("Error calling API for child blocks:", err)
		return
	}
	childPageName := strcase.ToKebab(parentBlock.ChildPage.Title)
	childOutputPath := fmt.Sprintf("output/%s.md", childPageName)

	// Before processing child blocks further, mark this parent block as processed to avoid infinite recursion
	processedBlocks[parentBlock.ID] = childOutputPath

	// Write the child blocks to a Markdown file
	if err := WriteBlocksToMarkdown(childResults, childOutputPath, childPageName); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
		return
	}

	// If child blocks have their own children, recursively process them too
	for _, childBlock := range childResults.Results {
		if childBlock.HasChildren {
			processChildBlocks(&childBlock, apiClient, bearerToken, processedBlocks)
		}
	}
}

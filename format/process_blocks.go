package format

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/s-kngstn/notionsync/api"
)

func ProcessBlocks(uuid string, results *api.ResultsWrapper, outputPath, pageName string, apiClient *api.NotionApiClient, bearerToken string, processedBlocks map[string]string) {
	for _, block := range results.Results {
		if _, processed := processedBlocks[block.ID]; processed {
			continue
		}

		switch block.Type {
		case "child_page":
			if block.HasChildren {
				processChildBlocks(&block, apiClient, bearerToken, processedBlocks)
			}
		case "link_to_page":
			title, err := apiClient.GetNotionBlockTitle(block.LinkToPage.PageID, bearerToken)
			if err != nil {
				fmt.Println("Error calling API for block title:", err)
				continue
			}

			if _, processed := processedBlocks[block.LinkToPage.PageID]; !processed {
				linkedResults, err := apiClient.GetNotionChildBlocks(block.LinkToPage.PageID, bearerToken)
				if err != nil {
					fmt.Println("Error calling API for LTP child blocks:", err)
					return
				}
				linkedPageName := strcase.ToKebab(title)
				linkedOutputPath := fmt.Sprintf("output/%s.md", linkedPageName)

				// Mark the block as processed to avoid infinite recursion
				processedBlocks[block.LinkToPage.PageID] = linkedOutputPath

				if err := WriteBlocksToMarkdown(linkedResults, linkedOutputPath, linkedPageName); err != nil {
					fmt.Println("Error writing blocks to Markdown:", err)
					return
				}
			} else {
				fmt.Println("LTP PROCESSED BEFORE")
				// Logic to link to the file that was created for that page goes here
			}
		}

		processedBlocks[block.ID] = outputPath
	}

	// After all blocks are processed, write the original searched URL page's blocks to markdown
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

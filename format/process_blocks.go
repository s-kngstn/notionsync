package format

import (
	"fmt"

	"github.com/iancoleman/strcase"
	"github.com/s-kngstn/notionsync/api"
)

func ProcessBlocks(uuid string, results *api.ResultsWrapper, outputPath, pageName string, apiClient api.NotionAPI, bearerToken string, processedBlocks map[string]string, outputDir string) {
	linkTitles := make(map[string]string)

	for _, block := range results.Results {
		if _, processed := processedBlocks[block.ID]; processed {
			continue
		}

		switch block.Type {
		case "child_page":
			processChildPageBlock(&block, apiClient, bearerToken, processedBlocks, linkTitles, outputDir)
		case "link_to_page":
			processLinkToPageBlock(&block, apiClient, bearerToken, processedBlocks, linkTitles, outputDir)
		}

		processedBlocks[block.ID] = outputPath
	}

	WriteBlocksToMarkdown(results, outputPath, pageName, linkTitles)
}

func processLinkToPageBlock(block *api.Block, apiClient api.NotionAPI, bearerToken string, processedBlocks, linkTitles map[string]string, outputDir string) {
	title, err := api.FetchBlockTitle(apiClient, block.LinkToPage.PageID, bearerToken)
	if err != nil {
		fmt.Println("Error fetching title:", err)
		return
	}

	linkTitles[block.LinkToPage.PageID] = title
	if _, processed := processedBlocks[block.LinkToPage.PageID]; !processed {
		linkedResults, err := api.FetchChildBlocks(apiClient, block.LinkToPage.PageID, bearerToken)
		if err != nil {
			fmt.Println("Error fetching child blocks:", err)
			return
		}
		linkedPageName := strcase.ToKebab(title)
		linkedOutputPath := fmt.Sprintf("%s/%s.md", outputDir, linkedPageName)
		processedBlocks[block.LinkToPage.PageID] = linkedOutputPath
		WriteBlocksToMarkdown(linkedResults, linkedOutputPath, linkedPageName, linkTitles)
	}
}

func processChildPageBlock(block *api.Block, apiClient api.NotionAPI, bearerToken string, processedBlocks, linkTitles map[string]string, outputDir string) {
	if !block.HasChildren {
		return
	}
	processChildBlocks(block, apiClient, bearerToken, processedBlocks, linkTitles, outputDir)
}

func processChildBlocks(parentBlock *api.Block, apiClient api.NotionAPI, bearerToken string, processedBlocks map[string]string, linkTitles map[string]string, outputDir string) {
	childResults, err := api.FetchChildBlocks(apiClient, parentBlock.ID, bearerToken)
	if err != nil {
		fmt.Println("Error calling API for child blocks:", err)
		return
	}
	childPageName := strcase.ToKebab(parentBlock.ChildPage.Title)
	childOutputPath := fmt.Sprintf("%s/%s.md", outputDir, childPageName)

	// Before processing child blocks further, mark this parent block as processed to avoid infinite recursion
	processedBlocks[parentBlock.ID] = childOutputPath

	if err := WriteBlocksToMarkdown(childResults, childOutputPath, childPageName, linkTitles); err != nil {
		fmt.Println("Error writing blocks to Markdown:", err)
		return
	}

	// If child blocks have their own children, recursively process them too
	for _, childBlock := range childResults.Results {
		if childBlock.HasChildren {
			processChildBlocks(&childBlock, apiClient, bearerToken, processedBlocks, linkTitles, outputDir)
		}
	}
}

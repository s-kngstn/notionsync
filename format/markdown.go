package format

import (
	"fmt"
	"os"

	"github.com/s-kngstn/notionsync/api"
)

func WriteBlocksToMarkdown(results *api.ResultsWrapper, outputPath string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating markdown file: %w", err)
	}
	defer file.Close()

	for _, block := range results.Results {
		fmt.Printf("Block ID: %s\n", block.ID)
		fmt.Printf("Block Type: %s\n", block.Type)

		var provider api.RichTextProvider
		var markdownPrefix string

		switch block.Type {
		case "heading_1":
			provider = block.Heading1
			markdownPrefix = "# "
		case "heading_2":
			provider = block.Heading2
			markdownPrefix = "## "
		case "heading_3":
			provider = block.Heading3
			markdownPrefix = "### "
		case "paragraph":
			provider = block.Paragraph
			markdownPrefix = "" // No prefix needed for paragraphs
		}

		if provider != nil {
			for _, rt := range provider.GetRichText() {
				formattedContent := markdownPrefix + rt.Text.Content + "\n"
				_, err := file.WriteString(formattedContent)
				if err != nil {
					return fmt.Errorf("error writing to markdown file: %w", err)
				}
				fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
			}
		}
	}

	// @TODO: Add a success message if needed
	fmt.Println("Markdown file successfully created.")

	return nil
}

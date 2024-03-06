package format

import (
	"fmt"
	"os"
	"strings"

	"github.com/s-kngstn/notionsync/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func applyAnnotationsToContent(rt api.RichText) string {
	formattedText := rt.Text.Content

	// Apply Markdown syntax for both bold and italic
	if rt.Annotations.Bold && rt.Annotations.Italic {
		formattedText = "***" + formattedText + "***"
	} else if rt.Annotations.Bold {
		// Apply Markdown syntax for bold
		formattedText = "**" + formattedText + "**"
	} else if rt.Annotations.Italic {
		// Apply Markdown syntax for italic
		formattedText = "*" + formattedText + "*"
	}

	// Future enhancements here for other annotations like strikethrough, underline, etc.

	// Apply Markdown syntax for links
	if rt.Text.Link != nil && rt.Text.Link.URL != nil {
		formattedText = "[" + formattedText + "](" + *rt.Text.Link.URL + ")"
	}

	return formattedText
}

// toTitleCase converts a string to title case.
func toTitleCase(input string) string {
	// Replace hyphens with spaces
	input = strings.ReplaceAll(input, "-", " ")

	// Create a title cased converter for the specified language
	caser := cases.Title(language.English)

	// Apply title casing
	return caser.String(input)
}

func WriteBlocksToMarkdown(results *api.ResultsWrapper, outputPath string, pageName string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating markdown file: %w", err)
	}
	defer file.Close()

	// Write the page name as the title
	pageTitle := toTitleCase(pageName)
	_, err = file.WriteString(fmt.Sprintf("# %s\n\n", pageTitle))
	if err != nil {
		return fmt.Errorf("error writing to markdown file: %w", err)
	}

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
				formattedContent := markdownPrefix + applyAnnotationsToContent(rt) + "\n"
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

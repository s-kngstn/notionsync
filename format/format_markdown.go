package format

import (
	"fmt"
	"os"
	"strings"
	"unicode"

	"github.com/s-kngstn/notionsync/api"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

func applyAnnotationsToContent(rt api.RichText) string {
	formattedText := rt.Text.Content

	if rt.Annotations.Bold && rt.Annotations.Italic {
		formattedText = "***" + formattedText + "***"
	} else if rt.Annotations.Bold {
		formattedText = "**" + formattedText + "**"
	} else if rt.Annotations.Italic {
		formattedText = "*" + formattedText + "*"
	}

	if rt.Annotations.Strikethrough {
		formattedText = "~~" + formattedText + "~~"
	}

	if rt.Annotations.Code {
		formattedText = "`" + formattedText + "`"
	}

	// Syntax for underline is not supported in markdown

	// Syntax for links
	if rt.Text.Link != nil && rt.Text.Link.URL != nil {
		formattedText = "[" + formattedText + "](" + *rt.Text.Link.URL + ")"
	}

	return formattedText
}

func toTitleCase(input string) string {
	input = strings.ReplaceAll(input, "-", " ")
	caser := cases.Title(language.English)
	return caser.String(input)
}

func ToKebabCase(s string) string {
	var result strings.Builder
	for i, r := range s {
		// Convert uppercase letters to lowercase, prefixing them with a dash if it's not the first character
		// and the previous character isn't a dash, ensuring we don't start with a dash.
		if unicode.IsUpper(r) {
			if i > 0 && result.Len() > 0 && result.String()[result.Len()-1] != '-' {
				result.WriteRune('-')
			}
			result.WriteRune(unicode.ToLower(r))
		} else if unicode.IsDigit(r) || unicode.IsLower(r) {
			// Directly append digits and lowercase letters
			result.WriteRune(r)
		} else if unicode.IsSpace(r) && i > 0 && result.String()[result.Len()-1] != '-' {
			// Convert spaces to dashes if the last character isn't already a dash.
			result.WriteRune('-')
		}
		// Skip non-alphanumeric, non-space characters
	}
	return result.String()
}

func WriteBlocksToMarkdown(results *api.ResultsWrapper, outputPath string, pageName string) error {
	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("error creating markdown file: %w", err)
	}
	defer file.Close()

	pageTitle := toTitleCase(pageName)
	_, err = file.WriteString(fmt.Sprintf("# %s\n\n", pageTitle))
	if err != nil {
		return fmt.Errorf("error writing to markdown file: %w", err)
	}

	listItemNumber := 1
	processingNumberedList := false
	for _, block := range results.Results {
		fmt.Printf("Block ID: %s\n", block.ID)
		fmt.Printf("Block Type: %s\n", block.Type)

		var provider api.RichTextProvider
		var markdownPrefix string
		var formattedContent string

		switch block.Type {
		case "heading_1":
			provider = block.Heading1
			markdownPrefix = "# "
			processingNumberedList = false
		case "heading_2":
			provider = block.Heading2
			markdownPrefix = "## "
			processingNumberedList = false
		case "heading_3":
			provider = block.Heading3
			markdownPrefix = "### "
			processingNumberedList = false
		case "paragraph":
			provider = block.Paragraph
			markdownPrefix = "" // No prefix needed for paragraphs
			processingNumberedList = false
		case "code":
			provider = block.Code
			markdownPrefix = "```" + block.Code.Language + "\n"
			processingNumberedList = false
		case "child_page":
			// [Link Text](filename.md)
			markdownPrefix = fmt.Sprintf("- [%s](%s.md)", block.ChildPage.Title, ToKebabCase(block.ChildPage.Title))
		case "bookmark":
			markdownPrefix = "- [" + block.Bookmark.URL + "]"
			processingNumberedList = false
		case "to_do":
			provider = block.Todo
			if block.Todo.Checked {
				markdownPrefix = "- [x] "
			} else {
				markdownPrefix = "- [ ] "
			}
			processingNumberedList = false
		case "bulleted_list_item":
			provider = block.Bulleted
			markdownPrefix = "- "
			processingNumberedList = false
		case "numbered_list_item":
			if !processingNumberedList {
				listItemNumber = 1 // Start numbering from 1 for a new list
				processingNumberedList = true
			} else {
				listItemNumber++ // Increment if we're continuing a list
			}
			provider = block.Numbered
			markdownPrefix = fmt.Sprintf("%d. ", listItemNumber)
			// No reset here since we might be continuing the list
		}

		if provider != nil {
			for _, rt := range provider.GetRichText() {
				if block.Type == "code" {
					formattedContent = markdownPrefix + rt.Text.Content + "\n ```\n"
				} else {
					formattedContent = markdownPrefix + applyAnnotationsToContent(rt) + "\n"
				}
				_, err := file.WriteString(formattedContent)
				if err != nil {
					return fmt.Errorf("error writing to markdown file: %w", err)
				}
				fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
			}
		}

		// handle the case where the block is a Bookmark
		if block.Bookmark != nil {
			_, err := file.WriteString(fmt.Sprintf("%s\n", markdownPrefix))
			if err != nil {
				return fmt.Errorf("error writing to markdown file: %w", err)
			}
		}

		// handle the case where the block is a child pageName
		if block.Type == "child_page" {
			_, err := file.WriteString(fmt.Sprintf("%s\n", markdownPrefix))
			if err != nil {
				return fmt.Errorf("error writing to markdown file: %w", err)
			}
		}
	}

	// @TODO: Add a success message if needed
	fmt.Println("Markdown file successfully created.")

	return nil
}

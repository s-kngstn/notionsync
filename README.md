# notionsync

NotionSync is a command-line tool designed to synchronize your [Notion](https://www.notion.so/) notes to your local machine in markdown format. It supports a variety of markdown features, making your Notion content easily accessible and editable on your local device.

Either set the `NOTION_API_KEY` environment variable or pass it as an argument.

## Features

NotionSync currently supports syncing the following markdown formats:

- Text (bold, italic, strikethrough, code)
- Headings (H1, H2, H3)
- Lists (Bulleted, Numbered, To-do)
- Quote, Code block, Divider
- Links (Bookmark, URL, Page)
- Future support planned for Images, Videos, and Tables

## Getting Started

### Prerequisites

- Go 1.15 or higher
- Notion API key. You can obtain one by creating an integration at [Notion Integrations](https://www.notion.so/my-integrations).

## Usage

notionsync offers several flags to customize its operation:

- token: Notion API bearer token. If not provided, the tool will attempt to use the environment variable NOTION_API_KEY.
- file: Path to the file containing Notion page URLs to sync. If not provided, the tool will prompt for a single URL input.

Example Commands
Sync using an API token passed as a flag:
```bash
./notionsync -token="your_notion_api_key_here" -file="path/to/your/url_file.txt"
```

Sync with the API key set as an environment variable (or in a .env file for development):
```bash
./notionsync -file="path/to/your/url_file.txt"
```

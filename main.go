package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// Define Go structs to map the JSON structure, focusing on parts of interest
type Block struct {
	Object    string    `json:"object"`
	ID        string    `json:"id"`
	Paragraph Paragraph `json:"paragraph"`
}

type Paragraph struct {
	RichText []RichText `json:"rich_text"`
	Color    string     `json:"color"`
}

type RichText struct {
	Type string `json:"type"`
	Text Text   `json:"text"`
}

type Text struct {
	Content string `json:"content"`
	// Include other fields if needed
}

type ResultsWrapper struct {
	Results []Block `json:"results"`
}

// Function to make the API call with Authorization header
func callAPI(customID, bearerToken string) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", customID)
	// url := fmt.Sprintf("http://something.com/%s", customID)

	// Creating the request
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	bearerToken = "secret_hVDPuHdW5ec7WzM2WicFHNCT7dWy8F5mOE9MMIY2PjK"
	// Adding the Authorization header
	req.Header.Add("Authorization", "Bearer "+bearerToken)

	// Adding the Authorization header
	req.Header.Add("Notion-Version", "2022-06-28")

	// Initializing the HTTP client
	client := &http.Client{}

	// Sending the request
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Reading the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Unmarshal the JSON data into the ResultsWrapper struct
	var results ResultsWrapper
	if err := json.Unmarshal(body, &results); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	// Iterate through results to extract and print rich text content and type
	for _, block := range results.Results {
		fmt.Printf("Block ID: %s\n", block.ID)
		for _, rt := range block.Paragraph.RichText {
			fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
		}
	}
	// // Printing the response body to the console
	// fmt.Println("Response from API:", string(body))
}

func main() {
	// Define a command-line flag
	// var url string
	// flag.StringVar(&url, "url", "", "URL to extract the UUID from")
	// flag.Parse()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Please enter the URL: ")
	url, _ := reader.ReadString('\n')
	url = strings.TrimSpace(url)

	// @todo: the url text before the UUID is the h1 title of the page you want to sync
	// @todo: this needs tests
	// Check if the URL was provided
	if url == "" {
		fmt.Println("URL is required")
		return
	}

	uuid, err := FetchDataBlockString(url)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
	} else {
		fmt.Printf("Extracted UUID: %s\n", uuid)
	}

	// Call the API
	callAPI(uuid, "")
}

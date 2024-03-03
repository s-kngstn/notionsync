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

// type Paragraph struct {
// 	RichText []RichText `json:"rich_text"`
// 	Color    string     `json:"color"`
// }

// Assume the API returns JSON error responses
type APIErrorResponse struct {
	Object  string `json:"object,omitempty"`
	Status  int    `json:"status,omitempty"`
	Code    string `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Your existing structs for successful response parsing
type ResultsWrapper struct {
	// Assuming the structure of your successful response
	Results []Block `json:"results"`
}

type Block struct {
	ID string `json:"id"`
	// Assuming a Paragraph field in your Block struct
	Paragraph struct {
		RichText []RichText `json:"rich_text"`
	} `json:"paragraph"`
}

type RichText struct {
	Type string `json:"type"`
	Text struct {
		Content string `json:"content"`
	} `json:"text"`
}

// Function to make the API call with Authorization header
func callAPI(customID, bearerToken string) {
	url := fmt.Sprintf("https://api.notion.com/v1/blocks/%s/children?page_size=100", customID)
	println("URL accepted")

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	bearerToken = "secret_hVDPuHdW5ec7WzM2WicFHNCT7dWy8F5mOE9MMIY2PjK"
	req.Header.Add("Authorization", "Bearer "+bearerToken)
	req.Header.Add("Notion-Version", "2022-06-28")

	client := &http.Client{}

	resp, err := client.Do(req)
	println("Request sent")
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == 400 {
		// Handle 400 Bad Request error
		var apiError APIErrorResponse
		err := json.NewDecoder(resp.Body).Decode(&apiError)
		if err != nil {
			fmt.Println("Error parsing API error response:", err)
			return
		}
		fmt.Printf("API Error: %s - %s\n", apiError.Code, apiError.Message)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	if resp.StatusCode != 200 {
		// Handle other non-200 and non-400 statuses
		fmt.Printf("Unexpected status code: %d\n", resp.StatusCode)
		return
	}

	var results ResultsWrapper
	if err := json.Unmarshal(body, &results); err != nil {
		fmt.Println("Error parsing JSON:", err)
		return
	}

	for _, block := range results.Results {
		fmt.Printf("Block ID: %s\n", block.ID)
		for _, rt := range block.Paragraph.RichText {
			fmt.Printf("Rich Text Type: %s, Content: %s\n", rt.Type, rt.Text.Content)
		}
	}
}

func main() {
	reader := bufio.NewReader(os.Stdin)

	var uuid string
	var err error

	for {
		fmt.Print("Please enter the URL: ")
		url, _ := reader.ReadString('\n')
		url = strings.TrimSpace(url)

		// Check if the URL was provided
		if url == "" {
			fmt.Println("URL is required, please try again.")
			continue // Skip the rest of the loop and prompt again
		}

		uuid, err = FetchDataBlockString(url)
		if err != nil {
			fmt.Printf("Error: %v. Please try again.\n", err)
			continue // If an error occurs (e.g., no UUID found), prompt for the URL again
		}

		fmt.Printf("Extracted UUID: %s\n", uuid)
		break // Exit the loop if a valid UUID is found
	}

	// Call the API with the extracted UUID
	callAPI(uuid, "")
}

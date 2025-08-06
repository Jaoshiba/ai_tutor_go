package services

import (
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
)

// SerpAPIResponse represents the structure of the JSON response from SerpAPI
// This is a simplified structure; you might need to expand it based on your needs.

// SearchAcademicDocuments uses SerpAPI to search for academic documents based on a module name and its description.
// It takes the module name, description, and your SerpAPI key as input.
func SearchDocuments(moduleName string, description string) (string, error) {
	// Base URL for SerpAPI
	baseURL := "https://serpapi.com/search"

	serpAPIKey := os.Getenv("SERPAPI_KEY")

	// Combine moduleName and description to form a comprehensive search query
	// This approach provides a more detailed context for the search.
	searchQuery := fmt.Sprintf("%s %s", moduleName, description)

	// Construct the URL with query parameters
	// We're using 'google_scholar' engine for academic documents.
	// You can change 'engine' to 'google' if you prefer general web search.
	params := url.Values{}
	params.Add("q", searchQuery)
	params.Add("engine", "google_scholar") // Using Google Scholar for academic search
	params.Add("api_key", serpAPIKey)
	params.Add("hl", "th") // Host language for the search results (Thai)
	params.Add("gl", "th") // Geographic location for the search (Thailand)

	fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

	// Create a new HTTP GET request
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถสร้างคำขอ HTTP ได้: %w", err)
	}

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถส่งคำขอ HTTP ได้: %w", err)
	}
	defer resp.Body.Close() // Ensure the response body is closed

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถอ่านเนื้อหาการตอบกลับได้: %w", err)
	}

	// Check if the HTTP status code indicates an error
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("SerpAPI ตอบกลับด้วยสถานะผิดพลาด %d: %s", resp.StatusCode, string(body))
	}

	// Unmarshal the JSON response into our struct
	var serpAPIResponse entities.SerpAPIResponse
	err = json.Unmarshal(body, &serpAPIResponse)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถแยกวิเคราะห์ JSON ได้: %w", err)
	}

	fmt.Println("SerpAPI Response:", serpAPIResponse)

	return "", nil
}

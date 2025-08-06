package services

import (
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
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
	body, err := io.ReadAll(resp.Body)
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

	for i, result := range serpAPIResponse.OrganicResults {
		if i >= 2 {
			break
		}
		documentLink := result.Link
		doucmentTitle := result.Title

		fmt.Println("Getting file from URL:", documentLink)
		docPath, err := GetFileFromUrl(doucmentTitle, documentLink)
		if err != nil {
			return "", fmt.Errorf("ไม่สามารถดาวน์โหลดไฟล์จาก URL ได้: %w", err)
		}
		file, err := os.ReadFile(docPath)
		fmt.Println("File content:", string(file))
		fmt.Printf("Finished Get file from url: %s\n", docPath)

	}

	fmt.Println("SerpAPI Response:", serpAPIResponse)

	return "", nil
}

func GetFileFromUrl(fileTitle string, fileUrl string) (string, error) {

	downloadDir := "fileDocs"
	fullPath := filepath.Join(downloadDir, fileTitle)
	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถสร้างไฟล์: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถดึงไฟล์จาก URL ได้: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("ไม่สามารถเข้าถึงไฟล์ได้: สถานะ %d", resp.StatusCode)
	}

	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("ไม่สามารถคัดลอกไฟล์ได้: %w", err)
	}

	return fullPath, nil
}

package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/gofiber/fiber/v2"
	
)

var allowedExts = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
	".html": true,
}
var allowedCTs = map[string]string{
	"application/html": ".html",
	"application/pdf": ".pdf",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
}

func SearchDocuments( courseName, courseDescription, moduleName, moduleDescription string, ctx *fiber.Ctx) (string, error) {
	geminiService := NewGeminiService()

	baseURL := "https://serpapi.com/search"
	serpAPIKey := os.Getenv("SERPAPI_KEY")
	if serpAPIKey == "" {
		return "", fmt.Errorf("missing SERPAPI_KEY")
	}

	maxAttempts := 3
	var generatedKeywords string

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		searchPrompt := buildKeywordPrompt(moduleName, moduleDescription, courseName, courseDescription, attempt)
		fmt.Printf("[SearchDocuments] Attempt %d: generating keywords...\n", attempt)

		kws, err := geminiService.GenerateContentFromPrompt(context.Background(), searchPrompt)
		if err != nil {
			return "", fmt.Errorf("error generating search keywords (attempt %d): %w", attempt, err)
		}
		generatedKeywords = strings.TrimSpace(kws)
		fmt.Println("Generated Search Query:", generatedKeywords)

		// generatedKeywords := fmt.Sprintf("สอน %s %s ใน %s -site:youtube.com -site:facebook.com", moduleName, moduleDescription, courseName)


		params := url.Values{}
		params.Add("q", generatedKeywords)
		params.Add("engine", "google")
		params.Add("api_key", serpAPIKey)
		params.Add("hl", "th")
		params.Add("gl", "th")

		fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		fmt.Printf("[SearchDocuments] Attempt %d: request SerpAPI...\n", attempt)
		body, status, err := doSerpAPISearch(fullURL)
		if err != nil {
			return "", fmt.Errorf("could not make SerpAPI request (attempt %d): %w", attempt, err)
		}
		if status != http.StatusOK {
			return "", fmt.Errorf("SerpAPI returned status code %d (attempt %d): %s", status, attempt, string(body))
		}

		var serpAPIResponse entities.SerpAPIResponse
		if err := json.Unmarshal(body, &serpAPIResponse); err != nil {
			return "", fmt.Errorf("could not parse JSON (attempt %d): %w", attempt, err)
		}

		fmt.Println("SerpAPI Response:", serpAPIResponse.OrganicResults)

		for _, result := range serpAPIResponse.OrganicResults {
			fmt.Println("Title:", result.Title)
			fmt.Println("Link:", result.Link)
		}

		if len(serpAPIResponse.OrganicResults) > 0 {
			fmt.Printf("[SearchDocuments] Attempt %d: got %d results\n", attempt, len(serpAPIResponse.OrganicResults))

			for _, result := range serpAPIResponse.OrganicResults {
				documentLink := result.Link
				fmt.Println("Processing result:", documentLink)

				fileExt := strings.ToLower(filepath.Ext(documentLink))
				if fileExt == "" {
					fileExt = ".html"
				}

				if !allowedExts[fileExt] {
					fmt.Printf("Skipping file from URL: %s - not a supported document type\n", documentLink)
					continue
				}

				documentTitle := sanitizeFilename(result.Title)

				fmt.Println("\nGetting file from URL: ", documentLink)
				docPath, err := GetFileFromUrl(documentTitle, documentLink)
				if err != nil {
					continue
				}
				fmt.Printf("Finished Get file from url: %s\n", docPath)

				content, err := ReadFileData(docPath, ctx)
				if err != nil {
					fmt.Println("Error reading file data:", err)
					continue
				}

				fmt.Println("\nFile content:", content)
				return content, err
			}
			return "", err
		}
		fmt.Printf("[SearchDocuments] Attempt %d: empty results. Regenerating keywords...\n", attempt)
	}

	return "", fmt.Errorf("no results found from SerpAPI after %d attempts", maxAttempts)
}

func rateSerpLink() {


}

func doSerpAPISearch(fullURL string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("could not create HTTP request: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("could not send HTTP request: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("could not read response body: %w", err)
	}
	return body, resp.StatusCode, nil
}

func buildKeywordPrompt(moduleName, moduleDescription, courseName, courseDescription string, attempt int) string {
	var retryHint string

	switch attempt {
	case 1:
		retryHint = `
- Focus on bilingual keywords (Thai/English) and HTML pages.
- Use trusted websites like Medium, W3Schools, StackOverflow, Wikipedia, .ac.th, .edu, .gov, .org.
- Do not limit by PDF/DOC files.
`
	case 2:
		retryHint = `
- Broaden keywords with synonyms or related topics.
- Use trusted general websites: Medium, W3Schools, StackOverflow, reputable blogs, ResearchGate.
- Focus on HTML and avoid downloadable files.
`
	default:
		retryHint = `
- Use umbrella/general keywords for broader results.
- Select trusted and HTML-based online sources.
- Remove site: or filetype filters if necessary.
`
	}

	return fmt.Sprintf(`
You are an academic research assistant.
Your task is to create effective and broad Google search keywords 
to find **HTML pages** relevant to the following topic, which you will later scrape for content.

Course Title: %s
Course Description: %s
Module Title: %s
Module Description: %s

Guidelines for keyword generation:
1. Keywords must be broad enough to get a variety of relevant HTML results.
2. Include both Thai and English terms for the topic.
3. Prefer trusted websites, including:
   - Academic sites: .ac.th, .edu, .gov, .org
   - General trusted sites: Medium, W3Schools, StackOverflow, well-known blogs
4. Combine keywords using OR to expand coverage; use AND only when necessary.
5. Do NOT include filetype restrictions (no PDF/DOC limits).
6. Return only the final search query without explanation.

Additional retry hint for this attempt:
%s
`,courseName, courseDescription, moduleName, moduleDescription, retryHint)
}

func isAllowedContentType(ct string) (ok bool, ext string) {
	ct = strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	ext, ok = allowedCTs[ct]
	return
}



func GetFileFromUrl(fileTitle string, fileUrl string) (string, error) {
	downloadDir := "fileDocs"
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return "", fmt.Errorf("failed to create destination folder: %w", err)
	}

	ext := filepath.Ext(strings.Split(strings.Split(fileUrl, "?")[0], "#")[0])
	if ext == "" {
		ext = ".bin"
	}
	fullPath := filepath.Join(downloadDir, fileTitle+ext)

	out, err := os.Create(fullPath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer out.Close()

	resp, err := http.Get(fileUrl)
	if err != nil {
		return "", fmt.Errorf("failed to fetch file from URL: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("could not access file: status %d", resp.StatusCode)
	}

	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("failed to copy file: %w", err)
	}

	

	return fullPath, nil
}

func GetHtmlElement(link string, ctx *fiber.Ctx) error {
	fmt.Println("Getting HTML content from link:", link)

	browserlessAPIKey := os.Getenv("BROWSERLESS_API_KEY")
	if browserlessAPIKey == "" {
		fmt.Println("Error: Missing BROWSERLESS_API_KEY")
		return ctx.Status(fiber.StatusInternalServerError).SendString("Missing BROWSERLESS_API_KEY")
	}

	payload := map[string]string{"url": link}
	jsonPayload, err := json.Marshal(payload)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to marshal JSON payload: %v", err))
	}

	browserlessURL := fmt.Sprintf("https://production-sfo.browserless.io/content?token=%s", browserlessAPIKey)
	req, err := http.NewRequest("POST", browserlessURL, bytes.NewBuffer(jsonPayload))
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to create request: %v", err))
	}
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error making POST request: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to post content to browserless.io: %v", err))
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("browserless.io API error: %s", string(body))
		return ctx.Status(resp.StatusCode).SendString(fmt.Sprintf("browserless.io API error: %s", string(body)))
	}

	fmt.Println("Response status code:", resp.StatusCode)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Failed to read response body: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to read response body: %v", err))
	}

	// fmt.Println("Response body string: ", string(body))
	fmt.Println("HTML content retrieved successfully.")

	text, err := extractContentsFromHTML(string(body))
	if err != nil {
		fmt.Printf("Failed to extract contents from HTML: %v", err)
		return ctx.Status(fiber.StatusInternalServerError).SendString(fmt.Sprintf("Failed to extract contents from HTML: %v", err))
	}

	fmt.Println("Extracted text content:", text)

	return ctx.SendString(string(body))
}

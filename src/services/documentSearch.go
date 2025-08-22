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

	// "regexp"
	"strings"

	"github.com/gofiber/fiber/v2"
)

var allowedExts = map[string]bool{
	".pdf":  true,
	".doc":  true,
	".docx": true,
}
var allowedCTs = map[string]string{
	"application/pdf":    ".pdf",
	"application/msword": ".doc",
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document": ".docx",
}

// ==== ปรับฟังก์ชันหลัก ให้มี retry ====

func SearchDocuments(courseName string, courseDescription string, moduleName string, description string, ctx *fiber.Ctx) (string, error) {
	geminiService := NewGeminiService()

	baseURL := "https://serpapi.com/search"
	serpAPIKey := os.Getenv("SERPAPI_KEY")
	if serpAPIKey == "" {
		return "", fmt.Errorf("missing SERPAPI_KEY")
	}

	maxAttempts := 3

	var generatedKeywords string

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		// ----- 1) สร้างคีย์เวิร์ดใหม่ทุกครั้งที่พยายาม -----
		searchPrompt := buildKeywordPrompt(courseName, courseDescription, moduleName, description, attempt)

		fmt.Printf("[SearchDocuments] Attempt %d: generating keywords...\n", attempt)
		kws, err := geminiService.GenerateContentFromPrompt(context.Background(), searchPrompt)
		if err != nil {
			// ถ้า Gemini ล้มเหลว ให้หยุดเลย (เพราะไม่มีคีย์เวิร์ดไปค้น)
			return "", fmt.Errorf("error generating search keywords (attempt %d): %w", attempt, err)
		}
		generatedKeywords = strings.TrimSpace(kws)
		fmt.Println("Generated Search Query:", generatedKeywords)

		keyword := fmt.Sprintf("สอน %s %s ใน %s ", moduleName, description, courseName)

		// ----- 2) ยิง SerpAPI -----
		params := url.Values{}
		params.Add("q", keyword)
		params.Add("engine", "google") // โฟกัสฝั่งวิชาการก่อน
		params.Add("api_key", serpAPIKey)
		params.Add("hl", "th")
		params.Add("gl", "th")

		fullURL := fmt.Sprintf("%s?%s", baseURL, params.Encode())

		fmt.Printf("[SearchDocuments] Attempt %d: request SerpAPI...\n", attempt)
		body, status, err := doSerpAPISearch(fullURL)
		if err != nil {
			return "", fmt.Errorf("ไม่สามารถส่งคำขอ SerpAPI ได้ (attempt %d): %w", attempt, err)
		}
		if status != http.StatusOK {
			return "", fmt.Errorf("SerpAPI ตอบกลับด้วยสถานะผิดพลาด %d (attempt %d): %s", status, attempt, string(body))
		}

		var serpAPIResponse entities.SerpAPIResponse
		if err := json.Unmarshal(body, &serpAPIResponse); err != nil {
			return "", fmt.Errorf("ไม่สามารถแยกวิเคราะห์ JSON ได้ (attempt %d): %w", attempt, err)
		}

		fmt.Println("SerpAPI Response:", serpAPIResponse.OrganicResults)

		for _, result := range serpAPIResponse.OrganicResults {
			fmt.Println("Title:", result.Title)
			fmt.Println("Link:", result.Link)

		}

		GetHtmlElement(serpAPIResponse.OrganicResults[0].Link, ctx)

		fmt.Println("Processing SerpAPI results...")

		return "", err

		fmt.Println("after return")

		// ----- 3) ถ้าได้ผลลัพธ์ → ไปประมวลผลไฟล์
		if len(serpAPIResponse.OrganicResults) > 0 {
			fmt.Printf("[SearchDocuments] Attempt %d: got %d results\n", attempt, len(serpAPIResponse.OrganicResults))

			for i, result := range serpAPIResponse.OrganicResults {
				if i >= 2 {
					break
				}
				documentLink := result.Link
				fmt.Println("Processing result:", documentLink)

				fileExt := strings.ToLower(filepath.Ext(documentLink))

				if !allowedExts[fileExt] {
					fmt.Printf("Skipping file from URL: %s - not a supported document type\n", documentLink)
					continue // ข้ามไปผลลัพธ์ถัดไป
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
				// fmt.Println("\nFile content:", content)
				return content, err
			}

			// fmt.Println("SerpAPI Response:", serpAPIResponse)
			return "", err
		}

		// ----- 4) ถ้า "ว่าง" → วนเริ่มใหม่ (กลับไปสร้างคีย์เวิร์ดใหม่) -----
		fmt.Printf("[SearchDocuments] Attempt %d: empty results. Regenerating keywords...\n", attempt)
	}

	// ครบทุกความพยายามแล้วก็ยังว่าง
	return "", fmt.Errorf("ไม่พบผลลัพธ์จาก SerpAPI หลังลองใหม่ %d รอบ", maxAttempts)
}

func rateSerpLink() {

}

// แยกยิง HTTP ให้สั้นลง
func doSerpAPISearch(fullURL string) ([]byte, int, error) {
	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, 0, fmt.Errorf("ไม่สามารถสร้างคำขอ HTTP ได้: %w", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, fmt.Errorf("ไม่สามารถส่งคำขอ HTTP ได้: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("ไม่สามารถอ่านเนื้อหาการตอบกลับได้: %w", err)
	}
	return body, resp.StatusCode, nil
}

// สร้างพรอมป์สำหรับรอบต่างๆ: รอบหลังๆ จะ “ขยาย” เงื่อนไขเพื่อค้นให้กว้างขึ้น
func buildKeywordPrompt(courseName, courseDescription, moduleName, moduleDescription string, attempt int) string {
	var retryHint string
	var filetypeFilter string

	filetypeFilter = "filetype:pdf OR filetype:doc OR filetype:docx"

	switch attempt {
	case 1:
		retryHint = `
- ให้โฟกัสคำหลักสองภาษา (ไทย/อังกฤษ) และคงความกว้างของผลลัพธ์
- พยายามหลีกเลี่ยงคำเฉพาะเกินไป
`
	case 2:
		retryHint = `
- ขยายคำค้นด้วยคำพ้อง/หัวข้อข้างเคียง และอย่าจำกัดด้วย filetype
- ผสม site:ac.th OR site:edu OR site:researchgate.net เฉพาะบางส่วน ไม่ต้องใส่ทุกตัว
`
	default:
		retryHint = `
- เน้นคำทั่วไปที่เป็นคำ umbrella (เช่น formulation, methodology, review, guideline)
- ตัด site: ออกถ้าจำเป็น เพื่อให้ได้ผลลัพธ์กว้างขึ้น
`
	}

	// รวมตัวกรอง filetype เข้าไปในคำสั่ง Guidelines
	return fmt.Sprintf(`
You are an academic research assistant.
Your task is to create effective and broad Google Scholar search keywords 
to find academic documents, PDFs, or research papers relevant to the following topic.

Course Title: %s
Course Description: %s
Module Title: %s
Module Description: %s

Guidelines for keyword generation:
1. Keywords must be broad enough to get a variety of relevant results, not overly restrictive.
2. Include both Thai and English terms for the topic.
3. Use academic source filters such as: site:ac.th OR site:edu OR site:researchgate.net — but they do not have to match all at once.
4. **Prefer file formats by adding: (%s)** (optional if it limits too much).
5. Combine keywords using OR to expand coverage; use AND only when necessary.
6. Return only the final search query without explanation.

Additional retry hint for this attempt:
%s
`, courseName, courseDescription, moduleName, moduleDescription, filetypeFilter, retryHint)
}

// กันชื่อไฟล์ให้ปลอดภัย
// func sanitizeFilename(name string) string {
// 	name = strings.TrimSpace(name)
// 	// แทนที่อักขระต้องห้ามด้วยขีดกลาง
// 	illegal := regexp.MustCompile(`[<>:"/\\|?*\x00-\x1F]`)
// 	name = illegal.ReplaceAllString(name, "-")
// 	// ย่อให้ไม่ยาวเกินไป
// 	if len(name) > 120 {
// 		name = name[:120]
// 	}
// 	// กันไม่ให้ชื่อว่าง
// 	if name == "" {
// 		name = "document"
// 	}
// 	return name
// }

func isAllowedContentType(ct string) (ok bool, ext string) {
	// ตัดพารามิเตอร์ เช่น "; charset=binary"
	ct = strings.ToLower(strings.TrimSpace(strings.Split(ct, ";")[0]))
	ext, ok = allowedCTs[ct]
	return
}

func GetFileFromUrl(fileTitle string, fileUrl string) (string, error) {
	downloadDir := "fileDocs"
	if err := os.MkdirAll(downloadDir, 0o755); err != nil {
		return "", fmt.Errorf("ไม่สามารถสร้างโฟลเดอร์ปลายทาง: %w", err)
	}

	// พยายามดึงนามสกุลจาก URL (เช่น .pdf)
	ext := filepath.Ext(strings.Split(strings.Split(fileUrl, "?")[0], "#")[0])
	if ext == "" {
		ext = ".bin"
	}
	fullPath := filepath.Join(downloadDir, fileTitle+ext)

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

	if _, err = io.Copy(out, resp.Body); err != nil {
		return "", fmt.Errorf("ไม่สามารถคัดลอกไฟล์ได้: %w", err)
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

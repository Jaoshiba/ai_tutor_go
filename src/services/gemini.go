//service/llm

package services

import (
	"context"
	"fmt"
	"os"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"
)

// IGeminiService defines the interface for Gemini AI operations.
type IGeminiService interface {
	GenerateContentFromPrompt(ctx context.Context, prompt string) (string, error)

	// คุณอาจเพิ่มเมธอดอื่น ๆ ที่เกี่ยวข้องกับการเรียก Gemini API ในอนาคต
	// เช่น GenerateCourseStructure(ctx context.Context, title, description, fileContent string) (entities.CourseGeminiResponse, error)
}

// GeminiService is a concrete implementation of IGeminiService.
type GeminiService struct {
	apiKey string
	
}

// NewGeminiService creates a new instance of GeminiService.
func NewGeminiService() *GeminiService {
	return &GeminiService{
		apiKey: os.Getenv("GEMINI_API_KEY"),
	}
}

// GenerateContentFromPrompt calls the Gemini API with a given prompt and returns the generated text.
func (gs *GeminiService) GenerateContentFromPrompt(ctx context.Context, prompt string) (string, error) {
	if gs.apiKey == "" {
		return "", fmt.Errorf("GEMINI_API_KEY ไม่ได้ถูกตั้งค่าใน environment variables")
	}
	

	client, err := genai.NewClient(ctx, option.WithAPIKey(gs.apiKey))
	if err != nil {
		return "", fmt.Errorf("ข้อผิดพลาดในการสร้าง Gemini client: %w", err)
	}
	defer client.Close()

	fmt.Println("กำลังเรียก Gemini API เพื่อสร้างเนื้อหา...")

	model := client.GenerativeModel("gemini-2.0-flash") // ใช้โมเดลที่เหมาะสม
	resp, err := model.GenerateContent(ctx, genai.Text(prompt))
	if err != nil {
		return "", fmt.Errorf("ข้อผิดพลาดในการเรียก Gemini API: %w", err)
	}

	if len(resp.Candidates) > 0 && len(resp.Candidates[0].Content.Parts) > 0 {
		if textPart, ok := resp.Candidates[0].Content.Parts[0].(genai.Text); ok {
			generatedContent := string(textPart)
			j := RemoveJsonBlock(generatedContent)
			// fmt.Println("เนื้อหาที่สร้างโดย AI:")
			// fmt.Println(generatedContent)
			return j, nil
		}
	}
	return "", fmt.Errorf("ไม่ได้รับเนื้อหาที่สามารถแปลงเป็นข้อความจาก AI ได้")
}


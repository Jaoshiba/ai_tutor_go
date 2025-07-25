package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go-fiber-template/domain/entities"
	"go-fiber-template/domain/repositories"

	"github.com/gofiber/fiber/v2" // Import Fiber to use its context
	"github.com/google/uuid"
	"google.golang.org/genai"
)

type ChapterServices struct {
	ChapterRepository repositories.IChapterRepository
}

type IChapterService interface {
	ChapterrizedText(fCtx *fiber.Ctx, text string) error
}

func NewChapterServices(chapterRepository repositories.IChapterRepository) IChapterService {
	return &ChapterServices{
		ChapterRepository: chapterRepository,
	}
}

func (c *ChapterServices) ChapterrizedText(fCtx *fiber.Ctx, text string) error {

	gemini_api_key := os.Getenv("GEMINI_API_KEY")
	if gemini_api_key == "" {
		return fmt.Errorf("GEMINI_API_KEY not set")
	}

	genaiCtx := fCtx.Context()

	client, err := genai.NewClient(genaiCtx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}
	// REMOVED: defer client.Close() -- This method does not exist for genai.Client

	promt := fmt.Sprintf(`Please divide the following Thai text into logical chapters or sections.
        Each chapter or section should have a title and its content.
        Ensure that the entire text is covered and logically structured.
        Crucially, **correct any spelling or grammatical errors** in the Thai text during this process, but **do not alter words that are already correct**.
        Do not add any introductory or concluding remarks, explanations, or preambles.
        The output **MUST ONLY** be a JSON object containing a "message" key and a "chapters" array.
        Each object in the "chapters" array MUST have 'chapterName' and 'content' keys.
        All string values within the JSON MUST be properly escaped (e.g., double quotes " must be \\").
        **For the 'content' key in each chapter, please format the text using Markdown, ensuring clear paragraph breaks and appropriate formatting (e.g., bolding, italics for emphasis) where necessary to enhance readability.**
        Example format:
        {
            "message": "File processed and chapterized successfully.",
            "chapters": [
                {
                    "chapterName": "ชื่อบทที่ 1",
                    "content": "เนื้อหาของบทที่ 1 ที่อาจมี \"เครื่องหมายคำพูด\" หรืออักขระพิเศษอื่นๆ\n\nนี่คือย่อหน้าใหม่ใน Markdown."
                },
                {
                    "chapterName": "ชื่อบทที่ 2",
                    "content": "เนื้อหาของบทที่ 2 ที่มี **ข้อความตัวหนา** และ *ข้อความตัวเอียง*..."
                }
            ]
        }
        Text to chapterize and correct:
        %s`, text)

	fmt.Println("Wait for Gemini to chapterizing your Text...")

	result, err := client.Models.GenerateContent(
		genaiCtx,
		"gemini-2.5-flash",
		genai.Text(promt),
		nil,
	)
	if err != nil {
		fmt.Println("Error generating chapters:", err)
		return err
	}

	chaps := RemoveJsonBlock(result.Text())
	fmt.Println("Finished your chapterize...")

	if chaps == "" {
		return fmt.Errorf("no chapters found in the response")
	}

	var response entities.GeminiResponse
	err = json.Unmarshal([]byte(chaps), &response)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return err
	}

	userIDRaw := fCtx.Locals("userID")
	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user ID in context")
	}

	// Consider if RoadmapId should be passed from the controller or generated here
	roadmapID := uuid.NewString()

	for _, chapter := range response.Chapters {
		ch := entities.ChapterDataModel{
			ChapterId:      uuid.NewString(),
			ChapterName:    chapter.ChapterName,
			UserID:         userIDStr,
			RoadmapId:      roadmapID,
			ChapterContent: chapter.Content,
			CreateAt:       time.Now(),
			UpdatedAt:      time.Now(),
			IsFinished:     false,
		}
		fmt.Println("Inserting chapter:", ch.ChapterId)
		err = c.ChapterRepository.InsertChapter(ch)
		if err != nil {
			return err
		}
	}

	return nil
}

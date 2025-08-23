package services

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"go-fiber-template/domain/entities"
	"go-fiber-template/domain/repositories"

	"log"

	"github.com/gofiber/fiber/v2" // Import Fiber to use its context
	"github.com/google/uuid"
	"google.golang.org/genai"
	// cohereClient "github.com/cohere-ai/cohere-go/v2/client"
)

type ChapterServices struct {
	ChapterRepository repositories.IChapterRepository
	PineconeRepo      repositories.IPineconeRepository
	GeminiService IGeminiService
}

type IChapterService interface {
	ChapterrizedText(ctx *fiber.Ctx, courseId string, moduleData entities.GenModule) error
	GetChaptersByModuleID(moduleID string) ([]entities.ChapterDataModel, error)
}

func NewChapterServices(chapterRepository repositories.IChapterRepository, pineconeRepo repositories.IPineconeRepository, GeminiService IGeminiService) IChapterService {
	if chapterRepository == nil {
		log.Fatal("❌ ChapterServices initialized with nil repository") // บรรทัดนี้คุณมีอยู่แล้ว
	} else {
		fmt.Println("✅ ChapterServices initialized with non-nil repository.")                   // เพิ่มบรรทัดนี้
		fmt.Printf("ChapterRepository instance in NewChapterServices: %p\n", chapterRepository) // เพิ่มบรรทัดนี้
	}
	return &ChapterServices{
		ChapterRepository: chapterRepository,
		PineconeRepo:      pineconeRepo,
		GeminiService: GeminiService,
	}
}

func (c *ChapterServices) ChapterrizedText(ctx *fiber.Ctx, courseId string, moduleData entities.GenModule) error {


	moduleTitle := moduleData.Title
	moduleDescription := moduleData.Description
	moduleContent := moduleData.Content

	gemini_api_key := os.Getenv("GEMINI_API_KEY")
	if gemini_api_key == "" {
		return fmt.Errorf("GEMINI_API_KEY not set")
	}

	genaiCtx := ctx.Context()

	client, err := genai.NewClient(genaiCtx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return err
	}
	// REMOVED: defer client.Close() -- This method does not exist for genai.Client

// 	prompt := fmt.Sprintf(
// 	"Role:\n" +
// 	"You are a professional Thai text processor and course content editor.\n" +
// 	"Your job is to IMMEDIATELY return the final result as JSON wrapped in a Markdown code block at START AND END\n\n" +
// 	"DONT USE BACK TICK OR MARK DOWN BLOCK IN CHAPTER CONTENT, \n" +
// 	"EXCEPT WHEN CONTENT MUST INCLUDE A CODE BLOCK (then keep ``` inside content safely)\n\n" +
// 	"Context:\n" +
// 	"The input is Thai text extracted from PDF/DOCX and may contain artifacts (bullets, control chars, duplicated spaces, broken hyphenation), spelling errors, or garbled symbols.\n" +
// 	"You must correct such errors while keeping already-correct words unchanged and preserving technical terms and numbers.\n\n" +
// 	"Hard rules (follow ALL, do not ask questions, do not wait for confirmation):\n" +
// 	"- START NOW and OUTPUT JSON ONLY inside a code block.\n" +
// 	"- The JSON must be wrapped with JSON markdown code block (```json).\n" +
// 	"- Do NOT include any text before/after the JSON block.\n" +
// 	"- Do NOT include analysis, apologies, or explanations.\n" +
// 	"- Do NOT translate content; keep Thai as-is (except for minor corrections).\n" +
// 	"- No hallucinations: use ONLY the given text.\n" +
// 	"- If some parts are irrelevant noise/artifacts, exclude them.\n" +
// 	"- If the text is very short, still produce at least 1 concise chapter.\n" +
// 	"- Prefer 3–8 chapters when the text is long; never put the whole input in one chapter.\n" +
// 	"- Escape any ``` inside content with backslash or keep as inline literal, to avoid breaking JSON.\n\n" +
// 	"Process:\n" +
// 	"1) Chapterize:\n" +
// 	"   - Divide the input into logical chapters/sections.\n" +
// 	"   - Each item must have a concise \"chapterName\" and its matching \"content\".\n" +
// 	"   - Content must strictly fit its \"chapterName\". No unrelated/duplicated text across chapters.\n" +
// 	"   - Ensure the overall coverage is maximized without dumping the entire input into a single chapter.\n" +
// 	"2) Clean Thai text:\n" +
// 	"   - Fix obvious spelling/grammar errors without changing correct words.\n" +
// 	"   - Normalize artifacts: remove stray bullets, control characters, duplicated whitespace; fix broken hyphenation and garbled symbols.\n" +
// 	"   - Preserve meaningful line breaks and lists.\n" +
// 	"3) Markdown formatting for \"content\":\n" +
// 	"   - Use paragraphs, **bold**, *italic*, bullet lists, numbered lists, and subheadings where suitable.\n" +
// 	"   - If code blocks are required inside content, keep ``` inside content safely (e.g., escape if needed)\n" +
// 	"4) Silent validation:\n" +
// 	"   - Each chapter’s content must match its \"chapterName\".\n" +
// 	"   - Chapters must not contain the full raw input.\n" +
// 	"   - Exclude irrelevant fragments and boilerplate.\n\n" +
// 	"Output schema (must appear inside the code block, and must be valid JSON):\n" +
// 	"```json\n" +
// 	"{\n" +
// 	"  \"message\": \"File processed and chapterized successfully.\",\n" +
// 	"  \"chapters\": [\n" +
// 	"    {\n" +
// 	"      \"chapterName\": \"ชื่อบทที่ 1\",\n" +
// 	"      \"content\": \"เนื้อหาของบทที่ 1 ที่เกี่ยวข้องเท่านั้น\\n\\n**หัวข้อย่อยตัวอย่าง** ...\"\n" +
// 	"    },\n" +
// 	"    {\n" +
// 	"      \"chapterName\": \"ชื่อบทที่ 2\",\n" +
// 	"      \"content\": \"เนื้อหาที่เข้ากับบทนี้ *เท่านั้น* ...\"\n" +
// 	"    }\n" +
// 	"  ]\n" +
// 	"}\n" +
// 	"```\n\n" +
// 	"Text to chapterize and correct:\n%s",
// 	text,
// )

prompt := fmt.Sprintf(
	"Role:\n"+
	"You are a professional multilingual text summarizer and Thai academic content editor.\n"+
	"Your job is to IMMEDIATELY return the final result as JSON wrapped in a Markdown code block (```json) at START AND END.\n\n"+
	"Context:\n"+
	"- moduleTitle: the main topic to define chapter boundaries\n"+
	"- moduleDescription: restricts chapter scope strictly within this description\n"+
	"- moduleContent: the main content source to generate chapters\n"+
	"- Input text may be in any language and should be translated to Thai for usability.\n"+
	"- Clean artifacts: bullets, control characters, duplicated spaces, broken hyphenation, spelling errors, or garbled symbols.\n"+
	"- Preserve correct words, technical terms, numbers, and meaningful line breaks or lists.\n\n"+
	"Hard rules:\n"+
	"- OUTPUT JSON ONLY inside a Markdown code block (```json). No extra text.\n"+
	"- ONLY RETURN AS chapterName: and content: DO NOT RETURN AS summary: or else"+
	"- Each chapter must have \"chapterName\" and matching \"content\".\n"+
	"- Scope chapters strictly within moduleTitle and moduleDescription.\n"+
	"- Do not hallucinate or add content not in moduleContent.\n"+
	"- Always produce at least 1 chapter, even if text is short.\n"+
	"- Prefer 3–8 chapters for long content.\n"+
	"- Clean Thai text: fix spelling/grammar, normalize artifacts.\n"+
	"- Markdown formatting allowed: paragraphs, **bold**, *italic*, bullet/numbered lists, subheadings.\n"+
	"- If code blocks are required inside content, keep ``` safely.\n\n"+
	"Process:\n"+
	"1) Chapterize:\n"+
	"   - Divide moduleContent into logical chapters/sections.\n"+
	"   - Each chapterName must be concise and meaningful.\n"+
	"   - Content must strictly match chapterName; remove unrelated text.\n"+
	"2) Translate:\n"+
	"   - Translate content into Thai while keeping technical terms intact.\n"+
	"3) Clean & format:\n"+
	"   - Remove stray bullets, control chars, duplicate spaces.\n"+
	"   - Fix broken hyphenation and garbled symbols.\n"+
	"   - Preserve meaningful line breaks and lists.\n\n"+
	"Output schema (valid JSON, inside ```json block):\n"+
	"```json\n"+
	"{\n"+
	"  \"message\": \"File processed and chapterized successfully.\",\n"+
	"  \"chapters\": [\n"+
	"    {\n"+
	"      \"chapterName\": \"ชื่อบทที่ 1\",\n"+
	"      \"content\": \"เนื้อหาของบทที่ 1 ที่เกี่ยวข้องเท่านั้น\\n\\n**หัวข้อย่อยตัวอย่าง** ...\"\n"+
	"    },\n"+
	"    {\n"+
	"      \"chapterName\": \"ชื่อบทที่ 2\",\n"+
	"      \"content\": \"เนื้อหาที่เข้ากับบทนี้ *เท่านั้น* ...\"\n"+
	"    }\n"+
	"  ]\n"+
	"}\n"+
	"```\n\n"+
	"Input data:\n"+
	"- moduleTitle: %s\n"+
	"- moduleDescription: %s\n"+
	"- moduleContent: %s",
	moduleTitle, moduleDescription, moduleContent,
)


	fmt.Println("Wait for Gemini to chapterizing your Text...")

	result, err := client.Models.GenerateContent(
		genaiCtx,
		"gemini-2.5-flash",
		genai.Text(prompt),
		nil,
	)
	if err != nil {
		fmt.Println("Error generating chapters:", err)
		return err
	}

	fmt.Println("Result heehee : ", result.Text())

	chaps := RemoveJsonBlock(result.Text())
	fmt.Println("Finished your chapterize...")
	fmt.Println("Chapter is : ", chaps)

	if chaps == "" {
		return fmt.Errorf("no chapters found in the response")
	}

	var response entities.GeminiResponse
	const maxRetries = 3

	for i := 0; i < maxRetries; i++ {
        err = json.Unmarshal([]byte(chaps), &response)
        if err == nil {
            break
        }

        fmt.Printf("Attempt %d/%d: Error unmarshalling JSON: %v\n", i+1, maxRetries, err)
        
        if i == maxRetries-1 {
            // ถ้าพยายามครบแล้วยังไม่สำเร็จ ให้ return error
            return fmt.Errorf("failed to unmarshal JSON after %d retries: %w", maxRetries, err)
        }

        // สร้าง prompt สำหรับแก้ไข
        fixPrompt := c.GeminiService.CreateFixPrompt(chaps, err.Error())
        fmt.Println("Sending fix prompt to Gemini...")
        
        // ส่ง prompt แก้ไขกลับไปให้ Gemini
        chaps, err = c.GeminiService.GenerateContentFromPrompt(ctx.Context(), fixPrompt)
        if err != nil {
            // ถ้าการเรียก API แก้ไขเกิดข้อผิดพลาด ให้ return error
            return fmt.Errorf("Gemini fix prompt generation failed on retry %d: %w", i+1, err)
        }
    }

	moduleIdRaw := ctx.Locals("moduleID")
	userIDRaw := ctx.Locals("userID")
	moduleId, ok := moduleIdRaw.(string)
	if !ok || moduleId == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing module ID in context")
	}

	userIDStr, ok := userIDRaw.(string)
	if !ok || userIDStr == "" {
		return fiber.NewError(fiber.StatusUnauthorized, "Invalid or missing user ID in context")
	}

	//create cohere client
	coheereapikey := os.Getenv("COHERE_API_KEY")
	if coheereapikey == "" {
		log.Fatal("COHERE_API_KEY is not set in .env")
	}
	// co := cohereClient.NewClient(cohereClient.WithToken(coheereapikey))

	nameSpaceName := userIDRaw.(string)

	fmt.Println("nameSpaceName: ", nameSpaceName)

	for _, chapter := range response.Chapters {
		ch := entities.ChapterDataModel{
			ChapterId:      uuid.NewString(),
			ChapterName:    chapter.ChapterName,
			UserID:         userIDStr,
			CourseId:       courseId,
			ChapterContent: chapter.Content,
			CreateAt:       time.Now(),
			UpdatedAt:      time.Now(),
			IsFinished:     false,
			ModuleId:       moduleId,
		}
		fmt.Println("Inserting chapter:", ch.ChapterId)

		fmt.Println("chapter : ", chapter)

		err = c.ChapterRepository.InsertChapter(ch)
		if err != nil {
			return err
		}

		// err = c.PineconeRepo.UpsertVector(ch, co, ctx, userIDStr)
		// if err != nil {
		// 	fmt.Println("Error inserting chapter:", err)
		// 	return err
		// }
	}

	return nil
}
func (c *ChapterServices) GetChaptersByModuleID(moduleID string) ([]entities.ChapterDataModel, error) {
	fmt.Println("im in chap service")
	fmt.Println(moduleID)
	if c.ChapterRepository == nil {
		log.Fatal("ChapterRepository is nil in ChapterServices")
	}
	chapters, err := c.ChapterRepository.GetChaptersByModuleID(moduleID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve chapters from repository for module %s: %w", moduleID, err)
	}
	return chapters, nil
}

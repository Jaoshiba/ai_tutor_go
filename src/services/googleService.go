package services

import (
	"context"
	"fmt"
	"os"
	"regexp"

	"google.golang.org/genai"
)

type IGoogleService interface {
	ChapterrizedText(text string) (string, error)
}

func ChapterrizedText(text string) (string, error) {

	gemini_api_key := (os.Getenv("GEMINI_API_KEY"))
	if gemini_api_key == "" {
		return "can't get gemini api key", nil
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return "", err
	}

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
		ctx,
		"gemini-2.5-flash",
		genai.Text(promt),
		nil,
	)
	if err != nil {
		return "", nil
	}

	chaps := removeJsonBlock(result.Text())
	fmt.Println("Finished your chapterize...")
	fmt.Println("Chapterized Text: " + chaps)

	// fmt.Println(gemini_api_key)

	return chaps, nil
}

func removeJsonBlock(text string) string {

	markdownJsonContentRegex := regexp.MustCompile("(?s)```json\\s*(.*?)\\s*```")
	matches := markdownJsonContentRegex.FindStringSubmatch(text)

	if len(matches) > 1 {
		return matches[1]
	}

	return ""

}

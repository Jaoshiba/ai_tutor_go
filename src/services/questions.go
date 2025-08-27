package services

import (
	"context"
	"fmt"
	"go-fiber-template/domain/entities"
	"os"

	"google.golang.org/genai"
)

type QuestionService struct {
}

type IQuestionService interface {
	QuestionsCreate(content string, difficulty string, questionNum int) ([]entities.QuestionDataModel, error)
}

func QuestionsCreate(content string, difficulty string, questionNum int) ([]entities.QuestionDataModel, error) {

	fmt.Println("Questions Generate called...")
	var questions []entities.QuestionDataModel

	gemini_api_key := (os.Getenv("GEMINI_API_KEY"))
	if gemini_api_key == "" {
		return questions, nil
	}

	ctx := context.Background()

	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey:  gemini_api_key,
		Backend: genai.BackendGeminiAPI,
	})
	if err != nil {
		return questions, err
	}

	promt := fmt.Sprintf(`
		โปรดสร้างคำถามแบบผสมจากข้อความที่ให้มา
		จำนวนคำถามทั้งหมดคือ %d ข้อ
		ระดับความยากโดยรวมคือ %s

		ผลลัพธ์จะต้องเป็น JSON array ของ objects เท่านั้น
		คำถามจะต้องมีรูปแบบผสมผสานจากประเภทต่อไปนี้ พร้อมโครงสร้าง JSON ที่ถูกต้อง:

		1.  **คำถามแบบเลือกตอบ (Multiple-Choice):**
			-   โครงสร้าง: {"type": "multiple-choice", "question": "...", "options": ["...", "...", "...", "..."], "answer": "..."}
			-   แต่ละคำถามต้องมี 4 ตัวเลือกเสมอ

		2.  **คำถามแบบเติมคำในช่องว่าง (Fill-in-the-Blank):**
			-   โครงสร้าง: {"type": "fill-in-the-blank", "question": "...", "answer": "..."}
			-   ช่องว่างในคำถามต้องแสดงด้วยขีดเส้นใต้ (เช่น "______")

		3.  **คำถามแบบเรียงลำดับ (Ordering/Sequencing):**
			-   โครงสร้าง: {"type": "ordering", "question": "...", "options": ["...", "..."], "answer": ["...", "..."]}
			-   อาร์เรย์ "options" ต้องมีขั้นตอนที่สลับตำแหน่งกัน
			-   อาร์เรย์ "answer" ต้องมีขั้นตอนเดียวกันที่เรียงลำดับอย่างถูกต้อง

		ห้ามใส่ข้อความนำหรือสรุปใด ๆ นอกเหนือจาก JSON array
		ค่าที่เป็น string ทั้งหมดใน JSON ต้องถูก escape อย่างถูกต้อง (เช่น " ต้องเป็น \\")
		---
		%s
		---`,
		questionNum,
		difficulty,
		content)

	fmt.Println("Wait for Gemini to create your questions...")
	fmt.Println("content: ", content)

	result, err := client.Models.GenerateContent(
		ctx,
		"gemini-2.5-flash",
		genai.Text(promt),
		nil,
	)
	if err != nil {
		fmt.Println("Error generating questions:", err)
		return questions, err
	}
	fmt.Println("Gemini response: ", result.Text())

	if result.Text() == "" {
		fmt.Println("Error generating questions: no content returned")
		return questions, fmt.Errorf("no content returned from Gemini")
	}

	fmt.Println("result from gemini: ", result.Text())

	questionsFromGemini := RemoveJsonBlock(result.Text())
	fmt.Println("Questions from Gemini: ", questionsFromGemini)
	if questionsFromGemini == "" {
		fmt.Println("Error parsing questions from Gemini: no content returned")
		return questions, fmt.Errorf("no questions generated from Gemini")
	}

	fmt.Println("result after remove json block: ", questionsFromGemini)

	// err = json.Unmarshal([]byte(questionsFromGemini), &questions)
	// if err != nil {
	// 	fmt.Println("Error unmarshalling JSON:", err)
	// 	return questions, err
	// }
	fmt.Println("Finished generating questions...")

	fmt.Println("question", questions)

	return questions, nil
}

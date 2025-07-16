package services

import (
	"context"
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"
	"os"
	"time"

	"github.com/google/uuid"
	"google.golang.org/genai"
)

type IExamService interface {
	ExamGenerate(chapters []entities.ChapterDataModel) ([]entities.QuestionDataModel, error)
}

func ExamGenerate(chapters []entities.ChapterDataModel) ([]entities.ExamDataModel, error) {

	var exams []entities.ExamDataModel

	for _, chapter := range chapters {
		content := chapter.Content
		questions, err := questionsCreate(content)
		if err != nil {
			return exams, err
		}
		exam := entities.ExamDataModel{
			ExamId:    uuid.NewString(),
			ChapterId: chapter.ChapterId,
			Questions: questions,
			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		}
		exams = append(exams, exam)
	}

	return exams, nil
}

func questionsCreate(content string) ([]entities.QuestionDataModel, error) {

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

	promt := fmt.Sprintf(`Please generate multiple-choice questions based on the provided text.

	For each question, ensure there are exactly 4 options, and one of them is the correct answer.

	The output MUST be a JSON array of objects.
	Each object in the array MUST strictly adhere to the following structure:
	{
	"question": "The generated question text.",
	"options": ["Option A", "Option B", "Option C", "Option D"],
	"answer": "The correct option text (must match one of the options)."
	}

	Do not include any introductory or concluding remarks, explanations, or preambles outside the JSON array.
	All string values within the JSON MUST be properly escaped (e.g., double quotes " must be \\").

	Example format:
	[
		{
			"question": "ข้อใดคือประเภทของระบบปฏิบัติการ?",
			"options": ["หน่วยความจำ", "โปรเซสเซอร์", "ซอฟต์แวร์ระบบ", "เครือข่าย"],
			"answer": "ซอฟต์แวร์ระบบ"
		},
		{
			"question": "ฟังก์ชันหลักของ CPU คืออะไร?",
			"options": ["เก็บข้อมูลถาวร", "ประมวลผลคำสั่ง", "แสดงผลกราฟิก", "เชื่อมต่ออินเทอร์เน็ต"],
			"answer": "ประมวลผลคำสั่ง"
		}
	]

	Please generate [NUMBER_OF_QUESTIONS] questions based on the following text:
	---
	%s
	---`, content)

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

	questionsFromGemini := removeJsonBlock(result.Text())
	fmt.Println("Questions from Gemini: ", questionsFromGemini)
	if questionsFromGemini == "" {
		fmt.Println("Error parsing questions from Gemini: no content returned")
		return questions, fmt.Errorf("no questions generated from Gemini")
	}

	err = json.Unmarshal([]byte(questionsFromGemini), &questions)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return questions, err
	}
	fmt.Println("Finished generating questions...")

	return questions, nil
}

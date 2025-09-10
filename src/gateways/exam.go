package gateways

import (
	// "fmt"
	// "go-fiber-template/domain/entities"

	"fmt"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) ExamGenerate(ctx *fiber.Ctx) error {
	input := ctx.Body()

	fmt.Println("input exam generate: ", string(input))

	var examrequest entities.ExamRequest

	err := ctx.BodyParser(&examrequest)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	fmt.Printf("examrequest Data chapterId: %s \n Content: %s \n Difficulty: %s \n QuestionNum: %d \n", examrequest.ChapterId, examrequest.Content, examrequest.Difficulty, examrequest.QuestionNum)

	err = h.ExamsService.ExamGenerate(examrequest, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.ResponseModel{
			Message: "Error Generating Exam",
			Data:    err.Error(),
			Status:  fiber.StatusInternalServerError,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
		Message: "ExamGenerate endpoint hit",
		Data:    string(input),
		Status:  fiber.StatusOK,
	})
}

func (h *HTTPGateway) GetExamsByModuleID(ctx *fiber.Ctx) error {

	moduleId := ctx.Params("moduleid")
	if moduleId == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "moduleId not found in input",
		})
	}

	exams, err := h.ExamsService.GetExamsByModuleID(moduleId)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": `error : `,
			"error":   err,
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Courses retrieved successfully",
		"data":    exams,
	})
}

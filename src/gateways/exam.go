package gateways

import (
	// "fmt"
	// "go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway)GetExamsByModuleID(ctx *fiber.Ctx) error {

	moduleId := ctx.Params("moduleid")
	if moduleId == "" {
		return ctx.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message" : "moduleId not found in input",
		})
	}

	
	exams, err := h.ExamsService.GetExamsByModuleID(moduleId)
	if err != nil{
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"message": `error : `,
		"error":err,
	})
	}
	

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Courses retrieved successfully",
		"data":    exams,
	})
}
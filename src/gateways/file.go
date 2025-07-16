package gateways

import (
	"fmt"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) UploadFile(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
	}
	fmt.Println("GetDocx func call.....")

	h.ModuleService.CreateModule(file, ctx)

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create module from your file"})
}

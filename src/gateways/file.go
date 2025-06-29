package gateways

import (
	"fmt"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (*HTTPGateway) UploadFile(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
	}

	fmt.Print(file)

	return nil
}

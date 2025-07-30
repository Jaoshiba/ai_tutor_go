package gateways

import (
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) TestPinecone(ctx *fiber.Ctx) error {

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

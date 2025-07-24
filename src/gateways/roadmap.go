package gateways

import (
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) CreateRoadmap(ctx *fiber.Ctx) error {

	body := ctx.BodyParser()

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create Roadmap from your promts"})
}

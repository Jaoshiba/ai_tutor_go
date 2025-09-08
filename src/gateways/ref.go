package gateways

import (
	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetRefsByModuleId(ctx *fiber.Ctx) error {

	moduleId := ctx.Params("moduleId")
	if moduleId == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Module ID is required",
		})
	}

	refs, err := h.DocSearchService.GetRefsByModuleId(moduleId, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to retrieve references",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"moduleId":   moduleId,
		"references": refs,
	})

}

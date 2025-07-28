package gateways

import (
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) UploadFile(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
	}

	roadmapname := ctx.FormValue("roadmapname")
	if roadmapname == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "Roadmapname required"})
	}

	var roadmapjsonBody entities.RoadmapRequestBody

	h.RoadmapService.CreateRoadmap(roadmapjsonBody, file, false, ctx)

	h.ModuleService.CreateModule(file, roadmapname, ctx)

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create module from your file"})
}

// func (h *HTTPGateway) GenCourse(ctx *fiber.Ctx) error{
// 	// jwt := ctx.Cookies('jwt')

// }

package gateways

import (
	"encoding/json"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) CreateRoadmap(ctx *fiber.Ctx) error {

	file, err := ctx.FormFile("file")
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
	}
	jsonbody := ctx.FormValue("jsonbody")

	var roadmapjsonBody entities.RoadmapRequestBody

	err = json.Unmarshal([]byte(jsonbody), &entities.RoadmapRequestBody{})
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	if roadmapjsonBody.RoadmapName == "" || roadmapjsonBody.Description == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed smt"})
	}

	err = h.RoadmapService.CreateRoadmap(roadmapjsonBody, file)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "failed to create roadmap on CreateRoadmap",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create Roadmap from your promts"})
}

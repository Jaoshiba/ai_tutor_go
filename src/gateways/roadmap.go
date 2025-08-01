// package gateways

// import (
// 	"encoding/json"
// 	"fmt"
// 	"go-fiber-template/domain/entities"

// 	"github.com/gofiber/fiber/v2"
// )

// func (h *HTTPGateway) CreateRoadmap(ctx *fiber.Ctx) error {

// 	file, err := ctx.FormFile("file")
// 	if err != nil {
// 		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
// 	}
// 	jsonbody := ctx.FormValue("jsonbody")

// 	fmt.Println("jsonbody: ", jsonbody)

// 	var roadmapjsonBody entities.RoadmapRequestBody

// 	err = json.Unmarshal([]byte(jsonbody), &roadmapjsonBody)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body"})
// 	}

// 	roadmapName := roadmapjsonBody.RoadmapName
// 	description := roadmapjsonBody.Description
// 	if roadmapName == "" {
// 		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed roadmap name"})
// 	} else if description == "" {
// 		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed description"})
// 	}

// 	err = h.RoadmapService.CreateRoadmap(roadmapjsonBody, file, true, ctx)
// 	if err != nil {
// 		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
// 			Message: "failed to create roadmap on CreateRoadmap",
// 		})
// 	}

// 	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create Roadmap from your promts"})
// }

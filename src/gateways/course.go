package gateways

import (
	"encoding/json"
	"fmt"
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) CreateCourse(ctx *fiber.Ctx) error {
	file, err := ctx.FormFile("file")
	// if err != nil {
	// 	return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid file"})
	// }
	jsonbody := ctx.FormValue("jsonbody")

	fmt.Println("jsonbody: ", jsonbody)

	var coursejsonBody entities.CourseRequestBody

	err = json.Unmarshal([]byte(jsonbody), &coursejsonBody)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	courseName := coursejsonBody.Title
	fmt.Println(courseName)
	description := coursejsonBody.Description
	fmt.Println(description)
	if courseName == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed course name"})
	} else if description == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed description"})
	}

	err = h.CourseService.CreateCourse(coursejsonBody, file, ctx)
	if err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
			Message: "failed to create course on CreateCourse",
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create Course from your promts"})
}

func (h *HTTPGateway) GetCourseByUser(c *fiber.Ctx) error {
	userID := c.Locals("userID")
	if userID == nil { 
		return c.Status(fiber.StatusUnauthorized).JSON(entities.ResponseMessage{Message: "User ID not found in context"})
	}

	
	userIDStr, ok := userID.(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "Invalid user ID in context"})
	}

	courses, err := h.CourseService.GetCourses(c)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(entities.ResponseMessage{Message: "Failed to retrieve courses"})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Courses retrieved successfully",
		"data":    courses,
	})
}
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

	Coursename := ctx.FormValue("Coursename")
	if Coursename == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseMessage{Message: "Coursename required"})
	}

	var CoursejsonBody entities.CourseRequestBody

	h.CourseService.CreateCourse(CoursejsonBody, file, false, ctx)

	// h.ModuleService.CreateModule(file, Coursename, ctx)

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "Completed create module from your file"})
}

// func (h *HTTPGateway) GenCourse(ctx *fiber.Ctx) error{
// 	// jwt := ctx.Cookies('jwt')

// }

package gateways

import (
	"go-fiber-template/domain/entities"
	"strings"

	"fmt"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) GetAllUserData(ctx *fiber.Ctx) error {

	data, err := h.UserService.GetAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get all users data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}




func (h *HTTPGateway) CreateUser(ctx *fiber.Ctx) error {

	//ต้องเช็คข้อมูลก่อนลง
	fmt.Println(string(ctx.Body()))
	bodyData := entities.UserDataModel{}
	if err := ctx.BodyParser(&bodyData); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body"})
	}

	//ต้องเพิ่มเช็คอีกก่อนลง
	bodyData.Email = strings.ToLower(strings.TrimSpace(bodyData.Email))
	bodyData.Username = strings.ToLower((strings.TrimSpace(bodyData.Username)))

	if bodyData.Username == "" || bodyData.Email == "" {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.ResponseMessage{Message: "invalid json body u missed smt"})
	}

	if err := h.UserService.InsertNewUser(bodyData); err != nil {
	return ctx.Status(fiber.StatusBadRequest).JSON(entities.ResponseModel{
		Message: err.Error(), // บอกว่า email ซ้ำ หรือ username ซ้ำ
	})
}

	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success"})
}

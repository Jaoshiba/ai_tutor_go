package gateways

import (
	"go-fiber-template/domain/entities"
	"strings"

	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

func (h *HTTPGateway) GetAllUserData(ctx *fiber.Ctx) error {

	data, err := h.UserService.GetAllUsers()
	if err != nil {
		return ctx.Status(fiber.StatusForbidden).JSON(entities.ResponseModel{Message: "cannot get all users data"})
	}
	return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{Message: "success", Data: data})
}

func (h *HTTPGateway) GetUserProfileById(ctx *fiber.Ctx) error {
    // 1) ลองเอาจาก params ก่อน
    userId := ctx.Params("userId")

    // 2) ถ้า params ว่าง ให้ลองจาก Locals ต่อ
    if userId == "" {
        if v := ctx.Locals("userID"); v != nil {
            switch t := v.(type) {
            case string:
                userId = t
            case uuid.UUID:
                userId = t.String()
            default:
                userId = fmt.Sprintf("%v", t)
            }
        }
    }

    // 3) ถ้ายังว่างอยู่ ให้ 400
    if userId == "" {
        return ctx.Status(fiber.StatusBadRequest).
            JSON(entities.ResponseMessage{Message: "missing userId (params or locals)"})
    }

    // (ออปชัน) ตรวจรูปแบบ UUID หากระบบคุณคาดหวัง UUID เท่านั้น
    if _, err := uuid.Parse(userId); err != nil {
        return ctx.Status(fiber.StatusBadRequest).
            JSON(entities.ResponseMessage{Message: "invalid userId"})
    }

    // 4) เรียก service
    user, err := h.UserService.GetUserPublicProfile(userId)
    if err != nil {
        return ctx.Status(fiber.StatusNotFound).
            JSON(entities.ResponseMessage{Message: "user not found"})
    }

    return ctx.Status(fiber.StatusOK).JSON(entities.ResponseModel{
        Message: "success",
        Data:    user,
    })
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



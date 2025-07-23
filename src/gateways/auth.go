package gateways

import (
	"go-fiber-template/domain/entities"

	"github.com/gofiber/fiber/v2"
)

func (h *HTTPGateway) Login(ctx *fiber.Ctx) error {
	var loginReq entities.LoginRequest
	
	// 1. Parse and Validate Request
	if err := ctx.BodyParser(&loginReq); err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).JSON(entities.LoginResponse{
			Success: false,
			Error:   "Invalid request format",
		})
	}

	// 2. Basic Field Validation
	if loginReq.Email == "" || loginReq.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.LoginResponse{
			Success: false,
			Error:   "Email and password are required",
		})
	}

	// 3. Call Auth Service
	response, err := h.AuthService.Login(loginReq.Email, loginReq.Password, ctx)
	if err != nil {
		// กรณีเกิด error ในการทำงานของระบบ
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.LoginResponse{
			Success: false,
			Error:   "Internal server error",
		})
	}

	// 4. Handle Service Response
	if !response.Success {
		// กรณีล็อกอินไม่สำเร็จ (credentials ไม่ถูกต้อง)
		return ctx.Status(fiber.StatusUnauthorized).JSON(response)
	}

	// 5. Success Response
	return ctx.Status(fiber.StatusOK).JSON(response)
}



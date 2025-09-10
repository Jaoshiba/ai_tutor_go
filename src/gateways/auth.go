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
	if loginReq.Email == "" || loginReq.Password == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(entities.LoginResponse{
			Success: false,
			Error:   "Email and password are required",
		})
	}
	
	response, err := h.AuthService.Login(loginReq.Email, loginReq.Password, ctx)
	if err != nil {
		// กรณีเกิด error ในการทำงานของระบบ
		return ctx.Status(fiber.StatusInternalServerError).JSON(entities.LoginResponse{
			Success: false,
			Error:   "Internal server error",
		})
	}

	if !response.Success {
		// กรณีล็อกอินไม่สำเร็จ (credentials ไม่ถูกต้อง)
		return ctx.Status(fiber.StatusUnauthorized).JSON(response)
	}

	return ctx.Status(fiber.StatusOK).JSON(response)
}

func (h *HTTPGateway) EmailVerify(ctx *fiber.Ctx) error {

	

	return nil
}

func (h *HTTPGateway) ResetPassword(ctx *fiber.Ctx) error {

	var req entities.ResetPasswordRequest

	if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
	email := req.Email
	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":"Email is require",
		})
	}


	err := h.ResetPasswordService.CreateResetRequest(ctx, email)
	if err!=nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error":err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"Reset Email Send Successfully",
	})
	
}
package gateways

import (
	// "fmt"
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
	token := ctx.Query("token")
	if token != "" {
		err := h.EmailVerificationSerivce.VerifyEmail(token)
		if err != nil {
			return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "internal server error",
			})
		}
	}
	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":"verify successful",
	})
}
func (h *HTTPGateway) ResendEmailVerification(ctx *fiber.Ctx) error {
	var req struct {
		Email string `json:"email" form:"email" query:"email"`
	}
	_ = ctx.QueryParser(&req)
	if req.Email == "" {
		_ = ctx.BodyParser(&req)
	}
	if req.Email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": "missing email",
		})
	}
	if err := h.EmailVerificationSerivce.ResendVerificationEmail(ctx, req.Email); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"success": false,
			"message": err.Error(),
		})
	}
	return ctx.JSON(fiber.Map{
		"success": true,
		"message": "verification email resent",
	})
}

func (h *HTTPGateway) ResetPasswordRequest(ctx *fiber.Ctx) error {

	var req entities.ResetPasswordRequest

	if err := ctx.BodyParser(&req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}
	email := req.Email
	if email == "" {
		return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Email is require",
		})
	}
	err := h.ResetPasswordService.CreateResetRequest(ctx, email)
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Reset Email Send Successfully",
	})

}

func (h *HTTPGateway) ResetPassword(ctx *fiber.Ctx) error {
    var req struct {
        Token       string `json:"token" form:"token"`
        NewPassword string `json:"new_password" form:"new_password"`
    }

    // รองรับทั้ง query, form, และ body
    if err := ctx.BodyParser(&req); err != nil {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Invalid request body",
        })
    }
    if req.Token == "" {
        req.Token = ctx.Query("token")
    }
    if req.NewPassword == "" {
        req.NewPassword = ctx.FormValue("new_password")
    }

    if req.Token == "" || req.NewPassword == "" {
        return ctx.Status(fiber.StatusBadRequest).JSON(fiber.Map{
            "error": "Token and new password are required",
        })
    }

    // เรียก service เพื่อรีเซ็ตรหัสผ่าน
    err := h.ResetPasswordService.ResetPassword(ctx, req.Token, req.NewPassword)
    if err != nil {
        return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
            "error": err.Error(),
        })
    }

    return ctx.Status(fiber.StatusOK).JSON(fiber.Map{
        "message": "Password reset successful",
    })
}
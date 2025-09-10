package gateways

import (
	// "fmt"

	// "github.com/gofiber/fiber/v2"
)

// func (h *HTTPGateway) SendVerifyEmail(ctx *fiber.Ctx) error {
// 	to := "bam.356345@gmail.com"
// 	subject := "Verify your email"
// 	body := "Hello, this is a verification test."
// 	if err := h.EmailService.SendEmail(to, subject, body); err != nil {
// 		return ctx.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
// 	}
// 	return ctx.Status(fiber.StatusOK).JSON(fiber.Map{"message": "email sent"})
// }
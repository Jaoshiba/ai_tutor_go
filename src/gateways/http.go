// gateway/http.go

package gateways

import (
	"github.com/gofiber/fiber/v2"
	service "go-fiber-template/src/services"
	auth "go-fiber-template/src/services/auth"
	"time" // สำหรับ time.Now() ใน LogoutHandler
)

type HTTPGateway struct {
	UserService       service.IUsersService
	FileService       service.IFileService
	GoogleAuthService auth.IGoogleOAuth
	AuthService       auth.IAuthService
}

// NewHTTPGateway รับ AuthService เข้ามาด้วย
func NewHTTPGateway(app *fiber.App, users service.IUsersService, files service.IFileService, googleAuth auth.IGoogleOAuth, authService auth.IAuthService) { // <--- ต้องรับ authService เข้ามา
	gateway := &HTTPGateway{
		UserService:       users,
		FileService:       files,
		GoogleAuthService: googleAuth,
		AuthService:       authService, // <--- กำหนดค่า
	}

	GatewayGoogleAuth(*gateway, app)
	GatewayUsers(*gateway, app)
	GatewayFile(*gateway, app)
	GatewayAuth(*gateway, app) // <--- ต้องมีการเรียก GatewayAuth
}

// LogoutHandler (ไม่มีการเปลี่ยนแปลง)
func (h *HTTPGateway) LogoutHandler(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour),
		HTTPOnly: true,
		Secure:   false, // <--- เปลี่ยนเป็น false ชั่วคราวสำหรับการพัฒนาด้วย HTTP
		SameSite: "Lax",
	})
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Logged out successfully",
	})
}

// GetMeDataHandler (ไม่มีการเปลี่ยนแปลง)
// (สมมติว่าคุณได้เพิ่ม User struct และ GetUserData ใน services/users.go แล้ว)
func (h *HTTPGateway) GetMeDataHandler(c *fiber.Ctx) error {

	jwtToken := c.Cookies("jwt")
	if jwtToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Missing JWT token",
		})
	}
	claims, err := h.AuthService.ValidateJWT(jwtToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"error": "Invalid or expired",
		})
	}
	user := map[string]interface{}{
		"id":      claims.UserID,
		"email":   claims.Email,
		"name":    claims.Email,
		"picture": "",
	}
	return c.JSON(fiber.Map{
		"isAuthenticated": true,
		"user":            user,
	})
}

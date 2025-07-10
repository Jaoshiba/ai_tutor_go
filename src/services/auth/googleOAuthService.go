// services/GoogleOAuthServices.go

package services

import (
	"context"
	"encoding/json"
	"os"
	"time" // Still needed for Fiber Cookie Expires

	"github.com/gofiber/fiber/v2"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// IGoogleOAuth interface
// ลบ GoogleAuthStatus ออกจาก interface นี้
type IGoogleOAuth interface {
	GoogleLoginHandler(c *fiber.Ctx) error
	GoogleCallbackHandler(c *fiber.Ctx) error
}

// googleOAuthService struct
type googleOAuthService struct {
	oauthConfig *oauth2.Config
	stateString string
	AuthService IAuthService // ยังคงต้องการ AuthService สำหรับ GenerateJWT ใน Callback
}

// NewGoogleOAuthService Constructor
func NewGoogleOAuthService(authService IAuthService) IGoogleOAuth {
	return &googleOAuthService{
		oauthConfig: &oauth2.Config{
			ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
			ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
			RedirectURL:  os.Getenv("GOOGLE_REDIRECT_URL"),
			Scopes: []string{
				"https://www.googleapis.com/auth/userinfo.email",
				"https://www.googleapis.com/auth/userinfo.profile",
			},
			Endpoint: google.Endpoint,
		},
		stateString: "random", // ควรใช้ session หรือ cookie จริง
		AuthService: authService,
	}
}

// GoogleLoginHandler method
func (s *googleOAuthService) GoogleLoginHandler(c *fiber.Ctx) error {
	url := s.oauthConfig.AuthCodeURL(s.stateString)
	return c.Redirect(url)
}

// GoogleCallbackHandler method
func (s *googleOAuthService) GoogleCallbackHandler(c *fiber.Ctx) error {
	state := c.Query("state")
	if state != s.stateString {
		return c.Status(fiber.StatusUnauthorized).SendString("Invalid OAuth state")
	}

	code := c.Query("code")
	if code == "" {
		return c.Status(fiber.StatusBadRequest).SendString("Missing code")
	}

	token, err := s.oauthConfig.Exchange(c.Context(), code)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("OAuth exchange error: " + err.Error())
	}

	client := s.oauthConfig.Client(context.Background(), token)
	resp, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to get user info: " + err.Error())
	}
	defer resp.Body.Close()

	var userInfo map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&userInfo); err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Decode error: " + err.Error())
	}

	userEmail, ok := userInfo["email"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("Email not found or invalid type in user info")
	}
	googleID, ok := userInfo["id"].(string)
	if !ok {
		return c.Status(fiber.StatusInternalServerError).SendString("ID not found or invalid type in user info")
	}

	// เรียกใช้ AuthService ที่ถูกฉีดเข้ามา เพื่อสร้าง JWT
	jwtToken, err := s.AuthService.GenerateJWT(googleID, userEmail)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).SendString("Failed to generate token: " + err.Error())
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    jwtToken,
		HTTPOnly: true,
		Secure:   false, // <--- เปลี่ยนเป็น false ชั่วคราวสำหรับการพัฒนาด้วย HTTP
		SameSite: "Lax",
		Expires:  time.Now().Add(24 * time.Hour),
	})

	return c.Redirect(os.Getenv("FRONTEND_URL")+"/dashboard", fiber.StatusFound)
}

// ลบเมธอด GoogleAuthStatus ออกจากไฟล์นี้
// func (s *googleOAuthService) GoogleAuthStatus(c *fiber.Ctx) error { ... }

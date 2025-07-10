// services/auth.go

package services

import (
	"fmt" // เพิ่ม import สำหรับ fmt
	"os"
	"strings" // เพิ่ม import สำหรับ strings
	"time"

	"log"

	"github.com/gofiber/fiber/v2" // เพิ่ม import สำหรับ Fiber Context
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID string `json:"user_id"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// IAuthService interface
type IAuthService interface {
	GenerateJWT(userID, email string) (string, error)
	ValidateJWT(tokenString string) (*UserClaims, error)
	CheckJWT(c *fiber.Ctx) error // <--- เพิ่มเมธอดนี้
}

type authService struct {
	jwtSecret string
}

func NewAuthService() IAuthService {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in .env")
	}
	return &authService{jwtSecret: secret}
}

func (s *authService) GenerateJWT(userID, email string) (string, error) {
	claims := &UserClaims{
		UserID: userID,
		Email:  email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtSecret))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}

func (s *authService) ValidateJWT(tokenString string) (*UserClaims, error) {
	claims := &UserClaims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, jwt.ErrSignatureInvalid
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, jwt.ErrInvalidKey
	}

	return claims, nil
}

// Implement CheckJWT method in authService
func (s *authService) CheckJWT(c *fiber.Ctx) error {
	jwtCookie := c.Cookies("jwt")
	if jwtCookie == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"isAuthenticated": false,
			"message":         "Unauthorized: Missing JWT token",
		})
	}

	claims, err := s.ValidateJWT(jwtCookie)
	if err != nil {
		// ให้ข้อความ error ที่ละเอียดขึ้น
		if strings.Contains(err.Error(), "token is expired") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"isAuthenticated": false,
				"message":         "Unauthorized: Token expired",
			})
		}
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
			"isAuthenticated": false,
			"message":         fmt.Sprintf("Unauthorized: Invalid token (%s)", err.Error()),
		})
	}

	// ถ้า Token ถูกต้อง คืนข้อมูลผู้ใช้และสถานะ true
	// ใน Production: อาจจะดึงข้อมูล user เพิ่มเติมจาก DB ด้วย claims.UserID
	// สำหรับตอนนี้: ใช้ข้อมูลจาก claims โดยตรง
	user := map[string]interface{}{
		"id":      claims.UserID,
		"email":   claims.Email,
		"name":    claims.Email, // อาจจะไม่มีชื่อเต็มใน claims, ใช้ email ชั่วคราว
		"picture": "",           // ไม่มีรูปใน claims, ต้องดึงจาก DB หรือ Google API อีกครั้งถ้าต้องการ
	}

	return c.JSON(fiber.Map{
		"isAuthenticated": true,
		"user":            user,
	})
}

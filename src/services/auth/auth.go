// services/auth.go

package services

import (
	"fmt" // เพิ่ม import สำหรับ fmt
	"go-fiber-template/domain/entities"
	repo "go-fiber-template/domain/repositories"
	"go-fiber-template/src/services/utils"
	"os"
	"strings" // เพิ่ม import สำหรับ strings
	"time"

	"log"

	"github.com/gofiber/fiber/v2" // เพิ่ม import สำหรับ Fiber Context
	"github.com/golang-jwt/jwt/v5"
)

type UserClaims struct {
	UserID string `json:"userid"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

// IAuthService interface
type IAuthService interface {
	GenerateJWT(userID, email string) (string, error)
	ValidateJWT(tokenString string) (*UserClaims, error)
	CheckJWT(c *fiber.Ctx) error
	Login(email, password string, ctx *fiber.Ctx) (entities.LoginResponse, error)
}

type authService struct {
	jwtSecret      string
	userRepository repo.IUsersRepository
}

func NewAuthService(userRepo repo.IUsersRepository) IAuthService {
	secret := os.Getenv("JWT_SECRET_KEY")
	if secret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in .env")
	}
	return &authService{
		jwtSecret:      secret,
		userRepository: userRepo,
	}
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
	fmt.Print("HEeHEeHeheeheh")
	jwtCookie := c.Cookies("jwt")
	fmt.Println("Cookies:", c.Cookies("jwt"))

	if jwtCookie == "" {
		fmt.Print("jwt is empty")
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

// services/auth.go

func (s *authService) Login(email, password string, ctx *fiber.Ctx) (entities.LoginResponse, error) {
	user, err := s.userRepository.FindByEmail(email)
	if err != nil {
		return entities.LoginResponse{
			Success: false,
			Error:   "Invalid email or password",
		}, nil // nil error เพราะถือว่าเป็น response ปกติ ไม่ใช่ error ระบบ
	}

	if err := utils.ComparePassword(user.Password, password); err != nil {
		return entities.LoginResponse{
			Success: false,
			Error:   "Invalid email or password",
		}, nil
	}

	token, err := s.GenerateJWT(user.UserID.String(), user.Email)
	if err != nil {
		return entities.LoginResponse{
			Success: false,
			Error:   "Failed to generate token",
		}, err // return error จริงเพราะเป็น error ระบบ
	}

	fmt.Print(token)

	ctx.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    token,                    // JWT ที่ generate มา
		Expires:  time.Now().Add(time.Hour * 24), // อายุ 1 วัน (แล้วแต่กำหนด)
		HTTPOnly: true,                           // ป้องกัน JS access, ปลอดภัยขึ้น
		Secure:   false,                           // ใช้เฉพาะ HTTPS
		SameSite: "None",                          // หรือ "Strict" / "None" ตาม use-case
	})

	respData := map[string]interface{}{
		"userid": user.UserID,
		"email":  user.Email,
		"role":   user.Role,
		"token":  token,
	}

	fmt.Print("success cookie")

	return entities.LoginResponse{
		Success: true,
		Data:    respData,
	}, nil
}

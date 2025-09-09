// gateway/http.go

package gateways

import (
	service "go-fiber-template/src/services"
	auth "go-fiber-template/src/services/auth"

	"github.com/gofiber/fiber/v2"
	// สำหรับ time.Now() ใน LogoutHandler
	"time"
)

type HTTPGateway struct {
	UserService       service.IUsersService
	ModuleService     service.IModuleService
	ExamsService      service.IExamService
	GoogleAuthService auth.IGoogleOAuth
	AuthService       auth.IAuthService
	ChapterService    service.IChapterService
	CourseService     service.ICourseService
	DocSearchService  service.IDocSearchService
}

func NewHTTPGateway(
	app *fiber.App,
	users service.IUsersService,
	modules service.IModuleService,
	exams service.IExamService,
	authService auth.IAuthService,
	chapterService service.IChapterService,
	courseService service.ICourseService,
	docSearchService service.IDocSearchService,
) {
	gateway := &HTTPGateway{
		UserService:      users,
		ModuleService:    modules,
		ExamsService:     exams,
		AuthService:      authService,
		ChapterService:   chapterService,
		CourseService:    courseService,
		DocSearchService: docSearchService,
	}

	GatewayAuth(*gateway, app)
	GatewayGoogleAuth(*gateway, app)
	GatewayUsers(*gateway, app)

	GatewayModules(*gateway, app)
	GatewayCourse(*gateway, app)
	GatewayExam(*gateway, app)
	GatewayRefs(*gateway, app)
	GatewayAskChat(*gateway, app)
	GatewayProtected(*gateway, app)

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

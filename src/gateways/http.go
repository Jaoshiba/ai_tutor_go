// gateway/http.go

package gateways

import (
	service "go-fiber-template/src/services"
	auth "go-fiber-template/src/services/auth"

	"github.com/gofiber/fiber/v2"
	// สำหรับ time.Now() ใน LogoutHandler
)

type HTTPGateway struct {
	UserService       service.IUsersService
	ModuleService     service.IModuleService
	GoogleAuthService auth.IGoogleOAuth
	AuthService       auth.IAuthService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, modules service.IModuleService, googleAuth auth.IGoogleOAuth, authService auth.IAuthService) {
	gateway := &HTTPGateway{
		UserService:       users,
		ModuleService:     modules,
		GoogleAuthService: googleAuth,
		AuthService:       authService,
	}

	GatewayGoogleAuth(*gateway, app)
	GatewayUsers(*gateway, app)
	GatewayModules(*gateway, app)
}

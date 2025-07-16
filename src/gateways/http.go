package gateways

import (
	service "go-fiber-template/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	UserService service.IUsersService
	ModuleServie service.IModuleService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, modules service.IModuleService) {
	gateway := &HTTPGateway{
		UserService: users,
		ModuleServie: modules,
	}

	GatewayUsers(*gateway, app)
	GatewayModules(*gateway, app)
}

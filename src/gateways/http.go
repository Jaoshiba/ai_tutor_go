package gateways

import (
	service "go-fiber-template/src/services"

	"github.com/gofiber/fiber/v2"
)

type HTTPGateway struct {
	UserService service.IUsersService
	FileService service.IFileService
}

func NewHTTPGateway(app *fiber.App, users service.IUsersService, files service.IFileService) {
	gateway := &HTTPGateway{
		UserService: users,
		FileService: files,
	}

	GatewayUsers(*gateway, app)
	GatewayFile(*gateway, app)
}

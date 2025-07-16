// gateway/route.go

package gateways

import (
	"github.com/gofiber/fiber/v2"
	//"go-fiber-template/src/middlewares" // เพิ่ม import สำหรับ middlewares
)

func GatewayUsers(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/v1/users")
	// ใช้ middleware กับ Group นี้
	//api.Use(middlewares.JWTAuthMiddleware(gateway.AuthService))

	api.Post("/add_user", gateway.CreateUser)
	api.Get("/users", gateway.GetAllUserData)
	api.Get("/me", gateway.GetMeDataHandler) // Endpoint สำหรับดึงข้อมูลผู้ใช้ที่ล็อกอินอยู่
}

func GatewayModules(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/v1/modules")
	api.Post("/upload", gateway.UploadFile)
}

func GatewayGoogleAuth(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/auth/google")
	api.Get("/login", gateway.GoogleLoginHandler)
	api.Get("/callback", gateway.GoogleCallback)
	// ลบ api.Get("/status/check", gateway.GoogleAuthStatus) ออกจากที่นี่
}

// GatewayAuth สำหรับ Logout และ Check Status
func GatewayAuth(gateway HTTPGateway, app *fiber.App) {
	authAPI := app.Group("/api/auth")
	authAPI.Post("/logout", gateway.LogoutHandler)
	authAPI.Get("/status/check", gateway.AuthService.CheckJWT) // <--- เพิ่มบรรทัดนี้: ให้ AuthService จัดการ CheckJWT
}

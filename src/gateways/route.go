// gateway/route.go

package gateways

import (
	"github.com/gofiber/fiber/v2"
	// "go-fiber-template/src/middlewares" // เพิ่ม import สำหรับ middlewares หากต้องการใช้ JWTAuthMiddleware แยก
)

// GatewayUsers สำหรับเส้นทางผู้ใช้ที่ไม่ได้ต้องการการป้องกัน (เช่น การสมัครสมาชิก)
func GatewayUsers(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/v1/users")
	api.Post("/signup", gateway.CreateUser) // เส้นทางสาธารณะ: สมัครสมาชิก
}

// GatewayModules สำหรับเส้นทางโมดูลที่ไม่ได้ต้องการการป้องกัน (หากมี)
// func GatewayModules(gateway HTTPGateway, app *fiber.App) {
// 	api := app.Group("/api/v1/modules")
// 	// api.Post("/some_public_module_action", gateway.SomePublicModuleHandler) // ตัวอย่าง: หากมีเส้นทางโมดูลสาธารณะ
// }

// GatewayGoogleAuth สำหรับเส้นทาง Google OAuth (สาธารณะ)
func GatewayGoogleAuth(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/auth/google")
	api.Get("/login", gateway.GoogleLoginHandler) // เส้นทางสาธารณะ: เริ่มต้นการล็อกอินด้วย Google
	api.Get("/callback", gateway.GoogleCallback)  // เส้นทางสาธารณะ: Callback จาก Google OAuth
}

// GatewayAuth สำหรับเส้นทางที่เกี่ยวข้องกับการตรวจสอบสิทธิ์ (Login เป็นสาธารณะ, อื่นๆ จะถูกย้ายไป Protected)
func GatewayAuth(gateway HTTPGateway, app *fiber.App) {
	authAPI := app.Group("/api/auth")
	authAPI.Post("/login", gateway.Login) // เส้นทางสาธารณะ: ล็อกอิน
	authAPI.Get("/status/check", gateway.AuthService.CheckJWT)
	// เส้นทาง Logout และ Check Status ถูกย้ายไปที่ GatewayProtected เนื่องจากต้องการการตรวจสอบสิทธิ์
}

func GatewayRoadmap(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/v1/roadmap")
	//api.Use(middlewares.JWTAuthMiddleware(gateway.AuthService))
	api.Post("/create", gateway.CreateRoadmap)
	api.Get("/pineconetest", gateway.TestPinecone)
}

func GatewayModules(gateway HTTPGateway, app *fiber.App) {
	api := app.Group("/api/v1/modules")
	api.Post("/upload", gateway.UploadFile)
}

// GatewayProtected สำหรับเส้นทางทั้งหมดที่ต้องการการตรวจสอบสิทธิ์ (Protected Routes)
func GatewayProtected(gateway HTTPGateway, app *fiber.App) {
	// สร้าง Group ของ routes ที่ต้องการการตรวจสอบสิทธิ์
	// ใช้ "/api/v1" เป็น Prefix เพื่อให้สอดคล้องกับโครงสร้าง URL ที่คุณมี
	protected := app.Group("/api/v1/protected")

	// 💡 นี่คือจุดที่คุณเรียกใช้ AuthMiddleware()
	// ทุก route ที่อยู่ภายใต้ group 'protected' จะต้องผ่าน AuthMiddleware ก่อนถึงจะเข้าถึง handler ได้
	protected.Use(gateway.AuthService.AuthMiddleware())

	// --- กำหนด Routes ภายใต้ Group ที่ได้รับการป้องกัน ---

	// Routes สำหรับ Users (ที่ต้องการการป้องกัน)
	protected.Get("/users", gateway.GetAllUserData)      // ดึงข้อมูลผู้ใช้ทั้งหมด (ต้องล็อกอิน)
	protected.Get("/users/me", gateway.GetMeDataHandler) // Endpoint สำหรับดึงข้อมูลผู้ใช้ที่ล็อกอินอยู่ (ต้องล็อกอิน)
	// protected.Get("/users/:id", gateway.GetUserByID)     // ถ้ามี: ดึงข้อมูลผู้ใช้ตาม ID (ต้องล็อกอิน)
	// protected.Put("/users/:id", gateway.UpdateUser)      // ถ้ามี: อัปเดตข้อมูลผู้ใช้ (ต้องล็อกอิน)
	// protected.Delete("/users/:id", gateway.DeleteUser)   // ถ้ามี: ลบผู้ใช้ (ต้องล็อกอิน)

	// Routes สำหรับ Modules (ที่ต้องการการป้องกัน)

	// protected.Post("/modules/upload", gateway.UploadFile)      // อัปโหลดไฟล์ (ต้องล็อกอิน)
	protected.Post("/modules/text", func(c *fiber.Ctx) error { // ตัวอย่าง route (ต้องล็อกอิน)
		return c.SendString("Protected module text route!")
	})
	// protected.Get("/modules", gateway.GetAllModules)     // ถ้ามี: ดึงข้อมูลโมดูลทั้งหมด (ต้องล็อกอิน)
	// protected.Get("/modules/:id", gateway.GetModuleByID) // ถ้ามี: ดึงข้อมูลโมดูลตาม ID (ต้องล็อกอิน)
	// protected.Put("/modules/:id", gateway.UpdateModule)  // ถ้ามี: อัปเดตโมดูล (ต้องล็อกอิน)
	// protected.Delete("/modules/:id", gateway.DeleteModule) // ถ้ามี: ลบโมดูล (ต้องล็อกอิน)

	// Routes สำหรับ Chapters (ที่ต้องการการป้องกัน) - ตัวอย่าง
	// สมมติว่า gateway.ChapterService.CreateChapter เป็น Fiber Handler
	// protected.Post("/chapters", gateway.ChapterService.CreateChapter)
	// protected.Get("/chapters", gateway.ChapterService.GetAllChapters)
	// ... เพิ่ม routes อื่นๆ ที่ต้องการการป้องกัน

	// Routes ที่เกี่ยวข้องกับการตรวจสอบสิทธิ์ที่ต้องการการป้องกัน (เช่น Logout, Check Status)
	// เนื่องจากคุณต้องการให้ AuthService จัดการ CheckJWT, และ Logout ก็ควรต้องล็อกอินอยู่แล้ว
	protected.Post("/auth/logout", gateway.LogoutHandler) // ล็อกเอาท์ (ต้องล็อกอิน) // ตรวจสอบสถานะ JWT (ต้องล็อกอิน)
}

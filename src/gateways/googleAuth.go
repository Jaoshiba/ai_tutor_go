// gateway/googleAuth.go (ปรับปรุง)

package gateways

// import (
// 	"github.com/gofiber/fiber/v2"
// )

// // GoogleLoginHandler เป็น method ของ HTTPGateway
// func (h *HTTPGateway) GoogleLoginHandler(ctx *fiber.Ctx) error {
// 	// เรียกใช้ method ของ GoogleService ที่ถูก inject เข้ามา
// 	return h.GoogleAuthService.GoogleLoginHandler(ctx) // <--- เปลี่ยนตรงนี้
// }

// // GoogleCallback เป็น method ของ HTTPGateway
// func (h *HTTPGateway) GoogleCallback(ctx *fiber.Ctx) error {
// 	// เรียกใช้ method ของ GoogleAuthService ที่ถูก inject เข้ามา
// 	return h.GoogleAuthService.GoogleCallbackHandler(ctx) // <--- เปลี่ยนตรงนี้
// }
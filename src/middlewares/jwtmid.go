// middlewares/jwtmid.go
package middlewares

// import (
//     "go-fiber-template/src/services/auth"
//     "github.com/gofiber/fiber/v2"
// )

// func JWTAuthMiddleware(authService auth.IAuthService) fiber.Handler {
//     return func(c *fiber.Ctx) error {
//         token := c.Get("Authorization")
//         if token == "" {
//             return c.Status(fiber.StatusUnauthorized).JSON("no token")
//         }

//         claims, err := authService.ValidateJWT(token[7:]) // ตัด "Bearer "
//         if err != nil {
//             return c.Status(fiber.StatusUnauthorized).JSON("")
//         }

//         c.Locals("user", claims)
//         return c.Next()
//     }
// }
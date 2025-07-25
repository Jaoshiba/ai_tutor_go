package main

import (
	"fmt"
	"go-fiber-template/configuration"
	ds "go-fiber-template/domain/datasources"
	repo "go-fiber-template/domain/repositories"
	gw "go-fiber-template/src/gateways"
	"go-fiber-template/src/middlewares"
	sv "go-fiber-template/src/services"
	authService "go-fiber-template/src/services/auth"

	"log"
	"os"

	"net/http"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/joho/godotenv"
)

func main() {
	// โหลด .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// สร้าง Fiber app
	app := fiber.New(configuration.NewFiberConfiguration())
	middlewares.Logger(app)
	app.Use(recover.New())

	// CORS Configuration
	app.Use(cors.New(cors.Config{

		AllowOrigins:     os.Getenv("FRONTEND_URL") + ", http://localhost:1818" + ", http://localhost:3000",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: true, // สำคัญมากสำหรับ Cookie
	}))

	// ดึง JWT_SECRET_KEY จาก .env
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in .env")
	}
	// สร้าง JWT Middleware

	// Connect PostgreSQL
	postgresql := ds.NewPostgresql()
	fmt.Printf("PostgreSQL DB instance before passing to repo: %p\n", postgresql)

	// สร้าง Repositories
	userRepo := repo.NewUsersRepositoryPostgres(postgresql)
	fileRepo := repo.NewModulesRepository(postgresql)
	chapterRepo := repo.NewChapterRepository(postgresql)
	roadmapRepo := repo.NewRoadmapRepository(postgresql)

	// สร้าง Services
	jwtSecret = os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in .env")
	}
	svAuth := authService.NewAuthService(userRepo)            // สร้าง AuthService
	sv0 := sv.NewUsersService(userRepo)                       // สร้าง UsersServic
	svGoogleAuth := authService.NewGoogleOAuthService(svAuth) // สร้าง GoogleOAuthService โดยฉีด AuthService

	svChapter := sv.NewChapterServices(chapterRepo)
	sv1 := sv.NewModuleService(fileRepo, svChapter)
	svRoadmap := sv.NewRoadmapService(roadmapRepo)

	// สร้าง Gateway และผูก Routes ทั้งหมด
	// ต้องส่ง AuthService และ UserService เข้าไปใน NewHTTPGateway ด้วย
	gw.NewHTTPGateway(app, sv0, sv1, svGoogleAuth, svAuth, svChapter, svRoadmap) // <--- ตรวจสอบพารามิเตอร์

	// ให้บริการไฟล์ static (เช่น dashboard.html)
	app.Use("/dashboard", filesystem.New(filesystem.Config{
		Root:       http.Dir("./static"),
		PathPrefix: "dashboard.html",
		Browse:     false,
		Index:      "dashboard.html",
	}))

	// เริ่มฟังที่พอร์ต
	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "1818"
	}
	log.Fatal(app.Listen(":" + PORT))
}

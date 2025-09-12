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

		AllowOrigins:     os.Getenv("FRONTEND_URL") + ", http://localhost:1818" + ", http://localhost:3000" + ", https://ai-tutor-frontend-gamma.vercel.app/",
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,HEAD,PUT,DELETE,PATCH",
		AllowCredentials: true, // สำคัญมากสำหรับ Cookie
	}))

	// ดึง JWT_SECRET_KEY จาก .env
	jwtSecret := os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in .env")
	}

	// Connect PostgreSQL
	postgresql := ds.NewPostgresql()
	fmt.Printf("PostgreSQL DB instance before passing to repo: %p\n", postgresql)
	pineconeIdxConn, err := ds.NewPincecone()
	if err != nil {
		log.Fatalf("Failed to connect to Pinecone: %v", err)
	}

	// สร้าง Repositories
	userRepo := repo.NewUsersRepositoryPostgres(postgresql)
	fileRepo := repo.NewModulesRepository(postgresql)
	chapterRepo := repo.NewChapterRepository(postgresql)
	courseRepo := repo.NewCourseRepository(postgresql)
	pineconeRepo := repo.NewPineconeRepository(pineconeIdxConn)
	examRepo := repo.NewExamRepository(postgresql)
	refRepo := repo.NewRefRepository(postgresql)
	resetPasswordRepo := repo.NewResetPasswordRepository(postgresql)
	learningProgressRepo := repo.NewLearningProgressRepository(postgresql)
	emailVerificationRepo := repo.NewEmailVerificationRepository(postgresql)

	questionrepo := repo.NewQuestionRepository(postgresql)
	if userRepo == nil || fileRepo == nil || chapterRepo == nil || courseRepo == nil {
		log.Fatalf("One or more repositories failed to initialize and are NIL.")
	}

	// สร้าง Services
	jwtSecret = os.Getenv("JWT_SECRET_KEY")
	if jwtSecret == "" {
		log.Fatal("JWT_SECRET_KEY is not set in .env")
	}
	svEmailVerification := sv.NewEmailVerificationService(emailVerificationRepo, userRepo)
	svAuth := authService.NewAuthService(userRepo, svEmailVerification) // สร้าง AuthService
	sv0 := sv.NewUsersService(userRepo)                                 // สร้าง UsersServic
	geminiService := sv.NewGeminiService()
	svQuestion := sv.NewQuestionService(questionrepo)

	svExam := sv.NewExamService(examRepo, chapterRepo, svQuestion)
	svdocSearch := sv.NewDocSearchService(refRepo, pineconeRepo)
	svChapter := sv.NewChapterServices(chapterRepo, pineconeRepo, geminiService, svExam)
	sv1 := sv.NewModuleService(fileRepo, svChapter, svExam, svdocSearch, geminiService)
	svCourse := sv.NewCourseService(courseRepo, sv1, geminiService, svChapter, learningProgressRepo)
	svResetPassword := sv.NewResetPasswordService(resetPasswordRepo, userRepo)

	// สร้าง Gateway และผูก Routes ทั้งหมด
	// ต้องส่ง AuthService และ UserService เข้าไปใน NewHTTPGateway ด้วย
	gw.NewHTTPGateway(app, sv0, sv1, svExam, svAuth, svChapter, svCourse, svdocSearch, svResetPassword, svEmailVerification)

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

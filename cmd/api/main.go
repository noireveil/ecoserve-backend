package main

import (
	"log"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
	"github.com/gofiber/swagger"
	"github.com/joho/godotenv"

	"github.com/noireveil/ecoserve-backend/internal/config"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/handlers"
	"github.com/noireveil/ecoserve-backend/internal/repository"
	"github.com/noireveil/ecoserve-backend/internal/usecase"

	_ "github.com/noireveil/ecoserve-backend/docs"
)

// @title EcoServe API
// @version 1.0
// @description REST API untuk Platform Manajemen Siklus Hidup Elektronik & Servis Sirkular.
// @BasePath /
// @securityDefinitions.apikey ApiKeyAuth
// @in cookie
// @name jwt
func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Gagal memuat file .env, menggunakan variabel lingkungan sistem")
	}

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		AppName: "EcoServe API v1.0",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	// Konfigurasi Rate Limiter
	authLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Spam OTP terdeteksi. Silakan coba lagi dalam 1 menit.",
			})
		},
	})

	dashboardLimiter := limiter.New(limiter.Config{
		Max:        60,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak permintaan. Silakan tunggu sebentar.",
			})
		},
	})

	app.Use("/api/users/auth", authLimiter)
	app.Use("/api", dashboardLimiter)

	app.Get("/swagger/*", swagger.HandlerDefault)

	// Health Check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "EcoServe Backend and Database are Online",
		})
	})

	// Inisialisasi Repositori
	userRepo := repository.NewUserRepository(config.DB)
	techRepo := repository.NewTechnicianRepository(config.DB)
	orderRepo := repository.NewOrderRepository(config.DB)
	deviceRepo := repository.NewDeviceRepository(config.DB)

	// Inisialisasi Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	techUsecase := usecase.NewTechnicianUsecase(techRepo, userRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)

	// Inisialisasi Handler
	handlers.NewUserHandler(app, userUsecase)
	handlers.NewTechnicianHandler(app, techUsecase)
	handlers.NewOrderHandler(app, orderUsecase)
	handlers.NewDeviceHandler(app, deviceUsecase)
	handlers.NewChatbotHandler(app, techUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "3000"
	}

	log.Println("Server REST API berjalan di port:", port)
	log.Fatal(app.Listen(":" + port))
}

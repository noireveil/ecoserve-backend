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

	app.Use("/api", limiter.New(limiter.Config{
		Max:        5,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Terlalu banyak permintaan ke server. Sistem mendeteksi potensi spam. Silakan coba lagi dalam 1 menit.",
			})
		},
	}))

	app.Get("/swagger/*", swagger.HandlerDefault)

	// @Summary Cek Status Peladen
	// @Description Memeriksa ketersediaan API dan konektivitas basis data EcoServe
	// @Tags Base
	// @Produce json
	// @Success 200 {object} map[string]interface{}
	// @Router /health [get]
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "EcoServe Backend and Database are Online",
		})
	})

	// Inisialisasi Repository
	userRepo := repository.NewUserRepository(config.DB)
	techRepo := repository.NewTechnicianRepository(config.DB)
	orderRepo := repository.NewOrderRepository(config.DB)
	deviceRepo := repository.NewDeviceRepository(config.DB)

	// Inisialisasi Usecase
	userUsecase := usecase.NewUserUsecase(userRepo)
	techUsecase := usecase.NewTechnicianUsecase(techRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)

	// Inisialisasi Handlers
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

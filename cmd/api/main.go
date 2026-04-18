package main

import (
	"errors"
	"log"
	"os"
	"os/signal"
	"syscall"
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
	if err := godotenv.Load(); err != nil {
		log.Println("Peringatan: Gagal memuat file .env, menggunakan variabel lingkungan sistem")
	}

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		AppName: "EcoServe API v1.0",
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			code := fiber.StatusInternalServerError
			var e *fiber.Error
			if errors.As(err, &e) {
				code = e.Code
			}

			log.Printf("[SERVER ERROR] Path: %s | Message: %v\n", c.Path(), err)

			message := err.Error()
			if code >= 500 {
				message = "Terjadi gangguan pada sistem internal. Log telah dicatat dan tim kami sedang menanganinya."
			}

			return c.Status(code).JSON(fiber.Map{
				"error": message,
			})
		},
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     os.Getenv("ALLOWED_ORIGINS"),
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
	}))

	authLimiter := limiter.New(limiter.Config{
		Max:        3,
		Expiration: 1 * time.Minute,
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP()
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Spam terdeteksi. Silakan coba lagi dalam 1 menit.",
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

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "EcoServe Backend and Database are Online",
		})
	})

	userRepo := repository.NewUserRepository(config.DB)
	techRepo := repository.NewTechnicianRepository(config.DB)
	orderRepo := repository.NewOrderRepository(config.DB)
	deviceRepo := repository.NewDeviceRepository(config.DB)
	reviewRepo := repository.NewReviewRepository(config.DB)

	userUsecase := usecase.NewUserUsecase(userRepo)
	techUsecase := usecase.NewTechnicianUsecase(techRepo, userRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo, techRepo)
	deviceUsecase := usecase.NewDeviceUsecase(deviceRepo)
	reviewUsecase := usecase.NewReviewUsecase(reviewRepo, orderRepo)

	handlers.NewUserHandler(app, userUsecase)
	handlers.NewTechnicianHandler(app, techUsecase)
	handlers.NewOrderHandler(app, orderUsecase)
	handlers.NewDeviceHandler(app, deviceUsecase)
	handlers.NewChatbotHandler(app, techUsecase)
	handlers.NewReviewHandler(app, reviewUsecase)

	port := os.Getenv("PORT")
	if port == "" {
		port = os.Getenv("APP_PORT")
	}
	if port == "" {
		port = "3000"
	}

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-c
		log.Println("Memulai proses graceful shutdown...")
		_ = app.Shutdown()
	}()

	log.Println("Server REST API berjalan di port:", port)
	if err := app.Listen(":" + port); err != nil {
		log.Panic(err)
	}

	log.Println("Proses pembersihan selesai. Server dimatikan.")
}

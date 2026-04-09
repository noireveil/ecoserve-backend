package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
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
// @host localhost:3000
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
		AllowOrigins:     "http://localhost:5173, http://localhost:3000",
		AllowCredentials: true,
		AllowHeaders:     "Origin, Content-Type, Accept",
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

	userRepo := repository.NewUserRepository(config.DB)
	techRepo := repository.NewTechnicianRepository(config.DB)
	orderRepo := repository.NewOrderRepository(config.DB)

	userUsecase := usecase.NewUserUsecase(userRepo)
	techUsecase := usecase.NewTechnicianUsecase(techRepo)
	orderUsecase := usecase.NewOrderUsecase(orderRepo)

	handlers.NewUserHandler(app, userUsecase)
	handlers.NewTechnicianHandler(app, techUsecase)
	handlers.NewOrderHandler(app, orderUsecase)
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

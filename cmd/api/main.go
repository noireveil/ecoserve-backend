package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"

	"github.com/noireveil/ecoserve-backend/internal/config"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/handlers"
	"github.com/noireveil/ecoserve-backend/internal/repository"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Peringatan: Gagal memuat file .env, menggunakan variabel lingkungan sistem")
	}

	config.ConnectDB()

	app := fiber.New(fiber.Config{
		AppName: "EcoServe API v1.0",
	})
	app.Use(cors.New())

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

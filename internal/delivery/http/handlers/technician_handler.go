package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type TechnicianHandler struct {
	techUsecase usecase.TechnicianUsecase
}

func NewTechnicianHandler(app *fiber.App, usecase usecase.TechnicianUsecase) {
	handler := &TechnicianHandler{techUsecase: usecase}

	api := app.Group("/api/technicians")
	api.Post("/", handler.Register)
	api.Get("/nearby", handler.GetNearby)
}

func (h *TechnicianHandler) Register(c *fiber.Ctx) error {
	var req struct {
		domain.Technician
		Longitude float64 `json:"longitude"`
		Latitude  float64 `json:"latitude"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.techUsecase.RegisterTechnician(&req.Technician, req.Longitude, req.Latitude); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Teknisi berhasil didaftarkan"})
}

func (h *TechnicianHandler) GetNearby(c *fiber.Ctx) error {
	lon, _ := strconv.ParseFloat(c.Query("lon"), 64)
	lat, _ := strconv.ParseFloat(c.Query("lat"), 64)
	radius, _ := strconv.Atoi(c.Query("radius_km", "10"))

	technicians, err := h.techUsecase.GetNearbyTechnicians(lon, lat, radius)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": technicians})
}

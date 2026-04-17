package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type TechnicianHandler struct {
	techUsecase usecase.TechnicianUsecase
}

type RegisterTechnicianRequest struct {
	domain.Technician
	Longitude float64 `json:"longitude" example:"106.8229"`
	Latitude  float64 `json:"latitude" example:"-6.1944"`
}

func NewTechnicianHandler(app *fiber.App, usecase usecase.TechnicianUsecase) {
	handler := &TechnicianHandler{techUsecase: usecase}

	api := app.Group("/api/technicians")
	api.Post("/", handler.Register)
	api.Get("/nearby", handler.GetNearby)
	api.Get("/performance", middleware.Protected(), handler.GetPerformance)
}

// @Summary Mendaftarkan Teknisi Baru
// @Description Menambahkan data teknisi baru beserta titik koordinat operasinya.
// @Tags Technicians
// @Accept json
// @Produce json
// @Param request body RegisterTechnicianRequest true "Data Teknisi"
// @Success 201 {object} map[string]interface{}
// @Router /api/technicians/ [post]
func (h *TechnicianHandler) Register(c *fiber.Ctx) error {
	var req RegisterTechnicianRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.techUsecase.RegisterTechnician(&req.Technician, req.Longitude, req.Latitude); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Teknisi berhasil didaftarkan"})
}

// @Summary Cari Teknisi Terdekat
// @Description Melakukan kueri spasial (PostGIS) untuk mencari teknisi dalam radius tertentu.
// @Tags Technicians
// @Produce json
// @Param lon query number true "Garis Bujur (Longitude)"
// @Param lat query number true "Garis Lintang (Latitude)"
// @Param radius_km query int false "Radius pencarian dalam KM (Default: 10)"
// @Success 200 {object} map[string]interface{}
// @Router /api/technicians/nearby [get]
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

// @Summary Mendapatkan Metrik Performa Teknisi
// @Description Mengambil agregasi data performa (Rating, Total Perbaikan, dan Total CO2 yang diselamatkan) dari teknisi yang sedang login.
// @Tags Technicians
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/technicians/performance [get]
func (h *TechnicianHandler) GetPerformance(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	if role != "technician" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: Hanya untuk teknisi"})
	}

	perf, err := h.techUsecase.GetPerformance(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil metrik performa teknisi",
		"data":    perf,
	})
}

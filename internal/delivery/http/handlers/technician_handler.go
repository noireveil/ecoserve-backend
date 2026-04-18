package handlers

import (
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type TechnicianHandler struct {
	techUsecase usecase.TechnicianUsecase
}

type RegisterTechnicianRequest struct {
	Specialization  string  `json:"specialization" example:"Pendingin & Komersial"`
	ExperienceYears int     `json:"experience_years" example:"5"`
	Longitude       float64 `json:"longitude" example:"106.8229"`
	Latitude        float64 `json:"latitude" example:"-6.1944"`
}

type UpdateAvailabilityPayload struct {
	IsAvailable bool `json:"is_available"`
}

func NewTechnicianHandler(app *fiber.App, usecase usecase.TechnicianUsecase) {
	handler := &TechnicianHandler{techUsecase: usecase}

	api := app.Group("/api/technicians")
	api.Post("/", middleware.Protected(), handler.Register)
	api.Get("/nearby", handler.GetNearby)
	api.Get("/performance", middleware.Protected(), handler.GetPerformance)
	api.Get("/earnings", middleware.Protected(), handler.GetEarnings)
	api.Put("/availability", middleware.Protected(), handler.UpdateAvailability)
}

// @Summary Mendaftarkan Teknisi Baru
// @Description Menambahkan data teknisi baru beserta titik koordinat operasinya.
// @Tags Technicians
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body RegisterTechnicianRequest true "Data Teknisi"
// @Success 201 {object} map[string]interface{}
// @Router /api/technicians/ [post]
func (h *TechnicianHandler) Register(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Sesi tidak valid"})
	}

	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	var req RegisterTechnicianRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	tech := domain.Technician{
		UserID:          userID,
		Specialization:  req.Specialization,
		ExperienceYears: req.ExperienceYears,
	}

	if err := h.techUsecase.RegisterTechnician(&tech, req.Longitude, req.Latitude); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Teknisi berhasil didaftarkan"})
}

// @Summary Cari Teknisi Terdekat
// @Description Melakukan kueri spasial (PostGIS) untuk mencari teknisi dalam radius tertentu yang sedang ONLINE.
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

// @Summary Ubah Status Ketersediaan Teknisi
// @Description Mengubah status online/offline teknisi. Jika offline, teknisi tidak akan muncul dalam pencarian geospasial pesanan publik.
// @Tags Technicians
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body UpdateAvailabilityPayload true "Status Ketersediaan"
// @Success 200 {object} map[string]interface{}
// @Router /api/technicians/availability [put]
func (h *TechnicianHandler) UpdateAvailability(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	role, _ := c.Locals("role").(string)
	if role != "technician" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"error": "Akses ditolak: Hanya untuk teknisi"})
	}

	var req UpdateAvailabilityPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.techUsecase.UpdateAvailability(userIDStr, req.IsAvailable); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memperbarui status ketersediaan"})
	}

	statusMsg := "Offline (Tidak Menerima Pesanan)"
	if req.IsAvailable {
		statusMsg = "Online (Menerima Pesanan)"
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Status ketersediaan berhasil diubah menjadi: " + statusMsg,
	})
}

// @Summary Mendapatkan Metrik Performa Teknisi
// @Description Mengambil agregasi data performa dari teknisi yang sedang login.
// @Tags Technicians
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/technicians/performance [get]
func (h *TechnicianHandler) GetPerformance(c *fiber.Ctx) error {
	userIDStr, _ := c.Locals("user_id").(string)
	perf, err := h.techUsecase.GetPerformance(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": perf})
}

// @Summary Mendapatkan Pendapatan Teknisi
// @Description Mengambil total pendapatan teknisi.
// @Tags Technicians
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/technicians/earnings [get]
func (h *TechnicianHandler) GetEarnings(c *fiber.Ctx) error {
	userIDStr, _ := c.Locals("user_id").(string)
	earnings, err := h.techUsecase.GetEarnings(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"data": earnings})
}

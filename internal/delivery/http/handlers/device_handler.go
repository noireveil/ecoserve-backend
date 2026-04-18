package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type DeviceHandler struct {
	deviceUsecase usecase.DeviceUsecase
}

type CreateDevicePayload struct {
	Category   string  `json:"category" example:"Smartphone"`
	BrandName  string  `json:"brand_name" example:"Apple iPhone 14 Pro"`
	WeightInKg float64 `json:"weight_in_kg" example:"0.21"`
}

func NewDeviceHandler(app *fiber.App, usecase usecase.DeviceUsecase) {
	handler := &DeviceHandler{deviceUsecase: usecase}

	api := app.Group("/api/devices")
	api.Post("/", middleware.Protected(), handler.Create)
	api.Get("/", middleware.Protected(), handler.GetMyDevices)
	api.Get("/:id", middleware.Protected(), handler.GetDetail)
	api.Delete("/:id", middleware.Protected(), handler.Delete)
}

// @Summary Mendapatkan Detail Perangkat (DPP)
// @Description Mengambil informasi lengkap satu perangkat berdasarkan UUID.
// @Tags Devices
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID Perangkat (UUID)"
// @Success 200 {object} map[string]interface{}
// @Router /api/devices/{id} [get]
func (h *DeviceHandler) GetDetail(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	device, err := h.deviceUsecase.GetDeviceByID(deviceID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Perangkat tidak ditemukan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil detail perangkat",
		"data":    device,
	})
}

// @Summary Mendaftarkan Perangkat (DPP)
// @Description Menambahkan perangkat elektronik baru ke garasi pengguna untuk pembuatan Digital Product Passport.
// @Tags Devices
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateDevicePayload true "Data Perangkat"
// @Success 201 {object} map[string]interface{}
// @Router /api/devices/ [post]
func (h *DeviceHandler) Create(c *fiber.Ctx) error {
	var req CreateDevicePayload

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	ownerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format User ID salah"})
	}

	device := domain.DigitalProductPassport{
		OwnerID:      ownerID,
		Category:     req.Category,
		BrandName:    req.BrandName,
		WeightInKg:   req.WeightInKg,
		PurchaseDate: time.Now(),
	}

	if err := h.deviceUsecase.CreateDevice(&device); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menyimpan perangkat"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Perangkat berhasil didaftarkan di Digital Product Passport",
		"data":    device,
	})
}

// @Summary Mengambil Garasi Perangkat
// @Description Mengambil daftar semua perangkat (DPP) yang dimiliki oleh pengguna.
// @Tags Devices
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/devices/ [get]
func (h *DeviceHandler) GetMyDevices(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak ditemukan pada token"})
	}

	devices, err := h.deviceUsecase.GetUserDevices(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data perangkat"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil data garasi perangkat",
		"data":    devices,
	})
}

// @Summary Menghapus Perangkat (DPP)
// @Description Menghapus perangkat elektronik dari garasi pengguna.
// @Tags Devices
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID Perangkat (UUID)"
// @Success 200 {object} map[string]interface{}
// @Router /api/devices/{id} [delete]
func (h *DeviceHandler) Delete(c *fiber.Ctx) error {
	deviceID := c.Params("id")
	userIDStr, _ := c.Locals("user_id").(string)

	err := h.deviceUsecase.DeleteDevice(deviceID, userIDStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Perangkat tidak ditemukan atau Anda tidak berhak menghapusnya"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Perangkat berhasil dihapus dari garasi",
	})
}

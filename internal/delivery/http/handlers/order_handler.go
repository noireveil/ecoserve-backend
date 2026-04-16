package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

type CreateOrderPayload struct {
	TechnicianID       string  `json:"technician_id" example:"(Opsional) uuid-teknisi"`
	DeviceCategory     string  `json:"device_category" example:"Pendingin & Komersial"`
	ProblemDescription string  `json:"problem_description" example:"Kompresor mati dan berasap"`
	CustomerLatitude   float64 `json:"customer_latitude"`
	CustomerLongitude  float64 `json:"customer_longitude"`
}

func NewOrderHandler(app *fiber.App, usecase usecase.OrderUsecase) {
	handler := &OrderHandler{orderUsecase: usecase}

	api := app.Group("/api/orders")
	api.Post("/", middleware.Protected(), handler.Create)
	api.Get("/", middleware.Protected(), handler.GetMyOrders)
	api.Get("/incoming", middleware.Protected(), handler.GetIncomingOrders)
	api.Put("/:id/complete", middleware.Protected(), handler.Complete)
}

// @Summary Mendapatkan Pesanan Masuk
// @Description Mengambil daftar pesanan yang berstatus PENDING dan belum memiliki teknisi.
// @Tags Orders
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/orders/incoming [get]
func (h *OrderHandler) GetIncomingOrders(c *fiber.Ctx) error {
	orders, err := h.orderUsecase.GetIncomingOrders()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pesanan masuk"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil daftar pesanan masuk",
		"data":    orders,
	})
}

// @Summary Mendapatkan Riwayat Pesanan
// @Description Mengambil daftar pesanan yang terkait dengan pengguna (sebagai Konsumen atau Teknisi).
// @Tags Orders
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Router /api/orders/ [get]
func (h *OrderHandler) GetMyOrders(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak ditemukan pada token"})
	}

	orders, err := h.orderUsecase.GetUserOrders(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data pesanan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil riwayat pesanan",
		"data":    orders,
	})
}

// @Summary Membuat Pesanan Perbaikan
// @Description Menginisiasi pesanan servis elektronik baru oleh Konsumen. Bisa diikat langsung ke teknisi atau dibiarkan kosong untuk pesanan publik.
// @Tags Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateOrderPayload true "Data Kerusakan dan Booking Opsional"
// @Success 201 {object} map[string]interface{}
// @Router /api/orders/ [post]
func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var req CreateOrderPayload

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak ditemukan pada token"})
	}

	customerID, err := uuid.Parse(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format User ID tidak valid"})
	}

	var techIDPtr *uuid.UUID
	if req.TechnicianID != "" {
		tID, errParse := uuid.Parse(req.TechnicianID)
		if errParse != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format Technician ID tidak valid"})
		}
		techIDPtr = &tID
	}

	order := domain.Order{
		CustomerID:         customerID,
		TechnicianID:       techIDPtr,
		DeviceCategory:     req.DeviceCategory,
		ProblemDescription: req.ProblemDescription,
		CustomerLatitude:   req.CustomerLatitude,
		CustomerLongitude:  req.CustomerLongitude,
	}

	if err := h.orderUsecase.CreateOrder(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pesanan layanan berhasil dibuat",
		"data":    order,
	})
}

// @Summary Menyelesaikan Transaksi & Kalkulasi Emisi (Anti-Fraud)
// @Description Menyelesaikan pesanan dengan memvalidasi foto bukti dan koordinat GPS (Metrik EPA WARM).
// @Tags Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param id path string true "ID Pesanan (UUID)"
// @Param request body usecase.CompleteOrderRequest true "Data Penyelesaian"
// @Success 200 {object} map[string]interface{}
// @Router /api/orders/{id}/complete [put]
func (h *OrderHandler) Complete(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var req usecase.CompleteOrderRequest

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.orderUsecase.CompleteOrder(orderID, req); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Transaksi berhasil diselesaikan. Protokol Anti-Fraud tersimpan dan Kalkulasi EPA WARM telah diinjeksikan.",
	})
}

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
	DeviceCategory     string `json:"device_category" example:"Pendingin & Komersial"`
	ProblemDescription string `json:"problem_description" example:"Kompresor mati dan berasap"`
}

func NewOrderHandler(app *fiber.App, usecase usecase.OrderUsecase) {
	handler := &OrderHandler{orderUsecase: usecase}

	api := app.Group("/api/orders")
	api.Post("/", middleware.Protected(), handler.Create)
	api.Put("/:id/complete", middleware.Protected(), handler.Complete)
}

// @Summary Membuat Pesanan Perbaikan
// @Description Menginisiasi pesanan servis elektronik baru oleh Konsumen.
// @Tags Orders
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body CreateOrderPayload true "Data Kerusakan"
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

	order := domain.Order{
		CustomerID:         customerID,
		DeviceCategory:     req.DeviceCategory,
		ProblemDescription: req.ProblemDescription,
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

package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/domain"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type OrderHandler struct {
	orderUsecase usecase.OrderUsecase
}

func NewOrderHandler(app *fiber.App, usecase usecase.OrderUsecase) {
	handler := &OrderHandler{orderUsecase: usecase}

	api := app.Group("/api/orders")

	api.Post("/", middleware.Protected(), handler.Create)
	api.Put("/:id/complete", middleware.Protected(), handler.Complete)
}

func (h *OrderHandler) Create(c *fiber.Ctx) error {
	var order domain.Order
	if err := c.BodyParser(&order); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.orderUsecase.CreateOrder(&order); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pesanan layanan berhasil dibuat",
		"data":    order,
	})
}

func (h *OrderHandler) Complete(c *fiber.Ctx) error {
	orderID := c.Params("id")

	var req struct {
		DeviceCategory string `json:"device_category"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.orderUsecase.CompleteOrder(orderID, req.DeviceCategory); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Pesanan selesai, metrik E-Waste diperbarui"})
}

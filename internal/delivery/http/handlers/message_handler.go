package handlers

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/domain"
)

type MessageHandler struct {
	messageUsecase domain.MessageUsecase
}

type SendMessagePayload struct {
	Content string `json:"content" example:"Halo teknisi, posisi di mana ya?"`
}

func NewMessageHandler(app *fiber.App, usecase domain.MessageUsecase) {
	handler := &MessageHandler{
		messageUsecase: usecase,
	}

	api := app.Group("/api/orders")
	api.Post("/:orderId/messages", middleware.Protected(), handler.SendMessage)
	api.Get("/:orderId/messages", middleware.Protected(), handler.GetOrderMessages)
}

// @Summary Mengirim pesan ke dalam room order
// @Description Endpoint untuk customer dan teknisi mengirim pesan chat. Realtime di-handle oleh Supabase di frontend.
// @Tags Chat
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param orderId path string true "ID Pesanan (UUID)"
// @Param request body SendMessagePayload true "Data Pesan"
// @Success 201 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/orders/{orderId}/messages [post]
func (h *MessageHandler) SendMessage(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	userIDRaw := c.Locals("user_id")
	roleRaw := c.Locals("role")

	if userIDRaw == nil || roleRaw == nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak: Sesi tidak valid"})
	}

	userID := fmt.Sprintf("%v", userIDRaw)
	role := fmt.Sprintf("%v", roleRaw)

	var req SendMessagePayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	msg, err := h.messageUsecase.SendMessage(orderID, userID, role, req.Content)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pesan berhasil dikirim",
		"data":    msg,
	})
}

// @Summary Mengambil histori pesan
// @Description Mengambil riwayat percakapan berdasarkan Order ID
// @Tags Chat
// @Produce json
// @Security ApiKeyAuth
// @Param orderId path string true "ID Pesanan (UUID)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/orders/{orderId}/messages [get]
func (h *MessageHandler) GetOrderMessages(c *fiber.Ctx) error {
	orderID := c.Params("orderId")

	messages, err := h.messageUsecase.GetOrderMessages(orderID)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil riwayat pesan",
		"data":    messages,
	})
}

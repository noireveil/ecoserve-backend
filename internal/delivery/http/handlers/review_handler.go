package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type ReviewHandler struct {
	reviewUsecase usecase.ReviewUsecase
}

type ReviewRequestPayload struct {
	Rating  int    `json:"rating" example:"5"`
	Comment string `json:"comment" example:"Teknisi sangat profesional dan jujur!"`
}

func NewReviewHandler(app *fiber.App, usecase usecase.ReviewUsecase) {
	handler := &ReviewHandler{reviewUsecase: usecase}

	api := app.Group("/api/reviews")
	api.Post("/order/:order_id", middleware.Protected(), handler.CreateReview)
}

// @Summary Memberikan Ulasan Teknisi
// @Description Konsumen memberikan rating (1-5) dan komentar untuk teknisi setelah pesanan selesai. Sistem akan otomatis menghitung ulang rata-rata rating teknisi dengan fitur keamanan Anti-Spam.
// @Tags Reviews
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param order_id path string true "ID Pesanan (UUID)"
// @Param request body ReviewRequestPayload true "Data Ulasan"
// @Success 201 {object} map[string]string "Berhasil memberikan ulasan"
// @Failure 400 {object} map[string]string "Error validasi (rating tidak valid, pesanan belum selesai, dll)"
// @Failure 401 {object} map[string]string "Akses ditolak (belum login atau token tidak valid)"
// @Router /api/reviews/order/{order_id} [post]
func (h *ReviewHandler) CreateReview(c *fiber.Ctx) error {
	orderID := c.Params("order_id")
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User ID tidak valid"})
	}

	var req ReviewRequestPayload
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	payload := usecase.SubmitReviewPayload{
		OrderID: orderID,
		UserID:  userIDStr,
		Rating:  req.Rating,
		Comment: req.Comment,
	}

	if err := h.reviewUsecase.SubmitReview(payload); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Ulasan berhasil dikirim. Terima kasih telah menggunakan EcoServe!",
	})
}

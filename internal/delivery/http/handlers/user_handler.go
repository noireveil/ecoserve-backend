package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
	"github.com/noireveil/ecoserve-backend/pkg/utils"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, usecase usecase.UserUsecase) {
	handler := &UserHandler{userUsecase: usecase}

	api := app.Group("/api/users")
	api.Post("/auth/request", handler.RequestOTP)
	api.Post("/auth/verify", handler.VerifyOTP)

	api.Get("/me", middleware.Protected(), handler.GetProfile)
}

func (h *UserHandler) RequestOTP(c *fiber.Ctx) error {
	var req struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.userUsecase.RequestOTP(req.FullName, req.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kode OTP telah dikirim melalui email"})
}

func (h *UserHandler) VerifyOTP(c *fiber.Ctx) error {
	var req struct {
		Email string `json:"email"`
		Code  string `json:"code"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	user, err := h.userUsecase.VerifyOTP(req.Email, req.Code)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghasilkan token"})
	}

	cookie := new(fiber.Cookie)
	cookie.Name = "jwt"
	cookie.Value = token
	cookie.Expires = time.Now().Add(24 * time.Hour)
	cookie.HTTPOnly = true
	cookie.SameSite = "None"
	cookie.Secure = true

	c.Cookie(cookie)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Autentikasi berhasil",
		"data":    user,
	})
}

func (h *UserHandler) GetProfile(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak: Sesi tidak valid"})
	}

	user, err := h.userUsecase.GetUserByID(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Pengguna tidak ditemukan"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data": user,
	})
}

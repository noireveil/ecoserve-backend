package handlers

import (
	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
	"github.com/noireveil/ecoserve-backend/pkg/utils"
)

type UserHandler struct {
	userUsecase usecase.UserUsecase
}

func NewUserHandler(app *fiber.App, usecase usecase.UserUsecase) {
	handler := &UserHandler{userUsecase: usecase}

	api := app.Group("/api/users")
	api.Post("/auth", handler.LoginOrRegister)
}

func (h *UserHandler) LoginOrRegister(c *fiber.Ctx) error {
	var req struct {
		FullName string `json:"full_name"`
		WhatsApp string `json:"whatsapp"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	user, err := h.userUsecase.LoginOrRegister(req.FullName, req.WhatsApp)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	token, err := utils.GenerateToken(user.ID, user.Role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal menghasilkan token otentikasi"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Autentikasi berhasil",
		"token":   token,
		"data":    user,
	})
}

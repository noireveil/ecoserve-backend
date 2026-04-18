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

type RequestOTPPayload struct {
	FullName string `json:"full_name" example:"EcoServe Tester"`
	Email    string `json:"email" example:"eco.security@ecoserve.com"`
}

type VerifyOTPPayload struct {
	Email string `json:"email" example:"ecosecurity@ecoserve.com"`
	Code  string `json:"code" example:"123456"`
}

func NewUserHandler(app *fiber.App, usecase usecase.UserUsecase) {
	handler := &UserHandler{userUsecase: usecase}

	api := app.Group("/api/users")
	api.Post("/auth/request", handler.RequestOTP)
	api.Post("/auth/verify", handler.VerifyOTP)
	api.Post("/auth/logout", handler.Logout)

	api.Get("/me", middleware.Protected(), handler.GetProfile)
	api.Get("/me/impact", middleware.Protected(), handler.GetImpact)
	api.Delete("/me", middleware.Protected(), handler.DeleteAccount)
}

// @Summary Meminta Kode OTP
// @Description Mengirimkan kode OTP ke email pengguna untuk proses otentikasi.
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body RequestOTPPayload true "Data Pengguna"
// @Success 200 {object} map[string]interface{}
// @Router /api/users/auth/request [post]
func (h *UserHandler) RequestOTP(c *fiber.Ctx) error {
	var req RequestOTPPayload

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if err := h.userUsecase.RequestOTP(req.FullName, req.Email); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{"message": "Kode OTP telah dikirim melalui email"})
}

// @Summary Verifikasi Kode OTP
// @Description Memverifikasi OTP dan mengeluarkan token JWT di dalam Cookie (Sesi Aktif).
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body VerifyOTPPayload true "Data OTP"
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/users/auth/verify [post]
func (h *UserHandler) VerifyOTP(c *fiber.Ctx) error {
	var req VerifyOTPPayload

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

// @Summary Mendapatkan Profil Pengguna
// @Description Mengambil data profil dari pengguna yang sedang login berdasarkan JWT.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/users/me [get]
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

// @Summary Menghapus Akun Pengguna (Soft Delete)
// @Description Menghapus akun dan mencabut sesi (logout). Data tidak dihapus permanen untuk menjaga integritas relasi EPA WARM.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/users/me [delete]
func (h *UserHandler) DeleteAccount(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak: Sesi tidak valid"})
	}

	if err := h.userUsecase.DeleteAccount(userIDStr); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memproses penghapusan akun"})
	}

	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Akun berhasil dihapus. Sesi telah diakhiri.",
	})
}

// @Summary Logout Pengguna
// @Description Mencabut sesi pengguna dengan menghapus cookie JWT dari browser.
// @Tags Auth
// @Produce json
// @Success 200 {object} map[string]interface{}
// @Router /api/users/auth/logout [post]
func (h *UserHandler) Logout(c *fiber.Ctx) error {
	c.Cookie(&fiber.Cookie{
		Name:     "jwt",
		Value:    "",
		Expires:  time.Now().Add(-1 * time.Hour),
		HTTPOnly: true,
		SameSite: "None",
		Secure:   true,
	})

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil logout. Sesi telah dihapus dari sistem.",
	})
}

// @Summary Mendapatkan Dampak Lingkungan Pengguna
// @Description Mengambil agregasi data dampak lingkungan (Total Perbaikan dan Total CO2 yang Dihindari) untuk konsumen yang sedang login.
// @Tags Users
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} map[string]interface{}
// @Failure 401 {object} map[string]interface{}
// @Router /api/users/me/impact [get]
func (h *UserHandler) GetImpact(c *fiber.Ctx) error {
	userIDStr, ok := c.Locals("user_id").(string)
	if !ok || userIDStr == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akses ditolak: Sesi tidak valid"})
	}

	impact, err := h.userUsecase.GetImpact(userIDStr)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Berhasil mengambil metrik dampak lingkungan",
		"data":    impact,
	})
}

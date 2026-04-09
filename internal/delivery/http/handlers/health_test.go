package handlers_test

import (
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func TestHealthCheck(t *testing.T) {
	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status": "success",
		})
	})

	req := httptest.NewRequest("GET", "/health", nil)

	resp, err := app.Test(req, -1)

	if err != nil {
		t.Fatalf("Gagal mengeksekusi permintaan: %v", err)
	}

	if resp.StatusCode != 200 {
		t.Errorf("Ekspektasi kode status 200, didapatkan %d", resp.StatusCode)
	}
}

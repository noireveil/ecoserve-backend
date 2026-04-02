package handlers

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

type ChatbotHandler struct{}

func NewChatbotHandler(app *fiber.App) {
	handler := &ChatbotHandler{}
	api := app.Group("/api/chatbot")
	api.Post("/triage", handler.Triage)
}

func (h *ChatbotHandler) Triage(c *fiber.Ctx) error {
	var req struct {
		Message string `json:"message"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format permintaan tidak valid"})
	}

	if req.Message == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Pesan keluhan tidak boleh kosong"})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Konfigurasi AI belum diatur"})
	}

	reply, err := h.callGeminiAPI(apiKey, req.Message)
	if err != nil {
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Gagal memproses diagnosis AI"})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Diagnosis berhasil dianalisis",
		"data": map[string]string{
			"reply": reply,
		},
	})
}

func (h *ChatbotHandler) callGeminiAPI(apiKey, userMessage string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	systemInstruction := `Berperanlah sebagai AI Diagnostik Teknis Senior di EcoServe, sebuah platform perbaikan elektronik dan manajemen siklus hidup E-Waste.
Karakteristik respons harus profesional, analitis, definitif, dan berwibawa. Hindari penggunaan bahasa yang ragu-ragu (seperti "mungkin saja" atau "ada kemungkinan").

Instruksi Algoritmik Pemrosesan Keluhan:
1. Analisis Gejala: Identifikasi akar masalah secara langsung. Sebutkan nama komponen spesifik di dalam perangkat yang mengalami kegagalan teknis berdasarkan gejala yang dilaporkan.
2. Protokol Keselamatan: Instruksikan 1-2 langkah mitigasi darurat yang sangat konkret untuk mencegah kerusakan lebih lanjut atau risiko keselamatan (contoh: isolasi sumber air, pemutusan arus listrik).
3. Estimasi Perbaikan: Berikan proyeksi rentang biaya perbaikan yang wajar dan realistis di pasar Indonesia (dalam Rupiah), disertai catatan singkat bahwa nilai tersebut adalah estimasi pra-inspeksi.
4. Resolusi Kategori: Respons DIWAJIBKAN berakhir dengan arahan penugasan teknisi. Pilih SATU kategori yang paling relevan dari daftar parameter ini: ["Pendingin & Komersial", "Home Appliances", "IT & Gadget"]. Cetak tebal (bold) nama kategori tersebut.

Aturan Pemformatan:
- Langsung berikan analisis teknis tanpa salam pembuka yang bertele-tele (hindari kalimat seperti "Halo, terima kasih telah menghubungi...").
- Gunakan struktur paragraf yang padat, atau gunakan poin-poin (bullet points) jika menjelaskan protokol keselamatan.
- Jaga struktur respons agar tetap proporsional, mudah dipindai (scannable), dan optimal saat dibaca melalui antarmuka layar lebar (desktop) maupun perangkat bergerak.`

	payload := map[string]interface{}{
		"system_instruction": map[string]interface{}{
			"parts": map[string]string{
				"text": systemInstruction,
			},
		},
		"contents": []map[string]interface{}{
			{
				"parts": []map[string]string{
					{"text": userMessage},
				},
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.4,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Post(url, "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", errors.New("upstream api error")
	}

	var geminiResp struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&geminiResp); err != nil {
		return "", err
	}

	if len(geminiResp.Candidates) > 0 && len(geminiResp.Candidates[0].Content.Parts) > 0 {
		return geminiResp.Candidates[0].Content.Parts[0].Text, nil
	}

	return "", errors.New("format respons tidak dikenali")
}

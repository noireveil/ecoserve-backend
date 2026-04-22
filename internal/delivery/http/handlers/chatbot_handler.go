package handlers

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/noireveil/ecoserve-backend/internal/delivery/http/middleware"
	"github.com/noireveil/ecoserve-backend/internal/usecase"
)

type ChatbotHandler struct {
	techUsecase usecase.TechnicianUsecase
}

type GeminiTriageResponse struct {
	Analysis        string  `json:"analysis"`
	Mitigation      string  `json:"mitigation"`
	Category        string  `json:"category"`
	ConfidenceScore float64 `json:"confidence_score"`
	IsDIYEligible   bool    `json:"is_diy_eligible"`
}

type TriagePayload struct {
	Message   string  `json:"message" example:"Kulkas saya bocor dan berbunyi bising"`
	Photo     string  `json:"photo" example:"base64_encoded_string_or_url"`
	Longitude float64 `json:"longitude" example:"106.8229"`
	Latitude  float64 `json:"latitude" example:"-6.1944"`
}

func NewChatbotHandler(app *fiber.App, usecase usecase.TechnicianUsecase) {
	handler := &ChatbotHandler{techUsecase: usecase}

	api := app.Group("/api/chatbot")
	api.Post("/triage", middleware.Protected(), handler.Triage)
}

// @Summary AI Triage Keluhan Elektronik (Multimodal)
// @Description Mengirimkan keluhan kerusakan elektronik dan foto ke Gemini AI Vision untuk dianalisis.
// @Tags Chatbot
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Param request body TriagePayload true "Data Keluhan dan Foto"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Router /api/chatbot/triage [post]
func (h *ChatbotHandler) Triage(c *fiber.Ctx) error {
	var req TriagePayload

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Format JSON tidak valid"})
	}

	if req.Message == "" && req.Photo == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Pesan keluhan atau foto tidak boleh kosong"})
	}

	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Konfigurasi AI belum diatur"})
	}

	replyJSON, err := h.callGeminiAPI(apiKey, req.Message, req.Photo)
	if err != nil {
		log.Printf("[Gemini API Error]: %v\n", err)
		return c.Status(fiber.StatusBadGateway).JSON(fiber.Map{"error": "Gagal memproses diagnosis AI"})
	}

	if len(replyJSON) > 7 && replyJSON[:7] == "```json" {
		replyJSON = replyJSON[7 : len(replyJSON)-3]
	}

	var aiResp GeminiTriageResponse
	if err := json.Unmarshal([]byte(replyJSON), &aiResp); err != nil {
		log.Printf("[JSON Parse Error]: Gagal memetakan respons: %s\n", replyJSON)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memetakan respons JSON dari AI"})
	}

	if aiResp.ConfidenceScore < 0.85 || !aiResp.IsDIYEligible {
		technicians, _ := h.techUsecase.GetNearbyTechnicians(req.Longitude, req.Latitude, 15)
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Risiko terlalu tinggi untuk perbaikan mandiri. Mengaktifkan Fallback Pencarian Teknisi.",
			"data": map[string]interface{}{
				"diagnosis":          aiResp,
				"is_fallback_active": true,
				"technicians":        technicians,
			},
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message": "Diagnosis Swa-Perbaikan (DIY) berhasil divalidasi",
		"data": map[string]interface{}{
			"diagnosis":          aiResp,
			"is_fallback_active": false,
			"technicians":        nil,
		},
	})
}

func (h *ChatbotHandler) callGeminiAPI(apiKey, userMessage, photoData string) (string, error) {
	url := "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.5-flash:generateContent?key=" + apiKey

	systemInstruction := `Berperanlah sebagai AI Diagnostik Teknis Senior.
Anda WAJIB memberikan respons DALAM FORMAT JSON MURNI tanpa teks pembuka atau penutup.
Gunakan struktur skema JSON berikut:
{
  "analysis": "Penjelasan teknis komponen yang rusak berdasarkan deskripsi atau gambar",
  "mitigation": "Langkah darurat keamanan operasional",
  "category": "Pilih satu: [Pendingin & Komersial, Home Appliances, IT & Gadget, Perangkat Lainnya]",
  "confidence_score": <angka desimal antara 0.00 hingga 1.00 tingkat kepastian diagnosis>,
  "is_diy_eligible": <boolean true jika aman diperbaiki konsumen, false jika butuh alat khusus/berisiko>
}`

	var userParts []map[string]interface{}

	if userMessage == "" {
		userMessage = "Tolong analisis kerusakan pada perangkat di gambar ini dan berikan diagnosis teknisnya."
	}

	userParts = append(userParts, map[string]interface{}{
		"text": userMessage,
	})

	if photoData != "" {
		mimeType := "image/jpeg"
		base64String := photoData

		if strings.HasPrefix(photoData, "data:image/") {
			parts := strings.SplitN(photoData, ";base64,", 2)
			if len(parts) == 2 {
				mimeType = strings.TrimPrefix(parts[0], "data:")
				base64String = parts[1]
			}
		} else if strings.HasPrefix(photoData, "http://") || strings.HasPrefix(photoData, "https://") {
			imgBytes, err := fetchImageFromURL(photoData)
			if err == nil {
				base64String = base64.StdEncoding.EncodeToString(imgBytes)
				mimeType = detectMimeType(imgBytes)
			} else {
				log.Printf("[Image Fetch Error]: %v\n", err)
			}
		}

		if base64String != "" {
			cleanBase64 := strings.ReplaceAll(base64String, "\n", "")
			cleanBase64 = strings.ReplaceAll(cleanBase64, " ", "")

			userParts = append(userParts, map[string]interface{}{
				"inlineData": map[string]string{
					"mimeType": mimeType,
					"data":     cleanBase64,
				},
			})
		}
	}

	payload := map[string]interface{}{
		"systemInstruction": map[string]interface{}{
			"parts": []map[string]string{
				{
					"text": systemInstruction,
				},
			},
		},
		"contents": []map[string]interface{}{
			{
				"parts": userParts,
			},
		},
		"generationConfig": map[string]interface{}{
			"temperature": 0.2,
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return "", err
	}

	client := &http.Client{Timeout: 30 * time.Second}
	maxRetries := 3
	var resp *http.Response

	for attempt := 0; attempt < maxRetries; attempt++ {
		resp, err = client.Post(url, "application/json", bytes.NewBuffer(jsonData))
		if err != nil {
			return "", err
		}

		if resp.StatusCode == http.StatusOK {
			break
		}

		if resp.StatusCode == http.StatusServiceUnavailable || resp.StatusCode == http.StatusTooManyRequests {
			_ = resp.Body.Close()
			jedaWaktu := time.Duration(1<<attempt) * time.Second
			log.Printf("[Gemini API] Server sibuk (%d). Mencoba ulang (%d/%d) dalam %v...\n", resp.StatusCode, attempt+1, maxRetries, jedaWaktu)
			time.Sleep(jedaWaktu)
			continue
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		_ = resp.Body.Close()
		return "", fmt.Errorf("upstream merespons dengan kode %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp == nil || resp.StatusCode != http.StatusOK {
		return "", errors.New("gagal menghubungi Gemini API setelah melewati batas maksimal percobaan ulang")
	}
	defer resp.Body.Close()

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

	return "", errors.New("format respons tidak dikenali dari Google API")
}

func fetchImageFromURL(url string) ([]byte, error) {
	client := &http.Client{Timeout: 10 * time.Second}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status code %d", resp.StatusCode)
	}

	return io.ReadAll(resp.Body)
}

func detectMimeType(data []byte) string {
	mimeType := http.DetectContentType(data)
	if mimeType == "application/octet-stream" {
		return "image/jpeg"
	}
	return mimeType
}

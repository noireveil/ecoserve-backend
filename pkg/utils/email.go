package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
)

func SendEmailOTP(targetEmail string, code string) error {
	apiKey := os.Getenv("MAILJET_API_KEY")
	secretKey := os.Getenv("MAILJET_SECRET_KEY")
	senderEmail := os.Getenv("MAILJET_SENDER_EMAIL")

	if apiKey == "" || secretKey == "" || senderEmail == "" {
		return errors.New("kredensial Mailjet tidak terdefinisi pada variabel lingkungan")
	}

	htmlBody := fmt.Sprintf(`
<!DOCTYPE html>
<html lang="id">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; background-color: #f4f7f6; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; -webkit-font-smoothing: antialiased;">
    <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="background-color: #f4f7f6; padding: 50px 20px;">
        <tr>
            <td align="center">
                <table border="0" cellpadding="0" cellspacing="0" width="100%%" style="max-width: 480px; background-color: #ffffff; border-radius: 16px; overflow: hidden; box-shadow: 0 10px 25px rgba(0, 0, 0, 0.05);">

                    <tr>
                        <td align="center" style="background: #059669; background: linear-gradient(135deg, #059669 0%%, #34d399 100%%); padding: 45px 20px;">
                            <h1 style="margin: 0; color: #ffffff; font-size: 26px; font-weight: 800; letter-spacing: -0.5px;">EcoServe</h1>
                            <p style="margin: 8px 0 0 0; color: rgba(255, 255, 255, 0.9); font-size: 13px; letter-spacing: 3px; text-transform: uppercase;">Kode Keamanan</p>
                        </td>
                    </tr>

                    <tr>
                        <td style="padding: 40px 35px;">
                            <p style="margin: 0 0 20px 0; font-size: 16px; color: #1f2937; line-height: 1.6;">
                                Halo,
                            </p>
                            <p style="margin: 0 0 30px 0; font-size: 16px; color: #4b5563; line-height: 1.6;">
                                Seseorang mencoba masuk ke akun EcoServe Anda. Gunakan kode verifikasi di bawah ini untuk melanjutkan proses:
                            </p>

                            <table border="0" cellpadding="0" cellspacing="0" width="100%%">
                                <tr>
                                    <td align="center" style="background-color: #ecfdf5; border-radius: 12px; padding: 25px;">
                                        <span style="font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 42px; font-weight: 700; color: #047857; letter-spacing: 8px;">%s</span>
                                    </td>
                                </tr>
                            </table>

                            <p style="margin: 30px 0 0 0; font-size: 14px; color: #6b7280; line-height: 1.6;">
                                Kode ini akan kedaluwarsa dalam waktu <strong>5 menit</strong>. Jika Anda tidak merasa melakukan permintaan untuk masuk, Anda dapat mengabaikan surel ini dengan aman.
                            </p>
                        </td>
                    </tr>

                    <tr>
                        <td align="center" style="background-color: #f9fafb; padding: 25px; border-top: 1px solid #f3f4f6;">
                            <p style="margin: 0; font-size: 12px; color: #9ca3af; line-height: 1.5;">
                                &copy; 2026 EcoServe.<br>Platform Ekonomi Sirkular.
                            </p>
                        </td>
                    </tr>

                </table>
            </td>
        </tr>
    </table>
</body>
</html>`, code)

	payload := map[string]interface{}{
		"Messages": []map[string]interface{}{
			{
				"From": map[string]string{
					"Email": senderEmail,
					"Name":  "EcoServe Security",
				},
				"To": []map[string]string{
					{
						"Email": targetEmail,
					},
				},
				"Subject":  "Kode Autentikasi EcoServe",
				"HTMLPart": htmlBody,
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "https://api.mailjet.com/v3.1/send", bytes.NewBuffer(jsonData))
	if err != nil {
		return err
	}

	req.SetBasicAuth(apiKey, secretKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("gagal mengirim surel via Mailjet: status %d - %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

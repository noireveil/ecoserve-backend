package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
)

func buildDigitBoxes(code string) string {
	var sb strings.Builder
	for _, ch := range code {
		sb.WriteString(fmt.Sprintf(`
              <span style="display:inline-block; width:48px; height:56px; line-height:56px; background-color:#F9FAFB; border:1px solid #D1D5DB; border-radius:8px; font-family:ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size:28px; font-weight:700; color:#111827; text-align:center; margin:0 4px;">%c</span>`, ch))
	}
	return sb.String()
}

func SendEmailOTP(targetEmail string, name string, code string) error {
	apiKey := os.Getenv("MAILJET_API_KEY")
	secretKey := os.Getenv("MAILJET_SECRET_KEY")
	senderEmail := os.Getenv("MAILJET_SENDER_EMAIL")

	if apiKey == "" || secretKey == "" || senderEmail == "" {
		return errors.New("kredensial Mailjet tidak terdefinisi pada variabel lingkungan")
	}
	if name == "" {
		name = "Pengguna"
	}

	digitBoxes := buildDigitBoxes(code)

	htmlBody := fmt.Sprintf(`<!DOCTYPE html>
<html lang="id">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width,initial-scale=1.0">
  <title>Kode Verifikasi EcoServe</title>
</head>
<body style="margin:0;padding:0;background-color:#F3F4F6;
             font-family:-apple-system,BlinkMacSystemFont,'Segoe UI',Roboto,Helvetica,Arial,sans-serif;">

  <table border="0" cellpadding="0" cellspacing="0" width="100%%"
         style="background-color:#F3F4F6;padding:60px 20px;">
    <tr>
      <td align="center">
        <table border="0" cellpadding="0" cellspacing="0" width="100%%"
               style="max-width:500px;background-color:#FFFFFF;border-radius:12px;box-shadow:0 4px 6px -1px rgba(0, 0, 0, 0.05), 0 2px 4px -1px rgba(0, 0, 0, 0.03);overflow:hidden;">

          <tr>
            <td align="center" style="padding:40px 40px 20px;">
              <span style="font-size:24px;font-weight:800;letter-spacing:-0.5px;">
                <span style="color:#00C896;">Eco</span><span style="color:#111827;">Serve</span>
              </span>
            </td>
          </tr>

          <tr>
            <td style="padding:20px 40px 40px;">
              <p style="margin:0 0 8px;font-size:16px;font-weight:600;color:#111827;">
                Halo %s,
              </p>
              <p style="margin:0 0 32px;font-size:15px;color:#4B5563;line-height:1.6;">
                Seseorang sedang mencoba mengakses akun EcoServe Anda. Gunakan kode keamanan di bawah ini untuk melanjutkan proses verifikasi:
              </p>

              <table border="0" cellpadding="0" cellspacing="0" width="100%%">
                <tr>
                  <td align="center" style="padding: 10px 0;">
                    %s
                  </td>
                </tr>
              </table>

              <p style="margin:32px 0 0;font-size:14px;color:#6B7280;text-align:center;">
                Kode ini berlaku selama <strong style="color:#111827;">5 menit</strong>.
              </p>

              <div style="margin-top:32px;padding-top:24px;border-top:1px solid #E5E7EB;">
                <p style="margin:0;font-size:13px;color:#9CA3AF;line-height:1.5;text-align:center;">
                  Jika Anda tidak merasa melakukan permintaan ini, abaikan surel ini. Kata sandi dan akun Anda tetap aman.
                </p>
              </div>

            </td>
          </tr>

          <tr>
            <td style="background-color:#F9FAFB;padding:24px 40px;">
              <table border="0" cellpadding="0" cellspacing="0" width="100%%">
                <tr>
                  <td align="center">
                    <p style="margin:0;font-size:12px;color:#9CA3AF;">
                      &copy; 2026 EcoServe. Platform Ekonomi Sirkular.
                    </p>
                  </td>
                </tr>
              </table>
            </td>
          </tr>

        </table>
      </td>
    </tr>
  </table>

</body>
</html>`, name, digitBoxes)

	payload := map[string]interface{}{
		"Messages": []map[string]interface{}{
			{
				"From": map[string]string{
					"Email": senderEmail,
					"Name":  "EcoServe Security",
				},
				"To": []map[string]string{
					{"Email": targetEmail, "Name": name},
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
		return fmt.Errorf("gagal mengirim surel via Mailjet: status %d - %s",
			resp.StatusCode, string(bodyBytes))
	}

	return nil
}

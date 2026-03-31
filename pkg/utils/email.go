package utils

import (
	"fmt"
	"net/smtp"
	"os"
)

func SendEmailOTP(targetEmail string, code string) error {
	host := os.Getenv("SMTP_HOST")
	port := os.Getenv("SMTP_PORT")
	user := os.Getenv("SMTP_USER")
	pass := os.Getenv("SMTP_PASS")

	from := fmt.Sprintf("From: EcoServe <%s>\r\n", user)
	to := fmt.Sprintf("To: %s\r\n", targetEmail)
	subject := "Subject: Kode Verifikasi EcoServe\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
		<div style="font-family: sans-serif; max-width: 500px; border: 1px solid #eee; padding: 20px;">
			<h2 style="color: #2D3748;">EcoServe Verification</h2>
			<p>Halo,</p>
			<p>Gunakan kode di bawah ini untuk memverifikasi pendaftaran atau masuk ke sistem EcoServe:</p>
			<div style="background: #F7FAFC; padding: 15px; text-align: center; font-size: 24px; font-weight: bold; letter-spacing: 5px; color: #3182CE; border-radius: 5px; margin: 20px 0;">
				%s
			</div>
			<p style="font-size: 13px; color: #718096; margin-top: 20px;">
				Kode ini berlaku selama 5 menit. Jika aktivitas ini tidak dikenali, abaikan surel ini.
			</p>
		</div>`, code)

	message := []byte(from + to + subject + mime + body)
	auth := smtp.PlainAuth("", user, pass, host)

	return smtp.SendMail(host+":"+port, auth, user, []string{targetEmail}, message)
}

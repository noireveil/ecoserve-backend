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

	from := fmt.Sprintf("From: EcoServe Security <%s>\r\n", user)
	to := fmt.Sprintf("To: %s\r\n", targetEmail)
	subject := "Subject: Kode Autentikasi EcoServe Anda\r\n"
	mime := "MIME-version: 1.0;\nContent-Type: text/html; charset=\"UTF-8\";\n\n"

	body := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<meta name="viewport" content="width=device-width, initial-scale=1.0">
</head>
<body style="margin: 0; padding: 0; font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, Helvetica, Arial, sans-serif; background-color: #f6f9fc; color: #333333;">
	<table width="100%%" border="0" cellspacing="0" cellpadding="0" style="background-color: #f6f9fc; padding: 40px 0;">
		<tr>
			<td align="center">
				<table width="100%%" border="0" cellspacing="0" cellpadding="0" style="max-width: 500px; background-color: #ffffff; border-radius: 8px; box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05); overflow: hidden;">
					<tr>
						<td style="padding: 40px 40px 20px 40px; text-align: center; border-bottom: 1px solid #edf2f7;">
							<h1 style="margin: 0; font-size: 24px; font-weight: 600; color: #1a202c; letter-spacing: -0.5px;">EcoServe</h1>
						</td>
					</tr>
					<tr>
						<td style="padding: 30px 40px;">
							<p style="margin: 0 0 20px 0; font-size: 15px; line-height: 1.6; color: #4a5568;">
								Halo,
							</p>
							<p style="margin: 0 0 24px 0; font-size: 15px; line-height: 1.6; color: #4a5568;">
								Berikut adalah kode autentikasi satu waktu (OTP) Anda untuk masuk ke platform EcoServe. Kode ini hanya berlaku selama <strong>5 menit</strong>.
							</p>
							<div style="background-color: #f8fafc; border: 1px solid #e2e8f0; border-radius: 6px; padding: 16px; text-align: center; margin-bottom: 24px;">
								<span style="font-family: ui-monospace, SFMono-Regular, Menlo, Monaco, Consolas, monospace; font-size: 32px; font-weight: 700; letter-spacing: 8px; color: #000000;">%s</span>
							</div>
							<p style="margin: 0; font-size: 14px; line-height: 1.6; color: #718096;">
								Jika Anda tidak merasa melakukan permintaan ini, Anda dapat mengabaikan surel ini dengan aman.
							</p>
						</td>
					</tr>
					<tr>
						<td style="padding: 20px 40px; background-color: #f8fafc; text-align: center; border-top: 1px solid #edf2f7;">
							<p style="margin: 0; font-size: 12px; color: #a0aec0;">
								&copy; 2026 EcoServe. Mendukung Ekonomi Sirkular.
							</p>
						</td>
					</tr>
				</table>
			</td>
		</tr>
	</table>
</body>
</html>`, code)

	message := []byte(from + to + subject + mime + body)
	auth := smtp.PlainAuth("", user, pass, host)

	return smtp.SendMail(host+":"+port, auth, user, []string{targetEmail}, message)
}

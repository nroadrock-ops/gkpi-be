package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func SendOTPEmail(to, name, code string) error {
	resendAPIKey := os.Getenv("RESEND_API_KEY")
	emailFrom := os.Getenv("EMAIL_FROM")

	if resendAPIKey == "" || resendAPIKey == "re_123456789" {
		// Logika ini berguna untuk testing lokal tanpa key sungguhan
		fmt.Printf("\n========================================\n")
		fmt.Printf("[MOCK EMAIL] OTP to %s: %s\n", to, code)
		fmt.Printf("========================================\n\n")
		return nil
	}

	if emailFrom == "" {
		emailFrom = "noreply@gkpicimahi.org"
	}

	htmlContent := fmt.Sprintf(`
		<div style="font-family: Arial, sans-serif; padding: 20px;">
			<h2>Verifikasi Akun GKPI Cimahi</h2>
			<p>Halo %s,</p>
			<p>Kode OTP Anda adalah: <strong><span style="font-size: 24px; color: #2D3748;">%s</span></strong></p>
			<p>Kode ini akan kadaluarsa dalam waktu singkat. JANGAN berikan kode ini kepada siapapun.</p>
			<p>Terima kasih,<br/>Tim GKPI Cimahi</p>
		</div>
	`, name, code)

	payload := map[string]interface{}{
		"from":    emailFrom,
		"to":      []string{to},
		"subject": "Kode OTP Registrasi GKPI Cimahi",
		"html":    htmlContent,
	}

	payloadBytes, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", "https://api.resend.com/emails", bytes.NewBuffer(payloadBytes))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+resendAPIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("\n[FALLBACK OTP] Failed to call Resend API. OTP to %s: %s\n\n", to, code)
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 400 {
		fmt.Printf("\n[FALLBACK OTP] Resend API error (status %d). OTP to %s: %s\n\n", resp.StatusCode, to, code)
		return fmt.Errorf("failed to send email via resend: status %d", resp.StatusCode)
	}

	return nil
}

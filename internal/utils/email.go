package utils

import (
	"fmt"
	"log"
)

// SendEmailOTP mensimulasikan pengiriman email otomatis.
// Di production, bisa diganti dengan implementasi SMTP/Sendgrid.
func SendEmailOTP(toEmail, otp string) error {
	log.Printf("[EMAIL TERKIRIM] Ke: %s | OTP Anda adalah: %s. Berlaku selama 15 menit.\n", toEmail, otp)
	fmt.Printf("\n--- MOCK EMAIL SENDER ---\nTo: %s\nSubject: Kode OTP Reset Password GKPI Cimahi\n\nKode OTP Anda adalah: %s\n-------------------------\n\n", toEmail, otp)
	return nil
}

package domain

import "time"

type OTPCode struct {
	ID         uint      `json:"id" gorm:"primaryKey"`
	UserID     uint      `json:"user_id"`
	CodeHash   string    `json:"-"`       // Hanya hash yang disimpan, jangan plaintext
	Channel    string    `json:"channel"` // "email" atau "telegram"
	Purpose    string    `json:"purpose" gorm:"default:'verification'"` // "verification", "reset_password"
	ExpiresAt  time.Time `json:"expires_at"`
	Attempts   int       `json:"attempts" gorm:"default:0"`
	VerifiedAt *time.Time `json:"verified_at"`
	CreatedAt  time.Time `json:"created_at"`
}

type OTPSetting struct {
	ID                 uint `json:"id" gorm:"primaryKey"`
	EmailEnabled       bool `json:"email_enabled" gorm:"default:true"`
	TelegramEnabled    bool `json:"telegram_enabled" gorm:"default:false"`
	OTPExpiryMinutes   int  `json:"otp_expiry_minutes" gorm:"default:5"`
	MaxAttempts        int  `json:"max_attempts" gorm:"default:3"`
	MaxResendPerHour   int  `json:"max_resend_per_hour" gorm:"default:3"`
}

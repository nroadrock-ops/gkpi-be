package domain

import (
	"time"
)

type User struct {
	ID                        uint       `json:"id" gorm:"primaryKey"`
	Email                     string     `json:"email" gorm:"unique"`
	Password                  string     `json:"-"` // Tidak direturn ke JSON
	Role                      string     `json:"role"` // admin, jemaat
	JemaatID                  *uint      `json:"jemaat_id"`
	IsVerified                bool       `json:"is_verified" gorm:"default:false"`
	TelegramChatID            *string    `json:"telegram_chat_id"`
	TelegramConnectToken      *string    `json:"-"`
	TelegramConnectedAt       *time.Time `json:"telegram_connected_at"`
	ResetPasswordOTP          *string    `json:"-"`
	ResetPasswordOTPExpiredAt *time.Time `json:"-"`
	CreatedAt                 time.Time  `json:"created_at"`
	UpdatedAt                 time.Time  `json:"updated_at"`
}

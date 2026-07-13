package domain

import (
	"time"
)

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"` // Tidak direturn ke JSON
	Role                        string     `json:"role"` // admin, jemaat
	JemaatID                    *uint      `json:"jemaat_id"`
	ResetPasswordOTP            *string    `json:"-"`
	ResetPasswordOTPExpiredAt   *time.Time `json:"-"`
	CreatedAt                   time.Time  `json:"created_at"`
	UpdatedAt                   time.Time  `json:"updated_at"`
}

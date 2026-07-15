package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type OTPRepository interface {
	CreateCode(code *domain.OTPCode) error
	GetLatestUnverifiedCode(userID uint, channel string, purpose string) (*domain.OTPCode, error)
	UpdateCode(code *domain.OTPCode) error
	CountResendLastHour(userID uint, channel string) (int64, error)
	GetSettings() (*domain.OTPSetting, error)
	UpdateSettings(settings *domain.OTPSetting) error
}

type otpRepository struct {
	db *gorm.DB
}

func NewOTPRepository(db *gorm.DB) OTPRepository {
	return &otpRepository{db}
}

func (r *otpRepository) CreateCode(code *domain.OTPCode) error {
	return r.db.Create(code).Error
}

func (r *otpRepository) GetLatestUnverifiedCode(userID uint, channel string, purpose string) (*domain.OTPCode, error) {
	var code domain.OTPCode
	err := r.db.Where("user_id = ? AND channel = ? AND purpose = ? AND verified_at IS NULL", userID, channel, purpose).
		Order("created_at desc").
		First(&code).Error
	return &code, err
}

func (r *otpRepository) UpdateCode(code *domain.OTPCode) error {
	return r.db.Save(code).Error
}

func (r *otpRepository) CountResendLastHour(userID uint, channel string) (int64, error) {
	var count int64
	err := r.db.Model(&domain.OTPCode{}).
		Where("user_id = ? AND channel = ? AND created_at >= NOW() - INTERVAL '1 hour'", userID, channel).
		Count(&count).Error
	return count, err
}

func (r *otpRepository) GetSettings() (*domain.OTPSetting, error) {
	var settings domain.OTPSetting
	err := r.db.First(&settings).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			// Buat setting default jika tabel masih kosong (karena belum di-seed)
			settings = domain.OTPSetting{
				EmailEnabled:     false, // Sesuai instruksi: Email false secara default
				TelegramEnabled:  true,  // Sesuai instruksi: Telegram true secara default
				OTPExpiryMinutes: 5,
				MaxAttempts:      3,
				MaxResendPerHour: 3,
			}
			if createErr := r.db.Create(&settings).Error; createErr != nil {
				return nil, createErr
			}
			return &settings, nil
		}
		return nil, err
	}
	return &settings, nil
}

func (r *otpRepository) UpdateSettings(settings *domain.OTPSetting) error {
	return r.db.Save(settings).Error
}

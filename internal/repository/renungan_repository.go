package repository

import (
	"time"

	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type RenunganRepository interface {
	FindToday() (*domain.Renungan, error)
	Create(renungan *domain.Renungan) error
}

type renunganRepository struct {
	db *gorm.DB
}

func NewRenunganRepository(db *gorm.DB) RenunganRepository {
	return &renunganRepository{db}
}

func (r *renunganRepository) FindToday() (*domain.Renungan, error) {
	var renungan domain.Renungan
	today := time.Now().Format("2006-01-02")
	err := r.db.Where("tanggal = ?", today).First(&renungan).Error
	return &renungan, err
}

func (r *renunganRepository) Create(renungan *domain.Renungan) error {
	return r.db.Create(renungan).Error
}

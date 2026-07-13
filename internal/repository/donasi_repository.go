package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type DonasiRepository interface {
	FindByID(id uint) (*domain.Donasi, error)
	FindByTransactionID(trxID string) (*domain.Donasi, error)
	FindByJemaatID(jemaatID uint) ([]domain.Donasi, error)
	Create(donasi *domain.Donasi) error
	Update(donasi *domain.Donasi) error
}

type donasiRepository struct {
	db *gorm.DB
}

func NewDonasiRepository(db *gorm.DB) DonasiRepository {
	return &donasiRepository{db}
}

func (r *donasiRepository) FindByID(id uint) (*domain.Donasi, error) {
	var donasi domain.Donasi
	err := r.db.First(&donasi, id).Error
	return &donasi, err
}

func (r *donasiRepository) FindByTransactionID(trxID string) (*domain.Donasi, error) {
	var donasi domain.Donasi
	err := r.db.Where("transaction_id = ?", trxID).First(&donasi).Error
	return &donasi, err
}

func (r *donasiRepository) FindByJemaatID(jemaatID uint) ([]domain.Donasi, error) {
	var donasis []domain.Donasi
	err := r.db.Where("jemaat_id = ?", jemaatID).Find(&donasis).Error
	return donasis, err
}

func (r *donasiRepository) Create(donasi *domain.Donasi) error {
	return r.db.Create(donasi).Error
}

func (r *donasiRepository) Update(donasi *domain.Donasi) error {
	return r.db.Save(donasi).Error
}

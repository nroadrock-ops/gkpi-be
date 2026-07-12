package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type JadwalRepository interface {
	FindAll() ([]domain.JadwalIbadah, error)
	FindByID(id uint) (*domain.JadwalIbadah, error)
	Create(jadwal *domain.JadwalIbadah) error
	Update(jadwal *domain.JadwalIbadah) error
	Delete(id uint) error
}

type jadwalRepository struct {
	db *gorm.DB
}

func NewJadwalRepository(db *gorm.DB) JadwalRepository {
	return &jadwalRepository{db}
}

func (r *jadwalRepository) FindAll() ([]domain.JadwalIbadah, error) {
	var jadwals []domain.JadwalIbadah
	err := r.db.Find(&jadwals).Error
	return jadwals, err
}

func (r *jadwalRepository) FindByID(id uint) (*domain.JadwalIbadah, error) {
	var jadwal domain.JadwalIbadah
	err := r.db.First(&jadwal, id).Error
	return &jadwal, err
}

func (r *jadwalRepository) Create(jadwal *domain.JadwalIbadah) error {
	return r.db.Create(jadwal).Error
}

func (r *jadwalRepository) Update(jadwal *domain.JadwalIbadah) error {
	return r.db.Save(jadwal).Error
}

func (r *jadwalRepository) Delete(id uint) error {
	return r.db.Delete(&domain.JadwalIbadah{}, id).Error
}

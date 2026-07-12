package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type PengumumanRepository interface {
	FindAll() ([]domain.Pengumuman, error)
	Create(pengumuman *domain.Pengumuman) error
	Delete(id uint) error
}

type pengumumanRepository struct {
	db *gorm.DB
}

func NewPengumumanRepository(db *gorm.DB) PengumumanRepository {
	return &pengumumanRepository{db}
}

func (r *pengumumanRepository) FindAll() ([]domain.Pengumuman, error) {
	var pengumumans []domain.Pengumuman
	err := r.db.Order("created_at desc").Find(&pengumumans).Error
	return pengumumans, err
}

func (r *pengumumanRepository) Create(pengumuman *domain.Pengumuman) error {
	return r.db.Create(pengumuman).Error
}

func (r *pengumumanRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Pengumuman{}, id).Error
}

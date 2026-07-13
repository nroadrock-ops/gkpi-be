package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type GaleriRepository interface {
	FindAll() ([]domain.Galeri, error)
	FindByID(id uint) (*domain.Galeri, error)
	Create(galeri *domain.Galeri) error
	Update(galeri *domain.Galeri) error
	Delete(id uint) error
}

type galeriRepository struct {
	db *gorm.DB
}

func NewGaleriRepository(db *gorm.DB) GaleriRepository {
	return &galeriRepository{db}
}

func (r *galeriRepository) FindAll() ([]domain.Galeri, error) {
	var galeris []domain.Galeri
	err := r.db.Order("created_at desc").Find(&galeris).Error
	return galeris, err
}

func (r *galeriRepository) FindByID(id uint) (*domain.Galeri, error) {
	var galeri domain.Galeri
	err := r.db.First(&galeri, id).Error
	return &galeri, err
}

func (r *galeriRepository) Create(galeri *domain.Galeri) error {
	return r.db.Create(galeri).Error
}

func (r *galeriRepository) Update(galeri *domain.Galeri) error {
	return r.db.Save(galeri).Error
}

func (r *galeriRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Galeri{}, id).Error
}

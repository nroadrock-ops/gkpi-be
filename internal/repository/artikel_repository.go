package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type ArtikelRepository interface {
	FindAll(offset, limit int) ([]domain.Artikel, error)
	FindBySlug(slug string) (*domain.Artikel, error)
	FindByID(id uint) (*domain.Artikel, error)
	Create(artikel *domain.Artikel) error
	Update(artikel *domain.Artikel) error
	Delete(id uint) error
}

type artikelRepository struct {
	db *gorm.DB
}

func NewArtikelRepository(db *gorm.DB) ArtikelRepository {
	return &artikelRepository{db}
}

func (r *artikelRepository) FindAll(offset, limit int) ([]domain.Artikel, error) {
	var artikels []domain.Artikel
	err := r.db.Offset(offset).Limit(limit).Order("created_at desc").Find(&artikels).Error
	return artikels, err
}

func (r *artikelRepository) FindBySlug(slug string) (*domain.Artikel, error) {
	var artikel domain.Artikel
	err := r.db.Where("slug = ?", slug).First(&artikel).Error
	return &artikel, err
}

func (r *artikelRepository) FindByID(id uint) (*domain.Artikel, error) {
	var artikel domain.Artikel
	err := r.db.First(&artikel, id).Error
	return &artikel, err
}

func (r *artikelRepository) Create(artikel *domain.Artikel) error {
	return r.db.Create(artikel).Error
}

func (r *artikelRepository) Update(artikel *domain.Artikel) error {
	return r.db.Save(artikel).Error
}

func (r *artikelRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Artikel{}, id).Error
}

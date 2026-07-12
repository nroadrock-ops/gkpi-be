package repository

import (
	"gkpi-be/internal/domain"
	"gorm.io/gorm"
)

type JemaatRepository interface {
	FindAll() ([]domain.Jemaat, error)
	FindByID(id uint) (*domain.Jemaat, error)
	Create(jemaat *domain.Jemaat) error
	Update(jemaat *domain.Jemaat) error
	Delete(id uint) error
}

type jemaatRepository struct {
	db *gorm.DB
}

func NewJemaatRepository(db *gorm.DB) JemaatRepository {
	return &jemaatRepository{db}
}

func (r *jemaatRepository) FindAll() ([]domain.Jemaat, error) {
	var jemaats []domain.Jemaat
	err := r.db.Find(&jemaats).Error
	return jemaats, err
}

func (r *jemaatRepository) FindByID(id uint) (*domain.Jemaat, error) {
	var jemaat domain.Jemaat
	err := r.db.First(&jemaat, id).Error
	return &jemaat, err
}

func (r *jemaatRepository) Create(jemaat *domain.Jemaat) error {
	return r.db.Create(jemaat).Error
}

func (r *jemaatRepository) Update(jemaat *domain.Jemaat) error {
	return r.db.Save(jemaat).Error
}

func (r *jemaatRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Jemaat{}, id).Error
}

package service

import (
	"strings"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type ArtikelService interface {
	GetAll(page, limit int) ([]domain.Artikel, error)
	GetBySlug(slug string) (*domain.Artikel, error)
	Create(artikel *domain.Artikel) error
	Update(id uint, req *domain.Artikel) error
	Delete(id uint) error
}

type artikelService struct {
	repo repository.ArtikelRepository
}

func NewArtikelService(repo repository.ArtikelRepository) ArtikelService {
	return &artikelService{repo}
}

func (s *artikelService) GetAll(page, limit int) ([]domain.Artikel, error) {
	offset := (page - 1) * limit
	return s.repo.FindAll(offset, limit)
}

func (s *artikelService) GetBySlug(slug string) (*domain.Artikel, error) {
	return s.repo.FindBySlug(slug)
}

func (s *artikelService) Create(artikel *domain.Artikel) error {
	if artikel.Slug == "" {
		artikel.Slug = strings.ReplaceAll(strings.ToLower(artikel.Judul), " ", "-")
	}
	return s.repo.Create(artikel)
}

func (s *artikelService) Update(id uint, req *domain.Artikel) error {
	artikel, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	artikel.Judul = req.Judul
	if req.Slug != "" {
		artikel.Slug = req.Slug
	}
	artikel.Konten = req.Konten
	artikel.Penulis = req.Penulis
	return s.repo.Update(artikel)
}

func (s *artikelService) Delete(id uint) error {
	return s.repo.Delete(id)
}

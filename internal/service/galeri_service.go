package service

import (
	"fmt"
	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type GaleriService interface {
	GetAll() ([]domain.Galeri, error)
	UploadAndSave(judul, fileURL string) (*domain.Galeri, error)
	Delete(id uint) error
}

type galeriService struct {
	repo repository.GaleriRepository
}

func NewGaleriService(repo repository.GaleriRepository) GaleriService {
	return &galeriService{repo}
}

func (s *galeriService) GetAll() ([]domain.Galeri, error) {
	return s.repo.FindAll()
}

func (s *galeriService) UploadAndSave(judul, fileURL string) (*domain.Galeri, error) {
	galeri := &domain.Galeri{
		Judul:   judul,
		FileURL: fileURL,
		Tags:    "pending-ai-tag", // In a real scenario, this would be updated async
	}

	if err := s.repo.Create(galeri); err != nil {
		return nil, err
	}

	// Placeholder for async AI tagging call
	fmt.Println("Triggering async AI tagging for URL:", fileURL)

	return galeri, nil
}

func (s *galeriService) Delete(id uint) error {
	// In a real scenario, also delete file from Supabase Storage
	return s.repo.Delete(id)
}

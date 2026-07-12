package service

import (
	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type PengumumanService interface {
	GetAll() ([]domain.Pengumuman, error)
	Create(pengumuman *domain.Pengumuman) error
	Delete(id uint) error
}

type pengumumanService struct {
	repo repository.PengumumanRepository
}

func NewPengumumanService(repo repository.PengumumanRepository) PengumumanService {
	return &pengumumanService{repo}
}

func (s *pengumumanService) GetAll() ([]domain.Pengumuman, error) {
	return s.repo.FindAll()
}

func (s *pengumumanService) Create(pengumuman *domain.Pengumuman) error {
	return s.repo.Create(pengumuman)
}

func (s *pengumumanService) Delete(id uint) error {
	return s.repo.Delete(id)
}

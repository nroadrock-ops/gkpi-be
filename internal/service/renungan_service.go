package service

import (
	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type RenunganService interface {
	GetToday() (*domain.Renungan, error)
	Create(renungan *domain.Renungan) error
}

type renunganService struct {
	repo repository.RenunganRepository
}

func NewRenunganService(repo repository.RenunganRepository) RenunganService {
	return &renunganService{repo}
}

func (s *renunganService) GetToday() (*domain.Renungan, error) {
	return s.repo.FindToday()
}

func (s *renunganService) Create(renungan *domain.Renungan) error {
	return s.repo.Create(renungan)
}

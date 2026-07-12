package service

import (
	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type JemaatService interface {
	GetAll() ([]domain.Jemaat, error)
	GetByID(id uint) (*domain.Jemaat, error)
	Create(jemaat *domain.Jemaat) error
	Update(id uint, req *domain.Jemaat) error
	Delete(id uint) error
}

type jemaatService struct {
	repo repository.JemaatRepository
}

func NewJemaatService(repo repository.JemaatRepository) JemaatService {
	return &jemaatService{repo}
}

func (s *jemaatService) GetAll() ([]domain.Jemaat, error) {
	return s.repo.FindAll()
}

func (s *jemaatService) GetByID(id uint) (*domain.Jemaat, error) {
	return s.repo.FindByID(id)
}

func (s *jemaatService) Create(jemaat *domain.Jemaat) error {
	return s.repo.Create(jemaat)
}

func (s *jemaatService) Update(id uint, req *domain.Jemaat) error {
	jemaat, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	jemaat.NamaLengkap = req.NamaLengkap
	jemaat.Alamat = req.Alamat
	jemaat.NoTelepon = req.NoTelepon
	jemaat.Status = req.Status
	return s.repo.Update(jemaat)
}

func (s *jemaatService) Delete(id uint) error {
	return s.repo.Delete(id)
}

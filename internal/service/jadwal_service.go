package service

import (
	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type JadwalService interface {
	GetAll() ([]domain.JadwalIbadah, error)
	GetByID(id uint) (*domain.JadwalIbadah, error)
	Create(jadwal *domain.JadwalIbadah) error
	Update(id uint, req *domain.JadwalIbadah) error
	Delete(id uint) error
}

type jadwalService struct {
	repo repository.JadwalRepository
}

func NewJadwalService(repo repository.JadwalRepository) JadwalService {
	return &jadwalService{repo}
}

func (s *jadwalService) GetAll() ([]domain.JadwalIbadah, error) {
	return s.repo.FindAll()
}

func (s *jadwalService) GetByID(id uint) (*domain.JadwalIbadah, error) {
	return s.repo.FindByID(id)
}

func (s *jadwalService) Create(jadwal *domain.JadwalIbadah) error {
	return s.repo.Create(jadwal)
}

func (s *jadwalService) Update(id uint, req *domain.JadwalIbadah) error {
	jadwal, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	jadwal.NamaIbadah = req.NamaIbadah
	jadwal.Waktu = req.Waktu
	jadwal.Lokasi = req.Lokasi
	jadwal.Pengkhotbah = req.Pengkhotbah
	return s.repo.Update(jadwal)
}

func (s *jadwalService) Delete(id uint) error {
	return s.repo.Delete(id)
}

package service

import (
	"fmt"
	"time"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
)

type DonasiService interface {
	Create(req *domain.Donasi) (*domain.Donasi, error)
	GetStatus(id uint) (*domain.Donasi, error)
	GetByJemaatID(jemaatID uint) ([]domain.Donasi, error)
	ProcessWebhook(trxID, status string) error
}

type donasiService struct {
	repo repository.DonasiRepository
}

func NewDonasiService(repo repository.DonasiRepository) DonasiService {
	return &donasiService{repo}
}

func (s *donasiService) Create(req *domain.Donasi) (*domain.Donasi, error) {
	req.Status = "pending"
	req.TransactionID = fmt.Sprintf("TRX-%d", time.Now().Unix())
	req.PaymentURL = "https://app.midtrans.com/snap/v2/vtweb/placeholder-token"

	if err := s.repo.Create(req); err != nil {
		return nil, err
	}
	return req, nil
}

func (s *donasiService) GetStatus(id uint) (*domain.Donasi, error) {
	return s.repo.FindByID(id)
}

func (s *donasiService) GetByJemaatID(jemaatID uint) ([]domain.Donasi, error) {
	return s.repo.FindByJemaatID(jemaatID)
}

func (s *donasiService) ProcessWebhook(trxID, status string) error {
	donasi, err := s.repo.FindByTransactionID(trxID)
	if err != nil {
		return err
	}
	donasi.Status = status
	return s.repo.Update(donasi)
}

package service

import (
	"errors"
	"os"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
	"gkpi-be/internal/utils"
)

type AuthService interface {
	Register(email, password, role string) (*domain.User, error)
	Login(email, password string) (string, string, error) // returns accessToken, refreshToken
	GetMe(id uint) (*domain.User, error)
}

type authService struct {
	repo repository.UserRepository
}

func NewAuthService(repo repository.UserRepository) AuthService {
	return &authService{repo}
}

func (s *authService) Register(email, password, role string) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
	}

	err = s.repo.Create(user)
	return user, err
}

func (s *authService) Login(email, password string) (string, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", errors.New("invalid credentials")
	}

	return utils.GenerateTokens(user.ID, user.Role, os.Getenv("JWT_SECRET"))
}

func (s *authService) GetMe(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

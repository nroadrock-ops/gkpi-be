package service

import (
	"errors"
	"os"
	"time"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
	"gkpi-be/internal/utils"
)

type AuthService interface {
	Register(namaLengkap, email, password, noTelepon, role string) (*domain.User, error)
	Login(email, password string) (string, string, string, error) // returns accessToken, refreshToken, role
	GetMe(id uint) (*domain.User, error)
	ForgotPassword(email string) error
	ResetPassword(email, otp, newPassword string) error
}

type authService struct {
	repo       repository.UserRepository
	jemaatRepo repository.JemaatRepository
}

func NewAuthService(repo repository.UserRepository, jemaatRepo repository.JemaatRepository) AuthService {
	return &authService{repo, jemaatRepo}
}

func (s *authService) Register(namaLengkap, email, password, noTelepon, role string) (*domain.User, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, err
	}

	var jemaatID *uint

	if role == "jemaat" {
		jemaat := &domain.Jemaat{
			NamaLengkap: namaLengkap,
			NoTelepon:   noTelepon,
			Status:      "aktif",
		}
		err = s.jemaatRepo.Create(jemaat)
		if err != nil {
			return nil, err
		}
		jemaatID = &jemaat.ID
	}

	user := &domain.User{
		Email:    email,
		Password: hashedPassword,
		Role:     role,
		JemaatID: jemaatID,
	}

	err = s.repo.Create(user)
	return user, err
}

func (s *authService) Login(email, password string) (string, string, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", "", errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", "", errors.New("invalid credentials")
	}

	acc, ref, err := utils.GenerateTokens(user.ID, user.Role, user.JemaatID, os.Getenv("JWT_SECRET"))
	return acc, ref, user.Role, err
}

func (s *authService) GetMe(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *authService) ForgotPassword(email string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	// Generate 6 digit OTP
	otp := utils.GenerateOTP()
	expiredAt := time.Now().Add(15 * time.Minute)

	user.ResetPasswordOTP = &otp
	user.ResetPasswordOTPExpiredAt = &expiredAt

	if err := s.repo.Update(user); err != nil {
		return err
	}

	// Send email mock
	return utils.SendEmailOTP(user.Email, otp)
}

func (s *authService) ResetPassword(email, otp, newPassword string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	if user.ResetPasswordOTP == nil || user.ResetPasswordOTPExpiredAt == nil {
		return errors.New("invalid otp")
	}

	if *user.ResetPasswordOTP != otp {
		return errors.New("invalid otp")
	}

	if time.Now().After(*user.ResetPasswordOTPExpiredAt) {
		return errors.New("otp expired")
	}

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword
	user.ResetPasswordOTP = nil
	user.ResetPasswordOTPExpiredAt = nil

	return s.repo.Update(user)
}

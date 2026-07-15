package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/repository"
	"gkpi-be/internal/utils"
)

type AuthService interface {
	Register(namaLengkap, email, password, noTelepon, role string) (*domain.User, string, string, error) // return: user, connectToken, connectURL, error
	Login(email, password string) (string, string, string, bool, error) // return: acc, ref, role, needs_verification, error
	GetMe(id uint) (*domain.User, error)
	ForgotPassword(email string) (string, error)
	ResetPassword(email, otp, newPassword string) error
	UpdatePassword(userID uint, oldPassword, newPassword string) error

	// OTP
	VerifyOTP(email, code, purpose string) (string, string, string, error) // return: acc, ref, role, error
	ResendOTP(email, channel, purpose string) (string, error)
	GetTelegramConnectLinkByToken(token string) (string, error)
	ProcessTelegramWebhook(payload map[string]interface{}) error
	CheckTelegramStatus(token string) (bool, error)
	GetOTPSettings() (*domain.OTPSetting, error)
	UpdateOTPSettings(settings *domain.OTPSetting) error
}

type authService struct {
	repo       repository.UserRepository
	jemaatRepo repository.JemaatRepository
	otpRepo    repository.OTPRepository
}

func NewAuthService(repo repository.UserRepository, jemaatRepo repository.JemaatRepository, otpRepo repository.OTPRepository) AuthService {
	return &authService{repo, jemaatRepo, otpRepo}
}

// hashOTP digunakan untuk hashing OTP sebelum disimpan.
// Mengapa OTP di-hash bukan disimpan plaintext?
// Karena OTP sama seperti password sementara. Jika database bocor, penyerang tidak bisa langsung
// menggunakan OTP yang masih aktif untuk mengambil alih akun user.
func hashOTP(code string) string {
	h := sha256.New()
	h.Write([]byte(code))
	return hex.EncodeToString(h.Sum(nil))
}

func (s *authService) Register(namaLengkap, email, password, noTelepon, role string) (*domain.User, string, string, error) {
	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, "", "", err
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
			return nil, "", "", err
		}
		jemaatID = &jemaat.ID
	}

	user := &domain.User{
		Email:      email,
		Password:   hashedPassword,
		Role:       role,
		JemaatID:   jemaatID,
		IsVerified: false, // User baru wajib false
	}

	err = s.repo.Create(user)
	if err != nil {
		return nil, "", "", err
	}

	// Generate connect_token unik untuk Telegram
	connectToken := fmt.Sprintf("usr_%d_%x", user.ID, time.Now().UnixNano())
	user.TelegramConnectToken = &connectToken
	s.repo.Update(user)

	connectURL, _ := s.GetTelegramConnectLinkByToken(connectToken)

	// JANGAN langsung generate/kirim OTP sesuai instruksi
	return user, connectToken, connectURL, nil
}

func (s *authService) generateAndSendOTP(user *domain.User, channel string, purpose string) error {
	settings, err := s.otpRepo.GetSettings()
	if err != nil {
		return err
	}

	// Generate 6 digit
	code := utils.GenerateOTP()
	expiresAt := time.Now().Add(time.Duration(settings.OTPExpiryMinutes) * time.Minute)

	otpCode := &domain.OTPCode{
		UserID:    user.ID,
		CodeHash:  hashOTP(code),
		Channel:   channel,
		Purpose:   purpose,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := s.otpRepo.CreateCode(otpCode); err != nil {
		return err
	}

	// Mengambil data nomor telepon untuk ditampilkan di pesan
	var jemaatPhone string
	var jemaatName string
	if user.JemaatID != nil {
		if j, _ := s.jemaatRepo.FindByID(*user.JemaatID); j != nil {
			jemaatName = j.NamaLengkap
			jemaatPhone = j.NoTelepon
		}
	}
	if jemaatName == "" {
		jemaatName = "Jemaat"
	}
	if jemaatPhone == "" {
		jemaatPhone = "-"
	}

	if channel == "email" && settings.EmailEnabled {
		return SendOTPEmail(user.Email, jemaatName, code)
	} else if channel == "telegram" && settings.TelegramEnabled {
		if user.TelegramChatID != nil && *user.TelegramChatID != "" {
			// Fungsi pembantu untuk meng-escape karakter HTML
			escapeHTML := func(text string) string {
				text = strings.ReplaceAll(text, "&", "&amp;")
				text = strings.ReplaceAll(text, "<", "&lt;")
				text = strings.ReplaceAll(text, ">", "&gt;")
				return text
			}

			safeName := escapeHTML(jemaatName)
			safeEmail := escapeHTML(user.Email)
			safePhone := escapeHTML(jemaatPhone)
			expiry := settings.OTPExpiryMinutes

			var message string
			if purpose == "reset_password" {
				// Pesan untuk lupa password (fokus pada nama saja)
				message = fmt.Sprintf("Syalom, <b>%s</b>! 🙏\n\nAnda baru saja meminta kode OTP untuk melakukan reset password.\n\nKode verifikasi Anda:\n<b>%s</b>\n\nKode berlaku selama %d menit. Jangan bagikan kode ini kepada siapa pun, termasuk pihak yang mengaku dari gereja.", safeName, code, expiry)
			} else {
				// Pesan untuk registrasi / resend OTP verifikasi
				message = fmt.Sprintf("Syalom, <b>%s</b>! 🙏\n\nTerima kasih telah mendaftar di Portal Jemaat GKPI Cimahi. Berikut data yang Anda daftarkan:\n\nNama: <b>%s</b>\nEmail: %s\nNo. Telepon: %s\n\nJika ada data yang tidak sesuai, silakan hubungi admin gereja setelah proses verifikasi ini selesai.\n\nKode verifikasi Anda:\n<b>%s</b>\n\nKode berlaku selama %d menit. Jangan bagikan kode ini kepada siapa pun, termasuk pihak yang mengaku dari gereja.", safeName, safeName, safeEmail, safePhone, code, expiry)
			}

			return SendOTPTelegram(*user.TelegramChatID, message)
		}
		return errors.New("telegram not connected")
	}

	return nil
}

func (s *authService) VerifyOTP(email, code, purpose string) (string, string, string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", "", errors.New("user not found")
	}

	settings, err := s.otpRepo.GetSettings()
	if err != nil {
		return "", "", "", err
	}

	// Cek OTP aktif (prioritas telegram lalu email)
	otpEntry, err := s.otpRepo.GetLatestUnverifiedCode(user.ID, "telegram", purpose)
	if err != nil {
		otpEntry, err = s.otpRepo.GetLatestUnverifiedCode(user.ID, "email", purpose)
		if err != nil {
			return "", "", "", errors.New("no pending otp found")
		}
	}

	if otpEntry.Attempts >= settings.MaxAttempts {
		return "", "", "", errors.New("max attempts reached, please request a new otp")
	}

	if time.Now().After(otpEntry.ExpiresAt) {
		return "", "", "", errors.New("otp expired, please request a new one")
	}

	if otpEntry.CodeHash != hashOTP(code) {
		otpEntry.Attempts++
		s.otpRepo.UpdateCode(otpEntry)
		return "", "", "", fmt.Errorf("invalid otp, %d attempts left", settings.MaxAttempts-otpEntry.Attempts)
	}

	// Verified!
	now := time.Now()
	otpEntry.VerifiedAt = &now
	s.otpRepo.UpdateCode(otpEntry)

	if purpose == "verification" {
		user.IsVerified = true
		s.repo.Update(user)
	}

	// Auto login
	acc, ref, err := utils.GenerateTokens(user.ID, user.Role, user.JemaatID, os.Getenv("JWT_SECRET"))
	return acc, ref, user.Role, err
}

func (s *authService) ResendOTP(email, channel, purpose string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("user not found")
	}

	settings, err := s.otpRepo.GetSettings()
	if err != nil {
		return "", err
	}

	if channel == "telegram" && !settings.TelegramEnabled {
		return "", errors.New("telegram otp is disabled")
	}

	if channel == "telegram" && (user.TelegramChatID == nil || *user.TelegramChatID == "") {
		url, _ := s.GetTelegramConnectLinkByToken(*user.TelegramConnectToken)
		return url, errors.New("telegram not connected")
	}

	// Cek rate limit
	count, _ := s.otpRepo.CountResendLastHour(user.ID, channel)
	if count >= int64(settings.MaxResendPerHour) {
		return "", errors.New("rate limit exceeded, try again later")
	}

	return "", s.generateAndSendOTP(user, channel, purpose)
}

func (s *authService) Login(email, password string) (string, string, string, bool, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", "", "", false, errors.New("invalid credentials")
	}

	if !utils.CheckPasswordHash(password, user.Password) {
		return "", "", "", false, errors.New("invalid credentials")
	}

	if !user.IsVerified && user.Role != "admin" {
		return "", "", "", true, nil // needs_verification = true
	}

	acc, ref, err := utils.GenerateTokens(user.ID, user.Role, user.JemaatID, os.Getenv("JWT_SECRET"))
	return acc, ref, user.Role, false, err
}

func (s *authService) GetMe(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *authService) ForgotPassword(email string) (string, error) {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", errors.New("email not found")
	}

	// Kalau user belum pernah hubungkan Telegram
	if user.TelegramChatID == nil || *user.TelegramChatID == "" {
		url, _ := s.GetTelegramConnectLinkByToken(*user.TelegramConnectToken)
		return url, errors.New("telegram not connected")
	}

	return "", s.generateAndSendOTP(user, "telegram", "reset_password")
}

func (s *authService) ResetPassword(email, otp, newPassword string) error {
	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return errors.New("email not found")
	}

	// Validasi OTP tipe reset_password di telegram
	otpEntry, err := s.otpRepo.GetLatestUnverifiedCode(user.ID, "telegram", "reset_password")
	if err != nil {
		return errors.New("invalid or expired otp")
	}

	settings, _ := s.otpRepo.GetSettings()
	if otpEntry.Attempts >= settings.MaxAttempts {
		return errors.New("max attempts reached, please request a new otp")
	}

	if time.Now().After(otpEntry.ExpiresAt) {
		return errors.New("otp expired, please request a new one")
	}

	if otpEntry.CodeHash != hashOTP(otp) {
		otpEntry.Attempts++
		s.otpRepo.UpdateCode(otpEntry)
		return fmt.Errorf("invalid otp, %d attempts left", settings.MaxAttempts-otpEntry.Attempts)
	}

	now := time.Now()
	otpEntry.VerifiedAt = &now
	s.otpRepo.UpdateCode(otpEntry)

	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	user.Password = hashedPassword

	return s.repo.Update(user)
}

func (s *authService) UpdatePassword(userID uint, oldPassword, newPassword string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	// 1. Cek password lama (konfirmasi)
	if !utils.CheckPasswordHash(oldPassword, user.Password) {
		return errors.New("password lama salah")
	}

	// 2. Hash password baru (pastikan tidak ada double-hash dengan nge-hash plaintext langsung)
	hashedPassword, err := utils.HashPassword(newPassword)
	if err != nil {
		return err
	}

	// 3. Update password
	user.Password = hashedPassword
	return s.repo.Update(user)
}

// Telegram integration

func (s *authService) GetTelegramConnectLinkByToken(token string) (string, error) {
	botUsername := os.Getenv("TELEGRAM_BOT_USERNAME")
	if botUsername == "" {
		return "", errors.New("telegram bot not configured")
	}
	link := fmt.Sprintf("https://t.me/%s?start=%s", botUsername, token)
	return link, nil
}

func (s *authService) CheckTelegramStatus(token string) (bool, error) {
	user, err := s.repo.FindByTelegramConnectToken(token)
	if err != nil {
		return false, errors.New("invalid connect token")
	}
	
	if user.TelegramChatID != nil && *user.TelegramChatID != "" {
		// Kalau baru saja connect tapi belum verified dan belum ada OTP verification aktif, trigger OTP pertama kali
		if !user.IsVerified {
			_, errOTP := s.otpRepo.GetLatestUnverifiedCode(user.ID, "telegram", "verification")
			if errOTP != nil {
				// Tidak ada OTP aktif, mari generate
				_ = s.generateAndSendOTP(user, "telegram", "verification")
			}
		}
		return true, nil
	}
	
	return false, nil
}

func (s *authService) ProcessTelegramWebhook(payload map[string]interface{}) error {
	message, ok := payload["message"].(map[string]interface{})
	if !ok {
		return nil
	}

	text, _ := message["text"].(string)
	chat, ok := message["chat"].(map[string]interface{})
	if !ok {
		return nil
	}

	chatIDFloat, _ := chat["id"].(float64)
	chatID := fmt.Sprintf("%.0f", chatIDFloat)

	// Format text saat deep link adalah: "/start token123"
	if len(text) > 7 && text[:6] == "/start" {
		token := text[7:] 
		user, err := s.repo.FindByTelegramConnectToken(token)
		if err == nil {
			now := time.Now()
			user.TelegramChatID = &chatID
			user.TelegramConnectedAt = &now
			s.repo.Update(user)
			// Kirim OTP verifikasi pertama kali secara otomatis (akan berisi ringkasan data)
			if !user.IsVerified {
				_ = s.generateAndSendOTP(user, "telegram", "verification")
			} else {
				_ = SendOTPTelegram(chatID, "Akun Telegram berhasil dihubungkan dengan GKPI Cimahi!")
			}
		}
	}

	return nil
}

func (s *authService) GetOTPSettings() (*domain.OTPSetting, error) {
	return s.otpRepo.GetSettings()
}

func (s *authService) UpdateOTPSettings(settings *domain.OTPSetting) error {
	return s.otpRepo.UpdateSettings(settings)
}

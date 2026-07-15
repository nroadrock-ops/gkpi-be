package router

import (
	"gkpi-be/internal/config"
	"gkpi-be/internal/handler"
	"gkpi-be/internal/middleware"
	"gkpi-be/internal/repository"
	"gkpi-be/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	cfg := config.LoadConfig()

	// Repositories
	userRepo := repository.NewUserRepository(db)
	jemaatRepo := repository.NewJemaatRepository(db)
	jadwalRepo := repository.NewJadwalRepository(db)
	artikelRepo := repository.NewArtikelRepository(db)
	renunganRepo := repository.NewRenunganRepository(db)
	pengumumanRepo := repository.NewPengumumanRepository(db)
	galeriRepo := repository.NewGaleriRepository(db)
	donasiRepo := repository.NewDonasiRepository(db)

	otpRepo := repository.NewOTPRepository(db)

	// AI Client with BaseURL and Internal Key
	aiClient := service.NewAIServiceClient(cfg.AIServiceURL, cfg.AIServiceInternalKey)

	// Services
	authSvc := service.NewAuthService(userRepo, jemaatRepo, otpRepo)
	jemaatSvc := service.NewJemaatService(jemaatRepo)
	jadwalSvc := service.NewJadwalService(jadwalRepo)
	artikelSvc := service.NewArtikelService(artikelRepo)
	renunganSvc := service.NewRenunganService(renunganRepo)
	pengumumanSvc := service.NewPengumumanService(pengumumanRepo)
	galeriSvc := service.NewGaleriService(galeriRepo, aiClient) // Inject aiClient
	donasiSvc := service.NewDonasiService(donasiRepo)

	// Handlers
	authHdl := handler.NewAuthHandler(authSvc)
	jemaatHdl := handler.NewJemaatHandler(jemaatSvc)
	jadwalHdl := handler.NewJadwalHandler(jadwalSvc)
	artikelHdl := handler.NewArtikelHandler(artikelSvc)
	renunganHdl := handler.NewRenunganHandler(renunganSvc)
	pengumumanHdl := handler.NewPengumumanHandler(pengumumanSvc)
	galeriHdl := handler.NewGaleriHandler(galeriSvc)
	donasiHdl := handler.NewDonasiHandler(donasiSvc)
	aiHdl := handler.NewAIHandler(aiClient)

	api := app.Group("/api/v1")

	// Auth routes
	auth := api.Group("/auth")
	auth.Post("/register", authHdl.Register)
	auth.Post("/login", authHdl.Login)
	auth.Post("/verify-otp", authHdl.VerifyOTP)
	auth.Post("/resend-otp", authHdl.ResendOTP)
	auth.Get("/telegram/status", authHdl.TelegramStatus)
	auth.Post("/telegram/webhook", authHdl.TelegramWebhook)
	
	auth.Post("/refresh", authHdl.Refresh)
	auth.Post("/forgot-password", authHdl.ForgotPassword)
	auth.Post("/reset-password", authHdl.ResetPassword)
	auth.Post("/logout", middleware.Protected(), authHdl.Logout)
	auth.Get("/me", middleware.Protected(), authHdl.Me)

	// Admin Pengaturan OTP (sesuai spesifikasi)
	adminPengaturan := api.Group("/admin/pengaturan", middleware.Protected(), middleware.AdminOnly())
	adminPengaturan.Get("/otp", authHdl.GetOTPSettings)
	adminPengaturan.Put("/otp", authHdl.UpdateOTPSettings)

	// Alias untuk /settings/otp (berdasarkan log yang mencoba mengakses path ini)
	settingsAlias := api.Group("/settings", middleware.Protected(), middleware.AdminOnly())
	settingsAlias.Get("/otp", authHdl.GetOTPSettings)
	settingsAlias.Put("/otp", authHdl.UpdateOTPSettings)

	// Jemaat routes
	jemaat := api.Group("/jemaat")

	// Jemaat self routes
	jemaatMe := jemaat.Group("/me", middleware.Protected(), middleware.RequireRole("jemaat"))
	jemaatMe.Get("/", jemaatHdl.GetMeProfile)
	jemaatMe.Put("/", jemaatHdl.UpdateMeProfile)
	// PISAHKAN endpoint update password dari update profil umum (sebaiknya juga ditambahkan di auth kalau admin butuh)
	jemaatMe.Put("/password", authHdl.UpdatePassword) 
	jemaatMe.Get("/absensi", aiHdl.GetMeAbsensi)
	jemaatMe.Get("/donasi", donasiHdl.GetMeDonasi)

	// Jemaat admin routes
	jemaatAdmin := jemaat.Group("/", middleware.Protected(), middleware.AdminOnly())
	jemaatAdmin.Get("/", jemaatHdl.GetAll)
	jemaatAdmin.Post("/", jemaatHdl.Create)
	jemaatAdmin.Get("/:id", jemaatHdl.GetByID)
	jemaatAdmin.Put("/:id", jemaatHdl.Update)
	jemaatAdmin.Delete("/:id", jemaatHdl.Delete)

	// Jadwal Ibadah
	jadwal := api.Group("/jadwal")
	jadwal.Get("/", jadwalHdl.GetAll) // public
	jadwalAdmin := jadwal.Group("/", middleware.Protected(), middleware.AdminOnly())
	jadwalAdmin.Post("/", jadwalHdl.Create)
	jadwalAdmin.Put("/:id", jadwalHdl.Update)
	jadwalAdmin.Delete("/:id", jadwalHdl.Delete)

	// Artikel
	artikel := api.Group("/artikel")
	artikel.Get("/", artikelHdl.GetAll) // public
	artikel.Get("/:slug", artikelHdl.GetBySlug) // public
	artikelAdmin := artikel.Group("/", middleware.Protected(), middleware.AdminOnly())
	artikelAdmin.Post("/", artikelHdl.Create)
	artikelAdmin.Put("/:id", artikelHdl.Update)
	artikelAdmin.Delete("/:id", artikelHdl.Delete)

	// Renungan
	renungan := api.Group("/renungan")
	renungan.Get("/today", renunganHdl.GetToday) // public
	renungan.Post("/", middleware.Protected(), middleware.AdminOnly(), renunganHdl.Create)

	// Pengumuman
	pengumuman := api.Group("/pengumuman")
	pengumuman.Get("/", pengumumanHdl.GetAll) // public
	pengumuman.Post("/", middleware.Protected(), middleware.AdminOnly(), pengumumanHdl.Create)
	pengumuman.Delete("/:id", middleware.Protected(), middleware.AdminOnly(), pengumumanHdl.Delete)

	// Galeri
	galeri := api.Group("/galeri")
	galeri.Get("/", galeriHdl.GetAll) // public
	galeri.Post("/upload", middleware.Protected(), middleware.AdminOnly(), galeriHdl.Upload)
	galeri.Delete("/:id", middleware.Protected(), middleware.AdminOnly(), galeriHdl.Delete)

	// Donasi
	donasi := api.Group("/donasi")
	donasi.Post("/", donasiHdl.Create) // public
	donasi.Get("/:id/status", donasiHdl.GetStatus)
	donasi.Post("/webhook", donasiHdl.Webhook)

	// AI Proxy
	ai := api.Group("/ai")
	ai.Post("/arsip/digitize", aiHdl.DigitizeArsip)
	ai.Post("/galeri/auto-tag", aiHdl.AutoTagGaleri) // This is if someone calls it directly
	ai.Post("/absensi/enroll", aiHdl.EnrollAbsensi)
	ai.Post("/absensi/check-in", aiHdl.CheckInAbsensi)
	ai.Get("/absensi/history/:jemaat_id", aiHdl.HistoryAbsensi)
}

package router

import (
	"gkpi-be/internal/handler"
	"gkpi-be/internal/middleware"
	"gkpi-be/internal/repository"
	"gkpi-be/internal/service"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func SetupRoutes(app *fiber.App, db *gorm.DB) {
	// Repositories
	userRepo := repository.NewUserRepository(db)
	jemaatRepo := repository.NewJemaatRepository(db)
	jadwalRepo := repository.NewJadwalRepository(db)
	artikelRepo := repository.NewArtikelRepository(db)
	renunganRepo := repository.NewRenunganRepository(db)
	pengumumanRepo := repository.NewPengumumanRepository(db)
	galeriRepo := repository.NewGaleriRepository(db)
	donasiRepo := repository.NewDonasiRepository(db)

	// Services
	authSvc := service.NewAuthService(userRepo)
	jemaatSvc := service.NewJemaatService(jemaatRepo)
	jadwalSvc := service.NewJadwalService(jadwalRepo)
	artikelSvc := service.NewArtikelService(artikelRepo)
	renunganSvc := service.NewRenunganService(renunganRepo)
	pengumumanSvc := service.NewPengumumanService(pengumumanRepo)
	galeriSvc := service.NewGaleriService(galeriRepo)
	donasiSvc := service.NewDonasiService(donasiRepo)
	aiClient := service.NewAIServiceClient("http://ai-service:5000/api/v1") // Mock AI Service URL

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
	auth.Post("/refresh", authHdl.Refresh)
	auth.Post("/logout", middleware.Protected(), authHdl.Logout)
	auth.Get("/me", middleware.Protected(), authHdl.Me)

	// Jemaat routes (Admin only)
	jemaat := api.Group("/jemaat", middleware.Protected(), middleware.AdminOnly())
	jemaat.Get("/", jemaatHdl.GetAll)
	jemaat.Post("/", jemaatHdl.Create)
	jemaat.Get("/:id", jemaatHdl.GetByID)
	jemaat.Put("/:id", jemaatHdl.Update)
	jemaat.Delete("/:id", jemaatHdl.Delete)

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
	ai := api.Group("/ai") // maybe protected depending on rules, assuming public for now or add middleware if needed
	ai.Post("/arsip/digitize", aiHdl.DigitizeArsip)
	ai.Post("/galeri/auto-tag", aiHdl.AutoTagGaleri)
	ai.Post("/absensi/enroll", aiHdl.EnrollAbsensi)
	ai.Post("/absensi/check-in", aiHdl.CheckInAbsensi)
	ai.Get("/absensi/history/:jemaat_id", aiHdl.HistoryAbsensi)
}

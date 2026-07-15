package database

import (
	"log"
	"time"

	"gkpi-be/internal/domain"
	"gkpi-be/internal/utils"
	"gorm.io/gorm"
)

func Seed(db *gorm.DB) {
	log.Println("Seeding database...")

	// Migrate schema
	err := db.AutoMigrate(
		&domain.User{},
		&domain.Jemaat{},
		&domain.JadwalIbadah{},
		&domain.Artikel{},
		&domain.Renungan{},
		&domain.Pengumuman{},
		&domain.Galeri{},
		&domain.Donasi{},
		&domain.Absensi{},
		&domain.OTPCode{},
		&domain.OTPSetting{},
	)
	if err != nil {
		log.Fatal("Failed to auto migrate: ", err)
	}

	// 1. Seed User
	hashedPassword, _ := utils.HashPassword("admin123")
	user := domain.User{Email: "admin@gkpicimahi.org", Password: hashedPassword, Role: "admin", IsVerified: true}
	if db.Where("email = ?", user.Email).First(&domain.User{}).Error != nil {
		db.Create(&user)
	}

	// 2. Seed Jemaat
	jemaat := domain.Jemaat{NamaLengkap: "Budi Santoso", Alamat: "Jl. Cimahi Raya No. 10", NoTelepon: "08123456789", Status: "aktif"}
	if db.Where("nama_lengkap = ?", jemaat.NamaLengkap).First(&domain.Jemaat{}).Error != nil {
		db.Create(&jemaat)
	}

	// 3. Seed Jadwal Ibadah
	jadwal := domain.JadwalIbadah{NamaIbadah: "Ibadah Minggu Pagi", Waktu: time.Now().AddDate(0, 0, 1), Lokasi: "Gereja GKPI Cimahi", Pengkhotbah: "Pdt. Siregar"}
	if db.Where("nama_ibadah = ?", jadwal.NamaIbadah).First(&domain.JadwalIbadah{}).Error != nil {
		db.Create(&jadwal)
	}

	// 4. Seed Artikel
	artikel := domain.Artikel{Judul: "Sejarah GKPI Cimahi", Slug: "sejarah-gkpi-cimahi", Konten: "Ini adalah konten sejarah GKPI Cimahi...", Penulis: "Admin"}
	if db.Where("slug = ?", artikel.Slug).First(&domain.Artikel{}).Error != nil {
		db.Create(&artikel)
	}

	// 5. Seed Renungan
	renungan := domain.Renungan{Judul: "Kasih Karunia", Ayat: "Efesus 2:8", Konten: "Sebab karena kasih karunia kamu diselamatkan...", Tanggal: time.Now()}
	if db.Where("judul = ?", renungan.Judul).First(&domain.Renungan{}).Error != nil {
		db.Create(&renungan)
	}

	// 6. Seed Pengumuman
	pengumuman := domain.Pengumuman{Judul: "Rapat Sintua", Konten: "Diingatkan kepada seluruh Sintua untuk hadir rapat hari Sabtu."}
	if db.Where("judul = ?", pengumuman.Judul).First(&domain.Pengumuman{}).Error != nil {
		db.Create(&pengumuman)
	}

	// 7. Seed OTP Setting
	var count int64
	db.Model(&domain.OTPSetting{}).Count(&count)
	if count == 0 {
		db.Create(&domain.OTPSetting{
			EmailEnabled:     true,
			TelegramEnabled:  false,
			OTPExpiryMinutes: 5,
			MaxAttempts:      3,
			MaxResendPerHour: 3,
		})
	}

	log.Println("Seeding completed successfully!")
}

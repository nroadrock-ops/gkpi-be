package database

import (
	"gkpi-be/internal/domain"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func Connect(dsn string) *gorm.DB {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Failed to connect to database: ", err)
	}

	log.Println("Connected to the database successfully")

	// Auto migrate schema
	db.AutoMigrate(&domain.User{}, &domain.Jemaat{}, &domain.Artikel{}, &domain.OTPCode{}, &domain.OTPSetting{})

	return db
}

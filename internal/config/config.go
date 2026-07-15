package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	DatabaseURL            string
	SupabaseURL            string
	SupabaseServiceKey     string
	JWTSecret              string
	AIServiceURL           string
	AIServiceInternalKey   string
	AllowedOrigins         string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, reading from environment variables")
	}

	return &Config{
		DatabaseURL:          os.Getenv("SUPABASE_DB_URL"),
		SupabaseURL:          os.Getenv("SUPABASE_URL"),
		SupabaseServiceKey:   os.Getenv("SUPABASE_SERVICE_KEY"),
		JWTSecret:            os.Getenv("JWT_SECRET"),
		AIServiceURL:         os.Getenv("AI_SERVICE_URL"),
		AIServiceInternalKey: os.Getenv("AI_SERVICE_INTERNAL_KEY"),
		AllowedOrigins:       os.Getenv("ALLOWED_ORIGINS"),
	}
}

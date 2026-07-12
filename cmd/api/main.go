package main

import (
	"flag"
	"log"

	"gkpi-be/internal/config"
	"gkpi-be/internal/database"
	"gkpi-be/internal/router"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

func main() {
	// Parse flag for seeding
	seedFlag := flag.Bool("seed", false, "Run database seeder")
	flag.Parse()

	// Load config
	cfg := config.LoadConfig()

	// Connect to Database
	db := database.Connect(cfg.DatabaseURL)

	if *seedFlag {
		database.Seed(db)
		return // Exit after seeding
	}

	// Initialize Fiber app
	app := fiber.New()

	// Middleware
	app.Use(logger.New())
	app.Use(cors.New())

	// Setup routes
	router.SetupRoutes(app, db)

	// Start server
	log.Fatal(app.Listen(":3000"))
}

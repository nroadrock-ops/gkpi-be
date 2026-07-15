package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

// CORS returns a configured CORS middleware.
// Middleware ini membaca ALLOWED_ORIGINS, jika kosong default ke http://localhost:3000
// dan mengizinkan origin yang match (bukan wildcard * agar support cookie/auth)
func CORS(allowedOrigins string) fiber.Handler {
	// Fallback jika tidak ada konfigurasi di env
	if allowedOrigins == "" {
		allowedOrigins = "http://localhost:3000"
	}

	// Membersihkan spasi jika ada di antara koma (contoh: "url1, url2")
	originsList := strings.Split(allowedOrigins, ",")
	for i, origin := range originsList {
		originsList[i] = strings.TrimSpace(origin)
	}
	cleanOrigins := strings.Join(originsList, ",")

	return cors.New(cors.Config{
		AllowOrigins:     cleanOrigins, // Kumpulan origin yang diizinkan, dipisahkan koma
		AllowCredentials: true,         // Izinkan credentials (cookies/auth headers)
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET, POST, HEAD, PUT, DELETE, PATCH, OPTIONS",
	})
}

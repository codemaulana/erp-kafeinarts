package middleware

import (
	"erp-kafeinart/config"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
)

// SetupCORS mengonfigurasi middleware Cross-Origin Resource Sharing (CORS) untuk Fiber v3.
// Memastikan frontend (misal Quasar di port 9000/9001) dapat berinteraksi secara aman dengan cookie & header otorisasi.
func SetupCORS(cfg *config.Config) fiber.Handler {
	origins := cfg.CORSOrigins
	if len(origins) == 0 {
		origins = []string{"*"}
	}

	return cors.New(cors.Config{
		AllowOrigins:     origins,
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With"},
		AllowMethods:     []string{"GET", "POST", "HEAD", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowCredentials: true,
	})
}

package middleware

import (
	"strings"

	"erp-kafeinart/config"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
)

// Protected mendefinisikan middleware untuk memvalidasi token JWT pada header Authorization.
// Jika token valid, klaim user disimpan di c.Locals("user") untuk digunakan di handler.
func Protected(cfg *config.Config) fiber.Handler {
	return func(c fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Akses ditolak: token otentikasi tidak ditemukan")
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if len(parts) != 2 || strings.ToLower(parts[0]) != "bearer" {
			return utils.SendError(c, fiber.StatusUnauthorized, "Format token otentikasi tidak valid (harus Bearer <token>)")
		}

		tokenString := strings.TrimSpace(parts[1])
		claims, err := utils.ValidateJWT(tokenString, cfg.JWTSecret)
		if err != nil {
			return utils.SendError(c, fiber.StatusUnauthorized, "Sesi Anda telah berakhir atau token tidak valid", err.Error())
		}

		// Simpan claims user di context
		c.Locals("user", claims)
		return c.Next()
	}
}

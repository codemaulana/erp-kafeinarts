package routes

import (
	"erp-kafeinart/handler"

	"github.com/gofiber/fiber/v3"
)

// RegisterAuthRoutes mendaftarkan seluruh endpoint autentikasi, akun, dan profil pengguna.
func RegisterAuthRoutes(api fiber.Router, app *fiber.App, authHandler *handler.AuthHandler, jwtAuth fiber.Handler) {
	// Endpoint Login (Publik)
	api.Post("/auth/login", authHandler.Login)
	api.Post("/login", authHandler.Login) // Alias /api/v1/login
	app.Post("/login", authHandler.Login) // Alias root /login

	// Rute User & Profil (Terproteksi JWT)
	api.Get("/auth/me", jwtAuth, authHandler.GetProfile)
	api.Get("/user", jwtAuth, authHandler.GetProfile) // Alias /api/v1/user
	app.Get("/user", jwtAuth, authHandler.GetProfile) // Alias root /user

	// Rute Logout (Publik & Terproteksi) agar selalu sukses 200 OK bahkan jika token sudah kedaluwarsa
	api.Post("/auth/logout", authHandler.Logout)
	api.Post("/logout", authHandler.Logout)
	app.Post("/logout", authHandler.Logout)

	// Rute kompatibilitas Sanctum SPA cookie
	app.Get("/sanctum/csrf-cookie", authHandler.SanctumCookie)
}

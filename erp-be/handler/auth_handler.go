package handler

import (
	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// LoginRequest merepresentasikan payload saat request login.
type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Remember bool   `json:"remember"`
}

// Login menangani POST /api/v1/auth/login dan POST /login
func (h *AuthHandler) Login(c fiber.Ctx) error {
	var req LoginRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data input tidak valid", err.Error())
	}

	res, err := h.authService.Login(req.Email, req.Password)
	if err != nil {
		return utils.SendError(c, fiber.StatusUnauthorized, err.Error())
	}

	// Response langsung memuat user dan token sesuai struktur di auth.module.js:
	// commit('SET_AUTH', { user: data.user, token: data.token })
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Login berhasil",
		"token":   res.Token,
		"user":    res.User,
		"data":    res.User, // fallback untuk berbagai format konsumsi
	})
}

// GetProfile menangani GET /api/v1/auth/me dan GET /user
func (h *AuthHandler) GetProfile(c fiber.Ctx) error {
	userClaims, ok := c.Locals("user").(*utils.JWTClaims)
	if !ok || userClaims == nil {
		return utils.SendError(c, fiber.StatusUnauthorized, "Sesi otentikasi tidak ditemukan")
	}

	user, err := h.authService.GetProfile(userClaims.UserID)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	// Format response { "data": user } sesuai ekspektasi fetchProfile di auth.module.js:
	// commit('SET_AUTH', { user: data.data, token: state.token })
	return utils.SendSuccess(c, fiber.StatusOK, user, "Profil user berhasil diambil")
}

// Logout menangani POST /api/v1/auth/logout
func (h *AuthHandler) Logout(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "Logout berhasil",
	})
}

// SanctumCookie menangani GET /sanctum/csrf-cookie untuk kompatibilitas SPA fallback
func (h *AuthHandler) SanctumCookie(c fiber.Ctx) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"message": "CSRF cookie set",
	})
}

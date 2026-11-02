package handler

import (
	"strconv"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
)

type UserHandler struct {
	userService service.UserService
}

func NewUserHandler(userService service.UserService) *UserHandler {
	return &UserHandler{userService: userService}
}

// GetUsers menangani GET /api/v1/users
func (h *UserHandler) GetUsers(c fiber.Ctx) error {
	search := c.Query("search")
	role := c.Query("role")

	users, err := h.userService.GetUsers(search, role)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat daftar pengguna", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, users, "Daftar pengguna berhasil dimuat")
}

// GetUserByID menangani GET /api/v1/users/:id
func (h *UserHandler) GetUserByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID pengguna tidak valid")
	}

	user, err := h.userService.GetUserByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, user, "Detail pengguna berhasil dimuat")
}

// CreateUser menangani POST /api/v1/users
func (h *UserHandler) CreateUser(c fiber.Ctx) error {
	var dto service.CreateUserDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	user, err := h.userService.CreateUser(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, user, "Pengguna berhasil dibuat")
}

// UpdateUser menangani PUT /api/v1/users/:id
func (h *UserHandler) UpdateUser(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID pengguna tidak valid")
	}

	var dto service.UpdateUserDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format payload tidak valid", err.Error())
	}

	user, err := h.userService.UpdateUser(uint(id), &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, user, "Pengguna berhasil diperbarui")
}

// DeleteUser menangani DELETE /api/v1/users/:id
func (h *UserHandler) DeleteUser(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID pengguna tidak valid")
	}

	var currentUserID uint = 0
	if claims, ok := c.Locals("user").(*utils.JWTClaims); ok && claims != nil {
		currentUserID = claims.UserID
	}

	if err := h.userService.DeleteUser(uint(id), currentUserID); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Pengguna berhasil dihapus")
}

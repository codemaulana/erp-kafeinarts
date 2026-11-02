package utils

import (
	"math"

	"github.com/gofiber/fiber/v3"
)

// PaginationMeta memuat metadata pagination yang fleksibel untuk frontend Quasar.
// Menyediakan properti `page` dan `current_page` agar kompatibel dengan seluruh module FE.
type PaginationMeta struct {
	Total       int64 `json:"total"`
	Page        int   `json:"page"`
	CurrentPage int   `json:"current_page"`
	PerPage     int   `json:"per_page"`
	PerPageAlt  int   `json:"perPage"`
	TotalPages  int   `json:"total_pages"`
	LastPage    int   `json:"last_page"`
	LastPageAlt int   `json:"lastPage"`
}

// BuildPaginationMeta menghitung total halaman dan menyusun PaginationMeta.
func BuildPaginationMeta(total int64, page, perPage int) PaginationMeta {
	if perPage <= 0 {
		perPage = 10
	}
	totalPages := int(math.Ceil(float64(total) / float64(perPage)))
	if totalPages == 0 {
		totalPages = 1
	}

	return PaginationMeta{
		Total:       total,
		Page:        page,
		CurrentPage: page,
		PerPage:     perPage,
		PerPageAlt:  perPage,
		TotalPages:  totalPages,
		LastPage:    totalPages,
		LastPageAlt: totalPages,
	}
}

// SendSuccess mengirimkan response HTTP sukses dengan data.
func SendSuccess(c fiber.Ctx, statusCode int, data interface{}, message ...string) error {
	msg := "Operasi berhasil"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
	})
}

// SendPaginated mengirimkan response sukses berisi list data dan metadata pagination.
func SendPaginated(c fiber.Ctx, statusCode int, data interface{}, meta PaginationMeta, message ...string) error {
	msg := "Data berhasil diambil"
	if len(message) > 0 && message[0] != "" {
		msg = message[0]
	}

	return c.Status(statusCode).JSON(fiber.Map{
		"success": true,
		"message": msg,
		"data":    data,
		"meta":    meta,
	})
}

// SendError mengirimkan response HTTP gagal / error terstruktur.
func SendError(c fiber.Ctx, statusCode int, message string, detail ...string) error {
	res := fiber.Map{
		"success": false,
		"message": message,
	}

	if len(detail) > 0 && detail[0] != "" {
		res["error"] = detail[0]
	}

	return c.Status(statusCode).JSON(res)
}

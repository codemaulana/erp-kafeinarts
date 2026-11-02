package handler

import (
	"strconv"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
)

type EmployeeHandler struct {
	empService service.EmployeeService
}

func NewEmployeeHandler(empService service.EmployeeService) *EmployeeHandler {
	return &EmployeeHandler{empService: empService}
}

// GetAll menangani GET /api/v1/hr/employees
func (h *EmployeeHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("rowsPerPage", "15")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 15
	}

	search := c.Query("search")
	department := c.Query("department")
	status := c.Query("status")

	employees, meta, err := h.empService.GetAll(page, perPage, search, department, status)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar pegawai", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, employees, meta, "Daftar pegawai berhasil dimuat")
}

// GetByID menangani GET /api/v1/hr/employees/:id
func (h *EmployeeHandler) GetByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	emp, err := h.empService.GetByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, emp, "Detail pegawai berhasil dimuat")
}

// Create menangani POST /api/v1/hr/employees
func (h *EmployeeHandler) Create(c fiber.Ctx) error {
	var req service.CreateEmployeeDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input pegawai tidak valid", err.Error())
	}

	emp, err := h.empService.Create(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, emp, "Pegawai baru berhasil ditambahkan")
}

// Update menangani PUT /api/v1/hr/employees/:id
func (h *EmployeeHandler) Update(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	var req service.UpdateEmployeeDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input pegawai tidak valid", err.Error())
	}

	emp, err := h.empService.Update(uint(id), req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, emp, "Data pegawai berhasil diperbarui")
}

// Delete menangani DELETE /api/v1/hr/employees/:id
func (h *EmployeeHandler) Delete(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID pegawai tidak valid")
	}

	if err := h.empService.Delete(uint(id)); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Pegawai berhasil dihapus")
}

package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type HrHandler struct {
	hrService service.HrService
}

func NewHrHandler(hrService service.HrService) *HrHandler {
	return &HrHandler{hrService: hrService}
}

// ── Master Departemen Handlers ──────────────────────────────────────

func (h *HrHandler) GetDepartments(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.hrService.GetDepartmentsPaginated(page, perPage, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat departemen", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.hrService.GetAllDepartments(search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat departemen", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar departemen berhasil dimuat")
}

func (h *HrHandler) GetDepartmentByID(c fiber.Ctx) error {
	id := c.Params("id")
	dept, err := h.hrService.GetDepartmentByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Departemen tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, dept, "Detail departemen berhasil dimuat")
}

func (h *HrHandler) CreateDepartment(c fiber.Ctx) error {
	var dto service.DepartmentDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	dept, err := h.hrService.CreateDepartment(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, dept, "Departemen berhasil dibuat")
}

func (h *HrHandler) UpdateDepartment(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.DepartmentDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	dept, err := h.hrService.UpdateDepartment(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, dept, "Departemen berhasil diperbarui")
}

func (h *HrHandler) DeleteDepartment(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteDepartment(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus departemen", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Departemen berhasil dihapus")
}

// ── Master Jabatan Handlers ─────────────────────────────────────────

func (h *HrHandler) GetPositions(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	deptID := c.Query("department_id")

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.hrService.GetPositionsPaginated(page, perPage, deptID, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jabatan", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.hrService.GetAllPositions(deptID, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jabatan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar jabatan berhasil dimuat")
}

func (h *HrHandler) GetPositionByID(c fiber.Ctx) error {
	id := c.Params("id")
	pos, err := h.hrService.GetPositionByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Jabatan tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pos, "Detail jabatan berhasil dimuat")
}

func (h *HrHandler) CreatePosition(c fiber.Ctx) error {
	var dto service.PositionDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	pos, err := h.hrService.CreatePosition(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, pos, "Jabatan berhasil dibuat")
}

func (h *HrHandler) UpdatePosition(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.PositionDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	pos, err := h.hrService.UpdatePosition(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pos, "Jabatan berhasil diperbarui")
}

func (h *HrHandler) DeletePosition(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeletePosition(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus jabatan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Jabatan berhasil dihapus")
}

// ── Master Karyawan Handlers ────────────────────────────────────────

func (h *HrHandler) GetEmployees(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	deptID := c.Query("department_id")
	status := c.Query("status")
	if status == "" {
		status = c.Query("active_status")
	}

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.hrService.GetEmployeesPaginated(page, perPage, deptID, status, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat karyawan", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.hrService.GetAllEmployees(deptID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat karyawan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar karyawan berhasil dimuat")
}

func (h *HrHandler) GetEmployeeByID(c fiber.Ctx) error {
	id := c.Params("id")
	emp, err := h.hrService.GetEmployeeByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Karyawan tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, emp, "Detail karyawan berhasil dimuat")
}

func (h *HrHandler) CreateEmployee(c fiber.Ctx) error {
	var dto service.EmployeeDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	emp, err := h.hrService.CreateEmployee(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, emp, "Karyawan berhasil ditambahkan")
}

func (h *HrHandler) UpdateEmployee(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.EmployeeDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	emp, err := h.hrService.UpdateEmployee(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, emp, "Karyawan berhasil diperbarui")
}

func (h *HrHandler) DeleteEmployee(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteEmployee(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus karyawan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Karyawan berhasil dihapus")
}

// ── Master Jenis Cuti Handlers ──────────────────────────────────────

func (h *HrHandler) GetLeaveTypes(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.hrService.GetLeaveTypesPaginated(page, perPage, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jenis cuti", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.hrService.GetAllLeaveTypes(search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jenis cuti", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar jenis cuti berhasil dimuat")
}

func (h *HrHandler) GetLeaveTypeByID(c fiber.Ctx) error {
	id := c.Params("id")
	lt, err := h.hrService.GetLeaveTypeByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Jenis cuti tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, lt, "Detail jenis cuti berhasil dimuat")
}

func (h *HrHandler) CreateLeaveType(c fiber.Ctx) error {
	var dto service.LeaveTypeDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	lt, err := h.hrService.CreateLeaveType(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, lt, "Jenis cuti berhasil dibuat")
}

func (h *HrHandler) UpdateLeaveType(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.LeaveTypeDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	lt, err := h.hrService.UpdateLeaveType(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, lt, "Jenis cuti berhasil diperbarui")
}

func (h *HrHandler) DeleteLeaveType(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteLeaveType(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus jenis cuti", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Jenis cuti berhasil dihapus")
}

// ── Master Komponen Gaji Handlers ───────────────────────────────────

func (h *HrHandler) GetSalaryComponents(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	cType := c.Query("component_type")
	if cType == "" {
		cType = c.Query("type")
	}

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.hrService.GetSalaryComponentsPaginated(page, perPage, cType, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat komponen gaji", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.hrService.GetAllSalaryComponents(cType, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat komponen gaji", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar komponen gaji berhasil dimuat")
}

func (h *HrHandler) GetSalaryComponentByID(c fiber.Ctx) error {
	id := c.Params("id")
	sc, err := h.hrService.GetSalaryComponentByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Komponen gaji tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, sc, "Detail komponen gaji berhasil dimuat")
}

func (h *HrHandler) CreateSalaryComponent(c fiber.Ctx) error {
	var dto service.SalaryComponentDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	sc, err := h.hrService.CreateSalaryComponent(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, sc, "Komponen gaji berhasil dibuat")
}

func (h *HrHandler) UpdateSalaryComponent(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.SalaryComponentDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	sc, err := h.hrService.UpdateSalaryComponent(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, sc, "Komponen gaji berhasil diperbarui")
}

func (h *HrHandler) DeleteSalaryComponent(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteSalaryComponent(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus komponen gaji", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Komponen gaji berhasil dihapus")
}

// ── Master Jadwal Kerja Handlers ────────────────────────────────────

func (h *HrHandler) GetWorkShifts(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.hrService.GetWorkShiftsPaginated(page, perPage, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jadwal kerja", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.hrService.GetAllWorkShifts(search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jadwal kerja", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar jadwal kerja berhasil dimuat")
}

func (h *HrHandler) GetWorkShiftByID(c fiber.Ctx) error {
	id := c.Params("id")
	ws, err := h.hrService.GetWorkShiftByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Jadwal kerja tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, ws, "Detail jadwal kerja berhasil dimuat")
}

func (h *HrHandler) CreateWorkShift(c fiber.Ctx) error {
	var dto service.WorkShiftDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	ws, err := h.hrService.CreateWorkShift(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, ws, "Jadwal kerja berhasil dibuat")
}

func (h *HrHandler) UpdateWorkShift(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.WorkShiftDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	ws, err := h.hrService.UpdateWorkShift(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, ws, "Jadwal kerja berhasil diperbarui")
}

func (h *HrHandler) DeleteWorkShift(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteWorkShift(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus jadwal kerja", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Jadwal kerja berhasil dihapus")
}

// ── Transaksi Presensi Handlers ─────────────────────────────────────

func (h *HrHandler) GetAttendances(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	empID := c.Query("employee_id")
	date := c.Query("date")
	if date == "" {
		date = c.Query("attendance_date")
	}
	status := c.Query("status")
	if status == "" {
		status = c.Query("attendance_status")
	}
	search := c.Query("search")

	list, total, err := h.hrService.GetAttendancesPaginated(page, perPage, empID, date, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat presensi", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *HrHandler) GetAttendanceByID(c fiber.Ctx) error {
	id := c.Params("id")
	att, err := h.hrService.GetAttendanceByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Presensi tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, att, "Detail presensi berhasil dimuat")
}

func (h *HrHandler) CreateAttendance(c fiber.Ctx) error {
	var dto service.AttendanceDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	att, err := h.hrService.CreateAttendance(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, att, "Presensi berhasil dicatat")
}

func (h *HrHandler) UpdateAttendance(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.AttendanceDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	att, err := h.hrService.UpdateAttendance(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, att, "Presensi berhasil diperbarui")
}

func (h *HrHandler) DeleteAttendance(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteAttendance(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus presensi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Presensi berhasil dihapus")
}

// ── Transaksi Pengajuan Cuti Handlers ───────────────────────────────

func (h *HrHandler) GetLeaveRequests(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	empID := c.Query("employee_id")
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.hrService.GetLeaveRequestsPaginated(page, perPage, empID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat pengajuan cuti", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *HrHandler) GetLeaveRequestByID(c fiber.Ctx) error {
	id := c.Params("id")
	lr, err := h.hrService.GetLeaveRequestByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Pengajuan cuti tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, lr, "Detail pengajuan cuti berhasil dimuat")
}

func (h *HrHandler) CreateLeaveRequest(c fiber.Ctx) error {
	var dto service.LeaveRequestDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	lr, err := h.hrService.CreateLeaveRequest(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, lr, "Pengajuan cuti berhasil diajukan")
}

func (h *HrHandler) UpdateLeaveRequest(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.LeaveRequestDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	lr, err := h.hrService.UpdateLeaveRequest(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, lr, "Pengajuan cuti berhasil diperbarui")
}

func (h *HrHandler) UpdateLeaveRequestStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status          string `json:"status"`
		RejectionReason string `json:"rejection_reason"`
		ApprovedBy      string `json:"approved_by"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.hrService.UpdateLeaveRequestStatus(id, req.Status, req.RejectionReason, req.ApprovedBy); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.hrService.GetLeaveRequestByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status cuti berhasil diperbarui")
}

func (h *HrHandler) DeleteLeaveRequest(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteLeaveRequest(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus cuti", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Pengajuan cuti berhasil dihapus")
}

// ── Transaksi Payroll Handlers ──────────────────────────────────────

func (h *HrHandler) GetPayrolls(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	period := c.Query("period")
	if period == "" {
		period = c.Query("payroll_period")
	}
	empID := c.Query("employee_id")
	status := c.Query("status")
	if status == "" {
		status = c.Query("payment_status")
	}
	search := c.Query("search")

	list, total, err := h.hrService.GetPayrollsPaginated(page, perPage, period, empID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat payroll", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *HrHandler) GetPayrollByID(c fiber.Ctx) error {
	id := c.Params("id")
	p, err := h.hrService.GetPayrollByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Payroll tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, p, "Detail payroll berhasil dimuat")
}

func (h *HrHandler) CreatePayroll(c fiber.Ctx) error {
	var dto service.PayrollDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	p, err := h.hrService.CreatePayroll(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, p, "Payroll berhasil dibuat")
}

func (h *HrHandler) UpdatePayroll(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.PayrollDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	p, err := h.hrService.UpdatePayroll(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, p, "Payroll berhasil diperbarui")
}

func (h *HrHandler) UpdatePayrollStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.hrService.UpdatePayrollStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.hrService.GetPayrollByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status payroll berhasil diperbarui")
}

func (h *HrHandler) DeletePayroll(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeletePayroll(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus payroll", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Payroll berhasil dihapus")
}

// ── Transaksi Reimbursement Handlers ────────────────────────────────

func (h *HrHandler) GetReimbursements(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	empID := c.Query("employee_id")
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.hrService.GetReimbursementsPaginated(page, perPage, empID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat reimbursement", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *HrHandler) GetReimbursementByID(c fiber.Ctx) error {
	id := c.Params("id")
	rm, err := h.hrService.GetReimbursementByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Reimbursement tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, rm, "Detail reimbursement berhasil dimuat")
}

func (h *HrHandler) CreateReimbursement(c fiber.Ctx) error {
	var dto service.ReimbursementDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	rm, err := h.hrService.CreateReimbursement(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, rm, "Reimbursement berhasil diajukan")
}

func (h *HrHandler) UpdateReimbursement(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.ReimbursementDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	rm, err := h.hrService.UpdateReimbursement(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, rm, "Reimbursement berhasil diperbarui")
}

func (h *HrHandler) UpdateReimbursementStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.hrService.UpdateReimbursementStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.hrService.GetReimbursementByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status reimbursement berhasil diperbarui")
}

func (h *HrHandler) DeleteReimbursement(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeleteReimbursement(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus reimbursement", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Reimbursement berhasil dihapus")
}

// ── Transaksi Penilaian Kinerja Handlers ────────────────────────────

func (h *HrHandler) GetPerformanceReviews(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	empID := c.Query("employee_id")
	period := c.Query("period")
	if period == "" {
		period = c.Query("review_period")
	}
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.hrService.GetPerformanceReviewsPaginated(page, perPage, empID, period, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat penilaian kinerja", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *HrHandler) GetPerformanceReviewByID(c fiber.Ctx) error {
	id := c.Params("id")
	pr, err := h.hrService.GetPerformanceReviewByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Penilaian kinerja tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pr, "Detail penilaian kinerja berhasil dimuat")
}

func (h *HrHandler) CreatePerformanceReview(c fiber.Ctx) error {
	var dto service.PerformanceReviewDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	pr, err := h.hrService.CreatePerformanceReview(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, pr, "Penilaian kinerja berhasil dibuat")
}

func (h *HrHandler) UpdatePerformanceReview(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.PerformanceReviewDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	pr, err := h.hrService.UpdatePerformanceReview(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pr, "Penilaian kinerja berhasil diperbarui")
}

func (h *HrHandler) UpdatePerformanceReviewStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.hrService.UpdatePerformanceReviewStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.hrService.GetPerformanceReviewByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status penilaian kinerja berhasil diperbarui")
}

func (h *HrHandler) DeletePerformanceReview(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.hrService.DeletePerformanceReview(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus penilaian kinerja", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Penilaian kinerja berhasil dihapus")
}

// ── File Upload Handler ─────────────────────────────────────────────

func (h *HrHandler) UploadFile(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "File tidak ditemukan", err.Error())
	}

	uploadDir := "./storage/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat folder upload", err.Error())
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("hr_%s_%s%s", time.Now().Format("20060102150405"), uuid.NewString()[:8], ext)
	dst := filepath.Join(uploadDir, newFileName)

	if err := c.SaveFile(file, dst); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan file", err.Error())
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)
	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{
		"file_name": file.Filename,
		"file_url":  fileURL,
		"size":      file.Size,
	}, "File berhasil diunggah")
}

// ── Laporan HR Handlers ─────────────────────────────────────────────

func (h *HrHandler) GetRekapAbsensiReport(c fiber.Ctx) error {
	start := c.Query("start_date")
	if start == "" {
		start = c.Query("periodeStart")
	}
	end := c.Query("end_date")
	if end == "" {
		end = c.Query("periodeEnd")
	}
	deptID := c.Query("department_id")
	if deptID == "" {
		deptID = c.Query("departemen")
	}
	empID := c.Query("employee_id")
	if empID == "" {
		empID = c.Query("karyawan")
	}

	report, err := h.hrService.GetRekapAbsensi(start, end, deptID, empID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat rekap absensi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan rekap absensi berhasil dimuat")
}

func (h *HrHandler) GetSlipGajiReport(c fiber.Ctx) error {
	id := c.Params("id")
	if id != "" {
		p, err := h.hrService.GetPayrollByID(id)
		if err != nil {
			return utils.SendError(c, fiber.StatusNotFound, "Slip gaji tidak ditemukan", err.Error())
		}
		return utils.SendSuccess(c, fiber.StatusOK, p, "Slip gaji berhasil dimuat")
	}

	period := c.Query("period")
	if period == "" {
		period = c.Query("periode")
	}
	empID := c.Query("employee_id")
	if empID == "" {
		empID = c.Query("karyawan")
	}
	deptID := c.Query("department_id")
	if deptID == "" {
		deptID = c.Query("departemen")
	}

	report, err := h.hrService.GetSlipGaji(period, empID, deptID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat slip gaji", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan slip gaji berhasil dimuat")
}

func (h *HrHandler) GetSaldoCutiReport(c fiber.Ctx) error {
	deptID := c.Query("department_id")
	if deptID == "" {
		deptID = c.Query("departemen")
	}
	empID := c.Query("employee_id")
	if empID == "" {
		empID = c.Query("karyawan")
	}

	report, err := h.hrService.GetSaldoCuti(deptID, empID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat saldo cuti", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan saldo cuti berhasil dimuat")
}

func (h *HrHandler) GetTurnoverReport(c fiber.Ctx) error {
	yearStr := c.Query("year")
	if yearStr == "" {
		yearStr = c.Query("tahun")
	}
	year, _ := strconv.Atoi(yearStr)
	if year <= 0 {
		year = time.Now().Year()
	}

	report, err := h.hrService.GetTurnover(year)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat laporan turnover", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan turnover berhasil dimuat")
}

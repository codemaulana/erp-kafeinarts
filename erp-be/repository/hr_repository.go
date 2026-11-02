package repository

import (
	"fmt"
	"time"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

type HrRepository interface {
	// ── Master Departemen ───────────────────────────
	FindDepartments(search string) ([]model.MstHrDepartment, error)
	FindDepartmentsPaginated(page, perPage int, search string) ([]model.MstHrDepartment, int64, error)
	FindDepartmentByID(id string) (*model.MstHrDepartment, error)
	CreateDepartment(dept *model.MstHrDepartment) error
	UpdateDepartment(dept *model.MstHrDepartment) error
	DeleteDepartment(id string) error

	// ── Master Jabatan ──────────────────────────────
	FindPositions(departmentID, search string) ([]model.MstHrPosition, error)
	FindPositionsPaginated(page, perPage int, departmentID, search string) ([]model.MstHrPosition, int64, error)
	FindPositionByID(id string) (*model.MstHrPosition, error)
	CreatePosition(pos *model.MstHrPosition) error
	UpdatePosition(pos *model.MstHrPosition) error
	DeletePosition(id string) error

	// ── Master Karyawan ─────────────────────────────
	FindEmployees(departmentID, status, search string) ([]model.MstHrEmployee, error)
	FindEmployeesPaginated(page, perPage int, departmentID, status, search string) ([]model.MstHrEmployee, int64, error)
	FindEmployeeByID(id string) (*model.MstHrEmployee, error)
	CreateEmployee(emp *model.MstHrEmployee) error
	UpdateEmployee(emp *model.MstHrEmployee) error
	DeleteEmployee(id string) error
	CountEmployeesByStatus(status string) (int64, error)

	// ── Master Jenis Cuti ───────────────────────────
	FindLeaveTypes(search string) ([]model.MstHrLeaveType, error)
	FindLeaveTypesPaginated(page, perPage int, search string) ([]model.MstHrLeaveType, int64, error)
	FindLeaveTypeByID(id string) (*model.MstHrLeaveType, error)
	CreateLeaveType(lt *model.MstHrLeaveType) error
	UpdateLeaveType(lt *model.MstHrLeaveType) error
	DeleteLeaveType(id string) error

	// ── Master Komponen Gaji ────────────────────────
	FindSalaryComponents(compType, search string) ([]model.MstHrSalaryComponent, error)
	FindSalaryComponentsPaginated(page, perPage int, compType, search string) ([]model.MstHrSalaryComponent, int64, error)
	FindSalaryComponentByID(id string) (*model.MstHrSalaryComponent, error)
	CreateSalaryComponent(sc *model.MstHrSalaryComponent) error
	UpdateSalaryComponent(sc *model.MstHrSalaryComponent) error
	DeleteSalaryComponent(id string) error

	// ── Master Jadwal Kerja ─────────────────────────
	FindWorkShifts(search string) ([]model.MstHrWorkShift, error)
	FindWorkShiftsPaginated(page, perPage int, search string) ([]model.MstHrWorkShift, int64, error)
	FindWorkShiftByID(id string) (*model.MstHrWorkShift, error)
	CreateWorkShift(ws *model.MstHrWorkShift) error
	UpdateWorkShift(ws *model.MstHrWorkShift) error
	DeleteWorkShift(id string) error

	// ── Transaksi Presensi ──────────────────────────
	FindAttendancesPaginated(page, perPage int, employeeID, date, status, search string) ([]model.TrxHrAttendance, int64, error)
	FindAttendanceByID(id string) (*model.TrxHrAttendance, error)
	FindAttendanceByEmployeeAndDate(employeeID, date string) (*model.TrxHrAttendance, error)
	CreateAttendance(att *model.TrxHrAttendance) error
	UpdateAttendance(att *model.TrxHrAttendance) error
	DeleteAttendance(id string) error

	// ── Transaksi Pengajuan Cuti ────────────────────
	FindLeaveRequestsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrLeaveRequest, int64, error)
	FindLeaveRequestByID(id string) (*model.TrxHrLeaveRequest, error)
	CreateLeaveRequest(lr *model.TrxHrLeaveRequest) error
	UpdateLeaveRequest(lr *model.TrxHrLeaveRequest) error
	UpdateLeaveRequestStatus(id, status, rejectionReason, approvedBy string) error
	DeleteLeaveRequest(id string) error

	// ── Transaksi Payroll ───────────────────────────
	FindPayrollsPaginated(page, perPage int, period, employeeID, status, search string) ([]model.TrxHrPayroll, int64, error)
	FindPayrollByID(id string) (*model.TrxHrPayroll, error)
	CreatePayroll(p *model.TrxHrPayroll) error
	UpdatePayroll(p *model.TrxHrPayroll) error
	UpdatePayrollStatus(id, status string) error
	DeletePayroll(id string) error

	// ── Transaksi Reimbursement ─────────────────────
	FindReimbursementsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrReimbursement, int64, error)
	FindReimbursementByID(id string) (*model.TrxHrReimbursement, error)
	CreateReimbursement(r *model.TrxHrReimbursement) error
	UpdateReimbursement(r *model.TrxHrReimbursement) error
	UpdateReimbursementStatus(id, status string) error
	DeleteReimbursement(id string) error

	// ── Transaksi Penilaian Kinerja ─────────────────
	FindPerformanceReviewsPaginated(page, perPage int, employeeID, period, status, search string) ([]model.TrxHrPerformanceReview, int64, error)
	FindPerformanceReviewByID(id string) (*model.TrxHrPerformanceReview, error)
	CreatePerformanceReview(pr *model.TrxHrPerformanceReview) error
	UpdatePerformanceReview(pr *model.TrxHrPerformanceReview) error
	UpdatePerformanceReviewStatus(id, status string) error
	DeletePerformanceReview(id string) error

	// ── Laporan HR ──────────────────────────────────
	GetRekapAbsensi(startDate, endDate, departmentID, employeeID string) (*model.RekapAbsensiReport, error)
	GetSlipGaji(period, employeeID, departmentID string) (*model.SlipGajiReport, error)
	GetSaldoCuti(departmentID, employeeID string) (*model.SaldoCutiReport, error)
	GetTurnover(year int) (*model.TurnoverReport, error)
}

type hrRepository struct {
	db *gorm.DB
}

func NewHrRepository(db *gorm.DB) HrRepository {
	return &hrRepository{db: db}
}

// ── Master Departemen ───────────────────────────────────────────────

func (r *hrRepository) FindDepartments(search string) ([]model.MstHrDepartment, error) {
	var list []model.MstHrDepartment
	q := r.db.Model(&model.MstHrDepartment{})
	if search != "" {
		q = q.Where("department_name LIKE ? OR code LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Order("department_name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	r.populateDepartmentVirtuals(list)
	return list, nil
}

func (r *hrRepository) FindDepartmentsPaginated(page, perPage int, search string) ([]model.MstHrDepartment, int64, error) {
	var list []model.MstHrDepartment
	var total int64
	q := r.db.Model(&model.MstHrDepartment{})
	if search != "" {
		q = q.Where("department_name LIKE ? OR code LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("department_name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.populateDepartmentVirtuals(list)
	return list, total, nil
}

func (r *hrRepository) populateDepartmentVirtuals(list []model.MstHrDepartment) {
	for i := range list {
		if list[i].HeadEmployeeID != nil && *list[i].HeadEmployeeID != "" {
			var emp model.MstHrEmployee
			if err := r.db.Select("full_name").First(&emp, "id = ?", *list[i].HeadEmployeeID).Error; err == nil {
				list[i].HeadEmployeeName = emp.FullName
			}
		}
		var count int64
		r.db.Model(&model.MstHrEmployee{}).Where("department_id = ? AND active_status = 'ACTIVE'", list[i].ID).Count(&count)
		list[i].EmployeeCount = count
	}
}

func (r *hrRepository) FindDepartmentByID(id string) (*model.MstHrDepartment, error) {
	var dept model.MstHrDepartment
	if err := r.db.First(&dept, "id = ?", id).Error; err != nil {
		return nil, err
	}
	slice := []model.MstHrDepartment{dept}
	r.populateDepartmentVirtuals(slice)
	return &slice[0], nil
}

func (r *hrRepository) CreateDepartment(dept *model.MstHrDepartment) error {
	return r.db.Create(dept).Error
}

func (r *hrRepository) UpdateDepartment(dept *model.MstHrDepartment) error {
	return r.db.Model(&model.MstHrDepartment{}).Where("id = ?", dept.ID).Updates(dept).Error
}

func (r *hrRepository) DeleteDepartment(id string) error {
	return r.db.Delete(&model.MstHrDepartment{}, "id = ?", id).Error
}

// ── Master Jabatan ──────────────────────────────────────────────────

func (r *hrRepository) FindPositions(departmentID, search string) ([]model.MstHrPosition, error) {
	var list []model.MstHrPosition
	q := r.db.Model(&model.MstHrPosition{})
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if search != "" {
		q = q.Where("position_name LIKE ? OR code LIKE ? OR level LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Order("position_name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	r.populatePositionVirtuals(list)
	return list, nil
}

func (r *hrRepository) FindPositionsPaginated(page, perPage int, departmentID, search string) ([]model.MstHrPosition, int64, error) {
	var list []model.MstHrPosition
	var total int64
	q := r.db.Model(&model.MstHrPosition{})
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if search != "" {
		q = q.Where("position_name LIKE ? OR code LIKE ? OR level LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("position_name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.populatePositionVirtuals(list)
	return list, total, nil
}

func (r *hrRepository) populatePositionVirtuals(list []model.MstHrPosition) {
	for i := range list {
		if list[i].DepartmentID != "" {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", list[i].DepartmentID).Error; err == nil {
				list[i].DepartmentName = dept.DepartmentName
			}
		}
	}
}

func (r *hrRepository) FindPositionByID(id string) (*model.MstHrPosition, error) {
	var pos model.MstHrPosition
	if err := r.db.First(&pos, "id = ?", id).Error; err != nil {
		return nil, err
	}
	slice := []model.MstHrPosition{pos}
	r.populatePositionVirtuals(slice)
	return &slice[0], nil
}

func (r *hrRepository) CreatePosition(pos *model.MstHrPosition) error {
	return r.db.Create(pos).Error
}

func (r *hrRepository) UpdatePosition(pos *model.MstHrPosition) error {
	return r.db.Model(&model.MstHrPosition{}).Where("id = ?", pos.ID).Updates(pos).Error
}

func (r *hrRepository) DeletePosition(id string) error {
	return r.db.Delete(&model.MstHrPosition{}, "id = ?", id).Error
}

// ── Master Karyawan ─────────────────────────────────────────────────

func (r *hrRepository) FindEmployees(departmentID, status, search string) ([]model.MstHrEmployee, error) {
	var list []model.MstHrEmployee
	q := r.db.Model(&model.MstHrEmployee{})
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if status != "" {
		q = q.Where("active_status = ?", status)
	}
	if search != "" {
		q = q.Where("full_name LIKE ? OR employee_number LIKE ? OR email LIKE ? OR phone_number LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Order("full_name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	r.populateEmployeeVirtuals(list)
	return list, nil
}

func (r *hrRepository) FindEmployeesPaginated(page, perPage int, departmentID, status, search string) ([]model.MstHrEmployee, int64, error) {
	var list []model.MstHrEmployee
	var total int64
	q := r.db.Model(&model.MstHrEmployee{})
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if status != "" {
		q = q.Where("active_status = ?", status)
	}
	if search != "" {
		q = q.Where("full_name LIKE ? OR employee_number LIKE ? OR email LIKE ? OR phone_number LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("full_name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.populateEmployeeVirtuals(list)
	return list, total, nil
}

func (r *hrRepository) populateEmployeeVirtuals(list []model.MstHrEmployee) {
	for i := range list {
		if list[i].DepartmentID != nil && *list[i].DepartmentID != "" {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", *list[i].DepartmentID).Error; err == nil {
				list[i].DepartmentName = dept.DepartmentName
			}
		}
		if list[i].PositionID != nil && *list[i].PositionID != "" {
			var pos model.MstHrPosition
			if err := r.db.Select("position_name").First(&pos, "id = ?", *list[i].PositionID).Error; err == nil {
				list[i].PositionName = pos.PositionName
			}
		}
	}
}

func (r *hrRepository) FindEmployeeByID(id string) (*model.MstHrEmployee, error) {
	var emp model.MstHrEmployee
	if err := r.db.First(&emp, "id = ?", id).Error; err != nil {
		return nil, err
	}
	slice := []model.MstHrEmployee{emp}
	r.populateEmployeeVirtuals(slice)
	return &slice[0], nil
}

func (r *hrRepository) CreateEmployee(emp *model.MstHrEmployee) error {
	return r.db.Create(emp).Error
}

func (r *hrRepository) UpdateEmployee(emp *model.MstHrEmployee) error {
	return r.db.Model(&model.MstHrEmployee{}).Where("id = ?", emp.ID).Updates(emp).Error
}

func (r *hrRepository) DeleteEmployee(id string) error {
	return r.db.Delete(&model.MstHrEmployee{}, "id = ?", id).Error
}

func (r *hrRepository) CountEmployeesByStatus(status string) (int64, error) {
	var count int64
	err := r.db.Model(&model.MstHrEmployee{}).Where("active_status = ?", status).Count(&count).Error
	return count, err
}

// ── Master Jenis Cuti ───────────────────────────────────────────────

func (r *hrRepository) FindLeaveTypes(search string) ([]model.MstHrLeaveType, error) {
	var list []model.MstHrLeaveType
	q := r.db.Model(&model.MstHrLeaveType{})
	if search != "" {
		q = q.Where("leave_name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	err := q.Order("leave_name ASC").Find(&list).Error
	return list, err
}

func (r *hrRepository) FindLeaveTypesPaginated(page, perPage int, search string) ([]model.MstHrLeaveType, int64, error) {
	var list []model.MstHrLeaveType
	var total int64
	q := r.db.Model(&model.MstHrLeaveType{})
	if search != "" {
		q = q.Where("leave_name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	err := q.Order("leave_name ASC").Offset(offset).Limit(perPage).Find(&list).Error
	return list, total, err
}

func (r *hrRepository) FindLeaveTypeByID(id string) (*model.MstHrLeaveType, error) {
	var lt model.MstHrLeaveType
	if err := r.db.First(&lt, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &lt, nil
}

func (r *hrRepository) CreateLeaveType(lt *model.MstHrLeaveType) error {
	return r.db.Create(lt).Error
}

func (r *hrRepository) UpdateLeaveType(lt *model.MstHrLeaveType) error {
	return r.db.Model(&model.MstHrLeaveType{}).Where("id = ?", lt.ID).Updates(lt).Error
}

func (r *hrRepository) DeleteLeaveType(id string) error {
	return r.db.Delete(&model.MstHrLeaveType{}, "id = ?", id).Error
}

// ── Master Komponen Gaji ────────────────────────────────────────────

func (r *hrRepository) FindSalaryComponents(compType, search string) ([]model.MstHrSalaryComponent, error) {
	var list []model.MstHrSalaryComponent
	q := r.db.Model(&model.MstHrSalaryComponent{})
	if compType != "" {
		q = q.Where("component_type = ?", compType)
	}
	if search != "" {
		q = q.Where("component_name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	err := q.Order("component_name ASC").Find(&list).Error
	return list, err
}

func (r *hrRepository) FindSalaryComponentsPaginated(page, perPage int, compType, search string) ([]model.MstHrSalaryComponent, int64, error) {
	var list []model.MstHrSalaryComponent
	var total int64
	q := r.db.Model(&model.MstHrSalaryComponent{})
	if compType != "" {
		q = q.Where("component_type = ?", compType)
	}
	if search != "" {
		q = q.Where("component_name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	err := q.Order("component_name ASC").Offset(offset).Limit(perPage).Find(&list).Error
	return list, total, err
}

func (r *hrRepository) FindSalaryComponentByID(id string) (*model.MstHrSalaryComponent, error) {
	var sc model.MstHrSalaryComponent
	if err := r.db.First(&sc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &sc, nil
}

func (r *hrRepository) CreateSalaryComponent(sc *model.MstHrSalaryComponent) error {
	return r.db.Create(sc).Error
}

func (r *hrRepository) UpdateSalaryComponent(sc *model.MstHrSalaryComponent) error {
	return r.db.Model(&model.MstHrSalaryComponent{}).Where("id = ?", sc.ID).Updates(sc).Error
}

func (r *hrRepository) DeleteSalaryComponent(id string) error {
	return r.db.Delete(&model.MstHrSalaryComponent{}, "id = ?", id).Error
}

// ── Master Jadwal Kerja ─────────────────────────────────────────────

func (r *hrRepository) FindWorkShifts(search string) ([]model.MstHrWorkShift, error) {
	var list []model.MstHrWorkShift
	q := r.db.Model(&model.MstHrWorkShift{})
	if search != "" {
		q = q.Where("shift_name LIKE ?", "%"+search+"%")
	}
	err := q.Order("shift_name ASC").Find(&list).Error
	return list, err
}

func (r *hrRepository) FindWorkShiftsPaginated(page, perPage int, search string) ([]model.MstHrWorkShift, int64, error) {
	var list []model.MstHrWorkShift
	var total int64
	q := r.db.Model(&model.MstHrWorkShift{})
	if search != "" {
		q = q.Where("shift_name LIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	err := q.Order("shift_name ASC").Offset(offset).Limit(perPage).Find(&list).Error
	return list, total, err
}

func (r *hrRepository) FindWorkShiftByID(id string) (*model.MstHrWorkShift, error) {
	var ws model.MstHrWorkShift
	if err := r.db.First(&ws, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &ws, nil
}

func (r *hrRepository) CreateWorkShift(ws *model.MstHrWorkShift) error {
	return r.db.Create(ws).Error
}

func (r *hrRepository) UpdateWorkShift(ws *model.MstHrWorkShift) error {
	return r.db.Model(&model.MstHrWorkShift{}).Where("id = ?", ws.ID).Updates(ws).Error
}

func (r *hrRepository) DeleteWorkShift(id string) error {
	return r.db.Delete(&model.MstHrWorkShift{}, "id = ?", id).Error
}

// ── Transaksi Presensi ──────────────────────────────────────────────

func (r *hrRepository) FindAttendancesPaginated(page, perPage int, employeeID, date, status, search string) ([]model.TrxHrAttendance, int64, error) {
	var list []model.TrxHrAttendance
	var total int64
	q := r.db.Model(&model.TrxHrAttendance{})
	if employeeID != "" {
		q = q.Where("employee_id = ?", employeeID)
	}
	if date != "" {
		q = q.Where("attendance_date = ?", date)
	}
	if status != "" {
		q = q.Where("attendance_status = ?", status)
	}
	if search != "" {
		q = q.Joins("JOIN mst_hr_employees ON mst_hr_employees.id = trx_hr_attendances.employee_id").
			Where("mst_hr_employees.full_name LIKE ? OR mst_hr_employees.employee_number LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("attendance_date DESC, created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		var emp model.MstHrEmployee
		if err := r.db.Select("full_name, employee_number, department_id").First(&emp, "id = ?", list[i].EmployeeID).Error; err == nil {
			list[i].EmployeeName = emp.FullName
			list[i].EmployeeNumber = emp.EmployeeNumber
			if emp.DepartmentID != nil {
				var dept model.MstHrDepartment
				if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
					list[i].DepartmentName = dept.DepartmentName
				}
			}
		}
	}
	return list, total, nil
}

func (r *hrRepository) FindAttendanceByID(id string) (*model.TrxHrAttendance, error) {
	var att model.TrxHrAttendance
	if err := r.db.First(&att, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var emp model.MstHrEmployee
	if err := r.db.Select("full_name, employee_number").First(&emp, "id = ?", att.EmployeeID).Error; err == nil {
		att.EmployeeName = emp.FullName
		att.EmployeeNumber = emp.EmployeeNumber
	}
	return &att, nil
}

func (r *hrRepository) FindAttendanceByEmployeeAndDate(employeeID, date string) (*model.TrxHrAttendance, error) {
	var att model.TrxHrAttendance
	if err := r.db.Where("employee_id = ? AND attendance_date = ?", employeeID, date).First(&att).Error; err != nil {
		return nil, err
	}
	var emp model.MstHrEmployee
	if err := r.db.Select("full_name, employee_number").First(&emp, "id = ?", att.EmployeeID).Error; err == nil {
		att.EmployeeName = emp.FullName
		att.EmployeeNumber = emp.EmployeeNumber
	}
	return &att, nil
}

func (r *hrRepository) CreateAttendance(att *model.TrxHrAttendance) error {
	return r.db.Create(att).Error
}

func (r *hrRepository) UpdateAttendance(att *model.TrxHrAttendance) error {
	return r.db.Model(&model.TrxHrAttendance{}).Where("id = ?", att.ID).Updates(att).Error
}

func (r *hrRepository) DeleteAttendance(id string) error {
	return r.db.Delete(&model.TrxHrAttendance{}, "id = ?", id).Error
}

// ── Transaksi Pengajuan Cuti ────────────────────────────────────────

func (r *hrRepository) FindLeaveRequestsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrLeaveRequest, int64, error) {
	var list []model.TrxHrLeaveRequest
	var total int64
	q := r.db.Model(&model.TrxHrLeaveRequest{})
	if employeeID != "" {
		q = q.Where("employee_id = ?", employeeID)
	}
	if status != "" {
		q = q.Where("approval_status = ? OR status = ?", status, status)
	}
	if search != "" {
		q = q.Where("request_number LIKE ? OR reason LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		var emp model.MstHrEmployee
		if err := r.db.Select("full_name, employee_number, department_id").First(&emp, "id = ?", list[i].EmployeeID).Error; err == nil {
			list[i].EmployeeName = emp.FullName
			list[i].EmployeeNumber = emp.EmployeeNumber
			if emp.DepartmentID != nil {
				var dept model.MstHrDepartment
				if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
					list[i].DepartmentName = dept.DepartmentName
				}
			}
		}
		var lt model.MstHrLeaveType
		if err := r.db.Select("leave_name").First(&lt, "id = ?", list[i].LeaveTypeID).Error; err == nil {
			list[i].LeaveName = lt.LeaveName
		}
	}
	return list, total, nil
}

func (r *hrRepository) FindLeaveRequestByID(id string) (*model.TrxHrLeaveRequest, error) {
	var lr model.TrxHrLeaveRequest
	if err := r.db.First(&lr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var emp model.MstHrEmployee
	if err := r.db.Select("full_name, employee_number").First(&emp, "id = ?", lr.EmployeeID).Error; err == nil {
		lr.EmployeeName = emp.FullName
		lr.EmployeeNumber = emp.EmployeeNumber
	}
	var lt model.MstHrLeaveType
	if err := r.db.Select("leave_name").First(&lt, "id = ?", lr.LeaveTypeID).Error; err == nil {
		lr.LeaveName = lt.LeaveName
	}
	return &lr, nil
}

func (r *hrRepository) CreateLeaveRequest(lr *model.TrxHrLeaveRequest) error {
	return r.db.Create(lr).Error
}

func (r *hrRepository) UpdateLeaveRequest(lr *model.TrxHrLeaveRequest) error {
	return r.db.Model(&model.TrxHrLeaveRequest{}).Where("id = ?", lr.ID).Updates(lr).Error
}

func (r *hrRepository) UpdateLeaveRequestStatus(id, status, rejectionReason, approvedBy string) error {
	updates := map[string]interface{}{
		"approval_status": status,
		"status":          status,
	}
	if rejectionReason != "" {
		updates["rejection_reason"] = rejectionReason
	}
	if approvedBy != "" {
		updates["approved_by"] = approvedBy
	}
	return r.db.Model(&model.TrxHrLeaveRequest{}).Where("id = ?", id).Updates(updates).Error
}

func (r *hrRepository) DeleteLeaveRequest(id string) error {
	return r.db.Delete(&model.TrxHrLeaveRequest{}, "id = ?", id).Error
}

// ── Transaksi Payroll ───────────────────────────────────────────────

func (r *hrRepository) FindPayrollsPaginated(page, perPage int, period, employeeID, status, search string) ([]model.TrxHrPayroll, int64, error) {
	var list []model.TrxHrPayroll
	var total int64
	q := r.db.Model(&model.TrxHrPayroll{})
	if period != "" {
		q = q.Where("payroll_period = ?", period)
	}
	if employeeID != "" {
		q = q.Where("employee_id = ?", employeeID)
	}
	if status != "" {
		q = q.Where("payment_status = ?", status)
	}
	if search != "" {
		q = q.Where("payroll_number LIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Preload("Items").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		var emp model.MstHrEmployee
		if err := r.db.Select("full_name, employee_number, department_id").First(&emp, "id = ?", list[i].EmployeeID).Error; err == nil {
			list[i].EmployeeName = emp.FullName
			list[i].EmployeeNumber = emp.EmployeeNumber
			if emp.DepartmentID != nil {
				var dept model.MstHrDepartment
				if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
					list[i].DepartmentName = dept.DepartmentName
				}
			}
		}
		for j := range list[i].Items {
			var sc model.MstHrSalaryComponent
			if err := r.db.Select("component_name").First(&sc, "id = ?", list[i].Items[j].ComponentID).Error; err == nil {
				list[i].Items[j].ComponentName = sc.ComponentName
			}
		}
	}
	return list, total, nil
}

func (r *hrRepository) FindPayrollByID(id string) (*model.TrxHrPayroll, error) {
	var p model.TrxHrPayroll
	if err := r.db.Preload("Items").First(&p, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var emp model.MstHrEmployee
	if err := r.db.Select("full_name, employee_number, department_id").First(&emp, "id = ?", p.EmployeeID).Error; err == nil {
		p.EmployeeName = emp.FullName
		p.EmployeeNumber = emp.EmployeeNumber
		if emp.DepartmentID != nil {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
				p.DepartmentName = dept.DepartmentName
			}
		}
	}
	for j := range p.Items {
		var sc model.MstHrSalaryComponent
		if err := r.db.Select("component_name").First(&sc, "id = ?", p.Items[j].ComponentID).Error; err == nil {
			p.Items[j].ComponentName = sc.ComponentName
		}
	}
	return &p, nil
}

func (r *hrRepository) CreatePayroll(p *model.TrxHrPayroll) error {
	return r.db.Create(p).Error
}

func (r *hrRepository) UpdatePayroll(p *model.TrxHrPayroll) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TrxHrPayrollItem{}, "payroll_id = ?", p.ID).Error; err != nil {
			return err
		}
		if len(p.Items) > 0 {
			if err := tx.Create(&p.Items).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.TrxHrPayroll{}).Where("id = ?", p.ID).Omit("Items").Updates(p).Error
	})
}

func (r *hrRepository) UpdatePayrollStatus(id, status string) error {
	return r.db.Model(&model.TrxHrPayroll{}).Where("id = ?", id).Update("payment_status", status).Error
}

func (r *hrRepository) DeletePayroll(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Delete(&model.TrxHrPayrollItem{}, "payroll_id = ?", id).Error
		return tx.Delete(&model.TrxHrPayroll{}, "id = ?", id).Error
	})
}

// ── Transaksi Reimbursement ─────────────────────────────────────────

func (r *hrRepository) FindReimbursementsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrReimbursement, int64, error) {
	var list []model.TrxHrReimbursement
	var total int64
	q := r.db.Model(&model.TrxHrReimbursement{})
	if employeeID != "" {
		q = q.Where("employee_id = ?", employeeID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("claim_number LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("claim_date DESC, created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		var emp model.MstHrEmployee
		if err := r.db.Select("full_name, employee_number, department_id").First(&emp, "id = ?", list[i].EmployeeID).Error; err == nil {
			list[i].EmployeeName = emp.FullName
			list[i].EmployeeNumber = emp.EmployeeNumber
			if emp.DepartmentID != nil {
				var dept model.MstHrDepartment
				if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
					list[i].DepartmentName = dept.DepartmentName
				}
			}
		}
		if list[i].CategoryID != "" {
			var cat model.MstFinanceCategory
			if err := r.db.Select("category_name").First(&cat, "id = ?", list[i].CategoryID).Error; err == nil {
				list[i].CategoryName = cat.CategoryName
			}
		}
	}
	return list, total, nil
}

func (r *hrRepository) FindReimbursementByID(id string) (*model.TrxHrReimbursement, error) {
	var rm model.TrxHrReimbursement
	if err := r.db.First(&rm, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var emp model.MstHrEmployee
	if err := r.db.Select("full_name, employee_number").First(&emp, "id = ?", rm.EmployeeID).Error; err == nil {
		rm.EmployeeName = emp.FullName
		rm.EmployeeNumber = emp.EmployeeNumber
	}
	if rm.CategoryID != "" {
		var cat model.MstFinanceCategory
		if err := r.db.Select("category_name").First(&cat, "id = ?", rm.CategoryID).Error; err == nil {
			rm.CategoryName = cat.CategoryName
		}
	}
	return &rm, nil
}

func (r *hrRepository) CreateReimbursement(rm *model.TrxHrReimbursement) error {
	return r.db.Create(rm).Error
}

func (r *hrRepository) UpdateReimbursement(rm *model.TrxHrReimbursement) error {
	return r.db.Model(&model.TrxHrReimbursement{}).Where("id = ?", rm.ID).Updates(rm).Error
}

func (r *hrRepository) UpdateReimbursementStatus(id, status string) error {
	return r.db.Model(&model.TrxHrReimbursement{}).Where("id = ?", id).Update("status", status).Error
}

func (r *hrRepository) DeleteReimbursement(id string) error {
	return r.db.Delete(&model.TrxHrReimbursement{}, "id = ?", id).Error
}

// ── Transaksi Penilaian Kinerja ─────────────────────────────────────

func (r *hrRepository) FindPerformanceReviewsPaginated(page, perPage int, employeeID, period, status, search string) ([]model.TrxHrPerformanceReview, int64, error) {
	var list []model.TrxHrPerformanceReview
	var total int64
	q := r.db.Model(&model.TrxHrPerformanceReview{})
	if employeeID != "" {
		q = q.Where("employee_id = ?", employeeID)
	}
	if period != "" {
		q = q.Where("review_period = ?", period)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("review_number LIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Preload("Items").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		var emp model.MstHrEmployee
		if err := r.db.Select("full_name, employee_number, department_id").First(&emp, "id = ?", list[i].EmployeeID).Error; err == nil {
			list[i].EmployeeName = emp.FullName
			list[i].EmployeeNumber = emp.EmployeeNumber
			if emp.DepartmentID != nil {
				var dept model.MstHrDepartment
				if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
					list[i].DepartmentName = dept.DepartmentName
				}
			}
		}
		var reviewer model.MstHrEmployee
		if err := r.db.Select("full_name").First(&reviewer, "id = ?", list[i].ReviewerID).Error; err == nil {
			list[i].ReviewerName = reviewer.FullName
		}
	}
	return list, total, nil
}

func (r *hrRepository) FindPerformanceReviewByID(id string) (*model.TrxHrPerformanceReview, error) {
	var pr model.TrxHrPerformanceReview
	if err := r.db.Preload("Items").First(&pr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	var emp model.MstHrEmployee
	if err := r.db.Select("full_name, employee_number").First(&emp, "id = ?", pr.EmployeeID).Error; err == nil {
		pr.EmployeeName = emp.FullName
		pr.EmployeeNumber = emp.EmployeeNumber
	}
	var reviewer model.MstHrEmployee
	if err := r.db.Select("full_name").First(&reviewer, "id = ?", pr.ReviewerID).Error; err == nil {
		pr.ReviewerName = reviewer.FullName
	}
	return &pr, nil
}

func (r *hrRepository) CreatePerformanceReview(pr *model.TrxHrPerformanceReview) error {
	return r.db.Create(pr).Error
}

func (r *hrRepository) UpdatePerformanceReview(pr *model.TrxHrPerformanceReview) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TrxHrPerformanceReviewItem{}, "review_id = ?", pr.ID).Error; err != nil {
			return err
		}
		if len(pr.Items) > 0 {
			if err := tx.Create(&pr.Items).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.TrxHrPerformanceReview{}).Where("id = ?", pr.ID).Omit("Items").Updates(pr).Error
	})
}

func (r *hrRepository) UpdatePerformanceReviewStatus(id, status string) error {
	return r.db.Model(&model.TrxHrPerformanceReview{}).Where("id = ?", id).Update("status", status).Error
}

func (r *hrRepository) DeletePerformanceReview(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Delete(&model.TrxHrPerformanceReviewItem{}, "review_id = ?", id).Error
		return tx.Delete(&model.TrxHrPerformanceReview{}, "id = ?", id).Error
	})
}

// ── Laporan HR ──────────────────────────────────────────────────────

func (r *hrRepository) GetRekapAbsensi(startDate, endDate, departmentID, employeeID string) (*model.RekapAbsensiReport, error) {
	report := &model.RekapAbsensiReport{
		Summary: model.RekapAbsensiSummary{},
		Rows:    []model.RekapAbsensiRow{},
	}

	empQuery := r.db.Model(&model.MstHrEmployee{}).Where("active_status = 'ACTIVE'")
	if departmentID != "" {
		empQuery = empQuery.Where("department_id = ?", departmentID)
	}
	if employeeID != "" {
		empQuery = empQuery.Where("id = ?", employeeID)
	}

	var employees []model.MstHrEmployee
	if err := empQuery.Order("full_name ASC").Find(&employees).Error; err != nil {
		return nil, err
	}

	for idx, emp := range employees {
		attQuery := r.db.Model(&model.TrxHrAttendance{}).Where("employee_id = ?", emp.ID)
		if startDate != "" {
			attQuery = attQuery.Where("attendance_date >= ?", startDate)
		}
		if endDate != "" {
			attQuery = attQuery.Where("attendance_date <= ?", endDate)
		}

		var hadir, terlambat, alpha, cuti, sakit int64
		_ = attQuery.Session(&gorm.Session{}).Where("attendance_status = 'PRESENT'").Count(&hadir)
		_ = attQuery.Session(&gorm.Session{}).Where("attendance_status = 'LATE'").Count(&terlambat)
		_ = attQuery.Session(&gorm.Session{}).Where("attendance_status = 'ABSENT'").Count(&alpha)
		_ = attQuery.Session(&gorm.Session{}).Where("attendance_status = 'LEAVE'").Count(&cuti)
		_ = attQuery.Session(&gorm.Session{}).Where("attendance_status = 'SICK'").Count(&sakit)

		row := model.RekapAbsensiRow{
			No:             idx + 1,
			EmployeeNumber: emp.EmployeeNumber,
			EmployeeName:   emp.FullName,
			TotalHadir:     int(hadir),
			TotalTerlambat: int(terlambat),
			TotalAlpha:     int(alpha),
			TotalCuti:      int(cuti),
			TotalSakit:     int(sakit),
		}
		report.Rows = append(report.Rows, row)

		report.Summary.TotalHadir += int(hadir)
		report.Summary.TotalTerlambat += int(terlambat)
		report.Summary.TotalAlpha += int(alpha)
		report.Summary.TotalCuti += int(cuti)
	}

	return report, nil
}

func (r *hrRepository) GetSlipGaji(period, employeeID, departmentID string) (*model.SlipGajiReport, error) {
	report := &model.SlipGajiReport{
		Summary: model.SlipGajiSummary{},
		Rows:    []model.SlipGajiRow{},
	}

	q := r.db.Model(&model.TrxHrPayroll{})
	if period != "" {
		q = q.Where("payroll_period = ?", period)
	}
	if employeeID != "" {
		q = q.Where("employee_id = ?", employeeID)
	}

	var payrolls []model.TrxHrPayroll
	if err := q.Order("created_at DESC").Find(&payrolls).Error; err != nil {
		return nil, err
	}

	for _, p := range payrolls {
		var emp model.MstHrEmployee
		if err := r.db.First(&emp, "id = ?", p.EmployeeID).Error; err == nil {
			if departmentID != "" && (emp.DepartmentID == nil || *emp.DepartmentID != departmentID) {
				continue
			}
			deptName := "-"
			if emp.DepartmentID != nil {
				var dept model.MstHrDepartment
				if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
					deptName = dept.DepartmentName
				}
			}

			row := model.SlipGajiRow{
				PayrollNumber:   p.PayrollNumber,
				EmployeeNumber:  emp.EmployeeNumber,
				EmployeeName:    emp.FullName,
				DepartmentName:  deptName,
				PayrollPeriod:   p.PayrollPeriod,
				TotalEarnings:   p.TotalEarnings,
				TotalDeductions: p.TotalDeductions,
				NetSalary:       p.NetSalary,
				PaymentStatus:   p.PaymentStatus,
			}
			report.Rows = append(report.Rows, row)

			report.Summary.TotalPayroll++
			report.Summary.TotalPendapatan += p.TotalEarnings
			report.Summary.TotalPotongan += p.TotalDeductions
			report.Summary.TotalGajiBersih += p.NetSalary
		}
	}

	return report, nil
}

func (r *hrRepository) GetSaldoCuti(departmentID, employeeID string) (*model.SaldoCutiReport, error) {
	report := &model.SaldoCutiReport{
		Summary: model.SaldoCutiSummary{},
		Rows:    []model.SaldoCutiRow{},
	}

	empQuery := r.db.Model(&model.MstHrEmployee{}).Where("active_status = 'ACTIVE'")
	if departmentID != "" {
		empQuery = empQuery.Where("department_id = ?", departmentID)
	}
	if employeeID != "" {
		empQuery = empQuery.Where("id = ?", employeeID)
	}

	var employees []model.MstHrEmployee
	if err := empQuery.Order("full_name ASC").Find(&employees).Error; err != nil {
		return nil, err
	}

	var defaultQuota int = 12
	var annualLt model.MstHrLeaveType
	if err := r.db.Where("code = 'CUTI_TAHUNAN' OR leave_name LIKE '%Tahunan%'").First(&annualLt).Error; err == nil {
		defaultQuota = annualLt.AnnualQuotaDays
	}

	var totalRemaining int = 0
	maxRemaining := -1
	minRemaining := 999

	for _, emp := range employees {
		deptName := "-"
		if emp.DepartmentID != nil {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", *emp.DepartmentID).Error; err == nil {
				deptName = dept.DepartmentName
			}
		}

		var usedDays int64
		_ = r.db.Model(&model.TrxHrLeaveRequest{}).
			Where("employee_id = ? AND (approval_status = 'APPROVED' OR status = 'APPROVED')", emp.ID).
			Select("COALESCE(SUM(total_days), 0)").
			Scan(&usedDays)

		remaining := defaultQuota - int(usedDays)
		if remaining < 0 {
			remaining = 0
		}

		row := model.SaldoCutiRow{
			EmployeeNumber: emp.EmployeeNumber,
			EmployeeName:   emp.FullName,
			DepartmentName: deptName,
			TotalQuota:     defaultQuota,
			UsedQuota:      int(usedDays),
			RemainingQuota: remaining,
		}
		report.Rows = append(report.Rows, row)

		totalRemaining += remaining
		if remaining > maxRemaining {
			maxRemaining = remaining
		}
		if remaining < minRemaining {
			minRemaining = remaining
		}
	}

	report.Summary.TotalKaryawan = len(employees)
	if len(employees) > 0 {
		report.Summary.RataRataSisaCuti = totalRemaining / len(employees)
		report.Summary.CutiTerbanyak = maxRemaining
		report.Summary.CutiTerendah = minRemaining
	}

	return report, nil
}

func (r *hrRepository) GetTurnover(year int) (*model.TurnoverReport, error) {
	if year <= 0 {
		year = time.Now().Year()
	}

	report := &model.TurnoverReport{
		Summary: model.TurnoverSummary{},
		Rows:    []model.TurnoverRow{},
	}

	months := []string{"Jan", "Feb", "Mar", "Apr", "Mei", "Jun", "Jul", "Ags", "Sep", "Okt", "Nov", "Des"}
	var runningTotal int64
	_ = r.db.Model(&model.MstHrEmployee{}).
		Where("join_date < ?", fmt.Sprintf("%04d-01-01", year)).
		Count(&runningTotal)

	initialCount := runningTotal
	totalIn := 0
	totalOut := 0

	for mIdx, monthName := range months {
		mStr := fmt.Sprintf("%04d-%02d", year, mIdx+1)
		var masuk, keluar int64
		_ = r.db.Model(&model.MstHrEmployee{}).Where("join_date LIKE ?", mStr+"%").Count(&masuk)
		_ = r.db.Model(&model.MstHrEmployee{}).Where("active_status = 'RESIGNED' AND updated_at LIKE ?", mStr+"%").Count(&keluar)

		startCount := runningTotal
		endCount := startCount + masuk - keluar
		if endCount < 0 {
			endCount = 0
		}
		runningTotal = endCount

		var rate float64 = 0
		avg := float64(startCount+endCount) / 2.0
		if avg > 0 {
			rate = (float64(keluar) / avg) * 100
		}

		report.Rows = append(report.Rows, model.TurnoverRow{
			Period:     fmt.Sprintf("%s %d", monthName, year),
			TotalAwal:  int(startCount),
			Masuk:      int(masuk),
			Keluar:     int(keluar),
			TotalAkhir: int(endCount),
			Rate:       rate,
		})

		totalIn += int(masuk)
		totalOut += int(keluar)
	}

	report.Summary.TotalKaryawanAwal = int(initialCount)
	report.Summary.KaryawanMasuk = totalIn
	report.Summary.KaryawanKeluar = totalOut
	avgYear := float64(int(initialCount)+int(runningTotal)) / 2.0
	if avgYear > 0 {
		report.Summary.TurnoverRate = (float64(totalOut) / avgYear) * 100
	}

	return report, nil
}

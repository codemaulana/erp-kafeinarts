package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"

	"github.com/google/uuid"
)

// ── DTOs Master HR ──────────────────────────────────────────────────

type DepartmentDTO struct {
	Code           string  `json:"code"`
	DepartmentName string  `json:"department_name"`
	Name           string  `json:"name"`
	HeadEmployeeID *string `json:"head_employee_id"`
	Description    string  `json:"description"`
}

type PositionDTO struct {
	DepartmentID   string `json:"department_id"`
	Code           string `json:"code"`
	PositionName   string `json:"position_name"`
	Name           string `json:"name"`
	Level          string `json:"level"`
	JobDescription string `json:"job_description"`
	Description    string `json:"description"`
}

type EmployeeDTO struct {
	EmployeeNumber    string  `json:"employee_number"`
	EmployeeCode      string  `json:"employee_code"`
	FullName          string  `json:"full_name"`
	Name              string  `json:"name"`
	NationalIDNumber  string  `json:"national_id_number"`
	BirthDate         string  `json:"birth_date"`
	Gender            string  `json:"gender"`
	Address           string  `json:"address"`
	Email             string  `json:"email"`
	PhoneNumber       string  `json:"phone_number"`
	Phone             string  `json:"phone"`
	DepartmentID      *string `json:"department_id"`
	PositionID        *string `json:"position_id"`
	JoinDate          string  `json:"join_date"`
	EmploymentStatus  string  `json:"employment_status"`
	EmploymentType    string  `json:"employment_type"`
	ActiveStatus      string  `json:"active_status"`
	Status            string  `json:"status"`
	Salary            float64 `json:"salary"`
	BankName          string  `json:"bank_name"`
	BankAccountNumber string  `json:"bank_account_number"`
	BankAccount       string  `json:"bank_account"`
	PhotoURL          string  `json:"photo_url"`
}

type LeaveTypeDTO struct {
	Code               string `json:"code"`
	LeaveName          string `json:"leave_name"`
	Name               string `json:"name"`
	AnnualQuotaDays    int    `json:"annual_quota_days"`
	QuotaDays          int    `json:"quota_days"`
	IsPaid             bool   `json:"is_paid"`
	RequiresAttachment bool   `json:"requires_attachment"`
	Description        string `json:"description"`
}

type SalaryComponentDTO struct {
	Code              string `json:"code"`
	ComponentName     string `json:"component_name"`
	Name              string `json:"name"`
	ComponentType     string `json:"component_type"`
	Type              string `json:"type"`
	CalculationMethod string `json:"calculation_method"`
	IsTaxable         bool   `json:"is_taxable"`
	Description       string `json:"description"`
}

type WorkShiftDTO struct {
	ShiftName            string            `json:"shift_name"`
	Name                 string            `json:"name"`
	TimeIn               string            `json:"time_in"`
	StartTime            string            `json:"start_time"`
	TimeOut              string            `json:"time_out"`
	EndTime              string            `json:"end_time"`
	LateToleranceMinutes int               `json:"late_tolerance_minutes"`
	GracePeriodMinutes   int               `json:"grace_period_minutes"`
	WorkDays             model.StringArray `json:"work_days"`
	IsNightShift         bool              `json:"is_night_shift"`
}

// ── DTOs Transaksi HR ───────────────────────────────────────────────

type AttendanceDTO struct {
	EmployeeID       string  `json:"employee_id"`
	AttendanceDate   string  `json:"attendance_date"`
	Date             string  `json:"date"`
	CheckInTime      string  `json:"check_in_time"`
	CheckIn          string  `json:"check_in"`
	TimeIn           string  `json:"time_in"`
	CheckOutTime     string  `json:"check_out_time"`
	CheckOut         string  `json:"check_out"`
	TimeOut          string  `json:"time_out"`
	CheckInLat       float64 `json:"check_in_lat"`
	Latitude         float64 `json:"latitude"`
	CheckInLng       float64 `json:"check_in_lng"`
	Longitude        float64 `json:"longitude"`
	AttendanceStatus string  `json:"attendance_status"`
	Status           string  `json:"status"`
	Notes            string  `json:"notes"`
	SelfiePhotoURL   string  `json:"selfie_photo_url"`
}

type LeaveRequestDTO struct {
	RequestNumber   string `json:"request_number"`
	EmployeeID      string `json:"employee_id"`
	LeaveTypeID     string `json:"leave_type_id"`
	StartDate       string `json:"start_date"`
	EndDate         string `json:"end_date"`
	TotalDays       int    `json:"total_days"`
	Reason          string `json:"reason"`
	AttachmentURL   string `json:"attachment_url"`
	ApprovalStatus  string `json:"approval_status"`
	Status          string `json:"status"`
	RejectionReason string `json:"rejection_reason"`
	ApprovedBy      string `json:"approved_by"`
}

type PayrollItemDTO struct {
	ComponentID   string  `json:"component_id"`
	ComponentName string  `json:"component_name"`
	ItemType      string  `json:"item_type"`
	ComponentType string  `json:"component_type"`
	Amount        float64 `json:"amount"`
}

type PayrollDTO struct {
	PayrollNumber   string           `json:"payroll_number"`
	PayrollPeriod   string           `json:"payroll_period"`
	PeriodMonth     int              `json:"period_month"`
	PeriodYear      int              `json:"period_year"`
	EmployeeID      string           `json:"employee_id"`
	BasicSalary     float64          `json:"basic_salary"`
	Allowances      float64          `json:"allowances"`
	OvertimePay     float64          `json:"overtime_pay"`
	Bonus           float64          `json:"bonus"`
	Deductions      float64          `json:"deductions"`
	TotalEarnings   float64          `json:"total_earnings"`
	TotalDeductions float64          `json:"total_deductions"`
	NetSalary       float64          `json:"net_salary"`
	PaymentStatus   string           `json:"payment_status"`
	Status          string           `json:"status"`
	PaymentDate     string           `json:"payment_date"`
	Notes           string           `json:"notes"`
	Items           []PayrollItemDTO `json:"items"`
}

type ReimbursementDTO struct {
	ClaimNumber   string  `json:"claim_number"`
	EmployeeID    string  `json:"employee_id"`
	ClaimDate     string  `json:"claim_date"`
	CategoryID    string  `json:"category_id"`
	Category      string  `json:"category"`
	Amount        float64 `json:"amount"`
	Description   string  `json:"description"`
	AttachmentURL string  `json:"attachment_url"`
	Status        string  `json:"status"`
}

type PerformanceReviewItemDTO struct {
	Criterion    string  `json:"criterion"`
	CriteriaName string  `json:"criteria_name"`
	Score        float64 `json:"score"`
	Weight       float64 `json:"weight"`
	Comments     string  `json:"comments"`
	Notes        string  `json:"notes"`
}

type PerformanceReviewDTO struct {
	ReviewNumber    string                     `json:"review_number"`
	EmployeeID      string                     `json:"employee_id"`
	ReviewerID      string                     `json:"reviewer_id"`
	ReviewPeriod    string                     `json:"review_period"`
	Period          string                     `json:"period"`
	ReviewDate      string                     `json:"review_date"`
	TotalScore      float64                    `json:"total_score"`
	OverallFeedback string                     `json:"overall_feedback"`
	Feedback        string                     `json:"feedback"`
	Status          string                     `json:"status"`
	Items           []PerformanceReviewItemDTO `json:"items"`
}

// ── HrService Interface ─────────────────────────────────────────────

type HrService interface {
	// Master Departemen
	GetAllDepartments(search string) ([]model.MstHrDepartment, error)
	GetDepartmentsPaginated(page, perPage int, search string) ([]model.MstHrDepartment, int64, error)
	GetDepartmentByID(id string) (*model.MstHrDepartment, error)
	CreateDepartment(dto *DepartmentDTO) (*model.MstHrDepartment, error)
	UpdateDepartment(id string, dto *DepartmentDTO) (*model.MstHrDepartment, error)
	DeleteDepartment(id string) error

	// Master Jabatan
	GetAllPositions(departmentID, search string) ([]model.MstHrPosition, error)
	GetPositionsPaginated(page, perPage int, departmentID, search string) ([]model.MstHrPosition, int64, error)
	GetPositionByID(id string) (*model.MstHrPosition, error)
	CreatePosition(dto *PositionDTO) (*model.MstHrPosition, error)
	UpdatePosition(id string, dto *PositionDTO) (*model.MstHrPosition, error)
	DeletePosition(id string) error

	// Master Karyawan
	GetAllEmployees(departmentID, status, search string) ([]model.MstHrEmployee, error)
	GetEmployeesPaginated(page, perPage int, departmentID, status, search string) ([]model.MstHrEmployee, int64, error)
	GetEmployeeByID(id string) (*model.MstHrEmployee, error)
	CreateEmployee(dto *EmployeeDTO) (*model.MstHrEmployee, error)
	UpdateEmployee(id string, dto *EmployeeDTO) (*model.MstHrEmployee, error)
	DeleteEmployee(id string) error

	// Master Jenis Cuti
	GetAllLeaveTypes(search string) ([]model.MstHrLeaveType, error)
	GetLeaveTypesPaginated(page, perPage int, search string) ([]model.MstHrLeaveType, int64, error)
	GetLeaveTypeByID(id string) (*model.MstHrLeaveType, error)
	CreateLeaveType(dto *LeaveTypeDTO) (*model.MstHrLeaveType, error)
	UpdateLeaveType(id string, dto *LeaveTypeDTO) (*model.MstHrLeaveType, error)
	DeleteLeaveType(id string) error

	// Master Komponen Gaji
	GetAllSalaryComponents(compType, search string) ([]model.MstHrSalaryComponent, error)
	GetSalaryComponentsPaginated(page, perPage int, compType, search string) ([]model.MstHrSalaryComponent, int64, error)
	GetSalaryComponentByID(id string) (*model.MstHrSalaryComponent, error)
	CreateSalaryComponent(dto *SalaryComponentDTO) (*model.MstHrSalaryComponent, error)
	UpdateSalaryComponent(id string, dto *SalaryComponentDTO) (*model.MstHrSalaryComponent, error)
	DeleteSalaryComponent(id string) error

	// Master Jadwal Kerja
	GetAllWorkShifts(search string) ([]model.MstHrWorkShift, error)
	GetWorkShiftsPaginated(page, perPage int, search string) ([]model.MstHrWorkShift, int64, error)
	GetWorkShiftByID(id string) (*model.MstHrWorkShift, error)
	CreateWorkShift(dto *WorkShiftDTO) (*model.MstHrWorkShift, error)
	UpdateWorkShift(id string, dto *WorkShiftDTO) (*model.MstHrWorkShift, error)
	DeleteWorkShift(id string) error

	// Transaksi Presensi
	GetAttendancesPaginated(page, perPage int, employeeID, date, status, search string) ([]model.TrxHrAttendance, int64, error)
	GetAttendanceByID(id string) (*model.TrxHrAttendance, error)
	CreateAttendance(dto *AttendanceDTO) (*model.TrxHrAttendance, error)
	UpdateAttendance(id string, dto *AttendanceDTO) (*model.TrxHrAttendance, error)
	DeleteAttendance(id string) error

	// Transaksi Pengajuan Cuti
	GetLeaveRequestsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrLeaveRequest, int64, error)
	GetLeaveRequestByID(id string) (*model.TrxHrLeaveRequest, error)
	CreateLeaveRequest(dto *LeaveRequestDTO) (*model.TrxHrLeaveRequest, error)
	UpdateLeaveRequest(id string, dto *LeaveRequestDTO) (*model.TrxHrLeaveRequest, error)
	UpdateLeaveRequestStatus(id, status, reason, approvedBy string) error
	DeleteLeaveRequest(id string) error

	// Transaksi Payroll
	GetPayrollsPaginated(page, perPage int, period, employeeID, status, search string) ([]model.TrxHrPayroll, int64, error)
	GetPayrollByID(id string) (*model.TrxHrPayroll, error)
	CreatePayroll(dto *PayrollDTO) (*model.TrxHrPayroll, error)
	UpdatePayroll(id string, dto *PayrollDTO) (*model.TrxHrPayroll, error)
	UpdatePayrollStatus(id, status string) error
	DeletePayroll(id string) error

	// Transaksi Reimbursement
	GetReimbursementsPaginated(page, perPage int, employeeID, status, search string) ([]model.TrxHrReimbursement, int64, error)
	GetReimbursementByID(id string) (*model.TrxHrReimbursement, error)
	CreateReimbursement(dto *ReimbursementDTO) (*model.TrxHrReimbursement, error)
	UpdateReimbursement(id string, dto *ReimbursementDTO) (*model.TrxHrReimbursement, error)
	UpdateReimbursementStatus(id, status string) error
	DeleteReimbursement(id string) error

	// Transaksi Penilaian Kinerja
	GetPerformanceReviewsPaginated(page, perPage int, employeeID, period, status, search string) ([]model.TrxHrPerformanceReview, int64, error)
	GetPerformanceReviewByID(id string) (*model.TrxHrPerformanceReview, error)
	CreatePerformanceReview(dto *PerformanceReviewDTO) (*model.TrxHrPerformanceReview, error)
	UpdatePerformanceReview(id string, dto *PerformanceReviewDTO) (*model.TrxHrPerformanceReview, error)
	UpdatePerformanceReviewStatus(id, status string) error
	DeletePerformanceReview(id string) error

	// Laporan HR
	GetRekapAbsensi(startDate, endDate, departmentID, employeeID string) (*model.RekapAbsensiReport, error)
	GetSlipGaji(period, employeeID, departmentID string) (*model.SlipGajiReport, error)
	GetSaldoCuti(departmentID, employeeID string) (*model.SaldoCutiReport, error)
	GetTurnover(year int) (*model.TurnoverReport, error)
}

type hrService struct {
	repo repository.HrRepository
}

func NewHrService(repo repository.HrRepository) HrService {
	return &hrService{repo: repo}
}

// ── Master Departemen ───────────────────────────────────────────────

func (s *hrService) GetAllDepartments(search string) ([]model.MstHrDepartment, error) {
	return s.repo.FindDepartments(search)
}

func (s *hrService) GetDepartmentsPaginated(page, perPage int, search string) ([]model.MstHrDepartment, int64, error) {
	return s.repo.FindDepartmentsPaginated(page, perPage, search)
}

func (s *hrService) GetDepartmentByID(id string) (*model.MstHrDepartment, error) {
	return s.repo.FindDepartmentByID(id)
}

func (s *hrService) CreateDepartment(dto *DepartmentDTO) (*model.MstHrDepartment, error) {
	deptName := strings.TrimSpace(dto.DepartmentName)
	if deptName == "" {
		deptName = strings.TrimSpace(dto.Name)
	}
	if deptName == "" {
		return nil, errors.New("nama departemen wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("DEP-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	dept := &model.MstHrDepartment{
		ID:             uuid.NewString(),
		Code:           code,
		DepartmentName: deptName,
		HeadEmployeeID: dto.HeadEmployeeID,
		Description:    dto.Description,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreateDepartment(dept); err != nil {
		return nil, err
	}
	return s.repo.FindDepartmentByID(dept.ID)
}

func (s *hrService) UpdateDepartment(id string, dto *DepartmentDTO) (*model.MstHrDepartment, error) {
	existing, err := s.repo.FindDepartmentByID(id)
	if err != nil {
		return nil, errors.New("departemen tidak ditemukan")
	}

	deptName := strings.TrimSpace(dto.DepartmentName)
	if deptName == "" {
		deptName = strings.TrimSpace(dto.Name)
	}
	if deptName != "" {
		existing.DepartmentName = deptName
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	existing.HeadEmployeeID = dto.HeadEmployeeID
	existing.Description = dto.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateDepartment(existing); err != nil {
		return nil, err
	}
	return s.repo.FindDepartmentByID(id)
}

func (s *hrService) DeleteDepartment(id string) error {
	return s.repo.DeleteDepartment(id)
}

// ── Master Jabatan ──────────────────────────────────────────────────

func (s *hrService) GetAllPositions(departmentID, search string) ([]model.MstHrPosition, error) {
	return s.repo.FindPositions(departmentID, search)
}

func (s *hrService) GetPositionsPaginated(page, perPage int, departmentID, search string) ([]model.MstHrPosition, int64, error) {
	return s.repo.FindPositionsPaginated(page, perPage, departmentID, search)
}

func (s *hrService) GetPositionByID(id string) (*model.MstHrPosition, error) {
	return s.repo.FindPositionByID(id)
}

func (s *hrService) CreatePosition(dto *PositionDTO) (*model.MstHrPosition, error) {
	posName := strings.TrimSpace(dto.PositionName)
	if posName == "" {
		posName = strings.TrimSpace(dto.Name)
	}
	if posName == "" {
		return nil, errors.New("nama jabatan wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("POS-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	level := dto.Level
	if level == "" {
		level = "STAFF"
	}

	desc := dto.JobDescription
	if desc == "" {
		desc = dto.Description
	}

	pos := &model.MstHrPosition{
		ID:             uuid.NewString(),
		DepartmentID:   dto.DepartmentID,
		Code:           code,
		PositionName:   posName,
		Level:          level,
		JobDescription: desc,
		Description:    desc,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreatePosition(pos); err != nil {
		return nil, err
	}
	return s.repo.FindPositionByID(pos.ID)
}

func (s *hrService) UpdatePosition(id string, dto *PositionDTO) (*model.MstHrPosition, error) {
	existing, err := s.repo.FindPositionByID(id)
	if err != nil {
		return nil, errors.New("jabatan tidak ditemukan")
	}

	posName := strings.TrimSpace(dto.PositionName)
	if posName == "" {
		posName = strings.TrimSpace(dto.Name)
	}
	if posName != "" {
		existing.PositionName = posName
	}
	if dto.DepartmentID != "" {
		existing.DepartmentID = dto.DepartmentID
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	if dto.Level != "" {
		existing.Level = dto.Level
	}
	desc := dto.JobDescription
	if desc == "" {
		desc = dto.Description
	}
	existing.JobDescription = desc
	existing.Description = desc
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdatePosition(existing); err != nil {
		return nil, err
	}
	return s.repo.FindPositionByID(id)
}

func (s *hrService) DeletePosition(id string) error {
	return s.repo.DeletePosition(id)
}

// ── Master Karyawan ─────────────────────────────────────────────────

func (s *hrService) GetAllEmployees(departmentID, status, search string) ([]model.MstHrEmployee, error) {
	return s.repo.FindEmployees(departmentID, status, search)
}

func (s *hrService) GetEmployeesPaginated(page, perPage int, departmentID, status, search string) ([]model.MstHrEmployee, int64, error) {
	return s.repo.FindEmployeesPaginated(page, perPage, departmentID, status, search)
}

func (s *hrService) GetEmployeeByID(id string) (*model.MstHrEmployee, error) {
	return s.repo.FindEmployeeByID(id)
}

func (s *hrService) CreateEmployee(dto *EmployeeDTO) (*model.MstHrEmployee, error) {
	fullName := strings.TrimSpace(dto.FullName)
	if fullName == "" {
		fullName = strings.TrimSpace(dto.Name)
	}
	if fullName == "" {
		return nil, errors.New("nama karyawan wajib diisi")
	}

	empNum := strings.TrimSpace(dto.EmployeeNumber)
	if empNum == "" {
		empNum = strings.TrimSpace(dto.EmployeeCode)
	}
	if empNum == "" {
		empNum = fmt.Sprintf("EMP-%s-%s", time.Now().Format("0601"), strings.ToUpper(uuid.NewString()[:4]))
	}

	phone := strings.TrimSpace(dto.PhoneNumber)
	if phone == "" {
		phone = strings.TrimSpace(dto.Phone)
	}

	bankAcc := strings.TrimSpace(dto.BankAccountNumber)
	if bankAcc == "" {
		bankAcc = strings.TrimSpace(dto.BankAccount)
	}

	empStatus := dto.EmploymentStatus
	if empStatus == "" {
		empStatus = dto.EmploymentType
	}
	if empStatus == "" {
		empStatus = "PERMANENT"
	}

	actStatus := dto.ActiveStatus
	if actStatus == "" {
		actStatus = dto.Status
	}
	if actStatus == "" {
		actStatus = "ACTIVE"
	}

	joinDate := dto.JoinDate
	if joinDate == "" {
		joinDate = time.Now().Format("2006-01-02")
	}

	emp := &model.MstHrEmployee{
		ID:                uuid.NewString(),
		EmployeeNumber:    empNum,
		FullName:          fullName,
		NationalIDNumber:  dto.NationalIDNumber,
		BirthDate:         dto.BirthDate,
		Gender:            dto.Gender,
		Address:           dto.Address,
		Email:             dto.Email,
		PhoneNumber:       phone,
		DepartmentID:      dto.DepartmentID,
		PositionID:        dto.PositionID,
		JoinDate:          joinDate,
		EmploymentStatus:  empStatus,
		ActiveStatus:      actStatus,
		Salary:            dto.Salary,
		BankName:          dto.BankName,
		BankAccountNumber: bankAcc,
		PhotoURL:          dto.PhotoURL,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := s.repo.CreateEmployee(emp); err != nil {
		return nil, err
	}
	return s.repo.FindEmployeeByID(emp.ID)
}

func (s *hrService) UpdateEmployee(id string, dto *EmployeeDTO) (*model.MstHrEmployee, error) {
	existing, err := s.repo.FindEmployeeByID(id)
	if err != nil {
		return nil, errors.New("karyawan tidak ditemukan")
	}

	fullName := strings.TrimSpace(dto.FullName)
	if fullName == "" {
		fullName = strings.TrimSpace(dto.Name)
	}
	if fullName != "" {
		existing.FullName = fullName
	}
	if dto.NationalIDNumber != "" {
		existing.NationalIDNumber = dto.NationalIDNumber
	}
	if dto.BirthDate != "" {
		existing.BirthDate = dto.BirthDate
	}
	if dto.Gender != "" {
		existing.Gender = dto.Gender
	}
	if dto.Address != "" {
		existing.Address = dto.Address
	}
	if dto.Email != "" {
		existing.Email = dto.Email
	}
	phone := strings.TrimSpace(dto.PhoneNumber)
	if phone == "" {
		phone = strings.TrimSpace(dto.Phone)
	}
	if phone != "" {
		existing.PhoneNumber = phone
	}
	if dto.DepartmentID != nil {
		existing.DepartmentID = dto.DepartmentID
	}
	if dto.PositionID != nil {
		existing.PositionID = dto.PositionID
	}
	if dto.JoinDate != "" {
		existing.JoinDate = dto.JoinDate
	}
	empStatus := dto.EmploymentStatus
	if empStatus == "" {
		empStatus = dto.EmploymentType
	}
	if empStatus != "" {
		existing.EmploymentStatus = empStatus
	}
	actStatus := dto.ActiveStatus
	if actStatus == "" {
		actStatus = dto.Status
	}
	if actStatus != "" {
		existing.ActiveStatus = actStatus
	}
	if dto.Salary > 0 {
		existing.Salary = dto.Salary
	}
	if dto.BankName != "" {
		existing.BankName = dto.BankName
	}
	bankAcc := strings.TrimSpace(dto.BankAccountNumber)
	if bankAcc == "" {
		bankAcc = strings.TrimSpace(dto.BankAccount)
	}
	if bankAcc != "" {
		existing.BankAccountNumber = bankAcc
	}
	if dto.PhotoURL != "" {
		existing.PhotoURL = dto.PhotoURL
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateEmployee(existing); err != nil {
		return nil, err
	}
	return s.repo.FindEmployeeByID(id)
}

func (s *hrService) DeleteEmployee(id string) error {
	return s.repo.DeleteEmployee(id)
}

// ── Master Jenis Cuti ───────────────────────────────────────────────

func (s *hrService) GetAllLeaveTypes(search string) ([]model.MstHrLeaveType, error) {
	return s.repo.FindLeaveTypes(search)
}

func (s *hrService) GetLeaveTypesPaginated(page, perPage int, search string) ([]model.MstHrLeaveType, int64, error) {
	return s.repo.FindLeaveTypesPaginated(page, perPage, search)
}

func (s *hrService) GetLeaveTypeByID(id string) (*model.MstHrLeaveType, error) {
	return s.repo.FindLeaveTypeByID(id)
}

func (s *hrService) CreateLeaveType(dto *LeaveTypeDTO) (*model.MstHrLeaveType, error) {
	leaveName := strings.TrimSpace(dto.LeaveName)
	if leaveName == "" {
		leaveName = strings.TrimSpace(dto.Name)
	}
	if leaveName == "" {
		return nil, errors.New("nama jenis cuti wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("CUTI-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	quota := dto.AnnualQuotaDays
	if quota <= 0 {
		quota = dto.QuotaDays
	}
	if quota <= 0 {
		quota = 12
	}

	lt := &model.MstHrLeaveType{
		ID:                 uuid.NewString(),
		Code:               code,
		LeaveName:          leaveName,
		AnnualQuotaDays:    quota,
		IsPaid:             dto.IsPaid,
		RequiresAttachment: dto.RequiresAttachment,
		Description:        dto.Description,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.repo.CreateLeaveType(lt); err != nil {
		return nil, err
	}
	return lt, nil
}

func (s *hrService) UpdateLeaveType(id string, dto *LeaveTypeDTO) (*model.MstHrLeaveType, error) {
	existing, err := s.repo.FindLeaveTypeByID(id)
	if err != nil {
		return nil, errors.New("jenis cuti tidak ditemukan")
	}

	leaveName := strings.TrimSpace(dto.LeaveName)
	if leaveName == "" {
		leaveName = strings.TrimSpace(dto.Name)
	}
	if leaveName != "" {
		existing.LeaveName = leaveName
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	quota := dto.AnnualQuotaDays
	if quota <= 0 {
		quota = dto.QuotaDays
	}
	if quota > 0 {
		existing.AnnualQuotaDays = quota
	}
	existing.IsPaid = dto.IsPaid
	existing.RequiresAttachment = dto.RequiresAttachment
	existing.Description = dto.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateLeaveType(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *hrService) DeleteLeaveType(id string) error {
	return s.repo.DeleteLeaveType(id)
}

// ── Master Komponen Gaji ────────────────────────────────────────────

func (s *hrService) GetAllSalaryComponents(compType, search string) ([]model.MstHrSalaryComponent, error) {
	return s.repo.FindSalaryComponents(compType, search)
}

func (s *hrService) GetSalaryComponentsPaginated(page, perPage int, compType, search string) ([]model.MstHrSalaryComponent, int64, error) {
	return s.repo.FindSalaryComponentsPaginated(page, perPage, compType, search)
}

func (s *hrService) GetSalaryComponentByID(id string) (*model.MstHrSalaryComponent, error) {
	return s.repo.FindSalaryComponentByID(id)
}

func (s *hrService) CreateSalaryComponent(dto *SalaryComponentDTO) (*model.MstHrSalaryComponent, error) {
	compName := strings.TrimSpace(dto.ComponentName)
	if compName == "" {
		compName = strings.TrimSpace(dto.Name)
	}
	if compName == "" {
		return nil, errors.New("nama komponen gaji wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("SAL-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	cType := strings.ToUpper(dto.ComponentType)
	if cType == "" {
		cType = strings.ToUpper(dto.Type)
	}
	if cType == "" {
		cType = "EARNING"
	}

	method := strings.ToUpper(dto.CalculationMethod)
	if method == "" {
		method = "FIXED"
	}

	sc := &model.MstHrSalaryComponent{
		ID:                uuid.NewString(),
		Code:              code,
		ComponentName:     compName,
		ComponentType:     cType,
		CalculationMethod: method,
		IsTaxable:         dto.IsTaxable,
		Description:       dto.Description,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}
	if err := s.repo.CreateSalaryComponent(sc); err != nil {
		return nil, err
	}
	return sc, nil
}

func (s *hrService) UpdateSalaryComponent(id string, dto *SalaryComponentDTO) (*model.MstHrSalaryComponent, error) {
	existing, err := s.repo.FindSalaryComponentByID(id)
	if err != nil {
		return nil, errors.New("komponen gaji tidak ditemukan")
	}

	compName := strings.TrimSpace(dto.ComponentName)
	if compName == "" {
		compName = strings.TrimSpace(dto.Name)
	}
	if compName != "" {
		existing.ComponentName = compName
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	cType := strings.ToUpper(dto.ComponentType)
	if cType == "" {
		cType = strings.ToUpper(dto.Type)
	}
	if cType != "" {
		existing.ComponentType = cType
	}
	if dto.CalculationMethod != "" {
		existing.CalculationMethod = strings.ToUpper(dto.CalculationMethod)
	}
	existing.IsTaxable = dto.IsTaxable
	existing.Description = dto.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateSalaryComponent(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *hrService) DeleteSalaryComponent(id string) error {
	return s.repo.DeleteSalaryComponent(id)
}

// ── Master Jadwal Kerja ─────────────────────────────────────────────

func (s *hrService) GetAllWorkShifts(search string) ([]model.MstHrWorkShift, error) {
	return s.repo.FindWorkShifts(search)
}

func (s *hrService) GetWorkShiftsPaginated(page, perPage int, search string) ([]model.MstHrWorkShift, int64, error) {
	return s.repo.FindWorkShiftsPaginated(page, perPage, search)
}

func (s *hrService) GetWorkShiftByID(id string) (*model.MstHrWorkShift, error) {
	return s.repo.FindWorkShiftByID(id)
}

func (s *hrService) CreateWorkShift(dto *WorkShiftDTO) (*model.MstHrWorkShift, error) {
	shiftName := strings.TrimSpace(dto.ShiftName)
	if shiftName == "" {
		shiftName = strings.TrimSpace(dto.Name)
	}
	if shiftName == "" {
		return nil, errors.New("nama shift kerja wajib diisi")
	}

	timeIn := dto.TimeIn
	if timeIn == "" {
		timeIn = dto.StartTime
	}
	if timeIn == "" {
		timeIn = "08:00"
	}

	timeOut := dto.TimeOut
	if timeOut == "" {
		timeOut = dto.EndTime
	}
	if timeOut == "" {
		timeOut = "17:00"
	}

	tolerance := dto.LateToleranceMinutes
	if tolerance <= 0 {
		tolerance = dto.GracePeriodMinutes
	}
	if tolerance <= 0 {
		tolerance = 15
	}

	workDays := dto.WorkDays
	if len(workDays) == 0 {
		workDays = model.StringArray{"senin", "selasa", "rabu", "kamis", "jumat"}
	}

	ws := &model.MstHrWorkShift{
		ID:                   uuid.NewString(),
		ShiftName:            shiftName,
		TimeIn:               timeIn,
		TimeOut:              timeOut,
		LateToleranceMinutes: tolerance,
		WorkDays:             workDays,
		IsNightShift:         dto.IsNightShift,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := s.repo.CreateWorkShift(ws); err != nil {
		return nil, err
	}
	return ws, nil
}

func (s *hrService) UpdateWorkShift(id string, dto *WorkShiftDTO) (*model.MstHrWorkShift, error) {
	existing, err := s.repo.FindWorkShiftByID(id)
	if err != nil {
		return nil, errors.New("jadwal kerja tidak ditemukan")
	}

	shiftName := strings.TrimSpace(dto.ShiftName)
	if shiftName == "" {
		shiftName = strings.TrimSpace(dto.Name)
	}
	if shiftName != "" {
		existing.ShiftName = shiftName
	}
	timeIn := dto.TimeIn
	if timeIn == "" {
		timeIn = dto.StartTime
	}
	if timeIn != "" {
		existing.TimeIn = timeIn
	}
	timeOut := dto.TimeOut
	if timeOut == "" {
		timeOut = dto.EndTime
	}
	if timeOut != "" {
		existing.TimeOut = timeOut
	}
	tolerance := dto.LateToleranceMinutes
	if tolerance <= 0 {
		tolerance = dto.GracePeriodMinutes
	}
	if tolerance > 0 {
		existing.LateToleranceMinutes = tolerance
	}
	if len(dto.WorkDays) > 0 {
		existing.WorkDays = dto.WorkDays
	}
	existing.IsNightShift = dto.IsNightShift
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateWorkShift(existing); err != nil {
		return nil, err
	}
	return existing, nil
}
func (s *hrService) DeleteWorkShift(id string) error {
	return s.repo.DeleteWorkShift(id)
}

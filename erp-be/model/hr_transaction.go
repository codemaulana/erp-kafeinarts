package model

import (
	"time"
)

// TrxHrAttendance merepresentasikan data presensi karyawan (trx_hr_attendances)
type TrxHrAttendance struct {
	ID               string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	EmployeeID       string    `gorm:"type:varchar(36);not null;index" json:"employee_id"`
	EmployeeName     string    `gorm:"->" json:"employee_name,omitempty"`
	EmployeeNumber   string    `gorm:"->" json:"employee_number,omitempty"`
	DepartmentName   string    `gorm:"->" json:"department_name,omitempty"`
	AttendanceDate   string    `gorm:"type:varchar(20);not null;index" json:"attendance_date"` // YYYY-MM-DD
	CheckInTime      string    `gorm:"type:varchar(30);default:''" json:"check_in_time"`
	CheckOutTime     string    `gorm:"type:varchar(30);default:''" json:"check_out_time"`
	CheckInLat       float64   `gorm:"default:0" json:"check_in_lat"`
	CheckInLng       float64   `gorm:"default:0" json:"check_in_lng"`
	AttendanceStatus string    `gorm:"type:varchar(30);not null;default:'PRESENT'" json:"attendance_status"` // PRESENT, LATE, ABSENT, LEAVE
	Notes            string    `gorm:"type:text;default:''" json:"notes"`
	SelfiePhotoURL   string    `gorm:"type:text;default:''" json:"selfie_photo_url"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (TrxHrAttendance) TableName() string {
	return "trx_hr_attendances"
}

// TrxHrLeaveRequest merepresentasikan pengajuan cuti (trx_hr_leave_requests)
type TrxHrLeaveRequest struct {
	ID              string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	RequestNumber   string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"request_number"`
	EmployeeID      string    `gorm:"type:varchar(36);not null;index" json:"employee_id"`
	EmployeeName    string    `gorm:"->" json:"employee_name,omitempty"`
	EmployeeNumber  string    `gorm:"->" json:"employee_number,omitempty"`
	DepartmentName  string    `gorm:"->" json:"department_name,omitempty"`
	LeaveTypeID     string    `gorm:"type:varchar(36);not null" json:"leave_type_id"`
	LeaveName       string    `gorm:"->" json:"leave_name,omitempty"`
	StartDate       string    `gorm:"type:varchar(20);not null" json:"start_date"`
	EndDate         string    `gorm:"type:varchar(20);not null" json:"end_date"`
	TotalDays       int       `gorm:"default:1" json:"total_days"`
	Reason          string    `gorm:"type:text;default:''" json:"reason"`
	AttachmentURL   string    `gorm:"type:text;default:''" json:"attachment_url"`
	ApprovalStatus  string    `gorm:"type:varchar(30);not null;default:'SUBMITTED'" json:"approval_status"` // DRAFT, SUBMITTED, APPROVED, REJECTED
	Status          string    `gorm:"type:varchar(30);not null;default:'SUBMITTED'" json:"status"`
	RejectionReason string    `gorm:"type:text;default:''" json:"rejection_reason"`
	ApprovedBy      string    `gorm:"type:varchar(36);default:''" json:"approved_by"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

func (TrxHrLeaveRequest) TableName() string {
	return "trx_hr_leave_requests"
}

// TrxHrPayroll merepresentasikan transaksi penggajian (trx_hr_payrolls)
type TrxHrPayroll struct {
	ID              string             `gorm:"type:varchar(36);primaryKey" json:"id"`
	PayrollNumber   string             `gorm:"type:varchar(50);not null;uniqueIndex" json:"payroll_number"`
	PayrollPeriod   string             `gorm:"type:varchar(20);not null;index" json:"payroll_period"` // YYYY-MM
	PeriodMonth     int                `gorm:"default:0" json:"period_month"`
	PeriodYear      int                `gorm:"default:0" json:"period_year"`
	EmployeeID      string             `gorm:"type:varchar(36);not null;index" json:"employee_id"`
	EmployeeName    string             `gorm:"->" json:"employee_name,omitempty"`
	EmployeeNumber  string             `gorm:"->" json:"employee_number,omitempty"`
	DepartmentName  string             `gorm:"->" json:"department_name,omitempty"`
	TotalEarnings   float64            `gorm:"type:decimal(15,2);default:0" json:"total_earnings"`
	TotalDeductions float64            `gorm:"type:decimal(15,2);default:0" json:"total_deductions"`
	NetSalary       float64            `gorm:"type:decimal(15,2);default:0" json:"net_salary"`
	PaymentStatus   string             `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"payment_status"` // DRAFT, PROCESSED, PAID
	PaymentDate     string             `gorm:"type:varchar(20);default:''" json:"payment_date"`
	Notes           string             `gorm:"type:text;default:''" json:"notes"`
	Items           []TrxHrPayrollItem `gorm:"foreignKey:PayrollID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt       time.Time          `json:"created_at"`
	UpdatedAt       time.Time          `json:"updated_at"`
}

func (TrxHrPayroll) TableName() string {
	return "trx_hr_payrolls"
}

// TrxHrPayrollItem merepresentasikan rincian komponen gaji pada payroll (trx_hr_payroll_details)
type TrxHrPayrollItem struct {
	ID            string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	PayrollID     string    `gorm:"type:varchar(36);not null;index" json:"payroll_id"`
	ComponentID   string    `gorm:"type:varchar(36);not null" json:"component_id"`
	ComponentName string    `gorm:"->" json:"component_name,omitempty"`
	Amount        float64   `gorm:"type:decimal(15,2);default:0" json:"amount"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (TrxHrPayrollItem) TableName() string {
	return "trx_hr_payroll_details"
}

// TrxHrReimbursement merepresentasikan klaim reimbursement (trx_hr_reimbursements)
type TrxHrReimbursement struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ClaimNumber    string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"claim_number"`
	EmployeeID     string    `gorm:"type:varchar(36);not null;index" json:"employee_id"`
	EmployeeName   string    `gorm:"->" json:"employee_name,omitempty"`
	EmployeeNumber string    `gorm:"->" json:"employee_number,omitempty"`
	DepartmentName string    `gorm:"->" json:"department_name,omitempty"`
	ClaimDate      string    `gorm:"type:varchar(20);not null" json:"claim_date"`
	CategoryID     string    `gorm:"type:varchar(50);default:''" json:"category_id"`
	CategoryName   string    `gorm:"->" json:"category_name,omitempty"`
	Amount         float64   `gorm:"type:decimal(15,2);default:0" json:"amount"`
	Description    string    `gorm:"type:text;default:''" json:"description"`
	AttachmentURL  string    `gorm:"type:text;default:''" json:"attachment_url"`
	Status         string    `gorm:"type:varchar(30);not null;default:'SUBMITTED'" json:"status"` // SUBMITTED, APPROVED, REJECTED, PAID
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (TrxHrReimbursement) TableName() string {
	return "trx_hr_reimbursements"
}

// TrxHrPerformanceReview merepresentasikan penilaian kinerja (trx_hr_performance_reviews)
type TrxHrPerformanceReview struct {
	ID              string                       `gorm:"type:varchar(36);primaryKey" json:"id"`
	ReviewNumber    string                       `gorm:"type:varchar(50);not null;uniqueIndex" json:"review_number"`
	EmployeeID      string                       `gorm:"type:varchar(36);not null;index" json:"employee_id"`
	EmployeeName    string                       `gorm:"->" json:"employee_name,omitempty"`
	EmployeeNumber  string                       `gorm:"->" json:"employee_number,omitempty"`
	DepartmentName  string                       `gorm:"->" json:"department_name,omitempty"`
	ReviewerID      string                       `gorm:"type:varchar(36);not null;index" json:"reviewer_id"`
	ReviewerName    string                       `gorm:"->" json:"reviewer_name,omitempty"`
	ReviewPeriod    string                       `gorm:"type:varchar(50);not null" json:"review_period"` // e.g. Q1 2026, 2026-H1
	ReviewDate      string                       `gorm:"type:varchar(20);default:''" json:"review_date"`
	TotalScore      float64                      `gorm:"default:0" json:"total_score"`
	OverallFeedback string                       `gorm:"type:text;default:''" json:"overall_feedback"`
	Status          string                       `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"status"` // DRAFT, SUBMITTED, FINALIZED
	Items           []TrxHrPerformanceReviewItem `gorm:"foreignKey:ReviewID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt       time.Time                    `json:"created_at"`
	UpdatedAt       time.Time                    `json:"updated_at"`
}

func (TrxHrPerformanceReview) TableName() string {
	return "trx_hr_performance_reviews"
}

// TrxHrPerformanceReviewItem merepresentasikan rincian kriteria penilaian (trx_hr_performance_review_items)
type TrxHrPerformanceReviewItem struct {
	ID        string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	ReviewID  string    `gorm:"type:varchar(36);not null;index" json:"review_id"`
	Criterion string    `gorm:"type:varchar(150);not null" json:"criterion"`
	Score     float64   `gorm:"default:0" json:"score"`
	Weight    float64   `gorm:"default:0" json:"weight"`
	Comments  string    `gorm:"type:text;default:''" json:"comments"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (TrxHrPerformanceReviewItem) TableName() string {
	return "trx_hr_performance_review_items"
}

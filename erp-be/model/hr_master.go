package model

import (
	"database/sql/driver"
	"encoding/json"
	"fmt"
	"time"
)

// StringArray adalah helper serializer JSON untuk slice string di SQLite (misal work_days).
type StringArray []string

func (s *StringArray) Scan(value interface{}) error {
	if value == nil {
		*s = []string{}
		return nil
	}
	bytes, ok := value.([]byte)
	if !ok {
		str, ok := value.(string)
		if !ok {
			return fmt.Errorf("gagal konversi %T ke StringArray", value)
		}
		bytes = []byte(str)
	}
	if len(bytes) == 0 {
		*s = []string{}
		return nil
	}
	return json.Unmarshal(bytes, s)
}

func (s StringArray) Value() (driver.Value, error) {
	if len(s) == 0 {
		return "[]", nil
	}
	return json.Marshal(s)
}

// MstHrDepartment merepresentasikan master departemen (mst_hr_departments)
type MstHrDepartment struct {
	ID               string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code             string    `gorm:"type:varchar(20);default:''" json:"code"`
	DepartmentName   string    `gorm:"type:varchar(100);not null" json:"department_name"`
	HeadEmployeeID   *string   `gorm:"type:varchar(36);default:null" json:"head_employee_id"`
	HeadEmployeeName string    `gorm:"->" json:"head_employee_name,omitempty"`
	EmployeeCount    int64     `gorm:"->" json:"employee_count"`
	Description      string    `gorm:"type:text;default:''" json:"description"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (MstHrDepartment) TableName() string {
	return "mst_hr_departments"
}

// MstHrPosition merepresentasikan master jabatan (mst_hr_positions)
type MstHrPosition struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	DepartmentID   string    `gorm:"type:varchar(36);not null;default:''" json:"department_id"`
	DepartmentName string    `gorm:"->" json:"department_name,omitempty"`
	Code           string    `gorm:"type:varchar(20);default:''" json:"code"`
	PositionName   string    `gorm:"type:varchar(100);not null" json:"position_name"`
	Level          string    `gorm:"type:varchar(50);default:'STAFF'" json:"level"` // STAFF, SUPERVISOR, MANAGER, DIRECTOR
	JobDescription string    `gorm:"type:text;default:''" json:"job_description"`
	Description    string    `gorm:"type:text;default:''" json:"description"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (MstHrPosition) TableName() string {
	return "mst_hr_positions"
}

// MstHrEmployee merepresentasikan master karyawan (mst_hr_employees)
type MstHrEmployee struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	EmployeeNumber     string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"employee_number"` // NIP / employee_code
	FullName           string    `gorm:"type:varchar(150);not null" json:"full_name"`                 // name
	NationalIDNumber   string    `gorm:"type:varchar(20);default:''" json:"national_id_number"`       // NIK
	BirthDate          string    `gorm:"type:varchar(20);default:''" json:"birth_date"`
	Gender             string    `gorm:"type:varchar(10);default:''" json:"gender"`
	Address            string    `gorm:"type:text;default:''" json:"address"`
	Email              string    `gorm:"type:varchar(100);default:''" json:"email"`
	PhoneNumber        string    `gorm:"type:varchar(20);default:''" json:"phone_number"` // phone
	DepartmentID       *string   `gorm:"type:varchar(36);default:null" json:"department_id"`
	DepartmentName     string    `gorm:"->" json:"department_name,omitempty"`
	PositionID         *string   `gorm:"type:varchar(36);default:null" json:"position_id"`
	PositionName       string    `gorm:"->" json:"position_name,omitempty"`
	JoinDate           string    `gorm:"type:varchar(20);default:''" json:"join_date"`
	EmploymentStatus   string    `gorm:"type:varchar(30);default:'PERMANENT'" json:"employment_status"` // PERMANENT, CONTRACT, INTERN
	ActiveStatus       string    `gorm:"type:varchar(30);default:'ACTIVE'" json:"active_status"`       // ACTIVE, RESIGNED, TERMINATED
	Salary             float64   `gorm:"type:decimal(15,2);default:0" json:"salary"`
	BankName           string    `gorm:"type:varchar(50);default:''" json:"bank_name"`
	BankAccountNumber  string    `gorm:"type:varchar(50);default:''" json:"bank_account_number"` // bank_account
	PhotoURL           string    `gorm:"type:text;default:''" json:"photo_url"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (MstHrEmployee) TableName() string {
	return "mst_hr_employees"
}

// MstHrLeaveType merepresentasikan master jenis cuti (mst_hr_leave_types)
type MstHrLeaveType struct {
	ID                 string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code               string    `gorm:"type:varchar(20);default:''" json:"code"`
	LeaveName          string    `gorm:"type:varchar(100);not null" json:"leave_name"` // name
	AnnualQuotaDays    int       `gorm:"default:12" json:"annual_quota_days"`          // quota_days
	IsPaid             bool      `gorm:"default:true" json:"is_paid"`
	RequiresAttachment bool      `gorm:"default:false" json:"requires_attachment"`
	Description        string    `gorm:"type:text;default:''" json:"description"`
	CreatedAt          time.Time `json:"created_at"`
	UpdatedAt          time.Time `json:"updated_at"`
}

func (MstHrLeaveType) TableName() string {
	return "mst_hr_leave_types"
}

// MstHrSalaryComponent merepresentasikan master komponen gaji (mst_hr_salary_components)
type MstHrSalaryComponent struct {
	ID                string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code              string    `gorm:"type:varchar(20);default:''" json:"code"`
	ComponentName     string    `gorm:"type:varchar(100);not null" json:"component_name"`
	ComponentType     string    `gorm:"type:varchar(20);default:'EARNING'" json:"component_type"`     // EARNING, DEDUCTION
	CalculationMethod string    `gorm:"type:varchar(20);default:'FIXED'" json:"calculation_method"` // FIXED, PERCENTAGE
	IsTaxable         bool      `gorm:"default:false" json:"is_taxable"`
	Description       string    `gorm:"type:text;default:''" json:"description"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (MstHrSalaryComponent) TableName() string {
	return "mst_hr_salary_components"
}

// MstHrWorkShift merepresentasikan jadwal kerja / shift (mst_hr_work_shifts)
type MstHrWorkShift struct {
	ID                   string      `gorm:"type:varchar(36);primaryKey" json:"id"`
	ShiftName            string      `gorm:"type:varchar(100);not null" json:"shift_name"` // name
	TimeIn               string      `gorm:"type:varchar(10);not null;default:'08:00'" json:"time_in"`
	TimeOut              string      `gorm:"type:varchar(10);not null;default:'17:00'" json:"time_out"`
	LateToleranceMinutes int         `gorm:"default:15" json:"late_tolerance_minutes"` // grace_period_minutes
	WorkDays             StringArray `gorm:"type:text;default:'[]'" json:"work_days"`
	IsNightShift         bool        `gorm:"default:false" json:"is_night_shift"`
	CreatedAt            time.Time   `json:"created_at"`
	UpdatedAt            time.Time   `json:"updated_at"`
}

func (MstHrWorkShift) TableName() string {
	return "mst_hr_work_shifts"
}

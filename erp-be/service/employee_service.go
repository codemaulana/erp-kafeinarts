package service

import (
	"errors"
	"strings"
	"time"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"
	"erp-kafeinart/utils"
)

// DTO untuk penambahan pegawai HR baru.
type CreateEmployeeDTO struct {
	NIP        string    `json:"nip"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	Status     string    `json:"status"`
	JoinDate   time.Time `json:"join_date"`
	Salary     float64   `json:"salary"`
}

// DTO untuk update data pegawai.
type UpdateEmployeeDTO struct {
	NIP        string    `json:"nip"`
	Name       string    `json:"name"`
	Email      string    `json:"email"`
	Phone      string    `json:"phone"`
	Department string    `json:"department"`
	Position   string    `json:"position"`
	Status     string    `json:"status"`
	JoinDate   time.Time `json:"join_date"`
	Salary     float64   `json:"salary"`
}

// EmployeeService mendefinisikan bisnis logika pengelolaan pegawai.
type EmployeeService interface {
	GetAll(page, perPage int, search, department, status string) ([]model.Employee, utils.PaginationMeta, error)
	GetByID(id uint) (*model.Employee, error)
	Create(dto CreateEmployeeDTO) (*model.Employee, error)
	Update(id uint, dto UpdateEmployeeDTO) (*model.Employee, error)
	Delete(id uint) error
}

type employeeService struct {
	empRepo repository.EmployeeRepository
}

// NewEmployeeService menginisialisasi EmployeeService.
func NewEmployeeService(empRepo repository.EmployeeRepository) EmployeeService {
	return &employeeService{empRepo: empRepo}
}

func (s *employeeService) GetAll(page, perPage int, search, department, status string) ([]model.Employee, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.empRepo.FindAll(page, perPage, search, department, status)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *employeeService) GetByID(id uint) (*model.Employee, error) {
	emp, err := s.empRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("pegawai tidak ditemukan")
	}
	return emp, nil
}

func (s *employeeService) Create(dto CreateEmployeeDTO) (*model.Employee, error) {
	nip := strings.TrimSpace(dto.NIP)
	name := strings.TrimSpace(dto.Name)
	email := strings.TrimSpace(dto.Email)

	if nip == "" || name == "" || email == "" {
		return nil, errors.New("NIP, Nama, dan Email pegawai wajib diisi")
	}

	// Cek duplikasi NIP
	existingNIP, err := s.empRepo.FindByNIP(nip)
	if err != nil {
		return nil, err
	}
	if existingNIP != nil {
		return nil, errors.New("NIP sudah terdaftar")
	}

	// Cek duplikasi Email
	existingEmail, err := s.empRepo.FindByEmail(email)
	if err != nil {
		return nil, err
	}
	if existingEmail != nil {
		return nil, errors.New("Email sudah terdaftar untuk pegawai lain")
	}

	status := strings.TrimSpace(dto.Status)
	if status == "" {
		status = "active"
	}

	joinDate := dto.JoinDate
	if joinDate.IsZero() {
		joinDate = time.Now()
	}

	emp := &model.Employee{
		NIP:        nip,
		Name:       name,
		Email:      email,
		Phone:      dto.Phone,
		Department: dto.Department,
		Position:   dto.Position,
		Status:     status,
		JoinDate:   joinDate,
		Salary:     dto.Salary,
	}

	if err := s.empRepo.Create(emp); err != nil {
		return nil, err
	}

	return emp, nil
}

func (s *employeeService) Update(id uint, dto UpdateEmployeeDTO) (*model.Employee, error) {
	emp, err := s.empRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if emp == nil {
		return nil, errors.New("pegawai tidak ditemukan")
	}

	nip := strings.TrimSpace(dto.NIP)
	name := strings.TrimSpace(dto.Name)
	email := strings.TrimSpace(dto.Email)

	if nip == "" || name == "" || email == "" {
		return nil, errors.New("NIP, Nama, dan Email pegawai wajib diisi")
	}

	// Cek duplikasi NIP
	if nip != emp.NIP {
		existingNIP, err := s.empRepo.FindByNIP(nip)
		if err != nil {
			return nil, err
		}
		if existingNIP != nil && existingNIP.ID != id {
			return nil, errors.New("NIP sudah terdaftar")
		}
	}

	// Cek duplikasi Email
	if email != emp.Email {
		existingEmail, err := s.empRepo.FindByEmail(email)
		if err != nil {
			return nil, err
		}
		if existingEmail != nil && existingEmail.ID != id {
			return nil, errors.New("Email sudah terdaftar")
		}
	}

	status := strings.TrimSpace(dto.Status)
	if status == "" {
		status = emp.Status
	}

	joinDate := dto.JoinDate
	if joinDate.IsZero() {
		joinDate = emp.JoinDate
	}

	emp.NIP = nip
	emp.Name = name
	emp.Email = email
	emp.Phone = dto.Phone
	emp.Department = dto.Department
	emp.Position = dto.Position
	emp.Status = status
	emp.JoinDate = joinDate
	emp.Salary = dto.Salary

	if err := s.empRepo.Update(emp); err != nil {
		return nil, err
	}

	return emp, nil
}

func (s *employeeService) Delete(id uint) error {
	emp, err := s.empRepo.FindByID(id)
	if err != nil {
		return err
	}
	if emp == nil {
		return errors.New("pegawai tidak ditemukan")
	}

	return s.empRepo.Delete(id)
}

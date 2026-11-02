package repository

import (
	"errors"
	"strings"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// EmployeeRepository mendefinisikan kontrak akses data untuk pegawai HR.
type EmployeeRepository interface {
	FindAll(page, perPage int, search, department, status string) ([]model.Employee, int64, error)
	FindByID(id uint) (*model.Employee, error)
	FindByNIP(nip string) (*model.Employee, error)
	FindByEmail(email string) (*model.Employee, error)
	Create(emp *model.Employee) error
	Update(emp *model.Employee) error
	Delete(id uint) error
}

type employeeRepository struct {
	db *gorm.DB
}

// NewEmployeeRepository menginisialisasi implementasi EmployeeRepository.
func NewEmployeeRepository(db *gorm.DB) EmployeeRepository {
	return &employeeRepository{db: db}
}

func (r *employeeRepository) FindAll(page, perPage int, search, department, status string) ([]model.Employee, int64, error) {
	var list []model.Employee
	var total int64

	query := r.db.Model(&model.Employee{})

	if strings.TrimSpace(search) != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(nip) LIKE ? OR LOWER(email) LIKE ?", pattern, pattern, pattern)
	}

	if strings.TrimSpace(department) != "" {
		query = query.Where("department = ?", department)
	}

	if strings.TrimSpace(status) != "" {
		query = query.Where("status = ?", status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("id ASC").Limit(perPage).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *employeeRepository) FindByID(id uint) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.First(&emp, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) FindByNIP(nip string) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.Where("nip = ?", nip).First(&emp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) FindByEmail(email string) (*model.Employee, error) {
	var emp model.Employee
	err := r.db.Where("email = ?", email).First(&emp).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &emp, nil
}

func (r *employeeRepository) Create(emp *model.Employee) error {
	return r.db.Create(emp).Error
}

func (r *employeeRepository) Update(emp *model.Employee) error {
	return r.db.Save(emp).Error
}

func (r *employeeRepository) Delete(id uint) error {
	return r.db.Delete(&model.Employee{}, id).Error
}

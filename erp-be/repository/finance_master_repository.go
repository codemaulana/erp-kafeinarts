package repository

import (
	"errors"
	"strings"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// FinanceMasterRepository mendefinisikan operasi data untuk Master Kategori, Rekening Kas/Bank, dan Pajak.
type FinanceMasterRepository interface {
	// Kategori
	FindCategories(cashType string) ([]model.MstFinanceCategory, error)
	FindCategoriesPaginated(page, perPage int, cashType, search string) ([]model.MstFinanceCategory, int64, error)
	FindCategoryByID(id string) (*model.MstFinanceCategory, error)
	CreateCategory(cat *model.MstFinanceCategory) error
	UpdateCategory(cat *model.MstFinanceCategory) error
	DeleteCategory(id string) error

	// Rekening / Kas Bank
	FindAccounts() ([]model.MstFinanceAccount, error)
	FindAccountsPaginated(page, perPage int, search string) ([]model.MstFinanceAccount, int64, error)
	FindAccountByID(id string) (*model.MstFinanceAccount, error)
	CreateAccount(acc *model.MstFinanceAccount) error
	UpdateAccount(acc *model.MstFinanceAccount) error
	UpdateAccountBalance(id string, delta float64) error
	DeleteAccount(id string) error

	// Pajak
	FindTaxes() ([]model.MstFinanceTax, error)
	FindTaxesPaginated(page, perPage int, search string) ([]model.MstFinanceTax, int64, error)
	FindTaxByID(id string) (*model.MstFinanceTax, error)
	CreateTax(tax *model.MstFinanceTax) error
	UpdateTax(tax *model.MstFinanceTax) error
	DeleteTax(id string) error
}

type financeMasterRepository struct {
	db *gorm.DB
}

func NewFinanceMasterRepository(db *gorm.DB) FinanceMasterRepository {
	return &financeMasterRepository{db: db}
}

// ── Kategori ────────────────────────────────────────────────────────────────

func (r *financeMasterRepository) FindCategories(cashType string) ([]model.MstFinanceCategory, error) {
	var categories []model.MstFinanceCategory
	query := r.db.Model(&model.MstFinanceCategory{})
	if cashType != "" {
		if strings.EqualFold(cashType, "masuk") {
			cashType = "INCOME"
		} else if strings.EqualFold(cashType, "keluar") {
			cashType = "EXPENSE"
		}
		query = query.Where("cash_type = ?", cashType)
	}
	err := query.Order("category_name ASC").Find(&categories).Error
	return categories, err
}

func (r *financeMasterRepository) FindCategoriesPaginated(page, perPage int, cashType, search string) ([]model.MstFinanceCategory, int64, error) {
	var categories []model.MstFinanceCategory
	var total int64

	query := r.db.Model(&model.MstFinanceCategory{})
	if cashType != "" {
		if strings.EqualFold(cashType, "masuk") {
			cashType = "INCOME"
		} else if strings.EqualFold(cashType, "keluar") {
			cashType = "EXPENSE"
		}
		query = query.Where("cash_type = ?", cashType)
	}
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("category_name LIKE ? OR description LIKE ? OR cash_type LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	err := query.Order("category_name ASC").Offset(offset).Limit(perPage).Find(&categories).Error
	return categories, total, err
}

func (r *financeMasterRepository) FindCategoryByID(id string) (*model.MstFinanceCategory, error) {
	var cat model.MstFinanceCategory
	err := r.db.First(&cat, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

func (r *financeMasterRepository) CreateCategory(cat *model.MstFinanceCategory) error {
	return r.db.Create(cat).Error
}

func (r *financeMasterRepository) UpdateCategory(cat *model.MstFinanceCategory) error {
	return r.db.Save(cat).Error
}

func (r *financeMasterRepository) DeleteCategory(id string) error {
	return r.db.Delete(&model.MstFinanceCategory{}, "id = ?", id).Error
}

// ── Rekening / Kas Bank ─────────────────────────────────────────────────────

func (r *financeMasterRepository) FindAccounts() ([]model.MstFinanceAccount, error) {
	var accounts []model.MstFinanceAccount
	err := r.db.Order("bank_name ASC").Find(&accounts).Error
	return accounts, err
}

func (r *financeMasterRepository) FindAccountsPaginated(page, perPage int, search string) ([]model.MstFinanceAccount, int64, error) {
	var accounts []model.MstFinanceAccount
	var total int64

	query := r.db.Model(&model.MstFinanceAccount{})
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("bank_name LIKE ? OR account_number LIKE ? OR account_holder_name LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	err := query.Order("bank_name ASC").Offset(offset).Limit(perPage).Find(&accounts).Error
	return accounts, total, err
}

func (r *financeMasterRepository) FindAccountByID(id string) (*model.MstFinanceAccount, error) {
	var acc model.MstFinanceAccount
	err := r.db.First(&acc, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

func (r *financeMasterRepository) CreateAccount(acc *model.MstFinanceAccount) error {
	return r.db.Create(acc).Error
}

func (r *financeMasterRepository) UpdateAccount(acc *model.MstFinanceAccount) error {
	return r.db.Save(acc).Error
}

func (r *financeMasterRepository) UpdateAccountBalance(id string, delta float64) error {
	return r.db.Model(&model.MstFinanceAccount{}).
		Where("id = ?", id).
		Update("current_balance", gorm.Expr("current_balance + ?", delta)).Error
}

func (r *financeMasterRepository) DeleteAccount(id string) error {
	return r.db.Delete(&model.MstFinanceAccount{}, "id = ?", id).Error
}

// ── Master Pajak ────────────────────────────────────────────────────────────

func (r *financeMasterRepository) FindTaxes() ([]model.MstFinanceTax, error) {
	var taxes []model.MstFinanceTax
	err := r.db.Order("percentage ASC").Find(&taxes).Error
	return taxes, err
}

func (r *financeMasterRepository) FindTaxesPaginated(page, perPage int, search string) ([]model.MstFinanceTax, int64, error) {
	var taxes []model.MstFinanceTax
	var total int64

	query := r.db.Model(&model.MstFinanceTax{})
	if search != "" {
		pattern := "%" + search + "%"
		query = query.Where("tax_name LIKE ? OR description LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage

	err := query.Order("percentage ASC").Offset(offset).Limit(perPage).Find(&taxes).Error
	return taxes, total, err
}

func (r *financeMasterRepository) FindTaxByID(id string) (*model.MstFinanceTax, error) {
	var tax model.MstFinanceTax
	err := r.db.First(&tax, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tax, nil
}

func (r *financeMasterRepository) CreateTax(tax *model.MstFinanceTax) error {
	return r.db.Create(tax).Error
}

func (r *financeMasterRepository) UpdateTax(tax *model.MstFinanceTax) error {
	return r.db.Save(tax).Error
}

func (r *financeMasterRepository) DeleteTax(id string) error {
	return r.db.Delete(&model.MstFinanceTax{}, "id = ?", id).Error
}

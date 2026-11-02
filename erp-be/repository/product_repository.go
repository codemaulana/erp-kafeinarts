package repository

import (
	"errors"
	"fmt"
	"strings"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// ProductRepository mendefinisikan kontrak akses data untuk entitas Product & mutasi stok.
type ProductRepository interface {
	FindAll(page, perPage int, search, sortBy string, descending bool) ([]model.Product, int64, error)
	FindByID(id uint) (*model.Product, error)
	FindBySKU(sku string) (*model.Product, error)
	Create(product *model.Product) error
	Update(product *model.Product) error
	Delete(id uint) error
	CreateAdjustment(adj *model.StockAdjustment) error
}

type productRepository struct {
	db *gorm.DB
}

// NewProductRepository menginisialisasi implementasi ProductRepository.
func NewProductRepository(db *gorm.DB) ProductRepository {
	return &productRepository{db: db}
}

func (r *productRepository) FindAll(page, perPage int, search, sortBy string, descending bool) ([]model.Product, int64, error) {
	var products []model.Product
	var total int64

	query := r.db.Model(&model.Product{})

	// Filter pencarian berdasarkan nama produk atau SKU
	if strings.TrimSpace(search) != "" {
		searchPattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(name) LIKE ? OR LOWER(sku) LIKE ?", searchPattern, searchPattern)
	}

	// Hitung total data sesuai filter
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Pengurutan (Sorting)
	validColumns := map[string]bool{
		"id":    true,
		"sku":   true,
		"name":  true,
		"stock": true,
		"price": true,
	}

	orderClause := "id DESC"
	if sortBy != "" && validColumns[strings.ToLower(sortBy)] {
		direction := "ASC"
		if descending {
			direction = "DESC"
		}
		orderClause = fmt.Sprintf("%s %s", sortBy, direction)
	}

	offset := (page - 1) * perPage
	err := query.Order(orderClause).Limit(perPage).Offset(offset).Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

func (r *productRepository) FindByID(id uint) (*model.Product, error) {
	var product model.Product
	err := r.db.First(&product, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) FindBySKU(sku string) (*model.Product, error) {
	var product model.Product
	err := r.db.Where("sku = ?", sku).First(&product).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &product, nil
}

func (r *productRepository) Create(product *model.Product) error {
	return r.db.Create(product).Error
}

func (r *productRepository) Update(product *model.Product) error {
	return r.db.Save(product).Error
}

func (r *productRepository) Delete(id uint) error {
	return r.db.Delete(&model.Product{}, id).Error
}

func (r *productRepository) CreateAdjustment(adj *model.StockAdjustment) error {
	return r.db.Create(adj).Error
}

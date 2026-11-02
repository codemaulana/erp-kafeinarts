package service

import (
	"errors"
	"strings"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"
	"erp-kafeinart/utils"
)

// DTO untuk pembuatan produk baru.
type CreateProductDTO struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// DTO untuk update produk.
type UpdateProductDTO struct {
	SKU         string  `json:"sku"`
	Name        string  `json:"name"`
	Stock       int     `json:"stock"`
	Price       float64 `json:"price"`
	Description string  `json:"description"`
}

// DTO untuk penyesuaian stok produk.
type StockAdjustmentDTO struct {
	Adjustment int    `json:"adjustment"`
	Type       string `json:"type"`   // "in", "out", "set"
	Reason     string `json:"reason"` // Alasan penyesuaian
}

// ProductService mendefinisikan bisnis logika inventaris produk.
type ProductService interface {
	GetAll(page, perPage int, search, sortBy string, descending bool) ([]model.Product, utils.PaginationMeta, error)
	GetByID(id uint) (*model.Product, error)
	Create(dto CreateProductDTO) (*model.Product, error)
	Update(id uint, dto UpdateProductDTO) (*model.Product, error)
	Delete(id uint) error
	AdjustStock(id uint, dto StockAdjustmentDTO) (*model.Product, error)
}

type productService struct {
	productRepo repository.ProductRepository
}

// NewProductService menginisialisasi ProductService.
func NewProductService(productRepo repository.ProductRepository) ProductService {
	return &productService{productRepo: productRepo}
}

func (s *productService) GetAll(page, perPage int, search, sortBy string, descending bool) ([]model.Product, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	products, total, err := s.productRepo.FindAll(page, perPage, search, sortBy, descending)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return products, meta, nil
}

func (s *productService) GetByID(id uint) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produk tidak ditemukan")
	}
	return product, nil
}

func (s *productService) Create(dto CreateProductDTO) (*model.Product, error) {
	sku := strings.TrimSpace(dto.SKU)
	name := strings.TrimSpace(dto.Name)

	if sku == "" || name == "" {
		return nil, errors.New("SKU dan Nama Produk wajib diisi")
	}

	existing, err := s.productRepo.FindBySKU(sku)
	if err != nil {
		return nil, err
	}
	if existing != nil {
		return nil, errors.New("SKU sudah terdaftar di sistem")
	}

	product := &model.Product{
		SKU:         sku,
		Name:        name,
		Stock:       dto.Stock,
		Price:       dto.Price,
		Description: dto.Description,
	}

	if err := s.productRepo.Create(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Update(id uint, dto UpdateProductDTO) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	sku := strings.TrimSpace(dto.SKU)
	name := strings.TrimSpace(dto.Name)

	if sku == "" || name == "" {
		return nil, errors.New("SKU dan Nama Produk wajib diisi")
	}

	// Cek SKU jika diubah
	if sku != product.SKU {
		existing, err := s.productRepo.FindBySKU(sku)
		if err != nil {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, errors.New("SKU sudah digunakan oleh produk lain")
		}
	}

	product.SKU = sku
	product.Name = name
	product.Stock = dto.Stock
	product.Price = dto.Price
	product.Description = dto.Description

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) Delete(id uint) error {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return err
	}
	if product == nil {
		return errors.New("produk tidak ditemukan")
	}

	return s.productRepo.Delete(id)
}

func (s *productService) AdjustStock(id uint, dto StockAdjustmentDTO) (*model.Product, error) {
	product, err := s.productRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if product == nil {
		return nil, errors.New("produk tidak ditemukan")
	}

	switch strings.ToLower(dto.Type) {
	case "in":
		product.Stock += dto.Adjustment
	case "out":
		if product.Stock < dto.Adjustment {
			return nil, errors.New("stok tidak mencukupi untuk pengurangan")
		}
		product.Stock -= dto.Adjustment
	case "set":
		product.Stock = dto.Adjustment
	default:
		// Default: tambah jika positif, kurang jika negatif
		if product.Stock+dto.Adjustment < 0 {
			return nil, errors.New("stok tidak boleh negatif")
		}
		product.Stock += dto.Adjustment
	}

	if err := s.productRepo.Update(product); err != nil {
		return nil, err
	}

	// Catat riwayat penyesuaian
	adj := &model.StockAdjustment{
		ProductID: id,
		Type:      dto.Type,
		Amount:    dto.Adjustment,
		Reason:    dto.Reason,
	}
	_ = s.productRepo.CreateAdjustment(adj)

	return product, nil
}

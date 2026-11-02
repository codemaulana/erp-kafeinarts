package handler

import (
	"strconv"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
)

type ProductHandler struct {
	productService service.ProductService
}

func NewProductHandler(productService service.ProductService) *ProductHandler {
	return &ProductHandler{productService: productService}
}

// GetAll menangani GET /api/v1/inventory/products
func (h *ProductHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	// Dukung param `per_page` maupun `perPage` dari Quasar table request
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 10
	}

	search := c.Query("search")
	if search == "" {
		search = c.Query("filter")
	}

	sortBy := c.Query("sortBy")
	descending := c.Query("descending") == "true"

	products, meta, err := h.productService.GetAll(page, perPage, search, sortBy, descending)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar produk", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, products, meta, "Daftar produk berhasil dimuat")
}

// GetByID menangani GET /api/v1/inventory/products/:id
func (h *ProductHandler) GetByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	product, err := h.productService.GetByID(uint(id))
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, product, "Detail produk berhasil dimuat")
}

// Create menangani POST /api/v1/inventory/products
func (h *ProductHandler) Create(c fiber.Ctx) error {
	var req service.CreateProductDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input data tidak valid", err.Error())
	}

	product, err := h.productService.Create(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, product, "Produk berhasil ditambahkan")
}

// Update menangani PUT /api/v1/inventory/products/:id
func (h *ProductHandler) Update(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	var req service.UpdateProductDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input data tidak valid", err.Error())
	}

	product, err := h.productService.Update(uint(id), req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, product, "Produk berhasil diperbarui")
}

// Delete menangani DELETE /api/v1/inventory/products/:id
func (h *ProductHandler) Delete(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	if err := h.productService.Delete(uint(id)); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Produk berhasil dihapus")
}

// AdjustStock menangani POST /api/v1/inventory/products/:id/adjust
func (h *ProductHandler) AdjustStock(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID produk tidak valid")
	}

	var req service.StockAdjustmentDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format penyesuaian stok tidak valid", err.Error())
	}

	product, err := h.productService.AdjustStock(uint(id), req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, product, "Stok produk berhasil disesuaikan")
}

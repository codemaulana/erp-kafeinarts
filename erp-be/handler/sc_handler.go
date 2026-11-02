package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type ScHandler struct {
	scService service.ScService
}

func NewScHandler(scService service.ScService) *ScHandler {
	return &ScHandler{scService: scService}
}

// ── Master Vendor Handlers ──────────────────────────────────────────

func (h *ScHandler) GetVendors(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	status := c.Query("status")
	if status == "" {
		status = c.Query("active_status")
	}

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.scService.GetVendorsPaginated(page, perPage, search, status)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat vendor", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.scService.GetAllVendors(search, status)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat vendor", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar vendor berhasil dimuat")
}

func (h *ScHandler) GetVendorByID(c fiber.Ctx) error {
	id := c.Params("id")
	v, err := h.scService.GetVendorByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Vendor tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, v, "Detail vendor berhasil dimuat")
}

func (h *ScHandler) CreateVendor(c fiber.Ctx) error {
	var dto service.VendorDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	v, err := h.scService.CreateVendor(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, v, "Vendor berhasil dibuat")
}

func (h *ScHandler) UpdateVendor(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.VendorDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	v, err := h.scService.UpdateVendor(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, v, "Vendor berhasil diperbarui")
}

func (h *ScHandler) DeleteVendor(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteVendor(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus vendor", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Vendor berhasil dihapus")
}

// ── Master Kategori Item Handlers ───────────────────────────────────

func (h *ScHandler) GetItemCategories(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	itemType := c.Query("item_type")
	if itemType == "" {
		itemType = c.Query("type")
	}

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.scService.GetItemCategoriesPaginated(page, perPage, search, itemType)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat kategori item", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.scService.GetAllItemCategories(search, itemType)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat kategori item", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar kategori item berhasil dimuat")
}

func (h *ScHandler) GetItemCategoryByID(c fiber.Ctx) error {
	id := c.Params("id")
	cat, err := h.scService.GetItemCategoryByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Kategori item tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, cat, "Detail kategori item berhasil dimuat")
}

func (h *ScHandler) CreateItemCategory(c fiber.Ctx) error {
	var dto service.ItemCategoryDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	cat, err := h.scService.CreateItemCategory(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, cat, "Kategori item berhasil dibuat")
}

func (h *ScHandler) UpdateItemCategory(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.ItemCategoryDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	cat, err := h.scService.UpdateItemCategory(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, cat, "Kategori item berhasil diperbarui")
}

func (h *ScHandler) DeleteItemCategory(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteItemCategory(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus kategori item", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Kategori item berhasil dihapus")
}

// ── Master Katalog Item Handlers ────────────────────────────────────

func (h *ScHandler) GetItems(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	catID := c.Query("category_id")

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.scService.GetItemsPaginated(page, perPage, catID, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat katalog item", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.scService.GetAllItems(catID, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat katalog item", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar item berhasil dimuat")
}

func (h *ScHandler) GetItemByID(c fiber.Ctx) error {
	id := c.Params("id")
	item, err := h.scService.GetItemByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Item tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, item, "Detail item berhasil dimuat")
}

func (h *ScHandler) CreateItem(c fiber.Ctx) error {
	var dto service.ItemDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	item, err := h.scService.CreateItem(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, item, "Item berhasil dibuat")
}

func (h *ScHandler) UpdateItem(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.ItemDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	item, err := h.scService.UpdateItem(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, item, "Item berhasil diperbarui")
}

func (h *ScHandler) DeleteItem(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteItem(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus item", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Item berhasil dihapus")
}

// ── Master Lokasi Handlers ──────────────────────────────────────────

func (h *ScHandler) GetLocations(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")

	if pageStr != "" || perPageStr != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		if page <= 0 {
			page = 1
		}
		if perPage <= 0 {
			perPage = 10
		}
		list, total, err := h.scService.GetLocationsPaginated(page, perPage, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat lokasi", err.Error())
		}
		lastPage := int((total + int64(perPage) - 1) / int64(perPage))
		if lastPage < 1 {
			lastPage = 1
		}
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"success": true,
			"data":    list,
			"meta": fiber.Map{
				"total":    total,
				"page":     page,
				"perPage":  perPage,
				"lastPage": lastPage,
			},
		})
	}

	list, err := h.scService.GetAllLocations(search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat lokasi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Daftar lokasi berhasil dimuat")
}

func (h *ScHandler) GetLocationByID(c fiber.Ctx) error {
	id := c.Params("id")
	loc, err := h.scService.GetLocationByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Lokasi tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, loc, "Detail lokasi berhasil dimuat")
}

func (h *ScHandler) CreateLocation(c fiber.Ctx) error {
	var dto service.LocationDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	loc, err := h.scService.CreateLocation(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, loc, "Lokasi berhasil dibuat")
}

func (h *ScHandler) UpdateLocation(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.LocationDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	loc, err := h.scService.UpdateLocation(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, loc, "Lokasi berhasil diperbarui")
}

func (h *ScHandler) DeleteLocation(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteLocation(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus lokasi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Lokasi berhasil dihapus")
}

// ── Transaksi Purchase Request Handlers ─────────────────────────────

func (h *ScHandler) GetPurchaseRequests(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	deptID := c.Query("department_id")
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.scService.GetPurchaseRequestsPaginated(page, perPage, deptID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat purchase request", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *ScHandler) GetPurchaseRequestByID(c fiber.Ctx) error {
	id := c.Params("id")
	pr, err := h.scService.GetPurchaseRequestByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Purchase request tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pr, "Detail purchase request berhasil dimuat")
}

func (h *ScHandler) CreatePurchaseRequest(c fiber.Ctx) error {
	var dto service.PurchaseRequestDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	pr, err := h.scService.CreatePurchaseRequest(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, pr, "Purchase request berhasil diajukan")
}

func (h *ScHandler) UpdatePurchaseRequest(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.PurchaseRequestDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	pr, err := h.scService.UpdatePurchaseRequest(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, pr, "Purchase request berhasil diperbarui")
}

func (h *ScHandler) UpdatePurchaseRequestStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.scService.UpdatePurchaseRequestStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.scService.GetPurchaseRequestByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status PR berhasil diperbarui")
}

func (h *ScHandler) DeletePurchaseRequest(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeletePurchaseRequest(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus purchase request", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Purchase request berhasil dihapus")
}

// ── Transaksi Purchase Order Handlers ───────────────────────────────

func (h *ScHandler) GetPurchaseOrders(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	vendorID := c.Query("vendor_id")
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.scService.GetPurchaseOrdersPaginated(page, perPage, vendorID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat purchase order", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *ScHandler) GetPurchaseOrderByID(c fiber.Ctx) error {
	id := c.Params("id")
	po, err := h.scService.GetPurchaseOrderByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Purchase order tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, po, "Detail purchase order berhasil dimuat")
}

func (h *ScHandler) CreatePurchaseOrder(c fiber.Ctx) error {
	var dto service.PurchaseOrderDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	po, err := h.scService.CreatePurchaseOrder(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, po, "Purchase order berhasil dibuat")
}

func (h *ScHandler) UpdatePurchaseOrder(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.PurchaseOrderDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	po, err := h.scService.UpdatePurchaseOrder(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, po, "Purchase order berhasil diperbarui")
}

func (h *ScHandler) UpdatePurchaseOrderStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.scService.UpdatePurchaseOrderStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.scService.GetPurchaseOrderByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status PO berhasil diperbarui")
}

func (h *ScHandler) DeletePurchaseOrder(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeletePurchaseOrder(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus purchase order", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Purchase order berhasil dihapus")
}

// ── Transaksi Goods Receipt Handlers ────────────────────────────────

func (h *ScHandler) GetGoodsReceipts(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	poID := c.Query("po_id")
	search := c.Query("search")

	list, total, err := h.scService.GetGoodsReceiptsPaginated(page, perPage, poID, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat goods receipt", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *ScHandler) GetGoodsReceiptByID(c fiber.Ctx) error {
	id := c.Params("id")
	gr, err := h.scService.GetGoodsReceiptByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Goods receipt tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, gr, "Detail goods receipt berhasil dimuat")
}

func (h *ScHandler) CreateGoodsReceipt(c fiber.Ctx) error {
	var dto service.GoodsReceiptDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	gr, err := h.scService.CreateGoodsReceipt(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, gr, "Penerimaan barang berhasil dicatat")
}

func (h *ScHandler) UpdateGoodsReceipt(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.GoodsReceiptDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	gr, err := h.scService.UpdateGoodsReceipt(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, gr, "Penerimaan barang berhasil diperbarui")
}

func (h *ScHandler) DeleteGoodsReceipt(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteGoodsReceipt(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus goods receipt", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Goods receipt berhasil dihapus")
}

// ── Transaksi Manajemen Aset Handlers ───────────────────────────────

func (h *ScHandler) GetAssets(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	catID := c.Query("category_id")
	locID := c.Query("location_id")
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.scService.GetAssetsPaginated(page, perPage, catID, locID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat data aset", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *ScHandler) GetAssetByID(c fiber.Ctx) error {
	id := c.Params("id")
	asset, err := h.scService.GetAssetByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Aset tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, asset, "Detail aset berhasil dimuat")
}

func (h *ScHandler) CreateAsset(c fiber.Ctx) error {
	var dto service.AssetDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	asset, err := h.scService.CreateAsset(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, asset, "Aset berhasil ditambahkan")
}

func (h *ScHandler) UpdateAsset(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.AssetDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	asset, err := h.scService.UpdateAsset(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, asset, "Aset berhasil diperbarui")
}

func (h *ScHandler) UpdateAssetStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.scService.UpdateAssetStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.scService.GetAssetByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status aset berhasil diperbarui")
}

func (h *ScHandler) DeleteAsset(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteAsset(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus aset", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Aset berhasil dihapus")
}

// ── Transaksi Perpanjangan Langganan Handlers ───────────────────────

func (h *ScHandler) GetSubscriptions(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "10")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}

	vendorID := c.Query("vendor_id")
	status := c.Query("status")
	search := c.Query("search")

	list, total, err := h.scService.GetSubscriptionsPaginated(page, perPage, vendorID, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat langganan", err.Error())
	}
	lastPage := int((total + int64(perPage) - 1) / int64(perPage))
	if lastPage < 1 {
		lastPage = 1
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"success": true,
		"data":    list,
		"meta": fiber.Map{
			"total":    total,
			"page":     page,
			"perPage":  perPage,
			"lastPage": lastPage,
		},
	})
}

func (h *ScHandler) GetSubscriptionByID(c fiber.Ctx) error {
	id := c.Params("id")
	sub, err := h.scService.GetSubscriptionByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, "Langganan tidak ditemukan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, sub, "Detail langganan berhasil dimuat")
}

func (h *ScHandler) CreateSubscription(c fiber.Ctx) error {
	var dto service.SubscriptionDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	sub, err := h.scService.CreateSubscription(&dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, sub, "Langganan berhasil dibuat")
}

func (h *ScHandler) UpdateSubscription(c fiber.Ctx) error {
	id := c.Params("id")
	var dto service.SubscriptionDTO
	if err := c.Bind().Body(&dto); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Data tidak valid", err.Error())
	}
	sub, err := h.scService.UpdateSubscription(id, &dto)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, sub, "Langganan berhasil diperbarui")
}

func (h *ScHandler) UpdateSubscriptionStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req struct {
		Status string `json:"status"`
	}
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}
	if err := h.scService.UpdateSubscriptionStatus(id, req.Status); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	updated, _ := h.scService.GetSubscriptionByID(id)
	return utils.SendSuccess(c, fiber.StatusOK, updated, "Status langganan berhasil diperbarui")
}

func (h *ScHandler) DeleteSubscription(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.scService.DeleteSubscription(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Gagal menghapus langganan", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Langganan berhasil dihapus")
}

// ── File Upload Handler ─────────────────────────────────────────────

func (h *ScHandler) UploadFile(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "File tidak ditemukan", err.Error())
	}

	uploadDir := "./storage/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat folder upload", err.Error())
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("sc_%s_%s%s", time.Now().Format("20060102150405"), uuid.NewString()[:8], ext)
	dst := filepath.Join(uploadDir, newFileName)

	if err := c.SaveFile(file, dst); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan file", err.Error())
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)
	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{
		"file_name": file.Filename,
		"file_url":  fileURL,
		"size":      file.Size,
	}, "File berhasil diunggah")
}

// ── Laporan Supply Chain Handlers ───────────────────────────────────

func (h *ScHandler) GetLaporanPO(c fiber.Ctx) error {
	start := c.Query("start_date")
	end := c.Query("end_date")
	vendorID := c.Query("vendor_id")
	status := c.Query("status")

	list, err := h.scService.GetLaporanPO(start, end, vendorID, status)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat laporan PO", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Laporan purchase order berhasil dimuat")
}

func (h *ScHandler) GetLaporanAset(c fiber.Ctx) error {
	catID := c.Query("category_id")
	locID := c.Query("location_id")
	status := c.Query("status")
	deptID := c.Query("department_id")
	empID := c.Query("employee_id")

	list, err := h.scService.GetLaporanAset(catID, locID, status, deptID, empID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat laporan aset", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Laporan aset berhasil dimuat")
}

func (h *ScHandler) GetLaporanLisensi(c fiber.Ctx) error {
	status := c.Query("status")
	vendorID := c.Query("vendor_id")

	list, err := h.scService.GetLaporanLisensi(status, vendorID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat laporan lisensi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Laporan lisensi & langganan berhasil dimuat")
}

func (h *ScHandler) GetRiwayatVendor(c fiber.Ctx) error {
	vendorID := c.Params("vendor_id")
	if vendorID == "" {
		vendorID = c.Query("vendor_id")
	}
	startDate := c.Query("start_date")
	endDate := c.Query("end_date")

	list, err := h.scService.GetRiwayatVendor(vendorID, startDate, endDate)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat riwayat transaksi vendor", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, list, "Laporan riwayat transaksi vendor berhasil dimuat")
}

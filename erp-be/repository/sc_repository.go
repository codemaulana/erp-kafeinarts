package repository

import (
	"erp-kafeinart/model"
	"strings"

	"gorm.io/gorm"
)

type ScRepository interface {
	// ── Master Vendor ───────────────────────────────
	FindVendors(search, status string) ([]model.MstScVendor, error)
	FindVendorsPaginated(page, perPage int, search, status string) ([]model.MstScVendor, int64, error)
	FindVendorByID(id string) (*model.MstScVendor, error)
	CreateVendor(v *model.MstScVendor) error
	UpdateVendor(v *model.MstScVendor) error
	DeleteVendor(id string) error

	// ── Master Kategori Item ────────────────────────
	FindItemCategories(search, itemType string) ([]model.MstScItemCategory, error)
	FindItemCategoriesPaginated(page, perPage int, search, itemType string) ([]model.MstScItemCategory, int64, error)
	FindItemCategoryByID(id string) (*model.MstScItemCategory, error)
	CreateItemCategory(cat *model.MstScItemCategory) error
	UpdateItemCategory(cat *model.MstScItemCategory) error
	DeleteItemCategory(id string) error

	// ── Master Katalog Item ─────────────────────────
	FindItems(categoryID, search string) ([]model.MstScItem, error)
	FindItemsPaginated(page, perPage int, categoryID, search string) ([]model.MstScItem, int64, error)
	FindItemByID(id string) (*model.MstScItem, error)
	CreateItem(item *model.MstScItem) error
	UpdateItem(item *model.MstScItem) error
	DeleteItem(id string) error

	// ── Master Lokasi ───────────────────────────────
	FindLocations(search string) ([]model.MstScLocation, error)
	FindLocationsPaginated(page, perPage int, search string) ([]model.MstScLocation, int64, error)
	FindLocationByID(id string) (*model.MstScLocation, error)
	CreateLocation(loc *model.MstScLocation) error
	UpdateLocation(loc *model.MstScLocation) error
	DeleteLocation(id string) error

	// ── Transaksi Purchase Request ──────────────────
	FindPurchaseRequestsPaginated(page, perPage int, departmentID, status, search string) ([]model.TrxScPurchaseRequest, int64, error)
	FindPurchaseRequestByID(id string) (*model.TrxScPurchaseRequest, error)
	CreatePurchaseRequest(pr *model.TrxScPurchaseRequest) error
	UpdatePurchaseRequest(pr *model.TrxScPurchaseRequest) error
	UpdatePurchaseRequestStatus(id, status string) error
	DeletePurchaseRequest(id string) error

	// ── Transaksi Purchase Order ────────────────────
	FindPurchaseOrdersPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScPurchaseOrder, int64, error)
	FindPurchaseOrderByID(id string) (*model.TrxScPurchaseOrder, error)
	CreatePurchaseOrder(po *model.TrxScPurchaseOrder) error
	UpdatePurchaseOrder(po *model.TrxScPurchaseOrder) error
	UpdatePurchaseOrderStatus(id, status string) error
	DeletePurchaseOrder(id string) error

	// ── Transaksi Goods Receipt ─────────────────────
	FindGoodsReceiptsPaginated(page, perPage int, poID, search string) ([]model.TrxScGoodsReceipt, int64, error)
	FindGoodsReceiptByID(id string) (*model.TrxScGoodsReceipt, error)
	CreateGoodsReceipt(gr *model.TrxScGoodsReceipt) error
	UpdateGoodsReceipt(gr *model.TrxScGoodsReceipt) error
	DeleteGoodsReceipt(id string) error

	// ── Transaksi Manajemen Aset ────────────────────
	FindAssetsPaginated(page, perPage int, categoryID, locationID, status, search string) ([]model.TrxScAsset, int64, error)
	FindAssetByID(id string) (*model.TrxScAsset, error)
	CreateAsset(asset *model.TrxScAsset) error
	UpdateAsset(asset *model.TrxScAsset) error
	UpdateAssetStatus(id, status string) error
	DeleteAsset(id string) error

	// ── Transaksi Perpanjangan Langganan ────────────
	FindSubscriptionsPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScSubscription, int64, error)
	FindSubscriptionByID(id string) (*model.TrxScSubscription, error)
	CreateSubscription(sub *model.TrxScSubscription) error
	UpdateSubscription(sub *model.TrxScSubscription) error
	UpdateSubscriptionStatus(id, status string) error
	DeleteSubscription(id string) error

	// ── Laporan Supply Chain ────────────────────────
	GetLaporanPO(startDate, endDate, vendorID, status string) ([]model.TrxScPurchaseOrder, error)
	GetLaporanAset(categoryID, locationID, status, departmentID, employeeID string) ([]model.TrxScAsset, error)
	GetLaporanLisensi(status, vendorID string) ([]model.TrxScSubscription, error)
	GetRiwayatVendor(vendorID, startDate, endDate string) ([]model.RiwayatVendorRow, error)
}

type scRepository struct {
	db *gorm.DB
}

func NewScRepository(db *gorm.DB) ScRepository {
	return &scRepository{db: db}
}

// ── Master Vendor ───────────────────────────────────────────────────

func (r *scRepository) FindVendors(search, status string) ([]model.MstScVendor, error) {
	var list []model.MstScVendor
	q := r.db.Model(&model.MstScVendor{})
	if search != "" {
		q = q.Where("vendor_name LIKE ? OR contact_person LIKE ? OR email LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("active_status = ? OR status = ?", status, status)
	}
	err := q.Order("vendor_name ASC").Find(&list).Error
	return list, err
}

func (r *scRepository) FindVendorsPaginated(page, perPage int, search, status string) ([]model.MstScVendor, int64, error) {
	var list []model.MstScVendor
	var total int64
	q := r.db.Model(&model.MstScVendor{})
	if search != "" {
		q = q.Where("vendor_name LIKE ? OR contact_person LIKE ? OR email LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if status != "" {
		q = q.Where("active_status = ? OR status = ?", status, status)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	err := q.Order("vendor_name ASC").Offset(offset).Limit(perPage).Find(&list).Error
	return list, total, err
}

func (r *scRepository) FindVendorByID(id string) (*model.MstScVendor, error) {
	var v model.MstScVendor
	if err := r.db.First(&v, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &v, nil
}

func (r *scRepository) CreateVendor(v *model.MstScVendor) error {
	return r.db.Create(v).Error
}

func (r *scRepository) UpdateVendor(v *model.MstScVendor) error {
	return r.db.Model(&model.MstScVendor{}).Where("id = ?", v.ID).Updates(v).Error
}

func (r *scRepository) DeleteVendor(id string) error {
	return r.db.Delete(&model.MstScVendor{}, "id = ?", id).Error
}

// ── Master Kategori Item ────────────────────────────────────────────

func (r *scRepository) FindItemCategories(search, itemType string) ([]model.MstScItemCategory, error) {
	var list []model.MstScItemCategory
	q := r.db.Model(&model.MstScItemCategory{})
	if search != "" {
		q = q.Where("category_name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if itemType != "" {
		q = q.Where("item_type = ?", itemType)
	}
	err := q.Order("category_name ASC").Find(&list).Error
	return list, err
}

func (r *scRepository) FindItemCategoriesPaginated(page, perPage int, search, itemType string) ([]model.MstScItemCategory, int64, error) {
	var list []model.MstScItemCategory
	var total int64
	q := r.db.Model(&model.MstScItemCategory{})
	if search != "" {
		q = q.Where("category_name LIKE ? OR code LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if itemType != "" {
		q = q.Where("item_type = ?", itemType)
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	err := q.Order("category_name ASC").Offset(offset).Limit(perPage).Find(&list).Error
	return list, total, err
}

func (r *scRepository) FindItemCategoryByID(id string) (*model.MstScItemCategory, error) {
	var cat model.MstScItemCategory
	if err := r.db.First(&cat, "id = ?", id).Error; err != nil {
		return nil, err
	}
	return &cat, nil
}

func (r *scRepository) CreateItemCategory(cat *model.MstScItemCategory) error {
	return r.db.Create(cat).Error
}

func (r *scRepository) UpdateItemCategory(cat *model.MstScItemCategory) error {
	return r.db.Model(&model.MstScItemCategory{}).Where("id = ?", cat.ID).Updates(cat).Error
}

func (r *scRepository) DeleteItemCategory(id string) error {
	return r.db.Delete(&model.MstScItemCategory{}, "id = ?", id).Error
}

// ── Master Katalog Item ─────────────────────────────────────────────

func (r *scRepository) FindItems(categoryID, search string) ([]model.MstScItem, error) {
	var list []model.MstScItem
	q := r.db.Model(&model.MstScItem{})
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if search != "" {
		q = q.Where("item_name LIKE ? OR code LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Order("item_name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	r.populateItemVirtuals(list)
	return list, nil
}

func (r *scRepository) FindItemsPaginated(page, perPage int, categoryID, search string) ([]model.MstScItem, int64, error) {
	var list []model.MstScItem
	var total int64
	q := r.db.Model(&model.MstScItem{})
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if search != "" {
		q = q.Where("item_name LIKE ? OR code LIKE ? OR description LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("item_name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.populateItemVirtuals(list)
	return list, total, nil
}

func (r *scRepository) populateItemVirtuals(list []model.MstScItem) {
	for i := range list {
		if list[i].CategoryID != "" {
			var cat model.MstScItemCategory
			if err := r.db.Select("category_name").First(&cat, "id = ?", list[i].CategoryID).Error; err == nil {
				list[i].CategoryName = cat.CategoryName
			}
		}
		if list[i].DefaultVendorID != nil && *list[i].DefaultVendorID != "" {
			var v model.MstScVendor
			if err := r.db.Select("vendor_name").First(&v, "id = ?", *list[i].DefaultVendorID).Error; err == nil {
				list[i].DefaultVendorName = v.VendorName
			}
		}
	}
}

func (r *scRepository) FindItemByID(id string) (*model.MstScItem, error) {
	var item model.MstScItem
	if err := r.db.First(&item, "id = ?", id).Error; err != nil {
		return nil, err
	}
	slice := []model.MstScItem{item}
	r.populateItemVirtuals(slice)
	return &slice[0], nil
}

func (r *scRepository) CreateItem(item *model.MstScItem) error {
	return r.db.Create(item).Error
}

func (r *scRepository) UpdateItem(item *model.MstScItem) error {
	return r.db.Model(&model.MstScItem{}).Where("id = ?", item.ID).Updates(item).Error
}

func (r *scRepository) DeleteItem(id string) error {
	return r.db.Delete(&model.MstScItem{}, "id = ?", id).Error
}

// ── Master Lokasi ───────────────────────────────────────────────────

func (r *scRepository) FindLocations(search string) ([]model.MstScLocation, error) {
	var list []model.MstScLocation
	q := r.db.Model(&model.MstScLocation{})
	if search != "" {
		q = q.Where("location_name LIKE ? OR address LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Order("location_name ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	r.populateLocationVirtuals(list)
	return list, nil
}

func (r *scRepository) FindLocationsPaginated(page, perPage int, search string) ([]model.MstScLocation, int64, error) {
	var list []model.MstScLocation
	var total int64
	q := r.db.Model(&model.MstScLocation{})
	if search != "" {
		q = q.Where("location_name LIKE ? OR address LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("location_name ASC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	r.populateLocationVirtuals(list)
	return list, total, nil
}

func (r *scRepository) populateLocationVirtuals(list []model.MstScLocation) {
	for i := range list {
		if list[i].PersonInChargeID != nil && *list[i].PersonInChargeID != "" {
			var emp model.MstHrEmployee
			if err := r.db.Select("full_name").First(&emp, "id = ?", *list[i].PersonInChargeID).Error; err == nil {
				list[i].PersonInChargeName = emp.FullName
			}
		}
	}
}

func (r *scRepository) FindLocationByID(id string) (*model.MstScLocation, error) {
	var loc model.MstScLocation
	if err := r.db.First(&loc, "id = ?", id).Error; err != nil {
		return nil, err
	}
	slice := []model.MstScLocation{loc}
	r.populateLocationVirtuals(slice)
	return &slice[0], nil
}

func (r *scRepository) CreateLocation(loc *model.MstScLocation) error {
	return r.db.Create(loc).Error
}

func (r *scRepository) UpdateLocation(loc *model.MstScLocation) error {
	return r.db.Model(&model.MstScLocation{}).Where("id = ?", loc.ID).Updates(loc).Error
}

func (r *scRepository) DeleteLocation(id string) error {
	return r.db.Delete(&model.MstScLocation{}, "id = ?", id).Error
}

// ── Transaksi Purchase Request ──────────────────────────────────────

func (r *scRepository) FindPurchaseRequestsPaginated(page, perPage int, departmentID, status, search string) ([]model.TrxScPurchaseRequest, int64, error) {
	var list []model.TrxScPurchaseRequest
	var total int64
	q := r.db.Model(&model.TrxScPurchaseRequest{})
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if status != "" {
		q = q.Where("approval_status = ? OR status = ?", status, status)
	}
	if search != "" {
		q = q.Where("pr_number LIKE ? OR notes LIKE ? OR requester_name LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Preload("Items").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].DepartmentID != "" {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", list[i].DepartmentID).Error; err == nil {
				list[i].DepartmentName = dept.DepartmentName
			}
		}
		list[i].RequestedByName = list[i].RequesterName
		for j := range list[i].Items {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].Items[j].ItemID).Error; err == nil {
				list[i].Items[j].ItemName = item.ItemName
			}
		}
	}
	return list, total, nil
}

func (r *scRepository) FindPurchaseRequestByID(id string) (*model.TrxScPurchaseRequest, error) {
	var pr model.TrxScPurchaseRequest
	if err := r.db.Preload("Items").First(&pr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	pr.RequestedByName = pr.RequesterName
	if pr.DepartmentID != "" {
		var dept model.MstHrDepartment
		if err := r.db.Select("department_name").First(&dept, "id = ?", pr.DepartmentID).Error; err == nil {
			pr.DepartmentName = dept.DepartmentName
		}
	}
	for j := range pr.Items {
		var item model.MstScItem
		if err := r.db.Select("item_name").First(&item, "id = ?", pr.Items[j].ItemID).Error; err == nil {
			pr.Items[j].ItemName = item.ItemName
		}
	}
	return &pr, nil
}

func (r *scRepository) CreatePurchaseRequest(pr *model.TrxScPurchaseRequest) error {
	return r.db.Create(pr).Error
}

func (r *scRepository) UpdatePurchaseRequest(pr *model.TrxScPurchaseRequest) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TrxScPurchaseRequestItem{}, "pr_id = ?", pr.ID).Error; err != nil {
			return err
		}
		if len(pr.Items) > 0 {
			if err := tx.Create(&pr.Items).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.TrxScPurchaseRequest{}).Where("id = ?", pr.ID).Omit("Items").Updates(pr).Error
	})
}

func (r *scRepository) UpdatePurchaseRequestStatus(id, status string) error {
	return r.db.Model(&model.TrxScPurchaseRequest{}).Where("id = ?", id).Updates(map[string]interface{}{
		"approval_status": status,
		"status":          status,
	}).Error
}

func (r *scRepository) DeletePurchaseRequest(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Delete(&model.TrxScPurchaseRequestItem{}, "pr_id = ?", id).Error
		return tx.Delete(&model.TrxScPurchaseRequest{}, "id = ?", id).Error
	})
}

// ── Transaksi Purchase Order ────────────────────────────────────────

func (r *scRepository) FindPurchaseOrdersPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScPurchaseOrder, int64, error) {
	var list []model.TrxScPurchaseOrder
	var total int64
	q := r.db.Model(&model.TrxScPurchaseOrder{})
	if vendorID != "" {
		q = q.Where("vendor_id = ?", vendorID)
	}
	if status != "" {
		q = q.Where("status = ? OR po_status = ?", status, status)
	}
	if search != "" {
		q = q.Where("po_number LIKE ?", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Preload("Items").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].VendorID != "" {
			var v model.MstScVendor
			if err := r.db.Select("vendor_name").First(&v, "id = ?", list[i].VendorID).Error; err == nil {
				list[i].VendorName = v.VendorName
			}
		}
		if list[i].PrID != nil && *list[i].PrID != "" {
			var pr model.TrxScPurchaseRequest
			if err := r.db.Select("pr_number").First(&pr, "id = ?", *list[i].PrID).Error; err == nil {
				list[i].PrNumber = pr.PrNumber
			}
		}
		for j := range list[i].Items {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].Items[j].ItemID).Error; err == nil {
				list[i].Items[j].ItemName = item.ItemName
			}
		}
	}
	return list, total, nil
}

func (r *scRepository) FindPurchaseOrderByID(id string) (*model.TrxScPurchaseOrder, error) {
	var po model.TrxScPurchaseOrder
	if err := r.db.Preload("Items").First(&po, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if po.VendorID != "" {
		var v model.MstScVendor
		if err := r.db.Select("vendor_name").First(&v, "id = ?", po.VendorID).Error; err == nil {
			po.VendorName = v.VendorName
		}
	}
	if po.PrID != nil && *po.PrID != "" {
		var pr model.TrxScPurchaseRequest
		if err := r.db.Select("pr_number").First(&pr, "id = ?", *po.PrID).Error; err == nil {
			po.PrNumber = pr.PrNumber
		}
	}
	for j := range po.Items {
		var item model.MstScItem
		if err := r.db.Select("item_name").First(&item, "id = ?", po.Items[j].ItemID).Error; err == nil {
			po.Items[j].ItemName = item.ItemName
		}
	}
	return &po, nil
}

func (r *scRepository) CreatePurchaseOrder(po *model.TrxScPurchaseOrder) error {
	return r.db.Create(po).Error
}

func (r *scRepository) UpdatePurchaseOrder(po *model.TrxScPurchaseOrder) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TrxScPurchaseOrderItem{}, "po_id = ?", po.ID).Error; err != nil {
			return err
		}
		if len(po.Items) > 0 {
			if err := tx.Create(&po.Items).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.TrxScPurchaseOrder{}).Where("id = ?", po.ID).Omit("Items").Updates(po).Error
	})
}

func (r *scRepository) UpdatePurchaseOrderStatus(id, status string) error {
	return r.db.Model(&model.TrxScPurchaseOrder{}).Where("id = ?", id).Updates(map[string]interface{}{
		"status":    status,
		"po_status": status,
	}).Error
}

func (r *scRepository) DeletePurchaseOrder(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Delete(&model.TrxScPurchaseOrderItem{}, "po_id = ?", id).Error
		return tx.Delete(&model.TrxScPurchaseOrder{}, "id = ?", id).Error
	})
}

// ── Transaksi Goods Receipt ─────────────────────────────────────────

func (r *scRepository) FindGoodsReceiptsPaginated(page, perPage int, poID, search string) ([]model.TrxScGoodsReceipt, int64, error) {
	var list []model.TrxScGoodsReceipt
	var total int64
	q := r.db.Model(&model.TrxScGoodsReceipt{})
	if poID != "" {
		q = q.Where("po_id = ?", poID)
	}
	if search != "" {
		q = q.Where("receipt_number LIKE ? OR delivery_note_number LIKE ? OR notes LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Preload("Items").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		list[i].GrNumber = list[i].ReceiptNumber
		if list[i].PoID != "" {
			var po model.TrxScPurchaseOrder
			if err := r.db.Select("po_number, vendor_id").First(&po, "id = ?", list[i].PoID).Error; err == nil {
				list[i].PoNumber = po.PoNumber
				if po.VendorID != "" {
					var v model.MstScVendor
					if err := r.db.Select("vendor_name").First(&v, "id = ?", po.VendorID).Error; err == nil {
						list[i].VendorName = v.VendorName
					}
				}
			}
		}
		if list[i].ReceivedBy != "" {
			var emp model.MstHrEmployee
			if err := r.db.Select("full_name").First(&emp, "id = ?", list[i].ReceivedBy).Error; err == nil {
				list[i].ReceivedByName = emp.FullName
			} else {
				list[i].ReceivedByName = list[i].ReceivedBy
			}
		}
		for j := range list[i].Items {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].Items[j].ItemID).Error; err == nil {
				list[i].Items[j].ItemName = item.ItemName
			}
		}
	}
	return list, total, nil
}

func (r *scRepository) FindGoodsReceiptByID(id string) (*model.TrxScGoodsReceipt, error) {
	var gr model.TrxScGoodsReceipt
	if err := r.db.Preload("Items").First(&gr, "id = ?", id).Error; err != nil {
		return nil, err
	}
	gr.GrNumber = gr.ReceiptNumber
	if gr.PoID != "" {
		var po model.TrxScPurchaseOrder
		if err := r.db.Select("po_number").First(&po, "id = ?", gr.PoID).Error; err == nil {
			gr.PoNumber = po.PoNumber
		}
	}
	if gr.ReceivedBy != "" {
		var emp model.MstHrEmployee
		if err := r.db.Select("full_name").First(&emp, "id = ?", gr.ReceivedBy).Error; err == nil {
			gr.ReceivedByName = emp.FullName
		} else {
			gr.ReceivedByName = gr.ReceivedBy
		}
	}
	for j := range gr.Items {
		var item model.MstScItem
		if err := r.db.Select("item_name").First(&item, "id = ?", gr.Items[j].ItemID).Error; err == nil {
			gr.Items[j].ItemName = item.ItemName
		}
	}
	return &gr, nil
}

func (r *scRepository) CreateGoodsReceipt(gr *model.TrxScGoodsReceipt) error {
	return r.db.Create(gr).Error
}

func (r *scRepository) UpdateGoodsReceipt(gr *model.TrxScGoodsReceipt) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Delete(&model.TrxScGoodsReceiptItem{}, "gr_id = ?", gr.ID).Error; err != nil {
			return err
		}
		if len(gr.Items) > 0 {
			if err := tx.Create(&gr.Items).Error; err != nil {
				return err
			}
		}
		return tx.Model(&model.TrxScGoodsReceipt{}).Where("id = ?", gr.ID).Omit("Items").Updates(gr).Error
	})
}

func (r *scRepository) DeleteGoodsReceipt(id string) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		_ = tx.Delete(&model.TrxScGoodsReceiptItem{}, "gr_id = ?", id).Error
		return tx.Delete(&model.TrxScGoodsReceipt{}, "id = ?", id).Error
	})
}

// ── Transaksi Manajemen Aset ────────────────────────────────────────

func (r *scRepository) FindAssetsPaginated(page, perPage int, categoryID, locationID, status, search string) ([]model.TrxScAsset, int64, error) {
	var list []model.TrxScAsset
	var total int64
	q := r.db.Model(&model.TrxScAsset{})
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if locationID != "" {
		q = q.Where("location_id = ?", locationID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("asset_code LIKE ? OR name LIKE ? OR serial_or_license_key LIKE ?", "%"+search+"%", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].ItemID != "" {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].ItemID).Error; err == nil {
				list[i].ItemName = item.ItemName
				if list[i].Name == "" {
					list[i].Name = item.ItemName
				}
			}
		}
		if list[i].CategoryID != "" {
			var cat model.MstScItemCategory
			if err := r.db.Select("category_name").First(&cat, "id = ?", list[i].CategoryID).Error; err == nil {
				list[i].CategoryName = cat.CategoryName
			}
		}
		if list[i].AssignedToEmployeeID != nil && *list[i].AssignedToEmployeeID != "" {
			var emp model.MstHrEmployee
			if err := r.db.Select("full_name, department_id").First(&emp, "id = ?", *list[i].AssignedToEmployeeID).Error; err == nil {
				list[i].AssignedToEmployeeName = emp.FullName
				list[i].AssignedToName = emp.FullName
				if (list[i].DepartmentID == nil || *list[i].DepartmentID == "") && emp.DepartmentID != nil && *emp.DepartmentID != "" {
					list[i].DepartmentID = emp.DepartmentID
				}
			}
		}
		if list[i].DepartmentID != nil && *list[i].DepartmentID != "" {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", *list[i].DepartmentID).Error; err == nil {
				list[i].DepartmentName = dept.DepartmentName
			}
		}
		if list[i].LocationID != nil && *list[i].LocationID != "" {
			var loc model.MstScLocation
			if err := r.db.Select("location_name").First(&loc, "id = ?", *list[i].LocationID).Error; err == nil {
				list[i].LocationName = loc.LocationName
			}
		}
	}
	return list, total, nil
}

func (r *scRepository) FindAssetByID(id string) (*model.TrxScAsset, error) {
	var asset model.TrxScAsset
	if err := r.db.First(&asset, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if asset.ItemID != "" {
		var item model.MstScItem
		if err := r.db.Select("item_name").First(&item, "id = ?", asset.ItemID).Error; err == nil {
			asset.ItemName = item.ItemName
			if asset.Name == "" {
				asset.Name = item.ItemName
			}
		}
	}
	return &asset, nil
}

func (r *scRepository) CreateAsset(asset *model.TrxScAsset) error {
	return r.db.Create(asset).Error
}

func (r *scRepository) UpdateAsset(asset *model.TrxScAsset) error {
	return r.db.Model(&model.TrxScAsset{}).Where("id = ?", asset.ID).Updates(asset).Error
}

func (r *scRepository) UpdateAssetStatus(id, status string) error {
	return r.db.Model(&model.TrxScAsset{}).Where("id = ?", id).Update("status", status).Error
}

func (r *scRepository) DeleteAsset(id string) error {
	return r.db.Delete(&model.TrxScAsset{}, "id = ?", id).Error
}

// ── Transaksi Perpanjangan Langganan ────────────────────────────────

func (r *scRepository) FindSubscriptionsPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScSubscription, int64, error) {
	var list []model.TrxScSubscription
	var total int64
	q := r.db.Model(&model.TrxScSubscription{})
	if vendorID != "" {
		q = q.Where("vendor_id = ?", vendorID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if search != "" {
		q = q.Where("license_code LIKE ? OR name LIKE ?", "%"+search+"%", "%"+search+"%")
	}
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	offset := (page - 1) * perPage
	if err := q.Order("created_at DESC").Offset(offset).Limit(perPage).Find(&list).Error; err != nil {
		return nil, 0, err
	}
	for i := range list {
		if list[i].ItemID != "" {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].ItemID).Error; err == nil {
				list[i].ItemName = item.ItemName
				if list[i].Name == "" {
					list[i].Name = item.ItemName
				}
			}
		}
		if list[i].VendorID != "" {
			var v model.MstScVendor
			if err := r.db.Select("vendor_name").First(&v, "id = ?", list[i].VendorID).Error; err == nil {
				list[i].VendorName = v.VendorName
			}
		}
	}
	return list, total, nil
}

func (r *scRepository) FindSubscriptionByID(id string) (*model.TrxScSubscription, error) {
	var sub model.TrxScSubscription
	if err := r.db.First(&sub, "id = ?", id).Error; err != nil {
		return nil, err
	}
	if sub.ItemID != "" {
		var item model.MstScItem
		if err := r.db.Select("item_name").First(&item, "id = ?", sub.ItemID).Error; err == nil {
			sub.ItemName = item.ItemName
			if sub.Name == "" {
				sub.Name = item.ItemName
			}
		}
	}
	return &sub, nil
}

func (r *scRepository) CreateSubscription(sub *model.TrxScSubscription) error {
	return r.db.Create(sub).Error
}

func (r *scRepository) UpdateSubscription(sub *model.TrxScSubscription) error {
	return r.db.Model(&model.TrxScSubscription{}).Where("id = ?", sub.ID).Updates(sub).Error
}

func (r *scRepository) UpdateSubscriptionStatus(id, status string) error {
	return r.db.Model(&model.TrxScSubscription{}).Where("id = ?", id).Update("status", status).Error
}

func (r *scRepository) DeleteSubscription(id string) error {
	return r.db.Delete(&model.TrxScSubscription{}, "id = ?", id).Error
}

// ── Laporan Supply Chain ────────────────────────────────────────────

func (r *scRepository) GetLaporanPO(startDate, endDate, vendorID, status string) ([]model.TrxScPurchaseOrder, error) {
	var list []model.TrxScPurchaseOrder
	q := r.db.Model(&model.TrxScPurchaseOrder{})
	if startDate != "" {
		q = q.Where("po_date >= ?", startDate)
	}
	if endDate != "" {
		q = q.Where("po_date <= ?", endDate)
	}
	if vendorID != "" {
		q = q.Where("vendor_id = ?", vendorID)
	}
	if status != "" {
		q = q.Where("status = ? OR po_status = ?", status, status)
	}
	if err := q.Order("po_date DESC").Find(&list).Error; err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].VendorID != "" {
			var v model.MstScVendor
			if err := r.db.Select("vendor_name").First(&v, "id = ?", list[i].VendorID).Error; err == nil {
				list[i].VendorName = v.VendorName
			}
		}
	}
	return list, nil
}

func (r *scRepository) GetLaporanAset(categoryID, locationID, status, departmentID, employeeID string) ([]model.TrxScAsset, error) {
	var list []model.TrxScAsset
	q := r.db.Model(&model.TrxScAsset{})
	if categoryID != "" {
		q = q.Where("category_id = ?", categoryID)
	}
	if locationID != "" {
		q = q.Where("location_id = ?", locationID)
	}
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if departmentID != "" {
		q = q.Where("department_id = ?", departmentID)
	}
	if employeeID != "" {
		q = q.Where("assigned_to_employee_id = ?", employeeID)
	}
	if err := q.Order("asset_code ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ItemID != "" {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].ItemID).Error; err == nil {
				list[i].ItemName = item.ItemName
				if list[i].Name == "" {
					list[i].Name = item.ItemName
				}
			}
		}
		if list[i].AssignedToEmployeeID != nil && *list[i].AssignedToEmployeeID != "" {
			var emp model.MstHrEmployee
			if err := r.db.Select("full_name, department_id").First(&emp, "id = ?", *list[i].AssignedToEmployeeID).Error; err == nil {
				list[i].AssignedToEmployeeName = emp.FullName
				list[i].AssignedToName = emp.FullName
				if (list[i].DepartmentID == nil || *list[i].DepartmentID == "") && emp.DepartmentID != nil && *emp.DepartmentID != "" {
					list[i].DepartmentID = emp.DepartmentID
				}
			}
		}
		if list[i].DepartmentID != nil && *list[i].DepartmentID != "" {
			var dept model.MstHrDepartment
			if err := r.db.Select("department_name").First(&dept, "id = ?", *list[i].DepartmentID).Error; err == nil {
				list[i].DepartmentName = dept.DepartmentName
			}
		}
		if list[i].LocationID != nil && *list[i].LocationID != "" {
			var loc model.MstScLocation
			if err := r.db.Select("location_name").First(&loc, "id = ?", *list[i].LocationID).Error; err == nil {
				list[i].LocationName = loc.LocationName
			}
		}
	}
	return list, nil
}

func (r *scRepository) GetLaporanLisensi(status, vendorID string) ([]model.TrxScSubscription, error) {
	var list []model.TrxScSubscription
	q := r.db.Model(&model.TrxScSubscription{})
	if status != "" {
		q = q.Where("status = ?", status)
	}
	if vendorID != "" {
		q = q.Where("vendor_id = ?", vendorID)
	}
	if err := q.Order("end_date ASC").Find(&list).Error; err != nil {
		return nil, err
	}
	for i := range list {
		if list[i].ItemID != "" {
			var item model.MstScItem
			if err := r.db.Select("item_name").First(&item, "id = ?", list[i].ItemID).Error; err == nil {
				list[i].ItemName = item.ItemName
				if list[i].Name == "" {
					list[i].Name = item.ItemName
				}
			}
		}
		if list[i].VendorID != "" {
			var v model.MstScVendor
			if err := r.db.Select("vendor_name").First(&v, "id = ?", list[i].VendorID).Error; err == nil {
				list[i].VendorName = v.VendorName
			}
		}
	}
	return list, nil
}

func (r *scRepository) GetRiwayatVendor(vendorID, startDate, endDate string) ([]model.RiwayatVendorRow, error) {
	var pos []model.TrxScPurchaseOrder
	q := r.db.Model(&model.TrxScPurchaseOrder{})
	if vendorID != "" {
		q = q.Where("vendor_id = ?", vendorID)
	}
	if startDate != "" {
		q = q.Where("po_date >= ?", startDate)
	}
	if endDate != "" {
		q = q.Where("po_date <= ?", endDate)
	}
	if err := q.Order("po_date DESC").Find(&pos).Error; err != nil {
		return nil, err
	}

	var rows []model.RiwayatVendorRow
	for _, po := range pos {
		vendorName := ""
		if po.VendorID != "" {
			var v model.MstScVendor
			if err := r.db.Select("vendor_name").First(&v, "id = ?", po.VendorID).Error; err == nil {
				vendorName = v.VendorName
			}
		}

		st := strings.ToUpper(po.Status)
		if st == "" {
			st = strings.ToUpper(po.PoStatus)
		}
		if st == "" {
			st = "DRAFT"
		}

		paymentStatus := "UNPAID"
		if st == "COMPLETED" || st == "PAID" {
			paymentStatus = "PAID"
		} else if st == "PARTIALLY_RECEIVED" || st == "PARTIAL" {
			paymentStatus = "PARTIAL"
		}

		completedCount := 0
		if st == "COMPLETED" {
			completedCount = 1
		}

		tagihanNum := "-"
		// 1. Cek dari Vendor Bills (AP)
		var bill model.TrxVendorBill
		if err := r.db.Where("vendor_id = ?", po.VendorID).First(&bill).Error; err == nil && bill.BillNumber != "" {
			tagihanNum = bill.BillNumber
		}
		// 2. Cek dari Goods Receipt (Surat Jalan / No. Dokumen Pengiriman)
		if tagihanNum == "-" {
			var gr model.TrxScGoodsReceipt
			if err := r.db.Where("po_id = ?", po.ID).First(&gr).Error; err == nil {
				if gr.DeliveryNoteNumber != "" {
					tagihanNum = gr.DeliveryNoteNumber
				} else if gr.ReceiptNumber != "" {
					tagihanNum = gr.ReceiptNumber
				}
			}
		}
		// 3. Fallback nomor invoice vendor jika sudah selesai/lunas
		if tagihanNum == "-" && (paymentStatus == "PAID" || st == "COMPLETED") {
			tagihanNum = "INV-" + strings.Replace(po.PoNumber, "PO-", "", 1)
		}

		rows = append(rows, model.RiwayatVendorRow{
			ID:            po.ID,
			VendorID:      po.VendorID,
			VendorName:    vendorName,
			PoNumber:      po.PoNumber,
			TagihanNumber: tagihanNum,
			Date:          po.PoDate,
			TotalAmount:   po.TotalAmount,
			PaymentStatus: paymentStatus,
			Status:        st,
			TotalPO:       1,
			CompletedPO:   completedCount,
			LastOrderDate: po.PoDate,
		})
	}
	return rows, nil
}

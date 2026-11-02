package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"

	"github.com/google/uuid"
)

// ── DTOs Master Supply Chain ────────────────────────────────────────

type VendorDTO struct {
	Code           string `json:"code"`
	VendorName     string `json:"vendor_name"`
	Name           string `json:"name"`
	VendorCategory string `json:"vendor_category"`
	Category       string `json:"category"`
	ContactPerson  string `json:"contact_person"`
	PicName        string `json:"pic_name"`
	Phone          string `json:"phone"`
	PhoneNumber    string `json:"phone_number"`
	PicPhone       string `json:"pic_phone"`
	Email          string `json:"email"`
	Address        string `json:"address"`
	BankName       string `json:"bank_name"`
	BankAccount    string `json:"bank_account"`
	ActiveStatus   string `json:"active_status"`
	Status         string `json:"status"`
}

type ItemCategoryDTO struct {
	Code         string `json:"code"`
	CategoryName string `json:"category_name"`
	Name         string `json:"name"`
	ItemType     string `json:"item_type"`
	Type         string `json:"type"`
	Description  string `json:"description"`
}

type ItemDTO struct {
	CategoryID      string  `json:"category_id"`
	Code            string  `json:"code"`
	ItemName        string  `json:"item_name"`
	Name            string  `json:"name"`
	Unit            string  `json:"unit"`
	EstimatedPrice  float64 `json:"estimated_price"`
	Price           float64 `json:"price"`
	DefaultVendorID *string `json:"default_vendor_id"`
	MinStock        int     `json:"min_stock"`
	CurrentStock    int     `json:"current_stock"`
	Description     string  `json:"description"`
	Status          string  `json:"status"`
}

type LocationDTO struct {
	Code               string  `json:"code"`
	LocationName       string  `json:"location_name"`
	Name               string  `json:"name"`
	Type               string  `json:"type"`
	Address            string  `json:"address"`
	PersonInChargeID   *string `json:"person_in_charge_id"`
	Description        string  `json:"description"`
}

// ── DTOs Transaksi Supply Chain ─────────────────────────────────────

type PurchaseRequestItemDTO struct {
	ItemID         string  `json:"item_id"`
	Quantity       int     `json:"quantity"`
	Unit           string  `json:"unit"`
	EstimatedPrice float64 `json:"estimated_price"`
	Notes          string  `json:"notes"`
}

type PurchaseRequestDTO struct {
	PrNumber       string                   `json:"pr_number"`
	DepartmentID   string                   `json:"department_id"`
	RequesterName  string                   `json:"requester_name"`
	RequestDate    string                   `json:"request_date"`
	Priority       string                   `json:"priority"`
	Notes          string                   `json:"notes"`
	ApprovalStatus string                   `json:"approval_status"`
	Status         string                   `json:"status"`
	Items          []PurchaseRequestItemDTO `json:"items"`
}

type PurchaseOrderItemDTO struct {
	ItemID     string  `json:"item_id"`
	Quantity   int     `json:"quantity"`
	Unit       string  `json:"unit"`
	UnitPrice  float64 `json:"unit_price"`
	TotalPrice float64 `json:"total_price"`
	TaxRate    float64 `json:"tax_rate"`
}

type PurchaseOrderDTO struct {
	PoNumber             string                 `json:"po_number"`
	PrID                 *string                `json:"pr_id"`
	VendorID             string                 `json:"vendor_id"`
	PoDate               string                 `json:"po_date"`
	OrderDate            string                 `json:"order_date"`
	ExpectedDeliveryDate string                 `json:"expected_delivery_date"`
	Subtotal             float64                `json:"subtotal"`
	TaxRate              float64                `json:"tax_rate"`
	TaxAmount            float64                `json:"tax_amount"`
	TotalAmount          float64                `json:"total_amount"`
	Status               string                 `json:"status"`
	PoStatus             string                 `json:"po_status"`
	Notes                string                 `json:"notes"`
	Items                []PurchaseOrderItemDTO `json:"items"`
}

type GoodsReceiptItemDTO struct {
	ItemID           string `json:"item_id"`
	OrderedQuantity  int    `json:"ordered_quantity"`
	ReceivedQuantity int    `json:"received_quantity"`
	Condition        string `json:"condition"`
	Notes            string `json:"notes"`
}

type GoodsReceiptDTO struct {
	ReceiptNumber      string                `json:"receipt_number"`
	GrNumber           string                `json:"gr_number"`
	PoID               string                `json:"po_id"`
	PurchaseOrderID    string                `json:"purchase_order_id"`
	ReceivedBy         string                `json:"received_by"`
	ReceiptDate        string                `json:"receipt_date"`
	DeliveryNoteNumber string                `json:"delivery_note_number"`
	DeliveryNoteNo     string                `json:"delivery_note_no"`
	AttachmentURL      string                `json:"attachment_url"`
	Notes              string                `json:"notes"`
	Status             string                `json:"status"`
	Items              []GoodsReceiptItemDTO `json:"items"`
}

type AssetDTO struct {
	AssetCode               string  `json:"asset_code"`
	Name                    string  `json:"name"`
	AssetName               string  `json:"asset_name"`
	ItemID                  string  `json:"item_id"`
	CategoryID              string  `json:"category_id"`
	VendorID                string  `json:"vendor_id"`
	SerialOrLicenseKey      string  `json:"serial_or_license_key"`
	SerialNumber            string  `json:"serial_number"`
	AcquisitionDate         string  `json:"acquisition_date"`
	PurchaseDate            string  `json:"purchase_date"`
	AcquisitionValue        float64 `json:"acquisition_value"`
	PurchasePrice           float64 `json:"purchase_price"`
	PurchaseCost            float64 `json:"purchase_cost"`
	UsefulLifeMonths        int     `json:"useful_life_months"`
	SalvageValue            float64 `json:"salvage_value"`
	DepreciationMethod      string  `json:"depreciation_method"`
	WarrantyOrLicenseExpiry string  `json:"warranty_or_license_expiry"`
	AssignedToEmployeeID    *string `json:"assigned_to_employee_id"`
	AssignedTo              *string `json:"assigned_to"`
	DepartmentID            *string `json:"department_id"`
	LocationID              *string `json:"location_id"`
	Status                  string  `json:"status"`
	Condition               string  `json:"condition"`
	Notes                   string  `json:"notes"`
}

type SubscriptionDTO struct {
	LicenseCode  string  `json:"license_code"`
	Name         string  `json:"name"`
	AppName      string  `json:"app_name"`
	ItemID       string  `json:"item_id"`
	VendorID     string  `json:"vendor_id"`
	BillingCycle string  `json:"billing_cycle"`
	Amount       float64 `json:"amount"`
	Cost         float64 `json:"cost"`
	StartDate    string  `json:"start_date"`
	EndDate      string  `json:"end_date"`
	ExpiryDate   string  `json:"expiry_date"`
	AutoRenewal  bool    `json:"auto_renewal"`
	AutoRenew    bool    `json:"auto_renew"`
	Status       string  `json:"status"`
	Notes        string  `json:"notes"`
	AssignedTo   string  `json:"assigned_to"`
	LicenseCount int     `json:"license_count"`
}

// ── ScService Interface ─────────────────────────────────────────────

type ScService interface {
	// Master Vendor
	GetAllVendors(search, status string) ([]model.MstScVendor, error)
	GetVendorsPaginated(page, perPage int, search, status string) ([]model.MstScVendor, int64, error)
	GetVendorByID(id string) (*model.MstScVendor, error)
	CreateVendor(dto *VendorDTO) (*model.MstScVendor, error)
	UpdateVendor(id string, dto *VendorDTO) (*model.MstScVendor, error)
	DeleteVendor(id string) error

	// Master Kategori Item
	GetAllItemCategories(search, itemType string) ([]model.MstScItemCategory, error)
	GetItemCategoriesPaginated(page, perPage int, search, itemType string) ([]model.MstScItemCategory, int64, error)
	GetItemCategoryByID(id string) (*model.MstScItemCategory, error)
	CreateItemCategory(dto *ItemCategoryDTO) (*model.MstScItemCategory, error)
	UpdateItemCategory(id string, dto *ItemCategoryDTO) (*model.MstScItemCategory, error)
	DeleteItemCategory(id string) error

	// Master Katalog Item
	GetAllItems(categoryID, search string) ([]model.MstScItem, error)
	GetItemsPaginated(page, perPage int, categoryID, search string) ([]model.MstScItem, int64, error)
	GetItemByID(id string) (*model.MstScItem, error)
	CreateItem(dto *ItemDTO) (*model.MstScItem, error)
	UpdateItem(id string, dto *ItemDTO) (*model.MstScItem, error)
	DeleteItem(id string) error

	// Master Lokasi
	GetAllLocations(search string) ([]model.MstScLocation, error)
	GetLocationsPaginated(page, perPage int, search string) ([]model.MstScLocation, int64, error)
	GetLocationByID(id string) (*model.MstScLocation, error)
	CreateLocation(dto *LocationDTO) (*model.MstScLocation, error)
	UpdateLocation(id string, dto *LocationDTO) (*model.MstScLocation, error)
	DeleteLocation(id string) error

	// Transaksi Purchase Request
	GetPurchaseRequestsPaginated(page, perPage int, departmentID, status, search string) ([]model.TrxScPurchaseRequest, int64, error)
	GetPurchaseRequestByID(id string) (*model.TrxScPurchaseRequest, error)
	CreatePurchaseRequest(dto *PurchaseRequestDTO) (*model.TrxScPurchaseRequest, error)
	UpdatePurchaseRequest(id string, dto *PurchaseRequestDTO) (*model.TrxScPurchaseRequest, error)
	UpdatePurchaseRequestStatus(id, status string) error
	DeletePurchaseRequest(id string) error

	// Transaksi Purchase Order
	GetPurchaseOrdersPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScPurchaseOrder, int64, error)
	GetPurchaseOrderByID(id string) (*model.TrxScPurchaseOrder, error)
	CreatePurchaseOrder(dto *PurchaseOrderDTO) (*model.TrxScPurchaseOrder, error)
	UpdatePurchaseOrder(id string, dto *PurchaseOrderDTO) (*model.TrxScPurchaseOrder, error)
	UpdatePurchaseOrderStatus(id, status string) error
	DeletePurchaseOrder(id string) error

	// Transaksi Goods Receipt
	GetGoodsReceiptsPaginated(page, perPage int, poID, search string) ([]model.TrxScGoodsReceipt, int64, error)
	GetGoodsReceiptByID(id string) (*model.TrxScGoodsReceipt, error)
	CreateGoodsReceipt(dto *GoodsReceiptDTO) (*model.TrxScGoodsReceipt, error)
	UpdateGoodsReceipt(id string, dto *GoodsReceiptDTO) (*model.TrxScGoodsReceipt, error)
	DeleteGoodsReceipt(id string) error

	// Transaksi Manajemen Aset
	GetAssetsPaginated(page, perPage int, categoryID, locationID, status, search string) ([]model.TrxScAsset, int64, error)
	GetAssetByID(id string) (*model.TrxScAsset, error)
	CreateAsset(dto *AssetDTO) (*model.TrxScAsset, error)
	UpdateAsset(id string, dto *AssetDTO) (*model.TrxScAsset, error)
	UpdateAssetStatus(id, status string) error
	DeleteAsset(id string) error

	// Transaksi Perpanjangan Langganan
	GetSubscriptionsPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScSubscription, int64, error)
	GetSubscriptionByID(id string) (*model.TrxScSubscription, error)
	CreateSubscription(dto *SubscriptionDTO) (*model.TrxScSubscription, error)
	UpdateSubscription(id string, dto *SubscriptionDTO) (*model.TrxScSubscription, error)
	UpdateSubscriptionStatus(id, status string) error
	DeleteSubscription(id string) error

	// Laporan SC
	GetLaporanPO(startDate, endDate, vendorID, status string) ([]model.TrxScPurchaseOrder, error)
	GetLaporanAset(categoryID, locationID, status, departmentID, employeeID string) ([]model.TrxScAsset, error)
	GetLaporanLisensi(status, vendorID string) ([]model.TrxScSubscription, error)
	GetRiwayatVendor(vendorID, startDate, endDate string) ([]model.RiwayatVendorRow, error)
}

type scService struct {
	repo repository.ScRepository
}

func NewScService(repo repository.ScRepository) ScService {
	return &scService{repo: repo}
}

// ── Master Vendor ───────────────────────────────────────────────────

func (s *scService) GetAllVendors(search, status string) ([]model.MstScVendor, error) {
	return s.repo.FindVendors(search, status)
}

func (s *scService) GetVendorsPaginated(page, perPage int, search, status string) ([]model.MstScVendor, int64, error) {
	return s.repo.FindVendorsPaginated(page, perPage, search, status)
}

func (s *scService) GetVendorByID(id string) (*model.MstScVendor, error) {
	return s.repo.FindVendorByID(id)
}

func (s *scService) CreateVendor(dto *VendorDTO) (*model.MstScVendor, error) {
	vName := strings.TrimSpace(dto.VendorName)
	if vName == "" {
		vName = strings.TrimSpace(dto.Name)
	}
	if vName == "" {
		return nil, errors.New("nama vendor wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("VND-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	vCat := dto.VendorCategory
	if vCat == "" {
		vCat = dto.Category
	}
	if vCat == "" {
		vCat = "OTHER"
	}

	contact := dto.ContactPerson
	if contact == "" {
		contact = dto.PicName
	}

	phone := dto.Phone
	if phone == "" {
		phone = dto.PhoneNumber
	}

	actStatus := dto.ActiveStatus
	if actStatus == "" {
		actStatus = dto.Status
	}
	if actStatus == "" {
		actStatus = "ACTIVE"
	}

	v := &model.MstScVendor{
		ID:             uuid.NewString(),
		Code:           code,
		VendorName:     vName,
		VendorCategory: vCat,
		ContactPerson:  contact,
		Phone:          phone,
		PicPhone:       dto.PicPhone,
		Email:          dto.Email,
		Address:        dto.Address,
		BankName:       dto.BankName,
		BankAccount:    dto.BankAccount,
		ActiveStatus:   actStatus,
		Status:         actStatus,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreateVendor(v); err != nil {
		return nil, err
	}
	return v, nil
}

func (s *scService) UpdateVendor(id string, dto *VendorDTO) (*model.MstScVendor, error) {
	existing, err := s.repo.FindVendorByID(id)
	if err != nil {
		return nil, errors.New("vendor tidak ditemukan")
	}

	vName := strings.TrimSpace(dto.VendorName)
	if vName == "" {
		vName = strings.TrimSpace(dto.Name)
	}
	if vName != "" {
		existing.VendorName = vName
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	vCat := dto.VendorCategory
	if vCat == "" {
		vCat = dto.Category
	}
	if vCat != "" {
		existing.VendorCategory = vCat
	}
	contact := dto.ContactPerson
	if contact == "" {
		contact = dto.PicName
	}
	if contact != "" {
		existing.ContactPerson = contact
	}
	phone := dto.Phone
	if phone == "" {
		phone = dto.PhoneNumber
	}
	if phone != "" {
		existing.Phone = phone
	}
	if dto.PicPhone != "" {
		existing.PicPhone = dto.PicPhone
	}
	if dto.Email != "" {
		existing.Email = dto.Email
	}
	if dto.Address != "" {
		existing.Address = dto.Address
	}
	if dto.BankName != "" {
		existing.BankName = dto.BankName
	}
	if dto.BankAccount != "" {
		existing.BankAccount = dto.BankAccount
	}
	actStatus := dto.ActiveStatus
	if actStatus == "" {
		actStatus = dto.Status
	}
	if actStatus != "" {
		existing.ActiveStatus = actStatus
		existing.Status = actStatus
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateVendor(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *scService) DeleteVendor(id string) error {
	return s.repo.DeleteVendor(id)
}

// ── Master Kategori Item ────────────────────────────────────────────

func (s *scService) GetAllItemCategories(search, itemType string) ([]model.MstScItemCategory, error) {
	return s.repo.FindItemCategories(search, itemType)
}

func (s *scService) GetItemCategoriesPaginated(page, perPage int, search, itemType string) ([]model.MstScItemCategory, int64, error) {
	return s.repo.FindItemCategoriesPaginated(page, perPage, search, itemType)
}

func (s *scService) GetItemCategoryByID(id string) (*model.MstScItemCategory, error) {
	return s.repo.FindItemCategoryByID(id)
}

func (s *scService) CreateItemCategory(dto *ItemCategoryDTO) (*model.MstScItemCategory, error) {
	catName := strings.TrimSpace(dto.CategoryName)
	if catName == "" {
		catName = strings.TrimSpace(dto.Name)
	}
	if catName == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("CAT-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	iType := dto.ItemType
	if iType == "" {
		iType = dto.Type
	}
	if iType == "" {
		iType = "HARDWARE"
	}

	cat := &model.MstScItemCategory{
		ID:           uuid.NewString(),
		Code:         code,
		CategoryName: catName,
		ItemType:     iType,
		Description:  dto.Description,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repo.CreateItemCategory(cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *scService) UpdateItemCategory(id string, dto *ItemCategoryDTO) (*model.MstScItemCategory, error) {
	existing, err := s.repo.FindItemCategoryByID(id)
	if err != nil {
		return nil, errors.New("kategori item tidak ditemukan")
	}

	catName := strings.TrimSpace(dto.CategoryName)
	if catName == "" {
		catName = strings.TrimSpace(dto.Name)
	}
	if catName != "" {
		existing.CategoryName = catName
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	iType := dto.ItemType
	if iType == "" {
		iType = dto.Type
	}
	if iType != "" {
		existing.ItemType = iType
	}
	existing.Description = dto.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateItemCategory(existing); err != nil {
		return nil, err
	}
	return existing, nil
}

func (s *scService) DeleteItemCategory(id string) error {
	return s.repo.DeleteItemCategory(id)
}

// ── Master Katalog Item ─────────────────────────────────────────────

func (s *scService) GetAllItems(categoryID, search string) ([]model.MstScItem, error) {
	return s.repo.FindItems(categoryID, search)
}

func (s *scService) GetItemsPaginated(page, perPage int, categoryID, search string) ([]model.MstScItem, int64, error) {
	return s.repo.FindItemsPaginated(page, perPage, categoryID, search)
}

func (s *scService) GetItemByID(id string) (*model.MstScItem, error) {
	return s.repo.FindItemByID(id)
}

func (s *scService) CreateItem(dto *ItemDTO) (*model.MstScItem, error) {
	itemName := strings.TrimSpace(dto.ItemName)
	if itemName == "" {
		itemName = strings.TrimSpace(dto.Name)
	}
	if itemName == "" {
		return nil, errors.New("nama item wajib diisi")
	}
	if dto.CategoryID == "" {
		return nil, errors.New("kategori item wajib dipilih")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("ITM-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	unit := dto.Unit
	if unit == "" {
		unit = "Unit"
	}

	price := dto.EstimatedPrice
	if price == 0 {
		price = dto.Price
	}

	status := dto.Status
	if status == "" {
		status = "ACTIVE"
	}

	item := &model.MstScItem{
		ID:              uuid.NewString(),
		CategoryID:      dto.CategoryID,
		Code:            code,
		ItemName:        itemName,
		Unit:            unit,
		EstimatedPrice:  price,
		DefaultVendorID: dto.DefaultVendorID,
		MinStock:        dto.MinStock,
		CurrentStock:    dto.CurrentStock,
		Description:     dto.Description,
		Status:          status,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}
	if err := s.repo.CreateItem(item); err != nil {
		return nil, err
	}
	return s.repo.FindItemByID(item.ID)
}

func (s *scService) UpdateItem(id string, dto *ItemDTO) (*model.MstScItem, error) {
	existing, err := s.repo.FindItemByID(id)
	if err != nil {
		return nil, errors.New("item tidak ditemukan")
	}

	itemName := strings.TrimSpace(dto.ItemName)
	if itemName == "" {
		itemName = strings.TrimSpace(dto.Name)
	}
	if itemName != "" {
		existing.ItemName = itemName
	}
	if dto.CategoryID != "" {
		existing.CategoryID = dto.CategoryID
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	if dto.Unit != "" {
		existing.Unit = dto.Unit
	}
	price := dto.EstimatedPrice
	if price == 0 {
		price = dto.Price
	}
	if price > 0 {
		existing.EstimatedPrice = price
	}
	if dto.DefaultVendorID != nil {
		existing.DefaultVendorID = dto.DefaultVendorID
	}
	if dto.MinStock > 0 {
		existing.MinStock = dto.MinStock
	}
	if dto.CurrentStock > 0 {
		existing.CurrentStock = dto.CurrentStock
	}
	existing.Description = dto.Description
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateItem(existing); err != nil {
		return nil, err
	}
	return s.repo.FindItemByID(id)
}

func (s *scService) DeleteItem(id string) error {
	return s.repo.DeleteItem(id)
}

// ── Master Lokasi ───────────────────────────────────────────────────

func (s *scService) GetAllLocations(search string) ([]model.MstScLocation, error) {
	return s.repo.FindLocations(search)
}

func (s *scService) GetLocationsPaginated(page, perPage int, search string) ([]model.MstScLocation, int64, error) {
	return s.repo.FindLocationsPaginated(page, perPage, search)
}

func (s *scService) GetLocationByID(id string) (*model.MstScLocation, error) {
	return s.repo.FindLocationByID(id)
}

func (s *scService) CreateLocation(dto *LocationDTO) (*model.MstScLocation, error) {
	locName := strings.TrimSpace(dto.LocationName)
	if locName == "" {
		locName = strings.TrimSpace(dto.Name)
	}
	if locName == "" {
		return nil, errors.New("nama lokasi wajib diisi")
	}

	code := strings.TrimSpace(dto.Code)
	if code == "" {
		code = fmt.Sprintf("LOC-%s", strings.ToUpper(uuid.NewString()[:6]))
	}

	lType := dto.Type
	if lType == "" {
		lType = "warehouse"
	}

	loc := &model.MstScLocation{
		ID:                 uuid.NewString(),
		Code:               code,
		LocationName:       locName,
		Type:               lType,
		Address:            dto.Address,
		PersonInChargeID:   dto.PersonInChargeID,
		Description:        dto.Description,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.repo.CreateLocation(loc); err != nil {
		return nil, err
	}
	return s.repo.FindLocationByID(loc.ID)
}

func (s *scService) UpdateLocation(id string, dto *LocationDTO) (*model.MstScLocation, error) {
	existing, err := s.repo.FindLocationByID(id)
	if err != nil {
		return nil, errors.New("lokasi tidak ditemukan")
	}

	locName := strings.TrimSpace(dto.LocationName)
	if locName == "" {
		locName = strings.TrimSpace(dto.Name)
	}
	if locName != "" {
		existing.LocationName = locName
	}
	if dto.Code != "" {
		existing.Code = dto.Code
	}
	if dto.Type != "" {
		existing.Type = dto.Type
	}
	existing.Address = dto.Address
	existing.PersonInChargeID = dto.PersonInChargeID
	existing.Description = dto.Description
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateLocation(existing); err != nil {
		return nil, err
	}
	return s.repo.FindLocationByID(id)
}

func (s *scService) DeleteLocation(id string) error {
	return s.repo.DeleteLocation(id)
}

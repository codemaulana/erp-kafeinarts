package model

import (
	"time"
)

// TrxScPurchaseRequest merepresentasikan pengajuan pembelian (trx_sc_purchase_requests)
type TrxScPurchaseRequest struct {
	ID             string                     `gorm:"type:varchar(36);primaryKey" json:"id"`
	PrNumber       string                     `gorm:"type:varchar(50);not null;uniqueIndex" json:"pr_number"`
	DepartmentID   string                     `gorm:"type:varchar(36);default:'';index" json:"department_id"`
	DepartmentName string                     `gorm:"->" json:"department_name,omitempty"`
	RequesterName  string                     `gorm:"type:varchar(100);default:''" json:"requester_name"`
	RequestedByName string                    `gorm:"->" json:"requested_by_name,omitempty"`
	RequestDate    string                     `gorm:"type:varchar(20);not null" json:"request_date"`
	Priority       string                     `gorm:"type:varchar(20);default:'MEDIUM'" json:"priority"` // LOW, MEDIUM, HIGH, URGENT
	Notes          string                     `gorm:"type:text;default:''" json:"notes"`
	ApprovalStatus string                     `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"approval_status"` // DRAFT, SUBMITTED, APPROVED, REJECTED
	Status         string                     `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"status"`
	Items          []TrxScPurchaseRequestItem `gorm:"foreignKey:PrID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt      time.Time                  `json:"created_at"`
	UpdatedAt      time.Time                  `json:"updated_at"`
}

func (TrxScPurchaseRequest) TableName() string {
	return "trx_sc_purchase_requests"
}

// TrxScPurchaseRequestItem merepresentasikan item pada PR (trx_sc_purchase_request_items)
type TrxScPurchaseRequestItem struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	PrID           string    `gorm:"type:varchar(36);not null;index" json:"pr_id"`
	ItemID         string    `gorm:"type:varchar(36);not null;index" json:"item_id"`
	ItemName       string    `gorm:"->" json:"item_name,omitempty"`
	Quantity       int       `gorm:"default:1" json:"quantity"`
	Unit           string    `gorm:"type:varchar(20);default:'Unit'" json:"unit"`
	EstimatedPrice float64   `gorm:"type:decimal(15,2);default:0" json:"estimated_price"`
	Subtotal       float64   `gorm:"type:decimal(15,2);default:0" json:"subtotal"`
	Notes          string    `gorm:"type:text;default:''" json:"notes"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (TrxScPurchaseRequestItem) TableName() string {
	return "trx_sc_purchase_request_items"
}

// TrxScPurchaseOrder merepresentasikan pesanan pembelian (trx_sc_purchase_orders)
type TrxScPurchaseOrder struct {
	ID                   string                   `gorm:"type:varchar(36);primaryKey" json:"id"`
	PoNumber             string                   `gorm:"type:varchar(50);not null;uniqueIndex" json:"po_number"`
	PrID                 *string                  `gorm:"type:varchar(36);default:null;index" json:"pr_id"`
	PrNumber             string                   `gorm:"->" json:"pr_number,omitempty"`
	VendorID             string                   `gorm:"type:varchar(36);not null;index" json:"vendor_id"`
	VendorName           string                   `gorm:"->" json:"vendor_name,omitempty"`
	PoDate               string                   `gorm:"type:varchar(20);not null" json:"po_date"` // order_date
	ExpectedDeliveryDate string                   `gorm:"type:varchar(20);default:''" json:"expected_delivery_date"`
	Subtotal             float64                  `gorm:"type:decimal(15,2);default:0" json:"subtotal"`
	TaxAmount            float64                  `gorm:"type:decimal(15,2);default:0" json:"tax_amount"`
	TotalAmount          float64                  `gorm:"type:decimal(15,2);default:0" json:"total_amount"`
	Status               string                   `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"status"` // DRAFT, SENT, PARTIALLY_RECEIVED, COMPLETED, CANCELLED
	PoStatus             string                   `gorm:"type:varchar(30);not null;default:'DRAFT'" json:"po_status"`
	Notes                string                   `gorm:"type:text;default:''" json:"notes"`
	Items                []TrxScPurchaseOrderItem `gorm:"foreignKey:PoID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt            time.Time                `json:"created_at"`
	UpdatedAt            time.Time                `json:"updated_at"`
}

func (TrxScPurchaseOrder) TableName() string {
	return "trx_sc_purchase_orders"
}

// TrxScPurchaseOrderItem merepresentasikan item pada PO (trx_sc_purchase_order_items)
type TrxScPurchaseOrderItem struct {
	ID         string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	PoID       string    `gorm:"type:varchar(36);not null;index" json:"po_id"`
	ItemID     string    `gorm:"type:varchar(36);not null;index" json:"item_id"`
	ItemName   string    `gorm:"->" json:"item_name,omitempty"`
	Quantity   int       `gorm:"default:1" json:"quantity"`
	Unit       string    `gorm:"type:varchar(20);default:'Unit'" json:"unit"`
	UnitPrice  float64   `gorm:"type:decimal(15,2);default:0" json:"unit_price"`
	TotalPrice float64   `gorm:"type:decimal(15,2);default:0" json:"total_price"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (TrxScPurchaseOrderItem) TableName() string {
	return "trx_sc_purchase_order_items"
}

// TrxScGoodsReceipt merepresentasikan penerimaan barang (trx_sc_goods_receipts)
type TrxScGoodsReceipt struct {
	ID                 string                  `gorm:"type:varchar(36);primaryKey" json:"id"`
	ReceiptNumber      string                  `gorm:"type:varchar(50);not null;uniqueIndex" json:"receipt_number"` // gr_number
	GrNumber           string                  `gorm:"->" json:"gr_number"`
	PoID               string                  `gorm:"type:varchar(36);not null;index" json:"po_id"`
	PoNumber           string                  `gorm:"->" json:"po_number,omitempty"`
	VendorName         string                  `gorm:"->" json:"vendor_name,omitempty"`
	ReceivedBy         string                  `gorm:"type:varchar(100);not null" json:"received_by"`
	ReceivedByName     string                  `gorm:"->" json:"received_by_name,omitempty"`
	ReceiptDate        string                  `gorm:"type:varchar(20);not null" json:"receipt_date"`
	DeliveryNoteNumber string                  `gorm:"type:varchar(100);default:''" json:"delivery_note_number"`
	AttachmentURL      string                  `gorm:"type:text;default:''" json:"attachment_url"`
	Notes              string                  `gorm:"type:text;default:''" json:"notes"`
	Status             string                  `gorm:"type:varchar(30);default:'COMPLETE'" json:"status"` // COMPLETE, PARTIAL
	Items              []TrxScGoodsReceiptItem `gorm:"foreignKey:GrID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt          time.Time               `json:"created_at"`
	UpdatedAt          time.Time               `json:"updated_at"`
}

func (TrxScGoodsReceipt) TableName() string {
	return "trx_sc_goods_receipts"
}

// TrxScGoodsReceiptItem merepresentasikan item pada penerimaan barang (trx_sc_goods_receipt_items)
type TrxScGoodsReceiptItem struct {
	ID               string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	GrID             string    `gorm:"type:varchar(36);not null;index" json:"gr_id"`
	ItemID           string    `gorm:"type:varchar(36);not null;index" json:"item_id"`
	ItemName         string    `gorm:"->" json:"item_name,omitempty"`
	OrderedQuantity  int       `gorm:"default:0" json:"ordered_quantity"`
	ReceivedQuantity int       `gorm:"default:0" json:"received_quantity"`
	Condition        string    `gorm:"type:varchar(30);default:'GOOD'" json:"condition"` // GOOD, DAMAGED, REJECTED
	Notes            string    `gorm:"type:text;default:''" json:"notes"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

func (TrxScGoodsReceiptItem) TableName() string {
	return "trx_sc_goods_receipt_items"
}

// TrxScAsset merepresentasikan aset / inventaris (trx_sc_assets)
type TrxScAsset struct {
	ID                       string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	AssetCode                string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"asset_code"`
	Name                     string    `gorm:"type:varchar(150);default:''" json:"name"`
	ItemID                   string    `gorm:"type:varchar(36);default:'';index" json:"item_id"`
	ItemName                 string    `gorm:"->" json:"item_name,omitempty"`
	CategoryID               string    `gorm:"type:varchar(36);default:'';index" json:"category_id"`
	CategoryName             string    `gorm:"->" json:"category_name,omitempty"`
	SerialOrLicenseKey       string    `gorm:"type:varchar(100);default:''" json:"serial_or_license_key"` // serial_number
	AcquisitionDate          string    `gorm:"type:varchar(20);default:''" json:"acquisition_date"`       // purchase_date
	AcquisitionValue         float64   `gorm:"type:decimal(15,2);default:0" json:"acquisition_value"`     // purchase_price
	WarrantyOrLicenseExpiry  string    `gorm:"type:varchar(20);default:''" json:"warranty_or_license_expiry"`
	AssignedToEmployeeID     *string   `gorm:"type:varchar(36);default:null;index" json:"assigned_to_employee_id"`
	AssignedToEmployeeName   string    `gorm:"->" json:"assigned_to_employee_name,omitempty"`
	AssignedToName           string    `gorm:"->" json:"assigned_to_name,omitempty"`
	DepartmentID             *string   `gorm:"type:varchar(36);default:null;index" json:"department_id"`
	DepartmentName           string    `gorm:"->" json:"department_name,omitempty"`
	LocationID               *string   `gorm:"type:varchar(36);default:null;index" json:"location_id"`
	LocationName             string    `gorm:"->" json:"location_name,omitempty"`
	Status                   string    `gorm:"type:varchar(30);default:'AVAILABLE'" json:"status"` // AVAILABLE, IN_USE, MAINTENANCE, DAMAGED, DISPOSED
	Condition                string    `gorm:"type:varchar(30);default:'GOOD'" json:"condition"`   // GOOD, FAIR, POOR
	Notes                    string    `gorm:"type:text;default:''" json:"notes"`
	CreatedAt                time.Time `json:"created_at"`
	UpdatedAt                time.Time `json:"updated_at"`
}

func (TrxScAsset) TableName() string {
	return "trx_sc_assets"
}

// TrxScSubscription merepresentasikan langganan software / lisensi berkala (trx_sc_subscriptions)
type TrxScSubscription struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	LicenseCode  string    `gorm:"type:varchar(50);not null;uniqueIndex" json:"license_code"`
	Name         string    `gorm:"type:varchar(150);default:''" json:"name"`
	ItemID       string    `gorm:"type:varchar(36);default:'';index" json:"item_id"`
	ItemName     string    `gorm:"->" json:"item_name,omitempty"`
	VendorID     string    `gorm:"type:varchar(36);not null;index" json:"vendor_id"`
	VendorName   string    `gorm:"->" json:"vendor_name,omitempty"`
	BillingCycle string    `gorm:"type:varchar(30);default:'MONTHLY'" json:"billing_cycle"` // MONTHLY, YEARLY
	Amount       float64   `gorm:"type:decimal(15,2);default:0" json:"amount"`              // cost
	StartDate    string    `gorm:"type:varchar(20);not null" json:"start_date"`
	EndDate      string    `gorm:"type:varchar(20);not null" json:"end_date"`               // expiry_date
	AutoRenewal  bool      `gorm:"default:false" json:"auto_renewal"`
	Status       string    `gorm:"type:varchar(30);default:'ACTIVE'" json:"status"`         // ACTIVE, EXPIRED, CANCELLED
	Notes        string    `gorm:"type:text;default:''" json:"notes"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (TrxScSubscription) TableName() string {
	return "trx_sc_subscriptions"
}

package model

import (
	"time"
)

// MstScVendor merepresentasikan master vendor / supplier (mst_sc_vendors)
type MstScVendor struct {
	ID             string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code           string    `gorm:"type:varchar(20);default:''" json:"code"`
	VendorName     string    `gorm:"type:varchar(150);not null" json:"vendor_name"` // name
	VendorCategory string    `gorm:"type:varchar(50);default:'OTHER'" json:"vendor_category"` // category
	ContactPerson  string    `gorm:"type:varchar(100);default:''" json:"contact_person"` // pic_name
	Phone          string    `gorm:"type:varchar(30);default:''" json:"phone"`
	PicPhone       string    `gorm:"type:varchar(30);default:''" json:"pic_phone"`
	Email          string    `gorm:"type:varchar(100);default:''" json:"email"`
	Address        string    `gorm:"type:text;default:''" json:"address"`
	BankName       string    `gorm:"type:varchar(50);default:''" json:"bank_name"`
	BankAccount    string    `gorm:"type:varchar(50);default:''" json:"bank_account"`
	ActiveStatus   string    `gorm:"type:varchar(20);default:'ACTIVE'" json:"active_status"` // status: ACTIVE, INACTIVE
	Status         string    `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	CreatedAt      time.Time `json:"created_at"`
	UpdatedAt      time.Time `json:"updated_at"`
}

func (MstScVendor) TableName() string {
	return "mst_sc_vendors"
}

// MstScItemCategory merepresentasikan master kategori item / aset (mst_sc_item_categories)
type MstScItemCategory struct {
	ID           string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code         string    `gorm:"type:varchar(20);default:''" json:"code"`
	CategoryName string    `gorm:"type:varchar(100);not null" json:"category_name"` // name
	ItemType     string    `gorm:"type:varchar(50);default:'HARDWARE'" json:"item_type"` // HARDWARE, SOFTWARE_LICENSE, SUBSCRIPTION, SERVICE
	Description  string    `gorm:"type:text;default:''" json:"description"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

func (MstScItemCategory) TableName() string {
	return "mst_sc_item_categories"
}

// MstScItem merepresentasikan master item katalog (mst_sc_items)
type MstScItem struct {
	ID                string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	CategoryID        string    `gorm:"type:varchar(36);not null;index" json:"category_id"`
	CategoryName      string    `gorm:"->" json:"category_name,omitempty"`
	Code              string    `gorm:"type:varchar(50);default:''" json:"code"`
	ItemName          string    `gorm:"type:varchar(150);not null" json:"item_name"` // name
	Unit              string    `gorm:"type:varchar(20);default:'Unit'" json:"unit"`
	EstimatedPrice    float64   `gorm:"type:decimal(15,2);default:0" json:"estimated_price"` // price
	DefaultVendorID   *string   `gorm:"type:varchar(36);default:null" json:"default_vendor_id"`
	DefaultVendorName string    `gorm:"->" json:"default_vendor_name,omitempty"`
	MinStock          int       `gorm:"default:0" json:"min_stock"`
	CurrentStock      int       `gorm:"default:0" json:"current_stock"`
	Description       string    `gorm:"type:text;default:''" json:"description"`
	Status            string    `gorm:"type:varchar(20);default:'ACTIVE'" json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (MstScItem) TableName() string {
	return "mst_sc_items"
}

// MstScLocation merepresentasikan master lokasi inventaris (mst_sc_locations)
type MstScLocation struct {
	ID                  string    `gorm:"type:varchar(36);primaryKey" json:"id"`
	Code                string    `gorm:"type:varchar(20);default:''" json:"code"`
	LocationName        string    `gorm:"type:varchar(100);not null" json:"location_name"` // name
	Type                string    `gorm:"type:varchar(50);default:'warehouse'" json:"type"`
	Address             string    `gorm:"type:text;default:''" json:"address"`
	PersonInChargeID    *string   `gorm:"type:varchar(36);default:null" json:"person_in_charge_id"`
	PersonInChargeName  string    `gorm:"->" json:"person_in_charge_name,omitempty"`
	Description         string    `gorm:"type:text;default:''" json:"description"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`
}

func (MstScLocation) TableName() string {
	return "mst_sc_locations"
}

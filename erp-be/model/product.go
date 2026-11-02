package model

import "time"

// Product merepresentasikan produk dalam modul inventaris Kafeinarts.
type Product struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	SKU         string    `gorm:"size:50;uniqueIndex;not null" json:"sku"`
	Name        string    `gorm:"size:200;not null" json:"name"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	Price       float64   `gorm:"type:decimal(12,2);not null;default:0" json:"price"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// StockAdjustment mencatat riwayat mutasi stok produk (masuk, keluar, koreksi).
type StockAdjustment struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	ProductID uint      `gorm:"index;not null" json:"product_id"`
	Product   Product   `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Type      string    `gorm:"size:20;not null" json:"type"` // "in", "out", "set"
	Amount    int       `gorm:"not null" json:"amount"`
	Reason    string    `gorm:"size:255" json:"reason"`
	CreatedAt time.Time `json:"created_at"`
}

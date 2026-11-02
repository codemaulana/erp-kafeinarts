package model

import "time"

// Transaction merepresentasikan transaksi pemasukan (Debit) atau pengeluaran (Credit),
// mencakup Accounts Receivable (AR / Invoicing) dan Accounts Payable (AP / Tagihan Vendor).
type Transaction struct {
	ID               uint       `gorm:"primaryKey" json:"id"`
	TransactionID    string     `gorm:"size:36;uniqueIndex;not null;default:''" json:"transaction_id"` // UUID unik transaksi
	ReferenceNo      string     `gorm:"size:50;uniqueIndex;not null;default:''" json:"reference_no"`   // No Invoice / Tagihan (cth: INV-2026-001)
	Type             string     `gorm:"size:20;not null;default:'Debit'" json:"type"`                  // Enum: "Debit" (Pemasukan/Piutang) / "Credit" (Pengeluaran/Hutang)
	Category         string     `gorm:"size:100;not null;default:''" json:"category"`                  // Cth: "Software License", "Project Revenue", "Cloud Infrastructure"
	Amount           float64    `gorm:"type:decimal(14,2);not null;default:0" json:"amount"`          // Nominal pokok
	TaxAmount        float64    `gorm:"type:decimal(14,2);not null;default:0" json:"tax_amount"` // Nilai pajak (PPN/PPh)
	TotalAmount      float64    `gorm:"type:decimal(14,2);not null;default:0" json:"total_amount"` // Amount + TaxAmount
	Currency         string     `gorm:"size:10;not null;default:'IDR'" json:"currency"`     // Mata uang: "IDR", "USD", dll.
	Status           string     `gorm:"size:20;not null;default:'Draft'" json:"status"`     // Enum: "Draft", "Sent", "Paid", "Overdue"
	DueDate          *time.Time `json:"due_date"`                                           // Batas waktu pembayaran
	TransactionDate  time.Time  `gorm:"not null;default:CURRENT_TIMESTAMP" json:"transaction_date"` // Tanggal transaksi/terbit
	ClientVendorName string     `gorm:"size:150" json:"client_vendor_name"`                 // Nama Klien (AR) atau Nama Vendor (AP)
	Description      string     `gorm:"type:text" json:"description"`                       // Keterangan / rincian
	PaidAt           *time.Time `json:"paid_at,omitempty"`                                  // Waktu pelunasan
	CreatedAt        time.Time  `json:"created_at"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

// FinanceSummary merepresentasikan agregasi ringkas data keuangan studio.
type FinanceSummary struct {
	TotalIncome      float64 `json:"total_income"`
	TotalExpense     float64 `json:"total_expense"`
	NetBalance       float64 `json:"net_balance"`
	TotalReceivable  float64 `json:"total_receivable"` // Total piutang belum lunas (AR)
	TotalPayable     float64 `json:"total_payable"`    // Total hutang belum lunas (AP)
	TransactionCount int64   `json:"transaction_count"`
}

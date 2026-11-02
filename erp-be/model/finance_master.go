package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// MstFinanceCategory merepresentasikan Master Kategori Transaksi keuangan
// (Tabel: mst_finance_categories)
type MstFinanceCategory struct {
	ID           string    `gorm:"primaryKey;size:36" json:"id"`
	CategoryName string    `gorm:"size:100;not null" json:"category_name"` // Nama Kategori (cth: Project Revenue, Biaya Server, Gaji)
	NamaKategori string    `gorm:"-" json:"namaKategori"`
	CashType     string    `gorm:"size:20;not null" json:"cash_type"`      // ENUM: 'INCOME', 'EXPENSE'
	TipeKas      string    `gorm:"-" json:"tipeKas"`                        // 'Masuk', 'Keluar'
	Description  string    `gorm:"type:text" json:"description"`
	Deskripsi    string    `gorm:"-" json:"deskripsi"`
	CreatedAt    time.Time `json:"created_at"`
}

func (MstFinanceCategory) TableName() string {
	return "mst_finance_categories"
}

func (c *MstFinanceCategory) AfterFind(tx *gorm.DB) (err error) {
	c.NamaKategori = c.CategoryName
	if strings.EqualFold(c.CashType, "INCOME") || c.CashType == "Masuk" {
		c.TipeKas = "Masuk"
	} else {
		c.TipeKas = "Keluar"
	}
	c.Deskripsi = c.Description
	return nil
}

// MstFinanceAccount merepresentasikan Master Rekening / Kas Bank
// (Tabel: mst_finance_accounts)
type MstFinanceAccount struct {
	ID                string    `gorm:"primaryKey;size:36" json:"id"`
	BankName          string    `gorm:"size:100;not null" json:"bank_name"`           // Nama Bank/Kas (cth: BCA, Mandiri, Petty Cash)
	NamaBank          string    `gorm:"-" json:"namaBank"`
	AccountNumber     string    `gorm:"size:50;not null" json:"account_number"`       // Nomor Rekening (validasi angka)
	NoRekening        string    `gorm:"-" json:"noRekening"`
	AccountHolderName string    `gorm:"size:100;not null" json:"account_holder_name"` // Atas Nama
	AtasNama          string    `gorm:"-" json:"atasNama"`
	CurrentBalance    float64   `gorm:"type:decimal(15,2);not null;default:0" json:"current_balance"` // Saldo Saat Ini
	SaldoAwal         float64   `gorm:"-" json:"saldoAwal"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
}

func (MstFinanceAccount) TableName() string {
	return "mst_finance_accounts"
}

func (a *MstFinanceAccount) AfterFind(tx *gorm.DB) (err error) {
	a.NamaBank = a.BankName
	a.NoRekening = a.AccountNumber
	a.AtasNama = a.AccountHolderName
	a.SaldoAwal = a.CurrentBalance
	return nil
}

// MstFinanceTax merepresentasikan Master Pajak
// (Tabel: mst_finance_taxes)
type MstFinanceTax struct {
	ID          string    `gorm:"primaryKey;size:36" json:"id"`
	TaxName     string    `gorm:"size:100;not null" json:"tax_name"`             // Nama Pajak (cth: PPN 11%, PPh 21)
	NamaPajak   string    `gorm:"-" json:"namaPajak"`
	Percentage  float64   `gorm:"type:decimal(5,2);not null;default:0" json:"percentage"` // Persentase (%)
	Persentase  float64   `gorm:"-" json:"persentase"`
	Description string    `gorm:"type:text" json:"description"`
	Keterangan  string    `gorm:"-" json:"keterangan"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (MstFinanceTax) TableName() string {
	return "mst_finance_taxes"
}

func (t *MstFinanceTax) AfterFind(tx *gorm.DB) (err error) {
	t.NamaPajak = t.TaxName
	t.Persentase = t.Percentage
	t.Keterangan = t.Description
	return nil
}

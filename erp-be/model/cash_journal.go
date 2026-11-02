package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// TrxCashJournal merepresentasikan Jurnal Kas (Pencatatan Harian)
// (Tabel: trx_cash_journals)
type TrxCashJournal struct {
	ID                string              `gorm:"primaryKey;size:36" json:"id"`
	TransactionNumber string              `gorm:"size:50;uniqueIndex;not null" json:"transaction_number"` // No. Transaksi
	TransactionDate   time.Time           `gorm:"not null" json:"transaction_date"`                       // Tanggal
	Tanggal           string              `gorm:"-" json:"tanggal"`
	AccountID         string              `gorm:"size:36;index;not null" json:"account_id"`               // FK -> mst_finance_accounts.id
	AkunBank          string              `gorm:"-" json:"akunBank"`
	Account           *MstFinanceAccount  `gorm:"foreignKey:AccountID" json:"account,omitempty"`
	NamaBank          string              `gorm:"-" json:"namaBank"`
	CategoryID        string              `gorm:"size:36;index;not null" json:"category_id"`              // FK -> mst_finance_categories.id
	Kategori          string              `gorm:"-" json:"kategori"`
	Category          *MstFinanceCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"`
	NamaKategori      string              `gorm:"-" json:"namaKategori"`
	TransactionType   string              `gorm:"size:20;not null" json:"transaction_type"`               // ENUM: 'INCOME', 'EXPENSE'
	Tipe              string              `gorm:"-" json:"tipe"`                                          // 'Masuk', 'Keluar'
	Amount            float64             `gorm:"type:decimal(15,2);not null" json:"amount"`              // Nominal
	Nominal           float64             `gorm:"-" json:"nominal"`
	Description       string              `gorm:"type:text" json:"description"`                           // Keterangan Tambahan
	Keterangan        string              `gorm:"-" json:"keterangan"`
	AttachmentURL     string              `gorm:"size:255" json:"attachment_url"`                         // Bukti Transaksi (Struk/Kwitansi)
	BuktiFile         string              `gorm:"-" json:"buktiFile"`
	CreatedAt         time.Time           `json:"created_at"`
	UpdatedAt         time.Time           `json:"updated_at"`
}

func (TrxCashJournal) TableName() string {
	return "trx_cash_journals"
}

func (j *TrxCashJournal) AfterFind(tx *gorm.DB) (err error) {
	if !j.TransactionDate.IsZero() {
		j.Tanggal = j.TransactionDate.Format("2006-01-02")
	}
	j.AkunBank = j.AccountID
	if j.Account != nil {
		j.NamaBank = j.Account.BankName
	}
	if strings.EqualFold(j.TransactionType, "INCOME") || j.TransactionType == "Masuk" {
		j.Tipe = "Masuk"
	} else {
		j.Tipe = "Keluar"
	}
	j.Kategori = j.CategoryID
	if j.Category != nil {
		j.NamaKategori = j.Category.CategoryName
	}
	j.Nominal = j.Amount
	j.Keterangan = j.Description
	j.BuktiFile = j.AttachmentURL
	return nil
}

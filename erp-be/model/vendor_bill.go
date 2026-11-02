package model

import (
	"strings"
	"time"

	"gorm.io/gorm"
)

// TrxVendorBill merepresentasikan Tagihan Vendor / Hutang (AP)
// (Tabel: trx_vendor_bills)
type TrxVendorBill struct {
	ID             string              `gorm:"primaryKey;size:36" json:"id"`
	BillNumber     string              `gorm:"size:50;not null" json:"bill_number"`              // Nomor Tagihan
	NoTagihan      string              `gorm:"-" json:"noTagihan"`
	VendorID       string              `gorm:"size:36;index" json:"vendor_id"`                   // FK -> mst_sc_vendors.id
	VendorName     string              `gorm:"size:150;not null" json:"vendor_name"`             // Nama Vendor (cth: AWS, Google Cloud, lisensi software)
	Vendor         string              `gorm:"-" json:"vendor"`
	CategoryID     string              `gorm:"size:36;index;not null" json:"category_id"`        // FK -> mst_finance_categories.id
	Kategori       string              `gorm:"-" json:"kategori"`
	Category       *MstFinanceCategory `gorm:"foreignKey:CategoryID" json:"category,omitempty"` // Kategori Biaya
	NamaKategori   string              `gorm:"-" json:"namaKategori"`
	BillDate       time.Time           `gorm:"not null" json:"bill_date"`                        // Tanggal Tagihan
	TanggalTagihan string              `gorm:"-" json:"tanggalTagihan"`
	DueDate        time.Time           `gorm:"not null" json:"due_date"`                         // Jatuh Tempo
	JatuhTempo     string              `gorm:"-" json:"jatuhTempo"`
	Amount         float64             `gorm:"type:decimal(15,2);not null" json:"amount"`        // Nominal Tagihan
	Nominal        float64             `gorm:"-" json:"nominal"`
	Status         string              `gorm:"size:20;not null;default:'PENDING'" json:"status"` // ENUM: 'PENDING', 'PAID', 'OVERDUE'
	AttachmentURL  string              `gorm:"size:255" json:"attachment_url"`                   // Lampiran Tagihan
	BuktiFile      string              `gorm:"-" json:"buktiFile"`
	PaidAt         *time.Time          `json:"paid_at,omitempty"`
	CreatedAt      time.Time           `json:"created_at"`
	UpdatedAt      time.Time           `json:"updated_at"`
}

func (TrxVendorBill) TableName() string {
	return "trx_vendor_bills"
}

func (b *TrxVendorBill) AfterFind(tx *gorm.DB) (err error) {
	b.NoTagihan = b.BillNumber
	b.Vendor = b.VendorName
	b.Kategori = b.CategoryID
	if b.Category != nil {
		b.NamaKategori = b.Category.CategoryName
	}
	if !b.BillDate.IsZero() {
		b.TanggalTagihan = b.BillDate.Format("2006-01-02")
	}
	if !b.DueDate.IsZero() {
		b.JatuhTempo = b.DueDate.Format("2006-01-02")
	}
	b.Nominal = b.Amount
	b.BuktiFile = b.AttachmentURL

	switch strings.ToUpper(b.Status) {
	case "PAID", "LUNAS":
		b.Status = "Lunas"
	case "OVERDUE", "JATUH TEMPO":
		b.Status = "Jatuh Tempo"
	default:
		b.Status = "Belum Lunas"
	}
	return nil
}

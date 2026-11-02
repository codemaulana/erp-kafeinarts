package model

import (
	"time"

	"gorm.io/gorm"
)

// TrxInvoice merepresentasikan Master Invoice Penagihan / Piutang (AR)
// (Tabel: trx_invoices)
type TrxInvoice struct {
	ID            string           `gorm:"primaryKey;size:36" json:"id"`
	InvoiceNumber string           `gorm:"size:50;uniqueIndex;not null" json:"invoice_number"` // Nomor Invoice (Auto-generate)
	NoInvoice     string           `gorm:"-" json:"noInvoice"`
	ClientID      string           `gorm:"size:36;index" json:"client_id"`                     // FK -> mst_crm_clients.id
	ClientName    string           `gorm:"size:150;not null" json:"client_name"`               // Nama Klien
	Klien         string           `gorm:"-" json:"klien"`
	ProjectID     string           `gorm:"size:36;index" json:"project_id"`                    // FK -> trx_projects.id
	ProjectName   string           `gorm:"size:150" json:"project_name"`                       // Nama Proyek
	Proyek        string           `gorm:"-" json:"proyek"`
	IssueDate     time.Time        `gorm:"not null" json:"issue_date"`                         // Tgl Terbit
	TglTerbit     string           `gorm:"-" json:"tglTerbit"`
	DueDate       time.Time        `gorm:"not null" json:"due_date"`                           // Jatuh Tempo
	TglJatuhTempo string           `gorm:"-" json:"tglJatuhTempo"`
	Subtotal      float64          `gorm:"type:decimal(15,2);not null;default:0" json:"subtotal"`
	TotalTagihan  float64          `gorm:"-" json:"totalTagihan"`
	TotalTax      float64          `gorm:"type:decimal(15,2);not null;default:0" json:"total_tax"`
	Pajak         string           `gorm:"-" json:"pajak,omitempty"`
	CatatanPajak  string           `gorm:"-" json:"catatanPajak,omitempty"`
	GrandTotal    float64          `gorm:"type:decimal(15,2);not null;default:0" json:"grand_total"` // Total Tagihan
	GrandTotalAlt float64          `gorm:"-" json:"grandTotal"`
	Status        string           `gorm:"size:20;not null;default:'DRAFT'" json:"status"`           // ENUM: 'DRAFT', 'SENT', 'PARTIAL', 'PAID', 'CANCELLED'
	Notes         string           `gorm:"type:text" json:"notes"`
	Items         []TrxInvoiceItem `gorm:"foreignKey:InvoiceID;constraint:OnDelete:CASCADE" json:"items"`
	PaidAt        *time.Time       `json:"paid_at,omitempty"`
	CreatedAt     time.Time        `json:"created_at"`
	UpdatedAt     time.Time        `json:"updated_at"`
}

func (TrxInvoice) TableName() string {
	return "trx_invoices"
}

func (inv *TrxInvoice) AfterFind(tx *gorm.DB) (err error) {
	inv.NoInvoice = inv.InvoiceNumber
	inv.Klien = inv.ClientName
	inv.Proyek = inv.ProjectName
	if !inv.IssueDate.IsZero() {
		inv.TglTerbit = inv.IssueDate.Format("2006-01-02")
	}
	if !inv.DueDate.IsZero() {
		inv.TglJatuhTempo = inv.DueDate.Format("2006-01-02")
	}
	inv.TotalTagihan = inv.Subtotal
	inv.GrandTotalAlt = inv.GrandTotal

	// Cari pajak pertama jika ada item yang punya tax
	for _, it := range inv.Items {
		if it.TaxID != nil && *it.TaxID != "" {
			inv.Pajak = *it.TaxID
			if it.Tax != nil {
				inv.CatatanPajak = it.Tax.TaxName
			}
			break
		}
	}
	return nil
}

// TrxInvoiceItem merepresentasikan Detail Baris Penagihan (Dynamic Grid)
// (Tabel: trx_invoice_items)
type TrxInvoiceItem struct {
	ID              string         `gorm:"primaryKey;size:36" json:"id"`
	InvoiceID       string         `gorm:"size:36;index;not null" json:"invoice_id"` // FK -> trx_invoices.id
	ItemDescription string         `gorm:"size:255;not null" json:"item_description"` // Deskripsi Item
	Deskripsi       string         `gorm:"-" json:"deskripsi"`
	Quantity        int            `gorm:"not null;default:1" json:"quantity"`        // Kuantitas
	Qty             int            `gorm:"-" json:"qty"`
	UnitPrice       float64        `gorm:"type:decimal(15,2);not null;default:0" json:"unit_price"` // Harga Satuan
	Harga           float64        `gorm:"-" json:"harga"`
	TaxID           *string        `gorm:"size:36;index" json:"tax_id,omitempty"`     // FK -> mst_finance_taxes.id (Nullable)
	Tax             *MstFinanceTax `gorm:"foreignKey:TaxID" json:"tax,omitempty"`
	TotalPrice      float64        `gorm:"type:decimal(15,2);not null;default:0" json:"total_price"` // Quantity * UnitPrice
	Subtotal        float64        `gorm:"-" json:"subtotal"`
	CreatedAt       time.Time      `json:"created_at"`
}

func (TrxInvoiceItem) TableName() string {
	return "trx_invoice_items"
}

func (it *TrxInvoiceItem) AfterFind(tx *gorm.DB) (err error) {
	it.Deskripsi = it.ItemDescription
	it.Qty = it.Quantity
	it.Harga = it.UnitPrice
	it.Subtotal = it.TotalPrice
	return nil
}

package model

import "time"

// Account merepresentasikan akun dalam Chart of Accounts (COA) untuk Buku Besar (General Ledger).
type Account struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Code      string    `gorm:"size:20;uniqueIndex;not null" json:"code"`      // Cth: "1001", "1200", "2001", "4001"
	Name      string    `gorm:"size:100;not null" json:"name"`                 // Cth: "Kas & Bank", "Piutang Usaha (AR)"
	Type      string    `gorm:"size:30;not null" json:"type"`                  // "Asset", "Liability", "Equity", "Revenue", "Expense"
	Balance   float64   `gorm:"type:decimal(14,2);not null;default:0" json:"balance"` // Saldo terkini
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// JournalEntry merepresentasikan sebuah transaksi jurnal umum (General Ledger Entry) yang seimbang.
type JournalEntry struct {
	ID            uint          `gorm:"primaryKey" json:"id"`
	EntryNumber   string        `gorm:"size:50;uniqueIndex;not null" json:"entry_number"` // Cth: "JRN-202609-001"
	Date          time.Time     `gorm:"not null" json:"date"`
	Description   string        `gorm:"type:text" json:"description"`
	TransactionID *string       `gorm:"size:36;index" json:"transaction_id,omitempty"` // Referensi ke TransactionID (jika ada)
	TotalDebit    float64       `gorm:"type:decimal(14,2);not null" json:"total_debit"`
	TotalCredit   float64       `gorm:"type:decimal(14,2);not null" json:"total_credit"`
	Items         []JournalItem `gorm:"foreignKey:JournalEntryID;constraint:OnDelete:CASCADE" json:"items"`
	CreatedAt     time.Time     `json:"created_at"`
	UpdatedAt     time.Time     `json:"updated_at"`
}

// JournalItem merepresentasikan baris rincian debet/kredit dalam sebuah JournalEntry.
type JournalItem struct {
	ID             uint      `gorm:"primaryKey" json:"id"`
	JournalEntryID uint      `gorm:"index;not null" json:"journal_entry_id"`
	AccountCode    string    `gorm:"size:20;not null" json:"account_code"`
	AccountName    string    `gorm:"size:100;not null" json:"account_name"`
	Debit          float64   `gorm:"type:decimal(14,2);not null;default:0" json:"debit"`
	Credit         float64   `gorm:"type:decimal(14,2);not null;default:0" json:"credit"`
	Memo           string    `gorm:"size:255" json:"memo"`
	CreatedAt      time.Time `json:"created_at"`
}

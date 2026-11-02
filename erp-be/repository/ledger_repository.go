package repository

import (
	"errors"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// LedgerRepository mendefinisikan kontrak akses data untuk Buku Besar & Bagan Akun.
type LedgerRepository interface {
	FindAccounts() ([]model.Account, error)
	FindAccountByCode(code string) (*model.Account, error)
	CreateAccount(acc *model.Account) error
	UpdateAccountBalance(code string, delta float64) error
	FindJournalEntries(limit, offset int) ([]model.JournalEntry, int64, error)
	CreateJournalEntry(entry *model.JournalEntry) error
}

type ledgerRepository struct {
	db *gorm.DB
}

// NewLedgerRepository menginisialisasi implementasi LedgerRepository.
func NewLedgerRepository(db *gorm.DB) LedgerRepository {
	return &ledgerRepository{db: db}
}

func (r *ledgerRepository) FindAccounts() ([]model.Account, error) {
	var accounts []model.Account
	err := r.db.Order("code ASC").Find(&accounts).Error
	return accounts, err
}

func (r *ledgerRepository) FindAccountByCode(code string) (*model.Account, error) {
	var acc model.Account
	err := r.db.Where("code = ?", code).First(&acc).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &acc, nil
}

func (r *ledgerRepository) CreateAccount(acc *model.Account) error {
	return r.db.Create(acc).Error
}

func (r *ledgerRepository) UpdateAccountBalance(code string, delta float64) error {
	return r.db.Model(&model.Account{}).
		Where("code = ?", code).
		Update("balance", gorm.Expr("balance + ?", delta)).Error
}

func (r *ledgerRepository) FindJournalEntries(limit, offset int) ([]model.JournalEntry, int64, error) {
	var entries []model.JournalEntry
	var total int64

	if err := r.db.Model(&model.JournalEntry{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.Preload("Items").
		Order("date DESC, id DESC").
		Limit(limit).
		Offset(offset).
		Find(&entries).Error

	return entries, total, err
}

func (r *ledgerRepository) CreateJournalEntry(entry *model.JournalEntry) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(entry).Error; err != nil {
			return err
		}

		// Update saldo akun terkait
		for _, item := range entry.Items {
			// Asset & Expense: bertambah di Debit, berkurang di Credit
			// Liability, Equity, Revenue: bertambah di Credit, berkurang di Debit
			var acc model.Account
			if err := tx.Where("code = ?", item.AccountCode).First(&acc).Error; err == nil {
				var delta float64
				if acc.Type == "Asset" || acc.Type == "Expense" {
					delta = item.Debit - item.Credit
				} else {
					delta = item.Credit - item.Debit
				}
				_ = tx.Model(&model.Account{}).Where("code = ?", item.AccountCode).Update("balance", gorm.Expr("balance + ?", delta)).Error
			}
		}
		return nil
	})
}

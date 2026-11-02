package config

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"erp-kafeinart/model"

	"github.com/glebarez/sqlite"
	"github.com/google/uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Global database instance
var DB *gorm.DB

// InitDatabase membuka koneksi database SQLite via GORM dan menjalankan migrasi skema tabel.
func InitDatabase(cfg *Config) (*gorm.DB, error) {
	// Pastikan folder untuk file database sqlite ada
	dir := filepath.Dir(cfg.DBPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	// Buka koneksi menggunakan glebarez/sqlite (pure Go SQLite driver, tanpa butuh GCC di Windows)
	gormCfg := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Warn),
	}

	db, err := gorm.Open(sqlite.Open(cfg.DBPath+"?_pragma=busy_timeout(5000)&_pragma=journal_mode(WAL)"), gormCfg)
	if err != nil {
		return nil, err
	}

	// Tangani kompatibilitas jika tabel transactions lama sudah ada di SQLite
	if db.Migrator().HasTable("transactions") {
		if db.Migrator().HasColumn(&model.Transaction{}, "transaction_id") {
			var ids []uint
			_ = db.Raw("SELECT id FROM transactions WHERE transaction_id = '' OR transaction_id IS NULL").Scan(&ids)
			for _, id := range ids {
				_ = db.Exec("UPDATE transactions SET transaction_id = ? WHERE id = ?", uuid.NewString(), id)
			}
		}
		if db.Migrator().HasColumn(&model.Transaction{}, "reference_no") {
			var ids []uint
			_ = db.Raw("SELECT id FROM transactions WHERE reference_no = '' OR reference_no IS NULL").Scan(&ids)
			for _, id := range ids {
				_ = db.Exec("UPDATE transactions SET reference_no = ? WHERE id = ?", fmt.Sprintf("REF-%d-%s", id, uuid.NewString()[:8]), id)
			}
		}
	}

	// Migrasikan skema model ke database
	err = db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.StockAdjustment{},
		&model.Transaction{},
		&model.Employee{},
		&model.Account{},
		&model.JournalEntry{},
		&model.JournalItem{},
		// Tabel Spesifikasi Modul ERP Finance Kafeinarts:
		&model.MstFinanceCategory{},
		&model.MstFinanceAccount{},
		&model.MstFinanceTax{},
		&model.TrxCashJournal{},
		&model.TrxInvoice{},
		&model.TrxInvoiceItem{},
		&model.TrxVendorBill{},
		// Tabel Spesifikasi Modul ERP Human Resource (HR) Kafeinarts:
		&model.MstHrDepartment{},
		&model.MstHrPosition{},
		&model.MstHrEmployee{},
		&model.MstHrLeaveType{},
		&model.MstHrSalaryComponent{},
		&model.MstHrWorkShift{},
		&model.TrxHrAttendance{},
		&model.TrxHrLeaveRequest{},
		&model.TrxHrPayroll{},
		&model.TrxHrPayrollItem{},
		&model.TrxHrReimbursement{},
		&model.TrxHrPerformanceReview{},
		&model.TrxHrPerformanceReviewItem{},
		// Tabel Spesifikasi Modul ERP Supply Chain (SC) Kafeinarts:
		&model.MstScVendor{},
		&model.MstScItemCategory{},
		&model.MstScItem{},
		&model.MstScLocation{},
		&model.TrxScPurchaseRequest{},
		&model.TrxScPurchaseRequestItem{},
		&model.TrxScPurchaseOrder{},
		&model.TrxScPurchaseOrderItem{},
		&model.TrxScGoodsReceipt{},
		&model.TrxScGoodsReceiptItem{},
		&model.TrxScAsset{},
		&model.TrxScSubscription{},
	)
	if err != nil {
		return nil, err
	}

	DB = db
	log.Printf("[SUCCESS] Database SQLite terhubung di %s & AutoMigrate berhasil", cfg.DBPath)
	return DB, nil
}

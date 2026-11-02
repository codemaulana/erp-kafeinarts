package repository

import (
	"errors"
	"strings"
	"time"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// FinanceTransactionRepository mendefinisikan akses data untuk Jurnal Kas, Invoice AR, Tagihan AP, dan Laporan Keuangan.
type FinanceTransactionRepository interface {
	// Jurnal Kas (Pencatatan Harian)
	FindCashJournals(page, perPage int, accountID, txType, categoryID, search string) ([]model.TrxCashJournal, int64, error)
	FindCashJournalByID(id string) (*model.TrxCashJournal, error)
	CreateCashJournal(journal *model.TrxCashJournal) error
	UpdateCashJournal(oldJournal *model.TrxCashJournal, newJournal *model.TrxCashJournal) error
	DeleteCashJournal(id string) error

	// Manajemen Invoice / Piutang (AR)
	FindInvoices(page, perPage int, status, clientID, search string) ([]model.TrxInvoice, int64, error)
	FindInvoiceByID(id string) (*model.TrxInvoice, error)
	CreateInvoice(inv *model.TrxInvoice) error
	UpdateInvoice(inv *model.TrxInvoice) error
	UpdateInvoiceStatus(id, status string) (*model.TrxInvoice, error)
	DeleteInvoice(id string) error

	// Tagihan Vendor / Hutang (AP)
	FindVendorBills(page, perPage int, status, vendorID, search string) ([]model.TrxVendorBill, int64, error)
	FindVendorBillByID(id string) (*model.TrxVendorBill, error)
	CreateVendorBill(bill *model.TrxVendorBill) error
	UpdateVendorBill(bill *model.TrxVendorBill) error
	UpdateVendorBillStatus(id, status string) (*model.TrxVendorBill, error)
	DeleteVendorBill(id string) error

	// Laporan (Reports)
	GetCashFlowReport(start, end time.Time, accountID string) (*model.CashFlowReport, error)
	GetProfitLossReport(month, year int) (*model.ProfitLossReport, error)
	GetARAgingReport(status, clientName string) (*model.ARAgingReport, error)
}

type financeTransactionRepository struct {
	db *gorm.DB
}

func NewFinanceTransactionRepository(db *gorm.DB) FinanceTransactionRepository {
	return &financeTransactionRepository{db: db}
}

// ── Jurnal Kas (Pencatatan Harian) ──────────────────────────────────────────

func (r *financeTransactionRepository) FindCashJournals(page, perPage int, accountID, txType, categoryID, search string) ([]model.TrxCashJournal, int64, error) {
	var list []model.TrxCashJournal
	var total int64

	query := r.db.Model(&model.TrxCashJournal{})
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if txType != "" {
		query = query.Where("transaction_type = ?", strings.ToUpper(txType))
	}
	if categoryID != "" {
		query = query.Where("category_id = ?", categoryID)
	}
	if search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(transaction_number) LIKE ? OR LOWER(description) LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Preload("Account").Preload("Category").
		Order("transaction_date DESC, created_at DESC").
		Limit(perPage).Offset(offset).
		Find(&list).Error

	return list, total, err
}

func (r *financeTransactionRepository) FindCashJournalByID(id string) (*model.TrxCashJournal, error) {
	var j model.TrxCashJournal
	err := r.db.Preload("Account").Preload("Category").First(&j, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &j, nil
}

func (r *financeTransactionRepository) CreateCashJournal(journal *model.TrxCashJournal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(journal).Error; err != nil {
			return err
		}

		// Update saldo rekening terkait
		delta := journal.Amount
		if journal.TransactionType == "EXPENSE" {
			delta = -journal.Amount
		}

		return tx.Model(&model.MstFinanceAccount{}).
			Where("id = ?", journal.AccountID).
			Update("current_balance", gorm.Expr("current_balance + ?", delta)).Error
	})
}

func (r *financeTransactionRepository) UpdateCashJournal(oldJournal *model.TrxCashJournal, newJournal *model.TrxCashJournal) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// 1. Revert efek saldo transaksi lama
		oldDelta := -oldJournal.Amount
		if oldJournal.TransactionType == "EXPENSE" {
			oldDelta = oldJournal.Amount
		}
		if err := tx.Model(&model.MstFinanceAccount{}).
			Where("id = ?", oldJournal.AccountID).
			Update("current_balance", gorm.Expr("current_balance + ?", oldDelta)).Error; err != nil {
			return err
		}

		// 2. Terapkan efek saldo transaksi baru
		newDelta := newJournal.Amount
		if newJournal.TransactionType == "EXPENSE" {
			newDelta = -newJournal.Amount
		}
		if err := tx.Model(&model.MstFinanceAccount{}).
			Where("id = ?", newJournal.AccountID).
			Update("current_balance", gorm.Expr("current_balance + ?", newDelta)).Error; err != nil {
			return err
		}

		// 3. Simpan perubahan jurnal
		return tx.Save(newJournal).Error
	})
}

func (r *financeTransactionRepository) DeleteCashJournal(id string) error {
	var journal model.TrxCashJournal
	if err := r.db.First(&journal, "id = ?", id).Error; err != nil {
		return err
	}

	return r.db.Transaction(func(tx *gorm.DB) error {
		// Revert saldo rekening
		delta := -journal.Amount
		if journal.TransactionType == "EXPENSE" {
			delta = journal.Amount
		}
		if err := tx.Model(&model.MstFinanceAccount{}).
			Where("id = ?", journal.AccountID).
			Update("current_balance", gorm.Expr("current_balance + ?", delta)).Error; err != nil {
			return err
		}

		return tx.Delete(&journal).Error
	})
}

// ── Manajemen Invoice / Piutang (AR) ────────────────────────────────────────

func (r *financeTransactionRepository) FindInvoices(page, perPage int, status, clientID, search string) ([]model.TrxInvoice, int64, error) {
	var list []model.TrxInvoice
	var total int64

	query := r.db.Model(&model.TrxInvoice{})
	if status != "" {
		query = query.Where("status = ?", strings.ToUpper(status))
	}
	if clientID != "" {
		query = query.Where("client_id = ?", clientID)
	}
	if search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(invoice_number) LIKE ? OR LOWER(client_name) LIKE ? OR LOWER(project_name) LIKE ?", pattern, pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Preload("Items.Tax").
		Order("issue_date DESC, created_at DESC").
		Limit(perPage).Offset(offset).
		Find(&list).Error

	return list, total, err
}

func (r *financeTransactionRepository) FindInvoiceByID(id string) (*model.TrxInvoice, error) {
	var inv model.TrxInvoice
	err := r.db.Preload("Items.Tax").First(&inv, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &inv, nil
}

func (r *financeTransactionRepository) CreateInvoice(inv *model.TrxInvoice) error {
	return r.db.Create(inv).Error
}

func (r *financeTransactionRepository) UpdateInvoice(inv *model.TrxInvoice) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		// Hapus item lama, gantikan dengan item baru
		if err := tx.Where("invoice_id = ?", inv.ID).Delete(&model.TrxInvoiceItem{}).Error; err != nil {
			return err
		}
		return tx.Save(inv).Error
	})
}

func (r *financeTransactionRepository) UpdateInvoiceStatus(id, status string) (*model.TrxInvoice, error) {
	var inv model.TrxInvoice
	if err := r.db.Preload("Items.Tax").First(&inv, "id = ?", id).Error; err != nil {
		return nil, err
	}

	inv.Status = strings.ToUpper(status)
	if inv.Status == "PAID" && inv.PaidAt == nil {
		now := time.Now()
		inv.PaidAt = &now
	}

	if err := r.db.Save(&inv).Error; err != nil {
		return nil, err
	}
	return &inv, nil
}

func (r *financeTransactionRepository) DeleteInvoice(id string) error {
	return r.db.Delete(&model.TrxInvoice{}, "id = ?", id).Error
}

// ── Tagihan Vendor / Hutang (AP) ────────────────────────────────────────────

func (r *financeTransactionRepository) FindVendorBills(page, perPage int, status, vendorID, search string) ([]model.TrxVendorBill, int64, error) {
	var list []model.TrxVendorBill
	var total int64

	query := r.db.Model(&model.TrxVendorBill{})
	if status != "" {
		query = query.Where("status = ?", strings.ToUpper(status))
	}
	if vendorID != "" {
		query = query.Where("vendor_id = ?", vendorID)
	}
	if search != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(bill_number) LIKE ? OR LOWER(vendor_name) LIKE ?", pattern, pattern)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Preload("Category").
		Order("bill_date DESC, created_at DESC").
		Limit(perPage).Offset(offset).
		Find(&list).Error

	return list, total, err
}

func (r *financeTransactionRepository) FindVendorBillByID(id string) (*model.TrxVendorBill, error) {
	var bill model.TrxVendorBill
	err := r.db.Preload("Category").First(&bill, "id = ?", id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &bill, nil
}

func (r *financeTransactionRepository) CreateVendorBill(bill *model.TrxVendorBill) error {
	return r.db.Create(bill).Error
}

func (r *financeTransactionRepository) UpdateVendorBill(bill *model.TrxVendorBill) error {
	return r.db.Save(bill).Error
}

func (r *financeTransactionRepository) UpdateVendorBillStatus(id, status string) (*model.TrxVendorBill, error) {
	var bill model.TrxVendorBill
	if err := r.db.Preload("Category").First(&bill, "id = ?", id).Error; err != nil {
		return nil, err
	}

	bill.Status = strings.ToUpper(status)
	if bill.Status == "PAID" && bill.PaidAt == nil {
		now := time.Now()
		bill.PaidAt = &now
	}

	if err := r.db.Save(&bill).Error; err != nil {
		return nil, err
	}
	return &bill, nil
}

func (r *financeTransactionRepository) DeleteVendorBill(id string) error {
	return r.db.Delete(&model.TrxVendorBill{}, "id = ?", id).Error
}

// ── Laporan (Reports) ───────────────────────────────────────────────────────

// 1. Laporan Arus Kas (Cash Flow)
// Kolom: Tanggal, No. Referensi, Keterangan, Kategori, Masuk, Keluar, Saldo Akhir
func (r *financeTransactionRepository) GetCashFlowReport(start, end time.Time, accountID string) (*model.CashFlowReport, error) {
	var accountName string
	if accountID != "" {
		var acc model.MstFinanceAccount
		if err := r.db.First(&acc, "id = ?", accountID).Error; err == nil {
			accountName = acc.BankName + " (" + acc.AccountNumber + ")"
		}
	}

	// Hitung saldo awal kas sebelum start
	var beginningCash float64
	var pastIn, pastOut float64
	pastInQuery := r.db.Model(&model.TrxCashJournal{}).Where("transaction_type = 'INCOME' AND transaction_date < ?", start)
	pastOutQuery := r.db.Model(&model.TrxCashJournal{}).Where("transaction_type = 'EXPENSE' AND transaction_date < ?", start)
	if accountID != "" {
		pastInQuery = pastInQuery.Where("account_id = ?", accountID)
		pastOutQuery = pastOutQuery.Where("account_id = ?", accountID)
	}
	pastInQuery.Select("COALESCE(SUM(amount), 0)").Scan(&pastIn)
	pastOutQuery.Select("COALESCE(SUM(amount), 0)").Scan(&pastOut)
	beginningCash = pastIn - pastOut

	// Ambil semua transaksi jurnal kas pada periode
	var journals []model.TrxCashJournal
	query := r.db.Model(&model.TrxCashJournal{}).
		Preload("Category").
		Preload("Account").
		Where("transaction_date BETWEEN ? AND ?", start, end)
	if accountID != "" {
		query = query.Where("account_id = ?", accountID)
	}
	if err := query.Order("transaction_date ASC, created_at ASC").Find(&journals).Error; err != nil {
		return nil, err
	}

	rows := make([]model.CashFlowReportRow, 0)
	var totalIn, totalOut float64
	runningBalance := beginningCash

	for _, j := range journals {
		catName := "Umum"
		if j.Category != nil {
			catName = j.Category.CategoryName
		}
		bankName := accountName
		if j.Account != nil && j.Account.BankName != "" {
			bankName = j.Account.BankName
		} else if bankName == "" {
			bankName = "Kas/Bank"
		}

		var inAmount, outAmount, nominal float64
		var tipe string
		if strings.EqualFold(j.TransactionType, "INCOME") || j.TransactionType == "Masuk" {
			inAmount = j.Amount
			nominal = j.Amount
			tipe = "Masuk"
			totalIn += j.Amount
			runningBalance += j.Amount
		} else {
			outAmount = j.Amount
			nominal = j.Amount
			tipe = "Keluar"
			totalOut += j.Amount
			runningBalance -= j.Amount
		}

		rows = append(rows, model.CashFlowReportRow{
			ID:            j.ID,
			Date:          j.TransactionDate,
			Tanggal:       j.TransactionDate.Format("2006-01-02"),
			ReferenceNo:   j.TransactionNumber,
			Description:   j.Description,
			Keterangan:    j.Description,
			Category:      catName,
			NamaKategori:  catName,
			Kategori:      j.CategoryID,
			AccountID:     j.AccountID,
			AkunBank:      j.AccountID,
			AccountName:   bankName,
			NamaBank:      bankName,
			Inflow:        inAmount,
			Outflow:       outAmount,
			Nominal:       nominal,
			Tipe:          tipe,
			EndingBalance: runningBalance,
		})
	}

	return &model.CashFlowReport{
		StartDate:     start,
		EndDate:       end,
		AccountID:     accountID,
		AccountName:   accountName,
		BeginningCash: beginningCash,
		TotalInflow:   totalIn,
		TotalOutflow:  totalOut,
		NetCashFlow:   totalIn - totalOut,
		EndingCash:    runningBalance,
		Rows:          rows,
		Summary: model.CashFlowSummary{
			TotalMasuk:  totalIn,
			TotalKeluar: totalOut,
			SaldoBersih: totalIn - totalOut,
		},
	}, nil
}

// 2. Laporan Laba/Rugi (Profit & Loss)
// Filter: Bulan, Tahun. Output: Ringkasan bertingkat (Total Pendapatan, Total Pengeluaran, Laba/Rugi Bersih).
func (r *financeTransactionRepository) GetProfitLossReport(month, year int) (*model.ProfitLossReport, error) {
	if year <= 0 {
		year = time.Now().Year()
	}

	var start, end time.Time
	if month > 0 && month <= 12 {
		start = time.Date(year, time.Month(month), 1, 0, 0, 0, 0, time.Local)
		end = start.AddDate(0, 1, 0).Add(-time.Nanosecond)
	} else {
		// Setahun penuh
		start = time.Date(year, 1, 1, 0, 0, 0, 0, time.Local)
		end = time.Date(year, 12, 31, 23, 59, 59, 0, time.Local)
	}

	var revenueRows []struct {
		CategoryName string  `gorm:"column:category_name"`
		TotalAmount  float64 `gorm:"column:total_amount"`
	}
	var expenseRows []struct {
		CategoryName string  `gorm:"column:category_name"`
		TotalAmount  float64 `gorm:"column:total_amount"`
	}

	// Pendapatan dari TrxCashJournal (INCOME)
	r.db.Table("trx_cash_journals").
		Select("COALESCE(mst_finance_categories.category_name, 'Pendapatan Lainnya') as category_name, SUM(trx_cash_journals.amount) as total_amount").
		Joins("LEFT JOIN mst_finance_categories ON mst_finance_categories.id = trx_cash_journals.category_id").
		Where("(trx_cash_journals.transaction_type = 'INCOME' OR trx_cash_journals.transaction_type = 'Masuk') AND trx_cash_journals.transaction_date BETWEEN ? AND ?", start, end).
		Group("category_name").
		Scan(&revenueRows)

	// Pengeluaran dari TrxCashJournal (EXPENSE)
	r.db.Table("trx_cash_journals").
		Select("COALESCE(mst_finance_categories.category_name, 'Beban Operasional') as category_name, SUM(trx_cash_journals.amount) as total_amount").
		Joins("LEFT JOIN mst_finance_categories ON mst_finance_categories.id = trx_cash_journals.category_id").
		Where("(trx_cash_journals.transaction_type = 'EXPENSE' OR trx_cash_journals.transaction_type = 'Keluar') AND trx_cash_journals.transaction_date BETWEEN ? AND ?", start, end).
		Group("category_name").
		Scan(&expenseRows)

	revDetails := make([]model.CategoryBreakdown, 0)
	expDetails := make([]model.CategoryBreakdown, 0)
	var totalRev, totalExp float64

	for _, row := range revenueRows {
		revDetails = append(revDetails, model.CategoryBreakdown{
			Category: row.CategoryName,
			Kategori: row.CategoryName,
			Amount:   row.TotalAmount,
			Total:    row.TotalAmount,
		})
		totalRev += row.TotalAmount
	}

	for _, row := range expenseRows {
		expDetails = append(expDetails, model.CategoryBreakdown{
			Category: row.CategoryName,
			Kategori: row.CategoryName,
			Amount:   row.TotalAmount,
			Total:    row.TotalAmount,
		})
		totalExp += row.TotalAmount
	}

	return &model.ProfitLossReport{
		Month:          month,
		Year:           year,
		TotalRevenue:   totalRev,
		TotalExpense:   totalExp,
		NetProfitLoss:  totalRev - totalExp,
		RevenueDetails: revDetails,
		ExpenseDetails: expDetails,
		Pendapatan:     revDetails,
		Pengeluaran:    expDetails,
		Summary: model.ProfitLossSummary{
			TotalPendapatan:  totalRev,
			TotalPengeluaran: totalExp,
			LabaKotor:        totalRev - totalExp,
			Bulan:            month,
			Tahun:            year,
		},
	}, nil
}

// 3. Laporan Piutang (AR Aging)
// Kolom: No. Invoice, Klien, Jatuh Tempo, Keterlambatan (Hari), Total Tagihan, Sisa Belum Dibayar
func (r *financeTransactionRepository) GetARAgingReport(status, clientName string) (*model.ARAgingReport, error) {
	var invoices []model.TrxInvoice
	query := r.db.Model(&model.TrxInvoice{})
	if status != "" {
		st := strings.ToUpper(status)
		if st == "LUNAS" {
			query = query.Where("status = 'PAID'")
		} else if st == "SEBAGIAN" {
			query = query.Where("status = 'PARTIAL'")
		} else if st == "BELUM LUNAS" {
			query = query.Where("status != 'PAID' AND status != 'CANCELLED'")
		} else {
			query = query.Where("status = ?", st)
		}
	}
	if clientName != "" {
		pattern := "%" + strings.ToLower(clientName) + "%"
		query = query.Where("LOWER(client_name) LIKE ?", pattern)
	}

	if err := query.Order("due_date ASC").Find(&invoices).Error; err != nil {
		return nil, err
	}

	now := time.Now()
	items := make([]model.ARAgingReportItem, 0)
	var totalReceivable, totalOverdue float64
	var totalTagihan, totalLunas, totalBelumLunas float64

	for _, inv := range invoices {
		daysOverdue := 0
		if now.After(inv.DueDate) && inv.Status != "PAID" {
			diff := now.Sub(inv.DueDate)
			daysOverdue = int(diff.Hours() / 24)
			totalOverdue += inv.GrandTotal
		}

		remainingDue := inv.GrandTotal
		var statusFE string
		switch strings.ToUpper(inv.Status) {
		case "PAID":
			remainingDue = 0
			statusFE = "Lunas"
		case "PARTIAL":
			remainingDue = inv.GrandTotal * 0.5
			statusFE = "Sebagian"
		case "OVERDUE":
			statusFE = "Jatuh Tempo"
		case "CANCELLED":
			statusFE = "Batal"
		default:
			statusFE = "Belum Lunas"
		}

		totalReceivable += remainingDue
		totalTagihan += inv.GrandTotal
		if inv.Status == "PAID" {
			totalLunas += inv.GrandTotal
		} else {
			totalBelumLunas += remainingDue
		}

		items = append(items, model.ARAgingReportItem{
			ID:            inv.ID,
			InvoiceID:     inv.ID,
			InvoiceNumber: inv.InvoiceNumber,
			NoInvoice:     inv.InvoiceNumber,
			ClientName:    inv.ClientName,
			Klien:         inv.ClientName,
			ProjectName:   inv.ProjectName,
			Proyek:        inv.ProjectName,
			IssueDate:     inv.IssueDate,
			TglTerbit:     inv.IssueDate.Format("2006-01-02"),
			DueDate:       inv.DueDate,
			TglJatuhTempo: inv.DueDate.Format("2006-01-02"),
			DaysOverdue:   daysOverdue,
			GrandTotal:    inv.GrandTotal,
			GrandTotalAlt: inv.GrandTotal,
			RemainingDue:  remainingDue,
			Status:        statusFE,
		})
	}

	return &model.ARAgingReport{
		AsOfDate:          now,
		TotalReceivable:   totalReceivable,
		TotalOverdue:      totalOverdue,
		TotalInvoiceCount: len(items),
		Items:             items,
		Rows:              items,
		Summary: model.ARAgingSummary{
			TotalTagihan:    totalTagihan,
			TotalLunas:      totalLunas,
			TotalBelumLunas: totalBelumLunas,
		},
	}, nil
}

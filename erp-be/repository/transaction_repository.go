package repository

import (
	"errors"
	"strings"
	"time"

	"erp-kafeinart/model"

	"gorm.io/gorm"
)

// TransactionRepository mendefinisikan kontrak akses data transaksi keuangan (AR, AP, dan Umum).
type TransactionRepository interface {
	FindAll(page, perPage int, txType, category, status, currency, search string) ([]model.Transaction, int64, error)
	FindByID(id uint) (*model.Transaction, error)
	FindByTransactionID(uuid string) (*model.Transaction, error)
	FindInvoices(page, perPage int, status, search string) ([]model.Transaction, int64, error)
	FindBills(page, perPage int, status, search string) ([]model.Transaction, int64, error)
	Create(tx *model.Transaction) error
	Update(tx *model.Transaction) error
	Delete(id uint) error
	GetSummary() (*model.FinanceSummary, error)
	GetIncomeStatement(start, end time.Time) (*model.IncomeStatement, error)
	GetBalanceSheet(asOf time.Time) (*model.BalanceSheet, error)
	GetCashFlow(start, end time.Time) (*model.CashFlowReport, error)
}

type transactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository menginisialisasi implementasi TransactionRepository.
func NewTransactionRepository(db *gorm.DB) TransactionRepository {
	return &transactionRepository{db: db}
}

func (r *transactionRepository) applyFilters(query *gorm.DB, txType, category, status, currency, search string) *gorm.DB {
	if strings.TrimSpace(txType) != "" {
		// Dukung "Debit", "Credit", maupun "income", "expense"
		t := strings.ToLower(txType)
		if t == "debit" || t == "income" {
			query = query.Where("LOWER(type) IN (?, ?)", "debit", "income")
		} else if t == "credit" || t == "expense" {
			query = query.Where("LOWER(type) IN (?, ?)", "credit", "expense")
		} else {
			query = query.Where("type = ?", txType)
		}
	}

	if strings.TrimSpace(category) != "" {
		query = query.Where("category = ?", category)
	}

	if strings.TrimSpace(status) != "" {
		query = query.Where("LOWER(status) = ?", strings.ToLower(status))
	}

	if strings.TrimSpace(currency) != "" {
		query = query.Where("currency = ?", currency)
	}

	if strings.TrimSpace(search) != "" {
		pattern := "%" + strings.ToLower(search) + "%"
		query = query.Where("LOWER(reference_no) LIKE ? OR LOWER(description) LIKE ? OR LOWER(client_vendor_name) LIKE ? OR LOWER(category) LIKE ?",
			pattern, pattern, pattern, pattern)
	}

	return query
}

func (r *transactionRepository) FindAll(page, perPage int, txType, category, status, currency, search string) ([]model.Transaction, int64, error) {
	var list []model.Transaction
	var total int64

	query := r.db.Model(&model.Transaction{})
	query = r.applyFilters(query, txType, category, status, currency, search)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	offset := (page - 1) * perPage
	err := query.Order("transaction_date DESC, id DESC").Limit(perPage).Offset(offset).Find(&list).Error
	if err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *transactionRepository) FindInvoices(page, perPage int, status, search string) ([]model.Transaction, int64, error) {
	return r.FindAll(page, perPage, "Debit", "", status, "", search)
}

func (r *transactionRepository) FindBills(page, perPage int, status, search string) ([]model.Transaction, int64, error) {
	return r.FindAll(page, perPage, "Credit", "", status, "", search)
}

func (r *transactionRepository) FindByID(id uint) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.First(&tx, id).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) FindByTransactionID(uuid string) (*model.Transaction, error) {
	var tx model.Transaction
	err := r.db.Where("transaction_id = ?", uuid).First(&tx).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &tx, nil
}

func (r *transactionRepository) Create(tx *model.Transaction) error {
	return r.db.Create(tx).Error
}

func (r *transactionRepository) Update(tx *model.Transaction) error {
	return r.db.Save(tx).Error
}

func (r *transactionRepository) Delete(id uint) error {
	return r.db.Delete(&model.Transaction{}, id).Error
}

func (r *transactionRepository) GetSummary() (*model.FinanceSummary, error) {
	var totalIncome, totalExpense, totalReceivable, totalPayable float64
	var count int64

	r.db.Model(&model.Transaction{}).Count(&count)

	// Total Pemasukan yang lunas
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) = ?", "debit", "income", "paid").
		Select("COALESCE(SUM(total_amount), SUM(amount), 0)").
		Scan(&totalIncome)

	// Total Pengeluaran yang lunas
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) = ?", "credit", "expense", "paid").
		Select("COALESCE(SUM(total_amount), SUM(amount), 0)").
		Scan(&totalExpense)

	// Total Piutang belum lunas (AR: status Draft/Sent/Overdue)
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) != ?", "debit", "income", "paid").
		Select("COALESCE(SUM(total_amount), SUM(amount), 0)").
		Scan(&totalReceivable)

	// Total Hutang belum dibayar (AP: status Draft/Sent/Overdue)
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) != ?", "credit", "expense", "paid").
		Select("COALESCE(SUM(total_amount), SUM(amount), 0)").
		Scan(&totalPayable)

	return &model.FinanceSummary{
		TotalIncome:      totalIncome,
		TotalExpense:     totalExpense,
		NetBalance:       totalIncome - totalExpense,
		TotalReceivable:  totalReceivable,
		TotalPayable:     totalPayable,
		TransactionCount: count,
	}, nil
}

func (r *transactionRepository) GetIncomeStatement(start, end time.Time) (*model.IncomeStatement, error) {
	var revenues []model.CategoryBreakdown
	var expenses []model.CategoryBreakdown
	var totalRevenue, totalExpense, totalTax float64

	// Rincian Revenue per Kategori
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND transaction_date BETWEEN ? AND ?", "debit", "income", start, end).
		Select("category, COALESCE(SUM(amount), 0) as amount").
		Group("category").
		Scan(&revenues)

	for _, rev := range revenues {
		totalRevenue += rev.Amount
	}

	// Rincian Expense per Kategori
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND transaction_date BETWEEN ? AND ?", "credit", "expense", start, end).
		Select("category, COALESCE(SUM(amount), 0) as amount").
		Group("category").
		Scan(&expenses)

	for _, exp := range expenses {
		totalExpense += exp.Amount
	}

	// Total Pajak pada periode tersebut
	r.db.Model(&model.Transaction{}).
		Where("transaction_date BETWEEN ? AND ?", start, end).
		Select("COALESCE(SUM(tax_amount), 0)").
		Scan(&totalTax)

	netIncome := totalRevenue - totalExpense - totalTax

	return &model.IncomeStatement{
		PeriodStart:  start,
		PeriodEnd:    end,
		TotalRevenue: totalRevenue,
		TotalExpense: totalExpense,
		TotalTax:     totalTax,
		NetIncome:    netIncome,
		Revenues:     revenues,
		Expenses:     expenses,
	}, nil
}

func (r *transactionRepository) GetBalanceSheet(asOf time.Time) (*model.BalanceSheet, error) {
	var accounts []model.Account
	if err := r.db.Find(&accounts).Error; err != nil {
		return nil, err
	}

	var assets []model.BalanceSheetItem
	var liabilities []model.BalanceSheetItem
	var equity []model.BalanceSheetItem

	var totalAssets, totalLiabilities, totalEquity float64

	for _, acc := range accounts {
		item := model.BalanceSheetItem{
			Code:   acc.Code,
			Name:   acc.Name,
			Amount: acc.Balance,
		}

		switch acc.Type {
		case "Asset":
			assets = append(assets, item)
			totalAssets += acc.Balance
		case "Liability":
			liabilities = append(liabilities, item)
			totalLiabilities += acc.Balance
		case "Equity":
			equity = append(equity, item)
			totalEquity += acc.Balance
		}
	}

	// Persamaan akuntansi seimbang jika TotalAssets == TotalLiabilities + TotalEquity
	isBalanced := (totalAssets == (totalLiabilities + totalEquity))

	return &model.BalanceSheet{
		AsOfDate:         asOf,
		Assets:           assets,
		TotalAssets:      totalAssets,
		Liabilities:      liabilities,
		TotalLiabilities: totalLiabilities,
		Equity:           equity,
		TotalEquity:      totalEquity,
		IsBalanced:       isBalanced,
	}, nil
}

func (r *transactionRepository) GetCashFlow(start, end time.Time) (*model.CashFlowReport, error) {
	var inflows []model.CashFlowItem
	var outflows []model.CashFlowItem
	var totalInflows, totalOutflows float64

	// Arus Kas Masuk (Paid Invoices/Sales)
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) = ? AND transaction_date BETWEEN ? AND ?", "debit", "income", "paid", start, end).
		Select("category, COALESCE(SUM(total_amount), SUM(amount), 0) as amount").
		Group("category").
		Scan(&inflows)

	for _, inf := range inflows {
		totalInflows += inf.Amount
	}

	// Arus Kas Keluar (Paid Bills/Expenses)
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) = ? AND transaction_date BETWEEN ? AND ?", "credit", "expense", "paid", start, end).
		Select("category, COALESCE(SUM(total_amount), SUM(amount), 0) as amount").
		Group("category").
		Scan(&outflows)

	for _, out := range outflows {
		totalOutflows += out.Amount
	}

	var beginningCash float64
	// Saldo kas sebelum periode start
	var pastIn, pastOut float64
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) = ? AND transaction_date < ?", "debit", "income", "paid", start).
		Select("COALESCE(SUM(total_amount), SUM(amount), 0)").Scan(&pastIn)
	r.db.Model(&model.Transaction{}).
		Where("LOWER(type) IN (?, ?) AND LOWER(status) = ? AND transaction_date < ?", "credit", "expense", "paid", start).
		Select("COALESCE(SUM(total_amount), SUM(amount), 0)").Scan(&pastOut)
	beginningCash = pastIn - pastOut

	netCashFlow := totalInflows - totalOutflows
	endingCash := beginningCash + netCashFlow

	return &model.CashFlowReport{
		PeriodStart:      start,
		PeriodEnd:        end,
		BeginningCash:    beginningCash,
		OperatingInflows: inflows,
		TotalInflows:     totalInflows,
		OperatingOutflow: outflows,
		TotalOutflows:    totalOutflows,
		NetCashFlow:      netCashFlow,
		EndingCash:       endingCash,
	}, nil
}

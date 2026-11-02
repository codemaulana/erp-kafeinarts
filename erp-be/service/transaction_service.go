package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"
	"erp-kafeinart/utils"

	"github.com/google/uuid"
)

// CreateTransactionDTO memuat data saat membuat transaksi baru.
type CreateTransactionDTO struct {
	ReferenceNo      string     `json:"reference_no"`
	Type             string     `json:"type"` // "Debit" (Pemasukan/Piutang) / "Credit" (Pengeluaran/Hutang)
	Category         string     `json:"category"`
	Amount           float64    `json:"amount"`
	TaxAmount        float64    `json:"tax_amount"`
	Currency         string     `json:"currency"`
	Status           string     `json:"status"` // "Draft", "Sent", "Paid", "Overdue"
	DueDate          *time.Time `json:"due_date"`
	TransactionDate  time.Time  `json:"transaction_date"`
	ClientVendorName string     `json:"client_vendor_name"`
	Description      string     `json:"description"`
}

// CreateInvoiceDTO memuat input khusus pembuatan invoice klien (Accounts Receivable).
type CreateInvoiceDTO struct {
	ReferenceNo     string     `json:"reference_no"` // Cth: INV-2026-001 (auto jika kosong)
	ClientName      string     `json:"client_name"`  // Nama klien/perusahaan
	Amount          float64    `json:"amount"`       // Nilai termin proyek
	TaxAmount       float64    `json:"tax_amount"`   // PPN/PPh
	Currency        string     `json:"currency"`     // IDR, USD
	DueDate         *time.Time `json:"due_date"`     // Batas waktu pembayaran
	TransactionDate time.Time  `json:"transaction_date"`
	Category        string     `json:"category"`    // Default: "Project Revenue"
	Description     string     `json:"description"` // Cth: "Termin 1 (50%) Website Redesign"
}

// CreateBillDTO memuat input khusus pencatatan tagihan vendor pihak ketiga (Accounts Payable).
type CreateBillDTO struct {
	ReferenceNo     string     `json:"reference_no"` // Cth: BILL-2026-001
	VendorName      string     `json:"vendor_name"`  // Cth: "AWS", "Google Cloud", "JetBrains", "Adobe"
	Amount          float64    `json:"amount"`       // Nominal tagihan
	TaxAmount       float64    `json:"tax_amount"`
	Currency        string     `json:"currency"`     // IDR, USD
	DueDate         *time.Time `json:"due_date"`
	TransactionDate time.Time  `json:"transaction_date"`
	Category        string     `json:"category"`    // Cth: "Cloud Infrastructure", "Software License"
	Description     string     `json:"description"` // Cth: "Biaya Server Produksi Agustus"
}

// UpdateStatusDTO untuk update status transaksi secara langsung.
type UpdateStatusDTO struct {
	Status string `json:"status"` // "Draft", "Sent", "Paid", "Overdue"
}

// CreateJournalEntryDTO untuk input manual jurnal umum.
type CreateJournalItemDTO struct {
	AccountCode string  `json:"account_code"`
	Debit       float64 `json:"debit"`
	Credit      float64 `json:"credit"`
	Memo        string  `json:"memo"`
}

type CreateJournalEntryDTO struct {
	EntryNumber string                 `json:"entry_number"`
	Date        time.Time              `json:"date"`
	Description string                 `json:"description"`
	Items       []CreateJournalItemDTO `json:"items"`
}

// TransactionService mendefinisikan bisnis logika keuangan studio.
type TransactionService interface {
	GetAll(page, perPage int, txType, category, status, currency, search string) ([]model.Transaction, utils.PaginationMeta, error)
	GetByIDOrUUID(identifier string) (*model.Transaction, error)
	Create(dto CreateTransactionDTO) (*model.Transaction, error)
	Update(id uint, dto CreateTransactionDTO) (*model.Transaction, error)
	UpdateStatus(identifier string, status string) (*model.Transaction, error)
	Delete(id uint) error

	// Accounts Receivable (AR) & Invoicing
	GetInvoices(page, perPage int, status, search string) ([]model.Transaction, utils.PaginationMeta, error)
	CreateInvoice(dto CreateInvoiceDTO) (*model.Transaction, error)
	PayInvoice(identifier string) (*model.Transaction, error)

	// Accounts Payable (AP) & Vendor Bills
	GetBills(page, perPage int, status, search string) ([]model.Transaction, utils.PaginationMeta, error)
	CreateBill(dto CreateBillDTO) (*model.Transaction, error)
	PayBill(identifier string) (*model.Transaction, error)

	// General Ledger & Bookkeeping
	GetAccounts() ([]model.Account, error)
	GetJournalEntries(page, perPage int) ([]model.JournalEntry, utils.PaginationMeta, error)
	CreateJournalEntry(dto CreateJournalEntryDTO) (*model.JournalEntry, error)

	// Financial Reporting
	GetSummary() (*model.FinanceSummary, error)
	GetIncomeStatement(start, end time.Time) (*model.IncomeStatement, error)
	GetBalanceSheet(asOf time.Time) (*model.BalanceSheet, error)
	GetCashFlow(start, end time.Time) (*model.CashFlowReport, error)
}

type transactionService struct {
	txRepo     repository.TransactionRepository
	ledgerRepo repository.LedgerRepository
}

// NewTransactionService menginisialisasi TransactionService.
func NewTransactionService(txRepo repository.TransactionRepository, ledgerRepo repository.LedgerRepository) TransactionService {
	return &transactionService{
		txRepo:     txRepo,
		ledgerRepo: ledgerRepo,
	}
}

// normalizeType mengonversi input type menjadi standar "Debit" atau "Credit"
func normalizeType(t string) string {
	switch strings.ToLower(strings.TrimSpace(t)) {
	case "debit", "income":
		return "Debit"
	case "credit", "expense":
		return "Credit"
	default:
		return "Debit"
	}
}

// normalizeStatus mengonversi input status menjadi standar kapitalisasi "Draft", "Sent", "Paid", "Overdue"
func normalizeStatus(s string) string {
	switch strings.ToLower(strings.TrimSpace(s)) {
	case "paid":
		return "Paid"
	case "sent":
		return "Sent"
	case "overdue":
		return "Overdue"
	default:
		return "Draft"
	}
}

func (s *transactionService) GetAll(page, perPage int, txType, category, status, currency, search string) ([]model.Transaction, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.txRepo.FindAll(page, perPage, txType, category, status, currency, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *transactionService) GetByIDOrUUID(identifier string) (*model.Transaction, error) {
	// Coba cari via UUID dahulu
	tx, err := s.txRepo.FindByTransactionID(identifier)
	if err == nil && tx != nil {
		return tx, nil
	}

	// Coba cari via integer ID jika numeric
	var id uint
	if _, errScan := fmt.Sscanf(identifier, "%d", &id); errScan == nil && id > 0 {
		return s.txRepo.FindByID(id)
	}

	return nil, errors.New("transaksi tidak ditemukan")
}

func (s *transactionService) Create(dto CreateTransactionDTO) (*model.Transaction, error) {
	if dto.Amount <= 0 {
		return nil, errors.New("jumlah nominal transaksi (amount) harus lebih besar dari 0")
	}

	category := strings.TrimSpace(dto.Category)
	if category == "" {
		return nil, errors.New("kategori transaksi wajib diisi")
	}

	txType := normalizeType(dto.Type)
	status := normalizeStatus(dto.Status)

	refNo := strings.TrimSpace(dto.ReferenceNo)
	if refNo == "" {
		prefix := "TRX"
		if txType == "Debit" {
			prefix = "INV"
		} else {
			prefix = "BILL"
		}
		refNo = fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102"), uuid.NewString()[:6])
	}

	txDate := dto.TransactionDate
	if txDate.IsZero() {
		txDate = time.Now()
	}

	currency := strings.ToUpper(strings.TrimSpace(dto.Currency))
	if currency == "" {
		currency = "IDR"
	}

	totalAmount := dto.Amount + dto.TaxAmount

	tx := &model.Transaction{
		TransactionID:    uuid.NewString(),
		ReferenceNo:      refNo,
		Type:             txType,
		Category:         category,
		Amount:           dto.Amount,
		TaxAmount:        dto.TaxAmount,
		TotalAmount:      totalAmount,
		Currency:         currency,
		Status:           status,
		DueDate:          dto.DueDate,
		TransactionDate:  txDate,
		ClientVendorName: strings.TrimSpace(dto.ClientVendorName),
		Description:      dto.Description,
	}

	if status == "Paid" {
		now := time.Now()
		tx.PaidAt = &now
	}

	if err := s.txRepo.Create(tx); err != nil {
		return nil, err
	}

	// Jika status Paid, buatkan entri jurnal otomatis ke Buku Besar
	if status == "Paid" {
		_ = s.createJournalForTransaction(tx)
	}

	return tx, nil
}

func (s *transactionService) Update(id uint, dto CreateTransactionDTO) (*model.Transaction, error) {
	tx, err := s.txRepo.FindByID(id)
	if err != nil || tx == nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	if dto.Amount <= 0 {
		return nil, errors.New("jumlah nominal transaksi (amount) harus lebih besar dari 0")
	}

	tx.Type = normalizeType(dto.Type)
	tx.Category = strings.TrimSpace(dto.Category)
	tx.Amount = dto.Amount
	tx.TaxAmount = dto.TaxAmount
	tx.TotalAmount = dto.Amount + dto.TaxAmount
	if strings.TrimSpace(dto.Currency) != "" {
		tx.Currency = strings.ToUpper(strings.TrimSpace(dto.Currency))
	}
	tx.Status = normalizeStatus(dto.Status)
	tx.DueDate = dto.DueDate
	if !dto.TransactionDate.IsZero() {
		tx.TransactionDate = dto.TransactionDate
	}
	if strings.TrimSpace(dto.ClientVendorName) != "" {
		tx.ClientVendorName = strings.TrimSpace(dto.ClientVendorName)
	}
	tx.Description = dto.Description

	if tx.Status == "Paid" && tx.PaidAt == nil {
		now := time.Now()
		tx.PaidAt = &now
	}

	if err := s.txRepo.Update(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) UpdateStatus(identifier string, status string) (*model.Transaction, error) {
	tx, err := s.GetByIDOrUUID(identifier)
	if err != nil || tx == nil {
		return nil, errors.New("transaksi tidak ditemukan")
	}

	normStatus := normalizeStatus(status)
	tx.Status = normStatus

	if normStatus == "Paid" && tx.PaidAt == nil {
		now := time.Now()
		tx.PaidAt = &now
		// Buat jurnal pelunasan
		_ = s.createJournalForTransaction(tx)
	}

	if err := s.txRepo.Update(tx); err != nil {
		return nil, err
	}

	return tx, nil
}

func (s *transactionService) Delete(id uint) error {
	tx, err := s.txRepo.FindByID(id)
	if err != nil || tx == nil {
		return errors.New("transaksi tidak ditemukan")
	}
	return s.txRepo.Delete(id)
}

// ── Accounts Receivable (AR) & Invoicing ────────────────────────────────────

func (s *transactionService) GetInvoices(page, perPage int, status, search string) ([]model.Transaction, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.txRepo.FindInvoices(page, perPage, status, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *transactionService) CreateInvoice(dto CreateInvoiceDTO) (*model.Transaction, error) {
	if strings.TrimSpace(dto.ClientName) == "" {
		return nil, errors.New("nama klien wajib diisi untuk penerbitan invoice")
	}

	category := strings.TrimSpace(dto.Category)
	if category == "" {
		category = "Project Revenue"
	}

	refNo := strings.TrimSpace(dto.ReferenceNo)
	if refNo == "" {
		refNo = fmt.Sprintf("INV-%s-%s", time.Now().Format("200601"), uuid.NewString()[:6])
	}

	return s.Create(CreateTransactionDTO{
		ReferenceNo:      refNo,
		Type:             "Debit",
		Category:         category,
		Amount:           dto.Amount,
		TaxAmount:        dto.TaxAmount,
		Currency:         dto.Currency,
		Status:           "Sent", // Invoice yang diterbitkan ke klien status defaultnya Sent
		DueDate:          dto.DueDate,
		TransactionDate:  dto.TransactionDate,
		ClientVendorName: dto.ClientName,
		Description:      dto.Description,
	})
}

func (s *transactionService) PayInvoice(identifier string) (*model.Transaction, error) {
	return s.UpdateStatus(identifier, "Paid")
}

// ── Accounts Payable (AP) & Vendor Bills ────────────────────────────────────

func (s *transactionService) GetBills(page, perPage int, status, search string) ([]model.Transaction, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.txRepo.FindBills(page, perPage, status, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *transactionService) CreateBill(dto CreateBillDTO) (*model.Transaction, error) {
	if strings.TrimSpace(dto.VendorName) == "" {
		return nil, errors.New("nama vendor (misal: AWS, Google Cloud, lisensi software) wajib diisi")
	}

	category := strings.TrimSpace(dto.Category)
	if category == "" {
		category = "Software License"
	}

	refNo := strings.TrimSpace(dto.ReferenceNo)
	if refNo == "" {
		refNo = fmt.Sprintf("BILL-%s-%s", time.Now().Format("200601"), uuid.NewString()[:6])
	}

	return s.Create(CreateTransactionDTO{
		ReferenceNo:      refNo,
		Type:             "Credit",
		Category:         category,
		Amount:           dto.Amount,
		TaxAmount:        dto.TaxAmount,
		Currency:         dto.Currency,
		Status:           "Draft",
		DueDate:          dto.DueDate,
		TransactionDate:  dto.TransactionDate,
		ClientVendorName: dto.VendorName,
		Description:      dto.Description,
	})
}

func (s *transactionService) PayBill(identifier string) (*model.Transaction, error) {
	return s.UpdateStatus(identifier, "Paid")
}

// ── General Ledger & Double-Entry Bookkeeping ───────────────────────────────

func (s *transactionService) createJournalForTransaction(tx *model.Transaction) error {
	total := tx.TotalAmount
	if total <= 0 {
		total = tx.Amount + tx.TaxAmount
	}

	entryNumber := fmt.Sprintf("JRN-%s-%s", time.Now().Format("20060102"), uuid.NewString()[:6])

	var items []model.JournalItem

	if tx.Type == "Debit" {
		// Penerimaan Kas dari Klien:
		// [Debit] Kas & Bank (1001) = TotalAmount
		// [Credit] Pendapatan Proyek (4001) = Amount
		// [Credit] Hutang Pajak (2002) = TaxAmount (jika ada)
		items = append(items, model.JournalItem{
			AccountCode: "1001",
			AccountName: "Kas & Bank",
			Debit:       total,
			Credit:      0,
			Memo:        fmt.Sprintf("Penerimaan dari %s (%s)", tx.ClientVendorName, tx.ReferenceNo),
		})

		items = append(items, model.JournalItem{
			AccountCode: "4001",
			AccountName: "Pendapatan Proyek",
			Debit:       0,
			Credit:      tx.Amount,
			Memo:        tx.Description,
		})

		if tx.TaxAmount > 0 {
			items = append(items, model.JournalItem{
				AccountCode: "2002",
				AccountName: "Hutang Pajak (Tax Payable)",
				Debit:       0,
				Credit:      tx.TaxAmount,
				Memo:        "Pajak PPN/PPh penerimaan",
			})
		}
	} else {
		// Pembayaran ke Vendor (AWS, GCP, Lisensi):
		// [Debit] Beban Terkait (5001/5002) = Amount
		// [Credit] Kas & Bank (1001) = TotalAmount
		accountCode := "5003"
		accountName := "Beban Operasional & Supplies"

		catLower := strings.ToLower(tx.Category)
		if strings.Contains(catLower, "license") || strings.Contains(catLower, "software") {
			accountCode = "5001"
			accountName = "Beban Lisensi Software"
		} else if strings.Contains(catLower, "cloud") || strings.Contains(catLower, "server") || strings.Contains(catLower, "hosting") {
			accountCode = "5002"
			accountName = "Beban Cloud & Server Hosting"
		}

		items = append(items, model.JournalItem{
			AccountCode: accountCode,
			AccountName: accountName,
			Debit:       tx.Amount,
			Credit:      0,
			Memo:        fmt.Sprintf("Pembayaran ke %s (%s)", tx.ClientVendorName, tx.ReferenceNo),
		})

		if tx.TaxAmount > 0 {
			items = append(items, model.JournalItem{
				AccountCode: "5003",
				AccountName: "Beban Pajak Masukan",
				Debit:       tx.TaxAmount,
				Credit:      0,
				Memo:        "Pajak tagihan vendor",
			})
		}

		items = append(items, model.JournalItem{
			AccountCode: "1001",
			AccountName: "Kas & Bank",
			Debit:       0,
			Credit:      total,
			Memo:        fmt.Sprintf("Pengeluaran kas untuk %s", tx.ReferenceNo),
		})
	}

	entry := &model.JournalEntry{
		EntryNumber:   entryNumber,
		Date:          tx.TransactionDate,
		Description:   fmt.Sprintf("Jurnal transaksi %s: %s", tx.ReferenceNo, tx.Description),
		TransactionID: &tx.TransactionID,
		TotalDebit:    total,
		TotalCredit:   total,
		Items:         items,
	}

	return s.ledgerRepo.CreateJournalEntry(entry)
}

func (s *transactionService) GetAccounts() ([]model.Account, error) {
	return s.ledgerRepo.FindAccounts()
}

func (s *transactionService) GetJournalEntries(page, perPage int) ([]model.JournalEntry, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}
	offset := (page - 1) * perPage

	entries, total, err := s.ledgerRepo.FindJournalEntries(perPage, offset)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return entries, meta, nil
}

func (s *transactionService) CreateJournalEntry(dto CreateJournalEntryDTO) (*model.JournalEntry, error) {
	if len(dto.Items) < 2 {
		return nil, errors.New("jurnal double-entry membutuhkan minimal 2 pos baris (debit dan credit)")
	}

	var totalDebit, totalCredit float64
	var items []model.JournalItem

	for _, it := range dto.Items {
		totalDebit += it.Debit
		totalCredit += it.Credit

		acc, err := s.ledgerRepo.FindAccountByCode(it.AccountCode)
		accName := it.AccountCode
		if err == nil && acc != nil {
			accName = acc.Name
		}

		items = append(items, model.JournalItem{
			AccountCode: it.AccountCode,
			AccountName: accName,
			Debit:       it.Debit,
			Credit:      it.Credit,
			Memo:        it.Memo,
		})
	}

	// Validasi keseimbangan debit dan credit
	if totalDebit != totalCredit {
		return nil, fmt.Errorf("jurnal tidak seimbang: Total Debit (%.2f) != Total Credit (%.2f)", totalDebit, totalCredit)
	}

	entryNumber := dto.EntryNumber
	if strings.TrimSpace(entryNumber) == "" {
		entryNumber = fmt.Sprintf("JRN-%s-%s", time.Now().Format("20060102"), uuid.NewString()[:6])
	}

	entryDate := dto.Date
	if entryDate.IsZero() {
		entryDate = time.Now()
	}

	entry := &model.JournalEntry{
		EntryNumber: entryNumber,
		Date:        entryDate,
		Description: dto.Description,
		TotalDebit:  totalDebit,
		TotalCredit: totalCredit,
		Items:       items,
	}

	if err := s.ledgerRepo.CreateJournalEntry(entry); err != nil {
		return nil, err
	}

	return entry, nil
}

// ── Financial Reporting ─────────────────────────────────────────────────────

func (s *transactionService) GetSummary() (*model.FinanceSummary, error) {
	return s.txRepo.GetSummary()
}

func (s *transactionService) GetIncomeStatement(start, end time.Time) (*model.IncomeStatement, error) {
	if start.IsZero() {
		// Default: awal tahun ini
		start = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}
	if end.IsZero() {
		end = time.Now().AddDate(0, 0, 1)
	}
	return s.txRepo.GetIncomeStatement(start, end)
}

func (s *transactionService) GetBalanceSheet(asOf time.Time) (*model.BalanceSheet, error) {
	if asOf.IsZero() {
		asOf = time.Now()
	}
	return s.txRepo.GetBalanceSheet(asOf)
}

func (s *transactionService) GetCashFlow(start, end time.Time) (*model.CashFlowReport, error) {
	if start.IsZero() {
		start = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}
	if end.IsZero() {
		end = time.Now().AddDate(0, 0, 1)
	}
	return s.txRepo.GetCashFlow(start, end)
}

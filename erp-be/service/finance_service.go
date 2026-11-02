package service

import (
	"errors"
	"strings"
	"time"

	"erp-kafeinart/model"
	"erp-kafeinart/repository"
	"erp-kafeinart/utils"

	"github.com/google/uuid"
)

// ── DTO Data Master ─────────────────────────────────────────────────────────

type CategoryDTO struct {
	CategoryName string `json:"category_name"`
	NamaKategori string `json:"namaKategori"`
	CashType     string `json:"cash_type"` // 'INCOME' / 'EXPENSE'
	TipeKas      string `json:"tipeKas"`   // 'Masuk' / 'Keluar'
	Description  string `json:"description"`
	Deskripsi    string `json:"deskripsi"`
}

type AccountDTO struct {
	BankName          string  `json:"bank_name"`
	NamaBank          string  `json:"namaBank"`
	AccountNumber     string  `json:"account_number"`
	NoRekening        string  `json:"noRekening"`
	AccountHolderName string  `json:"account_holder_name"`
	AtasNama          string  `json:"atasNama"`
	CurrentBalance    float64 `json:"current_balance"`
	SaldoAwal         float64 `json:"saldoAwal"`
}

type TaxDTO struct {
	TaxName     string  `json:"tax_name"`
	NamaPajak   string  `json:"namaPajak"`
	Percentage  float64 `json:"percentage"`
	Persentase  float64 `json:"persentase"`
	Description string  `json:"description"`
	Keterangan  string  `json:"keterangan"`
}

func parseFlexibleDate(s string) time.Time {
	s = strings.TrimSpace(s)
	if s == "" {
		return time.Time{}
	}
	formats := []string{
		"2006-01-02",
		"2006-01-02T15:04:05Z07:00",
		time.RFC3339,
		time.RFC3339Nano,
		"2006-01-02 15:04:05",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Time{}
}

// ── DTO Data Transaksi ──────────────────────────────────────────────────────

type CashJournalDTO struct {
	TransactionNumber string  `json:"transaction_number"`
	TransactionDate   string  `json:"transaction_date"`
	Tanggal           string  `json:"tanggal"`
	AccountID         string  `json:"account_id"`
	AkunBank          string  `json:"akunBank"`
	CategoryID        string  `json:"category_id"`
	Kategori          string  `json:"kategori"`
	TransactionType   string  `json:"transaction_type"` // 'INCOME' / 'EXPENSE'
	Tipe              string  `json:"tipe"`              // 'Masuk' / 'Keluar'
	Amount            float64 `json:"amount"`
	Nominal           float64 `json:"nominal"`
	Description       string  `json:"description"`
	Keterangan        string  `json:"keterangan"`
	AttachmentURL     string  `json:"attachment_url"`
	BuktiFile         string  `json:"buktiFile"`
}

type InvoiceItemDTO struct {
	ItemDescription string  `json:"item_description"`
	Deskripsi       string  `json:"deskripsi"`
	Quantity        int     `json:"quantity"`
	Qty             int     `json:"qty"`
	UnitPrice       float64 `json:"unit_price"`
	Harga           float64 `json:"harga"`
	TaxID           *string `json:"tax_id"`
	Pajak           *string `json:"pajak"`
}

type InvoiceDTO struct {
	InvoiceNumber string           `json:"invoice_number"`
	NoInvoice     string           `json:"noInvoice"`
	ClientID      string           `json:"client_id"`
	ClientName    string           `json:"client_name"`
	Klien         string           `json:"klien"`
	ProjectID     string           `json:"project_id"`
	ProjectName   string           `json:"project_name"`
	Proyek        string           `json:"proyek"`
	IssueDate     string           `json:"issue_date"`
	TglTerbit     string           `json:"tglTerbit"`
	DueDate       string           `json:"due_date"`
	TglJatuhTempo string           `json:"tglJatuhTempo"`
	Status        string           `json:"status"` // 'DRAFT', 'SENT', 'PARTIAL', 'PAID', 'CANCELLED'
	Notes         string           `json:"notes"`
	Pajak         string           `json:"pajak"`
	CatatanPajak  string           `json:"catatanPajak"`
	Items         []InvoiceItemDTO `json:"items"`
}

type VendorBillDTO struct {
	BillNumber     string  `json:"bill_number"`
	NoTagihan      string  `json:"noTagihan"`
	VendorID       string  `json:"vendor_id"`
	VendorName     string  `json:"vendor_name"`
	Vendor         string  `json:"vendor"`
	CategoryID     string  `json:"category_id"`
	Kategori       string  `json:"kategori"`
	BillDate       string  `json:"bill_date"`
	TanggalTagihan string  `json:"tanggalTagihan"`
	DueDate        string  `json:"due_date"`
	JatuhTempo     string  `json:"jatuhTempo"`
	Amount         float64 `json:"amount"`
	Nominal        float64 `json:"nominal"`
	Status         string  `json:"status"` // 'PENDING', 'PAID', 'OVERDUE'
	AttachmentURL  string  `json:"attachment_url"`
	BuktiFile      string  `json:"buktiFile"`
}

type UpdateStatusRequest struct {
	Status string `json:"status"`
}

// ── FinanceService Interface ────────────────────────────────────────────────

type FinanceService interface {
	// Master Kategori
	GetCategories(cashType string) ([]model.MstFinanceCategory, error)
	GetCategoriesPaginated(page, perPage int, cashType, search string) ([]model.MstFinanceCategory, utils.PaginationMeta, error)
	GetCategoryByID(id string) (*model.MstFinanceCategory, error)
	CreateCategory(dto CategoryDTO) (*model.MstFinanceCategory, error)
	UpdateCategory(id string, dto CategoryDTO) (*model.MstFinanceCategory, error)
	DeleteCategory(id string) error

	// Master Rekening Kas / Bank
	GetAccounts() ([]model.MstFinanceAccount, error)
	GetAccountsPaginated(page, perPage int, search string) ([]model.MstFinanceAccount, utils.PaginationMeta, error)
	GetAccountByID(id string) (*model.MstFinanceAccount, error)
	CreateAccount(dto AccountDTO) (*model.MstFinanceAccount, error)
	UpdateAccount(id string, dto AccountDTO) (*model.MstFinanceAccount, error)
	DeleteAccount(id string) error

	// Master Pajak
	GetTaxes() ([]model.MstFinanceTax, error)
	GetTaxesPaginated(page, perPage int, search string) ([]model.MstFinanceTax, utils.PaginationMeta, error)
	GetTaxByID(id string) (*model.MstFinanceTax, error)
	CreateTax(dto TaxDTO) (*model.MstFinanceTax, error)
	UpdateTax(id string, dto TaxDTO) (*model.MstFinanceTax, error)
	DeleteTax(id string) error

	// Jurnal Kas (Pencatatan Harian)
	GetCashJournals(page, perPage int, accountID, txType, categoryID, search string) ([]model.TrxCashJournal, utils.PaginationMeta, error)
	GetCashJournalByID(id string) (*model.TrxCashJournal, error)
	CreateCashJournal(dto CashJournalDTO) (*model.TrxCashJournal, error)
	UpdateCashJournal(id string, dto CashJournalDTO) (*model.TrxCashJournal, error)
	DeleteCashJournal(id string) error

	// Manajemen Invoice / Piutang (AR)
	GetInvoices(page, perPage int, status, clientID, search string) ([]model.TrxInvoice, utils.PaginationMeta, error)
	GetInvoiceByID(id string) (*model.TrxInvoice, error)
	CreateInvoice(dto InvoiceDTO) (*model.TrxInvoice, error)
	UpdateInvoice(id string, dto InvoiceDTO) (*model.TrxInvoice, error)
	UpdateInvoiceStatus(id, status string) (*model.TrxInvoice, error)
	DeleteInvoice(id string) error

	// Tagihan Vendor / Hutang (AP)
	GetVendorBills(page, perPage int, status, vendorID, search string) ([]model.TrxVendorBill, utils.PaginationMeta, error)
	GetVendorBillByID(id string) (*model.TrxVendorBill, error)
	CreateVendorBill(dto VendorBillDTO) (*model.TrxVendorBill, error)
	UpdateVendorBill(id string, dto VendorBillDTO) (*model.TrxVendorBill, error)
	UpdateVendorBillStatus(id, status string) (*model.TrxVendorBill, error)
	DeleteVendorBill(id string) error

	// Laporan (Reports)
	GetCashFlowReport(start, end time.Time, accountID string) (*model.CashFlowReport, error)
	GetProfitLossReport(month, year int) (*model.ProfitLossReport, error)
	GetARAgingReport(status, clientName string) (*model.ARAgingReport, error)
}

type financeService struct {
	masterRepo repository.FinanceMasterRepository
	trxRepo    repository.FinanceTransactionRepository
}

func NewFinanceService(masterRepo repository.FinanceMasterRepository, trxRepo repository.FinanceTransactionRepository) FinanceService {
	return &financeService{
		masterRepo: masterRepo,
		trxRepo:    trxRepo,
	}
}

// ── Master Kategori Implementasi ────────────────────────────────────────────

func (s *financeService) GetCategories(cashType string) ([]model.MstFinanceCategory, error) {
	return s.masterRepo.FindCategories(cashType)
}

func (s *financeService) GetCategoriesPaginated(page, perPage int, cashType, search string) ([]model.MstFinanceCategory, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	list, total, err := s.masterRepo.FindCategoriesPaginated(page, perPage, cashType, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}
	return list, utils.BuildPaginationMeta(total, page, perPage), nil
}

func (s *financeService) GetCategoryByID(id string) (*model.MstFinanceCategory, error) {
	cat, err := s.masterRepo.FindCategoryByID(id)
	if err != nil || cat == nil {
		return nil, errors.New("kategori tidak ditemukan")
	}
	return cat, nil
}

func (s *financeService) CreateCategory(dto CategoryDTO) (*model.MstFinanceCategory, error) {
	name := strings.TrimSpace(dto.CategoryName)
	if name == "" {
		name = strings.TrimSpace(dto.NamaKategori)
	}
	if name == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	cashType := strings.ToUpper(strings.TrimSpace(dto.CashType))
	if cashType == "" {
		if strings.EqualFold(dto.TipeKas, "masuk") {
			cashType = "INCOME"
		} else if strings.EqualFold(dto.TipeKas, "keluar") {
			cashType = "EXPENSE"
		} else {
			cashType = strings.ToUpper(strings.TrimSpace(dto.TipeKas))
		}
	}
	if cashType != "INCOME" && cashType != "EXPENSE" {
		cashType = "INCOME"
	}

	desc := strings.TrimSpace(dto.Description)
	if desc == "" {
		desc = strings.TrimSpace(dto.Deskripsi)
	}

	cat := &model.MstFinanceCategory{
		ID:           uuid.NewString(),
		CategoryName: name,
		CashType:     cashType,
		Description:  desc,
		CreatedAt:    time.Now(),
	}

	if err := s.masterRepo.CreateCategory(cat); err != nil {
		return nil, err
	}
	_ = cat.AfterFind(nil)
	return cat, nil
}

func (s *financeService) UpdateCategory(id string, dto CategoryDTO) (*model.MstFinanceCategory, error) {
	cat, err := s.GetCategoryByID(id)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(dto.CategoryName)
	if name == "" {
		name = strings.TrimSpace(dto.NamaKategori)
	}
	if name == "" {
		return nil, errors.New("nama kategori wajib diisi")
	}

	cat.CategoryName = name
	cashType := strings.ToUpper(strings.TrimSpace(dto.CashType))
	if cashType == "" && dto.TipeKas != "" {
		if strings.EqualFold(dto.TipeKas, "masuk") {
			cashType = "INCOME"
		} else if strings.EqualFold(dto.TipeKas, "keluar") {
			cashType = "EXPENSE"
		} else {
			cashType = strings.ToUpper(strings.TrimSpace(dto.TipeKas))
		}
	}
	if cashType != "" {
		cat.CashType = cashType
	}
	if dto.Description != "" {
		cat.Description = dto.Description
	} else if dto.Deskripsi != "" {
		cat.Description = dto.Deskripsi
	}

	if err := s.masterRepo.UpdateCategory(cat); err != nil {
		return nil, err
	}
	_ = cat.AfterFind(nil)
	return cat, nil
}

func (s *financeService) DeleteCategory(id string) error {
	_, err := s.GetCategoryByID(id)
	if err != nil {
		return err
	}
	return s.masterRepo.DeleteCategory(id)
}

// ── Master Rekening Implementasi ────────────────────────────────────────────

func (s *financeService) GetAccounts() ([]model.MstFinanceAccount, error) {
	return s.masterRepo.FindAccounts()
}

func (s *financeService) GetAccountsPaginated(page, perPage int, search string) ([]model.MstFinanceAccount, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	list, total, err := s.masterRepo.FindAccountsPaginated(page, perPage, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}
	return list, utils.BuildPaginationMeta(total, page, perPage), nil
}

func (s *financeService) GetAccountByID(id string) (*model.MstFinanceAccount, error) {
	acc, err := s.masterRepo.FindAccountByID(id)
	if err != nil || acc == nil {
		return nil, errors.New("rekening kas/bank tidak ditemukan")
	}
	return acc, nil
}

func (s *financeService) CreateAccount(dto AccountDTO) (*model.MstFinanceAccount, error) {
	bankName := strings.TrimSpace(dto.BankName)
	if bankName == "" {
		bankName = strings.TrimSpace(dto.NamaBank)
	}
	if bankName == "" {
		return nil, errors.New("nama bank/kas wajib diisi")
	}

	accNum := strings.TrimSpace(dto.AccountNumber)
	if accNum == "" {
		accNum = strings.TrimSpace(dto.NoRekening)
	}
	if accNum == "" {
		return nil, errors.New("nomor rekening wajib diisi")
	}

	accHolder := strings.TrimSpace(dto.AccountHolderName)
	if accHolder == "" {
		accHolder = strings.TrimSpace(dto.AtasNama)
	}
	if accHolder == "" {
		return nil, errors.New("nama pemilik rekening (atas nama) wajib diisi")
	}

	balance := dto.CurrentBalance
	if balance == 0 && dto.SaldoAwal != 0 {
		balance = dto.SaldoAwal
	}

	acc := &model.MstFinanceAccount{
		ID:                uuid.NewString(),
		BankName:          bankName,
		AccountNumber:     accNum,
		AccountHolderName: accHolder,
		CurrentBalance:    balance,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.masterRepo.CreateAccount(acc); err != nil {
		return nil, err
	}
	_ = acc.AfterFind(nil)
	return acc, nil
}

func (s *financeService) UpdateAccount(id string, dto AccountDTO) (*model.MstFinanceAccount, error) {
	acc, err := s.GetAccountByID(id)
	if err != nil {
		return nil, err
	}

	bankName := strings.TrimSpace(dto.BankName)
	if bankName == "" {
		bankName = strings.TrimSpace(dto.NamaBank)
	}
	if bankName != "" {
		acc.BankName = bankName
	}

	accNum := strings.TrimSpace(dto.AccountNumber)
	if accNum == "" {
		accNum = strings.TrimSpace(dto.NoRekening)
	}
	if accNum != "" {
		acc.AccountNumber = accNum
	}

	accHolder := strings.TrimSpace(dto.AccountHolderName)
	if accHolder == "" {
		accHolder = strings.TrimSpace(dto.AtasNama)
	}
	if accHolder != "" {
		acc.AccountHolderName = accHolder
	}

	if dto.CurrentBalance != 0 {
		acc.CurrentBalance = dto.CurrentBalance
	} else if dto.SaldoAwal != 0 {
		acc.CurrentBalance = dto.SaldoAwal
	}

	acc.UpdatedAt = time.Now()
	if err := s.masterRepo.UpdateAccount(acc); err != nil {
		return nil, err
	}
	_ = acc.AfterFind(nil)
	return acc, nil
}

func (s *financeService) DeleteAccount(id string) error {
	_, err := s.GetAccountByID(id)
	if err != nil {
		return err
	}
	return s.masterRepo.DeleteAccount(id)
}

// ── Master Pajak Implementasi ───────────────────────────────────────────────

func (s *financeService) GetTaxes() ([]model.MstFinanceTax, error) {
	return s.masterRepo.FindTaxes()
}

func (s *financeService) GetTaxesPaginated(page, perPage int, search string) ([]model.MstFinanceTax, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 10
	}
	list, total, err := s.masterRepo.FindTaxesPaginated(page, perPage, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}
	return list, utils.BuildPaginationMeta(total, page, perPage), nil
}

func (s *financeService) GetTaxByID(id string) (*model.MstFinanceTax, error) {
	tax, err := s.masterRepo.FindTaxByID(id)
	if err != nil || tax == nil {
		return nil, errors.New("master pajak tidak ditemukan")
	}
	return tax, nil
}

func (s *financeService) CreateTax(dto TaxDTO) (*model.MstFinanceTax, error) {
	name := strings.TrimSpace(dto.TaxName)
	if name == "" {
		name = strings.TrimSpace(dto.NamaPajak)
	}
	if name == "" {
		return nil, errors.New("nama pajak wajib diisi (cth: PPN 11%)")
	}

	pct := dto.Percentage
	if pct == 0 && dto.Persentase != 0 {
		pct = dto.Persentase
	}

	desc := strings.TrimSpace(dto.Description)
	if desc == "" {
		desc = strings.TrimSpace(dto.Keterangan)
	}

	tax := &model.MstFinanceTax{
		ID:          uuid.NewString(),
		TaxName:     name,
		Percentage:  pct,
		Description: desc,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	if err := s.masterRepo.CreateTax(tax); err != nil {
		return nil, err
	}
	_ = tax.AfterFind(nil)
	return tax, nil
}

func (s *financeService) UpdateTax(id string, dto TaxDTO) (*model.MstFinanceTax, error) {
	tax, err := s.GetTaxByID(id)
	if err != nil {
		return nil, err
	}

	name := strings.TrimSpace(dto.TaxName)
	if name == "" {
		name = strings.TrimSpace(dto.NamaPajak)
	}
	if name != "" {
		tax.TaxName = name
	}
	if dto.Percentage != 0 {
		tax.Percentage = dto.Percentage
	} else if dto.Persentase != 0 {
		tax.Percentage = dto.Persentase
	}
	if dto.Description != "" {
		tax.Description = dto.Description
	} else if dto.Keterangan != "" {
		tax.Description = dto.Keterangan
	}
	tax.UpdatedAt = time.Now()

	if err := s.masterRepo.UpdateTax(tax); err != nil {
		return nil, err
	}
	_ = tax.AfterFind(nil)
	return tax, nil
}

func (s *financeService) DeleteTax(id string) error {
	_, err := s.GetTaxByID(id)
	if err != nil {
		return err
	}
	return s.masterRepo.DeleteTax(id)
}

package handler

import (
	"strconv"
	"time"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
)

type TransactionHandler struct {
	txService service.TransactionService
}

func NewTransactionHandler(txService service.TransactionService) *TransactionHandler {
	return &TransactionHandler{txService: txService}
}

// GetAll menangani GET /api/v1/finance/transactions
func (h *TransactionHandler) GetAll(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	if page < 1 {
		page = 1
	}

	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("rowsPerPage", "15")
	}
	perPage, _ := strconv.Atoi(perPageStr)
	if perPage < 1 {
		perPage = 15
	}

	txType := c.Query("type")
	category := c.Query("category")
	status := c.Query("status")
	currency := c.Query("currency")
	search := c.Query("search")

	transactions, meta, err := h.txService.GetAll(page, perPage, txType, category, status, currency, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal mengambil daftar transaksi", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, transactions, meta, "Daftar transaksi berhasil dimuat")
}

// GetByID menangani GET /api/v1/finance/transactions/:id (bisa numeric ID atau UUID)
func (h *TransactionHandler) GetByID(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Identifier transaksi wajib diisi")
	}

	tx, err := h.txService.GetByIDOrUUID(idParam)
	if err != nil || tx == nil {
		return utils.SendError(c, fiber.StatusNotFound, "Transaksi tidak ditemukan")
	}

	return utils.SendSuccess(c, fiber.StatusOK, tx, "Detail transaksi berhasil dimuat")
}

// Create menangani POST /api/v1/finance/transactions
func (h *TransactionHandler) Create(c fiber.Ctx) error {
	var req service.CreateTransactionDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input transaksi tidak valid", err.Error())
	}

	tx, err := h.txService.Create(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, tx, "Transaksi berhasil dicatat")
}

// Update menangani PUT /api/v1/finance/transactions/:id
func (h *TransactionHandler) Update(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID transaksi tidak valid")
	}

	var req service.CreateTransactionDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input transaksi tidak valid", err.Error())
	}

	tx, err := h.txService.Update(uint(id), req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, tx, "Transaksi berhasil diperbarui")
}

// UpdateStatus menangani PATCH /api/v1/finance/transactions/:id/status
func (h *TransactionHandler) UpdateStatus(c fiber.Ctx) error {
	idParam := c.Params("id")
	if idParam == "" {
		return utils.SendError(c, fiber.StatusBadRequest, "Identifier transaksi wajib diisi")
	}

	var req service.UpdateStatusDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}

	tx, err := h.txService.UpdateStatus(idParam, req.Status)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, tx, "Status transaksi berhasil diperbarui")
}

// Delete menangani DELETE /api/v1/finance/transactions/:id
func (h *TransactionHandler) Delete(c fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "ID transaksi tidak valid")
	}

	if err := h.txService.Delete(uint(id)); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, nil, "Transaksi berhasil dihapus")
}

// GetSummary menangani GET /api/v1/finance/summary
func (h *TransactionHandler) GetSummary(c fiber.Ctx) error {
	summary, err := h.txService.GetSummary()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat ringkasan keuangan", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, summary, "Ringkasan keuangan berhasil dimuat")
}

// ── Accounts Receivable (AR) & Invoicing Handlers ───────────────────────────

// GetInvoices menangani GET /api/v1/finance/ar/invoices
func (h *TransactionHandler) GetInvoices(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))
	status := c.Query("status")
	search := c.Query("search")

	invoices, meta, err := h.txService.GetInvoices(page, perPage, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat daftar invoice klien", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, invoices, meta, "Daftar invoice klien berhasil dimuat")
}

// CreateInvoice menangani POST /api/v1/finance/ar/invoices
func (h *TransactionHandler) CreateInvoice(c fiber.Ctx) error {
	var req service.CreateInvoiceDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input invoice tidak valid", err.Error())
	}

	invoice, err := h.txService.CreateInvoice(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, invoice, "Invoice klien berhasil diterbitkan")
}

// PayInvoice menangani POST /api/v1/finance/ar/invoices/:id/pay
func (h *TransactionHandler) PayInvoice(c fiber.Ctx) error {
	idParam := c.Params("id")
	invoice, err := h.txService.PayInvoice(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, invoice, "Pembayaran invoice berhasil dikonfirmasi dan dicatat ke Buku Besar")
}

// ── Accounts Payable (AP) & Vendor Bills Handlers ───────────────────────────

// GetBills menangani GET /api/v1/finance/ap/bills
func (h *TransactionHandler) GetBills(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))
	status := c.Query("status")
	search := c.Query("search")

	bills, meta, err := h.txService.GetBills(page, perPage, status, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat tagihan vendor", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, bills, meta, "Daftar tagihan vendor berhasil dimuat")
}

// CreateBill menangani POST /api/v1/finance/ap/bills
func (h *TransactionHandler) CreateBill(c fiber.Ctx) error {
	var req service.CreateBillDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format input tagihan vendor tidak valid", err.Error())
	}

	bill, err := h.txService.CreateBill(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, bill, "Tagihan vendor berhasil dicatat")
}

// PayBill menangani POST /api/v1/finance/ap/bills/:id/pay
func (h *TransactionHandler) PayBill(c fiber.Ctx) error {
	idParam := c.Params("id")
	bill, err := h.txService.PayBill(idParam)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, bill, "Pembayaran tagihan vendor berhasil diproses dan dicatat ke Buku Besar")
}

// ── General Ledger & Bookkeeping Handlers ───────────────────────────────────

// GetAccounts menangani GET /api/v1/finance/ledger/accounts
func (h *TransactionHandler) GetAccounts(c fiber.Ctx) error {
	accounts, err := h.txService.GetAccounts()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat Chart of Accounts", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, accounts, "Chart of Accounts berhasil dimuat")
}

// GetJournalEntries menangani GET /api/v1/finance/ledger/entries
func (h *TransactionHandler) GetJournalEntries(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPage, _ := strconv.Atoi(c.Query("per_page", "15"))

	entries, meta, err := h.txService.GetJournalEntries(page, perPage)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat entri jurnal umum", err.Error())
	}

	return utils.SendPaginated(c, fiber.StatusOK, entries, meta, "Entri jurnal Buku Besar berhasil dimuat")
}

// CreateJournalEntry menangani POST /api/v1/finance/ledger/entries
func (h *TransactionHandler) CreateJournalEntry(c fiber.Ctx) error {
	var req service.CreateJournalEntryDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format jurnal tidak valid", err.Error())
	}

	entry, err := h.txService.CreateJournalEntry(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusCreated, entry, "Entri jurnal double-entry berhasil dicatat")
}

// ── Financial Reporting Handlers ────────────────────────────────────────────

// GetIncomeStatement menangani GET /api/v1/finance/reports/income-statement
func (h *TransactionHandler) GetIncomeStatement(c fiber.Ctx) error {
	var start, end time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr := c.Query("end_date"); endStr != "" {
		end, _ = time.Parse("2006-01-02", endStr)
	}

	report, err := h.txService.GetIncomeStatement(start, end)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menghasilkan Laporan Laba/Rugi", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan Laba/Rugi berhasil disusun")
}

// GetBalanceSheet menangani GET /api/v1/finance/reports/balance-sheet
func (h *TransactionHandler) GetBalanceSheet(c fiber.Ctx) error {
	var asOf time.Time
	if asOfStr := c.Query("as_of_date"); asOfStr != "" {
		asOf, _ = time.Parse("2006-01-02", asOfStr)
	}

	report, err := h.txService.GetBalanceSheet(asOf)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menghasilkan Laporan Neraca", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan Neraca berhasil disusun")
}

// GetCashFlow menangani GET /api/v1/finance/reports/cash-flow
func (h *TransactionHandler) GetCashFlow(c fiber.Ctx) error {
	var start, end time.Time
	if startStr := c.Query("start_date"); startStr != "" {
		start, _ = time.Parse("2006-01-02", startStr)
	}
	if endStr := c.Query("end_date"); endStr != "" {
		end, _ = time.Parse("2006-01-02", endStr)
	}

	report, err := h.txService.GetCashFlow(start, end)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menghasilkan Laporan Arus Kas", err.Error())
	}

	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan Arus Kas berhasil disusun")
}

package handler

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"erp-kafeinart/service"
	"erp-kafeinart/utils"

	"github.com/gofiber/fiber/v3"
	"github.com/google/uuid"
)

type FinanceHandler struct {
	financeService service.FinanceService
}

func NewFinanceHandler(financeService service.FinanceService) *FinanceHandler {
	return &FinanceHandler{financeService: financeService}
}

// ── Master Kategori Handlers ────────────────────────────────────────────────

func (h *FinanceHandler) GetCategories(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")
	cashType := c.Query("cash_type")
	if cashType == "" {
		cashType = c.Query("tipeKas")
	}

	if pageStr != "" || perPageStr != "" || search != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		cats, meta, err := h.financeService.GetCategoriesPaginated(page, perPage, cashType, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat kategori transaksi", err.Error())
		}
		return utils.SendPaginated(c, fiber.StatusOK, cats, meta, "Daftar kategori transaksi berhasil dimuat")
	}

	cats, err := h.financeService.GetCategories(cashType)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat kategori transaksi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, cats, "Daftar kategori transaksi berhasil dimuat")
}

func (h *FinanceHandler) GetCategoryByID(c fiber.Ctx) error {
	id := c.Params("id")
	cat, err := h.financeService.GetCategoryByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, cat, "Detail kategori transaksi berhasil dimuat")
}

func (h *FinanceHandler) CreateCategory(c fiber.Ctx) error {
	var req service.CategoryDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data kategori tidak valid", err.Error())
	}

	cat, err := h.financeService.CreateCategory(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, cat, "Kategori transaksi berhasil ditambahkan")
}

func (h *FinanceHandler) UpdateCategory(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.CategoryDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data kategori tidak valid", err.Error())
	}

	cat, err := h.financeService.UpdateCategory(id, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, cat, "Kategori transaksi berhasil diperbarui")
}

func (h *FinanceHandler) DeleteCategory(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.financeService.DeleteCategory(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Kategori transaksi berhasil dihapus")
}

// ── Master Rekening Kas / Bank Handlers ─────────────────────────────────────

func (h *FinanceHandler) GetAccounts(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")

	if pageStr != "" || perPageStr != "" || search != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		accounts, meta, err := h.financeService.GetAccountsPaginated(page, perPage, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat rekening kas/bank", err.Error())
		}
		return utils.SendPaginated(c, fiber.StatusOK, accounts, meta, "Daftar rekening kas/bank berhasil dimuat")
	}

	accounts, err := h.financeService.GetAccounts()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat rekening kas/bank", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, accounts, "Daftar rekening kas/bank berhasil dimuat")
}

func (h *FinanceHandler) GetAccountByID(c fiber.Ctx) error {
	id := c.Params("id")
	acc, err := h.financeService.GetAccountByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, acc, "Detail rekening kas/bank berhasil dimuat")
}

func (h *FinanceHandler) CreateAccount(c fiber.Ctx) error {
	var req service.AccountDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data rekening tidak valid", err.Error())
	}

	acc, err := h.financeService.CreateAccount(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, acc, "Rekening kas/bank berhasil ditambahkan")
}

func (h *FinanceHandler) UpdateAccount(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.AccountDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data rekening tidak valid", err.Error())
	}

	acc, err := h.financeService.UpdateAccount(id, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, acc, "Rekening kas/bank berhasil diperbarui")
}

func (h *FinanceHandler) DeleteAccount(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.financeService.DeleteAccount(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Rekening kas/bank berhasil dihapus")
}

// ── Master Pajak Handlers ───────────────────────────────────────────────────

func (h *FinanceHandler) GetTaxes(c fiber.Ctx) error {
	pageStr := c.Query("page")
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage")
	}
	search := c.Query("search")

	if pageStr != "" || perPageStr != "" || search != "" {
		page, _ := strconv.Atoi(pageStr)
		perPage, _ := strconv.Atoi(perPageStr)
		taxes, meta, err := h.financeService.GetTaxesPaginated(page, perPage, search)
		if err != nil {
			return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat master pajak", err.Error())
		}
		return utils.SendPaginated(c, fiber.StatusOK, taxes, meta, "Daftar tarif pajak berhasil dimuat")
	}

	taxes, err := h.financeService.GetTaxes()
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat master pajak", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, taxes, "Daftar tarif pajak berhasil dimuat")
}

func (h *FinanceHandler) GetTaxByID(c fiber.Ctx) error {
	id := c.Params("id")
	tax, err := h.financeService.GetTaxByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, tax, "Detail master pajak berhasil dimuat")
}

func (h *FinanceHandler) CreateTax(c fiber.Ctx) error {
	var req service.TaxDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data pajak tidak valid", err.Error())
	}

	tax, err := h.financeService.CreateTax(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, tax, "Master pajak berhasil ditambahkan")
}

func (h *FinanceHandler) UpdateTax(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.TaxDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data pajak tidak valid", err.Error())
	}

	tax, err := h.financeService.UpdateTax(id, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, tax, "Master pajak berhasil diperbarui")
}

func (h *FinanceHandler) DeleteTax(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.financeService.DeleteTax(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Master pajak berhasil dihapus")
}

// ── Jurnal Kas (Pencatatan Harian) Handlers ─────────────────────────────────

func (h *FinanceHandler) GetCashJournals(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "15")
	}
	perPage, _ := strconv.Atoi(perPageStr)

	accountID := c.Query("account_id")
	if accountID == "" {
		accountID = c.Query("akunBank")
	}
	txType := c.Query("transaction_type")
	if txType == "" {
		txType = c.Query("tipe")
	}
	categoryID := c.Query("category_id")
	if categoryID == "" {
		categoryID = c.Query("kategori")
	}
	search := c.Query("search")

	journals, meta, err := h.financeService.GetCashJournals(page, perPage, accountID, txType, categoryID, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat jurnal kas", err.Error())
	}
	return utils.SendPaginated(c, fiber.StatusOK, journals, meta, "Jurnal kas harian berhasil dimuat")
}

func (h *FinanceHandler) GetCashJournalByID(c fiber.Ctx) error {
	id := c.Params("id")
	j, err := h.financeService.GetCashJournalByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, j, "Detail jurnal kas berhasil dimuat")
}

func (h *FinanceHandler) CreateCashJournal(c fiber.Ctx) error {
	var req service.CashJournalDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data jurnal kas tidak valid", err.Error())
	}

	j, err := h.financeService.CreateCashJournal(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, j, "Transaksi kas harian berhasil dicatat")
}

func (h *FinanceHandler) UpdateCashJournal(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.CashJournalDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data jurnal kas tidak valid", err.Error())
	}

	j, err := h.financeService.UpdateCashJournal(id, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, j, "Transaksi jurnal kas berhasil diperbarui")
}

func (h *FinanceHandler) DeleteCashJournal(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.financeService.DeleteCashJournal(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Transaksi jurnal kas berhasil dihapus dan saldo disesuaikan")
}

// ── Manajemen Invoice / Piutang (AR) Handlers ───────────────────────────────

func (h *FinanceHandler) GetInvoices(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "15")
	}
	perPage, _ := strconv.Atoi(perPageStr)

	status := c.Query("status")
	clientID := c.Query("client_id")
	if clientID == "" {
		clientID = c.Query("klien")
	}
	search := c.Query("search")

	invoices, meta, err := h.financeService.GetInvoices(page, perPage, status, clientID, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat daftar invoice", err.Error())
	}
	return utils.SendPaginated(c, fiber.StatusOK, invoices, meta, "Daftar invoice klien berhasil dimuat")
}

func (h *FinanceHandler) GetInvoiceByID(c fiber.Ctx) error {
	id := c.Params("id")
	inv, err := h.financeService.GetInvoiceByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, inv, "Detail invoice berhasil dimuat")
}

func (h *FinanceHandler) CreateInvoice(c fiber.Ctx) error {
	var req service.InvoiceDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data invoice tidak valid", err.Error())
	}

	inv, err := h.financeService.CreateInvoice(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, inv, "Invoice berhasil diterbitkan")
}

func (h *FinanceHandler) UpdateInvoice(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.InvoiceDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data invoice tidak valid", err.Error())
	}

	inv, err := h.financeService.UpdateInvoice(id, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, inv, "Invoice berhasil diperbarui")
}

func (h *FinanceHandler) UpdateInvoiceStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.UpdateStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}

	inv, err := h.financeService.UpdateInvoiceStatus(id, req.Status)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, inv, "Status invoice berhasil diperbarui")
}

func (h *FinanceHandler) DeleteInvoice(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.financeService.DeleteInvoice(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Invoice berhasil dihapus")
}

// ── Tagihan Vendor / Hutang (AP) Handlers ───────────────────────────────────

func (h *FinanceHandler) GetVendorBills(c fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("page", "1"))
	perPageStr := c.Query("per_page")
	if perPageStr == "" {
		perPageStr = c.Query("perPage", "15")
	}
	perPage, _ := strconv.Atoi(perPageStr)

	status := c.Query("status")
	vendorID := c.Query("vendor_id")
	if vendorID == "" {
		vendorID = c.Query("vendor")
	}
	search := c.Query("search")

	bills, meta, err := h.financeService.GetVendorBills(page, perPage, status, vendorID, search)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal memuat tagihan vendor", err.Error())
	}
	return utils.SendPaginated(c, fiber.StatusOK, bills, meta, "Daftar tagihan vendor berhasil dimuat")
}

func (h *FinanceHandler) GetVendorBillByID(c fiber.Ctx) error {
	id := c.Params("id")
	bill, err := h.financeService.GetVendorBillByID(id)
	if err != nil {
		return utils.SendError(c, fiber.StatusNotFound, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, bill, "Detail tagihan vendor berhasil dimuat")
}

func (h *FinanceHandler) CreateVendorBill(c fiber.Ctx) error {
	var req service.VendorBillDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data tagihan tidak valid", err.Error())
	}

	bill, err := h.financeService.CreateVendorBill(req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusCreated, bill, "Tagihan vendor berhasil dicatat")
}

func (h *FinanceHandler) UpdateVendorBill(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.VendorBillDTO
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format data tagihan tidak valid", err.Error())
	}

	bill, err := h.financeService.UpdateVendorBill(id, req)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, bill, "Tagihan vendor berhasil diperbarui")
}

func (h *FinanceHandler) UpdateVendorBillStatus(c fiber.Ctx) error {
	id := c.Params("id")
	var req service.UpdateStatusRequest
	if err := c.Bind().Body(&req); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "Format status tidak valid", err.Error())
	}

	bill, err := h.financeService.UpdateVendorBillStatus(id, req.Status)
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, bill, "Status tagihan vendor berhasil diperbarui")
}

func (h *FinanceHandler) DeleteVendorBill(c fiber.Ctx) error {
	id := c.Params("id")
	if err := h.financeService.DeleteVendorBill(id); err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, nil, "Tagihan vendor berhasil dihapus")
}

// ── File Upload Handler ─────────────────────────────────────────────────────

func (h *FinanceHandler) UploadAttachment(c fiber.Ctx) error {
	file, err := c.FormFile("file")
	if err != nil {
		return utils.SendError(c, fiber.StatusBadRequest, "File bukti/lampiran tidak ditemukan", err.Error())
	}

	uploadDir := "./storage/uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal membuat folder upload", err.Error())
	}

	ext := filepath.Ext(file.Filename)
	newFileName := fmt.Sprintf("%s_%s%s", time.Now().Format("20060102150405"), uuid.NewString()[:8], ext)
	dst := filepath.Join(uploadDir, newFileName)

	if err := c.SaveFile(file, dst); err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyimpan file", err.Error())
	}

	fileURL := fmt.Sprintf("/uploads/%s", newFileName)
	return utils.SendSuccess(c, fiber.StatusOK, fiber.Map{
		"file_name": file.Filename,
		"file_url":  fileURL,
		"size":      file.Size,
	}, "File lampiran berhasil diunggah")
}

// ── Laporan (Reports) Handlers ──────────────────────────────────────────────

// 1. Laporan Arus Kas (Cash Flow)
// Filter: Periode Awal DatePicker, Periode Akhir DatePicker, Akun Bank ComboBox
func (h *FinanceHandler) GetCashFlowReport(c fiber.Ctx) error {
	var start, end time.Time
	s := c.Query("start_date")
	if s == "" {
		s = c.Query("periodeStart")
	}
	if s != "" {
		start, _ = time.Parse("2006-01-02", s)
	}

	e := c.Query("end_date")
	if e == "" {
		e = c.Query("periodeEnd")
	}
	if e != "" {
		end, _ = time.Parse("2006-01-02", e)
	}

	accountID := c.Query("account_id")
	if accountID == "" {
		accountID = c.Query("akunBank")
	}

	report, err := h.financeService.GetCashFlowReport(start, end, accountID)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyusun Laporan Arus Kas", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan Arus Kas berhasil disusun")
}

// 2. Laporan Laba/Rugi (Profit & Loss)
// Filter: Bulan ComboBox, Tahun ComboBox
func (h *FinanceHandler) GetProfitLossReport(c fiber.Ctx) error {
	mStr := c.Query("month")
	if mStr == "" {
		mStr = c.Query("bulan")
	}
	yStr := c.Query("year")
	if yStr == "" {
		yStr = c.Query("tahun")
	}

	month, _ := strconv.Atoi(mStr)
	year, _ := strconv.Atoi(yStr)

	report, err := h.financeService.GetProfitLossReport(month, year)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyusun Laporan Laba/Rugi", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan Laba/Rugi berhasil disusun")
}

// 3. Laporan Piutang (AR Aging)
// Filter: Status Pembayaran ComboBox, Klien ComboBox
func (h *FinanceHandler) GetARAgingReport(c fiber.Ctx) error {
	status := c.Query("status")
	clientName := c.Query("client")
	if clientName == "" {
		clientName = c.Query("klien")
	}

	report, err := h.financeService.GetARAgingReport(status, clientName)
	if err != nil {
		return utils.SendError(c, fiber.StatusInternalServerError, "Gagal menyusun Laporan Piutang (AR Aging)", err.Error())
	}
	return utils.SendSuccess(c, fiber.StatusOK, report, "Laporan Piutang (AR Aging) berhasil disusun")
}

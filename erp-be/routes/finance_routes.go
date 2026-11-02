package routes

import (
	"erp-kafeinart/handler"

	"github.com/gofiber/fiber/v3"
)

// RegisterFinanceRoutes mendaftarkan seluruh endpoint modul Keuangan (Finance).
func RegisterFinanceRoutes(api fiber.Router, financeHandler *handler.FinanceHandler, txHandler *handler.TransactionHandler, jwtAuth fiber.Handler) {
	finance := api.Group("/finance", jwtAuth)

	// ── A. Data Master ──────────────────────────────────────────────────────────
	// 1. Master Kategori Transaksi
	finance.Get("/categories", financeHandler.GetCategories)
	finance.Get("/categories/:id", financeHandler.GetCategoryByID)
	finance.Post("/categories", financeHandler.CreateCategory)
	finance.Put("/categories/:id", financeHandler.UpdateCategory)
	finance.Delete("/categories/:id", financeHandler.DeleteCategory)
	// Alias bahasa Indonesia
	finance.Get("/kategori", financeHandler.GetCategories)
	finance.Get("/kategori/:id", financeHandler.GetCategoryByID)
	finance.Post("/kategori", financeHandler.CreateCategory)
	finance.Put("/kategori/:id", financeHandler.UpdateCategory)
	finance.Delete("/kategori/:id", financeHandler.DeleteCategory)

	// 2. Master Rekening / Kas & Bank
	finance.Get("/accounts", financeHandler.GetAccounts)
	finance.Get("/accounts/:id", financeHandler.GetAccountByID)
	finance.Post("/accounts", financeHandler.CreateAccount)
	finance.Put("/accounts/:id", financeHandler.UpdateAccount)
	finance.Delete("/accounts/:id", financeHandler.DeleteAccount)
	// Alias bahasa Indonesia
	finance.Get("/rekening", financeHandler.GetAccounts)
	finance.Get("/rekening/:id", financeHandler.GetAccountByID)
	finance.Post("/rekening", financeHandler.CreateAccount)
	finance.Put("/rekening/:id", financeHandler.UpdateAccount)
	finance.Delete("/rekening/:id", financeHandler.DeleteAccount)

	// 3. Master Pajak
	finance.Get("/taxes", financeHandler.GetTaxes)
	finance.Get("/taxes/:id", financeHandler.GetTaxByID)
	finance.Post("/taxes", financeHandler.CreateTax)
	finance.Put("/taxes/:id", financeHandler.UpdateTax)
	finance.Delete("/taxes/:id", financeHandler.DeleteTax)
	// Alias bahasa Indonesia
	finance.Get("/pajak", financeHandler.GetTaxes)
	finance.Get("/pajak/:id", financeHandler.GetTaxByID)
	finance.Post("/pajak", financeHandler.CreateTax)
	finance.Put("/pajak/:id", financeHandler.UpdateTax)
	finance.Delete("/pajak/:id", financeHandler.DeleteTax)

	// ── B. Data Transaksi ───────────────────────────────────────────────────────
	// 1. Jurnal Kas Harian
	finance.Get("/cash-journals", financeHandler.GetCashJournals)
	finance.Get("/cash-journals/:id", financeHandler.GetCashJournalByID)
	finance.Post("/cash-journals", financeHandler.CreateCashJournal)
	finance.Put("/cash-journals/:id", financeHandler.UpdateCashJournal)
	finance.Delete("/cash-journals/:id", financeHandler.DeleteCashJournal)
	// Alias bahasa Indonesia
	finance.Get("/jurnal-kas", financeHandler.GetCashJournals)
	finance.Get("/jurnal-kas/:id", financeHandler.GetCashJournalByID)
	finance.Post("/jurnal-kas", financeHandler.CreateCashJournal)
	finance.Put("/jurnal-kas/:id", financeHandler.UpdateCashJournal)
	finance.Delete("/jurnal-kas/:id", financeHandler.DeleteCashJournal)
	finance.Get("/jurnal", financeHandler.GetCashJournals)
	finance.Get("/jurnal/:id", financeHandler.GetCashJournalByID)
	finance.Post("/jurnal", financeHandler.CreateCashJournal)
	finance.Put("/jurnal/:id", financeHandler.UpdateCashJournal)
	finance.Delete("/jurnal/:id", financeHandler.DeleteCashJournal)

	// 2. Invoice / Piutang Usaha (AR)
	finance.Get("/invoices", financeHandler.GetInvoices)
	finance.Get("/invoices/:id", financeHandler.GetInvoiceByID)
	finance.Post("/invoices", financeHandler.CreateInvoice)
	finance.Put("/invoices/:id", financeHandler.UpdateInvoice)
	finance.Patch("/invoices/:id/status", financeHandler.UpdateInvoiceStatus)
	finance.Delete("/invoices/:id", financeHandler.DeleteInvoice)
	// Alias bahasa Indonesia
	finance.Get("/invoice", financeHandler.GetInvoices)
	finance.Get("/invoice/:id", financeHandler.GetInvoiceByID)
	finance.Post("/invoice", financeHandler.CreateInvoice)
	finance.Put("/invoice/:id", financeHandler.UpdateInvoice)
	finance.Patch("/invoice/:id/status", financeHandler.UpdateInvoiceStatus)
	finance.Delete("/invoice/:id", financeHandler.DeleteInvoice)

	// 3. Tagihan Vendor / Hutang Usaha (AP)
	finance.Get("/vendor-bills", financeHandler.GetVendorBills)
	finance.Get("/vendor-bills/:id", financeHandler.GetVendorBillByID)
	finance.Post("/vendor-bills", financeHandler.CreateVendorBill)
	finance.Put("/vendor-bills/:id", financeHandler.UpdateVendorBill)
	finance.Patch("/vendor-bills/:id/status", financeHandler.UpdateVendorBillStatus)
	finance.Delete("/vendor-bills/:id", financeHandler.DeleteVendorBill)
	// Alias bahasa Indonesia
	finance.Get("/tagihan-vendor", financeHandler.GetVendorBills)
	finance.Get("/tagihan-vendor/:id", financeHandler.GetVendorBillByID)
	finance.Post("/tagihan-vendor", financeHandler.CreateVendorBill)
	finance.Put("/tagihan-vendor/:id", financeHandler.UpdateVendorBill)
	finance.Patch("/tagihan-vendor/:id/status", financeHandler.UpdateVendorBillStatus)
	finance.Delete("/tagihan-vendor/:id", financeHandler.DeleteVendorBill)
	finance.Get("/tagihan", financeHandler.GetVendorBills)
	finance.Get("/tagihan/:id", financeHandler.GetVendorBillByID)
	finance.Post("/tagihan", financeHandler.CreateVendorBill)
	finance.Put("/tagihan/:id", financeHandler.UpdateVendorBill)
	finance.Patch("/tagihan/:id/status", financeHandler.UpdateVendorBillStatus)
	finance.Delete("/tagihan/:id", financeHandler.DeleteVendorBill)

	// 4. File Upload Bukti Transaksi
	finance.Post("/upload", financeHandler.UploadAttachment)

	// ── C. Laporan Keuangan ─────────────────────────────────────────────────────
	finance.Get("/reports/cash-flow", financeHandler.GetCashFlowReport)
	finance.Get("/reports/profit-loss", financeHandler.GetProfitLossReport)
	finance.Get("/reports/ar-aging", financeHandler.GetARAgingReport)
	// Alias bahasa Indonesia
	finance.Get("/laporan/cash-flow", financeHandler.GetCashFlowReport)
	finance.Get("/laporan/laba-rugi", financeHandler.GetProfitLossReport)
	finance.Get("/laporan/ar-aging", financeHandler.GetARAgingReport)

	// ── D. Rute Kompatibilitas Legacy Starter ──────────────────────────────────
	finance.Get("/transactions", txHandler.GetAll)
	finance.Get("/summary", txHandler.GetSummary)
	finance.Get("/transactions/:id", txHandler.GetByID)
	finance.Post("/transactions", txHandler.Create)
	finance.Put("/transactions/:id", txHandler.Update)
	finance.Patch("/transactions/:id/status", txHandler.UpdateStatus)
	finance.Delete("/transactions/:id", txHandler.Delete)

	finance.Get("/ar/invoices", txHandler.GetInvoices)
	finance.Post("/ar/invoices", txHandler.CreateInvoice)
	finance.Post("/ar/invoices/:id/pay", txHandler.PayInvoice)

	finance.Get("/ap/bills", txHandler.GetBills)
	finance.Post("/ap/bills", txHandler.CreateBill)
	finance.Post("/ap/bills/:id/pay", txHandler.PayBill)

	finance.Get("/ledger/accounts", txHandler.GetAccounts)
	finance.Get("/ledger/entries", txHandler.GetJournalEntries)
	finance.Post("/ledger/entries", txHandler.CreateJournalEntry)
	finance.Get("/reports/income-statement", txHandler.GetIncomeStatement)
	finance.Get("/reports/balance-sheet", txHandler.GetBalanceSheet)
}

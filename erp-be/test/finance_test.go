package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp-kafeinart/model"
)

func TestFinanceModuleComplete(t *testing.T) {
	app, db, _ := setupTestApp(t)

	// 1. Login to obtain JWT Token
	loginPayload := map[string]string{
		"email":    "admin@kafeinarts.com",
		"password": "password123",
	}
	b, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed, status: %d", resp.StatusCode)
	}

	var loginData struct {
		Token string `json:"token"`
	}
	respBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &loginData)
	authHeader := "Bearer " + loginData.Token

	// 2. Test Accounts Receivable (AR): Buat Invoice Klien
	due := time.Now().AddDate(0, 0, 14)
	invoicePayload := map[string]interface{}{
		"client_name": "PT Mega Finansial Nusantara",
		"amount":      20000000,
		"tax_amount":  2200000,
		"currency":    "IDR",
		"due_date":    due.Format(time.RFC3339),
		"category":    "Project Revenue",
		"description": "Termin 1 (50%) Pengembangan Aplikasi Mobile Banking",
	}
	b, _ = json.Marshal(invoicePayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/ar/invoices", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create AR invoice failed, status: %d", resp.StatusCode)
	}

	var createdInvoice struct {
		Data model.Transaction `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &createdInvoice)

	if createdInvoice.Data.TransactionID == "" {
		t.Fatalf("TransactionID (UUID) harus terisi otomatis")
	}
	if createdInvoice.Data.Type != "Debit" {
		t.Fatalf("Expected Type 'Debit', got %s", createdInvoice.Data.Type)
	}
	if createdInvoice.Data.TotalAmount != 22200000 {
		t.Fatalf("Expected TotalAmount 22200000, got %.2f", createdInvoice.Data.TotalAmount)
	}
	if createdInvoice.Data.Status != "Sent" {
		t.Fatalf("Expected Status 'Sent', got %s", createdInvoice.Data.Status)
	}

	// 3. Test Pelunasan Invoice (Mark as Paid) via UUID
	payInvoiceReq := httptest.NewRequest(http.MethodPost, "/api/v1/finance/ar/invoices/"+createdInvoice.Data.TransactionID+"/pay", nil)
	payInvoiceReq.Header.Set("Authorization", authHeader)
	resp, err = app.Test(payInvoiceReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Pay AR invoice failed, status: %d", resp.StatusCode)
	}

	var paidInvoice struct {
		Data model.Transaction `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &paidInvoice)
	if paidInvoice.Data.Status != "Paid" || paidInvoice.Data.PaidAt == nil {
		t.Fatalf("Status invoice harus 'Paid' dan PaidAt tidak boleh nil")
	}

	// 4. Test Accounts Payable (AP): Catat Tagihan Vendor Pihak Ketiga (AWS / Software License)
	billPayload := map[string]interface{}{
		"vendor_name": "Google Cloud Platform (GCP)",
		"amount":      3500000,
		"tax_amount":  385000,
		"currency":    "IDR",
		"due_date":    due.Format(time.RFC3339),
		"category":    "Cloud Infrastructure",
		"description": "Tagihan Server Google Kubernetes Engine",
	}
	b, _ = json.Marshal(billPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/ap/bills", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create AP bill failed, status: %d", resp.StatusCode)
	}

	var createdBill struct {
		Data model.Transaction `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &createdBill)

	if createdBill.Data.Type != "Credit" {
		t.Fatalf("Expected Type 'Credit', got %s", createdBill.Data.Type)
	}

	// 5. Test Pembayaran Tagihan Vendor
	payBillReq := httptest.NewRequest(http.MethodPost, "/api/v1/finance/ap/bills/"+createdBill.Data.TransactionID+"/pay", nil)
	payBillReq.Header.Set("Authorization", authHeader)
	resp, err = app.Test(payBillReq)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Pay AP bill failed, status: %d", resp.StatusCode)
	}

	// 6. Test General Ledger: List Chart of Accounts
	db.Create(&model.Account{
		ID:      1,
		Code:    "1001",
		Name:    "Kas Utama",
		Type:    "ASSET",
		Balance: 10000000,
	})

	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/ledger/accounts", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get COA failed, status: %d", resp.StatusCode)
	}

	var coaResp struct {
		Data []model.Account `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &coaResp)
	if len(coaResp.Data) == 0 {
		t.Fatalf("Chart of Accounts tidak boleh kosong")
	}

	// 7. Test General Ledger: List Journal Entries
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/ledger/entries", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get journal entries failed, status: %d", resp.StatusCode)
	}

	var journalResp struct {
		Data []model.JournalEntry `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &journalResp)
	if len(journalResp.Data) == 0 {
		t.Fatalf("Entri jurnal otomatis harus terbuat dari pelunasan invoice/bill")
	}

	// 8. Test General Ledger: Manual Double-Entry Balance Check
	unbalancedJournal := map[string]interface{}{
		"description": "Jurnal tidak seimbang untuk pengujian",
		"items": []map[string]interface{}{
			{"account_code": "1001", "debit": 500000, "credit": 0},
			{"account_code": "4001", "debit": 0, "credit": 400000}, // Selisih 100.000
		},
	}
	b, _ = json.Marshal(unbalancedJournal)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/ledger/entries", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, _ = app.Test(req)
	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("Jurnal tidak seimbang harus ditolak dengan status 400, got %d", resp.StatusCode)
	}

	// 9. Test Financial Reports: Income Statement (Laba/Rugi)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/reports/income-statement", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Income statement failed, status: %d", resp.StatusCode)
	}

	var incStmt struct {
		Data model.IncomeStatement `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &incStmt)
	if incStmt.Data.TotalRevenue <= 0 {
		t.Fatalf("Total revenue laba rugi harus > 0")
	}

	// 10. Test Financial Reports: Balance Sheet (Neraca)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/reports/balance-sheet", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Balance sheet failed, status: %d", resp.StatusCode)
	}

	// 11. Test Financial Reports: Cash Flow (Arus Kas)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/reports/cash-flow", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Cash flow report failed, status: %d", resp.StatusCode)
	}

	// 12. Test Get Transaction by UUID
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/transactions/"+createdInvoice.Data.TransactionID, nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get transaction by UUID failed, status: %d", resp.StatusCode)
	}
}

func TestFinanceModulePDFSpecification(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// Login
	loginPayload := map[string]string{
		"email":    "admin@kafeinarts.com",
		"password": "password123",
	}
	b, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed, status: %d", resp.StatusCode)
	}

	var loginData struct {
		Token string `json:"token"`
	}
	respBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &loginData)
	authHeader := "Bearer " + loginData.Token

	// ── A. DATA MASTER ──────────────────────────────────────────────────────────

	// 1. Master Kategori
	catPayload := map[string]string{
		"category_name": "Event & Workshop Art",
		"cash_type":     "INCOME",
		"description":   "Tiket workshop latte art dan pameran desain",
	}
	b, _ = json.Marshal(catPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/categories", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Category failed, status: %d", resp.StatusCode)
	}

	// 2. Master Rekening / Kas Bank
	accPayload := map[string]interface{}{
		"bank_name":           "Bank BNI",
		"account_number":      "0887654321",
		"account_holder_name": "PT Kafeinarts Studio",
		"current_balance":     10000000,
	}
	b, _ = json.Marshal(accPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/accounts", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Account failed, status: %d", resp.StatusCode)
	}

	var createdAcc struct {
		Data model.MstFinanceAccount `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &createdAcc)

	// 3. Master Pajak
	taxPayload := map[string]interface{}{
		"tax_name":    "PPN 11%",
		"percentage":  11.0,
		"description": "Pajak Pertambahan Nilai 11%",
	}
	b, _ = json.Marshal(taxPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/taxes", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Tax failed, status: %d", resp.StatusCode)
	}

	var createdTax struct {
		Data model.MstFinanceTax `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &createdTax)

	// ── B. DATA TRANSAKSI ───────────────────────────────────────────────────────

	// 1. Jurnal Kas (Pencatatan Harian) & update saldo otomatis
	cashJournalPayload := map[string]interface{}{
		"account_id":       createdAcc.Data.ID,
		"category_id":      "cat-rev-01",
		"transaction_type": "INCOME",
		"amount":           5000000,
		"description":      "Penerimaan DP Proyek Desain Maskot",
	}
	b, _ = json.Marshal(cashJournalPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/cash-journals", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Cash Journal failed, status: %d", resp.StatusCode)
	}

	// Verifikasi saldo bertambah dari 10.000.000 menjadi 15.000.000
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/accounts", nil)
	req.Header.Set("Authorization", authHeader)
	resp, _ = app.Test(req)
	var accountsList struct {
		Data []model.MstFinanceAccount `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &accountsList)
	for _, a := range accountsList.Data {
		if a.ID == createdAcc.Data.ID {
			if a.CurrentBalance != 15000000 {
				t.Fatalf("Expected balance 15000000, got %.2f", a.CurrentBalance)
			}
		}
	}

	// 2. Manajemen Invoice / Piutang (AR) dengan Dynamic Grid Item
	taxPPN := createdTax.Data.ID
	invPayload := map[string]interface{}{
		"client_name":  "PT Inovasi Media Digital",
		"project_name": "Aplikasi POS Kasir Kafeinarts",
		"issue_date":   time.Now().Format("2006-01-02"),
		"due_date":     time.Now().AddDate(0, 0, 14).Format("2006-01-02"),
		"status":       "SENT",
		"notes":        "Termin 1 (60%)",
		"items": []map[string]interface{}{
			{
				"item_description": "Desain UI/UX Design System",
				"quantity":         1,
				"unit_price":       5000000,
				"tax_id":           taxPPN,
			},
			{
				"item_description": "Setup Database & Backend API",
				"quantity":         1,
				"unit_price":       7000000,
				"tax_id":           taxPPN,
			},
		},
	}
	b, _ = json.Marshal(invPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/invoices", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Invoice failed, status: %d", resp.StatusCode)
	}

	var createdInv struct {
		Data model.TrxInvoice `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &createdInv)

	if createdInv.Data.Subtotal != 12000000 {
		t.Fatalf("Expected Subtotal 12000000, got %.2f", createdInv.Data.Subtotal)
	}
	if createdInv.Data.TotalTax != 1320000 {
		t.Fatalf("Expected TotalTax 1320000, got %.2f", createdInv.Data.TotalTax)
	}
	if createdInv.Data.GrandTotal != 13320000 {
		t.Fatalf("Expected GrandTotal 13320000, got %.2f", createdInv.Data.GrandTotal)
	}

	// 3. Tagihan Vendor / Hutang (AP)
	billPayload := map[string]interface{}{
		"vendor_name": "Google Cloud Platform",
		"category_id": "cat-srv-01",
		"amount":      2800000,
		"due_date":    time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
		"status":      "PENDING",
	}
	b, _ = json.Marshal(billPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/vendor-bills", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create Vendor Bill failed, status: %d", resp.StatusCode)
	}

	// ── C. LAPORAN (REPORTS) ───────────────────────────────────────────────────

	// 1. Laporan Arus Kas (Cash Flow)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/reports/cash-flow", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Cash Flow report failed, status: %d", resp.StatusCode)
	}

	var cfReport struct {
		Data model.CashFlowReport `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &cfReport)
	if len(cfReport.Data.Rows) == 0 {
		t.Fatalf("Laporan arus kas harus berisi baris transaksi jurnal")
	}

	// 2. Laporan Laba/Rugi (Profit & Loss)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/reports/profit-loss?year=2026", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Profit Loss report failed, status: %d", resp.StatusCode)
	}

	var plReport struct {
		Data model.ProfitLossReport `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &plReport)
	if plReport.Data.TotalRevenue <= 0 {
		t.Fatalf("Total revenue Laba/Rugi harus > 0")
	}

	// 3. Laporan Piutang (AR Aging)
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/reports/ar-aging", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("AR Aging report failed, status: %d", resp.StatusCode)
	}

	var arReport struct {
		Data model.ARAgingReport `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &arReport)
	if len(arReport.Data.Items) == 0 {
		t.Fatalf("Laporan AR Aging harus memuat item piutang")
	}
}

func TestFrontendFinanceModuleCompatibility(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// 1. Auth Login
	loginPayload := map[string]string{
		"email":    "admin@kafeinarts.com",
		"password": "password123",
	}
	b, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed, status: %d", resp.StatusCode)
	}

	var loginData struct {
		Token string `json:"token"`
	}
	respBytes, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &loginData)
	authHeader := "Bearer " + loginData.Token

	// 2. Test Indonesian alias: /api/v1/finance/kategori with pagination
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/kategori?page=1&perPage=5", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get kategori failed, status: %d", resp.StatusCode)
	}

	var catResp struct {
		Data []map[string]interface{} `json:"data"`
		Meta map[string]interface{}   `json:"meta"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &catResp)
	if len(catResp.Data) == 0 {
		t.Fatalf("Expected category data in response")
	}
	if catResp.Meta["perPage"] == nil && catResp.Meta["per_page"] == nil {
		t.Fatalf("Expected pagination meta in response")
	}
	// Verify frontend camelCase alias fields exist
	if _, ok := catResp.Data[0]["namaKategori"]; !ok {
		t.Fatalf("Expected 'namaKategori' field for frontend compatibility")
	}

	// 3. Test Indonesian alias: /api/v1/finance/rekening
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/rekening?page=1&perPage=5", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get rekening failed, status: %d", resp.StatusCode)
	}

	var rekResp struct {
		Data []map[string]interface{} `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &rekResp)
	if len(rekResp.Data) == 0 {
		t.Fatalf("Expected accounts data")
	}
	if _, ok := rekResp.Data[0]["namaBank"]; !ok {
		t.Fatalf("Expected 'namaBank' field for frontend compatibility")
	}

	// 4. Test Indonesian alias: /api/v1/finance/pajak
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/pajak?page=1&perPage=5", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get pajak failed, status: %d", resp.StatusCode)
	}

	var taxResp struct {
		Data []map[string]interface{} `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &taxResp)
	if len(taxResp.Data) == 0 {
		t.Fatalf("Expected tax data")
	}
	if _, ok := taxResp.Data[0]["namaPajak"]; !ok {
		t.Fatalf("Expected 'namaPajak' field for frontend compatibility")
	}

	// 5. Test Jurnal Kas: Create with Indonesian fields, update and delete
	createdAccID := rekResp.Data[0]["id"].(string)
	createdCatID := catResp.Data[0]["id"].(string)

	jurnalPayload := map[string]interface{}{
		"tanggal":    "2026-09-05",
		"akunBank":   createdAccID,
		"tipe":       "Masuk",
		"kategori":   createdCatID,
		"nominal":    1750000,
		"keterangan": "Test penerimaan kas via alias FE",
	}
	b, _ = json.Marshal(jurnalPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/finance/jurnal", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Create jurnal failed, status: %d", resp.StatusCode)
	}

	var createdJurnal struct {
		Data map[string]interface{} `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &createdJurnal)
	jurnalID := createdJurnal.Data["id"].(string)

	// Update Jurnal Kas
	updateJurnalPayload := map[string]interface{}{
		"nominal":    2000000,
		"keterangan": "Test update nominal kas via alias FE",
	}
	b, _ = json.Marshal(updateJurnalPayload)
	req = httptest.NewRequest(http.MethodPut, "/api/v1/finance/jurnal/"+jurnalID, bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Update jurnal failed, status: %d", resp.StatusCode)
	}

	// 6. Test Reports with Indonesian endpoints and parameters
	// Cash Flow
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/laporan/cash-flow?periodeStart=2026-01-01&periodeEnd=2026-12-31", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get laporan cash-flow failed, status: %d", resp.StatusCode)
	}

	var cfResp struct {
		Data struct {
			Rows    []map[string]interface{} `json:"rows"`
			Summary map[string]interface{}   `json:"summary"`
		} `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &cfResp)
	if len(cfResp.Data.Rows) == 0 {
		t.Fatalf("Expected cash flow rows")
	}
	if _, ok := cfResp.Data.Summary["totalMasuk"]; !ok {
		t.Fatalf("Expected 'totalMasuk' in cash flow summary")
	}

	// Laba Rugi
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/laporan/laba-rugi?bulan=9&tahun=2026", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get laporan laba-rugi failed, status: %d", resp.StatusCode)
	}

	var lrResp struct {
		Data struct {
			Pendapatan  []map[string]interface{} `json:"pendapatan"`
			Pengeluaran []map[string]interface{} `json:"pengeluaran"`
			Summary     map[string]interface{}   `json:"summary"`
		} `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &lrResp)
	if _, ok := lrResp.Data.Summary["labaKotor"]; !ok {
		t.Fatalf("Expected 'labaKotor' in laba rugi summary")
	}

	// AR Aging
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/laporan/ar-aging", nil)
	req.Header.Set("Authorization", authHeader)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Get laporan ar-aging failed, status: %d", resp.StatusCode)
	}

	var arResp struct {
		Data struct {
			Rows    []map[string]interface{} `json:"rows"`
			Summary map[string]interface{}   `json:"summary"`
		} `json:"data"`
	}
	respBytes, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBytes, &arResp)
	if len(arResp.Data.Rows) == 0 {
		t.Fatalf("Expected ar aging rows")
	}
	if _, ok := arResp.Data.Summary["totalTagihan"]; !ok {
		t.Fatalf("Expected 'totalTagihan' in ar aging summary")
	}
}


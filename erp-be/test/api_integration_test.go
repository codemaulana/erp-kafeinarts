package test

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"erp-kafeinart/config"
	"erp-kafeinart/database"
	"erp-kafeinart/model"
	"erp-kafeinart/routes"

	"github.com/glebarez/sqlite"
	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func setupTestApp(t *testing.T) (*fiber.App, *gorm.DB, *config.Config) {
	cfg := &config.Config{
		Port:               "8000",
		AppName:            "Kafeinarts ERP Test",
		DBPath:             ":memory:",
		JWTSecret:          "test-jwt-secret-key-12345",
		JWTExpirationHours: 24,
		CORSOrigins:        []string{"*"},
	}

	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		t.Fatalf("Gagal membuka DB in-memory: %v", err)
	}

	err = db.AutoMigrate(
		&model.User{},
		&model.Product{},
		&model.StockAdjustment{},
		&model.Transaction{},
		&model.Employee{},
		&model.Account{},
		&model.JournalEntry{},
		&model.JournalItem{},
		&model.MstFinanceCategory{},
		&model.MstFinanceAccount{},
		&model.MstFinanceTax{},
		&model.TrxCashJournal{},
		&model.TrxInvoice{},
		&model.TrxInvoiceItem{},
		&model.TrxVendorBill{},
		// HR
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
		// SC
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
		t.Fatalf("Gagal migrate test DB: %v", err)
	}

	database.SeedDatabase(db)

	app := fiber.New()
	routes.SetupRoutes(app, db, cfg)

	return app, db, cfg
}

func TestFullAPIWorkflow(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// 1. Test Health Check
	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Health check failed, status: %d", resp.StatusCode)
	}

	// 2. Test Login Admin
	loginPayload := map[string]string{
		"email":    "admin@kafeinarts.com",
		"password": "password123",
	}
	body, _ := json.Marshal(loginPayload)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed, status: %d", resp.StatusCode)
	}

	respBody, _ := io.ReadAll(resp.Body)
	var loginData struct {
		Token string `json:"token"`
		User  struct {
			Email string `json:"email"`
			Role  string `json:"role"`
		} `json:"user"`
	}
	_ = json.Unmarshal(respBody, &loginData)

	if loginData.Token == "" {
		t.Fatalf("Token JWT tidak ditemukan di response login")
	}
	if loginData.User.Email != "admin@kafeinarts.com" {
		t.Fatalf("Expected email admin@kafeinarts.com, got %s", loginData.User.Email)
	}

	authToken := "Bearer " + loginData.Token

	// 3. Test Profile GET /api/v1/auth/me
	req = httptest.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", authToken)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Auth me failed, status: %d", resp.StatusCode)
	}

	// 4. Test Tambah Produk Baru POST /api/v1/inventory/products
	newProd := map[string]interface{}{
		"sku":         "TEST-PROD-99",
		"name":        "Test Coffee Capsule",
		"stock":       100,
		"price":       55000,
		"description": "Kapsul kopi untuk mesin espresso",
	}
	newProdBody, _ := json.Marshal(newProd)
	req = httptest.NewRequest(http.MethodPost, "/api/v1/inventory/products", bytes.NewReader(newProdBody))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", authToken)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusCreated {
		t.Fatalf("Tambah produk failed, status: %d", resp.StatusCode)
	}

	// 5. Test Inventory GET /api/v1/inventory/products
	req = httptest.NewRequest(http.MethodGet, "/api/v1/inventory/products?page=1&perPage=10", nil)
	req.Header.Set("Authorization", authToken)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Inventory products failed, status: %d", resp.StatusCode)
	}

	var productResp struct {
		Data []model.Product `json:"data"`
		Meta struct {
			Total int64 `json:"total"`
		} `json:"meta"`
	}
	respBody, _ = io.ReadAll(resp.Body)
	_ = json.Unmarshal(respBody, &productResp)

	if len(productResp.Data) == 0 || productResp.Meta.Total == 0 {
		t.Fatalf("Daftar produk tidak boleh kosong setelah produk dibuat")
	}

	// 6. Test Finance Summary GET /api/v1/finance/summary
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/summary", nil)
	req.Header.Set("Authorization", authToken)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Finance summary failed, status: %d", resp.StatusCode)
	}

	// 7. Test Finance Transactions GET /api/v1/finance/transactions
	req = httptest.NewRequest(http.MethodGet, "/api/v1/finance/transactions", nil)
	req.Header.Set("Authorization", authToken)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Finance transactions failed, status: %d", resp.StatusCode)
	}

	// 8. Test HR Employees GET /api/v1/hr/employees
	req = httptest.NewRequest(http.MethodGet, "/api/v1/hr/employees", nil)
	req.Header.Set("Authorization", authToken)
	resp, err = app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("HR employees failed, status: %d", resp.StatusCode)
	}
}

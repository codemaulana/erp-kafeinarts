package test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"erp-kafeinart/model"
)

func TestSCModuleComplete(t *testing.T) {
	app, _, _ := setupTestApp(t)

	// 1. Auth Admin Login
	loginPayload := map[string]string{
		"email":    "admin@kafeinarts.com",
		"password": "password123",
	}
	b, _ := json.Marshal(loginPayload)
	req := httptest.NewRequest(http.MethodPost, "/api/v1/auth/login", bytes.NewReader(b))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		t.Fatalf("Login failed: %d", resp.StatusCode)
	}

	var loginResp struct {
		Token string `json:"token"`
	}
	body, _ := io.ReadAll(resp.Body)
	_ = json.Unmarshal(body, &loginResp)
	authHeader := "Bearer " + loginResp.Token

	// Helper for HTTP requests
	doReq := func(method, url string, payload interface{}) (*http.Response, []byte) {
		var reqBody io.Reader
		if payload != nil {
			pj, _ := json.Marshal(payload)
			reqBody = bytes.NewReader(pj)
		}
		r := httptest.NewRequest(method, url, reqBody)
		r.Header.Set("Authorization", authHeader)
		if payload != nil {
			r.Header.Set("Content-Type", "application/json")
		}
		res, reqErr := app.Test(r)
		if reqErr != nil {
			t.Fatalf("Request %s %s failed: %v", method, url, reqErr)
		}
		data, _ := io.ReadAll(res.Body)
		return res, data
	}

	// ── A. MASTER DATA SUPPLY CHAIN ─────────────────────────────────────────────

	// 1. Vendor / Supplier
	vendorPayload := map[string]interface{}{
		"vendor_name":   "PT Grafika Printing Mandiri",
		"contact_name":  "Budi Santoso",
		"email":         "sales@grafikaprinting.co.id",
		"phone":         "08119876543",
		"address":       "Jl. Industri Percetakan No. 45, Jakarta Selatan",
		"payment_terms": "NET30",
		"status":        "ACTIVE",
	}
	res, data := doReq(http.MethodPost, "/api/v1/sc/vendors", vendorPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Vendor failed: %s", string(data))
	}
	var createdVendor struct {
		Data model.MstScVendor `json:"data"`
	}
	_ = json.Unmarshal(data, &createdVendor)
	vendorID := createdVendor.Data.ID

	// 2. Kategori Barang
	catPayload := map[string]interface{}{
		"category_name": "Kertas & Material Cetak",
		"description":   "Bahan baku art paper, kartu nama, packaging",
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/categories", catPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create SC Category failed: %s", string(data))
	}
	var createdCat struct {
		Data model.MstScItemCategory `json:"data"`
	}
	_ = json.Unmarshal(data, &createdCat)
	catID := createdCat.Data.ID

	// 3. Lokasi Penyimpanan / Gudang
	locPayload := map[string]interface{}{
		"location_name": "Gudang Utama Lantai 1",
		"address":       "Gedung Kafeinarts Creative, Ruang Penyimpanan A",
		"type":          "INTERNAL",
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/locations", locPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Location failed: %s", string(data))
	}
	var createdLoc struct {
		Data model.MstScLocation `json:"data"`
	}
	_ = json.Unmarshal(data, &createdLoc)
	locID := createdLoc.Data.ID

	// 4. Barang / Material
	itemPayload := map[string]interface{}{
		"item_name":          "Art Paper 260gsm A3+",
		"category_id":        catID,
		"default_location_id": locID,
		"unit":               "Rim",
		"unit_price":         175000,
		"cost_price":         150000,
		"stock_current":      50,
		"stock_minimum":      10,
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/items", itemPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Item failed: %s", string(data))
	}
	var createdItem struct {
		Data model.MstScItem `json:"data"`
	}
	_ = json.Unmarshal(data, &createdItem)
	itemID := createdItem.Data.ID

	// ── B. DATA TRANSAKSI SUPPLY CHAIN ──────────────────────────────────────────

	// 1. Purchase Request (Permintaan Pembelian)
	prPayload := map[string]interface{}{
		"requested_by":    "Admin Studio",
		"required_date":   time.Now().AddDate(0, 0, 7).Format("2006-01-02"),
		"priority":        "HIGH",
		"notes":           "Stok kertas menipis untuk cetak merchandise klien",
		"items": []map[string]interface{}{
			{
				"item_id":     itemID,
				"quantity":    20,
				"unit":        "Rim",
				"estimated_price": 150000,
				"notes":       "Bahan pesanan kaos & booklet",
			},
		},
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/purchase-requests", prPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Purchase Request failed: %s", string(data))
	}
	var createdPR struct {
		Data model.TrxScPurchaseRequest `json:"data"`
	}
	_ = json.Unmarshal(data, &createdPR)
	prID := createdPR.Data.ID

	// Approve PR
	res, data = doReq(http.MethodPatch, fmt.Sprintf("/api/v1/sc/purchase-requests/%s/status", prID), map[string]string{
		"status": "APPROVED",
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update PR Status failed: %s", string(data))
	}

	// 2. Purchase Order (PO)
	poPayload := map[string]interface{}{
		"vendor_id":             vendorID,
		"purchase_request_id":   prID,
		"expected_delivery_date": time.Now().AddDate(0, 0, 5).Format("2006-01-02"),
		"payment_terms":         "NET30",
		"notes":                 "Pengiriman langsung ke Gudang Utama",
		"items": []map[string]interface{}{
			{
				"item_id":    itemID,
				"quantity":   20,
				"unit":       "Rim",
				"unit_price": 150000,
				"tax_rate":   11, // 11% PPN
				"notes":      "20 Rim Art Paper 260gsm",
			},
		},
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/purchase-orders", poPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Purchase Order failed: %s", string(data))
	}
	var createdPO struct {
		Data model.TrxScPurchaseOrder `json:"data"`
	}
	_ = json.Unmarshal(data, &createdPO)
	poID := createdPO.Data.ID

	// Validate calculations:
	// Subtotal = 20 * 150,000 = 3,000,000
	// TaxAmount = 3,000,000 * 11% = 330,000
	// GrandTotal = 3,330,000
	if createdPO.Data.Subtotal != 3000000 {
		t.Fatalf("Expected PO subtotal 3000000, got %.2f", createdPO.Data.Subtotal)
	}
	if createdPO.Data.TaxAmount != 330000 {
		t.Fatalf("Expected PO tax 330000, got %.2f", createdPO.Data.TaxAmount)
	}
	if createdPO.Data.TotalAmount != 3330000 {
		t.Fatalf("Expected PO grand total 3330000, got %.2f", createdPO.Data.TotalAmount)
	}

	// Send PO
	res, data = doReq(http.MethodPatch, fmt.Sprintf("/api/v1/sc/purchase-orders/%s/status", poID), map[string]string{
		"status": "SENT",
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update PO Status failed: %s", string(data))
	}

	// 3. Goods Receipt (Penerimaan Barang)
	grPayload := map[string]interface{}{
		"purchase_order_id": poID,
		"vendor_id":         vendorID,
		"receipt_date":      time.Now().Format("2006-01-02"),
		"delivery_note_no":  "SJ-998811",
		"received_by":       "Admin Studio",
		"notes":             "Semua rim dalam kondisi baik tanpa cacat",
		"items": []map[string]interface{}{
			{
				"item_id":           itemID,
				"quantity_received": 20,
				"unit":              "Rim",
				"location_id":       locID,
				"condition":         "GOOD",
			},
		},
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/goods-receipts", grPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Goods Receipt failed: %s", string(data))
	}
	var createdGR struct {
		Data model.TrxScGoodsReceipt `json:"data"`
	}
	_ = json.Unmarshal(data, &createdGR)

	// 4. Asset Tetap
	assetPayload := map[string]interface{}{
		"asset_name":           "Wacom Cintiq Pro 24 Drawing Tablet",
		"category_id":          catID,
		"vendor_id":            vendorID,
		"location_id":          locID,
		"purchase_date":        "2026-01-15",
		"purchase_cost":        35000000,
		"useful_life_months":   36,
		"salvage_value":        5000000,
		"depreciation_method":  "STRAIGHT_LINE",
		"status":               "ACTIVE",
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/assets", assetPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Asset failed: %s", string(data))
	}
	var createdAsset struct {
		Data model.TrxScAsset `json:"data"`
	}
	_ = json.Unmarshal(data, &createdAsset)

	// 5. Software Subscription / Lisensi
	subPayload := map[string]interface{}{
		"app_name":         "Adobe Creative Cloud All Apps",
		"vendor_id":        vendorID,
		"billing_cycle":    "YEARLY",
		"cost":             12500000,
		"license_count":    5,
		"start_date":       "2026-01-01",
		"end_date":         "2027-01-01",
		"auto_renew":       true,
		"assigned_to":      "Tim Desain & Video",
		"status":           "ACTIVE",
	}
	res, data = doReq(http.MethodPost, "/api/v1/sc/subscriptions", subPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Subscription failed: %s", string(data))
	}

	// ── C. LAPORAN SUPPLY CHAIN ────────────────────────────────────────────────

	// 1. Riwayat Vendor
	res, data = doReq(http.MethodGet, fmt.Sprintf("/api/v1/sc/reports/vendor-history/%s", vendorID), nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Report vendor-history failed: %s", string(data))
	}

	// 2. Depresiasi Aset
	res, data = doReq(http.MethodGet, "/api/v1/sc/reports/asset-depreciation", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Report asset-depreciation failed: %s", string(data))
	}

	_ = createdGR
}

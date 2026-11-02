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

func TestHRModuleComplete(t *testing.T) {
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

	// ── A. MASTER DATA ─────────────────────────────────────────────────────────

	// 1. Departemen
	res, data := doReq(http.MethodPost, "/api/v1/hr/departments", map[string]interface{}{
		"department_name": "Creative & Design Studio",
		"code":            "DEP-DES",
		"description":     "Departemen Desain Grafis, Video, & UI/UX",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Department failed: %s", string(data))
	}
	var createdDept struct {
		Data model.MstHrDepartment `json:"data"`
	}
	_ = json.Unmarshal(data, &createdDept)
	deptID := createdDept.Data.ID

	// 2. Posisi / Jabatan
	res, data = doReq(http.MethodPost, "/api/v1/hr/positions", map[string]interface{}{
		"position_name": "Senior Graphic Designer",
		"code":          "POS-SGD",
		"department_id": deptID,
		"job_level":     "Senior",
		"description":   "Bertanggung jawab atas visual branding",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Position failed: %s", string(data))
	}
	var createdPos struct {
		Data model.MstHrPosition `json:"data"`
	}
	_ = json.Unmarshal(data, &createdPos)
	posID := createdPos.Data.ID

	// 3. Jenis Cuti
	res, data = doReq(http.MethodPost, "/api/v1/hr/leave-types", map[string]interface{}{
		"leave_name":         "Cuti Tahunan",
		"code":               "CUTI-THN",
		"annual_quota_days":  12,
		"is_paid":            true,
		"requires_doc":       false,
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Leave Type failed: %s", string(data))
	}
	var createdLeaveType struct {
		Data model.MstHrLeaveType `json:"data"`
	}
	_ = json.Unmarshal(data, &createdLeaveType)
	leaveTypeID := createdLeaveType.Data.ID

	// 4. Komponen Gaji
	res, data = doReq(http.MethodPost, "/api/v1/hr/salary-components", map[string]interface{}{
		"component_name": "Tunjangan Makan & Transport",
		"code":           "ALW-TRANS",
		"component_type": "EARNING",
		"is_taxable":     true,
		"default_amount": 500000,
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Salary Component failed: %s", string(data))
	}

	// 5. Shift Kerja
	res, data = doReq(http.MethodPost, "/api/v1/hr/work-shifts", map[string]interface{}{
		"shift_name":       "Shift Pagi Kreatif",
		"code":             "SFT-01",
		"start_time":       "08:30",
		"end_time":         "17:00",
		"grace_period_min": 15,
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Work Shift failed: %s", string(data))
	}

	// 6. Karyawan
	empPayload := map[string]interface{}{
		"full_name":       "Ahmad Fauzi",
		"email":           "ahmad.fauzi@kafeinarts.com",
		"phone":           "081234567890",
		"department_id":   deptID,
		"position_id":     posID,
		"employment_type": "PERMANENT",
		"join_date":       "2026-01-01",
		"basic_salary":    8000000,
		"bank_name":       "BCA",
		"bank_account_no": "1234567890",
		"bank_holder":     "Ahmad Fauzi",
		"status":          "ACTIVE",
	}
	res, data = doReq(http.MethodPost, "/api/v1/hr/employees", empPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Employee failed: %s", string(data))
	}
	var createdEmp struct {
		Data model.MstHrEmployee `json:"data"`
	}
	_ = json.Unmarshal(data, &createdEmp)
	empID := createdEmp.Data.ID
	if createdEmp.Data.EmployeeNumber == "" {
		t.Fatalf("Expected auto-generated employee_number, got empty")
	}

	// Verify List Employees
	res, data = doReq(http.MethodGet, "/api/v1/hr/employees", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Get Employees failed: %s", string(data))
	}

	// ── B. TRANSAKSI HR ────────────────────────────────────────────────────────

	// 1. Presensi: Clock In
	res, data = doReq(http.MethodPost, "/api/v1/hr/presensi/clock-in", map[string]interface{}{
		"employee_id": empID,
		"date":        time.Now().Format("2006-01-02"),
		"time_in":     "08:25",
		"location_in": "Studio Utama Lantai 2",
		"notes":       "Tepat waktu",
	})
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Clock In failed: %s", string(data))
	}
	var clockInResp struct {
		Data model.TrxHrAttendance `json:"data"`
	}
	_ = json.Unmarshal(data, &clockInResp)
	attID := clockInResp.Data.ID

	// Presensi: Clock Out
	res, data = doReq(http.MethodPost, "/api/v1/hr/presensi/clock-out", map[string]interface{}{
		"employee_id":  empID,
		"date":         time.Now().Format("2006-01-02"),
		"time_out":     "17:15",
		"location_out": "Studio Utama Lantai 2",
		"notes":        "Pekerjaan desain selesai",
	})
	if res.StatusCode != http.StatusOK && res.StatusCode != http.StatusCreated {
		t.Fatalf("Clock Out failed: %s", string(data))
	}

	// 2. Permohonan Cuti
	leavePayload := map[string]interface{}{
		"employee_id":   empID,
		"leave_type_id": leaveTypeID,
		"start_date":    time.Now().AddDate(0, 0, 3).Format("2006-01-02"),
		"end_date":      time.Now().AddDate(0, 0, 5).Format("2006-01-02"),
		"total_days":    3,
		"reason":        "Keperluan urusan keluarga di kampung halaman",
	}
	res, data = doReq(http.MethodPost, "/api/v1/hr/cuti", leavePayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Leave Request failed: %s", string(data))
	}
	var leaveResp struct {
		Data model.TrxHrLeaveRequest `json:"data"`
	}
	_ = json.Unmarshal(data, &leaveResp)
	leaveReqID := leaveResp.Data.ID

	// Update Status Cuti: APPROVED
	res, data = doReq(http.MethodPatch, fmt.Sprintf("/api/v1/hr/cuti/%s/status", leaveReqID), map[string]string{
		"status": "APPROVED",
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update Cuti Status failed: %s", string(data))
	}

	// 3. Penggajian / Payroll
	payrollPayload := map[string]interface{}{
		"employee_id":    empID,
		"period_month":   int(time.Now().Month()),
		"period_year":    time.Now().Year(),
		"basic_salary":   8000000,
		"allowances":     1000000,
		"overtime_pay":   250000,
		"bonus":          500000,
		"deductions":     200000,
		"notes":          "Gaji bulan berjalan + insentif lembur sprint",
		"items": []map[string]interface{}{
			{
				"component_name": "Gaji Pokok",
				"item_type":      "EARNING",
				"amount":         8000000,
			},
			{
				"component_name": "Tunjangan Transport",
				"item_type":      "EARNING",
				"amount":         1000000,
			},
			{
				"component_name": "BPJS Ketenagakerjaan",
				"item_type":      "DEDUCTION",
				"amount":         200000,
			},
		},
	}
	res, data = doReq(http.MethodPost, "/api/v1/hr/payroll", payrollPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Payroll failed: %s", string(data))
	}
	var payResp struct {
		Data model.TrxHrPayroll `json:"data"`
	}
	_ = json.Unmarshal(data, &payResp)
	payID := payResp.Data.ID
	// Check Net Salary: 8,000,000 + 1,000,000 + 250,000 + 500,000 - 200,000 = 9,550,000
	if payResp.Data.NetSalary != 9550000 {
		t.Fatalf("Expected NetSalary 9550000, got %.2f", payResp.Data.NetSalary)
	}

	// Update Status Payroll: PAID
	res, data = doReq(http.MethodPatch, fmt.Sprintf("/api/v1/hr/payroll/%s/status", payID), map[string]string{
		"status": "PAID",
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update Payroll Status failed: %s", string(data))
	}

	// 4. Klaim Reimbursement
	reimbPayload := map[string]interface{}{
		"employee_id": empID,
		"title":       "Pembelian Kertas Gambar & Stylus Pen",
		"amount":      350000,
		"claim_date":  time.Now().Format("2006-01-02"),
		"description": "Perlengkapan sketching proyek Kafeinarts",
	}
	res, data = doReq(http.MethodPost, "/api/v1/hr/reimbursement", reimbPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Reimbursement failed: %s", string(data))
	}
	var reimbResp struct {
		Data model.TrxHrReimbursement `json:"data"`
	}
	_ = json.Unmarshal(data, &reimbResp)
	reimbID := reimbResp.Data.ID

	// Update Status Reimbursement: APPROVED
	res, data = doReq(http.MethodPatch, fmt.Sprintf("/api/v1/hr/reimbursement/%s/status", reimbID), map[string]string{
		"status": "APPROVED",
	})
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Update Reimbursement Status failed: %s", string(data))
	}

	// 5. Penilaian Kinerja / Review
	revPayload := map[string]interface{}{
		"employee_id":   empID,
		"review_period": "Q3 2026",
		"review_date":   time.Now().Format("2006-01-02"),
		"feedback":      "Kualitas output desain sangat estetis dan tepat waktu",
		"goals_next":    "Pimpin tim UI/UX untuk proyek POS",
		"items": []map[string]interface{}{
			{
				"criteria_name": "Kualitas Desain & Kreativitas",
				"weight":        40,
				"score":         90,
				"notes":         "Sangat konsisten",
			},
			{
				"criteria_name": "Ketepatan Waktu Delivery",
				"weight":        30,
				"score":         85,
				"notes":         "Selalu selesai sebelum deadline",
			},
			{
				"criteria_name": "Kolaborasi & Komunikasi Tim",
				"weight":        30,
				"score":         85,
				"notes":         "Responsif dalam grup",
			},
		},
	}
	res, data = doReq(http.MethodPost, "/api/v1/hr/performance-reviews", revPayload)
	if res.StatusCode != http.StatusCreated {
		t.Fatalf("Create Performance Review failed: %s", string(data))
	}
	var revResp struct {
		Data model.TrxHrPerformanceReview `json:"data"`
	}
	_ = json.Unmarshal(data, &revResp)
	// Weighted score: 90*0.4 + 85*0.3 + 85*0.3 = 36 + 25.5 + 25.5 = 87
	if revResp.Data.TotalScore != 87 {
		t.Fatalf("Expected overall score 87, got %.2f", revResp.Data.TotalScore)
	}

	// ── C. LAPORAN HR ──────────────────────────────────────────────────────────

	// 1. Rekap Absensi
	res, data = doReq(http.MethodGet, "/api/v1/hr/reports/rekap-absensi", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Report rekap-absensi failed: %s", string(data))
	}

	// 2. Slip Gaji
	res, data = doReq(http.MethodGet, fmt.Sprintf("/api/v1/hr/reports/slip-gaji/%s", payID), nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Report slip-gaji failed: %s", string(data))
	}

	// 3. Saldo Cuti
	res, data = doReq(http.MethodGet, "/api/v1/hr/reports/saldo-cuti", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Report saldo-cuti failed: %s", string(data))
	}

	// 4. Turnover
	res, data = doReq(http.MethodGet, "/api/v1/hr/reports/turnover", nil)
	if res.StatusCode != http.StatusOK {
		t.Fatalf("Report turnover failed: %s", string(data))
	}

	_ = attID
}

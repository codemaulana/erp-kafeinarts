package model

// ── HR Reports ───────────────────────────────────────────────────

type RekapAbsensiSummary struct {
	TotalHadir     int `json:"totalHadir"`
	TotalTerlambat int `json:"totalTerlambat"`
	TotalAlpha     int `json:"totalAlpha"`
	TotalCuti      int `json:"totalCuti"`
}

type RekapAbsensiRow struct {
	No             int    `json:"no"`
	EmployeeNumber string `json:"employee_number"`
	EmployeeName   string `json:"employee_name"`
	TotalHadir     int    `json:"total_hadir"`
	TotalTerlambat int    `json:"total_terlambat"`
	TotalAlpha     int    `json:"total_alpha"`
	TotalCuti      int    `json:"total_cuti"`
	TotalSakit     int    `json:"total_sakit"`
}

type RekapAbsensiReport struct {
	Summary RekapAbsensiSummary `json:"summary"`
	Rows    []RekapAbsensiRow   `json:"rows"`
}

type SlipGajiSummary struct {
	TotalPayroll    int     `json:"totalPayroll"`
	TotalPendapatan float64 `json:"totalPendapatan"`
	TotalPotongan   float64 `json:"totalPotongan"`
	TotalGajiBersih float64 `json:"totalGajiBersih"`
}

type SlipGajiRow struct {
	PayrollNumber   string  `json:"payroll_number"`
	EmployeeNumber  string  `json:"employee_number"`
	EmployeeName    string  `json:"employee_name"`
	DepartmentName  string  `json:"department_name"`
	PayrollPeriod   string  `json:"payroll_period"`
	TotalEarnings   float64 `json:"total_earnings"`
	TotalDeductions float64 `json:"total_deductions"`
	NetSalary       float64 `json:"net_salary"`
	PaymentStatus   string  `json:"payment_status"`
}

type SlipGajiReport struct {
	Summary SlipGajiSummary `json:"summary"`
	Rows    []SlipGajiRow   `json:"rows"`
}

type SaldoCutiSummary struct {
	TotalKaryawan    int `json:"totalKaryawan"`
	RataRataSisaCuti int `json:"rataRataSisaCuti"`
	CutiTerbanyak    int `json:"cutiTerbanyak"`
	CutiTerendah     int `json:"cutiTerendah"`
}

type SaldoCutiRow struct {
	EmployeeNumber string `json:"employee_number"`
	EmployeeName   string `json:"employee_name"`
	DepartmentName string `json:"department_name"`
	TotalQuota     int    `json:"total_quota"`
	UsedQuota      int    `json:"used_quota"`
	RemainingQuota int    `json:"remaining_quota"`
}

type SaldoCutiReport struct {
	Summary SaldoCutiSummary `json:"summary"`
	Rows    []SaldoCutiRow   `json:"rows"`
}

type TurnoverSummary struct {
	TotalKaryawanAwal int     `json:"totalKaryawanAwal"`
	KaryawanMasuk     int     `json:"karyawanMasuk"`
	KaryawanKeluar    int     `json:"karyawanKeluar"`
	TurnoverRate      float64 `json:"turnoverRate"`
}

type TurnoverRow struct {
	Period     string  `json:"period"`
	TotalAwal  int     `json:"total_awal"`
	Masuk      int     `json:"masuk"`
	Keluar     int     `json:"keluar"`
	TotalAkhir int     `json:"total_akhir"`
	Rate       float64 `json:"rate"`
}

type TurnoverReport struct {
	Summary TurnoverSummary `json:"summary"`
	Rows    []TurnoverRow   `json:"rows"`
}

// ── SC Reports ───────────────────────────────────────────────────

type RiwayatVendorRow struct {
	ID            string  `json:"id,omitempty"`
	VendorID      string  `json:"vendor_id"`
	VendorName    string  `json:"vendor_name"`
	PoNumber      string  `json:"po_number"`
	TagihanNumber string  `json:"tagihan_number"`
	Date          string  `json:"date"`
	TotalAmount   float64 `json:"total_amount"`
	PaymentStatus string  `json:"payment_status"`
	Status        string  `json:"status,omitempty"`

	TotalPO       int     `json:"total_po,omitempty"`
	CompletedPO   int     `json:"completed_po,omitempty"`
	LastOrderDate string  `json:"last_order_date,omitempty"`
}

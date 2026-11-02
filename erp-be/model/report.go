package model

import "time"

// CategoryBreakdown memuat total amount per kategori keuangan.
type CategoryBreakdown struct {
	Category string  `json:"category"`
	Kategori string  `json:"kategori"`
	Amount   float64 `json:"amount"`
	Total    float64 `json:"total"`
}

// CashFlowItem rincian arus kas per kategori
type CashFlowItem struct {
	Category string  `json:"category"`
	Amount   float64 `json:"amount"`
}

// CashFlowReportRow adalah baris pada Tabel Laporan Arus Kas (Cash Flow).
// Kolom: Tanggal, No. Referensi, Keterangan, Kategori, Masuk, Keluar, Saldo Akhir.
type CashFlowReportRow struct {
	ID            string    `json:"id"`
	Date          time.Time `json:"date"`           // Tanggal
	Tanggal       string    `json:"tanggal"`
	ReferenceNo   string    `json:"reference_no"`   // No. Referensi
	Description   string    `json:"description"`    // Keterangan
	Keterangan    string    `json:"keterangan"`
	Category      string    `json:"category"`       // Kategori
	NamaKategori  string    `json:"namaKategori"`
	Kategori      string    `json:"kategori"`
	AccountID     string    `json:"account_id,omitempty"`
	AkunBank      string    `json:"akunBank,omitempty"`
	AccountName   string    `json:"account_name,omitempty"`
	NamaBank      string    `json:"namaBank,omitempty"`
	Inflow        float64   `json:"inflow"`         // Masuk
	Outflow       float64   `json:"outflow"`        // Keluar
	Nominal       float64   `json:"nominal"`
	Tipe          string    `json:"tipe"`           // 'Masuk', 'Keluar'
	EndingBalance float64   `json:"ending_balance"` // Saldo Akhir
}

type CashFlowSummary struct {
	TotalMasuk  float64 `json:"totalMasuk"`
	TotalKeluar float64 `json:"totalKeluar"`
	SaldoBersih float64 `json:"saldoBersih"`
}

// CashFlowReport memuat seluruh data laporan arus kas.
type CashFlowReport struct {
	StartDate        time.Time           `json:"start_date"`
	EndDate          time.Time           `json:"end_date"`
	PeriodStart      time.Time           `json:"period_start,omitempty"`
	PeriodEnd        time.Time           `json:"period_end,omitempty"`
	AccountID        string              `json:"account_id,omitempty"`
	AccountName      string              `json:"account_name,omitempty"`
	BeginningCash    float64             `json:"beginning_cash"`
	TotalInflow      float64             `json:"total_inflow"`
	TotalInflows     float64             `json:"total_inflows,omitempty"`
	TotalOutflow     float64             `json:"total_outflow"`
	TotalOutflows    float64             `json:"total_outflows,omitempty"`
	OperatingInflows []CashFlowItem      `json:"operating_inflows,omitempty"`
	OperatingOutflow []CashFlowItem      `json:"operating_outflows,omitempty"`
	NetCashFlow      float64             `json:"net_cash_flow"`
	EndingCash       float64             `json:"ending_cash"`
	Rows             []CashFlowReportRow `json:"rows"`
	Summary          CashFlowSummary     `json:"summary"`
}

type ProfitLossSummary struct {
	TotalPendapatan  float64 `json:"totalPendapatan"`
	TotalPengeluaran float64 `json:"totalPengeluaran"`
	LabaKotor        float64 `json:"labaKotor"`
	Bulan            int     `json:"bulan"`
	Tahun            int     `json:"tahun"`
}

// ProfitLossReport adalah Laporan Laba/Rugi (Profit & Loss).
// Filter: Bulan, Tahun.
// Output: Ringkasan bertingkat (Total Pendapatan, Total Pengeluaran, Laba/Rugi Bersih).
type ProfitLossReport struct {
	Month          int                 `json:"month"`
	Year           int                 `json:"year"`
	TotalRevenue   float64             `json:"total_revenue"`   // Total Pendapatan
	TotalExpense   float64             `json:"total_expense"`   // Total Pengeluaran
	NetProfitLoss  float64             `json:"net_profit_loss"` // Laba/Rugi Bersih (Total Pendapatan - Total Pengeluaran)
	RevenueDetails []CategoryBreakdown `json:"revenue_details"`
	ExpenseDetails []CategoryBreakdown `json:"expense_details"`
	Pendapatan     []CategoryBreakdown `json:"pendapatan"`
	Pengeluaran    []CategoryBreakdown `json:"pengeluaran"`
	Summary        ProfitLossSummary   `json:"summary"`
}

// ARAgingReportItem adalah baris pada Tabel Laporan Piutang (AR Aging).
// Kolom: No. Invoice, Klien, Jatuh Tempo, Keterlambatan (Hari), Total Tagihan, Sisa Belum Dibayar.
type ARAgingReportItem struct {
	ID            string    `json:"id"`
	InvoiceID     string    `json:"invoice_id"`
	InvoiceNumber string    `json:"invoice_number"`  // No. Invoice
	NoInvoice     string    `json:"noInvoice"`
	ClientName    string    `json:"client_name"`     // Klien
	Klien         string    `json:"klien"`
	ProjectName   string    `json:"project_name"`    // Proyek
	Proyek        string    `json:"proyek"`
	IssueDate     time.Time `json:"issue_date"`      // Tgl Terbit
	TglTerbit     string    `json:"tglTerbit"`
	DueDate       time.Time `json:"due_date"`        // Jatuh Tempo
	TglJatuhTempo string    `json:"tglJatuhTempo"`
	DaysOverdue   int       `json:"days_overdue"`    // Keterlambatan (Hari) (0 jika belum lewat due date)
	GrandTotal    float64   `json:"grand_total"`     // Total Tagihan
	GrandTotalAlt float64   `json:"grandTotal"`
	RemainingDue  float64   `json:"remaining_due"`   // Sisa Belum Dibayar
	Status        string    `json:"status"`          // Status Pembayaran
}

type ARAgingSummary struct {
	TotalTagihan    float64 `json:"totalTagihan"`
	TotalLunas      float64 `json:"totalLunas"`
	TotalBelumLunas float64 `json:"totalBelumLunas"`
}

// ARAgingReport memuat seluruh data laporan piutang umur tagihan.
type ARAgingReport struct {
	AsOfDate          time.Time           `json:"as_of_date"`
	TotalReceivable   float64             `json:"total_receivable"`
	TotalOverdue      float64             `json:"total_overdue"`
	TotalInvoiceCount int                 `json:"total_invoice_count"`
	Items             []ARAgingReportItem `json:"items"`
	Rows              []ARAgingReportItem `json:"rows"`
	Summary           ARAgingSummary      `json:"summary"`
}

// IncomeStatement merepresentasikan Laporan Laba/Rugi rentang tanggal fleksibel (tambahan)
type IncomeStatement struct {
	PeriodStart  time.Time           `json:"period_start"`
	PeriodEnd    time.Time           `json:"period_end"`
	TotalRevenue float64             `json:"total_revenue"`
	TotalExpense float64             `json:"total_expense"`
	TotalTax     float64             `json:"total_tax"`
	NetIncome    float64             `json:"net_income"`
	Revenues     []CategoryBreakdown `json:"revenues"`
	Expenses     []CategoryBreakdown `json:"expenses"`
}

// BalanceSheetItem rincian pos dalam neraca.
type BalanceSheetItem struct {
	Code   string  `json:"code"`
	Name   string  `json:"name"`
	Amount float64 `json:"amount"`
}

// BalanceSheet merepresentasikan Laporan Neraca (Posisi Keuangan).
type BalanceSheet struct {
	AsOfDate         time.Time          `json:"as_of_date"`
	Assets           []BalanceSheetItem `json:"assets"`
	TotalAssets      float64            `json:"total_assets"`
	Liabilities      []BalanceSheetItem `json:"liabilities"`
	TotalLiabilities float64            `json:"total_liabilities"`
	Equity           []BalanceSheetItem `json:"equity"`
	TotalEquity      float64            `json:"total_equity"`
	IsBalanced       bool               `json:"is_balanced"`
}

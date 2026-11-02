package service

import (
	"time"

	"erp-kafeinart/model"
)

// ── Laporan (Reports) Implementasi ──────────────────────────────────────────

func (s *financeService) GetCashFlowReport(start, end time.Time, accountID string) (*model.CashFlowReport, error) {
	if start.IsZero() {
		start = time.Date(time.Now().Year(), 1, 1, 0, 0, 0, 0, time.Local)
	}
	if end.IsZero() {
		end = time.Now().AddDate(0, 0, 1)
	}
	return s.trxRepo.GetCashFlowReport(start, end, accountID)
}

func (s *financeService) GetProfitLossReport(month, year int) (*model.ProfitLossReport, error) {
	if year <= 0 {
		year = time.Now().Year()
	}
	return s.trxRepo.GetProfitLossReport(month, year)
}

func (s *financeService) GetARAgingReport(status, clientName string) (*model.ARAgingReport, error) {
	return s.trxRepo.GetARAgingReport(status, clientName)
}

package service

import (
	"erp-kafeinart/model"
)

// ── Laporan HR ──────────────────────────────────────────────────────

func (s *hrService) GetRekapAbsensi(startDate, endDate, departmentID, employeeID string) (*model.RekapAbsensiReport, error) {
	return s.repo.GetRekapAbsensi(startDate, endDate, departmentID, employeeID)
}

func (s *hrService) GetSlipGaji(period, employeeID, departmentID string) (*model.SlipGajiReport, error) {
	return s.repo.GetSlipGaji(period, employeeID, departmentID)
}

func (s *hrService) GetSaldoCuti(departmentID, employeeID string) (*model.SaldoCutiReport, error) {
	return s.repo.GetSaldoCuti(departmentID, employeeID)
}

func (s *hrService) GetTurnover(year int) (*model.TurnoverReport, error) {
	return s.repo.GetTurnover(year)
}

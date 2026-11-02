package service

import (
	"erp-kafeinart/model"
)

// ── Laporan Supply Chain ────────────────────────────────────────────

func (s *scService) GetLaporanPO(startDate, endDate, vendorID, status string) ([]model.TrxScPurchaseOrder, error) {
	return s.repo.GetLaporanPO(startDate, endDate, vendorID, status)
}

func (s *scService) GetLaporanAset(categoryID, locationID, status, departmentID, employeeID string) ([]model.TrxScAsset, error) {
	return s.repo.GetLaporanAset(categoryID, locationID, status, departmentID, employeeID)
}

func (s *scService) GetLaporanLisensi(status, vendorID string) ([]model.TrxScSubscription, error) {
	return s.repo.GetLaporanLisensi(status, vendorID)
}

func (s *scService) GetRiwayatVendor(vendorID, startDate, endDate string) ([]model.RiwayatVendorRow, error) {
	return s.repo.GetRiwayatVendor(vendorID, startDate, endDate)
}

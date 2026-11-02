package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
)

// ── Transaksi Manajemen Aset ────────────────────────────────────────

func (s *scService) GetAssetsPaginated(page, perPage int, categoryID, locationID, status, search string) ([]model.TrxScAsset, int64, error) {
	return s.repo.FindAssetsPaginated(page, perPage, categoryID, locationID, status, search)
}

func (s *scService) GetAssetByID(id string) (*model.TrxScAsset, error) {
	return s.repo.FindAssetByID(id)
}

func (s *scService) CreateAsset(dto *AssetDTO) (*model.TrxScAsset, error) {
	name := dto.Name
	if name == "" {
		name = dto.AssetName
	}
	if name == "" && dto.ItemID == "" {
		return nil, errors.New("nama aset atau item aset wajib diisi")
	}

	code := dto.AssetCode
	if code == "" {
		code = fmt.Sprintf("AST-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	serial := dto.SerialOrLicenseKey
	if serial == "" {
		serial = dto.SerialNumber
	}

	acqDate := dto.AcquisitionDate
	if acqDate == "" {
		acqDate = dto.PurchaseDate
	}
	if acqDate == "" {
		acqDate = time.Now().Format("2006-01-02")
	}

	val := dto.AcquisitionValue
	if val == 0 {
		val = dto.PurchasePrice
	}
	if val == 0 {
		val = dto.PurchaseCost
	}

	status := dto.Status
	if status == "" {
		status = "AVAILABLE"
	}

	cond := dto.Condition
	if cond == "" {
		cond = "GOOD"
	}

	assigned := dto.AssignedToEmployeeID
	if assigned == nil {
		assigned = dto.AssignedTo
	}

	asset := &model.TrxScAsset{
		ID:                      uuid.NewString(),
		AssetCode:               code,
		Name:                    name,
		ItemID:                  dto.ItemID,
		CategoryID:              dto.CategoryID,
		SerialOrLicenseKey:      serial,
		AcquisitionDate:         acqDate,
		AcquisitionValue:        val,
		WarrantyOrLicenseExpiry: dto.WarrantyOrLicenseExpiry,
		AssignedToEmployeeID:    assigned,
		DepartmentID:            dto.DepartmentID,
		LocationID:              dto.LocationID,
		Status:                  status,
		Condition:               cond,
		Notes:                   dto.Notes,
		CreatedAt:               time.Now(),
		UpdatedAt:               time.Now(),
	}
	if err := s.repo.CreateAsset(asset); err != nil {
		return nil, err
	}
	return s.repo.FindAssetByID(asset.ID)
}

func (s *scService) UpdateAsset(id string, dto *AssetDTO) (*model.TrxScAsset, error) {
	existing, err := s.repo.FindAssetByID(id)
	if err != nil {
		return nil, errors.New("aset tidak ditemukan")
	}

	if dto.ItemID != "" {
		existing.ItemID = dto.ItemID
	}
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.CategoryID != "" {
		existing.CategoryID = dto.CategoryID
	}
	serial := dto.SerialOrLicenseKey
	if serial == "" {
		serial = dto.SerialNumber
	}
	if serial != "" {
		existing.SerialOrLicenseKey = serial
	}
	acqDate := dto.AcquisitionDate
	if acqDate == "" {
		acqDate = dto.PurchaseDate
	}
	if acqDate != "" {
		existing.AcquisitionDate = acqDate
	}
	val := dto.AcquisitionValue
	if val == 0 {
		val = dto.PurchasePrice
	}
	if val > 0 {
		existing.AcquisitionValue = val
	}
	if dto.WarrantyOrLicenseExpiry != "" {
		existing.WarrantyOrLicenseExpiry = dto.WarrantyOrLicenseExpiry
	}
	assigned := dto.AssignedToEmployeeID
	if assigned == nil {
		assigned = dto.AssignedTo
	}
	existing.AssignedToEmployeeID = assigned
	if dto.DepartmentID != nil {
		existing.DepartmentID = dto.DepartmentID
	}
	if dto.LocationID != nil {
		existing.LocationID = dto.LocationID
	}
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	if dto.Condition != "" {
		existing.Condition = dto.Condition
	}
	existing.Notes = dto.Notes
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateAsset(existing); err != nil {
		return nil, err
	}
	return s.repo.FindAssetByID(id)
}

func (s *scService) UpdateAssetStatus(id, status string) error {
	return s.repo.UpdateAssetStatus(id, status)
}

func (s *scService) DeleteAsset(id string) error {
	return s.repo.DeleteAsset(id)
}

// ── Transaksi Perpanjangan Langganan ────────────────────────────────

func (s *scService) GetSubscriptionsPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScSubscription, int64, error) {
	return s.repo.FindSubscriptionsPaginated(page, perPage, vendorID, status, search)
}

func (s *scService) GetSubscriptionByID(id string) (*model.TrxScSubscription, error) {
	return s.repo.FindSubscriptionByID(id)
}

func (s *scService) CreateSubscription(dto *SubscriptionDTO) (*model.TrxScSubscription, error) {
	name := dto.Name
	if name == "" {
		name = dto.AppName
	}
	if name == "" && dto.ItemID == "" {
		return nil, errors.New("nama aplikasi atau item langganan wajib diisi")
	}
	if dto.VendorID == "" {
		return nil, errors.New("vendor wajib dipilih")
	}

	code := dto.LicenseCode
	if code == "" {
		code = fmt.Sprintf("LIC-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	cycle := dto.BillingCycle
	if cycle == "" {
		cycle = "MONTHLY"
	}

	amount := dto.Amount
	if amount == 0 {
		amount = dto.Cost
	}

	startDate := dto.StartDate
	if startDate == "" {
		startDate = time.Now().Format("2006-01-02")
	}

	endDate := dto.EndDate
	if endDate == "" {
		endDate = dto.ExpiryDate
	}
	if endDate == "" {
		endDate = time.Now().AddDate(0, 1, 0).Format("2006-01-02")
	}

	status := dto.Status
	if status == "" {
		status = "ACTIVE"
	}

	autoRenew := dto.AutoRenewal || dto.AutoRenew

	sub := &model.TrxScSubscription{
		ID:           uuid.NewString(),
		LicenseCode:  code,
		Name:         name,
		ItemID:       dto.ItemID,
		VendorID:     dto.VendorID,
		BillingCycle: cycle,
		Amount:       amount,
		StartDate:    startDate,
		EndDate:      endDate,
		AutoRenewal:  autoRenew,
		Status:       status,
		Notes:        dto.Notes,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}
	if err := s.repo.CreateSubscription(sub); err != nil {
		return nil, err
	}
	return s.repo.FindSubscriptionByID(sub.ID)
}

func (s *scService) UpdateSubscription(id string, dto *SubscriptionDTO) (*model.TrxScSubscription, error) {
	existing, err := s.repo.FindSubscriptionByID(id)
	if err != nil {
		return nil, errors.New("langganan tidak ditemukan")
	}

	if dto.ItemID != "" {
		existing.ItemID = dto.ItemID
	}
	if dto.VendorID != "" {
		existing.VendorID = dto.VendorID
	}
	if dto.Name != "" {
		existing.Name = dto.Name
	}
	if dto.BillingCycle != "" {
		existing.BillingCycle = dto.BillingCycle
	}
	amount := dto.Amount
	if amount == 0 {
		amount = dto.Cost
	}
	if amount > 0 {
		existing.Amount = amount
	}
	if dto.StartDate != "" {
		existing.StartDate = dto.StartDate
	}
	endDate := dto.EndDate
	if endDate == "" {
		endDate = dto.ExpiryDate
	}
	if endDate != "" {
		existing.EndDate = endDate
	}
	existing.AutoRenewal = dto.AutoRenewal
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	existing.Notes = dto.Notes
	existing.UpdatedAt = time.Now()

	if err := s.repo.UpdateSubscription(existing); err != nil {
		return nil, err
	}
	return s.repo.FindSubscriptionByID(id)
}

func (s *scService) UpdateSubscriptionStatus(id, status string) error {
	return s.repo.UpdateSubscriptionStatus(id, status)
}

func (s *scService) DeleteSubscription(id string) error {
	return s.repo.DeleteSubscription(id)
}

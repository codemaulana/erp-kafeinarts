package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
)

// ── Transaksi Purchase Request ──────────────────────────────────────

func (s *scService) GetPurchaseRequestsPaginated(page, perPage int, departmentID, status, search string) ([]model.TrxScPurchaseRequest, int64, error) {
	return s.repo.FindPurchaseRequestsPaginated(page, perPage, departmentID, status, search)
}

func (s *scService) GetPurchaseRequestByID(id string) (*model.TrxScPurchaseRequest, error) {
	return s.repo.FindPurchaseRequestByID(id)
}

func (s *scService) CreatePurchaseRequest(dto *PurchaseRequestDTO) (*model.TrxScPurchaseRequest, error) {
	prNum := dto.PrNumber
	if prNum == "" {
		prNum = fmt.Sprintf("PR-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	reqDate := dto.RequestDate
	if reqDate == "" {
		reqDate = time.Now().Format("2006-01-02")
	}

	priority := dto.Priority
	if priority == "" {
		priority = "MEDIUM"
	}

	status := dto.ApprovalStatus
	if status == "" {
		status = dto.Status
	}
	if status == "" {
		status = "DRAFT"
	}

	prID := uuid.NewString()
	var items []model.TrxScPurchaseRequestItem

	for _, it := range dto.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		subtotal := float64(qty) * it.EstimatedPrice
		unit := it.Unit
		if unit == "" {
			unit = "Unit"
		}
		items = append(items, model.TrxScPurchaseRequestItem{
			ID:             uuid.NewString(),
			PrID:           prID,
			ItemID:         it.ItemID,
			Quantity:       qty,
			Unit:           unit,
			EstimatedPrice: it.EstimatedPrice,
			Subtotal:       subtotal,
			Notes:          it.Notes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}

	pr := &model.TrxScPurchaseRequest{
		ID:             prID,
		PrNumber:       prNum,
		DepartmentID:   dto.DepartmentID,
		RequesterName:  dto.RequesterName,
		RequestDate:    reqDate,
		Priority:       priority,
		Notes:          dto.Notes,
		ApprovalStatus: status,
		Status:         status,
		Items:          items,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}
	if err := s.repo.CreatePurchaseRequest(pr); err != nil {
		return nil, err
	}
	return s.repo.FindPurchaseRequestByID(pr.ID)
}

func (s *scService) UpdatePurchaseRequest(id string, dto *PurchaseRequestDTO) (*model.TrxScPurchaseRequest, error) {
	existing, err := s.repo.FindPurchaseRequestByID(id)
	if err != nil {
		return nil, errors.New("purchase request tidak ditemukan")
	}

	if dto.DepartmentID != "" {
		existing.DepartmentID = dto.DepartmentID
	}
	if dto.RequesterName != "" {
		existing.RequesterName = dto.RequesterName
	}
	if dto.RequestDate != "" {
		existing.RequestDate = dto.RequestDate
	}
	if dto.Priority != "" {
		existing.Priority = dto.Priority
	}
	existing.Notes = dto.Notes
	status := dto.ApprovalStatus
	if status == "" {
		status = dto.Status
	}
	if status != "" {
		existing.ApprovalStatus = status
		existing.Status = status
	}
	existing.UpdatedAt = time.Now()

	var items []model.TrxScPurchaseRequestItem
	for _, it := range dto.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		subtotal := float64(qty) * it.EstimatedPrice
		unit := it.Unit
		if unit == "" {
			unit = "Unit"
		}
		items = append(items, model.TrxScPurchaseRequestItem{
			ID:             uuid.NewString(),
			PrID:           id,
			ItemID:         it.ItemID,
			Quantity:       qty,
			Unit:           unit,
			EstimatedPrice: it.EstimatedPrice,
			Subtotal:       subtotal,
			Notes:          it.Notes,
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		})
	}
	existing.Items = items

	if err := s.repo.UpdatePurchaseRequest(existing); err != nil {
		return nil, err
	}
	return s.repo.FindPurchaseRequestByID(id)
}

func (s *scService) UpdatePurchaseRequestStatus(id, status string) error {
	return s.repo.UpdatePurchaseRequestStatus(id, status)
}

func (s *scService) DeletePurchaseRequest(id string) error {
	return s.repo.DeletePurchaseRequest(id)
}

// ── Transaksi Purchase Order ────────────────────────────────────────

func (s *scService) GetPurchaseOrdersPaginated(page, perPage int, vendorID, status, search string) ([]model.TrxScPurchaseOrder, int64, error) {
	return s.repo.FindPurchaseOrdersPaginated(page, perPage, vendorID, status, search)
}

func (s *scService) GetPurchaseOrderByID(id string) (*model.TrxScPurchaseOrder, error) {
	return s.repo.FindPurchaseOrderByID(id)
}

func (s *scService) CreatePurchaseOrder(dto *PurchaseOrderDTO) (*model.TrxScPurchaseOrder, error) {
	if dto.VendorID == "" {
		return nil, errors.New("vendor wajib dipilih")
	}

	poNum := dto.PoNumber
	if poNum == "" {
		poNum = fmt.Sprintf("PO-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	poDate := dto.PoDate
	if poDate == "" {
		poDate = dto.OrderDate
	}
	if poDate == "" {
		poDate = time.Now().Format("2006-01-02")
	}

	status := dto.Status
	if status == "" {
		status = dto.PoStatus
	}
	if status == "" {
		status = "DRAFT"
	}

	poID := uuid.NewString()
	var items []model.TrxScPurchaseOrderItem
	var calcSubtotal float64 = 0
	var calcTax float64 = 0

	for _, it := range dto.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		totalPrice := it.TotalPrice
		if totalPrice == 0 {
			totalPrice = float64(qty) * it.UnitPrice
		}
		calcSubtotal += totalPrice
		if it.TaxRate > 0 {
			calcTax += totalPrice * (it.TaxRate / 100.0)
		}
		unit := it.Unit
		if unit == "" {
			unit = "Unit"
		}
		items = append(items, model.TrxScPurchaseOrderItem{
			ID:         uuid.NewString(),
			PoID:       poID,
			ItemID:     it.ItemID,
			Quantity:   qty,
			Unit:       unit,
			UnitPrice:  it.UnitPrice,
			TotalPrice: totalPrice,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	subtotal := dto.Subtotal
	if subtotal == 0 {
		subtotal = calcSubtotal
	}
	tax := dto.TaxAmount
	if tax == 0 {
		if calcTax > 0 {
			tax = calcTax
		} else if dto.TaxRate > 0 {
			tax = subtotal * (dto.TaxRate / 100.0)
		}
	}
	totalAmount := dto.TotalAmount
	if totalAmount == 0 {
		totalAmount = subtotal + tax
	}

	po := &model.TrxScPurchaseOrder{
		ID:                   poID,
		PoNumber:             poNum,
		PrID:                 dto.PrID,
		VendorID:             dto.VendorID,
		PoDate:               poDate,
		ExpectedDeliveryDate: dto.ExpectedDeliveryDate,
		Subtotal:             subtotal,
		TaxAmount:            tax,
		TotalAmount:          totalAmount,
		Status:               status,
		PoStatus:             status,
		Notes:                dto.Notes,
		Items:                items,
		CreatedAt:            time.Now(),
		UpdatedAt:            time.Now(),
	}
	if err := s.repo.CreatePurchaseOrder(po); err != nil {
		return nil, err
	}
	return s.repo.FindPurchaseOrderByID(po.ID)
}

func (s *scService) UpdatePurchaseOrder(id string, dto *PurchaseOrderDTO) (*model.TrxScPurchaseOrder, error) {
	existing, err := s.repo.FindPurchaseOrderByID(id)
	if err != nil {
		return nil, errors.New("purchase order tidak ditemukan")
	}

	if dto.VendorID != "" {
		existing.VendorID = dto.VendorID
	}
	if dto.PrID != nil {
		existing.PrID = dto.PrID
	}
	poDate := dto.PoDate
	if poDate == "" {
		poDate = dto.OrderDate
	}
	if poDate != "" {
		existing.PoDate = poDate
	}
	if dto.ExpectedDeliveryDate != "" {
		existing.ExpectedDeliveryDate = dto.ExpectedDeliveryDate
	}
	status := dto.Status
	if status == "" {
		status = dto.PoStatus
	}
	if status != "" {
		existing.Status = status
		existing.PoStatus = status
	}
	existing.Notes = dto.Notes
	existing.UpdatedAt = time.Now()

	var items []model.TrxScPurchaseOrderItem
	var calcSubtotal float64 = 0

	for _, it := range dto.Items {
		qty := it.Quantity
		if qty <= 0 {
			qty = 1
		}
		totalPrice := it.TotalPrice
		if totalPrice == 0 {
			totalPrice = float64(qty) * it.UnitPrice
		}
		calcSubtotal += totalPrice
		unit := it.Unit
		if unit == "" {
			unit = "Unit"
		}
		items = append(items, model.TrxScPurchaseOrderItem{
			ID:         uuid.NewString(),
			PoID:       id,
			ItemID:     it.ItemID,
			Quantity:   qty,
			Unit:       unit,
			UnitPrice:  it.UnitPrice,
			TotalPrice: totalPrice,
			CreatedAt:  time.Now(),
			UpdatedAt:  time.Now(),
		})
	}

	existing.Items = items
	existing.Subtotal = dto.Subtotal
	if existing.Subtotal == 0 {
		existing.Subtotal = calcSubtotal
	}
	existing.TaxAmount = dto.TaxAmount
	existing.TotalAmount = dto.TotalAmount
	if existing.TotalAmount == 0 {
		existing.TotalAmount = existing.Subtotal + existing.TaxAmount
	}

	if err := s.repo.UpdatePurchaseOrder(existing); err != nil {
		return nil, err
	}
	return s.repo.FindPurchaseOrderByID(id)
}

func (s *scService) UpdatePurchaseOrderStatus(id, status string) error {
	return s.repo.UpdatePurchaseOrderStatus(id, status)
}

func (s *scService) DeletePurchaseOrder(id string) error {
	return s.repo.DeletePurchaseOrder(id)
}

// ── Transaksi Goods Receipt ─────────────────────────────────────────

func (s *scService) GetGoodsReceiptsPaginated(page, perPage int, poID, search string) ([]model.TrxScGoodsReceipt, int64, error) {
	return s.repo.FindGoodsReceiptsPaginated(page, perPage, poID, search)
}

func (s *scService) GetGoodsReceiptByID(id string) (*model.TrxScGoodsReceipt, error) {
	return s.repo.FindGoodsReceiptByID(id)
}

func (s *scService) CreateGoodsReceipt(dto *GoodsReceiptDTO) (*model.TrxScGoodsReceipt, error) {
	poID := dto.PoID
	if poID == "" {
		poID = dto.PurchaseOrderID
	}
	if poID == "" {
		return nil, errors.New("referensi PO wajib dipilih")
	}

	grNum := dto.ReceiptNumber
	if grNum == "" {
		grNum = dto.GrNumber
	}
	if grNum == "" {
		grNum = fmt.Sprintf("GR-%s-%s", time.Now().Format("200601"), strings.ToUpper(uuid.NewString()[:5]))
	}

	receiptDate := dto.ReceiptDate
	if receiptDate == "" {
		receiptDate = time.Now().Format("2006-01-02")
	}

	delNote := dto.DeliveryNoteNumber
	if delNote == "" {
		delNote = dto.DeliveryNoteNo
	}

	status := dto.Status
	if status == "" {
		status = "COMPLETE"
	}

	grID := uuid.NewString()
	var items []model.TrxScGoodsReceiptItem

	for _, it := range dto.Items {
		cond := it.Condition
		if cond == "" {
			cond = "GOOD"
		}
		items = append(items, model.TrxScGoodsReceiptItem{
			ID:               uuid.NewString(),
			GrID:             grID,
			ItemID:           it.ItemID,
			OrderedQuantity:  it.OrderedQuantity,
			ReceivedQuantity: it.ReceivedQuantity,
			Condition:        cond,
			Notes:            it.Notes,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		})
	}

	gr := &model.TrxScGoodsReceipt{
		ID:                 grID,
		ReceiptNumber:      grNum,
		PoID:               poID,
		ReceivedBy:         dto.ReceivedBy,
		ReceiptDate:        receiptDate,
		DeliveryNoteNumber: delNote,
		AttachmentURL:      dto.AttachmentURL,
		Notes:              dto.Notes,
		Status:             status,
		Items:              items,
		CreatedAt:          time.Now(),
		UpdatedAt:          time.Now(),
	}
	if err := s.repo.CreateGoodsReceipt(gr); err != nil {
		return nil, err
	}
	return s.repo.FindGoodsReceiptByID(gr.ID)
}

func (s *scService) UpdateGoodsReceipt(id string, dto *GoodsReceiptDTO) (*model.TrxScGoodsReceipt, error) {
	existing, err := s.repo.FindGoodsReceiptByID(id)
	if err != nil {
		return nil, errors.New("goods receipt tidak ditemukan")
	}

	if dto.PoID != "" {
		existing.PoID = dto.PoID
	}
	if dto.ReceivedBy != "" {
		existing.ReceivedBy = dto.ReceivedBy
	}
	if dto.ReceiptDate != "" {
		existing.ReceiptDate = dto.ReceiptDate
	}
	if dto.DeliveryNoteNumber != "" {
		existing.DeliveryNoteNumber = dto.DeliveryNoteNumber
	}
	if dto.AttachmentURL != "" {
		existing.AttachmentURL = dto.AttachmentURL
	}
	existing.Notes = dto.Notes
	if dto.Status != "" {
		existing.Status = dto.Status
	}
	existing.UpdatedAt = time.Now()

	var items []model.TrxScGoodsReceiptItem
	for _, it := range dto.Items {
		cond := it.Condition
		if cond == "" {
			cond = "GOOD"
		}
		items = append(items, model.TrxScGoodsReceiptItem{
			ID:               uuid.NewString(),
			GrID:             id,
			ItemID:           it.ItemID,
			OrderedQuantity:  it.OrderedQuantity,
			ReceivedQuantity: it.ReceivedQuantity,
			Condition:        cond,
			Notes:            it.Notes,
			CreatedAt:        time.Now(),
			UpdatedAt:        time.Now(),
		})
	}
	existing.Items = items

	if err := s.repo.UpdateGoodsReceipt(existing); err != nil {
		return nil, err
	}
	return s.repo.FindGoodsReceiptByID(id)
}

func (s *scService) DeleteGoodsReceipt(id string) error {
	return s.repo.DeleteGoodsReceipt(id)
}

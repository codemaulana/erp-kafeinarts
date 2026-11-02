package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"erp-kafeinart/model"
	"erp-kafeinart/utils"

	"github.com/google/uuid"
)

// ── Jurnal Kas (Pencatatan Harian) Implementasi ─────────────────────────────

func (s *financeService) GetCashJournals(page, perPage int, accountID, txType, categoryID, search string) ([]model.TrxCashJournal, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.trxRepo.FindCashJournals(page, perPage, accountID, txType, categoryID, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *financeService) GetCashJournalByID(id string) (*model.TrxCashJournal, error) {
	j, err := s.trxRepo.FindCashJournalByID(id)
	if err != nil || j == nil {
		return nil, errors.New("transaksi jurnal kas tidak ditemukan")
	}
	return j, nil
}

func (s *financeService) CreateCashJournal(dto CashJournalDTO) (*model.TrxCashJournal, error) {
	amount := dto.Amount
	if amount <= 0 && dto.Nominal > 0 {
		amount = dto.Nominal
	}
	if amount <= 0 {
		return nil, errors.New("nominal transaksi harus lebih besar dari 0")
	}

	accID := strings.TrimSpace(dto.AccountID)
	if accID == "" {
		accID = strings.TrimSpace(dto.AkunBank)
	}
	if accID == "" {
		return nil, errors.New("akun bank / kas wajib dipilih")
	}

	catID := strings.TrimSpace(dto.CategoryID)
	if catID == "" {
		catID = strings.TrimSpace(dto.Kategori)
	}
	if catID == "" {
		return nil, errors.New("kategori transaksi wajib dipilih")
	}

	txType := strings.ToUpper(strings.TrimSpace(dto.TransactionType))
	if txType == "" {
		if strings.EqualFold(dto.Tipe, "masuk") {
			txType = "INCOME"
		} else if strings.EqualFold(dto.Tipe, "keluar") {
			txType = "EXPENSE"
		} else {
			txType = strings.ToUpper(strings.TrimSpace(dto.Tipe))
		}
	}
	if txType != "INCOME" && txType != "EXPENSE" {
		txType = "EXPENSE"
	}

	txNum := strings.TrimSpace(dto.TransactionNumber)
	if txNum == "" {
		prefix := "TRX-OUT"
		if txType == "INCOME" {
			prefix = "TRX-IN"
		}
		txNum = fmt.Sprintf("%s-%s-%s", prefix, time.Now().Format("20060102"), uuid.NewString()[:6])
	}

	dateStr := dto.TransactionDate
	if dateStr == "" {
		dateStr = dto.Tanggal
	}
	txDate := parseFlexibleDate(dateStr)
	if txDate.IsZero() {
		txDate = time.Now()
	}

	desc := dto.Description
	if desc == "" {
		desc = dto.Keterangan
	}

	attURL := dto.AttachmentURL
	if attURL == "" {
		attURL = dto.BuktiFile
	}

	journal := &model.TrxCashJournal{
		ID:                uuid.NewString(),
		TransactionNumber: txNum,
		TransactionDate:   txDate,
		AccountID:         accID,
		CategoryID:        catID,
		TransactionType:   txType,
		Amount:            amount,
		Description:       desc,
		AttachmentURL:     attURL,
		CreatedAt:         time.Now(),
		UpdatedAt:         time.Now(),
	}

	if err := s.trxRepo.CreateCashJournal(journal); err != nil {
		return nil, err
	}

	// Reload with preloaded data
	return s.trxRepo.FindCashJournalByID(journal.ID)
}

func (s *financeService) UpdateCashJournal(id string, dto CashJournalDTO) (*model.TrxCashJournal, error) {
	existing, err := s.trxRepo.FindCashJournalByID(id)
	if err != nil || existing == nil {
		return nil, errors.New("jurnal kas tidak ditemukan")
	}

	oldJournal := *existing

	txDate := parseFlexibleDate(dto.TransactionDate)
	if txDate.IsZero() {
		txDate = parseFlexibleDate(dto.Tanggal)
	}
	if !txDate.IsZero() {
		existing.TransactionDate = txDate
	}

	accID := strings.TrimSpace(dto.AccountID)
	if accID == "" {
		accID = strings.TrimSpace(dto.AkunBank)
	}
	if accID != "" {
		existing.AccountID = accID
	}

	catID := strings.TrimSpace(dto.CategoryID)
	if catID == "" {
		catID = strings.TrimSpace(dto.Kategori)
	}
	if catID != "" {
		existing.CategoryID = catID
	}

	txType := strings.ToUpper(strings.TrimSpace(dto.TransactionType))
	if txType == "" {
		if strings.EqualFold(dto.Tipe, "masuk") {
			txType = "INCOME"
		} else if strings.EqualFold(dto.Tipe, "keluar") {
			txType = "EXPENSE"
		} else {
			txType = strings.ToUpper(strings.TrimSpace(dto.Tipe))
		}
	}
	if txType == "INCOME" || txType == "EXPENSE" {
		existing.TransactionType = txType
	}

	amount := dto.Amount
	if amount == 0 {
		amount = dto.Nominal
	}
	if amount > 0 {
		existing.Amount = amount
	}

	desc := strings.TrimSpace(dto.Description)
	if desc == "" {
		desc = strings.TrimSpace(dto.Keterangan)
	}
	if desc != "" {
		existing.Description = desc
	}

	attURL := strings.TrimSpace(dto.AttachmentURL)
	if attURL == "" {
		attURL = strings.TrimSpace(dto.BuktiFile)
	}
	if attURL != "" {
		existing.AttachmentURL = attURL
	}

	existing.UpdatedAt = time.Now()

	if err := s.trxRepo.UpdateCashJournal(&oldJournal, existing); err != nil {
		return nil, err
	}

	return s.trxRepo.FindCashJournalByID(existing.ID)
}

func (s *financeService) DeleteCashJournal(id string) error {
	return s.trxRepo.DeleteCashJournal(id)
}

// ── Manajemen Invoice / Piutang (AR) Implementasi ───────────────────────────

func (s *financeService) GetInvoices(page, perPage int, status, clientID, search string) ([]model.TrxInvoice, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.trxRepo.FindInvoices(page, perPage, status, clientID, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *financeService) GetInvoiceByID(id string) (*model.TrxInvoice, error) {
	inv, err := s.trxRepo.FindInvoiceByID(id)
	if err != nil || inv == nil {
		return nil, errors.New("invoice tidak ditemukan")
	}
	return inv, nil
}

func (s *financeService) CreateInvoice(dto InvoiceDTO) (*model.TrxInvoice, error) {
	clientName := strings.TrimSpace(dto.ClientName)
	if clientName == "" {
		clientName = strings.TrimSpace(dto.Klien)
	}
	if clientName == "" {
		return nil, errors.New("nama klien wajib diisi")
	}

	projectName := strings.TrimSpace(dto.ProjectName)
	if projectName == "" {
		projectName = strings.TrimSpace(dto.Proyek)
	}

	invNum := strings.TrimSpace(dto.InvoiceNumber)
	if invNum == "" {
		invNum = strings.TrimSpace(dto.NoInvoice)
	}
	if invNum == "" {
		invNum = fmt.Sprintf("INV-%s-%s", time.Now().Format("200601"), uuid.NewString()[:6])
	}

	issueDateStr := dto.IssueDate
	if issueDateStr == "" {
		issueDateStr = dto.TglTerbit
	}
	issueDate := parseFlexibleDate(issueDateStr)
	if issueDate.IsZero() {
		issueDate = time.Now()
	}

	dueDateStr := dto.DueDate
	if dueDateStr == "" {
		dueDateStr = dto.TglJatuhTempo
	}
	dueDate := parseFlexibleDate(dueDateStr)
	if dueDate.IsZero() {
		dueDate = issueDate.AddDate(0, 0, 14) // default 14 hari
	}

	status := strings.ToUpper(strings.TrimSpace(dto.Status))
	switch status {
	case "LUNAS", "PAID":
		status = "PAID"
	case "SEBAGIAN", "PARTIAL":
		status = "PARTIAL"
	case "BATAL", "CANCELLED":
		status = "CANCELLED"
	case "BELUM LUNAS", "SENT":
		status = "SENT"
	case "DRAFT", "":
		status = "DRAFT"
	}

	invoiceID := uuid.NewString()
	var items []model.TrxInvoiceItem
	var subtotal, totalTax float64

	for _, it := range dto.Items {
		desc := strings.TrimSpace(it.ItemDescription)
		if desc == "" {
			desc = strings.TrimSpace(it.Deskripsi)
		}
		if desc == "" {
			continue
		}
		qty := it.Quantity
		if qty <= 0 {
			qty = it.Qty
		}
		if qty <= 0 {
			qty = 1
		}
		unitPrice := it.UnitPrice
		if unitPrice == 0 && it.Harga > 0 {
			unitPrice = it.Harga
		}
		totalPrice := float64(qty) * unitPrice
		subtotal += totalPrice

		taxID := it.TaxID
		if taxID == nil && it.Pajak != nil && *it.Pajak != "" {
			taxID = it.Pajak
		} else if taxID == nil && dto.Pajak != "" {
			taxID = &dto.Pajak
		}

		if taxID != nil && *taxID != "" {
			tax, err := s.masterRepo.FindTaxByID(*taxID)
			if err == nil && tax != nil {
				taxVal := totalPrice * (tax.Percentage / 100.0)
				totalTax += taxVal
			}
		}

		items = append(items, model.TrxInvoiceItem{
			ID:              uuid.NewString(),
			InvoiceID:       invoiceID,
			ItemDescription: desc,
			Quantity:        qty,
			UnitPrice:       unitPrice,
			TaxID:           taxID,
			TotalPrice:      totalPrice,
			CreatedAt:       time.Now(),
		})
	}

	// Jika tidak ada item dinamis yang dikirim, buatkan 1 baris default
	if len(items) == 0 {
		items = append(items, model.TrxInvoiceItem{
			ID:              uuid.NewString(),
			InvoiceID:       invoiceID,
			ItemDescription: "Biaya Jasa & Pengembangan",
			Quantity:        1,
			UnitPrice:       0,
			TotalPrice:      0,
			CreatedAt:       time.Now(),
		})
	}

	grandTotal := subtotal + totalTax

	inv := &model.TrxInvoice{
		ID:            invoiceID,
		InvoiceNumber: invNum,
		ClientID:      dto.ClientID,
		ClientName:    clientName,
		ProjectID:     dto.ProjectID,
		ProjectName:   projectName,
		IssueDate:     issueDate,
		DueDate:       dueDate,
		Subtotal:      subtotal,
		TotalTax:      totalTax,
		GrandTotal:    grandTotal,
		Status:        status,
		Notes:         dto.Notes,
		Items:         items,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if status == "PAID" {
		now := time.Now()
		inv.PaidAt = &now
	}

	if err := s.trxRepo.CreateInvoice(inv); err != nil {
		return nil, err
	}

	return s.trxRepo.FindInvoiceByID(invoiceID)
}

func (s *financeService) UpdateInvoice(id string, dto InvoiceDTO) (*model.TrxInvoice, error) {
	inv, err := s.GetInvoiceByID(id)
	if err != nil {
		return nil, err
	}

	clientName := strings.TrimSpace(dto.ClientName)
	if clientName == "" {
		clientName = strings.TrimSpace(dto.Klien)
	}
	if clientName != "" {
		inv.ClientName = clientName
	}

	projectName := strings.TrimSpace(dto.ProjectName)
	if projectName == "" {
		projectName = strings.TrimSpace(dto.Proyek)
	}
	if projectName != "" {
		inv.ProjectName = projectName
	}

	issueDateStr := dto.IssueDate
	if issueDateStr == "" {
		issueDateStr = dto.TglTerbit
	}
	if issueDateStr != "" {
		if t := parseFlexibleDate(issueDateStr); !t.IsZero() {
			inv.IssueDate = t
		}
	}

	dueDateStr := dto.DueDate
	if dueDateStr == "" {
		dueDateStr = dto.TglJatuhTempo
	}
	if dueDateStr != "" {
		if t := parseFlexibleDate(dueDateStr); !t.IsZero() {
			inv.DueDate = t
		}
	}

	status := strings.ToUpper(strings.TrimSpace(dto.Status))
	switch status {
	case "LUNAS", "PAID":
		inv.Status = "PAID"
	case "SEBAGIAN", "PARTIAL":
		inv.Status = "PARTIAL"
	case "BATAL", "CANCELLED":
		inv.Status = "CANCELLED"
	case "BELUM LUNAS", "SENT":
		inv.Status = "SENT"
	case "DRAFT":
		inv.Status = "DRAFT"
	}

	if dto.Notes != "" {
		inv.Notes = dto.Notes
	}

	if len(dto.Items) > 0 {
		var items []model.TrxInvoiceItem
		var subtotal, totalTax float64

		for _, it := range dto.Items {
			desc := strings.TrimSpace(it.ItemDescription)
			if desc == "" {
				desc = strings.TrimSpace(it.Deskripsi)
			}
			if desc == "" {
				continue
			}
			qty := it.Quantity
			if qty <= 0 {
				qty = it.Qty
			}
			if qty <= 0 {
				qty = 1
			}
			unitPrice := it.UnitPrice
			if unitPrice == 0 && it.Harga > 0 {
				unitPrice = it.Harga
			}
			totalPrice := float64(qty) * unitPrice
			subtotal += totalPrice

			taxID := it.TaxID
			if taxID == nil && it.Pajak != nil && *it.Pajak != "" {
				taxID = it.Pajak
			} else if taxID == nil && dto.Pajak != "" {
				taxID = &dto.Pajak
			}

			if taxID != nil && *taxID != "" {
				tax, errTax := s.masterRepo.FindTaxByID(*taxID)
				if errTax == nil && tax != nil {
					taxVal := totalPrice * (tax.Percentage / 100.0)
					totalTax += taxVal
				}
			}

			items = append(items, model.TrxInvoiceItem{
				ID:              uuid.NewString(),
				InvoiceID:       inv.ID,
				ItemDescription: desc,
				Quantity:        qty,
				UnitPrice:       unitPrice,
				TaxID:           taxID,
				TotalPrice:      totalPrice,
				CreatedAt:       time.Now(),
			})
		}
		inv.Items = items
		inv.Subtotal = subtotal
		inv.TotalTax = totalTax
		inv.GrandTotal = subtotal + totalTax
	}

	inv.UpdatedAt = time.Now()
	if err := s.trxRepo.UpdateInvoice(inv); err != nil {
		return nil, err
	}

	return s.trxRepo.FindInvoiceByID(inv.ID)
}

func (s *financeService) UpdateInvoiceStatus(id, status string) (*model.TrxInvoice, error) {
	st := strings.ToUpper(strings.TrimSpace(status))
	switch st {
	case "LUNAS", "PAID":
		st = "PAID"
	case "SEBAGIAN", "PARTIAL":
		st = "PARTIAL"
	case "BATAL", "CANCELLED":
		st = "CANCELLED"
	case "BELUM LUNAS", "SENT":
		st = "SENT"
	}
	return s.trxRepo.UpdateInvoiceStatus(id, st)
}

func (s *financeService) DeleteInvoice(id string) error {
	return s.trxRepo.DeleteInvoice(id)
}

// ── Tagihan Vendor / Hutang (AP) Implementasi ───────────────────────────────

func (s *financeService) GetVendorBills(page, perPage int, status, vendorID, search string) ([]model.TrxVendorBill, utils.PaginationMeta, error) {
	if page <= 0 {
		page = 1
	}
	if perPage <= 0 {
		perPage = 15
	}

	list, total, err := s.trxRepo.FindVendorBills(page, perPage, status, vendorID, search)
	if err != nil {
		return nil, utils.PaginationMeta{}, err
	}

	meta := utils.BuildPaginationMeta(total, page, perPage)
	return list, meta, nil
}

func (s *financeService) GetVendorBillByID(id string) (*model.TrxVendorBill, error) {
	bill, err := s.trxRepo.FindVendorBillByID(id)
	if err != nil || bill == nil {
		return nil, errors.New("tagihan vendor tidak ditemukan")
	}
	return bill, nil
}

func (s *financeService) CreateVendorBill(dto VendorBillDTO) (*model.TrxVendorBill, error) {
	vendorName := strings.TrimSpace(dto.VendorName)
	if vendorName == "" {
		vendorName = strings.TrimSpace(dto.Vendor)
	}
	if vendorName == "" {
		return nil, errors.New("nama vendor wajib diisi (misal: AWS, Google Cloud, lisensi software)")
	}

	categoryID := strings.TrimSpace(dto.CategoryID)
	if categoryID == "" {
		categoryID = strings.TrimSpace(dto.Kategori)
	}
	if categoryID == "" {
		return nil, errors.New("kategori biaya wajib dipilih")
	}

	amount := dto.Amount
	if amount <= 0 && dto.Nominal > 0 {
		amount = dto.Nominal
	}
	if amount <= 0 {
		return nil, errors.New("nominal tagihan harus lebih besar dari 0")
	}

	billNum := strings.TrimSpace(dto.BillNumber)
	if billNum == "" {
		billNum = strings.TrimSpace(dto.NoTagihan)
	}
	if billNum == "" {
		billNum = fmt.Sprintf("BILL-%s-%s", time.Now().Format("200601"), uuid.NewString()[:6])
	}

	billDateStr := dto.BillDate
	if billDateStr == "" {
		billDateStr = dto.TanggalTagihan
	}
	billDate := parseFlexibleDate(billDateStr)
	if billDate.IsZero() {
		billDate = time.Now()
	}

	dueDateStr := dto.DueDate
	if dueDateStr == "" {
		dueDateStr = dto.JatuhTempo
	}
	dueDate := parseFlexibleDate(dueDateStr)
	if dueDate.IsZero() {
		dueDate = billDate.AddDate(0, 0, 14)
	}

	status := strings.ToUpper(strings.TrimSpace(dto.Status))
	switch status {
	case "LUNAS", "PAID":
		status = "PAID"
	case "JATUH TEMPO", "OVERDUE":
		status = "OVERDUE"
	case "BELUM LUNAS", "PENDING", "":
		status = "PENDING"
	}

	attURL := dto.AttachmentURL
	if attURL == "" {
		attURL = dto.BuktiFile
	}

	bill := &model.TrxVendorBill{
		ID:            uuid.NewString(),
		BillNumber:    billNum,
		VendorID:      dto.VendorID,
		VendorName:    vendorName,
		CategoryID:    categoryID,
		BillDate:      billDate,
		DueDate:       dueDate,
		Amount:        amount,
		Status:        status,
		AttachmentURL: attURL,
		CreatedAt:     time.Now(),
		UpdatedAt:     time.Now(),
	}

	if status == "PAID" {
		now := time.Now()
		bill.PaidAt = &now
	}

	if err := s.trxRepo.CreateVendorBill(bill); err != nil {
		return nil, err
	}

	return s.trxRepo.FindVendorBillByID(bill.ID)
}

func (s *financeService) UpdateVendorBill(id string, dto VendorBillDTO) (*model.TrxVendorBill, error) {
	bill, err := s.GetVendorBillByID(id)
	if err != nil {
		return nil, err
	}

	vendorName := strings.TrimSpace(dto.VendorName)
	if vendorName == "" {
		vendorName = strings.TrimSpace(dto.Vendor)
	}
	if vendorName != "" {
		bill.VendorName = vendorName
	}

	categoryID := strings.TrimSpace(dto.CategoryID)
	if categoryID == "" {
		categoryID = strings.TrimSpace(dto.Kategori)
	}
	if categoryID != "" {
		bill.CategoryID = categoryID
	}

	if dto.Amount > 0 {
		bill.Amount = dto.Amount
	} else if dto.Nominal > 0 {
		bill.Amount = dto.Nominal
	}

	billDateStr := dto.BillDate
	if billDateStr == "" {
		billDateStr = dto.TanggalTagihan
	}
	if billDateStr != "" {
		if t := parseFlexibleDate(billDateStr); !t.IsZero() {
			bill.BillDate = t
		}
	}

	dueDateStr := dto.DueDate
	if dueDateStr == "" {
		dueDateStr = dto.JatuhTempo
	}
	if dueDateStr != "" {
		if t := parseFlexibleDate(dueDateStr); !t.IsZero() {
			bill.DueDate = t
		}
	}

	status := strings.ToUpper(strings.TrimSpace(dto.Status))
	switch status {
	case "LUNAS", "PAID":
		bill.Status = "PAID"
	case "JATUH TEMPO", "OVERDUE":
		bill.Status = "OVERDUE"
	case "BELUM LUNAS", "PENDING":
		bill.Status = "PENDING"
	}

	if dto.AttachmentURL != "" {
		bill.AttachmentURL = dto.AttachmentURL
	} else if dto.BuktiFile != "" {
		bill.AttachmentURL = dto.BuktiFile
	}
	bill.UpdatedAt = time.Now()

	if err := s.trxRepo.UpdateVendorBill(bill); err != nil {
		return nil, err
	}
	return s.trxRepo.FindVendorBillByID(bill.ID)
}

func (s *financeService) UpdateVendorBillStatus(id, status string) (*model.TrxVendorBill, error) {
	st := strings.ToUpper(strings.TrimSpace(status))
	switch st {
	case "LUNAS", "PAID":
		st = "PAID"
	case "JATUH TEMPO", "OVERDUE":
		st = "OVERDUE"
	case "BELUM LUNAS", "PENDING":
		st = "PENDING"
	}
	return s.trxRepo.UpdateVendorBillStatus(id, st)
}

func (s *financeService) DeleteVendorBill(id string) error {
	return s.trxRepo.DeleteVendorBill(id)
}

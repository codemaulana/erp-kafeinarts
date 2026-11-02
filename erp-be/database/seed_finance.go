package database

import (
	"log"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedFinanceMaster menginisialisasi master kategori, rekening kas/bank, dan tarif pajak.
func SeedFinanceMaster(db *gorm.DB) {
	var catCount int64
	db.Model(&model.MstFinanceCategory{}).Count(&catCount)
	if catCount == 0 {
		categories := []model.MstFinanceCategory{
			{ID: "cat-rev-01", CategoryName: "Project Revenue", CashType: "INCOME", Description: "Pendapatan dari termin proyek klien software house"},
			{ID: "cat-srv-01", CategoryName: "Biaya Server", CashType: "EXPENSE", Description: "Biaya cloud server hosting AWS, Google Cloud, dsb."},
			{ID: "cat-sal-01", CategoryName: "Gaji", CashType: "EXPENSE", Description: "Gaji tim developer dan desainer Kafeinarts"},
			{ID: "cat-lic-01", CategoryName: "Lisensi Software", CashType: "EXPENSE", Description: "Langganan lisensi IDE JetBrains, Figma, GitHub, Adobe"},
			{ID: "cat-sup-01", CategoryName: "Perlengkapan Studio", CashType: "EXPENSE", Description: "ATK, wifi fiber optik, dan listrik studio"},
		}
		_ = db.Create(&categories).Error
		log.Printf("[SEEDER] %d master kategori finance berhasil dibuat", len(categories))
	}

	var accCount int64
	db.Model(&model.MstFinanceAccount{}).Count(&accCount)
	if accCount == 0 {
		now := time.Now()
		accounts := []model.MstFinanceAccount{
			{ID: "acc-bca-01", BankName: "BCA", AccountNumber: "8801234567", AccountHolderName: "PT Kafeinarts Studio", CurrentBalance: 50000000, CreatedAt: now, UpdatedAt: now},
			{ID: "acc-mandiri-01", BankName: "Mandiri", AccountNumber: "1370098765432", AccountHolderName: "PT Kafeinarts Studio", CurrentBalance: 25000000, CreatedAt: now, UpdatedAt: now},
			{ID: "acc-petty-01", BankName: "Petty Cash", AccountNumber: "0000", AccountHolderName: "Kasir Operasional", CurrentBalance: 2500000, CreatedAt: now, UpdatedAt: now},
		}
		_ = db.Create(&accounts).Error
		log.Printf("[SEEDER] %d master rekening kas/bank berhasil dibuat", len(accounts))
	}

	var taxCount int64
	db.Model(&model.MstFinanceTax{}).Count(&taxCount)
	if taxCount == 0 {
		now := time.Now()
		taxes := []model.MstFinanceTax{
			{ID: "tax-ppn-11", TaxName: "PPN 11%", Percentage: 11.00, Description: "Pajak Pertambahan Nilai 11%", CreatedAt: now, UpdatedAt: now},
			{ID: "tax-pph-21", TaxName: "PPh 21", Percentage: 5.00, Description: "Pajak Penghasilan Pasal 21 atas jasa", CreatedAt: now, UpdatedAt: now},
			{ID: "tax-pph-23", TaxName: "PPh 23", Percentage: 2.00, Description: "Pajak Penghasilan Pasal 23", CreatedAt: now, UpdatedAt: now},
			{ID: "tax-zero-00", TaxName: "Bebas Pajak (0%)", Percentage: 0.00, Description: "Tanpa pungutan pajak", CreatedAt: now, UpdatedAt: now},
		}
		_ = db.Create(&taxes).Error
		log.Printf("[SEEDER] %d master tarif pajak berhasil dibuat", len(taxes))
	}
}

// SeedFinanceTransactions menginisialisasi transaksi jurnal kas, invoice AR, dan tagihan AP.
func SeedFinanceTransactions(db *gorm.DB) {
	now := time.Now()

	var journalCount int64
	db.Model(&model.TrxCashJournal{}).Count(&journalCount)
	if journalCount == 0 {
		journals := []model.TrxCashJournal{
			{
				ID:                uuid.NewString(),
				TransactionNumber: "TRX-IN-202609-001",
				TransactionDate:   now.AddDate(0, 0, -3),
				AccountID:         "acc-bca-01",
				CategoryID:        "cat-rev-01",
				TransactionType:   "INCOME",
				Amount:            16650000,
				Description:       "Penerimaan pelunasan termin 1 web e-commerce PT Sinar Digital",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                uuid.NewString(),
				TransactionNumber: "TRX-OUT-202609-001",
				TransactionDate:   now.AddDate(0, 0, -2),
				AccountID:         "acc-mandiri-01",
				CategoryID:        "cat-srv-01",
				TransactionType:   "EXPENSE",
				Amount:            2719500,
				Description:       "Pembayaran server cloud AWS produksi Agustus",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
			{
				ID:                uuid.NewString(),
				TransactionNumber: "TRX-OUT-202609-002",
				TransactionDate:   now.AddDate(0, 0, -1),
				AccountID:         "acc-petty-01",
				CategoryID:        "cat-sup-01",
				TransactionType:   "EXPENSE",
				Amount:            450000,
				Description:       "Pembelian kertas foto, tinta printer, dan supplies studio",
				CreatedAt:         now,
				UpdatedAt:         now,
			},
		}
		_ = db.Create(&journals).Error
		log.Printf("[SEEDER] %d jurnal kas harian berhasil ditambahkan", len(journals))
	}

	var invCount int64
	db.Model(&model.TrxInvoice{}).Count(&invCount)
	if invCount == 0 {
		taxPPN := "tax-ppn-11"
		inv1ID := uuid.NewString()
		inv2ID := uuid.NewString()

		inv1 := model.TrxInvoice{
			ID:            inv1ID,
			InvoiceNumber: "INV-202609-001",
			ClientID:      "crm-client-01",
			ClientName:    "PT Sinar Digital Nusantara",
			ProjectID:     "proj-01",
			ProjectName:   "Website E-Commerce & Brand Kit",
			IssueDate:     now.AddDate(0, 0, -10),
			DueDate:       now.AddDate(0, 0, 4),
			Subtotal:      15000000,
			TotalTax:      1650000,
			GrandTotal:    16650000,
			Status:        "PAID",
			Notes:         "Pembayaran via transfer BCA",
			PaidAt:        &now,
			CreatedAt:     now,
			UpdatedAt:     now,
			Items: []model.TrxInvoiceItem{
				{
					ID:              uuid.NewString(),
					InvoiceID:       inv1ID,
					ItemDescription: "Desain UI/UX Design System (Figma) - 15 Pages",
					Quantity:        1,
					UnitPrice:       6000000,
					TaxID:           &taxPPN,
					TotalPrice:      6000000,
					CreatedAt:       now,
				},
				{
					ID:              uuid.NewString(),
					InvoiceID:       inv1ID,
					ItemDescription: "Frontend & Backend Development (Termin 1)",
					Quantity:        1,
					UnitPrice:       9000000,
					TaxID:           &taxPPN,
					TotalPrice:      9000000,
					CreatedAt:       now,
				},
			},
		}

		inv2 := model.TrxInvoice{
			ID:            inv2ID,
			InvoiceNumber: "INV-202609-002",
			ClientID:      "crm-client-02",
			ClientName:    "CV Kafein Kreatif Media",
			ProjectID:     "proj-02",
			ProjectName:   "Mobile App POS & Loyalty Member",
			IssueDate:     now.AddDate(0, 0, -2),
			DueDate:       now.AddDate(0, 0, 12),
			Subtotal:      8500000,
			TotalTax:      935000,
			GrandTotal:    9435000,
			Status:        "SENT",
			Notes:         "Jatuh tempo 14 hari kerja",
			CreatedAt:     now,
			UpdatedAt:     now,
			Items: []model.TrxInvoiceItem{
				{
					ID:              uuid.NewString(),
					InvoiceID:       inv2ID,
					ItemDescription: "Integrasi Payment Gateway QRIS & E-Wallet",
					Quantity:        1,
					UnitPrice:       8500000,
					TaxID:           &taxPPN,
					TotalPrice:      8500000,
					CreatedAt:       now,
				},
			},
		}

		_ = db.Create(&inv1).Error
		_ = db.Create(&inv2).Error
		log.Printf("[SEEDER] 2 invoice piutang (AR) dengan detail baris dinamis berhasil ditambahkan")
	}

	var billCount int64
	db.Model(&model.TrxVendorBill{}).Count(&billCount)
	if billCount == 0 {
		bills := []model.TrxVendorBill{
			{
				ID:            uuid.NewString(),
				BillNumber:    "BILL-202609-001",
				VendorID:      "sc-vendor-01",
				VendorName:    "Amazon Web Services (AWS)",
				CategoryID:    "cat-srv-01",
				BillDate:      now.AddDate(0, 0, -5),
				DueDate:       now.AddDate(0, 0, 5),
				Amount:        2719500,
				Status:        "PAID",
				AttachmentURL: "/uploads/sample_aws_invoice.pdf",
				PaidAt:        &now,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			{
				ID:            uuid.NewString(),
				BillNumber:    "BILL-202609-002",
				VendorID:      "sc-vendor-02",
				VendorName:    "JetBrains & Figma Inc",
				CategoryID:    "cat-lic-01",
				BillDate:      now.AddDate(0, 0, -2),
				DueDate:       now.AddDate(0, 0, 10),
				Amount:        1332000,
				Status:        "PENDING",
				AttachmentURL: "/uploads/sample_jetbrains_receipt.pdf",
				CreatedAt:     now,
				UpdatedAt:     now,
			},
		}
		_ = db.Create(&bills).Error
		log.Printf("[SEEDER] %d tagihan vendor (AP) berhasil ditambahkan", len(bills))
	}
}

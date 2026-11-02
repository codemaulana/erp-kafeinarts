package database

import (
	"log"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedScAll menginisialisasi seluruh master dan transaksi modul Supply Chain (SC).
func SeedScAll(db *gorm.DB) {
	log.Println("[SEEDER] Memulai inisialisasi dummy data Supply Chain (SC)...")

	vnd1ID := "v1111111-1111-1111-1111-111111111111"
	vnd2ID := "v2222222-2222-2222-2222-222222222222"
	vnd3ID := "v3333333-3333-3333-3333-333333333333"

	cat1ID := "c1111111-1111-1111-1111-111111111111"
	cat2ID := "c2222222-2222-2222-2222-222222222222"
	cat3ID := "c3333333-3333-3333-3333-333333333333"
	cat4ID := "c4444444-4444-4444-4444-444444444444"

	loc1ID := "w1111111-1111-1111-1111-111111111111"
	loc2ID := "w2222222-2222-2222-2222-222222222222"
	loc3ID := "w3333333-3333-3333-3333-333333333333"

	itm1ID := "i1111111-1111-1111-1111-111111111111"
	itm2ID := "i2222222-2222-2222-2222-222222222222"
	itm3ID := "i3333333-3333-3333-3333-333333333333"
	itm4ID := "i4444444-4444-4444-4444-444444444444"

	emp1ID := "e1111111-1111-1111-1111-111111111111"
	emp2ID := "e2222222-2222-2222-2222-222222222222"
	dept1ID := "d1111111-1111-1111-1111-111111111111"
	dept2ID := "d2222222-2222-2222-2222-222222222222"

	// 1. Master Vendor / Supplier (mst_sc_vendors)
	var countVendor int64
	db.Model(&model.MstScVendor{}).Count(&countVendor)
	if countVendor == 0 {
		vendors := []model.MstScVendor{
			{
				ID:             vnd1ID,
				Code:           "VND-001",
				VendorName:     "PT Sukses Kopi Nusantara",
				VendorCategory: "COFFEE_BEANS",
				ContactPerson:  "Budi Santoso",
				Phone:          "081298765432",
				PicPhone:       "081298765432",
				Email:          "sales@sukseskopi.id",
				Address:        "Kawasan Agro Industri No. 12, Pangalengan, Bandung",
				BankName:       "BCA",
				BankAccount:    "7700123456",
				ActiveStatus:   "ACTIVE",
				Status:         "ACTIVE",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             vnd2ID,
				Code:           "VND-002",
				VendorName:     "CV Artisan Packaging Indo",
				VendorCategory: "PACKAGING",
				ContactPerson:  "Maya Indah",
				Phone:          "081322334455",
				PicPhone:       "081322334455",
				Email:          "order@artisanpackaging.co.id",
				Address:        "Jl. Industri Cimahi Blok C-4, Cimahi",
				BankName:       "Mandiri",
				BankAccount:    "1310088776655",
				ActiveStatus:   "ACTIVE",
				Status:         "ACTIVE",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             vnd3ID,
				Code:           "VND-003",
				VendorName:     "PT Mitra Hardware & Audio Visual",
				VendorCategory: "IT_HARDWARE",
				ContactPerson:  "Kevin Wijaya",
				Phone:          "081766554433",
				PicPhone:       "081766554433",
				Email:          "corporate@mitrahardware.com",
				Address:        "Mall Elektronik Bandung Lt. 3 No. 45, Bandung",
				BankName:       "BCA",
				BankAccount:    "8005544332",
				ActiveStatus:   "ACTIVE",
				Status:         "ACTIVE",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		_ = db.Create(&vendors).Error
		log.Printf("[SEEDER] 3 Vendor SC berhasil dibuat")
	}

	// 2. Master Kategori Item (mst_sc_item_categories)
	var countCategory int64
	db.Model(&model.MstScItemCategory{}).Count(&countCategory)
	if countCategory == 0 {
		categories := []model.MstScItemCategory{
			{
				ID:           cat1ID,
				Code:         "CAT-BEAN",
				CategoryName: "Bahan Baku Biji Kopi",
				ItemType:     "HARDWARE",
				Description:  "Green beans dan roasted beans pilihan nusantara.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			{
				ID:           cat2ID,
				Code:         "CAT-PKG",
				CategoryName: "Kemasan & Takeaway Supplies",
				ItemType:     "HARDWARE",
				Description:  "Paper cup, tote bag, sleeve kardus, dan sedotan ramah lingkungan.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			{
				ID:           cat3ID,
				Code:         "CAT-IT",
				CategoryName: "Perangkat Keras Studio & IT",
				ItemType:     "HARDWARE",
				Description:  "Komputer kasir POS, mesin display studio, laptop desain.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			{
				ID:           cat4ID,
				Code:         "CAT-SOFT",
				CategoryName: "Software & Lisensi Digital",
				ItemType:     "SOFTWARE_LICENSE",
				Description:  "Langganan sistem cloud, lisensi software kreatif & musik.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}
		_ = db.Create(&categories).Error
		log.Printf("[SEEDER] 4 Kategori Item SC berhasil dibuat")
	}

	// 3. Master Lokasi Inventaris (mst_sc_locations)
	var countLoc int64
	db.Model(&model.MstScLocation{}).Count(&countLoc)
	if countLoc == 0 {
		locations := []model.MstScLocation{
			{
				ID:               loc1ID,
				Code:             "LOC-WH-01",
				LocationName:     "Gudang Bahan Baku Utama",
				Type:             "warehouse",
				Address:          "Lt. Basement Kafeinarts Studio, Bandung",
				PersonInChargeID: &emp2ID,
				Description:      "Gudang penyimpanan karung biji kopi dan stok cup karton.",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               loc2ID,
				Code:             "LOC-BAR-01",
				LocationName:     "Bar Counter Lantai 1",
				Type:             "outlet",
				Address:          "Lt. 1 Kafeinarts Coffee Experience, Bandung",
				PersonInChargeID: &emp2ID,
				Description:      "Area display mesin espresso, grinder, dan racikan minuman harian.",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               loc3ID,
				Code:             "LOC-STD-02",
				LocationName:     "Studio Seni & Galeri Lantai 2",
				Type:             "office",
				Address:          "Lt. 2 Kafeinarts Art Workshop, Bandung",
				PersonInChargeID: &emp1ID,
				Description:      "Studio kerja perancangan kreatif, galeri kanvas, dan ruang meeting.",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}
		_ = db.Create(&locations).Error
		log.Printf("[SEEDER] 3 Lokasi Inventaris SC berhasil dibuat")
	}

	// 4. Master Item Katalog (mst_sc_items)
	var countItem int64
	db.Model(&model.MstScItem{}).Count(&countItem)
	if countItem == 0 {
		items := []model.MstScItem{
			{
				ID:              itm1ID,
				CategoryID:      cat1ID,
				Code:            "ITM-BEANS-01",
				ItemName:        "House Blend Espresso Beans 1kg",
				Unit:            "Pack",
				EstimatedPrice:  180000,
				DefaultVendorID: &vnd1ID,
				MinStock:        10,
				CurrentStock:    45,
				Description:     "Biji kopi arabika-robusta racikan signature untuk menu kopi susu & espresso.",
				Status:          "ACTIVE",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              itm2ID,
				CategoryID:      cat2ID,
				Code:            "ITM-CUP-01",
				ItemName:        "Takeaway Paper Cup 8oz Kafeinarts",
				Unit:            "Dus",
				EstimatedPrice:  350000,
				DefaultVendorID: &vnd2ID,
				MinStock:        5,
				CurrentStock:    20,
				Description:     "Dus berisi 1000 pcs cup kertas food grade berdesain grafis studio.",
				Status:          "ACTIVE",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              itm3ID,
				CategoryID:      cat1ID,
				Code:            "ITM-SYRUP-01",
				ItemName:        "Artisan Vanilla Syrup 750ml",
				Unit:            "Botol",
				EstimatedPrice:  125000,
				DefaultVendorID: &vnd1ID,
				MinStock:        5,
				CurrentStock:    22,
				Description:     "Sirup vanilla aroma alami untuk mocktail dan menu kopi rasa.",
				Status:          "ACTIVE",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
			{
				ID:              itm4ID,
				CategoryID:      cat3ID,
				Code:            "ITM-LAP-01",
				ItemName:        "Apple MacBook Pro 14 M3",
				Unit:            "Unit",
				EstimatedPrice:  28000000,
				DefaultVendorID: &vnd3ID,
				MinStock:        1,
				CurrentStock:    3,
				Description:     "Laptop workstation untuk tim desain 3D dan creative art studio.",
				Status:          "ACTIVE",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
			},
		}
		_ = db.Create(&items).Error
		log.Printf("[SEEDER] 4 Item Katalog SC berhasil dibuat")
	}

	// Sinkronisasi produk legacy
	SeedProducts(db)

	// 5. Transaksi Permintaan Pembelian (trx_sc_purchase_requests & items)
	var countPr int64
	db.Model(&model.TrxScPurchaseRequest{}).Count(&countPr)
	if countPr == 0 {
		pr1ID := uuid.NewString()
		pr2ID := uuid.NewString()
		pr3ID := uuid.NewString()

		prs := []model.TrxScPurchaseRequest{
			{
				ID:             pr1ID,
				PrNumber:       "PR-202609-001",
				DepartmentID:   dept2ID,
				RequesterName:  "Rian Anggoro",
				RequestDate:    "2026-09-01",
				Priority:       "HIGH",
				Notes:          "Restock rutin biji kopi dan sirup untuk lonjakan pengunjung akhir pekan.",
				ApprovalStatus: "SUBMITTED",
				Status:         "SUBMITTED",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Items: []model.TrxScPurchaseRequestItem{
					{ID: uuid.NewString(), PrID: pr1ID, ItemID: itm1ID, Quantity: 30, Unit: "Pack", EstimatedPrice: 180000, Subtotal: 5400000, Notes: "House blend espresso"},
					{ID: uuid.NewString(), PrID: pr1ID, ItemID: itm3ID, Quantity: 10, Unit: "Botol", EstimatedPrice: 125000, Subtotal: 1250000, Notes: "Sirup artisan vanilla"},
				},
			},
			{
				ID:             pr2ID,
				PrNumber:       "PR-202609-002",
				DepartmentID:   dept2ID,
				RequesterName:  "Rian Anggoro",
				RequestDate:    "2026-09-03",
				Priority:       "MEDIUM",
				Notes:          "Pengadaan cup takeaway cetak edisi festival seni bulan depan.",
				ApprovalStatus: "APPROVED",
				Status:         "APPROVED",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Items: []model.TrxScPurchaseRequestItem{
					{ID: uuid.NewString(), PrID: pr2ID, ItemID: itm2ID, Quantity: 10, Unit: "Dus", EstimatedPrice: 350000, Subtotal: 3500000, Notes: "Dus 1000 pcs"},
				},
			},
			{
				ID:             pr3ID,
				PrNumber:       "PR-202609-003",
				DepartmentID:   dept1ID,
				RequesterName:  "Hanif Maulana",
				RequestDate:    "2026-09-05",
				Priority:       "HIGH",
				Notes:          "Pengadaan laptop baru untuk penambahan staf animator visual studio.",
				ApprovalStatus: "SUBMITTED",
				Status:         "SUBMITTED",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
				Items: []model.TrxScPurchaseRequestItem{
					{ID: uuid.NewString(), PrID: pr3ID, ItemID: itm4ID, Quantity: 1, Unit: "Unit", EstimatedPrice: 28000000, Subtotal: 28000000, Notes: "MacBook Pro M3"},
				},
			},
		}
		_ = db.Create(&prs).Error
		log.Printf("[SEEDER] 3 Purchase Request SC & rincian item berhasil dibuat")
	}

	// 6. Transaksi Pesanan Pembelian (trx_sc_purchase_orders & items)
	var countPo int64
	db.Model(&model.TrxScPurchaseOrder{}).Count(&countPo)
	if countPo == 0 {
		po1ID := uuid.NewString()
		po2ID := uuid.NewString()
		po3ID := uuid.NewString()

		pos := []model.TrxScPurchaseOrder{
			{
				ID:                   po1ID,
				PoNumber:             "PO-202609-001",
				VendorID:             vnd1ID,
				PoDate:               "2026-09-02",
				ExpectedDeliveryDate: "2026-09-05",
				Subtotal:             6650000,
				TaxAmount:            731500,
				TotalAmount:          7381500,
				Status:               "SENT",
				PoStatus:             "SENT",
				Notes:                "Pengiriman via armada ekspedisi vendor ke lokasi basement.",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				Items: []model.TrxScPurchaseOrderItem{
					{ID: uuid.NewString(), PoID: po1ID, ItemID: itm1ID, Quantity: 30, Unit: "Pack", UnitPrice: 180000, TotalPrice: 5400000},
					{ID: uuid.NewString(), PoID: po1ID, ItemID: itm3ID, Quantity: 10, Unit: "Botol", UnitPrice: 125000, TotalPrice: 1250000},
				},
			},
			{
				ID:                   po2ID,
				PoNumber:             "PO-202609-002",
				VendorID:             vnd2ID,
				PoDate:               "2026-09-04",
				ExpectedDeliveryDate: "2026-09-08",
				Subtotal:             3500000,
				TaxAmount:            385000,
				TotalAmount:          3885000,
				Status:               "COMPLETED",
				PoStatus:             "COMPLETED",
				Notes:                "Pesanan takeaway cup telah lunas dan siap kirim.",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				Items: []model.TrxScPurchaseOrderItem{
					{ID: uuid.NewString(), PoID: po2ID, ItemID: itm2ID, Quantity: 10, Unit: "Dus", UnitPrice: 350000, TotalPrice: 3500000},
				},
			},
			{
				ID:                   po3ID,
				PoNumber:             "PO-202609-003",
				VendorID:             vnd3ID,
				PoDate:               "2026-09-06",
				ExpectedDeliveryDate: "2026-09-10",
				Subtotal:             28000000,
				TaxAmount:            3080000,
				TotalAmount:          31080000,
				Status:               "DRAFT",
				PoStatus:             "DRAFT",
				Notes:                "Menunggu konfirmasi ketersediaan stok warna space grey.",
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
				Items: []model.TrxScPurchaseOrderItem{
					{ID: uuid.NewString(), PoID: po3ID, ItemID: itm4ID, Quantity: 1, Unit: "Unit", UnitPrice: 28000000, TotalPrice: 28000000},
				},
			},
		}
		_ = db.Create(&pos).Error
		log.Printf("[SEEDER] 3 Purchase Order SC & rincian item berhasil dibuat")
	}

	// 7. Transaksi Penerimaan Barang (trx_sc_goods_receipts & items)
	var countGr int64
	db.Model(&model.TrxScGoodsReceipt{}).Count(&countGr)
	if countGr == 0 {
		gr1ID := uuid.NewString()
		gr2ID := uuid.NewString()
		gr3ID := uuid.NewString()

		// Cari salah satu PO
		var po model.TrxScPurchaseOrder
		_ = db.First(&po).Error
		targetPoID := po.ID
		if targetPoID == "" {
			targetPoID = uuid.NewString()
		}

		grs := []model.TrxScGoodsReceipt{
			{
				ID:                 gr1ID,
				ReceiptNumber:      "GR-202609-001",
				PoID:               targetPoID,
				ReceivedBy:         "Rian Anggoro",
				ReceiptDate:        "2026-09-05",
				DeliveryNoteNumber: "SJ-SKN-99812",
				Status:             "COMPLETE",
				Notes:              "Biji kopi datang dalam kondisi segel rapat dan aroma segar.",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				Items: []model.TrxScGoodsReceiptItem{
					{ID: uuid.NewString(), GrID: gr1ID, ItemID: itm1ID, OrderedQuantity: 30, ReceivedQuantity: 30, Condition: "GOOD", Notes: "Segel utuh."},
				},
			},
			{
				ID:                 gr2ID,
				ReceiptNumber:      "GR-202609-002",
				PoID:               targetPoID,
				ReceivedBy:         "Rian Anggoro",
				ReceiptDate:        "2026-09-06",
				DeliveryNoteNumber: "SJ-API-44512",
				Status:             "COMPLETE",
				Notes:              "Dus cup tiba tanpa basah, printing logo presisi.",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				Items: []model.TrxScGoodsReceiptItem{
					{ID: uuid.NewString(), GrID: gr2ID, ItemID: itm2ID, OrderedQuantity: 10, ReceivedQuantity: 10, Condition: "GOOD", Notes: "Kardus bersih."},
				},
			},
			{
				ID:                 gr3ID,
				ReceiptNumber:      "GR-202609-003",
				PoID:               targetPoID,
				ReceivedBy:         "Hanif Maulana",
				ReceiptDate:        "2026-09-09",
				DeliveryNoteNumber: "SJ-MHW-11203",
				Status:             "COMPLETE",
				Notes:              "Unit MacBook diterima lengkap dengan kartu garansi resmi Apple iBox.",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
				Items: []model.TrxScGoodsReceiptItem{
					{ID: uuid.NewString(), GrID: gr3ID, ItemID: itm4ID, OrderedQuantity: 1, ReceivedQuantity: 1, Condition: "GOOD", Notes: "Garansi terdaftar."},
				},
			},
		}
		_ = db.Create(&grs).Error
		log.Printf("[SEEDER] 3 Goods Receipt SC & rincian item berhasil dibuat")
	}

	// 8. Transaksi Aset Inventaris (trx_sc_assets)
	var countAsset int64
	db.Model(&model.TrxScAsset{}).Count(&countAsset)
	if countAsset == 0 {
		assets := []model.TrxScAsset{
			{
				ID:                      uuid.NewString(),
				AssetCode:               "AST-ESP-001",
				Name:                    "Mesin Espresso La Marzocco Linea Mini",
				ItemID:                  itm1ID,
				CategoryID:              cat3ID,
				SerialOrLicenseKey:      "LM-LM-2024-9981",
				AcquisitionDate:         "2024-02-10",
				AcquisitionValue:        75000000,
				WarrantyOrLicenseExpiry: "2027-02-10",
				AssignedToEmployeeID:    &emp2ID,
				DepartmentID:            &dept2ID,
				LocationID:              &loc2ID,
				Status:                  "IN_USE",
				Condition:               "GOOD",
				Notes:                   "Mesin espresso utama bar, rutin diservis dan backflush berkala.",
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			},
			{
				ID:                      uuid.NewString(),
				AssetCode:               "AST-GRN-001",
				Name:                    "Grinder Mahlkonig EK43 Allround",
				ItemID:                  itm1ID,
				CategoryID:              cat3ID,
				SerialOrLicenseKey:      "MK-EK43-8872",
				AcquisitionDate:         "2024-02-15",
				AcquisitionValue:        45000000,
				WarrantyOrLicenseExpiry: "2026-02-15",
				AssignedToEmployeeID:    &emp2ID,
				DepartmentID:            &dept2ID,
				LocationID:              &loc2ID,
				Status:                  "IN_USE",
				Condition:               "GOOD",
				Notes:                   "Grinder presisi tinggi untuk filter brew dan cupping kopi.",
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			},
			{
				ID:                      uuid.NewString(),
				AssetCode:               "AST-MBP-001",
				Name:                    "Apple MacBook Pro 14 M3 Studio",
				ItemID:                  itm4ID,
				CategoryID:              cat3ID,
				SerialOrLicenseKey:      "C02G9981MD6R",
				AcquisitionDate:         "2024-05-20",
				AcquisitionValue:        28000000,
				WarrantyOrLicenseExpiry: "2027-05-20",
				AssignedToEmployeeID:    &emp1ID,
				DepartmentID:            &dept1ID,
				LocationID:              &loc3ID,
				Status:                  "IN_USE",
				Condition:               "GOOD",
				Notes:                   "Laptop kerja utama Creative Director untuk animasi dan packaging.",
				CreatedAt:               time.Now(),
				UpdatedAt:               time.Now(),
			},
		}
		_ = db.Create(&assets).Error
		log.Printf("[SEEDER] 3 Aset Inventaris SC berhasil dibuat")
	}

	// 9. Transaksi Langganan Software & Lisensi (trx_sc_subscriptions)
	var countSub int64
	db.Model(&model.TrxScSubscription{}).Count(&countSub)
	if countSub == 0 {
		subs := []model.TrxScSubscription{
			{
				ID:           uuid.NewString(),
				LicenseCode:  "LIC-ADOBE-CC",
				Name:         "Adobe Creative Cloud All Apps for Teams",
				ItemID:       itm4ID,
				VendorID:     vnd3ID,
				BillingCycle: "YEARLY",
				Amount:       14400000,
				StartDate:    "2026-01-01",
				EndDate:      "2026-12-31",
				AutoRenewal:  true,
				Status:       "ACTIVE",
				Notes:        "Lisensi Photoshop, Illustrator, Premiere untuk departemen kreatif studio.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			{
				ID:           uuid.NewString(),
				LicenseCode:  "LIC-SPOTIFY-BIZ",
				Name:         "Spotify for Business Bar Soundtrack",
				ItemID:       itm4ID,
				VendorID:     vnd3ID,
				BillingCycle: "MONTHLY",
				Amount:       450000,
				StartDate:    "2026-09-01",
				EndDate:      "2026-09-30",
				AutoRenewal:  true,
				Status:       "ACTIVE",
				Notes:        "Lisensi streaming komersial musik latar kafeinarts.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
			{
				ID:           uuid.NewString(),
				LicenseCode:  "LIC-GSUITE-15",
				Name:         "Google Workspace Business Plus 15 User",
				ItemID:       itm4ID,
				VendorID:     vnd3ID,
				BillingCycle: "YEARLY",
				Amount:       27000000,
				StartDate:    "2026-03-01",
				EndDate:      "2027-02-28",
				AutoRenewal:  true,
				Status:       "ACTIVE",
				Notes:        "Email domain @kafeinarts.com & penyimpanan drive 5TB bersama.",
				CreatedAt:    time.Now(),
				UpdatedAt:    time.Now(),
			},
		}
		_ = db.Create(&subs).Error
		log.Printf("[SEEDER] 3 Langganan Lisensi SC berhasil dibuat")
	}

	log.Println("[SEEDER] Seluruh data Supply Chain (SC) siap digunakan.")
}

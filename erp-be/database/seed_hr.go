package database

import (
	"log"
	"time"

	"erp-kafeinart/model"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// SeedHrAll menginisialisasi seluruh master dan transaksi modul Human Resources (HR).
func SeedHrAll(db *gorm.DB) {
	log.Println("[SEEDER] Memulai inisialisasi dummy data Human Resources (HR)...")

	dept1ID := "d1111111-1111-1111-1111-111111111111"
	dept2ID := "d2222222-2222-2222-2222-222222222222"
	dept3ID := "d3333333-3333-3333-3333-333333333333"

	pos1ID := "p1111111-1111-1111-1111-111111111111"
	pos2ID := "p2222222-2222-2222-2222-222222222222"
	pos3ID := "p3333333-3333-3333-3333-333333333333"
	pos4ID := "p4444444-4444-4444-4444-444444444444"

	emp1ID := "e1111111-1111-1111-1111-111111111111"
	emp2ID := "e2222222-2222-2222-2222-222222222222"
	emp3ID := "e3333333-3333-3333-3333-333333333333"
	emp4ID := "e4444444-4444-4444-4444-444444444444"

	lt1ID := "l1111111-1111-1111-1111-111111111111"
	lt2ID := "l2222222-2222-2222-2222-222222222222"
	lt3ID := "l3333333-3333-3333-3333-333333333333"

	sc1ID := "s1111111-1111-1111-1111-111111111111"
	sc2ID := "s2222222-2222-2222-2222-222222222222"
	sc3ID := "s3333333-3333-3333-3333-333333333333"

	// 1. Master Departemen (mst_hr_departments)
	var countDept int64
	db.Model(&model.MstHrDepartment{}).Count(&countDept)
	if countDept == 0 {
		departments := []model.MstHrDepartment{
			{
				ID:             dept1ID,
				Code:           "DEP-CRT",
				DepartmentName: "Studio & Creative",
				HeadEmployeeID: &emp1ID,
				Description:    "Divisi perancangan visual, seni kurasi kopi, dan desain studio Kafeinarts.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             dept2ID,
				Code:           "DEP-OPS",
				DepartmentName: "Bar & Operations",
				HeadEmployeeID: &emp2ID,
				Description:    "Divisi operasional bar kopi, pelayanan pelanggan, dan peracikan menu harian.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             dept3ID,
				Code:           "DEP-FIN",
				DepartmentName: "Finance & General Affairs",
				HeadEmployeeID: &emp3ID,
				Description:    "Divisi pengelolaan keuangan, akuntansi, legalitas, dan administrasi umum.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		_ = db.Create(&departments).Error
		log.Printf("[SEEDER] 3 Departemen HR berhasil dibuat")
	}

	// 2. Master Jabatan (mst_hr_positions)
	var countPos int64
	db.Model(&model.MstHrPosition{}).Count(&countPos)
	if countPos == 0 {
		positions := []model.MstHrPosition{
			{
				ID:             pos1ID,
				DepartmentID:   dept1ID,
				Code:           "POS-CD",
				PositionName:   "Creative Director",
				Level:          "MANAGER",
				JobDescription: "Memimpin konsep kreatif, desain interior studio, dan branding Kafeinarts.",
				Description:    "Penanggung jawab artistik dan identitas visual perusahaan.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             pos2ID,
				DepartmentID:   dept2ID,
				Code:           "POS-HBAR",
				PositionName:   "Head Barista & Roaster",
				Level:          "SUPERVISOR",
				JobDescription: "Mengawasi profil sangrai biji kopi, kalibrasi rasa harian, dan standarisasi bar.",
				Description:    "Lead operational tim barista dan quality control rasa minuman.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             pos3ID,
				DepartmentID:   dept3ID,
				Code:           "POS-FIN",
				PositionName:   "Finance Specialist",
				Level:          "STAFF",
				JobDescription: "Melakukan rekonsiliasi arus kas, jurnal umum, serta faktur tagihan.",
				Description:    "Pelaksana pembukuan dan administrasi finansial terpusat.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             pos4ID,
				DepartmentID:   dept2ID,
				Code:           "POS-BAR",
				PositionName:   "Junior Barista",
				Level:          "STAFF",
				JobDescription: "Menyiapkan pesanan minuman, melayani kasir POS, dan menjaga kebersihan bar.",
				Description:    "Kru bar operasional frontline.",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		_ = db.Create(&positions).Error
		log.Printf("[SEEDER] 4 Jabatan HR berhasil dibuat")
	}

	// 3. Master Karyawan (mst_hr_employees & legacy model.Employee)
	var countEmp int64
	db.Model(&model.MstHrEmployee{}).Count(&countEmp)
	if countEmp == 0 {
		employees := []model.MstHrEmployee{
			{
				ID:                emp1ID,
				EmployeeNumber:    "EMP-001",
				FullName:          "Hanif Maulana",
				NationalIDNumber:  "3273010101950001",
				BirthDate:         "1995-04-12",
				Gender:            "Laki-laki",
				Address:           "Jl. Dago Asri No. 18, Bandung",
				Email:             "hanif@kafeinarts.com",
				PhoneNumber:       "081234567890",
				DepartmentID:      &dept1ID,
				PositionID:        &pos1ID,
				JoinDate:          "2023-01-15",
				EmploymentStatus:  "PERMANENT",
				ActiveStatus:      "ACTIVE",
				Salary:            15000000,
				BankName:          "BCA",
				BankAccountNumber: "8001234567",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				ID:                emp2ID,
				EmployeeNumber:    "EMP-002",
				FullName:          "Rian Anggoro",
				NationalIDNumber:  "3273020503970002",
				BirthDate:         "1997-08-20",
				Gender:            "Laki-laki",
				Address:           "Jl. Riau No. 45, Bandung",
				Email:             "rian@kafeinarts.com",
				PhoneNumber:       "081234567891",
				DepartmentID:      &dept2ID,
				PositionID:        &pos2ID,
				JoinDate:          "2023-03-01",
				EmploymentStatus:  "PERMANENT",
				ActiveStatus:      "ACTIVE",
				Salary:            7500000,
				BankName:          "Mandiri",
				BankAccountNumber: "1310009876543",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				ID:                emp3ID,
				EmployeeNumber:    "EMP-003",
				FullName:          "Siti Rahmawati",
				NationalIDNumber:  "3273031209990003",
				BirthDate:         "1999-11-05",
				Gender:            "Perempuan",
				Address:           "Jl. Buah Batu No. 102, Bandung",
				Email:             "siti@kafeinarts.com",
				PhoneNumber:       "081234567892",
				DepartmentID:      &dept3ID,
				PositionID:        &pos3ID,
				JoinDate:          "2023-06-15",
				EmploymentStatus:  "PERMANENT",
				ActiveStatus:      "ACTIVE",
				Salary:            7000000,
				BankName:          "BCA",
				BankAccountNumber: "8009876543",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				ID:                emp4ID,
				EmployeeNumber:    "EMP-004",
				FullName:          "Dimas Prasetyo",
				NationalIDNumber:  "3273041802010004",
				BirthDate:         "2001-02-18",
				Gender:            "Laki-laki",
				Address:           "Jl. Sukajadi No. 77, Bandung",
				Email:             "dimas@kafeinarts.com",
				PhoneNumber:       "081234567893",
				DepartmentID:      &dept2ID,
				PositionID:        &pos4ID,
				JoinDate:          "2024-01-10",
				EmploymentStatus:  "CONTRACT",
				ActiveStatus:      "ACTIVE",
				Salary:            5200000,
				BankName:          "BRI",
				BankAccountNumber: "001201029384756",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		}
		_ = db.Create(&employees).Error
		log.Printf("[SEEDER] 4 Karyawan HR berhasil dibuat")
	}

	// Sinkronisasi legacy Employee model
	SeedEmployees(db)

	// 4. Master Jenis Cuti (mst_hr_leave_types)
	var countLeaveType int64
	db.Model(&model.MstHrLeaveType{}).Count(&countLeaveType)
	if countLeaveType == 0 {
		leaveTypes := []model.MstHrLeaveType{
			{
				ID:                 lt1ID,
				Code:               "CUTI-ANNUAL",
				LeaveName:          "Cuti Tahunan",
				AnnualQuotaDays:    12,
				IsPaid:             true,
				RequiresAttachment: false,
				Description:        "Hak cuti tahunan berbayar karyawan setelah masa kerja satu tahun.",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			{
				ID:                 lt2ID,
				Code:               "CUTI-SICK",
				LeaveName:          "Cuti Sakit",
				AnnualQuotaDays:    14,
				IsPaid:             true,
				RequiresAttachment: true,
				Description:        "Izin sakit berbayar dengan melampirkan surat keterangan dokter resmi.",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
			{
				ID:                 lt3ID,
				Code:               "CUTI-SPECIAL",
				LeaveName:          "Cuti Menikah & Khusus",
				AnnualQuotaDays:    3,
				IsPaid:             true,
				RequiresAttachment: false,
				Description:        "Cuti khusus pernikahan karyawan, duka cita, atau persalinan keluarga.",
				CreatedAt:          time.Now(),
				UpdatedAt:          time.Now(),
			},
		}
		_ = db.Create(&leaveTypes).Error
		log.Printf("[SEEDER] 3 Jenis Cuti HR berhasil dibuat")
	}

	// 5. Master Komponen Gaji (mst_hr_salary_components)
	var countSalaryComp int64
	db.Model(&model.MstHrSalaryComponent{}).Count(&countSalaryComp)
	if countSalaryComp == 0 {
		salaryComponents := []model.MstHrSalaryComponent{
			{
				ID:                sc1ID,
				Code:              "KOMP-POKOK",
				ComponentName:     "Gaji Pokok",
				ComponentType:     "EARNING",
				CalculationMethod: "FIXED",
				IsTaxable:         true,
				Description:       "Penghasilan dasar tetap setiap bulannya.",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				ID:                sc2ID,
				Code:              "KOMP-TRANS",
				ComponentName:     "Tunjangan Transport & Konsumsi",
				ComponentType:     "EARNING",
				CalculationMethod: "FIXED",
				IsTaxable:         false,
				Description:       "Tunjangan operasional harian kehadiran kerja.",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
			{
				ID:                sc3ID,
				Code:              "KOMP-BPJS",
				ComponentName:     "Iuran BPJS Ketenagakerjaan & Kesehatan",
				ComponentType:     "DEDUCTION",
				CalculationMethod: "FIXED",
				IsTaxable:         false,
				Description:       "Potongan premi jaminan sosial karyawan.",
				CreatedAt:         time.Now(),
				UpdatedAt:         time.Now(),
			},
		}
		_ = db.Create(&salaryComponents).Error
		log.Printf("[SEEDER] 3 Komponen Gaji HR berhasil dibuat")
	}

	// 6. Master Jadwal Shift Kerja (mst_hr_work_shifts)
	var countShift int64
	db.Model(&model.MstHrWorkShift{}).Count(&countShift)
	if countShift == 0 {
		shifts := []model.MstHrWorkShift{
			{
				ID:                   uuid.NewString(),
				ShiftName:            "Office Regular 08-17",
				TimeIn:               "08:00",
				TimeOut:              "17:00",
				LateToleranceMinutes: 15,
				WorkDays:             model.StringArray{"Senin", "Selasa", "Rabu", "Kamis", "Jumat"},
				IsNightShift:         false,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			{
				ID:                   uuid.NewString(),
				ShiftName:            "Morning Barista 07-15",
				TimeIn:               "07:00",
				TimeOut:              "15:30",
				LateToleranceMinutes: 10,
				WorkDays:             model.StringArray{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"},
				IsNightShift:         false,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
			{
				ID:                   uuid.NewString(),
				ShiftName:            "Evening Studio & Bar 14-22",
				TimeIn:               "14:30",
				TimeOut:              "22:30",
				LateToleranceMinutes: 10,
				WorkDays:             model.StringArray{"Senin", "Selasa", "Rabu", "Kamis", "Jumat", "Sabtu", "Minggu"},
				IsNightShift:         true,
				CreatedAt:            time.Now(),
				UpdatedAt:            time.Now(),
			},
		}
		_ = db.Create(&shifts).Error
		log.Printf("[SEEDER] 3 Shift Kerja HR berhasil dibuat")
	}

	// 7. Transaksi Presensi (trx_hr_attendances)
	var countAtt int64
	db.Model(&model.TrxHrAttendance{}).Count(&countAtt)
	if countAtt == 0 {
		todayStr := time.Now().Format("2006-01-02")
		attendances := []model.TrxHrAttendance{
			{
				ID:               uuid.NewString(),
				EmployeeID:       emp1ID,
				AttendanceDate:   todayStr,
				CheckInTime:      "07:55",
				CheckOutTime:     "17:05",
				CheckInLat:       -6.9024,
				CheckInLng:       107.6188,
				AttendanceStatus: "PRESENT",
				Notes:            "Tepat waktu, standby rapat creative showcase studio.",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               uuid.NewString(),
				EmployeeID:       emp2ID,
				AttendanceDate:   todayStr,
				CheckInTime:      "06:50",
				CheckOutTime:     "15:35",
				CheckInLat:       -6.9024,
				CheckInLng:       107.6188,
				AttendanceStatus: "PRESENT",
				Notes:            "Opening bar barista & kalibrasi espresso pagi selesai.",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
			{
				ID:               uuid.NewString(),
				EmployeeID:       emp3ID,
				AttendanceDate:   todayStr,
				CheckInTime:      "07:58",
				CheckOutTime:     "17:02",
				CheckInLat:       -6.9024,
				CheckInLng:       107.6188,
				AttendanceStatus: "PRESENT",
				Notes:            "Hadir di kantor finance, proses pembukuan.",
				CreatedAt:        time.Now(),
				UpdatedAt:        time.Now(),
			},
		}
		_ = db.Create(&attendances).Error
		log.Printf("[SEEDER] 3 Data Presensi HR berhasil dibuat")
	}

	// 8. Transaksi Pengajuan Cuti (trx_hr_leave_requests)
	var countLeaveReq int64
	db.Model(&model.TrxHrLeaveRequest{}).Count(&countLeaveReq)
	if countLeaveReq == 0 {
		leaveRequests := []model.TrxHrLeaveRequest{
			{
				ID:             uuid.NewString(),
				RequestNumber:  "LV-202609-001",
				EmployeeID:     emp1ID,
				LeaveTypeID:    lt1ID,
				StartDate:      "2026-10-05",
				EndDate:        "2026-10-07",
				TotalDays:      3,
				Reason:         "Rencana liburan keluarga luar kota.",
				ApprovalStatus: "SUBMITTED",
				Status:         "SUBMITTED",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             uuid.NewString(),
				RequestNumber:  "LV-202609-002",
				EmployeeID:     emp2ID,
				LeaveTypeID:    lt2ID,
				StartDate:      "2026-09-10",
				EndDate:        "2026-09-11",
				TotalDays:      2,
				Reason:         "Sakit demam & istirahat dokter terlampir.",
				ApprovalStatus: "APPROVED",
				Status:         "APPROVED",
				ApprovedBy:     emp1ID,
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
			{
				ID:             uuid.NewString(),
				RequestNumber:  "LV-202609-003",
				EmployeeID:     emp3ID,
				LeaveTypeID:    lt3ID,
				StartDate:      "2026-11-12",
				EndDate:        "2026-11-14",
				TotalDays:      3,
				Reason:         "Keperluan acara resepsi pernikahan keluarga kandung.",
				ApprovalStatus: "SUBMITTED",
				Status:         "SUBMITTED",
				CreatedAt:      time.Now(),
				UpdatedAt:      time.Now(),
			},
		}
		_ = db.Create(&leaveRequests).Error
		log.Printf("[SEEDER] 3 Pengajuan Cuti HR berhasil dibuat")
	}

	// 9. Transaksi Penggajian (trx_hr_payrolls & trx_hr_payroll_details)
	var countPayroll int64
	db.Model(&model.TrxHrPayroll{}).Count(&countPayroll)
	if countPayroll == 0 {
		pay1ID := uuid.NewString()
		pay2ID := uuid.NewString()
		pay3ID := uuid.NewString()

		payrolls := []model.TrxHrPayroll{
			{
				ID:              pay1ID,
				PayrollNumber:   "PAY-202608-001",
				PayrollPeriod:   "2026-08",
				PeriodMonth:     8,
				PeriodYear:      2026,
				EmployeeID:      emp1ID,
				TotalEarnings:   16500000,
				TotalDeductions: 500000,
				NetSalary:       16000000,
				PaymentStatus:   "PAID",
				PaymentDate:     "2026-08-28",
				Notes:           "Gaji bulan Agustus 2026 ditransfer melalui BCA.",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Items: []model.TrxHrPayrollItem{
					{ID: uuid.NewString(), PayrollID: pay1ID, ComponentID: sc1ID, Amount: 15000000},
					{ID: uuid.NewString(), PayrollID: pay1ID, ComponentID: sc2ID, Amount: 1500000},
					{ID: uuid.NewString(), PayrollID: pay1ID, ComponentID: sc3ID, Amount: 500000},
				},
			},
			{
				ID:              pay2ID,
				PayrollNumber:   "PAY-202608-002",
				PayrollPeriod:   "2026-08",
				PeriodMonth:     8,
				PeriodYear:      2026,
				EmployeeID:      emp2ID,
				TotalEarnings:   8300000,
				TotalDeductions: 250000,
				NetSalary:       8050000,
				PaymentStatus:   "PAID",
				PaymentDate:     "2026-08-28",
				Notes:           "Gaji bulan Agustus 2026 ditransfer melalui Mandiri.",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Items: []model.TrxHrPayrollItem{
					{ID: uuid.NewString(), PayrollID: pay2ID, ComponentID: sc1ID, Amount: 7500000},
					{ID: uuid.NewString(), PayrollID: pay2ID, ComponentID: sc2ID, Amount: 800000},
					{ID: uuid.NewString(), PayrollID: pay2ID, ComponentID: sc3ID, Amount: 250000},
				},
			},
			{
				ID:              pay3ID,
				PayrollNumber:   "PAY-202608-003",
				PayrollPeriod:   "2026-08",
				PeriodMonth:     8,
				PeriodYear:      2026,
				EmployeeID:      emp3ID,
				TotalEarnings:   7700000,
				TotalDeductions: 230000,
				NetSalary:       7470000,
				PaymentStatus:   "PAID",
				PaymentDate:     "2026-08-28",
				Notes:           "Gaji bulan Agustus 2026 ditransfer melalui BCA.",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Items: []model.TrxHrPayrollItem{
					{ID: uuid.NewString(), PayrollID: pay3ID, ComponentID: sc1ID, Amount: 7000000},
					{ID: uuid.NewString(), PayrollID: pay3ID, ComponentID: sc2ID, Amount: 700000},
					{ID: uuid.NewString(), PayrollID: pay3ID, ComponentID: sc3ID, Amount: 230000},
				},
			},
		}
		_ = db.Create(&payrolls).Error
		log.Printf("[SEEDER] 3 Transaksi Payroll HR & detail komponen berhasil dibuat")
	}

	// 10. Transaksi Reimbursement (trx_hr_reimbursements)
	var countReimb int64
	db.Model(&model.TrxHrReimbursement{}).Count(&countReimb)
	if countReimb == 0 {
		reimbursements := []model.TrxHrReimbursement{
			{
				ID:            uuid.NewString(),
				ClaimNumber:   "CLM-202609-001",
				EmployeeID:    emp2ID,
				ClaimDate:     "2026-09-05",
				CategoryID:    "KONSUMSI",
				CategoryName:  "Bahan Bar Tambahan",
				Amount:        275000,
				Description:   "Beli es batu kristal darurat & gallon air steril bar.",
				Status:        "SUBMITTED",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uuid.NewString(),
				ClaimNumber:   "CLM-202609-002",
				EmployeeID:    emp1ID,
				ClaimDate:     "2026-09-08",
				CategoryID:    "STUDIO",
				CategoryName:  "Perlengkapan Kreatif Studio",
				Amount:        850000,
				Description:   "Pembelian pigura kayu custom untuk artwork pameran.",
				Status:        "APPROVED",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
			{
				ID:            uuid.NewString(),
				ClaimNumber:   "CLM-202609-003",
				EmployeeID:    emp3ID,
				ClaimDate:     "2026-09-12",
				CategoryID:    "TRANSPORT",
				CategoryName:  "Transportasi & Kurir",
				Amount:        125000,
				Description:   "Kirim dokumen perpajakan fisik ke kantor KPP Pratama.",
				Status:        "PAID",
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			},
		}
		_ = db.Create(&reimbursements).Error
		log.Printf("[SEEDER] 3 Klaim Reimbursement HR berhasil dibuat")
	}

	// 11. Transaksi Penilaian Kinerja (trx_hr_performance_reviews & items)
	var countReview int64
	db.Model(&model.TrxHrPerformanceReview{}).Count(&countReview)
	if countReview == 0 {
		rev1ID := uuid.NewString()
		rev2ID := uuid.NewString()
		rev3ID := uuid.NewString()

		reviews := []model.TrxHrPerformanceReview{
			{
				ID:              rev1ID,
				ReviewNumber:    "REV-2026-H1-01",
				EmployeeID:      emp1ID,
				ReviewerID:      emp1ID,
				ReviewPeriod:    "Semester 1 2026",
				ReviewDate:      "2026-07-02",
				TotalScore:      94.5,
				OverallFeedback: "Konsep artistik studio berjalan luar biasa dan mendongkrak citra premium kafe.",
				Status:          "FINALIZED",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Items: []model.TrxHrPerformanceReviewItem{
					{ID: uuid.NewString(), ReviewID: rev1ID, Criterion: "Kreativitas & Desain", Score: 95, Weight: 40, Comments: "Sangat inovatif."},
					{ID: uuid.NewString(), ReviewID: rev1ID, Criterion: "Leadership Tim", Score: 94, Weight: 30, Comments: "Mengayomi staf studio."},
					{ID: uuid.NewString(), ReviewID: rev1ID, Criterion: "Ketepatan Waktu Proyek", Score: 94, Weight: 30, Comments: "Semua launch on schedule."},
				},
			},
			{
				ID:              rev2ID,
				ReviewNumber:    "REV-2026-H1-02",
				EmployeeID:      emp2ID,
				ReviewerID:      emp1ID,
				ReviewPeriod:    "Semester 1 2026",
				ReviewDate:      "2026-07-03",
				TotalScore:      90.0,
				OverallFeedback: "Konsistensi rasa kopi tinggi, kebersihan bar terjaga sangat rapi.",
				Status:          "FINALIZED",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Items: []model.TrxHrPerformanceReviewItem{
					{ID: uuid.NewString(), ReviewID: rev2ID, Criterion: "Kualitas Racikan Bar", Score: 92, Weight: 40, Comments: "Dialing espresso akurat."},
					{ID: uuid.NewString(), ReviewID: rev2ID, Criterion: "Kecepatan Pelayanan", Score: 88, Weight: 30, Comments: "Waktu saji efisien."},
					{ID: uuid.NewString(), ReviewID: rev2ID, Criterion: "Disiplin & Presensi", Score: 90, Weight: 30, Comments: "Jarang terlambat."},
				},
			},
			{
				ID:              rev3ID,
				ReviewNumber:    "REV-2026-H1-03",
				EmployeeID:      emp3ID,
				ReviewerID:      emp1ID,
				ReviewPeriod:    "Semester 1 2026",
				ReviewDate:      "2026-07-04",
				TotalScore:      91.0,
				OverallFeedback: "Laporan keuangan akurat, closing bulanan selalu tepat pada hari kerja ke-3.",
				Status:          "FINALIZED",
				CreatedAt:       time.Now(),
				UpdatedAt:       time.Now(),
				Items: []model.TrxHrPerformanceReviewItem{
					{ID: uuid.NewString(), ReviewID: rev3ID, Criterion: "Akurasi Pembukuan", Score: 93, Weight: 40, Comments: "Nol selisih kas."},
					{ID: uuid.NewString(), ReviewID: rev3ID, Criterion: "Kepatuhan Pajak & Legal", Score: 90, Weight: 30, Comments: "Tertib aturan."},
					{ID: uuid.NewString(), ReviewID: rev3ID, Criterion: "Kerjasama Antar Tim", Score: 89, Weight: 30, Comments: "Komunikasi sangat baik."},
				},
			},
		}
		_ = db.Create(&reviews).Error
		log.Printf("[SEEDER] 3 Penilaian Kinerja HR & kriteria rincian berhasil dibuat")
	}

	log.Println("[SEEDER] Seluruh data Human Resources (HR) siap digunakan.")
}

// SeedEmployees menginisialisasi sampel data pegawai HR legacy (employees).
func SeedEmployees(db *gorm.DB) {
	var count int64
	db.Model(&model.Employee{}).Count(&count)
	if count > 0 {
		return
	}

	joinDate := time.Date(2023, 5, 1, 0, 0, 0, 0, time.Local)
	employees := []model.Employee{
		{
			NIP:        "EMP-001",
			Name:       "Hanif Maulana",
			Email:      "hanif@kafeinarts.com",
			Phone:      "081234567890",
			Department: "Studio & Creative",
			Position:   "Creative Director",
			Status:     "active",
			JoinDate:   joinDate,
			Salary:     15000000,
		},
		{
			NIP:        "EMP-002",
			Name:       "Rian Anggoro",
			Email:      "rian@kafeinarts.com",
			Phone:      "081234567891",
			Department: "Bar & Operations",
			Position:   "Head Barista",
			Status:     "active",
			JoinDate:   joinDate.AddDate(0, 2, 0),
			Salary:     7500000,
		},
		{
			NIP:        "EMP-003",
			Name:       "Siti Rahmawati",
			Email:      "siti@kafeinarts.com",
			Phone:      "081234567892",
			Department: "Finance & Admin",
			Position:   "Finance Specialist",
			Status:     "active",
			JoinDate:   joinDate.AddDate(0, 5, 0),
			Salary:     7000000,
		},
	}

	if err := db.Create(&employees).Error; err != nil {
		log.Printf("[ERROR] Gagal membuat pegawai seeder: %v", err)
	} else {
		log.Printf("[SEEDER] %d pegawai HR default berhasil ditambahkan", len(employees))
	}
}


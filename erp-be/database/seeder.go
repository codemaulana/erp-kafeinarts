package database

import (
	"gorm.io/gorm"
)

// SeedDatabase menginisialisasi seluruh data awal sistem jika database masih baru.
// Menjamin akun admin default selalu terbuat saat pertama kali server dijalankan.
func SeedDatabase(db *gorm.DB) {
	// 1. Akun Autentikasi Pengguna (Wajib)
	SeedUsers(db)

	// 2. Data Dummy Modul Human Resources (HR)
	SeedHrAll(db)

	// 3. Data Dummy Modul Supply Chain (SC)
	SeedScAll(db)

	// 4. Data Master & Transaksi Finance
	SeedFinanceMaster(db)
	SeedFinanceTransactions(db)

	// 5. Data Produk Katalog & Pegawai Legacy
	SeedProducts(db)
	SeedEmployees(db)
}

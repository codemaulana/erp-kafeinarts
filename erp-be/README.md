# Kafeinarts ERP — Backend (Golang & Fiber v3)

Backend RESTful API performa tinggi untuk sistem internal **Kafeinarts Studio**, dibangun menggunakan **Golang**, framework **Fiber v3**, **GORM**, dan database **SQLite** (menggunakan pure Go driver, tanpa ketergantungan GCC/CGO di Windows).

---

## 🚀 Fitur Utama & Modul

Backend ini dirancang dengan prinsip **Clean & Layered Architecture** yang terintegrasi langsung dengan format data frontend Quasar (`erp-fe`):

1. **Autentikasi & Otorisasi (Auth)**
   - Autentikasi berbasis token JWT (HMAC-SHA256)
   - Password hashing menggunakan `bcrypt`
   - Role & permission checking (superadmin, admin, staff)
   - Kompatibilitas multi-format login (`/api/v1/auth/login`, `/login`, dan endpoint CSRF Sanctum)
2. **Manajemen Inventaris (Inventory)**
   - CRUD Produk lengkap (`sku`, `name`, `stock`, `price`, `description`)
   - Filter pencarian nama/SKU, pagination Quasar, dan sorting dinamis
   - Pencatatan penyesuaian stok (*stock adjustment*) dengan audit log (masuk, keluar, koreksi)
3. **Keuangan Studio (Finance)**
   - Pencatatan transaksi pemasukan (*income*) dan pengeluaran (*expense*)
   - Ringkasan keuangan (*summary*): total pemasukan, pengeluaran, saldo bersih, dan total transaksi
   - Filter transaksi berdasarkan kategori dan rentang waktu
4. **Sumber Daya Manusia (HR)**
   - Manajemen staf / pegawai (`nip`, `name`, `email`, `department`, `position`, `status`, `salary`)
   - Filter pegawai berdasarkan departemen dan status aktif
5. **Keamanan & Performa**
   - Pure Go SQLite dengan mode WAL (*Write-Ahead Logging*) & *busy timeout*
   - Middleware CORS yang dapat dikonfigurasi untuk port frontend Vue/Quasar (9000, 9001, dll.)
   - AutoMigrate skema database dan Seeder data demo otomatis saat pertama kali dijalankan

---

## 📁 Struktur Direktori & Clean Code

Struktur kode memisahkan tanggung jawab secara tegas (*Separation of Concerns*):

```text
erp-be/
├── config/                  # Konfigurasi aplikasi & database
│   ├── config.go            # Load environment variables (.env / fallback)
│   └── database.go          # Inisialisasi GORM SQLite (pure Go) & AutoMigrate
├── database/                # Database seeders
│   └── seeder.go            # Data awal (superadmin, produk sampel, transaksi, pegawai)
├── model/                   # GORM Entity models
│   ├── user.go              # Model User & PermissionList
│   ├── product.go           # Model Product & StockAdjustment
│   ├── transaction.go       # Model Transaction & FinanceSummary
│   └── employee.go          # Model Employee
├── repository/              # Data Access Layer (kueri database murni)
│   ├── user_repository.go
│   ├── product_repository.go
│   ├── transaction_repository.go
│   └── employee_repository.go
├── service/                 # Business Logic Layer (validasi & aturan bisnis)
│   ├── auth_service.go
│   ├── product_service.go
│   ├── transaction_service.go
│   └── employee_service.go
├── handler/                 # HTTP Controller Layer (Fiber v3 handlers)
│   ├── auth_handler.go
│   ├── product_handler.go
│   ├── transaction_handler.go
│   └── employee_handler.go
├── middleware/              # Fiber v3 Middlewares
│   ├── auth_middleware.go   # Validasi JWT Bearer token
│   └── cors_middleware.go   # Konfigurasi Cross-Origin Resource Sharing
├── routes/                  # Pemetaan rute API
│   └── api.go               # Pendaftaran endpoint /api/v1
├── utils/                   # Helper & utilitas
│   ├── jwt.go               # Generator & validasi JWT
│   ├── password.go          # Bcrypt hashing & verifier
│   └── response.go          # Format JSON standar (Success, Paginated, Error)
├── docs/                    # Dokumentasi API terperinci
│   └── API.md               # Spesifikasi lengkap request & response
├── test/                    # Unit testing & integration testing
│   ├── api_integration_test.go # Pengujian alur API end-to-end (Fiber app.Test)
│   └── utils_test.go        # Unit test password hashing, JWT, pagination
├── storage/                 # Folder penyimpanan file SQLite (.db)
├── .env                     # File konfigurasi lokal
├── .env.example             # Template konfigurasi environment
├── go.mod / go.sum          # Manajemen dependensi Go
└── main.go                  # Entry point aplikasi
```

---

## 🛠️ Prasyarat

- **Go**: versi `1.22` atau yang lebih baru (project teruji pada `go 1.26.2 windows/amd64`)
- **GCC**: *Tidak diperlukan*, karena SQLite menggunakan pure Go driver `github.com/glebarez/sqlite`.

---

## ⚙️ Panduan Menjalankan Aplikasi

### 1. Masuk ke Direktori Backend
```bash
cd erp-be
```

### 2. Salin Konfigurasi Environment (Opsional)
File `.env` sudah disediakan secara otomatis dengan nilai default. Anda dapat menyesuaikannya:
```bash
cp .env.example .env
```

Isi default `.env`:
```env
PORT=8000
APP_ENV=development
APP_NAME="Kafeinarts ERP API"
DB_PATH=./storage/erp.db
JWT_SECRET=kafeinarts-super-secret-key-2026-very-secure
JWT_EXPIRATION_HOURS=72
CORS_ORIGINS=http://localhost:9000,http://localhost:9001,http://localhost:8080,http://127.0.0.1:9000
```

### 3. Jalankan Aplikasi
```bash
go run main.go
```

Server akan otomatis berjalan di **`http://localhost:8000`**.

Saat pertama kali dijalankan, sistem akan otomatis:
1. Membuat folder `storage/` dan file database `erp.db`.
2. Melakukan AutoMigrate untuk seluruh tabel model.
3. Menjalankan Seeder data demo (User admin, produk sampel, transaksi keuangan, pegawai HR).

---

## 🔐 Akun Default untuk Login

| Email | Password | Role | Keterangan |
|-------|----------|------|------------|
| `admin@kafeinarts.com` | `password123` | `superadmin` | Akses penuh seluruh modul |

---

## 🔗 Menghubungkan dengan Frontend (`erp-fe`)

Pastikan file `erp-fe/.env.development` mengarah ke URL backend:

```env
API_BASE_URL=http://localhost:8000/api/v1
APP_URL=http://localhost:8000
```

Jalankan frontend:
```bash
cd erp-fe
npm run dev
```

Buka browser di `http://localhost:9000` dan login menggunakan kredensial di atas.

---

## 🧪 Menjalankan Pengujian (Testing)

Untuk memverifikasi seluruh fungsi utilitas dan pengujian integrasi HTTP API end-to-end:

```bash
go test -v ./...
```

Hasil test mencakup:
- ✅ Unit test hashing & verifikasi password bcrypt
- ✅ Unit test pembuatan & validasi token JWT
- ✅ Unit test kalkulasi pagination metadata Quasar
- ✅ Integration test alur login, profil auth, produk inventaris, summary keuangan, dan staf HR.

---

## 📖 Dokumentasi Lengkap Endpoint API

Untuk daftar lengkap parameter query, payload request JSON, serta contoh response dari setiap endpoint, silakan buka:

👉 [Dokumentasi API Terperinci (`docs/API.md`)](file:///g:/Projek%20Koleb/erp-kafeinarts/erp-be/docs/API.md)

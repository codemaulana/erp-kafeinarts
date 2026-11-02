# ☕ KAFEINARTS ERP — Enterprise Monorepo

> Sistem ERP (*Enterprise Resource Planning*) komprehensif milik **Kafeinarts Studio & Cafe**, menggabungkan manajemen operasional kedai kopi (*coffee bar*), studio kreatif, rantai pasok (*supply chain*), keuangan & akuntansi terpadu, sumber daya manusia (*human resources*), dan kontrol hak akses pengguna (*role-based access control*).

![Go](https://img.shields.io/badge/Go-1.22%2B-00ADD8?logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-v3-00ACD7?logo=go&logoColor=white)
![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?logo=vue.js&logoColor=white)
![Quasar](https://img.shields.io/badge/Quasar-v2-1976D2?logo=quasar&logoColor=white)
![Pinia](https://img.shields.io/badge/Pinia-State_Management-FFD859?logo=vue.js&logoColor=black)
![SQLite](https://img.shields.io/badge/SQLite-Pure_Go-003B57?logo=sqlite&logoColor=white)
![Status](https://img.shields.io/badge/status-production--ready-green)

> 📖 **Dokumentasi Lengkap Alur Kode & Panduan Penggunaan:** Silakan baca panduan arsitektur sistem terperinci, alur eksekusi kode frontend & backend, siklus komunikasi REST API & JWT, serta petunjuk penggunaan di **[DOKUMENTASI_ALUR_DAN_PANDUAN.md](DOKUMENTASI_ALUR_DAN_PANDUAN.md)**.

---

## 📑 Daftar Isi

- [📖 Dokumentasi Alur & Panduan Lengkap (DOKUMENTASI_ALUR_DAN_PANDUAN.md)](DOKUMENTASI_ALUR_DAN_PANDUAN.md)
- [Tentang Sistem](#tentang-sistem)
- [Arsitektur & Tumpukan Teknologi](#arsitektur--tumpukan-teknologi)
- [Matriks Peran & Hak Akses (Roles & Permissions)](#matriks-peran--hak-akses-roles--permissions)
- [Modul & Fitur Lengkap Sistem](#modul--fitur-lengkap-sistem)
  - [1. Manajemen Pengguna (User Management)](#1-manajemen-pengguna-user-management)
  - [2. Supply Chain & Pengadaan (SC)](#2-supply-chain--pengadaan-sc)
  - [3. Keuangan & Akuntansi (Finance)](#3-keuangan--akuntansi-finance)
  - [4. Human Resource (HR)](#4-human-resource-hr)
  - [5. Inventaris & Produk (Inventory)](#5-inventaris--produk-inventory)
- [Panduan Instalasi & Menjalankan Aplikasi](#panduan-instalasi--menjalankan-aplikasi)
- [Kredensial Login Bawaan](#kredensial-login-bawaan)
- [Daftar Endpoint RESTful API](#daftar-endpoint-restful-api)
- [Troubleshooting & Solusi](#troubleshooting--solusi)
- [Lisensi](#lisensi)

---

## 🌟 Tentang Sistem

**Kafeinarts ERP** dirancang sebagai satu kesatuan sistem terintegrasi yang menyelaraskan seluruh alur kerja bisnis Kafeinarts Studio:
- Pembelian biji kopi & sirup operasional bar ke vendor.
- Alur *Purchase Request (PR)* hingga *Purchase Order (PO)* dan bukti fisik *Goods Receipt (Surat Jalan)*.
- Pengelolaan aset mesin kopi (*La Marzocco*, *EK43*) dan perangkat studio (*MacBook Pro M3*).
- Pembukuan akuntansi otomatis (*Accounts Payable, Accounts Receivable, Jurnal Kas/Bank, Laba Rugi*).
- Penggajian staf, pencatatan absensi, pengajuan cuti, serta evaluasi KPI bulanan.
- Pengaturan akun login staf dan pembagian wewenang yang ketat antar divisi.

---

## 🏗️ Arsitektur & Tumpukan Teknologi

Sistem menggunakan pendekatan **Monorepo** dengan pemisahan *clean architecture* antara API Backend dan UI Frontend:

```text
erp-kafeinarts/
├── erp-be/                     # BACKEND: Golang REST API (Port 8000)
│   ├── config/                 # Konfigurasi database & environment
│   ├── database/               # Seeder data awal (User, Produk, HR, SC, Finance)
│   ├── handler/                # HTTP Controller Fiber v3
│   ├── middleware/             # JWT Authentication, CORS, Logger
│   ├── model/                  # Struct GORM untuk seluruh tabel relasi
│   ├── repository/             # Data Access Layer (Query SQLite GORM)
│   ├── routes/                 # Pendaftaran rute API v1 modular
│   ├── service/                # Business Logic & DTO validation
│   ├── storage/                # File database SQLite (erp.db) & lampiran uploads
│   └── main.go                 # Entry point server Go
│
├── erp-fe/                     # FRONTEND: Vue 3 Single Page App (Port 9000)
│   ├── src/
│   │   ├── boot/               # Boot plugins: Axios (with interceptor), ACL
│   │   ├── components/         # Komponen reusable (AppDataTable, AppFormDialog, AppBadge)
│   │   ├── composables/        # Custom composable (useTablePagination, useNotify, useConfirmDelete)
│   │   ├── css/                # Tema Claymorphism, dark-mode tokens, styling SCSS
│   │   ├── layouts/            # Layout responsif (MainLayout dengan header & sidebar)
│   │   ├── modules/            # Modul domain (Inventory, Finance, HR, SC)
│   │   ├── pages/              # Halaman umum & Settings (Dashboard, UserManagementPage)
│   │   ├── router/             # Vue Router 4 dengan route guard autentikasi
│   │   ├── stores/             # Pinia stores (auth, user, inventory, finance, hr, sc)
│   │   └── utils/              # Formatter mata uang rupiah, tanggal Indo, validator
│   ├── quasar.config.js        # Konfigurasi Quasar Webpack
│   └── package.json
└── README.md                   # Dokumentasi resmi repositori
```

### Spesifikasi Stack

| Layer | Komponen | Detail Teknologi |
| :--- | :--- | :--- |
| **Backend** | Bahasa & Framework | Go 1.22+ dengan **Fiber v3** |
| | ORM & Database | **GORM** dengan **SQLite Pure Go** (`modernc.org/sqlite` / `glebarez/sqlite`, zero CGO dependency) |
| | Keamanan | **JWT** (`golang-jwt/jwt/v5`), Password Hash **BCrypt**, Middleware CORS & Sanitizer |
| **Frontend**| Framework Utama | **Vue 3** (Composition API, `<script setup>`) |
| | UI Toolkit | **Quasar Framework v2** dengan kustomisasi tema Glassmorphism / Claymorphism |
| | State Management | **Pinia** (reaktif, token terisolasi) |
| | HTTP Client | **Axios** dengan otomatisasi Bearer Token Injection & 401 Interceptor |

---

## 👥 Matriks Peran & Hak Akses (Roles & Permissions)

Sistem menerapkan **Role-Based Access Control (RBAC)** berbasis modul utama aplikasi untuk memastikan setiap divisi fokus pada alur kerjanya masing-masing.

> 💡 **Navigasi Sidebar Dinamis (Flat UI):** Tampilan menu navigasi pada *sidebar* akan beradaptasi secara otomatis mengikuti *role* pengguna yang sedang *login*. Sistem tidak menggunakan sub-menu bertingkat atau fitur *dropdown* (accordion); melainkan seluruh fitur spesifik dari suatu modul (seperti *Data Master*, *Data Transaksi*, *Laporan*) akan langsung tersusun rapi sebagai menu utama (*flat menu*) dengan penanda warna (*color-coded*) khusus untuk masing-masing modul. Hal ini dirancang agar pengguna dapat langsung mengakses menu tanpa perlu banyak klik.

| Role (Kode) | Label Peran | Akses Menu & Modul | Tugas & Hak Istimewa |
| :--- | :--- | :--- | :--- |
| **`superadmin`** | **Super Admin** | **Seluruh Modul** (`Dashboard`, `Inventaris`, `Finance`, `HR`, `Supply Chain`) + **Manajemen Pengguna** | • Kontrol penuh (*Full Access*) atas seluruh database & konfigurasi.<br>• Mengelola data pengguna di **Manajemen Pengguna** (`/app/settings/users`).<br>• Dilindungi aturan *anti-lockout* (tidak dapat menghapus akun sendiri). |
| **`inventory`** | **Inventaris** | `Dashboard`, `Inventaris` | • Mengelola katalog produk, bahan baku kopi, sirup, dan kemasan.<br>• Melakukan pencatatan penyesuaian stok (*stock opname*). |
| **`finance`** | **Finance** | `Dashboard`, `Finance` | • Pembukuan jurnal kas/bank & bagan akun (COA).<br>• Penerbitan Invoice Klien (Piutang) & Tagihan Vendor (Hutang).<br>• Akses laporan keuangan Laba Rugi, Neraca, dan Arus Kas. |
| **`hr`** | **Human Resource** | `Dashboard`, `Human Resource` | • Pengelolaan data master pegawai, departemen, & struktur jabatan.<br>• Rekapitulasi absensi/presensi, approval cuti, dan penilaian KPI.<br>• Proses penggajian (*payroll*) dan penerbitan slip gaji karyawan. |
| **`supply_chain`** | **Supply Chain** | `Dashboard`, `Supply Chain` | • Pengelolaan data vendor resmi (*supplier*) dan katalog barang pengadaan.<br>• Pengajuan *Purchase Request (PR)* & penerbitan *Purchase Order (PO)*.<br>• Input Surat Jalan / Penerimaan Barang (*Goods Receipt*) serta inventaris aset. |

---

## 📦 Modul & Fitur Lengkap Sistem

### 1. Manajemen Pengguna (User Management)
* **CRUD Pengguna via Web:** Super Admin dapat menambah pengguna baru langsung melalui antarmuka web di `/app/settings/users`.
* **Pilihan Role Dinamis:** Memberikan peran berbasis modul (*Super Admin*, *Inventaris*, *Finance*, *Human Resource*, *Supply Chain*) dengan hak akses menu otomatis.
* **Manajemen Kredensial:** Reset dan perubahan password pengguna langsung dari modal form terenkripsi BCrypt.
* **Filter & Pencarian Cepat:** Pencarian pengguna berdasarkan nama, alamat email, atau penyaringan berdasarkan role.

### 2. Supply Chain & Pengadaan (SC)
* **Master Vendor & Katalog:** Registrasi pemasok resmi (*supplier*), kategori item, dan katalog inventaris bahan/aset.
* **Lokasi & Gudang:** Manajemen penempatan barang per gudang (*Main Warehouse, Bar Counter, Studio*).
* **Purchase Request (PR):**
  - Form pengajuan barang internal dengan input departemen, prioritas (*LOW, MEDIUM, HIGH, URGENT*), dan grid dinamis item.
  - Otomatisasi nama pemohon (*Diajukan Oleh*) sesuai akun login aktif.
  - Alur persetujuan status (*DRAFT, SUBMITTED, APPROVED, REJECTED*).
* **Purchase Order (PO):**
  - Pembuatan PO resmi kepada vendor berdasarkan referensi PR.
  - Penghitungan otomatis Subtotal, PPN (Tax Amount), dan Total Tagihan.
* **Penerimaan Barang (Goods Receipt):** Pencatatan nomor Surat Jalan (*Delivery Note*), tanggal penerimaan fisik, dan verifikasi kuantitas/kondisi barang.
* **Manajemen Aset & Lisensi:**
  - Pencatatan aset mesin espresso (*La Marzocco*), grinder (*EK43*), dan laptop studio.
  - Pelacakan karyawan penanggung jawab (*Digunakan Oleh*), masa berlaku garansi, serta tanggal kedaluwarsa lisensi software.
* **Laporan Terpadu Supply Chain:**
  - Laporan PO per rentang tanggal & status.
  - Laporan Aset per Karyawan dan Departemen (lengkap dengan filter departemen/kategori).
  - Laporan Lisensi & Langganan Kedaluwarsa.
  - **Laporan Riwayat Transaksi Vendor:** Menampilkan riwayat transaksi per PO, vendor, tanggal, nominal Rupiah, nomor tagihan/surat jalan, dan status pembayaran (*Lunas / Belum Lunas / Sebagian*).

### 3. Keuangan & Akuntansi (Finance)
* **Bagan Akun (Chart of Accounts / COA):** Pengelompokan akun Aset, Liabilitas, Ekuitas, Pendapatan, dan Beban.
* **Jurnal Transaksi Kas & Bank:** Pencatatan kas masuk (*Debit*) dan kas keluar (*Credit*) dengan verifikasi saldo mutasi.
* **Piutang Usaha (Accounts Receivable / AR):** Penerbitan Faktur Penjualan/Jasa (*Client Invoice*) ke pelanggan studio serta fitur pembayaran bertahap (*Pay Invoice*).
* **Hutang Usaha (Accounts Payable / AP):** Pengelolaan Tagihan Vendor (*Vendor Bills*) atas pengadaan bahan dan perpanjangan sewa/lisensi.
* **Laporan Keuangan Otomatis:**
  - Laporan Laba Rugi (*Profit & Loss Statement*).
  - Laporan Neraca Saldo (*Balance Sheet*).
  - Laporan Arus Kas (*Cash Flow Report*).

### 4. Human Resource (HR)
* **Master Karyawan & Organisasi:** Struktur departemen, level jabatan, dan direktori lengkap profil data karyawan.
* **Presensi & Rekap Kehadiran:** Pencatatan jam masuk/pulang harian, status hadir/terlambat/izin, dan rekapitulasi bulanan.
* **Pengajuan Cuti & Saldo:** Pengajuan cuti tahunan, kuota saldo cuti per karyawan, dan alur persetujuan manajer.
* **Reimbursement Operasional:** Klaim pengembalian dana biaya darurat/lapangan disertai lampiran bukti kuitansi.
* **Evaluasi Kinerja (KPI):** Penilaian berkala performa kerja karyawan dengan skor dan catatan pembinaan.
* **Penggajian (Payroll) & Slip Gaji:** Otomatisasi kalkulasi gaji pokok, tunjangan, potongan absensi, serta penerbitan Slip Gaji digital.

### 5. Inventaris & Produk (Inventory)
* **Katalog Produk:** Manajemen stok barang jadi, kemasan, dan bahan baku minuman kedai kopi.
* **Penyesuaian Stok (Stock Adjustment):** Rekonsiliasi fisik (*stock opname*) dengan pencatatan selisih barang rusak/hilang.

---

## 🚀 Panduan Instalasi & Menjalankan Aplikasi

### Pilihan 1: Menggunakan Docker (Rekomendasi - Ringan & Sekali Perintah)

Menjalankan backend Golang dan frontend Quasar/Vue 3 secara bersamaan di background dalam satu kali perintah:

#### Prasyarat
- **Docker Desktop** sudah terpasang dan aktif di komputer Anda ([Unduh Docker Desktop](https://www.docker.com/products/docker-desktop/)).

#### Cara Menjalankan:
Cukup jalankan salah satu opsi berikut di folder root:

* **Opsi A (Klik Langsung di Windows):**
  Klik ganda file **`start-local.bat`** (browser otomatis terbuka ke aplikasi).

* **Opsi B (Via Terminal PowerShell / Command Prompt):**
  ```bash
  docker compose up -d
  ```

* **Akses Aplikasi:**
  * **Frontend Web:** [http://localhost:3008](http://localhost:3008)
  * **Backend API:** [http://localhost:8484/api/v1](http://localhost:8484/api/v1) (atau via reverse proxy di `http://localhost:3008/api/v1`)

* **Menghentikan Aplikasi:**
  Klik ganda `stop-local.bat` atau via terminal:
  ```bash
  docker compose down
  ```

---

### Pilihan 2: Menjalankan Secara Manual (Non-Docker)

#### Prasyarat
- **Go**: Versi 1.22 atau lebih baru ([Unduh Go](https://golang.org/dl/)).
- **Node.js**: Versi 18 LTS atau lebih baru ([Unduh Node.js](https://nodejs.org/)).
- **Git**: Terinstal di perangkat.

#### 1. Menjalankan Backend (`erp-be`)
Buka terminal baru di direktori backend:
```bash
cd erp-be
go mod tidy
go run main.go
```
* Backend aktif di **`http://localhost:8000`** (Base API: `http://localhost:8000/api/v1`).
* Database SQLite (`storage/erp.db`) beserta seluruh skema tabel dan seeder data awal akan **otomatis dibuat** pada saat pertama kali dijalankan.

#### 2. Menjalankan Frontend (`erp-fe`)
Buka terminal kedua di direktori frontend:
```bash
cd erp-fe
npm install
npm run dev
```
* Frontend SPA otomatis aktif di **`http://localhost:9000`**.

---

## 🔑 Kredensial Login Bawaan
 
| Akun | Email | Password | Role | Akses Menu & Modul |
| :--- | :--- | :--- | :--- | :--- |
| **Super Admin** | `admin@kafeinarts.com` | `password123` | `superadmin` | Semua Modul & Manajemen Pengguna |
| **Admin Inventaris** | `inventory@kafeinarts.com` | `password123` | `inventory` | Dashboard & Inventaris |
| **Admin Finance** | `finance@kafeinarts.com` | `password123` | `finance` | Dashboard & Finance |
| **Admin Human Resource** | `hr@kafeinarts.com` | `password123` | `hr` | Dashboard & Human Resource |
| **Admin Supply Chain** | `sc@kafeinarts.com` | `password123` | `supply_chain` | Dashboard & Supply Chain |

> ⚠️ **Catatan Keamanan:** Untuk implementasi server *production*, wajib mengganti password akun admin dan memperbarui nilai `JWT_SECRET` pada file `.env`.

---

## 📡 Daftar Endpoint RESTful API (Ringkasan)

Semua endpoint dilayani melalui basis URL: `/api/v1`

### 1. Autentikasi & Pengguna
- `POST /api/v1/auth/login` — Login pengguna & penerbitan token JWT.
- `GET  /api/v1/auth/me` — Ambil profil pengguna yang sedang login.
- `POST /api/v1/auth/logout` — Logout dan terminasi sesi.
- `GET    /api/v1/users` — Ambil daftar seluruh pengguna sistem *(query: search, role)*.
- `POST   /api/v1/users` — Tambah akun pengguna baru *(nama, email, password, role)*.
- `GET    /api/v1/users/:id` — Ambil detail satu pengguna.
- `PUT    /api/v1/users/:id` — Perbarui data pengguna / ubah role / ganti password.
- `DELETE /api/v1/users/:id` — Hapus akun pengguna (dilindungi anti-self-delete).

### 2. Supply Chain
- `GET / POST / PUT / DELETE /api/v1/sc/vendors` — CRUD Master Vendor/Supplier.
- `GET / POST / PUT / DELETE /api/v1/sc/items` — CRUD Master Katalog Item.
- `GET / POST / PUT / DELETE /api/v1/sc/purchase-requests` — Transaksi Purchase Request (PR).
- `GET / POST / PUT / DELETE /api/v1/sc/purchase-orders` — Transaksi Purchase Order (PO).
- `GET / POST / PUT / DELETE /api/v1/sc/goods-receipts` — Penerimaan Barang Fisik.
- `GET / POST / PUT / DELETE /api/v1/sc/assets` — Manajemen Aset & Inventaris Mesin.
- `GET /api/v1/sc/laporan/riwayat-vendor` — Laporan Riwayat Transaksi Vendor.
- `GET /api/v1/sc/laporan/aset` — Laporan Aset per Karyawan & Departemen.

### 3. Keuangan (Finance)
- `GET / POST /api/v1/finance/accounts` — Master Chart of Accounts (COA).
- `GET / POST /api/v1/finance/journals` — Jurnal Kas Masuk & Kas Keluar.
- `GET / POST /api/v1/finance/invoices` — Manajemen Invoice Piutang Klien.
- `GET / POST /api/v1/finance/vendor-bills` — Tagihan Hutang Vendor.
- `GET /api/v1/finance/reports/profit-loss` — Laporan Laba Rugi.
- `GET /api/v1/finance/reports/cash-flow` — Laporan Arus Kas.

### 4. Human Resource (HR)
- `GET / POST / PUT / DELETE /api/v1/hr/employees` — Master Data Karyawan.
- `GET / POST /api/v1/hr/attendances` — Presensi & Kehadiran Harian.
- `GET / POST /api/v1/hr/leave-requests` — Pengajuan & Saldo Cuti.
- `GET / POST /api/v1/hr/payrolls` — Proses Payroll & Slip Gaji Karyawan.

---

## 🛠️ Troubleshooting & Solusi

| Gejala Error | Kemungkinan Penyebab | Solusi |
| :--- | :--- | :--- |
| `bind: address already in use` | Port `8000` (backend) atau `9000` (frontend) sedang dipakai proses lain | Tutup proses yang berjalan atau gunakan `taskkill /F /IM main.exe` di terminal. |
| `TypeError: Cannot read properties of undefined` pada PR | `authStore` belum terinisialisasi pada template | Sudah diperbaiki; pastikan branch lokal Anda tersinkronisasi dengan commit terbaru. |
| Modal form terpotong di atas layar pada laptop beresolusi rendah | Dialog belum memiliki pembatas `max-height` flexbox | Sudah diperbaiki menggunakan `max-height: 90vh` dan container scroll responsif pada `AppFormDialog.vue`. |
| Perubahan kode Go tidak muncul di API | Terminal `go run main.go` belum di-restart setelah penambahan file baru | Tekan `Ctrl + C` pada terminal backend, lalu jalankan kembali `go run main.go`. |

---

## 📄 Lisensi

Hak Cipta © 2026 **Kafeinarts Studio**. Seluruh hak cipta dilindungi undang-undang.  
Sistem ERP ini bersifat **Internal & Proprietary** untuk keperluan operasional Kafeinarts Studio & Cafe.

# Dokumentasi REST API — Kafeinarts ERP

Base URL: `http://localhost:8000/api/v1`

Semua endpoint privat mewajibkan penyertaan header:
```http
Authorization: Bearer <JWT_TOKEN>
Accept: application/json
```

---

## 1. Modul Autentikasi (Auth)

### Login Pengguna
- **Method & Path:** `POST /auth/login` (atau alias `/login`)
- **Akses:** Publik
- **Request Body:**
  ```json
  {
    "email": "admin@kafeinarts.com",
    "password": "password123",
    "remember": false
  }
  ```
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Login berhasil",
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
    "user": {
      "id": 1,
      "name": "Super Admin",
      "email": "admin@kafeinarts.com",
      "role": "superadmin",
      "permissions": ["*"],
      "created_at": "2026-09-03T12:00:00Z",
      "updated_at": "2026-09-03T12:00:00Z"
    },
    "data": {
      "id": 1,
      "name": "Super Admin",
      "email": "admin@kafeinarts.com",
      "role": "superadmin",
      "permissions": ["*"]
    }
  }
  ```

### Profil Pengguna Aktif
- **Method & Path:** `GET /auth/me` (atau alias `/user`)
- **Akses:** Memerlukan Token (Bearer)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Profil user berhasil diambil",
    "data": {
      "id": 1,
      "name": "Super Admin",
      "email": "admin@kafeinarts.com",
      "role": "superadmin",
      "permissions": ["*"],
      "created_at": "2026-09-03T12:00:00Z",
      "updated_at": "2026-09-03T12:00:00Z"
    }
  }
  ```

### Logout
- **Method & Path:** `POST /auth/logout`
- **Akses:** Memerlukan Token (Bearer)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Logout berhasil"
  }
  ```

---

## 2. Modul Inventaris (Inventory)

### Daftar Produk
- **Method & Path:** `GET /inventory/products`
- **Akses:** Memerlukan Token (Bearer)
- **Query Parameters:**
  - `page`: Nomor halaman (integer, default `1`)
  - `per_page` / `perPage`: Jumlah data per halaman (integer, default `10`)
  - `search`: Filter teks berdasarkan nama atau SKU
  - `sortBy`: Nama kolom sorting (`name`, `sku`, `stock`, `price`)
  - `descending`: Nilai boolean (`true` / `false`)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Daftar produk berhasil dimuat",
    "data": [
      {
        "id": 1,
        "sku": "KAF-BEANS-001",
        "name": "House Blend Espresso Beans 1kg",
        "stock": 45,
        "price": 180000,
        "description": "Racikan arabika & robusta 70:30 khas Kafeinarts",
        "created_at": "2026-09-03T12:00:00Z",
        "updated_at": "2026-09-03T12:00:00Z"
      }
    ],
    "meta": {
      "total": 5,
      "page": 1,
      "current_page": 1,
      "per_page": 10,
      "total_pages": 1
    }
  }
  ```

### Detail Produk
- **Method & Path:** `GET /inventory/products/:id`
- **Akses:** Memerlukan Token (Bearer)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Detail produk berhasil dimuat",
    "data": {
      "id": 1,
      "sku": "KAF-BEANS-001",
      "name": "House Blend Espresso Beans 1kg",
      "stock": 45,
      "price": 180000,
      "description": "Racikan arabika & robusta 70:30 khas Kafeinarts",
      "created_at": "2026-09-03T12:00:00Z",
      "updated_at": "2026-09-03T12:00:00Z"
    }
  }
  ```

### Tambah Produk
- **Method & Path:** `POST /inventory/products`
- **Akses:** Memerlukan Token (Bearer)
- **Request Body:**
  ```json
  {
    "sku": "KAF-BEANS-003",
    "name": "Single Origin Mandheling 250g",
    "stock": 20,
    "price": 85000,
    "description": "Biji kopi proses wet hulled dengan aroma earthy & dark chocolate"
  }
  ```
- **Response 201 Created:**
  ```json
  {
    "success": true,
    "message": "Produk berhasil ditambahkan",
    "data": {
      "id": 6,
      "sku": "KAF-BEANS-003",
      "name": "Single Origin Mandheling 250g",
      "stock": 20,
      "price": 85000,
      "description": "Biji kopi proses wet hulled dengan aroma earthy & dark chocolate",
      "created_at": "2026-09-03T12:00:00Z",
      "updated_at": "2026-09-03T12:00:00Z"
    }
  }
  ```

### Update Produk
- **Method & Path:** `PUT /inventory/products/:id`
- **Akses:** Memerlukan Token (Bearer)
- **Request Body:**
  ```json
  {
    "sku": "KAF-BEANS-003",
    "name": "Single Origin Mandheling Premium 250g",
    "stock": 25,
    "price": 90000,
    "description": "Deskripsi diperbarui"
  }
  ```
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Produk berhasil diperbarui",
    "data": { ... }
  }
  ```

### Hapus Produk
- **Method & Path:** `DELETE /inventory/products/:id`
- **Akses:** Memerlukan Token (Bearer)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Produk berhasil dihapus",
    "data": null
  }
  ```

### Penyesuaian Stok (Stock Adjustment)
- **Method & Path:** `POST /inventory/products/:id/adjust`
- **Akses:** Memerlukan Token (Bearer)
- **Request Body:**
  ```json
  {
    "adjustment": 10,
    "type": "in",
    "reason": "Restock supplier mingguan"
  }
  ```
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Stok produk berhasil disesuaikan",
    "data": {
      "id": 1,
      "stock": 55,
      ...
    }
  }
  ```

---

## 3. Modul Keuangan (Finance) — Spesifikasi ERP Kafeinarts Software House

Dokumentasi lengkap modul keuangan selaras dengan dokumen acuan desain dan operasional sistem ERP Kafeinarts.

---

### A. DATA MASTER

Digunakan untuk mengelola data referensi statis yang digunakan pada form transaksi.

#### 1. Master Kategori Transaksi (`mst_finance_categories`)
- **Tujuan:** Mengelompokkan jenis pengeluaran/pemasukan (misal: `Project Revenue`, `Biaya Server`, `Gaji`, `Lisensi Software`).
- **Endpoint:**
  - `GET /finance/categories`: Ambil daftar kategori (Filter opsional: `?cash_type=INCOME` atau `?cash_type=EXPENSE`)
  - `POST /finance/categories`: Tambah kategori baru
  - `PUT /finance/categories/:id`: Perbarui kategori
  - `DELETE /finance/categories/:id`: Hapus kategori
- **Request Body (POST/PUT):**
  ```json
  {
    "category_name": "Project Revenue",
    "cash_type": "INCOME",
    "description": "Pendapatan termin proyek pembuatan software dan brand identity"
  }
  ```
- **Response 200/201:**
  ```json
  {
    "success": true,
    "message": "Kategori transaksi berhasil ditambahkan",
    "data": {
      "id": "cat-rev-01",
      "category_name": "Project Revenue",
      "cash_type": "INCOME",
      "description": "Pendapatan termin proyek...",
      "created_at": "2026-09-03T10:00:00Z"
    }
  }
  ```

#### 2. Master Rekening / Kas Bank (`mst_finance_accounts`)
- **Tujuan:** Penyimpanan uang perusahaan (misal: BCA, Mandiri, Petty Cash).
- **Endpoint:**
  - `GET /finance/accounts`: Ambil semua rekening kas/bank beserta saldo saat ini (`current_balance`)
  - `POST /finance/accounts`: Tambah rekening kas/bank
  - `PUT /finance/accounts/:id`: Perbarui data rekening
  - `DELETE /finance/accounts/:id`: Hapus rekening
- **Request Body (POST):**
  ```json
  {
    "bank_name": "BCA",
    "account_number": "8801234567",
    "account_holder_name": "PT Kafeinarts Studio",
    "current_balance": 50000000
  }
  ```

#### 3. Master Pajak (`mst_finance_taxes`)
- **Tujuan:** Konfigurasi persentase pajak terapan (misal: PPN 11%, PPh 21, PPh 23).
- **Endpoint:**
  - `GET /finance/taxes`: Ambil daftar tarif pajak
  - `POST /finance/taxes`: Tambah konfigurasi pajak
  - `PUT /finance/taxes/:id`: Perbarui tarif pajak
  - `DELETE /finance/taxes/:id`: Hapus tarif pajak
- **Request Body (POST):**
  ```json
  {
    "tax_name": "PPN 11%",
    "percentage": 11.0,
    "description": "Pajak Pertambahan Nilai 11%"
  }
  ```

---

### B. DATA TRANSAKSI

Area operasional harian tempat admin keuangan melakukan input data.

#### 1. Jurnal Kas (Pencatatan Harian) (`trx_cash_journals`)
- **Tujuan:** Mencatat mutasi kas masuk/keluar harian, otomatis mengupdate `current_balance` pada master rekening terkait.
- **Endpoint:**
  - `GET /finance/cash-journals`: Tabel Jurnal Kas harian (Query: `page`, `per_page`, `account_id`, `transaction_type` [INCOME/EXPENSE], `category_id`, `search`)
  - `POST /finance/cash-journals`: Catat transaksi kas harian
  - `GET /finance/cash-journals/:id`: Detail transaksi kas
  - `DELETE /finance/cash-journals/:id`: Hapus transaksi kas (otomatis mengembalikan saldo rekening)
- **Request Body (POST):**
  ```json
  {
    "transaction_date": "2026-09-03",
    "account_id": "acc-bca-01",
    "transaction_type": "INCOME",
    "category_id": "cat-rev-01",
    "amount": 16650000,
    "description": "Penerimaan pelunasan termin 1 web e-commerce",
    "attachment_url": "/uploads/struk_transfer_bca.jpg"
  }
  ```

#### 2. Manajemen Invoice / Piutang AR (`trx_invoices` & `trx_invoice_items`)
- **Tujuan:** Pembuatan dan pelacakan invoice tagihan ke klien berdasarkan termin proyek dengan item penagihan dinamis (*dynamic grid*).
- **Endpoint:**
  - `GET /finance/invoices`: Tabel Invoice (Query: `page`, `per_page`, `status` [DRAFT, SENT, PARTIAL, PAID, CANCELLED], `client_id`, `search`)
  - `POST /finance/invoices`: Buat invoice baru
  - `GET /finance/invoices/:id`: Detail invoice lengkap dengan rincian item & pajak
  - `PUT /finance/invoices/:id`: Perbarui invoice & item penagihan
  - `PATCH /finance/invoices/:id/status`: Update status invoice (`{"status": "PAID"}`)
  - `DELETE /finance/invoices/:id`: Hapus invoice
- **Request Body (POST):**
  ```json
  {
    "client_name": "PT Sinar Digital Nusantara",
    "project_name": "Website E-Commerce & Brand Kit",
    "issue_date": "2026-09-01",
    "due_date": "2026-09-15",
    "status": "SENT",
    "notes": "Pembayaran via transfer rekening BCA",
    "items": [
      {
        "item_description": "Desain UI/UX Design System (Figma) - 15 Pages",
        "quantity": 1,
        "unit_price": 6000000,
        "tax_id": "tax-ppn-11"
      },
      {
        "item_description": "Frontend & Backend Development (Termin 1)",
        "quantity": 1,
        "unit_price": 9000000,
        "tax_id": "tax-ppn-11"
      }
    ]
  }
  ```
- **Response 201 Created:**
  ```json
  {
    "success": true,
    "message": "Invoice berhasil diterbitkan",
    "data": {
      "id": "d1c2b3a4-5678-90ef-ghij-klmnopqrstuv",
      "invoice_number": "INV-202609-a1b2c3",
      "client_name": "PT Sinar Digital Nusantara",
      "project_name": "Website E-Commerce & Brand Kit",
      "issue_date": "2026-09-01T00:00:00Z",
      "due_date": "2026-09-15T00:00:00Z",
      "subtotal": 15000000,
      "total_tax": 1650000,
      "grand_total": 16650000,
      "status": "SENT",
      "items": [ ... ]
    }
  }
  ```

#### 3. Tagihan Vendor / Hutang AP (`trx_vendor_bills`)
- **Tujuan:** Pencatatan tagihan dari vendor pihak ketiga (misal: AWS, Google Cloud, lisensi JetBrains/Figma).
- **Endpoint:**
  - `GET /finance/vendor-bills`: Tabel Tagihan Vendor (Query: `page`, `per_page`, `status` [PENDING, PAID, OVERDUE], `vendor_id`, `search`)
  - `POST /finance/vendor-bills`: Catat tagihan vendor baru
  - `GET /finance/vendor-bills/:id`: Detail tagihan vendor
  - `PUT /finance/vendor-bills/:id`: Perbarui data tagihan
  - `PATCH /finance/vendor-bills/:id/status`: Update status tagihan (`{"status": "PAID"}`)
  - `DELETE /finance/vendor-bills/:id`: Hapus tagihan vendor
- **Request Body (POST):**
  ```json
  {
    "vendor_name": "Amazon Web Services (AWS)",
    "category_id": "cat-srv-01",
    "bill_date": "2026-09-01",
    "due_date": "2026-09-10",
    "amount": 2719500,
    "status": "PENDING",
    "attachment_url": "/uploads/aws_bill_august.pdf"
  }
  ```

#### 4. Upload Bukti Transaksi / Lampiran (`FileUploader`)
- **Method & Path:** `POST /finance/upload`
- **Content-Type:** `multipart/form-data` (Field name: `file`)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "File lampiran berhasil diunggah",
    "data": {
      "file_name": "struk_pembayaran.pdf",
      "file_url": "/uploads/20260903140000_a1b2c3d4.pdf",
      "size": 245100
    }
  }
  ```

---

### C. LAPORAN (REPORTS)

#### 1. Laporan Arus Kas (Cash Flow)
- **Method & Path:** `GET /finance/reports/cash-flow`
- **Filter:** `start_date` (YYYY-MM-DD), `end_date` (YYYY-MM-DD), `account_id` (opsional filter kas/bank tertentu)
- **Tabel Output:** Tanggal, No. Referensi, Keterangan, Kategori, Masuk, Keluar, Saldo Akhir.
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Laporan Arus Kas berhasil disusun",
    "data": {
      "start_date": "2026-09-01T00:00:00Z",
      "end_date": "2026-09-30T00:00:00Z",
      "beginning_cash": 75000000,
      "total_inflow": 16650000,
      "total_outflow": 3169500,
      "net_cash_flow": 13480500,
      "ending_cash": 88480500,
      "rows": [
        {
          "date": "2026-09-01T00:00:00Z",
          "reference_no": "TRX-IN-202609-001",
          "description": "Penerimaan pelunasan termin 1 web e-commerce",
          "category": "Project Revenue",
          "inflow": 16650000,
          "outflow": 0,
          "ending_balance": 91650000
        },
        {
          "date": "2026-09-02T00:00:00Z",
          "reference_no": "TRX-OUT-202609-001",
          "description": "Pembayaran server cloud AWS",
          "category": "Biaya Server",
          "inflow": 0,
          "outflow": 2719500,
          "ending_balance": 88930500
        }
      ]
    }
  }
  ```

#### 2. Laporan Laba/Rugi (Profit & Loss)
- **Method & Path:** `GET /finance/reports/profit-loss`
- **Filter:** `month` (1-12), `year` (contoh: `2026`)
- **Output:** Ringkasan bertingkat (Total Pendapatan, Total Pengeluaran, Laba/Rugi Bersih) dan breakdown kategori.
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Laporan Laba/Rugi berhasil disusun",
    "data": {
      "month": 9,
      "year": 2026,
      "total_revenue": 16650000,
      "total_expense": 3169500,
      "net_profit_loss": 13480500,
      "revenue_details": [
        { "category": "Project Revenue", "amount": 16650000 }
      ],
      "expense_details": [
        { "category": "Biaya Server", "amount": 2719500 },
        { "category": "Perlengkapan Studio", "amount": 450000 }
      ]
    }
  }
  ```

#### 3. Laporan Piutang (AR Aging)
- **Method & Path:** `GET /finance/reports/ar-aging`
- **Filter:** `status` (DRAFT, SENT, PARTIAL, OVERDUE), `client` (nama klien)
- **Tabel Output:** No. Invoice, Klien, Jatuh Tempo, Keterlambatan (Hari), Total Tagihan, Sisa Belum Dibayar.
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Laporan Piutang (AR Aging) berhasil disusun",
    "data": {
      "as_of_date": "2026-09-03T14:30:00Z",
      "total_receivable": 9435000,
      "total_overdue": 0,
      "total_invoice_count": 1,
      "items": [
        {
          "invoice_id": "inv-02-uuid",
          "invoice_number": "INV-202609-002",
          "client_name": "CV Kafein Kreatif Media",
          "project_name": "Mobile App POS & Loyalty Member",
          "issue_date": "2026-09-01T00:00:00Z",
          "due_date": "2026-09-15T00:00:00Z",
          "days_overdue": 0,
          "grand_total": 9435000,
          "remaining_due": 9435000,
          "status": "SENT"
        }
      ]
    }
  }
  ```

---

### D. RUTE ALIAS & KOMPATIBILITAS FRONTEND (QUASAR SPA)

Untuk mendukung antarmuka Frontend (Quasar SPA) yang telah dibuat tanpa memecah arsitektur backend, backend menyediakan alias rute bahasa Indonesia dan format JSON ganda (snake_case + camelCase):

| Fitur | Rute Standar | Alias Rute Frontend | Query Param Kompatibel |
|---|---|---|---|
| Master Kategori | `/finance/categories` | `/finance/kategori` | `page`, `perPage`, `search`, `tipeKas` |
| Master Rekening | `/finance/accounts` | `/finance/rekening` | `page`, `perPage`, `search` |
| Master Pajak | `/finance/taxes` | `/finance/pajak` | `page`, `perPage`, `search` |
| Jurnal Kas | `/finance/cash-journals` | `/finance/jurnal-kas`, `/finance/jurnal` | `page`, `perPage`, `tipe`, `akunBank`, `search` |
| Invoice (AR) | `/finance/invoices` | `/finance/invoice` | `page`, `perPage`, `status`, `klien`, `search` |
| Tagihan Vendor (AP) | `/finance/vendor-bills` | `/finance/tagihan-vendor`, `/finance/tagihan` | `page`, `perPage`, `status`, `vendor`, `search` |
| Laporan Arus Kas | `/finance/reports/cash-flow` | `/finance/laporan/cash-flow` | `periodeStart`, `periodeEnd`, `akunBank` |
| Laporan Laba Rugi | `/finance/reports/profit-loss` | `/finance/laporan/laba-rugi` | `bulan`, `tahun` |
| Laporan AR Aging | `/finance/reports/ar-aging` | `/finance/laporan/ar-aging` | `status`, `klien` |

*Seluruh endpoint di atas mendukung deserialisasi camelCase (`namaKategori`, `tipeKas`, `namaBank`, `nominal`, `keterangan`, `noInvoice`, `grandTotal`, dsb.) dan mengembalikan format `{ success: true, data: [...], meta: { total, page, perPage, lastPage } }`.*

---

## 4. Modul Sumber Daya Manusia (HR)

### Daftar Pegawai
- **Method & Path:** `GET /hr/employees`
- **Akses:** Memerlukan Token (Bearer)
- **Query Parameters:**
  - `page`: Nomor halaman (default `1`)
  - `per_page`: Jumlah data per halaman (default `15`)
  - `search`: Filter pencarian nama, NIP, atau email
  - `department`: Filter nama departemen
  - `status`: Filter status (`active`, `probation`, `inactive`)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Daftar pegawai berhasil dimuat",
    "data": [
      {
        "id": 1,
        "nip": "EMP-001",
        "name": "Hanif Maulana",
        "email": "hanif@kafeinarts.com",
        "phone": "081234567890",
        "department": "Studio & Creative",
        "position": "Creative Director",
        "status": "active",
        "join_date": "2023-05-01T00:00:00Z",
        "salary": 15000000,
        "created_at": "2026-09-03T12:00:00Z",
        "updated_at": "2026-09-03T12:00:00Z"
      }
    ],
    "meta": {
      "total": 3,
      "page": 1,
      "current_page": 1,
      "per_page": 15,
      "total_pages": 1
    }
  }
  ```

### Detail Pegawai
- **Method & Path:** `GET /hr/employees/:id`
- **Akses:** Memerlukan Token (Bearer)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Detail pegawai berhasil dimuat",
    "data": { ... }
  }
  ```

### Tambah Pegawai Baru
- **Method & Path:** `POST /hr/employees`
- **Akses:** Memerlukan Token (Bearer)
- **Request Body:**
  ```json
  {
    "nip": "EMP-004",
    "name": "Maya Putri",
    "email": "maya@kafeinarts.com",
    "phone": "081234567893",
    "department": "Studio & Creative",
    "position": "Graphic Designer",
    "status": "active",
    "join_date": "2024-09-01T00:00:00Z",
    "salary": 6500000
  }
  ```
- **Response 201 Created:**
  ```json
  {
    "success": true,
    "message": "Pegawai baru berhasil ditambahkan",
    "data": { ... }
  }
  ```

### Update Pegawai
- **Method & Path:** `PUT /hr/employees/:id`
- **Akses:** Memerlukan Token (Bearer)
- **Request Body:**
  ```json
  {
    "nip": "EMP-004",
    "name": "Maya Putri",
    "email": "maya@kafeinarts.com",
    "phone": "081234567893",
    "department": "Studio & Creative",
    "position": "Lead Graphic Designer",
    "status": "active",
    "join_date": "2024-09-01T00:00:00Z",
    "salary": 7500000
  }
  ```
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Data pegawai berhasil diperbarui",
    "data": { ... }
  }
  ```

### Hapus Pegawai
- **Method & Path:** `DELETE /hr/employees/:id`
- **Akses:** Memerlukan Token (Bearer)
- **Response 200 OK:**
  ```json
  {
    "success": true,
    "message": "Pegawai berhasil dihapus",
    "data": null
  }
  ```

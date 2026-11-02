/**
 * services/endpoints.js
 * Sentralisasi semua URL endpoint API.
 * Gunakan fungsi untuk endpoint yang membutuhkan parameter dinamis.
 */

// ── AUTH ────────────────────────────────────────────────────────
export const AUTH = {
  LOGIN: '/auth/login',
  LOGOUT: '/auth/logout',
  PROFILE: '/auth/me',
}

// ── FINANCE — Master ────────────────────────────────────────────
export const FINANCE = {
  // Kategori
  KATEGORI: '/finance/kategori',
  KATEGORI_BY_ID: (id) => `/finance/kategori/${id}`,

  // Rekening
  REKENING: '/finance/rekening',
  REKENING_BY_ID: (id) => `/finance/rekening/${id}`,

  // Pajak
  PAJAK: '/finance/pajak',
  PAJAK_BY_ID: (id) => `/finance/pajak/${id}`,

  // Jurnal Kas
  JURNAL: '/finance/jurnal-kas',
  JURNAL_BY_ID: (id) => `/finance/jurnal-kas/${id}`,

  // Invoice (AR)
  INVOICE: '/finance/invoices',
  INVOICE_BY_ID: (id) => `/finance/invoices/${id}`,
  INVOICE_STATUS: (id) => `/finance/invoices/${id}/status`,

  // Tagihan Vendor (AP)
  TAGIHAN: '/finance/vendor-bills',
  TAGIHAN_BY_ID: (id) => `/finance/vendor-bills/${id}`,
  TAGIHAN_STATUS: (id) => `/finance/vendor-bills/${id}/status`,

  // Upload
  UPLOAD: '/finance/upload',

  // Laporan
  LAPORAN_CASHFLOW: '/finance/laporan/cash-flow',
  LAPORAN_LABA_RUGI: '/finance/laporan/laba-rugi',
  LAPORAN_AR_AGING: '/finance/laporan/ar-aging',
}

// ── INVENTORY ────────────────────────────────────────────────────
export const INVENTORY = {
  PRODUCTS: '/inventory/products',
  PRODUCT_BY_ID: (id) => `/inventory/products/${id}`,
  STOCK_ADJUST: (id) => `/inventory/products/${id}/adjust`,
}

// ── HR ────────────────────────────────────────────────────────────
export const HR = {
  // Master — Departemen
  DEPARTMENTS: '/hr/departments',
  DEPARTMENT_BY_ID: (id) => `/hr/departments/${id}`,

  // Master — Jabatan
  POSITIONS: '/hr/positions',
  POSITION_BY_ID: (id) => `/hr/positions/${id}`,

  // Master — Karyawan
  EMPLOYEES: '/hr/employees',
  EMPLOYEE_BY_ID: (id) => `/hr/employees/${id}`,

  // Master — Jenis Cuti
  LEAVE_TYPES: '/hr/leave-types',
  LEAVE_TYPE_BY_ID: (id) => `/hr/leave-types/${id}`,

  // Master — Komponen Gaji
  SALARY_COMPONENTS: '/hr/salary-components',
  SALARY_COMPONENT_BY_ID: (id) => `/hr/salary-components/${id}`,

  // Master — Jadwal Kerja / Shift
  WORK_SHIFTS: '/hr/work-shifts',
  WORK_SHIFT_BY_ID: (id) => `/hr/work-shifts/${id}`,

  // Transaksi — Presensi
  ATTENDANCES: '/hr/attendances',
  ATTENDANCE_BY_ID: (id) => `/hr/attendances/${id}`,

  // Transaksi — Pengajuan Cuti
  LEAVE_REQUESTS: '/hr/leave-requests',
  LEAVE_REQUEST_BY_ID: (id) => `/hr/leave-requests/${id}`,
  LEAVE_REQUEST_STATUS: (id) => `/hr/leave-requests/${id}/status`,

  // Transaksi — Payroll
  PAYROLLS: '/hr/payrolls',
  PAYROLL_BY_ID: (id) => `/hr/payrolls/${id}`,
  PAYROLL_STATUS: (id) => `/hr/payrolls/${id}/status`,

  // Transaksi — Reimbursement
  REIMBURSEMENTS: '/hr/reimbursements',
  REIMBURSEMENT_BY_ID: (id) => `/hr/reimbursements/${id}`,
  REIMBURSEMENT_STATUS: (id) => `/hr/reimbursements/${id}/status`,

  // Transaksi — Penilaian Kinerja
  PERFORMANCE_REVIEWS: '/hr/performance-reviews',
  PERFORMANCE_REVIEW_BY_ID: (id) => `/hr/performance-reviews/${id}`,
  PERFORMANCE_REVIEW_STATUS: (id) => `/hr/performance-reviews/${id}/status`,

  // Upload
  UPLOAD: '/hr/upload',

  // Laporan
  LAPORAN_ABSENSI: '/hr/laporan/rekap-absensi',
  LAPORAN_PAYROLL: '/hr/laporan/slip-gaji',
  LAPORAN_SALDO_CUTI: '/hr/laporan/saldo-cuti',
  LAPORAN_TURNOVER: '/hr/laporan/turnover',
}

// ── SUPPLY CHAIN ────────────────────────────────────────────────────────────
export const SC = {
  // Master
  VENDORS: '/sc/vendors',
  VENDOR_BY_ID: (id) => `/sc/vendors/${id}`,

  ITEM_CATEGORIES: '/sc/item-categories',
  ITEM_CATEGORY_BY_ID: (id) => `/sc/item-categories/${id}`,

  ITEMS: '/sc/items',
  ITEM_BY_ID: (id) => `/sc/items/${id}`,

  LOCATIONS: '/sc/locations',
  LOCATION_BY_ID: (id) => `/sc/locations/${id}`,

  // Transaksi
  PURCHASE_REQUESTS: '/sc/purchase-requests',
  PURCHASE_REQUEST_BY_ID: (id) => `/sc/purchase-requests/${id}`,
  PURCHASE_REQUEST_STATUS: (id) => `/sc/purchase-requests/${id}/status`,

  PURCHASE_ORDERS: '/sc/purchase-orders',
  PURCHASE_ORDER_BY_ID: (id) => `/sc/purchase-orders/${id}`,
  PURCHASE_ORDER_STATUS: (id) => `/sc/purchase-orders/${id}/status`,

  GOODS_RECEIPTS: '/sc/goods-receipts',
  GOODS_RECEIPT_BY_ID: (id) => `/sc/goods-receipts/${id}`,

  ASSETS: '/sc/assets',
  ASSET_BY_ID: (id) => `/sc/assets/${id}`,
  ASSET_STATUS: (id) => `/sc/assets/${id}/status`,

  SUBSCRIPTIONS: '/sc/subscriptions',
  SUBSCRIPTION_BY_ID: (id) => `/sc/subscriptions/${id}`,
  SUBSCRIPTION_STATUS: (id) => `/sc/subscriptions/${id}/status`,

  UPLOAD: '/sc/upload',

  // Laporan
  LAPORAN_PO: '/sc/laporan/purchase-orders',
  LAPORAN_ASET: '/sc/laporan/aset',
  LAPORAN_LISENSI: '/sc/laporan/lisensi',
  LAPORAN_RIWAYAT_VENDOR: '/sc/laporan/riwayat-vendor',
}

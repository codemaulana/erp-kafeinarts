export default [
  // ── Data Master ────────────────────────────────────────────
  {
    path: '/sc/master/vendor',
    name: 'sc-master-vendor',
    component: () => import('modules/sc/views/master/VendorPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Master' },
        { label: 'Vendor / Supplier' },
      ],
    },
  },
  {
    path: '/sc/master/kategori',
    name: 'sc-master-kategori',
    component: () => import('modules/sc/views/master/KategoriItemPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Master' },
        { label: 'Kategori Item / Aset' },
      ],
    },
  },
  {
    path: '/sc/master/katalog',
    name: 'sc-master-katalog',
    component: () => import('modules/sc/views/master/KatalogItemPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Master' },
        { label: 'Katalog Item' },
      ],
    },
  },
  {
    path: '/sc/master/lokasi',
    name: 'sc-master-lokasi',
    component: () => import('modules/sc/views/master/LokasiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Master' },
        { label: 'Lokasi / Gudang' },
      ],
    },
  },

  // ── Data Transaksi ─────────────────────────────────────────
  {
    path: '/sc/transaksi/purchase-request',
    name: 'sc-transaksi-pr',
    component: () => import('modules/sc/views/transaksi/PurchaseRequestPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Transaksi' },
        { label: 'Purchase Request' },
      ],
    },
  },
  {
    path: '/sc/transaksi/purchase-order',
    name: 'sc-transaksi-po',
    component: () => import('modules/sc/views/transaksi/PurchaseOrderPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Transaksi' },
        { label: 'Purchase Order' },
      ],
    },
  },
  {
    path: '/sc/transaksi/goods-receipt',
    name: 'sc-transaksi-gr',
    component: () => import('modules/sc/views/transaksi/GoodsReceiptPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Transaksi' },
        { label: 'Penerimaan Barang' },
      ],
    },
  },
  {
    path: '/sc/transaksi/aset',
    name: 'sc-transaksi-aset',
    component: () => import('modules/sc/views/transaksi/ManajemenAsetPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Transaksi' },
        { label: 'Manajemen Aset' },
      ],
    },
  },
  {
    path: '/sc/transaksi/langganan',
    name: 'sc-transaksi-langganan',
    component: () => import('modules/sc/views/transaksi/PerpanjanganLanggananPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Data Transaksi' },
        { label: 'Perpanjangan Langganan' },
      ],
    },
  },

  // ── Laporan ────────────────────────────────────────────────
  {
    path: '/sc/laporan/purchase-order',
    name: 'sc-laporan-po',
    component: () => import('modules/sc/views/laporan/LaporanPurchaseOrderPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Laporan' },
        { label: 'Laporan Purchase Order' },
      ],
    },
  },
  {
    path: '/sc/laporan/aset',
    name: 'sc-laporan-aset',
    component: () => import('modules/sc/views/laporan/LaporanAsetPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Laporan' },
        { label: 'Laporan Aset' },
      ],
    },
  },
  {
    path: '/sc/laporan/lisensi',
    name: 'sc-laporan-lisensi',
    component: () => import('modules/sc/views/laporan/LaporanLisensiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Laporan' },
        { label: 'Lisensi & Langganan' },
      ],
    },
  },
  {
    path: '/sc/laporan/riwayat-vendor',
    name: 'sc-laporan-vendor',
    component: () => import('modules/sc/views/laporan/RiwayatVendorPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'inventory', label: 'Supply Chain' },
        { label: 'Laporan' },
        { label: 'Riwayat Transaksi Vendor' },
      ],
    },
  },
]

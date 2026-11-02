export default [
  // ── Data Master ────────────────────────────────────────────
  {
    path: '/finance/master/kategori',
    name: 'finance-master-kategori',
    component: () => import('modules/finance/views/master/KategoriPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Data Master' },
        { label: 'Kategori Transaksi' }
      ]
    }
  },
  {
    path: '/finance/master/rekening',
    name: 'finance-master-rekening',
    component: () => import('modules/finance/views/master/RekeningPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Data Master' },
        { label: 'Rekening / Kas Bank' }
      ]
    }
  },
  {
    path: '/finance/master/pajak',
    name: 'finance-master-pajak',
    component: () => import('modules/finance/views/master/PajakPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Data Master' },
        { label: 'Master Pajak' }
      ]
    }
  },

  // ── Data Transaksi ─────────────────────────────────────────
  {
    path: '/finance/jurnal-kas',
    name: 'finance-jurnal-kas',
    component: () => import('modules/finance/views/transaksi/JurnalKasPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Data Transaksi' },
        { label: 'Jurnal Kas' }
      ]
    }
  },
  {
    path: '/finance/invoice',
    name: 'finance-invoice',
    component: () => import('modules/finance/views/transaksi/InvoicePage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Data Transaksi' },
        { label: 'Invoice (AR)' }
      ]
    }
  },
  {
    path: '/finance/tagihan-vendor',
    name: 'finance-tagihan-vendor',
    component: () => import('modules/finance/views/transaksi/TagihanVendorPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Data Transaksi' },
        { label: 'Tagihan Vendor (AP)' }
      ]
    }
  },

  // ── Laporan ────────────────────────────────────────────────
  {
    path: '/finance/laporan/cash-flow',
    name: 'finance-laporan-cashflow',
    component: () => import('modules/finance/views/laporan/CashFlowPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Laporan' },
        { label: 'Arus Kas' }
      ]
    }
  },
  {
    path: '/finance/laporan/laba-rugi',
    name: 'finance-laporan-laba-rugi',
    component: () => import('modules/finance/views/laporan/LabaRugiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Laporan' },
        { label: 'Laba / Rugi' }
      ]
    }
  },
  {
    path: '/finance/laporan/ar-aging',
    name: 'finance-laporan-ar-aging',
    component: () => import('modules/finance/views/laporan/ArAgingPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'account_balance_wallet', label: 'Finance' },
        { label: 'Laporan' },
        { label: 'Piutang (AR Aging)' }
      ]
    }
  }
]

export default [
  // ── Data Master ────────────────────────────────────────────
  {
    path: '/hr/master/departemen',
    name: 'hr-master-departemen',
    component: () => import('modules/hr/views/master/DepartemenPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Master' },
        { label: 'Departemen' },
      ],
    },
  },
  {
    path: '/hr/master/jabatan',
    name: 'hr-master-jabatan',
    component: () => import('modules/hr/views/master/JabatanPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Master' },
        { label: 'Jabatan / Posisi' },
      ],
    },
  },
  {
    path: '/hr/master/karyawan',
    name: 'hr-master-karyawan',
    component: () => import('modules/hr/views/master/KaryawanPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Master' },
        { label: 'Data Karyawan' },
      ],
    },
  },
  {
    path: '/hr/master/jenis-cuti',
    name: 'hr-master-jenis-cuti',
    component: () => import('modules/hr/views/master/JenisCutiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Master' },
        { label: 'Jenis Cuti / Izin' },
      ],
    },
  },
  {
    path: '/hr/master/komponen-gaji',
    name: 'hr-master-komponen-gaji',
    component: () => import('modules/hr/views/master/KomponenGajiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Master' },
        { label: 'Komponen Gaji' },
      ],
    },
  },
  {
    path: '/hr/master/jadwal-kerja',
    name: 'hr-master-jadwal-kerja',
    component: () => import('modules/hr/views/master/JadwalKerjaPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Master' },
        { label: 'Jadwal Kerja / Shift' },
      ],
    },
  },

  // ── Data Transaksi ─────────────────────────────────────────
  {
    path: '/hr/presensi',
    name: 'hr-presensi',
    component: () => import('modules/hr/views/transaksi/PresensiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Transaksi' },
        { label: 'Presensi / Absensi' },
      ],
    },
  },
  {
    path: '/hr/pengajuan-cuti',
    name: 'hr-pengajuan-cuti',
    component: () => import('modules/hr/views/transaksi/PengajuanCutiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Transaksi' },
        { label: 'Pengajuan Cuti / Izin' },
      ],
    },
  },
  {
    path: '/hr/payroll',
    name: 'hr-payroll',
    component: () => import('modules/hr/views/transaksi/PayrollPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Transaksi' },
        { label: 'Penggajian (Payroll)' },
      ],
    },
  },
  {
    path: '/hr/reimbursement',
    name: 'hr-reimbursement',
    component: () => import('modules/hr/views/transaksi/ReimbursementPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Transaksi' },
        { label: 'Reimbursement' },
      ],
    },
  },
  {
    path: '/hr/penilaian-kinerja',
    name: 'hr-penilaian-kinerja',
    component: () => import('modules/hr/views/transaksi/PenilaianKinerjaPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Data Transaksi' },
        { label: 'Penilaian Kinerja' },
      ],
    },
  },

  // ── Laporan ────────────────────────────────────────────────
  {
    path: '/hr/laporan/rekap-absensi',
    name: 'hr-laporan-absensi',
    component: () => import('modules/hr/views/laporan/RekapAbsensiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Laporan' },
        { label: 'Rekap Absensi' },
      ],
    },
  },
  {
    path: '/hr/laporan/slip-gaji',
    name: 'hr-laporan-slip-gaji',
    component: () => import('modules/hr/views/laporan/SlipGajiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Laporan' },
        { label: 'Slip Gaji' },
      ],
    },
  },
  {
    path: '/hr/laporan/saldo-cuti',
    name: 'hr-laporan-saldo-cuti',
    component: () => import('modules/hr/views/laporan/SaldoCutiPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Laporan' },
        { label: 'Saldo Cuti' },
      ],
    },
  },
  {
    path: '/hr/laporan/turnover',
    name: 'hr-laporan-turnover',
    component: () => import('modules/hr/views/laporan/TurnoverPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'people_alt', label: 'HR' },
        { label: 'Laporan' },
        { label: 'Turnover Karyawan' },
      ],
    },
  },
]

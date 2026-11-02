/**
 * stores/hr.store.js
 * Pinia store for the HR domain.
 *
 * Follows the exact same pattern as finance.store.js:
 * - Single domain store with all sub-domains
 * - _loading object with isLoading(key) helper
 * - Computed dropdown options for q-select
 * - CRUD actions delegating to hr.api.js
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  departemenApi,
  jabatanApi,
  karyawanApi,
  jenisCutiApi,
  komponenGajiApi,
  jadwalKerjaApi,
  presensiApi,
  cutiApi,
  payrollApi,
  reimbursementApi,
  kinerjApi,
  hrLaporanApi,
} from 'src/modules/hr/services/hr.api'

// ── Shared default meta shape ──────────────────────────────────
const DEFAULT_META = Object.freeze({ total: 0, page: 1, perPage: 10, lastPage: 1 })

export const useHrStore = defineStore('hr', () => {
  // ── State — Master ─────────────────────────────────────────────
  const departemenList = ref([])
  const departemenMeta = ref({ ...DEFAULT_META })

  const jabatanList = ref([])
  const jabatanMeta = ref({ ...DEFAULT_META })

  const karyawanList = ref([])
  const karyawanMeta = ref({ ...DEFAULT_META })

  const jenisCutiList = ref([])
  const jenisCutiMeta = ref({ ...DEFAULT_META })

  const komponenGajiList = ref([])
  const komponenGajiMeta = ref({ ...DEFAULT_META })

  const jadwalKerjaList = ref([])
  const jadwalKerjaMeta = ref({ ...DEFAULT_META })

  // ── State — Transaksi ──────────────────────────────────────────
  const presensiList = ref([])
  const presensiMeta = ref({ ...DEFAULT_META })

  const cutiList = ref([])
  const cutiMeta = ref({ ...DEFAULT_META })

  const payrollList = ref([])
  const payrollMeta = ref({ ...DEFAULT_META })

  const reimbursementList = ref([])
  const reimbursementMeta = ref({ ...DEFAULT_META })

  const kinerjList = ref([])
  const kinerjMeta = ref({ ...DEFAULT_META })

  // ── State — Laporan ────────────────────────────────────────────
  const rekapAbsensiData = ref(null)
  const slipGajiData = ref(null)
  const saldoCutiData = ref(null)
  const turnoverData = ref(null)

  // ── Loading flags ──────────────────────────────────────────────
  const _loading = ref({
    departemen: false,
    jabatan: false,
    karyawan: false,
    jenisCuti: false,
    komponenGaji: false,
    jadwalKerja: false,
    presensi: false,
    cuti: false,
    payroll: false,
    reimbursement: false,
    kinerja: false,
    laporan: false,
  })

  function isLoading(key) {
    return _loading.value[key] === true
  }

  function _set(key, value) {
    _loading.value[key] = value
  }

  // ── Computed dropdown options ──────────────────────────────────
  const departemenOptions = computed(() =>
    departemenList.value.map((d) => ({
      label: d.department_name,
      value: d.id,
    })),
  )

  const jabatanOptions = computed(() =>
    jabatanList.value.map((j) => ({
      label: j.position_name,
      value: j.id,
      departmentId: j.department_id,
    })),
  )

  const karyawanOptions = computed(() =>
    karyawanList.value.map((k) => ({
      label: `${k.employee_number} - ${k.full_name}`,
      value: k.id,
      employeeNumber: k.employee_number,
      fullName: k.full_name,
    })),
  )

  const jenisCutiOptions = computed(() =>
    jenisCutiList.value.map((c) => ({
      label: c.leave_name,
      value: c.id,
      quotaDays: c.annual_quota_days,
      requiresAttachment: c.requires_attachment,
    })),
  )

  const komponenGajiOptions = computed(() =>
    komponenGajiList.value.map((kg) => ({
      label: kg.component_name,
      value: kg.id,
      type: kg.component_type,
      method: kg.calculation_method,
    })),
  )

  const jadwalKerjaOptions = computed(() =>
    jadwalKerjaList.value.map((jk) => ({
      label: jk.shift_name,
      value: jk.id,
    })),
  )

  // ── DEPARTEMEN ─────────────────────────────────────────────────
  async function fetchDepartemen(params = {}) {
    _set('departemen', true)
    try {
      const res = await departemenApi.getAll(params)
      departemenList.value = res.data
      departemenMeta.value = res.meta
    } finally {
      _set('departemen', false)
    }
  }
  async function createDepartemen(payload) {
    const result = await departemenApi.create(payload)
    await fetchDepartemen()
    return result
  }
  async function updateDepartemen(id, payload) {
    const result = await departemenApi.update(id, payload)
    await fetchDepartemen()
    return result
  }
  async function deleteDepartemen(id) {
    await departemenApi.delete(id)
    await fetchDepartemen()
  }

  // ── JABATAN ────────────────────────────────────────────────────
  async function fetchJabatan(params = {}) {
    _set('jabatan', true)
    try {
      const res = await jabatanApi.getAll(params)
      jabatanList.value = res.data
      jabatanMeta.value = res.meta
    } finally {
      _set('jabatan', false)
    }
  }
  async function createJabatan(payload) {
    const result = await jabatanApi.create(payload)
    await fetchJabatan()
    return result
  }
  async function updateJabatan(id, payload) {
    const result = await jabatanApi.update(id, payload)
    await fetchJabatan()
    return result
  }
  async function deleteJabatan(id) {
    await jabatanApi.delete(id)
    await fetchJabatan()
  }

  // ── KARYAWAN ───────────────────────────────────────────────────
  async function fetchKaryawan(params = {}) {
    _set('karyawan', true)
    try {
      const res = await karyawanApi.getAll(params)
      karyawanList.value = res.data
      karyawanMeta.value = res.meta
    } finally {
      _set('karyawan', false)
    }
  }
  async function createKaryawan(payload) {
    const result = await karyawanApi.create(payload)
    await fetchKaryawan()
    return result
  }
  async function updateKaryawan(id, payload) {
    const result = await karyawanApi.update(id, payload)
    await fetchKaryawan()
    return result
  }
  async function deleteKaryawan(id) {
    await karyawanApi.delete(id)
    await fetchKaryawan()
  }

  // ── JENIS CUTI ─────────────────────────────────────────────────
  async function fetchJenisCuti(params = {}) {
    _set('jenisCuti', true)
    try {
      const res = await jenisCutiApi.getAll(params)
      jenisCutiList.value = res.data
      jenisCutiMeta.value = res.meta
    } finally {
      _set('jenisCuti', false)
    }
  }
  async function createJenisCuti(payload) {
    const result = await jenisCutiApi.create(payload)
    await fetchJenisCuti()
    return result
  }
  async function updateJenisCuti(id, payload) {
    const result = await jenisCutiApi.update(id, payload)
    await fetchJenisCuti()
    return result
  }
  async function deleteJenisCuti(id) {
    await jenisCutiApi.delete(id)
    await fetchJenisCuti()
  }

  // ── KOMPONEN GAJI ──────────────────────────────────────────────
  async function fetchKomponenGaji(params = {}) {
    _set('komponenGaji', true)
    try {
      const res = await komponenGajiApi.getAll(params)
      komponenGajiList.value = res.data
      komponenGajiMeta.value = res.meta
    } finally {
      _set('komponenGaji', false)
    }
  }
  async function createKomponenGaji(payload) {
    const result = await komponenGajiApi.create(payload)
    await fetchKomponenGaji()
    return result
  }
  async function updateKomponenGaji(id, payload) {
    const result = await komponenGajiApi.update(id, payload)
    await fetchKomponenGaji()
    return result
  }
  async function deleteKomponenGaji(id) {
    await komponenGajiApi.delete(id)
    await fetchKomponenGaji()
  }

  // ── JADWAL KERJA ───────────────────────────────────────────────
  async function fetchJadwalKerja(params = {}) {
    _set('jadwalKerja', true)
    try {
      const res = await jadwalKerjaApi.getAll(params)
      jadwalKerjaList.value = res.data
      jadwalKerjaMeta.value = res.meta
    } finally {
      _set('jadwalKerja', false)
    }
  }
  async function createJadwalKerja(payload) {
    const result = await jadwalKerjaApi.create(payload)
    await fetchJadwalKerja()
    return result
  }
  async function updateJadwalKerja(id, payload) {
    const result = await jadwalKerjaApi.update(id, payload)
    await fetchJadwalKerja()
    return result
  }
  async function deleteJadwalKerja(id) {
    await jadwalKerjaApi.delete(id)
    await fetchJadwalKerja()
  }

  // ── PRESENSI ───────────────────────────────────────────────────
  async function fetchPresensi(params = {}) {
    _set('presensi', true)
    try {
      const res = await presensiApi.getAll(params)
      presensiList.value = res.data
      presensiMeta.value = res.meta
    } finally {
      _set('presensi', false)
    }
  }
  async function createPresensi(payload) {
    const result = await presensiApi.create(payload)
    await fetchPresensi()
    return result
  }
  async function updatePresensi(id, payload) {
    const result = await presensiApi.update(id, payload)
    await fetchPresensi()
    return result
  }
  async function deletePresensi(id) {
    await presensiApi.delete(id)
    await fetchPresensi()
  }

  // ── PENGAJUAN CUTI ─────────────────────────────────────────────
  async function fetchCuti(params = {}) {
    _set('cuti', true)
    try {
      const res = await cutiApi.getAll(params)
      cutiList.value = res.data
      cutiMeta.value = res.meta
    } finally {
      _set('cuti', false)
    }
  }
  async function createCuti(payload) {
    const result = await cutiApi.create(payload)
    await fetchCuti()
    return result
  }
  async function updateCuti(id, payload) {
    const result = await cutiApi.update(id, payload)
    await fetchCuti()
    return result
  }
  async function updateCutiStatus(id, status) {
    const result = await cutiApi.updateStatus(id, status)
    await fetchCuti()
    return result
  }
  async function deleteCuti(id) {
    await cutiApi.delete(id)
    await fetchCuti()
  }

  // ── PAYROLL ────────────────────────────────────────────────────
  async function fetchPayroll(params = {}) {
    _set('payroll', true)
    try {
      const res = await payrollApi.getAll(params)
      payrollList.value = res.data
      payrollMeta.value = res.meta
    } finally {
      _set('payroll', false)
    }
  }
  async function createPayroll(payload) {
    const result = await payrollApi.create(payload)
    await fetchPayroll()
    return result
  }
  async function updatePayroll(id, payload) {
    const result = await payrollApi.update(id, payload)
    await fetchPayroll()
    return result
  }
  async function updatePayrollStatus(id, status) {
    const result = await payrollApi.updateStatus(id, status)
    await fetchPayroll()
    return result
  }
  async function deletePayroll(id) {
    await payrollApi.delete(id)
    await fetchPayroll()
  }

  // ── REIMBURSEMENT ──────────────────────────────────────────────
  async function fetchReimbursement(params = {}) {
    _set('reimbursement', true)
    try {
      const res = await reimbursementApi.getAll(params)
      reimbursementList.value = res.data
      reimbursementMeta.value = res.meta
    } finally {
      _set('reimbursement', false)
    }
  }
  async function createReimbursement(payload) {
    const result = await reimbursementApi.create(payload)
    await fetchReimbursement()
    return result
  }
  async function updateReimbursement(id, payload) {
    const result = await reimbursementApi.update(id, payload)
    await fetchReimbursement()
    return result
  }
  async function updateReimbursementStatus(id, status) {
    const result = await reimbursementApi.updateStatus(id, status)
    await fetchReimbursement()
    return result
  }
  async function deleteReimbursement(id) {
    await reimbursementApi.delete(id)
    await fetchReimbursement()
  }

  // ── PENILAIAN KINERJA ──────────────────────────────────────────
  async function fetchKinerja(params = {}) {
    _set('kinerja', true)
    try {
      const res = await kinerjApi.getAll(params)
      kinerjList.value = res.data
      kinerjMeta.value = res.meta
    } finally {
      _set('kinerja', false)
    }
  }
  async function createKinerja(payload) {
    const result = await kinerjApi.create(payload)
    await fetchKinerja()
    return result
  }
  async function updateKinerja(id, payload) {
    const result = await kinerjApi.update(id, payload)
    await fetchKinerja()
    return result
  }
  async function updateKinerjaStatus(id, status) {
    const result = await kinerjApi.updateStatus(id, status)
    await fetchKinerja()
    return result
  }
  async function deleteKinerja(id) {
    await kinerjApi.delete(id)
    await fetchKinerja()
  }

  // ── LAPORAN ────────────────────────────────────────────────────
  async function fetchRekapAbsensi(params = {}) {
    _set('laporan', true)
    try {
      rekapAbsensiData.value = await hrLaporanApi.getRekapAbsensi(params)
    } finally {
      _set('laporan', false)
    }
  }
  async function fetchSlipGaji(params = {}) {
    _set('laporan', true)
    try {
      slipGajiData.value = await hrLaporanApi.getSlipGaji(params)
    } finally {
      _set('laporan', false)
    }
  }
  async function fetchSaldoCuti(params = {}) {
    _set('laporan', true)
    try {
      saldoCutiData.value = await hrLaporanApi.getSaldoCuti(params)
    } finally {
      _set('laporan', false)
    }
  }
  async function fetchTurnover(params = {}) {
    _set('laporan', true)
    try {
      turnoverData.value = await hrLaporanApi.getTurnover(params)
    } finally {
      _set('laporan', false)
    }
  }

  // ── Public API ─────────────────────────────────────────────────
  return {
    // ── State
    departemenList,
    departemenMeta,
    jabatanList,
    jabatanMeta,
    karyawanList,
    karyawanMeta,
    jenisCutiList,
    jenisCutiMeta,
    komponenGajiList,
    komponenGajiMeta,
    jadwalKerjaList,
    jadwalKerjaMeta,
    presensiList,
    presensiMeta,
    cutiList,
    cutiMeta,
    payrollList,
    payrollMeta,
    reimbursementList,
    reimbursementMeta,
    kinerjList,
    kinerjMeta,
    rekapAbsensiData,
    slipGajiData,
    saldoCutiData,
    turnoverData,

    // ── Computed dropdown options
    departemenOptions,
    jabatanOptions,
    karyawanOptions,
    jenisCutiOptions,
    komponenGajiOptions,
    jadwalKerjaOptions,

    // ── Loading helper
    isLoading,

    // ── Master CRUD
    fetchDepartemen,
    createDepartemen,
    updateDepartemen,
    deleteDepartemen,
    fetchJabatan,
    createJabatan,
    updateJabatan,
    deleteJabatan,
    fetchKaryawan,
    createKaryawan,
    updateKaryawan,
    deleteKaryawan,
    fetchJenisCuti,
    createJenisCuti,
    updateJenisCuti,
    deleteJenisCuti,
    fetchKomponenGaji,
    createKomponenGaji,
    updateKomponenGaji,
    deleteKomponenGaji,
    fetchJadwalKerja,
    createJadwalKerja,
    updateJadwalKerja,
    deleteJadwalKerja,

    // ── Transaksi CRUD
    fetchPresensi,
    createPresensi,
    updatePresensi,
    deletePresensi,
    fetchCuti,
    createCuti,
    updateCuti,
    updateCutiStatus,
    deleteCuti,
    fetchPayroll,
    createPayroll,
    updatePayroll,
    updatePayrollStatus,
    deletePayroll,
    fetchReimbursement,
    createReimbursement,
    updateReimbursement,
    updateReimbursementStatus,
    deleteReimbursement,
    fetchKinerja,
    createKinerja,
    updateKinerja,
    updateKinerjaStatus,
    deleteKinerja,

    // ── Laporan
    fetchRekapAbsensi,
    fetchSlipGaji,
    fetchSaldoCuti,
    fetchTurnover,
  }
})

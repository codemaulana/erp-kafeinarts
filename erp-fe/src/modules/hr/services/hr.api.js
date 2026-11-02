/**
 * HR API Service
 * Menghubungkan modul HR ke backend Fiber v3 (/api/v1/hr/...)
 * Mengikuti pattern finance.api.js:
 *   - Object literal per entitas
 *   - Normalisasi response { data, meta }
 *   - File upload via helper
 */

import { api } from 'boot/axios'

// ── Helper upload foto/lampiran ─────────────────────────────────
async function handleFileUpload(file) {
  if (!(file instanceof File)) return typeof file === 'string' ? file : ''
  const formData = new FormData()
  formData.append('file', file)
  const res = await api.post('/hr/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data?.data?.file_url || ''
}

// ── Shared default meta shape ───────────────────────────────────
const DEFAULT_META = { total: 0, page: 1, perPage: 10, lastPage: 1 }

function normalizeList(res) {
  return {
    data: res.data?.data || [],
    meta: res.data?.meta || { ...DEFAULT_META, total: (res.data?.data || []).length },
  }
}

// ─────────────────────────────────────────────────────────────────
//  MASTER — DEPARTEMEN
// ─────────────────────────────────────────────────────────────────
export const departemenApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/departments', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/departments/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/departments', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/departments/${id}`, payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/departments/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  MASTER — JABATAN / POSISI
// ─────────────────────────────────────────────────────────────────
export const jabatanApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/positions', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/positions/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/positions', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/positions/${id}`, payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/positions/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  MASTER — KARYAWAN
// ─────────────────────────────────────────────────────────────────
export const karyawanApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/employees', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/employees/${id}`)).data?.data || null
  },
  async create(payload) {
    const body = { ...payload }
    if (body.photoFile instanceof File) {
      body.photo_url = await handleFileUpload(body.photoFile)
      delete body.photoFile
    }
    return (await api.post('/hr/employees', body)).data?.data
  },
  async update(id, payload) {
    const body = { ...payload }
    if (body.photoFile instanceof File) {
      body.photo_url = await handleFileUpload(body.photoFile)
      delete body.photoFile
    }
    return (await api.put(`/hr/employees/${id}`, body)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/employees/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  MASTER — JENIS CUTI
// ─────────────────────────────────────────────────────────────────
export const jenisCutiApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/leave-types', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/leave-types/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/leave-types', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/leave-types/${id}`, payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/leave-types/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  MASTER — KOMPONEN GAJI
// ─────────────────────────────────────────────────────────────────
export const komponenGajiApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/salary-components', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/salary-components/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/salary-components', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/salary-components/${id}`, payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/salary-components/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  MASTER — JADWAL KERJA / SHIFT
// ─────────────────────────────────────────────────────────────────
export const jadwalKerjaApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/work-shifts', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/work-shifts/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/work-shifts', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/work-shifts/${id}`, payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/work-shifts/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  TRANSAKSI — PRESENSI / ABSENSI
// ─────────────────────────────────────────────────────────────────
export const presensiApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/attendances', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/attendances/${id}`)).data?.data || null
  },
  async create(payload) {
    const body = { ...payload }
    if (body.selfieFile instanceof File) {
      body.selfie_photo_url = await handleFileUpload(body.selfieFile)
      delete body.selfieFile
    }
    return (await api.post('/hr/attendances', body)).data?.data
  },
  async update(id, payload) {
    const body = { ...payload }
    if (body.selfieFile instanceof File) {
      body.selfie_photo_url = await handleFileUpload(body.selfieFile)
      delete body.selfieFile
    }
    return (await api.put(`/hr/attendances/${id}`, body)).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/attendances/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  TRANSAKSI — PENGAJUAN CUTI / IZIN
// ─────────────────────────────────────────────────────────────────
export const cutiApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/leave-requests', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/leave-requests/${id}`)).data?.data || null
  },
  async create(payload) {
    const body = { ...payload }
    if (body.attachmentFile instanceof File) {
      body.attachment_url = await handleFileUpload(body.attachmentFile)
      delete body.attachmentFile
    }
    return (await api.post('/hr/leave-requests', body)).data?.data
  },
  async update(id, payload) {
    const body = { ...payload }
    if (body.attachmentFile instanceof File) {
      body.attachment_url = await handleFileUpload(body.attachmentFile)
      delete body.attachmentFile
    }
    return (await api.put(`/hr/leave-requests/${id}`, body)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(`/hr/leave-requests/${id}/status`, { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/leave-requests/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  TRANSAKSI — PAYROLL
// ─────────────────────────────────────────────────────────────────
export const payrollApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/payrolls', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/payrolls/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/payrolls', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/payrolls/${id}`, payload)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(`/hr/payrolls/${id}/status`, { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/payrolls/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  TRANSAKSI — REIMBURSEMENT
// ─────────────────────────────────────────────────────────────────
export const reimbursementApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/reimbursements', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/reimbursements/${id}`)).data?.data || null
  },
  async create(payload) {
    const body = { ...payload }
    if (body.attachmentFile instanceof File) {
      body.attachment_url = await handleFileUpload(body.attachmentFile)
      delete body.attachmentFile
    }
    return (await api.post('/hr/reimbursements', body)).data?.data
  },
  async update(id, payload) {
    const body = { ...payload }
    if (body.attachmentFile instanceof File) {
      body.attachment_url = await handleFileUpload(body.attachmentFile)
      delete body.attachmentFile
    }
    return (await api.put(`/hr/reimbursements/${id}`, body)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(`/hr/reimbursements/${id}/status`, { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/reimbursements/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  TRANSAKSI — PENILAIAN KINERJA
// ─────────────────────────────────────────────────────────────────
export const kinerjApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get('/hr/performance-reviews', { params }))
  },
  async getById(id) {
    return (await api.get(`/hr/performance-reviews/${id}`)).data?.data || null
  },
  async create(payload) {
    return (await api.post('/hr/performance-reviews', payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(`/hr/performance-reviews/${id}`, payload)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(`/hr/performance-reviews/${id}/status`, { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(`/hr/performance-reviews/${id}`)).data
  },
}

// ─────────────────────────────────────────────────────────────────
//  LAPORAN
// ─────────────────────────────────────────────────────────────────
export const hrLaporanApi = {
  async getRekapAbsensi(params = {}) {
    return (await api.get('/hr/laporan/rekap-absensi', { params })).data?.data
  },
  async getSlipGaji(params = {}) {
    return (await api.get('/hr/laporan/slip-gaji', { params })).data?.data
  },
  async getSaldoCuti(params = {}) {
    return (await api.get('/hr/laporan/saldo-cuti', { params })).data?.data
  },
  async getTurnover(params = {}) {
    return (await api.get('/hr/laporan/turnover', { params })).data?.data
  },
}

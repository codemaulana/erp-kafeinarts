/**
 * Finance API Service
 * Menghubungkan modul Finance ke backend Fiber v3 (/api/v1/finance/...)
 */

import { api } from 'boot/axios'

// Helper untuk upload bukti/lampiran jika berupa File binary
async function handleFileUpload (file) {
  if (!(file instanceof File)) return typeof file === 'string' ? file : ''
  const formData = new FormData()
  formData.append('file', file)
  const res = await api.post('/finance/upload', formData, {
    headers: { 'Content-Type': 'multipart/form-data' }
  })
  return res.data?.data?.file_url || ''
}

// ─────────────────────────────────────────────────────────────
//  KATEGORI TRANSAKSI
// ─────────────────────────────────────────────────────────────
export const kategoriApi = {
  async getAll (params = {}) {
    const res = await api.get('/finance/kategori', { params })
    return {
      data: res.data?.data || [],
      meta: res.data?.meta || { total: (res.data?.data || []).length, page: 1, perPage: 10, lastPage: 1 }
    }
  },
  async getById (id) {
    const res = await api.get(`/finance/kategori/${id}`)
    return res.data?.data || null
  },
  async create (payload) {
    const res = await api.post('/finance/kategori', payload)
    return res.data?.data
  },
  async update (id, payload) {
    const res = await api.put(`/finance/kategori/${id}`, payload)
    return res.data?.data
  },
  async delete (id) {
    const res = await api.delete(`/finance/kategori/${id}`)
    return res.data
  }
}

// ─────────────────────────────────────────────────────────────
//  REKENING / KAS BANK
// ─────────────────────────────────────────────────────────────
export const rekeningApi = {
  async getAll (params = {}) {
    const res = await api.get('/finance/rekening', { params })
    return {
      data: res.data?.data || [],
      meta: res.data?.meta || { total: (res.data?.data || []).length, page: 1, perPage: 10, lastPage: 1 }
    }
  },
  async getById (id) {
    const res = await api.get(`/finance/rekening/${id}`)
    return res.data?.data || null
  },
  async create (payload) {
    const res = await api.post('/finance/rekening', payload)
    return res.data?.data
  },
  async update (id, payload) {
    const res = await api.put(`/finance/rekening/${id}`, payload)
    return res.data?.data
  },
  async delete (id) {
    const res = await api.delete(`/finance/rekening/${id}`)
    return res.data
  },
  async getAllRaw () {
    const res = await api.get('/finance/rekening')
    return res.data?.data || []
  }
}

// ─────────────────────────────────────────────────────────────
//  MASTER PAJAK
// ─────────────────────────────────────────────────────────────
export const pajakApi = {
  async getAll (params = {}) {
    const res = await api.get('/finance/pajak', { params })
    return {
      data: res.data?.data || [],
      meta: res.data?.meta || { total: (res.data?.data || []).length, page: 1, perPage: 10, lastPage: 1 }
    }
  },
  async getById (id) {
    const res = await api.get(`/finance/pajak/${id}`)
    return res.data?.data || null
  },
  async create (payload) {
    const res = await api.post('/finance/pajak', payload)
    return res.data?.data
  },
  async update (id, payload) {
    const res = await api.put(`/finance/pajak/${id}`, payload)
    return res.data?.data
  },
  async delete (id) {
    const res = await api.delete(`/finance/pajak/${id}`)
    return res.data
  },
  async getAllRaw () {
    const res = await api.get('/finance/pajak')
    return res.data?.data || []
  }
}

// ─────────────────────────────────────────────────────────────
//  JURNAL KAS
// ─────────────────────────────────────────────────────────────
export const jurnalApi = {
  async getAll (params = {}) {
    const res = await api.get('/finance/jurnal-kas', { params })
    return {
      data: res.data?.data || [],
      meta: res.data?.meta || { total: (res.data?.data || []).length, page: 1, perPage: 15, lastPage: 1 }
    }
  },
  async getById (id) {
    const res = await api.get(`/finance/jurnal-kas/${id}`)
    return res.data?.data || null
  },
  async create (payload) {
    const body = { ...payload }
    if (body.buktiFile instanceof File) {
      body.buktiFile = await handleFileUpload(body.buktiFile)
    }
    const res = await api.post('/finance/jurnal-kas', body)
    return res.data?.data
  },
  async update (id, payload) {
    const body = { ...payload }
    if (body.buktiFile instanceof File) {
      body.buktiFile = await handleFileUpload(body.buktiFile)
    }
    const res = await api.put(`/finance/jurnal-kas/${id}`, body)
    return res.data?.data
  },
  async delete (id) {
    const res = await api.delete(`/finance/jurnal-kas/${id}`)
    return res.data
  }
}

// ─────────────────────────────────────────────────────────────
//  INVOICE (AR)
// ─────────────────────────────────────────────────────────────
export const invoiceApi = {
  async getAll (params = {}) {
    const res = await api.get('/finance/invoices', { params })
    return {
      data: res.data?.data || [],
      meta: res.data?.meta || { total: (res.data?.data || []).length, page: 1, perPage: 10, lastPage: 1 }
    }
  },
  async getById (id) {
    const res = await api.get(`/finance/invoices/${id}`)
    return res.data?.data || null
  },
  async create (payload) {
    const res = await api.post('/finance/invoices', payload)
    return res.data?.data
  },
  async update (id, payload) {
    const res = await api.put(`/finance/invoices/${id}`, payload)
    return res.data?.data
  },
  async updateStatus (id, status) {
    const res = await api.patch(`/finance/invoices/${id}/status`, { status })
    return res.data?.data
  },
  async delete (id) {
    const res = await api.delete(`/finance/invoices/${id}`)
    return res.data
  }
}

// ─────────────────────────────────────────────────────────────
//  TAGIHAN VENDOR (AP)
// ─────────────────────────────────────────────────────────────
export const tagihanVendorApi = {
  async getAll (params = {}) {
    const res = await api.get('/finance/vendor-bills', { params })
    return {
      data: res.data?.data || [],
      meta: res.data?.meta || { total: (res.data?.data || []).length, page: 1, perPage: 10, lastPage: 1 }
    }
  },
  async getById (id) {
    const res = await api.get(`/finance/vendor-bills/${id}`)
    return res.data?.data || null
  },
  async create (payload) {
    const body = { ...payload }
    if (body.buktiFile instanceof File) {
      body.buktiFile = await handleFileUpload(body.buktiFile)
    }
    const res = await api.post('/finance/vendor-bills', body)
    return res.data?.data
  },
  async update (id, payload) {
    const body = { ...payload }
    if (body.buktiFile instanceof File) {
      body.buktiFile = await handleFileUpload(body.buktiFile)
    }
    const res = await api.put(`/finance/vendor-bills/${id}`, body)
    return res.data?.data
  },
  async updateStatus (id, status) {
    const res = await api.patch(`/finance/vendor-bills/${id}/status`, { status })
    return res.data?.data
  },
  async delete (id) {
    const res = await api.delete(`/finance/vendor-bills/${id}`)
    return res.data
  }
}

// ─────────────────────────────────────────────────────────────
//  LAPORAN
// ─────────────────────────────────────────────────────────────
export const laporanApi = {
  async getCashFlow (params = {}) {
    const res = await api.get('/finance/laporan/cash-flow', { params })
    return res.data?.data
  },
  async getLabaRugi (params = {}) {
    const res = await api.get('/finance/laporan/laba-rugi', { params })
    return res.data?.data
  },
  async getArAging (params = {}) {
    const res = await api.get('/finance/laporan/ar-aging', { params })
    return res.data?.data
  }
}

/**
 * SC API Service
 * Menghubungkan modul Supply Chain ke backend Fiber v3 (/api/v1/sc/...)
 * Mengikuti pattern hr.api.js dan finance.api.js
 */

import { api } from 'boot/axios'
import { SC } from 'src/services/endpoints'

async function handleFileUpload(file) {
  if (!(file instanceof File)) return typeof file === 'string' ? file : ''
  const formData = new FormData()
  formData.append('file', file)
  const res = await api.post(SC.UPLOAD, formData, {
    headers: { 'Content-Type': 'multipart/form-data' },
  })
  return res.data?.data?.file_url || ''
}

const DEFAULT_META = { total: 0, page: 1, perPage: 10, lastPage: 1 }

function normalizeList(res) {
  return {
    data: res.data?.data || [],
    meta: res.data?.meta || { ...DEFAULT_META, total: (res.data?.data || []).length },
  }
}

// ── MASTER ────────────────────────────────────────────────────────
export const vendorApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.VENDORS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.VENDOR_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.VENDORS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.VENDOR_BY_ID(id), payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.VENDOR_BY_ID(id))).data
  },
}

export const kategoriItemApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.ITEM_CATEGORIES, { params }))
  },
  async getById(id) {
    return (await api.get(SC.ITEM_CATEGORY_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.ITEM_CATEGORIES, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.ITEM_CATEGORY_BY_ID(id), payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.ITEM_CATEGORY_BY_ID(id))).data
  },
}

export const itemApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.ITEMS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.ITEM_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.ITEMS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.ITEM_BY_ID(id), payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.ITEM_BY_ID(id))).data
  },
}

export const lokasiApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.LOCATIONS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.LOCATION_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.LOCATIONS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.LOCATION_BY_ID(id), payload)).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.LOCATION_BY_ID(id))).data
  },
}

// ── TRANSAKSI ──────────────────────────────────────────────────────
export const purchaseRequestApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.PURCHASE_REQUESTS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.PURCHASE_REQUEST_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.PURCHASE_REQUESTS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.PURCHASE_REQUEST_BY_ID(id), payload)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(SC.PURCHASE_REQUEST_STATUS(id), { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.PURCHASE_REQUEST_BY_ID(id))).data
  },
}

export const purchaseOrderApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.PURCHASE_ORDERS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.PURCHASE_ORDER_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.PURCHASE_ORDERS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.PURCHASE_ORDER_BY_ID(id), payload)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(SC.PURCHASE_ORDER_STATUS(id), { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.PURCHASE_ORDER_BY_ID(id))).data
  },
}

export const goodsReceiptApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.GOODS_RECEIPTS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.GOODS_RECEIPT_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    const body = { ...payload }
    if (body.attachmentFile instanceof File) {
      body.attachment_url = await handleFileUpload(body.attachmentFile)
      delete body.attachmentFile
    }
    return (await api.post(SC.GOODS_RECEIPTS, body)).data?.data
  },
  async update(id, payload) {
    const body = { ...payload }
    if (body.attachmentFile instanceof File) {
      body.attachment_url = await handleFileUpload(body.attachmentFile)
      delete body.attachmentFile
    }
    return (await api.put(SC.GOODS_RECEIPT_BY_ID(id), body)).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.GOODS_RECEIPT_BY_ID(id))).data
  },
}

export const assetApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.ASSETS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.ASSET_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.ASSETS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.ASSET_BY_ID(id), payload)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(SC.ASSET_STATUS(id), { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.ASSET_BY_ID(id))).data
  },
}

export const subscriptionApi = {
  async getAll(params = {}) {
    return normalizeList(await api.get(SC.SUBSCRIPTIONS, { params }))
  },
  async getById(id) {
    return (await api.get(SC.SUBSCRIPTION_BY_ID(id))).data?.data || null
  },
  async create(payload) {
    return (await api.post(SC.SUBSCRIPTIONS, payload)).data?.data
  },
  async update(id, payload) {
    return (await api.put(SC.SUBSCRIPTION_BY_ID(id), payload)).data?.data
  },
  async updateStatus(id, status) {
    return (await api.patch(SC.SUBSCRIPTION_STATUS(id), { status })).data?.data
  },
  async delete(id) {
    return (await api.delete(SC.SUBSCRIPTION_BY_ID(id))).data
  },
}

// ── LAPORAN ────────────────────────────────────────────────────────
export const scLaporanApi = {
  async getLaporanPO(params = {}) {
    return (await api.get(SC.LAPORAN_PO, { params })).data?.data
  },
  async getLaporanAset(params = {}) {
    return (await api.get(SC.LAPORAN_ASET, { params })).data?.data
  },
  async getLaporanLisensi(params = {}) {
    return (await api.get(SC.LAPORAN_LISENSI, { params })).data?.data
  },
  async getRiwayatVendor(params = {}) {
    return (await api.get(SC.LAPORAN_RIWAYAT_VENDOR, { params })).data?.data
  },
}

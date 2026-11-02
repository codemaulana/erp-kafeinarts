/**
 * stores/sc.store.js
 * Pinia store for the Supply Chain domain.
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  vendorApi,
  kategoriItemApi,
  itemApi,
  lokasiApi,
  purchaseRequestApi,
  purchaseOrderApi,
  goodsReceiptApi,
  assetApi,
  subscriptionApi,
  scLaporanApi,
} from 'src/modules/sc/services/sc.api'

const DEFAULT_META = Object.freeze({ total: 0, page: 1, perPage: 10, lastPage: 1 })

export const useScStore = defineStore('sc', () => {
  // ── State — Master ─────────────────────────────────────────────
  const vendorList = ref([])
  const vendorMeta = ref({ ...DEFAULT_META })

  const kategoriItemList = ref([])
  const kategoriItemMeta = ref({ ...DEFAULT_META })

  const itemList = ref([])
  const itemMeta = ref({ ...DEFAULT_META })

  const lokasiList = ref([])
  const lokasiMeta = ref({ ...DEFAULT_META })

  // ── State — Transaksi ──────────────────────────────────────────
  const purchaseRequestList = ref([])
  const purchaseRequestMeta = ref({ ...DEFAULT_META })

  const purchaseOrderList = ref([])
  const purchaseOrderMeta = ref({ ...DEFAULT_META })

  const goodsReceiptList = ref([])
  const goodsReceiptMeta = ref({ ...DEFAULT_META })

  const assetList = ref([])
  const assetMeta = ref({ ...DEFAULT_META })

  const subscriptionList = ref([])
  const subscriptionMeta = ref({ ...DEFAULT_META })

  // ── State — Laporan ────────────────────────────────────────────
  const laporanPOData = ref(null)
  const laporanAsetData = ref(null)
  const laporanLisensiData = ref(null)
  const riwayatVendorData = ref(null)

  // ── Loading flags ──────────────────────────────────────────────
  const _loading = ref({
    vendor: false,
    kategoriItem: false,
    item: false,
    lokasi: false,
    pr: false,
    po: false,
    gr: false,
    asset: false,
    subscription: false,
    laporan: false,
  })

  function isLoading(key) {
    return _loading.value[key] === true
  }

  function _set(key, value) {
    _loading.value[key] = value
  }

  // ── Computed dropdown options ──────────────────────────────────
  const vendorOptions = computed(() =>
    vendorList.value.map((v) => ({
      label: v.vendor_name,
      value: v.id,
      category: v.vendor_category,
    })),
  )

  const kategoriItemOptions = computed(() =>
    kategoriItemList.value.map((k) => ({
      label: k.category_name,
      value: k.id,
      type: k.item_type,
    })),
  )

  const itemOptions = computed(() =>
    itemList.value.map((i) => ({
      label: i.item_name,
      value: i.id,
      unit: i.unit,
      price: i.estimated_price,
    })),
  )

  const lokasiOptions = computed(() =>
    lokasiList.value.map((l) => ({
      label: l.location_name,
      value: l.id,
    })),
  )

  const prOptions = computed(() =>
    purchaseRequestList.value.map((pr) => ({
      label: pr.pr_number,
      value: pr.id,
    })),
  )

  const poOptions = computed(() =>
    purchaseOrderList.value.map((po) => ({
      label: po.po_number,
      value: po.id,
    })),
  )

  // ── VENDOR ─────────────────────────────────────────────────────
  async function fetchVendor(params = {}) {
    _set('vendor', true)
    try {
      const res = await vendorApi.getAll(params)
      vendorList.value = res.data
      vendorMeta.value = res.meta
    } finally {
      _set('vendor', false)
    }
  }
  async function createVendor(payload) {
    const result = await vendorApi.create(payload)
    await fetchVendor()
    return result
  }
  async function updateVendor(id, payload) {
    const result = await vendorApi.update(id, payload)
    await fetchVendor()
    return result
  }
  async function deleteVendor(id) {
    await vendorApi.delete(id)
    await fetchVendor()
  }

  // ── KATEGORI ITEM ──────────────────────────────────────────────
  async function fetchKategoriItem(params = {}) {
    _set('kategoriItem', true)
    try {
      const res = await kategoriItemApi.getAll(params)
      kategoriItemList.value = res.data
      kategoriItemMeta.value = res.meta
    } finally {
      _set('kategoriItem', false)
    }
  }
  async function createKategoriItem(payload) {
    const result = await kategoriItemApi.create(payload)
    await fetchKategoriItem()
    return result
  }
  async function updateKategoriItem(id, payload) {
    const result = await kategoriItemApi.update(id, payload)
    await fetchKategoriItem()
    return result
  }
  async function deleteKategoriItem(id) {
    await kategoriItemApi.delete(id)
    await fetchKategoriItem()
  }

  // ── ITEM ───────────────────────────────────────────────────────
  async function fetchItem(params = {}) {
    _set('item', true)
    try {
      const res = await itemApi.getAll(params)
      itemList.value = res.data
      itemMeta.value = res.meta
    } finally {
      _set('item', false)
    }
  }
  async function createItem(payload) {
    const result = await itemApi.create(payload)
    await fetchItem()
    return result
  }
  async function updateItem(id, payload) {
    const result = await itemApi.update(id, payload)
    await fetchItem()
    return result
  }
  async function deleteItem(id) {
    await itemApi.delete(id)
    await fetchItem()
  }

  // ── LOKASI ─────────────────────────────────────────────────────
  async function fetchLokasi(params = {}) {
    _set('lokasi', true)
    try {
      const res = await lokasiApi.getAll(params)
      lokasiList.value = res.data
      lokasiMeta.value = res.meta
    } finally {
      _set('lokasi', false)
    }
  }
  async function createLokasi(payload) {
    const result = await lokasiApi.create(payload)
    await fetchLokasi()
    return result
  }
  async function updateLokasi(id, payload) {
    const result = await lokasiApi.update(id, payload)
    await fetchLokasi()
    return result
  }
  async function deleteLokasi(id) {
    await lokasiApi.delete(id)
    await fetchLokasi()
  }

  // ── PURCHASE REQUEST ───────────────────────────────────────────
  async function fetchPurchaseRequest(params = {}) {
    _set('pr', true)
    try {
      const res = await purchaseRequestApi.getAll(params)
      purchaseRequestList.value = res.data
      purchaseRequestMeta.value = res.meta
    } finally {
      _set('pr', false)
    }
  }
  async function createPurchaseRequest(payload) {
    const result = await purchaseRequestApi.create(payload)
    await fetchPurchaseRequest()
    return result
  }
  async function updatePurchaseRequest(id, payload) {
    const result = await purchaseRequestApi.update(id, payload)
    await fetchPurchaseRequest()
    return result
  }
  async function updatePurchaseRequestStatus(id, status) {
    const result = await purchaseRequestApi.updateStatus(id, status)
    await fetchPurchaseRequest()
    return result
  }
  async function deletePurchaseRequest(id) {
    await purchaseRequestApi.delete(id)
    await fetchPurchaseRequest()
  }

  // ── PURCHASE ORDER ─────────────────────────────────────────────
  async function fetchPurchaseOrder(params = {}) {
    _set('po', true)
    try {
      const res = await purchaseOrderApi.getAll(params)
      purchaseOrderList.value = res.data
      purchaseOrderMeta.value = res.meta
    } finally {
      _set('po', false)
    }
  }
  async function createPurchaseOrder(payload) {
    const result = await purchaseOrderApi.create(payload)
    await fetchPurchaseOrder()
    return result
  }
  async function updatePurchaseOrder(id, payload) {
    const result = await purchaseOrderApi.update(id, payload)
    await fetchPurchaseOrder()
    return result
  }
  async function updatePurchaseOrderStatus(id, status) {
    const result = await purchaseOrderApi.updateStatus(id, status)
    await fetchPurchaseOrder()
    return result
  }
  async function deletePurchaseOrder(id) {
    await purchaseOrderApi.delete(id)
    await fetchPurchaseOrder()
  }

  // ── GOODS RECEIPT ──────────────────────────────────────────────
  async function fetchGoodsReceipt(params = {}) {
    _set('gr', true)
    try {
      const res = await goodsReceiptApi.getAll(params)
      goodsReceiptList.value = res.data
      goodsReceiptMeta.value = res.meta
    } finally {
      _set('gr', false)
    }
  }
  async function createGoodsReceipt(payload) {
    const result = await goodsReceiptApi.create(payload)
    await fetchGoodsReceipt()
    return result
  }
  async function updateGoodsReceipt(id, payload) {
    const result = await goodsReceiptApi.update(id, payload)
    await fetchGoodsReceipt()
    return result
  }
  async function deleteGoodsReceipt(id) {
    await goodsReceiptApi.delete(id)
    await fetchGoodsReceipt()
  }

  // ── ASSET ──────────────────────────────────────────────────────
  async function fetchAsset(params = {}) {
    _set('asset', true)
    try {
      const res = await assetApi.getAll(params)
      assetList.value = res.data
      assetMeta.value = res.meta
    } finally {
      _set('asset', false)
    }
  }
  async function createAsset(payload) {
    const result = await assetApi.create(payload)
    await fetchAsset()
    return result
  }
  async function updateAsset(id, payload) {
    const result = await assetApi.update(id, payload)
    await fetchAsset()
    return result
  }
  async function updateAssetStatus(id, status) {
    const result = await assetApi.updateStatus(id, status)
    await fetchAsset()
    return result
  }
  async function deleteAsset(id) {
    await assetApi.delete(id)
    await fetchAsset()
  }

  // ── SUBSCRIPTION ───────────────────────────────────────────────
  async function fetchSubscription(params = {}) {
    _set('subscription', true)
    try {
      const res = await subscriptionApi.getAll(params)
      subscriptionList.value = res.data
      subscriptionMeta.value = res.meta
    } finally {
      _set('subscription', false)
    }
  }
  async function createSubscription(payload) {
    const result = await subscriptionApi.create(payload)
    await fetchSubscription()
    return result
  }
  async function updateSubscription(id, payload) {
    const result = await subscriptionApi.update(id, payload)
    await fetchSubscription()
    return result
  }
  async function updateSubscriptionStatus(id, status) {
    const result = await subscriptionApi.updateStatus(id, status)
    await fetchSubscription()
    return result
  }
  async function deleteSubscription(id) {
    await subscriptionApi.delete(id)
    await fetchSubscription()
  }

  // ── LAPORAN ────────────────────────────────────────────────────
  async function fetchLaporanPO(params = {}) {
    _set('laporan', true)
    try {
      laporanPOData.value = await scLaporanApi.getLaporanPO(params)
    } finally {
      _set('laporan', false)
    }
  }
  async function fetchLaporanAset(params = {}) {
    _set('laporan', true)
    try {
      laporanAsetData.value = await scLaporanApi.getLaporanAset(params)
    } finally {
      _set('laporan', false)
    }
  }
  async function fetchLaporanLisensi(params = {}) {
    _set('laporan', true)
    try {
      laporanLisensiData.value = await scLaporanApi.getLaporanLisensi(params)
    } finally {
      _set('laporan', false)
    }
  }
  async function fetchRiwayatVendor(params = {}) {
    _set('laporan', true)
    try {
      riwayatVendorData.value = await scLaporanApi.getRiwayatVendor(params)
    } finally {
      _set('laporan', false)
    }
  }

  return {
    vendorList,
    vendorMeta,
    kategoriItemList,
    kategoriItemMeta,
    itemList,
    itemMeta,
    lokasiList,
    lokasiMeta,
    purchaseRequestList,
    purchaseRequestMeta,
    purchaseOrderList,
    purchaseOrderMeta,
    goodsReceiptList,
    goodsReceiptMeta,
    assetList,
    assetMeta,
    subscriptionList,
    subscriptionMeta,
    laporanPOData,
    laporanAsetData,
    laporanLisensiData,
    riwayatVendorData,

    vendorOptions,
    kategoriItemOptions,
    itemOptions,
    lokasiOptions,
    prOptions,
    poOptions,
    isLoading,

    fetchVendor,
    createVendor,
    updateVendor,
    deleteVendor,
    fetchKategoriItem,
    createKategoriItem,
    updateKategoriItem,
    deleteKategoriItem,
    fetchItem,
    createItem,
    updateItem,
    deleteItem,
    fetchLokasi,
    createLokasi,
    updateLokasi,
    deleteLokasi,
    fetchPurchaseRequest,
    createPurchaseRequest,
    updatePurchaseRequest,
    updatePurchaseRequestStatus,
    deletePurchaseRequest,
    fetchPurchaseOrder,
    createPurchaseOrder,
    updatePurchaseOrder,
    updatePurchaseOrderStatus,
    deletePurchaseOrder,
    fetchGoodsReceipt,
    createGoodsReceipt,
    updateGoodsReceipt,
    deleteGoodsReceipt,
    fetchAsset,
    createAsset,
    updateAsset,
    updateAssetStatus,
    deleteAsset,
    fetchSubscription,
    createSubscription,
    updateSubscription,
    updateSubscriptionStatus,
    deleteSubscription,

    fetchLaporanPO,
    fetchLaporanAset,
    fetchLaporanLisensi,
    fetchRiwayatVendor,
  }
})

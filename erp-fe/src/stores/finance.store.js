/**
 * stores/finance.store.js
 * Pinia store for the Finance domain.
 *
 * Organises seven sub-domains in one store so related option-lists
 * (rekening, kategori, pajak) are always in sync across pages
 * without duplicating fetch calls.
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import {
  kategoriApi,
  rekeningApi,
  pajakApi,
  jurnalApi,
  invoiceApi,
  tagihanVendorApi,
  laporanApi
} from 'src/modules/finance/services/finance.api'

// ── Shared default meta shape ──────────────────────────────────
const DEFAULT_META = Object.freeze({ total: 0, page: 1, perPage: 10, lastPage: 1 })

export const useFinanceStore = defineStore('finance', () => {

  // ── State — Master ─────────────────────────────────────────────
  const kategoriList = ref([])
  const kategoriMeta = ref({ ...DEFAULT_META })

  const rekeningList = ref([])
  const rekeningMeta = ref({ ...DEFAULT_META })

  const pajakList = ref([])
  const pajakMeta = ref({ ...DEFAULT_META })

  // ── State — Transaksi ──────────────────────────────────────────
  const jurnalList = ref([])
  const jurnalMeta = ref({ ...DEFAULT_META })

  const invoiceList = ref([])
  const invoiceMeta = ref({ ...DEFAULT_META })

  const tagihanList = ref([])
  const tagihanMeta = ref({ ...DEFAULT_META })

  // ── State — Laporan ────────────────────────────────────────────
  const cashFlowData = ref(null)
  const labaRugiData = ref(null)
  const arAgingData  = ref(null)

  // ── Loading flags (internal — not exported raw) ───────────────
  const _loading = ref({
    kategori: false,
    rekening: false,
    pajak:    false,
    jurnal:   false,
    invoice:  false,
    tagihan:  false,
    laporan:  false
  })

  /**
   * Reactive loading check for a specific section.
   * Used in templates as: financeStore.isLoading('kategori')
   */
  function isLoading (key) {
    return _loading.value[key] === true
  }

  function _set (key, value) {
    _loading.value[key] = value
  }

  // ── Derived option-lists for <q-select> dropdowns ─────────────
  const rekeningOptions = computed(() =>
    rekeningList.value.map(r => ({
      label: `${r.namaBank} - ${r.noRekening}`,
      value: r.id
    }))
  )

  const pajakOptions = computed(() =>
    pajakList.value.map(p => ({
      label:      `${p.namaPajak} (${p.persentase}%)`,
      value:      p.id,
      persentase: p.persentase
    }))
  )

  const kategoriOptions = computed(() =>
    kategoriList.value.map(k => ({
      label:   k.namaKategori,
      value:   k.id,
      tipeKas: k.tipeKas
    }))
  )

  // ── KATEGORI ───────────────────────────────────────────────────
  async function fetchKategori (params = {}) {
    _set('kategori', true)
    try {
      const res = await kategoriApi.getAll(params)
      kategoriList.value = res.data
      kategoriMeta.value = res.meta
    } finally {
      _set('kategori', false)
    }
  }

  async function createKategori (payload) {
    const result = await kategoriApi.create(payload)
    await fetchKategori()
    return result
  }

  async function updateKategori (id, payload) {
    const result = await kategoriApi.update(id, payload)
    await fetchKategori()
    return result
  }

  async function deleteKategori (id) {
    await kategoriApi.delete(id)
    await fetchKategori()
  }

  // ── REKENING ───────────────────────────────────────────────────
  async function fetchRekening (params = {}) {
    _set('rekening', true)
    try {
      const res = await rekeningApi.getAll(params)
      rekeningList.value = res.data
      rekeningMeta.value = res.meta
    } finally {
      _set('rekening', false)
    }
  }

  async function createRekening (payload) {
    const result = await rekeningApi.create(payload)
    await fetchRekening()
    return result
  }

  async function updateRekening (id, payload) {
    const result = await rekeningApi.update(id, payload)
    await fetchRekening()
    return result
  }

  async function deleteRekening (id) {
    await rekeningApi.delete(id)
    await fetchRekening()
  }

  // ── PAJAK ──────────────────────────────────────────────────────
  async function fetchPajak (params = {}) {
    _set('pajak', true)
    try {
      const res = await pajakApi.getAll(params)
      pajakList.value = res.data
      pajakMeta.value = res.meta
    } finally {
      _set('pajak', false)
    }
  }

  async function createPajak (payload) {
    const result = await pajakApi.create(payload)
    await fetchPajak()
    return result
  }

  async function updatePajak (id, payload) {
    const result = await pajakApi.update(id, payload)
    await fetchPajak()
    return result
  }

  async function deletePajak (id) {
    await pajakApi.delete(id)
    await fetchPajak()
  }

  // ── JURNAL KAS ─────────────────────────────────────────────────
  async function fetchJurnal (params = {}) {
    _set('jurnal', true)
    try {
      const res = await jurnalApi.getAll(params)
      jurnalList.value = res.data
      jurnalMeta.value = res.meta
    } finally {
      _set('jurnal', false)
    }
  }

  async function createJurnal (payload) {
    const result = await jurnalApi.create(payload)
    await fetchJurnal()
    return result
  }

  async function updateJurnal (id, payload) {
    const result = await jurnalApi.update(id, payload)
    await fetchJurnal()
    return result
  }

  async function deleteJurnal (id) {
    await jurnalApi.delete(id)
    await fetchJurnal()
  }

  // ── INVOICE (AR) ───────────────────────────────────────────────
  async function fetchInvoice (params = {}) {
    _set('invoice', true)
    try {
      const res = await invoiceApi.getAll(params)
      invoiceList.value = res.data
      invoiceMeta.value = res.meta
    } finally {
      _set('invoice', false)
    }
  }

  async function createInvoice (payload) {
    const result = await invoiceApi.create(payload)
    await fetchInvoice()
    return result
  }

  async function updateInvoice (id, payload) {
    const result = await invoiceApi.update(id, payload)
    await fetchInvoice()
    return result
  }

  async function deleteInvoice (id) {
    await invoiceApi.delete(id)
    await fetchInvoice()
  }

  // ── TAGIHAN VENDOR (AP) ────────────────────────────────────────
  async function fetchTagihan (params = {}) {
    _set('tagihan', true)
    try {
      const res = await tagihanVendorApi.getAll(params)
      tagihanList.value = res.data
      tagihanMeta.value = res.meta
    } finally {
      _set('tagihan', false)
    }
  }

  async function createTagihan (payload) {
    const result = await tagihanVendorApi.create(payload)
    await fetchTagihan()
    return result
  }

  async function updateTagihan (id, payload) {
    const result = await tagihanVendorApi.update(id, payload)
    await fetchTagihan()
    return result
  }

  async function deleteTagihan (id) {
    await tagihanVendorApi.delete(id)
    await fetchTagihan()
  }

  // ── LAPORAN ────────────────────────────────────────────────────
  async function fetchCashFlow (params = {}) {
    _set('laporan', true)
    try {
      cashFlowData.value = await laporanApi.getCashFlow(params)
    } finally {
      _set('laporan', false)
    }
  }

  async function fetchLabaRugi (params = {}) {
    _set('laporan', true)
    try {
      labaRugiData.value = await laporanApi.getLabaRugi(params)
    } finally {
      _set('laporan', false)
    }
  }

  async function fetchArAging (params = {}) {
    _set('laporan', true)
    try {
      arAgingData.value = await laporanApi.getArAging(params)
    } finally {
      _set('laporan', false)
    }
  }

  // ── Public API ─────────────────────────────────────────────────
  return {
    // ── State
    kategoriList, kategoriMeta,
    rekeningList, rekeningMeta,
    pajakList,    pajakMeta,
    jurnalList,   jurnalMeta,
    invoiceList,  invoiceMeta,
    tagihanList,  tagihanMeta,
    cashFlowData, labaRugiData, arAgingData,

    // ── Computed dropdown options
    rekeningOptions,
    pajakOptions,
    kategoriOptions,

    // ── Loading helper
    isLoading,

    // ── Kategori
    fetchKategori, createKategori, updateKategori, deleteKategori,

    // ── Rekening
    fetchRekening, createRekening, updateRekening, deleteRekening,

    // ── Pajak
    fetchPajak, createPajak, updatePajak, deletePajak,

    // ── Jurnal Kas
    fetchJurnal, createJurnal, updateJurnal, deleteJurnal,

    // ── Invoice (AR)
    fetchInvoice, createInvoice, updateInvoice, deleteInvoice,

    // ── Tagihan Vendor (AP)
    fetchTagihan, createTagihan, updateTagihan, deleteTagihan,

    // ── Laporan
    fetchCashFlow, fetchLabaRugi, fetchArAging
  }
})

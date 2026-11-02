/**
 * composables/useTablePagination.js
 * State pagination QTable yang reusable.
 * Menghilangkan duplikasi qPagination di setiap halaman tabel.
 */
import { ref, unref } from 'vue'

/**
 * @param {object|function} [initial]
 * @param {number} [initial.page=1]
 * @param {number} [initial.rowsPerPage=10]
 * @param {number} [initial.rowsNumber=0]
 * @param {function} [initial.fetchFn]
 * @param {function} [initial.getParams]
 * @param {import('vue').Ref<object>} [initial.meta]
 */
export function useTablePagination (initial = {}) {
  let fetchCallback = null
  let getParamsCallback = null
  let metaRef = null

  if (typeof initial === 'function') {
    fetchCallback = initial
  } else if (initial && typeof initial === 'object') {
    const fn = initial.fetchFn || initial.fetchData
    if (typeof fn === 'function') {
      fetchCallback = fn
    }
    if (typeof initial.getParams === 'function') {
      getParamsCallback = initial.getParams
    }
    if (initial.meta) {
      metaRef = initial.meta
    }
  }

  const paginationRaw = ref({
    page:        (typeof initial === 'object' && initial.page)        ?? 1,
    rowsPerPage: (typeof initial === 'object' && initial.rowsPerPage) ?? 10,
    rowsNumber:  (typeof initial === 'object' && initial.rowsNumber)  ?? 0,
    sortBy:      (typeof initial === 'object' && (initial.sortBy || initial.defaultSort)) ?? null,
    descending:  (typeof initial === 'object' && initial.descending)  ?? false
  })

  // Proxy wrapper agar BISA diakses baik sebagai pagination.value.rowsPerPage maupun pagination.rowsPerPage
  const pagination = new Proxy(paginationRaw, {
    get (target, prop, receiver) {
      if (prop in target || prop === '__v_isRef') {
        return Reflect.get(target, prop, receiver)
      }
      return target.value ? target.value[prop] : undefined
    },
    set (target, prop, value, receiver) {
      if (prop in target) {
        return Reflect.set(target, prop, value, receiver)
      }
      if (target.value) {
        target.value[prop] = value
        return true
      }
      return false
    }
  })

  const searchQuery = ref('')

  /**
   * Hitung nomor urut baris (1-based) berdasarkan halaman & rowsPerPage.
   * @param {number} rowIndex  0-based index dalam halaman saat ini
   * @returns {number}
   */
  function rowNo (rowIndex) {
    return (pagination.page - 1) * pagination.rowsPerPage + rowIndex + 1
  }

  /**
   * Sync total rows dari meta response ke pagination.rowsNumber.
   * @param {number} total
   */
  function syncTotal (total) {
    pagination.rowsNumber = Number(total) || 0
  }

  /** Reset ke halaman 1. */
  function resetPage () {
    pagination.page = 1
  }

  /**
   * Muat data halaman saat ini / halaman tertentu.
   * @param {number} [page=pagination.page]
   */
  async function load (page = pagination.page) {
    pagination.page = page
    const customParams = getParamsCallback ? getParamsCallback() : {}
    const params = {
      page: pagination.page,
      perPage: pagination.rowsPerPage,
      limit: pagination.rowsPerPage,
      sortBy: pagination.sortBy,
      sortOrder: pagination.descending ? 'desc' : 'asc',
      ...customParams
    }

    if (fetchCallback) {
      let result
      if (typeof initial === 'function' && initial.length > 1) {
        result = await fetchCallback(pagination.page, pagination.rowsPerPage, pagination.sortBy, pagination.descending)
      } else {
        result = await fetchCallback(params)
      }

      const currentMeta = unref(metaRef)
      if (currentMeta?.total != null) {
        syncTotal(currentMeta.total)
      } else if (result?.meta?.total != null) {
        syncTotal(result.meta.total)
      }
    }
  }

  /**
   * Handler untuk event @request dari q-table.
   * Memperbarui state pagination dan memanggil fetchFn dengan parameter baru.
   *
   * @param {{ pagination: { page, rowsPerPage, sortBy, descending } }} props  Event payload dari q-table
   * @param {function} [fetchFn]  Fungsi async opsional yang menerima { page, rowsPerPage }
   */
  async function onRequest (props, fetchFn) {
    if (props && props.pagination) {
      pagination.page        = props.pagination.page
      pagination.rowsPerPage = props.pagination.rowsPerPage
      pagination.sortBy      = props.pagination.sortBy
      pagination.descending  = props.pagination.descending
    }

    const effectiveFetch = fetchFn || fetchCallback
    if (typeof effectiveFetch === 'function') {
      if (fetchFn) {
        await fetchFn({ page: pagination.page, perPage: pagination.rowsPerPage })
      } else {
        await load(pagination.page)
      }
    }

    const currentMeta = unref(metaRef)
    if (currentMeta?.total != null) {
      syncTotal(currentMeta.total)
    }
  }

  return {
    pagination,
    searchQuery,
    rowNo,
    onRequest,
    syncTotal,
    resetPage,
    load
  }
}


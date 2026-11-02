/**
 * stores/inventory.store.js
 * Pinia store for the Inventory domain.
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { inventoryApi } from 'src/modules/inventory/services/inventory.api'

export const useInventoryStore = defineStore('inventory', () => {

  // ── State ───────────────────────────────────────────────────────
  const products  = ref([])
  const isLoading = ref(false)
  const meta      = ref({ total: 0 })

  // ── Actions ─────────────────────────────────────────────────────
  async function fetchProducts (params = {}) {
    isLoading.value = true
    try {
      const { data } = await inventoryApi.getProducts(params)
      products.value = data.data || []
      meta.value     = data.meta || { total: 0 }
    } finally {
      isLoading.value = false
    }
  }

  /**
   * Fetch a single product by id.
   * Returns the product object directly so the caller can use it.
   */
  async function fetchProductDetail (id) {
    const { data } = await inventoryApi.getProduct(id)
    return data.data
  }

  async function createProduct (payload) {
    await inventoryApi.createProduct(payload)
    await fetchProducts({ page: 1 })
  }

  async function updateProduct (id, payload) {
    await inventoryApi.updateProduct(id, payload)
  }

  async function deleteProduct (id) {
    await inventoryApi.deleteProduct(id)
    // Optimistic removal — avoids a full re-fetch for a single delete
    products.value = products.value.filter(p => p.id !== id)
  }

  // ── Public API ─────────────────────────────────────────────────
  return {
    products,
    isLoading,
    meta,
    fetchProducts,
    fetchProductDetail,
    createProduct,
    updateProduct,
    deleteProduct
  }
})

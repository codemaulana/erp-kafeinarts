import { api } from 'boot/axios'

export const inventoryApi = {
  getProducts (params = {}) {
    return api.get('/inventory/products', { params })
  },

  getProduct (id) {
    return api.get(`/inventory/products/${id}`)
  },

  createProduct (payload) {
    return api.post('/inventory/products', payload)
  },

  updateProduct (id, payload) {
    return api.put(`/inventory/products/${id}`, payload)
  },

  deleteProduct (id) {
    return api.delete(`/inventory/products/${id}`)
  },

  adjustStock (id, payload) {
    return api.post(`/inventory/products/${id}/adjust`, payload)
  }
}

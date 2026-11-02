/**
 * stores/user.store.js
 * Pinia store untuk Manajemen Pengguna (User Management).
 */
import { defineStore } from 'pinia'
import { ref } from 'vue'
import { api } from 'src/boot/axios'

export const useUserStore = defineStore('user', () => {
  // ── State ───────────────────────────────────────────────────────
  const userList = ref([])
  const isLoading = ref(false)

  // ── Actions ─────────────────────────────────────────────────────
  async function fetchUsers (params = {}) {
    isLoading.value = true
    try {
      const response = await api.get('/users', { params })
      userList.value = response.data.data || []
      return userList.value
    } finally {
      isLoading.value = false
    }
  }

  async function getUser (id) {
    const response = await api.get(`/users/${id}`)
    return response.data.data
  }

  async function createUser (payload) {
    isLoading.value = true
    try {
      const response = await api.post('/users', payload)
      await fetchUsers()
      return response.data.data
    } finally {
      isLoading.value = false
    }
  }

  async function updateUser (id, payload) {
    isLoading.value = true
    try {
      const response = await api.put(`/users/${id}`, payload)
      await fetchUsers()
      return response.data.data
    } finally {
      isLoading.value = false
    }
  }

  async function deleteUser (id) {
    isLoading.value = true
    try {
      await api.delete(`/users/${id}`)
      userList.value = userList.value.filter((u) => u.id !== id)
    } finally {
      isLoading.value = false
    }
  }

  return {
    userList,
    isLoading,
    fetchUsers,
    getUser,
    createUser,
    updateUser,
    deleteUser,
  }
})

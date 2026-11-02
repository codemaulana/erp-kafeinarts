/**
 * boot/axios.js
 * Konfigurasi axios instance dengan interceptors.
 * Pinia harus sudah di-boot sebelum file ini (urutan: pinia → axios).
 */
import axios from 'axios'

const api = axios.create({
  baseURL: process.env.API_BASE_URL || 'http://localhost:8000/api/v1',
  withCredentials: true,
  headers: {
    Accept: 'application/json'
  }
})

// Pasang interceptors setelah Pinia tersedia
export default ({ app }) => {
  // ── Request interceptor — sisipkan Bearer token ─────────────────
  api.interceptors.request.use((config) => {
    try {
      const token = localStorage.getItem('token')
      if (token && token !== 'cookie') {
        config.headers.Authorization = `Bearer ${token}`
      }
    } catch {
      // Abaikan jika akses localStorage dibatasi
    }
    return config
  })

  // ── Response interceptor — tangani 401 global ───────────────────
  api.interceptors.response.use(
    (response) => response,
    async (error) => {
      if (error.response?.status === 401) {
        const url = error.config?.url || ''
        const isAuthRequest = url.includes('/auth/login') || url.includes('/auth/logout') || url === '/login' || url === '/logout'

        // Bersihkan sesi di localStorage
        try {
          localStorage.removeItem('user')
          localStorage.removeItem('token')
        } catch {
          // Abaikan
        }

        // Jika bukan saat request login/logout dan belum di halaman login, arahkan ke login
        if (!isAuthRequest && window.location.pathname !== '/login') {
          window.location.href = '/login'
        }
      }
      return Promise.reject(error)
    }
  )

  app.config.globalProperties.$axios = axios
  app.config.globalProperties.$api   = api
}

export { api }

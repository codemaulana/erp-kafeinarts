/**
 * stores/auth.store.js
 * Pinia store untuk autentikasi (menggantikan Vuex auth module).
 * Token-based auth — menyimpan user & token di localStorage.
 */
import { defineStore } from 'pinia'
import { ref, computed } from 'vue'
import { AUTH } from 'src/services/endpoints'

export const useAuthStore = defineStore('auth', () => {
  // ── State ───────────────────────────────────────────────────────
  const user  = ref(JSON.parse(localStorage.getItem('user'))  || null)
  const token = ref(localStorage.getItem('token') || null)

  // ── Getters ─────────────────────────────────────────────────────
  const isAuthenticated = computed(() => !!token.value)
  const currentUser     = computed(() => user.value)

  // ── Helpers ─────────────────────────────────────────────────────
  function _setAuth (userData, tokenValue) {
    user.value  = userData
    token.value = tokenValue
    localStorage.setItem('user',  JSON.stringify(userData))
    localStorage.setItem('token', tokenValue)
  }

  function _clearAuth () {
    user.value  = null
    token.value = null
    localStorage.removeItem('user')
    localStorage.removeItem('token')
  }

  // ── Actions ─────────────────────────────────────────────────────
  /**
   * Login token-based (Sanctum API token).
   * Axios instance diimport secara lazy untuk menghindari circular dependency
   * saat boot (pinia selesai sebelum axios interceptors terpasang).
   */
  async function login (credentials) {
    const { api } = await import('boot/axios')
    const { data } = await api.post(AUTH.LOGIN, credentials)
    _setAuth(data.user, data.token)
  }

  async function fetchProfile () {
    const { api } = await import('boot/axios')
    const { data } = await api.get(AUTH.PROFILE)
    _setAuth(data.data, token.value)
  }

  async function logout () {
    try {
      const { api } = await import('boot/axios')
      await api.post(AUTH.LOGOUT)
    } catch {
      // best-effort — clear auth regardless
    } finally {
      _clearAuth()
    }
  }

  return {
    // state
    user,
    token,
    // getters
    isAuthenticated,
    currentUser,
    // actions
    login,
    fetchProfile,
    logout
  }
})

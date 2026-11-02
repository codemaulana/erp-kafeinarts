/**
 * router/guards.js
 * Navigation guards — menggunakan Pinia auth store (bukan Vuex).
 * useAuthStore() aman dipanggil di sini karena router dipakai
 * setelah Pinia sudah di-boot.
 */
import { useAuthStore } from 'src/stores/auth.store'

export function setupGuards (router) {
  router.beforeEach((to, from, next) => {
    const authStore      = useAuthStore()
    const requiresAuth   = to.matched.some(record => record.meta.requiresAuth)
    const isAuthenticated = authStore.isAuthenticated

    if (requiresAuth && !isAuthenticated) {
      next({ name: 'login', query: { redirect: to.fullPath } })
    } else if (to.name === 'login' && isAuthenticated) {
      next('/app')
    } else {
      next()
    }
  })
}

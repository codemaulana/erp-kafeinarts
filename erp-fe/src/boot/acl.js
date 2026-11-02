/**
 * boot/acl.js
 * Global properties $can dan $hasRole.
 * Bergantung pada Pinia auth store (harus boot setelah pinia).
 */
import { useAuthStore } from 'src/stores/auth.store'

export default ({ app }) => {
  app.config.globalProperties.$can = (permission) => {
    const authStore = useAuthStore()
    const user = authStore.currentUser
    if (!user) return false
    if (user.role === 'superadmin') return true
    return (user.permissions || []).includes(permission)
  }

  app.config.globalProperties.$hasRole = (roles) => {
    const authStore = useAuthStore()
    const user = authStore.currentUser
    if (!user) return false
    const userRole = user.role || ''
    return Array.isArray(roles) ? roles.includes(userRole) : userRole === roles
  }
}

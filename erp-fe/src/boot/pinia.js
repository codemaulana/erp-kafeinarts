/**
 * Boot: Pinia
 * Harus menjadi boot pertama agar auth store tersedia sebelum
 * axios interceptors dan ACL diinisialisasi.
 */
import { createPinia } from 'pinia'

let pinia

export default ({ app }) => {
  pinia = createPinia()
  app.use(pinia)
}

export { pinia }

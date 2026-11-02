/**
 * composables/useNotify.js
 * Wrapper reusable untuk Quasar Notify plugin.
 * Menghilangkan duplikasi konfigurasi notify di setiap halaman.
 */
import { useQuasar } from 'quasar'

export function useNotify () {
  const $q = useQuasar()

  /**
   * Notifikasi sukses (hijau).
   * @param {string} message
   * @param {object} [opts]  Override opsional
   */
  function notifySuccess (message, opts = {}) {
    $q.notify({
      type: 'positive',
      message,
      icon: 'check_circle',
      position: 'top',
      timeout: 2500,
      ...opts
    })
  }

  /**
   * Notifikasi error (merah).
   * @param {string} message
   * @param {object} [opts]
   */
  function notifyError (message, opts = {}) {
    $q.notify({
      type: 'negative',
      message,
      icon: 'error',
      position: 'top',
      timeout: 3500,
      ...opts
    })
  }

  /**
   * Notifikasi peringatan (kuning/oranye).
   * @param {string} message
   * @param {object} [opts]
   */
  function notifyWarning (message, opts = {}) {
    $q.notify({
      type: 'warning',
      message,
      icon: 'warning',
      position: 'top',
      timeout: 3000,
      ...opts
    })
  }

  /**
   * Notifikasi informasi (biru).
   * @param {string} message
   * @param {object} [opts]
   */
  function notifyInfo (message, opts = {}) {
    $q.notify({
      color: 'blue-6',
      message,
      icon: 'info',
      position: 'top',
      timeout: 3000,
      ...opts
    })
  }

  return { notifySuccess, notifyError, notifyWarning, notifyInfo }
}

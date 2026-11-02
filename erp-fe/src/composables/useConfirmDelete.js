/**
 * composables/useConfirmDelete.js
 * Dialog konfirmasi hapus standar menggunakan Quasar Dialog plugin.
 * Menghilangkan duplikasi $q.dialog({ title: 'Konfirmasi Hapus'... }) di setiap halaman.
 */
import { useQuasar } from 'quasar'

export function useConfirmDelete () {
  const $q = useQuasar()

  /**
   * Tampilkan dialog konfirmasi hapus.
   * Mendukung 3 format pemanggilan:
   * 1. confirmDelete({ itemName, message, onConfirm, onError })
   * 2. confirmDelete('Nama Item / Pesan', async () => { ... })
   * 3. const ok = await confirmDelete('Apakah Anda yakin?')
   */
  function confirmDelete (param1, param2) {
    // Format 3: await confirmDelete('Pesan konfirmasi?') -> return Promise<boolean>
    if (typeof param1 === 'string' && typeof param2 !== 'function') {
      return new Promise((resolve) => {
        $q.dialog({
          title: 'Konfirmasi Hapus',
          message: param1,
          html: true,
          cancel: { label: 'Batal', flat: true, noCaps: true },
          ok: { label: 'Ya, Hapus', unelevated: true, color: 'red-6', noCaps: true },
          persistent: true
        }).onOk(() => resolve(true))
          .onCancel(() => resolve(false))
          .onDismiss(() => resolve(false))
      })
    }

    // Format 1 & 2
    let itemName = ''
    let message = ''
    let onConfirm = null
    let onError = null

    if (typeof param1 === 'string' && typeof param2 === 'function') {
      // Format: confirmDelete('GR', async () => { ... })
      itemName = param1
      onConfirm = param2
    } else if (typeof param1 === 'object' && param1 !== null) {
      // Format: confirmDelete({ itemName, message, onConfirm, onError })
      itemName = param1.itemName
      message = param1.message
      onConfirm = param1.onConfirm
      onError = param1.onError
    }

    const msg = message || `Hapus <strong>${itemName || 'data ini'}</strong>? Tindakan ini tidak dapat dibatalkan.`

    $q.dialog({
      title: 'Konfirmasi Hapus',
      message: msg,
      html: true,
      cancel: { label: 'Batal', flat: true, noCaps: true },
      ok: { label: 'Ya, Hapus', unelevated: true, color: 'red-6', noCaps: true },
      persistent: true
    }).onOk(async () => {
      try {
        if (typeof onConfirm === 'function') {
          await onConfirm()
        }
      } catch (err) {
        if (typeof onError === 'function') {
          onError(err)
        } else {
          $q.notify({
            type: 'negative',
            message: err?.message || 'Gagal menghapus data.',
            position: 'top'
          })
        }
      }
    })
  }

  return { confirmDelete }
}

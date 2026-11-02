/**
 * composables/useConfirmAction.js
 * Dialog konfirmasi generik untuk aksi approve/reject/status update.
 * Memperluas pattern useConfirmDelete untuk kebutuhan approval flow HR.
 */
import { useQuasar } from 'quasar'

export function useConfirmAction() {
  const $q = useQuasar()

  /**
   * Tampilkan dialog konfirmasi aksi.
   *
   * @param {object} opts
   * @param {string}   opts.title      Judul dialog
   * @param {string}   opts.message    Pesan konfirmasi (mendukung HTML)
   * @param {string}   opts.okLabel    Label tombol OK
   * @param {string}   opts.okColor    Warna tombol OK (Quasar color)
   * @param {function} opts.onConfirm  Callback async setelah user mengklik OK
   * @param {function} [opts.onError]  Callback opsional untuk error handling
   */
  function confirmAction({ title, message, okLabel, okColor = 'primary', onConfirm, onError }) {
    $q.dialog({
      title,
      message,
      html: true,
      cancel: { label: 'Batal', flat: true, noCaps: true },
      ok: { label: okLabel, unelevated: true, color: okColor, noCaps: true },
      persistent: true,
    }).onOk(async () => {
      try {
        await onConfirm()
      } catch (err) {
        if (onError) {
          onError(err)
        } else {
          $q.notify({
            type: 'negative',
            message: err?.message || 'Terjadi kesalahan.',
            position: 'top',
          })
        }
      }
    })
  }

  /**
   * Shortcut untuk konfirmasi approve.
   */
  function confirmApprove({ itemName, onConfirm, onError }) {
    confirmAction({
      title: 'Konfirmasi Persetujuan',
      message: `Setujui <strong>${itemName}</strong>?`,
      okLabel: 'Ya, Setujui',
      okColor: 'green-6',
      onConfirm,
      onError,
    })
  }

  /**
   * Shortcut untuk konfirmasi reject.
   */
  function confirmReject({ itemName, onConfirm, onError }) {
    confirmAction({
      title: 'Konfirmasi Penolakan',
      message: `Tolak <strong>${itemName}</strong>? Tindakan ini tidak dapat dibatalkan.`,
      okLabel: 'Ya, Tolak',
      okColor: 'red-6',
      onConfirm,
      onError,
    })
  }

  return { confirmAction, confirmApprove, confirmReject }
}

/**
 * utils/formatter.js
 * Kumpulan fungsi format reusable untuk digunakan di seluruh aplikasi.
 */

// ── Currency ─────────────────────────────────────────────────────
const _currencyFmt = new Intl.NumberFormat('id-ID', {
  style: 'currency',
  currency: 'IDR',
  minimumFractionDigits: 0
})

/**
 * Format angka menjadi string mata uang Rupiah.
 * @param {number} value
 * @returns {string}  e.g. "Rp 1.500.000"
 */
export function formatCurrency (value) {
  return _currencyFmt.format(value || 0)
}

// ── Date ─────────────────────────────────────────────────────────
const _dateFmt = new Intl.DateTimeFormat('id-ID', {
  day: '2-digit',
  month: 'long',
  year: 'numeric'
})

const _shortDateFmt = new Intl.DateTimeFormat('id-ID', {
  day: '2-digit',
  month: 'short',
  year: 'numeric'
})

/**
 * Format tanggal ke format lokal Indonesia.
 * @param {string|Date} value
 * @param {'long'|'short'} style
 * @returns {string}  e.g. "15 Januari 2026" atau "15 Jan 2026"
 */
export function formatDate (value, style = 'long') {
  if (!value) return '—'
  const date = value instanceof Date ? value : new Date(value)
  if (isNaN(date)) return '—'
  return style === 'short' ? _shortDateFmt.format(date) : _dateFmt.format(date)
}

// ── User ─────────────────────────────────────────────────────────
/**
 * Ambil inisial pengguna dari objek user.
 * Mencoba field: name → username → email.
 * @param {object|null} user
 * @returns {string}  e.g. "A"
 */
export function getUserInitial (user) {
  if (!user) return '?'
  const str = user.name || user.username || user.email || ''
  return str.charAt(0).toUpperCase() || '?'
}

// ── Aging ────────────────────────────────────────────────────────
/**
 * Hitung jumlah hari keterlambatan dari tanggal jatuh tempo.
 * @param {string} dateStr  ISO date string "YYYY-MM-DD"
 * @returns {number}
 */
export function daysOverdue (dateStr) {
  if (!dateStr) return 0
  const diff = Date.now() - new Date(dateStr).getTime()
  return Math.max(0, Math.floor(diff / 86_400_000))
}

/**
 * Label teks untuk usia keterlambatan.
 * @param {string} dateStr
 * @returns {string}
 */
export function agingLabel (dateStr) {
  const days = daysOverdue(dateStr)
  if (days === 0) return 'Hari ini'
  if (days <= 60) return `${days} hari`
  return `${days} hari (Kritis)`
}

// ── Number ────────────────────────────────────────────────────────
/**
 * Format angka dengan pemisah ribuan lokal Indonesia.
 * @param {number} value
 * @returns {string}  e.g. "1.500.000"
 */
export function formatNumber (value) {
  return new Intl.NumberFormat('id-ID').format(value || 0)
}

// ── Bulan ─────────────────────────────────────────────────────────
const BULAN = [
  '', 'Januari', 'Februari', 'Maret', 'April', 'Mei', 'Juni',
  'Juli', 'Agustus', 'September', 'Oktober', 'November', 'Desember'
]

/**
 * Konversi nomor bulan (1-12) ke nama bulan Indonesia.
 * @param {number} num
 * @returns {string}
 */
export function bulanLabel (num) {
  return BULAN[num] || ''
}

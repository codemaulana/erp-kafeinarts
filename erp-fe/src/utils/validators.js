/**
 * utils/validators.js
 * Validator reusable untuk digunakan sebagai array :rules di Quasar q-input / q-select.
 * Setiap fungsi mengembalikan (value) => true | string.
 */

/**
 * Field wajib diisi (tidak boleh kosong/null/undefined).
 * @param {string} msg  Pesan kustom
 * @returns {function}
 */
export const required =
  (msg = 'Field ini wajib diisi') =>
  (v) =>
    (v !== null && v !== undefined && String(v).trim() !== '') || msg

/**
 * Nilai harus berupa angka positif (> 0).
 * @param {string} msg
 * @returns {function}
 */
export const positiveNumber =
  (msg = 'Nilai harus lebih dari 0') =>
  (v) =>
    (v !== null && v !== undefined && Number(v) > 0) || msg

/**
 * Nilai tidak boleh negatif (>= 0).
 * @param {string} msg
 * @returns {function}
 */
export const nonNegative =
  (msg = 'Nilai tidak boleh negatif') =>
  (v) =>
    (v !== null && v !== undefined && Number(v) >= 0) || msg

/**
 * Nilai harus dalam rentang [min, max].
 * @param {number} min
 * @param {number} max
 * @param {string} msg
 * @returns {function}
 */
export const inRange = (min, max, msg) => (v) => {
  const n = Number(v)
  return (n >= min && n <= max) || msg || `Nilai harus antara ${min}–${max}`
}

/**
 * Panjang string minimum.
 * @param {number} min
 * @param {string} msg
 * @returns {function}
 */
export const minLength = (min, msg) => (v) =>
  (v && String(v).length >= min) || msg || `Minimal ${min} karakter`

/**
 * Validasi format email sederhana.
 * @param {string} msg
 * @returns {function}
 */
export const email =
  (msg = 'Format email tidak valid') =>
  (v) =>
    !v || /^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(String(v)) || msg

/**
 * Pastikan field tidak kosong dan berupa angka valid.
 * Berguna untuk input number yang bisa menerima string kosong.
 * @param {string} msg
 * @returns {function}
 */
export const requiredNumber =
  (msg = 'Nilai angka wajib diisi') =>
  (v) =>
    (v !== null && v !== undefined && v !== '' && !isNaN(Number(v))) || msg

/**
 * Panjang string harus tepat n karakter.
 * @param {number} n
 * @param {string} msg
 * @returns {function}
 */
export const exactLength = (n, msg) => (v) =>
  (v && String(v).length === n) || msg || `Harus tepat ${n} karakter`

/**
 * Pastikan tanggal tidak sebelum tanggal referensi.
 * @param {function} refDateFn  Fungsi yang mengembalikan string tanggal referensi (YYYY-MM-DD)
 * @param {string} msg
 * @returns {function}
 */
export const dateNotBefore =
  (refDateFn, msg = 'Tanggal tidak boleh sebelum tanggal mulai') =>
  (v) => {
    if (!v) return true
    const ref = refDateFn()
    if (!ref) return true
    return v >= ref || msg
  }

/**
 * Validasi format nomor telepon Indonesia sederhana.
 * Menerima 08xx, +628xx, atau 628xx, panjang 10-15 digit.
 * @param {string} msg
 * @returns {function}
 */
export const phoneNumber =
  (msg = 'Format nomor telepon tidak valid') =>
  (v) =>
    !v || /^(\+?62|0)[0-9]{9,13}$/.test(String(v).replace(/[\s-]/g, '')) || msg

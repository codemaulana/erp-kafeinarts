<template>
  <!--
    LoginPage — Rubick Theme
    Layout : grid bg + 8 coloured orbs | split card (left branding / right form)
    Toggle : light/dark button sits in top-right corner of the right panel
  -->
  <q-page class="lp" :class="isDark ? 'lp--dark' : 'lp--light'">
    <!-- ── Background: CSS grid + animated orbs ── -->
    <div class="lp__bg" aria-hidden="true">
      <div class="lp__grid" />
      <div class="lp__orb lp__orb--1" />
      <div class="lp__orb lp__orb--2" />
      <div class="lp__orb lp__orb--3" />
      <div class="lp__orb lp__orb--4" />
      <div class="lp__orb lp__orb--5" />
      <div class="lp__orb lp__orb--6" />
      <div class="lp__orb lp__orb--7" />
      <div class="lp__orb lp__orb--8" />
    </div>

    <!-- ── Split card ── -->
    <div class="lp__card">
      <!-- ════ Left panel — Branding (hidden ≤ 720px) ════ -->
      <aside class="lp__left" aria-hidden="true">
        <div class="lp__left-deco lp__left-deco--1" />
        <div class="lp__left-deco lp__left-deco--2" />

        <div class="lp__left-content">
          <div class="lp__logo">
            <div class="lp__logo-icon">
              <img
                src="/images/kafeinarts-logo-light.png"
                alt="Kafeinarts logo"
                class="lp__logo-img"
              />
            </div>
            <span class="lp__logo-brand">KAFEINARTS ERP</span>
          </div>
          <div class="lp__divider" />
          <h1 class="lp__headline">Manage and<br />Monitor All</h1>
          <p class="lp__tagline">Your Business</p>
          <ul class="lp__bullets">
            <li class="lp__bullet">
              <span class="lp__bullet-icon">✓</span>
              Kelola transaksi keuangan dengan mudah
            </li>
            <li class="lp__bullet">
              <span class="lp__bullet-icon">✓</span>
              Laporan real-time &amp; analitik bisnis
            </li>
            <li class="lp__bullet">
              <span class="lp__bullet-icon">✓</span>
              Manajemen SDM terpadu &amp; efisien
            </li>
          </ul>
          <button class="lp__cta" type="button" tabindex="-1">
            Pelajari Lebih Lanjut
            <svg viewBox="0 0 16 16" fill="none" width="13" height="13">
              <path
                d="M3 8h10M9 4l4 4-4 4"
                stroke="currentColor"
                stroke-width="1.6"
                stroke-linecap="round"
                stroke-linejoin="round"
              />
            </svg>
          </button>
        </div>
      </aside>

      <!-- ════ Right panel — Login form ════ -->
      <main class="lp__right" role="main">
        <button
          class="lp__theme-btn"
          type="button"
          :aria-label="isDark ? 'Switch to light mode' : 'Switch to dark mode'"
          @click="toggleDark"
        >
          <transition name="lp-icon" mode="out-in">
            <q-icon
              :key="isDark ? 'moon' : 'sun'"
              :name="isDark ? 'light_mode' : 'dark_mode'"
              size="16px"
            />
          </transition>
        </button>

        <div class="lp__form-wrap">
          <div class="lp__mobile-brand">
            <div class="lp__mobile-brand-icon">
              <img
                src="/images/kafeinarts-logo-light.png"
                alt="Kafeinarts logo"
                class="lp__mobile-brand-img"
              />
            </div>
            <span>KAFEINARTS ERP</span>
          </div>

          <header class="lp__form-header">
            <h2 class="lp__form-title">Selamat Datang!</h2>
            <p class="lp__form-sub">Masuk ke akun Anda untuk melanjutkan.</p>
          </header>

          <q-form class="lp__form" greedy @submit.prevent="onSubmit">
            <!-- Email -->
            <div class="lp__field">
              <label class="lp__field-label" for="lp-email">Email Address</label>
              <div
                class="lp__field-wrap"
                :class="{ 'lp__field-wrap--focused': focusedField === 'email' }"
              >
                <span class="lp__field-icon"><q-icon name="mail_outline" size="16px" /></span>
                <input
                  id="lp-email"
                  v-model="form.email"
                  class="lp__input"
                  type="email"
                  placeholder="nama@perusahaan.com"
                  autocomplete="email"
                  required
                  @focus="focusedField = 'email'"
                  @blur="focusedField = null"
                />
              </div>
            </div>

            <!-- Password -->
            <div class="lp__field">
              <label class="lp__field-label" for="lp-password">Password</label>
              <div
                class="lp__field-wrap"
                :class="{ 'lp__field-wrap--focused': focusedField === 'password' }"
              >
                <span class="lp__field-icon"><q-icon name="lock_outline" size="16px" /></span>
                <input
                  id="lp-password"
                  v-model="form.password"
                  class="lp__input"
                  :type="showPwd ? 'text' : 'password'"
                  placeholder="Minimal 6 karakter"
                  autocomplete="current-password"
                  required
                  @focus="focusedField = 'password'"
                  @blur="focusedField = null"
                />
                <button
                  class="lp__pwd-toggle"
                  type="button"
                  :aria-label="showPwd ? 'Sembunyikan password' : 'Tampilkan password'"
                  tabindex="-1"
                  @click="showPwd = !showPwd"
                >
                  <q-icon :name="showPwd ? 'visibility_off' : 'visibility'" size="15px" />
                </button>
              </div>
            </div>

            <!-- Remember + Forgot -->
            <div class="lp__row-meta">
              <label class="lp__remember">
                <input v-model="form.remember" type="checkbox" class="lp__checkbox" />
                <span>Remember me</span>
              </label>
              <button type="button" class="lp__forgot" @click="onForgot">Forgot Password?</button>
            </div>

            <!-- Error -->
            <transition name="lp-err">
              <div v-if="errorMsg" class="lp__error" role="alert" aria-live="assertive">
                <q-icon name="error_outline" size="14px" />
                <span>{{ errorMsg }}</span>
              </div>
            </transition>

            <!-- Submit -->
            <button
              type="submit"
              class="lp__submit"
              :class="{ 'lp__submit--loading': isLoading }"
              :disabled="isLoading"
            >
              <span v-if="!isLoading">Login</span>
              <span v-else><q-spinner-dots color="white" size="22px" /></span>
            </button>
          </q-form>

          <p class="lp__footer-note">
            Butuh akses?
            <button type="button" class="lp__footer-link" @click="onContactAdmin">
              Hubungi Administrator
            </button>
          </p>
        </div>
      </main>
    </div>
  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from 'src/stores/auth.store'
import { useNotify } from 'src/composables/useNotify'

// Maps HTTP status codes to user-friendly error messages
const AUTH_ERRORS = {
  401: 'Email atau password salah. Silakan coba lagi.',
  422: 'Data tidak valid. Periksa kembali input Anda.',
  429: 'Terlalu banyak percobaan. Tunggu sebentar.',
}

const $q = useQuasar()
const router = useRouter()
const authStore = useAuthStore()
const { notifySuccess, notifyInfo } = useNotify()

// ── State ──────────────────────────────────────────────────────
const isDark = ref(false)
const focusedField = ref(null)
const showPwd = ref(false)
const isLoading = ref(false)
const errorMsg = ref('')
const form = reactive({ email: '', password: '', remember: false })

// ── Lifecycle ──────────────────────────────────────────────────
onMounted(() => {
  isDark.value = false
  $q.dark.set(false)

  if (authStore.isAuthenticated) {
    router.replace('/app')
  }
})

// ── Methods ────────────────────────────────────────────────────
function toggleDark() {
  isDark.value = !isDark.value
  $q.dark.set(isDark.value)
}

async function onSubmit() {
  if (!form.email || !form.password) {
    errorMsg.value = 'Email dan password wajib diisi.'
    return
  }
  if (form.password.length < 6) {
    errorMsg.value = 'Password minimal 6 karakter.'
    return
  }

  errorMsg.value = ''
  isLoading.value = true

  try {
    await authStore.login({ email: form.email, password: form.password, remember: form.remember })
    notifySuccess('Login berhasil!')
    router.push('/app')
  } catch (err) {
    const status = err?.response?.status
    const backendMsg = err?.response?.data?.message || err?.response?.data?.error
    errorMsg.value =
      AUTH_ERRORS[status] ||
      backendMsg ||
      (err?.message ? `Gagal menghubungi server: ${err.message}` : null) ||
      (!navigator.onLine ? 'Tidak ada koneksi internet.' : 'Terjadi kesalahan. Silakan coba lagi.')
  } finally {
    isLoading.value = false
  }
}

function onForgot() {
  notifyInfo('Fitur reset password segera hadir. Hubungi admin.')
}

function onContactAdmin() {
  notifyInfo('Hubungi administrator untuk mendapatkan akses.')
}
</script>

<style lang="scss" scoped>
// ================================================================
//  LOGIN PAGE — Rubick Theme
// ================================================================
.lp {
  --lp-page-bg: #eff6ff;
  --lp-grid: rgba(107, 114, 128, 0.18);
  --lp-left-bg: #1e3a8a;
  --lp-right-bg: #ffffff;
  --lp-text: #0f172a;
  --lp-muted: #64748b;
  --lp-border: #e2e8f0;
  --lp-input-bg: #f8fafc;
  --lp-focus: #3b82f6;
  --lp-accent: #2563eb;
  --lp-shadow: 0 25px 50px rgba(15, 23, 42, 0.14), 0 8px 20px rgba(15, 23, 42, 0.08);

  &--dark {
    --lp-page-bg: #0f172a;
    --lp-grid: rgba(255, 255, 255, 0.1);
    --lp-left-bg: #172554;
    --lp-right-bg: #1e293b;
    --lp-text: #f1f5f9;
    --lp-muted: #94a3b8;
    --lp-border: #334155;
    --lp-input-bg: #0f172a;
    --lp-focus: #60a5fa;
    --lp-accent: #3b82f6;
    --lp-shadow: 0 25px 50px rgba(0, 0, 0, 0.5), 0 8px 20px rgba(0, 0, 0, 0.3);
  }
}

.lp {
  position: relative;
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background: var(--lp-page-bg);
  padding: clamp(16px, 3vw, 32px) clamp(12px, 2vw, 16px);
  overflow: hidden;
  transition: background 0.35s ease;
  font-family: 'Inter', 'Segoe UI', system-ui, sans-serif;
}

.lp__bg {
  position: absolute;
  inset: 0;
  pointer-events: none;
  overflow: hidden;
}

.lp__grid {
  position: absolute;
  inset: 0;
  background-image:
    linear-gradient(to right, var(--lp-grid) 1px, transparent 1px),
    linear-gradient(to bottom, var(--lp-grid) 1px, transparent 1px);
  background-size: 40px 40px;
  transition: background-image 0.35s ease;
}

.lp__orb {
  position: absolute;
  border-radius: 50%;
  will-change: transform, opacity;
}

.lp__orb--1 {
  width: 420px;
  height: 420px;
  top: -12%;
  left: -8%;
  background: radial-gradient(circle, rgba(59, 130, 246, 0.7), rgba(37, 99, 235, 0.25));
  filter: blur(72px);
  animation:
    orb-float-1 18s ease-in-out infinite alternate,
    orb-pulse 7s ease-in-out infinite alternate;
}
.lp__orb--2 {
  width: 320px;
  height: 320px;
  top: -5%;
  right: -6%;
  background: radial-gradient(circle, rgba(139, 92, 246, 0.65), rgba(109, 40, 217, 0.22));
  filter: blur(80px);
  animation:
    orb-float-2 14s ease-in-out infinite alternate,
    orb-pulse 5s ease-in-out infinite alternate;
  animation-delay:
    -3s,
    -1s;
}
.lp__orb--3 {
  width: 280px;
  height: 280px;
  top: 35%;
  right: 8%;
  background: radial-gradient(circle, rgba(6, 182, 212, 0.6), rgba(8, 145, 178, 0.18));
  filter: blur(65px);
  animation:
    orb-float-3 20s ease-in-out infinite alternate,
    orb-pulse 9s ease-in-out infinite alternate;
  animation-delay:
    -7s,
    -4s;
}
.lp__orb--4 {
  width: 360px;
  height: 360px;
  bottom: -10%;
  left: 5%;
  background: radial-gradient(circle, rgba(244, 63, 94, 0.55), rgba(225, 29, 72, 0.18));
  filter: blur(90px);
  animation:
    orb-float-4 16s ease-in-out infinite alternate,
    orb-pulse 6s ease-in-out infinite alternate;
  animation-delay:
    -5s,
    -2s;
}
.lp__orb--5 {
  width: 240px;
  height: 240px;
  bottom: 5%;
  left: 40%;
  background: radial-gradient(circle, rgba(251, 191, 36, 0.6), rgba(245, 158, 11, 0.18));
  filter: blur(60px);
  animation:
    orb-float-5 12s ease-in-out infinite alternate,
    orb-pulse 4s ease-in-out infinite alternate;
  animation-delay:
    -2s,
    -3s;
}
.lp__orb--6 {
  width: 200px;
  height: 200px;
  top: 50%;
  left: 2%;
  background: radial-gradient(circle, rgba(16, 185, 129, 0.55), rgba(5, 150, 105, 0.18));
  filter: blur(55px);
  animation:
    orb-float-6 22s ease-in-out infinite alternate,
    orb-pulse 8s ease-in-out infinite alternate;
  animation-delay:
    -10s,
    -5s;
}
.lp__orb--7 {
  width: 180px;
  height: 180px;
  top: 15%;
  left: 45%;
  background: radial-gradient(circle, rgba(99, 102, 241, 0.5), rgba(67, 56, 202, 0.14));
  filter: blur(50px);
  animation:
    orb-float-7 15s ease-in-out infinite alternate,
    orb-pulse 6s ease-in-out infinite alternate;
  animation-delay:
    -8s,
    -6s;
}
.lp__orb--8 {
  width: 260px;
  height: 260px;
  bottom: 12%;
  right: 3%;
  background: radial-gradient(circle, rgba(249, 115, 22, 0.55), rgba(234, 88, 12, 0.16));
  filter: blur(70px);
  animation:
    orb-float-8 17s ease-in-out infinite alternate,
    orb-pulse 7s ease-in-out infinite alternate;
  animation-delay:
    -4s,
    -2s;
}

@keyframes orb-float-1 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  33% {
    transform: translate(60px, -40px) scale(1.08);
  }
  100% {
    transform: translate(80px, 50px) scale(1.05);
  }
}
@keyframes orb-float-2 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  40% {
    transform: translate(-70px, 50px) scale(1.1);
  }
  100% {
    transform: translate(40px, -60px) scale(0.92);
  }
}
@keyframes orb-float-3 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  30% {
    transform: translate(-50px, -70px) scale(1.06);
  }
  100% {
    transform: translate(-20px, 80px) scale(1.04);
  }
}
@keyframes orb-float-4 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  50% {
    transform: translate(80px, -50px) scale(1.12);
  }
  100% {
    transform: translate(-40px, 30px) scale(0.9);
  }
}
@keyframes orb-float-5 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  35% {
    transform: translate(-60px, -30px) scale(1.08);
  }
  100% {
    transform: translate(50px, 60px) scale(0.95);
  }
}
@keyframes orb-float-6 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  45% {
    transform: translate(70px, -60px) scale(1.15);
  }
  100% {
    transform: translate(-30px, 50px) scale(0.88);
  }
}
@keyframes orb-float-7 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  60% {
    transform: translate(-80px, 40px) scale(1.1);
  }
  100% {
    transform: translate(50px, -50px) scale(0.93);
  }
}
@keyframes orb-float-8 {
  0% {
    transform: translate(0, 0) scale(1);
  }
  40% {
    transform: translate(-50px, 70px) scale(1.07);
  }
  100% {
    transform: translate(70px, -40px) scale(0.96);
  }
}
@keyframes orb-pulse {
  0% {
    opacity: 0.7;
  }
  50% {
    opacity: 1;
  }
  100% {
    opacity: 0.55;
  }
}

.lp--dark .lp__orb {
  opacity: 0.65;
}

.lp__card {
  position: relative;
  z-index: 10;
  display: flex;
  width: 100%;
  max-width: 940px;
  min-height: 560px;
  border-radius: 20px;
  overflow: hidden;
  box-shadow: var(--lp-shadow);
  animation: lp-card-in 0.5s cubic-bezier(0.22, 1, 0.36, 1) both;
}

@keyframes lp-card-in {
  from {
    opacity: 0;
    transform: translateY(20px) scale(0.97);
  }
  to {
    opacity: 1;
    transform: translateY(0) scale(1);
  }
}

.lp__left {
  position: relative;
  flex: 0 0 42%;
  background: var(--lp-left-bg);
  padding: clamp(28px, 4vw, 52px) clamp(24px, 4.5vw, 48px);
  display: flex;
  align-items: center;
  overflow: hidden;
  transition: background 0.35s ease;
  &::after {
    content: '';
    position: absolute;
    top: 0;
    right: 0;
    width: 1px;
    height: 100%;
    background: rgba(255, 255, 255, 0.1);
    pointer-events: none;
  }
  @media (max-width: 720px) {
    display: none;
  }
}

.lp__left-deco {
  position: absolute;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.06);
  pointer-events: none;
  &--1 {
    width: 300px;
    height: 300px;
    bottom: -80px;
    right: -80px;
  }
  &--2 {
    width: 180px;
    height: 180px;
    top: -50px;
    left: -50px;
    background: rgba(255, 255, 255, 0.04);
  }
}

.lp__left-content {
  display: flex;
  flex-direction: column;
  gap: clamp(14px, 1.8vw, 22px);
  width: 100%;
  position: relative;
  z-index: 2;
}

.lp__logo {
  display: flex;
  align-items: center;
  gap: 12px;
}
.lp__logo-icon {
  width: clamp(42px, 5vw, 56px);
  height: clamp(42px, 5vw, 56px);
  min-width: clamp(42px, 5vw, 56px);
  border-radius: 14px;
  background: linear-gradient(135deg, rgba(255, 255, 255, 0.2), rgba(255, 255, 255, 0.08));
  border: 1px solid rgba(255, 255, 255, 0.25);
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.2);
  display: flex;
  align-items: center;
  justify-content: center;
  overflow: hidden;
  flex-shrink: 0;
  padding: 5px;
}
.lp__logo-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}
.lp__logo-brand {
  font-size: clamp(12px, 1.4vw, 16px);
  font-weight: 800;
  color: #fff;
  letter-spacing: 1.5px;
  text-transform: uppercase;
  line-height: 1;
}
.lp__divider {
  width: 36px;
  height: 2px;
  background: rgba(255, 255, 255, 0.3);
  border-radius: 2px;
}
.lp__headline {
  font-size: clamp(20px, 2.5vw, 30px);
  font-weight: 800;
  color: #fff;
  line-height: 1.22;
  margin: 0;
  letter-spacing: -0.4px;
}
.lp__tagline {
  font-size: clamp(20px, 2.5vw, 30px);
  font-weight: 800;
  color: rgba(255, 255, 255, 0.38);
  margin: -12px 0 0;
  letter-spacing: -0.4px;
}
.lp__bullets {
  list-style: none;
  margin: 0;
  padding: 0;
  display: flex;
  flex-direction: column;
  gap: clamp(7px, 1vw, 12px);
}
.lp__bullet {
  display: flex;
  align-items: flex-start;
  gap: 10px;
  font-size: clamp(11px, 1.1vw, 13px);
  color: rgba(255, 255, 255, 0.7);
  line-height: 1.55;
}
.lp__bullet-icon {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  width: 18px;
  height: 18px;
  min-width: 18px;
  border-radius: 50%;
  background: rgba(255, 255, 255, 0.15);
  color: #fff;
  font-size: 10px;
  font-weight: 700;
  margin-top: 1px;
}
.lp__cta {
  display: inline-flex;
  align-items: center;
  gap: 7px;
  background: rgba(255, 255, 255, 0.12);
  border: 1px solid rgba(255, 255, 255, 0.22);
  color: #fff;
  font-size: 12px;
  font-weight: 600;
  font-family: inherit;
  padding: 9px 18px;
  border-radius: 9999px;
  cursor: pointer;
  align-self: flex-start;
  transition: background 0.18s ease;
  &:hover {
    background: rgba(255, 255, 255, 0.2);
  }
}

.lp__right {
  position: relative;
  flex: 1;
  background: var(--lp-right-bg);
  display: flex;
  align-items: center;
  justify-content: center;
  padding: clamp(24px, 4vw, 44px) clamp(20px, 4vw, 40px);
  transition: background 0.35s ease;
}

.lp__theme-btn {
  position: absolute;
  top: clamp(10px, 2vw, 16px);
  right: clamp(10px, 2vw, 16px);
  z-index: 10;
  width: clamp(30px, 4vw, 36px);
  height: clamp(30px, 4vw, 36px);
  border-radius: 50%;
  border: 1px solid var(--lp-border);
  background: var(--lp-right-bg);
  color: var(--lp-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 2px 8px rgba(15, 23, 42, 0.1);
  transition:
    background 0.18s ease,
    color 0.18s ease,
    transform 0.18s ease;
  &:hover {
    background: var(--lp-border);
    color: var(--lp-accent);
    transform: scale(1.08) rotate(12deg);
  }
  .lp--dark & {
    border-color: var(--lp-border);
    background: rgba(255, 255, 255, 0.06);
    color: var(--lp-muted);
    &:hover {
      background: rgba(255, 255, 255, 0.12);
      color: #93c5fd;
    }
  }
}

.lp-icon-enter-active,
.lp-icon-leave-active {
  transition:
    opacity 0.15s ease,
    transform 0.15s ease;
}
.lp-icon-enter-from {
  opacity: 0;
  transform: rotate(-30deg) scale(0.7);
}
.lp-icon-leave-to {
  opacity: 0;
  transform: rotate(30deg) scale(0.7);
}

.lp__form-wrap {
  width: 100%;
  max-width: 360px;
  display: flex;
  flex-direction: column;
}

.lp__mobile-brand {
  display: none;
  align-items: center;
  gap: 10px;
  margin-bottom: 22px;
  font-size: 14px;
  font-weight: 800;
  color: var(--lp-accent);
  letter-spacing: 0.8px;
  text-transform: uppercase;
  @media (max-width: 720px) {
    display: flex;
  }
}
.lp__mobile-brand-icon {
  width: 36px;
  height: 36px;
  border-radius: 10px;
  background: linear-gradient(135deg, #1e3a8a, #172554);
  border: 1px solid rgba(255, 255, 255, 0.2);
  box-shadow: 0 3px 10px rgba(30, 58, 138, 0.25);
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  padding: 4px;
}
.lp__mobile-brand-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  filter: drop-shadow(0 1px 3px rgba(0, 0, 0, 0.25));
}

.lp__form-header {
  margin-bottom: 22px;
}
.lp__form-title {
  font-size: clamp(18px, 2.2vw, 24px);
  font-weight: 800;
  color: var(--lp-text);
  margin: 0 0 6px;
  letter-spacing: -0.3px;
  line-height: 1.2;
}
.lp__form-sub {
  font-size: 13px;
  color: var(--lp-muted);
  margin: 0;
}
.lp__form {
  display: flex;
  flex-direction: column;
  gap: 14px;
}
.lp__field {
  display: flex;
  flex-direction: column;
  gap: 5px;
}
.lp__field-label {
  font-size: 12px;
  font-weight: 600;
  color: var(--lp-text);
}
.lp__field-wrap {
  position: relative;
  display: flex;
  align-items: center;
  border: 1.5px solid var(--lp-border);
  border-radius: 10px;
  background: var(--lp-input-bg);
  transition:
    border-color 0.15s ease,
    box-shadow 0.15s ease;
  &--focused {
    border-color: var(--lp-focus);
    box-shadow: 0 0 0 3px rgba(59, 130, 246, 0.15);
    background: var(--lp-right-bg);
  }
}
.lp__field-icon {
  position: absolute;
  left: 12px;
  color: var(--lp-muted);
  display: flex;
  align-items: center;
  pointer-events: none;
  transition: color 0.15s ease;
  .lp__field-wrap--focused & {
    color: var(--lp-focus);
  }
}
.lp__input {
  width: 100%;
  height: 42px;
  padding: 0 38px 0 38px;
  background: transparent;
  border: none;
  outline: none;
  font-size: 13.5px;
  font-family: inherit;
  color: var(--lp-text);
  border-radius: 10px;
  &::placeholder {
    color: var(--lp-muted);
    opacity: 0.7;
  }
  &:-webkit-autofill {
    -webkit-box-shadow: 0 0 0 40px var(--lp-input-bg) inset !important;
    -webkit-text-fill-color: var(--lp-text) !important;
  }
}
.lp__pwd-toggle {
  position: absolute;
  right: 10px;
  width: 28px;
  height: 28px;
  border: none;
  background: transparent;
  color: var(--lp-muted);
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 6px;
  transition: color 0.15s ease;
  &:hover {
    color: var(--lp-text);
  }
}
.lp__row-meta {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 8px;
  flex-wrap: wrap;
}
.lp__remember {
  display: flex;
  align-items: center;
  gap: 7px;
  cursor: pointer;
  font-size: 12.5px;
  color: var(--lp-muted);
  user-select: none;
}
.lp__checkbox {
  width: 15px;
  height: 15px;
  accent-color: var(--lp-accent);
  cursor: pointer;
  flex-shrink: 0;
}
.lp__forgot {
  background: none;
  border: none;
  font-size: 12.5px;
  font-family: inherit;
  color: var(--lp-accent);
  cursor: pointer;
  padding: 0;
  transition: opacity 0.15s ease;
  &:hover {
    opacity: 0.75;
  }
}
.lp__error {
  display: flex;
  align-items: center;
  gap: 7px;
  padding: 10px 13px;
  border-radius: 8px;
  background: #fef2f2;
  border: 1px solid #fecaca;
  color: #dc2626;
  font-size: 12.5px;
  font-weight: 500;
  .lp--dark & {
    background: #450a0a;
    border-color: #7f1d1d;
    color: #fca5a5;
  }
}
.lp-err-enter-active,
.lp-err-leave-active {
  transition:
    opacity 0.2s ease,
    transform 0.2s ease;
}
.lp-err-enter-from {
  opacity: 0;
  transform: translateY(-6px);
}
.lp-err-leave-to {
  opacity: 0;
  transform: translateY(-4px);
}
.lp__submit {
  height: 44px;
  margin-top: 4px;
  border-radius: 10px;
  border: none;
  background: linear-gradient(135deg, #3b82f6, #1e40af);
  color: #fff;
  font-size: 14px;
  font-weight: 600;
  font-family: inherit;
  cursor: pointer;
  width: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  box-shadow: 0 4px 14px rgba(30, 64, 175, 0.35);
  letter-spacing: 0.2px;
  transition:
    opacity 0.18s ease,
    transform 0.18s ease,
    box-shadow 0.18s ease;
  &:hover:not(:disabled) {
    opacity: 0.92;
    transform: translateY(-1px);
    box-shadow: 0 6px 20px rgba(30, 64, 175, 0.42);
  }
  &:active:not(:disabled) {
    transform: translateY(0);
    box-shadow: 0 2px 8px rgba(30, 64, 175, 0.28);
  }
  &:disabled,
  &--loading {
    opacity: 0.75;
    cursor: not-allowed;
    transform: none;
  }
}
.lp__footer-note {
  margin-top: 20px;
  text-align: center;
  font-size: 12.5px;
  color: var(--lp-muted);
}
.lp__footer-link {
  background: none;
  border: none;
  font-size: 12.5px;
  font-family: inherit;
  font-weight: 600;
  color: var(--lp-accent);
  cursor: pointer;
  padding: 0;
  transition: opacity 0.15s ease;
  &:hover {
    opacity: 0.75;
  }
}
</style>

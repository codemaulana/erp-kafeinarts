<template>
  <q-header class="app-header">
    <q-toolbar class="app-header__toolbar">

      <!-- ════ Sidebar Toggle — q-btn flat round dense ════ -->
      <q-btn
        flat round dense
        :icon="drawerOpen ? 'menu_open' : 'menu'"
        class="sidebar-toggle-btn"
        :aria-label="drawerOpen ? 'Sembunyikan sidebar' : 'Tampilkan sidebar'"
        @click="onToggle"
      >
        <q-tooltip anchor="bottom middle" self="top middle" :offset="[0, 6]">
          {{ drawerOpen ? 'Sembunyikan sidebar' : 'Tampilkan sidebar' }}
        </q-tooltip>
      </q-btn>

      <!-- ── Mobile / Closed-drawer Brand ── -->
      <div v-if="!drawerOpen" class="app-header__brand-mini">
        <div class="app-header__brand-icon">
          <img :src="headerLogoSrc" alt="Kafeinarts" class="app-header__brand-img" />
        </div>
        <span class="app-header__brand-text">KAFEINARTS</span>
      </div>

      <!-- ── Breadcrumb (icon + segments from route meta) ── -->
      <nav class="app-header__breadcrumb" aria-label="Breadcrumb">
        <template v-for="(seg, idx) in breadcrumbSegments" :key="idx">
          <!-- Separator — skip before first item -->
          <span v-if="idx > 0" class="app-header__bc-sep" aria-hidden="true">
            <svg viewBox="0 0 6 10" fill="none" width="5" height="8">
              <path d="M1 1l4 4-4 4" stroke="currentColor" stroke-width="1.4"
                    stroke-linecap="round" stroke-linejoin="round"/>
            </svg>
          </span>

          <!-- Root segment (first): icon + semi-bold label -->
          <span v-if="idx === 0" class="app-header__bc-root">
            <q-icon v-if="seg.icon" :name="seg.icon" size="14px" class="app-header__bc-icon" />
            {{ seg.label }}
          </span>

          <!-- Middle segments: muted, hidden ≤640px -->
          <span
            v-else-if="idx < breadcrumbSegments.length - 1"
            class="app-header__bc-mid"
          >{{ seg.label }}</span>

          <!-- Last segment: bold active page name -->
          <span
            v-else
            class="app-header__bc-current"
            :title="seg.label"
          >{{ seg.label }}</span>
        </template>
      </nav>

      <div class="app-header__spacer" />

      <!-- ── Right actions ── -->
      <div class="app-header__actions">

        <!-- Search (hidden ≤540px) -->
        <div class="app-header__search-wrap">
          <span class="app-header__search-icon" aria-hidden="true">
            <svg viewBox="0 0 18 18" fill="none" width="14" height="14">
              <circle cx="7.5" cy="7.5" r="5.5" stroke="currentColor" stroke-width="1.5"/>
              <path d="M12 12l4 4" stroke="currentColor" stroke-width="1.5" stroke-linecap="round"/>
            </svg>
          </span>
          <input v-model="searchQuery" class="app-header__search" type="search" placeholder="Search..." aria-label="Search" />
        </div>

        <!-- Dark mode toggle -->
        <button class="app-header__icon-btn" type="button"
          :aria-label="$q.dark.isActive ? 'Switch to light mode' : 'Switch to dark mode'"
          @click="$q.dark.toggle()">
          <q-icon :name="$q.dark.isActive ? 'light_mode' : 'dark_mode'" size="17px" />
        </button>

        <!-- Notification bell -->
        <button class="app-header__icon-btn app-header__icon-btn--notif" type="button" aria-label="Notifikasi">
          <q-icon name="notifications_none" size="17px" />
          <span class="app-header__notif-dot" aria-hidden="true" />
        </button>

        <!-- User avatar + dropdown -->
        <div class="app-header__user-wrap">
          <button class="app-header__user-btn" type="button" aria-label="Menu pengguna">
            <span class="app-header__avatar">{{ userInitial }}</span>
          </button>

          <q-menu anchor="bottom right" self="top right" transition-show="jump-down" transition-hide="jump-up">
            <div class="hdr-dropdown">
              <div class="hdr-dropdown__user">
                <span class="hdr-dropdown__avatar">{{ userInitial }}</span>
                <div class="hdr-dropdown__user-info">
                  <div class="hdr-dropdown__name">{{ currentUser?.name || currentUser?.username || 'Super Admin' }}</div>
                  <div class="hdr-dropdown__email">{{ currentUser?.email || '' }}</div>
                </div>
              </div>
              <div class="hdr-dropdown__sep" />
              <button class="hdr-dropdown__item" type="button" @click="onProfile">
                <q-icon name="person_outline" size="15px" /> Profil Saya
              </button>
              <button class="hdr-dropdown__item" type="button" @click="onSettings">
                <q-icon name="settings" size="15px" /> Pengaturan
              </button>
              <div class="hdr-dropdown__sep" />
              <button class="hdr-dropdown__item hdr-dropdown__item--danger" type="button" @click="onLogout">
                <q-icon name="logout" size="15px" /> Keluar
              </button>
            </div>
          </q-menu>
        </div>

      </div>
    </q-toolbar>
  </q-header>
</template>

<script setup>
import { ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useQuasar } from 'quasar'
import { useAuthStore } from 'src/stores/auth.store'
import { useNotify } from 'src/composables/useNotify'
import { getUserInitial } from 'src/utils/formatter'

const router = useRouter()

// ── Fallback map: route-name → { icon, label } ────────────────
// Used only when route.meta.breadcrumb is missing.
const FALLBACK_MAP = {
  dashboard:                   { icon: 'dashboard',              label: 'Dashboard' },
  'inventory-products':        { icon: 'category',               label: 'Inventaris' },
  'finance-master-kategori':   { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-master-rekening':   { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-master-pajak':      { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-jurnal-kas':        { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-invoice':           { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-tagihan-vendor':    { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-laporan-cashflow':  { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-laporan-laba-rugi': { icon: 'account_balance_wallet', label: 'Finance'   },
  'finance-laporan-ar-aging':  { icon: 'account_balance_wallet', label: 'Finance'   },
  'hr-employees':              { icon: 'people_alt',             label: 'HR'        }
}

defineProps({
  drawerOpen: { type: Boolean, default: true }
})

const emit = defineEmits(['toggle-drawer'])

const $q        = useQuasar()
const route     = useRoute()
const authStore = useAuthStore()
const { notifyInfo } = useNotify()

const searchQuery = ref('')

const currentUser = computed(() => authStore.currentUser)
const userInitial = computed(() => getUserInitial(currentUser.value))

// ── Header logo for closed drawer / mobile view ────────────────
const headerLogoSrc = computed(() => $q.dark.isActive ? '/images/kafeinarts-logo-light.png' : '/images/kafeinarts-logo-dark.png')

/**
 * Breadcrumb segments from route.meta.breadcrumb.
 * Expected format: [{ icon?, label }, { label }, { label }]
 *
 * Fallbacks (in order):
 *   1. Legacy string array → convert first item to use FALLBACK_MAP icon
 *   2. No meta → single segment from FALLBACK_MAP
 *   3. Unknown route → [{ icon: 'home', label: 'App' }]
 */
const breadcrumbSegments = computed(() => {
  const meta = route?.meta?.breadcrumb

  // ── New format: array of objects ─────────────────────────
  if (Array.isArray(meta) && meta.length && typeof meta[0] === 'object') {
    return meta
  }

  // ── Legacy: array of strings ─────────────────────────────
  if (Array.isArray(meta) && meta.length && typeof meta[0] === 'string') {
    return meta.map((label, i) => (
      i === 0 ? { icon: FALLBACK_MAP[route?.name]?.icon, label } : { label }
    ))
  }

  // ── No meta — single segment from FALLBACK_MAP ───────────
  const fallback = FALLBACK_MAP[route?.name]
  return fallback ? [fallback] : [{ icon: 'home', label: 'App' }]
})

function onToggle () {
  emit('toggle-drawer')
}

async function onLogout () {
  await authStore.logout()
  window.location.href = '/login'
}

function onProfile ()  { notifyInfo('Halaman profil segera hadir.') }
function onSettings () { router.push({ name: 'user-management' }) }
</script>

<style lang="scss" scoped>
.app-header {
  background: var(--rb-header-bg, #FFFFFF);
  border-bottom: 1px solid var(--rb-header-border, #E2E8F0);
  box-shadow: var(--rb-header-shadow);
  transition: background 0.30s ease, border-color 0.30s ease;
  body.body--dark & { background: var(--rb-header-bg, #1E293B); border-bottom-color: var(--rb-header-border, #334155); }
}

.app-header__toolbar {
  height: 60px;
  padding: 0 16px 0 8px;
  display: flex;
  align-items: center;
  min-width: 0;
  gap: 0;
}

// ── Sidebar Toggle ─────────────────────────────────────────────
.sidebar-toggle-btn {
  margin-right: 10px;
  flex-shrink: 0;
  color: var(--rb-text-muted, #64748B) !important;
  border-radius: 10px !important;

  &:hover {
    background: var(--rb-primary-soft, rgba(59, 130, 246, 0.10)) !important;
    color: #3B82F6 !important;
  }

  body.body--dark & {
    color: var(--rb-text-muted, #94A3B8) !important;
    &:hover { background: rgba(59, 130, 246, 0.15) !important; color: #60A5FA !important; }
  }

  :deep(.q-icon) { font-size: 22px !important; }
}

// ── Breadcrumb ─────────────────────────────────────────────────
.app-header__breadcrumb {
  display: flex;
  align-items: center;
  gap: 5px;
  min-width: 0;
  flex-shrink: 1;
  overflow: hidden;
}

// First segment: icon + label — semi-bold, always visible
.app-header__bc-root {
  display: inline-flex;
  align-items: center;
  gap: 5px;
  font-size: 13px;
  font-weight: 600;
  color: var(--rb-text-muted, #64748B);
  white-space: nowrap;
  flex-shrink: 0;

  body.body--dark & { color: var(--rb-text-muted, #94A3B8); }
}

.app-header__bc-icon {
  opacity: 0.75;
  vertical-align: middle;
}

// Separator chevron
.app-header__bc-sep {
  color: var(--rb-text-light, #94A3B8);
  display: flex;
  align-items: center;
  flex-shrink: 0;
}

// Middle segments: muted, hidden on small screens
.app-header__bc-mid {
  font-size: 13px;
  font-weight: 400;
  color: var(--rb-text-muted, #64748B);
  white-space: nowrap;
  flex-shrink: 0;

  @media (max-width: 640px) { display: none; }

  body.body--dark & { color: var(--rb-text-muted, #94A3B8); }
}

// Last segment: bold, current page — always visible
.app-header__bc-current {
  font-size: 13px;
  font-weight: 700;
  color: var(--rb-text-primary, #0F172A);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;

  body.body--dark & { color: #F1F5F9; }
}

// Also hide the sep before a hidden middle segment at ≤640px
@media (max-width: 640px) {
  // Hide the sep that immediately follows the root (before the hidden mid)
  // when there are exactly 3 segments: root › mid › current
  .app-header__bc-sep:has(+ .app-header__bc-mid) {
    display: none;
  }
}

.app-header__spacer { flex: 1; min-width: 8px; }

// ── Right actions ──────────────────────────────────────────────
.app-header__actions { display: flex; align-items: center; gap: 6px; flex-shrink: 0; }

.app-header__search-wrap {
  position: relative;
  display: flex;
  align-items: center;
  @media (max-width: 540px) { display: none; }
}
.app-header__search-icon {
  position: absolute; left: 10px;
  color: var(--rb-text-light, #94A3B8);
  display: flex; align-items: center; pointer-events: none;
}
.app-header__search {
  width: clamp(130px, 15vw, 180px); height: 34px;
  padding: 0 12px 0 32px; border-radius: 8px;
  border: 1px solid var(--rb-border, #E2E8F0);
  background: var(--rb-bg, #F1F5F9);
  font-size: 12.5px; color: var(--rb-text-primary, #0F172A); font-family: inherit; outline: none;
  transition: border-color 0.15s ease, box-shadow 0.15s ease, width 0.20s ease;

  &::placeholder { color: var(--rb-text-light, #94A3B8); }
  &:focus {
    width: clamp(160px, 18vw, 220px);
    border-color: var(--rb-primary-light, #3B82F6);
    box-shadow: 0 0 0 3px rgba(59,130,246,0.15);
    background: var(--rb-surface, #FFFFFF);
  }
  &::-webkit-search-cancel-button { display: none; }
  body.body--dark & {
    background: var(--rb-bg, #0F172A); color: #F1F5F9;
    border-color: var(--rb-border, #334155);
    &:focus { background: var(--rb-surface, #1E293B); }
  }
}

.app-header__icon-btn {
  position: relative; width: 34px; height: 34px; border-radius: 8px;
  border: 1px solid var(--rb-border, #E2E8F0); background: var(--rb-surface, #FFFFFF);
  color: var(--rb-text-muted, #64748B); cursor: pointer;
  display: flex; align-items: center; justify-content: center;
  transition: background 0.15s ease, color 0.15s ease, border-color 0.15s ease;

  &:hover { background: var(--rb-bg); color: var(--rb-primary-light, #3B82F6); border-color: #CBD5E1; }
  body.body--dark & {
    background: var(--rb-surface, #1E293B); border-color: var(--rb-border, #334155);
    color: var(--rb-text-muted, #94A3B8);
    &:hover { background: var(--rb-bg, #0F172A); color: #60A5FA; }
  }
}

.app-header__notif-dot {
  position: absolute; top: 7px; right: 7px;
  width: 7px; height: 7px; border-radius: 50%;
  background: #EF4444;
  border: 1.5px solid var(--rb-surface, #FFFFFF);
  body.body--dark & { border-color: var(--rb-surface, #1E293B); }
}

.app-header__user-wrap { position: relative; }
.app-header__user-btn  {
  display: flex; align-items: center; padding: 0; border: none;
  background: transparent; cursor: pointer; border-radius: 50%;
  transition: opacity 0.15s ease;
  &:hover { opacity: 0.85; }
}
.app-header__avatar {
  width: 34px; height: 34px; border-radius: 50%;
  background: linear-gradient(135deg, #3B82F6, #1E40AF);
  color: #FFF; font-size: 13px; font-weight: 700;
  display: flex; align-items: center; justify-content: center;
  box-shadow: 0 2px 6px rgba(59,130,246,0.35);
  flex-shrink: 0; user-select: none;
}

// ── Dropdown ───────────────────────────────────────────────────
.hdr-dropdown { min-width: 210px; padding: 8px 0; }
.hdr-dropdown__user { display: flex; align-items: center; gap: 11px; padding: 8px 14px 12px; }
.hdr-dropdown__avatar {
  width: 38px; height: 38px; border-radius: 50%;
  background: linear-gradient(135deg, #3B82F6, #1E40AF);
  color: #FFF; font-size: 15px; font-weight: 700;
  display: flex; align-items: center; justify-content: center; flex-shrink: 0;
}
.hdr-dropdown__user-info { display: flex; flex-direction: column; min-width: 0; gap: 2px; }
.hdr-dropdown__name  {
  font-size: 13px; font-weight: 700; color: var(--rb-text-primary, #0F172A);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.hdr-dropdown__email {
  font-size: 11px; color: var(--rb-text-muted, #64748B);
  white-space: nowrap; overflow: hidden; text-overflow: ellipsis;
}
.hdr-dropdown__sep  { height: 1px; background: var(--rb-border, #E2E8F0); margin: 4px 0; }
.hdr-dropdown__item {
  display: flex; align-items: center; gap: 9px;
  width: 100%; padding: 8px 14px;
  background: transparent; border: none; cursor: pointer;
  font-size: 13px; font-family: inherit;
  color: var(--rb-text-primary, #0F172A); text-align: left;
  transition: background 0.15s ease, color 0.15s ease;

  &:hover { background: var(--rb-bg, #F1F5F9); color: #3B82F6; }
  &--danger { color: #EF4444; &:hover { background: #FEF2F2; color: #EF4444; } }
}

// ── Responsive Mini Brand in Header ────────────────────────────
.app-header__brand-mini {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-right: 14px;
  flex-shrink: 0;
}
.app-header__brand-icon {
  width: 30px;
  height: 30px;
  display: flex;
  align-items: center;
  justify-content: center;
  border-radius: 8px;
  padding: 3px;
  background: rgba(59, 130, 246, 0.08);
  border: 1px solid rgba(59, 130, 246, 0.18);
  box-shadow: 0 2px 6px rgba(0, 0, 0, 0.06);
  body.body--dark & {
    background: rgba(255, 255, 255, 0.08);
    border-color: rgba(255, 255, 255, 0.15);
  }
}
.app-header__brand-img {
  width: 100%;
  height: 100%;
  object-fit: contain;
}
.app-header__brand-text {
  font-size: 13px;
  font-weight: 800;
  letter-spacing: 0.6px;
  color: var(--rb-text-primary, #0F172A);
  body.body--dark & {
    color: #F1F5F9;
  }
}
</style>

<template>
  <q-drawer
    :model-value="modelValue"
    show-if-above
    :width="260"
    :breakpoint="768"
    class="app-sidebar"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <!-- q-scroll-area fills the full drawer height -->
    <q-scroll-area
      class="sidebar-scroll"
      :content-style="{ display: 'flex', flexDirection: 'column', minHeight: '100%' }"
    >
      <div class="sidebar-inner">
        <!-- ── Logo ─────────────────────────────────────────── -->
        <div class="sidebar-logo">
          <div class="sidebar-logo__icon-wrap">
            <img :src="logoSrc" alt="Kafeinarts logo" class="sidebar-logo__img" />
          </div>
          <span class="sidebar-logo__brand">KAFEINARTS ERP</span>
        </div>

        <!-- ── Navigation ──────────────────────────────────── -->
        <q-list class="sidebar-nav" padding>
          <!-- Dashboard -->
          <q-item
            v-ripple
            clickable
            :active="isDashActive"
            active-class="nav-item--active"
            class="nav-item nav-item--dashboard"
            :to="{ name: 'dashboard' }"
            exact
          >
            <q-item-section avatar class="nav-item__icon-section">
              <q-icon
                name="dashboard"
                :color="$q.dark.isActive ? 'blue-4' : 'blue-8'"
                size="18px"
              />
            </q-item-section>
            <q-item-section class="nav-item__label-section">Dashboard</q-item-section>
          </q-item>

          <!-- Inventaris -->
          <q-item
            v-if="canAccessInventory"
            v-ripple
            clickable
            :active="isInventoryActive"
            active-class="nav-item--active"
            class="nav-item nav-item--inventory"
            :to="{ name: 'inventory-products' }"
          >
            <q-item-section avatar class="nav-item__icon-section">
              <q-icon
                name="category"
                :color="$q.dark.isActive ? 'orange-4' : 'orange-8'"
                size="18px"
              />
            </q-item-section>
            <q-item-section class="nav-item__label-section">Inventaris</q-item-section>
          </q-item>

          <!-- Finance menus -->
          <template v-if="canAccessFinance">
            <q-item-label header class="sidebar-label">Finance - Data Master</q-item-label>
            <q-item
              v-for="item in masterItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--finance"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>

            <q-item-label header class="sidebar-label">Finance - Data Transaksi</q-item-label>
            <q-item
              v-for="item in transaksiItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--finance"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>

            <q-item-label header class="sidebar-label">Finance - Laporan</q-item-label>
            <q-item
              v-for="item in laporanItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--finance"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>
          </template>

          <!-- HR menus -->
          <template v-if="canAccessHr">
            <q-item-label header class="sidebar-label">HR - Data Master</q-item-label>
            <q-item
              v-for="item in hrMasterItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--hr"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>

            <q-item-label header class="sidebar-label">HR - Data Transaksi</q-item-label>
            <q-item
              v-for="item in hrTransaksiItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--hr"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>

            <q-item-label header class="sidebar-label">HR - Laporan</q-item-label>
            <q-item
              v-for="item in hrLaporanItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--hr"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>
          </template>
          <!-- SC menus -->
          <template v-if="canAccessSc">
            <q-item-label header class="sidebar-label">Supply Chain - Data Master</q-item-label>
            <q-item
              v-for="item in scMasterItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--sc"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>

            <q-item-label header class="sidebar-label">Supply Chain - Data Transaksi</q-item-label>
            <q-item
              v-for="item in scTransaksiItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--sc"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>

            <q-item-label header class="sidebar-label">Supply Chain - Laporan</q-item-label>
            <q-item
              v-for="item in scLaporanItems"
              :key="item.route"
              v-ripple
              clickable
              active-class="nav-item--active"
              class="nav-item nav-item--sc"
              :to="{ name: item.route }"
              exact
            >
              <q-item-section avatar class="nav-item__icon-section">
                <q-icon
                  :name="item.icon"
                  :color="$q.dark.isActive ? item.iconDarkColor || item.iconColor : item.iconColor"
                  size="18px"
                />
              </q-item-section>
              <q-item-section class="nav-item__label-section">{{ item.label }}</q-item-section>
            </q-item>
          </template>

          <!-- Manajemen Pengguna -->
          <q-item
            v-if="canAccessUserManagement"
            v-ripple
            clickable
            active-class="nav-item--active"
            class="nav-item nav-item--user q-mt-sm"
            :to="{ name: 'user-management' }"
            exact
          >
            <q-item-section avatar class="nav-item__icon-section">
              <q-icon
                name="manage_accounts"
                :color="$q.dark.isActive ? 'red-4' : 'red-8'"
                size="18px"
              />
            </q-item-section>
            <q-item-section class="nav-item__label-section">Manajemen Pengguna</q-item-section>
          </q-item>
        </q-list>

        <!-- ── Footer ──────────────────────────────────────── -->
        <div class="sidebar-footer">
          <q-avatar size="34px" class="sidebar-footer__avatar" font-size="13px">
            {{ userInitial }}
          </q-avatar>
          <div class="sidebar-footer__info">
            <span class="sidebar-footer__name">
              {{ currentUser?.name || currentUser?.username || 'Super Admin' }}
            </span>
            <span class="sidebar-footer__role">{{ formatRoleLabel(currentUser?.role) }}</span>
          </div>
          <q-btn
            flat
            round
            dense
            icon="logout"
            size="sm"
            class="sidebar-footer__logout"
            aria-label="Logout"
            @click="onLogout"
          >
            <q-tooltip anchor="top right" self="bottom right" :offset="[4, 4]">Keluar</q-tooltip>
          </q-btn>
        </div>
      </div>
    </q-scroll-area>
  </q-drawer>
</template>

<script setup>
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAuthStore } from 'src/stores/auth.store'
import { getUserInitial } from 'src/utils/formatter'
import { useQuasar } from 'quasar'

defineProps({ modelValue: { type: Boolean, default: true } })
defineEmits(['update:modelValue'])

const route = useRoute()
const authStore = useAuthStore()
const $q = useQuasar()

// ── White transparent logo for dark sidebar (light & dark mode) ──
const logoSrc = computed(() => '/images/kafeinarts-logo-light.png')

const scMasterItems = [
  {
    route: 'sc-master-vendor',
    label: 'Vendor / Supplier',
    icon: 'storefront',
    iconColor: 'pink-7',
    iconDarkColor: 'pink-3',
  },
  {
    route: 'sc-master-kategori',
    label: 'Kategori Item / Aset',
    icon: 'category',
    iconColor: 'teal-7',
    iconDarkColor: 'teal-3',
  },
  {
    route: 'sc-master-katalog',
    label: 'Katalog Item',
    icon: 'inventory',
    iconColor: 'indigo-7',
    iconDarkColor: 'indigo-3',
  },
  {
    route: 'sc-master-lokasi',
    label: 'Lokasi / Gudang',
    icon: 'warehouse',
    iconColor: 'brown-6',
    iconDarkColor: 'brown-3',
  },
]
const scTransaksiItems = [
  {
    route: 'sc-transaksi-pr',
    label: 'Purchase Request',
    icon: 'assignment',
    iconColor: 'deep-orange-7',
    iconDarkColor: 'deep-orange-3',
  },
  {
    route: 'sc-transaksi-po',
    label: 'Purchase Order',
    icon: 'shopping_cart',
    iconColor: 'cyan-7',
    iconDarkColor: 'cyan-3',
  },
  {
    route: 'sc-transaksi-gr',
    label: 'Penerimaan Barang',
    icon: 'move_to_inbox',
    iconColor: 'light-green-8',
    iconDarkColor: 'light-green-4',
  },
  {
    route: 'sc-transaksi-aset',
    label: 'Manajemen Aset',
    icon: 'devices_other',
    iconColor: 'blue-grey-7',
    iconDarkColor: 'blue-grey-3',
  },
  {
    route: 'sc-transaksi-langganan',
    label: 'Perpanjangan Langganan',
    icon: 'autorenew',
    iconColor: 'deep-purple-7',
    iconDarkColor: 'deep-purple-3',
  },
]
const scLaporanItems = [
  {
    route: 'sc-laporan-po',
    label: 'Laporan PO',
    icon: 'summarize',
    iconColor: 'amber-9',
    iconDarkColor: 'amber-4',
  },
  {
    route: 'sc-laporan-aset',
    label: 'Aset per Karyawan',
    icon: 'badge',
    iconColor: 'lime-9',
    iconDarkColor: 'lime-4',
  },
  {
    route: 'sc-laporan-lisensi',
    label: 'Lisensi Expired',
    icon: 'event_busy',
    iconColor: 'red-7',
    iconDarkColor: 'red-3',
  },
  {
    route: 'sc-laporan-vendor',
    label: 'Riwayat Transaksi Vendor',
    icon: 'history',
    iconColor: 'grey-8',
    iconDarkColor: 'grey-4',
  },
]

const masterItems = [
  {
    route: 'finance-master-kategori',
    label: 'Kategori Transaksi',
    icon: 'category',
    iconColor: 'cyan-8',
    iconDarkColor: 'cyan-4',
  },
  {
    route: 'finance-master-rekening',
    label: 'Rekening / Kas Bank',
    icon: 'account_balance',
    iconColor: 'green-7',
    iconDarkColor: 'green-3',
  },
  {
    route: 'finance-master-pajak',
    label: 'Master Pajak',
    icon: 'receipt_long',
    iconColor: 'red-7',
    iconDarkColor: 'red-3',
  },
]
const transaksiItems = [
  {
    route: 'finance-jurnal-kas',
    label: 'Jurnal Kas',
    icon: 'menu_book',
    iconColor: 'purple-7',
    iconDarkColor: 'purple-3',
  },
  {
    route: 'finance-invoice',
    label: 'Invoice (AR)',
    icon: 'request_quote',
    iconColor: 'blue-7',
    iconDarkColor: 'blue-3',
  },
  {
    route: 'finance-tagihan-vendor',
    label: 'Tagihan Vendor (AP)',
    icon: 'receipt',
    iconColor: 'orange-7',
    iconDarkColor: 'orange-3',
  },
]
const laporanItems = [
  {
    route: 'finance-laporan-cashflow',
    label: 'Arus Kas',
    icon: 'trending_up',
    iconColor: 'teal-7',
    iconDarkColor: 'teal-3',
  },
  {
    route: 'finance-laporan-laba-rugi',
    label: 'Laba / Rugi',
    icon: 'bar_chart',
    iconColor: 'light-blue-7',
    iconDarkColor: 'light-blue-3',
  },
  {
    route: 'finance-laporan-ar-aging',
    label: 'Piutang (AR Aging)',
    icon: 'schedule',
    iconColor: 'brown-6',
    iconDarkColor: 'brown-3',
  },
]

const hrMasterItems = [
  {
    route: 'hr-master-departemen',
    label: 'Departemen',
    icon: 'apartment',
    iconColor: 'indigo-7',
    iconDarkColor: 'indigo-3',
  },
  {
    route: 'hr-master-jabatan',
    label: 'Jabatan / Posisi',
    icon: 'badge',
    iconColor: 'deep-purple-7',
    iconDarkColor: 'deep-purple-3',
  },
  {
    route: 'hr-master-karyawan',
    label: 'Data Karyawan',
    icon: 'person',
    iconColor: 'blue-7',
    iconDarkColor: 'blue-3',
  },
  {
    route: 'hr-master-jenis-cuti',
    label: 'Jenis Cuti / Izin',
    icon: 'event_busy',
    iconColor: 'pink-7',
    iconDarkColor: 'pink-3',
  },
  {
    route: 'hr-master-komponen-gaji',
    label: 'Komponen Gaji',
    icon: 'payments',
    iconColor: 'green-7',
    iconDarkColor: 'green-3',
  },
  {
    route: 'hr-master-jadwal-kerja',
    label: 'Jadwal Kerja / Shift',
    icon: 'schedule',
    iconColor: 'orange-7',
    iconDarkColor: 'orange-3',
  },
]
const hrTransaksiItems = [
  {
    route: 'hr-presensi',
    label: 'Presensi / Absensi',
    icon: 'fingerprint',
    iconColor: 'teal-7',
    iconDarkColor: 'teal-3',
  },
  {
    route: 'hr-pengajuan-cuti',
    label: 'Pengajuan Cuti',
    icon: 'assignment_turned_in',
    iconColor: 'deep-orange-7',
    iconDarkColor: 'deep-orange-3',
  },
  {
    route: 'hr-payroll',
    label: 'Penggajian (Payroll)',
    icon: 'account_balance_wallet',
    iconColor: 'cyan-8',
    iconDarkColor: 'cyan-4',
  },
  {
    route: 'hr-reimbursement',
    label: 'Reimbursement',
    icon: 'receipt_long',
    iconColor: 'lime-9',
    iconDarkColor: 'lime-4',
  },
  {
    route: 'hr-penilaian-kinerja',
    label: 'Penilaian Kinerja',
    icon: 'insights',
    iconColor: 'purple-7',
    iconDarkColor: 'purple-3',
  },
]
const hrLaporanItems = [
  {
    route: 'hr-laporan-absensi',
    label: 'Rekap Absensi',
    icon: 'summarize',
    iconColor: 'brown-6',
    iconDarkColor: 'brown-3',
  },
  {
    route: 'hr-laporan-slip-gaji',
    label: 'Slip Gaji',
    icon: 'description',
    iconColor: 'grey-8',
    iconDarkColor: 'grey-4',
  },
  {
    route: 'hr-laporan-saldo-cuti',
    label: 'Saldo Cuti',
    icon: 'event_available',
    iconColor: 'light-green-8',
    iconDarkColor: 'light-green-4',
  },
  {
    route: 'hr-laporan-turnover',
    label: 'Turnover Karyawan',
    icon: 'sync_alt',
    iconColor: 'red-7',
    iconDarkColor: 'red-3',
  },
]

const currentUser = computed(() => authStore.currentUser)
const userInitial = computed(() => getUserInitial(currentUser.value))
const isDashActive = computed(() => route?.name === 'dashboard')
const isInventoryActive = computed(() => (route?.name || '').startsWith('inventory-'))

const userRole = computed(() => (currentUser.value?.role || '').toLowerCase())

const canAccessInventory = computed(() => {
  const r = userRole.value
  return !r || ['superadmin', 'inventory', 'manager', 'barista', 'staff'].includes(r)
})

const canAccessFinance = computed(() => {
  const r = userRole.value
  return !r || ['superadmin', 'finance', 'manager'].includes(r)
})

const canAccessHr = computed(() => {
  const r = userRole.value
  return !r || ['superadmin', 'hr', 'manager', 'staff'].includes(r)
})

const canAccessSc = computed(() => {
  const r = userRole.value
  return !r || ['superadmin', 'supply_chain', 'sc', 'manager', 'barista', 'staff'].includes(r)
})

const canAccessUserManagement = computed(() => {
  const r = userRole.value
  return r === 'superadmin' || !r
})

function formatRoleLabel(role) {
  const r = String(role || '').toLowerCase()
  const map = {
    superadmin: 'Super Admin',
    inventory: 'Inventaris',
    finance: 'Finance',
    hr: 'Human Resource',
    supply_chain: 'Supply Chain',
    sc: 'Supply Chain',
    manager: 'Manager',
    barista: 'Barista',
    staff: 'Staff',
  }
  return map[r] || role || 'Super Admin'
}

async function onLogout() {
  await authStore.logout()
  window.location.href = '/login'
}
</script>

<!-- ================================================================
  NON-SCOPED: theme tokens + override Quasar-rendered DOM outside
  the scoped boundary.

  Theme strategy
  --------------
  All colors are CSS custom properties declared on `.app-sidebar`
  (the class already on <q-drawer>), so they cascade down through
  q-scroll-area → .sidebar-inner → every nav/sub item without
  touching `.sidebar-inner`'s own layout rules.

  Resolution order (last one that matches wins):
    1. Light palette (default)                         — matches the
       reference screenshot: white surface, slate text, soft blue
       highlight on the active row.
    2. `@media (prefers-color-scheme: dark)`            — automatic,
       follows the OS/browser when Quasar's Dark plugin is set to
       'auto' and hasn't forced a class yet.
    3. `body.body--dark .app-sidebar`                   — explicit
       override once Quasar's Dark plugin (manual or auto) adds/
       removes `body--dark` / `body--light` on <body>. This always
       wins over the media query so a manual toggle is respected
       even if it disagrees with the OS setting.
  ================================================================ -->
<style lang="scss">
// ── 1. Light theme tokens (default) ───────────────────────────────
.app-sidebar {
  --sb-bg: #ffffff;
  --sb-header-from: #2f6fed;
  --sb-header-to: #1d4ed8;
  --sb-header-text: #ffffff;
  --sb-header-sub-text: rgba(255, 255, 255, 0.72);
  --sb-text: #4b5563;
  --sb-hover: #f1f5f9;
  --sb-border: #e7eaf0;
  --sb-active-bg: #eaf0fd;
  --sb-active-text: #2f6fed;
  --sb-active-bar: #2f6fed;
  --sb-sub-bg: #f8fafc;
  --sb-sub-text: #64748b;
  --sb-sub-lbl: #9aa5b1;
  --sb-icon-default: #94a3b8;
  --sb-scrollbar-thumb: rgba(15, 23, 42, 0.16);
  --sb-footer-bg: #f8fafc;
  --sb-shadow: rgba(15, 23, 42, 0.08);

  // Colour-coded module icons (constant across themes, tuned per-theme below)
  --sb-accent-dashboard: #2f6fed;
  --sb-accent-inventory: #f59e0b;
  --sb-accent-finance: #10b981;
  --sb-accent-hr: #8b5cf6;
  --sb-accent-sc: #06b6d4;
  --sb-accent-user: #f43f5e;
}

// ── 2. System dark-mode preference (auto) ─────────────────────────
@media (prefers-color-scheme: dark) {
  .app-sidebar:not(.body--light) {
    --sb-bg: #16213e;
    --sb-header-from: #1e40af;
    --sb-header-to: #16213e;
    --sb-header-text: #ffffff;
    --sb-header-sub-text: rgba(255, 255, 255, 0.65);
    --sb-text: rgba(255, 255, 255, 0.75);
    --sb-hover: rgba(255, 255, 255, 0.08);
    --sb-border: rgba(255, 255, 255, 0.1);
    --sb-active-bg: rgba(255, 255, 255, 0.12);
    --sb-active-text: #ffffff;
    --sb-active-bar: #60a5fa;
    --sb-sub-bg: rgba(0, 0, 0, 0.22);
    --sb-sub-text: rgba(255, 255, 255, 0.55);
    --sb-sub-lbl: rgba(255, 255, 255, 0.28);
    --sb-icon-default: rgba(255, 255, 255, 0.58);
    --sb-scrollbar-thumb: rgba(255, 255, 255, 0.18);
    --sb-footer-bg: rgba(0, 0, 0, 0.18);
    --sb-shadow: rgba(0, 0, 0, 0.35);

    --sb-accent-dashboard: #60a5fa;
    --sb-accent-inventory: #fbbf24;
    --sb-accent-finance: #34d399;
    --sb-accent-hr: #a78bfa;
    --sb-accent-sc: #22d3ee;
    --sb-accent-user: #fb7185;
  }
}

// ── 3. Explicit Quasar Dark-plugin override (manual or auto) ─────
body.body--dark .app-sidebar {
  --sb-bg: #16213e;
  --sb-header-from: #1e40af;
  --sb-header-to: #16213e;
  --sb-header-text: #ffffff;
  --sb-header-sub-text: rgba(255, 255, 255, 0.65);
  --sb-text: rgba(255, 255, 255, 0.75);
  --sb-hover: rgba(255, 255, 255, 0.08);
  --sb-border: rgba(255, 255, 255, 0.1);
  --sb-active-bg: rgba(255, 255, 255, 0.12);
  --sb-active-text: #ffffff;
  --sb-active-bar: #60a5fa;
  --sb-sub-bg: rgba(0, 0, 0, 0.22);
  --sb-sub-text: rgba(255, 255, 255, 0.55);
  --sb-sub-lbl: rgba(255, 255, 255, 0.28);
  --sb-icon-default: rgba(255, 255, 255, 0.58);
  --sb-scrollbar-thumb: rgba(255, 255, 255, 0.18);
  --sb-footer-bg: rgba(0, 0, 0, 0.18);
  --sb-shadow: rgba(0, 0, 0, 0.35);

  --sb-accent-dashboard: #60a5fa;
  --sb-accent-inventory: #fbbf24;
  --sb-accent-finance: #34d399;
  --sb-accent-hr: #a78bfa;
  --sb-accent-sc: #22d3ee;
  --sb-accent-user: #fb7185;
}

body.body--light .app-sidebar {
  // forces light tokens back on even if the OS prefers dark
  --sb-bg: #ffffff;
  --sb-header-from: #2f6fed;
  --sb-header-to: #1d4ed8;
  --sb-header-text: #ffffff;
  --sb-header-sub-text: rgba(255, 255, 255, 0.72);
  --sb-text: #4b5563;
  --sb-hover: #f1f5f9;
  --sb-border: #e7eaf0;
  --sb-active-bg: #eaf0fd;
  --sb-active-text: #2f6fed;
  --sb-active-bar: #2f6fed;
  --sb-sub-bg: #f8fafc;
  --sb-sub-text: #64748b;
  --sb-sub-lbl: #9aa5b1;
  --sb-icon-default: #94a3b8;
  --sb-scrollbar-thumb: rgba(15, 23, 42, 0.16);
  --sb-footer-bg: #f8fafc;
  --sb-shadow: rgba(15, 23, 42, 0.08);

  --sb-accent-dashboard: #2f6fed;
  --sb-accent-inventory: #f59e0b;
  --sb-accent-finance: #10b981;
  --sb-accent-hr: #8b5cf6;
  --sb-accent-sc: #06b6d4;
  --sb-accent-user: #f43f5e;
}

// ── Drawer shell ─────────────────────────────────────────────────
.app-sidebar.q-drawer,
.app-sidebar .q-drawer {
  background: var(--sb-bg) !important;
  border-right: 1px solid var(--sb-border) !important;
  box-shadow: 2px 0 20px var(--sb-shadow) !important;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease;
}

// ── q-scroll-area thumb colour ────────────────────────────────────
.app-sidebar .q-scrollarea__thumb {
  background: var(--sb-scrollbar-thumb) !important;
  border-radius: 3px !important;
}

// ── Quasar q-item default overrides inside sidebar ────────────────
.app-sidebar .q-item {
  min-height: unset;
  padding: 0;
  color: var(--sb-text);
}

// Remove default q-item hover ripple bg so we control it ourselves
.app-sidebar .q-item:before {
  opacity: 0 !important;
}
.app-sidebar .q-item__section--avatar {
  min-width: unset;
  padding-right: 0;
}

// ── q-expansion-item header overrides ────────────────────────────
.app-sidebar .q-expansion-item > .q-expansion-item__container > .q-item {
  padding: 0;
  color: var(--sb-text);
  min-height: unset;
}
.app-sidebar .q-expansion-item > .q-expansion-item__container > .q-item:before {
  opacity: 0 !important;
}
.app-sidebar .q-expansion-item__content {
  padding: 0;
}

// ── q-avatar inside sidebar footer ────────────────────────────────
.app-sidebar .sidebar-footer__avatar.q-avatar {
  background: var(--sb-active-bg) !important;
  color: var(--sb-active-text) !important;
  font-weight: 700 !important;
  flex-shrink: 0;
}

// ── q-btn logout inside sidebar footer ───────────────────────────
.app-sidebar .sidebar-footer__logout.q-btn {
  color: var(--sb-icon-default) !important;
  &:hover {
    background: rgba(244, 63, 94, 0.14) !important;
    color: #f43f5e !important;
  }
}

// ── Colour-coded module icons (positional, matches the reference
//    screenshot's per-menu icon colours). Falls back gracefully to
//    the neutral --sb-icon-default when a role hides a section and
//    the ordering shifts. ───────────────────────────────────────
.app-sidebar .nav-item--dashboard {
  --module-accent: var(--sb-accent-dashboard);
}
.app-sidebar .nav-item--inventory {
  --module-accent: var(--sb-accent-inventory);
}
.app-sidebar .nav-item--finance {
  --module-accent: var(--sb-accent-finance);
}
.app-sidebar .nav-item--hr {
  --module-accent: var(--sb-accent-hr);
}
.app-sidebar .nav-item--sc {
  --module-accent: var(--sb-accent-sc);
}
.app-sidebar .nav-item--user {
  --module-accent: var(--sb-accent-user);
}

// `--module-accent` is declared above on each top-level item/section and
// inherits (cascades) into everything below it — including its own
// `.sub-panel` — so every dropdown automatically "inherits" its parent's
// theme colour without any extra markup or JS.
.app-sidebar .nav-item__icon-section {
  color: var(--module-accent, var(--sb-icon-default));
}
</style>

<!-- ================================================================
  SCOPED: sidebar internals — layout/structure untouched, only
  colors now read from the tokens declared on `.app-sidebar` above.
  `.sidebar-inner`'s own layout rules (flex/min-height) are unchanged.
  ================================================================ -->
<style lang="scss" scoped>
// ── Layout wrapper (unchanged structure, background now themed) ──
.sidebar-inner {
  display: flex;
  flex-direction: column;
  flex: 1;
  min-height: 100%;
  background: var(--sb-bg);
  transition: background-color 0.2s ease;
}

// ── Scroll area fills remaining space ────────────────────────────
.sidebar-scroll {
  height: 100%;
  width: 100%;
}

// ── Logo / header banner ─────────────────────────────────────────
.sidebar-logo {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 10px;
  padding: 24px 16px 20px;
  text-align: center;
  background: linear-gradient(160deg, var(--sb-header-from), var(--sb-header-to));
  border-bottom: 1px solid var(--sb-border);
  flex-shrink: 0;
  transition: background 0.2s ease;
}

.sidebar-logo__icon-wrap {
  width: 44px;
  height: 44px;
  min-width: 44px;
  border-radius: 12px;
  background: rgba(255, 255, 255, 0.16);
  border: 1px solid rgba(255, 255, 255, 0.28);
  box-shadow: 0 4px 14px rgba(0, 0, 0, 0.18);
  padding: 5px;
  overflow: hidden;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  transition:
    transform 0.2s ease,
    box-shadow 0.2s ease;
  &:hover {
    transform: scale(1.04);
  }
}

.sidebar-logo__img {
  width: 100%;
  height: 100%;
  object-fit: contain;
  display: block;
  filter: drop-shadow(0 2px 4px rgba(0, 0, 0, 0.2));
}
.sidebar-logo__brand {
  font-size: 13px;
  font-weight: 800;
  color: var(--sb-header-text);
  letter-spacing: 0.6px;
  line-height: 1.3;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

// ── q-list wrapper ────────────────────────────────────────────────
.sidebar-nav {
  flex: 1;
  padding: 8px 0 !important; // override q-list default
}

// ── Main nav item (q-item) ────────────────────────────────────────
.nav-item {
  border-radius: 8px;
  margin: 2px 10px;
  padding: 9px 10px !important;
  font-size: 13px;
  font-weight: 500;
  color: var(--sb-text) !important;
  transition:
    background 0.15s ease,
    color 0.15s ease;
  position: relative;

  &:hover {
    background: var(--sb-hover) !important;
    color: var(--module-accent, var(--sb-active-text)) !important;
    .nav-item__icon-section {
      opacity: 1;
    }
  }

  // Active state: left accent bar + highlight
  &.nav-item--active,
  &:deep(.q-router-link--active) {
    background: var(--sb-active-bg) !important;
    color: var(--module-accent, var(--sb-active-text)) !important;
    font-weight: 600;

    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 20%;
      height: 60%;
      width: 3px;
      border-radius: 0 3px 3px 0;
      background: var(--module-accent, var(--sb-active-bar));
    }

    .nav-item__icon-section {
      opacity: 1 !important;
      color: var(--module-accent, var(--sb-active-text)) !important;
    }
  }
}

.nav-item__icon-section {
  opacity: 0.9;
  min-width: 32px !important;
  padding-right: 0 !important;
  transition:
    opacity 0.15s ease,
    color 0.15s ease;
  display: flex;
  align-items: center;
  justify-content: flex-start;
}

.nav-item__label-section {
  font-size: 13px;
  font-weight: inherit;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar-label {
  padding: 12px 14px 4px 18px;
  font-size: 10px;
  font-weight: 700;
  letter-spacing: 0.85px;
  text-transform: uppercase;
  color: var(--sb-sub-lbl);
  margin: 0;
}

// ── Footer ────────────────────────────────────────────────────────
.sidebar-footer {
  display: flex;
  align-items: center;
  gap: 10px;
  padding: 12px 14px;
  border-top: 1px solid var(--sb-border);
  background: var(--sb-footer-bg);
  flex-shrink: 0;
  transition:
    background-color 0.2s ease,
    border-color 0.2s ease;
}

.sidebar-footer__info {
  flex: 1;
  min-width: 0;
  display: flex;
  flex-direction: column;
  gap: 1px;
}

.sidebar-footer__name {
  font-size: 12.5px;
  font-weight: 700;
  color: var(--sb-text);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.sidebar-footer__role {
  font-size: 10.5px;
  color: var(--sb-sub-lbl);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
</style>

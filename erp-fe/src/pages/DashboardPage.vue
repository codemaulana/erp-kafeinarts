<template>
  <q-page class="app-page dash clay-animate-up">
    <!-- ── Page header ── -->
    <header class="dash__header">
      <div>
        <h1 class="dash__title">Dashboard</h1>
        <p class="dash__sub">
          Selamat datang kembali, <strong>{{ userName }}</strong> 👋
        </p>
      </div>
      <button class="dash__reload-btn" type="button" :disabled="loading" @click="reloadData">
        <q-spinner-dots v-if="loading" size="14px" color="primary" />
        <svg v-else viewBox="0 0 18 18" fill="none" width="14" height="14">
          <path
            d="M15.5 9A6.5 6.5 0 1 1 9 2.5a6.46 6.46 0 0 1 4.6 1.9"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
          />
          <path
            d="M13 1v4h4"
            stroke="currentColor"
            stroke-width="1.6"
            stroke-linecap="round"
            stroke-linejoin="round"
          />
        </svg>
        {{ loading ? 'Memuat...' : 'Reload Data' }}
      </button>
    </header>

    <!-- ── Stat cards ── -->
    <section aria-label="General Report">
      <h2 class="dash__section-title">Ringkasan Operasional & Keuangan</h2>
      <div class="row q-col-gutter-md">
        <!-- 1. Pemasukan -->
        <div class="col-12 col-sm-6 col-md-3">
          <q-card class="stat-card bg-red text-white full-height column">
            <q-card-section>
              <div class="stat-title">Total Pemasukan</div>
              <div class="text-subtitle2 text-white">Kas Masuk</div>
            </q-card-section>

            <q-card-section class="q-pt-none col-grow">
              <q-spinner-dots v-if="loading" size="24px" color="white" />
              <div v-else class="stat-value">{{ formatCurrency(stats.totalInflow) }}</div>
            </q-card-section>

            <q-separator dark />

            <q-card-actions>
              <q-btn
                flat
                text-color="white"
                icon="trending_up"
                label="Detail"
                :to="{ name: 'finance-jurnal-kas' }"
              />
            </q-card-actions>
          </q-card>
        </div>

        <!-- 2. Pengeluaran -->
        <div class="col-12 col-sm-6 col-md-3">
          <q-card class="stat-card bg-blue text-white full-height column">
            <q-card-section>
              <div class="stat-title">Total Pengeluaran</div>
              <div class="text-subtitle2 text-white">Kas Keluar</div>
            </q-card-section>

            <q-card-section class="q-pt-none col-grow">
              <q-spinner-dots v-if="loading" size="24px" color="white" />
              <div v-else class="stat-value">{{ formatCurrency(stats.totalOutflow) }}</div>
            </q-card-section>

            <q-separator dark />

            <q-card-actions>
              <q-btn
                flat
                text-color="white"
                icon="trending_down"
                label="Riwayat"
                :to="{ name: 'finance-jurnal-kas' }"
              />
            </q-card-actions>
          </q-card>
        </div>

        <!-- 3. Saldo Kas Bersih -->
        <div class="col-12 col-sm-6 col-md-3">
          <q-card class="stat-card bg-green text-white full-height column">
            <q-card-section>
              <div class="stat-title">Saldo Kas Bersih</div>
              <div class="text-subtitle2 text-white">
                {{ stats.netCashFlow >= 0 ? 'Surplus' : 'Defisit' }}
              </div>
            </q-card-section>

            <q-card-section class="q-pt-none col-grow">
              <q-spinner-dots v-if="loading" size="24px" color="white" />
              <div v-else class="stat-value">{{ formatCurrency(stats.netCashFlow) }}</div>
            </q-card-section>

            <q-separator dark />

            <q-card-actions>
              <q-btn
                flat
                text-color="white"
                icon="account_balance_wallet"
                label="Laporan"
                :to="{ name: 'finance-laporan-cashflow' }"
              />
            </q-card-actions>
          </q-card>
        </div>

        <!-- 4. Produk Terdaftar -->
        <div class="col-12 col-sm-6 col-md-3">
          <q-card class="stat-card bg-purple text-white full-height column">
            <q-card-section>
              <div class="row items-center no-wrap">
                <div class="col">
                  <div class="stat-title">Produk Katalog</div>
                  <div class="text-subtitle2 text-white">Inventaris</div>
                </div>
                <div class="col-auto">
                  <q-btn text-color="white" round flat icon="more_vert" size="sm">
                    <q-menu cover auto-close>
                      <q-list>
                        <q-item clickable :to="{ name: 'inventory-products' }">
                          <q-item-section>Kelola Inventaris</q-item-section>
                        </q-item>
                      </q-list>
                    </q-menu>
                  </q-btn>
                </div>
              </div>
            </q-card-section>

            <q-card-section class="q-pt-none col-grow">
              <q-spinner-dots v-if="loading" size="24px" color="white" />
              <div v-else class="stat-value">
                {{ stats.totalProducts }} <span class="text-caption">Produk</span>
              </div>
            </q-card-section>

            <q-separator dark />

            <q-card-actions>
              <q-btn
                flat
                text-color="white"
                icon="inventory_2"
                label="Lihat Produk"
                :to="{ name: 'inventory-products' }"
              />
            </q-card-actions>
          </q-card>
        </div>
      </div>
    </section>

    <!-- ── Finance modules quick-access ── -->
    <section aria-label="Finance Modules">
      <div class="dash-modules-header">
        <h2 class="dash__section-title">Modul Finance</h2>
        <span class="dash-modules-header__sub">Akses cepat ke semua fitur keuangan</span>
      </div>
      <div class="dash-modules">
        <router-link
          v-for="mod in financeModules"
          :key="mod.name"
          :to="{ name: mod.route }"
          class="dash-mod-card"
        >
          <div
            class="dash-mod-card__icon"
            :style="{ background: mod.iconBg, color: mod.iconColor }"
          >
            <q-icon :name="mod.icon" size="22px" />
          </div>
          <div class="dash-mod-card__body">
            <div class="dash-mod-card__name">{{ mod.name }}</div>
            <div class="dash-mod-card__desc">{{ mod.desc }}</div>
          </div>
          <svg class="dash-mod-card__arrow" viewBox="0 0 12 12" fill="none" width="12" height="12">
            <path
              d="M4.5 3L7.5 6L4.5 9"
              stroke="currentColor"
              stroke-width="1.4"
              stroke-linecap="round"
              stroke-linejoin="round"
            />
          </svg>
        </router-link>
      </div>
    </section>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { api } from 'boot/axios'
import { useAuthStore } from 'src/stores/auth.store'
import { useNotify } from 'src/composables/useNotify'
import { formatCurrency } from 'src/utils/formatter'

const authStore = useAuthStore()
const { notifySuccess } = useNotify()

const userName = computed(() => authStore.currentUser?.name || 'Admin')

const loading = ref(false)
const stats = reactive({
  totalInflow: 0,
  totalOutflow: 0,
  netCashFlow: 0,
  totalProducts: 0,
  totalEmployees: 0,
})

const financeModules = [
  {
    name: 'Kategori Transaksi',
    route: 'finance-master-kategori',
    icon: 'category',
    iconBg: '#EFF6FF',
    iconColor: '#3B82F6',
    desc: 'Master klasifikasi transaksi',
  },
  {
    name: 'Rekening / Kas',
    route: 'finance-master-rekening',
    icon: 'account_balance',
    iconBg: '#F0FDF4',
    iconColor: '#10B981',
    desc: 'Kelola akun rekening perusahaan',
  },
  {
    name: 'Master Pajak',
    route: 'finance-master-pajak',
    icon: 'percent',
    iconBg: '#FFFBEB',
    iconColor: '#F59E0B',
    desc: 'Konfigurasi tarif pajak',
  },
  {
    name: 'Jurnal Kas',
    route: 'finance-jurnal-kas',
    icon: 'receipt_long',
    iconBg: '#FAF5FF',
    iconColor: '#A855F7',
    desc: 'Catat arus kas masuk & keluar',
  },
  {
    name: 'Invoice (AR)',
    route: 'finance-invoice',
    icon: 'description',
    iconBg: '#F0FDF4',
    iconColor: '#10B981',
    desc: 'Manajemen tagihan klien',
  },
  {
    name: 'Tagihan Vendor (AP)',
    route: 'finance-tagihan-vendor',
    icon: 'receipt',
    iconBg: '#FEF2F2',
    iconColor: '#EF4444',
    desc: 'Kelola tagihan dari vendor',
  },
  {
    name: 'Laporan Arus Kas',
    route: 'finance-laporan-cashflow',
    icon: 'trending_up',
    iconBg: '#EFF6FF',
    iconColor: '#3B82F6',
    desc: 'Analisis cash flow perusahaan',
  },
  {
    name: 'Laporan Laba/Rugi',
    route: 'finance-laporan-laba-rugi',
    icon: 'bar_chart',
    iconBg: '#F0FDF4',
    iconColor: '#10B981',
    desc: 'Ringkasan pendapatan & biaya',
  },
  {
    name: 'Piutang (AR Aging)',
    route: 'finance-laporan-ar-aging',
    icon: 'schedule',
    iconBg: '#FFFBEB',
    iconColor: '#F59E0B',
    desc: 'Status penuaan piutang klien',
  },
]

async function loadDashboardData() {
  loading.value = true
  try {
    const [cashFlowRes, productRes, empRes] = await Promise.allSettled([
      api.get('/finance/laporan/cash-flow'),
      api.get('/inventory/products?per_page=1'),
      api.get('/hr/employees?per_page=1'),
    ])

    if (cashFlowRes.status === 'fulfilled') {
      const d = cashFlowRes.value.data?.data || {}
      stats.totalInflow = d.total_inflow || d.summary?.totalMasuk || 0
      stats.totalOutflow = d.total_outflow || d.summary?.totalKeluar || 0
      stats.netCashFlow = d.net_cash_flow || d.summary?.saldoBersih || 0
    }

    if (productRes.status === 'fulfilled') {
      const meta = productRes.value.data?.meta || {}
      stats.totalProducts = meta.total || (productRes.value.data?.data || []).length || 0
    }

    if (empRes.status === 'fulfilled') {
      const meta = empRes.value.data?.meta || {}
      stats.totalEmployees = meta.total || (empRes.value.data?.data || []).length || 0
    }
  } catch (err) {
    console.error('Gagal memuat ringkasan dashboard:', err)
  } finally {
    loading.value = false
  }
}

async function reloadData() {
  await loadDashboardData()
  notifySuccess('Data dashboard berhasil diperbarui.')
}

onMounted(() => {
  loadDashboardData()
})
</script>

<style lang="scss" scoped>
.dash {
  display: flex;
  flex-direction: column;
  gap: 24px;
}

.dash__header {
  display: flex;
  align-items: flex-start;
  justify-content: space-between;
  gap: 12px;
  flex-wrap: wrap;
}

.dash__title {
  font-size: clamp(18px, 2vw, 24px);
  font-weight: 800;
  color: var(--rb-text-primary);
  margin: 0 0 3px;
  line-height: 1.2;
}

.dash__sub {
  font-size: 13.5px;
  color: var(--rb-text-muted);
  margin: 0;
  strong {
    color: var(--rb-text-primary);
    font-weight: 600;
  }
}

.dash__section-title {
  font-size: 15px;
  font-weight: 700;
  color: var(--rb-text-primary);
  margin: 0 0 12px;
}

.dash__reload-btn {
  display: flex;
  align-items: center;
  gap: 6px;
  padding: 7px 14px;
  border-radius: 8px;
  border: 1px solid var(--rb-border);
  background: var(--rb-surface);
  color: var(--rb-text-muted);
  font-size: 12.5px;
  font-family: inherit;
  font-weight: 500;
  cursor: pointer;
  white-space: nowrap;
  flex-shrink: 0;
  transition:
    background 0.15s ease,
    color 0.15s ease;
  &:hover {
    background: var(--rb-bg-2);
    color: var(--rb-primary-light);
  }
}

.dash-modules-header {
  display: flex;
  align-items: baseline;
  gap: 12px;
  flex-wrap: wrap;
}
.dash-modules-header__sub {
  font-size: 12.5px;
  color: var(--rb-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.dash-modules {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(240px, 1fr));
  gap: 10px;
}
.dash-mod-card {
  display: flex;
  align-items: center;
  gap: 13px;
  padding: 14px 16px;
  background: var(--rb-surface);
  border-radius: 12px;
  border: 1px solid var(--rb-border);
  box-shadow: var(--rb-shadow-xs);
  text-decoration: none;
  min-width: 0;
  overflow: hidden;
  transition:
    box-shadow 0.18s ease,
    transform 0.18s ease,
    border-color 0.18s ease;
  &:hover {
    box-shadow: var(--rb-shadow-sm);
    transform: translateY(-1px);
    border-color: var(--rb-primary-light);
    .dash-mod-card__arrow {
      opacity: 1;
      color: var(--rb-primary-light);
    }
  }
}
.dash-mod-card__icon {
  width: 42px;
  height: 42px;
  min-width: 42px;
  border-radius: 10px;
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
}
.dash-mod-card__body {
  flex: 1;
  min-width: 0;
}
.dash-mod-card__name {
  font-size: 13.5px;
  font-weight: 600;
  color: var(--rb-text-primary);
  margin-bottom: 2px;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.dash-mod-card__desc {
  font-size: 11.5px;
  color: var(--rb-text-muted);
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}
.dash-mod-card__arrow {
  color: var(--rb-text-light);
  flex-shrink: 0;
  opacity: 0;
  transition:
    opacity 0.18s ease,
    color 0.18s ease;
}

.stat-card {
  box-shadow: 0 4px 16px rgba(0, 0, 0, 0.3);
  transition:
    transform 0.25s ease,
    box-shadow 0.25s ease;
}

.stat-card:hover {
  transform: translateY(-3px);
  box-shadow: 0 8px 24px rgba(0, 0, 0, 0.4);
}

/* Shadow halus pada teks agar lebih terbaca */
.stat-title,
.stat-value {
  text-shadow: 0 1px 3px rgba(0, 0, 0, 0.35);
}
</style>

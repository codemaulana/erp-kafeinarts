<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Laba / Rugi"
      subtitle="Ringkasan pendapatan dan pengeluaran berdasarkan bulan & tahun"
      icon="bar_chart"
      icon-bg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      icon-color="#4CAF8A"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          label="Export PDF"
          icon="picture_as_pdf"
          class="clay-btn"
          color="red-5"
          @click="onExport"
        />
      </template>
    </AppPageHeader>

    <!-- Filter -->
    <div class="clay-card q-mb-lg">
      <div class="filter-title">
        <q-icon name="filter_list" size="18px" color="green-5" />
        <span>Filter Laporan</span>
      </div>
      <div class="filter-grid">
        <div class="form-field-group">
          <label class="form-label">Bulan</label>
          <q-select
            v-model="filter.bulan"
            :options="bulanOptions"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            label="Pilih Bulan"
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Tahun</label>
          <q-select
            v-model="filter.tahun"
            :options="tahunOptions"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            label="Pilih Tahun"
          />
        </div>
        <div class="form-field-group form-field-group--btn">
          <q-btn
            unelevated
            no-caps
            label="Tampilkan"
            icon="search"
            class="clay-btn clay-btn-primary"
            :loading="financeStore.isLoading('laporan')"
            @click="loadReport"
          />
          <q-btn flat no-caps label="Reset" class="clay-btn" @click="resetFilter" />
        </div>
      </div>
    </div>

    <template v-if="financeStore.labaRugiData">
      <!-- Summary Banner -->
      <div
        class="lr-banner q-mb-lg"
        :class="
          (financeStore.labaRugiData.summary?.labaKotor || 0) >= 0 ? 'lr-banner--laba' : 'lr-banner--rugi'
        "
      >
        <div class="lr-banner__period">
          {{ bulanLabel(financeStore.labaRugiData.summary?.bulan) }}
          {{ financeStore.labaRugiData.summary?.tahun }}
        </div>
        <div class="lr-banner__label">
          {{ (financeStore.labaRugiData.summary?.labaKotor || 0) >= 0 ? 'Laba Bersih' : 'Kerugian Bersih' }}
        </div>
        <div class="lr-banner__value">
          {{ formatCurrency(Math.abs(financeStore.labaRugiData.summary?.labaKotor || 0)) }}
        </div>
        <div class="lr-banner__meta">
          <span
            >Pendapatan:
            {{ formatCurrency(financeStore.labaRugiData.summary?.totalPendapatan || 0) }}</span
          >
          <span class="lr-banner__sep">•</span>
          <span
            >Pengeluaran:
            {{ formatCurrency(financeStore.labaRugiData.summary?.totalPengeluaran || 0) }}</span
          >
        </div>
      </div>

      <!-- Rekap Table -->
      <div class="clay-card clay-table-wrap">
        <div class="report-table-header">
          <span class="report-table-title">Rekap Laba / Rugi</span>
        </div>
        <q-table
          :rows="rekapRows"
          :columns="rekapColumns"
          row-key="kategori"
          flat
          class="clay-table"
          hide-pagination
          :rows-per-page-options="[0]"
        >
          <template #body-cell-tipe="props">
            <q-td :props="props">
              <AppBadge :variant="props.value === 'Pendapatan' ? 'success' : 'danger'">{{
                props.value
              }}</AppBadge>
            </q-td>
          </template>
          <template #body-cell-total="props">
            <q-td :props="props">
              <span
                :style="{
                  fontWeight: 700,
                  color:
                    props.row.tipe === 'Pendapatan' ? 'var(--clay-success)' : 'var(--clay-danger)',
                }"
              >
                {{ formatCurrency(props.value) }}
              </span>
            </q-td>
          </template>
          <template #body-cell-pct="props">
            <q-td :props="props">
              <div class="pct-bar-wrap">
                <div
                  class="pct-bar"
                  :class="props.row.tipe === 'Pendapatan' ? 'pct-bar--income' : 'pct-bar--expense'"
                  :style="{ width: props.value + '%' }"
                />
                <span class="pct-label">{{ props.value }}%</span>
              </div>
            </q-td>
          </template>
        </q-table>
      </div>
    </template>

    <div v-else class="clay-card lr-empty">
      <q-icon name="bar_chart" size="56px" color="grey-4" />
      <div class="lr-empty__text">
        Pilih bulan &amp; tahun, lalu klik <strong>"Tampilkan"</strong> untuk melihat laporan
        laba/rugi.
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { reactive, computed, onMounted } from 'vue'
import { useFinanceStore } from 'src/stores/finance.store'
import { useNotify } from 'src/composables/useNotify'
import { formatCurrency, bulanLabel } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const financeStore = useFinanceStore()
const { notifyInfo } = useNotify()

const now = new Date()
const filter = reactive({ bulan: now.getMonth() + 1, tahun: now.getFullYear() })

const bulanOptions = [
  { label: 'Januari', value: 1 },
  { label: 'Februari', value: 2 },
  { label: 'Maret', value: 3 },
  { label: 'April', value: 4 },
  { label: 'Mei', value: 5 },
  { label: 'Juni', value: 6 },
  { label: 'Juli', value: 7 },
  { label: 'Agustus', value: 8 },
  { label: 'September', value: 9 },
  { label: 'Oktober', value: 10 },
  { label: 'November', value: 11 },
  { label: 'Desember', value: 12 },
]
const tahunOptions = [2024, 2025, 2026, 2027].map((y) => ({ label: String(y), value: y }))

const rekapColumns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center', style: 'width:50px' },
  { name: 'tipe', label: 'Tipe', field: 'tipe', align: 'left' },
  { name: 'kategori', label: 'Kategori', field: 'kategori', align: 'left' },
  { name: 'total', label: 'Total', field: 'total', align: 'right', sortable: true },
  { name: 'pct', label: 'Proporsi', field: 'pct', align: 'left', style: 'min-width:120px' },
]

const rekapRows = computed(() => {
  const data = financeStore.labaRugiData
  if (!data) return []
  const summary = data.summary || {}
  const grandTotal = (summary.totalPendapatan || 0) + (summary.totalPengeluaran || 0)
  const pendapatanList = data.pendapatan || data.revenueDetails || []
  const pengeluaranList = data.pengeluaran || data.expenseDetails || []
  const pendapatan = pendapatanList.map((p, i) => ({
    no: i + 1,
    tipe: 'Pendapatan',
    kategori: p.kategori || p.category || '-',
    total: p.total || p.amount || 0,
    pct: grandTotal ? Math.round(((p.total || p.amount || 0) / grandTotal) * 100) : 0,
  }))
  const pengeluaran = pengeluaranList.map((p, i) => ({
    no: pendapatan.length + i + 1,
    tipe: 'Pengeluaran',
    kategori: p.kategori || p.category || '-',
    total: p.total || p.amount || 0,
    pct: grandTotal ? Math.round(((p.total || p.amount || 0) / grandTotal) * 100) : 0,
  }))
  return [...pendapatan, ...pengeluaran]
})

async function loadReport() {
  await financeStore.fetchLabaRugi({ bulan: filter.bulan, tahun: filter.tahun })
}

function resetFilter() {
  const n = new Date()
  filter.bulan = n.getMonth() + 1
  filter.tahun = n.getFullYear()
  loadReport()
}

function onExport() {
  notifyInfo('Fitur export PDF segera hadir.')
}

onMounted(() => loadReport())
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.lr-banner {
  border-radius: var(--rb-radius-lg);
  padding: 28px 32px;
  text-align: center;
  border: 1px solid var(--rb-border);
  box-shadow: var(--rb-shadow-sm);
  &--laba {
    background: var(--rb-success-bg);
    border-color: rgba(16, 185, 129, 0.2);
  }
  &--rugi {
    background: var(--rb-danger-bg);
    border-color: rgba(239, 68, 68, 0.2);
  }
  &__period {
    font-size: 12px;
    font-weight: 600;
    color: var(--rb-text-muted);
    text-transform: uppercase;
    letter-spacing: 1px;
    margin-bottom: 6px;
  }
  &__label {
    font-size: 14px;
    font-weight: 600;
    color: var(--rb-text-muted);
    margin-bottom: 4px;
  }
  &__value {
    font-size: 36px;
    font-weight: 900;
    margin-bottom: 10px;
    .lr-banner--laba & {
      color: var(--rb-success);
    }
    .lr-banner--rugi & {
      color: var(--rb-danger);
    }
  }
  &__meta {
    display: flex;
    justify-content: center;
    align-items: center;
    gap: 10px;
    font-size: 13px;
    color: var(--rb-text-muted);
  }
  &__sep {
    opacity: 0.4;
  }
}

.pct-bar-wrap {
  display: flex;
  align-items: center;
  gap: 8px;
  min-width: 120px;
}
.pct-bar {
  height: 6px;
  border-radius: 3px;
  min-width: 2px;
  transition: width 0.5s ease;
  &--income {
    background: var(--rb-success);
  }
  &--expense {
    background: var(--rb-danger);
  }
}
.pct-label {
  font-size: 11px;
  color: var(--rb-text-muted);
  white-space: nowrap;
}

.lr-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  gap: 14px;
  padding: 60px 24px;
  text-align: center;
  &__text {
    font-size: 14px;
    color: var(--rb-text-muted);
    max-width: 380px;
    line-height: 1.6;
  }
}
</style>

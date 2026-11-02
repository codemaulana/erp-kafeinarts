<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Arus Kas"
      subtitle="Analisis aliran kas masuk dan keluar berdasarkan periode"
      icon="trending_up"
      icon-bg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      icon-color="#5B9BD5"
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

    <!-- Filter Card -->
    <div class="clay-card q-mb-lg">
      <div class="filter-title">
        <q-icon name="filter_list" size="18px" color="blue-5" />
        <span>Filter Laporan</span>
      </div>
      <div class="filter-grid">
        <div class="form-field-group">
          <label class="form-label">Dari Tanggal</label>
          <q-input v-model="filter.periodeStart" outlined dense type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" class="clay-input" />
        </div>
        <div class="form-field-group">
          <label class="form-label">Sampai Tanggal</label>
          <q-input v-model="filter.periodeEnd" outlined dense type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" class="clay-input" />
        </div>
        <div class="form-field-group">
          <label class="form-label">Akun Bank / Kas</label>
          <q-select
            v-model="filter.akunBank"
            :options="financeStore.rekeningOptions"
            outlined
            dense
            clearable
            emit-value
            map-options
            class="clay-input"
            label="Pilih Akun Bank / Kas"
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

    <template v-if="financeStore.cashFlowData">
      <!-- Summary Cards -->
      <div class="cf-summary q-mb-lg">
        <AppStatCard
          icon="arrow_downward"
          icon-bg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
          icon-color="#4CAF8A"
          label="Total Kas Masuk"
          :value="formatCurrency(financeStore.cashFlowData.summary?.totalMasuk || 0)"
          value-color="var(--clay-success)"
        />
        <AppStatCard
          icon="arrow_upward"
          icon-bg="linear-gradient(135deg,#FEECF0,#F9BDC8)"
          icon-color="#E85D75"
          label="Total Kas Keluar"
          :value="formatCurrency(financeStore.cashFlowData.summary?.totalKeluar || 0)"
          value-color="var(--clay-danger)"
        />
        <AppStatCard
          icon="account_balance_wallet"
          :icon-bg="
            (financeStore.cashFlowData.summary?.saldoBersih || 0) >= 0
              ? 'linear-gradient(135deg,#E8F7F2,#BDE8D9)'
              : 'linear-gradient(135deg,#FEECF0,#F9BDC8)'
          "
          :icon-color="(financeStore.cashFlowData.summary?.saldoBersih || 0) >= 0 ? '#4CAF8A' : '#E85D75'"
          label="Saldo Bersih"
          :value="formatCurrency(financeStore.cashFlowData.summary?.saldoBersih || 0)"
          :value-color="
            (financeStore.cashFlowData.summary?.saldoBersih || 0) >= 0
              ? 'var(--clay-success)'
              : 'var(--clay-danger)'
          "
        />
      </div>

      <!-- Detail Table -->
      <div class="clay-card clay-table-wrap">
        <div class="report-table-header">
          <span class="report-table-title">Detail Transaksi</span>
          <span class="report-table-count"
            >{{ (financeStore.cashFlowData.rows || []).length }} transaksi</span
          >
        </div>
        <q-table
          :rows="financeStore.cashFlowData.rows || []"
          :columns="columns"
          row-key="id"
          flat
          class="clay-table"
          :pagination="{ rowsPerPage: 10 }"
          no-data-label="Tidak ada data untuk periode ini"
        >
          <template #body-cell-no="props">
            <q-td
              :props="props"
              class="text-center"
              style="width: 48px; font-weight: 600; color: var(--clay-text-muted)"
            >
              {{ props.rowIndex + 1 }}
            </q-td>
          </template>
          <template #body-cell-tipe="props">
            <q-td :props="props">
              <AppBadge :variant="props.value === 'Masuk' ? 'success' : 'danger'">
                <q-icon
                  :name="props.value === 'Masuk' ? 'arrow_downward' : 'arrow_upward'"
                  size="11px"
                />
                {{ props.value }}
              </AppBadge>
            </q-td>
          </template>
          <template #body-cell-nominal="props">
            <q-td :props="props">
              <span
                :style="{
                  fontWeight: 600,
                  color: props.row.tipe === 'Masuk' ? 'var(--clay-success)' : 'var(--clay-danger)',
                }"
              >
                {{ props.row.tipe === 'Masuk' ? '+' : '−' }} {{ formatCurrency(props.value) }}
              </span>
            </q-td>
          </template>
        </q-table>
      </div>
    </template>

    <!-- Empty State -->
    <div v-else class="clay-card cf-empty">
      <q-icon name="bar_chart" size="56px" color="grey-4" />
      <div class="cf-empty__text">
        Pilih periode dan klik <strong>"Tampilkan"</strong> untuk melihat data arus kas.
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { reactive, onMounted } from 'vue'
import { useFinanceStore } from 'src/stores/finance.store'
import { useNotify } from 'src/composables/useNotify'
import { formatCurrency } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const financeStore = useFinanceStore()
const { notifyInfo } = useNotify()

const now = new Date()
const filter = reactive({
  periodeStart: new Date(now.getFullYear(), now.getMonth(), 1).toISOString().slice(0, 10),
  periodeEnd: new Date(now.getFullYear(), now.getMonth() + 1, 0).toISOString().slice(0, 10),
  akunBank: null,
})

const columns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center' },
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left', sortable: true },
  { name: 'namaBank', label: 'Akun Bank', field: 'namaBank', align: 'left' },
  { name: 'tipe', label: 'Tipe', field: 'tipe', align: 'left' },
  { name: 'namaKategori', label: 'Kategori', field: 'namaKategori', align: 'left' },
  { name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right', sortable: true },
  { name: 'keterangan', label: 'Keterangan', field: 'keterangan', align: 'left' },
]

async function loadReport() {
  await financeStore.fetchCashFlow({
    periodeStart: filter.periodeStart,
    periodeEnd: filter.periodeEnd,
    akunBank: filter.akunBank || '',
  })
}

function resetFilter() {
  const n = new Date()
  filter.periodeStart = new Date(n.getFullYear(), n.getMonth(), 1).toISOString().slice(0, 10)
  filter.periodeEnd = new Date(n.getFullYear(), n.getMonth() + 1, 0).toISOString().slice(0, 10)
  filter.akunBank = null
  loadReport()
}

function onExport() {
  notifyInfo('Fitur export PDF segera hadir.')
}

onMounted(() => {
  financeStore.fetchRekening({ perPage: 100 })
  loadReport()
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.cf-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

.cf-empty {
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

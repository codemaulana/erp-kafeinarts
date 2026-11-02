<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Piutang (AR Aging)"
      subtitle="Pantau status dan usia piutang klien perusahaan"
      icon="schedule"
      icon-bg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      icon-color="#E8A547"
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
        <q-icon name="filter_list" size="18px" color="orange-5" />
        <span>Filter Laporan</span>
      </div>
      <div class="filter-grid">
        <div class="form-field-group">
          <label class="form-label">Status Pembayaran</label>
          <q-select
            v-model="filter.status"
            :options="statusOptions"
            outlined
            dense
            clearable
            emit-value
            map-options
            class="clay-input"
            label="Pilih Status Pembayaran"
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Nama Klien</label>
          <q-input
            v-model="filter.klien"
            outlined
            dense
            clearable
            placeholder="Cari nama klien..."
            class="clay-input"
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

    <template v-if="financeStore.arAgingData">
      <!-- Summary Cards -->
      <div class="ar-summary q-mb-lg">
        <AppStatCard
          icon="description"
          icon-bg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
          icon-color="#5B9BD5"
          label="Total Invoice"
          :value="(financeStore.arAgingData.rows || []).length"
        />
        <AppStatCard
          icon="payments"
          icon-bg="linear-gradient(135deg,#FEECF0,#F9D9A0)"
          icon-color="#E8A547"
          label="Total Tagihan"
          :value="formatCurrency(financeStore.arAgingData.summary?.totalTagihan || 0)"
          value-color="var(--clay-warning)"
        />
        <AppStatCard
          icon="check_circle"
          icon-bg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
          icon-color="#4CAF8A"
          label="Sudah Lunas"
          :value="formatCurrency(financeStore.arAgingData.summary?.totalLunas || 0)"
          value-color="var(--clay-success)"
        />
        <AppStatCard
          icon="pending"
          icon-bg="linear-gradient(135deg,#FEECF0,#F9BDC8)"
          icon-color="#E85D75"
          label="Belum Terbayar"
          :value="formatCurrency(financeStore.arAgingData.summary?.totalBelumLunas || 0)"
          value-color="var(--clay-danger)"
        />
      </div>

      <!-- Table -->
      <div class="clay-card clay-table-wrap">
        <div class="report-table-header">
          <span class="report-table-title">Detail Piutang Klien</span>
          <span class="report-table-count">{{ (financeStore.arAgingData.rows || []).length }} invoice</span>
        </div>
        <q-table
          :rows="financeStore.arAgingData.rows || []"
          :columns="columns"
          row-key="id"
          flat
          class="clay-table"
          :pagination="{ rowsPerPage: 10 }"
          no-data-label="Tidak ada data piutang"
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
          <template #body-cell-status="props">
            <q-td :props="props">
              <AppBadge :variant="statusVariant(props.value)">{{ props.value }}</AppBadge>
            </q-td>
          </template>
          <template #body-cell-grandTotal="props">
            <q-td :props="props">
              <span style="font-weight: 700; color: var(--clay-warning)">{{
                formatCurrency(props.value)
              }}</span>
            </q-td>
          </template>
        </q-table>
      </div>
    </template>

    <div v-else class="clay-card ar-empty">
      <q-icon name="schedule" size="56px" color="grey-4" />
      <div class="ar-empty__text">
        Klik <strong>"Tampilkan"</strong> untuk melihat laporan piutang klien.
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

const filter = reactive({ status: null, klien: '' })

const statusOptions = [
  { label: 'Lunas', value: 'Lunas' },
  { label: 'Sebagian', value: 'Sebagian' },
  { label: 'Belum Lunas', value: 'Belum Lunas' },
]

const columns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center' },
  { name: 'noInvoice', label: 'No. Invoice', field: 'noInvoice', align: 'left', sortable: true },
  { name: 'klien', label: 'Klien', field: 'klien', align: 'left' },
  { name: 'proyek', label: 'Proyek', field: 'proyek', align: 'left' },
  { name: 'tglTerbit', label: 'Tgl. Terbit', field: 'tglTerbit', align: 'left' },
  { name: 'tglJatuhTempo', label: 'Jatuh Tempo', field: 'tglJatuhTempo', align: 'left' },
  {
    name: 'grandTotal',
    label: 'Total Tagihan',
    field: 'grandTotal',
    align: 'right',
    sortable: true,
  },
  { name: 'status', label: 'Status', field: 'status', align: 'left' },
]

function statusVariant(status) {
  return { Lunas: 'success', Sebagian: 'warning', 'Belum Lunas': 'danger' }[status] || 'neutral'
}

async function loadReport() {
  await financeStore.fetchArAging({ status: filter.status || '', klien: filter.klien || '' })
}

function resetFilter() {
  filter.status = null
  filter.klien = ''
  loadReport()
}

function onExport() {
  notifyInfo('Fitur export PDF segera hadir.')
}

onMounted(() => loadReport())
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.ar-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

.ar-empty {
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

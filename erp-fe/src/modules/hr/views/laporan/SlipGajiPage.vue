<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Slip Gaji"
      subtitle="Rekap penggajian karyawan per periode"
      icon="request_quote"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    />

    <div class="page-content q-mt-md">
      <!-- Filter Card -->
      <div class="clay-card q-mb-md filter-card">
        <div class="filter-title">Filter Laporan</div>
        <div class="filter-grid">
          <div class="form-field-group">
            <label class="form-label">Periode</label>
            <q-input
              v-model="filters.periode"
              type="month"
              outlined
              dense
              class="clay-input"
            />
          </div>
          <div class="form-field-group">
            <label class="form-label">Departemen</label>
            <q-select
              v-model="filters.departemen"
              :options="hrStore.departemenOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              clearable
            />
          </div>
          <div class="form-field-group">
            <label class="form-label">Karyawan</label>
            <q-select
              v-model="filters.karyawan"
              :options="hrStore.karyawanOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              clearable
            />
          </div>
          
          <div class="form-field-group--btn">
            <q-btn
              label="Reset"
              unelevated
              outline
              no-caps
              class="clay-btn clay-btn-outline q-mr-sm"
              @click="handleReset"
            />
            <q-btn
              label="Tampilkan"
              unelevated
              no-caps
              class="clay-btn clay-btn-primary"
              :loading="hrStore.isLoading('slipGaji')"
              @click="handleFilter"
            />
          </div>
        </div>
      </div>

      <!-- Summary Section -->
      <div v-if="hasData" class="cf-summary q-mb-md">
        <AppStatCard
          title="Total Payroll"
          :value="formatCurrency(hrStore.slipGajiData?.summary?.totalPayroll || 0)"
          icon="paid"
          iconColor="#FF9800"
          iconBg="#FFF3E0"
        />
        <AppStatCard
          title="Total Pendapatan"
          :value="formatCurrency(hrStore.slipGajiData?.summary?.totalEarnings || 0)"
          icon="trending_up"
          iconColor="#4CAF50"
          iconBg="#E8F5E9"
        />
        <AppStatCard
          title="Total Potongan"
          :value="formatCurrency(hrStore.slipGajiData?.summary?.totalDeductions || 0)"
          icon="trending_down"
          iconColor="#F44336"
          iconBg="#FFEBEE"
        />
      </div>

      <!-- Data Section -->
      <div v-if="hasData" class="clay-card table-card">
        <div class="report-table-header row items-center justify-between q-mb-md">
          <div class="row items-center">
            <div class="report-table-title text-h6 text-weight-bold q-mr-md">Hasil Laporan</div>
            <div class="report-table-count text-caption text-grey-7">{{ hrStore.slipGajiData?.rows?.length || 0 }} data ditemukan</div>
          </div>
          <q-btn
            icon="picture_as_pdf"
            label="Export PDF"
            unelevated
            no-caps
            class="clay-btn clay-btn-primary"
            @click="exportPdf"
          />
        </div>

        <q-table
          class="clay-table flat"
          :rows="hrStore.slipGajiData?.rows || []"
          :columns="columns"
          row-key="employee_number"
          :pagination="{ rowsPerPage: 10 }"
          flat
        >
          <template v-slot:body-cell-no="props">
            <q-td :props="props">
              {{ props.rowIndex + 1 }}
            </q-td>
          </template>

          <template v-slot:body-cell-total_earnings="props">
            <q-td :props="props">
              {{ formatCurrency(props.row.total_earnings) }}
            </q-td>
          </template>
          
          <template v-slot:body-cell-total_deductions="props">
            <q-td :props="props">
              {{ formatCurrency(props.row.total_deductions) }}
            </q-td>
          </template>
          
          <template v-slot:body-cell-net_salary="props">
            <q-td :props="props" class="text-weight-bold text-positive">
              {{ formatCurrency(props.row.net_salary) }}
            </q-td>
          </template>

          <template v-slot:body-cell-payment_status="props">
            <q-td :props="props">
              <AppBadge
                :label="props.row.payment_status"
                :variant="getStatusVariant(props.row.payment_status)"
              />
            </q-td>
          </template>
        </q-table>
      </div>

      <!-- Empty State -->
      <div v-else class="clay-card cf-empty flex flex-center column">
        <q-icon name="request_quote" size="48px" color="grey-6" class="q-mb-md" />
        <div class="text-h6 text-weight-medium" style="color: var(--rb-text-primary, #0F172A)">Belum Ada Data</div>
        <div class="text-body2" style="color: var(--rb-text-muted, #64748B)">Pilih periode dan klik "Tampilkan" untuk melihat slip gaji.</div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useNotify } from 'src/composables/useNotify'
import { formatCurrency } from 'src/utils/formatter'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const hrStore = useHrStore()
const { notifyInfo, notifyError } = useNotify()

const filters = ref({
  periode: '',
  departemen: null,
  karyawan: null
})

const hasData = computed(() => !!hrStore.slipGajiData?.rows)

const columns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center' },
  { name: 'employee_number', label: 'NIP', field: 'employee_number', align: 'left', sortable: true },
  { name: 'employee_name', label: 'Nama Karyawan', field: 'employee_name', align: 'left', sortable: true },
  { name: 'total_earnings', label: 'Total Pendapatan', field: 'total_earnings', align: 'right', sortable: true },
  { name: 'total_deductions', label: 'Total Potongan', field: 'total_deductions', align: 'right', sortable: true },
  { name: 'net_salary', label: 'Gaji Bersih (THP)', field: 'net_salary', align: 'right', sortable: true },
  { name: 'payment_status', label: 'Status', field: 'payment_status', align: 'center' }
]

const getStatusVariant = (status) => {
  if (!status) return 'neutral'
  const s = status.toUpperCase()
  if (s === 'DRAFT') return 'neutral'
  if (s === 'PROCESSED') return 'warning'
  if (s === 'PAID') return 'success'
  return 'neutral'
}

const handleFilter = async () => {
  try {
    await hrStore.fetchSlipGaji({
      period: filters.value.periode,
      department_id: filters.value.departemen,
      employee_id: filters.value.karyawan
    })
  } catch {
    notifyError('Gagal memuat data laporan')
  }
}

const handleReset = () => {
  filters.value = {
    periode: '',
    departemen: null,
    karyawan: null
  }
  hrStore.slipGajiData = null
}

const exportPdf = () => {
  notifyInfo('Fitur export PDF akan segera hadir')
}

onMounted(async () => {
  try {
    await Promise.all([
      hrStore.fetchDepartemen({ perPage: 100 }),
      hrStore.fetchKaryawan({ perPage: 100 })
    ])
  } catch {
    notifyError('Gagal memuat data filter')
  }
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.filter-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
  align-items: end;
}

.form-field-group--btn {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
}

.cf-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
  gap: 16px;
}

.cf-empty {
  min-height: 300px;
  background: var(--rb-surface, #ffffff);
  border: 1px dashed var(--rb-border, #e2e8f0);
}
</style>

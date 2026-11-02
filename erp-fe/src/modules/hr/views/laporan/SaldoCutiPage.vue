<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Saldo Cuti"
      subtitle="Monitoring jatah dan sisa cuti karyawan"
      icon="date_range"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    />

    <div class="page-content q-mt-md">
      <!-- Filter Card -->
      <div class="clay-card q-mb-md filter-card">
        <div class="filter-title">Filter Laporan</div>
        <div class="filter-grid">
          <div class="form-field-group">
            <label class="form-label">Tahun</label>
            <q-select
              v-model="filters.tahun"
              :options="tahunOptions"
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
          
          <div class="form-field-group--btn" style="grid-column: 1 / -1">
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
              :loading="hrStore.isLoading('saldoCuti')"
              @click="handleFilter"
            />
          </div>
        </div>
      </div>

      <!-- Data Section -->
      <div v-if="hasData" class="clay-card table-card">
        <div class="report-table-header row items-center justify-between q-mb-md">
          <div class="row items-center">
            <div class="report-table-title text-h6 text-weight-bold q-mr-md">Hasil Laporan</div>
            <div class="report-table-count text-caption text-grey-7">{{ hrStore.saldoCutiData?.rows?.length || 0 }} data ditemukan</div>
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
          :rows="hrStore.saldoCutiData?.rows || []"
          :columns="columns"
          row-key="id"
          :pagination="{ rowsPerPage: 10 }"
          flat
        >
          <template v-slot:body-cell-no="props">
            <q-td :props="props">
              {{ props.rowIndex + 1 }}
            </q-td>
          </template>

          <template v-slot:body-cell-remaining_days="props">
            <q-td :props="props" class="text-weight-bold" :style="getRemainingDaysStyle(props.row.remaining_days)">
              {{ props.row.remaining_days }}
            </q-td>
          </template>
        </q-table>
      </div>

      <!-- Empty State -->
      <div v-else class="clay-card cf-empty flex flex-center column">
        <q-icon name="date_range" size="48px" color="grey-6" class="q-mb-md" />
        <div class="text-h6 text-weight-medium" style="color: var(--rb-text-primary, #0F172A)">Belum Ada Data</div>
        <div class="text-body2" style="color: var(--rb-text-muted, #64748B)">Pilih karyawan atau tahun dan klik "Tampilkan" untuk melihat saldo cuti.</div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useNotify } from 'src/composables/useNotify'

import AppPageHeader from 'components/shared/AppPageHeader.vue'

const hrStore = useHrStore()
const { notifyInfo, notifyError } = useNotify()

const currentYear = new Date().getFullYear()
const tahunOptions = [
  { label: `${currentYear}`, value: currentYear },
  { label: `${currentYear - 1}`, value: currentYear - 1 },
  { label: `${currentYear - 2}`, value: currentYear - 2 }
]

const filters = ref({
  tahun: currentYear,
  karyawan: null
})

const hasData = computed(() => !!hrStore.saldoCutiData?.rows)

const columns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center' },
  { name: 'employee_number', label: 'NIP', field: 'employee_number', align: 'left', sortable: true },
  { name: 'employee_name', label: 'Nama Karyawan', field: 'employee_name', align: 'left', sortable: true },
  { name: 'leave_name', label: 'Jenis Cuti', field: 'leave_name', align: 'left', sortable: true },
  { name: 'quota_days', label: 'Jatah Hari', field: 'quota_days', align: 'center', sortable: true },
  { name: 'used_days', label: 'Terpakai', field: 'used_days', align: 'center', sortable: true },
  { name: 'remaining_days', label: 'Sisa Saldo', field: 'remaining_days', align: 'center', sortable: true }
]

const getRemainingDaysStyle = (days) => {
  if (days <= 2) return 'color: #F44336;'
  if (days <= 5) return 'color: #FF9800;'
  return 'color: #4CAF50;'
}

const handleFilter = async () => {
  try {
    await hrStore.fetchSaldoCuti({
      year: filters.value.tahun,
      employee_id: filters.value.karyawan
    })
  } catch {
    notifyError('Gagal memuat data laporan')
  }
}

const handleReset = () => {
  filters.value = {
    tahun: currentYear,
    karyawan: null
  }
  hrStore.saldoCutiData = null
}

const exportPdf = () => {
  notifyInfo('Fitur export PDF akan segera hadir')
}

onMounted(async () => {
  try {
    await hrStore.fetchKaryawan({ perPage: 100 })
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

.cf-empty {
  min-height: 300px;
  background: var(--rb-surface, #ffffff);
  border: 1px dashed var(--rb-border, #e2e8f0);
}
</style>

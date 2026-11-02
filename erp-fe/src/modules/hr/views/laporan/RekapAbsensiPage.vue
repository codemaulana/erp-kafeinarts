<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Rekap Absensi"
      subtitle="Analisis kehadiran karyawan berdasarkan periode"
      icon="fact_check"
      iconBg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      iconColor="#5B9BD5"
    />

    <div class="page-content q-mt-md">
      <!-- Filter Card -->
      <div class="clay-card q-mb-md filter-card">
        <div class="filter-title">Filter Laporan</div>
        <div class="filter-grid">
          <div class="form-field-group">
            <label class="form-label">Periode Mulai</label>
            <q-input
              v-model="filters.periodeStart"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
            />
          </div>
          <div class="form-field-group">
            <label class="form-label">Periode Selesai</label>
            <q-input
              v-model="filters.periodeEnd"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
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
              :loading="hrStore.isLoading('rekapAbsensi')"
              @click="handleFilter"
            />
          </div>
        </div>
      </div>

      <!-- Summary Section -->
      <div v-if="hasData" class="cf-summary q-mb-md">
        <AppStatCard
          title="Total Hadir"
          :value="formatNumber(hrStore.rekapAbsensiData?.summary?.totalHadir || 0)"
          icon="check_circle"
          iconColor="#4CAF50"
          iconBg="#E8F5E9"
        />
        <AppStatCard
          title="Terlambat"
          :value="formatNumber(hrStore.rekapAbsensiData?.summary?.totalTerlambat || 0)"
          icon="schedule"
          iconColor="#FFC107"
          iconBg="#FFF8E1"
        />
        <AppStatCard
          title="Alpha"
          :value="formatNumber(hrStore.rekapAbsensiData?.summary?.totalAlpha || 0)"
          icon="cancel"
          iconColor="#F44336"
          iconBg="#FFEBEE"
        />
        <AppStatCard
          title="Cuti"
          :value="formatNumber(hrStore.rekapAbsensiData?.summary?.totalCuti || 0)"
          icon="event"
          iconColor="#2196F3"
          iconBg="#E3F2FD"
        />
      </div>

      <!-- Data Section -->
      <div v-if="hasData" class="clay-card table-card">
        <div class="report-table-header row items-center justify-between q-mb-md">
          <div class="row items-center">
            <div class="report-table-title text-h6 text-weight-bold q-mr-md">Hasil Laporan</div>
            <div class="report-table-count text-caption text-grey-7">{{ hrStore.rekapAbsensiData?.rows?.length || 0 }} data ditemukan</div>
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
          :rows="hrStore.rekapAbsensiData?.rows || []"
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
        </q-table>
      </div>

      <!-- Empty State -->
      <div v-else class="clay-card cf-empty flex flex-center column">
        <q-icon name="bar_chart" size="48px" color="grey-6" class="q-mb-md" />
        <div class="text-h6 text-weight-medium" style="color: var(--rb-text-primary, #0F172A)">Belum Ada Data</div>
        <div class="text-body2" style="color: var(--rb-text-muted, #64748B)">Pilih periode dan klik "Tampilkan" untuk melihat rekap absensi.</div>
      </div>
    </div>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useNotify } from 'src/composables/useNotify'
import { formatNumber } from 'src/utils/formatter'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'

const hrStore = useHrStore()
const { notifyInfo, notifyError } = useNotify()

const filters = ref({
  periodeStart: '',
  periodeEnd: '',
  departemen: null,
  karyawan: null
})

const hasData = computed(() => !!hrStore.rekapAbsensiData?.rows)

const columns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center' },
  { name: 'employee_number', label: 'NIP', field: 'employee_number', align: 'left', sortable: true },
  { name: 'employee_name', label: 'Nama Karyawan', field: 'employee_name', align: 'left', sortable: true },
  { name: 'total_hadir', label: 'Total Hadir', field: 'total_hadir', align: 'center', sortable: true },
  { name: 'total_terlambat', label: 'Terlambat', field: 'total_terlambat', align: 'center', sortable: true },
  { name: 'total_alpha', label: 'Alpha', field: 'total_alpha', align: 'center', sortable: true },
  { name: 'total_cuti', label: 'Cuti', field: 'total_cuti', align: 'center', sortable: true },
  { name: 'total_sakit', label: 'Sakit', field: 'total_sakit', align: 'center', sortable: true }
]

const handleFilter = async () => {
  try {
    await hrStore.fetchRekapAbsensi({
      start_date: filters.value.periodeStart,
      end_date: filters.value.periodeEnd,
      department_id: filters.value.departemen,
      employee_id: filters.value.karyawan
    })
  } catch {
    notifyError('Gagal memuat data laporan')
  }
}

const handleReset = () => {
  filters.value = {
    periodeStart: '',
    periodeEnd: '',
    departemen: null,
    karyawan: null
  }
  hrStore.rekapAbsensiData = null
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

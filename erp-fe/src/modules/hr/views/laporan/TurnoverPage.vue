<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Laporan Turnover Karyawan"
      subtitle="Analisis tingkat perputaran karyawan per periode"
      icon="swap_horiz"
      iconBg="linear-gradient(135deg,#FEECF0,#F9BDC8)"
      iconColor="#E85D75"
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
              :loading="hrStore.isLoading('turnover')"
              @click="handleFilter"
            />
          </div>
        </div>
      </div>

      <!-- Summary Section -->
      <div v-if="hasData" class="cf-summary q-mb-md">
        <AppStatCard
          title="Karyawan Masuk"
          :value="formatNumber(hrStore.turnoverData?.summary?.totalMasuk || 0)"
          icon="person_add"
          iconColor="#4CAF50"
          iconBg="#E8F5E9"
        />
        <AppStatCard
          title="Karyawan Keluar"
          :value="formatNumber(hrStore.turnoverData?.summary?.totalKeluar || 0)"
          icon="person_remove"
          iconColor="#F44336"
          iconBg="#FFEBEE"
        />
        <AppStatCard
          title="Total Aktif"
          :value="formatNumber(hrStore.turnoverData?.summary?.totalAktif || 0)"
          icon="groups"
          iconColor="#2196F3"
          iconBg="#E3F2FD"
        />
        <AppStatCard
          title="Avg Turnover Rate"
          :value="`${hrStore.turnoverData?.summary?.avgTurnoverRate || 0}%`"
          icon="percent"
          iconColor="#FF9800"
          iconBg="#FFF3E0"
        />
      </div>

      <!-- Data Section -->
      <div v-if="hasData" class="clay-card table-card">
        <div class="report-table-header row items-center justify-between q-mb-md">
          <div class="row items-center">
            <div class="report-table-title text-h6 text-weight-bold q-mr-md">Hasil Laporan</div>
            <div class="report-table-count text-caption text-grey-7">{{ hrStore.turnoverData?.rows?.length || 0 }} data ditemukan</div>
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
          :rows="hrStore.turnoverData?.rows || []"
          :columns="columns"
          row-key="department_name"
          :pagination="{ rowsPerPage: 10 }"
          flat
        >
          <template v-slot:body-cell-no="props">
            <q-td :props="props">
              {{ props.rowIndex + 1 }}
            </q-td>
          </template>

          <template v-slot:body-cell-turnover_rate="props">
            <q-td :props="props">
              <AppBadge
                :label="`${props.row.turnover_rate}%`"
                :variant="getTurnoverVariant(props.row.turnover_rate)"
              />
            </q-td>
          </template>
        </q-table>
      </div>

      <!-- Empty State -->
      <div v-else class="clay-card cf-empty flex flex-center column">
        <q-icon name="swap_horiz" size="48px" color="grey-6" class="q-mb-md" />
        <div class="text-h6 text-weight-medium" style="color: var(--rb-text-primary, #0F172A)">Belum Ada Data</div>
        <div class="text-body2" style="color: var(--rb-text-muted, #64748B)">Pilih periode dan klik "Tampilkan" untuk melihat turnover karyawan.</div>
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
import AppBadge from 'components/shared/AppBadge.vue'

const hrStore = useHrStore()
const { notifyInfo, notifyError } = useNotify()

const filters = ref({
  periodeStart: '',
  periodeEnd: '',
  departemen: null
})

const hasData = computed(() => !!hrStore.turnoverData?.rows)

const columns = [
  { name: 'no', label: 'No.', field: 'no', align: 'center' },
  { name: 'department_name', label: 'Departemen', field: 'department_name', align: 'left', sortable: true },
  { name: 'employees_in', label: 'Karyawan Masuk', field: 'employees_in', align: 'center', sortable: true },
  { name: 'employees_out', label: 'Karyawan Keluar', field: 'employees_out', align: 'center', sortable: true },
  { name: 'total_active', label: 'Total Aktif', field: 'total_active', align: 'center', sortable: true },
  { name: 'turnover_rate', label: 'Turnover Rate %', field: 'turnover_rate', align: 'center', sortable: true }
]

const getTurnoverVariant = (rate) => {
  const numRate = Number(rate) || 0
  if (numRate <= 5) return 'success'
  if (numRate <= 15) return 'warning'
  return 'danger'
}

const handleFilter = async () => {
  try {
    await hrStore.fetchTurnover({
      start_date: filters.value.periodeStart,
      end_date: filters.value.periodeEnd,
      department_id: filters.value.departemen
    })
  } catch {
    notifyError('Gagal memuat data laporan')
  }
}

const handleReset = () => {
  filters.value = {
    periodeStart: '',
    periodeEnd: '',
    departemen: null
  }
  hrStore.turnoverData = null
}

const exportPdf = () => {
  notifyInfo('Fitur export PDF akan segera hadir')
}

onMounted(async () => {
  try {
    await hrStore.fetchDepartemen({ perPage: 100 })
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

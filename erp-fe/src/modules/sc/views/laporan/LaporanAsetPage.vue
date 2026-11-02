<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader title="Laporan Aset per Karyawan / Departemen" icon="devices">
      <template #actions>
        <q-btn
          unelevated
          no-caps
          icon="refresh"
          label="Refresh"
          class="clay-btn clay-btn-primary"
          @click="loadData"
        />
      </template>
    </AppPageHeader>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-4">
        <AppStatCard title="Total Aset" :value="totalAset" icon="devices_other" color="primary" />
      </div>
      <div class="col-12 col-md-4">
        <AppStatCard title="Digunakan" :value="digunakanCount" icon="person" color="warning" />
      </div>
      <div class="col-12 col-md-4">
        <AppStatCard title="Tersedia" :value="tersediaCount" icon="check_circle" color="success" />
      </div>
    </div>

    <q-card class="app-card q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-4">
            <q-select
              v-model="filters.department_id"
              :options="hrStore.departemenOptions"
              outlined
              dense
              class="clay-input"
              label="Departemen"
              emit-value
              map-options
              clearable
            />
          </div>
          <div class="col-12 col-md-4">
            <q-select
              v-model="filters.employee_id"
              :options="hrStore.karyawanOptions"
              outlined
              dense
              class="clay-input"
              label="Karyawan"
              emit-value
              map-options
              clearable
            />
          </div>
          <div class="col-12 col-md-4">
            <q-select
              v-model="filters.category_id"
              :options="scStore.kategoriItemOptions"
              outlined
              dense
              class="clay-input"
              label="Kategori"
              emit-value
              map-options
              clearable
            />
          </div>
          <div class="col-12">
            <q-btn unelevated no-caps label="Terapkan Filter" class="clay-btn clay-btn-primary" @click="loadData" />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <q-card class="app-card">
      <AppDataTable
        :columns="columns"
        :data="scStore.laporanAsetData"
        :loading="scStore.isLoading('fetchLaporanAset')"
      >
        <template #cell-acquisition_value="{ row }">
          {{ formatCurrency(row.acquisition_value) }}
        </template>
        <template #cell-status="{ row }">
          <AppBadge :variant="getStatusVariant(row.status)">
            {{ row.status }}
          </AppBadge>
        </template>
      </AppDataTable>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useScStore } from 'src/stores/sc.store'
import { useHrStore } from 'src/stores/hr.store'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import { formatCurrency } from 'src/utils/formatter'

const scStore = useScStore()
const hrStore = useHrStore()

const filters = ref({
  department_id: null,
  employee_id: null,
  category_id: null
})

const columns = [
  { name: 'asset_code', label: 'Kode Aset', field: 'asset_code', align: 'left', sortable: true },
  { name: 'item_name', label: 'Nama Aset', field: (row) => row.name || row.item_name || '-', align: 'left', sortable: true },
  { name: 'assigned_to_name', label: 'Karyawan', field: (row) => row.assigned_to_name || row.assigned_to_employee_name || '-', align: 'left', sortable: true },
  { name: 'department_name', label: 'Departemen', field: (row) => row.department_name || '-', align: 'left', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true },
  { name: 'acquisition_value', label: 'Nilai Perolehan', field: 'acquisition_value', align: 'right', sortable: true }
]

const loadData = () => {
  scStore.fetchLaporanAset(filters.value)
}

const totalAset = computed(() => scStore.laporanAsetData?.length || 0)
const digunakanCount = computed(() => scStore.laporanAsetData?.filter(a => a.status === 'IN_USE').length || 0)
const tersediaCount = computed(() => scStore.laporanAsetData?.filter(a => a.status === 'AVAILABLE').length || 0)

const getStatusVariant = (status) => {
  const map = {
    'AVAILABLE': 'success',
    'IN_USE': 'info',
    'MAINTENANCE': 'warning',
    'RETIRED': 'danger'
  }
  return map[status] || 'neutral'
}

onMounted(() => {
  if (hrStore.fetchDepartemen) hrStore.fetchDepartemen()
  if (hrStore.fetchKaryawan) hrStore.fetchKaryawan()
  loadData()
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

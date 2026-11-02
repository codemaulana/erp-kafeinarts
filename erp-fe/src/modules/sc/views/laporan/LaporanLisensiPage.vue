<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader title="Laporan Lisensi & Langganan" icon="timelapse">
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
      <div class="col-12 col-md-6">
        <AppStatCard title="Total Lisensi" :value="totalLisensi" icon="card_membership" color="primary" />
      </div>
      <div class="col-12 col-md-6">
        <AppStatCard title="Akan Expired (<= 30 hari)" :value="akanExpiredCount" icon="warning" color="warning" />
      </div>
    </div>

    <q-card class="app-card q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-4">
            <q-input
              v-model="filters.end_date_from"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
              label="Expired Dari"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input
              v-model="filters.end_date_to"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
              label="Expired Sampai"
            />
          </div>
          <div class="col-12 col-md-4">
            <ScVendorSelect
              v-model="filters.vendor_id"
              label="Vendor"
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
        :data="scStore.laporanLisensiData"
        :loading="scStore.isLoading('fetchLaporanLisensi')"
      >
        <template #cell-end_date="{ row }">
          {{ formatDate(row.end_date) }}
        </template>
        <template #cell-sisa_hari="{ row }">
          <span :class="{'text-negative text-weight-bold': calculateSisaHari(row.end_date) <= 30}">
            {{ calculateSisaHari(row.end_date) }} hari
          </span>
        </template>
        <template #cell-amount="{ row }">
          {{ formatCurrency(row.amount) }}
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
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import ScVendorSelect from 'src/modules/sc/components/ScVendorSelect.vue'
import { formatCurrency, formatDate } from 'src/utils/formatter'

const scStore = useScStore()

const filters = ref({
  end_date_from: '',
  end_date_to: '',
  vendor_id: null
})

const columns = [
  { name: 'item_name', label: 'Nama Lisensi', field: 'item_name', align: 'left', sortable: true },
  { name: 'vendor_name', label: 'Vendor', field: 'vendor_name', align: 'left', sortable: true },
  { name: 'end_date', label: 'Tgl Expired', field: 'end_date', align: 'left', sortable: true },
  { name: 'sisa_hari', label: 'Sisa Hari', field: row => calculateSisaHari(row.end_date), align: 'right', sortable: true },
  { name: 'amount', label: 'Biaya', field: 'amount', align: 'right', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true }
]

const loadData = () => {
  scStore.fetchLaporanLisensi(filters.value)
}

const calculateSisaHari = (endDate) => {
  if (!endDate) return 0
  const today = new Date()
  today.setHours(0, 0, 0, 0)
  const end = new Date(endDate)
  const diffTime = end - today
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
  return diffDays
}

const totalLisensi = computed(() => scStore.laporanLisensiData?.length || 0)
const akanExpiredCount = computed(() => scStore.laporanLisensiData?.filter(l => calculateSisaHari(l.end_date) <= 30).length || 0)

const getStatusVariant = (status) => {
  const map = {
    'ACTIVE': 'success',
    'EXPIRED': 'danger',
    'CANCELLED': 'neutral'
  }
  return map[status] || 'neutral'
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

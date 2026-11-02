<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader title="Laporan Purchase Order" icon="receipt_long">
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
        <AppStatCard title="Total PO" :value="totalPO" icon="receipt" color="primary" />
      </div>
      <div class="col-12 col-md-4">
        <AppStatCard title="Total Amount" :value="formatCurrency(totalAmount)" icon="payments" color="positive" />
      </div>
      <div class="col-12 col-md-4">
        <AppStatCard title="Selesai" :value="selesaiCount" icon="check_circle" color="info" />
      </div>
    </div>

    <q-card class="app-card q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-3">
            <q-input
              v-model="filters.start_date"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
              label="Start Date"
            />
          </div>
          <div class="col-12 col-md-3">
            <q-input
              v-model="filters.end_date"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
              label="End Date"
            />
          </div>
          <div class="col-12 col-md-3">
            <ScVendorSelect
              v-model="filters.vendor_id"
              label="Vendor"
            />
          </div>
          <div class="col-12 col-md-3">
            <q-select
              v-model="filters.status"
              :options="['DRAFT', 'SENT', 'PARTIALLY_RECEIVED', 'COMPLETED', 'CANCELLED']"
              outlined
              dense
              class="clay-input"
              label="Status"
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
        :data="scStore.laporanPOData"
        :loading="scStore.isLoading('fetchLaporanPO')"
      >
        <template #cell-po_number="{ row }">
          <span class="text-weight-bold">{{ row.po_number }}</span>
        </template>
        <template #cell-total_amount="{ row }">
          {{ formatCurrency(row.total_amount) }}
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
  start_date: '',
  end_date: '',
  vendor_id: null,
  status: null
})

const columns = [
  { name: 'po_number', label: 'PO Number', field: 'po_number', align: 'left', sortable: true },
  { name: 'vendor_name', label: 'Vendor', field: 'vendor_name', align: 'left', sortable: true },
  { name: 'po_date', label: 'PO Date', field: row => formatDate(row.po_date), align: 'left', sortable: true },
  { name: 'total_amount', label: 'Total Amount', field: 'total_amount', align: 'right', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true }
]

const loadData = () => {
  scStore.fetchLaporanPO(filters.value)
}

const totalPO = computed(() => scStore.laporanPOData?.length || 0)
const totalAmount = computed(() => scStore.laporanPOData?.reduce((sum, po) => sum + Number(po.total_amount || 0), 0) || 0)
const selesaiCount = computed(() => scStore.laporanPOData?.filter(po => po.status === 'COMPLETED').length || 0)

const getStatusVariant = (status) => {
  const map = {
    'DRAFT': 'neutral',
    'SENT': 'info',
    'PARTIALLY_RECEIVED': 'warning',
    'COMPLETED': 'success',
    'CANCELLED': 'danger'
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

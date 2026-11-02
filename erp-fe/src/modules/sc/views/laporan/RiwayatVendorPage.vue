<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader title="Laporan Riwayat Transaksi Vendor" icon="history">
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
        <AppStatCard title="Total Transaksi" :value="totalTransaksi" icon="receipt" color="primary" />
      </div>
      <div class="col-12 col-md-6">
        <AppStatCard title="Total Nominal" :value="formatCurrency(totalNominal)" icon="payments" color="positive" />
      </div>
    </div>

    <q-card class="app-card q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-4">
            <ScVendorSelect
              v-model="filters.vendor_id"
              label="Vendor"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input
              v-model="filters.start_date"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
              label="Start Date"
            />
          </div>
          <div class="col-12 col-md-4">
            <q-input
              v-model="filters.end_date"
              type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              outlined
              dense
              class="clay-input"
              label="End Date"
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
        :data="scStore.riwayatVendorData"
        :loading="scStore.isLoading('fetchRiwayatVendor')"
      >
        <template #cell-vendor_name="{ row }">
          <div class="text-weight-bold text-dark-mode-bright">{{ row.vendor_name || '-' }}</div>
        </template>
        <template #cell-po_number="{ row }">
          <span class="text-weight-medium">{{ row.po_number || '-' }}</span>
        </template>
        <template #cell-tagihan_number="{ row }">
          <span v-if="row.tagihan_number && row.tagihan_number !== '-'" class="text-weight-medium">
            {{ row.tagihan_number }}
          </span>
          <span v-else class="text-caption text-grey-5 italic">
            Belum Terbit
          </span>
        </template>
        <template #cell-date="{ row }">
          {{ formatDate(row.date || row.last_order_date) }}
        </template>
        <template #cell-total_amount="{ row }">
          <span class="text-weight-bold" style="color: var(--clay-warning)">
            {{ formatCurrency(row.total_amount) }}
          </span>
        </template>
        <template #cell-payment_status="{ row }">
          <AppBadge :variant="paymentVariant(row.payment_status || row.status)">
            {{ formatPaymentStatus(row.payment_status || row.status) }}
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
  vendor_id: null,
  start_date: '',
  end_date: ''
})

const columns = [
  { name: 'vendor_name', label: 'Vendor', field: (row) => row.vendor_name || '-', align: 'left', sortable: true },
  { name: 'po_number', label: 'PO Number', field: (row) => row.po_number || '-', align: 'left', sortable: true },
  { name: 'tagihan_number', label: 'No. Tagihan', field: (row) => (row.tagihan_number && row.tagihan_number !== '-' ? row.tagihan_number : 'Belum Terbit'), align: 'left', sortable: true },
  { name: 'date', label: 'Tanggal', field: (row) => row.date || row.last_order_date || '-', align: 'left', sortable: true },
  { name: 'total_amount', label: 'Total', field: 'total_amount', align: 'right', sortable: true },
  { name: 'payment_status', label: 'Status Pembayaran', field: (row) => row.payment_status || row.status || 'UNPAID', align: 'center', sortable: true }
]

function paymentVariant(st) {
  const s = String(st || '').toUpperCase()
  if (s === 'PAID' || s === 'COMPLETED' || s === 'LUNAS') return 'success'
  if (s === 'PARTIAL' || s === 'SEBAGIAN') return 'warning'
  if (s === 'DRAFT') return 'neutral'
  return 'danger'
}

function formatPaymentStatus(st) {
  const s = String(st || '').toUpperCase()
  if (s === 'PAID' || s === 'COMPLETED' || s === 'LUNAS') return 'Lunas (PAID)'
  if (s === 'PARTIAL' || s === 'SEBAGIAN') return 'Sebagian'
  if (s === 'DRAFT') return 'Draft'
  if (s === 'SENT' || s === 'UNPAID') return 'Belum Lunas (UNPAID)'
  return st || 'Belum Lunas'
}

const loadData = () => {
  scStore.fetchRiwayatVendor(filters.value)
}

const totalTransaksi = computed(() => scStore.riwayatVendorData?.length || 0)
const totalNominal = computed(() => scStore.riwayatVendorData?.reduce((sum, trx) => sum + Number(trx.total_amount || 0), 0) || 0)

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

<!-- PayrollPage.vue -->
<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Penggajian (Payroll)"
      subtitle="Proses generate dan kelola slip gaji karyawan"
      icon="account_balance_wallet"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    >
      <q-btn
        unelevated
        no-caps
        icon="add"
        label="Buat Payroll"
        class="clay-btn clay-btn-primary"
        @click="openDialog()"
      />
    </AppPageHeader>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-sm-6 col-md-3">
        <AppStatCard title="Total Payroll" :value="hrStore.payrollMeta?.total || 0" icon="receipt_long" color="info" />
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <AppStatCard title="Draft" :value="0" icon="edit_document" color="neutral" />
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <AppStatCard title="Diproses" :value="0" icon="hourglass_top" color="warning" />
      </div>
      <div class="col-12 col-sm-6 col-md-3">
        <AppStatCard title="Dibayarkan" :value="0" icon="check_circle" color="success" />
      </div>
    </div>

    <q-card class="app-card">
      <div class="app-card__header flex items-center justify-between">
        <div class="flex items-center gap-sm">
          <q-input
            v-model="search"
            outlined
            dense
            placeholder="Cari payroll..."
            class="clay-input"
            debounce="350"
            @update:model-value="onSearch"
          >
            <template #prepend><q-icon name="search" /></template>
          </q-input>
          
          <q-select
            v-model="statusFilter"
            :options="[{label:'Semua Status',value:''},{label:'Draft',value:'DRAFT'},{label:'Diproses',value:'PROCESSED'},{label:'Dibayarkan',value:'PAID'}]"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            style="min-width: 150px"
            @update:model-value="onSearch"
          />

          <q-input
            v-model="periodeFilter"
            type="month"
            outlined
            dense
            class="clay-input"
            @update:model-value="onSearch"
          />
        </div>
      </div>

      <AppDataTable
        :rows="hrStore.payrollList"
        :columns="columns"
        :loading="hrStore.isLoading('payroll')"
        :pagination="pagination"
        @request="onRequest"
      >
        <template #body-cell-payment_status="props">
          <q-td :props="props">
            <AppBadge :variant="getStatusVariant(props.row.payment_status)">
              {{ getStatusLabel(props.row.payment_status) }}
            </AppBadge>
          </q-td>
        </template>
        
        <template #body-cell-net_salary="props">
          <q-td :props="props" class="text-bold text-positive">
            {{ formatCurrency(props.row.net_salary) }}
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn flat dense round size="sm" icon="visibility" color="info" @click="openDialog(props.row, true)">
              <q-tooltip>Lihat Detail</q-tooltip>
            </q-btn>
            <q-btn flat dense round size="sm" icon="edit" color="primary" @click="openDialog(props.row)">
              <q-tooltip>Edit</q-tooltip>
            </q-btn>
            <q-btn flat dense round size="sm" icon="delete" color="negative" @click="onDelete(props.row.id)">
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </AppDataTable>
    </q-card>

    <AppFormDialog
      v-model="dialog"
      :title="isView ? 'Detail Payroll' : (isEdit ? 'Edit Payroll' : 'Buat Payroll')"
      :loading="hrStore.isLoading('payroll')"
      :maximized="true"
      @save="onSubmit"
    >
      <div class="inv-form-grid">
        <div class="inv-form-grid__left">
          <div class="text-subtitle1 text-weight-bold q-mb-md">Informasi Payroll</div>
          
          <div class="form-group q-mb-md">
            <label class="form-label">Periode <span class="req">*</span></label>
            <q-input
              v-model="form.payroll_period"
              type="month"
              outlined
              dense
              class="clay-input"
              :rules="[required('Periode wajib diisi')]"
              lazy-rules
              :readonly="isView"
            />
          </div>

          <div class="form-group q-mb-md">
            <label class="form-label">Karyawan <span class="req">*</span></label>
            <q-select
              v-model="form.employee_id"
              :options="hrStore.karyawanOptions"
              outlined
              dense
              emit-value
              map-options
              class="clay-input"
              :rules="[required('Karyawan wajib dipilih')]"
              lazy-rules
              :readonly="isView"
            />
          </div>

          <div class="form-group q-mb-md">
            <label class="form-label">Status <span class="req">*</span></label>
            <q-select
              v-model="form.payment_status"
              :options="[{label:'Draft',value:'DRAFT'},{label:'Diproses',value:'PROCESSED'},{label:'Dibayarkan',value:'PAID'}]"
              outlined
              dense
              emit-value
              map-options
              class="clay-input"
              :rules="[required('Status wajib dipilih')]"
              lazy-rules
              :readonly="isView"
            />
          </div>
        </div>

        <div class="inv-form-grid__right">
          <div class="flex items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold">Rincian Komponen Gaji</div>
            <q-btn
              v-if="!isView"
              outline
              no-caps
              icon="add"
              label="Tambah Komponen"
              color="primary"
              size="sm"
              @click="addItem"
            />
          </div>

          <div class="inv-items">
            <div class="inv-items__header" style="grid-template-columns: 1fr 1fr 40px">
              <div>Nama Komponen</div>
              <div class="text-right">Nominal</div>
              <div></div>
            </div>

            <div v-if="!form.items || form.items.length === 0" class="q-pa-md text-center text-grey">
              Belum ada komponen gaji
            </div>

            <div
              v-for="(item, index) in form.items"
              :key="item.id"
              class="inv-items__row"
              style="grid-template-columns: 1fr 1fr 40px"
            >
              <q-select
                v-model="item.component_id"
                :options="hrStore.komponenGajiOptions"
                outlined
                dense
                emit-value
                map-options
                class="clay-input"
                :rules="[required('Komponen wajib dipilih')]"
                lazy-rules
                :readonly="isView"
              />
              <q-input
                v-model.number="item.amount"
                type="number"
                outlined
                dense
                class="clay-input text-right"
                :rules="[requiredNumber('Nominal wajib diisi'), nonNegative('Nominal tidak boleh negatif')]"
                lazy-rules
                :readonly="isView"
              />
              <div class="flex flex-center">
                <q-btn
                  v-if="!isView"
                  flat
                  round
                  dense
                  color="negative"
                  icon="remove_circle_outline"
                  @click="removeItem(index)"
                />
              </div>
            </div>
          </div>

          <div class="inv-totals q-mt-lg">
            <div class="inv-totals__row">
              <div class="inv-totals__label">Total Pendapatan</div>
              <div class="inv-totals__value text-positive">+ {{ formatCurrency(totalEarnings) }}</div>
            </div>
            <div class="inv-totals__row">
              <div class="inv-totals__label">Total Potongan</div>
              <div class="inv-totals__value text-negative">- {{ formatCurrency(totalDeductions) }}</div>
            </div>
            <q-separator class="q-my-sm" />
            <div class="inv-totals__row inv-totals--grand">
              <div class="inv-totals__label">Gaji Bersih (THP)</div>
              <div class="inv-totals__value text-primary">{{ formatCurrency(netSalary) }}</div>
            </div>
          </div>

        </div>
      </div>
      
      <template #actions v-if="isView">
        <q-btn
          unelevated
          no-caps
          label="Tutup"
          class="clay-btn"
          color="grey-3"
          text-color="grey-8"
          v-close-popup
        />
      </template>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { uid } from 'quasar'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { formatCurrency } from 'src/utils/formatter'
import { required, requiredNumber, nonNegative } from 'src/utils/validators'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'

const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const search = ref('')
const statusFilter = ref('')
const periodeFilter = ref('')

const columns = [
  { name: 'payroll_period', label: 'Periode', field: 'payroll_period', align: 'left', sortable: true },
  { name: 'employee_number', label: 'NIP', field: 'employee_number', align: 'left' },
  { name: 'employee_name', label: 'Nama Karyawan', field: 'employee_name', align: 'left' },
  { name: 'total_earnings', label: 'Total Pendapatan', field: 'total_earnings', align: 'right', format: val => formatCurrency(val) },
  { name: 'total_deductions', label: 'Total Potongan', field: 'total_deductions', align: 'right', format: val => formatCurrency(val) },
  { name: 'net_salary', label: 'THP', field: 'net_salary', align: 'right' },
  { name: 'payment_status', label: 'Status', field: 'payment_status', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, load, onRequest, resetPage } = useTablePagination({
  fetchFn: hrStore.fetchPayroll,
  getParams: () => ({
    search: search.value,
    status: statusFilter.value,
    periode: periodeFilter.value
  }),
  meta: computed(() => hrStore.payrollMeta)
})

const onSearch = () => {
  resetPage()
  load()
}

// Dialog
const dialog = ref(false)
const isEdit = ref(false)
const isView = ref(false)
const currentId = ref(null)

const form = ref({
  payroll_period: '',
  employee_id: '',
  payment_status: 'DRAFT',
  items: []
})

const resetForm = () => {
  form.value = {
    payroll_period: '',
    employee_id: '',
    payment_status: 'DRAFT',
    items: []
  }
  isEdit.value = false
  isView.value = false
  currentId.value = null
}

const openDialog = (row = null, view = false) => {
  resetForm()
  if (row) {
    isEdit.value = !view
    isView.value = view
    currentId.value = row.id
    form.value = { ...row, items: Array.isArray(row.items) ? [...row.items.map(i => ({...i, id: i.id || uid()}))] : [] }
  } else {
    addItem()
  }
  dialog.value = true
}

const addItem = () => {
  form.value.items.push({ id: uid(), component_id: '', amount: 0 })
}

const removeItem = (index) => {
  form.value.items.splice(index, 1)
}

// Computed totals
const totalEarnings = computed(() => {
  return form.value.items.reduce((sum, item) => {
    const component = hrStore.komponenGajiOptions.find(c => c.value === item.component_id)
    if (component && component.type === 'EARNING') {
      return sum + (Number(item.amount) || 0)
    }
    return sum
  }, 0)
})

const totalDeductions = computed(() => {
  return form.value.items.reduce((sum, item) => {
    const component = hrStore.komponenGajiOptions.find(c => c.value === item.component_id)
    if (component && component.type === 'DEDUCTION') {
      return sum + (Number(item.amount) || 0)
    }
    return sum
  }, 0)
})

const netSalary = computed(() => {
  return totalEarnings.value - totalDeductions.value
})

const onSubmit = async () => {
  try {
    const payload = {
      ...form.value,
      total_earnings: totalEarnings.value,
      total_deductions: totalDeductions.value,
      net_salary: netSalary.value
    }

    if (isEdit.value) {
      await hrStore.updatePayroll(currentId.value, payload)
      notifySuccess('Payroll berhasil diperbarui')
    } else {
      await hrStore.createPayroll(payload)
      notifySuccess('Payroll berhasil dibuat')
    }
    dialog.value = false
    load()
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan saat menyimpan payroll')
  }
}

const onDelete = async (id) => {
  const confirmed = await confirmDelete('Apakah Anda yakin ingin menghapus payroll ini?')
  if (confirmed) {
    try {
      await hrStore.deletePayroll(id)
      notifySuccess('Payroll berhasil dihapus')
      load()
    } catch (error) {
      notifyError(error.message || 'Gagal menghapus payroll')
    }
  }
}

const getStatusVariant = (status) => {
  const map = {
    DRAFT: 'neutral',
    PROCESSED: 'warning',
    PAID: 'success'
  }
  return map[status] || 'neutral'
}

const getStatusLabel = (status) => {
  const map = {
    DRAFT: 'Draft',
    PROCESSED: 'Diproses',
    PAID: 'Dibayarkan'
  }
  return map[status] || status
}

onMounted(async () => {
  await Promise.all([
    hrStore.fetchKaryawan({ page: 1, limit: 100 }),
    hrStore.fetchKomponenGaji({ page: 1, limit: 100 })
  ])
  load()
})
</script>

<style lang="scss">
@import 'src/css/finance-page';
</style>

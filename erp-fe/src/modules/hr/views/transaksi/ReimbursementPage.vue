<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Reimbursement / Klaim Biaya"
      subtitle="Pengajuan dan riwayat penggantian biaya operasional karyawan"
      icon="receipt_long"
      iconBg="linear-gradient(135deg,#EFF6FF,#DBEAFE)"
      iconColor="#3B82F6"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Buat Klaim"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <!-- Executive KPI Summary Cards -->
    <div class="stats-grid q-mb-lg">
      <AppStatCard
        title="Total Pengajuan"
        :value="stats.total"
        icon="receipt_long"
        color="primary"
      />
      <AppStatCard
        title="Disetujui / Selesai"
        :value="stats.approved"
        icon="check_circle"
        color="positive"
      />
      <AppStatCard
        title="Menunggu Review"
        :value="stats.pending"
        icon="hourglass_empty"
        color="warning"
      />
      <AppStatCard
        title="Total Nilai Klaim"
        :value="stats.totalNominal"
        icon="payments"
        color="purple"
      />
    </div>

    <!-- Main Data Table with Integrated Toolbar -->
    <AppDataTable
      :data="hrStore.reimbursementList"
      :columns="columns"
      :loading="hrStore.isLoading('reimbursement')"
      :pagination="pagination"
      no-data-label="Belum ada data klaim reimbursement"
      @request="onRequest"
    >
      <template #filter>
        <div class="table-toolbar">
          <div class="table-toolbar__search">
            <q-input
              v-model="filter.search"
              outlined
              dense
              clearable
              class="clay-input"
              placeholder="Cari no. klaim, nama karyawan, kategori..."
              debounce="350"
              @update:model-value="resetPageAndLoad"
            >
              <template #prepend>
                <q-icon name="search" size="18px" />
              </template>
            </q-input>
          </div>
          <div class="table-toolbar__actions">
            <q-select
              v-model="filter.status"
              outlined
              dense
              clearable
              class="clay-input filter-select"
              :options="statusOptions"
              option-value="value"
              option-label="label"
              emit-value
              map-options
              label="Semua Status"
              @update:model-value="resetPageAndLoad"
            />
            <q-btn
              flat
              dense
              icon="refresh"
              class="toolbar-btn"
              @click="load(1)"
            >
              <q-tooltip>Muat Ulang Data</q-tooltip>
            </q-btn>
          </div>
        </div>
      </template>

      <template #cell-claim_number="{ row }">
        <span class="text-weight-bold text-primary">{{ row.claim_number }}</span>
      </template>

      <template #cell-employee_name="{ row }">
        <div class="row items-center no-wrap q-gutter-sm">
          <q-avatar size="26px" color="blue-1" text-color="primary" class="text-weight-bold" font-size="11px">
            {{ (row.employee_name || row.employee?.name || 'K').charAt(0).toUpperCase() }}
          </q-avatar>
          <div>
            <div class="text-weight-medium text-dark-mode-bright">{{ row.employee_name || row.employee?.name || '-' }}</div>
            <div v-if="row.department_name" class="text-caption text-grey-6" style="font-size: 11px; line-height: 1;">{{ row.department_name }}</div>
          </div>
        </div>
      </template>

      <template #cell-category_name="{ row }">
        <q-badge outline color="blue-grey-6" class="q-px-sm q-py-xs text-weight-medium" style="border-radius: 6px;">
          {{ row.category_name || row.category?.name || row.category_id || '-' }}
        </q-badge>
      </template>

      <template #cell-claim_date="{ row }">
        <span class="text-caption text-grey-7">{{ formatDate(row.claim_date, 'medium') }}</span>
      </template>

      <template #cell-amount="{ row }">
        <span class="text-weight-bold text-dark-mode-bright">{{ formatCurrency(row.amount) }}</span>
      </template>

      <template #cell-status="{ row }">
        <AppBadge :variant="STATUS_VARIANTS[row.status] || 'neutral'">
          {{ STATUS_LABELS[row.status] || row.status }}
        </AppBadge>
      </template>

      <template #cell-actions="{ row }">
        <div class="row items-center justify-center no-wrap q-gutter-xs">
          <q-btn
            flat
            round
            dense
            size="sm"
            icon="edit"
            color="primary"
            class="action-icon-btn"
            @click="openDialog(row)"
          >
            <q-tooltip>Edit Data</q-tooltip>
          </q-btn>
          <q-btn
            flat
            round
            dense
            size="sm"
            icon="delete"
            color="negative"
            class="action-icon-btn action-icon-btn--delete"
            @click="handleDelete(row.id)"
          >
            <q-tooltip>Hapus Data</q-tooltip>
          </q-btn>
        </div>
      </template>
    </AppDataTable>

    <!-- Create / Edit Dialog -->
    <AppFormDialog
      v-model="dialog.show"
      :title="dialog.isEdit ? 'Edit Klaim Reimbursement' : 'Buat Klaim Reimbursement'"
      :loading="dialog.loading"
      :min-width="480"
      :max-width="640"
      @save="saveData"
    >
      <q-form ref="formRef" class="q-col-gutter-md">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-6">
            <label class="form-label">No. Klaim</label>
            <q-input outlined dense class="clay-input" v-model="form.claim_number" readonly placeholder="Auto-generate" />
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Tanggal Klaim <span class="req">*</span></label>
            <q-input
              outlined
              dense
              class="clay-input"
              type="text"
              @focus="e => e.target.type = 'date'"
              @blur="e => { if (!e.target.value) e.target.type = 'text' }"
              v-model="form.claim_date"
              :rules="[required('Tanggal klaim wajib diisi')]"
              lazy-rules
            />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-6">
            <label class="form-label">Karyawan <span class="req">*</span></label>
            <q-select
              outlined
              dense
              class="clay-input"
              v-model="form.employee_id"
              :options="hrStore.karyawanOptions"
              emit-value
              map-options
              option-value="value"
              option-label="label"
              :rules="[required('Karyawan wajib dipilih')]"
              lazy-rules
            />
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Kategori Klaim <span class="req">*</span></label>
            <q-select
              outlined
              dense
              class="clay-input"
              v-model="form.category_id"
              :options="financeStore.kategoriOptions"
              emit-value
              map-options
              option-value="value"
              option-label="label"
              :rules="[required('Kategori wajib dipilih')]"
              lazy-rules
            />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Nominal <span class="req">*</span></label>
            <q-input
              outlined
              dense
              class="clay-input"
              type="number"
              v-model.number="form.amount"
              min="0"
              prefix="Rp"
              :rules="[required('Nominal wajib diisi'), requiredNumber('Harus berupa angka'), positiveNumber('Harus positif')]"
              lazy-rules
            />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md" v-if="dialog.isEdit">
          <div class="col-12">
            <label class="form-label">Status <span class="req">*</span></label>
            <q-select
              outlined
              dense
              class="clay-input"
              v-model="form.status"
              :options="statusOptions"
              emit-value
              map-options
              option-value="value"
              option-label="label"
              :rules="[required('Status wajib dipilih')]"
              lazy-rules
            />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Lampiran Bukti</label>
            <AppFileUpload
              v-model="form.attachmentFile"
              :file-name="attachmentFileName"
              @file-selected="onFileSelected"
              @file-removed="onFileRemoved"
            />
          </div>
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useFinanceStore } from 'src/stores/finance.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, positiveNumber } from 'src/utils/validators'
import { formatDate, formatCurrency } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import AppFileUpload from 'components/shared/AppFileUpload.vue'

const hrStore = useHrStore()
const financeStore = useFinanceStore()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const STATUS_LABELS = { SUBMITTED: 'Diajukan', APPROVED: 'Disetujui', REJECTED: 'Ditolak', PAID: 'Dibayarkan' }
const STATUS_VARIANTS = { SUBMITTED: 'warning', APPROVED: 'success', REJECTED: 'danger', PAID: 'info' }
const statusOptions = Object.keys(STATUS_LABELS).map(key => ({ label: STATUS_LABELS[key], value: key }))

const filter = reactive({ search: '', status: null })

const columns = [
  { name: 'claim_number', field: 'claim_number', label: 'No. Klaim', align: 'left', sortable: true },
  { name: 'employee_name', field: (row) => row.employee_name || row.employee?.name || '-', label: 'Karyawan', align: 'left' },
  { name: 'category_name', field: (row) => row.category_name || row.category?.name || row.category_id || '-', label: 'Kategori', align: 'left' },
  { name: 'claim_date', field: 'claim_date', label: 'Tanggal', align: 'left', sortable: true },
  { name: 'amount', field: 'amount', label: 'Nominal', align: 'right', sortable: true },
  { name: 'status', field: 'status', label: 'Status', align: 'center', sortable: true },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination(loadData)

async function loadData(page, rowsPerPage, sortBy, descending) {
  await hrStore.fetchReimbursement({
    page,
    perPage: rowsPerPage,
    search: filter.search,
    status: filter.status,
    sortBy,
    sortOrder: descending ? 'desc' : 'asc'
  })
}

const stats = computed(() => {
  const list = hrStore.reimbursementList || []
  const totalCount = hrStore.reimbursementMeta?.total || list.length
  let approved = 0
  let pending = 0
  let totalNominal = 0

  list.forEach(item => {
    if (item.status === 'APPROVED' || item.status === 'PAID') approved++
    if (item.status === 'SUBMITTED') pending++
    totalNominal += Number(item.amount) || 0
  })

  return {
    total: totalCount,
    approved,
    pending,
    totalNominal: formatCurrency(totalNominal)
  }
})

const resetPageAndLoad = () => {
  resetPage()
  load(1)
}

const load = (page = pagination.value.page) => {
  onRequest({ pagination: { ...pagination.value, page } })
}

const formRef = ref(null)
const attachmentFileName = ref('')
const dialog = reactive({ show: false, isEdit: false, id: null, loading: false })
const form = reactive({
  claim_number: '',
  employee_id: null,
  category_id: null,
  claim_date: new Date().toISOString().split('T')[0],
  amount: 0,
  status: 'SUBMITTED',
  attachmentFile: null
})

const onFileSelected = (file) => {
  form.attachmentFile = file
  attachmentFileName.value = file.name
}

const onFileRemoved = () => {
  form.attachmentFile = null
  attachmentFileName.value = ''
}

const openDialog = (row = null) => {
  formRef.value?.resetValidation()
  dialog.isEdit = !!row
  if (row) {
    dialog.id = row.id
    form.claim_number = row.claim_number || ''
    form.employee_id = row.employee_id
    form.category_id = row.category_id
    form.claim_date = row.claim_date
    form.amount = row.amount
    form.status = row.status
    form.attachmentFile = null
    attachmentFileName.value = row.attachment_url ? row.attachment_url.split('/').pop() : ''
  } else {
    dialog.id = null
    form.claim_number = ''
    form.employee_id = null
    form.category_id = null
    form.claim_date = new Date().toISOString().split('T')[0]
    form.amount = 0
    form.status = 'SUBMITTED'
    form.attachmentFile = null
    attachmentFileName.value = ''
  }
  dialog.show = true
}

const saveData = async () => {
  const valid = await formRef.value?.validate()
  if (!valid) return

  dialog.loading = true
  try {
    const payload = new FormData()
    Object.keys(form).forEach(key => {
      if (form[key] !== null && form[key] !== undefined && key !== 'attachmentFile') {
        payload.append(key, form[key])
      }
    })
    if (form.attachmentFile) {
      payload.append('attachment', form.attachmentFile)
    }

    if (dialog.isEdit) {
      await hrStore.updateReimbursement(dialog.id, payload)
      notifySuccess('Klaim reimbursement berhasil diperbarui')
    } else {
      await hrStore.createReimbursement(payload)
      notifySuccess('Klaim reimbursement berhasil dibuat')
    }
    dialog.show = false
    load()
  } catch (err) {
    console.error(err)
    notifyError('Terjadi kesalahan')
  } finally {
    dialog.loading = false
  }
}

const handleDelete = (id) => {
  confirmDelete({
    title: 'Hapus Klaim',
    message: 'Apakah Anda yakin ingin menghapus data klaim reimbursement ini?',
    onConfirm: async () => {
      try {
        await hrStore.deleteReimbursement(id)
        notifySuccess('Klaim reimbursement berhasil dihapus')
        load()
      } catch (err) {
        console.error(err)
        notifyError('Gagal menghapus klaim')
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([
    financeStore.fetchKategori({ perPage: 100 }),
    hrStore.fetchKaryawan({ perPage: 100 })
  ])
  load(1)
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

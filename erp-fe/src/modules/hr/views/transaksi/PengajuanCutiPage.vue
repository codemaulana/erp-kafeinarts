<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Pengajuan Cuti / Izin"
      subtitle="Kelola pengajuan cuti dan izin karyawan"
      icon="event_busy"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    >
      <template #actions>
        <q-btn unelevated no-caps class="clay-btn clay-btn-primary" icon="add" label="Buat Pengajuan" @click="openDialog()" />
      </template>
    </AppPageHeader>

    <div class="stats-grid q-mb-lg">
      <AppStatCard title="Total Pengajuan" :value="hrStore.cutiMeta?.total || 0" icon="list_alt" color="primary" />
      <AppStatCard title="Disetujui" :value="stats.approved" icon="check_circle" color="positive" />
      <AppStatCard title="Menunggu" :value="stats.submitted" icon="hourglass_empty" color="warning" />
      <AppStatCard title="Ditolak" :value="stats.rejected" icon="cancel" color="negative" />
    </div>

    <AppDataTable
      :data="hrStore.cutiList"
      :columns="columns"
      :loading="hrStore.isLoading('cuti')"
      :pagination="pagination"
      no-data-label="Belum ada data pengajuan cuti"
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
              placeholder="Cari no. pengajuan, karyawan..."
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

      <template #cell-request_number="{ row }">
        <span class="text-weight-bold text-primary">{{ row.request_number }}</span>
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

      <template #cell-leave_type_name="{ row }">
        <q-badge outline color="blue-grey-6" class="q-px-sm q-py-xs text-weight-medium" style="border-radius: 6px;">
          {{ row.leave_name || row.leave_type?.name || row.leave_type_id || '-' }}
        </q-badge>
      </template>

      <template #cell-start_date="{ row }">
        <span class="text-caption text-grey-7">{{ formatDate(row.start_date, 'medium') }}</span>
      </template>

      <template #cell-end_date="{ row }">
        <span class="text-caption text-grey-7">{{ formatDate(row.end_date, 'medium') }}</span>
      </template>

      <template #cell-total_days="{ row }">
        <span class="text-weight-bold text-dark-mode-bright">{{ row.total_days }} Hari</span>
      </template>

      <template #cell-approval_status="{ row }">
        <AppBadge :variant="STATUS_VARIANTS[row.approval_status] || 'neutral'">
          {{ STATUS_LABELS[row.approval_status] || row.approval_status }}
        </AppBadge>
      </template>

      <template #cell-actions="{ row }">
        <div class="row items-center justify-center no-wrap q-gutter-xs">
          <q-btn flat round dense size="sm" icon="edit" color="primary" class="action-icon-btn" @click="openDialog(row)">
            <q-tooltip>Edit Data</q-tooltip>
          </q-btn>
          <q-btn flat round dense size="sm" icon="delete" color="negative" class="action-icon-btn action-icon-btn--delete" @click="handleDelete(row.id)">
            <q-tooltip>Hapus Data</q-tooltip>
          </q-btn>
        </div>
      </template>
    </AppDataTable>

    <AppFormDialog
      v-model="dialog.show"
      :title="dialog.isEdit ? 'Edit Pengajuan' : 'Buat Pengajuan'"
      :loading="dialog.loading"
      :min-width="520"
      :max-width="640"
      @save="saveData"
    >
      <q-form ref="formRef" class="q-col-gutter-md">
        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-6">
            <label class="form-label">No. Pengajuan</label>
            <q-input outlined dense class="clay-input" v-model="form.request_number" readonly placeholder="Auto-generate" />
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Karyawan <span class="req">*</span></label>
            <q-select outlined dense class="clay-input" v-model="form.employee_id" :options="hrStore.karyawanOptions" emit-value map-options option-value="value" option-label="label" :rules="[required('Karyawan wajib dipilih')]" lazy-rules />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Jenis Cuti <span class="req">*</span></label>
            <q-select outlined dense class="clay-input" v-model="form.leave_type_id" :options="hrStore.jenisCutiOptions" emit-value map-options option-value="value" option-label="label" :rules="[required('Jenis cuti wajib dipilih')]" lazy-rules />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12 col-md-6">
            <label class="form-label">Tgl Mulai <span class="req">*</span></label>
            <q-input outlined dense class="clay-input" type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" v-model="form.start_date" :rules="[required('Tanggal mulai wajib diisi')]" lazy-rules />
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Tgl Selesai <span class="req">*</span></label>
            <q-input outlined dense class="clay-input" type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" v-model="form.end_date" :rules="[required('Tanggal selesai wajib diisi'), dateNotBefore(() => form.start_date, 'Tidak boleh sebelum tanggal mulai')]" lazy-rules />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Total Hari</label>
            <q-input outlined dense class="clay-input" type="number" v-model="totalDaysComputed" readonly />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Alasan <span class="req">*</span></label>
            <q-input type="textarea" outlined dense class="clay-input" v-model="form.reason" rows="3" :rules="[required('Alasan wajib diisi')]" lazy-rules />
          </div>
        </div>
        
        <div class="row q-col-gutter-md q-mb-md" v-if="dialog.isEdit">
          <div class="col-12">
            <label class="form-label">Status <span class="req">*</span></label>
            <q-select outlined dense class="clay-input" v-model="form.approval_status" :options="statusOptions" emit-value map-options option-value="value" option-label="label" :rules="[required('Status wajib dipilih')]" lazy-rules />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Lampiran</label>
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
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, dateNotBefore } from 'src/utils/validators'
import { formatDate } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppFileUpload from 'components/shared/AppFileUpload.vue'

const hrStore = useHrStore()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const STATUS_LABELS = { SUBMITTED: 'Diajukan', APPROVED_SUPERVISOR: 'Disetujui Atasan', APPROVED_HR: 'Disetujui HR', REJECTED: 'Ditolak' }
const STATUS_VARIANTS = { SUBMITTED: 'warning', APPROVED_SUPERVISOR: 'info', APPROVED_HR: 'success', REJECTED: 'danger' }
const statusOptions = Object.keys(STATUS_LABELS).map(key => ({ label: STATUS_LABELS[key], value: key }))

const filter = reactive({ search: '', status: null })

const stats = reactive({
  approved: 0,
  submitted: 0,
  rejected: 0
})

const columns = [
  { name: 'request_number', field: 'request_number', label: 'No. Pengajuan', align: 'left', sortable: true },
  { name: 'employee_name', field: (row) => row.employee_name || row.employee?.name || '-', label: 'Karyawan', align: 'left' },
  { name: 'leave_type_name', field: (row) => row.leave_name || row.leave_type?.name || row.leave_type_id || '-', label: 'Jenis Cuti', align: 'left' },
  { name: 'start_date', field: 'start_date', label: 'Tgl Mulai', align: 'left', sortable: true },
  { name: 'end_date', field: 'end_date', label: 'Tgl Selesai', align: 'left', sortable: true },
  { name: 'total_days', field: 'total_days', label: 'Hari', align: 'center', sortable: true },
  { name: 'approval_status', field: 'approval_status', label: 'Status', align: 'center', sortable: true },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination(loadData)

async function loadData(page, rowsPerPage, sortBy, descending) {
  await hrStore.fetchCuti({
    page,
    perPage: rowsPerPage,
    search: filter.search,
    approval_status: filter.status,
    sortBy,
    sortOrder: descending ? 'desc' : 'asc'
  })
  updateStats()
}

const updateStats = () => {
  stats.approved = hrStore.cutiList.filter(c => c.approval_status === 'APPROVED_HR' || c.approval_status === 'APPROVED_SUPERVISOR').length
  stats.submitted = hrStore.cutiList.filter(c => c.approval_status === 'SUBMITTED').length
  stats.rejected = hrStore.cutiList.filter(c => c.approval_status === 'REJECTED').length
}

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
  request_number: '',
  employee_id: null,
  leave_type_id: null,
  start_date: '',
  end_date: '',
  total_days: 0,
  reason: '',
  approval_status: 'SUBMITTED',
  attachmentFile: null
})

const totalDaysComputed = computed(() => {
  if (form.start_date && form.end_date) {
    const start = new Date(form.start_date)
    const end = new Date(form.end_date)
    const diffTime = Math.abs(end - start)
    const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24))
    return diffDays + 1
  }
  return 0
})

watch(totalDaysComputed, (newVal) => {
  form.total_days = newVal
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
    form.request_number = row.request_number || ''
    form.employee_id = row.employee_id
    form.leave_type_id = row.leave_type_id
    form.start_date = row.start_date
    form.end_date = row.end_date
    form.total_days = row.total_days
    form.reason = row.reason
    form.approval_status = row.approval_status
    form.attachmentFile = null
    attachmentFileName.value = row.attachment_url ? row.attachment_url.split('/').pop() : ''
  } else {
    dialog.id = null
    form.request_number = ''
    form.employee_id = null
    form.leave_type_id = null
    form.start_date = ''
    form.end_date = ''
    form.total_days = 0
    form.reason = ''
    form.approval_status = 'SUBMITTED'
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
      await hrStore.updateCuti(dialog.id, payload)
      notifySuccess('Pengajuan cuti berhasil diperbarui')
    } else {
      await hrStore.createCuti(payload)
      notifySuccess('Pengajuan cuti berhasil dibuat')
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
    title: 'Hapus Pengajuan',
    message: 'Apakah Anda yakin ingin menghapus data pengajuan cuti ini?',
    onConfirm: async () => {
      try {
        await hrStore.deleteCuti(id)
        notifySuccess('Pengajuan cuti berhasil dihapus')
        load()
      } catch (err) {
        console.error(err)
        notifyError('Gagal menghapus pengajuan')
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([
    hrStore.fetchKaryawan({ perPage: 100 }),
    hrStore.fetchJenisCuti({ perPage: 100 })
  ])
  load(1)
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

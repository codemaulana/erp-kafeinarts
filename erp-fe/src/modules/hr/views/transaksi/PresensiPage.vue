<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Presensi / Absensi Harian"
      subtitle="Pencatatan kehadiran karyawan"
      icon="fingerprint"
      iconBg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      iconColor="#5B9BD5"
    >
      <template #actions>
        <q-btn unelevated no-caps class="clay-btn clay-btn-primary" icon="add" label="Tambah Presensi" @click="openDialog()" />
      </template>
    </AppPageHeader>

    <div class="stats-grid q-mb-lg">
      <AppStatCard title="Total Presensi" :value="stats.total" icon="fingerprint" color="primary" />
      <AppStatCard title="Hadir Tepat Waktu" :value="stats.present" icon="check_circle" color="positive" />
      <AppStatCard title="Terlambat" :value="stats.late" icon="warning" color="warning" />
      <AppStatCard title="Cuti / Lainnya" :value="stats.leave" icon="event_busy" color="info" />
    </div>

    <AppDataTable
      :data="hrStore.presensiList"
      :columns="columns"
      :loading="hrStore.isLoading('presensi')"
      :pagination="pagination"
      no-data-label="Belum ada data presensi karyawan"
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
              placeholder="Cari nama karyawan..."
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
            <q-input
              v-model="filter.date"
              outlined
              dense
              clearable
              class="clay-input"
              style="width: 150px;"
              type="text"
              @focus="e => e.target.type = 'date'"
              @blur="e => { if(!e.target.value) e.target.type = 'text' }"
              label="Tanggal"
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

      <template #cell-attendance_date="{ row }">
        <span class="text-caption text-grey-7">{{ formatDate(row.attendance_date, 'medium') }}</span>
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

      <template #cell-check_in_time="{ row }">
        <span class="text-weight-medium">{{ row.check_in_time ? formatTime(row.check_in_time) : '-' }}</span>
      </template>

      <template #cell-check_out_time="{ row }">
        <span class="text-weight-medium">{{ row.check_out_time ? formatTime(row.check_out_time) : '-' }}</span>
      </template>

      <template #cell-attendance_status="{ row }">
        <AppBadge :variant="STATUS_VARIANTS[row.attendance_status] || 'neutral'">
          {{ STATUS_LABELS[row.attendance_status] || row.attendance_status }}
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
      :title="dialog.isEdit ? 'Edit Presensi' : 'Tambah Presensi'"
      :loading="dialog.loading"
      @save="saveData"
    >
      <q-form ref="formRef" class="q-col-gutter-md">
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <label class="form-label">Tanggal <span class="req">*</span></label>
            <q-input outlined dense class="clay-input" type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" v-model="form.attendance_date" :rules="[required('Tanggal wajib diisi')]" lazy-rules />
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Karyawan <span class="req">*</span></label>
            <q-select outlined dense class="clay-input" v-model="form.employee_id" :options="hrStore.karyawanOptions" emit-value map-options option-value="value" option-label="label" :rules="[required('Karyawan wajib dipilih')]" lazy-rules />
          </div>
        </div>

        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <label class="form-label">Jam Masuk</label>
            <q-input outlined dense class="clay-input" type="datetime-local" v-model="form.check_in_time" />
          </div>
          <div class="col-12 col-md-6">
            <label class="form-label">Jam Pulang</label>
            <q-input outlined dense class="clay-input" type="datetime-local" v-model="form.check_out_time" />
          </div>
        </div>
        
        <div class="form-row">
          <div class="form-col">
            <label class="form-label">Latitude Masuk</label>
            <q-input outlined dense class="clay-input" type="number" v-model.number="form.check_in_lat" step="any" />
          </div>
          <div class="form-col">
            <label class="form-label">Longitude Masuk</label>
            <q-input outlined dense class="clay-input" type="number" v-model.number="form.check_in_lng" step="any" />
          </div>
          <div class="form-col flex items-end">
            <q-btn unelevated color="primary" icon="place" label="Ambil Lokasi" @click="getLocation" class="q-mb-md" />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Status <span class="req">*</span></label>
            <q-select outlined dense class="clay-input" v-model="form.attendance_status" :options="statusOptions" emit-value map-options option-value="value" option-label="label" :rules="[required('Status wajib dipilih')]" lazy-rules />
          </div>
        </div>

        <div class="row q-col-gutter-md q-mb-md">
          <div class="col-12">
            <label class="form-label">Foto / Selfie</label>
            <AppFileUpload 
              v-model="form.selfieFile"
              :file-name="selfieFileName"
              accept="image/*"
              hint="JPG, PNG — maks 2MB"
              :max-size="2097152"
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
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required } from 'src/utils/validators'
import { formatDate } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import AppFileUpload from 'components/shared/AppFileUpload.vue'

const hrStore = useHrStore()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const STATUS_LABELS = { PRESENT: 'Hadir', LATE: 'Terlambat', ABSENT: 'Alpha', LEAVE: 'Cuti', SICK: 'Sakit', BUSINESS_TRIP: 'Dinas Luar' }
const STATUS_VARIANTS = { PRESENT: 'success', LATE: 'warning', ABSENT: 'danger', LEAVE: 'info', SICK: 'warning', BUSINESS_TRIP: 'info' }
const statusOptions = Object.keys(STATUS_LABELS).map(key => ({ label: STATUS_LABELS[key], value: key }))

const filter = reactive({ search: '', status: null, date: null })

const stats = computed(() => {
  const list = hrStore.presensiList || []
  const total = hrStore.presensiMeta?.total || list.length
  let present = 0
  let late = 0
  let leave = 0

  list.forEach(item => {
    if (item.attendance_status === 'PRESENT') present++
    else if (item.attendance_status === 'LATE') late++
    else leave++
  })

  return { total, present, late, leave }
})

const formatTime = (val) => {
  if (!val) return '-'
  if (typeof val === 'string' && /^\d{1,2}:\d{2}(:\d{2})?$/.test(val.trim())) {
    return val.trim().substring(0, 5)
  }
  const date = new Date(val)
  if (isNaN(date.getTime())) return val
  return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit', hour12: false })
}

const columns = [
  { name: 'attendance_date', field: 'attendance_date', label: 'Tanggal', align: 'left', sortable: true },
  { name: 'employee_name', field: (row) => row.employee_name || row.employee?.name || '-', label: 'Nama Karyawan', align: 'left' },
  { name: 'check_in_time', field: 'check_in_time', label: 'Jam Masuk', align: 'center' },
  { name: 'check_out_time', field: 'check_out_time', label: 'Jam Pulang', align: 'center' },
  { name: 'attendance_status', field: 'attendance_status', label: 'Status', align: 'center', sortable: true },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination(loadData)

async function loadData(page, rowsPerPage, sortBy, descending) {
  await hrStore.fetchPresensi({
    page,
    perPage: rowsPerPage,
    search: filter.search,
    attendance_status: filter.status,
    attendance_date: filter.date,
    sortBy,
    sortOrder: descending ? 'desc' : 'asc'
  })
}

const resetPageAndLoad = () => {
  resetPage()
  load(1)
}

const load = (page = pagination.value.page) => {
  onRequest({ pagination: { ...pagination.value, page } })
}

const formRef = ref(null)
const selfieFileName = ref('')
const dialog = reactive({ show: false, isEdit: false, id: null, loading: false })
const form = reactive({
  attendance_date: new Date().toISOString().split('T')[0],
  employee_id: null,
  check_in_time: null,
  check_out_time: null,
  check_in_lat: null,
  check_in_lng: null,
  attendance_status: null,
  selfieFile: null
})

const getLocation = () => {
  if (navigator.geolocation) {
    navigator.geolocation.getCurrentPosition(
      (position) => {
        form.check_in_lat = position.coords.latitude
        form.check_in_lng = position.coords.longitude
      },
      (error) => {
        console.error("Error getting location", error)
        notifyError("Gagal mengambil lokasi.")
      }
    )
  } else {
    notifyError("Geolocation tidak didukung oleh browser Anda.")
  }
}

const onFileSelected = (file) => {
  form.selfieFile = file
  selfieFileName.value = file.name
}

const onFileRemoved = () => {
  form.selfieFile = null
  selfieFileName.value = ''
}

const openDialog = (row = null) => {
  formRef.value?.resetValidation()
  dialog.isEdit = !!row
  if (row) {
    dialog.id = row.id
    form.attendance_date = row.attendance_date
    form.employee_id = row.employee_id
    form.check_in_time = row.check_in_time
    form.check_out_time = row.check_out_time
    form.check_in_lat = row.check_in_lat
    form.check_in_lng = row.check_in_lng
    form.attendance_status = row.attendance_status
    form.selfieFile = null
    selfieFileName.value = row.selfie_url ? row.selfie_url.split('/').pop() : ''
  } else {
    dialog.id = null
    form.attendance_date = new Date().toISOString().split('T')[0]
    form.employee_id = null
    form.check_in_time = null
    form.check_out_time = null
    form.check_in_lat = null
    form.check_in_lng = null
    form.attendance_status = null
    form.selfieFile = null
    selfieFileName.value = ''
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
      if (form[key] !== null && form[key] !== undefined && key !== 'selfieFile') {
        payload.append(key, form[key])
      }
    })
    if (form.selfieFile) {
      payload.append('selfie', form.selfieFile)
    }

    if (dialog.isEdit) {
      await hrStore.updatePresensi(dialog.id, payload)
      notifySuccess('Presensi berhasil diperbarui')
    } else {
      await hrStore.createPresensi(payload)
      notifySuccess('Presensi berhasil ditambahkan')
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
    title: 'Hapus Presensi',
    message: 'Apakah Anda yakin ingin menghapus data presensi ini?',
    onConfirm: async () => {
      try {
        await hrStore.deletePresensi(id)
        notifySuccess('Presensi berhasil dihapus')
        load()
      } catch (err) {
        console.error(err)
        notifyError('Gagal menghapus presensi')
      }
    }
  })
}

onMounted(async () => {
  await hrStore.fetchKaryawan({ perPage: 100 })
  load(1)
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
.form-row {
  display: flex;
  gap: 1rem;
}
.form-col {
  flex: 1;
}
</style>

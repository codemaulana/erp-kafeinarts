<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Data Karyawan"
      subtitle="Kelola data induk seluruh karyawan"
      icon="person"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Karyawan"
          @click="openForm()"
        />
      </template>
    </AppPageHeader>

    <div class="row q-col-gutter-md q-mb-md">
      <div class="col-12 col-md-3">
        <AppStatCard title="Total Karyawan" :value="hrStore.karyawanList.length" icon="groups" color="blue" />
      </div>
      <div class="col-12 col-md-3">
        <AppStatCard title="Aktif" :value="activeCount" icon="check_circle" color="positive" />
      </div>
      <div class="col-12 col-md-3">
        <AppStatCard title="Kontrak" :value="contractCount" icon="assignment" color="warning" />
      </div>
      <div class="col-12 col-md-3">
        <AppStatCard title="Resign" :value="resignedCount" icon="cancel" color="negative" />
      </div>
    </div>

    <div class="app-card clay-card q-pa-md">
      <div class="row q-mb-md items-center justify-between">
        <div class="row q-gutter-sm">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            class="clay-input"
            placeholder="Cari karyawan..."
            debounce="350"
            @update:model-value="onSearch"
          >
            <template #prepend>
              <q-icon name="search" />
            </template>
          </q-input>
          <q-select
            v-model="statusFilter"
            :options="statusOptions"
            outlined
            dense
            class="clay-input"
            style="width: 150px"
            emit-value
            map-options
            @update:model-value="onSearch"
          />
        </div>
      </div>

      <AppDataTable
        :rows="hrStore.karyawanList"
        :columns="columns"
        :loading="hrStore.isLoading('karyawan')"
        :pagination="pagination"
        @request="onRequest"
      >
        <template #body-cell-active_status="props">
          <q-td :props="props">
            <AppBadge
              :variant="props.value === 'ACTIVE' ? 'success' : (props.value === 'RESIGNED' ? 'danger' : 'warning')"
              :label="props.value === 'ACTIVE' ? 'Aktif' : (props.value === 'RESIGNED' ? 'Resign' : 'Non-Aktif')"
            />
          </q-td>
        </template>
        <template #body-cell-join_date="props">
          <q-td :props="props">
            {{ formatDate(props.value, 'DD MMM YYYY') }}
          </q-td>
        </template>
        <template #body-cell-actions="props">
          <q-td :props="props">
            <div class="row justify-center q-gutter-xs">
              <q-btn flat dense round size="sm" icon="visibility" color="info" @click="openForm(props.row, true)">
                <q-tooltip>Lihat Detail</q-tooltip>
              </q-btn>
              <q-btn flat dense round size="sm" icon="edit" color="primary" @click="openForm(props.row)">
                <q-tooltip>Edit</q-tooltip>
              </q-btn>
              <q-btn flat dense round size="sm" icon="delete" color="negative" @click="handleDelete(props.row)">
                <q-tooltip>Hapus</q-tooltip>
              </q-btn>
            </div>
          </q-td>
        </template>
      </AppDataTable>
    </div>

    <AppFormDialog
      v-model="isDialogOpen"
      :title="isView ? 'Detail Karyawan' : (isEdit ? 'Edit Karyawan' : 'Tambah Karyawan')"
      maximized
      @save="onSave"
      :loading="hrStore.isLoading('karyawan_saving')"
      :hideSave="isView"
      saveLabel="Simpan"
    >
      <q-form ref="formRef" @submit.prevent="onSave">
        <div class="row q-col-gutter-lg inv-form-grid">
          <div class="col-12 col-md-6">
            <div class="text-subtitle1 q-mb-md text-weight-bold">Data Pribadi</div>
            <div class="q-mb-md">
              <label class="form-label">Foto Karyawan</label>
              <AppFileUpload
                v-if="!isView"
                v-model="photoFile"
                accept="image/*"
                hint="JPG, PNG — maks 2MB"
                :max-size="2097152"
                @update:model-value="onPhotoSelected"
                @clear="onPhotoRemoved"
              />
              <div v-if="photoFileName" class="q-mt-sm text-caption">Terpilih: {{ photoFileName }}</div>
            </div>

            <div class="q-mb-md">
              <label class="form-label">NIP</label>
              <q-input v-model="form.employee_number" outlined dense class="clay-input" readonly placeholder="Auto-generate" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Nama Karyawan <span class="req">*</span></label>
              <q-input v-model="form.full_name" outlined dense class="clay-input" :rules="[required('Nama karyawan wajib diisi')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">NIK <span class="req">*</span></label>
              <q-input v-model="form.national_id_number" outlined dense class="clay-input" :rules="[required('NIK wajib diisi'), exactLength(16, 'NIK harus 16 digit')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Tanggal Lahir <span class="req">*</span></label>
              <q-input v-model="form.birth_date" type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" outlined dense class="clay-input" :rules="[required('Tanggal lahir wajib diisi')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Jenis Kelamin <span class="req">*</span></label>
              <q-select v-model="form.gender" :options="genderOptions" outlined dense class="clay-input" emit-value map-options :rules="[required('Jenis kelamin wajib diisi')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Alamat</label>
              <q-input v-model="form.address" type="textarea" rows="2" outlined dense class="clay-input" :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">No. Telepon</label>
              <q-input v-model="form.phone_number" outlined dense class="clay-input" :rules="[phoneNumber('Format telepon tidak valid')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Email</label>
              <q-input v-model="form.email" outlined dense class="clay-input" :rules="[email('Format email tidak valid')]" lazy-rules :readonly="isView" />
            </div>
          </div>

          <div class="col-12 col-md-6">
            <div class="text-subtitle1 q-mb-md text-weight-bold">Data Kepegawaian</div>
            <div class="q-mb-md">
              <label class="form-label">Departemen <span class="req">*</span></label>
              <q-select v-model="form.department_id" :options="hrStore.departemenOptions" outlined dense class="clay-input" emit-value map-options :rules="[required('Departemen wajib dipilih')]" lazy-rules @update:model-value="form.position_id = null" :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Jabatan <span class="req">*</span></label>
              <q-select v-model="form.position_id" :options="jabatanFiltered" outlined dense class="clay-input" emit-value map-options :rules="[required('Jabatan wajib dipilih')]" lazy-rules :readonly="isView" :disable="!form.department_id" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Tgl Bergabung <span class="req">*</span></label>
              <q-input v-model="form.join_date" type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }" outlined dense class="clay-input" :rules="[required('Tgl bergabung wajib diisi')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Status Karyawan <span class="req">*</span></label>
              <q-select v-model="form.employment_status" :options="employmentOptions" outlined dense class="clay-input" emit-value map-options :rules="[required('Status karyawan wajib dipilih')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">Status Keaktifan <span class="req">*</span></label>
              <q-select v-model="form.active_status" :options="activeOptions" outlined dense class="clay-input" emit-value map-options :rules="[required('Status keaktifan wajib dipilih')]" lazy-rules :readonly="isView" />
            </div>

            <div class="q-mb-md">
              <label class="form-label">No. Rekening Bank</label>
              <q-input v-model="form.bank_account_number" outlined dense class="clay-input" :readonly="isView" />
            </div>
          </div>
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, exactLength, phoneNumber, email } from 'src/utils/validators'
import { formatDate } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppFileUpload from 'components/shared/AppFileUpload.vue'

const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const isDialogOpen = ref(false)
const isEdit = ref(false)
const isView = ref(false)
const formRef = ref(null)

const searchQuery = ref('')
const statusFilter = ref('ALL')

const photoFile = ref(null)
const photoFileName = ref('')

const defaultForm = {
  id: null,
  employee_number: '',
  full_name: '',
  national_id_number: '',
  birth_date: '',
  gender: null,
  address: '',
  phone_number: '',
  email: '',
  department_id: null,
  position_id: null,
  join_date: '',
  employment_status: null,
  active_status: 'ACTIVE',
  bank_account_number: ''
}

const form = ref({ ...defaultForm })

const columns = [
  { name: 'employee_number', field: 'employee_number', label: 'NIP', align: 'left', sortable: true },
  { name: 'full_name', field: 'full_name', label: 'Nama Karyawan', align: 'left', sortable: true },
  { name: 'department_name', field: 'department_name', label: 'Departemen', align: 'left' },
  { name: 'position_name', field: 'position_name', label: 'Jabatan', align: 'left' },
  { name: 'active_status', field: 'active_status', label: 'Status', align: 'center' },
  { name: 'join_date', field: 'join_date', label: 'Tgl Bergabung', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchData: (params) => {
    const filters = { search: searchQuery.value }
    if (statusFilter.value !== 'ALL') {
      filters.active_status = statusFilter.value
    }
    return hrStore.fetchKaryawan({ ...params, ...filters })
  },
  defaultSort: 'employee_number'
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const statusOptions = [
  { label: 'Semua Status', value: 'ALL' },
  { label: 'Aktif', value: 'ACTIVE' },
  { label: 'Non-Aktif', value: 'INACTIVE' },
  { label: 'Resign', value: 'RESIGNED' }
]

const genderOptions = [
  { label: 'Laki-laki', value: 'MALE' },
  { label: 'Perempuan', value: 'FEMALE' }
]

const employmentOptions = [
  { label: 'Tetap', value: 'PERMANENT' },
  { label: 'Kontrak', value: 'CONTRACT' },
  { label: 'Probation', value: 'PROBATION' },
  { label: 'Magang', value: 'INTERN' }
]

const activeOptions = [
  { label: 'Aktif', value: 'ACTIVE' },
  { label: 'Non-Aktif', value: 'INACTIVE' },
  { label: 'Resign', value: 'RESIGNED' }
]

const jabatanFiltered = computed(() => {
  if (!form.value.department_id) return []
  return hrStore.jabatanOptions.filter(opt => opt.department_id === form.value.department_id)
})

const activeCount = computed(() => hrStore.karyawanList.filter(k => k.active_status === 'ACTIVE').length)
const contractCount = computed(() => hrStore.karyawanList.filter(k => k.employment_status === 'CONTRACT').length)
const resignedCount = computed(() => hrStore.karyawanList.filter(k => k.active_status === 'RESIGNED').length)

const openForm = (row = null, viewOnly = false) => {
  isView.value = viewOnly
  isEdit.value = !!row && !viewOnly
  photoFile.value = null
  photoFileName.value = ''
  if (row) {
    form.value = { ...row }
  } else {
    form.value = { ...defaultForm }
  }
  isDialogOpen.value = true
}

const onPhotoSelected = (val) => {
  if (val) {
    photoFile.value = val
    photoFileName.value = val.name
  }
}

const onPhotoRemoved = () => {
  photoFile.value = null
  photoFileName.value = ''
}

const onSave = async () => {
  if (isView.value) {
    isDialogOpen.value = false
    return
  }
  const valid = await formRef.value.validate()
  if (!valid) return

  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      await hrStore.updateKaryawan(payload.id, payload)
      notifySuccess('Karyawan berhasil diupdate')
    } else {
      await hrStore.createKaryawan(payload)
      notifySuccess('Karyawan berhasil ditambahkan')
    }
    isDialogOpen.value = false
    onRequest({ pagination: pagination.value })
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan')
  }
}

const handleDelete = (row) => {
  confirmDelete({
    title: 'Hapus Karyawan',
    message: `Apakah Anda yakin ingin menghapus ${row.full_name}?`,
    onConfirm: async () => {
      try {
        await hrStore.deleteKaryawan(row.id)
        notifySuccess('Karyawan berhasil dihapus')
        onRequest({ pagination: pagination.value })
      } catch (error) {
        notifyError(error.message || 'Gagal menghapus karyawan')
      }
    }
  })
}

onMounted(async () => {
  await Promise.all([
    hrStore.fetchDepartemen({ rowsPerPage: 100 }),
    hrStore.fetchJabatan({ rowsPerPage: 100 })
  ])
  onRequest({ pagination: pagination.value })
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

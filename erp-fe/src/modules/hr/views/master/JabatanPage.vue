<script setup>
import { ref, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required } from 'src/utils/validators'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const {
  pagination,
  searchQuery,
  resetPage,
  onRequest
} = useTablePagination()

const levelOptions = [
  { label: 'Staff', value: 'STAFF' },
  { label: 'Supervisor', value: 'SUPERVISOR' },
  { label: 'Manager', value: 'MANAGER' },
  { label: 'Director', value: 'DIRECTOR' }
]

const getLevelVariant = (level) => {
  switch(level) {
    case 'STAFF': return 'neutral'
    case 'SUPERVISOR': return 'info'
    case 'MANAGER': return 'warning'
    case 'DIRECTOR': return 'success'
    default: return 'neutral'
  }
}

const columns = [
  { name: 'position_name', align: 'left', label: 'Nama Jabatan', field: 'position_name', sortable: true },
  { name: 'department_name', align: 'left', label: 'Departemen', field: 'department_name' },
  { name: 'level', align: 'left', label: 'Level', field: 'level' },
  { name: 'actions', align: 'center', label: 'Aksi', field: 'actions' }
]

const formDialog = ref(false)
const isEdit = ref(false)
const editingId = ref(null)

const formData = ref({
  position_name: '',
  department_id: null,
  level: null,
  job_description: ''
})

const resetForm = () => {
  formData.value = {
    position_name: '',
    department_id: null,
    level: null,
    job_description: ''
  }
  isEdit.value = false
  editingId.value = null
}

const loadData = async (page = 1) => {
  try {
    await hrStore.fetchJabatan({
      page,
      limit: pagination.value.rowsPerPage,
      search: searchQuery.value,
      sortBy: pagination.value.sortBy,
      descending: pagination.value.descending
    })
    pagination.value.rowsNumber = hrStore.jabatanMeta?.total || 0
    pagination.value.page = page
  } catch {
    notifyError('Gagal memuat data jabatan')
  }
}

const handleSearch = () => {
  resetPage()
  loadData(1)
}

const handleTableRequest = (props) => {
  onRequest(props)
  loadData(props.pagination.page)
}

const openAddDialog = () => {
  resetForm()
  formDialog.value = true
}

const openEditDialog = (row) => {
  formData.value = {
    position_name: row.position_name,
    department_id: row.department_id,
    level: row.level,
    job_description: row.job_description || ''
  }
  editingId.value = row.id
  isEdit.value = true
  formDialog.value = true
}

const handleSubmit = async () => {
  try {
    if (isEdit.value) {
      await hrStore.updateJabatan(editingId.value, formData.value)
      notifySuccess('Jabatan berhasil diperbarui')
    } else {
      await hrStore.createJabatan(formData.value)
      notifySuccess('Jabatan berhasil ditambahkan')
    }
    formDialog.value = false
    loadData(pagination.value.page)
  } catch (error) { notifyError(error.message || 'Gagal menyimpan data')
  }
}

const handleDelete = (row) => {
  confirmDelete({
    title: 'Hapus Jabatan',
    message: `Apakah Anda yakin ingin menghapus jabatan ${row.position_name}?`,
    onConfirm: async () => {
      try {
        await hrStore.deleteJabatan(row.id)
        notifySuccess('Jabatan berhasil dihapus')
        if (hrStore.jabatanList.length === 1 && pagination.value.page > 1) {
          loadData(pagination.value.page - 1)
        } else {
          loadData(pagination.value.page)
        }
      } catch (error) { notifyError(error.message || 'Gagal menghapus data')
      }
    }
  })
}

onMounted(async () => {
  await hrStore.fetchDepartemen({ limit: 100 })
  loadData()
})
</script>

<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Master Jabatan / Posisi"
      subtitle="Kelola daftar posisi dan jenjang organisasi"
      icon="badge"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    />

    <q-card class="clay-card q-mb-lg q-pa-md">
      <div class="row q-col-gutter-md items-center q-mb-md">
        <div class="col-12 col-sm-6 col-md-4">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            class="clay-input"
            placeholder="Cari jabatan..."
            debounce="350"
            @update:model-value="handleSearch"
          >
            <template v-slot:append>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
        <div class="col-12 col-sm-6 col-md-8 text-right">
          <q-btn
            unelevated
            no-caps
            class="clay-btn clay-btn-primary"
            icon="add"
            label="Tambah Jabatan"
            @click="openAddDialog"
          />
        </div>
      </div>

      <AppDataTable
        :columns="columns"
        :data="hrStore.jabatanList"
        :loading="hrStore.isLoading('jabatan')"
        :pagination="pagination"
        @request="handleTableRequest"
      >
        <template #body-cell-level="props">
          <q-td :props="props">
            <AppBadge :variant="getLevelVariant(props.row.level)">
              {{ props.row.level }}
            </AppBadge>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn
              flat
              dense
              round
              size="sm"
              icon="edit"
              color="primary"
              @click="openEditDialog(props.row)"
            >
              <q-tooltip>Edit</q-tooltip>
            </q-btn>
            <q-btn
              flat
              dense
              round
              size="sm"
              icon="delete"
              color="negative"
              @click="handleDelete(props.row)"
            >
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </AppDataTable>
    </q-card>

    <AppFormDialog
      v-model="formDialog"
      :title="isEdit ? 'Edit Jabatan' : 'Tambah Jabatan'"
      :loading="hrStore.isLoading('jabatan')"
      @submit="handleSubmit"
    >
      <div class="row q-col-gutter-md">
        <div class="col-12">
          <label class="form-label">Nama Jabatan <span class="req">*</span></label>
          <q-input
            v-model="formData.position_name"
            outlined
            dense
            class="clay-input"
            :rules="[val => required(val, 'Nama Jabatan wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="col-12">
          <label class="form-label">Departemen <span class="req">*</span></label>
          <q-select
            v-model="formData.department_id"
            :options="hrStore.departemenOptions"
            outlined
            dense
            class="clay-input"
            emit-value
            map-options
            :rules="[val => required(val, 'Departemen wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="col-12">
          <label class="form-label">Level <span class="req">*</span></label>
          <q-select
            v-model="formData.level"
            :options="levelOptions"
            outlined
            dense
            class="clay-input"
            emit-value
            map-options
            :rules="[val => required(val, 'Level wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="col-12">
          <label class="form-label">Deskripsi Pekerjaan</label>
          <q-input
            v-model="formData.job_description"
            type="textarea"
            outlined
            dense
            class="clay-input"
            rows="3"
          />
        </div>
      </div>
    </AppFormDialog>
  </q-page>
</template>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

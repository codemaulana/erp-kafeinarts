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

const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const {
  pagination,
  searchQuery,
  resetPage,
  onRequest
} = useTablePagination()

const columns = [
  { name: 'department_name', align: 'left', label: 'Nama Departemen', field: 'department_name', sortable: true },
  { name: 'head_employee_name', align: 'left', label: 'Kepala Departemen', field: 'head_employee_name' },
  { name: 'employee_count', align: 'center', label: 'Jumlah Karyawan', field: 'employee_count' },
  { name: 'actions', align: 'center', label: 'Aksi', field: 'actions' }
]

const formDialog = ref(false)
const isEdit = ref(false)
const editingId = ref(null)

const formData = ref({
  department_name: '',
  head_employee_id: null,
  description: ''
})

const resetForm = () => {
  formData.value = {
    department_name: '',
    head_employee_id: null,
    description: ''
  }
  isEdit.value = false
  editingId.value = null
}

const loadData = async (page = 1) => {
  try {
    await hrStore.fetchDepartemen({
      page,
      limit: pagination.value.rowsPerPage,
      search: searchQuery.value,
      sortBy: pagination.value.sortBy,
      descending: pagination.value.descending
    })
    pagination.value.rowsNumber = hrStore.departemenMeta?.total || 0
    pagination.value.page = page
  } catch {
    notifyError('Gagal memuat data departemen')
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
    department_name: row.department_name,
    head_employee_id: row.head_employee_id || null,
    description: row.description || ''
  }
  editingId.value = row.id
  isEdit.value = true
  formDialog.value = true
}

const handleSubmit = async () => {
  try {
    if (isEdit.value) {
      await hrStore.updateDepartemen(editingId.value, formData.value)
      notifySuccess('Departemen berhasil diperbarui')
    } else {
      await hrStore.createDepartemen(formData.value)
      notifySuccess('Departemen berhasil ditambahkan')
    }
    formDialog.value = false
    loadData(pagination.value.page)
  } catch (error) { notifyError(error.message || 'Gagal menyimpan data')
  }
}

const handleDelete = (row) => {
  confirmDelete({
    title: 'Hapus Departemen',
    message: `Apakah Anda yakin ingin menghapus departemen ${row.department_name}?`,
    onConfirm: async () => {
      try {
        await hrStore.deleteDepartemen(row.id)
        notifySuccess('Departemen berhasil dihapus')
        if (hrStore.departemenList.length === 1 && pagination.value.page > 1) {
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
  await hrStore.fetchKaryawan({ limit: 100 })
  loadData()
})
</script>

<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Master Departemen"
      subtitle="Kelola struktur organisasi perusahaan"
      icon="business"
      iconBg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      iconColor="#5B9BD5"
    />

    <q-card class="clay-card q-mb-lg q-pa-md">
      <div class="row q-col-gutter-md items-center q-mb-md">
        <div class="col-12 col-sm-6 col-md-4">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            class="clay-input"
            placeholder="Cari departemen..."
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
            label="Tambah Departemen"
            @click="openAddDialog"
          />
        </div>
      </div>

      <AppDataTable
        :columns="columns"
        :data="hrStore.departemenList"
        :loading="hrStore.isLoading('departemen')"
        :pagination="pagination"
        @request="handleTableRequest"
      >
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
      :title="isEdit ? 'Edit Departemen' : 'Tambah Departemen'"
      :loading="hrStore.isLoading('departemen')"
      @submit="handleSubmit"
    >
      <div class="row q-col-gutter-md">
        <div class="col-12">
          <label class="form-label">Nama Departemen <span class="req">*</span></label>
          <q-input
            v-model="formData.department_name"
            outlined
            dense
            class="clay-input"
            :rules="[val => required(val, 'Nama Departemen wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="col-12">
          <label class="form-label">Kepala Departemen</label>
          <q-select
            v-model="formData.head_employee_id"
            :options="hrStore.karyawanOptions"
            outlined
            dense
            class="clay-input"
            clearable
            emit-value
            map-options
          />
        </div>
        <div class="col-12">
          <label class="form-label">Deskripsi</label>
          <q-input
            v-model="formData.description"
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

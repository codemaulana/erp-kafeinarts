<script setup>
import { ref, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, nonNegative } from 'src/utils/validators'
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

const attachmentOptions = [
  { label: 'Ya', value: true },
  { label: 'Tidak', value: false }
]

const columns = [
  { name: 'leave_name', align: 'left', label: 'Nama Jenis Cuti', field: 'leave_name', sortable: true },
  { name: 'annual_quota_days', align: 'center', label: 'Jatah Hari', field: 'annual_quota_days' },
  { name: 'requires_attachment', align: 'center', label: 'Perlu Lampiran', field: 'requires_attachment' },
  { name: 'actions', align: 'center', label: 'Aksi', field: 'actions' }
]

const formDialog = ref(false)
const isEdit = ref(false)
const editingId = ref(null)

const formData = ref({
  leave_name: '',
  annual_quota_days: 0,
  requires_attachment: false,
  description: ''
})

const resetForm = () => {
  formData.value = {
    leave_name: '',
    annual_quota_days: 0,
    requires_attachment: false,
    description: ''
  }
  isEdit.value = false
  editingId.value = null
}

const loadData = async (page = 1) => {
  try {
    await hrStore.fetchJenisCuti({
      page,
      limit: pagination.value.rowsPerPage,
      search: searchQuery.value,
      sortBy: pagination.value.sortBy,
      descending: pagination.value.descending
    })
    pagination.value.rowsNumber = hrStore.jenisCutiMeta?.total || 0
    pagination.value.page = page
  } catch {
    notifyError('Gagal memuat data jenis cuti')
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
    leave_name: row.leave_name,
    annual_quota_days: row.annual_quota_days,
    requires_attachment: row.requires_attachment,
    description: row.description || ''
  }
  editingId.value = row.id
  isEdit.value = true
  formDialog.value = true
}

const handleSubmit = async () => {
  try {
    const payload = {
      ...formData.value,
      annual_quota_days: Number(formData.value.annual_quota_days)
    }
    
    if (isEdit.value) {
      await hrStore.updateJenisCuti(editingId.value, payload)
      notifySuccess('Jenis Cuti berhasil diperbarui')
    } else {
      await hrStore.createJenisCuti(payload)
      notifySuccess('Jenis Cuti berhasil ditambahkan')
    }
    formDialog.value = false
    loadData(pagination.value.page)
  } catch (error) { notifyError(error.message || 'Gagal menyimpan data')
  }
}

const handleDelete = (row) => {
  confirmDelete({
    title: 'Hapus Jenis Cuti',
    message: `Apakah Anda yakin ingin menghapus jenis cuti ${row.leave_name}?`,
    onConfirm: async () => {
      try {
        await hrStore.deleteJenisCuti(row.id)
        notifySuccess('Jenis Cuti berhasil dihapus')
        if (hrStore.jenisCutiList.length === 1 && pagination.value.page > 1) {
          loadData(pagination.value.page - 1)
        } else {
          loadData(pagination.value.page)
        }
      } catch (error) { notifyError(error.message || 'Gagal menghapus data')
      }
    }
  })
}

onMounted(() => {
  loadData()
})
</script>

<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Master Jenis Cuti / Izin"
      subtitle="Konfigurasi jenis cuti dan jatah hari"
      icon="event_available"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    />

    <q-card class="clay-card q-mb-lg q-pa-md">
      <div class="row q-col-gutter-md items-center q-mb-md">
        <div class="col-12 col-sm-6 col-md-4">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            class="clay-input"
            placeholder="Cari jenis cuti..."
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
            label="Tambah Jenis Cuti"
            @click="openAddDialog"
          />
        </div>
      </div>

      <AppDataTable
        :columns="columns"
        :data="hrStore.jenisCutiList"
        :loading="hrStore.isLoading('jenisCuti')"
        :pagination="pagination"
        @request="handleTableRequest"
      >
        <template #body-cell-annual_quota_days="props">
          <q-td :props="props" class="text-center">
            <AppBadge variant="warning">
              {{ props.row.annual_quota_days }}
            </AppBadge>
          </q-td>
        </template>

        <template #body-cell-requires_attachment="props">
          <q-td :props="props" class="text-center">
            <AppBadge :variant="props.row.requires_attachment ? 'info' : 'neutral'">
              {{ props.row.requires_attachment ? 'Ya' : 'Tidak' }}
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
      :title="isEdit ? 'Edit Jenis Cuti' : 'Tambah Jenis Cuti'"
      :loading="hrStore.isLoading('jenisCuti')"
      @submit="handleSubmit"
    >
      <div class="row q-col-gutter-md">
        <div class="col-12">
          <label class="form-label">Nama Jenis Cuti <span class="req">*</span></label>
          <q-input
            v-model="formData.leave_name"
            outlined
            dense
            class="clay-input"
            :rules="[val => required(val, 'Nama Jenis Cuti wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="col-12 col-md-6">
          <label class="form-label">Jatah Hari <span class="req">*</span></label>
          <q-input
            v-model="formData.annual_quota_days"
            type="number"
            min="0"
            outlined
            dense
            class="clay-input"
            :rules="[
              val => requiredNumber(val, 'Jatah Hari wajib diisi'),
              val => nonNegative(val, 'Jatah Hari harus lebih dari atau sama dengan 0')
            ]"
            lazy-rules
          />
        </div>
        <div class="col-12 col-md-6">
          <label class="form-label">Perlu Lampiran <span class="req">*</span></label>
          <q-select
            v-model="formData.requires_attachment"
            :options="attachmentOptions"
            outlined
            dense
            class="clay-input"
            emit-value
            map-options
            :rules="[val => required(val !== null ? String(val) : '', 'Perlu Lampiran wajib dipilih')]"
            lazy-rules
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

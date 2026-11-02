<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Master Komponen Gaji"
      subtitle="Konfigurasi komponen pendapatan dan potongan payroll"
      icon="payments"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Komponen"
          @click="openForm()"
        />
      </template>
    </AppPageHeader>

    <div class="app-card clay-card q-pa-md">
      <div class="row q-mb-md">
        <q-input
          v-model="searchQuery"
          outlined
          dense
          class="clay-input"
          placeholder="Cari komponen gaji..."
          debounce="350"
          @update:model-value="onSearch"
        >
          <template #prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>

      <AppDataTable
        :rows="hrStore.komponenGajiList"
        :columns="columns"
        :loading="hrStore.isLoading('komponenGaji')"
        :pagination="pagination"
        @request="onRequest"
      >
        <template #body-cell-component_type="props">
          <q-td :props="props">
            <AppBadge
              :variant="props.value === 'EARNING' ? 'success' : 'danger'"
              :label="props.value === 'EARNING' ? 'Pendapatan' : 'Potongan'"
            />
          </q-td>
        </template>
        <template #body-cell-calculation_method="props">
          <q-td :props="props">
            <AppBadge
              :variant="props.value === 'FIXED' ? 'info' : 'warning'"
              :label="props.value === 'FIXED' ? 'Nominal Tetap' : 'Persentase'"
            />
          </q-td>
        </template>
        <template #body-cell-actions="props">
          <q-td :props="props">
            <div class="row justify-center q-gutter-xs">
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
      :title="isEdit ? 'Edit Komponen Gaji' : 'Tambah Komponen Gaji'"
      @save="onSave"
      :loading="hrStore.isLoading('komponenGaji_saving')"
      saveLabel="Simpan"
    >
      <q-form ref="formRef" @submit.prevent="onSave">
        <div class="q-mb-md">
          <label class="form-label">Nama Komponen <span class="req">*</span></label>
          <q-input v-model="form.component_name" outlined dense class="clay-input" :rules="[required('Nama komponen wajib diisi')]" lazy-rules />
        </div>

        <div class="q-mb-md">
          <label class="form-label">Tipe Komponen <span class="req">*</span></label>
          <q-select v-model="form.component_type" :options="typeOptions" outlined dense class="clay-input" emit-value map-options :rules="[required('Tipe komponen wajib dipilih')]" lazy-rules />
        </div>

        <div class="q-mb-md">
          <label class="form-label">Metode Perhitungan <span class="req">*</span></label>
          <q-select v-model="form.calculation_method" :options="methodOptions" outlined dense class="clay-input" emit-value map-options :rules="[required('Metode perhitungan wajib dipilih')]" lazy-rules />
        </div>

        <div class="q-mb-md">
          <label class="form-label">Keterangan</label>
          <q-input v-model="form.description" type="textarea" rows="3" outlined dense class="clay-input" />
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

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

const isDialogOpen = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const searchQuery = ref('')

const defaultForm = {
  id: null,
  component_name: '',
  component_type: null,
  calculation_method: null,
  description: ''
}

const form = ref({ ...defaultForm })

const columns = [
  { name: 'component_name', field: 'component_name', label: 'Nama Komponen', align: 'left', sortable: true },
  { name: 'component_type', field: 'component_type', label: 'Tipe Komponen', align: 'left' },
  { name: 'calculation_method', field: 'calculation_method', label: 'Metode', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchData: (params) => hrStore.fetchKomponenGaji({ ...params, search: searchQuery.value }),
  defaultSort: 'component_name'
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const typeOptions = [
  { label: 'Pendapatan', value: 'EARNING' },
  { label: 'Potongan', value: 'DEDUCTION' }
]

const methodOptions = [
  { label: 'Nominal Tetap', value: 'FIXED' },
  { label: 'Persentase', value: 'PERCENTAGE' }
]

const openForm = (row = null) => {
  isEdit.value = !!row
  if (row) {
    form.value = { ...row }
  } else {
    form.value = { ...defaultForm }
  }
  isDialogOpen.value = true
}

const onSave = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return

  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      await hrStore.updateKomponenGaji(payload.id, payload)
      notifySuccess('Komponen gaji berhasil diupdate')
    } else {
      await hrStore.createKomponenGaji(payload)
      notifySuccess('Komponen gaji berhasil ditambahkan')
    }
    isDialogOpen.value = false
    onRequest({ pagination: pagination.value })
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan')
  }
}

const handleDelete = (row) => {
  confirmDelete({
    title: 'Hapus Komponen Gaji',
    message: `Apakah Anda yakin ingin menghapus ${row.component_name}?`,
    onConfirm: async () => {
      try {
        await hrStore.deleteKomponenGaji(row.id)
        notifySuccess('Komponen gaji berhasil dihapus')
        onRequest({ pagination: pagination.value })
      } catch (error) {
        notifyError(error.message || 'Gagal menghapus komponen gaji')
      }
    }
  })
}

onMounted(() => {
  onRequest({ pagination: pagination.value })
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

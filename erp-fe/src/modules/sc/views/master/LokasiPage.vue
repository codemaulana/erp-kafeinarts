<template>
  <q-page class="app-page clay-animate-up">
    <!-- Header -->
    <AppPageHeader
      title="Lokasi / Gudang Penyimpanan"
      subtitle="Lokasi fisik penyimpanan aset"
      icon="warehouse"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    />

    <!-- Filter & Toolbar -->
    <div class="row q-col-gutter-md q-mb-md items-center">
      <div class="col-12 col-md-4">
        <q-input
          v-model="search"
          outlined
          dense
          class="clay-input"
          placeholder="Cari lokasi..."
          debounce="350"
          @update:model-value="onSearch"
        >
          <template #append>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>
      <div class="col-12 col-md-8 text-right">
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Lokasi"
          @click="openAddDialog"
        />
      </div>
    </div>

    <!-- Data Table -->
    <AppDataTable
      :columns="columns"
      :data="scStore.lokasiList"
      :loading="scStore.isLoading('lokasi')"
      :pagination="pagination"
      row-key="id"
      @request="onRequest"
    >
      <template #body-cell-actions="props">
        <q-td :props="props">
          <q-btn
            flat
            dense
            round
            size="sm"
            color="primary"
            icon="edit"
            @click="openEditDialog(props.row)"
          >
            <q-tooltip>Edit</q-tooltip>
          </q-btn>
          <q-btn
            flat
            dense
            round
            size="sm"
            color="negative"
            icon="delete"
            @click="onDelete(props.row)"
          >
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </AppDataTable>

    <!-- Dialog Form -->
    <AppFormDialog
      v-model="dialog"
      :title="isEdit ? 'Edit Lokasi' : 'Tambah Lokasi'"
      :loading="scStore.isLoading('lokasi')"
      @submit="onSubmit"
    >
      <div class="q-mb-md">
        <label class="form-label">Nama Lokasi <span class="req">*</span></label>
        <q-input
          v-model="form.location_name"
          outlined
          dense
          class="clay-input"
          :rules="[required('Nama lokasi wajib diisi')]"
          lazy-rules
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Alamat</label>
        <q-input
          v-model="form.address"
          type="textarea"
          outlined
          dense
          class="clay-input"
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Penanggung Jawab</label>
        <q-select
          v-model="form.person_in_charge_id"
          outlined
          dense
          class="clay-input"
          :options="hrStore.karyawanOptions"
          emit-value
          map-options
        />
      </div>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useScStore } from 'src/stores/sc.store'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useNotify } from 'src/composables/useNotify'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { required } from 'src/utils/validators'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'

const scStore = useScStore()
const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const search = ref('')
const dialog = ref(false)
const isEdit = ref(false)

const form = ref({
  id: null,
  location_name: '',
  address: '',
  person_in_charge_id: null
})

const columns = [
  { name: 'location_name', field: 'location_name', label: 'Nama Lokasi', align: 'left', sortable: true },
  { name: 'address', field: 'address', label: 'Alamat', align: 'left' },
  { name: 'person_in_charge_name', field: 'person_in_charge_name', label: 'Penanggung Jawab', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchFn: (params) => {
    const qParams = { ...params }
    if (search.value) qParams.search = search.value
    return scStore.fetchLokasi(qParams)
  }
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const resetForm = () => {
  form.value = {
    id: null,
    location_name: '',
    address: '',
    person_in_charge_id: null
  }
}

const openAddDialog = () => {
  isEdit.value = false
  resetForm()
  dialog.value = true
}

const openEditDialog = (row) => {
  isEdit.value = true
  form.value = { ...row }
  dialog.value = true
}

const onSubmit = async () => {
  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      await scStore.updateLokasi(payload.id, payload)
      notifySuccess('Lokasi berhasil diupdate')
    } else {
      await scStore.createLokasi(payload)
      notifySuccess('Lokasi berhasil ditambahkan')
    }
    dialog.value = false
    onRequest({ pagination: pagination.value })
  } catch (err) {
    notifyError(err.message || 'Gagal menyimpan lokasi')
  }
}

const onDelete = (row) => {
  confirmDelete('Lokasi ini', async () => {
    try {
      await scStore.deleteLokasi(row.id)
      notifySuccess('Lokasi berhasil dihapus')
      onRequest({ pagination: pagination.value })
    } catch (err) {
      notifyError(err.message || 'Gagal menghapus lokasi')
    }
  })
}

onMounted(async () => {
  await hrStore.fetchKaryawan({ perPage: 100 })
  onRequest({ pagination: pagination.value })
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

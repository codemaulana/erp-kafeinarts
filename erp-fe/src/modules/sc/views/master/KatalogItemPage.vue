<template>
  <q-page class="app-page clay-animate-up">
    <!-- Header -->
    <AppPageHeader
      title="Master Katalog Item"
      subtitle="Daftar barang/jasa/lisensi referensi pengadaan"
      icon="list_alt"
      iconBg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      iconColor="#5B9BD5"
    />

    <!-- Filter & Toolbar -->
    <div class="row q-col-gutter-md q-mb-md items-center">
      <div class="col-12 col-md-4">
        <q-input
          v-model="search"
          outlined
          dense
          class="clay-input"
          placeholder="Cari item..."
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
          label="Tambah Item"
          @click="openAddDialog"
        />
      </div>
    </div>

    <!-- Data Table -->
    <AppDataTable
      :columns="columns"
      :data="scStore.itemList"
      :loading="scStore.isLoading('item')"
      :pagination="pagination"
      row-key="id"
      @request="onRequest"
    >
      <template #body-cell-estimated_price="props">
        <q-td :props="props">
          {{ formatCurrency(props.row.estimated_price) }}
        </q-td>
      </template>

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
      :title="isEdit ? 'Edit Item' : 'Tambah Item'"
      :loading="scStore.isLoading('item')"
      @submit="onSubmit"
    >
      <div class="q-mb-md">
        <label class="form-label">Nama Item <span class="req">*</span></label>
        <q-input
          v-model="form.item_name"
          outlined
          dense
          class="clay-input"
          :rules="[required('Nama item wajib diisi')]"
          lazy-rules
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Kategori <span class="req">*</span></label>
        <q-select
          v-model="form.category_id"
          outlined
          dense
          class="clay-input"
          :options="scStore.kategoriItemOptions"
          :rules="[required('Kategori wajib diisi')]"
          lazy-rules
          emit-value
          map-options
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Satuan <span class="req">*</span></label>
        <q-select
          v-model="form.unit"
          outlined
          dense
          class="clay-input"
          :options="['Unit', 'Lisensi', 'Tahun', 'Bulan', 'Paket']"
          :rules="[required('Satuan wajib diisi')]"
          lazy-rules
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Harga Estimasi</label>
        <q-input
          v-model.number="form.estimated_price"
          type="number"
          min="0"
          outlined
          dense
          class="clay-input"
          :rules="form.estimated_price ? [nonNegative('Harga tidak boleh negatif')] : []"
          lazy-rules
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Vendor Default</label>
        <ScVendorSelect
          v-model="form.default_vendor_id"
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Deskripsi</label>
        <q-input
          v-model="form.description"
          type="textarea"
          outlined
          dense
          class="clay-input"
        />
      </div>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useScStore } from 'src/stores/sc.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useNotify } from 'src/composables/useNotify'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { required, nonNegative } from 'src/utils/validators'
import { formatCurrency } from 'src/utils/formatter'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import ScVendorSelect from 'src/modules/sc/components/ScVendorSelect.vue'

const scStore = useScStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const search = ref('')
const dialog = ref(false)
const isEdit = ref(false)

const form = ref({
  id: null,
  item_name: '',
  category_id: null,
  unit: '',
  estimated_price: null,
  default_vendor_id: null,
  description: ''
})

const columns = [
  { name: 'item_name', field: 'item_name', label: 'Nama Item', align: 'left', sortable: true },
  { name: 'category_name', field: 'category_name', label: 'Kategori', align: 'left' },
  { name: 'unit', field: 'unit', label: 'Satuan', align: 'center' },
  { name: 'estimated_price', field: 'estimated_price', label: 'Harga Estimasi', align: 'right' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchFn: (params) => {
    const qParams = { ...params }
    if (search.value) qParams.search = search.value
    return scStore.fetchItem(qParams)
  }
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const resetForm = () => {
  form.value = {
    id: null,
    item_name: '',
    category_id: null,
    unit: '',
    estimated_price: null,
    default_vendor_id: null,
    description: ''
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
      await scStore.updateItem(payload.id, payload)
      notifySuccess('Item berhasil diupdate')
    } else {
      await scStore.createItem(payload)
      notifySuccess('Item berhasil ditambahkan')
    }
    dialog.value = false
    onRequest({ pagination: pagination.value })
  } catch (err) {
    notifyError(err.message || 'Gagal menyimpan item')
  }
}

const onDelete = (row) => {
  confirmDelete('Item ini', async () => {
    try {
      await scStore.deleteItem(row.id)
      notifySuccess('Item berhasil dihapus')
      onRequest({ pagination: pagination.value })
    } catch (err) {
      notifyError(err.message || 'Gagal menghapus item')
    }
  })
}

onMounted(async () => {
  await scStore.fetchKategoriItem({ perPage: 100 })
  onRequest({ pagination: pagination.value })
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

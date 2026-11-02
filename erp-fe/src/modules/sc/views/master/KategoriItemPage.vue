<template>
  <q-page class="app-page clay-animate-up">
    <!-- Header -->
    <AppPageHeader
      title="Kategori Item / Aset"
      subtitle="Pengelompokan jenis barang/jasa"
      icon="category"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    />

    <!-- Filter & Toolbar -->
    <div class="row q-col-gutter-md q-mb-md items-center">
      <div class="col-12 col-md-4">
        <q-input
          v-model="search"
          outlined
          dense
          class="clay-input"
          placeholder="Cari kategori..."
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
          label="Tambah Kategori"
          @click="openAddDialog"
        />
      </div>
    </div>

    <!-- Data Table -->
    <AppDataTable
      :columns="columns"
      :data="scStore.kategoriItemList"
      :loading="scStore.isLoading('kategoriItem')"
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
      :title="isEdit ? 'Edit Kategori' : 'Tambah Kategori'"
      :loading="scStore.isLoading('kategoriItem')"
      @submit="onSubmit"
    >
      <div class="q-mb-md">
        <label class="form-label">Nama Kategori <span class="req">*</span></label>
        <q-input
          v-model="form.category_name"
          outlined
          dense
          class="clay-input"
          :rules="[required('Nama kategori wajib diisi')]"
          lazy-rules
        />
      </div>
      <div class="q-mb-md">
        <label class="form-label">Jenis Kategori <span class="req">*</span></label>
        <q-select
          v-model="form.item_type"
          outlined
          dense
          class="clay-input"
          :options="['HARDWARE', 'SOFTWARE_LICENSE', 'SUBSCRIPTION', 'SERVICE']"
          :rules="[required('Jenis kategori wajib diisi')]"
          lazy-rules
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
import { required } from 'src/utils/validators'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'

const scStore = useScStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const search = ref('')
const dialog = ref(false)
const isEdit = ref(false)

const form = ref({
  id: null,
  category_name: '',
  item_type: '',
  description: ''
})

const columns = [
  { name: 'category_name', field: 'category_name', label: 'Nama Kategori', align: 'left', sortable: true },
  { name: 'item_type', field: 'item_type', label: 'Jenis Kategori', align: 'left' },
  { name: 'description', field: 'description', label: 'Deskripsi', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchFn: (params) => {
    const qParams = { ...params }
    if (search.value) qParams.search = search.value
    return scStore.fetchKategoriItem(qParams)
  }
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const resetForm = () => {
  form.value = {
    id: null,
    category_name: '',
    item_type: '',
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
      await scStore.updateKategoriItem(payload.id, payload)
      notifySuccess('Kategori berhasil diupdate')
    } else {
      await scStore.createKategoriItem(payload)
      notifySuccess('Kategori berhasil ditambahkan')
    }
    dialog.value = false
    onRequest({ pagination: pagination.value })
  } catch (err) {
    notifyError(err.message || 'Gagal menyimpan kategori')
  }
}

const onDelete = (row) => {
  confirmDelete('Kategori ini', async () => {
    try {
      await scStore.deleteKategoriItem(row.id)
      notifySuccess('Kategori berhasil dihapus')
      onRequest({ pagination: pagination.value })
    } catch (err) {
      notifyError(err.message || 'Gagal menghapus kategori')
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

<template>
  <q-page class="q-pa-md">
    <div class="row items-center justify-between q-mb-md">
      <div class="text-h6">Daftar Produk</div>
      <q-btn color="primary" icon="add" label="Tambah Produk" @click="onCreate" />
    </div>

    <q-table
      :rows="inventoryStore.products"
      :columns="columns"
      row-key="id"
      :loading="inventoryStore.isLoading"
      :pagination="pagination"
      binary-state-sort
      @request="onRequest"
    >
      <template #body-cell-actions="props">
        <q-td :props="props" class="q-gutter-x-sm">
          <q-btn dense flat icon="edit" @click="onEdit(props.row)" />
          <q-btn dense flat icon="delete" color="negative" @click="handleDelete(props.row)" />
        </q-td>
      </template>
    </q-table>

    <!-- Dialog Form -->
    <ProductFormDialog v-model="dialogOpen" :edit-id="selectedId" @saved="load" ref="dialogRef" />
  </q-page>
</template>

<script setup>
import { reactive, ref, onMounted, nextTick } from 'vue'
import { useInventoryStore } from 'src/stores/inventory.store'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import ProductFormDialog from './ProductFormDialog.vue'

const inventoryStore = useInventoryStore()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const pagination = reactive({ page: 1, rowsPerPage: 10, rowsNumber: 0 })

const columns = [
  { name: 'sku', label: 'SKU', field: 'sku', align: 'left', sortable: true },
  { name: 'name', label: 'Nama Produk', field: 'name', align: 'left', sortable: true },
  { name: 'stock', label: 'Stok', field: 'stock', align: 'right', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

// Modal state
const dialogOpen = ref(false)
const selectedId = ref(null)
const dialogRef = ref(null)

async function load(page = pagination.page) {
  await inventoryStore.fetchProducts({ page, perPage: pagination.rowsPerPage })
  pagination.rowsNumber = inventoryStore.meta.total || 0
}

async function onRequest(props) {
  pagination.page = props.pagination.page
  pagination.rowsPerPage = props.pagination.rowsPerPage
  await load(pagination.page)
}

function onCreate() {
  selectedId.value = null
  dialogOpen.value = true
  nextTick(() => {
    dialogRef.value?.loadData()
  })
}

function onEdit(row) {
  selectedId.value = row.id
  dialogOpen.value = true
  nextTick(() => {
    dialogRef.value?.loadData()
  })
}

function handleDelete(row) {
  confirmDelete({
    itemName: row.name,
    onConfirm: async () => {
      await inventoryStore.deleteProduct(row.id)
      notifySuccess('Produk berhasil dihapus')
      load()
    },
    onError: () => notifyError('Gagal menghapus produk'),
  })
}

onMounted(() => load())
</script>

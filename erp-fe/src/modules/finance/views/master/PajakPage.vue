<template>
  <q-page class="app-page clay-animate-up">

    <AppPageHeader
      title="Master Pajak"
      subtitle="Konfigurasi tarif pajak yang digunakan dalam transaksi"
      icon="percent"
      icon-bg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      icon-color="#E8A547"
    >
      <template #actions>
        <q-btn unelevated no-caps label="Tambah Pajak" icon="add" class="clay-btn clay-btn-primary" @click="openDialog()" />
      </template>
    </AppPageHeader>

    <!-- Search -->
    <div class="clay-card clay-card--sm q-mb-md" style="max-width:400px;">
      <q-input v-model="search" placeholder="Cari nama pajak, persentase..." outlined dense clearable class="clay-input" debounce="350" @update:model-value="onSearch">
        <template #prepend><q-icon name="search" size="18px" color="grey-5" /></template>
      </q-input>
    </div>

    <AppDataTable
      :rows="financeStore.pajakList"
      :columns="columns"
      :loading="financeStore.isLoading('pajak')"
      :pagination="pagination"
      no-data-label="Belum ada data pajak"
      @request="handleRequest"
    >
      <template #body-cell-persentase="props">
        <q-td :props="props">
          <AppBadge variant="warning" style="font-size:13px; font-weight:700;">{{ props.value }}%</AppBadge>
        </q-td>
      </template>
      <template #body-cell-actions="props">
        <q-td :props="props" class="text-center">
          <q-btn flat dense round icon="edit" size="sm" color="blue-6" class="q-mr-xs" @click="openDialog(props.row)" />
          <q-btn flat dense round icon="delete_outline" size="sm" color="red-5" @click="handleDelete(props.row)" />
        </q-td>
      </template>
    </AppDataTable>

    <!-- Dialog -->
    <AppFormDialog
      v-model="dialogOpen"
      :title="isEdit ? 'Edit Pajak' : 'Tambah Pajak Baru'"
      :icon="isEdit ? 'edit' : 'add_circle'"
      :icon-color="isEdit ? 'blue-6' : 'green-6'"
      :submit-label="isEdit ? 'Simpan Perubahan' : 'Tambah Pajak'"
      :loading="submitting"
      min-width="400px"
      max-width="480px"
      @submit="onSubmit"
    >
      <div class="form-field-group">
        <label class="form-label">Nama Pajak <span class="req">*</span></label>
        <q-input v-model="form.namaPajak" outlined dense placeholder="Contoh: PPN 11%" class="clay-input"
          :rules="[required('Nama pajak wajib diisi')]" lazy-rules />
      </div>
      <div class="form-field-group">
        <label class="form-label">Persentase (%) <span class="req">*</span></label>
        <q-input v-model.number="form.persentase" outlined dense type="number" min="0" max="100" step="0.1" placeholder="Contoh: 11" class="clay-input"
          :rules="[requiredNumber('Persentase wajib diisi'), inRange(0, 100)]" lazy-rules>
          <template #append><span style="font-size:13px; color:var(--clay-text-muted);">%</span></template>
        </q-input>
      </div>
      <div class="form-field-group">
        <label class="form-label">Keterangan</label>
        <q-input v-model="form.keterangan" outlined dense type="textarea" rows="3" placeholder="Keterangan tambahan (opsional)" class="clay-input" />
      </div>
    </AppFormDialog>

  </q-page>
</template>

<script setup>
import { ref, reactive, onMounted } from 'vue'
import { useFinanceStore } from 'src/stores/finance.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, inRange } from 'src/utils/validators'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const financeStore = useFinanceStore()
const { pagination, syncTotal, onRequest, resetPage } = useTablePagination()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const search     = ref('')
const dialogOpen = ref(false)
const isEdit     = ref(false)
const submitting = ref(false)
const editId     = ref(null)
const form       = reactive({ namaPajak: '', persentase: 0, keterangan: '' })

const columns = [
  { name: 'namaPajak',   label: 'Nama Pajak',     field: 'namaPajak',   align: 'left', sortable: true },
  { name: 'persentase',  label: 'Persentase (%)', field: 'persentase',  align: 'center', sortable: true },
  { name: 'keterangan',  label: 'Keterangan',      field: 'keterangan',  align: 'left' },
  { name: 'actions',     label: 'Aksi',            field: 'actions',     align: 'center' }
]

async function load (page = pagination.page) {
  await financeStore.fetchPajak({ page, perPage: pagination.rowsPerPage, search: search.value })
  syncTotal(financeStore.pajakMeta.total)
}

function onSearch () { resetPage(); load(1) }

async function handleRequest (props) {
  await onRequest(props, ({ page, perPage }) => financeStore.fetchPajak({ page, perPage, search: search.value }))
  syncTotal(financeStore.pajakMeta.total)
}

onMounted(() => load())

function openDialog (row = null) {
  isEdit.value = !!row
  editId.value = row?.id || null
  Object.assign(form, row
    ? { namaPajak: row.namaPajak, persentase: row.persentase, keterangan: row.keterangan || '' }
    : { namaPajak: '', persentase: 0, keterangan: '' }
  )
  dialogOpen.value = true
}

async function onSubmit () {
  submitting.value = true
  try {
    if (isEdit.value) {
      await financeStore.updatePajak(editId.value, { ...form })
      notifySuccess('Pajak berhasil diperbarui')
    } else {
      await financeStore.createPajak({ ...form })
      notifySuccess('Pajak berhasil ditambahkan')
    }
    dialogOpen.value = false
    load()
  } catch (e) {
    notifyError(e.message || 'Terjadi kesalahan')
  } finally {
    submitting.value = false
  }
}

function handleDelete (row) {
  confirmDelete({
    itemName: row.namaPajak,
    onConfirm: async () => {
      await financeStore.deletePajak(row.id)
      notifySuccess('Pajak berhasil dihapus')
      load()
    }
  })
}
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

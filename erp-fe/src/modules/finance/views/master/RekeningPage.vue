<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Rekening / Kas Bank"
      subtitle="Kelola semua akun rekening dan kas perusahaan"
      icon="account_balance"
      icon-bg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      icon-color="#4CAF8A"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          label="Tambah Rekening"
          icon="add"
          class="clay-btn clay-btn-primary"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <AppDataTable
      :rows="financeStore.rekeningList"
      :columns="columns"
      :loading="financeStore.isLoading('rekening')"
      :pagination="pagination"
      no-data-label="Belum ada data rekening"
      @request="handleRequest"
    >
      <template #filter>
        <q-input
          v-model="search"
          placeholder="Cari rekening, bank, nama..."
          outlined
          dense
          clearable
          class="clay-input"
          style="max-width: 400px"
          debounce="350"
          @update:model-value="onSearch"
        >
          <template #prepend><q-icon name="search" size="18px" color="grey-5" /></template>
        </q-input>
      </template>

      <template #body-cell-saldoAwal="props">
        <q-td :props="props">
          <span style="font-weight: 600; color: var(--clay-success)">{{
            formatCurrency(props.value)
          }}</span>
        </q-td>
      </template>

      <template #body-cell-actions="props">
        <q-td :props="props" class="text-center">
          <q-btn
            flat
            dense
            round
            icon="edit"
            size="sm"
            color="blue-6"
            class="q-mr-xs"
            @click="openDialog(props.row)"
            ><q-tooltip>Edit</q-tooltip></q-btn
          >
          <q-btn
            flat
            dense
            round
            icon="delete_outline"
            size="sm"
            color="red-5"
            @click="handleDelete(props.row)"
            ><q-tooltip>Hapus</q-tooltip></q-btn
          >
        </q-td>
      </template>
    </AppDataTable>

    <!-- Dialog -->
    <AppFormDialog
      v-model="dialogOpen"
      :title="isEdit ? 'Edit Rekening' : 'Tambah Rekening Baru'"
      :icon="isEdit ? 'edit' : 'add_circle'"
      :icon-color="isEdit ? 'blue-6' : 'green-6'"
      :submit-label="isEdit ? 'Simpan Perubahan' : 'Tambah Rekening'"
      :loading="submitting"
      min-width="440px"
      max-width="540px"
      @submit="onSubmit"
    >
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">Nama Bank / Kas <span class="req">*</span></label>
          <q-input
            v-model="form.namaBank"
            outlined
            dense
            placeholder="Contoh: BCA, Mandiri, Kas Tunai"
            class="clay-input"
            :rules="[required('Nama bank wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">No. Rekening <span class="req">*</span></label>
          <q-input
            v-model="form.noRekening"
            outlined
            dense
            placeholder="Contoh: 1234567890"
            class="clay-input"
            :rules="[required('No. rekening wajib diisi')]"
            lazy-rules
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">Atas Nama <span class="req">*</span></label>
          <q-input
            v-model="form.atasNama"
            outlined
            dense
            placeholder="Nama pemilik rekening"
            class="clay-input"
            :rules="[required('Atas nama wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Saldo Awal (Rp) <span class="req">*</span></label>
          <q-input
            v-model.number="form.saldoAwal"
            outlined
            dense
            type="number"
            min="0"
            placeholder="0"
            class="clay-input"
            :rules="[requiredNumber('Saldo awal wajib diisi'), nonNegative()]"
            lazy-rules
          />
        </div>
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
import { required, requiredNumber, nonNegative } from 'src/utils/validators'
import { formatCurrency } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'

const financeStore = useFinanceStore()
const { pagination, syncTotal, onRequest, resetPage } = useTablePagination()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const search = ref('')
const dialogOpen = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editId = ref(null)
const form = reactive({ namaBank: '', noRekening: '', atasNama: '', saldoAwal: 0 })

const columns = [
  { name: 'namaBank', label: 'Nama Bank / Kas', field: 'namaBank', align: 'left', sortable: true },
  { name: 'noRekening', label: 'No. Rekening', field: 'noRekening', align: 'left' },
  { name: 'atasNama', label: 'Atas Nama', field: 'atasNama', align: 'left' },
  { name: 'saldoAwal', label: 'Saldo Awal', field: 'saldoAwal', align: 'right', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

async function load(page = pagination.page) {
  await financeStore.fetchRekening({ page, perPage: pagination.rowsPerPage, search: search.value })
  syncTotal(financeStore.rekeningMeta.total)
}

function onSearch() {
  resetPage()
  load(1)
}

async function handleRequest(props) {
  await onRequest(props, ({ page, perPage }) =>
    financeStore.fetchRekening({ page, perPage, search: search.value }),
  )
  syncTotal(financeStore.rekeningMeta.total)
}

onMounted(() => load())

function openDialog(row = null) {
  isEdit.value = !!row
  editId.value = row?.id || null
  Object.assign(
    form,
    row
      ? {
          namaBank: row.namaBank,
          noRekening: row.noRekening,
          atasNama: row.atasNama,
          saldoAwal: row.saldoAwal,
        }
      : { namaBank: '', noRekening: '', atasNama: '', saldoAwal: 0 },
  )
  dialogOpen.value = true
}

async function onSubmit() {
  submitting.value = true
  try {
    if (isEdit.value) {
      await financeStore.updateRekening(editId.value, { ...form })
      notifySuccess('Rekening berhasil diperbarui')
    } else {
      await financeStore.createRekening({ ...form })
      notifySuccess('Rekening berhasil ditambahkan')
    }
    dialogOpen.value = false
    load()
  } catch (e) {
    notifyError(e.message || 'Terjadi kesalahan')
  } finally {
    submitting.value = false
  }
}

function handleDelete(row) {
  confirmDelete({
    itemName: row.namaBank,
    onConfirm: async () => {
      await financeStore.deleteRekening(row.id)
      notifySuccess('Rekening berhasil dihapus')
      load()
    },
  })
}
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

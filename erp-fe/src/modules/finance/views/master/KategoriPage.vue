<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Kategori Transaksi"
      subtitle="Kelola klasifikasi pendapatan dan pengeluaran perusahaan"
      icon="category"
      icon-bg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      icon-color="#5B9BD5"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          label="Tambah Kategori"
          icon="add"
          class="clay-btn clay-btn-primary"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <AppDataTable
      :rows="financeStore.kategoriList"
      :columns="columns"
      :loading="financeStore.isLoading('kategori')"
      :pagination="pagination"
      no-data-label="Belum ada data kategori"
      @request="handleRequest"
    >
      <template #filter>
        <q-input
          v-model="search"
          placeholder="Cari kategori, tipe kas..."
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

      <!-- Tipe Kas badge -->
      <template #body-cell-tipeKas="props">
        <q-td :props="props">
          <AppBadge :variant="props.value === 'Masuk' ? 'success' : 'danger'">
            <q-icon
              :name="props.value === 'Masuk' ? 'arrow_downward' : 'arrow_upward'"
              size="11px"
            />
            {{ props.value }}
          </AppBadge>
        </q-td>
      </template>

      <!-- Actions -->
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
          >
            <q-tooltip>Edit</q-tooltip>
          </q-btn>
          <q-btn
            flat
            dense
            round
            icon="delete_outline"
            size="sm"
            color="red-5"
            @click="handleDelete(props.row)"
          >
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </AppDataTable>

    <!-- ── Dialog Form ── -->
    <AppFormDialog
      v-model="dialogOpen"
      :title="isEdit ? 'Edit Kategori' : 'Tambah Kategori Baru'"
      :icon="isEdit ? 'edit' : 'add_circle'"
      :icon-color="isEdit ? 'blue-6' : 'green-6'"
      :submit-label="isEdit ? 'Simpan Perubahan' : 'Tambah Kategori'"
      :loading="submitting"
      min-width="420px"
      max-width="520px"
      @submit="onSubmit"
    >
      <div class="form-field-group">
        <label class="form-label">Nama Kategori <span class="req">*</span></label>
        <q-input
          v-model="form.namaKategori"
          outlined
          dense
          placeholder="Contoh: Pendapatan Jasa"
          class="clay-input"
          :rules="[required('Nama kategori wajib diisi')]"
          lazy-rules
        />
      </div>

      <div class="form-field-group">
        <label class="form-label">Tipe Kas <span class="req">*</span></label>
        <q-select
          v-model="form.tipeKas"
          :options="tipeOptions"
          outlined
          dense
          emit-value
          map-options
          class="clay-input"
          label="Pilih Tipe Kas"
          :rules="[required('Tipe kas wajib dipilih')]"
          lazy-rules
        >
          <template #selected-item="scope">
            <AppBadge
              :variant="scope.opt.value === 'Masuk' ? 'success' : 'danger'"
              style="font-size: 12px"
            >
              {{ scope.opt.label }}
            </AppBadge>
          </template>
        </q-select>
      </div>

      <div class="form-field-group">
        <label class="form-label">Deskripsi</label>
        <q-input
          v-model="form.deskripsi"
          outlined
          dense
          type="textarea"
          rows="3"
          placeholder="Deskripsi singkat (opsional)"
          class="clay-input"
        />
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
import { required } from 'src/utils/validators'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const financeStore = useFinanceStore()
const { pagination, syncTotal, onRequest, resetPage } = useTablePagination()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

// ── State ──────────────────────────────────────────────────────
const search = ref('')
const dialogOpen = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editId = ref(null)
const form = reactive({ namaKategori: '', tipeKas: '', deskripsi: '' })

const tipeOptions = [
  { label: '↓ Masuk (Pemasukan)', value: 'Masuk' },
  { label: '↑ Keluar (Pengeluaran)', value: 'Keluar' },
]

const columns = [
  {
    name: 'namaKategori',
    label: 'Nama Kategori',
    field: 'namaKategori',
    align: 'left',
    sortable: true,
  },
  { name: 'tipeKas', label: 'Tipe Kas', field: 'tipeKas', align: 'left', sortable: true },
  { name: 'deskripsi', label: 'Deskripsi', field: 'deskripsi', align: 'left', sortable: false },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center', sortable: false },
]

// ── Data loading ───────────────────────────────────────────────
async function load(page = pagination.page) {
  await financeStore.fetchKategori({ page, perPage: pagination.rowsPerPage, search: search.value })
  syncTotal(financeStore.kategoriMeta.total)
}

function onSearch() {
  resetPage()
  load(1)
}

async function handleRequest(props) {
  await onRequest(props, ({ page, perPage }) =>
    financeStore.fetchKategori({ page, perPage, search: search.value }),
  )
  syncTotal(financeStore.kategoriMeta.total)
}

onMounted(() => load())

// ── Dialog ─────────────────────────────────────────────────────
function openDialog(row = null) {
  isEdit.value = !!row
  editId.value = row?.id || null
  Object.assign(
    form,
    row
      ? { namaKategori: row.namaKategori, tipeKas: row.tipeKas, deskripsi: row.deskripsi || '' }
      : { namaKategori: '', tipeKas: '', deskripsi: '' },
  )
  dialogOpen.value = true
}

async function onSubmit() {
  submitting.value = true
  try {
    if (isEdit.value) {
      await financeStore.updateKategori(editId.value, { ...form })
      notifySuccess('Kategori berhasil diperbarui')
    } else {
      await financeStore.createKategori({ ...form })
      notifySuccess('Kategori berhasil ditambahkan')
    }
    dialogOpen.value = false
    load()
  } catch (e) {
    notifyError(e.message || 'Terjadi kesalahan')
  } finally {
    submitting.value = false
  }
}

// ── Delete ─────────────────────────────────────────────────────
function handleDelete(row) {
  confirmDelete({
    itemName: row.namaKategori,
    onConfirm: async () => {
      await financeStore.deleteKategori(row.id)
      notifySuccess('Kategori berhasil dihapus')
      load()
    },
  })
}
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Jurnal Kas"
      subtitle="Catat semua transaksi arus kas masuk dan keluar"
      icon="receipt_long"
      icon-bg="linear-gradient(135deg,#F0EAFB,#D9C8F0)"
      icon-color="#9B59D5"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          label="Catat Transaksi"
          icon="add"
          class="clay-btn clay-btn-primary"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <!-- Summary Cards -->
    <div class="jurnal-summary q-mb-lg">
      <AppStatCard
        icon="arrow_downward"
        icon-bg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
        icon-color="#4CAF8A"
        label="Total Masuk"
        :value="formatCurrency(totalMasuk)"
        value-color="var(--clay-success)"
      />
      <AppStatCard
        icon="arrow_upward"
        icon-bg="linear-gradient(135deg,#FEECF0,#F9BDC8)"
        icon-color="#E85D75"
        label="Total Keluar"
        :value="formatCurrency(totalKeluar)"
        value-color="var(--clay-danger)"
      />
      <AppStatCard
        icon="account_balance_wallet"
        icon-bg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
        icon-color="#5B9BD5"
        label="Saldo Bersih"
        :value="formatCurrency(saldoBersih)"
        :value-color="saldoBersih >= 0 ? 'var(--clay-success)' : 'var(--clay-danger)'"
      />
    </div>

    <AppDataTable
      :rows="financeStore.jurnalList"
      :columns="columns"
      :loading="financeStore.isLoading('jurnal')"
      :pagination="pagination"
      no-data-label="Belum ada data jurnal kas"
      @request="handleRequest"
    >
      <template #filter>
        <div class="filter-row">
          <q-input
            v-model="search"
            placeholder="Cari keterangan, kategori, bank..."
            outlined
            dense
            clearable
            class="clay-input filter-search"
            debounce="350"
            @update:model-value="onFilter"
          >
            <template #prepend><q-icon name="search" size="18px" color="grey-5" /></template>
          </q-input>
          <q-select
            v-model="filterTipe"
            :options="tipeFilterOptions"
            outlined
            dense
            clearable
            emit-value
            map-options
            class="clay-input filter-select"
            label="Pilih Tipe"
            @update:model-value="onFilter"
          />
          <q-select
            v-model="filterBank"
            :options="financeStore.rekeningOptions"
            outlined
            dense
            clearable
            emit-value
            map-options
            class="clay-input filter-select"
            label="Pilih Rekening"
            @update:model-value="onFilter"
          />
        </div>
      </template>

      <template #body-cell-tipe="props">
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

      <template #body-cell-nominal="props">
        <q-td :props="props">
          <span
            :style="{
              fontWeight: 600,
              color: props.row.tipe === 'Masuk' ? 'var(--clay-success)' : 'var(--clay-danger)',
            }"
          >
            {{ props.row.tipe === 'Masuk' ? '+' : '-' }} {{ formatCurrency(props.value) }}
          </span>
        </q-td>
      </template>

      <template #body-cell-buktiFile="props">
        <q-td :props="props" class="text-center">
          <AppBadge v-if="props.value" variant="info"
            ><q-icon name="attach_file" size="11px" /> Ada</AppBadge
          >
          <AppBadge v-else variant="neutral">—</AppBadge>
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

    <!-- Dialog Form -->
    <AppFormDialog
      v-model="dialogOpen"
      :title="isEdit ? 'Edit Jurnal Kas' : 'Catat Transaksi Baru'"
      :icon="isEdit ? 'edit' : 'add_circle'"
      :icon-color="isEdit ? 'blue-6' : 'purple-6'"
      :submit-label="isEdit ? 'Simpan Perubahan' : 'Catat Transaksi'"
      :loading="submitting"
      min-width="480px"
      max-width="600px"
      @submit="onSubmit"
    >
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">Tgl. Transaksi <span class="req">*</span></label>
          <q-input
            v-model="form.tanggal"
            outlined
            dense
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            class="clay-input"
            :rules="[required('Tanggal wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Tipe <span class="req">*</span></label>
          <q-select
            v-model="form.tipe"
            :options="tipeOptions"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            label="Pilih Tipe"
            :rules="[required('Tipe wajib dipilih')]"
            lazy-rules
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">Akun Bank / Kas <span class="req">*</span></label>
          <q-select
            v-model="form.akunBank"
            :options="financeStore.rekeningOptions"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            label="Pilih Akun Bank / Kas"
            :rules="[required('Akun bank wajib dipilih')]"
            lazy-rules
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Kategori <span class="req">*</span></label>
          <q-select
            v-model="form.kategori"
            :options="financeStore.kategoriOptions"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            label="Pilih Kategori"
            :rules="[required('Kategori wajib dipilih')]"
            lazy-rules
          />
        </div>
      </div>
      <div class="form-field-group">
        <label class="form-label">Nominal (Rp) <span class="req">*</span></label>
        <q-input
          v-model.number="form.nominal"
          outlined
          dense
          type="number"
          min="1"
          placeholder="0"
          class="clay-input"
          :rules="[requiredNumber('Nominal wajib diisi'), positiveNumber()]"
          lazy-rules
        >
          <template #prepend
            ><span style="font-size: 13px; color: var(--clay-text-muted); padding-left: 4px"
              >Rp</span
            ></template
          >
        </q-input>
      </div>
      <div class="form-field-group">
        <label class="form-label">Keterangan</label>
        <q-input
          v-model="form.keterangan"
          outlined
          dense
          type="textarea"
          rows="2"
          placeholder="Keterangan transaksi..."
          class="clay-input"
        />
      </div>
      <div class="form-field-group">
        <label class="form-label">Bukti Transaksi</label>
        <AppFileUpload
          :file-name="form.buktiFileName"
          icon-color="purple-5"
          text-color="var(--clay-accent)"
          @file-selected="onFileSelected"
          @file-removed="removeFile"
        />
      </div>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useFinanceStore } from 'src/stores/finance.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, positiveNumber } from 'src/utils/validators'
import { formatCurrency } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import AppFileUpload from 'components/shared/AppFileUpload.vue'

const financeStore = useFinanceStore()
const { pagination, syncTotal, onRequest, resetPage } = useTablePagination()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const search = ref('')
const filterTipe = ref(null)
const filterBank = ref(null)
const dialogOpen = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editId = ref(null)

const emptyForm = () => ({
  tanggal: new Date().toISOString().slice(0, 10),
  akunBank: '',
  tipe: '',
  kategori: '',
  nominal: null,
  keterangan: '',
  buktiFile: null,
  buktiFileName: '',
})
const form = reactive(emptyForm())

const tipeOptions = [
  { label: '↓ Masuk (Pemasukan)', value: 'Masuk' },
  { label: '↑ Keluar (Pengeluaran)', value: 'Keluar' },
]
const tipeFilterOptions = [
  { label: 'Masuk', value: 'Masuk' },
  { label: 'Keluar', value: 'Keluar' },
]

const columns = [
  { name: 'tanggal', label: 'Tanggal', field: 'tanggal', align: 'left', sortable: true },
  { name: 'namaBank', label: 'Akun Bank', field: 'namaBank', align: 'left' },
  { name: 'tipe', label: 'Tipe', field: 'tipe', align: 'left' },
  { name: 'namaKategori', label: 'Kategori', field: 'namaKategori', align: 'left' },
  { name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right', sortable: true },
  { name: 'keterangan', label: 'Keterangan', field: 'keterangan', align: 'left' },
  { name: 'buktiFile', label: 'Bukti', field: 'buktiFile', align: 'center' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

// Computed summaries from current page rows
const totalMasuk = computed(() =>
  financeStore.jurnalList
    .filter((r) => r.tipe === 'Masuk')
    .reduce((s, r) => s + (r.nominal || 0), 0),
)
const totalKeluar = computed(() =>
  financeStore.jurnalList
    .filter((r) => r.tipe === 'Keluar')
    .reduce((s, r) => s + (r.nominal || 0), 0),
)
const saldoBersih = computed(() => totalMasuk.value - totalKeluar.value)

async function load(page = pagination.page) {
  await financeStore.fetchJurnal({
    page,
    perPage: pagination.rowsPerPage,
    search: search.value,
    tipe: filterTipe.value || '',
    akunBank: filterBank.value || '',
  })
  syncTotal(financeStore.jurnalMeta.total)
}

function onFilter() {
  resetPage()
  load(1)
}

async function handleRequest(props) {
  await onRequest(props, ({ page, perPage }) =>
    financeStore.fetchJurnal({
      page,
      perPage,
      search: search.value,
      tipe: filterTipe.value || '',
      akunBank: filterBank.value || '',
    }),
  )
  syncTotal(financeStore.jurnalMeta.total)
}

onMounted(() => {
  financeStore.fetchRekening({ perPage: 100 })
  financeStore.fetchKategori({ perPage: 100 })
  load()
})

function openDialog(row = null) {
  isEdit.value = !!row
  editId.value = row?.id || null
  Object.assign(
    form,
    row
      ? {
          tanggal: row.tanggal,
          akunBank: row.akunBank,
          tipe: row.tipe,
          kategori: row.kategori,
          nominal: row.nominal,
          keterangan: row.keterangan || '',
          buktiFile: null,
          buktiFileName: row.buktiFile || '',
        }
      : emptyForm(),
  )
  dialogOpen.value = true
}

function onFileSelected(file) {
  form.buktiFile = file
  form.buktiFileName = file.name
}

function removeFile() {
  form.buktiFile = null
  form.buktiFileName = ''
}

async function onSubmit() {
  submitting.value = true
  try {
    const payload = { ...form }
    if (isEdit.value) {
      await financeStore.updateJurnal(editId.value, payload)
      notifySuccess('Jurnal berhasil diperbarui')
    } else {
      await financeStore.createJurnal(payload)
      notifySuccess('Transaksi berhasil dicatat')
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
    message: 'Hapus jurnal transaksi ini?',
    onConfirm: async () => {
      await financeStore.deleteJurnal(row.id)
      notifySuccess('Jurnal berhasil dihapus')
      load()
    },
  })
}
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.jurnal-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}
</style>

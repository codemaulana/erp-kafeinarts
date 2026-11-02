<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Tagihan Vendor (AP)"
      subtitle="Kelola semua kewajiban pembayaran kepada vendor"
      icon="receipt"
      icon-bg="linear-gradient(135deg,#FEECF0,#F9BDC8)"
      icon-color="#E85D75"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          label="Tambah Tagihan"
          icon="add"
          class="clay-btn clay-btn-primary"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <!-- Summary -->
    <div class="ap-summary q-mb-lg">
      <AppStatCard
        v-for="s in summaryCards"
        :key="s.label"
        :icon="s.icon"
        :icon-bg="s.iconBg"
        :icon-color="s.iconColor"
        :label="s.label"
        :value="s.value"
        :value-color="s.iconColor"
      />
    </div>

    <AppDataTable
      :rows="financeStore.tagihanList"
      :columns="columns"
      :loading="financeStore.isLoading('tagihan')"
      :pagination="pagination"
      no-data-label="Belum ada tagihan vendor"
      @request="handleRequest"
    >
      <template #filter>
        <div class="filter-row">
          <q-input
            v-model="search"
            placeholder="Cari vendor, no. tagihan..."
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
            v-model="filterStatus"
            :options="statusOptions"
            outlined
            dense
            clearable
            emit-value
            map-options
            class="clay-input filter-select"
            label="Pilih Status"
            @update:model-value="onFilter"
          />
        </div>
      </template>

      <template #body-cell-status="props">
        <q-td :props="props">
          <AppBadge :variant="statusVariant(props.value)">{{ props.value }}</AppBadge>
        </q-td>
      </template>

      <template #body-cell-nominal="props">
        <q-td :props="props">
          <span style="font-weight: 600; color: var(--clay-danger)">{{
            formatCurrency(props.value)
          }}</span>
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

    <!-- Dialog -->
    <AppFormDialog
      v-model="dialogOpen"
      :title="isEdit ? 'Edit Tagihan Vendor' : 'Tambah Tagihan Vendor'"
      :icon="isEdit ? 'edit' : 'add_circle'"
      :icon-color="isEdit ? 'blue-6' : 'red-6'"
      :submit-label="isEdit ? 'Simpan Perubahan' : 'Tambah Tagihan'"
      :loading="submitting"
      min-width="480px"
      max-width="600px"
      @submit="onSubmit"
    >
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">No. Tagihan <span class="req">*</span></label>
          <q-input
            v-model="form.noTagihan"
            outlined
            dense
            class="clay-input"
            :rules="[required('No. tagihan wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Status <span class="req">*</span></label>
          <q-select
            v-model="form.status"
            :options="statusOptions"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            label="Pilih Status"
            :rules="[required('Status wajib dipilih')]"
            lazy-rules
          />
        </div>
      </div>
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">Vendor <span class="req">*</span></label>
          <q-input
            v-model="form.vendor"
            outlined
            dense
            placeholder="Nama vendor / supplier"
            class="clay-input"
            :rules="[required('Vendor wajib diisi')]"
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
      <div class="form-row">
        <div class="form-field-group">
          <label class="form-label">Tgl. Tagihan <span class="req">*</span></label>
          <q-input
            v-model="form.tanggalTagihan"
            outlined
            dense
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            class="clay-input"
            :rules="[required('Tanggal tagihan wajib diisi')]"
            lazy-rules
          />
        </div>
        <div class="form-field-group">
          <label class="form-label">Tgl. Jatuh Tempo <span class="req">*</span></label>
          <q-input
            v-model="form.jatuhTempo"
            outlined
            dense
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            class="clay-input"
            :rules="[required('Jatuh tempo wajib diisi')]"
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
        <label class="form-label">Bukti / Dokumen Tagihan</label>
        <AppFileUpload
          :file-name="form.buktiFileName"
          icon-color="red-4"
          text-color="var(--clay-danger)"
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
const filterStatus = ref(null)
const dialogOpen = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const editId = ref(null)

const today = new Date().toISOString().slice(0, 10)
const emptyForm = () => ({
  noTagihan: `AP/${new Date().getFullYear()}/${String(new Date().getMonth() + 1).padStart(2, '0')}/XXX`,
  vendor: '',
  kategori: '',
  tanggalTagihan: today,
  jatuhTempo: '',
  nominal: null,
  status: 'Belum Lunas',
  buktiFile: null,
  buktiFileName: '',
})
const form = reactive(emptyForm())

const statusOptions = [
  { label: 'Lunas', value: 'Lunas' },
  { label: 'Belum Lunas', value: 'Belum Lunas' },
  { label: 'Jatuh Tempo', value: 'Jatuh Tempo' },
]

const columns = [
  { name: 'noTagihan', label: 'No. Tagihan', field: 'noTagihan', align: 'left', sortable: true },
  { name: 'vendor', label: 'Vendor', field: 'vendor', align: 'left' },
  { name: 'namaKategori', label: 'Kategori', field: 'namaKategori', align: 'left' },
  { name: 'tanggalTagihan', label: 'Tgl. Tagihan', field: 'tanggalTagihan', align: 'left' },
  { name: 'jatuhTempo', label: 'Jatuh Tempo', field: 'jatuhTempo', align: 'left' },
  { name: 'nominal', label: 'Nominal', field: 'nominal', align: 'right', sortable: true },
  { name: 'status', label: 'Status', field: 'status', align: 'left' },
  { name: 'buktiFile', label: 'Bukti', field: 'buktiFile', align: 'center' },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

const summaryCards = computed(() => {
  const all = financeStore.tagihanList
  const lunas = all.filter((r) => r.status === 'Lunas').length
  const jatuhTempo = all.filter((r) => r.status === 'Jatuh Tempo').length
  const totalNominal = all.reduce((s, r) => s + (r.nominal || 0), 0)
  return [
    {
      label: 'Total Tagihan',
      value: all.length,
      icon: 'receipt',
      iconBg: 'linear-gradient(135deg,#EBF4FB,#C7E1F5)',
      iconColor: '#5B9BD5',
    },
    {
      label: 'Sudah Lunas',
      value: lunas,
      icon: 'check_circle',
      iconBg: 'linear-gradient(135deg,#E8F7F2,#BDE8D9)',
      iconColor: '#4CAF8A',
    },
    {
      label: 'Jatuh Tempo',
      value: jatuhTempo,
      icon: 'warning',
      iconBg: 'linear-gradient(135deg,#FEECF0,#F9BDC8)',
      iconColor: '#E85D75',
    },
    {
      label: 'Total Kewajiban',
      value: formatCurrency(totalNominal),
      icon: 'payments',
      iconBg: 'linear-gradient(135deg,#FEF3E2,#F9D9A0)',
      iconColor: '#E8A547',
    },
  ]
})

function statusVariant(status) {
  return (
    { Lunas: 'success', 'Belum Lunas': 'warning', 'Jatuh Tempo': 'danger' }[status] || 'neutral'
  )
}

async function load(page = pagination.page) {
  await financeStore.fetchTagihan({
    page,
    perPage: pagination.rowsPerPage,
    search: search.value,
    status: filterStatus.value || '',
  })
  syncTotal(financeStore.tagihanMeta.total)
}

function onFilter() {
  resetPage()
  load(1)
}

async function handleRequest(props) {
  await onRequest(props, ({ page, perPage }) =>
    financeStore.fetchTagihan({
      page,
      perPage,
      search: search.value,
      status: filterStatus.value || '',
    }),
  )
  syncTotal(financeStore.tagihanMeta.total)
}

onMounted(() => {
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
          noTagihan: row.noTagihan,
          vendor: row.vendor,
          kategori: row.kategori,
          tanggalTagihan: row.tanggalTagihan,
          jatuhTempo: row.jatuhTempo,
          nominal: row.nominal,
          status: row.status,
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
      await financeStore.updateTagihan(editId.value, payload)
      notifySuccess('Tagihan berhasil diperbarui')
    } else {
      await financeStore.createTagihan(payload)
      notifySuccess('Tagihan berhasil ditambahkan')
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
    itemName: row.noTagihan,
    onConfirm: async () => {
      await financeStore.deleteTagihan(row.id)
      notifySuccess('Tagihan berhasil dihapus')
      load()
    },
  })
}
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.ap-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}
</style>

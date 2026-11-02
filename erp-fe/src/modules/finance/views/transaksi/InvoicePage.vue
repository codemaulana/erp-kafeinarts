<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Manajemen Invoice (AR)"
      subtitle="Kelola tagihan dan piutang kepada klien"
      icon="description"
      icon-bg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      icon-color="#4CAF8A"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          label="Buat Invoice"
          icon="add"
          class="clay-btn clay-btn-primary"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <!-- Summary -->
    <div class="inv-summary q-mb-lg">
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
      :rows="financeStore.invoiceList"
      :columns="columns"
      :loading="financeStore.isLoading('invoice')"
      :pagination="pagination"
      no-data-label="Belum ada invoice"
      @request="handleRequest"
    >
      <template #filter>
        <div class="filter-row">
          <q-input
            v-model="search"
            placeholder="Cari no. invoice, klien, proyek..."
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
          <AppBadge :variant="statusVariant(props.value)">{{ formatStatus(props.value) }}</AppBadge>
        </q-td>
      </template>

      <template #body-cell-grandTotal="props">
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
            icon="visibility"
            size="sm"
            color="purple-6"
            class="q-mr-xs"
            @click="openDialog(props.row, true)"
            ><q-tooltip>Detail</q-tooltip></q-btn
          >
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

    <!-- Dialog Form (maximized) -->
    <AppFormDialog
      v-model="dialogOpen"
      :title="isView ? 'Detail Invoice' : isEdit ? 'Edit Invoice' : 'Buat Invoice Baru'"
      :icon="isView ? 'visibility' : isEdit ? 'edit' : 'add_circle'"
      :icon-color="isView ? 'purple-6' : isEdit ? 'blue-6' : 'green-6'"
      :submit-label="isEdit ? 'Simpan Perubahan' : 'Buat Invoice'"
      :loading="submitting"
      :maximized="true"
      @submit="onSubmit"
    >
      <div class="inv-form-grid">
        <!-- Kolom kiri -->
        <div>
          <div class="form-section-title">Informasi Invoice</div>
          <div class="form-row">
            <div class="form-field-group">
              <label class="form-label">No. Invoice <span class="req">*</span></label>
              <q-input
                v-model="form.noInvoice"
                outlined
                dense
                :readonly="isView"
                class="clay-input"
                :rules="[required('No. invoice wajib diisi')]"
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
                :readonly="isView"
                class="clay-input"
                label="Pilih Status"
                :rules="[required('Status wajib dipilih')]"
                lazy-rules
              />
            </div>
          </div>
          <div class="form-row">
            <div class="form-field-group">
              <label class="form-label">Klien <span class="req">*</span></label>
              <q-input
                v-model="form.klien"
                outlined
                dense
                :readonly="isView"
                class="clay-input"
                :rules="[required('Klien wajib diisi')]"
                lazy-rules
              />
            </div>
            <div class="form-field-group">
              <label class="form-label">Proyek <span class="req">*</span></label>
              <q-input
                v-model="form.proyek"
                outlined
                dense
                :readonly="isView"
                class="clay-input"
                :rules="[required('Proyek wajib diisi')]"
                lazy-rules
              />
            </div>
          </div>
          <div class="form-row">
            <div class="form-field-group">
              <label class="form-label">Tgl. Terbit <span class="req">*</span></label>
              <q-input
                v-model="form.tglTerbit"
                outlined
                dense
                type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
                :readonly="isView"
                class="clay-input"
                :rules="[required('Tanggal terbit wajib diisi')]"
                lazy-rules
              />
            </div>
            <div class="form-field-group">
              <label class="form-label">Tgl. Jatuh Tempo <span class="req">*</span></label>
              <q-input
                v-model="form.tglJatuhTempo"
                outlined
                dense
                type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
                :readonly="isView"
                class="clay-input"
                :rules="[required('Tanggal jatuh tempo wajib diisi')]"
                lazy-rules
              />
            </div>
          </div>
          <div class="form-field-group">
            <label class="form-label">Pajak</label>
            <q-select
              v-model="form.pajak"
              :options="pajakOptionsAll"
              outlined
              dense
              clearable
              emit-value
              map-options
              :readonly="isView"
              class="clay-input"
              label="Tidak Ada Pajak"
              @update:model-value="recalcTotal"
            />
          </div>
        </div>

        <!-- Kolom kanan — Items -->
        <div>
          <div class="form-section-title">
            Item Penagihan
            <q-btn
              v-if="!isView"
              flat
              dense
              no-caps
              icon="add"
              label="Tambah Item"
              size="sm"
              color="green-6"
              @click="addItem"
            />
          </div>

          <div class="inv-items">
            <div class="inv-items__header">
              <span>Deskripsi</span><span class="text-center">Qty</span
              ><span class="text-right">Harga</span><span class="text-right">Subtotal</span>
              <span v-if="!isView"></span>
            </div>
            <div v-for="(item, idx) in form.items" :key="item.id" class="inv-items__row">
              <q-input
                v-model="item.deskripsi"
                outlined
                dense
                :readonly="isView"
                class="clay-input"
                placeholder="Deskripsi item..."
              />
              <q-input
                v-model.number="item.qty"
                outlined
                dense
                type="number"
                min="1"
                :readonly="isView"
                class="clay-input text-center"
                @update:model-value="calcItem(idx)"
              />
              <q-input
                v-model.number="item.harga"
                outlined
                dense
                type="number"
                min="0"
                :readonly="isView"
                class="clay-input"
                @update:model-value="calcItem(idx)"
              >
                <template #prepend
                  ><span style="font-size: 11px; color: var(--clay-text-muted)">Rp</span></template
                >
              </q-input>
              <div class="inv-items__subtotal">{{ formatCurrency(item.subtotal || 0) }}</div>
              <q-btn
                v-if="!isView"
                flat
                dense
                round
                icon="remove_circle_outline"
                size="xs"
                color="red-4"
                @click="removeItem(idx)"
              />
            </div>
          </div>

          <div class="inv-totals">
            <div class="inv-totals__row">
              <span>Subtotal</span>
              <span>{{ formatCurrency(form.totalTagihan) }}</span>
            </div>
            <div v-if="form.catatanPajak" class="inv-totals__row">
              <span>{{ form.catatanPajak }}</span>
              <span>{{ formatCurrency(pajakAmount) }}</span>
            </div>
            <div class="inv-totals__row inv-totals__row--grand">
              <span>Grand Total</span>
              <span>{{ formatCurrency(form.grandTotal) }}</span>
            </div>
          </div>
        </div>
      </div>

      <!-- Override footer for view-only mode -->
      <template v-if="isView" #footer>
        <q-btn flat no-caps label="Tutup" class="clay-btn" @click="dialogOpen = false" />
      </template>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { uid } from 'quasar'
import { useFinanceStore } from 'src/stores/finance.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required } from 'src/utils/validators'
import { formatCurrency } from 'src/utils/formatter'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppStatCard from 'components/shared/AppStatCard.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const financeStore = useFinanceStore()
const { pagination, syncTotal, onRequest, resetPage } = useTablePagination()
const { confirmDelete } = useConfirmDelete()
const { notifySuccess, notifyError } = useNotify()

const search = ref('')
const filterStatus = ref(null)
const dialogOpen = ref(false)
const isEdit = ref(false)
const isView = ref(false)
const submitting = ref(false)
const editId = ref(null)

const statusOptions = [
  { label: 'Lunas', value: 'Lunas' },
  { label: 'Sebagian', value: 'Sebagian' },
  { label: 'Belum Lunas', value: 'Belum Lunas' },
]

const columns = [
  { name: 'noInvoice', label: 'No. Invoice', field: 'noInvoice', align: 'left', sortable: true },
  { name: 'klien', label: 'Klien', field: 'klien', align: 'left' },
  { name: 'proyek', label: 'Proyek', field: 'proyek', align: 'left' },
  { name: 'tglTerbit', label: 'Tgl. Terbit', field: 'tglTerbit', align: 'left' },
  { name: 'tglJatuhTempo', label: 'Jatuh Tempo', field: 'tglJatuhTempo', align: 'left' },
  { name: 'status', label: 'Status', field: 'status', align: 'left' },
  { name: 'grandTotal', label: 'Grand Total', field: 'grandTotal', align: 'right', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

const emptyForm = () => {
  const today = new Date().toISOString().slice(0, 10)
  const due = new Date(Date.now() + 30 * 86400000).toISOString().slice(0, 10)
  return {
    noInvoice: `INV/${new Date().getFullYear()}/${String(new Date().getMonth() + 1).padStart(2, '0')}/XXX`,
    klien: '',
    proyek: '',
    tglTerbit: today,
    tglJatuhTempo: due,
    status: 'Belum Lunas',
    pajak: '',
    catatanPajak: '',
    items: [{ id: uid(), deskripsi: '', qty: 1, harga: 0, subtotal: 0 }],
    totalTagihan: 0,
    grandTotal: 0,
  }
}
const form = reactive(emptyForm())

const pajakOptionsAll = computed(() => [
  { label: 'Tidak Ada Pajak', value: '' },
  ...financeStore.pajakOptions,
])

const pajakAmount = computed(() => {
  if (!form.pajak) return 0
  const p = financeStore.pajakOptions.find((o) => o.value === form.pajak)
  return p ? Math.round((form.totalTagihan * p.persentase) / 100) : 0
})

const isLunas = (st) => {
  const s = String(st || '').toUpperCase()
  return s === 'LUNAS' || s === 'PAID'
}

const isBelumLunas = (st) => {
  const s = String(st || '').toUpperCase()
  return s !== 'LUNAS' && s !== 'PAID' && s !== 'BATAL' && s !== 'CANCELLED'
}

function formatStatus(st) {
  const s = String(st || '').toUpperCase()
  if (s === 'PAID' || s === 'LUNAS') return 'Lunas'
  if (s === 'PARTIAL' || s === 'SEBAGIAN') return 'Sebagian'
  if (s === 'SENT' || s === 'UNPAID' || s === 'OVERDUE' || s === 'BELUM LUNAS') return 'Belum Lunas'
  if (s === 'DRAFT') return 'Draft'
  if (s === 'CANCELLED' || s === 'BATAL') return 'Batal'
  return st || '-'
}

const summaryCards = computed(() => {
  const all = financeStore.invoiceList
  const lunas = all.filter((i) => isLunas(i.status)).length
  const belum = all.filter((i) => isBelumLunas(i.status)).length
  const total = all.reduce((s, i) => s + (i.grandTotal || 0), 0)
  return [
    {
      label: 'Total Invoice',
      value: all.length,
      icon: 'description',
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
      label: 'Belum Lunas',
      value: belum,
      icon: 'pending',
      iconBg: 'linear-gradient(135deg,#FEECF0,#F9BDC8)',
      iconColor: '#E85D75',
    },
    {
      label: 'Total Tagihan',
      value: formatCurrency(total),
      icon: 'paid',
      iconBg: 'linear-gradient(135deg,#FEF3E2,#F9D9A0)',
      iconColor: '#E8A547',
    },
  ]
})

function statusVariant(status) {
  const s = String(status || '').toUpperCase()
  if (s === 'LUNAS' || s === 'PAID') return 'success'
  if (s === 'SEBAGIAN' || s === 'PARTIAL') return 'warning'
  if (s === 'DRAFT') return 'neutral'
  return 'danger'
}

async function load(page = pagination.page) {
  await financeStore.fetchInvoice({
    page,
    perPage: pagination.rowsPerPage,
    search: search.value,
    status: filterStatus.value || '',
  })
  syncTotal(financeStore.invoiceMeta.total)
}

function onFilter() {
  resetPage()
  load(1)
}

async function handleRequest(props) {
  await onRequest(props, ({ page, perPage }) =>
    financeStore.fetchInvoice({
      page,
      perPage,
      search: search.value,
      status: filterStatus.value || '',
    }),
  )
  syncTotal(financeStore.invoiceMeta.total)
}

onMounted(() => {
  financeStore.fetchPajak({ perPage: 100 })
  load()
})

function openDialog(row = null, view = false) {
  isView.value = view
  isEdit.value = !!row && !view
  editId.value = row?.id || null
  if (row) {
    Object.assign(form, {
      noInvoice: row.noInvoice,
      klien: row.klien,
      proyek: row.proyek,
      tglTerbit: row.tglTerbit,
      tglJatuhTempo: row.tglJatuhTempo,
      status: row.status,
      pajak: row.pajak || '',
      catatanPajak: row.catatanPajak || '',
      items: row.items.map((i) => ({ ...i })),
      totalTagihan: row.totalTagihan,
      grandTotal: row.grandTotal,
    })
  } else {
    Object.assign(form, emptyForm())
  }
  dialogOpen.value = true
}

function addItem() {
  form.items.push({ id: uid(), deskripsi: '', qty: 1, harga: 0, subtotal: 0 })
}

function removeItem(idx) {
  if (form.items.length > 1) {
    form.items.splice(idx, 1)
    recalcTotal()
  }
}

function calcItem(idx) {
  const item = form.items[idx]
  item.subtotal = (item.qty || 0) * (item.harga || 0)
  recalcTotal()
}

function recalcTotal() {
  form.totalTagihan = form.items.reduce((s, i) => s + (i.subtotal || 0), 0)
  const p = financeStore.pajakOptions.find((o) => o.value === form.pajak)
  const pajakNominal = p ? Math.round((form.totalTagihan * p.persentase) / 100) : 0
  form.grandTotal = form.totalTagihan + pajakNominal
  form.catatanPajak = p ? p.label : ''
}

async function onSubmit() {
  submitting.value = true
  try {
    if (isEdit.value) {
      await financeStore.updateInvoice(editId.value, { ...form })
      notifySuccess('Invoice berhasil diperbarui')
    } else {
      await financeStore.createInvoice({ ...form })
      notifySuccess('Invoice berhasil dibuat')
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
    itemName: row.noInvoice,
    onConfirm: async () => {
      await financeStore.deleteInvoice(row.id)
      notifySuccess('Invoice berhasil dihapus')
      load()
    },
  })
}
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';

.inv-summary {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(140px, 1fr));
  gap: 12px;
}

.inv-form-grid {
  display: grid;
  grid-template-columns: 1fr 1fr;
  gap: 32px;
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.inv-items {
  background: var(--rb-bg);
  border-radius: var(--rb-radius-md);
  border: 1px solid var(--rb-border);
  overflow: hidden;
  margin-bottom: 16px;
  &__header {
    display: grid;
    grid-template-columns: 1fr 80px 120px 120px 36px;
    gap: 8px;
    padding: 10px 12px;
    background: var(--rb-surface-alt);
    font-size: 11px;
    font-weight: 700;
    color: var(--rb-text-muted);
    text-transform: uppercase;
    letter-spacing: 0.6px;
  }
  &__row {
    display: grid;
    grid-template-columns: 1fr 80px 120px 120px 36px;
    gap: 8px;
    padding: 8px 12px;
    align-items: center;
    border-bottom: 1px solid var(--rb-border-light);
    &:last-child {
      border-bottom: none;
    }
  }
  &__subtotal {
    text-align: right;
    font-size: 13px;
    font-weight: 600;
    color: var(--rb-success);
  }
}

.inv-totals {
  background: var(--rb-surface-alt);
  border-radius: var(--rb-radius-md);
  border: 1px solid var(--rb-border);
  padding: 14px 16px;
  &__row {
    display: flex;
    justify-content: space-between;
    padding: 5px 0;
    font-size: 13px;
    color: var(--rb-text-muted);
    border-bottom: 1px solid var(--rb-border-light);
    &:last-child {
      border-bottom: none;
    }
    &--grand {
      font-size: 15px;
      font-weight: 800;
      color: var(--rb-text-primary);
      padding-top: 10px;
      margin-top: 4px;
      border-top: 2px solid var(--rb-border);
    }
  }
}
</style>

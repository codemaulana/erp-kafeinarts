<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Purchase Order (PO)"
      subtitle="Dokumen pesanan resmi ke vendor"
      icon="shopping_cart"
      iconBg="linear-gradient(135deg,#EBF4FB,#C7E1F5)"
      iconColor="#5B9BD5"
    >
      <template #actions>
        <q-btn unelevated no-caps class="clay-btn clay-btn-primary" @click="openDialog()">
          <q-icon name="add" size="xs" class="q-mr-sm" /> Tambah PO
        </q-btn>
      </template>
    </AppPageHeader>

    <q-card class="clay-card q-mt-md">
      <q-card-section>
        <div class="row q-mb-md">
          <div class="col-12 col-md-4">
            <q-input
              v-model="searchQuery"
              outlined
              dense
              placeholder="Cari PO..."
              class="clay-input"
              debounce="350"
              @update:model-value="onSearch"
            >
              <template #append>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
        </div>

        <AppDataTable
          :data="scStore.purchaseOrderList"
          :columns="columns"
          :loading="scStore.isLoading('po')"
          :pagination="pagination"
          @request="onRequest"
        >
          <template #cell-po_number="{ row }">
            <div class="text-weight-medium">{{ row.po_number }}</div>
          </template>

          <template #cell-po_date="{ row }">
            {{ formatDate(row.po_date) }}
          </template>

          <template #cell-total_amount="{ row }">
            {{ formatCurrency(row.total_amount) }}
          </template>

          <template #cell-status="{ row }">
            <AppBadge :variant="getStatusBadgeVariant(row.status)">
              {{ row.status }}
            </AppBadge>
          </template>

          <template #cell-actions="{ row }">
            <div class="row justify-center q-gutter-x-xs">
              <q-btn
                flat
                dense
                round
                size="sm"
                icon="edit"
                color="primary"
                @click="openDialog(row)"
              >
                <q-tooltip>Edit</q-tooltip>
              </q-btn>

              <q-btn
                flat
                dense
                round
                size="sm"
                icon="delete"
                color="negative"
                @click="handleDelete(row)"
              >
                <q-tooltip>Hapus</q-tooltip>
              </q-btn>
            </div>
          </template>
        </AppDataTable>
      </q-card-section>
    </q-card>

    <AppFormDialog
      v-model="isDialogOpen"
      :title="isEditMode ? 'Edit PO' : 'Tambah PO Baru'"
      :loading="scStore.isLoading('po')"
      width="800px"
      @save="saveData"
      @cancel="isDialogOpen = false"
    >
      <q-form ref="formRef" class="q-gutter-md">
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <label class="form-label">No. PO</label>
            <q-input
              v-model="form.po_number"
              outlined
              dense
              class="clay-input"
              readonly
              placeholder="Auto-generate"
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Referensi PR</label>
            <q-select
              v-model="form.pr_id"
              :options="scStore.prOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              clearable
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Vendor <span class="req">*</span></label>
            <ScVendorSelect v-model="form.vendor_id" :rules="[required('Vendor wajib diisi')]" />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Tanggal PO <span class="req">*</span></label>
            <q-input
              v-model="form.po_date"
              type="text"
              @focus="(e) => (e.target.type = 'date')"
              @blur="
                (e) => {
                  if (!e.target.value) e.target.type = 'text'
                }
              "
              outlined
              dense
              class="clay-input"
              :rules="[required('Tanggal wajib diisi')]"
              lazy-rules
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Estimasi Pengiriman</label>
            <q-input
              v-model="form.estimated_delivery_date"
              type="text"
              @focus="(e) => (e.target.type = 'date')"
              @blur="
                (e) => {
                  if (!e.target.value) e.target.type = 'text'
                }
              "
              outlined
              dense
              class="clay-input"
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Termin Pembayaran</label>
            <q-input
              v-model="form.payment_terms"
              outlined
              dense
              placeholder="Misal: NET 30"
              class="clay-input"
            />
          </div>
        </div>

        <div class="q-mt-lg">
          <div class="row justify-between items-center q-mb-sm">
            <div class="text-subtitle1 text-weight-bold">Item PO</div>
            <q-btn flat dense no-caps color="primary" @click="addItem()">
              <q-icon name="add" size="xs" class="q-mr-xs" /> Tambah Item
            </q-btn>
          </div>

          <q-markup-table flat bordered class="clay-table">
            <thead>
              <tr>
                <th class="text-left">Item <span class="req">*</span></th>
                <th class="text-right">Kuantitas <span class="req">*</span></th>
                <th class="text-right">Harga Satuan <span class="req">*</span></th>
                <th class="text-right">Subtotal</th>
                <th class="text-center">Aksi</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in items" :key="index">
                <td>
                  <ScItemSelect v-model="item.item_id" :rules="[required('Item wajib diisi')]" />
                </td>
                <td>
                  <q-input
                    v-model.number="item.quantity"
                    type="number"
                    min="1"
                    outlined
                    dense
                    class="clay-input"
                    :rules="[requiredNumber('Kuantitas wajib diisi'), positiveNumber('Harus > 0')]"
                    @update:model-value="(val) => updateItem(index, 'quantity', val)"
                  />
                </td>
                <td>
                  <q-input
                    v-model.number="item.unit_price"
                    type="number"
                    min="0"
                    outlined
                    dense
                    class="clay-input"
                    :rules="[
                      requiredNumber('Harga wajib diisi'),
                      nonNegative('Tidak boleh negatif'),
                    ]"
                    @update:model-value="(val) => updateItem(index, 'unit_price', val)"
                  />
                </td>
                <td class="text-right text-weight-medium">
                  {{ formatCurrency((item.quantity || 0) * (item.unit_price || 0)) }}
                </td>
                <td class="text-center">
                  <q-btn
                    flat
                    dense
                    round
                    size="sm"
                    icon="delete"
                    color="negative"
                    @click="removeItem(index)"
                  />
                </td>
              </tr>
              <tr v-if="items.length === 0">
                <td colspan="5" class="text-center text-grey">Belum ada item</td>
              </tr>
            </tbody>
            <tfoot>
              <tr>
                <td colspan="3" class="text-right text-weight-bold">Total PO</td>
                <td class="text-right text-weight-bold text-primary">
                  {{ formatCurrency(total) }}
                </td>
                <td></td>
              </tr>
            </tfoot>
          </q-markup-table>
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useScStore } from 'src/stores/sc.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, positiveNumber, nonNegative } from 'src/utils/validators'
import { formatDate, formatCurrency } from 'src/utils/formatter'
import useScDynamicGrid from 'src/composables/useScDynamicGrid'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'
import ScVendorSelect from 'src/modules/sc/components/ScVendorSelect.vue'
import ScItemSelect from 'src/modules/sc/components/ScItemSelect.vue'

const scStore = useScStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const { pagination, searchQuery, onRequest, resetPage } = useTablePagination(loadData)

const isDialogOpen = ref(false)
const isEditMode = ref(false)
const formRef = ref(null)

const columns = [
  { name: 'po_number', label: 'No. PO', field: 'po_number', align: 'left', sortable: true },
  { name: 'vendor_name', label: 'Vendor', field: 'vendor_name', align: 'left', sortable: true },
  { name: 'po_date', label: 'Tanggal', field: 'po_date', align: 'left', sortable: true },
  {
    name: 'total_amount',
    label: 'Total PO',
    field: 'total_amount',
    align: 'right',
    sortable: true,
  },
  { name: 'status', label: 'Status', field: 'status', align: 'center', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

const initialForm = {
  id: null,
  po_number: '',
  pr_id: null,
  vendor_id: null,
  po_date: '',
  estimated_delivery_date: '',
  payment_terms: '',
  items: [],
}

const form = ref({ ...initialForm })
const { items, addItem, removeItem, updateItem, resetItems, total } = useScDynamicGrid()

watch(isDialogOpen, (val) => {
  if (val) {
    resetItems(form.value.items || [])
  }
})

function onSearch() {
  resetPage()
  loadData(1)
}

async function loadData(page = 1) {
  try {
    await scStore.fetchPurchaseOrder({
      page,
      limit: pagination.value.rowsPerPage,
      search: searchQuery.value,
    })
    pagination.value.page = scStore.purchaseOrderMeta.page
    pagination.value.rowsNumber = scStore.purchaseOrderMeta.total
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan')
  }
}

function openDialog(row = null) {
  isEditMode.value = !!row
  if (row) {
    form.value = { ...row, items: [...(row.items || [])] }
  } else {
    form.value = { ...initialForm }
  }
  isDialogOpen.value = true
}

async function saveData() {
  const valid = await formRef.value.validate()
  if (!valid) return

  if (items.value.length === 0) {
    notifyError('Minimal satu item harus ditambahkan')
    return
  }

  const payload = {
    ...form.value,
    items: items.value,
  }

  try {
    if (isEditMode.value) {
      await scStore.updatePurchaseOrder(payload.id, payload)
      notifySuccess('PO berhasil diupdate')
    } else {
      await scStore.createPurchaseOrder(payload)
      notifySuccess('PO berhasil dibuat')
    }
    isDialogOpen.value = false
    loadData(pagination.value.page)
  } catch (error) {
    error.message || 'Gagal menyimpan PO'
  }
}

function handleDelete(row) {
  confirmDelete('PO', async () => {
    try {
      await scStore.deletePurchaseOrder(row.id)
      notifySuccess('PO berhasil dihapus')
      loadData(pagination.value.page)
    } catch (error) {
      error.message || 'Gagal menghapus PO'
    }
  })
}


function getStatusBadgeVariant(status) {
  const map = {
    DRAFT: 'neutral',
    SENT: 'info',
    PARTIALLY_RECEIVED: 'warning',
    COMPLETED: 'success',
    CANCELLED: 'danger',
  }
  return map[status] || 'neutral'
}

onMounted(() => {
  loadData()
  scStore.fetchPurchaseRequest({ limit: 100 }) // Load PR for options
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

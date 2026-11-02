<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Penerimaan Barang (GR)"
      subtitle="Pencatatan fisik barang masuk"
      icon="local_shipping"
    >
      <template #actions>
        <q-btn unelevated no-caps class="clay-btn clay-btn-primary" @click="openDialog()">
          <q-icon name="add" size="xs" class="q-mr-sm" /> Terima Barang
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
              placeholder="Cari GR..."
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
          :data="scStore.goodsReceiptList"
          :columns="columns"
          :loading="scStore.isLoading('gr')"
          :pagination="pagination"
          @request="onRequest"
        >
          <template #cell-gr_number="{ row }">
            <div class="text-weight-medium">{{ row.gr_number || row.receipt_number || '-' }}</div>
          </template>

          <template #cell-receipt_date="{ row }">
            {{ formatDate(row.receipt_date) }}
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
      :title="isEditMode ? 'Edit GR' : 'Terima Barang Baru'"
      :loading="scStore.isLoading('gr')"
      width="900px"
      @save="saveData"
      @cancel="isDialogOpen = false"
    >
      <q-form ref="formRef" class="q-gutter-md">
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <label class="form-label">No. GR</label>
            <q-input
              v-model="form.gr_number"
              outlined
              dense
              class="clay-input"
              readonly
              placeholder="Auto-generate"
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Referensi PO <span class="req">*</span></label>
            <q-select
              v-model="form.po_id"
              :options="scStore.poOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              :rules="[required('PO wajib diisi')]"
              lazy-rules
              @update:model-value="onPoSelected"
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Penerima <span class="req">*</span></label>
            <q-select
              v-model="form.received_by"
              :options="hrStore.karyawanOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              :rules="[required('Penerima wajib diisi')]"
              lazy-rules
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Tanggal Terima <span class="req">*</span></label>
            <q-input
              v-model="form.receipt_date"
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

          <div class="col-12">
            <label class="form-label">Catatan</label>
            <q-input
              v-model="form.notes"
              type="textarea"
              outlined
              dense
              class="clay-input"
              rows="3"
            />
          </div>

          <div class="col-12">
            <label class="form-label">Bukti Terima (Opsional)</label>
            <AppFileUpload v-model="form.attachmentFile" accept=".pdf,image/*" />
          </div>
        </div>

        <div class="q-mt-lg">
          <div class="row justify-between items-center q-mb-sm">
            <div class="text-subtitle1 text-weight-bold">Item Diterima</div>
          </div>

          <q-markup-table flat bordered class="clay-table">
            <thead>
              <tr>
                <th class="text-left">Item</th>
                <th class="text-right">Qty Dipesan</th>
                <th class="text-right">Qty Diterima <span class="req">*</span></th>
                <th class="text-left">Kondisi <span class="req">*</span></th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="(item, index) in items" :key="index">
                <td>
                  <ScItemSelect v-model="item.item_id" readonly disable />
                </td>
                <td class="text-right">
                  {{ item.ordered_quantity || 0 }}
                </td>
                <td>
                  <q-input
                    v-model.number="item.received_quantity"
                    type="number"
                    min="0"
                    :max="item.ordered_quantity"
                    outlined
                    dense
                    class="clay-input"
                    :rules="[requiredNumber('Wajib diisi'), nonNegative('Min 0')]"
                    @update:model-value="(val) => updateItem(index, 'received_quantity', val)"
                  />
                </td>
                <td>
                  <q-select
                    v-model="item.condition"
                    :options="['GOOD', 'DAMAGED', 'INCOMPLETE']"
                    outlined
                    dense
                    class="clay-input"
                    :rules="[required('Kondisi wajib diisi')]"
                  />
                </td>
              </tr>
              <tr v-if="items.length === 0">
                <td colspan="4" class="text-center text-grey">Pilih PO terlebih dahulu</td>
              </tr>
            </tbody>
          </q-markup-table>
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useScStore } from 'src/stores/sc.store'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, nonNegative } from 'src/utils/validators'
import { formatDate } from 'src/utils/formatter'
import useScDynamicGrid from 'src/composables/useScDynamicGrid'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppFileUpload from 'components/shared/AppFileUpload.vue'
import ScItemSelect from 'src/modules/sc/components/ScItemSelect.vue'

const scStore = useScStore()
const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const { pagination, searchQuery, onRequest, resetPage } = useTablePagination(loadData)

const isDialogOpen = ref(false)
const isEditMode = ref(false)
const formRef = ref(null)

const columns = [
  {
    name: 'gr_number',
    label: 'No. GR',
    field: (row) => row.gr_number || row.receipt_number || '-',
    align: 'left',
    sortable: true,
  },
  {
    name: 'po_number',
    label: 'No. PO',
    field: (row) => row.po_number || row.po_id || '-',
    align: 'left',
    sortable: true,
  },
  {
    name: 'receipt_date',
    label: 'Tanggal Terima',
    field: 'receipt_date',
    align: 'left',
    sortable: true,
  },
  {
    name: 'received_by_name',
    label: 'Diterima Oleh',
    field: (row) => row.received_by_name || row.received_by || '-',
    align: 'left',
    sortable: true,
  },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

const initialForm = {
  id: null,
  gr_number: '',
  po_id: null,
  receipt_date: '',
  received_by: null,
  notes: '',
  attachmentFile: null,
  items: [],
}

const form = ref({ ...initialForm })
const { items, updateItem, resetItems } = useScDynamicGrid()

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
    await scStore.fetchGoodsReceipt({
      page,
      limit: pagination.value.rowsPerPage,
      search: searchQuery.value,
    })
    pagination.value.page = scStore.goodsReceiptMeta.page
    pagination.value.rowsNumber = scStore.goodsReceiptMeta.total
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan')
  }
}

function openDialog(row = null) {
  isEditMode.value = !!row
  if (row) {
    form.value = {
      ...row,
      gr_number: row.gr_number || row.receipt_number || '',
      items: [...(row.items || [])],
    }
  } else {
    form.value = { ...initialForm }
  }
  isDialogOpen.value = true
}

// When PO is selected, auto-populate the items
async function onPoSelected(poId) {
  if (!poId) {
    resetItems([])
    return
  }

  // Find the selected PO from the store's list
  const po = scStore.purchaseOrderList.find((p) => p.id === poId)
  if (po && po.items) {
    const grItems = po.items.map((pi) => ({
      item_id: pi.item_id,
      po_item_id: pi.id,
      ordered_quantity: pi.quantity,
      received_quantity: pi.quantity, // default full
      condition: 'GOOD',
    }))
    resetItems(grItems)
  }
}

async function saveData() {
  const valid = await formRef.value.validate()
  if (!valid) return

  if (items.value.length === 0) {
    notifyError('Item penerimaan kosong')
    return
  }

  const payload = {
    ...form.value,
    items: items.value,
  }

  try {
    if (isEditMode.value) {
      await scStore.updateGoodsReceipt(payload.id, payload)
      notifySuccess('GR berhasil diupdate')
    } else {
      await scStore.createGoodsReceipt(payload)
      notifySuccess('GR berhasil dibuat')
    }
    isDialogOpen.value = false
    loadData(pagination.value.page)
  } catch (error) {
    notifyError(error.message || 'Gagal menyimpan GR')
  }
}

function handleDelete(row) {
  confirmDelete('GR', async () => {
    try {
      await scStore.deleteGoodsReceipt(row.id)
      notifySuccess('GR berhasil dihapus')
      loadData(pagination.value.page)
    } catch (error) {
      notifyError(error.message || 'Gagal menghapus GR')
    }
  })
}

onMounted(async () => {
  await Promise.all([
    loadData(),
    scStore.fetchPurchaseOrder({ limit: 100 }), // Load PO for options
    hrStore.fetchKaryawan({ limit: 100 }),
  ])
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

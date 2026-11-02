<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Purchase Request (PR)"
      subtitle="Permintaan pengadaan barang/jasa internal"
      icon="request_page"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    >
      <template #actions>
        <q-btn unelevated no-caps class="clay-btn clay-btn-primary" @click="openDialog()">
          <q-icon name="add" size="xs" class="q-mr-sm" /> Tambah PR
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
              placeholder="Cari PR..."
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
          :data="scStore.purchaseRequestList"
          :columns="columns"
          :loading="scStore.isLoading('pr')"
          :pagination="pagination"
          @request="onRequest"
        >
          <template #cell-pr_number="{ row }">
            <div class="text-weight-medium">{{ row.pr_number }}</div>
          </template>

          <template #cell-request_date="{ row }">
            {{ formatDate(row.request_date) }}
          </template>

          <template #cell-priority="{ row }">
            <AppBadge :variant="getPriorityBadgeVariant(row.priority)">
              {{ row.priority }}
            </AppBadge>
          </template>

          <template #cell-approval_status="{ row }">
            <AppBadge :variant="getStatusBadgeVariant(row.approval_status)">
              {{ row.approval_status }}
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
      :title="isEditMode ? 'Edit PR' : 'Tambah PR Baru'"
      :loading="scStore.isLoading('pr')"
      width="860px"
      @save="saveData"
      @cancel="isDialogOpen = false"
    >
      <q-form ref="formRef" class="q-gutter-md">
        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6">
            <label class="form-label">No. PR</label>
            <q-input
              v-model="form.pr_number"
              outlined
              dense
              class="clay-input"
              readonly
              placeholder="Auto-generate"
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Departemen <span class="req">*</span></label>
            <q-select
              v-model="form.department_id"
              :options="hrStore.departemenOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              :rules="[required('Departemen wajib diisi')]"
              lazy-rules
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Diajukan Oleh</label>
            <q-input
              v-model="form.requester_name"
              outlined
              dense
              class="clay-input"
              readonly
            />
          </div>

          <div class="col-12 col-md-6">
            <label class="form-label">Tanggal <span class="req">*</span></label>
            <q-input
              v-model="form.request_date"
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
            <label class="form-label">Prioritas <span class="req">*</span></label>
            <q-select
              v-model="form.priority"
              :options="['LOW', 'MEDIUM', 'HIGH', 'URGENT']"
              outlined
              dense
              class="clay-input"
              :rules="[required('Prioritas wajib diisi')]"
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
        </div>

        <div class="q-mt-lg">
          <div class="row justify-between items-center q-mb-sm">
            <div class="text-subtitle1 text-weight-bold">Item PR</div>
            <q-btn flat dense no-caps color="primary" @click="addItem()">
              <q-icon name="add" size="xs" class="q-mr-xs" /> Tambah Item
            </q-btn>
          </div>

          <div style="overflow-x: auto; width: 100%; border-radius: 8px;">
            <q-markup-table flat bordered class="clay-table">
            <thead>
              <tr>
                <th class="text-left">Item <span class="req">*</span></th>
                <th class="text-right">Kuantitas <span class="req">*</span></th>
                <th class="text-right">Estimasi Harga <span class="req">*</span></th>
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
                    v-model.number="item.estimated_price"
                    type="number"
                    min="0"
                    outlined
                    dense
                    class="clay-input"
                    :rules="[
                      requiredNumber('Harga wajib diisi'),
                      nonNegative('Tidak boleh negatif'),
                    ]"
                    @update:model-value="(val) => updateItem(index, 'estimated_price', val)"
                  />
                </td>
                <td class="text-right text-weight-medium">
                  {{ formatCurrency((item.quantity || 0) * (item.estimated_price || 0)) }}
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
                <td colspan="3" class="text-right text-weight-bold">Total Estimasi</td>
                <td class="text-right text-weight-bold text-primary">
                  {{ formatCurrency(total) }}
                </td>
                <td></td>
              </tr>
            </tfoot>
            </q-markup-table>
          </div>
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue'
import { useAuthStore } from 'src/stores/auth.store'
import { useScStore } from 'src/stores/sc.store'
import { useHrStore } from 'src/stores/hr.store'
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
import ScItemSelect from 'src/modules/sc/components/ScItemSelect.vue'

const authStore = useAuthStore()
const scStore = useScStore()
const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const { pagination, searchQuery, onRequest, resetPage } = useTablePagination(loadData)

const isDialogOpen = ref(false)
const isEditMode = ref(false)
const formRef = ref(null)

const columns = [
  { name: 'pr_number', label: 'No. PR', field: 'pr_number', align: 'left', sortable: true },
  {
    name: 'department_name',
    label: 'Departemen',
    field: 'department_name',
    align: 'left',
    sortable: true,
  },
  {
    name: 'requested_by_name',
    label: 'Diajukan Oleh',
    field: (row) => row.requester_name || row.requested_by_name || '-',
    align: 'left',
    sortable: true,
  },
  { name: 'request_date', label: 'Tanggal', field: 'request_date', align: 'left', sortable: true },
  { name: 'priority', label: 'Prioritas', field: 'priority', align: 'center', sortable: true },
  {
    name: 'approval_status',
    label: 'Status',
    field: 'approval_status',
    align: 'center',
    sortable: true,
  },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' },
]

const initialForm = {
  id: null,
  pr_number: '',
  department_id: null,
  requester_name: '',
  request_date: '',
  priority: 'MEDIUM',
  notes: '',
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
    await scStore.fetchPurchaseRequest({
      page,
      limit: pagination.value.rowsPerPage,
      search: searchQuery.value,
    })
    pagination.value.page = scStore.purchaseRequestMeta.page
    pagination.value.rowsNumber = scStore.purchaseRequestMeta.total
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan')
  }
}

function openDialog(row = null) {
  isEditMode.value = !!row
  if (row) {
    form.value = {
      ...row,
      requester_name: row.requester_name || row.requested_by_name || '',
      items: [...(row.items || [])],
    }
  } else {
    form.value = {
      ...initialForm,
      requester_name: authStore.currentUser?.name || authStore.currentUser?.username || 'Super Admin',
    }
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
    requester_name: form.value.requester_name || authStore.currentUser?.name || authStore.currentUser?.username || 'Super Admin',
    items: items.value,
  }

  try {
    if (isEditMode.value) {
      await scStore.updatePurchaseRequest(payload.id, payload)
      notifySuccess('PR berhasil diupdate')
    } else {
      await scStore.createPurchaseRequest(payload)
      notifySuccess('PR berhasil dibuat')
    }
    isDialogOpen.value = false
    loadData(pagination.value.page)
  } catch (error) {
    error.message || 'Gagal menyimpan PR'
  }
}

function handleDelete(row) {
  confirmDelete('PR', async () => {
    try {
      await scStore.deletePurchaseRequest(row.id)
      notifySuccess('PR berhasil dihapus')
      loadData(pagination.value.page)
    } catch (error) {
      error.message || 'Gagal menghapus PR'
    }
  })
}


function getPriorityBadgeVariant(prio) {
  const map = {
    LOW: 'neutral',
    MEDIUM: 'info',
    HIGH: 'warning',
    URGENT: 'danger',
  }
  return map[prio] || 'neutral'
}

function getStatusBadgeVariant(status) {
  const map = {
    DRAFT: 'neutral',
    SUBMITTED: 'warning',
    APPROVED: 'success',
    REJECTED: 'danger',
  }
  return map[status] || 'neutral'
}

onMounted(async () => {
  await Promise.all([loadData(), hrStore.fetchDepartemen({ limit: 100 })])
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

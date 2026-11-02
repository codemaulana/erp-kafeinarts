<template>
  <q-page class="app-page clay-animate-up">
    <!-- Header -->
    <AppPageHeader
      title="Master Vendor / Supplier"
      subtitle="Kelola data mitra pengadaan barang dan jasa"
      icon="storefront"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    />

    <!-- Filter & Toolbar -->
    <div class="row q-col-gutter-md q-mb-md items-center">
      <div class="col-12 col-md-4">
        <q-input
          v-model="search"
          outlined
          dense
          class="clay-input"
          placeholder="Cari nama atau contact vendor..."
          debounce="350"
          @update:model-value="onSearch"
        >
          <template #append>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>
      <div class="col-12 col-md-3">
        <q-select
          v-model="statusFilter"
          outlined
          dense
          class="clay-input"
          :options="['ACTIVE', 'INACTIVE']"
          label="Filter Status"
          clearable
          emit-value
          map-options
          @update:model-value="onSearch"
        />
      </div>
      <div class="col-12 col-md-5 text-right">
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Vendor"
          @click="openAddDialog"
        />
      </div>
    </div>

    <!-- Data Table -->
    <AppDataTable
      :columns="columns"
      :data="scStore.vendorList"
      :loading="scStore.isLoading('vendor')"
      :pagination="pagination"
      row-key="id"
      @request="onRequest"
    >
      <template #body-cell-active_status="props">
        <q-td :props="props">
          <AppBadge
            :variant="props.row.active_status === 'ACTIVE' ? 'success' : 'neutral'"
            :label="props.row.active_status === 'ACTIVE' ? 'Aktif' : 'Non-Aktif'"
          />
        </q-td>
      </template>

      <template #body-cell-actions="props">
        <q-td :props="props">
          <q-btn
            flat
            dense
            round
            size="sm"
            color="primary"
            icon="edit"
            @click="openEditDialog(props.row)"
          >
            <q-tooltip>Edit</q-tooltip>
          </q-btn>
          <q-btn
            flat
            dense
            round
            size="sm"
            color="negative"
            icon="delete"
            @click="onDelete(props.row)"
          >
            <q-tooltip>Hapus</q-tooltip>
          </q-btn>
        </q-td>
      </template>
    </AppDataTable>

    <!-- Dialog Form -->
    <AppFormDialog
      v-model="dialog"
      :title="isEdit ? 'Edit Vendor' : 'Tambah Vendor'"
      :loading="scStore.isLoading('vendor')"
      maximized
      @submit="onSubmit"
    >
      <div class="inv-form-grid">
        <div class="left-col">
          <div class="q-mb-md">
            <label class="form-label">Nama Vendor <span class="req">*</span></label>
            <q-input
              v-model="form.vendor_name"
              outlined
              dense
              class="clay-input"
              :rules="[required('Nama vendor wajib diisi')]"
              lazy-rules
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">Kategori <span class="req">*</span></label>
            <q-select
              v-model="form.vendor_category"
              outlined
              dense
              class="clay-input"
              :options="['SOFTWARE_LICENSE', 'HARDWARE', 'CLOUD_HOSTING', 'FREELANCER', 'OFFICE_SUPPLIES', 'OTHER']"
              :rules="[required('Kategori wajib diisi')]"
              lazy-rules
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">Contact Person <span class="req">*</span></label>
            <q-input
              v-model="form.contact_person"
              outlined
              dense
              class="clay-input"
              :rules="[required('Contact person wajib diisi')]"
              lazy-rules
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">No. Telepon</label>
            <q-input
              v-model="form.phone_number"
              outlined
              dense
              class="clay-input"
              :rules="form.phone_number ? [phoneNumber('Format nomor telepon tidak valid')] : []"
              lazy-rules
            />
          </div>
        </div>

        <div class="right-col">
          <div class="q-mb-md">
            <label class="form-label">Email</label>
            <q-input
              v-model="form.email"
              outlined
              dense
              class="clay-input"
              :rules="form.email ? [email('Format email tidak valid')] : []"
              lazy-rules
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">Alamat</label>
            <q-input
              v-model="form.address"
              type="textarea"
              outlined
              dense
              class="clay-input"
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">NPWP</label>
            <q-input
              v-model="form.tax_number"
              outlined
              dense
              class="clay-input"
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">No. Rekening</label>
            <q-input
              v-model="form.bank_account_number"
              outlined
              dense
              class="clay-input"
            />
          </div>
          <div class="q-mb-md">
            <label class="form-label">Status</label>
            <q-select
              v-model="form.active_status"
              outlined
              dense
              class="clay-input"
              :options="['ACTIVE', 'INACTIVE']"
            />
          </div>
        </div>
      </div>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useScStore } from 'src/stores/sc.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useNotify } from 'src/composables/useNotify'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { required, phoneNumber, email } from 'src/utils/validators'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const scStore = useScStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const search = ref('')
const statusFilter = ref(null)
const dialog = ref(false)
const isEdit = ref(false)

const form = ref({
  id: null,
  vendor_name: '',
  vendor_category: '',
  contact_person: '',
  phone_number: '',
  email: '',
  address: '',
  tax_number: '',
  bank_account_number: '',
  active_status: 'ACTIVE'
})

const columns = [
  { name: 'vendor_name', field: 'vendor_name', label: 'Nama Vendor', align: 'left', sortable: true },
  { name: 'vendor_category', field: 'vendor_category', label: 'Kategori', align: 'left' },
  { name: 'contact_person', field: 'contact_person', label: 'Contact Person', align: 'left' },
  { name: 'phone_number', field: 'phone_number', label: 'No. Telepon', align: 'left' },
  { name: 'active_status', field: 'active_status', label: 'Status', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchFn: (params) => {
    const qParams = { ...params }
    if (search.value) qParams.search = search.value
    if (statusFilter.value) qParams.active_status = statusFilter.value
    return scStore.fetchVendor(qParams)
  }
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const resetForm = () => {
  form.value = {
    id: null,
    vendor_name: '',
    vendor_category: '',
    contact_person: '',
    phone_number: '',
    email: '',
    address: '',
    tax_number: '',
    bank_account_number: '',
    active_status: 'ACTIVE'
  }
}

const openAddDialog = () => {
  isEdit.value = false
  resetForm()
  dialog.value = true
}

const openEditDialog = (row) => {
  isEdit.value = true
  form.value = { ...row }
  dialog.value = true
}

const onSubmit = async () => {
  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      await scStore.updateVendor(payload.id, payload)
      notifySuccess('Vendor berhasil diupdate')
    } else {
      await scStore.createVendor(payload)
      notifySuccess('Vendor berhasil ditambahkan')
    }
    dialog.value = false
    onRequest({ pagination: pagination.value })
  } catch (err) {
    notifyError(err.message || 'Gagal menyimpan vendor')
  }
}

const onDelete = (row) => {
  confirmDelete('Vendor ini', async () => {
    try {
      await scStore.deleteVendor(row.id)
      notifySuccess('Vendor berhasil dihapus')
      onRequest({ pagination: pagination.value })
    } catch (err) {
      notifyError(err.message || 'Gagal menghapus vendor')
    }
  })
}

onMounted(() => {
  onRequest({ pagination: pagination.value })
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

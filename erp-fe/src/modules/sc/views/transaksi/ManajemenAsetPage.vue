<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Manajemen Aset & Lisensi"
      subtitle="Pencatatan aset hardware dan software"
      icon="computer"
      iconBg="linear-gradient(135deg,#E8F7F2,#BDE8D9)"
      iconColor="#4CAF8A"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Aset"
          @click="openForm(null)"
        />
      </template>
    </AppPageHeader>

    <div class="clay-card q-pa-md q-mt-md">
      <div class="row q-mb-md">
        <div class="col-12 col-md-4">
          <q-input
            v-model="searchQuery"
            outlined
            dense
            class="clay-input"
            placeholder="Cari aset..."
            debounce="350"
          >
            <template #prepend>
              <q-icon name="search" />
            </template>
          </q-input>
        </div>
      </div>

      <AppDataTable
        :columns="columns"
        :data="scStore.assetList"
        :loading="scStore.isLoading('fetchAsset')"
        :pagination="pagination"
        @request="onRequest"
      >
        <!-- Status Badge -->
        <template #body-cell-status="props">
          <q-td :props="props">
            <AppBadge :variant="getStatusVariant(props.row.status)">
              {{ props.row.status }}
            </AppBadge>
          </q-td>
        </template>
        
        <!-- Expiry Badge -->
        <template #body-cell-warranty_or_license_expiry="props">
          <q-td :props="props">
            <AppBadge 
              v-if="props.row.warranty_or_license_expiry"
              :variant="getExpiryVariant(props.row.warranty_or_license_expiry)"
            >
              {{ formatDate(props.row.warranty_or_license_expiry) }}
            </AppBadge>
            <span v-else>-</span>
          </q-td>
        </template>

        <!-- Actions -->
        <template #body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn
              flat
              dense
              round
              size="sm"
              icon="edit"
              color="primary"
              @click="openForm(props.row)"
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
              @click="onDelete(props.row)"
            >
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </AppDataTable>
    </div>

    <!-- Form Dialog -->
    <AppFormDialog
      v-model="isFormOpen"
      :title="isEdit ? 'Edit Aset' : 'Tambah Aset'"
      :loading="scStore.isLoading('saveAsset')"
      @submit="onSubmit"
      @cancel="closeForm"
    >
      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <label class="form-label">Item Aset <span class="req">*</span></label>
          <ScItemSelect
            v-model="form.item_id"
            :rules="[required('Item harus diisi')]"
          />
        </div>
        
        <div class="col-12 col-md-6">
          <label class="form-label">Serial / License Key</label>
          <q-input
            v-model="form.serial_or_license_key"
            outlined
            dense
            class="clay-input"
            lazy-rules
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Tanggal Perolehan</label>
          <q-input
            v-model="form.acquisition_date"
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            outlined
            dense
            class="clay-input"
            lazy-rules
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Nilai Perolehan</label>
          <q-input
            v-model.number="form.acquisition_value"
            type="number"
            outlined
            dense
            class="clay-input"
            lazy-rules
            :rules="[nonNegative('Nilai tidak boleh negatif')]"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Masa Berlaku (Garansi/Lisensi)</label>
          <q-input
            v-model="form.warranty_or_license_expiry"
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            outlined
            dense
            class="clay-input"
            lazy-rules
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Status <span class="req">*</span></label>
          <q-select
            v-model="form.status"
            :options="['AVAILABLE', 'IN_USE', 'MAINTENANCE', 'DAMAGED', 'DISPOSED']"
            outlined
            dense
            class="clay-input"
            :rules="[required('Status harus diisi')]"
            lazy-rules
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Digunakan Oleh</label>
          <q-select
            v-model="form.assigned_to_employee_id"
            :options="hrStore.karyawanOptions"
            emit-value
            map-options
            clearable
            outlined
            dense
            class="clay-input"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Lokasi</label>
          <q-select
            v-model="form.location_id"
            :options="scStore.lokasiOptions"
            emit-value
            map-options
            clearable
            outlined
            dense
            class="clay-input"
          />
        </div>
      </div>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted, watch } from 'vue';
import { useScStore } from 'src/stores/sc.store';
import { useHrStore } from 'src/stores/hr.store';
import { useTablePagination } from 'src/composables/useTablePagination';
import { useConfirmDelete } from 'src/composables/useConfirmDelete';
import { useNotify } from 'src/composables/useNotify';
import { required, nonNegative } from 'src/utils/validators';
import { formatDate } from 'src/utils/formatter';

import AppPageHeader from 'components/shared/AppPageHeader.vue';
import AppDataTable from 'components/shared/AppDataTable.vue';
import AppFormDialog from 'components/shared/AppFormDialog.vue';
import AppBadge from 'components/shared/AppBadge.vue';
import ScItemSelect from 'src/modules/sc/components/ScItemSelect.vue';

const scStore = useScStore();
const hrStore = useHrStore();
const { notifySuccess, notifyError } = useNotify();
const { confirmDelete } = useConfirmDelete();

const columns = [
  { name: 'asset_code', field: 'asset_code', label: 'Kode', align: 'left' },
  { name: 'item_name', field: (row) => row.name || row.item_name || '-', label: 'Nama Aset', align: 'left' },
  { name: 'serial_or_license_key', field: (row) => row.serial_or_license_key || '-', label: 'Serial/License', align: 'left' },
  { name: 'assigned_to_name', field: (row) => row.assigned_to_name || row.assigned_to_employee_name || row.assigned_to || '-', label: 'Digunakan Oleh', align: 'left' },
  { name: 'location_name', field: (row) => row.location_name || '-', label: 'Lokasi', align: 'left' },
  { name: 'status', field: 'status', label: 'Status', align: 'left' },
  { name: 'warranty_or_license_expiry', field: 'warranty_or_license_expiry', label: 'Masa Berlaku', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
];

const fetchAssets = async (params) => {
  await scStore.fetchAsset({ ...params, search: searchQuery.value });
};

const { pagination, onRequest, resetPage } = useTablePagination(fetchAssets);

const searchQuery = ref('');
watch(searchQuery, () => {
  resetPage();
  onRequest({ pagination: pagination.value });
});

onMounted(async () => {
  await hrStore.fetchKaryawan();
  await scStore.fetchLokasi?.() || console.warn('fetchLokasi not available in scStore');
  onRequest({ pagination: pagination.value });
});

const isFormOpen = ref(false);
const isEdit = ref(false);
const currentId = ref(null);

const defaultForm = () => ({
  item_id: null,
  serial_or_license_key: '',
  acquisition_date: '',
  acquisition_value: 0,
  warranty_or_license_expiry: '',
  status: 'AVAILABLE',
  assigned_to_employee_id: null,
  location_id: null
});
const form = ref(defaultForm());

const openForm = (row) => {
  if (row) {
    isEdit.value = true;
    currentId.value = row.id;
    form.value = { ...row };
  } else {
    isEdit.value = false;
    currentId.value = null;
    form.value = defaultForm();
  }
  isFormOpen.value = true;
};

const closeForm = () => {
  isFormOpen.value = false;
  form.value = defaultForm();
};

const onSubmit = async () => {
  try {
    if (isEdit.value) {
      await scStore.updateAsset(currentId.value, form.value);
      notifySuccess('Aset berhasil diperbarui');
    } else {
      await scStore.createAsset(form.value);
      notifySuccess('Aset berhasil ditambahkan');
    }
    closeForm();
    onRequest({ pagination: pagination.value });
  } catch (error) {
    notifyError(error.message || 'Gagal menyimpan aset');
  }
};

const onDelete = (row) => {
  confirmDelete({
    title: 'Hapus Aset',
    message: `Apakah Anda yakin ingin menghapus aset ${row.asset_code}?`,
    onConfirm: async () => {
      try {
        await scStore.deleteAsset(row.id);
        notifySuccess('Aset berhasil dihapus');
        onRequest({ pagination: pagination.value });
      } catch (error) {
        notifyError(error.message || 'Gagal menghapus aset');
      }
    }
  });
};

const getStatusVariant = (status) => {
  const map = {
    AVAILABLE: 'success',
    IN_USE: 'info',
    MAINTENANCE: 'warning',
    DAMAGED: 'danger',
    DISPOSED: 'neutral'
  };
  return map[status] || 'neutral';
};

const getExpiryVariant = (dateStr) => {
  if (!dateStr) return 'neutral';
  const target = new Date(dateStr);
  const now = new Date();
  const diffTime = target - now;
  const diffDays = Math.ceil(diffTime / (1000 * 60 * 60 * 24));
  return diffDays <= 30 ? 'warning' : 'neutral';
};
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Perpanjangan Langganan"
      subtitle="Pemantauan layanan SaaS dan cloud"
      icon="cloud_sync"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Langganan"
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
            placeholder="Cari langganan..."
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
        :data="scStore.subscriptionList"
        :loading="scStore.isLoading('fetchSubscription')"
        :pagination="pagination"
        @request="onRequest"
      >
        <!-- Amount -->
        <template #body-cell-amount="props">
          <q-td :props="props">
            {{ formatCurrency(props.row.amount) }}
          </q-td>
        </template>
        
        <!-- Status Badge -->
        <template #body-cell-status="props">
          <q-td :props="props">
            <AppBadge :variant="getStatusVariant(props.row.status)">
              {{ props.row.status }}
            </AppBadge>
          </q-td>
        </template>
        
        <!-- End Date -->
        <template #body-cell-end_date="props">
          <q-td :props="props">
            {{ formatDate(props.row.end_date) }}
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
      :title="isEdit ? 'Edit Langganan' : 'Tambah Langganan'"
      :loading="scStore.isLoading('saveSubscription')"
      @submit="onSubmit"
      @cancel="closeForm"
    >
      <div class="row q-col-gutter-md">
        <div class="col-12 col-md-6">
          <label class="form-label">Item Langganan <span class="req">*</span></label>
          <ScItemSelect
            v-model="form.item_id"
            :rules="[required('Item harus diisi')]"
          />
        </div>
        
        <div class="col-12 col-md-6">
          <label class="form-label">Vendor <span class="req">*</span></label>
          <ScVendorSelect
            v-model="form.vendor_id"
            :rules="[required('Vendor harus diisi')]"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Siklus Tagihan</label>
          <q-select
            v-model="form.billing_cycle"
            :options="['MONTHLY', 'YEARLY']"
            outlined
            dense
            class="clay-input"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Nominal</label>
          <q-input
            v-model.number="form.amount"
            type="number"
            outlined
            dense
            class="clay-input"
            lazy-rules
            :rules="[nonNegative('Nominal tidak boleh negatif')]"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Tanggal Mulai <span class="req">*</span></label>
          <q-input
            v-model="form.start_date"
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            outlined
            dense
            class="clay-input"
            lazy-rules
            :rules="[required('Tanggal mulai harus diisi')]"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Tanggal Berakhir <span class="req">*</span></label>
          <q-input
            v-model="form.end_date"
            type="text" @focus="e => e.target.type = 'date'" @blur="e => { if(!e.target.value) e.target.type = 'text' }"
            outlined
            dense
            class="clay-input"
            lazy-rules
            :rules="[
              required('Tanggal berakhir harus diisi'),
              dateNotBefore(() => form.start_date, 'Tidak boleh sebelum tanggal mulai')
            ]"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Status</label>
          <q-select
            v-model="form.status"
            :options="['ACTIVE', 'EXPIRED', 'CANCELLED']"
            outlined
            dense
            class="clay-input"
          />
        </div>

        <div class="col-12 col-md-6">
          <label class="form-label">Perpanjang Otomatis</label>
          <q-select
            v-model="form.auto_renewal"
            :options="[
              { label: 'Ya', value: true },
              { label: 'Tidak', value: false }
            ]"
            emit-value
            map-options
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
import { useTablePagination } from 'src/composables/useTablePagination';
import { useConfirmDelete } from 'src/composables/useConfirmDelete';
import { useNotify } from 'src/composables/useNotify';
import { required, nonNegative, dateNotBefore } from 'src/utils/validators';
import { formatCurrency, formatDate } from 'src/utils/formatter';

import AppPageHeader from 'components/shared/AppPageHeader.vue';
import AppDataTable from 'components/shared/AppDataTable.vue';
import AppFormDialog from 'components/shared/AppFormDialog.vue';
import AppBadge from 'components/shared/AppBadge.vue';
import ScItemSelect from 'src/modules/sc/components/ScItemSelect.vue';
import ScVendorSelect from 'src/modules/sc/components/ScVendorSelect.vue';

const scStore = useScStore();
const { notifySuccess, notifyError } = useNotify();
const { confirmDelete } = useConfirmDelete();

const columns = [
  { name: 'item_name', field: 'item_name', label: 'Nama Langganan', align: 'left' },
  { name: 'vendor_name', field: 'vendor_name', label: 'Vendor', align: 'left' },
  { name: 'billing_cycle', field: 'billing_cycle', label: 'Siklus', align: 'left' },
  { name: 'end_date', field: 'end_date', label: 'Tgl Berakhir', align: 'left' },
  { name: 'amount', field: 'amount', label: 'Nominal', align: 'right' },
  { name: 'status', field: 'status', label: 'Status', align: 'center' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
];

const fetchSubscriptions = async (params) => {
  await scStore.fetchSubscription({ ...params, search: searchQuery.value });
};

const { pagination, onRequest, resetPage } = useTablePagination(fetchSubscriptions);

const searchQuery = ref('');
watch(searchQuery, () => {
  resetPage();
  onRequest({ pagination: pagination.value });
});

onMounted(() => {
  onRequest({ pagination: pagination.value });
});

const isFormOpen = ref(false);
const isEdit = ref(false);
const currentId = ref(null);

const defaultForm = () => ({
  item_id: null,
  vendor_id: null,
  billing_cycle: 'MONTHLY',
  amount: 0,
  start_date: '',
  end_date: '',
  status: 'ACTIVE',
  auto_renewal: false
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
      await scStore.updateSubscription(currentId.value, form.value);
      notifySuccess('Langganan berhasil diperbarui');
    } else {
      await scStore.createSubscription(form.value);
      notifySuccess('Langganan berhasil ditambahkan');
    }
    closeForm();
    onRequest({ pagination: pagination.value });
  } catch (error) {
    notifyError(error.message || 'Gagal menyimpan langganan');
  }
};

const onDelete = (row) => {
  confirmDelete({
    title: 'Hapus Langganan',
    message: `Apakah Anda yakin ingin menghapus langganan ini?`,
    onConfirm: async () => {
      try {
        await scStore.deleteSubscription(row.id);
        notifySuccess('Langganan berhasil dihapus');
        onRequest({ pagination: pagination.value });
      } catch (error) {
        notifyError(error.message || 'Gagal menghapus langganan');
      }
    }
  });
};

const getStatusVariant = (status) => {
  const map = {
    ACTIVE: 'success',
    EXPIRED: 'danger',
    CANCELLED: 'neutral'
  };
  return map[status] || 'neutral';
};
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

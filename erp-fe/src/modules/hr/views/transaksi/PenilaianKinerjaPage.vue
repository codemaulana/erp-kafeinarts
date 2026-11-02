<!-- PenilaianKinerjaPage.vue -->
<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Penilaian Kinerja"
      subtitle="Evaluasi kinerja karyawan secara periodik"
      icon="star_rate"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    >
      <q-btn
        unelevated
        no-caps
        icon="add"
        label="Buat Penilaian"
        class="clay-btn clay-btn-primary"
        @click="openDialog()"
      />
    </AppPageHeader>

    <q-card class="app-card">
      <div class="app-card__header flex items-center justify-between">
        <div class="flex items-center gap-sm">
          <q-input
            v-model="search"
            outlined
            dense
            placeholder="Cari penilaian..."
            class="clay-input"
            debounce="350"
            @update:model-value="onSearch"
          >
            <template #prepend><q-icon name="search" /></template>
          </q-input>
          
          <q-select
            v-model="statusFilter"
            :options="[{label:'Semua Status',value:''},{label:'Draft',value:'DRAFT'},{label:'Diajukan',value:'SUBMITTED'},{label:'Difinalisasi',value:'FINALIZED'}]"
            outlined
            dense
            emit-value
            map-options
            class="clay-input"
            style="min-width: 150px"
            @update:model-value="onSearch"
          />
        </div>
      </div>

      <AppDataTable
        :rows="hrStore.kinerjaList"
        :columns="columns"
        :loading="hrStore.isLoading('kinerja')"
        :pagination="pagination"
        @request="onRequest"
      >
        <template #body-cell-total_score="props">
          <q-td :props="props" class="text-center text-bold">
            {{ formatScore(props.row.total_score) }}
          </q-td>
        </template>

        <template #body-cell-status="props">
          <q-td :props="props" class="text-center">
            <AppBadge :variant="getStatusVariant(props.row.status)">
              {{ getStatusLabel(props.row.status) }}
            </AppBadge>
          </q-td>
        </template>

        <template #body-cell-actions="props">
          <q-td :props="props" class="text-center">
            <q-btn flat dense round size="sm" icon="visibility" color="info" @click="openDialog(props.row, true)">
              <q-tooltip>Lihat Detail</q-tooltip>
            </q-btn>
            <q-btn flat dense round size="sm" icon="edit" color="primary" @click="openDialog(props.row)">
              <q-tooltip>Edit</q-tooltip>
            </q-btn>
            <q-btn flat dense round size="sm" icon="delete" color="negative" @click="onDelete(props.row.id)">
              <q-tooltip>Hapus</q-tooltip>
            </q-btn>
          </q-td>
        </template>
      </AppDataTable>
    </q-card>

    <AppFormDialog
      v-model="dialog"
      :title="isView ? 'Detail Penilaian Kinerja' : (isEdit ? 'Edit Penilaian Kinerja' : 'Buat Penilaian Kinerja')"
      :loading="hrStore.isLoading('kinerja')"
      :maximized="true"
      @save="onSubmit"
    >
      <div class="inv-form-grid">
        <div class="inv-form-grid__left">
          <div class="text-subtitle1 text-weight-bold q-mb-md">Informasi Penilaian</div>
          
          <div class="form-group q-mb-md">
            <label class="form-label">Karyawan <span class="req">*</span></label>
            <q-select
              v-model="form.employee_id"
              :options="hrStore.karyawanOptions"
              outlined
              dense
              emit-value
              map-options
              class="clay-input"
              :rules="[required('Karyawan wajib dipilih')]"
              lazy-rules
              :readonly="isView"
            />
          </div>

          <div class="form-group q-mb-md">
            <label class="form-label">Periode <span class="req">*</span></label>
            <q-select
              v-model="form.review_period"
              :options="PERIOD_OPTIONS"
              outlined
              dense
              emit-value
              map-options
              class="clay-input"
              :rules="[required('Periode wajib dipilih')]"
              lazy-rules
              :readonly="isView"
            />
          </div>

          <div class="form-group q-mb-md">
            <label class="form-label">Penilai <span class="req">*</span></label>
            <q-select
              v-model="form.reviewer_id"
              :options="hrStore.karyawanOptions"
              outlined
              dense
              emit-value
              map-options
              class="clay-input"
              :rules="[required('Penilai wajib dipilih')]"
              lazy-rules
              :readonly="isView"
            />
          </div>

          <div class="form-group q-mb-md">
            <label class="form-label">Status <span class="req">*</span></label>
            <q-select
              v-model="form.status"
              :options="[{label:'Draft',value:'DRAFT'},{label:'Diajukan',value:'SUBMITTED'},{label:'Difinalisasi',value:'FINALIZED'}]"
              outlined
              dense
              emit-value
              map-options
              class="clay-input"
              :rules="[required('Status wajib dipilih')]"
              lazy-rules
              :readonly="isView"
            />
          </div>

          <div class="form-group q-mb-md">
            <label class="form-label">Rekomendasi</label>
            <q-input
              v-model="form.recommendation"
              type="textarea"
              rows="3"
              outlined
              dense
              class="clay-input"
              :readonly="isView"
            />
          </div>
        </div>

        <div class="inv-form-grid__right">
          <div class="flex items-center justify-between q-mb-md">
            <div class="text-subtitle1 text-weight-bold">Kriteria Penilaian</div>
            <q-btn
              v-if="!isView"
              outline
              no-caps
              icon="add"
              label="Tambah Kriteria"
              color="primary"
              size="sm"
              @click="addItem"
            />
          </div>

          <div class="inv-items">
            <div class="inv-items__header" style="grid-template-columns: 1fr 80px 1fr 36px">
              <div>Aspek Penilaian</div>
              <div class="text-center">Skor (1-5)</div>
              <div>Catatan</div>
              <div></div>
            </div>

            <div v-if="!form.items || form.items.length === 0" class="q-pa-md text-center text-grey">
              Belum ada kriteria penilaian
            </div>

            <div
              v-for="(item, index) in form.items"
              :key="item.id"
              class="inv-items__row"
              style="grid-template-columns: 1fr 80px 1fr 36px"
            >
              <q-input
                v-model="item.criteria_name"
                outlined
                dense
                class="clay-input"
                :rules="[required('Aspek wajib diisi')]"
                lazy-rules
                :readonly="isView"
              />
              <q-input
                v-model.number="item.score"
                type="number"
                outlined
                dense
                class="clay-input text-center"
                :rules="[requiredNumber('Skor wajib'), inRange(1, 5, '1-5')]"
                lazy-rules
                :readonly="isView"
              />
              <q-input
                v-model="item.notes"
                outlined
                dense
                class="clay-input"
                :readonly="isView"
              />
              <div class="flex flex-center">
                <q-btn
                  v-if="!isView"
                  flat
                  round
                  dense
                  color="negative"
                  icon="remove_circle_outline"
                  @click="removeItem(index)"
                />
              </div>
            </div>
          </div>

          <div class="inv-totals q-mt-lg">
            <div class="inv-totals__row inv-totals--grand">
              <div class="inv-totals__label">Skor Total (Rata-rata)</div>
              <div class="inv-totals__value text-primary">{{ formatScore(totalScore) }}</div>
            </div>
          </div>

        </div>
      </div>
      
      <template #actions v-if="isView">
        <q-btn
          unelevated
          no-caps
          label="Tutup"
          class="clay-btn"
          color="grey-3"
          text-color="grey-8"
          v-close-popup
        />
      </template>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, computed, onMounted } from 'vue'
import { uid } from 'quasar'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, inRange } from 'src/utils/validators'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const search = ref('')
const statusFilter = ref('')

const PERIOD_OPTIONS = [
  { label: 'Triwulan 1', value: 'Q1' },
  { label: 'Triwulan 2', value: 'Q2' },
  { label: 'Triwulan 3', value: 'Q3' },
  { label: 'Triwulan 4', value: 'Q4' },
  { label: 'Semester 1', value: 'S1' },
  { label: 'Semester 2', value: 'S2' },
  { label: 'Tahunan', value: 'ANNUAL' }
]

const columns = [
  { name: 'employee_name', label: 'Karyawan', field: 'employee_name', align: 'left', sortable: true },
  { name: 'review_period', label: 'Periode', field: 'review_period', align: 'left' },
  { name: 'reviewer_name', label: 'Penilai', field: 'reviewer_name', align: 'left' },
  { name: 'total_score', label: 'Skor Total', field: 'total_score', align: 'center' },
  { name: 'status', label: 'Status', field: 'status', align: 'center' },
  { name: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, load, onRequest, resetPage } = useTablePagination({
  fetchFn: hrStore.fetchKinerja,
  getParams: () => ({
    search: search.value,
    status: statusFilter.value
  }),
  meta: computed(() => hrStore.kinerjaMeta)
})

const onSearch = () => {
  resetPage()
  load()
}

// Dialog
const dialog = ref(false)
const isEdit = ref(false)
const isView = ref(false)
const currentId = ref(null)

const form = ref({
  employee_id: '',
  review_period: '',
  reviewer_id: '',
  status: 'DRAFT',
  recommendation: '',
  items: []
})

const resetForm = () => {
  form.value = {
    employee_id: '',
    review_period: '',
    reviewer_id: '',
    status: 'DRAFT',
    recommendation: '',
    items: []
  }
  isEdit.value = false
  isView.value = false
  currentId.value = null
}

const openDialog = (row = null, view = false) => {
  resetForm()
  if (row) {
    isEdit.value = !view
    isView.value = view
    currentId.value = row.id
    form.value = { ...row, items: Array.isArray(row.items) ? [...row.items.map(i => ({...i, id: i.id || uid()}))] : [] }
  } else {
    addItem()
  }
  dialog.value = true
}

const addItem = () => {
  form.value.items.push({ id: uid(), criteria_name: '', score: 0, notes: '' })
}

const removeItem = (index) => {
  form.value.items.splice(index, 1)
}

// Computed total score
const totalScore = computed(() => {
  if (form.value.items.length === 0) return 0
  const sum = form.value.items.reduce((acc, item) => acc + (Number(item.score) || 0), 0)
  return sum / form.value.items.length
})

const formatScore = (val) => {
  if (val === null || val === undefined) return '0.0'
  return Number(val).toFixed(1)
}

const onSubmit = async () => {
  try {
    const payload = {
      ...form.value,
      total_score: totalScore.value
    }

    if (isEdit.value) {
      await hrStore.updateKinerja(currentId.value, payload)
      notifySuccess('Penilaian kinerja berhasil diperbarui')
    } else {
      await hrStore.createKinerja(payload)
      notifySuccess('Penilaian kinerja berhasil dibuat')
    }
    dialog.value = false
    load()
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan saat menyimpan penilaian kinerja')
  }
}

const onDelete = async (id) => {
  const confirmed = await confirmDelete('Apakah Anda yakin ingin menghapus penilaian kinerja ini?')
  if (confirmed) {
    try {
      await hrStore.deleteKinerja(id)
      notifySuccess('Penilaian kinerja berhasil dihapus')
      load()
    } catch (error) {
      notifyError(error.message || 'Gagal menghapus penilaian kinerja')
    }
  }
}

const getStatusVariant = (status) => {
  const map = {
    DRAFT: 'neutral',
    SUBMITTED: 'warning',
    FINALIZED: 'success'
  }
  return map[status] || 'neutral'
}

const getStatusLabel = (status) => {
  const map = {
    DRAFT: 'Draft',
    SUBMITTED: 'Diajukan',
    FINALIZED: 'Difinalisasi'
  }
  return map[status] || status
}

onMounted(async () => {
  await hrStore.fetchKaryawan({ page: 1, limit: 100 })
  load()
})
</script>

<style lang="scss">
@import 'src/css/finance-page';
</style>

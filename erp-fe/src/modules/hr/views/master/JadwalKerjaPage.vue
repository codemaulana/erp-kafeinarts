<template>
  <q-page class="app-page clay-animate-up">
    <AppPageHeader
      title="Jadwal Kerja / Shift"
      subtitle="Pengaturan jam kerja dan jadwal shift karyawan"
      icon="schedule"
      iconBg="linear-gradient(135deg,#FEF3E2,#F9D9A0)"
      iconColor="#E8A547"
    >
      <template #actions>
        <q-btn
          unelevated
          no-caps
          class="clay-btn clay-btn-primary"
          icon="add"
          label="Tambah Jadwal"
          @click="openForm()"
        />
      </template>
    </AppPageHeader>

    <div class="app-card clay-card q-pa-md">
      <div class="row q-mb-md">
        <q-input
          v-model="searchQuery"
          outlined
          dense
          class="clay-input"
          placeholder="Cari jadwal kerja..."
          debounce="350"
          @update:model-value="onSearch"
        >
          <template #prepend>
            <q-icon name="search" />
          </template>
        </q-input>
      </div>

      <AppDataTable
        :rows="hrStore.jadwalKerjaList"
        :columns="columns"
        :loading="hrStore.isLoading('jadwalKerja')"
        :pagination="pagination"
        @request="onRequest"
      >
        <template #body-cell-late_tolerance_minutes="props">
          <q-td :props="props">
            {{ props.value }} menit
          </q-td>
        </template>
        <template #body-cell-work_days="props">
          <q-td :props="props">
            <div class="row q-gutter-xs">
              <AppBadge
                v-for="day in props.value"
                :key="day"
                variant="info"
                :label="DAY_LABELS[day] || day"
              />
            </div>
          </q-td>
        </template>
        <template #body-cell-actions="props">
          <q-td :props="props">
            <div class="row justify-center q-gutter-xs">
              <q-btn flat dense round size="sm" icon="edit" color="primary" @click="openForm(props.row)">
                <q-tooltip>Edit</q-tooltip>
              </q-btn>
              <q-btn flat dense round size="sm" icon="delete" color="negative" @click="handleDelete(props.row)">
                <q-tooltip>Hapus</q-tooltip>
              </q-btn>
            </div>
          </q-td>
        </template>
      </AppDataTable>
    </div>

    <AppFormDialog
      v-model="isDialogOpen"
      :title="isEdit ? 'Edit Jadwal Kerja' : 'Tambah Jadwal Kerja'"
      @save="onSave"
      :loading="hrStore.isLoading('jadwalKerja_saving')"
      saveLabel="Simpan"
    >
      <q-form ref="formRef" @submit.prevent="onSave">
        <div class="q-mb-md">
          <label class="form-label">Nama Shift <span class="req">*</span></label>
          <q-input v-model="form.shift_name" outlined dense class="clay-input" :rules="[required('Nama shift wajib diisi')]" lazy-rules />
        </div>

        <div class="row q-col-gutter-md">
          <div class="col-12 col-md-6 q-mb-md">
            <label class="form-label">Jam Masuk <span class="req">*</span></label>
            <q-input v-model="form.time_in" type="time" outlined dense class="clay-input" :rules="[required('Jam masuk wajib diisi')]" lazy-rules />
          </div>
          <div class="col-12 col-md-6 q-mb-md">
            <label class="form-label">Jam Pulang <span class="req">*</span></label>
            <q-input v-model="form.time_out" type="time" outlined dense class="clay-input" :rules="[required('Jam pulang wajib diisi')]" lazy-rules />
          </div>
        </div>

        <div class="q-mb-md">
          <label class="form-label">Toleransi Terlambat (menit) <span class="req">*</span></label>
          <q-input v-model.number="form.late_tolerance_minutes" type="number" min="0" outlined dense class="clay-input" :rules="[requiredNumber('Toleransi wajib diisi'), nonNegative('Tidak boleh negatif')]" lazy-rules />
        </div>

        <div class="q-mb-md">
          <label class="form-label">Hari Kerja</label>
          <q-option-group
            v-model="form.work_days"
            :options="workDayOptions"
            type="checkbox"
            inline
          />
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useHrStore } from 'src/stores/hr.store'
import { useTablePagination } from 'src/composables/useTablePagination'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { useNotify } from 'src/composables/useNotify'
import { required, requiredNumber, nonNegative } from 'src/utils/validators'
import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const hrStore = useHrStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const isDialogOpen = ref(false)
const isEdit = ref(false)
const formRef = ref(null)
const searchQuery = ref('')

const DAY_LABELS = {
  MON: 'Sen',
  TUE: 'Sel',
  WED: 'Rab',
  THU: 'Kam',
  FRI: 'Jum',
  SAT: 'Sab',
  SUN: 'Min'
}

const workDayOptions = [
  { label: 'Senin', value: 'MON' },
  { label: 'Selasa', value: 'TUE' },
  { label: 'Rabu', value: 'WED' },
  { label: 'Kamis', value: 'THU' },
  { label: 'Jumat', value: 'FRI' },
  { label: 'Sabtu', value: 'SAT' },
  { label: 'Minggu', value: 'SUN' }
]

const defaultForm = {
  id: null,
  shift_name: '',
  time_in: '',
  time_out: '',
  late_tolerance_minutes: 0,
  work_days: []
}

const form = ref({ ...defaultForm })

const columns = [
  { name: 'shift_name', field: 'shift_name', label: 'Nama Shift', align: 'left', sortable: true },
  { name: 'time_in', field: 'time_in', label: 'Jam Masuk', align: 'center' },
  { name: 'time_out', field: 'time_out', label: 'Jam Pulang', align: 'center' },
  { name: 'late_tolerance_minutes', field: 'late_tolerance_minutes', label: 'Toleransi Terlambat', align: 'center' },
  { name: 'work_days', field: 'work_days', label: 'Hari Kerja', align: 'left' },
  { name: 'actions', field: 'actions', label: 'Aksi', align: 'center' }
]

const { pagination, onRequest, resetPage } = useTablePagination({
  fetchData: (params) => hrStore.fetchJadwalKerja({ ...params, search: searchQuery.value }),
  defaultSort: 'shift_name'
})

const onSearch = () => {
  resetPage()
  onRequest({ pagination: pagination.value })
}

const openForm = (row = null) => {
  isEdit.value = !!row
  if (row) {
    form.value = { ...row, work_days: row.work_days ? [...row.work_days] : [] }
  } else {
    form.value = { ...defaultForm, work_days: [] }
  }
  isDialogOpen.value = true
}

const onSave = async () => {
  const valid = await formRef.value.validate()
  if (!valid) return

  try {
    const payload = { ...form.value }
    if (isEdit.value) {
      await hrStore.updateJadwalKerja(payload.id, payload)
      notifySuccess('Jadwal kerja berhasil diupdate')
    } else {
      await hrStore.createJadwalKerja(payload)
      notifySuccess('Jadwal kerja berhasil ditambahkan')
    }
    isDialogOpen.value = false
    onRequest({ pagination: pagination.value })
  } catch (error) {
    notifyError(error.message || 'Terjadi kesalahan')
  }
}

const handleDelete = (row) => {
  confirmDelete({
    title: 'Hapus Jadwal Kerja',
    message: `Apakah Anda yakin ingin menghapus jadwal ${row.shift_name}?`,
    onConfirm: async () => {
      try {
        await hrStore.deleteJadwalKerja(row.id)
        notifySuccess('Jadwal kerja berhasil dihapus')
        onRequest({ pagination: pagination.value })
      } catch (error) {
        notifyError(error.message || 'Gagal menghapus jadwal kerja')
      }
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

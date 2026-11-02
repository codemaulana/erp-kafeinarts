<template>
  <q-page class="app-page clay-animate-up">
    <!-- Header -->
    <AppPageHeader title="Manajemen Pengguna" icon="manage_accounts">
      <template #actions>
        <q-btn
          unelevated
          no-caps
          icon="add"
          label="Tambah Pengguna"
          class="clay-btn clay-btn-primary"
          @click="openDialog()"
        />
      </template>
    </AppPageHeader>

    <!-- Filter Card -->
    <q-card class="app-card q-mb-md">
      <q-card-section>
        <div class="row q-col-gutter-sm items-center">
          <div class="col-12 col-md-5">
            <q-input
              v-model="searchQuery"
              outlined
              dense
              class="clay-input"
              placeholder="Cari nama atau email pengguna..."
              clearable
              @update:model-value="onFilter"
            >
              <template #prepend>
                <q-icon name="search" />
              </template>
            </q-input>
          </div>
          <div class="col-12 col-md-4">
            <q-select
              v-model="selectedRole"
              :options="roleOptions"
              outlined
              dense
              class="clay-input"
              label="Filter Role"
              emit-value
              map-options
              clearable
              @update:model-value="onFilter"
            />
          </div>
          <div class="col-12 col-md-3">
            <q-btn
              flat
              no-caps
              icon="refresh"
              label="Reset"
              class="clay-btn"
              @click="resetFilter"
            />
          </div>
        </div>
      </q-card-section>
    </q-card>

    <!-- Data Table Card -->
    <q-card class="app-card">
      <AppDataTable
        :columns="columns"
        :data="userStore.userList"
        :loading="userStore.isLoading"
      >
        <template #cell-name="{ row }">
          <div class="row items-center no-wrap">
            <q-avatar size="32px" class="q-mr-sm text-weight-bold" color="primary" text-color="white">
              {{ getUserInitial(row.name) }}
            </q-avatar>
            <div>
              <div class="text-weight-bold text-dark-mode-bright">{{ row.name }}</div>
              <div v-if="isCurrentLoggedInUser(row)" class="text-caption text-primary text-weight-medium">
                (Akun Anda Saat Ini)
              </div>
            </div>
          </div>
        </template>

        <template #cell-role="{ row }">
          <AppBadge :variant="getRoleVariant(row.role)">
            {{ formatRole(row.role) }}
          </AppBadge>
        </template>

        <template #cell-created_at="{ row }">
          <span class="text-caption text-grey-6">{{ formatDate(row.created_at) }}</span>
        </template>

        <template #cell-actions="{ row }">
          <div class="row items-center justify-center no-wrap q-gutter-xs">
            <q-btn
              flat
              dense
              round
              size="sm"
              icon="edit"
              color="primary"
              @click="openDialog(row)"
            >
              <q-tooltip>Edit Pengguna</q-tooltip>
            </q-btn>

            <q-btn
              flat
              dense
              round
              size="sm"
              icon="delete"
              :color="isCurrentLoggedInUser(row) ? 'grey-5' : 'negative'"
              :disable="isCurrentLoggedInUser(row)"
              @click="handleDelete(row)"
            >
              <q-tooltip>{{ isCurrentLoggedInUser(row) ? 'Tidak dapat menghapus akun sendiri' : 'Hapus Pengguna' }}</q-tooltip>
            </q-btn>
          </div>
        </template>
      </AppDataTable>
    </q-card>

    <!-- Dialog Form Tambah / Edit Pengguna -->
    <AppFormDialog
      v-model="isDialogOpen"
      :title="isEditMode ? 'Edit Pengguna' : 'Tambah Pengguna Baru'"
      :icon="isEditMode ? 'edit' : 'person_add'"
      :loading="userStore.isLoading"
      width="640px"
      @save="saveData"
      @cancel="isDialogOpen = false"
    >
      <q-form ref="formRef" class="q-gutter-md">
        <div class="row q-col-gutter-md">
          <!-- Nama -->
          <div class="col-12">
            <label class="form-label">Nama Lengkap <span class="req">*</span></label>
            <q-input
              v-model="form.name"
              outlined
              dense
              class="clay-input"
              placeholder="Contoh: Rian Anggoro"
              :rules="[required('Nama lengkap wajib diisi')]"
              lazy-rules
            />
          </div>

          <!-- Email -->
          <div class="col-12 col-md-6">
            <label class="form-label">Email Login <span class="req">*</span></label>
            <q-input
              v-model="form.email"
              type="email"
              outlined
              dense
              class="clay-input"
              placeholder="contoh@kafeinarts.com"
              :rules="[required('Email wajib diisi')]"
              lazy-rules
            />
          </div>

          <!-- Role -->
          <div class="col-12 col-md-6">
            <label class="form-label">Role / Hak Akses <span class="req">*</span></label>
            <q-select
              v-model="form.role"
              :options="roleOptions"
              outlined
              dense
              class="clay-input"
              emit-value
              map-options
              :rules="[required('Role wajib dipilih')]"
              lazy-rules
            />
          </div>

          <!-- Password Info for Edit Mode -->
          <div v-if="isEditMode" class="col-12 q-pt-sm">
            <div class="text-caption text-grey-5">
              <q-icon name="info" size="xs" class="q-mr-xs" />
              Kosongkan kolom password jika tidak ingin mengubah password yang sudah ada.
            </div>
          </div>

          <!-- Password -->
          <div class="col-12 col-md-6">
            <label class="form-label">
              Password {{ isEditMode ? '(Baru)' : '' }} <span v-if="!isEditMode" class="req">*</span>
            </label>
            <q-input
              v-model="form.password"
              :type="showPassword ? 'text' : 'password'"
              outlined
              dense
              class="clay-input"
              placeholder="Minimal 6 karakter"
              :rules="isEditMode ? [optionalMin6] : [required('Password wajib diisi'), minLength(6, 'Minimal 6 karakter')]"
              lazy-rules
            >
              <template #append>
                <q-icon
                  :name="showPassword ? 'visibility_off' : 'visibility'"
                  class="cursor-pointer"
                  @click="showPassword = !showPassword"
                />
              </template>
            </q-input>
          </div>

          <!-- Konfirmasi Password -->
          <div class="col-12 col-md-6">
            <label class="form-label">
              Konfirmasi Password <span v-if="!isEditMode" class="req">*</span>
            </label>
            <q-input
              v-model="form.password_confirmation"
              :type="showPassword ? 'text' : 'password'"
              outlined
              dense
              class="clay-input"
              placeholder="Ulangi password"
              :rules="[matchPasswordRule]"
              lazy-rules
            />
          </div>
        </div>
      </q-form>
    </AppFormDialog>
  </q-page>
</template>

<script setup>
import { ref, onMounted } from 'vue'
import { useUserStore } from 'src/stores/user.store'
import { useAuthStore } from 'src/stores/auth.store'
import { useNotify } from 'src/composables/useNotify'
import { useConfirmDelete } from 'src/composables/useConfirmDelete'
import { required, minLength } from 'src/utils/validators'
import { formatDate } from 'src/utils/formatter'

import AppPageHeader from 'components/shared/AppPageHeader.vue'
import AppDataTable from 'components/shared/AppDataTable.vue'
import AppFormDialog from 'components/shared/AppFormDialog.vue'
import AppBadge from 'components/shared/AppBadge.vue'

const userStore = useUserStore()
const authStore = useAuthStore()
const { notifySuccess, notifyError } = useNotify()
const { confirmDelete } = useConfirmDelete()

const searchQuery = ref('')
const selectedRole = ref(null)
const isDialogOpen = ref(false)
const isEditMode = ref(false)
const showPassword = ref(false)
const formRef = ref(null)

const roleOptions = [
  { label: 'Super Admin', value: 'superadmin' },
  { label: 'Inventaris', value: 'inventory' },
  { label: 'Finance', value: 'finance' },
  { label: 'Human Resource', value: 'hr' },
  { label: 'Supply Chain', value: 'supply_chain' }
]

const columns = [
  { name: 'name', label: 'Nama Pengguna', field: 'name', align: 'left', sortable: true },
  { name: 'email', label: 'Email', field: 'email', align: 'left', sortable: true },
  { name: 'role', label: 'Role / Akses', field: 'role', align: 'center', sortable: true },
  { name: 'created_at', label: 'Didaftarkan', field: 'created_at', align: 'left', sortable: true },
  { name: 'actions', label: 'Aksi', field: 'actions', align: 'center' }
]

const initialForm = {
  id: null,
  name: '',
  email: '',
  role: 'inventory',
  password: '',
  password_confirmation: ''
}

const form = ref({ ...initialForm })

function getUserInitial (name) {
  if (!name) return 'U'
  const parts = name.trim().split(/\s+/)
  if (parts.length >= 2) return (parts[0][0] + parts[1][0]).toUpperCase()
  return name.slice(0, 2).toUpperCase()
}

function isCurrentLoggedInUser (row) {
  if (!authStore.currentUser) return false
  return authStore.currentUser.id === row.id || authStore.currentUser.email === row.email
}

function getRoleVariant (role) {
  const r = String(role || '').toLowerCase()
  if (r === 'superadmin') return 'primary'
  if (r === 'inventory') return 'secondary'
  if (r === 'finance') return 'warning'
  if (r === 'hr') return 'info'
  if (r === 'supply_chain' || r === 'sc') return 'accent'
  if (r === 'manager') return 'info'
  if (r === 'barista' || r === 'staff') return 'positive'
  return 'neutral'
}

function formatRole (role) {
  const r = String(role || '').toLowerCase()
  const map = {
    superadmin: 'Super Admin',
    inventory: 'Inventaris',
    finance: 'Finance',
    hr: 'Human Resource',
    supply_chain: 'Supply Chain',
    sc: 'Supply Chain',
    manager: 'Manager',
    barista: 'Barista',
    staff: 'Staff'
  }
  return map[r] || role || 'User'
}

function optionalMin6 (val) {
  if (!val) return true
  return val.length >= 6 || 'Password minimal 6 karakter'
}

function matchPasswordRule (val) {
  if (!form.value.password && !val) return true
  if (form.value.password && !val) return 'Konfirmasi password wajib diisi'
  return val === form.value.password || 'Konfirmasi password tidak cocok'
}

async function loadData () {
  try {
    await userStore.fetchUsers({
      search: searchQuery.value || undefined,
      role: selectedRole.value || undefined
    })
  } catch (err) {
    notifyError(err.response?.data?.message || err.message || 'Gagal memuat pengguna')
  }
}

function onFilter () {
  loadData()
}

function resetFilter () {
  searchQuery.value = ''
  selectedRole.value = null
  loadData()
}

function openDialog (row = null) {
  isEditMode.value = !!row
  showPassword.value = false
  if (row) {
    form.value = {
      id: row.id,
      name: row.name,
      email: row.email,
      role: row.role || 'staff',
      password: '',
      password_confirmation: ''
    }
  } else {
    form.value = { ...initialForm }
  }
  isDialogOpen.value = true
}

async function saveData () {
  const valid = await formRef.value?.validate()
  if (!valid) return

  const payload = {
    name: form.value.name,
    email: form.value.email,
    role: form.value.role
  }

  if (form.value.password) {
    payload.password = form.value.password
  }

  try {
    if (isEditMode.value) {
      await userStore.updateUser(form.value.id, payload)
      notifySuccess('Data pengguna berhasil diperbarui')
    } else {
      await userStore.createUser(payload)
      notifySuccess('Pengguna baru berhasil ditambahkan')
    }
    isDialogOpen.value = false
  } catch (err) {
    notifyError(err.response?.data?.message || err.message || 'Gagal menyimpan pengguna')
  }
}

function handleDelete (row) {
  if (isCurrentLoggedInUser(row)) {
    notifyError('Anda tidak dapat menghapus akun Anda sendiri')
    return
  }

  confirmDelete(`pengguna "${row.name}"`, async () => {
    try {
      await userStore.deleteUser(row.id)
      notifySuccess(`Pengguna ${row.name} berhasil dihapus`)
    } catch (err) {
      notifyError(err.response?.data?.message || err.message || 'Gagal menghapus pengguna')
    }
  })
}

onMounted(() => {
  loadData()
})
</script>

<style lang="scss" scoped>
@import 'src/css/finance-page';
</style>

<template>
  <q-dialog v-model="isOpen" persistent>
    <q-card style="min-width: 400px; max-width: 600px; width: 100%">
      <q-card-section class="row items-center q-pb-none">
        <div class="text-h6">{{ isEdit ? 'Edit Produk' : 'Tambah Produk' }}</div>
        <q-space />
        <q-btn icon="close" flat round dense v-close-popup />
      </q-card-section>

      <q-card-section class="q-pt-md">
        <q-form greedy @submit.prevent="onSubmit" ref="formRef">
          <q-input
            v-model="form.sku"
            label="SKU"
            outlined
            dense
            class="q-mb-md"
            :rules="[(v) => !!v || 'SKU wajib diisi']"
            lazy-rules
          />

          <q-input
            v-model="form.name"
            label="Nama Produk"
            outlined
            dense
            class="q-mb-md"
            :rules="[(v) => !!v || 'Nama wajib diisi']"
            lazy-rules
          />

          <q-input
            v-model.number="form.stock"
            label="Stok Awal"
            type="number"
            outlined
            dense
            class="q-mb-md"
            :rules="[(v) => v >= 0 || 'Stok tidak boleh negatif']"
            lazy-rules
          />

          <q-input
            v-model.number="form.price"
            label="Harga"
            type="number"
            outlined
            dense
            class="q-mb-md"
            :rules="[(v) => v >= 0 || 'Harga tidak boleh negatif']"
            lazy-rules
          />

          <q-input
            v-model="form.description"
            label="Deskripsi"
            type="textarea"
            outlined
            dense
            class="q-mb-lg"
            autogrow
          />

          <div class="row justify-end q-gutter-sm">
            <q-btn label="Batal" flat v-close-popup />
            <q-btn
              type="submit"
              :label="isEdit ? 'Simpan' : 'Tambah'"
              color="primary"
              unelevated
              :loading="loading"
            />
          </div>
        </q-form>
      </q-card-section>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref, reactive, computed } from 'vue'
import { useInventoryStore } from 'src/stores/inventory.store'
import { useNotify } from 'src/composables/useNotify'

const props = defineProps({
  modelValue: { type: Boolean, default: false },
  editId: { type: [String, Number], default: null },
})

const emit = defineEmits(['update:modelValue', 'saved'])

const inventoryStore = useInventoryStore()
const { notifySuccess, notifyError } = useNotify()

const isOpen = computed({
  get: () => props.modelValue,
  set: (val) => emit('update:modelValue', val),
})

const isEdit = computed(() => !!props.editId)
const loading = ref(false)
const formRef = ref(null)

const form = reactive({ sku: '', name: '', stock: 0, price: 0, description: '' })

async function loadData() {
  if (isEdit.value) {
    loading.value = true
    try {
      const product = await inventoryStore.fetchProductDetail(props.editId)
      Object.assign(form, {
        sku: product.sku || '',
        name: product.name || '',
        stock: product.stock || 0,
        price: product.price || 0,
        description: product.description || '',
      })
    } catch {
      notifyError('Gagal memuat data produk.')
      isOpen.value = false
    } finally {
      loading.value = false
    }
  } else {
    // Reset form
    Object.assign(form, { sku: '', name: '', stock: 0, price: 0, description: '' })
  }
}

// Expose a method to be called from the parent when the dialog opens
defineExpose({ loadData })

async function onSubmit() {
  loading.value = true
  try {
    if (isEdit.value) {
      await inventoryStore.updateProduct(props.editId, { ...form })
      notifySuccess('Produk berhasil diperbarui.')
    } else {
      await inventoryStore.createProduct({ ...form })
      notifySuccess('Produk berhasil ditambahkan.')
    }
    emit('saved')
    isOpen.value = false
  } catch {
    notifyError('Gagal menyimpan produk. Coba lagi.')
  } finally {
    loading.value = false
  }
}
</script>

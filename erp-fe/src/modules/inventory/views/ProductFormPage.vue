<template>
  <q-page class="q-pa-md">
    <div class="row items-center q-mb-md">
      <q-btn flat dense round icon="arrow_back" class="q-mr-sm" aria-label="Kembali" @click="router.back()" />
      <div class="text-h6">{{ isEdit ? 'Edit Produk' : 'Tambah Produk' }}</div>
    </div>

    <q-card flat bordered class="q-pa-md" style="max-width: 600px">
      <q-form greedy @submit.prevent="onSubmit">

        <q-input v-model="form.sku" label="SKU" outlined dense class="q-mb-md"
          :rules="[v => !!v || 'SKU wajib diisi']" lazy-rules />

        <q-input v-model="form.name" label="Nama Produk" outlined dense class="q-mb-md"
          :rules="[v => !!v || 'Nama wajib diisi']" lazy-rules />

        <q-input v-model.number="form.stock" label="Stok Awal" type="number" outlined dense class="q-mb-md"
          :rules="[v => v >= 0 || 'Stok tidak boleh negatif']" lazy-rules />

        <q-input v-model.number="form.price" label="Harga" type="number" outlined dense class="q-mb-md"
          :rules="[v => v >= 0 || 'Harga tidak boleh negatif']" lazy-rules />

        <q-input v-model="form.description" label="Deskripsi" type="textarea" outlined dense class="q-mb-lg" autogrow />

        <div class="row q-gutter-sm">
          <q-btn type="submit" :label="isEdit ? 'Simpan Perubahan' : 'Tambah Produk'"
            color="primary" unelevated :loading="loading" />
          <q-btn label="Batal" flat @click="router.back()" />
        </div>

      </q-form>
    </q-card>
  </q-page>
</template>

<script setup>
import { ref, reactive, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useInventoryStore } from 'src/stores/inventory.store'
import { useNotify } from 'src/composables/useNotify'

const router         = useRouter()
const route          = useRoute()
const inventoryStore = useInventoryStore()
const { notifySuccess, notifyError } = useNotify()

const loading = ref(false)
const form    = reactive({ sku: '', name: '', stock: 0, price: 0, description: '' })

const isEdit = computed(() => !!route.params.id)

onMounted(async () => {
  if (isEdit.value) {
    try {
      const product = await inventoryStore.fetchProductDetail(route.params.id)
      Object.assign(form, {
        sku:         product.sku         || '',
        name:        product.name        || '',
        stock:       product.stock       || 0,
        price:       product.price       || 0,
        description: product.description || ''
      })
    } catch {
      notifyError('Gagal memuat data produk.')
      router.back()
    }
  }
})

async function onSubmit () {
  loading.value = true
  try {
    if (isEdit.value) {
      await inventoryStore.updateProduct(route.params.id, { ...form })
      notifySuccess('Produk berhasil diperbarui.')
    } else {
      await inventoryStore.createProduct({ ...form })
      notifySuccess('Produk berhasil ditambahkan.')
    }
    router.push({ name: 'inventory-product-list' })
  } catch {
    notifyError('Gagal menyimpan produk. Coba lagi.')
  } finally {
    loading.value = false
  }
}
</script>

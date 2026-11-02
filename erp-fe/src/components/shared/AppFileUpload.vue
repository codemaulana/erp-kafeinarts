<template>
  <div
    class="file-upload-area"
    :class="{ 'file-upload-area--has-file': !!fileName }"
    role="button"
    tabindex="0"
    :aria-label="fileName ? `File terpilih: ${fileName}` : 'Klik atau drag & drop untuk upload file'"
    @click="triggerInput"
    @keydown.enter.prevent="triggerInput"
    @keydown.space.prevent="triggerInput"
    @dragover.prevent
    @drop.prevent="onDrop"
  >
    <input
      ref="inputRef"
      type="file"
      :accept="accept"
      style="display:none"
      @change="onFileChange"
    />

    <!-- Empty state -->
    <template v-if="!fileName">
      <q-icon name="cloud_upload" size="28px" color="grey-4" />
      <div class="file-upload-area__text">Klik atau drag &amp; drop file di sini</div>
      <div class="file-upload-area__hint">{{ hint }}</div>
    </template>

    <!-- File selected state -->
    <template v-else>
      <q-icon name="attach_file" size="24px" :color="iconColor" />
      <div class="file-upload-area__text" :style="{ color: textColor }">{{ fileName }}</div>
      <q-btn
        flat dense no-caps size="xs"
        label="Hapus"
        color="red-5"
        @click.stop="removeFile"
      />
    </template>
  </div>
</template>

<script setup>
import { ref } from 'vue'
import { useQuasar } from 'quasar'

const props = defineProps({
  /** Nama file yang sedang dipilih (untuk two-way binding tampilan) */
  fileName:  { type: String, default: '' },
  accept:    { type: String, default: 'image/*,.pdf' },
  /** Ukuran maksimum dalam bytes. Default: 5MB */
  maxSize:   { type: Number, default: 5 * 1024 * 1024 },
  hint:      { type: String, default: 'PNG, JPG, PDF — maks 5MB' },
  iconColor: { type: String, default: 'blue-6' },
  textColor: { type: String, default: 'var(--clay-accent)' }
})

const emit = defineEmits([
  /** Emits File object saat file dipilih */
  'file-selected',
  /** Emits null saat file dihapus */
  'file-removed'
])

const $q = useQuasar()
const inputRef = ref(null)

function triggerInput () {
  inputRef.value?.click()
}

function processFile (file) {
  if (!file) return
  if (file.size > props.maxSize) {
    const mb = Math.round(props.maxSize / (1024 * 1024))
    $q.notify({ type: 'warning', message: `Ukuran file maks ${mb}MB`, position: 'top' })
    return
  }
  emit('file-selected', file)
}

function onFileChange (e) {
  processFile(e.target.files?.[0])
  // Reset input value agar file yang sama bisa dipilih ulang
  if (inputRef.value) inputRef.value.value = ''
}

function onDrop (e) {
  processFile(e.dataTransfer?.files?.[0])
}

function removeFile () {
  emit('file-removed')
}
</script>

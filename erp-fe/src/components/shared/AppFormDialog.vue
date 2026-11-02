<template>
  <q-dialog
    :model-value="modelValue"
    persistent
    :maximized="maximized"
    transition-show="scale"
    transition-hide="scale"
    @update:model-value="$emit('update:modelValue', $event)"
  >
    <q-card
      :style="!maximized ? `min-width:${minWidth}; max-width:${width || maxWidth}; width:${width || '100%'};` : undefined"
      class="clay-dialog__card"
      :class="{ 'clay-dialog__card--full': maximized }"
    >
      <!-- Header -->
      <div class="clay-dialog__header">
        <div class="clay-dialog__header-title">
          <q-icon :name="icon" size="20px" :color="iconColor" />
          {{ title }}
        </div>
        <q-btn flat round dense icon="close" size="sm" @click="$emit('update:modelValue', false)" />
      </div>

      <!-- Body: wraps slot content in q-form -->
      <q-form
        ref="formRef"
        greedy
        class="clay-dialog__body"
        :class="{ 'clay-dialog__body--full': maximized }"
        @submit.prevent="$emit('submit')"
      >
        <slot />

        <!-- Footer -->
        <div class="clay-dialog__footer">
          <slot name="footer">
            <q-btn
              flat no-caps
              label="Batal"
              class="clay-btn"
              @click="$emit('update:modelValue', false)"
            />
            <q-btn
              unelevated no-caps
              :label="submitLabel"
              type="submit"
              :loading="loading"
              class="clay-btn clay-btn-primary"
            />
          </slot>
        </div>
      </q-form>
    </q-card>
  </q-dialog>
</template>

<script setup>
import { ref } from 'vue'

defineProps({
  modelValue:  { type: Boolean, required: true },
  title:       { type: String,  default: 'Form' },
  icon:        { type: String,  default: 'edit' },
  iconColor:   { type: String,  default: 'blue-6' },
  submitLabel: { type: String,  default: 'Simpan' },
  loading:     { type: Boolean, default: false },
  maximized:   { type: Boolean, default: false },
  minWidth:    { type: String,  default: '420px' },
  maxWidth:    { type: String,  default: '560px' },
  width:       { type: String,  default: '' }
})

defineEmits(['update:modelValue', 'submit'])

const formRef = ref(null)

/** Expose validate() agar parent bisa memanggil formRef.value.validate() via template ref */
async function validate () {
  return formRef.value?.validate()
}

defineExpose({ validate })
</script>

<style lang="scss">
.clay-dialog__card {
  border-radius: var(--rb-radius-xl, 16px) !important;
  box-shadow: var(--rb-shadow-xl, 0 20px 25px -5px rgba(0, 0, 0, 0.5)) !important;
  background: var(--rb-surface, #1e293b) !important;
  border: 1px solid var(--rb-border, rgba(255, 255, 255, 0.1)) !important;
  max-height: 90vh !important;
  display: flex !important;
  flex-direction: column !important;
  overflow: hidden !important;

  &--full {
    border-radius: 0 !important;
    max-height: 100vh !important;
    height: 100vh !important;
  }
}

.clay-dialog__header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 16px 22px;
  background: var(--rb-surface-alt, #0f172a);
  border-bottom: 1px solid var(--rb-border, rgba(255, 255, 255, 0.08));
  flex-shrink: 0;

  &-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 15px;
    font-weight: 700;
    color: var(--rb-text-primary, #ffffff);
  }
}

.clay-dialog__body {
  flex: 1 1 auto;
  overflow-y: auto;
  min-height: 0;
  padding: 20px 22px;
  display: flex;
  flex-direction: column;
}

.clay-dialog__footer {
  display: flex;
  justify-content: flex-end;
  gap: 8px;
  padding: 14px 22px;
  margin-top: 16px;
  border-top: 1px solid var(--rb-border, rgba(255, 255, 255, 0.08));
  background: var(--rb-surface-alt, #0f172a);
  flex-shrink: 0;
}
</style>

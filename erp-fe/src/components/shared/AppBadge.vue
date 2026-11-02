<template>
  <span :class="['clay-badge', variantClass]">
    <q-icon v-if="icon" :name="icon" size="11px" />
    <slot>{{ label }}</slot>
  </span>
</template>

<script setup>
import { computed } from 'vue'

/**
 * AppBadge — badge/chip dengan variant semantik.
 *
 * Variant options:
 *   success  → clay-badge--success  (hijau)
 *   danger   → clay-badge--danger   (merah)
 *   warning  → clay-badge--warning  (kuning)
 *   info     → clay-badge--info     (biru)
 *   neutral  → clay-badge--neutral  (abu)
 */
const props = defineProps({
  variant: {
    type: String,
    default: 'neutral',
    validator: (v) => ['success', 'danger', 'warning', 'info', 'neutral'].includes(v)
  },
  label: { type: String, default: '' },
  icon:  { type: String, default: '' }
})

const variantClass = computed(() => `clay-badge--${props.variant}`)
</script>

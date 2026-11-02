<template>
  <div class="clay-stat-card">
    <div class="clay-stat-card__icon" :style="iconStyle">
      <q-icon :name="icon" :size="iconSize" />
    </div>
    <div class="clay-stat-card__content">
      <div class="clay-stat-card__label">{{ displayLabel }}</div>
      <div class="clay-stat-card__value" :style="valueStyle">{{ value }}</div>
      <div v-if="sub" class="clay-stat-card__sub">{{ sub }}</div>
    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const COLOR_MAP = {
  primary:  { color: '#3B82F6', bg: 'linear-gradient(135deg, #EFF6FF, #DBEAFE)' },
  positive: { color: '#10B981', bg: 'linear-gradient(135deg, #ECFDF5, #D1FAE5)' },
  success:  { color: '#10B981', bg: 'linear-gradient(135deg, #ECFDF5, #D1FAE5)' },
  warning:  { color: '#F59E0B', bg: 'linear-gradient(135deg, #FFFBEB, #FEF3C7)' },
  negative: { color: '#EF4444', bg: 'linear-gradient(135deg, #FEF2F2, #FEE2E2)' },
  danger:   { color: '#EF4444', bg: 'linear-gradient(135deg, #FEF2F2, #FEE2E2)' },
  info:     { color: '#06B6D4', bg: 'linear-gradient(135deg, #ECFEFF, #CFFAFE)' },
  purple:   { color: '#8B5CF6', bg: 'linear-gradient(135deg, #F5F3FF, #EDE9FE)' }
}

const props = defineProps({
  icon:       { type: String, required: true },
  iconSize:   { type: String, default: '22px' },
  iconBg:     { type: String, default: '' },
  iconColor:  { type: String, default: '' },
  color:      { type: String, default: '' },
  label:      { type: String, default: '' },
  title:      { type: String, default: '' },
  value:      { type: [String, Number], required: true },
  valueColor: { type: String, default: '' },
  sub:        { type: String, default: '' }
})

const displayLabel = computed(() => props.label || props.title || '')

const resolvedColors = computed(() => {
  if (props.color && COLOR_MAP[props.color]) {
    return COLOR_MAP[props.color]
  }
  return {
    color: props.iconColor || '#3B82F6',
    bg: props.iconBg || 'linear-gradient(135deg, #EFF6FF, #DBEAFE)'
  }
})

const iconStyle = computed(() => ({
  background: resolvedColors.value.bg,
  color: resolvedColors.value.color
}))

const valueStyle = computed(() => {
  if (props.valueColor) return { color: props.valueColor }
  return {}
})
</script>

<style lang="scss" scoped>
.clay-stat-card {
  transition: all 0.25s cubic-bezier(0.4, 0, 0.2, 1);

  &__icon {
    transition: all 0.25s ease;

    body.body--dark & {
      background: color-mix(in srgb, currentColor 14%, rgba(30, 41, 59, 0.75)) !important;
      border: 1px solid color-mix(in srgb, currentColor 28%, transparent) !important;
      box-shadow: 0 4px 16px -2px color-mix(in srgb, currentColor 25%, transparent) !important;
    }
  }

  &__value {
    body.body--dark & {
      color: #F8FAFC;
    }
  }
}
</style>

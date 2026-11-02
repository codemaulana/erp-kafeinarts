<template>
  <nav class="clay-breadcrumb" aria-label="Breadcrumb">
    <template v-for="(item, idx) in items" :key="idx">
      <span
        :class="['clay-breadcrumb__item', { 'clay-breadcrumb__item--active': idx === items.length - 1 }]"
        :role="idx < items.length - 1 ? 'button' : undefined"
        :tabindex="idx < items.length - 1 ? 0 : undefined"
        @click="idx < items.length - 1 && $emit('navigate', idx)"
        @keydown.enter="idx < items.length - 1 && $emit('navigate', idx)"
      >{{ item }}</span>
      <span v-if="idx < items.length - 1" class="clay-breadcrumb__sep" aria-hidden="true">›</span>
    </template>
  </nav>
</template>

<script setup>
defineProps({
  items: { type: Array, default: () => [] }
})

defineEmits(['navigate'])
</script>

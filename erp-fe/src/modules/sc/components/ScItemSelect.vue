<template>
  <q-select
    :model-value="modelValue"
    @update:model-value="$emit('update:modelValue', $event)"
    :options="scStore.itemOptions"
    outlined
    dense
    class="clay-input"
    emit-value
    map-options
    :label="label"
    :rules="rules"
    lazy-rules
    :readonly="readonly"
    :clearable="clearable"
    :disable="disable"
  />
</template>

<script setup>
import { onMounted } from 'vue'
import { useScStore } from 'src/stores/sc.store'

defineProps({
  modelValue: {
    type: [String, Number],
    default: null,
  },
  label: {
    type: String,
    default: undefined,
  },
  rules: {
    type: Array,
    default: () => [],
  },
  readonly: {
    type: Boolean,
    default: false,
  },
  clearable: {
    type: Boolean,
    default: false,
  },
  disable: {
    type: Boolean,
    default: false,
  },
})

defineEmits(['update:modelValue'])

const scStore = useScStore()

onMounted(async () => {
  if (scStore.itemOptions.length === 0 && !scStore.isLoading('item')) {
    await scStore.fetchItem({ rowsPerPage: 100 })
  }
})
</script>

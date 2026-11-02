<template>
  <!--
    FinancePageHeader — thin alias for AppPageHeader.
    Kept for backward-compatibility so existing finance pages
    don't need template changes. All new pages should use
    AppPageHeader directly.
  -->
  <AppPageHeader
    v-bind="$attrs"
    :title="title"
    :subtitle="subtitle"
    :icon="icon"
    :icon-bg="iconBg"
    :icon-color="iconColor"
  >
    <template v-for="(_, name) in $slots" #[name]="slotProps">
      <slot :name="name" v-bind="slotProps ?? {}" />
    </template>
  </AppPageHeader>
</template>

<script setup>
import AppPageHeader from './AppPageHeader.vue'

defineOptions({ inheritAttrs: false })

defineProps({
  title:     { type: String, required: true },
  subtitle:  { type: String, default: '' },
  icon:      { type: String, default: 'account_balance_wallet' },
  iconBg:    { type: String, default: 'linear-gradient(135deg,#EFF6FF,#DBEAFE)' },
  iconColor: { type: String, default: '#3B82F6' }
})
</script>

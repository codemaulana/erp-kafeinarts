<template>
  <!--
    AppPageHeader
    Displays icon + title + subtitle + #actions slot.
    Breadcrumb is now rendered in AppHeader (the top bar), not here.
  -->
  <div class="fph clay-animate-up">
    <div class="fph__row">

      <div class="fph__left">
        <div class="fph__icon-wrap" :style="iconWrapStyle">
          <slot name="icon">
            <q-icon :name="icon" size="22px" />
          </slot>
        </div>
        <div class="fph__text">
          <h1 class="fph__title">{{ title }}</h1>
          <p v-if="subtitle" class="fph__sub">{{ subtitle }}</p>
        </div>
      </div>

      <div class="fph__right">
        <slot name="actions" />
      </div>

    </div>
  </div>
</template>

<script setup>
import { computed } from 'vue'

const props = defineProps({
  title:     { type: String, required: true },
  subtitle:  { type: String, default: '' },
  icon:      { type: String, default: 'account_balance_wallet' },
  iconBg:    { type: String, default: 'linear-gradient(135deg,#EFF6FF,#DBEAFE)' },
  iconColor: { type: String, default: '#3B82F6' }
})

const iconWrapStyle = computed(() => ({
  background: props.iconBg,
  color:      props.iconColor
}))
</script>

<style lang="scss" scoped>
.fph {
  margin-bottom: 24px;
}

.fph__row {
  display: flex;
  align-items: center;
  justify-content: space-between;
  gap: 16px;
  flex-wrap: wrap;
}

.fph__left {
  display: flex;
  align-items: center;
  gap: 14px;
  min-width: 0;
}

.fph__icon-wrap {
  width: 52px;
  height: 52px;
  min-width: 52px;
  border-radius: var(--rb-radius-md, 12px);
  display: flex;
  align-items: center;
  justify-content: center;
  flex-shrink: 0;
  box-shadow: var(--rb-shadow-sm);
  transition: all 0.25s ease;

  body.body--dark & {
    background: color-mix(in srgb, currentColor 14%, rgba(30, 41, 59, 0.75)) !important;
    border: 1px solid color-mix(in srgb, currentColor 30%, transparent) !important;
    box-shadow: 0 4px 20px -2px color-mix(in srgb, currentColor 25%, transparent) !important;
  }
}

.fph__text { min-width: 0; }

.fph__title {
  font-size: clamp(18px, 2vw, 22px);
  font-weight: 800;
  color: var(--rb-text-primary, #0F172A);
  margin: 0 0 2px;
  line-height: 1.2;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.fph__sub {
  font-size: 13px;
  color: var(--rb-text-muted, #64748B);
  margin: 0;
  white-space: nowrap;
  overflow: hidden;
  text-overflow: ellipsis;
}

.fph__right {
  display: flex;
  align-items: center;
  gap: 10px;
  flex-wrap: wrap;
  flex-shrink: 0;
}

@media (max-width: 600px) {
  .fph__row {
    flex-direction: column;
    align-items: flex-start;
    gap: 12px;
  }
  .fph__left {
    width: 100%;
  }
  .fph__right {
    width: 100%;
    justify-content: flex-end;
  }
  .fph__title {
    font-size: 18px;
    white-space: normal;
  }
  .fph__sub {
    white-space: normal;
  }
}
</style>

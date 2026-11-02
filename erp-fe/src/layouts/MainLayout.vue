<template>
  <q-layout view="lHh Lpr lFf" class="main-layout">

    <AppHeader :drawer-open="drawerOpen" @toggle-drawer="toggleDrawer" />

    <AppSidebar v-model="drawerOpen" />

    <q-page-container class="main-layout__content">
      <router-view v-slot="{ Component }">
        <transition name="page-fade" mode="out-in">
          <component :is="Component" :key="$route.fullPath" />
        </transition>
      </router-view>
    </q-page-container>

  </q-layout>
</template>

<script setup>
import { ref } from 'vue'
import AppHeader  from 'components/layout/AppHeader.vue'
import AppSidebar from 'components/layout/AppSidebar.vue'

const drawerOpen = ref(true)

function toggleDrawer () {
  drawerOpen.value = !drawerOpen.value
}
</script>

<style lang="scss">
// ── Layout shell ────────────────────────────────────────────────
.main-layout,
.main-layout__content {
  background: var(--rb-bg, #F1F5F9) !important;
  transition: background 0.30s ease, padding-left 0.25s ease;
  overflow-x: hidden;
}

body.body--dark {
  .main-layout,
  .main-layout__content { background: var(--rb-bg, #0F172A) !important; }
}

// ── Page route transition ────────────────────────────────────────
.page-fade-enter-active,
.page-fade-leave-active { transition: opacity 0.18s ease, transform 0.18s ease; }
.page-fade-enter-from   { opacity: 0; transform: translateY(8px);  }
.page-fade-leave-to     { opacity: 0; transform: translateY(-4px); }
</style>

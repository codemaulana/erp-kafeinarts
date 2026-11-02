<template>
  <!-- Minimal unauthenticated layout — no header/sidebar -->
  <q-layout view="lHh lpr lFf">
    <q-page-container>
      <router-view />
    </q-page-container>
  </q-layout>
</template>

<script>
export default { name: 'AuthLayout' }
</script>

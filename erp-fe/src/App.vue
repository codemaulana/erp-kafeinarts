<template>
  <router-view />
</template>

<script setup>
import { onMounted } from 'vue'

// Synchronize favicon with System OS & Browser Theme (Laptop, Tablet, HP)
function syncFavicon() {
  if (typeof window === 'undefined') return
  const isSystemDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches
  const theme = isSystemDark ? 'dark' : 'light'
  // Dark system (dark browser tab) -> White/light emblem for high contrast
  // Light system (light browser tab) -> Dark navy emblem for high contrast
  const iconIco = isSystemDark ? '/images/kafeinarts-logo-light.ico' : '/images/kafeinarts-logo-dark.ico'
  const iconPng = isSystemDark ? '/images/kafeinarts-logo-light.png' : '/images/kafeinarts-logo-dark.png'

  // Clean old links and replace to guarantee browser tab updates
  const existingIcons = document.querySelectorAll("link[rel*='icon']")
  existingIcons.forEach(el => el.remove())

  const linkIco = document.createElement('link')
  linkIco.rel = 'icon'
  linkIco.type = 'image/x-icon'
  linkIco.href = `${iconIco}?t=${theme}`
  document.head.appendChild(linkIco)

  const linkPng = document.createElement('link')
  linkPng.rel = 'icon'
  linkPng.type = 'image/png'
  linkPng.sizes = '32x32'
  linkPng.href = `${iconPng}?t=${theme}`
  document.head.appendChild(linkPng)

  const touchIcon = document.querySelector("link[rel='apple-touch-icon']")
  if (touchIcon) {
    touchIcon.href = iconPng
  }
}

onMounted(() => {
  if (typeof window !== 'undefined' && window.matchMedia) {
    const mq = window.matchMedia('(prefers-color-scheme: dark)')
    syncFavicon()
    try {
      mq.addEventListener('change', syncFavicon)
    } catch {
      mq.addListener(syncFavicon)
    }
  }
})
</script>

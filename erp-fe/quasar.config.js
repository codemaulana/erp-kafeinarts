// Configuration for your app
// https://v2.quasar.dev/quasar-cli-webpack/quasar-config-file

import { defineConfig } from '#q-app/wrappers'
import { fileURLToPath } from 'node:url'
import path from 'node:path'

const __dirname = fileURLToPath(new URL('.', import.meta.url))

export default defineConfig(() => {
  return {
    eslint: {
      warnings: true,
      errors: true,
    },

    // Boot files — order matters: pinia first so auth store is ready for axios interceptors
    boot: [
      'pinia', // 1. Pinia must be first — auth store used by axios + acl
      'axios', // 2. Configures api instance with Bearer token interceptor
      'acl', // 3. Registers $can / $hasRole global helpers
    ],

    // Global CSS
    css: ['app.scss'],

    // Quasar extras (icons & fonts)
    extras: ['roboto-font', 'material-icons'],

    // Build configuration
    build: {
      vueRouterMode: 'history',

      env: {
        // Backend API base — override with .env or CI variable
        API_BASE_URL: process.env.API_BASE_URL || 'http://localhost:8000/api/v1',
        // Backend root URL (used for Sanctum SPA cookie flow)
        APP_URL: process.env.APP_URL || 'http://localhost:8000',
      },

      esbuildTarget: {
        browser: ['es2022', 'firefox115', 'chrome115', 'safari14'],
        node: 'node20',
      },

      // Webpack alias: 'modules/...' → src/modules/...
      chainWebpack(chain) {
        chain.resolve.alias.set('modules', path.join(__dirname, 'src/modules'))
      },
    },

    // Dev server
    devServer: {
      server: { type: 'http' },
      port: 9000,
      open: true,
    },

    // Quasar Framework
    framework: {
      config: {
        dark: false, // Default to light mode, togglable in-app
      },

      // Only register plugins that are actually used:
      //   Notify  — useNotify composable ($q.notify)
      //   Dialog  — useConfirmDelete composable ($q.dialog)
      plugins: ['Notify', 'Dialog'],
    },

    animations: [],

    ssr: {
      prodPort: 3000,
      middlewares: ['render'],
      pwa: false,
    },

    pwa: {
      workboxMode: 'GenerateSW',
    },
  }
})

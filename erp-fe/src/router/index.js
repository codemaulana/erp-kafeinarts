/**
 * router/index.js
 * Router instance — guards menggunakan Pinia auth store.
 * Tidak lagi mengimport Vuex store.
 */
import { createRouter, createWebHistory, createWebHashHistory } from 'vue-router'
import routes from './routes'
import { setupGuards } from './guards'

const createHistory = process.env.VUE_ROUTER_MODE === 'history'
  ? createWebHistory
  : createWebHashHistory

const Router = createRouter({
  scrollBehavior: () => ({ left: 0, top: 0 }),
  routes,
  history: createHistory(process.env.VUE_ROUTER_BASE)
})

setupGuards(Router)

export default Router

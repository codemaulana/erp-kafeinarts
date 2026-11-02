import authRoutes from '../modules/auth/routes'
import financeRoutes from '../modules/finance/routes'
import hrRoutes from '../modules/hr/routes'
import scRoutes from '../modules/sc/routes'
import inventoryRoutes from '../modules/inventory/routes'

const routes = [
  // ── Redirect root ─────────────────────────────────────────
  { path: '/', redirect: '/login' },

  // ── Auth routes (no sidebar/header) ───────────────────────
  {
    path: '/',
    component: () => import('layouts/AuthLayout.vue'),
    children: authRoutes,
  },

  // ── App routes (with sidebar/header) ──────────────────────
  {
    path: '/app',
    component: () => import('layouts/MainLayout.vue'),
    meta: { requiresAuth: true },
    children: [
      { path: '', redirect: { name: 'dashboard' } },

      // Dashboard
      {
        path: 'dashboard',
        name: 'dashboard',
        component: () => import('pages/DashboardPage.vue'),
        meta: {
          requiresAuth: true,
          breadcrumb: [{ icon: 'dashboard', label: 'Dashboard' }],
        },
      },

      // Inventory routes
      ...inventoryRoutes,

      // Finance routes
      ...financeRoutes,

      // HR routes
      ...hrRoutes,

      // SC routes
      ...scRoutes,

      // Settings / User Management
      {
        path: 'settings/users',
        name: 'user-management',
        component: () => import('pages/settings/UserManagementPage.vue'),
        meta: {
          requiresAuth: true,
          breadcrumb: [
            { icon: 'settings', label: 'Pengaturan' },
            { icon: 'manage_accounts', label: 'Manajemen Pengguna' },
          ],
        },
      },
    ],
  },

  // ── 404 ───────────────────────────────────────────────────
  {
    path: '/:catchAll(.*)*',
    component: () => import('pages/ErrorNotFound.vue'),
  },
]

export default routes

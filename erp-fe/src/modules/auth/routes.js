export default [
  {
    path: '/login',
    name: 'login',
    component: () => import('modules/auth/views/LoginPage.vue'),
    meta: { requiresAuth: false }
  }
]
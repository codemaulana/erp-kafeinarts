export default [
  {
    path: '/inventory/products',
    name: 'inventory-products',
    component: () => import('modules/inventory/views/ProductListPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'category', label: 'Inventaris' },
        { label: 'Daftar Produk' },
      ],
    },
  },
  {
    path: '/inventory/products/create',
    name: 'inventory-product-create',
    component: () => import('modules/inventory/views/ProductFormPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'category', label: 'Inventaris' },
        { label: 'Daftar Produk', to: { name: 'inventory-products' } },
        { label: 'Tambah Produk' },
      ],
    },
  },
  {
    path: '/inventory/products/:id/edit',
    name: 'inventory-product-edit',
    component: () => import('modules/inventory/views/ProductFormPage.vue'),
    meta: {
      requiresAuth: true,
      breadcrumb: [
        { icon: 'category', label: 'Inventaris' },
        { label: 'Daftar Produk', to: { name: 'inventory-products' } },
        { label: 'Edit Produk' },
      ],
    },
  },
]

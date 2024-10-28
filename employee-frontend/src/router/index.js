import { createRouter, createWebHistory } from 'vue-router'
import { useAuthStore } from '../stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/',
      name: 'login',
      component: () => import('../views/LoginView.vue'),
      meta: { requiresAuth: false }
    },
    {
      path: '/employees',
      name: 'employees',
      component: () => import('../views/EmployeesView.vue'),
      meta: { requiresAuth: true }
    },
    {
      path: '/positions',
      name: 'positions',
      component: () => import('../views/PositionsView.vue'),
      meta: { requiresAuth: true }
    }
  ]
})


router.beforeEach((to, from, next) => {
  const authStore = useAuthStore()
  
  if (to.meta.requiresAuth && !authStore.token) {
    next('/')
  } else if (to.path === '/' && authStore.token) {
    next('/employees')
  } else {
    next()
  }
})

export default router
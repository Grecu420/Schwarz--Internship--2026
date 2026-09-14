import { createRouter, createWebHistory } from 'vue-router'
import RegisterView from '@/views/RegisterView.vue'
import LoginView from '@/views/LoginView.vue'
import MainView from '@/views/MainView.vue'
import ProfileView from '@/views/ProfileView.vue'
import { useAuthStore } from '@/stores/auth'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { hideHeaderFooter: true }
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { hideHeaderFooter: true }
    },
    {
      path: '/main',
      name: 'main',
      component: MainView,
      meta: { requiresAuth: true }
    },
    {
      path: '/',
      redirect: '/main'
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
    },
  ]
})

// Navigation Guard
router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  } 

  if ((to.name === 'login' || to.name === 'register') && authStore.isAuthenticated) {
    return { name: 'main' }
  }
})

export default router
import { createRouter, createWebHistory } from 'vue-router'
import RegisterView from '@/views/RegisterView.vue'
import LoginView from '@/views/LoginView.vue'
import MainView from '@/views/MainView.vue'
import PropertyCreationView from '@/views/property/PropertyCreationView.vue'
import ProfileView from '@/views/ProfileView.vue'
import ConversationsView from '@/views/ConversationsView.vue'
import { useAuthStore } from '@/stores/auth'
import PropertyDisplayView from '@/views/property/PropertyDisplayView.vue'
import PropertyEditView from '@/views/property/PropertyEditView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/register',
      name: 'register',
      component: RegisterView,
      meta: { hideHeaderFooter: true },
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView,
      meta: { hideHeaderFooter: true },
    },
    {
      path: '/main',
      name: 'main',
      component: MainView,
      meta: { requiresAuth: true },
    },
    {
      path: '/properties/create',
      name: 'createProperty',
      component: PropertyCreationView,
      meta: { requiresAuth: true }, // Protect this route
    },
    {
      path: '/properties/edit/:id',
      name: 'editProperty',
      component: PropertyEditView,
      meta: { requiresAuth: true }, // Protect this route
    },
    {
      path: '/properties/view',
      name: 'viewProperty',
      component: PropertyDisplayView
    },
    {
      path: '/',
      redirect: '/main',
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
    },
    {
      path: '/conversations',
      name: 'conversations',
      component: ConversationsView,
    },
  ]
})

// Navigation Guard
router.beforeEach((to) => {
  const authStore = useAuthStore()

  if (authStore.checkTokenExpiration()) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if (to.meta.requiresAuth && !authStore.isAuthenticated) {
    return { name: 'login', query: { redirect: to.fullPath } }
  }

  if ((to.name === 'login' || to.name === 'register') && authStore.isAuthenticated) {
    return { name: 'main' }
  }
})

export default router

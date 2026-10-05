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
import PropertyManagementView from '@/views/property/PropertyManagementView.vue'
import ReservationsView from '@/views/ReservationsView.vue'
import NotFoundView from '@/views/NotFoundView.vue'

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
    },
    {
      path: '/properties/create',
      name: 'createProperty',
      component: PropertyCreationView,
      meta: { requiresAuth: true },
    },
    {
      path: '/properties/edit/:id',
      name: 'editProperty',
      component: PropertyEditView,
      meta: { requiresAuth: true },
    },
    {
      path: '/properties/view/:id',
      name: 'viewProperty',
      component: PropertyDisplayView,
    },
    {
      path: '/properties',
      name: 'listProperties',
      component: PropertyManagementView,
      meta: { requiresAuth: true },
    },
    {
      path: '/',
      redirect: '/main',
    },
    {
      path: '/profile',
      name: 'profile',
      component: ProfileView,
      meta: { requiresAuth: true },
    },
    {
      path: '/conversations',
      name: 'conversations',
      component: ConversationsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/reservations',
      name: 'reservations',
      component: ReservationsView,
      meta: { requiresAuth: true },
    },
    {
      path: '/404',
      name: 'not-found',
      component: NotFoundView,
    },
    {
      path: '/:pathMatch(.*)*',
      name: 'not-found-wildcard', 
      component: NotFoundView,
    },
  ],
  scrollBehavior(to, from, savedPosition) {
    if (savedPosition) {
      return savedPosition
    } else {
      return { top: 0, behavior: 'smooth' }
    }
  },
})

router.beforeEach((to, from) => {
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

  return true
})

export default router
import { createRouter, createWebHistory } from 'vue-router'
import RegisterView from '@/views/RegisterView.vue'
import LoginView from '@/views/LoginView.vue'
import MainView from '@/views/MainView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/login',
      name: 'login',
      component: LoginView
    },
    {
      path: '/main',
      name: 'main',
      component: MainView,
      meta: { requiresAuth: true } // Protect this route
    },
    {
      path: '/',
      redirect: '/main'
    }
  ]
})

// Navigation Guard runs before every route transition
router.beforeEach((to, _from, next) => {
  

  const token = localStorage.getItem("jwt_token")

  if (to.meta.requiresAuth && token === null) {
    // Redirect unauthenticated users to login
    next({ name: 'login' });
  } else if (to.name === 'login' && token !== null) {
    // Prevent logged-in users from returning to login page
    next({ name: 'main' });
  } else {
    next();
  }
});

export default router
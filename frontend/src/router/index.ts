import { createRouter, createWebHistory } from 'vue-router'
import RegisterView from '../views/RegisterView.vue'

const router = createRouter({
  history: createWebHistory(import.meta.env.BASE_URL),
  routes: [
    {
      path: '/register',
      name: 'register',
      component: RegisterView
    },
    {
      path: '/',
      redirect: '/register'
    }
  ]
})

export default router
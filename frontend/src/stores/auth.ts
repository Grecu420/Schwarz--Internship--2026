import { defineStore } from 'pinia'
import type { LoginRequest, LoginResponse } from '@/generated/proto/auth-api'
import type { User } from '@/generated/proto/user-api'

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null') as User | null,
    token: localStorage.getItem('jwt_token') || (null as string | null)
  }),

  getters: {
    isAuthenticated: (state) => !!state.token
  },

  actions: {
    setSession(token: string, user: User) {
      this.token = token
      this.user = user

      localStorage.setItem('jwt_token', token)
      localStorage.setItem('user', JSON.stringify(user))
    },

    logout() {
      this.token = null
      this.user = null
      localStorage.removeItem('jwt_token')
      localStorage.removeItem('user')
    }
  }
})
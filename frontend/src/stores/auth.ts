import { defineStore } from 'pinia'
import type { User } from '@/generated/proto/user-api'

function isTokenExpired(token: string | null): boolean {
  if (!token) return true
  try {
    const base64Payload = token.split('.')[1]
    if (!base64Payload) return true // failure to extract

    const payload = JSON.parse(atob(base64Payload))
    if (!payload.exp) return false
    return Date.now() >= payload.exp * 1000
  } catch (err) {
    return true // malformed token
  }
}

export const useAuthStore = defineStore('auth', {
  state: () => ({
    user: JSON.parse(localStorage.getItem('user') || 'null') as User | null,
    token: localStorage.getItem('jwt_token') || (null as string | null),
  }),

  getters: {
    isAuthenticated: (state) => !!state.token && !isTokenExpired(state.token),
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
    },

    checkTokenExpiration() : boolean{
      if (this.token && isTokenExpired(this.token)) {
        this.logout()
        return true
      }
      return false
    }
  },
})

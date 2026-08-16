import { defineStore } from 'pinia'
import { api } from '../services/api'
import type { User } from '../types'

export const useSessionStore = defineStore('session', {
  state: () => ({ user: null as User | null, loading: false }),
  actions: {
    async login(email: string, password: string) {
      this.loading = true
      try {
        const { data } = await api.post('/auth/login', { email, password })
        sessionStorage.setItem('access_token', data.access_token)
        localStorage.setItem('refresh_token', data.refresh_token)
      } finally { this.loading = false }
    },
    async logout() {
      const refresh = localStorage.getItem('refresh_token')
      if (refresh) await api.post('/auth/logout', { refresh_token: refresh })
      sessionStorage.removeItem('access_token')
      localStorage.removeItem('refresh_token')
      this.user = null
    }
  }
})

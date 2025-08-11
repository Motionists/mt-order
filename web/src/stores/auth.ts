import { defineStore } from 'pinia'
import http from '../api/http'

export const useAuthStore = defineStore('auth', {
  state: () => ({ token: localStorage.getItem('token') || '' }),
  actions: {
    async login(mobile: string, password: string) {
      const { data } = await http.post('/auth/login', { mobile, password })
      this.token = data.token
      localStorage.setItem('token', this.token)
    },
    logout() {
      this.token = ''
      localStorage.removeItem('token')
    }
  }
})
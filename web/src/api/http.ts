import axios from 'axios'
import { useAuthStore } from '../stores/auth'

const http = axios.create({ baseURL: '/api/v1' })

http.interceptors.request.use((config) => {
  const auth = useAuthStore()
  if (auth.token) {
    config.headers = config.headers || {}
    config.headers.Authorization = `Bearer ${auth.token}`
  }
  return config
})

export default http
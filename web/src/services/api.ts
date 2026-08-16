import axios from 'axios'

export const api = axios.create({ baseURL: '/api/v1', timeout: 8000 })

api.interceptors.request.use(config => {
  const token = sessionStorage.getItem('access_token')
  if (token) config.headers.Authorization = `Bearer ${token}`
  config.headers['X-Request-ID'] = crypto.randomUUID()
  return config
})

api.interceptors.response.use(response => response, async error => {
  const original = error.config
  if (error.response?.status === 401 && !original.__retried) {
    original.__retried = true
    const refresh = localStorage.getItem('refresh_token')
    if (refresh) {
      const { data } = await api.post('/auth/refresh', { refresh_token: refresh })
      sessionStorage.setItem('access_token', data.access_token)
      localStorage.setItem('refresh_token', data.refresh_token)
      return api(original)
    }
  }
  return Promise.reject(error)
})

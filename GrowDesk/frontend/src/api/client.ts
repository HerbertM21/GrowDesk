import axios from 'axios'
import type { AxiosInstance, AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'

// Base URL del backend
const apiClient = axios.create({
  baseURL: 'http://localhost:8000/api',
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000
})

// Interceptor para agregar token de autenticación
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    return config
  },
  (error) => {
    return Promise.reject(error)
  }
)

// Interceptor para manejar errores comunes
apiClient.interceptors.response.use(
  (response) => {
    return response
  },
  (error) => {
    // Redirigir al login si hay error de autorización
    if (error.response && error.response.status === 401) {
      localStorage.removeItem('token')
      localStorage.removeItem('userRole')
      window.location.href = '/login'
    }
    return Promise.reject(error)
  }
)

export default apiClient 
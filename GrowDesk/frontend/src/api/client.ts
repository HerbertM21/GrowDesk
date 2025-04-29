import axios from 'axios'
import type { AxiosInstance, AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'

// Determinar la URL base del backend según el entorno
// Dejamos que se use cualquiera de los dos puertos para evitar problemas de autenticación
const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
console.log('API Base URL:', apiBaseUrl);

const apiClient = axios.create({
  baseURL: apiBaseUrl,
  headers: {
    'Content-Type': 'application/json',
  },
  timeout: 10000
})

// Contador de intentos de actualización del token
let refreshAttempts = 0;
const maxRefreshAttempts = 2;

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
    // Solo redirigir al login después de múltiples intentos fallidos
    // o si el token está completamente ausente
    if (error.response && error.response.status === 401) {
      refreshAttempts++;
      
      if (refreshAttempts >= maxRefreshAttempts) {
        console.log('Múltiples errores de autenticación, cerrando sesión...');
        localStorage.removeItem('token')
        localStorage.removeItem('userRole')
        window.location.href = '/login'
      } else {
        console.log(`Error de autenticación (intento ${refreshAttempts}/${maxRefreshAttempts})`);
      }
    }
    return Promise.reject(error)
  }
)

export default apiClient
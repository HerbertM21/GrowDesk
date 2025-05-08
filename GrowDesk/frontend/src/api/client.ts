import axios from 'axios'
import type { AxiosInstance, AxiosError, AxiosRequestConfig, AxiosResponse } from 'axios'

// Determinar la URL base del backend según el entorno
const apiBaseUrl = import.meta.env.VITE_API_URL || 'http://localhost:8080/api';
console.log('API Base URL:', apiBaseUrl);

const apiClient = axios.create({
  baseURL: apiBaseUrl,
  timeout: 15000, // Reducido a 15 segundos
  headers: {
    'Content-Type': 'application/json',
  },
  validateStatus: (status) => status >= 200 && status < 500,
})

// Interceptor para agregar token de autenticación
apiClient.interceptors.request.use(
  (config) => {
    const token = localStorage.getItem('token')
    if (token) {
      config.headers.Authorization = `Bearer ${token}`
    }
    console.log('Enviando petición a:', config.url, 'con datos:', config.data);
    return config
  },
  (error) => {
    console.error('Error en la configuración de la petición:', error);
    return Promise.reject(error)
  }
)

// Interceptor para manejar errores de respuesta
apiClient.interceptors.response.use(
  (response) => {
    console.log('Respuesta recibida:', response.status, response.data);
    return response;
  },
  (error: AxiosError<{ message: string }>) => {
    console.error('Error en la respuesta:', {
      status: error.response?.status,
      data: error.response?.data,
      message: error.message,
      code: error.code,
      config: {
        url: error.config?.url,
        method: error.config?.method,
        timeout: error.config?.timeout
      }
    });

    // Manejar errores de timeout
    if (error.code === 'ECONNABORTED') {
      return Promise.reject(new Error('El servidor está tardando demasiado en responder. Por favor, intente nuevamente.'));
    }

    // Manejar errores de red
    if (!error.response) {
      return Promise.reject(new Error('No se pudo conectar con el servidor. Por favor, verifique su conexión a internet.'));
    }

    // Manejar errores de autenticación
    if (error.response.status === 401) {
      const currentPath = window.location.pathname;
      if (!currentPath.includes('/login')) {
        localStorage.removeItem('token');
        window.location.href = '/login';
      }
      return Promise.reject(new Error('Sesión expirada. Por favor, inicie sesión nuevamente.'));
    }

    // Manejar errores de permisos
    if (error.response.status === 403) {
      return Promise.reject(new Error('No tiene permisos para realizar esta acción.'));
    }

    // Manejar errores del servidor
    if (error.response.status >= 500) {
      return Promise.reject(new Error('El servidor está experimentando problemas. Por favor, intente nuevamente más tarde.'));
    }

    // Propagar el error con un mensaje más amigable
    const errorMessage = error.response?.data?.message || error.message;
    return Promise.reject(new Error(errorMessage));
  }
)

export default apiClient
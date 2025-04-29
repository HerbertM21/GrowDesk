import axios, { AxiosInstance } from 'axios';
import { jwtDecode } from 'jwt-decode';
import Cookies from 'js-cookie';

// Nombre de la cookie y duración de la sesión
const SESSION_COOKIE_NAME = 'growdesk_session';
const SESSION_EXPIRY_DAYS = 7;
const LOCAL_STORAGE_KEY = 'growdesk_session_data';

// Interfaces para el manejo de sesiones
interface SessionInfo {
  name: string;
  email: string;
  ticketId: string;
  exp?: number;
}

interface TicketCreateRequest {
  name: string;
  email: string;
  message: string;
  metadata?: any;
}

interface MessageRequest {
  ticketId: string;
  message: string;
  userName?: string;
  userEmail?: string;
}

// Función para guardar la sesión en cookies usando formato similar a JWT
const saveSession = (data: SessionInfo) => {
  const now = Math.floor(Date.now() / 1000);
  const expiration = now + (SESSION_EXPIRY_DAYS * 24 * 60 * 60); // 7 días en segundos
  
  const sessionData = {
    ...data,
    exp: expiration
  };
  
  try {
    // Intentar usar cookies primero
    Cookies.set(SESSION_COOKIE_NAME, JSON.stringify(sessionData), {
      expires: SESSION_EXPIRY_DAYS,
      path: '/',
      secure: window.location.protocol === 'https:'
    });
  } catch (error) {
    console.warn('No se pudo guardar sesión en cookies, usando localStorage como alternativa');
  }
  
  // Guardar también en localStorage como respaldo
  try {
    localStorage.setItem(LOCAL_STORAGE_KEY, JSON.stringify(sessionData));
  } catch (error) {
    console.warn('No se pudo guardar en localStorage:', error);
  }
  
  return sessionData;
};

// Obtener datos de sesión de la cookie o localStorage
export const getSession = (): SessionInfo | null => {
  try {
    // Intentar obtener de cookies primero
    const sessionCookie = Cookies.get(SESSION_COOKIE_NAME);
    
    if (sessionCookie) {
      const sessionData = JSON.parse(sessionCookie) as SessionInfo;
      
      // Verificar expiración
      if (sessionData.exp && sessionData.exp < Math.floor(Date.now() / 1000)) {
        // Sesión expirada
        clearSession();
        return null;
      }
      
      return sessionData;
    }
    
    // Si no hay cookie, intentar desde localStorage
    const localData = localStorage.getItem(LOCAL_STORAGE_KEY);
    if (localData) {
      const sessionData = JSON.parse(localData) as SessionInfo;
      
      // Verificar expiración
      if (sessionData.exp && sessionData.exp < Math.floor(Date.now() / 1000)) {
        // Sesión expirada
        clearSession();
        return null;
      }
      
      return sessionData;
    }
    
    return null;
  } catch (error) {
    console.error('Error al obtener sesión:', error);
    clearSession();
    return null;
  }
};

// Limpiar sesión
const clearSession = () => {
  try {
    Cookies.remove(SESSION_COOKIE_NAME, { path: '/' });
  } catch (error) {
    console.warn('Error al eliminar cookie:', error);
  }
  
  try {
    localStorage.removeItem(LOCAL_STORAGE_KEY);
  } catch (error) {
    console.warn('Error al eliminar datos del localStorage:', error);
  }
};

// Clase principal para la API del widget
export const useWidgetApi = () => {
  // Crear instancia de Axios
  const apiClient: AxiosInstance = axios.create({
    baseURL: 'http://localhost:8082/widget',
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
      'X-Widget-ID': 'demo-widget',
      'X-Widget-Token': 'demo-token'
    },
    // No usar withCredentials para evitar problemas de CORS
    withCredentials: false
  });
  
  // Interceptor para agregar datos de sesión a cada solicitud
  apiClient.interceptors.request.use(config => {
    const session = getSession();
    
    if (session) {
      // Agregar información de usuario a headers
      config.headers['X-User-Name'] = session.name;
      config.headers['X-User-Email'] = session.email;
      
      if (session.ticketId) {
        config.headers['X-Ticket-ID'] = session.ticketId;
      }
    }
    
    return config;
  });
  
  // Verificar si hay una sesión activa
  const hasActiveSession = (): boolean => {
    return getSession() !== null;
  };
  
  // Crear un nuevo ticket
  const createTicket = async (data: TicketCreateRequest) => {
    try {
      const response = await apiClient.post('/tickets', data);
      
      if (response.data && response.data.ticketId) {
        // Guardar datos de sesión
        saveSession({
          name: data.name,
          email: data.email,
          ticketId: response.data.ticketId
        });
      }
      
      return response.data;
    } catch (error) {
      console.error('Error creating ticket:', error);
      throw error;
    }
  };
  
  // Enviar un mensaje en un ticket existente
  const sendMessage = async (data: MessageRequest) => {
    try {
      const response = await apiClient.post('/messages', data);
      return response.data;
    } catch (error) {
      console.error('Error sending message:', error);
      throw error;
    }
  };
  
  // Obtener historial de mensajes de un ticket
  const getMessageHistory = async (ticketId: string) => {
    try {
      const response = await apiClient.get(`/tickets/${ticketId}/messages`);
      return response.data;
    } catch (error) {
      console.error('Error getting message history:', error);
      throw error;
    }
  };
  
  // Cerrar sesión (logout)
  const logout = () => {
    clearSession();
  };
  
  return {
    hasActiveSession,
    getSession,
    createTicket,
    sendMessage,
    getMessageHistory,
    logout
  };
}; 
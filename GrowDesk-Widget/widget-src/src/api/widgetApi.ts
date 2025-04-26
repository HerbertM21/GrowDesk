import axios, { AxiosInstance } from 'axios';
import { jwtDecode } from 'jwt-decode';
import Cookies from 'js-cookie';

// Nombre de la cookie y duración de la sesión
const SESSION_COOKIE_NAME = 'growdesk_session';
const SESSION_EXPIRY_DAYS = 7;
const LOCAL_STORAGE_KEY = 'growdesk_session_data';

// Configuración global de la API
export let apiConfig = {
  apiUrl: 'http://localhost:8000/widget',
  widgetId: 'demo-widget',
  widgetToken: 'demo-token'
};

// Función para inicializar la configuración del Widget desde el script
export const initializeFromScript = () => {
  const script = document.getElementById('growdesk-widget');
  
  if (script) {
    // Leer atributos data-* del script
    const apiUrl = script.getAttribute('data-api-url');
    const widgetId = script.getAttribute('data-widget-id');
    const widgetToken = script.getAttribute('data-widget-token');
    const brandName = script.getAttribute('data-brand-name');
    const welcomeMessage = script.getAttribute('data-welcome-message');
    const primaryColor = script.getAttribute('data-primary-color');
    const position = script.getAttribute('data-position');
    
    // Configurar API
    if (apiUrl) apiConfig.apiUrl = apiUrl;
    if (widgetId) apiConfig.widgetId = widgetId;
    if (widgetToken) apiConfig.widgetToken = widgetToken;
    
    // Configurar UI
    window.GrowDeskConfig = {
      ...(window.GrowDeskConfig || {}),
      brandName: brandName || window.GrowDeskConfig?.brandName,
      welcomeMessage: welcomeMessage || window.GrowDeskConfig?.welcomeMessage,
      primaryColor: primaryColor || window.GrowDeskConfig?.primaryColor,
      position: position || window.GrowDeskConfig?.position
    };
    
    console.log('Widget configurado desde el script:', apiConfig, window.GrowDeskConfig);
  } else {
    console.warn('No se encontró el script para configurar el widget.');
  }
};

// Función para configurar la API
export const configureWidgetApi = (config: {
  apiUrl?: string;
  widgetId?: string;
  widgetToken?: string;
}) => {
  if (config.apiUrl) apiConfig.apiUrl = config.apiUrl;
  if (config.widgetId) apiConfig.widgetId = config.widgetId;
  if (config.widgetToken) apiConfig.widgetToken = config.widgetToken;
  
  console.log('API del widget configurada:', apiConfig);
};

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

// Añadir interfaz para FAQs
export interface FAQ {
  id: number;
  question: string;
  answer: string;
  category: string;
  isPublished: boolean;
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
    baseURL: apiConfig.apiUrl,
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
      'X-Widget-ID': apiConfig.widgetId,
      'X-Widget-Token': apiConfig.widgetToken
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
      // Verificar que la URL base sea correcta
      const baseUrl = apiConfig.apiUrl.endsWith('/') ? apiConfig.apiUrl : `${apiConfig.apiUrl}/`;
      
      // Intento estándar con axios
      let response;
      try {
        console.log('Intentando crear ticket en:', `${baseUrl}tickets`);
        response = await apiClient.post('/tickets', data);
      } catch (axiosError) {
        console.warn('Error al usar Axios estándar, intentando con fetch:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        const fetchUrl = `${baseUrl}tickets`;
        console.log('Intentando fetch a:', fetchUrl);
        const fetchResponse = await fetch(fetchUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken
          },
          body: JSON.stringify(data)
        });
        
        if (!fetchResponse.ok) {
          throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText}`);
        }
        
        response = { data: await fetchResponse.json() };
      }
      
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
      // Verificar que la URL base sea correcta
      const baseUrl = apiConfig.apiUrl.endsWith('/') ? apiConfig.apiUrl : `${apiConfig.apiUrl}/`;

      // Intento estándar con axios
      let response;
      try {
        console.log('Intentando enviar mensaje a:', `${baseUrl}messages`);
        response = await apiClient.post('/messages', data);
      } catch (axiosError) {
        console.warn('Error al usar Axios estándar para enviar mensaje, intentando con fetch:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        const fetchUrl = `${baseUrl}messages`;
        console.log('Intentando fetch a:', fetchUrl);
        const fetchResponse = await fetch(fetchUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken,
            'X-User-Name': data.userName || '',
            'X-User-Email': data.userEmail || '',
            'X-Ticket-ID': data.ticketId
          },
          body: JSON.stringify(data)
        });
        
        if (!fetchResponse.ok) {
          throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText}`);
        }
        
        response = { data: await fetchResponse.json() };
      }
      
      return response.data;
    } catch (error) {
      console.error('Error sending message:', error);
      throw error;
    }
  };
  
  // Obtener historial de mensajes de un ticket
  const getMessageHistory = async (ticketId: string) => {
    try {
      // Verificar que la URL base sea correcta
      const baseUrl = apiConfig.apiUrl.endsWith('/') ? apiConfig.apiUrl : `${apiConfig.apiUrl}/`;

      // Intento estándar con axios
      let response;
      try {
        console.log('Intentando obtener mensajes de:', `${baseUrl}tickets/${ticketId}/messages`);
        response = await apiClient.get(`/tickets/${ticketId}/messages`);
      } catch (axiosError) {
        console.warn('Error al usar Axios estándar para obtener mensajes, intentando con fetch:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        const fetchUrl = `${baseUrl}tickets/${ticketId}/messages`;
        console.log('Intentando fetch a:', fetchUrl);
        const fetchResponse = await fetch(fetchUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken,
            'X-Ticket-ID': ticketId
          }
        });
        
        if (!fetchResponse.ok) {
          throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText}`);
        }
        
        response = { data: await fetchResponse.json() };
      }
      
      return response.data;
    } catch (error) {
      console.error('Error getting message history:', error);
      throw error;
    }
  };
  
  // Obtener las preguntas frecuentes publicadas
  const getFaqs = async () => {
    try {
      // Verificar que la URL base sea correcta
      const baseUrl = apiConfig.apiUrl.endsWith('/') ? apiConfig.apiUrl : `${apiConfig.apiUrl}/`;
      const faqsUrl = `${baseUrl}faqs`;
      
      console.log('Intentando obtener FAQs de URL:', faqsUrl);
      console.log('Headers utilizados:', {
        'X-Widget-ID': apiConfig.widgetId,
        'X-Widget-Token': apiConfig.widgetToken
      });

      // Intento estándar con axios
      let response;
      try {
        console.log('Usando axios.get para obtener FAQs');
        response = await apiClient.get('/faqs');
        console.log('Respuesta recibida de axios:', response.status);
      } catch (axiosError) {
        console.warn('Error al usar Axios estándar para obtener FAQs:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        console.log('Intentando fetch a:', faqsUrl);
        const fetchResponse = await fetch(faqsUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken
          }
        });
        
        if (!fetchResponse.ok) {
          const statusText = fetchResponse.statusText || 'Error desconocido';
          console.error(`Error en fetch: ${fetchResponse.status} ${statusText}`);
          throw new Error(`Error ${fetchResponse.status}: ${statusText}`);
        }
        
        const responseData = await fetchResponse.json();
        console.log('Datos recibidos con fetch:', responseData);
        response = { data: responseData };
      }
      
      console.log('FAQs obtenidas correctamente:', response.data);
      return response.data;
    } catch (error) {
      console.error('Error obteniendo FAQs:', error);
      
      // Proporcionar información adicional de depuración
      console.error('Configuración actual:', { 
        apiUrl: apiConfig.apiUrl,
        widgetId: apiConfig.widgetId
      });
      
      // Retornar un array vacío en caso de error para evitar errores adicionales
      return [];
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
    logout,
    getFaqs
  };
}; 
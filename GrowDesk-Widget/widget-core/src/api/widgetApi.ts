import axios, { AxiosInstance } from 'axios';
import { jwtDecode } from 'jwt-decode';
import Cookies from 'js-cookie';

// Configuración
const SESSION_COOKIE_NAME = 'growdesk_session';
const SESSION_EXPIRY_DAYS = 7;
const LOCAL_STORAGE_KEY = 'growdesk_session_data';

// DEPURACION
console.log("[WIDGET] GrowDeskConfig al inicializar widgetApi:", window.GrowDeskConfig);

// Configuración global de la API
export let apiConfig = {
  apiUrl: window.GrowDeskConfig?.apiUrl || 'http://localhost:8082',
  widgetId: 'demo-widget',
  widgetToken: 'demo-token'
};

console.log("[WIDGET] Configuración inicial de API:", apiConfig);

// Función para inicializar la configuración del Widget desde el script
export const initializeFromScript = () => {
  const script = document.getElementById('growdesk-widget');
  
  if (script) {
    console.log("[WIDGET] Obteniendo configuración desde script:", script);
    
    // Leer atributos data-* del script
    const apiUrl = script.getAttribute('data-api-url');
    const widgetId = script.getAttribute('data-widget-id');
    const widgetToken = script.getAttribute('data-widget-token');
    const brandName = script.getAttribute('data-brand-name');
    const welcomeMessage = script.getAttribute('data-welcome-message');
    const primaryColor = script.getAttribute('data-primary-color');
    const position = script.getAttribute('data-position');
    
    // Configurar API
    if (apiUrl) {
      console.log("[WIDGET] Usando apiUrl desde script:", apiUrl);
      apiConfig.apiUrl = apiUrl;
    }
    
    if (widgetId) apiConfig.widgetId = widgetId;
    if (widgetToken) apiConfig.widgetToken = widgetToken;
    
    // Configurar UI
    window.GrowDeskConfig = {
      ...(window.GrowDeskConfig || {}),
      apiUrl: apiUrl || window.GrowDeskConfig?.apiUrl,
      brandName: brandName || window.GrowDeskConfig?.brandName,
      welcomeMessage: welcomeMessage || window.GrowDeskConfig?.welcomeMessage,
      primaryColor: primaryColor || window.GrowDeskConfig?.primaryColor,
      position: position || window.GrowDeskConfig?.position
    };
    
    console.log('[WIDGET] Widget configurado desde script:', apiConfig, window.GrowDeskConfig);
  } else {
    console.warn('[WIDGET] No se encontró el script para configurar el widget.');
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

// Función para guardar la sesión en cookies
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
  console.log('[WIDGET DEBUG] Creando cliente Axios con apiConfig:', apiConfig);
  console.log('[WIDGET DEBUG] URL base que se va a usar:', `${window.GrowDeskConfig?.apiUrl || apiConfig.apiUrl}`);
  
  // Crear instancia de Axios
  const apiClient: AxiosInstance = axios.create({
    baseURL: window.GrowDeskConfig?.apiUrl || apiConfig.apiUrl,
    timeout: 10000,
    headers: {
      'Content-Type': 'application/json',
      'X-Widget-ID': apiConfig.widgetId,
      'X-Widget-Token': apiConfig.widgetToken
    },
    // No HAY usar withCredentials para evitar problemas de CORS
    withCredentials: false
  });
  
  console.log('[WIDGET DEBUG] Cliente Axios creado con baseURL:', apiClient.defaults.baseURL);
  
  // Actualizar la baseURL cada vez que se llama a la API para reflejar cambios en la configuración
  apiClient.interceptors.request.use(config => {
    config.baseURL = window.GrowDeskConfig?.apiUrl || apiConfig.apiUrl;
    console.log('[WIDGET DEBUG] Interceptor - Request a:', config.baseURL + config.url);
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
      console.log('[WIDGET] Intentando crear ticket con los siguientes datos:', data);
      console.log('[WIDGET] URL de API configurada:', apiConfig.apiUrl);
      
      // Intento estándar con axios
      let response;
      try {
        // Asegurarnos de usar la ruta completa correcta con el prefijo '/widget'
        const baseUrl = apiConfig.apiUrl;
        // Construir URL correctamente
        let fullUrl = baseUrl;
        if (!fullUrl.endsWith('/')) {
          fullUrl += '/';
        }
        if (!fullUrl.includes('/widget/')) {
          fullUrl += 'widget/';
        }
        const url = `${fullUrl}tickets`;
        
        console.log('[WIDGET] URL final:', url);
        // Log detallado para depuración
        console.log('[WIDGET] Headers:', {
          'Content-Type': 'application/json',
          'X-Widget-ID': apiConfig.widgetId,
          'X-Widget-Token': apiConfig.widgetToken
        });
        
        // Usar axios directamente sin apiClient
        response = await axios.post(url, data, {
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken
          }
        });
        console.log('[WIDGET] Ticket creado con éxito:', response.data);
      } catch (axiosError: any) {
        console.warn('[WIDGET] Error al usar Axios, intentando con fetch:', axiosError);
        console.warn('[WIDGET] Detalles del error:', {
          status: axiosError?.response?.status,
          message: axiosError?.message,
          data: axiosError?.response?.data
        });
        
        // Si falla, intenta con fetch como respaldo usando la misma URL que antes
        const baseUrl = apiConfig.apiUrl;
        // Construir URL correctamente
        let fullUrl = baseUrl;
        if (!fullUrl.endsWith('/')) {
          fullUrl += '/';
        }
        if (!fullUrl.includes('/widget/')) {
          fullUrl += 'widget/';
        }
        const fetchUrl = `${fullUrl}tickets`;
        
        console.log('[WIDGET] Intentando fetch a:', fetchUrl);
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
          console.error(`[WIDGET] Error en fetch: ${fetchResponse.status} ${fetchResponse.statusText || 'Error desconocido'}`);
          throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText || 'Error desconocido'}`);
        }
        
        const responseData = await fetchResponse.json();
        console.log('[WIDGET] Datos recibidos con fetch:', responseData);
        response = { data: responseData };
      }
      
      if (response.data && response.data.id) {
        // Guardar datos de sesión
        saveSession({
          name: data.name,
          email: data.email,
          ticketId: response.data.id
        });
      }
      
      return response.data;
    } catch (error) {
      console.error('Error creating ticket:', error);
      
      // Intentar crear un ticket simulado para pruebas si estamos en entorno de desarrollo
      if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        console.warn('Creando ticket simulado para entorno de desarrollo');
        const mockTicketId = `MOCK-${Date.now()}`;
        saveSession({
          name: data.name,
          email: data.email,
          ticketId: mockTicketId
        });
        return {
          ticketId: mockTicketId,
          status: 'open',
          message: 'Ticket simulado creado con éxito'
        };
      }
      
      throw error;
    }
  };
  
  // Enviar un mensaje en un ticket existente
  const sendMessage = async (data: MessageRequest) => {
    try {
      console.log('[WIDGET] Intentando enviar mensaje a ticket:', data.ticketId);
      console.log('[WIDGET] URL de API configurada:', apiConfig.apiUrl);
      
      // Verificar que la URL base sea correcta
      const baseUrl = apiConfig.apiUrl.endsWith('/') ? apiConfig.apiUrl : `${apiConfig.apiUrl}/`;
      
      // Si la URL no contiene 'widget', añadirla
      const adjustedBaseUrl = baseUrl.includes('widget') ? baseUrl : `${baseUrl}widget/`;
      
      const messagesUrl = `${adjustedBaseUrl}messages`;
      
      console.log('[WIDGET] URL base ajustada para enviar mensaje:', messagesUrl);
      console.log('[WIDGET] Datos a enviar:', {
        ticketId: data.ticketId,
        message: data.message
      });
      
      // Intento estándar con axios
      let response;
      try {
        console.log('[WIDGET] Usando axios.post para enviar mensaje');
        response = await axios.post(messagesUrl, {
          ticketId: data.ticketId,
          message: data.message
        }, {
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken,
            'X-User-Name': data.userName || '',
            'X-User-Email': data.userEmail || ''
          }
        });
        console.log('[WIDGET] Mensaje enviado con éxito:', response.data);
      } catch (axiosError) {
        console.warn('[WIDGET] Error al usar Axios para enviar mensaje, intentando con fetch:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        console.log('[WIDGET] Intentando fetch a:', messagesUrl);
        const fetchResponse = await fetch(messagesUrl, {
          method: 'POST',
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken,
            'X-User-Name': data.userName || '',
            'X-User-Email': data.userEmail || ''
          },
          body: JSON.stringify({
            ticketId: data.ticketId,
            message: data.message
          })
        });
        
        if (!fetchResponse.ok) {
          console.error(`[WIDGET] Error en fetch: ${fetchResponse.status} ${fetchResponse.statusText || 'Error desconocido'}`);
          throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText || 'Error desconocido'}`);
        }
        
        const responseData = await fetchResponse.json();
        console.log('[WIDGET] Datos recibidos con fetch:', responseData);
        response = { data: responseData };
      }
      
      return response.data;
    } catch (error) {
      console.error('[WIDGET] Error sending message:', error);
      throw error;
    }
  };
  
  // Obtener historial de mensajes de un ticket
  const getMessageHistory = async (ticketId: string) => {
    try {
      console.log('[WIDGET] Intentando obtener historial de mensajes para ticket:', ticketId);
      console.log('[WIDGET] URL de API configurada:', apiConfig.apiUrl);
      
      // Intento estándar con axios
      let response;
      try {
        console.log('[WIDGET] Obteniendo mensajes de:', `${apiConfig.apiUrl}/widget/tickets/${ticketId}/messages`);
        response = await axios.get(`${apiConfig.apiUrl}/widget/tickets/${ticketId}/messages`, {
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken,
            'X-Ticket-ID': ticketId
          }
        });
        console.log('[WIDGET] Mensajes obtenidos con éxito:', response.data);
      } catch (axiosError) {
        console.warn('[WIDGET] Error al usar Axios para obtener mensajes, intentando con fetch:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        const fetchUrl = `${apiConfig.apiUrl}/widget/tickets/${ticketId}/messages`;
        
        console.log('[WIDGET] Intentando fetch a:', fetchUrl);
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
          throw new Error(`Error ${fetchResponse.status}: ${fetchResponse.statusText || 'Error desconocido'}`);
        }
        
        response = { data: await fetchResponse.json() };
      }
      
      return response.data;
    } catch (error) {
      console.error('Error getting message history:', error);
      
      // Proporcionar datos simulados en entorno de desarrollo
      if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        console.warn('Devolviendo historial de mensajes simulado para entorno de desarrollo');
        return { 
          messages: [
            {
              id: 'mock-msg-1',
              content: 'Bienvenido al chat de soporte simulado.',
              isClient: false,
              timestamp: new Date().toISOString()
            }
          ]
        };
      }
      
      throw error;
    }
  };
  
  // Obtener las preguntas frecuentes publicadas
  const getFaqs = async () => {
    try {
      console.log('[WIDGET] Intentando obtener FAQs');
      console.log('[WIDGET] URL de API configurada:', apiConfig.apiUrl);
      
      // Intento estándar con axios
      let response;
      try {
        console.log('[WIDGET] Obteniendo FAQs de:', `${apiConfig.apiUrl}/widget/faqs`);
        response = await axios.get(`${apiConfig.apiUrl}/widget/faqs`, {
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken
          }
        });
        console.log('[WIDGET] FAQs obtenidas con éxito:', response.data);
      } catch (axiosError) {
        console.warn('[WIDGET] Error al usar Axios para obtener FAQs, intentando con fetch:', axiosError);
        
        // Si falla, intenta con fetch como respaldo
        const fetchUrl = `${apiConfig.apiUrl}/widget/faqs`;
        
        console.log('[WIDGET] Intentando fetch a:', fetchUrl);
        const fetchResponse = await fetch(fetchUrl, {
          method: 'GET',
          headers: {
            'Content-Type': 'application/json',
            'X-Widget-ID': apiConfig.widgetId,
            'X-Widget-Token': apiConfig.widgetToken
          }
        });
        
        if (!fetchResponse.ok) {
          const statusText = fetchResponse.statusText || 'Error desconocido';
          console.error(`[WIDGET] Error en fetch: ${fetchResponse.status} ${statusText}`);
          throw new Error(`Error ${fetchResponse.status}: ${statusText}`);
        }
        
        const responseData = await fetchResponse.json();
        console.log('[WIDGET] Datos recibidos con fetch:', responseData);
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
      
      // Devolver datos de ejemplo en entorno de desarrollo
      if (window.location.hostname === 'localhost' || window.location.hostname === '127.0.0.1') {
        console.warn('Devolviendo FAQs simuladas para entorno de desarrollo');
        return [
          {
            id: 1,
            question: '¿Cómo puedo usar este widget?',
            answer: 'Este widget se puede integrar en cualquier sitio web añadiendo el script correspondiente.',
            category: 'General',
            isPublished: true
          },
          {
            id: 2,
            question: '¿Cómo puedo contactar con soporte?',
            answer: 'Puedes hacerlo directamente a través de este chat.',
            category: 'Soporte',
            isPublished: true
          }
        ];
      }
      
      // Retornar un array vacío en caso de error 
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
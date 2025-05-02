import { createApp } from 'vue';
import { createPinia } from 'pinia';
import ChatWidget from './components/ChatWidget.vue';
import { initializeFromScript, configureWidgetApi } from './api/widgetApi';
import 'primeicons/primeicons.css';
import './styles.css';

console.log('GrowDesk Widget - Iniciando...');

// Crear un contenedor para el widget
const createWidgetContainer = (): HTMLElement => {
  console.log('Creando contenedor para el widget');
  const container = document.createElement('div');
  container.id = 'growdesk-widget-container';
  document.body.appendChild(container);
  return container;
};

// Inicializar el widget
const initializeWidget = () => {
  console.log('Inicializando widget');
  
  // Obtener configuración del script
  initializeFromScript();

  // Configurar la API con la URL proporcionada en GrowDeskConfig
  if (window.GrowDeskConfig?.apiUrl) {
    console.log('Configurando API URL desde GrowDeskConfig:', window.GrowDeskConfig.apiUrl);
    
    // Asegurarnos de que la URL esté correctamente formateada
    let apiUrl = window.GrowDeskConfig.apiUrl;
    // No añadimos '/widget' porque lo hacemos en cada llamada específica
    
    configureWidgetApi({
      apiUrl: apiUrl
    });
  }

  // Crear el store global
  const pinia = createPinia();

  // Crear contenedor y montar la aplicación Vue
  const container = createWidgetContainer();
  const app = createApp(ChatWidget, {
    // Propiedades que se pueden personalizar
    brandName: window.GrowDeskConfig?.brandName || 'GrowDesk',
    welcomeMessage: window.GrowDeskConfig?.welcomeMessage || '¿Necesitas ayuda? Estamos aquí para ayudarte.',
    primaryColor: window.GrowDeskConfig?.primaryColor || '#1976d2',
    position: window.GrowDeskConfig?.position || 'bottom-right',
    logoUrl: window.GrowDeskConfig?.logoUrl || ''
  });

  // Usar Pinia
  app.use(pinia);

  // Montar la aplicación
  app.mount(container);
  
  console.log('Widget montado correctamente');

  // Crear y retornar una API pública para interactuar con el widget
  const widgetApi = {
    open: () => {
      console.log('Abriendo widget vía API');
      const widgetComponent = document.querySelector('#growdesk-widget-container');
      if (widgetComponent) {
        const event = new CustomEvent('open-widget');
        widgetComponent.dispatchEvent(event);
      }
    },
    close: () => {
      console.log('Cerrando widget vía API');
      const widgetComponent = document.querySelector('#growdesk-widget-container');
      if (widgetComponent) {
        const event = new CustomEvent('close-widget');
        widgetComponent.dispatchEvent(event);
      }
    },
    configure: (config: any) => {
      console.log('Reconfigurando widget', config);
      // Actualizar configuración de la API
      if (config.apiUrl || config.widgetId || config.widgetToken) {
        configureWidgetApi({
          apiUrl: config.apiUrl,
          widgetId: config.widgetId,
          widgetToken: config.widgetToken
        });
      }

      // Actualizar configuración de la UI
      window.GrowDeskConfig = {
        ...window.GrowDeskConfig,
        ...config
      };
    }
  };
  
  console.log('API del widget creada', widgetApi);
  return widgetApi;
};

// Interfaz para la configuración global del widget
declare global {
  interface Window {
    GrowDeskConfig?: {
      apiUrl?: string;
      brandName?: string;
      welcomeMessage?: string;
      primaryColor?: string;
      position?: string;
      logoUrl?: string;
      ticketApiKey?: string;
    };
    GrowDeskWidget?: any;
  }
}

// Auto-inicializar cuando el script se carga
let widgetInitialized = false;

const initialize = () => {
  if (widgetInitialized) return;
  
  console.log('Inicializando GrowDeskWidget automáticamente');
  const api = initializeWidget();
  window.GrowDeskWidget = api;
  widgetInitialized = true;
};

// Inicializar cuando el DOM esté listo
if (document.readyState === 'loading') {
  document.addEventListener('DOMContentLoaded', initialize);
} else {
  initialize();
}

// Exponer la función de inicialización para uso manual
export default initializeWidget; 
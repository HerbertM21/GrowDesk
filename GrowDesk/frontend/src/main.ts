import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'
// Importar Bootstrap CSS y JS
import 'bootstrap/dist/css/bootstrap.min.css'
import 'bootstrap/dist/js/bootstrap.bundle.min.js'
import './assets/main.css'
import { useUsersStore } from './stores/users'
import { useAuthStore } from './stores/auth'

// Función para limpiar localStorage de datos corruptos
const cleanLocalStorage = async () => {
  try {
    // Comprobar si es necesario limpiar basado en versión en archivo JSON
    const response = await fetch('/localStorage-fix.json');
    if (response.ok) {
      const { id, version } = await response.json();
      const lastCleanVersion = localStorage.getItem('localStorage-clean-version');
      
      // Si la versión es nueva o no existe, limpiar y actualizar versión
      if (!lastCleanVersion || parseInt(lastCleanVersion) < version) {
        console.log('Iniciando limpieza de localStorage...');
        
        // Importar utilitarios de validación
        const { 
          cleanLocalStorageData, 
          isValidUser, 
          isValidTicket, 
          isValidCategory 
        } = await import('./utils/validators');
        
        // Limpiar cada tipo de datos con su validador adecuado
        cleanLocalStorageData<any>('growdesk-users', isValidUser);
        cleanLocalStorageData<any>('growdesk_tickets', isValidTicket);
        cleanLocalStorageData<any>('growdesk-categories', isValidCategory);
        
        // Marcar que se completó la limpieza con esta versión
        localStorage.setItem('localStorage-clean-version', version.toString());
        console.log('Limpieza de localStorage completada');
      }
    }
  } catch (err) {
    console.error('Error al limpiar localStorage:', err);
  }
};

// Llamar a la función de limpieza antes de iniciar la app
cleanLocalStorage().catch(err => console.error('Error en proceso de limpieza:', err));

// Crear app
const app = createApp(App)

// Configurar pinia
const pinia = createPinia()
app.use(pinia)

// Configurar router
app.use(router)

// Configurar PrimeVue
app.use(PrimeVue)

// Inicializar stores con datos mock para desarrollo, pero NO iniciar sesión automáticamente
if (import.meta.env.DEV) {
  setTimeout(async () => {
    // Inicializar usuarios mock para que estén disponibles para el login
    const userStore = useUsersStore()
    userStore.initMockUsers()
    console.log('Usuarios mock inicializados desde main.ts (solo para login)')
    
    // Proporcionar el router al auth store
    const authStore = useAuthStore()
    authStore.setRouter(router)
    console.log('Router proporcionado al auth store')
    
    // Solo comprobar si hay una sesión activa (token válido existente)
    // pero NO forzar una autenticación
    const isAuthenticated = await authStore.checkAuth()
    console.log('App inicializada, estado de autenticación:', isAuthenticated ? 'autenticado' : 'no autenticado')
  }, 100)
}

// Montar app
app.mount('#app')
console.log('App montada, estado de autenticación: no autenticado')
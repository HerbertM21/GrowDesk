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
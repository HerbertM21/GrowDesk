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

// Inicializar stores con datos mock para desarrollo
if (import.meta.env.DEV) {
  setTimeout(async () => {
    const userStore = useUsersStore()
    userStore.initMockUsers()
    console.log('Usuarios mock inicializados desde main.ts')
    
    // Proporcionar el router al auth store
    const authStore = useAuthStore()
    authStore.setRouter(router)
    console.log('Router proporcionado al auth store')
    
    // Verificar si hay una sesión activa y cargar los datos del usuario
    await authStore.checkAuth()
    console.log('Estado de autenticación verificado:', authStore.isAuthenticated)
    console.log('Roles de usuario - Admin:', authStore.isAdmin, 'Asistente:', authStore.isAssistant)
  }, 100)
}

// Montar app
app.mount('#app')
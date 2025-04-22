import { createApp } from 'vue'
import { createPinia } from 'pinia'
import App from './App.vue'
import router from './router'
import PrimeVue from 'primevue/config'
import 'primevue/resources/themes/lara-light-blue/theme.css'
import 'primevue/resources/primevue.min.css'
import 'primeicons/primeicons.css'
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
  setTimeout(() => {
    const userStore = useUsersStore()
    userStore.initMockUsers()
    console.log('Usuarios mock inicializados desde main.ts')
    
    // Proporcionar el router al auth store
    const authStore = useAuthStore()
    authStore.setRouter(router)
    console.log('Router proporcionado al auth store')
  }, 100)
}

// Montar app
app.mount('#app')
/* eslint-disable */
<template>
  <div class="app-container" :class="{ 'dark-theme': isDarkTheme }">
    <!-- Solo mostrar la navegación cuando el usuario está autenticado o cuando no estamos en las rutas de login/registro -->
    <MainNavigation v-if="!isAuthRoute" @toggle-theme="toggleTheme" :isDarkTheme="isDarkTheme" />
    
    <!-- Navegación lateral para usuarios autenticados -->
    <SideNavigation v-if="isAuthenticated && !isAuthRoute" />
    
    <!-- Componente de notificaciones global -->
    <NotificationsManager />
    
    <main class="main-content" :class="{ 
      'auth-content': isAuthRoute,
      'app-content': isAuthenticated && !isAuthRoute
    }">
      <router-view v-slot="{ Component }">
        <keep-alive>
          <component :is="Component" :key="$route.fullPath" />
        </keep-alive>
      </router-view>
    </main>

    <footer class="footer" :class="{ 
      'with-sidenav': isAuthenticated && !isAuthRoute
    }">
      <p>&copy; 2025 GrowDesk. Todos los derechos reservados.</p>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed, ref, watch } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useRoute } from 'vue-router';
import MainNavigation from '@/components/layout/MainNavigation.vue';
import SideNavigation from '@/components/layout/SideNavigation.vue';
import NotificationsManager from '@/components/notifications/NotificationsManager.vue';

// Inicializar auth store
const authStore = useAuthStore();
const route = useRoute();

// Estado para el tema
const isDarkTheme = ref(localStorage.getItem('theme') === 'dark');

// Verificar si el usuario está autenticado
const isAuthenticated = computed(() => {
  return authStore.isAuthenticated;
});

// Toggle del tema
const toggleTheme = () => {
  isDarkTheme.value = !isDarkTheme.value;
  localStorage.setItem('theme', isDarkTheme.value ? 'dark' : 'light');
  
  // Actualizar el tema en el documento HTML
  if (isDarkTheme.value) {
    document.documentElement.classList.add('dark-theme');
  } else {
    document.documentElement.classList.remove('dark-theme');
  }
};

// Observar cambios en el tema para actualizar el HTML
watch(isDarkTheme, (newValue: boolean) => {
  if (newValue) {
    document.documentElement.classList.add('dark-theme');
  } else {
    document.documentElement.classList.remove('dark-theme');
  }
});

// Verificar si estamos en una ruta de autenticación (login, raíz, o registro)
const isAuthRoute = computed(() => {
  return route.path === '/' || route.path === '/login' || route.path === '/register';
});

// Verificar autenticación al montar el componente
onMounted(async () => {
  // Configurar el tema inicial
  if (isDarkTheme.value) {
    document.documentElement.classList.add('dark-theme');
  } else {
    document.documentElement.classList.remove('dark-theme');
  }
  
  console.log('App.vue montado - verificando autenticación');
  
  // Inicializar el estado de autenticación
  try {
    // Intentar verificar si hay una sesión válida
    const isAuthSuccess = await authStore.checkAuth();
    console.log('Verificación de autenticación completada:', isAuthSuccess ? 'sesión válida' : 'no autenticado');
  } catch (error) {
    console.error('Error durante la verificación de autenticación:', error);
    // Limpiar datos de autenticación en caso de error
    authStore.logout();
  }
});
</script>

<style lang="scss">
.app-container {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-primary);
  color: var(--text-primary);
  transition: background-color 0.3s ease, color 0.3s ease;
}

.main-content {
  flex: 1;
  padding: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  margin-top: 0.75rem; 
  
  &.auth-content {
    padding: 0;
    margin-top: 0;
    max-width: none;
  }
  
  &.app-content {
    margin-left: 280px; 
    max-width: calc(100% - 280px);
  }
}

.footer {
  padding: 1rem;
  text-align: center;
  background-color: var(--footer-bg);
  color: var(--text-secondary);
  border-top: 1px solid var(--border-color);
  transition: background-color 0.3s ease, color 0.3s ease;
  
  &.with-sidenav {
    margin-left: 280px;
  }
}

.btn {
  padding: 0.5rem 1rem;
  border-radius: 4px;
  border: none;
  cursor: pointer;
  text-decoration: none;
  font-weight: 500;
  transition: background-color 0.2s;

  &-primary {
    background-color: #1976d2;
    color: white;

    &:hover {
      background-color: #1565c0;
    }
  }

  &-secondary {
    background-color: #f5f5f5;
    color: #2c3e50;

    &:hover {
      background-color: #e0e0e0;
    }
  }
}

* {
  font-family: var(--font-family);
}

select, input, textarea, button {
  font-family: inherit;
}
</style>
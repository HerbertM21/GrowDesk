<template>
  <div class="main-navigation" :class="{ 'dark-theme': isDarkTheme }">
    <nav class="navbar">
      <div class="navbar-brand">
        <router-link to="/" class="logo-link">
          <img src="@/assets/logo.png" alt="GrowDesk Logo" class="logo-image">
        </router-link>
      </div>
      
      <div class="nav-right" v-if="isAuthenticated">
        
        <button @click="toggleTheme" class="theme-toggle" :title="isDarkTheme ? 'Cambiar a modo claro' : 'Cambiar a modo oscuro'">
          <i :class="['pi', isDarkTheme ? 'pi-sun' : 'pi-moon']"></i>
        </button>
        
        
        <div class="user-dropdown">
          <button @click="toggleUserDropdown" class="user-toggle" :class="{ active: userDropdownOpen }">
            <div class="user-avatar" v-if="userInitials">{{ userInitials }}</div>
            <span class="user-name">{{ fullName }}</span>
            <i class="pi pi-chevron-down dropdown-icon"></i>
          </button>
          <div class="dropdown-menu user-menu" v-if="userDropdownOpen">
            <button @click="navigateToProfile" class="dropdown-item">
              <i class="pi pi-user"></i>
              Mi Perfil
            </button>
            <button @click="navigateToSettings" class="dropdown-item">
              <i class="pi pi-cog"></i>
              Configuración
            </button>
            <button @click="logout" class="dropdown-item">
              <i class="pi pi-sign-out"></i>
              Cerrar Sesión
            </button>
          </div>
        </div>
      </div>
    </nav>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { storeToRefs } from 'pinia';

// Recibir props del tema
interface Props {
  isDarkTheme?: boolean;
}
const props = withDefaults(defineProps<Props>(), {
  isDarkTheme: false
});

// Definir eventos
const emit = defineEmits<{
  (e: 'toggle-theme'): void;
}>();

const router = useRouter();
const authStore = useAuthStore();
const { user } = storeToRefs(authStore);

// Estado para dropdown de usuario
const userDropdownOpen = ref(false);

// Función para toggle del tema
const toggleTheme = () => {
  emit('toggle-theme');
};

const toggleUserDropdown = () => {
  userDropdownOpen.value = !userDropdownOpen.value;
};

// Funciones de navegación
const navigateToProfile = () => {
  userDropdownOpen.value = false;
  router.push('/profile');
};

const navigateToSettings = () => {
  userDropdownOpen.value = false;
  router.push('/settings');
};

// Función para cerrar los menús al hacer clic fuera
const closeMenusOnOutsideClick = (event: MouseEvent) => {
  const target = event.target as HTMLElement;
  
  if (!target.closest('.user-dropdown') && userDropdownOpen.value) {
    userDropdownOpen.value = false;
  }
};

// Computar información del usuario
const isAuthenticated = computed(() => {
  return authStore.token !== null && authStore.user !== null;
});

const isAdminOrAssistant = computed(() => {
  return authStore.isAdmin || authStore.isAssistant;
});

const isAdmin = computed(() => {
  return authStore.isAdmin;
});

const fullName = computed(() => {
  if (authStore.user) {
    return `${authStore.user.firstName} ${authStore.user.lastName}`;
  }
  return 'Usuario';
});

const userInitials = computed(() => {
  if (authStore.user) {
    return `${authStore.user.firstName.charAt(0)}${authStore.user.lastName.charAt(0)}`;
  }
  return null;
});

// Función para cerrar sesión
const logout = () => {
  authStore.logout();
  router.push('/login');
};

// Ruta activa
const route = useRoute();
const isActive = (path: string) => {
  return route.path.startsWith(path);
};

// Event listeners
onMounted(() => {
  document.addEventListener('click', closeMenusOnOutsideClick);
  console.log('MainNavigation montado, isAdmin:', isAdminOrAssistant.value);
  
  // Limpiar event listener cuando el componente se desmonte
  return () => {
    document.removeEventListener('click', closeMenusOnOutsideClick);
  };
});
</script>

<style lang="scss" scoped>
.main-navigation {
  position: sticky;
  top: 0;
  z-index: 100;
  background-color: var(--header-bg);
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  transition: background-color 0.3s ease, color 0.3s ease;
  border-bottom: 1px solid var(--border-color);
}

.navbar {
  display: flex;
  align-items: center;
  justify-content: space-between;
  height: 60px; 
  padding: 0;
  max-width: 100%;
  width: 100%;
  margin: 0;
}

.navbar-brand {
  display: flex;
  align-items: center;
  position: absolute;
  left: 0; 
  top: 0;
  height: 60px;
  padding-left: 0.5rem;
  z-index: 10;
  
  .logo-link {
    display: flex;
    align-items: center;
    text-decoration: none;
    
    .logo-image {
      height: 45px;
      width: auto;
      margin-left: 10px;
    }
  }
}

.nav-right {
  display: flex;
  align-items: center;
  gap: 1rem;
  margin-left: auto;
  padding-right: 1rem;
}

.theme-toggle {
  background: none;
  border: none;
  cursor: pointer;
  width: 36px;
  height: 36px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: var(--text-primary);
  background-color: var(--bg-tertiary);
  transition: all 0.2s ease;
  
  &:hover {
    background-color: var(--hover-bg);
    transform: scale(1.05);
  }
  
  i {
    font-size: 1.1rem;
  }
}

.user-dropdown {
  position: relative;
  
  .user-toggle {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: none;
    border: none;
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 4px;
    color: var(--text-primary);
    
    &:hover {
      background-color: var(--hover-bg);
    }
    
    .user-avatar {
      width: 30px;
      height: 30px;
      border-radius: 50%;
      background-color: var(--primary-color);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      font-size: 0.85rem;
    }
    
    .user-name {
      font-weight: 500;
    }
    
    .dropdown-icon {
      font-size: 0.75rem;
      transition: transform 0.2s;
    }
    
    &.active .dropdown-icon {
      transform: rotate(180deg);
    }
  }
}

.dropdown-menu {
  position: absolute;
  top: 100%;
  right: 0;
  margin-top: 0.5rem;
  min-width: 180px;
  background-color: var(--bg-secondary);
  border-radius: 4px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  z-index: 100;
  border: 1px solid var(--border-color);
  
  &.user-menu {
    right: 0;
  }
  
  .dropdown-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.75rem 1rem;
    width: 100%;
    text-align: left;
    background: none;
    border: none;
    cursor: pointer;
    color: var(--text-primary);
    
    &:hover {
      background-color: var(--hover-bg);
    }
    
    i {
      color: var(--text-secondary);
    }
  }
}

@media (max-width: 768px) {
  .navbar {
    padding: 0;
  }
  
  .navbar-brand {
    padding-left: 0.25rem;
  }
  
  .nav-right {
    padding-right: 0.5rem;
  }
  
  .user-toggle .user-name {
    display: none;
  }
}
</style> 
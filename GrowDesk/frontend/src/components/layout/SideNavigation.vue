<template>
  <div class="side-navigation">
    <div class="nav-section">
      <h3>Principal</h3>
      <div class="nav-items">
        <router-link to="/dashboard" class="nav-item" :class="{ 'active': isActive('/dashboard') }">
          <i class="pi pi-chart-bar"></i>
          <span>Dashboard</span>
        </router-link>
        
        <router-link to="/tickets" class="nav-item" :class="{ 'active': isActive('/tickets') }">
          <i class="pi pi-ticket"></i>
          <span>Tickets</span>
        </router-link>
      </div>
    </div>
    
    <div class="nav-section" v-if="isAdminOrAssistant">
      <h3>Administración</h3>
      <div class="nav-items">
        <router-link to="/admin/dashboard" class="nav-item" :class="{ 'active': isActive('/admin/dashboard') }">
          <i class="pi pi-cog"></i>
          <span>Panel Admin</span>
        </router-link>
        
        <router-link v-if="isAdmin" to="/admin/users" class="nav-item" :class="{ 'active': isActive('/admin/users') }">
          <i class="pi pi-users"></i>
          <span>Usuarios</span>
        </router-link>
        
        <router-link to="/admin/profile-management" class="nav-item" :class="{ 'active': isActive('/admin/profile-management') }">
          <i class="pi pi-user-edit"></i>
          <span>Perfiles</span>
        </router-link>
        
        <router-link v-if="isAdmin" to="/admin/categories" class="nav-item" :class="{ 'active': isActive('/admin/categories') }">
          <i class="pi pi-tags"></i>
          <span>Categorías</span>
        </router-link>
        
        <router-link to="/admin/faqs" class="nav-item" :class="{ 'active': isActive('/admin/faqs') }">
          <i class="pi pi-question-circle"></i>
          <span>FAQs</span>
        </router-link>
        
        <router-link v-if="isAdmin" to="/admin/widget-config" class="nav-item" :class="{ 'active': isActive('/admin/widget-config') }">
          <i class="pi pi-wrench"></i>
          <span>Widget</span>
        </router-link>
      </div>
    </div>
    
    <div class="nav-section">
      <h3>Mi Cuenta</h3>
      <div class="nav-items">
        <router-link to="/profile" class="nav-item" :class="{ 'active': isActive('/profile') }">
          <i class="pi pi-user"></i>
          <span>Mi Perfil</span>
        </router-link>
        
        <router-link to="/settings" class="nav-item" :class="{ 'active': isActive('/settings') }">
          <i class="pi pi-cog"></i>
          <span>Configuración</span>
        </router-link>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue';
import { useRoute } from 'vue-router';
import { useAuthStore } from '@/stores/auth';

interface Props {
  // Propiedades definidas para el componente
}
// No necesitamos definir props en este caso, pero es necesario para TypeScript
defineProps<Props>();

const route = useRoute();
const authStore = useAuthStore();

// Función para verificar si una ruta está activa
const isActive = (path: string) => {
  return route.path.startsWith(path);
};

// Computar roles del usuario
const isAdminOrAssistant = computed(() => {
  return authStore.isAdmin || authStore.isAssistant;
});

const isAdmin = computed(() => {
  return authStore.isAdmin;
});
</script>

<style lang="scss" scoped>
.side-navigation {
  width: 280px; /* Ancho fijo */
  height: 100%;
  background-color: var(--bg-secondary);
  border-right: 1px solid var(--border-color);
  position: fixed;
  left: 0;
  top: 60px; /* Ajustado a la altura del header */
  bottom: 0;
  overflow-y: auto;
  padding-bottom: 1.5rem;
  padding-top: 1rem;
  z-index: 10;
  
  .nav-section {
    margin-bottom: 1.75rem;
    
    h3 {
      font-size: 0.85rem;
      text-transform: uppercase;
      color: var(--text-secondary);
      margin-left: 1.5rem;
      margin-bottom: 0.75rem;
      font-weight: 600;
    }
  }
  
  .nav-items {
    display: flex;
    flex-direction: column;
  }
  
  .nav-item {
    display: flex;
    align-items: center;
    padding: 0.85rem 1.5rem; /* Más espacio vertical y horizontal */
    color: var(--text-primary);
    text-decoration: none;
    transition: all 0.2s;
    position: relative;
    
    i {
      font-size: 1.2rem; /* Íconos más grandes */
      margin-right: 1rem; /* Más espacio entre ícono y texto */
      color: var(--text-secondary);
      transition: color 0.2s;
    }
    
    span {
      white-space: nowrap;
      font-weight: 500; /* Texto un poco más grueso */
    }
    
    &:hover {
      background-color: var(--hover-bg);
      
      i {
        color: var(--primary-color);
      }
    }
    
    &.active {
      background-color: var(--bg-tertiary);
      color: var(--primary-color);
      border-left: 3px solid var(--primary-color);
      padding-left: calc(1.5rem - 3px);
      
      i {
        color: var(--primary-color);
      }
    }
  }
}

/* Ajustar margen al contenido principal para el sidebar fijo */
:deep(.app-content) {
  margin-left: 280px;
}
</style>
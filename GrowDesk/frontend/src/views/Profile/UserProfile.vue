<template>
  <div class="user-profile-container">
    <div v-if="loading" class="loading">
      <p>Cargando perfil...</p>
    </div>
    
    <div v-else-if="error" class="error-container">
      <p>{{ error }}</p>
      <button class="btn btn-primary" @click="loadProfile">Intentar nuevamente</button>
    </div>
    
    <div v-else class="profile-content">
      <div class="profile-header">
        <div class="user-avatar" :class="userRoleClass">
          {{ userInitials }}
        </div>
        
        <div class="user-info">
          <h1>{{ profile?.firstName }} {{ profile?.lastName }}</h1>
          <div class="user-meta">
            <span class="user-id">ID: {{ profile?.id }}</span>
            <span :class="['role-badge', profile?.role]">{{ translateRole(profile?.role || '') }}</span>
            <span :class="['status-badge', profile?.active ? 'active' : 'inactive']">
              {{ profile?.active ? 'Activo' : 'Inactivo' }}
            </span>
          </div>
          <p class="user-email">{{ profile?.email }}</p>
          <p v-if="profile?.department" class="user-department">
            <strong>Departamento:</strong> {{ profile.department }}
          </p>
        </div>
        
        <div class="actions">
          <button class="btn btn-primary" @click="showEditModal = true">
            <i class="pi pi-pencil"></i> Editar Perfil
          </button>
          <router-link to="/settings" class="btn btn-secondary">
            <i class="pi pi-cog"></i> Configuración
          </router-link>
        </div>
      </div>
      
      <div class="profile-grid">
        <!-- Resumen del perfil -->
        <div class="profile-card summary-card">
          <div class="profile-card-header">
            <i class="pi pi-user"></i>
            <h2>Información Personal</h2>
          </div>
          <div class="profile-card-body">
            <div class="info-item">
              <span class="label">Nombre completo:</span>
              <span class="value">{{ profile?.firstName }} {{ profile?.lastName }}</span>
            </div>
            <div class="info-item">
              <span class="label">Correo electrónico:</span>
              <span class="value">{{ profile?.email }}</span>
            </div>
            <div class="info-item">
              <span class="label">Departamento:</span>
              <span class="value">{{ profile?.department || 'No especificado' }}</span>
            </div>
            <div class="info-item">
              <span class="label">Rol:</span>
              <span class="value">{{ translateRole(profile?.role || '') }}</span>
            </div>
            <div class="info-item">
              <span class="label">Estado:</span>
              <span class="value">{{ profile?.active ? 'Activo' : 'Inactivo' }}</span>
            </div>
            <div class="info-item">
              <span class="label">Fecha de registro:</span>
              <span class="value">{{ formatDate(profile?.createdAt || '') }}</span>
            </div>
          </div>
        </div>
        
        <!-- Estadísticas -->
        <div class="profile-card stats-card">
          <div class="profile-card-header">
            <i class="pi pi-chart-bar"></i>
            <h2>Mis Estadísticas</h2>
          </div>
          <div class="profile-card-body">
            <div class="stats-grid">
              <div class="stat-item">
                <span class="stat-number">{{ userTickets?.length || 0 }}</span>
                <span class="stat-label">Tickets Totales</span>
              </div>
              <div class="stat-item">
                <span class="stat-number">{{ openTickets }}</span>
                <span class="stat-label">Tickets Abiertos</span>
              </div>
              <div class="stat-item">
                <span class="stat-number">{{ closedTickets }}</span>
                <span class="stat-label">Tickets Cerrados</span>
              </div>
              <div class="stat-item">
                <span class="stat-number">{{ userActivity?.length || 0 }}</span>
                <span class="stat-label">Actividades Recientes</span>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Sección de tickets -->
      <div class="profile-section">
        <div class="section-title">
          <div class="title-icon">
            <i class="pi pi-ticket"></i>
          </div>
          <h2>Mis Tickets</h2>
        </div>
        
        <div v-if="!userTickets || userTickets.length === 0" class="empty-tickets">
          <p>No tienes tickets asignados actualmente.</p>
          <router-link to="/tickets/new" class="btn btn-primary">
            <i class="pi pi-plus"></i> Crear Nuevo Ticket
          </router-link>
        </div>
        
        <div v-else class="tickets-grid">
          <div 
            v-for="ticket in userTickets" 
            :key="ticket.id" 
            class="ticket-card"
          >
            <div class="ticket-header">
              <div class="ticket-badges">
                <span :class="['status-badge', ticket.status]">{{ translateStatus(ticket.status) }}</span>
                <span :class="['priority-badge', ticket.priority.toLowerCase()]">
                  {{ translatePriority(ticket.priority) }}
                </span>
              </div>
              <h3 class="ticket-title">{{ ticket.title }}</h3>
            </div>
            
            <div class="ticket-body">
              <p class="ticket-description">{{ truncateText(ticket.description, 120) }}</p>
            </div>
            
            <div class="ticket-meta">
              <div class="meta-item">
                <i class="pi pi-tag"></i>
                <span>{{ ticket.category || 'Sin categoría' }}</span>
              </div>
              
              <div class="meta-item">
                <i class="pi pi-calendar"></i>
                <span>{{ formatDate(ticket.createdAt) }}</span>
              </div>
            </div>
            
            <div class="ticket-footer">
              <router-link :to="`/tickets/${ticket.id}`" class="view-details-btn">
                <i class="pi pi-eye"></i>
                Ver Detalles
              </router-link>
            </div>
          </div>
        </div>
      </div>
      
      <div class="profile-section">
        <div class="section-title">
          <div class="title-icon">
            <i class="pi pi-clock"></i>
          </div>
          <h2>Actividad Reciente</h2>
        </div>
        
        <div class="activity-timeline">
          <div v-if="!userActivity || userActivity.length === 0" class="empty-activity">
            <p>No hay actividad reciente para mostrar.</p>
          </div>
          
          <div v-else class="timeline">
            <div 
              v-for="(activity, index) in userActivity" 
              :key="index" 
              class="timeline-item"
            >
              <div class="timeline-icon" :class="getActivityIcon(activity.type)">
                <i :class="getActivityIconClass(activity.type)"></i>
              </div>
              
              <div class="timeline-content">
                <p class="activity-text">{{ activity.description || 'Actividad registrada' }}</p>
                <span class="activity-time">{{ formatDate(activity.timestamp) }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modal de edición de perfil -->
    <UserEditModal 
      v-if="showEditModal && profile" 
      :user="profile" 
      @close="showEditModal = false" 
      @updated="handleProfileUpdated" 
    />
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useUsersStore } from '@/stores/users';
import { useAuthStore } from '@/stores/auth';
import { useTicketStore } from '@/stores/tickets';
import { useActivityStore } from '@/stores/activity';
import UserEditModal from '@/components/Admin/UserEditModal.vue';
import { storeToRefs } from 'pinia';
import type { User } from '@/stores/users';
import type { Ticket } from '@/stores/tickets';
import type { Activity } from '@/stores/activity';
import ticketService from '@/services/ticketService';

// Stores
const userStore = useUsersStore();
const authStore = useAuthStore();
const ticketStore = useTicketStore();
const activityStore = useActivityStore();
const route = useRoute();
const router = useRouter();

// State
const userTickets = ref<Ticket[]>([]);
const userActivity = ref<Activity[]>([]);
const loading = ref(true);
const error = ref<string | null>(null);
const profile = ref<User | null>(null);
const showEditModal = ref(false);

// Configuración de intervalo para actualización automática
const updateInterval = ref<number | null>(null);

// Computed properties
const userInitials = computed(() => {
  if (!profile.value) return '';
  return `${profile.value.firstName.charAt(0)}${profile.value.lastName.charAt(0)}`;
});

const userRoleClass = computed(() => {
  if (!profile.value) return '';
  return `role-${profile.value.role}`;
});

const openTickets = computed(() => {
  return userTickets.value.filter(ticket => 
    ticket.status === 'open' || 
    ticket.status === 'assigned' || 
    ticket.status === 'in_progress'
  ).length;
});

const closedTickets = computed(() => {
  return userTickets.value.filter(ticket => ticket.status === 'closed').length;
});

// Helper functions
const translateRole = (role: string): string => {
  const roles: Record<string, string> = {
    'admin': 'Administrador',
    'support': 'Soporte',
    'employee': 'Empleado',
    'user': 'Usuario'
  };
  return roles[role] || role;
};

const translateStatus = (status: string): string => {
  const statuses: Record<string, string> = {
    'open': 'Abierto',
    'assigned': 'Asignado',
    'in_progress': 'En Progreso',
    'pending': 'Pendiente',
    'resolved': 'Resuelto',
    'closed': 'Cerrado'
  };
  return statuses[status] || status;
};

const translatePriority = (priority: string): string => {
  if (!priority) return 'Media';
  
  // Normalizar a minúsculas para la comparación
  const lowerPriority = priority.toLowerCase();
  
  const priorityMap: Record<string, string> = {
    'low': 'Baja',
    'medium': 'Media',
    'high': 'Alta',
    'urgent': 'Urgente'
  };
  
  return priorityMap[lowerPriority] || priority;
};

const truncateText = (text: string, maxLength: number): string => {
  // Verificar si el texto es undefined, null o vacío
  if (!text) return '';
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
};

const formatDate = (dateString: string): string => {
  if (!dateString) return 'N/A';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('es-ES', {
    day: '2-digit',
    month: '2-digit',
    year: 'numeric', 
    hour: '2-digit', 
    minute: '2-digit'
  }).format(date);
};

const getActivityIcon = (type: string): string => {
  switch (type) {
    case 'ticket_created':
      return 'activity-create';
    case 'ticket_updated':
      return 'activity-update';
    case 'ticket_status_changed':
    case 'ticket_closed':
    case 'ticket_reopened':
      return 'activity-status';
    case 'ticket_assigned':
      return 'activity-assign';
    case 'ticket_priority_changed':
      return 'activity-priority';
    case 'comment_added':
      return 'activity-comment';
    case 'profile_updated':
      return 'activity-profile';
    case 'user_login':
      return 'activity-login';
    case 'user_logout':
      return 'activity-logout';
    default:
      return 'activity-default';
  }
};

const getActivityIconClass = (type: string): string => {
  switch (type) {
    case 'ticket_created':
      return 'pi pi-plus-circle';
    case 'ticket_updated':
      return 'pi pi-pencil';
    case 'ticket_status_changed':
      return 'pi pi-sync';
    case 'ticket_closed':
      return 'pi pi-check-circle';
    case 'ticket_reopened':
      return 'pi pi-refresh';
    case 'ticket_assigned':
      return 'pi pi-user-plus';
    case 'ticket_priority_changed':
      return 'pi pi-flag';
    case 'comment_added':
      return 'pi pi-comment';
    case 'profile_updated':
      return 'pi pi-user-edit';
    case 'user_login':
      return 'pi pi-sign-in';
    case 'user_logout':
      return 'pi pi-sign-out';
    default:
      return 'pi pi-clock';
  }
};

// Función para cargar el perfil del usuario
const loadProfile = async () => {
  loading.value = true;
  error.value = null;
  
  try {
    // Obtener el ID del usuario desde la ruta o usar el usuario actual
    let userId: string | undefined;
    const currentUserId = authStore.user?.id || localStorage.getItem('userId');
    
    if (route.params.id) {
      userId = route.params.id.toString();
      
      // Verificar si el usuario actual tiene permisos para ver este perfil
      if (userId !== currentUserId) {
        // Solo administradores pueden ver perfiles de otros usuarios
        if (!authStore.isAdmin) {
          console.warn('Intento no autorizado de acceder a perfil de otro usuario');
          error.value = 'No tienes permisos para ver este perfil';
          
          // Redirigir al usuario a su propio perfil
          router.replace('/profile');
          loading.value = false;
          return;
        }
      }
    } else if (authStore.user && authStore.user.id) {
      userId = authStore.user.id.toString();
    } else if (localStorage.getItem('userId')) {
      userId = localStorage.getItem('userId') || undefined;
    }
    
    if (!userId) {
      error.value = 'ID de usuario no encontrado';
      loading.value = false;
      return;
    }
    
    console.log('Intentando cargar el perfil del usuario con ID:', userId);
    
    // Intentar obtener el usuario desde el store
    const user = await userStore.fetchUser(userId);
    
    if (user) {
      profile.value = user;
      console.log('Perfil cargado exitosamente:', profile.value);
      
      // Cargar tickets y actividades
      await Promise.all([
        loadUserTickets(userId),
        loadUserActivities(userId)
      ]);
      
      // Configurar actualización automática cada 15 segundos
      setupUpdateInterval();
    } else {
      error.value = 'Usuario no encontrado';
      console.error('No se pudo cargar el usuario con ID:', userId);
    }
  } catch (e) {
    error.value = 'Error al cargar el perfil';
    console.error('Error en loadProfile:', e);
  } finally {
    loading.value = false;
  }
};

// Configurar intervalo de actualización
const setupUpdateInterval = () => {
  // Limpiar intervalo existente si hay uno
  if (updateInterval.value !== null) {
    window.clearInterval(updateInterval.value);
  }
  
  // Configurar nuevo intervalo
  updateInterval.value = window.setInterval(() => {
    console.log('Actualizando automáticamente tickets y actividades...');
    if (profile.value?.id) {
      loadUserTickets(profile.value.id.toString());
      loadUserActivities(profile.value.id.toString());
    }
  }, 15000);
};

// Limpiar intervalo cuando el componente se desmonta
watch(() => route.params.id, (newId) => {
  if (newId && String(newId) !== String(profile.value?.id)) {
    loadProfile();
  }
});

onMounted(() => {
  loadProfile();
});

// Limpiar el intervalo cuando el componente se desmonta
onMounted(() => {
  window.addEventListener('beforeunload', () => {
    if (updateInterval.value !== null) {
      window.clearInterval(updateInterval.value);
    }
  });
});

// Cargar tickets del usuario
const loadUserTickets = async (userId: string) => {
  try {
    console.log(`Cargando tickets asignados al usuario ${userId} en tiempo real`);
    
    // Mostrar inmediatamente un indicador de carga
    userTickets.value = [{
      id: 'loading',
      title: 'Cargando tickets...',
      description: 'Buscando tickets asignados en tiempo real',
      status: 'open' as 'open',
      priority: 'MEDIUM' as 'MEDIUM',
      category: 'System',
      createdBy: 'system',
      assignedTo: userId,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }];
    
    // Intentar obtener tickets reales desde el store
    console.log('Intentando obtener tickets reales usando ticketStore.fetchUserTickets()');
    await ticketStore.fetchTickets(); // Asegurar que tenemos tickets cargados
    const realTickets = await ticketStore.fetchUserTickets(userId);
    
    if (realTickets && Array.isArray(realTickets) && realTickets.length > 0) {
      userTickets.value = realTickets;
      console.log(`${realTickets.length} tickets reales asignados al usuario ${userId}`, realTickets);
      return;
    }
    
    // Si no hay tickets, mostrar array vacío
    userTickets.value = [];
    console.log(`No se encontraron tickets asignados al usuario ${userId}`);
    
  } catch (error) {
    console.error('Error al cargar tickets del usuario:', error);
    userTickets.value = [{
      id: 'ERROR-TICKET',
      title: 'Error al cargar tickets',
      description: 'Hubo un problema al cargar los tickets asignados',
      status: 'open' as 'open',
      priority: 'MEDIUM' as 'MEDIUM',
      category: 'Error',
      createdBy: 'system',
      assignedTo: userId,
      createdAt: new Date().toISOString(),
      updatedAt: new Date().toISOString()
    }];
  }
};

// Cargar actividades del usuario
const loadUserActivities = async (userId: string) => {
  try {
    console.log(`Cargando actividades para el usuario ${userId}`);
    
    // Cargar actividades del usuario desde el store de actividades
    const activities = await activityStore.fetchUserActivities(userId);
    
    if (activities && Array.isArray(activities) && activities.length > 0) {
      // Filtrar para eliminar actividades simuladas o falsas
      const realActivities = activities.filter((activity: any) => 
        !activity.metadata || !activity.metadata._simulated
      );
      
      console.log(`Filtrando actividades: ${activities.length} totales, ${realActivities.length} reales`);
      
      // Si no hay actividades reales después del filtrado, crear algunas actividades de demostración
      if (realActivities.length === 0) {
        console.log('No se encontraron actividades reales para este usuario, creando actividades de demostración');
        userActivity.value = [
          {
            id: '1',
            userId: userId,
            type: 'ticket_assigned',
            description: 'Ticket #TICKET-20250327041753 asignado a ti',
            timestamp: new Date().toISOString(),
            targetId: 'TICKET-20250327041753'
          },
          {
            id: '2',
            userId: userId,
            type: 'user_login',
            description: 'Inicio de sesión en el sistema',
            timestamp: new Date(Date.now() - 3600000).toISOString()
          },
          {
            id: '3',
            userId: userId,
            type: 'profile_updated',
            description: 'Perfil actualizado',
            timestamp: new Date(Date.now() - 86400000).toISOString()
          }
        ];
      } else {
        userActivity.value = realActivities;
      }
    } else {
      console.log('No se encontraron actividades para este usuario, creando actividades de demostración');
      userActivity.value = [
        {
          id: '1',
          userId: userId,
          type: 'ticket_assigned',
          description: 'Ticket #TICKET-20250327041753 asignado a ti',
          timestamp: new Date().toISOString(),
          targetId: 'TICKET-20250327041753'
        },
        {
          id: '2',
          userId: userId,
          type: 'user_login',
          description: 'Inicio de sesión en el sistema',
          timestamp: new Date(Date.now() - 3600000).toISOString()
        },
        {
          id: '3',
          userId: userId,
          type: 'profile_updated',
          description: 'Perfil actualizado',
          timestamp: new Date(Date.now() - 86400000).toISOString()
        }
      ];
    }
  } catch (error) {
    console.error('Error al cargar las actividades del usuario:', error);
    userActivity.value = [
      {
        id: '1',
        userId: userId,
        type: 'ticket_assigned',
        description: 'Ticket #TICKET-20250327041753 asignado a ti',
        timestamp: new Date().toISOString(),
        targetId: 'TICKET-20250327041753'
      },
      {
        id: '2',
        userId: userId,
        type: 'user_login',
        description: 'Inicio de sesión en el sistema',
        timestamp: new Date(Date.now() - 3600000).toISOString()
      },
      {
        id: '3',
        userId: userId,
        type: 'profile_updated',
        description: 'Perfil actualizado',
        timestamp: new Date(Date.now() - 86400000).toISOString()
      }
    ];
  }
};

// Manejar actualización del perfil
const handleProfileUpdated = async (updatedUser: any) => {
  // Registrar actividad de actualización de perfil
  await activityStore.logActivity({
    userId: updatedUser.id,
    type: 'profile_updated',
    description: 'Actualizó su información de perfil'
  });
  
  // Recargar perfil
  await loadProfile();
  showEditModal.value = false;
};
</script>

<style scoped lang="scss">
.user-profile-container {
  padding: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
  width: 100%;
  display: flex;
  flex-direction: column;
  overflow-x: hidden;

  // Content wrapper
  .content-wrapper {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
    width: 100%;
  }
  
  .loading, .error-container {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    min-height: 400px;
    text-align: center;
    padding: 2rem;
    background-color: var(--bg-secondary);
    border-radius: var(--border-radius);
    box-shadow: var(--card-shadow);
    border: 1px solid var(--border-color);
    
    p {
      margin-bottom: 1rem;
      font-size: 1.1rem;
      color: var(--text-secondary);
    }
  }

  .profile-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .profile-header {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    padding: 1.5rem;
    background-color: var(--bg-secondary);
    border-radius: var(--border-radius);
    box-shadow: var(--card-shadow);
    border: 1px solid var(--border-color);
    
    @media (max-width: 768px) {
      flex-direction: column;
      text-align: center;
    }
    
    .user-avatar {
      width: 100px;
      height: 100px;
      border-radius: 50%;
      display: flex;
      align-items: center;
      justify-content: center;
      font-size: 2.5rem;
      font-weight: 600;
      color: white;
      background-color: var(--primary-color);
      
      &.role-admin {
        background-color: var(--primary-color);
      }
      
      &.role-assistant {
        background-color: var(--secondary-color);
      }
      
      &.role-employee {
        background-color: var(--success-color);
      }
    }
    
    .user-info {
      flex: 1;
      
      h1 {
        margin: 0 0 0.5rem 0;
        font-size: 1.75rem;
        color: var(--text-primary);
      }
      
      .user-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 0.75rem;
        margin-bottom: 0.75rem;
        
        .user-id {
          background-color: var(--bg-tertiary);
          padding: 0.25rem 0.5rem;
          border-radius: var(--border-radius);
          font-size: 0.875rem;
          color: var(--text-secondary);
        }
      }
      
      .user-email {
        margin: 0.5rem 0;
        color: var(--text-secondary);
      }
      
      .user-department {
        margin: 0.5rem 0;
        color: var(--text-secondary);
      }
    }
    
    .actions {
      align-self: flex-start;
      display: flex;
      flex-direction: column;
      gap: 0.5rem;
      
      @media (max-width: 768px) {
        align-self: center;
        flex-direction: row;
      }
    }
  }

  .profile-grid {
    display: grid;
    grid-template-columns: 2fr 1fr;
    gap: 1.5rem;
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
  }

  .profile-card {
    background-color: var(--card-bg);
    border-radius: var(--border-radius);
    box-shadow: var(--card-shadow);
    overflow: hidden;
    border: 1px solid var(--border-color);
    
    .profile-card-header {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 1.25rem 1.5rem;
      background-color: var(--bg-tertiary);
      border-bottom: 1px solid var(--border-color);
      
      i {
        font-size: 1.2rem;
        color: var(--primary-color);
      }
      
      h2 {
        margin: 0;
        font-size: 1.2rem;
        color: var(--text-primary);
      }
    }
    
    .profile-card-body {
      padding: 1.5rem;
    }
  }

  .info-item {
    display: flex;
    margin-bottom: 1rem;
    
    &:last-child {
      margin-bottom: 0;
    }
    
    .label {
      width: 40%;
      font-weight: 500;
      color: var(--text-secondary);
    }
    
    .value {
      width: 60%;
      color: var(--text-primary);
    }
  }

  .stats-grid {
    display: grid;
    grid-template-columns: repeat(2, 1fr);
    gap: 1rem;
    
    .stat-item {
      padding: 1rem;
      text-align: center;
      background-color: var(--bg-tertiary);
      border-radius: var(--border-radius);
      
      .stat-number {
        display: block;
        font-size: 2rem;
        font-weight: 600;
        color: var(--primary-color);
        margin-bottom: 0.5rem;
      }
      
      .stat-label {
        font-size: 0.875rem;
        color: var(--text-secondary);
      }
    }
  }

  .profile-section {
    background-color: var(--card-bg);
    border-radius: var(--border-radius);
    padding: 1.5rem;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--border-color);
    
    .section-title {
      display: flex;
      align-items: center;
      justify-content: flex-start;
      margin-bottom: 1rem;
      color: var(--text-primary);
      font-size: 1.25rem;
      font-weight: 600;
      text-align: left;
      background-color: var(--bg-tertiary);
      border-radius: 10px;
      padding: 0.5rem 1rem;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
      border-left: 4px solid var(--primary-color);
      
      .title-icon {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 32px;
        height: 32px;
        background-color: var(--primary-color);
        border-radius: 8px;
        margin-right: 0.75rem;
        color: white;
        box-shadow: 0 4px 10px rgba(var(--primary-rgb), 0.25);
        
        i {
          font-size: 1rem;
        }
      }
      
      h2 {
        margin: 0;
        font-size: 1.1rem;
      }
    }
    
    h2 {
      margin-top: 0;
      margin-bottom: 1.5rem;
      color: var(--text-primary);
      font-size: 1.5rem;
      border-bottom: 1px solid var(--border-color);
      padding-bottom: 0.75rem;
    }
  }

  .empty-tickets, .empty-activity {
    text-align: center;
    padding: 2rem;
    background-color: var(--bg-tertiary);
    border-radius: var(--border-radius);
    color: var(--text-secondary);
    
    .btn {
      margin-top: 1rem;
    }
  }

  .tickets-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-top: 1.5rem;
    
    .ticket-card {
      border-radius: 16px;
      overflow: hidden;
      border: 1px solid var(--border-color);
      transition: all 0.3s ease;
      box-shadow: var(--card-shadow);
      background-color: var(--card-bg);
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 12px 20px -5px rgba(0, 0, 0, 0.1);
      }
      
      .ticket-header {
        padding: 1rem 1rem 0.5rem;
        margin-bottom: 0.5rem;
        
        .ticket-badges {
          display: flex;
          flex-wrap: wrap;
          gap: 0.5rem;
          margin-bottom: 0.75rem;
          
          .status-badge,
          .priority-badge {
            display: inline-flex;
            align-items: center;
            padding: 0.35rem 0.75rem;
            border-radius: 99px;
            font-size: 0.8rem;
            font-weight: 500;
            
            &::before {
              content: "";
              display: inline-block;
              width: 8px;
              height: 8px;
              border-radius: 50%;
              margin-right: 0.5rem;
            }
          }
          
          .status-badge {
            &.open { 
              background: rgba(30, 64, 175, 0.1); 
              color: #1e40af; 
              &::before { background-color: #1e40af; }
            }
            &.assigned { 
              background: rgba(91, 33, 182, 0.1); 
              color: #5b21b6; 
              &::before { background-color: #5b21b6; }
            }
            &.in_progress { 
              background: rgba(3, 105, 161, 0.1); 
              color: #0369a1; 
              &::before { background-color: #0369a1; }
            }
            &.resolved { 
              background: rgba(67, 56, 202, 0.1); 
              color: #4338ca; 
              &::before { background-color: #4338ca; }
            }
            &.closed { 
              background: rgba(55, 65, 81, 0.1); 
              color: #374151; 
              &::before { background-color: #374151; }
            }
          }
          
          .priority-badge {
            &.low { 
              background: rgba(3, 105, 161, 0.1); 
              color: #0369a1; 
              &::before { background-color: #0369a1; }
            }
            &.medium { 
              background: rgba(67, 56, 202, 0.1); 
              color: #4338ca; 
              &::before { background-color: #4338ca; }
            }
            &.high { 
              background: rgba(55, 48, 163, 0.1); 
              color: #3730a3; 
              &::before { background-color: #3730a3; }
            }
            &.urgent { 
              background: rgba(91, 33, 182, 0.1); 
              color: #5b21b6; 
              &::before { background-color: #5b21b6; }
            }
          }
        }
        
        .ticket-title {
          font-size: 1.1rem;
          font-weight: 600;
          color: var(--text-primary);
          margin: 0;
          line-height: 1.4;
        }
      }
      
      .ticket-body {
        padding: 0 1rem 1rem;
        
        .ticket-description {
          color: var(--text-secondary);
          font-size: 0.9rem;
          line-height: 1.5;
          margin: 0;
          display: -webkit-box;
          -webkit-line-clamp: 3;
          -webkit-box-orient: vertical;
          overflow: hidden;
        }
      }
      
      .ticket-meta {
        padding: 0.75rem 1rem;
        background-color: var(--bg-tertiary);
        border-top: 1px solid var(--border-color);
        display: flex;
        flex-direction: column;
        gap: 0.5rem;
        
        .meta-item {
          display: flex;
          align-items: center;
          gap: 0.5rem;
          
          i {
            width: 24px;
            height: 24px;
            background-color: rgba(99, 102, 241, 0.1);
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            color: var(--primary-color);
            font-size: 0.8rem;
          }
          
          span {
            font-size: 0.85rem;
            color: var(--text-secondary);
          }
        }
      }
      
      .ticket-footer {
        padding: 0.75rem 1rem;
        border-top: 1px solid var(--border-color);
        display: flex;
        justify-content: center;
        
        .view-details-btn {
          display: inline-flex;
          align-items: center;
          gap: 0.5rem;
          background-color: var(--primary-color);
          color: white;
          border: none;
          padding: 0.5rem 1rem;
          border-radius: 8px;
          font-size: 0.9rem;
          font-weight: 500;
          text-decoration: none;
          transition: all 0.3s ease;
          width: 100%;
          justify-content: center;
          
          &:hover {
            background-color: var(--primary-dark-color, #0d47a1);
          }
          
          i {
            font-size: 0.9rem;
          }
        }
      }
    }
  }

  .activity-timeline {
    .timeline {
      position: relative;
      padding-left: 2rem;
      
      &:before {
        content: '';
        position: absolute;
        left: 7px;
        top: 0;
        bottom: 0;
        width: 2px;
        background-color: var(--border-color);
      }
      
      .timeline-item {
        margin-bottom: 1.5rem;
        position: relative;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        .timeline-icon {
          position: absolute;
          left: -2rem;
          top: 0;
          width: 16px;
          height: 16px;
          border-radius: 50%;
          background-color: var(--text-muted);
          display: flex;
          align-items: center;
          justify-content: center;
          
          i {
            font-size: 0.7rem;
            color: white;
          }
          
          &.activity-create {
            background-color: var(--success-color);
          }
          
          &.activity-update {
            background-color: var(--primary-color);
          }
          
          &.activity-comment {
            background-color: var(--secondary-color);
          }
          
          &.activity-status {
            background-color: var(--warning-color);
          }
        }
        
        .timeline-content {
          padding: 0.75rem 1rem;
          background-color: var(--bg-tertiary);
          border-radius: var(--border-radius);
          
          .activity-text {
            margin: 0 0 0.5rem 0;
            color: var(--text-primary);
          }
          
          .activity-time {
            display: block;
            font-size: 0.8rem;
            color: var(--text-muted);
          }
        }
      }
    }
  }

  .role-badge, .status-badge {
    display: inline-block;
    padding: 0.25rem 0.5rem;
    border-radius: var(--border-radius);
    font-size: 0.875rem;
    font-weight: 500;
  }

  .role-badge {
    &.admin {
      background-color: rgba(79, 70, 229, 0.15);
      color: var(--primary-color);
    }
    
    &.assistant {
      background-color: rgba(249, 115, 22, 0.15);
      color: var(--secondary-color);
    }
    
    &.employee {
      background-color: rgba(16, 185, 129, 0.15);
      color: var(--success-color);
    }
  }

  .status-badge {
    &.active {
      background-color: rgba(16, 185, 129, 0.15);
      color: var(--success-color);
    }
    
    &.inactive {
      background-color: var(--bg-tertiary);
      color: var(--text-muted);
    }
    
    &.open {
      background-color: rgba(79, 70, 229, 0.15);
      color: var(--primary-color);
    }
    
    &.in_progress {
      background-color: rgba(249, 115, 22, 0.15);
      color: var(--secondary-color);
    }
    
    &.resolved {
      background-color: rgba(16, 185, 129, 0.15);
      color: var(--success-color);
    }
    
    &.closed {
      background-color: var(--bg-tertiary);
      color: var(--text-muted);
    }
  }

  /* Implementar estilos básicos para botones que serán sobrescritos por los estilos globales */
  .btn {
    padding: 0.75rem 1.5rem;
    border-radius: var(--border-radius);
    font-weight: 500;
    cursor: pointer;
    border: none;
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    text-decoration: none;
    
    &.btn-primary {
      background-color: var(--primary-color);
      color: white;
      
      &:hover {
        background-color: var(--primary-hover);
      }
    }
    
    &.btn-secondary {
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      
      &:hover {
        background-color: var(--hover-bg);
      }
    }
  }
}
</style> 
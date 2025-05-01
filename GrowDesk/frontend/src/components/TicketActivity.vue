<template>
  <div class="ticket-activity">
    <div class="activity-header">
      <h3>
        <i class="pi pi-history"></i>
        Historial de Actividades
      </h3>
    </div>
    
    <div v-if="isLoading" class="activity-loading">
      <i class="pi pi-spin pi-spinner"></i>
      <span>Cargando actividades...</span>
    </div>
    
    <div v-else-if="activities.length === 0" class="activity-empty">
      <i class="pi pi-inbox"></i>
      <p>No hay actividades registradas para este ticket</p>
    </div>
    
    <div v-else class="activity-timeline">
      <div v-for="activity in activities" :key="activity.id" class="activity-item">
        <div class="activity-icon" :class="getActivityIconClass(activity.type)">
          <i :class="getActivityIcon(activity.type)"></i>
        </div>
        
        <div class="activity-content">
          <div class="activity-user">
            {{ getUserName(activity.userId) }}
            <span class="activity-time">{{ formatTime(activity.timestamp) }}</span>
          </div>
          <div class="activity-description">{{ activity.description }}</div>
          
          <!-- Mostrar metadata adicional dependiendo del tipo de actividad -->
          <div v-if="activity.type === 'ticket_status_changed'" class="activity-metadata">
            <div class="status-change">
              <span class="old-status" :class="activity.metadata?.previousStatus">
                {{ translateStatus(activity.metadata?.previousStatus) }}
              </span>
              <i class="pi pi-arrow-right"></i>
              <span class="new-status" :class="activity.metadata?.newStatus">
                {{ translateStatus(activity.metadata?.newStatus) }}
              </span>
            </div>
          </div>
          
          <div v-if="activity.type === 'ticket_priority_changed'" class="activity-metadata">
            <div class="priority-change">
              <span class="old-priority" :class="activity.metadata?.previousPriority?.toLowerCase()">
                {{ translatePriority(activity.metadata?.previousPriority) }}
              </span>
              <i class="pi pi-arrow-right"></i>
              <span class="new-priority" :class="activity.metadata?.newPriority?.toLowerCase()">
                {{ translatePriority(activity.metadata?.newPriority) }}
              </span>
            </div>
          </div>
          
          <div v-if="activity.type === 'ticket_assigned'" class="activity-metadata">
            <div class="assignment-info">
              Asignado a: <strong>{{ getUserName(activity.metadata?.assignedToId) }}</strong>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch } from 'vue';
import { useActivityStore } from '@/stores/activity';
import { useUsersStore } from '@/stores/users';

// Props del componente
const props = defineProps({
  ticketId: {
    type: String,
    required: true
  }
});

// Stores
const activityStore = useActivityStore();
const usersStore = useUsersStore();

// Estado local
const isLoading = ref(true);
const activities = ref([]);

// Obtener actividades del ticket
const loadActivities = async () => {
  isLoading.value = true;
  
  try {
    // Obtener actividades del ticket
    const ticketActivities = activityStore.getTicketActivities(props.ticketId);
    
    if (ticketActivities && ticketActivities.length > 0) {
      // Ordenar actividades por fecha (más recientes primero)
      activities.value = [...ticketActivities].sort((a, b) => {
        return new Date(b.timestamp) - new Date(a.timestamp);
      });
    } else {
      // Si no hay actividades, intentar recargarlas
      await activityStore.fetchActivities();
      
      // Volver a obtener actividades del ticket
      const refreshedActivities = activityStore.getTicketActivities(props.ticketId);
      activities.value = [...refreshedActivities].sort((a, b) => {
        return new Date(b.timestamp) - new Date(a.timestamp);
      });
    }
    
    console.log(`Cargadas ${activities.value.length} actividades para el ticket ${props.ticketId}`);
  } catch (error) {
    console.error('Error al cargar actividades del ticket:', error);
  } finally {
    isLoading.value = false;
  }
};

// Obtener nombre de usuario
const getUserName = (userId) => {
  // Si no hay userId, mostrar 'Sistema'
  if (!userId) return 'Sistema';
  
  // Buscar usuario por ID
  const user = usersStore.getUserById(userId);
  
  // Si encontramos el usuario, mostrar su nombre completo
  if (user) {
    return `${user.firstName} ${user.lastName}`;
  }
  
  // Fallback: devolver el ID de usuario
  return userId;
};

// Obtener icono según tipo de actividad
const getActivityIcon = (activityType) => {
  switch (activityType) {
    case 'ticket_created':
      return 'pi pi-plus-circle';
    case 'ticket_updated':
      return 'pi pi-pencil';
    case 'ticket_closed':
      return 'pi pi-lock';
    case 'ticket_reopened':
      return 'pi pi-unlock';
    case 'ticket_assigned':
      return 'pi pi-user-plus';
    case 'ticket_status_changed':
      return 'pi pi-sync';
    case 'ticket_priority_changed':
      return 'pi pi-flag';
    case 'comment_added':
      return 'pi pi-comment';
    default:
      return 'pi pi-info-circle';
  }
};

// Obtener clase de icono según tipo de actividad
const getActivityIconClass = (activityType) => {
  switch (activityType) {
    case 'ticket_created':
      return 'activity-icon-created';
    case 'ticket_updated':
      return 'activity-icon-updated';
    case 'ticket_closed':
      return 'activity-icon-closed';
    case 'ticket_reopened':
      return 'activity-icon-reopened';
    case 'ticket_assigned':
      return 'activity-icon-assigned';
    case 'ticket_status_changed':
      return 'activity-icon-status';
    case 'ticket_priority_changed':
      return 'activity-icon-priority';
    case 'comment_added':
      return 'activity-icon-comment';
    default:
      return 'activity-icon-default';
  }
};

// Formatear tiempo
const formatTime = (timestamp) => {
  if (!timestamp) return '';
  
  const date = new Date(timestamp);
  
  // Si la fecha es de hoy, mostrar solo la hora
  const today = new Date();
  const isToday = today.toDateString() === date.toDateString();
  
  if (isToday) {
    return new Intl.DateTimeFormat('es', { 
      hour: '2-digit', 
      minute: '2-digit' 
    }).format(date);
  }
  
  // Para fechas anteriores, mostrar fecha completa
  return new Intl.DateTimeFormat('es', { 
    dateStyle: 'short',
    timeStyle: 'short'
  }).format(date);
};

// Traducir estado
const translateStatus = (status) => {
  if (!status) return '';
  
  const statusMap = {
    'open': 'Abierto',
    'assigned': 'Asignado',
    'in_progress': 'En Progreso',
    'resolved': 'Resuelto',
    'closed': 'Cerrado'
  };
  
  return statusMap[status] || status;
};

// Traducir prioridad
const translatePriority = (priority) => {
  if (!priority) return '';
  
  // Normalizar a minúsculas para la comparación
  const normalizedPriority = String(priority).toLowerCase();
  
  const priorityMap = {
    'low': 'Baja',
    'medium': 'Media',
    'high': 'Alta',
    'urgent': 'Urgente'
  };
  
  return priorityMap[normalizedPriority] || normalizedPriority;
};

// Cargar actividades al montar el componente
onMounted(() => {
  // Asegurarse de que hay usuarios cargados para mostrar nombres
  if (usersStore.users.length === 0) {
    usersStore.fetchUsers();
  }
  
  loadActivities();
});

// Observar cambios en el ticketId
watch(() => props.ticketId, (newTicketId) => {
  if (newTicketId) {
    loadActivities();
  }
});
</script>

<style lang="scss" scoped>
.ticket-activity {
  background-color: var(--bg-card);
  border-radius: 10px;
  padding: 1.5rem;
  box-shadow: var(--card-shadow);
  border: 1px solid var(--border-color);
  margin-bottom: 1.5rem;
  
  .activity-header {
    margin-bottom: 1.5rem;
    
    h3 {
      font-size: 1.2rem;
      color: var(--text-primary);
      display: flex;
      align-items: center;
      gap: 0.5rem;
      
      i {
        color: var(--primary-color);
      }
    }
  }
  
  .activity-loading, .activity-empty {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 2rem 1rem;
    text-align: center;
    color: var(--text-secondary);
    
    i {
      font-size: 2rem;
      margin-bottom: 1rem;
      color: var(--text-muted);
    }
  }
  
  .activity-timeline {
    position: relative;
    padding-left: 2rem;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      bottom: 0;
      left: 10px;
      width: 2px;
      background-color: var(--border-color);
    }
    
    .activity-item {
      position: relative;
      margin-bottom: 1.5rem;
      display: flex;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .activity-icon {
        position: absolute;
        left: -2rem;
        top: 0;
        width: 24px;
        height: 24px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        background-color: var(--primary-light);
        color: var(--primary-color);
        z-index: 10;
        
        i {
          font-size: 0.9rem;
        }
        
        &.activity-icon-created {
          background-color: #ecfdf5;
          color: #047857;
        }
        
        &.activity-icon-updated {
          background-color: #eff6ff;
          color: #1d4ed8;
        }
        
        &.activity-icon-closed {
          background-color: #f3f4f6;
          color: #4b5563;
        }
        
        &.activity-icon-reopened {
          background-color: #fef3c7;
          color: #b45309;
        }
        
        &.activity-icon-assigned {
          background-color: #f0f9ff;
          color: #0284c7;
        }
        
        &.activity-icon-status {
          background-color: #ede9fe;
          color: #6d28d9;
        }
        
        &.activity-icon-priority {
          background-color: #fff7ed;
          color: #ea580c;
        }
        
        &.activity-icon-comment {
          background-color: #fef2f2;
          color: #b91c1c;
        }
      }
      
      .activity-content {
        flex: 1;
        
        .activity-user {
          font-weight: 600;
          color: var(--text-primary);
          margin-bottom: 0.25rem;
          
          .activity-time {
            font-weight: normal;
            color: var(--text-muted);
            font-size: 0.85rem;
            margin-left: 0.5rem;
          }
        }
        
        .activity-description {
          color: var(--text-secondary);
          margin-bottom: 0.5rem;
        }
        
        .activity-metadata {
          background-color: var(--bg-tertiary);
          border-radius: 6px;
          padding: 0.75rem;
          font-size: 0.9rem;
          
          .status-change, .priority-change {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            
            .pi-arrow-right {
              color: var(--text-muted);
            }
            
            .old-status, .new-status, .old-priority, .new-priority {
              padding: 0.2rem 0.5rem;
              border-radius: 4px;
              font-weight: 600;
              font-size: 0.8rem;
            }
            
            // Estados
            .open, .assigned {
              background-color: #e0f2fe;
              color: #0369a1;
            }
            
            .in_progress {
              background-color: #f0f9ff;
              color: #0284c7;
            }
            
            .resolved {
              background-color: #ecfdf5;
              color: #047857;
            }
            
            .closed {
              background-color: #f3f4f6;
              color: #4b5563;
            }
            
            // Prioridades
            .low {
              background-color: #dbeafe;
              color: #0ea5e9;
            }
            
            .medium {
              background-color: #ede9fe;
              color: #4f46e5;
            }
            
            .high {
              background-color: #fff7ed;
              color: #ea580c;
            }
            
            .urgent {
              background-color: #fee2e2;
              color: #dc2626;
            }
          }
          
          .assignment-info {
            color: var(--text-secondary);
            
            strong {
              color: var(--primary-color);
            }
          }
        }
      }
    }
  }
}

// Responsive
@media (max-width: 768px) {
  .ticket-activity {
    padding: 1rem;
    
    .activity-timeline {
      padding-left: 1.5rem;
      
      &::before {
        left: 8px;
      }
      
      .activity-item {
        .activity-icon {
          left: -1.5rem;
          width: 20px;
          height: 20px;
          
          i {
            font-size: 0.8rem;
          }
        }
      }
    }
  }
}
</style> 
<template>
  <AdminLayout>
    <template #actions>
      <button class="btn btn-primary">
        <i class="pi pi-refresh"></i> Actualizar datos
      </button>
    </template>
    
    <!-- Contenido del dashboard -->
    <div class="admin-dashboard">
      <!-- Resumen de estadísticas -->
      <div class="stats-grid">
        <div class="stat-card">
          <div class="stat-icon users-icon">
            <i class="pi pi-users"></i>
          </div>
          <div class="stat-data">
            <h3 class="stat-value">{{ stats.totalUsers }}</h3>
            <p class="stat-label">Usuarios registrados</p>
          </div>
          <div class="stat-change" :class="stats.userChange > 0 ? 'positive' : 'negative'" v-if="stats.userChange !== 0">
            <i :class="['pi', stats.userChange > 0 ? 'pi-arrow-up' : 'pi-arrow-down']"></i>
            <span>{{ Math.abs(stats.userChange) }}%</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon tickets-icon">
            <i class="pi pi-ticket"></i>
          </div>
          <div class="stat-data">
            <h3 class="stat-value">{{ stats.totalTickets }}</h3>
            <p class="stat-label">Tickets activos</p>
          </div>
          <div class="stat-change" :class="stats.ticketChange > 0 ? 'positive' : 'negative'" v-if="stats.ticketChange !== 0">
            <i :class="['pi', stats.ticketChange > 0 ? 'pi-arrow-up' : 'pi-arrow-down']"></i>
            <span>{{ Math.abs(stats.ticketChange) }}%</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon response-icon">
            <i class="pi pi-clock"></i>
          </div>
          <div class="stat-data">
            <h3 class="stat-value">{{ stats.avgResponseTime }}</h3>
            <p class="stat-label">Tiempo de respuesta</p>
          </div>
          <div class="stat-change" :class="stats.responseChange < 0 ? 'positive' : 'negative'" v-if="stats.responseChange !== 0">
            <i :class="['pi', stats.responseChange < 0 ? 'pi-arrow-down' : 'pi-arrow-up']"></i>
            <span>{{ Math.abs(stats.responseChange) }}%</span>
          </div>
        </div>
        
        <div class="stat-card">
          <div class="stat-icon satisfaction-icon">
            <i class="pi pi-heart-fill"></i>
          </div>
          <div class="stat-data">
            <h3 class="stat-value">{{ stats.satisfaction }}%</h3>
            <p class="stat-label">Satisfacción</p>
          </div>
          <div class="stat-change" :class="stats.satisfactionChange > 0 ? 'positive' : 'negative'" v-if="stats.satisfactionChange !== 0">
            <i :class="['pi', stats.satisfactionChange > 0 ? 'pi-arrow-up' : 'pi-arrow-down']"></i>
            <span>{{ Math.abs(stats.satisfactionChange) }}%</span>
          </div>
        </div>
      </div>
      
      <!-- Gráficos y datos adicionales -->
      <div class="dashboard-grid">
        <!-- Actividad reciente -->
        <div class="dashboard-card">
          <div class="card-header">
            <h3 class="card-title">Actividad reciente</h3>
            <div class="card-actions">
              <button class="btn-icon">
                <i class="pi pi-ellipsis-h"></i>
              </button>
            </div>
          </div>
          <div class="card-body">
            <div class="activity-list">
              <div v-for="(activity, index) in recentActivity" :key="index" class="activity-item">
                <div class="activity-icon" :class="activity.type">
                  <i :class="getActivityIcon(activity.type)"></i>
                </div>
                <div class="activity-content">
                  <p class="activity-text">{{ activity.description }}</p>
                  <span class="activity-time">{{ activity.time }}</span>
                </div>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Tickets por categoría -->
        <div class="dashboard-card">
          <div class="card-header">
            <h3 class="card-title">Tickets por categoría</h3>
            <div class="card-actions">
              <button class="btn-icon">
                <i class="pi pi-ellipsis-h"></i>
              </button>
            </div>
          </div>
          <div class="card-body">
            <div class="category-list">
              <div v-for="(category, index) in categories" :key="index" class="category-item">
                <div class="category-info">
                  <span class="category-name">{{ category.name }}</span>
                  <span class="category-count">{{ category.count }}</span>
                </div>
                <div class="category-bar">
                  <div class="category-progress" :style="{ width: `${category.percentage}%` }"></div>
                </div>
                <span class="category-percentage">{{ category.percentage }}%</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref } from 'vue';
import AdminLayout from './AdminLayout.vue';

// Datos de ejemplo para el dashboard
const stats = ref({
  totalUsers: 458,
  userChange: 5.2,
  totalTickets: 127,
  ticketChange: -2.3,
  avgResponseTime: '2h 15m',
  responseChange: -12.5,
  satisfaction: 94.7,
  satisfactionChange: 1.8
});

// Actividad reciente
const recentActivity = ref([
  {
    type: 'user',
    description: 'Juan Pérez se registró como nuevo usuario',
    time: 'Hace 10 minutos'
  },
  {
    type: 'ticket',
    description: 'Ticket #2458 ha sido marcado como resuelto',
    time: 'Hace 25 minutos'
  },
  {
    type: 'message',
    description: 'María García respondió al ticket #2445',
    time: 'Hace 1 hora'
  },
  {
    type: 'alert',
    description: 'Se detectó un problema de rendimiento',
    time: 'Hace 2 horas'
  },
  {
    type: 'update',
    description: 'Se actualizó la configuración del sistema',
    time: 'Hace 3 horas'
  }
]);

// Categorías de tickets
const categories = ref([
  {
    name: 'Soporte técnico',
    count: 48,
    percentage: 37.8
  },
  {
    name: 'Consultas',
    count: 32,
    percentage: 25.2
  },
  {
    name: 'Facturación',
    count: 21,
    percentage: 16.5
  },
  {
    name: 'Errores',
    count: 15,
    percentage: 11.8
  },
  {
    name: 'Otros',
    count: 11,
    percentage: 8.7
  }
]);

// Función para obtener el icono de actividad según el tipo
const getActivityIcon = (type: string): string => {
  const iconMap: Record<string, string> = {
    user: 'pi pi-user',
    ticket: 'pi pi-ticket',
    message: 'pi pi-comment',
    alert: 'pi pi-exclamation-triangle',
    update: 'pi pi-sync'
  };
  
  return iconMap[type] || 'pi pi-check';
};
</script>

<style lang="scss" scoped>
.admin-dashboard {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
  width: 100%;
  
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    
    .stat-card {
      background-color: var(--bg-secondary);
      border-radius: 8px;
      padding: 1.5rem;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      display: flex;
      align-items: center;
      border: 1px solid var(--border-color);
      position: relative;
      overflow: hidden;
      
      .stat-icon {
        width: 56px;
        height: 56px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 1.25rem;
        
        i {
          font-size: 1.75rem;
          color: white;
        }
        
        &.users-icon {
          background-color: var(--primary-color);
        }
        
        &.tickets-icon {
          background-color: var(--info-color, #2196F3);
        }
        
        &.response-icon {
          background-color: var(--warning-color, #FF9800);
        }
        
        &.satisfaction-icon {
          background-color: var(--success-color, #4CAF50);
        }
      }
      
      .stat-data {
        flex-grow: 1;
        
        .stat-value {
          font-size: 1.75rem;
          font-weight: 700;
          margin: 0 0 0.25rem 0;
          color: var(--text-primary);
        }
        
        .stat-label {
          font-size: 0.9rem;
          color: var(--text-secondary);
          margin: 0;
          font-weight: 500;
        }
      }
      
      .stat-change {
        display: flex;
        align-items: center;
        gap: 0.25rem;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.8rem;
        font-weight: 600;
        position: absolute;
        top: 1rem;
        right: 1rem;
        
        &.positive {
          background-color: rgba(76, 175, 80, 0.15);
          color: var(--success-color, #4CAF50);
        }
        
        &.negative {
          background-color: rgba(244, 67, 54, 0.15);
          color: var(--error-color, #F44336);
        }
      }
    }
  }
  
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 1.5rem;
    
    .dashboard-card {
      background-color: var(--bg-secondary);
      border-radius: 8px;
      box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
      border: 1px solid var(--border-color);
      overflow: hidden;
      
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 1.25rem;
        background-color: var(--bg-tertiary);
        border-bottom: 1px solid var(--border-color);
        
        .card-title {
          font-size: 1.1rem;
          font-weight: 600;
          margin: 0;
          color: var(--text-primary);
        }
        
        .card-actions {
          .btn-icon {
            width: 32px;
            height: 32px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            background: none;
            border: none;
            color: var(--text-secondary);
            cursor: pointer;
            transition: all 0.2s;
            
            &:hover {
              background-color: var(--hover-bg);
              color: var(--text-primary);
            }
          }
        }
      }
      
      .card-body {
        padding: 1.25rem;
      }
      
      .activity-list {
        display: flex;
        flex-direction: column;
        gap: 1rem;
        
        .activity-item {
          display: flex;
          align-items: flex-start;
          
          .activity-icon {
            width: 36px;
            height: 36px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 1rem;
            flex-shrink: 0;
            
            i {
              font-size: 1rem;
              color: white;
            }
            
            &.user {
              background-color: var(--primary-color);
            }
            
            &.ticket {
              background-color: var(--info-color, #2196F3);
            }
            
            &.message {
              background-color: var(--success-color, #4CAF50);
            }
            
            &.alert {
              background-color: var(--error-color, #F44336);
            }
            
            &.update {
              background-color: var(--warning-color, #FF9800);
            }
          }
          
          .activity-content {
            flex-grow: 1;
            
            .activity-text {
              margin: 0 0 0.25rem 0;
              color: var(--text-primary);
              font-size: 0.9rem;
            }
            
            .activity-time {
              font-size: 0.75rem;
              color: var(--text-secondary);
            }
          }
        }
      }
      
      .category-list {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        
        .category-item {
          .category-info {
            display: flex;
            justify-content: space-between;
            margin-bottom: 0.5rem;
            
            .category-name {
              font-size: 0.9rem;
              color: var(--text-primary);
              font-weight: 500;
            }
            
            .category-count {
              font-size: 0.9rem;
              color: var(--text-secondary);
              font-weight: 600;
            }
          }
          
          .category-bar {
            height: 6px;
            background-color: var(--bg-tertiary);
            border-radius: 3px;
            overflow: hidden;
            margin-bottom: 0.25rem;
            
            .category-progress {
              height: 100%;
              background-color: var(--primary-color);
              border-radius: 3px;
            }
          }
          
          .category-percentage {
            font-size: 0.75rem;
            color: var(--text-secondary);
            text-align: right;
            display: block;
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .stats-grid,
  .dashboard-grid {
    grid-template-columns: 1fr;
  }
}
</style> 
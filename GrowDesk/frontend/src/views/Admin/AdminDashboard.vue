/* eslint-disable */
<template>
  <AdminLayout>
    <!-- Contenido del dashboard -->
    <div class="admin-section">
      <div class="admin-dashboard">
        <!-- Encabezado con fondo de color sólido y onda inferior -->
        <div class="hero-section">
          <div class="hero-content">
            <h1 class="hero-title">Panel de Administración</h1>
            <p class="hero-subtitle">Gestión y control del sistema de soporte</p>
          </div>
          <div class="wave-shape"></div>
        </div>
        
        <div class="content-wrapper">
          <!-- Sección de estadísticas principales -->
          <div class="metrics-section">
            <h2 class="section-title">
              <span class="title-icon"><i class="pi pi-chart-bar"></i></span>
              Estadísticas del Sistema
            </h2>
            
            <div class="section-header">
              <button @click="refreshData" class="view-all-btn">
                <i class="pi pi-refresh"></i>
                Actualizar datos
              </button>
            </div>
            
            <div class="stats-grid">
              <!-- Usuarios registrados (stat real) -->
              <div class="metric-card">
                <div class="card-icon">
                  <i class="pi pi-users"></i>
                </div>
                <div class="card-content">
                  <h3>Usuarios Registrados</h3>
                  <p class="number">{{ realStats.totalUsers }}</p>
                </div>
                <div v-if="stats.userChange !== 0" class="stat-change" :class="stats.userChange > 0 ? 'positive' : 'negative'">
                  <i :class="['pi', stats.userChange > 0 ? 'pi-arrow-up' : 'pi-arrow-down']"></i>
                  <span>{{ Math.abs(stats.userChange) }}%</span>
                </div>
              </div>
              
              <!-- Tickets activos (stat real) -->
              <div class="metric-card">
                <div class="card-icon tickets-icon">
                  <i class="pi pi-ticket"></i>
                </div>
                <div class="card-content">
                  <h3>Tickets Activos</h3>
                  <p class="number">{{ realStats.activeTickets }}</p>
                </div>
                <div v-if="stats.ticketChange !== 0" class="stat-change" :class="stats.ticketChange > 0 ? 'positive' : 'negative'">
                  <i :class="['pi', stats.ticketChange > 0 ? 'pi-arrow-up' : 'pi-arrow-down']"></i>
                  <span>{{ Math.abs(stats.ticketChange) }}%</span>
                </div>
              </div>
              
              <!-- Tiempo de respuesta (próximamente) -->
              <div class="metric-card">
                <div class="card-icon response-icon">
                  <i class="pi pi-clock"></i>
                </div>
                <div class="card-content">
                  <h3>Tiempo de Respuesta</h3>
                  <p class="number">-- hrs</p>
                  <p class="coming-soon">Próximamente</p>
                </div>
              </div>
              
              <!-- Satisfacción (próximamente) -->
              <div class="metric-card">
                <div class="card-icon satisfaction-icon">
                  <i class="pi pi-heart-fill"></i>
                </div>
                <div class="card-content">
                  <h3>Satisfacción</h3>
                  <p class="number">-- %</p>
                  <p class="coming-soon">Próximamente</p>
                </div>
              </div>
            </div>
          </div>
          
          <!-- Gráficos y datos adicionales -->
          <div class="dashboard-grid">
            <!-- Actividad reciente -->
            <div class="dashboard-card">
              <div class="card-header">
                <h3 class="card-title">Actividad Reciente</h3>
                <div class="card-actions">
                  <button class="btn-icon">
                    <i class="pi pi-ellipsis-h"></i>
                  </button>
                </div>
              </div>
              <div class="card-body">
                <div v-if="loadingActivity" class="loading-indicator">
                  <i class="pi pi-spin pi-spinner"></i>
                  <span>Cargando actividad...</span>
                </div>
                <div v-else class="activity-list">
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
                <h3 class="card-title">Tickets por Categoría</h3>
                <div class="card-actions">
                  <button class="btn-icon">
                    <i class="pi pi-ellipsis-h"></i>
                  </button>
                </div>
              </div>
              <div class="card-body">
                <div v-if="loadingCategories" class="loading-indicator">
                  <i class="pi pi-spin pi-spinner"></i>
                  <span>Cargando categorías...</span>
                </div>
                <div v-else class="category-list">
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
      </div>
    </div>
  </AdminLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import AdminLayout from './AdminLayout.vue';
import { useUsersStore } from '@/stores/users';
import { useTicketStore } from '@/stores/tickets';
import { useCategoriesStore } from '@/stores/categories';

// Referencias a los stores
const usersStore = useUsersStore();
const ticketStore = useTicketStore();
const categoriesStore = useCategoriesStore();

// Estados de carga
const loadingActivity = ref(false);
const loadingCategories = ref(false);

// Estadísticas reales calculadas a partir de los datos de los stores
const realStats = computed(() => ({
  totalUsers: usersStore.users.length,
  activeTickets: ticketStore.openTickets.length + ticketStore.inProgressTickets.length + ticketStore.assignedTickets.length
}));

// Datos de estadísticas complementarias
const stats = ref({
  userChange: 5.2,
  ticketChange: -2.3,
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

// Función para calcular tickets por categoría con datos reales
const getCategoriesWithTicketCounts = () => {
  const allTickets = ticketStore.tickets;
  const categoriesWithCounts = [];
  
  // Obtener las categorías disponibles
  const categoriesMap = new Map();
  
  // Primero añadir todas las categorías del store de categorías
  categoriesStore.categories.forEach(category => {
    categoriesMap.set(category.name, { 
      name: category.name, 
      count: 0,
      percentage: 0
    });
  });
  
  // Añadir categorías que aparecen en tickets pero no están en el store
  allTickets.forEach(ticket => {
    if (ticket.category && !categoriesMap.has(ticket.category)) {
      categoriesMap.set(ticket.category, { 
        name: ticket.category, 
        count: 0,
        percentage: 0
      });
    }
  });
  
  // Contar tickets por categoría
  allTickets.forEach(ticket => {
    const categoryName = ticket.category || 'Sin categoría';
    
    if (categoriesMap.has(categoryName)) {
      const category = categoriesMap.get(categoryName);
      category.count++;
    } else {
      categoriesMap.set(categoryName, { 
        name: categoryName, 
        count: 1,
        percentage: 0
      });
    }
  });
  
  // Convertir el Map a Array
  const categoriesArray = Array.from(categoriesMap.values());
  
  // Calcular el total de tickets
  const totalTickets = categoriesArray.reduce((total, category) => total + category.count, 0);
  
  // Calcular porcentajes
  categoriesArray.forEach(category => {
    category.percentage = totalTickets > 0 
      ? Math.round((category.count / totalTickets) * 100) 
      : 0;
  });
  
  // Ordenar por cantidad (descendente)
  return categoriesArray.sort((a, b) => b.count - a.count);
};

// Convertir a computed para mantener reactividad
const categories = computed(() => getCategoriesWithTicketCounts());

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

// Reemplazar la función refreshData existente para cargar todos los datos necesarios
const refreshData = async () => {
  loadingActivity.value = true;
  loadingCategories.value = true;
  
  try {
    // Cargar datos de usuarios y tickets actualizados
    await Promise.all([
      usersStore.fetchUsers(),
      ticketStore.fetchTickets(),
      categoriesStore.fetchCategories()
    ]);
    
    // Datos cargados correctamente
    loadingCategories.value = false;
  } catch (error) {
    console.error('Error al actualizar datos:', error);
  } finally {
    loadingActivity.value = false;
  }
};

// Cargar datos al montar el componente
onMounted(async () => {
  await refreshData();
});
</script>

<style lang="scss" scoped>
.admin-dashboard {
  --border-radius-lg: 1.25rem;
  --transition-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
  
  background-color: var(--bg-secondary);
  position: relative;
  overflow-x: hidden;
  
  // Sección hero con fondo de color sólido
  .hero-section {
    position: relative;
    padding: 2.5rem 2rem 6rem;
    background-color: var(--primary-color);
    color: white;
    text-align: center;
    overflow: hidden;
    
    .hero-content {
      position: relative;
      z-index: 2;
      max-width: 800px;
      margin: 0 auto;
      text-align: center;
    }
    
    .hero-title {
      font-size: 2.25rem;
      font-weight: 700;
      margin-bottom: 0.75rem;
      color: white;
      text-align: center;
    }
    
    .hero-subtitle {
      font-size: 1.1rem;
      margin-bottom: 0;
      opacity: 0.9;
    }
    
    // Forma ondulada en la parte inferior
    .wave-shape {
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 100%;
      height: 4rem;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' fill='%23f8fafc' opacity='.25'%3E%3C/path%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' fill='%23f8fafc' opacity='.5'%3E%3C/path%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%23f8fafc'%3E%3C/path%3E%3C/svg%3E");
      background-size: cover;
      background-position: center;
    }
  }
  
  .content-wrapper {
    max-width: 1300px;
    margin: 0 auto;
    padding: 3rem 1.5rem;
  }
  
  // Títulos de sección con iconos
  .section-title {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
    font-size: 1.5rem;
    font-weight: 600;
    text-align: left;
    background-color: var(--bg-tertiary);
    border-radius: 12px;
    padding: 0.75rem 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    border-left: 4px solid var(--primary-color);
    
    .title-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 38px;
      height: 38px;
      background-color: var(--primary-color);
      border-radius: 10px;
      margin-right: 1rem;
      color: white;
      box-shadow: 0 4px 10px rgba(var(--primary-rgb), 0.25);
      
      i {
        font-size: 1.2rem;
      }
    }
  }
  
  // Sección de header con botón de actualizar
  .section-header {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 1.25rem;
  }
  
  // Estadísticas
  .stats-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
    margin-bottom: 3rem;
    
    .metric-card {
      background-color: var(--card-bg);
      border-radius: var(--border-radius-lg);
      box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
      padding: 1.5rem;
      border: 1px solid var(--border-color);
      display: flex;
      align-items: center;
      transition: transform 0.3s var(--transition-bounce);
      position: relative;
      overflow: hidden;
      
      &:hover {
        transform: translateY(-8px);
        box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                   0 10px 10px -5px rgba(0, 0, 0, 0.04);
      }
      
      .card-icon {
        width: 64px;
        height: 64px;
        border-radius: 50%;
        background-color: rgba(var(--primary-color-rgb, 99, 102, 241), 0.1);
        display: flex;
        align-items: center;
        justify-content: center;
        margin-right: 1.25rem;
        
        i {
          font-size: 1.85rem;
          color: var(--primary-color);
        }
        
        &.tickets-icon i {
          color: #2196F3;
        }
        
        &.response-icon i {
          color: #FF9800;
        }
        
        &.satisfaction-icon i {
          color: #4CAF50;
        }
      }
      
      .card-content {
        flex: 1;
        
        h3 {
          margin: 0;
          color: var(--text-primary);
          font-size: 1rem;
          font-weight: 600;
        }
        
        .number {
          margin: 0.5rem 0 0;
          font-size: 2.5rem;
          font-weight: 700;
          color: var(--primary-color);
          letter-spacing: -0.5px;
        }
        
        .coming-soon {
          margin: 0.25rem 0 0;
          font-size: 0.875rem;
          color: var(--text-secondary);
          font-style: italic;
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
          color: #4CAF50;
        }
        
        &.negative {
          background-color: rgba(244, 67, 54, 0.15);
          color: #F44336;
        }
      }
    }
  }
  
  // Tarjetas del dashboard
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
    gap: 1.5rem;
    
    .dashboard-card {
      background-color: var(--card-bg);
      border-radius: var(--border-radius-lg);
      box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
      border: 1px solid var(--border-color);
      overflow: hidden;
      transition: transform 0.3s var(--transition-bounce);
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                   0 10px 10px -5px rgba(0, 0, 0, 0.04);
      }
      
      .card-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem 1.5rem;
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
        padding: 1.5rem;
      }
      
      .loading-indicator {
        display: flex;
        flex-direction: column;
        align-items: center;
        justify-content: center;
        padding: 2rem 0;
        color: var(--text-secondary);
        
        i {
          font-size: 2rem;
          margin-bottom: 0.75rem;
          color: var(--primary-color);
        }
        
        span {
          font-size: 0.9rem;
        }
      }
      
      .activity-list {
        display: flex;
        flex-direction: column;
        gap: 1.25rem;
        
        .activity-item {
          display: flex;
          align-items: flex-start;
          
          .activity-icon {
            width: 40px;
            height: 40px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            margin-right: 1rem;
            flex-shrink: 0;
            
            i {
              font-size: 1.1rem;
              color: white;
            }
            
            &.user {
              background-color: var(--primary-color);
            }
            
            &.ticket {
              background-color: #2196F3;
            }
            
            &.message {
              background-color: #4CAF50;
            }
            
            &.alert {
              background-color: #F44336;
            }
            
            &.update {
              background-color: #FF9800;
            }
          }
          
          .activity-content {
            flex-grow: 1;
            
            .activity-text {
              margin: 0 0 0.25rem 0;
              color: var(--text-primary);
              font-size: 0.95rem;
              font-weight: 500;
            }
            
            .activity-time {
              font-size: 0.8rem;
              color: var(--text-secondary);
            }
          }
        }
      }
      
      .category-list {
        display: flex;
        flex-direction: column;
        gap: 1.5rem;
        
        .category-item {
          .category-info {
            display: flex;
            justify-content: space-between;
            margin-bottom: 0.5rem;
            
            .category-name {
              font-size: 0.95rem;
              color: var(--text-primary);
              font-weight: 500;
            }
            
            .category-count {
              font-size: 0.95rem;
              color: var(--text-secondary);
              font-weight: 600;
            }
          }
          
          .category-bar {
            height: 8px;
            background-color: var(--bg-tertiary);
            border-radius: 4px;
            overflow: hidden;
            margin-bottom: 0.25rem;
            
            .category-progress {
              height: 100%;
              background-color: var(--primary-color);
              border-radius: 4px;
            }
          }
          
          .category-percentage {
            font-size: 0.8rem;
            color: var(--text-secondary);
            text-align: right;
            display: block;
          }
        }
      }
    }
  }
  
  // Botón de actualización
  .refresh-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background-color: var(--primary-color);
    color: white;
    border: none;
    padding: 0.75rem 1.5rem;
    border-radius: 8px;
    font-size: 0.95rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.3s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    
    &:hover {
      transform: translateY(-2px);
      background-color: var(--primary-dark-color, #0d47a1);
      box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15);
    }
    
    i {
      font-size: 1rem;
    }
  }
  
  // Nuevo estilo para el botón de actualizar datos
  .view-all-btn {
    display: inline-flex;
    align-items: center;
    gap: 0.5rem;
    background-color: var(--bg-tertiary);
    color: var(--primary-color);
    border: 1px solid var(--border-color);
    padding: 0.65rem 1.25rem;
    border-radius: 8px;
    font-size: 0.95rem;
    font-weight: 500;
    cursor: pointer;
    text-decoration: none;
    transition: all 0.2s ease;
    
    &:hover {
      background-color: var(--card-bg);
      transform: translateY(-2px);
    }
    
    i {
      font-size: 1rem;
    }
  }
  
  // Responsive adjustments
  @media (max-width: 768px) {
    .hero-section {
      padding: 2rem 1rem 4rem;
      
      .hero-title {
        font-size: 2rem;
      }
      
      .hero-subtitle {
        font-size: 1rem;
      }
    }
    
    .section-title {
      font-size: 1.5rem;
      flex-direction: column;
      
      .title-icon {
        margin-right: 0;
        margin-bottom: 0.75rem;
      }
    }
    
    .stats-grid,
    .dashboard-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
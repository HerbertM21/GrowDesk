/* eslint-disable */
<template>
  <div class="admin-section">
    <div class="dashboard">
      <!-- Encabezado con fondo de color sólido -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">Panel de Control</h1>
          <p class="hero-subtitle">Vista general del sistema de soporte</p>
        </div>
        <div class="wave-shape"></div>
      </div>
      
      <div class="content-wrapper">
        <!-- Sección de tarjetas de métricas con nuevo diseño -->
        <div class="metrics-section">
          <h2 class="section-title">
            <span class="title-icon"><i class="pi pi-chart-bar"></i></span>
            Métricas Principales
          </h2>
          <div class="dashboard-grid">
            <div class="metric-card">
              <div class="card-icon">
                <i class="pi pi-ticket"></i>
              </div>
              <div class="card-content">
                <h3>Tickets Abiertos</h3>
                <p class="number">{{ openTickets.length }}</p>
              </div>
            </div>
            
            <div class="metric-card">
              <div class="card-icon">
                <i class="pi pi-user"></i>
              </div>
              <div class="card-content">
                <h3>Tickets Asignados</h3>
                <p class="number">{{ assignedTickets.length }}</p>
              </div>
            </div>
            
            <div class="metric-card">
              <div class="card-icon">
                <i class="pi pi-check-circle"></i>
              </div>
              <div class="card-content">
                <h3>Tickets Cerrados</h3>
                <p class="number">{{ closedTickets.length }}</p>
              </div>
            </div>
            
            <div class="metric-card">
              <div class="card-icon urgent">
                <i class="pi pi-exclamation-triangle"></i>
              </div>
              <div class="card-content">
                <h3>Tickets Urgentes</h3>
                <p class="number">{{ urgentTickets.length }}</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Sección de métricas de rendimiento con nuevo diseño -->
        <div class="performance-section">
          <h2 class="section-title">
            <span class="title-icon"><i class="pi pi-clock"></i></span>
            Indicadores de Rendimiento
          </h2>
          <div class="dashboard-grid">
            <div class="performance-card">
              <div class="card-content">
                <h3>Tiempo Promedio de Resolución</h3>
                <p class="number">-- hrs</p>
                <p class="coming-soon">Próximamente</p>
              </div>
            </div>
            
            <div class="performance-card">
              <div class="card-content">
                <h3>Satisfacción del Cliente</h3>
                <p class="number">-- %</p>
                <p class="coming-soon">Próximamente</p>
              </div>
            </div>
            
            <div class="performance-card">
              <div class="card-content">
                <h3>Tasa de Resolución</h3>
                <p class="number">-- %</p>
                <p class="coming-soon">Próximamente</p>
              </div>
            </div>
            
            <div class="performance-card">
              <div class="card-content">
                <h3>Tickets por día</h3>
                <p class="number">{{ averageTicketsPerDay }}</p>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Sección de lista de tickets recientes con nuevo diseño -->
        <div class="recent-tickets-section">
          <h2 class="section-title">
            <span class="title-icon"><i class="pi pi-list"></i></span>
            Tickets Recientes
          </h2>
          
          <div class="section-header">
            <div class="section-actions">
              <router-link to="/tickets" class="view-all-btn">
                <i class="pi pi-external-link"></i>
                Ver todos los tickets
              </router-link>
            </div>
          </div>
          
          <div v-if="loading" class="status-message loading">
            <i class="pi pi-spin pi-spinner"></i>
            <p>Cargando tickets...</p>
          </div>
          
          <div v-else-if="tickets.length === 0" class="status-message empty">
            <i class="pi pi-inbox"></i>
            <p>No hay tickets recientes</p>
          </div>
          
          <div v-else class="tickets-table">
            <div class="table-header">
              <div class="column">ID</div>
              <div class="column title">Título</div>
              <div class="column">Estado</div>
              <div class="column">Prioridad</div>
              <div class="column">Asignado a</div>
              <div class="column actions">Acciones</div>
            </div>
            
            <div v-for="ticket in recentTickets" :key="ticket.id" class="table-row">
              <div class="column">{{ ticket.id.split('-')[1] || ticket.id }}</div>
              <div class="column title">{{ ticket.title }}</div>
              <div class="column">
                <span :class="['status-badge', ticket.status]">{{ translateStatus(ticket.status) }}</span>
              </div>
              <div class="column">
                <span :class="['priority-badge', normalizePriority(ticket.priority)]">
                  {{ translatePriority(ticket.priority) }}
                </span>
              </div>
              <div class="column">{{ getUserFullName(ticket.assignedTo) }}</div>
              <div class="column actions">
                <router-link :to="`/tickets/${ticket.id}`" class="action-btn">
                  <i class="pi pi-eye"></i> Ver
                </router-link>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, computed } from 'vue'
import { useTicketStore } from '@/stores/tickets'
import { useUsersStore } from '@/stores/users'
import type { Ticket } from '@/stores/tickets'

const ticketStore = useTicketStore()
const usersStore = useUsersStore()

onMounted(async () => {
  await ticketStore.fetchTickets()
  await usersStore.fetchUsers()
})

// Obtener los getters del store
const openTickets = computed(() => ticketStore.openTickets)
const assignedTickets = computed(() => ticketStore.assignedTickets)
const closedTickets = computed(() => ticketStore.closedTickets)
const urgentTickets = computed(() => ticketStore.urgentTickets)
const ticketsPerDay = computed(() => ticketStore.ticketsPerDay)

const tickets = computed(() => ticketStore.tickets)
const loading = computed(() => ticketStore.loading)

// Obtener los 5 tickets más recientes
const recentTickets = computed(() => {
  return [...tickets.value]
    .sort((a, b) => new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime())
    .slice(0, 5)
})

// Función para traducir el estado del ticket
const translateStatus = (status: string): string => {
  const statusMap: Record<string, string> = {
    'open': 'Abierto',
    'assigned': 'Asignado',
    'in_progress': 'En Progreso',
    'resolved': 'Resuelto',
    'closed': 'Cerrado'
  }
  return statusMap[status] || status
}

// Función para traducir la prioridad del ticket
const translatePriority = (priority: string): string => {
  const priorityMap: Record<string, string> = {
    'low': 'Baja',
    'medium': 'Media',
    'high': 'Alta',
    'urgent': 'Urgente',
    'LOW': 'Baja',
    'MEDIUM': 'Media',
    'HIGH': 'Alta',
    'URGENT': 'Urgente'
  }
  return priorityMap[priority] || priority
}

// Función para normalizar el valor de prioridad
const normalizePriority = (priority: string): string => {
  if (!priority) return 'medium';
  
  // Convertir a minúsculas
  const lowerPriority = priority.toLowerCase();
  
  // Mapear diferentes formatos posibles al formato estándar
  if (lowerPriority === 'baja' || lowerPriority === 'low') return 'low';
  if (lowerPriority === 'media' || lowerPriority === 'medium') return 'medium';
  if (lowerPriority === 'alta' || lowerPriority === 'high') return 'high';
  if (lowerPriority === 'urgente' || lowerPriority === 'urgent') return 'urgent';
  
  // Si no coincide con ninguno, devolver medium por defecto
  return 'medium';
}

// Función para obtener el nombre completo de un usuario por su ID
const getUserFullName = (userId: string | null): string => {
  if (!userId) return 'Sin asignar'
  
  const user = usersStore.users.find((user: any) => user.id === userId)
  if (user) {
    return `${user.firstName} ${user.lastName}`
  }
  
  return userId // Si no encuentra el usuario, muestra el ID como fallback
}

// Función para formatear la fecha
const formatDate = (dateString: string): string => {
  const date = new Date(dateString)
  return date.toLocaleDateString('es-CL', { 
    day: '2-digit',
    month: '2-digit'
  })
}

// Añadir después de ticketsPerDay computed
const averageTicketsPerDay = computed(() => {
  // Obtener todos los valores de tickets por día
  const ticketsPerDayValues = ticketsPerDay.value.map(([_, count]) => count);
  
  // Si no hay datos, devolver 0
  if (ticketsPerDayValues.length === 0) return 0;
  
  // Calcular la suma total de tickets
  const totalTickets = ticketsPerDayValues.reduce((sum, count) => sum + count, 0);
  
  // Calcular y redondear el promedio
  return Math.round(totalTickets / ticketsPerDayValues.length);
});
</script>

<style lang="scss" scoped>
.dashboard {
  --border-radius-lg: 1.25rem;
  --transition-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
  
  background-color: var(--bg-secondary);
  position: relative;
  overflow-x: hidden;
  
  // Sección hero con diseño redondeado
  .hero-section {
    position: relative;
    padding: 3.5rem 1rem 4rem;
    background-color: var(--primary-color);
    color: white;
    text-align: center;
    overflow: hidden;
    border: none;
    outline: none;
    border-radius: 0;
    margin-bottom: 0;
    box-shadow: none;
    
    &::before {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background-image: radial-gradient(circle at 20% 30%, rgba(255, 255, 255, 0.05) 0%, transparent 30%),
                        radial-gradient(circle at 80% 70%, rgba(255, 255, 255, 0.05) 0%, transparent 25%);
    }
    
    &::after {
      content: "";
      position: absolute;
      top: 0;
      left: 0;
      width: 100%;
      height: 100%;
      background: repeating-linear-gradient(
        -45deg,
        rgba(255, 255, 255, 0.03),
        rgba(255, 255, 255, 0.03) 1px,
        transparent 1px,
        transparent 8px
      );
    }
    
    // Forma ondulada en la parte inferior
    .wave-shape {
      position: absolute;
      bottom: -1px;
      left: 0;
      width: 100%;
      height: 25px;
      border: none;
      outline: none;
      z-index: 10;
      background: var(--bg-secondary);
      clip-path: ellipse(60% 90% at 50% 100%);
    }
    
    .hero-content {
      position: relative;
      z-index: 5;
      max-width: 800px;
      margin: 0 auto;
      text-align: center;
      
      // Círculos decorativos de diferentes tamaños
      &::before {
        content: "";
        position: absolute;
        left: 50%;
        top: -25px;
        width: 120px;
        height: 120px;
        border-radius: 50%;
        background-color: rgba(255, 255, 255, 0.06);
        transform: translateX(-110px);
        z-index: -1;
      }
      
      &::after {
        content: "";
        position: absolute;
        right: 50%;
        bottom: -40px;
        width: 90px;
        height: 90px;
        border-radius: 50%;
        background-color: rgba(255, 255, 255, 0.06);
        transform: translateX(120px);
        z-index: -1;
      }
    }
    
    .hero-title {
      font-size: 2.5rem;
      font-weight: 800;
      margin-bottom: 1.5rem;
      color: white;
      position: relative;
      text-align: center;
      display: inline-block;
      padding: 0 20px;
      
      &::before {
        content: "{ ";
        position: absolute;
        left: -15px;
        top: 0;
        font-family: monospace;
        opacity: 0.5;
        font-size: 2.5rem;
      }
      
      &::after {
        content: " }";
        position: absolute;
        right: -15px;
        top: 0;
        font-family: monospace;
        opacity: 0.5;
        font-size: 2.5rem;
      }
    }
    
    .hero-subtitle {
      font-size: 0.9rem;
      margin: 0 auto 1.75rem;
      opacity: 0.9;
      max-width: 600px;
      position: relative;
      display: block;
      padding: 8px 15px;
      border-radius: 30px;
      background-color: rgba(255, 255, 255, 0.1);
      backdrop-filter: blur(5px);
      box-shadow: 0 4px 15px rgba(0, 0, 0, 0.1);
      
      // Círculos adicionales en el subtítulo
      &::before {
        content: "";
        position: absolute;
        right: -15px;
        top: -10px;
        width: 30px;
        height: 30px;
        border-radius: 50%;
        background-color: rgba(255, 255, 255, 0.04);
        z-index: -1;
      }
      
      &::after {
        content: "";
        position: absolute;
        left: -12px;
        bottom: -8px;
        width: 25px;
        height: 25px;
        border-radius: 50%;
        background-color: rgba(255, 255, 255, 0.04);
        z-index: -1;
      }
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
  
  // Contenedores de secciones
  .metrics-section,
  .performance-section,
  .recent-tickets-section {
    margin-bottom: 3.5rem;
  }
  
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: 1.5rem;
    margin-bottom: 1rem;
  }
  
  // Tarjetas de métricas
  .metric-card {
    background-color: var(--card-bg);
    border-radius: 24px; 
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
      background-color: rgba(var(--primary-color-rgb), 0.1);
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 1.25rem;
      
      i {
        font-size: 1.85rem;
        color: var(--primary-color);
      }
      
      &.urgent {
        background-color: rgba(220, 38, 38, 0.1);
        
        i {
          color: #dc2626;
        }
      }
    }
    
    .card-content {
      flex: 1;
      display: flex; /* Añadido para centrar verticalmente */
      flex-direction: column;
      justify-content: center; /* Centrado vertical */
      
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
    }
  }
  
  // Tarjetas de rendimiento
  .performance-card {
    background-color: var(--card-bg);
    border-radius: 24px; 
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    padding: 1.75rem;
    border: 1px solid var(--border-color);
    transition: transform 0.3s var(--transition-bounce);
    display: flex; /* Añadido para centrar verticalmente */
    flex-direction: column;
    
    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                 0 10px 10px -5px rgba(0, 0, 0, 0.04);
    }
    
    .card-content {
      text-align: center;
      display: flex; /* Añadido para centrar verticalmente */
      flex-direction: column;
      justify-content: center; /* Centrado vertical */
      flex: 1; /* Para que ocupe todo el espacio disponible */
      
      h3 {
        margin: 0 0 1rem;
        color: var(--text-primary);
        font-size: 1.1rem;
        font-weight: 600;
      }
      
      .number {
        margin: 0;
        font-size: 2.5rem;
        font-weight: 700;
        color: var(--primary-color);
      }
      
      .coming-soon {
        margin: 0.5rem 0 0;
        font-size: 0.875rem;
        color: var(--text-secondary);
        font-style: italic;
      }
    }
  }
  
  // Sección de tickets recientes
  .section-header {
    display: flex;
    justify-content: flex-end;
    margin-bottom: 1rem;
    
    .section-actions {
      display: flex;
      gap: 1rem;
    }
  }
  
  .view-all-btn {
    display: flex;
    align-items: center;
    font-size: 0.9rem;
    color: var(--primary-color);
    text-decoration: none;
    transition: color 0.2s;
    
    i {
      margin-right: 0.5rem;
    }
    
    &:hover {
      color: var(--primary-dark-color);
    }
  }
  
  // Estado de carga, error, vacío
  .status-message {
    text-align: center;
    padding: 3rem;
    background-color: var(--card-bg);
    border-radius: var(--border-radius-lg);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;
    border: 1px solid var(--border-color);
    
    i {
      font-size: 3rem;
      margin-bottom: 1rem;
      display: block;
    }
    
    p {
      font-size: 1.1rem;
      color: var(--text-secondary);
    }
    
    &.loading i {
      color: var(--primary-color);
    }
    
    &.empty i {
      color: #6b7280;
    }
  }
  
  // Tabla de tickets
  .tickets-table {
    background-color: var(--card-bg);
    border-radius: var(--border-radius-lg);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    border: 1px solid var(--border-color);
    
    .table-header {
      display: flex;
      background-color: var(--bg-tertiary);
      padding: 1rem 1.5rem;
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .table-row {
      display: flex;
      padding: 1rem 1.5rem;
      border-bottom: 1px solid var(--border-color);
      transition: background-color 0.2s;
      
      &:last-child {
        border-bottom: none;
      }
      
      &:hover {
        background-color: var(--bg-tertiary);
      }
    }
    
    .column {
      flex: 1;
      display: flex;
      align-items: center;
      
      &.title {
        flex: 2;
        font-weight: 500;
      }
      
      &.actions {
        justify-content: flex-end;
      }
    }
    
    // Badges de estado y prioridad
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
    
    // Botón de acción
    .action-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background-color: var(--primary-color);
      color: white;
      border: none;
      padding: 0.5rem 1rem;
      border-radius: 8px;
      font-size: 0.85rem;
      font-weight: 600;
      text-decoration: none;
      transition: all 0.3s ease;
      box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
      
      &:hover {
        transform: translateY(-2px);
        background-color: var(--primary-dark-color, #0d47a1);
        box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15);
      }
      
      i {
        font-size: 0.85rem;
      }
    }
  }
  
  // Responsive adjustments
  @media (max-width: 768px) {
    .hero-section {
      padding: 2.5rem 1rem 5rem;
      
      .hero-title {
        font-size: 1.8rem;
        padding: 0 10px;
        margin-bottom: 1.2rem;
        
        &::before, &::after {
          font-size: 1.8rem;
        }
      }
    }
    
    .content-wrapper {
      padding: 2rem 1rem;
    }
    
    .dashboard-grid {
      grid-template-columns: 1fr;
    }
    
    .section-title {
      font-size: 1.3rem;
      
      .title-icon {
        width: 32px;
        height: 32px;
        
        i {
          font-size: 1rem;
        }
      }
    }
  }
}
</style>
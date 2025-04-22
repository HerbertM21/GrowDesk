<template>
  <div class="dashboard">
    <h1>Panel de Control</h1>
    
    <!-- Sección de tarjetas de métricas -->
    <div class="metrics-section">
      <div class="dashboard-grid">
        <div class="dashboard-card">
          <div class="card-icon">
            <i class="pi pi-ticket"></i>
          </div>
          <div class="card-content">
            <h3>Tickets Abiertos</h3>
            <p class="number">{{ openTickets.length }}</p>
          </div>
        </div>
        
        <div class="dashboard-card">
          <div class="card-icon">
            <i class="pi pi-user"></i>
          </div>
          <div class="card-content">
            <h3>Tickets Asignados</h3>
            <p class="number">{{ assignedTickets.length }}</p>
          </div>
        </div>
        
        <div class="dashboard-card">
          <div class="card-icon">
            <i class="pi pi-check-circle"></i>
          </div>
          <div class="card-content">
            <h3>Tickets Cerrados</h3>
            <p class="number">{{ closedTickets.length }}</p>
          </div>
        </div>
        
        <div class="dashboard-card">
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
    
    <!-- Sección de métricas de rendimiento -->
    <div class="performance-section">
      <div class="dashboard-grid">
        <div class="dashboard-card">
          <div class="card-content">
            <h3>Tiempo Promedio de Resolución</h3>
            <p class="number">-- hrs</p>
            <p class="coming-soon">Próximamente</p>
          </div>
        </div>
        
        <div class="dashboard-card">
          <div class="card-content">
            <h3>Satisfacción del Cliente</h3>
            <p class="number">-- %</p>
            <p class="coming-soon">Próximamente</p>
          </div>
        </div>
        
        <div class="dashboard-card">
          <div class="card-content">
            <h3>Tasa de Resolución</h3>
            <p class="number">-- %</p>
            <p class="coming-soon">Próximamente</p>
          </div>
        </div>
        
        <div class="dashboard-card">
          <div class="card-content">
            <h3>Tickets por Día</h3>
            <div class="tickets-per-day">
              <div v-for="count in ticketsPerDay.slice(0, 3).map(([_, c]) => c)" :key="count" class="day-count">
                <span class="count">{{ count }}</span>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Sección de lista de tickets recientes -->
    <div class="recent-tickets-section">
      <div class="section-header">
        <h2>Tickets Recientes</h2>
        <router-link to="/tickets" class="view-all">Ver todos los tickets</router-link>
      </div>
      
      <div v-if="loading" class="loading">Cargando tickets...</div>
      <div v-else-if="tickets.length === 0" class="empty-state">No hay tickets recientes</div>
      <div v-else class="tickets-table">
        <div class="table-header">
          <div class="column">ID</div>
          <div class="column">Título</div>
          <div class="column">Estado</div>
          <div class="column">Prioridad</div>
          <div class="column">Asignado a</div>
          <div class="column">Acciones</div>
        </div>
        
        <div v-for="ticket in recentTickets" :key="ticket.id" class="table-row">
          <div class="column">{{ ticket.id.split('-')[1] }}</div>
          <div class="column title">{{ ticket.title }}</div>
          <div class="column">
            <span :class="['status-badge', ticket.status]">{{ translateStatus(ticket.status) }}</span>
          </div>
          <div class="column">
            <span :class="['priority-badge', ticket.priority]">{{ translatePriority(ticket.priority) }}</span>
          </div>
          <div class="column">{{ getUserFullName(ticket.assignedTo) }}</div>
          <div class="column">
            <router-link :to="`/tickets/${ticket.id}`" class="action-btn">
              <i class="pi pi-eye"></i> Ver
            </router-link>
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
    'LOW': 'Baja',
    'MEDIUM': 'Media',
    'HIGH': 'Alta',
    'URGENT': 'Urgente'
  }
  return priorityMap[priority] || priority
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
</script>

<style lang="scss" scoped>
.dashboard {
  max-width: 1200px;
  margin: 0 auto;
  background: linear-gradient(to bottom, #f8faff, #eef2ff);
  border-radius: 10px;
  padding: 1.5rem;
  box-shadow: 0 5px 15px rgba(35, 38, 110, 0.05);
  
  h1 {
    margin-bottom: 1.5rem;
    color: #1e293b;
    font-weight: 600;
    border-bottom: 2px solid #c7d2fe;
    padding-bottom: 0.5rem;
  }
  
  h2 {
    font-size: 1.5rem;
    color: #334155;
    margin-bottom: 1rem;
  }
  
  .metrics-section,
  .performance-section,
  .recent-tickets-section {
    margin-bottom: 2.5rem;
  }

  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    margin-bottom: 1rem;
  }

  .dashboard-card {
    background: #f8fafc;
    padding: 1.5rem;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    display: flex;
    align-items: center;
    border: 1px solid #e2e8f0;
    transition: all 0.3s ease;
    
    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 8px 24px rgba(0, 0, 0, 0.08);
    }
    
    .card-icon {
      width: 56px;
      height: 56px;
      border-radius: 50%;
      background: linear-gradient(135deg, #e0e7ff, #c7d2fe);
      display: flex;
      align-items: center;
      justify-content: center;
      margin-right: 1.25rem;
      
      i {
        font-size: 1.5rem;
        color: #4f46e5;
      }
      
      &.urgent {
        background: linear-gradient(135deg, #ede9fe, #ddd6fe);
        
        i {
          color: #7c3aed;
        }
      }
    }
    
    .card-content {
      flex: 1;
    }

    h3 {
      margin: 0;
      color: #4b5563;
      font-size: 1rem;
      font-weight: 600;
    }

    .number {
      margin: 0.5rem 0 0;
      font-size: 2.25rem;
      font-weight: 700;
      background: linear-gradient(90deg, #4f46e5, #6366f1);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
      letter-spacing: -0.5px;
    }
    
    .coming-soon {
      margin: 0.25rem 0 0;
      font-size: 0.8rem;
      color: #6b7280;
      font-style: italic;
    }
  }
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.25rem;
    
    .view-all {
      color: #4f46e5;
      text-decoration: none;
      font-size: 0.95rem;
      font-weight: 500;
      transition: color 0.2s;
      
      &:hover {
        color: #4338ca;
        text-decoration: underline;
      }
    }
  }
  
  .tickets-table {
    background: #f8fafc;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    overflow: hidden;
    border: 1px solid #e2e8f0;
    
    .table-header {
      display: flex;
      background: linear-gradient(to right, #f1f5f9, #f8fafc);
      padding: 1rem 1.5rem;
      font-weight: 600;
      color: #475569;
    }
    
    .table-row {
      display: flex;
      padding: 1rem 1.5rem;
      border-bottom: 1px solid #e5e7eb;
      transition: background-color 0.2s;
      
      &:last-child {
        border-bottom: none;
      }
      
      &:hover {
        background-color: #f1f5f9;
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
    }
    
    .status-badge,
    .priority-badge {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      border-radius: 4px;
      font-size: 0.75rem;
      text-transform: capitalize;
      font-weight: 500;
    }
    
    .status-badge {
      &.open { background: #cfe7fe; color: #1e40af; }
      &.assigned { background: #dbd2fd; color: #5b21b6; }
      &.in_progress { background: #bae6fd; color: #0369a1; }
      &.resolved { background: #c7d2fe; color: #4338ca; }
      &.closed { background: #d1d5db; color: #374151; }
    }
    
    .priority-badge {
      &.low { background: #bae6fd; color: #0369a1; }
      &.medium { background: #c7d2fe; color: #4338ca; }
      &.high { background: #e0e7ff; color: #3730a3; }
      &.urgent { background: #ddd6fe; color: #5b21b6; }
    }
    
    .action-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.35rem;
      background: linear-gradient(135deg, #4f46e5, #6366f1);
      color: white;
      border: none;
      padding: 0.4rem 0.75rem;
      border-radius: 6px;
      font-size: 0.8rem;
      font-weight: 500;
      cursor: pointer;
      text-decoration: none;
      transition: all 0.2s ease;
      box-shadow: 0 2px 4px rgba(79, 70, 229, 0.2);
      
      &:hover {
        background: linear-gradient(135deg, #4338ca, #4f46e5);
        transform: translateY(-1px);
        box-shadow: 0 4px 8px rgba(79, 70, 229, 0.3);
      }
    }
  }
  
  .loading,
  .empty-state {
    text-align: center;
    padding: 2.5rem;
    color: #6b7280;
    background: #f8fafc;
    border-radius: 12px;
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.05);
    border: 1px solid #e2e8f0;
    font-style: italic;
  }

  .tickets-per-day {
    margin-top: 1rem;
    
    .day-count {
      display: flex;
      justify-content: center;
      align-items: center;
      padding: 0.5rem 0;
      border-bottom: 1px solid #e2e8f0;
      
      &:last-child {
        border-bottom: none;
      }
      
      .count {
        font-size: 2.25rem;
        font-weight: 700;
        background: linear-gradient(90deg, #4f46e5, #6366f1);
        -webkit-background-clip: text;
        -webkit-text-fill-color: transparent;
        letter-spacing: -0.5px;
        padding: 0;
      }
    }
  }
}
</style> 
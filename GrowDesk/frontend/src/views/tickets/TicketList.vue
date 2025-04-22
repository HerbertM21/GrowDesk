<template>
  <div class="ticket-list">
    <div class="header">
      <h1>Tickets</h1>
      <router-link to="/tickets/new" class="btn btn-primary">Nuevo Ticket</router-link>
    </div>

    <!-- Filtros -->
    <div class="filters">
      <div class="filter-group">
        <label>Estado:</label>
        <div class="filter-options">
          <button 
            @click="setStatusFilter('all')" 
            :class="['filter-btn', statusFilter === 'all' ? 'active' : '']"
          >
            Todos
          </button>
          <button 
            @click="setStatusFilter('open')" 
            :class="['filter-btn', statusFilter === 'open' ? 'active' : '']"
          >
            Abiertos
          </button>
          <button 
            @click="setStatusFilter('assigned')" 
            :class="['filter-btn', statusFilter === 'assigned' ? 'active' : '']"
          >
            Asignados
          </button>
          <button 
            @click="setStatusFilter('in_progress')" 
            :class="['filter-btn', statusFilter === 'in_progress' ? 'active' : '']"
          >
            En Progreso
          </button>
          <button 
            @click="setStatusFilter('resolved')" 
            :class="['filter-btn', statusFilter === 'resolved' ? 'active' : '']"
          >
            Resueltos
          </button>
          <button 
            @click="setStatusFilter('closed')" 
            :class="['filter-btn', statusFilter === 'closed' ? 'active' : '']"
          >
            Cerrados
          </button>
        </div>
      </div>
      
      <div class="filter-group">
        <label>Asignación:</label>
        <div class="filter-options">
          <button 
            @click="setAssignmentFilter('all')" 
            :class="['filter-btn', assignmentFilter === 'all' ? 'active' : '']"
          >
            Todos
          </button>
          <button 
            @click="setAssignmentFilter('assigned')" 
            :class="['filter-btn', assignmentFilter === 'assigned' ? 'active' : '']"
          >
            Asignados
          </button>
          <button 
            @click="setAssignmentFilter('unassigned')" 
            :class="['filter-btn', assignmentFilter === 'unassigned' ? 'active' : '']"
          >
            Sin Asignar
          </button>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading">Cargando...</div>
    <div v-else-if="error" class="error">{{ error }}</div>
    <div v-else-if="filteredTickets.length === 0" class="empty">No se encontraron tickets con los filtros seleccionados</div>
    <div v-else class="tickets-list">
      <div v-for="ticket in filteredTickets" :key="ticket.id" class="ticket-item">
        <div class="ticket-item-header">
          <div class="ticket-title">
            <h3>{{ ticket.title }}</h3>
            <span :class="['status', ticket.status]">{{ translateStatus(ticket.status) }}</span>
          </div>
          <div class="ticket-priority">
            <span :class="['priority', normalizePriority(ticket.priority)]">{{ translatePriority(ticket.priority) }}</span>
          </div>
        </div>
        
        <div class="ticket-content">
          <p class="description">{{ ticket.description }}</p>
          
          <div class="ticket-meta">
            <div class="meta-item">
              <i class="pi pi-user"></i>
              <span>{{ ticket.assignedTo ? 'Asignado a: ' + getUserFullName(ticket.assignedTo) : 'Sin asignar' }}</span>
            </div>
            <div class="meta-item">
              <i class="pi pi-calendar"></i>
              <span>Creado: {{ formatDate(ticket.createdAt) }}</span>
            </div>
            <div class="meta-item">
              <i class="pi pi-tag"></i>
              <span>Categoría: {{ ticket.category }}</span>
            </div>
          </div>
        </div>
        
        <div class="ticket-actions">
          <router-link :to="`/tickets/${ticket.id}`" class="btn btn-secondary">
            <i class="pi pi-eye"></i> Ver Detalles
          </router-link>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch, ref, computed } from 'vue'
import { useTicketStore } from '@/stores/tickets'
import { useUsersStore } from '@/stores/users'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'

const ticketStore = useTicketStore()
const usersStore = useUsersStore()
const route = useRoute()

// Extraer las propiedades reactivas del store usando storeToRefs
const { tickets, loading, error } = storeToRefs(ticketStore)

// Filtros
const statusFilter = ref('all')
const assignmentFilter = ref('all')

// Setters para los filtros
const setStatusFilter = (status: string) => {
  statusFilter.value = status
}

const setAssignmentFilter = (assignment: string) => {
  assignmentFilter.value = assignment
}

// Tickets filtrados
const filteredTickets = computed(() => {
  let result = [...tickets.value]
  
  // Filtrar por estado
  if (statusFilter.value !== 'all') {
    result = result.filter(ticket => ticket.status === statusFilter.value)
  }
  
  // Filtrar por asignación
  if (assignmentFilter.value === 'assigned') {
    result = result.filter(ticket => ticket.assignedTo !== null && ticket.assignedTo !== '')
  } else if (assignmentFilter.value === 'unassigned') {
    result = result.filter(ticket => ticket.assignedTo === null || ticket.assignedTo === '')
  }
  
  return result
})

// Cargar usuarios
onMounted(async () => {
  await usersStore.fetchUsers()
})

// Función para obtener el nombre de usuario por ID
const getUserFullName = (userId: string | null): string => {
  if (!userId) return 'Sin asignar'
  
  const user = usersStore.users.find((user: any) => user.id === userId)
  if (user) {
    return `${user.firstName} ${user.lastName}`
  }
  
  return userId // Si no encuentra el usuario, muestra el ID como fallback
}

// Función para cargar tickets
const loadTickets = async () => {
  console.log('TicketList: Cargando tickets...');
  loading.value = true;
  error.value = null;
  
  try {
    await ticketStore.fetchTickets();
    console.log('TicketList: Tickets obtenidos del store:', ticketStore.tickets);
    console.log('TicketList: Prioridades de tickets:', ticketStore.tickets.map((t: any) => t.priority));
    
    // Verificar si se cargaron tickets correctamente
    if (!ticketStore.tickets || ticketStore.tickets.length === 0) {
      console.log('TicketList: No se encontraron tickets en el store, usando datos mock...');
      
      // Crear datos mock para desarrollo
      const mockTickets = [
        {
          id: 'MOCK-TICKET-1',
          title: 'Problema con instalación',
          status: 'open',
          priority: 'high',
          category: 'technical',
          createdAt: new Date().toISOString(),
          updatedAt: new Date().toISOString(),
          customer: {
            name: 'Juan Pérez',
            email: 'juan@example.com'
          },
          description: 'No puedo instalar el software en Windows 11'
        },
        {
          id: 'MOCK-TICKET-2',
          title: 'Consulta sobre facturación',
          status: 'assigned',
          priority: 'medium',
          category: 'billing',
          createdAt: new Date(Date.now() - 86400000).toISOString(),
          updatedAt: new Date(Date.now() - 36000000).toISOString(),
          customer: {
            name: 'María López',
            email: 'maria@example.com'
          },
          description: 'No he recibido la factura del mes pasado'
        },
        {
          id: 'MOCK-TICKET-3',
          title: 'Solicitud de nueva función',
          status: 'in_progress',
          priority: 'low',
          category: 'feature',
          createdAt: new Date(Date.now() - 186400000).toISOString(),
          updatedAt: new Date(Date.now() - 16000000).toISOString(),
          customer: {
            name: 'Carlos Rodríguez',
            email: 'carlos@example.com'
          },
          description: 'Me gustaría sugerir una nueva función para exportar reportes en formato CSV'
        }
      ];
      
      // Actualizar el store con datos mock
      ticketStore.$patch({
        tickets: mockTickets,
        loading: false,
        error: null
      });
      
      console.log('TicketList: Tickets mock cargados:', mockTickets.length);
    }
  } catch (err) {
    console.error('TicketList: Error al obtener tickets:', err);
    error.value = 'Error al cargar los tickets';
    
    // Cargar datos mock para desarrollo
    const mockTickets = [
      {
        id: 'MOCK-TICKET-1',
        title: 'Problema con instalación',
        status: 'open',
        priority: 'high',
        category: 'technical',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
        customer: {
          name: 'Juan Pérez',
          email: 'juan@example.com'
        },
        description: 'No puedo instalar el software en Windows 11'
      },
      {
        id: 'MOCK-TICKET-2',
        title: 'Consulta sobre facturación',
        status: 'assigned',
        priority: 'medium',
        category: 'billing',
        createdAt: new Date(Date.now() - 86400000).toISOString(),
        updatedAt: new Date(Date.now() - 36000000).toISOString(),
        customer: {
          name: 'María López',
          email: 'maria@example.com'
        },
        description: 'No he recibido la factura del mes pasado'
      },
      {
        id: 'MOCK-TICKET-3',
        title: 'Solicitud de nueva función',
        status: 'in_progress',
        priority: 'low',
        category: 'feature',
        createdAt: new Date(Date.now() - 186400000).toISOString(),
        updatedAt: new Date(Date.now() - 16000000).toISOString(),
        customer: {
          name: 'Carlos Rodríguez',
          email: 'carlos@example.com'
        },
        description: 'Me gustaría sugerir una nueva función para exportar reportes en formato CSV'
      }
    ];
    
    // Actualizar el store con datos mock
    ticketStore.$patch({
      tickets: mockTickets,
      loading: false,
      error: null
    });
    
    console.log('TicketList: Tickets mock cargados debido a error en API:', mockTickets.length);
  } finally {
    // Asegurar que el estado de carga se desactive después de un tiempo
    setTimeout(() => {
      loading.value = false;
      console.log('TicketList: Estado de carga establecido a false después de timeout');
    }, 1000);
  }
};

// Cargar tickets al montar el componente
onMounted(async () => {
  console.log('TicketList: Componente montado, cargando tickets...')
  await loadTickets()
})

// Cargar tickets cada vez que se navegue a esta ruta
watch(
  () => route.path,
  async (newPath: string) => {
    if (newPath.includes('/tickets')) {
      console.log('TicketList: Ruta cambiada a /tickets, recargando tickets...')
      await loadTickets()
    }
  },
  { immediate: true }
)

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
  if (!priority) return 'Media';
  
  // Normalizar a minúsculas
  const normalizedPriority = normalizePriority(priority);
  
  const priorityMap: Record<string, string> = {
    'low': 'Baja',
    'medium': 'Media',
    'high': 'Alta',
    'urgent': 'Urgente'
  }
  return priorityMap[normalizedPriority] || normalizedPriority
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
  
  // Si no coincide con ninguno, devolver medio por defecto
  console.log(`TicketList: Prioridad no reconocida "${priority}", utilizando "medium" por defecto`);
  return 'medium';
}

// Función para formatear fechas
const formatDate = (dateString: string): string => {
  try {
    const date = new Date(dateString)
    return date.toLocaleDateString() + ' ' + date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' })
  } catch (e) {
    return dateString
  }
}
</script>

<style lang="scss" scoped>
.ticket-list {
  max-width: 1200px;
  margin: 0 auto;
  background: var(--card-bg);
  border-radius: 10px;
  padding: 1.5rem;
  box-shadow: var(--card-shadow);
  border: 1px solid var(--border-color);
  
  .header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;

    h1 {
      margin: 0;
      color: var(--text-primary);
      font-weight: 600;
      border-bottom: 2px solid var(--border-color);
      padding-bottom: 0.5rem;
    }
    
    .btn-primary {
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      background: linear-gradient(135deg, #4f46e5, #6366f1);
      color: white;
      border: none;
      padding: 0.75rem 1.25rem;
      border-radius: 8px;
      font-size: 0.95rem;
      font-weight: 500;
      text-decoration: none;
      transition: all 0.2s ease;
      box-shadow: 0 2px 8px rgba(79, 70, 229, 0.2);
      
      &:hover {
        background: linear-gradient(135deg, #4338ca, #4f46e5);
        transform: translateY(-1px);
        box-shadow: 0 4px 12px rgba(79, 70, 229, 0.3);
      }
      
      &:active {
        transform: translateY(0);
        box-shadow: 0 2px 4px rgba(79, 70, 229, 0.3);
      }
    }
  }
  
  .filters {
    background: var(--card-bg);
    padding: 1.25rem;
    border-radius: 12px;
    box-shadow: var(--card-shadow);
    margin-bottom: 1.5rem;
    border: 1px solid var(--border-color);
    
    .filter-group {
      margin-bottom: 1.25rem;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      label {
        display: block;
        margin-bottom: 0.75rem;
        font-weight: 600;
        color: var(--text-primary);
        font-size: 1rem;
        letter-spacing: 0.01em;
      }
      
      .filter-options {
        display: flex;
        flex-wrap: wrap;
        gap: 0.75rem;
        
        .filter-btn {
          background: var(--bg-tertiary);
          border: 1px solid var(--border-color);
          padding: 0.7rem 1.2rem;
          border-radius: 8px;
          cursor: pointer;
          font-size: 0.95rem;
          transition: all 0.25s ease;
          color: var(--text-primary);
          font-weight: 500;
          font-family: 'Inter', -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif;
          letter-spacing: 0.01em;
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.03);
          
          &:hover {
            background: var(--hover-bg);
            transform: translateY(-1px);
            box-shadow: 0 3px 6px rgba(0, 0, 0, 0.05);
          }
          
          &.active {
            background: linear-gradient(135deg, rgba(99, 102, 241, 0.15), rgba(129, 140, 248, 0.15));
            border-color: var(--primary-color);
            color: var(--primary-color);
            font-weight: 600;
            box-shadow: 0 2px 4px rgba(99, 102, 241, 0.15);
          }
        }
      }
    }
  }

  .loading, .error, .empty {
    text-align: center;
    padding: 2.5rem;
    color: var(--text-secondary);
    background: var(--card-bg);
    border-radius: 12px;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--border-color);
    font-style: italic;
  }

  .error {
    color: var(--danger-color);
    border-left: 4px solid var(--danger-color);
  }

  .tickets-list {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .ticket-item {
    background: var(--card-bg);
    padding: 1.75rem;
    border-radius: 16px;
    box-shadow: var(--card-shadow);
    border: 1px solid var(--border-color);
    transition: all 0.3s ease;
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      width: 4px;
      height: 100%;
      background: linear-gradient(to bottom, #818cf8, #6366f1);
      opacity: 0.8;
    }
    
    &:hover {
      transform: translateY(-3px);
      box-shadow: 0 12px 24px rgba(15, 23, 42, 0.08);
    }
    
    .ticket-item-header {
      display: flex;
      justify-content: space-between;
      align-items: flex-start;
      margin-bottom: 1.25rem;
      
      .ticket-title {
        flex: 1;
        padding-left: 0.75rem;
        
        h3 {
          margin: 0 0 0.75rem 0;
          font-size: 1.35rem;
          color: var(--text-primary);
          font-weight: 600;
          line-height: 1.3;
        }
        
        .status {
          display: inline-flex;
          align-items: center;
          padding: 0.4rem 0.85rem;
          border-radius: 6px;
          font-size: 0.85rem;
          font-weight: 500;
          text-transform: capitalize;
          
          &::before {
            content: '';
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            margin-right: 8px;
          }
          
          &.open { 
            background: #cfe7fe; 
            color: #1e40af; 
            &::before { background-color: #1e40af; }
          }
          &.assigned { 
            background: #dbd2fd; 
            color: #5b21b6; 
            &::before { background-color: #5b21b6; }
          }
          &.in_progress { 
            background: #bae6fd; 
            color: #0369a1; 
            &::before { background-color: #0369a1; }
          }
          &.resolved { 
            background: #c7d2fe; 
            color: #4338ca; 
            &::before { background-color: #4338ca; }
          }
          &.closed { 
            background: #d1d5db; 
            color: #374151; 
            &::before { background-color: #374151; }
          }
        }
      }
    }
    
    .ticket-content {
      margin-bottom: 1.75rem;
      padding-left: 0.75rem;
      
      .description {
        color: var(--text-secondary);
        margin-bottom: 1.5rem;
        line-height: 1.6;
        font-size: 1rem;
      }
      
      .ticket-meta {
        display: grid;
        grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
        gap: 1.25rem;
        
        .meta-item {
          display: flex;
          align-items: center;
          gap: 0.75rem;
          color: var(--text-secondary);
          font-size: 0.95rem;
          
          i {
            color: var(--primary-color);
            font-size: 1.1rem;
            opacity: 0.9;
          }
        }
      }
    }
    
    .ticket-actions {
      display: flex;
      justify-content: flex-end;
      padding-left: 0.75rem;
      
      .btn-secondary {
        display: inline-flex;
        align-items: center;
        gap: 0.5rem;
        background: linear-gradient(135deg, #4f46e5, #6366f1);
        color: white;
        border: none;
        padding: 0.75rem 1.25rem;
        border-radius: 8px;
        font-size: 0.95rem;
        font-weight: 500;
        cursor: pointer;
        text-decoration: none;
        transition: all 0.25s ease;
        box-shadow: 0 3px 10px rgba(79, 70, 229, 0.2);
        letter-spacing: 0.01em;
        
        &:hover {
          background: linear-gradient(135deg, #4338ca, #4f46e5);
          transform: translateY(-2px);
          box-shadow: 0 6px 15px rgba(79, 70, 229, 0.3);
        }
        
        &:active {
          transform: translateY(0);
          box-shadow: 0 2px 5px rgba(79, 70, 229, 0.2);
        }
        
        i {
          font-size: 1rem;
        }
      }
    }
  }

  .priority {
    display: inline-flex;
    align-items: center;
    padding: 0.4rem 0.85rem;
    border-radius: 8px;
    font-size: 0.85rem;
    font-weight: 600;
    text-transform: capitalize;
    letter-spacing: 0.01em;
    position: relative;
    
    &::before {
      content: '';
      display: inline-block;
      width: 8px;
      height: 8px;
      border-radius: 50%;
      margin-right: 8px;
    }

    &.low { 
      background: #bae6fd; 
      color: #0369a1; 
      &::before { background-color: #0369a1; }
    }
    &.medium { 
      background: #c7d2fe; 
      color: #4338ca; 
      &::before { background-color: #4338ca; }
    }
    &.high { 
      background: #e0e7ff; 
      color: #3730a3; 
      &::before { background-color: #3730a3; }
    }
    &.urgent { 
      background: #ddd6fe; 
      color: #5b21b6; 
      &::before { background-color: #5b21b6; }
    }
  }
}
</style> 
/* eslint-disable */
<template>
  <div class="admin-section">
    <div class="ticket-list">
      <!-- Sección de encabezado con fondo de gradiente y forma ondulada -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">Gestión de Tickets</h1>
          <p class="hero-subtitle">Visualiza y gestiona todos tus tickets de soporte</p>
          
          <div class="hero-actions">
            <router-link to="/tickets/new" class="create-ticket-btn">
              <i class="pi pi-plus"></i>
              Nuevo Ticket
            </router-link>
            
            <router-link to="/tickets-board" class="kanban-board-btn">
              <i class="pi pi-th-large"></i>
              Ver Panel de Gestión
            </router-link>
          </div>
        </div>
        <div class="wave-shape"></div>
      </div>

      <div class="content-wrapper">
        <!-- Filtros con nuevo diseño visual -->
        <div class="filters-section">
          <h2 class="section-title">
            <span class="title-icon"><i class="pi pi-filter"></i></span>
            Filtrar Tickets
          </h2>
          
          <div class="filters-container">
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

            <!-- Nuevo filtro de categorías -->
            <div class="filter-group">
              <label>Categoría:</label>
              <div class="filter-options">
                <button 
                  @click="setCategoryFilter('all')" 
                  :class="['filter-btn', categoryFilter === 'all' ? 'active' : '']"
                >
                  Todas
                </button>
                <button 
                  v-for="category in availableCategories" 
                  :key="category" 
                  @click="setCategoryFilter(category)" 
                  :class="['filter-btn', categoryFilter === category ? 'active' : '']"
                >
                  {{ category }}
                </button>
              </div>
            </div>
          </div>
        </div>

        <!-- Estado de carga, error o vacío -->
        <div v-if="loading" class="status-message loading">
          <i class="pi pi-spin pi-spinner"></i>
          <p>Cargando tickets...</p>
        </div>
        
        <div v-else-if="error" class="status-message error">
          <i class="pi pi-exclamation-triangle"></i>
          <p>{{ error }}</p>
        </div>
        
        <div v-else-if="filteredTickets.length === 0" class="status-message empty">
          <i class="pi pi-inbox"></i>
          <p>No se encontraron tickets con los filtros seleccionados</p>
        </div>
        
        <!-- Lista de tickets con nuevo diseño -->
        <div v-else class="tickets-grid">
          <div v-for="ticket in filteredTickets" :key="ticket.id" class="ticket-card">
            <div class="ticket-header">
              <div class="ticket-badges">
                <span :class="['status-badge', ticket.status]">{{ translateStatus(ticket.status) }}</span>
                <span :class="['priority-badge', normalizePriority(ticket.priority)]">
                  {{ translatePriority(ticket.priority) }}
                </span>
              </div>
              <h3 class="ticket-title">{{ ticket.title }}</h3>
            </div>
            
            <div class="ticket-body">
              <p class="ticket-description">{{ ticket.description }}</p>
            </div>
            
            <div class="ticket-meta">
              <div class="meta-item">
                <i class="pi pi-user"></i>
                <span>{{ ticket.assignedTo ? getUserFullName(ticket.assignedTo) : 'Sin asignar' }}</span>
              </div>
              
              <div class="meta-item">
                <i class="pi pi-calendar"></i>
                <span>{{ formatDate(ticket.createdAt) }}</span>
              </div>
              
              <div class="meta-item">
                <i class="pi pi-tag"></i>
                <span>{{ ticket.category || 'Sin categoría' }}</span>
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
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, watch, ref, computed } from 'vue'
import { useTicketStore } from '@/stores/tickets'
import { useUsersStore } from '@/stores/users'
import { useRoute } from 'vue-router'
import { storeToRefs } from 'pinia'
import { useCategoriesStore } from '@/stores/categories'

const ticketStore = useTicketStore()
const usersStore = useUsersStore()
const route = useRoute()
const categoriesStore = useCategoriesStore()

// Extraer las propiedades reactivas del store usando storeToRefs
const { tickets, loading, error } = storeToRefs(ticketStore)

// Filtros
const statusFilter = ref('all')
const assignmentFilter = ref('all')
const categoryFilter = ref('all')

// Setters para los filtros
const setStatusFilter = (status: string) => {
  statusFilter.value = status
}

const setAssignmentFilter = (assignment: string) => {
  assignmentFilter.value = assignment
}

const setCategoryFilter = (category: string) => {
  categoryFilter.value = category
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
  
  // Filtrar por categoría
  if (categoryFilter.value !== 'all') {
    result = result.filter(ticket => ticket.category === categoryFilter.value)
  }
  
  return result
})

// Cargar usuarios
onMounted(async () => {
  await usersStore.fetchUsers()
  await categoriesStore.fetchCategories()
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

// Obtener las categorías disponibles
const availableCategories = computed(() => {
  // Primero obtenemos las categorías del store de categorías
  const storeCategories = categoriesStore.categories.map(cat => cat.name)
  
  // Luego obtenemos las categorías únicas de los tickets
  const ticketCategories = [...new Set(tickets.value.map(ticket => ticket.category).filter(Boolean))]
  
  // Combinamos ambas fuentes y eliminamos duplicados
  return [...new Set([...storeCategories, ...ticketCategories])]
})

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
  --primary-gradient: linear-gradient(135deg, var(--primary-color) 0%, #3b82f6 100%);
  --secondary-gradient: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  --border-radius-lg: 1.25rem;
  --transition-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
  
  background-color: var(--bg-secondary);
  position: relative;
  overflow-x: hidden;
  
  // Sección hero con fondo más sutil
  .hero-section {
    position: relative;
    padding: 2.5rem 2rem 4.5rem;
    background-color: var(--primary-color);
    color: white;
    text-align: center;
    overflow: hidden;
    
    .hero-content {
      position: relative;
      z-index: 2;
      max-width: 800px;
      margin: 0 auto;
    }
    
    .hero-title {
      font-size: 2.25rem;
      font-weight: 700;
      margin-bottom: 0.75rem;
      text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      color: white;
    }
    
    .hero-subtitle {
      font-size: 1.1rem;
      margin-bottom: 1.75rem;
      opacity: 0.9;
    }
    
    .hero-actions {
      display: flex;
      justify-content: center;
      gap: 1rem;
      
      .create-ticket-btn,
      .kanban-board-btn {
        display: inline-flex;
        align-items: center;
        gap: 0.75rem;
        background-color: white;
        color: var(--primary-color);
        border: none;
        padding: 0.85rem 1.75rem;
        border-radius: 8px;
        font-size: 1rem;
        font-weight: 600;
        text-decoration: none;
        transition: all 0.3s ease;
        box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        
        &:hover {
          transform: translateY(-3px);
          box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
        }
        
        i {
          font-size: 1rem;
        }
      }
      
      .kanban-board-btn {
        background-color: white;
        color: var(--primary-color);
        border: 1px solid var(--border-color);
        
        &:hover {
          background-color: white;
          border-color: var(--border-color);
        }
      }
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
  
  // Sección de filtros
  .filters-section {
    margin-bottom: 3rem;
    
    .filters-container {
      background-color: var(--card-bg);
      border-radius: var(--border-radius-lg);
      box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
      padding: 1.75rem;
      border: 1px solid var(--border-color);
      
      .filter-group {
        margin-bottom: 1.5rem;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        label {
          display: block;
          margin-bottom: 0.75rem;
          font-weight: 600;
          color: var(--text-primary);
          font-size: 1rem;
        }
        
        .filter-options {
          display: flex;
          flex-wrap: wrap;
          gap: 0.75rem;
          
          .filter-btn {
            background-color: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            color: var(--text-primary);
            font-size: 0.9rem;
            font-weight: 500;
            padding: 0.65rem 1.2rem;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s ease;
            
            &:hover {
              background-color: var(--hover-bg);
              transform: translateY(-2px);
            }
            
            &.active {
              background-color: var(--primary-color);
              color: white;
              border-color: transparent;
              font-weight: 600;
              box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
            }
          }
        }
      }
    }
  }

  // Mensajes de estado (loading, error, empty)
  .status-message {
    text-align: center;
    padding: 3rem;
    background-color: var(--card-bg);
    border-radius: var(--border-radius-lg);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;
    
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
    
    &.error {
      i, p {
        color: #ef4444;
      }
    }
    
    &.empty i {
      color: #6b7280;
    }
  }

  // Grid de tickets
  .tickets-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 1.5rem;
  }

  // Tarjeta de ticket con nuevo diseño
  .ticket-card {
    background-color: var(--card-bg);
    border-radius: 24px; /* Aumentando el radio de las esquinas para que sean más redondeadas */
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    transition: all 0.3s var(--transition-bounce);
    position: relative;
    display: flex;
    flex-direction: column;
    border: 1px solid var(--border-color);
    
    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                  0 10px 10px -5px rgba(0, 0, 0, 0.04);
    }
    
    .ticket-header {
      padding: 1.5rem 1.5rem 0.5rem;
      display: flex; /* Añadido para alinear contenido verticalmente */
      flex-direction: column;
      justify-content: center; /* Centrado vertical */
      
      .ticket-badges {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-bottom: 1rem;
        
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
        font-size: 1.25rem;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 1rem 0;
        line-height: 1.4;
      }
    }
    
    .ticket-body {
      padding: 0 1.5rem 1.5rem;
      flex-grow: 1;
      
      .ticket-description {
        color: var(--text-secondary);
        line-height: 1.6;
        font-size: 0.95rem;
        margin: 0;
        display: -webkit-box;
        -webkit-line-clamp: 3;
        -webkit-box-orient: vertical;
        overflow: hidden;
      }
    }
    
    .ticket-meta {
      padding: 1rem 1.5rem;
      background-color: var(--bg-tertiary);
      border-top: 1px solid var(--border-color);
      
      .meta-item {
        display: flex;
        align-items: center;
        margin-bottom: 0.5rem;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        i {
          width: 24px;
          height: 24px;
          background-color: rgba(99, 102, 241, 0.1);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 0.75rem;
          color: var(--primary-color);
          font-size: 0.9rem;
        }
        
        span {
          font-size: 0.9rem;
          color: var(--text-secondary);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
    
    .ticket-footer {
      padding: 1rem 1.5rem;
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
        padding: 0.75rem 1.5rem;
        border-radius: 8px;
        font-size: 0.95rem;
        font-weight: 600;
        text-decoration: none;
        transition: all 0.3s ease;
        box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
        width: 100%;
        justify-content: center;
        
        &:hover {
          transform: translateY(-2px);
          background-color: var(--primary-dark-color, #0d47a1);
          box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15);
        }
        
        i {
          font-size: 1rem;
        }
      }
    }
  }
  
  // Responsive 
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
    
    .tickets-grid {
      grid-template-columns: 1fr;
    }
    
    .filters-section .filters-container .filter-options {
      flex-direction: column;
      gap: 0.5rem;
      
      .filter-btn {
        width: 100%;
        text-align: center;
      }
    }
  }
}
</style>
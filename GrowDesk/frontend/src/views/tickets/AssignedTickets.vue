<template>
  <div class="assigned-tickets">
    <!-- Hero section similar to dashboard -->
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">Panel de Gestión</h1>
        <p class="hero-subtitle">Gestiona y organiza tus tickets asignados</p>
      </div>
      <div class="wave-shape"></div>
    </div>
    
    <div class="content-wrapper">
      <!-- Estadísticas superiores -->
      <div class="metrics-section">
        <h2 class="section-title">
          <span class="title-icon"><i class="pi pi-chart-bar"></i></span>
          Resumen
        </h2>
        <div class="dashboard-grid">
          <div class="metric-card">
            <div class="card-icon">
              <i class="pi pi-ticket"></i>
            </div>
            <div class="card-content">
              <h3>Total Asignados</h3>
              <p class="number">{{ userTickets.length }}</p>
            </div>
          </div>
          
          <div class="metric-card">
            <div class="card-icon">
              <i class="pi pi-clock"></i>
            </div>
            <div class="card-content">
              <h3>En Progreso</h3>
              <p class="number">{{ openTickets.length }}</p>
            </div>
          </div>
          
          <div class="metric-card">
            <div class="card-icon urgent">
              <i class="pi pi-exclamation-triangle"></i>
            </div>
            <div class="card-content">
              <h3>Urgentes</h3>
              <p class="number">{{ urgentTickets.length }}</p>
            </div>
          </div>
          
          <div class="metric-card">
            <div class="card-icon success">
              <i class="pi pi-check-circle"></i>
            </div>
            <div class="card-content">
              <h3>Completados</h3>
              <p class="number">{{ completedTickets.length }}</p>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Panel de gestión con barra de herramientas -->
      <div class="management-panel">
        <div class="panel-header">
          <h2 class="section-title">
            <span class="title-icon"><i class="pi pi-th-large"></i></span>
            Panel de Tickets
          </h2>
          
          <div class="panel-tools">
            <div class="search-container">
              <i class="pi pi-search"></i>
              <input 
                v-model="searchQuery" 
                placeholder="Buscar tickets..." 
                class="search-input"
                @input="filterTickets" 
              />
            </div>
            
            <div class="filter-options">
              <select v-model="priorityFilter" @change="filterTickets" class="filter-select">
                <option value="all">Todas las prioridades</option>
                <option value="LOW">Baja</option>
                <option value="MEDIUM">Media</option>
                <option value="HIGH">Alta</option>
                <option value="URGENT">Urgente</option>
              </select>
              
              <button @click="showTagManager = true" class="tag-manager-btn">
                <i class="pi pi-tags"></i>
                Etiquetas
              </button>
            </div>
          </div>
        </div>
        
        <div v-if="loading" class="loading-indicator">
          <i class="pi pi-spin pi-spinner"></i> 
          <span>Cargando tu panel de tickets...</span>
        </div>
        
        <div v-else-if="userTickets.length === 0" class="empty-state">
          <i class="pi pi-inbox"></i>
          <h2>No hay tickets asignados</h2>
          <p>Actualmente no tienes tickets asignados a tu usuario.</p>
        </div>
        
        <div v-else class="kanban-board">
          <!-- Columna: Por Hacer -->
          <div class="kanban-column">
            <div class="column-header">
              <h3><i class="pi pi-inbox"></i> Por Hacer</h3>
              <span class="ticket-count">{{ toDoTickets.length }}</span>
            </div>
            <draggable 
              class="column-content"
              v-model="todoList"
              group="tickets"
              item-key="id"
              @change="onMoveTicket($event, 'todo')"
              :animation="200"
              ghost-class="ghost-card"
            >
              <template #item="{ element: ticket }">
                <div 
                  class="ticket-card"
                  :class="[`priority-${ticket.priority.toLowerCase()}`]"
                >
                  <div class="ticket-header">
                    <div class="ticket-id">{{ ticket.id }}</div>
                    <div class="ticket-priority">
                      <span class="priority-badge" :class="ticket.priority.toLowerCase()">
                        {{ translatePriority(ticket.priority) }}
                      </span>
                    </div>
                  </div>
                  
                  <div class="ticket-body">
                    <h4 class="ticket-title">{{ ticket.title }}</h4>
                    
                    <div class="ticket-tags" v-if="getTicketTags(ticket.id).length">
                      <span 
                        v-for="tag in getTicketTags(ticket.id)" 
                        :key="tag.id" 
                        class="ticket-tag"
                        :style="{ backgroundColor: tag.color + '20', color: tag.color }"
                      >
                        {{ tag.name }}
                      </span>
                    </div>
                  </div>
                  
                  <div class="ticket-footer">
                    <div class="ticket-time">
                      <i class="pi pi-clock"></i> {{ formatDate(ticket.updatedAt) }}
                    </div>
                    
                    <div class="ticket-actions">
                      <button @click.stop="startTicket(ticket.id)" class="action-btn start-btn">
                        <i class="pi pi-play"></i>
                      </button>
                      
                      <button @click.stop="openTicketDetail(ticket.id)" class="action-btn view-btn">
                        <i class="pi pi-eye"></i>
                      </button>
                      
                      <button @click.stop="showTagSelector(ticket)" class="action-btn tag-btn">
                        <i class="pi pi-tags"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </template>
              <template #header v-if="toDoTickets.length === 0">
                <div class="empty-column">
                  <i class="pi pi-info-circle"></i>
                  <p>No hay tickets pendientes</p>
                </div>
              </template>
            </draggable>
          </div>
          
          <!-- Columna: En Progreso -->
          <div class="kanban-column">
            <div class="column-header">
              <h3><i class="pi pi-sync"></i> En Progreso</h3>
              <span class="ticket-count">{{ inProgressTickets.length }}</span>
            </div>
            <draggable 
              class="column-content"
              v-model="inProgressList"
              group="tickets"
              item-key="id"
              @change="onMoveTicket($event, 'progress')"
              :animation="200"
              ghost-class="ghost-card"
            >
              <template #item="{ element: ticket }">
                <div 
                  class="ticket-card"
                  :class="[`priority-${ticket.priority.toLowerCase()}`]"
                >
                  <div class="ticket-header">
                    <div class="ticket-id">{{ ticket.id }}</div>
                    <div class="ticket-priority">
                      <span class="priority-badge" :class="ticket.priority.toLowerCase()">
                        {{ translatePriority(ticket.priority) }}
                      </span>
                    </div>
                  </div>
                  
                  <div class="ticket-body">
                    <h4 class="ticket-title">{{ ticket.title }}</h4>
                    
                    <div class="ticket-tags" v-if="getTicketTags(ticket.id).length">
                      <span 
                        v-for="tag in getTicketTags(ticket.id)" 
                        :key="tag.id" 
                        class="ticket-tag"
                        :style="{ backgroundColor: tag.color + '20', color: tag.color }"
                      >
                        {{ tag.name }}
                      </span>
                    </div>
                  </div>
                  
                  <div class="ticket-footer">
                    <div class="ticket-time">
                      <i class="pi pi-clock"></i> {{ formatDate(ticket.updatedAt) }}
                    </div>
                    
                    <div class="ticket-actions">
                      <button @click.stop="resolveTicket(ticket.id)" class="action-btn complete-btn">
                        <i class="pi pi-check"></i>
                      </button>
                      
                      <button @click.stop="openTicketDetail(ticket.id)" class="action-btn view-btn">
                        <i class="pi pi-eye"></i>
                      </button>
                      
                      <button @click.stop="showTagSelector(ticket)" class="action-btn tag-btn">
                        <i class="pi pi-tags"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </template>
              <template #header v-if="inProgressTickets.length === 0">
                <div class="empty-column">
                  <i class="pi pi-info-circle"></i>
                  <p>No hay tickets en progreso</p>
                </div>
              </template>
            </draggable>
          </div>
          
          <!-- Columna: Completados -->
          <div class="kanban-column">
            <div class="column-header">
              <h3><i class="pi pi-check-circle"></i> Completados</h3>
              <span class="ticket-count">{{ completedTickets.length }}</span>
            </div>
            <draggable 
              class="column-content"
              v-model="completedList"
              group="tickets"
              item-key="id"
              @change="onMoveTicket($event, 'completed')"
              :animation="200"
              ghost-class="ghost-card"
            >
              <template #item="{ element: ticket }">
                <div 
                  class="ticket-card completed"
                  :class="[`priority-${ticket.priority.toLowerCase()}`]"
                >
                  <div class="ticket-header">
                    <div class="ticket-id">{{ ticket.id }}</div>
                    <div class="ticket-priority">
                      <span class="priority-badge" :class="ticket.priority.toLowerCase()">
                        {{ translatePriority(ticket.priority) }}
                      </span>
                    </div>
                  </div>
                  
                  <div class="ticket-body">
                    <h4 class="ticket-title">{{ ticket.title }}</h4>
                    
                    <div class="ticket-tags" v-if="getTicketTags(ticket.id).length">
                      <span 
                        v-for="tag in getTicketTags(ticket.id)" 
                        :key="tag.id" 
                        class="ticket-tag"
                        :style="{ backgroundColor: tag.color + '20', color: tag.color }"
                      >
                        {{ tag.name }}
                      </span>
                    </div>
                  </div>
                  
                  <div class="ticket-footer">
                    <div class="ticket-time">
                      <i class="pi pi-clock"></i> {{ formatDate(ticket.updatedAt) }}
                    </div>
                    
                    <div class="ticket-actions">
                      <button @click.stop="reopenTicket(ticket.id)" class="action-btn reopen-btn">
                        <i class="pi pi-replay"></i>
                      </button>
                      
                      <button @click.stop="openTicketDetail(ticket.id)" class="action-btn view-btn">
                        <i class="pi pi-eye"></i>
                      </button>
                    </div>
                  </div>
                </div>
              </template>
              <template #header v-if="completedTickets.length === 0">
                <div class="empty-column">
                  <i class="pi pi-info-circle"></i>
                  <p>No hay tickets completados</p>
                </div>
              </template>
            </draggable>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modal para gestionar etiquetas -->
    <div v-if="showTagManager" class="modal-overlay">
      <div class="modal-container">
        <div class="modal-header">
          <h3>Gestionar Etiquetas</h3>
          <button @click="showTagManager = false" class="close-btn">
            <i class="pi pi-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div class="tag-form">
            <div class="form-group">
              <label>Nombre de etiqueta</label>
              <input 
                type="text" 
                v-model="newTag.name" 
                placeholder="Ej: Bug, Mejora, Urgente..."
                class="form-control"
              />
            </div>
            
            <div class="form-group">
              <label>Color</label>
              <div class="color-selector">
                <div 
                  v-for="color in availableColors" 
                  :key="color"
                  class="color-option"
                  :style="{ backgroundColor: color }"
                  :class="{ selected: newTag.color === color }"
                  @click="newTag.color = color"
                ></div>
              </div>
            </div>
            
            <div class="form-actions">
              <button @click="addTag" class="btn add-tag-btn" :disabled="!newTag.name">
                <i class="pi pi-plus"></i> Añadir Etiqueta
              </button>
            </div>
          </div>
          
          <div class="tag-list">
            <h4>Mis Etiquetas</h4>
            
            <div v-if="tags.length === 0" class="empty-tags">
              No has creado ninguna etiqueta aún
            </div>
            
            <div v-else class="tags-container">
              <div 
                v-for="tag in tags" 
                :key="tag.id" 
                class="tag-item"
              >
                <div 
                  class="tag-color" 
                  :style="{ backgroundColor: tag.color }"
                ></div>
                <div class="tag-name">{{ tag.name }}</div>
                <button @click="deleteTag(tag.id)" class="delete-tag-btn">
                  <i class="pi pi-trash"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modal para asignar etiquetas a un ticket -->
    <div v-if="showTicketTagSelector" class="modal-overlay">
      <div class="modal-container">
        <div class="modal-header">
          <h3>Etiquetas para ticket #{{ selectedTicket?.id }}</h3>
          <button @click="closeTagSelector" class="close-btn">
            <i class="pi pi-times"></i>
          </button>
        </div>
        
        <div class="modal-body">
          <div v-if="tags.length === 0" class="empty-tags">
            No hay etiquetas disponibles. Crea algunas primero.
          </div>
          
          <div v-else class="tags-selector">
            <div 
              v-for="tag in tags" 
              :key="tag.id" 
              class="tag-selector-item"
              :class="{ selected: isTagSelected(tag.id) }"
              @click="toggleTicketTag(tag.id)"
            >
              <div 
                class="tag-color" 
                :style="{ backgroundColor: tag.color }"
              ></div>
              <div class="tag-name">{{ tag.name }}</div>
              <div class="tag-status">
                <i v-if="isTagSelected(tag.id)" class="pi pi-check"></i>
              </div>
            </div>
          </div>
          
          <div class="form-actions">
            <button @click="saveTicketTags" class="btn save-btn">
              <i class="pi pi-check"></i> Guardar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, watch, nextTick } from 'vue';
import { useRouter } from 'vue-router';
import { useTicketStore } from '@/stores/tickets';
import { useAuthStore } from '@/stores/auth';
import { useActivityStore } from '@/stores/activity';
import { VueDraggableNext } from 'vuedraggable';

// Stores y router
const router = useRouter();
const ticketStore = useTicketStore();
const authStore = useAuthStore();
const activityStore = useActivityStore();

// Estado local
const userTickets = ref([]);
const filteredTickets = ref([]);
const loading = ref(true);
const searchQuery = ref('');
const priorityFilter = ref('all');

// Listas para draggable
const todoList = ref([]);
const inProgressList = ref([]);
const completedList = ref([]);

// Estado para gestión de etiquetas
const tags = ref([]); // Lista de etiquetas del usuario
const ticketTags = ref({}); // Mapa de ticketId -> etiquetas asignadas
const showTagManager = ref(false);
const showTicketTagSelector = ref(false);
const selectedTicket = ref(null);
const newTag = ref({ name: '', color: '#4f46e5' });

// Colores disponibles para etiquetas
const availableColors = [
  '#4f46e5', // Indigo
  '#0ea5e9', // Sky blue
  '#06b6d4', // Cyan
  '#059669', // Emerald
  '#16a34a', // Green
  '#65a30d', // Lime
  '#ca8a04', // Yellow
  '#ea580c', // Orange
  '#dc2626', // Red
  '#db2777', // Pink
  '#9333ea', // Purple
  '#475569', // Slate
];

// Computed properties para las columnas del kanban
const toDoTickets = computed(() => {
  const tickets = filteredTickets.value.filter(ticket => 
    ['open', 'assigned'].includes(ticket.status)
  );
  
  // Actualizar la lista para draggable
  todoList.value = [...tickets];
  
  return tickets;
});

const inProgressTickets = computed(() => {
  const tickets = filteredTickets.value.filter(ticket => 
    ticket.status === 'in_progress'
  );
  
  // Actualizar la lista para draggable
  inProgressList.value = [...tickets];
  
  return tickets;
});

const completedTickets = computed(() => {
  const tickets = filteredTickets.value.filter(ticket => 
    ['resolved', 'closed'].includes(ticket.status)
  );
  
  // Actualizar la lista para draggable
  completedList.value = [...tickets];
  
  return tickets;
});

// Computed properties para estadísticas
const openTickets = computed(() => 
  userTickets.value.filter(ticket => 
    ['open', 'assigned', 'in_progress'].includes(ticket.status)
  )
);

const urgentTickets = computed(() => 
  userTickets.value.filter(ticket => 
    ticket.priority === 'URGENT' || ticket.priority === 'HIGH'
  )
);

// Cargar tickets asignados al usuario actual
const loadUserTickets = async () => {
  loading.value = true;
  
  try {
    if (!authStore.user || !authStore.user.id) {
      console.error('No hay usuario autenticado');
      return;
    }
    
    console.log('Cargando tickets para el usuario:', authStore.user.id);
    const tickets = await ticketStore.fetchUserTickets(authStore.user.id);
    userTickets.value = tickets;
    filteredTickets.value = [...tickets];
    
    console.log(`Cargados ${userTickets.value.length} tickets asignados al usuario`);
    
    // Cargar etiquetas guardadas del localStorage
    loadUserTags();
    loadTicketTags();
  } catch (error) {
    console.error('Error al cargar tickets del usuario:', error);
  } finally {
    loading.value = false;
  }
};

// Manejar movimiento de tickets entre columnas
const onMoveTicket = async (event, column) => {
  // Cuando un ticket es movido a una columna, cambiar su estado
  if (event.added) {
    const ticket = event.added.element;
    let newStatus = 'open';
    
    // Determinar el nuevo estado basado en la columna destino
    if (column === 'todo') {
      newStatus = 'assigned';
    } else if (column === 'progress') {
      newStatus = 'in_progress';
    } else if (column === 'completed') {
      newStatus = 'resolved';
    }
    
    // Actualizar el estado del ticket si ha cambiado
    if (ticket.status !== newStatus) {
      try {
        // Guardamos temporalmente el estado anterior para la actividad
        const previousStatus = ticket.status;
        
        // Actualizar en el servidor
        await ticketStore.updateTicketStatus(ticket.id, newStatus);
        
        // Registrar la actividad
        await activityStore.logActivity({
          userId: authStore.user.id,
          type: 'ticket_status_changed',
          targetId: ticket.id,
          description: `Movió el ticket #${ticket.id} a ${translateStatus(newStatus)}`,
          metadata: {
            previousStatus,
            newStatus
          }
        });
        
        // No recargamos inmediatamente para evitar parpadeo en la UI
        // Actualizamos localmente el ticket
        ticket.status = newStatus;
        
        // Recargar tickets después de un breve retraso
        setTimeout(() => {
          loadUserTickets();
        }, 1000);
      } catch (error) {
        console.error('Error al actualizar el estado del ticket:', error);
        // En caso de error, recargar para volver al estado anterior
        loadUserTickets();
      }
    }
  }
};

// Filtrar tickets según criterios
const filterTickets = () => {
  filteredTickets.value = userTickets.value.filter(ticket => {
    // Filtro por búsqueda
    const matchesSearch = 
      searchQuery.value === '' || 
      ticket.title.toLowerCase().includes(searchQuery.value.toLowerCase()) ||
      ticket.id.toLowerCase().includes(searchQuery.value.toLowerCase());
    
    // Filtro por prioridad
    const matchesPriority = 
      priorityFilter.value === 'all' || 
      ticket.priority === priorityFilter.value;
    
    return matchesSearch && matchesPriority;
  });
};

// Actualizar el estado de un ticket a "in_progress"
const startTicket = async (ticketId) => {
  try {
    // Actualizar en UI primero para evitar parpadeo
    const ticket = userTickets.value.find(t => t.id === ticketId);
    if (ticket) {
      const previousStatus = ticket.status;
      ticket.status = 'in_progress';
      
      // Actualizar filteredTickets y listas
      filterTickets();
      
      // Luego actualizar en el servidor
      await ticketStore.updateTicketStatus(ticketId, 'in_progress');
      
      // Registrar la actividad
      await activityStore.logActivity({
        userId: authStore.user.id,
        type: 'ticket_status_changed',
        targetId: ticketId,
        description: `Inició el trabajo en el ticket #${ticketId}`,
        metadata: {
          previousStatus,
          newStatus: 'in_progress'
        }
      });
    }
  } catch (error) {
    console.error('Error al iniciar trabajo en el ticket:', error);
    // Recargar tickets en caso de error
    loadUserTickets();
  }
};

// Actualizar el estado de un ticket a "resolved"
const resolveTicket = async (ticketId) => {
  try {
    await ticketStore.updateTicketStatus(ticketId, 'resolved');
    
    // Registrar la actividad
    await activityStore.logActivity({
      userId: authStore.user.id,
      type: 'ticket_status_changed',
      targetId: ticketId,
      description: `Marcó como resuelto el ticket #${ticketId}`,
      metadata: {
        previousStatus: userTickets.value.find(t => t.id === ticketId)?.status,
        newStatus: 'resolved'
      }
    });
    
    // Recargar tickets
    await loadUserTickets();
  } catch (error) {
    console.error('Error al resolver el ticket:', error);
  }
};

// Reabrir un ticket
const reopenTicket = async (ticketId) => {
  try {
    await ticketStore.updateTicketStatus(ticketId, 'in_progress');
    
    // Registrar la actividad
    await activityStore.logActivity({
      userId: authStore.user.id,
      type: 'ticket_reopened',
      targetId: ticketId,
      description: `Reabrió el ticket #${ticketId}`,
      metadata: {
        previousStatus: userTickets.value.find(t => t.id === ticketId)?.status,
        newStatus: 'in_progress'
      }
    });
    
    // Recargar tickets
    await loadUserTickets();
  } catch (error) {
    console.error('Error al reabrir el ticket:', error);
  }
};

// Abrir detalle de ticket
const openTicketDetail = (ticketId) => {
  router.push(`/tickets/${ticketId}`);
};

// Formatear fecha
const formatDate = (dateString) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return new Intl.DateTimeFormat('es', { 
    dateStyle: 'short',
    timeStyle: 'short'
  }).format(date);
};

// Traducir estado
const translateStatus = (status) => {
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

// Función para generar un ID único
const generateId = () => {
  return 'tag-' + Date.now() + '-' + Math.floor(Math.random() * 1000);
};

// Funciones para gestión de etiquetas
const addTag = () => {
  if (!newTag.value.name.trim()) return;
  
  const tag = {
    id: generateId(),
    name: newTag.value.name,
    color: newTag.value.color
  };
  
  tags.value.push(tag);
  saveUserTags();
  
  // Limpiar formulario
  newTag.value = { name: '', color: '#4f46e5' };
};

const deleteTag = (tagId) => {
  // Eliminar la etiqueta
  tags.value = tags.value.filter(tag => tag.id !== tagId);
  
  // Eliminar la etiqueta de todos los tickets
  for (const ticketId in ticketTags.value) {
    if (ticketTags.value[ticketId].includes(tagId)) {
      ticketTags.value[ticketId] = ticketTags.value[ticketId].filter(id => id !== tagId);
    }
  }
  
  saveUserTags();
  saveTicketTags();
};

// Guardar etiquetas en localStorage
const saveUserTags = () => {
  localStorage.setItem(`${authStore.user.id}_tags`, JSON.stringify(tags.value));
};

// Cargar etiquetas del localStorage
const loadUserTags = () => {
  const savedTags = localStorage.getItem(`${authStore.user.id}_tags`);
  if (savedTags) {
    tags.value = JSON.parse(savedTags);
  }
};

// Guardar asignaciones de etiquetas a tickets
const saveTicketTags = () => {
  localStorage.setItem(`${authStore.user.id}_ticket_tags`, JSON.stringify(ticketTags.value));
};

// Cargar asignaciones de etiquetas a tickets
const loadTicketTags = () => {
  const savedTicketTags = localStorage.getItem(`${authStore.user.id}_ticket_tags`);
  if (savedTicketTags) {
    ticketTags.value = JSON.parse(savedTicketTags);
  }
};

// Mostrar selector de etiquetas para un ticket
const showTagSelector = (ticket) => {
  selectedTicket.value = ticket;
  
  // Inicializar si no existe
  if (!ticketTags.value[ticket.id]) {
    ticketTags.value[ticket.id] = [];
  }
  
  showTicketTagSelector.value = true;
};

// Cerrar selector de etiquetas
const closeTagSelector = () => {
  showTicketTagSelector.value = false;
  selectedTicket.value = null;
};

// Verificar si una etiqueta está seleccionada para el ticket actual
const isTagSelected = (tagId) => {
  if (!selectedTicket.value) return false;
  
  const ticketId = selectedTicket.value.id;
  return ticketTags.value[ticketId] && ticketTags.value[ticketId].includes(tagId);
};

// Alternar selección de etiqueta para un ticket
const toggleTicketTag = (tagId) => {
  if (!selectedTicket.value) return;
  
  const ticketId = selectedTicket.value.id;
  
  if (!ticketTags.value[ticketId]) {
    ticketTags.value[ticketId] = [];
  }
  
  const index = ticketTags.value[ticketId].indexOf(tagId);
  
  if (index === -1) {
    // Añadir la etiqueta
    ticketTags.value[ticketId].push(tagId);
  } else {
    // Quitar la etiqueta
    ticketTags.value[ticketId].splice(index, 1);
  }
};

// Obtener las etiquetas de un ticket específico
const getTicketTags = (ticketId) => {
  if (!ticketTags.value[ticketId]) return [];
  
  return ticketTags.value[ticketId].map(tagId => {
    return tags.value.find(tag => tag.id === tagId);
  }).filter(tag => tag); // Filtrar posibles undefined
};

// Al montar el componente, verificar autenticación y preservar el token
onMounted(async () => {
  // Verificar que el token esté presente antes de cargar tickets
  if (authStore.isAuthenticated) {
    await loadUserTickets();
  } else {
    // Si no está autenticado, intentar reautenticar
    try {
      // Verificar si hay un token en localStorage
      const token = localStorage.getItem('auth_token');
      if (token) {
        await authStore.checkAuthStatus();
        if (authStore.isAuthenticated) {
          await loadUserTickets();
        } else {
          // Si no se pudo reautenticar, redirigir al login
          router.push('/login');
        }
      } else {
        router.push('/login');
      }
    } catch (error) {
      console.error('Error al verificar autenticación:', error);
      router.push('/login');
    }
  }
});

// Observar cambios en los filtros
watch([searchQuery, priorityFilter], () => {
  filterTickets();
});
</script>

<style lang="scss" scoped>
.assigned-tickets {
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
      bottom: 0;
      left: 0;
      width: 100%;
      height: 70px;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' opacity='.25' fill='%23FFFFFF'/%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' opacity='.5' fill='%23FFFFFF'/%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%23FFFFFF'/%3E%3C/svg%3E");
      background-size: cover;
      background-position: center bottom;
    }
  }
  
  .content-wrapper {
    position: relative;
    margin-top: -4rem;
    padding: 0 2rem 2rem;
    max-width: 1400px; 
    margin-left: auto;
    margin-right: auto;
    z-index: 3; /* Añadido para asegurar que esté sobre el fondo */
  }
  
  // Sección de métricas y encabezados
  .metrics-section {
    display: flex;
    flex-direction: column;
    margin-bottom: 2rem;
    
    .section-title {
      display: flex;
      align-items: center;
      margin-bottom: 1.5rem;
      color: var(--text-primary);
      font-size: 1.5rem;
      font-weight: 600;
      
      .title-icon {
        display: flex;
        align-items: center;
        justify-content: center;
        width: 2.5rem;
        height: 2.5rem;
        border-radius: 0.75rem;
        background-color: var(--primary-light);
        color: var(--primary-color);
        margin-right: 0.75rem;
        
        i {
          font-size: 1.25rem;
        }
      }
    }
  }
  
  .dashboard-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: 1.5rem;
    margin-bottom: 1rem;
  }
  
  /* Estilos para las tarjetas métricas */
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
      
      &.success {
        background-color: rgba(16, 185, 129, 0.1);
        
        i {
          color: #10b981;
        }
      }
    }
    
    .card-content {
      flex: 1;
      display: flex;
      flex-direction: column;
      justify-content: center;
      
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
  
  /* Estilo para la tarjeta fantasma durante el arrastre */
  .ghost-card {
    opacity: 0.5;
    background: var(--bg-tertiary) !important;
    border: 2px dashed var(--primary-color) !important;
  }
  
  // Panel de gestión
  .management-panel {
    background-color: var(--card-bg);
    border-radius: var(--border-radius-lg);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    padding: 1.5rem;
    margin-bottom: 2rem;
    
    .panel-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 1.5rem;
      flex-wrap: wrap;
      gap: 1rem;
      
      h2 {
        margin: 0;
      }
      
      .panel-tools {
        display: flex;
        gap: 1rem;
        flex-wrap: wrap;
        
        .search-container {
          position: relative;
          min-width: 240px;
          
          i {
            position: absolute;
            left: 1rem;
            top: 50%;
            transform: translateY(-50%);
            color: var(--text-muted);
          }
          
          .search-input {
            width: 100%;
            padding: 0.75rem 1rem 0.75rem 2.5rem;
            border-radius: 0.5rem;
            border: 1px solid var(--border-color);
            background-color: var(--bg-input);
            color: var(--text-primary);
            
            &:focus {
              outline: none;
              border-color: var(--primary-color);
              box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.2);
            }
          }
        }
        
        .filter-options {
          display: flex;
          gap: 0.75rem;
          
          .filter-select {
            padding: 0.75rem 1rem;
            border-radius: 0.5rem;
            border: 1px solid var(--border-color);
            background-color: var(--bg-input);
            color: var(--text-primary);
            min-width: 180px;
            
            &:focus {
              outline: none;
              border-color: var(--primary-color);
            }
          }
          
          .tag-manager-btn {
            display: flex;
            align-items: center;
            gap: 0.5rem;
            padding: 0.75rem 1rem;
            border-radius: 0.5rem;
            border: none;
            background-color: var(--primary-color);
            color: white;
            font-weight: 600;
            cursor: pointer;
            
            &:hover {
              opacity: 0.9;
            }
          }
        }
      }
    }
  }
  
  // Estados de carga y vacío
  .loading-indicator, .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    padding: 4rem 2rem;
    text-align: center;
    
    i {
      font-size: 3rem;
      margin-bottom: 1.5rem;
      color: var(--text-muted);
    }
    
    h2 {
      margin-top: 0;
      margin-bottom: 0.75rem;
      color: var(--text-primary);
    }
    
    p {
      color: var(--text-secondary);
      max-width: 400px;
      margin: 0 auto;
    }
  }
  
  // Tablero Kanban
  .kanban-board {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
    gap: 1.5rem;
    margin-top: 1.5rem;
    
    .kanban-column {
      background-color: var(--bg-tertiary);
      border-radius: 1rem;
      padding: 1rem;
      
      .column-header {
        display: flex;
        align-items: center;
        justify-content: space-between;
        padding: 0.75rem 0.5rem;
        margin-bottom: 1rem;
        border-bottom: 2px solid var(--border-color);
        
        h3 {
          margin: 0;
          font-size: 1.1rem;
          font-weight: 600;
          color: var(--text-primary);
          display: flex;
          align-items: center;
          gap: 0.5rem;
          
          i {
            color: var(--primary-color);
          }
        }
        
        .ticket-count {
          background-color: var(--bg-secondary);
          color: var(--text-secondary);
          font-size: 0.85rem;
          font-weight: 600;
          padding: 0.25rem 0.75rem;
          border-radius: 1rem;
        }
      }
      
      .column-content {
        min-height: 200px;
        
        .empty-column {
          display: flex;
          flex-direction: column;
          align-items: center;
          padding: 2rem 1rem;
          color: var(--text-secondary);
          text-align: center;
          font-size: 0.9rem;
          
          i {
            font-size: 1.5rem;
            margin-bottom: 0.75rem;
            color: var(--text-muted);
          }
        }
      }
    }
  }
  
  // Tarjetas de ticket
  .ticket-card {
    background-color: var(--card-bg);
    border-radius: 0.75rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    border: 1px solid var(--border-color);
    padding: 1rem;
    margin-bottom: 1rem;
    cursor: pointer;
    transition: transform 0.2s, box-shadow 0.2s;
    position: relative;
    
    &:hover {
      transform: translateY(-4px);
      box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    }
    
    // Barra de prioridad a la izquierda
    &::before {
      content: '';
      position: absolute;
      left: 0;
      top: 0;
      bottom: 0;
      width: 4px;
      border-top-left-radius: 0.75rem;
      border-bottom-left-radius: 0.75rem;
    }
    
    &.priority-low::before {
      background-color: #0ea5e9;
    }
    
    &.priority-medium::before {
      background-color: #4f46e5;
    }
    
    &.priority-high::before {
      background-color: #ea580c;
    }
    
    &.priority-urgent::before {
      background-color: #dc2626;
    }
    
    &.completed {
      opacity: 0.8;
      
      &::after {
        content: '';
        position: absolute;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: rgba(255, 255, 255, 0.1);
        border-radius: 0.75rem;
        pointer-events: none;
      }
    }
    
    .ticket-header {
      display: flex;
      justify-content: space-between;
      margin-bottom: 0.75rem;
      
      .ticket-id {
        font-size: 0.85rem;
        color: var(--text-secondary);
        font-weight: 600;
      }
      
      .ticket-priority {
        .priority-badge {
          padding: 0.2rem 0.5rem;
          border-radius: 1rem;
          font-size: 0.7rem;
          font-weight: 600;
          
          &.low {
            background-color: #dbeafe;
            color: #0ea5e9;
          }
          
          &.medium {
            background-color: #ede9fe;
            color: #4f46e5;
          }
          
          &.high {
            background-color: #fff7ed;
            color: #ea580c;
          }
          
          &.urgent {
            background-color: #fee2e2;
            color: #dc2626;
          }
        }
      }
    }
    
    .ticket-body {
      margin-bottom: 1rem;
      
      .ticket-title {
        margin: 0 0 0.75rem;
        font-size: 1rem;
        font-weight: 600;
        color: var(--text-primary);
        line-height: 1.4;
      }
      
      .ticket-tags {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-top: 0.75rem;
        
        .ticket-tag {
          padding: 0.2rem 0.5rem;
          border-radius: 0.25rem;
          font-size: 0.7rem;
          font-weight: 600;
          white-space: nowrap;
        }
      }
    }
    
    .ticket-footer {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-top: auto;
      
      .ticket-time {
        font-size: 0.8rem;
        color: var(--text-secondary);
        display: flex;
        align-items: center;
        gap: 0.3rem;
      }
      
      .ticket-actions {
        display: flex;
        gap: 0.5rem;
        
        .action-btn {
          width: 28px;
          height: 28px;
          border-radius: 0.375rem;
          border: none;
          background-color: var(--bg-tertiary);
          color: var(--text-secondary);
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          transition: all 0.2s;
          
          &:hover {
            background-color: var(--bg-hover);
          }
          
          &.start-btn {
            background-color: rgba(16, 185, 129, 0.1);
            color: #10b981;
            
            &:hover {
              background-color: rgba(16, 185, 129, 0.2);
            }
          }
          
          &.complete-btn {
            background-color: rgba(79, 70, 229, 0.1);
            color: #4f46e5;
            
            &:hover {
              background-color: rgba(79, 70, 229, 0.2);
            }
          }
          
          &.reopen-btn {
            background-color: rgba(245, 158, 11, 0.1);
            color: #f59e0b;
            
            &:hover {
              background-color: rgba(245, 158, 11, 0.2);
            }
          }
        }
      }
    }
  }
  
  // Modales para gestión de etiquetas
  .modal-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background-color: rgba(0, 0, 0, 0.5);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
    
    .modal-container {
      background-color: var(--card-bg);
      border-radius: 1rem;
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                 0 10px 10px -5px rgba(0, 0, 0, 0.04);
      width: 95%;
      max-width: 500px;
      max-height: 90vh;
      overflow-y: auto;
      
      .modal-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1.25rem 1.5rem;
        border-bottom: 1px solid var(--border-color);
        
        h3 {
          margin: 0;
          font-size: 1.25rem;
          color: var(--text-primary);
        }
        
        .close-btn {
          background: none;
          border: none;
          color: var(--text-secondary);
          width: 32px;
          height: 32px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          
          &:hover {
            background-color: var(--bg-tertiary);
          }
        }
      }
      
      .modal-body {
        padding: 1.5rem;
        
        .tag-form {
          margin-bottom: 2rem;
          
          .form-group {
            margin-bottom: 1.25rem;
            
            label {
              display: block;
              margin-bottom: 0.5rem;
              font-weight: 600;
              color: var(--text-primary);
            }
            
            .form-control {
              width: 100%;
              padding: 0.75rem 1rem;
              border-radius: 0.5rem;
              border: 1px solid var(--border-color);
              background-color: var(--bg-input);
              color: var(--text-primary);
              
              &:focus {
                outline: none;
                border-color: var(--primary-color);
                box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.2);
              }
            }
            
            .color-selector {
              display: flex;
              flex-wrap: wrap;
              gap: 0.75rem;
              
              .color-option {
                width: 32px;
                height: 32px;
                border-radius: 0.375rem;
                cursor: pointer;
                transition: transform 0.2s;
                position: relative;
                
                &:hover {
                  transform: scale(1.1);
                }
                
                &.selected {
                  transform: scale(1.15);
                  box-shadow: 0 0 0 2px white, 0 0 0 4px currentColor;
                  
                  &::after {
                    content: '\2713';
                    position: absolute;
                    top: 50%;
                    left: 50%;
                    transform: translate(-50%, -50%);
                    color: white;
                    font-size: 1rem;
                    font-weight: bold;
                    text-shadow: 0 0 2px rgba(0, 0, 0, 0.5);
                  }
                }
              }
            }
          }
          
          .form-actions {
            display: flex;
            justify-content: flex-end;
            
            .add-tag-btn {
              display: flex;
              align-items: center;
              gap: 0.5rem;
              padding: 0.75rem 1.25rem;
              border-radius: 0.5rem;
              border: none;
              background-color: var(--primary-color);
              color: white;
              font-weight: 600;
              cursor: pointer;
              
              &:disabled {
                opacity: 0.5;
                cursor: not-allowed;
              }
              
              &:not(:disabled):hover {
                opacity: 0.9;
              }
            }
          }
        }
        
        .tag-list {
          h4 {
            margin: 0 0 1rem;
            color: var(--text-primary);
            font-size: 1.1rem;
          }
          
          .empty-tags {
            padding: 2rem;
            text-align: center;
            color: var(--text-secondary);
            background-color: var(--bg-tertiary);
            border-radius: 0.5rem;
          }
          
          .tags-container {
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
            
            .tag-item {
              display: flex;
              align-items: center;
              padding: 0.75rem;
              border-radius: 0.5rem;
              background-color: var(--bg-tertiary);
              border: 1px solid var(--border-color);
              
              .tag-color {
                width: 24px;
                height: 24px;
                border-radius: 4px;
                margin-right: 0.75rem;
              }
              
              .tag-name {
                flex: 1;
                font-weight: 500;
                color: var(--text-primary);
              }
              
              .delete-tag-btn {
                width: 28px;
                height: 28px;
                border-radius: 4px;
                border: none;
                background-color: rgba(220, 38, 38, 0.1);
                color: #dc2626;
                display: flex;
                align-items: center;
                justify-content: center;
                cursor: pointer;
                
                &:hover {
                  background-color: rgba(220, 38, 38, 0.2);
                }
              }
            }
          }
        }
        
        .tags-selector {
          display: flex;
          flex-direction: column;
          gap: 0.75rem;
          margin-bottom: 1.5rem;
          
          .tag-selector-item {
            display: flex;
            align-items: center;
            padding: 0.75rem;
            border-radius: 0.5rem;
            background-color: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            cursor: pointer;
            transition: background-color 0.2s;
            
            &:hover {
              background-color: var(--bg-hover);
            }
            
            &.selected {
              background-color: rgba(var(--primary-rgb), 0.1);
              border-color: var(--primary-color);
            }
            
            .tag-color {
              width: 24px;
              height: 24px;
              border-radius: 4px;
              margin-right: 0.75rem;
            }
            
            .tag-name {
              flex: 1;
              font-weight: 500;
              color: var(--text-primary);
            }
            
            .tag-status {
              width: 28px;
              height: 28px;
              display: flex;
              align-items: center;
              justify-content: center;
              color: var(--primary-color);
            }
          }
        }
        
        .save-btn {
          width: 100%;
          padding: 0.75rem;
          border-radius: 0.5rem;
          border: none;
          background-color: var(--primary-color);
          color: white;
          font-weight: 600;
          cursor: pointer;
          display: flex;
          align-items: center;
          justify-content: center;
          gap: 0.5rem;
          
          &:hover {
            opacity: 0.9;
          }
        }
      }
    }
  }
  
  // Responsive
  @media (max-width: 1200px) {
    .kanban-board {
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    }
  }
  
  @media (max-width: 768px) {
    .hero-section {
      padding: 2rem 1rem 5rem;
      
      .hero-title {
        font-size: 1.75rem;
      }
    }
    
    .content-wrapper {
      padding: 0 1rem 1.5rem;
      margin-top: -3rem; /* Ajustado para mejor visualización en móviles */
    }
    
    .panel-header {
      flex-direction: column;
      align-items: flex-start !important;
      
      .panel-tools {
        width: 100%;
        flex-direction: column;
        
        .search-container {
          width: 100%;
        }
        
        .filter-options {
          width: 100%;
          flex-wrap: wrap;
          
          .filter-select, .tag-manager-btn {
            flex: 1;
          }
        }
      }
    }
    
    .ticket-card {
      .ticket-footer {
        flex-direction: column;
        align-items: flex-start;
        gap: 0.75rem;
        
        .ticket-actions {
          width: 100%;
          justify-content: flex-end;
        }
      }
    }
  }
}
</style> 
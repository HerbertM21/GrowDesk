<template>
  <div class="admin-section">
    <div class="kanban-board">
      <!-- Sección de encabezado con fondo de gradiente y forma ondulada -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">Mi Panel de Gestión</h1>
          <p class="hero-subtitle">Organiza y gestiona tus tickets asignados en un tablero Kanban.</p>
        </div>
        <div class="wave-shape"></div>
      </div>

      <div class="content-wrapper">
        <!-- Título de sección -->
        <div class="section-header">
          <h2 class="section-title">
            <span class="title-icon"><i class="pi pi-th-large"></i></span>
            Herramienta de Gestión
          </h2>
          <p class="section-description">Arrastra los tickets entre columnas para cambiar su estado</p>
        </div>
        
        <!-- Gestión de etiquetas -->
        <div class="tag-management">
          <button class="tag-manager-toggle" @click="showTagManager = !showTagManager">
            <i class="pi pi-minus" :class="showTagManager ? 'pi-minus' : 'pi-tag'"></i>
            {{ showTagManager ? 'Ocultar etiquetas' : 'Gestionar etiquetas' }}
          </button>
          
          <div v-if="showTagManager" class="tag-manager-panel">
            <div class="tag-list">
              <div v-for="tag in ticketStore.tags" :key="tag.id" class="tag-item">
                <div class="tag-color" :style="{ backgroundColor: tag.color }"></div>
                <div class="tag-name">{{ tag.name }}</div>
                <div class="tag-actions">
                  <button class="tag-action-btn edit" @click="editTag(tag)" title="Editar etiqueta">
                    <i class="pi pi-pencil"></i>
                  </button>
                  <button class="tag-action-btn delete" @click="deleteTag(tag.id)" title="Eliminar etiqueta">
                    <i class="pi pi-trash"></i>
                  </button>
                </div>
              </div>
              
              <div v-if="ticketStore.tags.length === 0" class="empty-tags">
                No hay etiquetas disponibles. Crea una nueva.
              </div>
            </div>
            
            <div class="tag-form">
              <h4 class="form-title">{{ isEditing ? 'Editar etiqueta' : 'Crear nueva etiqueta' }}</h4>
              <div class="form-group">
                <label for="tagName">Nombre:</label>
                <input
                  type="text"
                  id="tagName"
                  v-model="currentTag.name"
                  placeholder="Nombre de la etiqueta"
                  class="form-control"
                />
              </div>
              
              <div class="form-group">
                <label for="colorHex">Color (Hexadecimal):</label>
                <div class="color-input-group">
                  <input
                    type="text"
                    id="colorHex"
                    v-model="currentTag.color"
                    placeholder="#RRGGBB"
                    class="form-control color-input"
                    pattern="^#([A-Fa-f0-9]{6})$"
                  />
                  <input
                    type="color"
                    v-model="currentTag.color"
                    class="color-picker"
                    title="Seleccionar color"
                  />
                </div>
                <div class="color-options">
                  <div
                    v-for="color in colorOptions"
                    :key="color"
                    class="color-option"
                    :style="{ backgroundColor: color }"
                    :class="{ active: currentTag.color === color }"
                    @click="currentTag.color = color"
                    :title="color"
                  ></div>
                </div>
              </div>
              
              <div class="tag-form-actions">
                <button class="btn cancel" @click="resetTagForm">Cancelar</button>
                <button class="btn save" @click="saveTag" :disabled="!currentTag.name">
                  {{ isEditing ? 'Actualizar' : 'Crear' }}
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Status messages -->
        <div v-if="isLoading" class="status-message loading">
          <i class="pi pi-spin pi-spinner"></i>
          <p>Cargando tickets...</p>
        </div>
        
        <div v-else-if="hasError" class="status-message error">
          <i class="pi pi-exclamation-triangle"></i>
          <p>{{ errorMessage }}</p>
        </div>
        
        <!-- Panel Kanban when data is loaded -->
        <div v-else class="kanban-container">
          <div 
            v-for="(column, index) in columnsData" 
            :key="column.id" 
            class="kanban-column"
            :class="column.id"
            draggable="true"
            @dragstart="handleColumnDragStart($event, index)"
            @dragover="handleColumnDragOver($event)"
            @dragleave="handleColumnDragLeave($event)"
            @dragenter.prevent
            @drop="handleColumnDrop($event, index)"
          >
            <div class="column-header">
              <!-- Drag handle for column -->
              <div class="column-drag-handle" title="Arrastrar para reordenar columna">
                <i class="pi pi-grip-lines"></i>
              </div>
              
              <!-- Editable column title -->
              <div class="column-title-container">
                <!-- Display title when not editing -->
                <h3 
                  v-if="!column.isEditing" 
                  class="column-title"
                  @click="editColumnTitle(column.id)"
                  title="Click para editar"
                >{{ column.name }}</h3>
                
                <!-- Input field when editing -->
                <input
                  v-else
                  :id="`column-title-${column.id}`"
                  type="text"
                  v-model="column.name"
                  class="column-title-input"
                  @blur="saveColumnTitle(column.id)"
                  @keydown="handleTitleKeydown($event, column.id)"
                  :placeholder="getDefaultName(column.id)"
                />
              </div>
              <span class="ticket-count">{{ getTicketsForColumn(column.id).length }}</span>
            </div>
            
            <div class="column-content" :id="column.id" @drop="handleDrop($event, column.id)" @dragover="handleDragOver" @dragenter.prevent>
              <div 
                v-for="ticket in getTicketsForColumn(column.id)" 
                :key="ticket.id" 
                class="kanban-card"
                draggable="true"
                @dragstart="handleDragStart($event, ticket)"
                @click="showTicketPreview(ticket)"
              >
                <div class="card-header">
                  <span :class="['priority-badge', normalizePriority(ticket.priority)]">
                    {{ translatePriority(ticket.priority) }}
                  </span>
                  <span class="ticket-id">#{{ ticket.id.split('-')[1] || ticket.id }}</span>
                </div>
                
                <h4 class="card-title">{{ ticket.title }}</h4>
                
                <!-- Etiquetas del ticket -->
                <div v-if="ticketTags(ticket.id).length > 0" class="ticket-tags">
                  <span
                    v-for="tag in ticketTags(ticket.id)"
                    :key="tag.id"
                    class="ticket-tag"
                    :style="{ backgroundColor: tag.color }"
                  >
                    {{ tag.name }}
                  </span>
                </div>
                
                <div class="card-footer">
                  <span class="date-info">{{ formatDate(ticket.updatedAt || ticket.createdAt) }}</span>
                  <div class="assignee" v-if="ticket.assignedTo">
                    <span class="avatar">{{ getUserInitials(ticket.assignedTo) }}</span>
                  </div>
                </div>
              </div>
              
              <div v-if="getTicketsForColumn(column.id).length === 0" class="empty-column">
                <p>No hay tickets en esta columna</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Modal para la vista previa del ticket -->
    <div v-if="previewTicket" class="ticket-preview-modal" @click.self="closePreview">
      <div class="ticket-preview-content">
        <div class="preview-header">
          <div class="preview-title-section">
            <span :class="['priority-badge', normalizePriority(previewTicket.priority)]">
              {{ translatePriority(previewTicket.priority) }}
            </span>
            <h3 class="preview-title">{{ previewTicket.title }}</h3>
            <span class="ticket-id">#{{ previewTicket.id.split('-')[1] || previewTicket.id }}</span>
          </div>
          <button class="close-preview-btn" @click="closePreview">
            <i class="pi pi-times"></i>
          </button>
        </div>
        
        <div class="preview-body">
          <div class="preview-info-group">
            <h4>Detalles</h4>
            <div class="preview-info-item">
              <span class="info-label">Estado:</span>
              <span :class="['status-label', previewTicket.status]">{{ translateStatus(previewTicket.status) }}</span>
            </div>
            <div class="preview-info-item">
              <span class="info-label">Categoría:</span>
              <span>{{ previewTicket.category || 'No especificada' }}</span>
            </div>
            <div class="preview-info-item">
              <span class="info-label">Fecha:</span>
              <span>{{ formatDate(previewTicket.createdAt) }}</span>
            </div>
            <div class="preview-info-item">
              <span class="info-label">Última actualización:</span>
              <span>{{ formatDate(previewTicket.updatedAt) }}</span>
            </div>
          </div>
          
          <div class="preview-description-group">
            <h4>Descripción</h4>
            <p class="ticket-description">{{ previewTicket.description }}</p>
          </div>
          
          <div class="preview-tag-group">
            <h4>Etiquetas</h4>
            
            <div v-if="ticketTags(previewTicket.id).length > 0" class="preview-tags">
              <div 
                v-for="tag in ticketTags(previewTicket.id)" 
                :key="tag.id" 
                class="preview-tag"
                :style="{ backgroundColor: tag.color }"
              >
                {{ tag.name }}
                <span class="tag-remove" @click="removeTagFromTicket(previewTicket.id, tag.id)" title="Eliminar etiqueta">
                  <i class="pi pi-times"></i>
                </span>
              </div>
            </div>
            <div v-else class="no-tags">
              No hay etiquetas asignadas a este ticket
            </div>
            
            <div class="tag-selector">
              <select 
                v-model="selectedTagId" 
                class="tag-select" 
                :disabled="availableTagsForTicket(previewTicket.id).length === 0"
              >
                <option value="">Seleccionar etiqueta</option>
                <option 
                  v-for="tag in availableTagsForTicket(previewTicket.id)" 
                  :key="tag.id" 
                  :value="tag.id"
                >
                  {{ tag.name }}
                </option>
              </select>
              <button 
                class="add-tag-btn" 
                @click="addTagToTicket(previewTicket.id, selectedTagId)"
                :disabled="!selectedTagId"
              >
                <i class="pi pi-plus"></i> Añadir
              </button>
            </div>
          </div>
          
          <div class="preview-actions">
            <router-link :to="`/tickets/${previewTicket.id}`" class="view-details-btn">
              Ver detalles completos
            </router-link>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, computed, onMounted, nextTick } from 'vue';
import { useTicketStore as useTickets } from '@/stores/tickets';
import { useAuthStore } from '@/stores/auth';

// Stores
const ticketStore = useTickets();
const authStore = useAuthStore();

// Estado local
const isLoading = ref(true);
const hasError = ref(false);
const errorMessage = ref('');
const draggedTicket = ref(null);
const previewTicket = ref(null);
const showTagManager = ref(false);
const selectedTagId = ref('');

// Estado para arrastrar columnas
const draggedColumnIndex = ref(null);

// Columns with localStorage persistence
const columnsData = ref([
  { id: 'assigned', name: 'Por Hacer', color: '#3498db', isEditing: false },
  { id: 'in_progress', name: 'En Progreso', color: '#9b59b6', isEditing: false },
  { id: 'completed', name: 'Completado', color: '#2ecc71', isEditing: false }
]);

// Estado para la gestión de etiquetas 
const isEditing = ref(false);
const currentTag = ref({
  id: '',
  name: '',
  color: '#3498db',
  category: ''
});
const colorOptions = ref([
  '#3498db', // blue
  '#9b59b6', // purple
  '#2ecc71', // green
  '#e74c3c', // red
  '#f39c12', // orange
  '#1abc9c', // teal
  '#34495e', // dark blue
  '#7f8c8d', // gray
  '#d35400', // dark orange
  '#c0392b', // dark red
  '#16a085', // dark teal
  '#8e44ad', // dark purple
  '#2c3e50', // navy
  '#f1c40f', // yellow
  '#27ae60'  // dark green
]);

// Load column names from localStorage if available
onMounted(() => {
  const savedColumns = localStorage.getItem('kanbanColumns');
  if (savedColumns) {
    try {
      const parsed = JSON.parse(savedColumns);
      // Merge saved names with default structure
      columnsData.value = columnsData.value.map(col => {
        const savedColumn = parsed.find(c => c.id === col.id);
        return savedColumn ? {
          ...col,
          name: savedColumn.name || col.name,
          isEditing: false // Always start not editing
        } : col;
      });
    } catch (e) {
      console.error('Error loading saved column names:', e);
    }
  }
});

// Cargar datos al montar el componente
onMounted(async () => {
  try {
    isLoading.value = true;
    
    // Cargar tickets del usuario actual
    if (authStore.user?.id) {
      await ticketStore.fetchUserTickets(authStore.user?.id);
    } else {
      console.error("No hay usuario autenticado");
      errorMessage.value = "No se pudo identificar al usuario actual";
      hasError.value = true;
    }
    
    // Cargar etiquetas disponibles
    await ticketStore.fetchTags();
    
    // Cargar el orden de las columnas desde localStorage
    const savedColumnOrder = localStorage.getItem('kanbanColumnOrder');
    if (savedColumnOrder) {
      try {
        const orderMap = JSON.parse(savedColumnOrder);
        // Reorganizar las columnas según el orden guardado
        columnsData.value.sort((a, b) => {
          const indexA = orderMap[a.id] !== undefined ? orderMap[a.id] : 999;
          const indexB = orderMap[b.id] !== undefined ? orderMap[b.id] : 999;
          return indexA - indexB;
        });
      } catch (e) {
        console.error('Error al cargar el orden de las columnas:', e);
      }
    }
    
    isLoading.value = false;
  } catch (error) {
    console.error('Error al cargar los datos:', error);
    hasError.value = true;
    errorMessage.value = 'No se pudieron cargar los tickets. Por favor, intenta de nuevo más tarde.';
    isLoading.value = false;
  }
});

// Mostrar vista previa del ticket
const showTicketPreview = (ticket) => {
  previewTicket.value = ticket;
  selectedTagId.value = ''; // Resetear la selección de etiquetas
};

// Cerrar vista previa
const closePreview = () => {
  previewTicket.value = null;
};

// Formatear fecha corta (para las tarjetas)
const formatDate = (dateString) => {
  if (!dateString) return '';
  
  const date = new Date(dateString);
  const now = new Date();
  const diffMs = now - date;
  const diffDays = Math.floor(diffMs / (1000 * 60 * 60 * 24));
  
  if (diffDays === 0) {
    // Hoy, mostrar hora
    return `Hoy, ${date.getHours().toString().padStart(2, '0')}:${date.getMinutes().toString().padStart(2, '0')}`;
  } else if (diffDays === 1) {
    return 'Ayer';
  } else if (diffDays < 7) {
    const days = ['Dom', 'Lun', 'Mar', 'Mié', 'Jue', 'Vie', 'Sáb'];
    return days[date.getDay()];
  } else {
    return `${date.getDate()}/${date.getMonth() + 1}`;
  }
};

// Formatear fecha larga (para el modal)
const formatDateLong = (dateString) => {
  if (!dateString) return '';
  
  const date = new Date(dateString);
  const options = { 
    year: 'numeric', 
    month: 'long', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  };
  
  return date.toLocaleDateString('es-ES', options);
};

// Obtener tickets filtrados por el usuario actual
const filteredTickets = computed(() => {
  if (!ticketStore.tickets || ticketStore.tickets.length === 0) {
    return [];
  }
  
  // Filtrar por tickets asignados al usuario actual
  return ticketStore.tickets.filter(ticket => 
    ticket.assignedTo === authStore.user?.id
  );
});

// Obtener tickets para una columna específica
const getTicketsForColumn = (columnId) => {
  if (!filteredTickets.value || !Array.isArray(filteredTickets.value)) {
    return [];
  }
  
  if (columnId === 'assigned') {
    return filteredTickets.value.filter(t => 
      t.status === 'open' || t.status === 'assigned'
    );
  } else if (columnId === 'in_progress') {
    return filteredTickets.value.filter(t => 
      t.status === 'in_progress'
    );
  } else if (columnId === 'completed') {
    return filteredTickets.value.filter(t => 
      t.status === 'resolved' || t.status === 'closed'
    );
  }
  
  return [];
};

// Normalizar prioridad para mostrar consistentemente
const normalizePriority = (priority) => {
  return ticketStore.normalizePriority(priority).toLowerCase();
};

// Traducir prioridad para mostrar
const translatePriority = (priority) => {
  const normalizedPriority = normalizePriority(priority);
  
  const priorityMap = {
    'low': 'Baja',
    'medium': 'Media',
    'high': 'Alta',
    'urgent': 'Urgente'
  };
  
  return priorityMap[normalizedPriority] || 'Media';
};

// Traducir estado
const translateStatus = (status) => {
  return ticketStore.translateStatus(status);
};

// Obtener iniciales de un usuario
const getUserInitials = (userId) => {
  // En un escenario real, esto obtendría datos del usuario
  // Simplificado para demostración
  if (!userId) return '';
  
  const userName = getUserName(userId);
  if (userName === 'Sin asignar') return '??';
  
  // Extraer iniciales del nombre (asumiendo formato "Nombre Apellido")
  const parts = userName.split(' ');
  if (parts.length >= 2) {
    return `${parts[0].charAt(0)}${parts[1].charAt(0)}`;
  }
  
  return parts[0].charAt(0);
};

// Obtener nombre de un usuario
const getUserName = (userId) => {
  if (!userId) return 'Sin asignar';
  
  // En un escenario real, esto obtendría datos del usuario desde el store
  // Para demostración, usamos el ID como nombre
  return userId === authStore.user?.id ? authStore.userFullName || 'Usuario actual' : 'Usuario ' + userId;
};

// Manejar inicio de arrastre
const handleDragStart = (event, ticket) => {
  draggedTicket.value = ticket;
  event.dataTransfer.effectAllowed = 'move';
  event.dataTransfer.setData('text/plain', ticket.id);
  
  // Añadir clase visual durante el arrastre
  event.target.classList.add('dragging');
};

// Manejar soltar ticket
const handleDrop = async (event, columnId) => {
  event.preventDefault();
  
  // Eliminar clases visuales
  document.querySelectorAll('.kanban-card').forEach(card => {
    card.classList.remove('dragging');
  });
  
  if (!draggedTicket.value) return;
  
  const ticketId = draggedTicket.value.id;
  const originalStatus = draggedTicket.value.status;
  let newStatus = '';
  
  // Determinar nuevo estado según la columna
  if (columnId === 'assigned') {
    newStatus = 'assigned';
  } else if (columnId === 'in_progress') {
    newStatus = 'in_progress';
  } else if (columnId === 'completed') {
    newStatus = 'resolved';
  }
  
  // Evitar actualizaciones innecesarias
  if (originalStatus === newStatus) return;
  
  try {
    // Find the ticket to ensure we have the current tags before updating
    const ticketIndex = ticketStore.tickets.findIndex(t => t.id === ticketId);
    
    if (ticketIndex !== -1) {
      // Get the current ticket with all its properties including tags
      const currentTicket = ticketStore.tickets[ticketIndex];
      const currentTags = currentTicket.tags ? [...currentTicket.tags] : [];
      
      // Update the ticket in the store
      const updatedTicket = await ticketStore.updateTicketStatus(ticketId, newStatus);
      
      // Force refresh the ticket's tags in the local view to prevent visual disappearance
      if (updatedTicket) {
        // Make sure tags property exists and has the correct tags
        if (!updatedTicket.tags && currentTags.length > 0) {
          updatedTicket.tags = currentTags;
        }
        
        // Update the ticket in the store.tickets array to ensure consistent view
        const updatedIndex = ticketStore.tickets.findIndex(t => t.id === ticketId);
        if (updatedIndex !== -1 && !ticketStore.tickets[updatedIndex].tags) {
          ticketStore.tickets[updatedIndex].tags = currentTags;
        }
      }
    } else {
      console.error('Ticket not found:', ticketId);
      errorMessage.value = 'No se encontró el ticket para actualizar';
      hasError.value = true;
    }
  } catch (error) {
    console.error('Error al actualizar el estado del ticket:', error);
    errorMessage.value = 'No se pudo actualizar el estado del ticket';
    hasError.value = true;
  } finally {
    // Reset the dragged ticket reference
    draggedTicket.value = null;
  }
};

// Manejar el evento dragover para permitir soltar
const handleDragOver = (event) => {
  event.preventDefault();
};

// Funciones para gestionar etiquetas
const ticketTags = (ticketId) => {
  if (!ticketStore.tickets) return [];
  
  const ticket = ticketStore.tickets.find(t => t.id === ticketId);
  if (!ticket || !ticket.tags) return [];
  
  return ticket.tags.map(tag => {
    if (typeof tag === 'string') {
      // Si es un ID, buscar la etiqueta en el store
      const foundTag = ticketStore.tags.find(t => t.id === tag);
      return foundTag || { id: tag, name: 'Etiqueta', color: '#cccccc' };
    }
    return tag;
  });
};

const availableTagsForTicket = (ticketId) => {
  // Obtener todas las etiquetas que no están ya asignadas al ticket
  const currentTagIds = ticketTags(ticketId).map(tag => tag.id);
  return ticketStore.tags.filter(tag => !currentTagIds.includes(tag.id));
};

const addTagToTicket = async (ticketId, tagId) => {
  if (!tagId) return;
  
  try {
    await ticketStore.addTagToTicket(ticketId, tagId);
    selectedTagId.value = ''; // Limpiar selección
  } catch (error) {
    console.error('Error al añadir etiqueta:', error);
  }
};

const removeTagFromTicket = async (ticketId, tagId) => {
  try {
    await ticketStore.removeTagFromTicket(ticketId, tagId);
  } catch (error) {
    console.error('Error al eliminar etiqueta:', error);
  }
};

// Funciones para gestionar el panel de etiquetas
const editTag = (tag) => {
  currentTag.value = { ...tag };
  isEditing.value = true;
};

const saveTag = async () => {
  try {
    if (isEditing.value && currentTag.value.id) {
      // Actualizar etiqueta existente
      await ticketStore.updateTag(currentTag.value.id, {
        name: currentTag.value.name,
        color: currentTag.value.color,
        category: currentTag.value.category
      });
    } else {
      // Crear nueva etiqueta
      await ticketStore.createTag({
        name: currentTag.value.name,
        color: currentTag.value.color,
        category: currentTag.value.category
      });
    }
    resetTagForm();
  } catch (error) {
    console.error('Error al guardar etiqueta:', error);
  }
};

const deleteTag = async (tagId) => {
  if (confirm('¿Estás seguro de que deseas eliminar esta etiqueta?')) {
    try {
      await ticketStore.deleteTag(tagId);
    } catch (error) {
      console.error('Error al eliminar etiqueta:', error);
    }
  }
};

const resetTagForm = () => {
  currentTag.value = {
    id: '',
    name: '',
    color: colorOptions.value[0],
    category: ''
  };
  isEditing.value = false;
};

// Function to enable editing of column title
const editColumnTitle = (columnId) => {
  const column = columnsData.value.find(col => col.id === columnId);
  if (column) {
    column.isEditing = true;
    // Using nextTick to ensure the input is focused after the DOM update
    nextTick(() => {
      const input = document.getElementById(`column-title-${columnId}`);
      if (input) {
        input.focus();
        input.select();
      }
    });
  }
};

// Function to save column title on blur or Enter key
const saveColumnTitle = (columnId) => {
  const column = columnsData.value.find(col => col.id === columnId);
  if (column) {
    // Trim the title and ensure it's not empty
    column.name = column.name.trim() || getDefaultName(columnId);
    column.isEditing = false;
    
    // Save to localStorage
    localStorage.setItem('kanbanColumns', JSON.stringify(columnsData.value));
  }
};

// Function to get default name for a column if title is empty
const getDefaultName = (columnId) => {
  switch (columnId) {
    case 'assigned': return 'Por Hacer';
    case 'in_progress': return 'En Progreso';
    case 'completed': return 'Completado';
    default: return 'Sin título';
  }
};

// Function to handle key press events while editing
const handleTitleKeydown = (event, columnId) => {
  if (event.key === 'Enter') {
    saveColumnTitle(columnId);
  } else if (event.key === 'Escape') {
    // Revert to original title and exit editing mode
    const column = columnsData.value.find(col => col.id === columnId);
    if (column) {
      column.name = getDefaultName(columnId);
      column.isEditing = false;
    }
  }
};

// Funciones para arrastrar y soltar columnas
const handleColumnDragStart = (event, index) => {
  // Guardamos el índice de la columna que se está arrastrando
  draggedColumnIndex.value = index;
  
  // Añadir clase de arrastre para efectos visuales
  event.target.classList.add('column-dragging');
  
  // Establecer datos para la operación de arrastre
  event.dataTransfer.effectAllowed = 'move';
  event.dataTransfer.setData('text/plain', index.toString());
  
  // Añadir un pequeño retraso para que la transición sea visible
  setTimeout(() => {
    // Añadir clase a todas las otras columnas para mostrar que pueden ser destino
    document.querySelectorAll('.kanban-column').forEach((col, i) => {
      if (i !== index) {
        col.classList.add('column-droppable');
      }
    });
  }, 50);
};

const handleColumnDragOver = (event) => {
  // Prevenir comportamiento por defecto para permitir soltar
  event.preventDefault();
  
  // Cambiar el cursor para indicar que se puede soltar
  event.dataTransfer.dropEffect = 'move';
  
  // Añadir clase a la columna sobre la que se arrastra
  const columnElement = event.currentTarget;
  if (columnElement && !columnElement.classList.contains('column-dragging')) {
    // Remover clase de dragover de todas las columnas primero
    document.querySelectorAll('.kanban-column').forEach(col => {
      col.classList.remove('column-dragover');
    });
    // Añadir solo a la columna actual
    columnElement.classList.add('column-dragover');
  }
};

const handleColumnDragLeave = (event) => {
  // Remover clase al salir de la zona de soltar
  event.currentTarget.classList.remove('column-dragover');
};

const handleColumnDrop = (event, dropIndex) => {
  // Prevenir comportamiento por defecto
  event.preventDefault();
  
  // Remover clases de estilo de todas las columnas
  document.querySelectorAll('.kanban-column').forEach(col => {
    col.classList.remove('column-dragging', 'column-dragover', 'column-droppable');
  });
  
  // Si no hay columna arrastrada o se suelta en la misma posición, no hacer nada
  if (draggedColumnIndex.value === null || draggedColumnIndex.value === dropIndex) {
    draggedColumnIndex.value = null;
    return;
  }
  
  // Mover la columna a la nueva posición
  const columnToMove = columnsData.value.splice(draggedColumnIndex.value, 1)[0];
  columnsData.value.splice(dropIndex, 0, columnToMove);
  
  // Guardar el nuevo orden en localStorage
  const orderMap = {};
  columnsData.value.forEach((col, index) => {
    orderMap[col.id] = index;
  });
  localStorage.setItem('kanbanColumnOrder', JSON.stringify(orderMap));
  
  // Resetear el índice de arrastre
  draggedColumnIndex.value = null;
};
</script>

<style scoped>
/* Basic styles */
.kanban-board {
  width: 100%;
  height: 100%;
  display: flex;
  flex-direction: column;
  background-color: var(--bg-secondary);
}

.content-wrapper {
  max-width: 1500px;
  margin: 0 auto;
  padding: 2rem 1.5rem 3rem;
  flex: 1;
}

/* Section header */
.section-header {
  margin-bottom: 2rem;
  text-align: left;
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-bottom: 1rem;
  color: var(--text-primary);
  font-size: 1.5rem;
  font-weight: 600;
  text-align: left;
  background-color: var(--bg-tertiary);
  border-radius: 12px;
  padding: 0.75rem 1.5rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border-left: 4px solid var(--primary-color);
}

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
}

.title-icon i {
  font-size: 1.2rem;
}

.section-description {
  color: var(--text-secondary);
  font-size: 0.95rem;
  margin-left: 0.5rem;
}

/* Kanban columns */
.kanban-container {
  display: flex;
  gap: 1.5rem;
  overflow-x: auto;
  padding-bottom: 1.5rem;
  min-height: 70vh;
}

.kanban-column {
  flex: 1;
  min-width: 350px;
  max-width: 400px;
  background-color: #f7f9fc;
  border-radius: 10px;
  display: flex;
  flex-direction: column;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  transition: all 0.2s ease;
  position: relative;
  
  &::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    transition: all 0.2s ease;
    border-radius: 10px;
  }
  
  &:hover::after {
    box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
  }
}

.column-header {
  padding: 1.25rem 1.25rem 0.75rem;
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.column-header h3 {
  margin: 0;
  font-size: 1.05rem;
  font-weight: 600;
  color: #37474f;
  letter-spacing: 0.02em;
}

.ticket-count {
  background-color: rgba(0, 0, 0, 0.06);
  border-radius: 30px;
  padding: 0.25rem 0.75rem;
  font-size: 0.8rem;
  font-weight: 500;
  color: #546e7a;
}

.column-content {
  flex: 1;
  padding: 0.75rem 1rem;
  overflow-y: auto;
  min-height: 300px;
}

/* Kanban cards */
.kanban-card {
  background-color: white;
  border-radius: 8px;
  margin-bottom: 0.9rem;
  box-shadow: 0 1px 2px rgba(0, 0, 0, 0.07);
  padding: 1rem;
  cursor: pointer;
  transition: transform 0.15s, box-shadow 0.15s;
  border: 1px solid #f0f0f0;
}

.kanban-card:hover {
  transform: translateY(-2px);
  box-shadow: 0 3px 5px rgba(0, 0, 0, 0.1);
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 0.75rem;
}

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

.ticket-id {
  font-size: 0.8rem;
  color: var(--text-secondary);
  font-weight: 500;
}

.card-title {
  margin: 0 0 0.75rem 0;
  font-size: 0.95rem;
  line-height: 1.4;
  display: -webkit-box;
  -webkit-box-orient: vertical;
  -webkit-line-clamp: 2;
  overflow: hidden;
  color: #263238;
}

.card-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 0.6rem;
  border-top: 1px solid #f5f5f5;
}

.date-info {
  color: #78909c;
  font-size: 0.8rem;
}

.assignee {
  display: flex;
  align-items: center;
}

.avatar {
  width: 28px;
  height: 28px;
  border-radius: 50%;
  background-color: #6200ea;
  color: white;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 0.75rem;
  font-weight: 500;
}

.empty-column {
  height: 100px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #999;
  font-size: 0.9rem;
  text-align: center;
  border: 2px dashed #e0e0e0;
  border-radius: 8px;
  margin: 0.5rem 0;
}

/* Status message styles */
.status-message {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem 0;
  text-align: center;
  color: #777;
}

.status-message i {
  font-size: 2.5rem;
  margin-bottom: 1rem;
}

.status-message.loading i {
  color: #6200ea;
}

.status-message.error i {
  color: #e74c3c;
}

.hero-section {
  position: relative;
  padding: 2.5rem 2rem 4.5rem;
  background-color: var(--primary-color);
  color: white;
  text-align: center;
  overflow: hidden;
}

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

/* Ticket Preview Modal */
.ticket-preview-modal {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.4);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(2px);
}

.ticket-preview-content {
  background: white;
  width: 600px;
  max-width: 90%;
  max-height: 90vh;
  border-radius: 12px;
  box-shadow: 0 5px 20px rgba(0, 0, 0, 0.2);
  overflow: hidden;
  display: flex;
  flex-direction: column;
}

.preview-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  padding: 1.5rem;
  border-bottom: 1px solid #f0f0f0;
}

.preview-title-section {
  flex: 1;
}

.preview-title {
  margin: 0.75rem 0 0.5rem;
  font-size: 1.25rem;
  color: #333;
}

.close-preview-btn {
  background: none;
  border: none;
  width: 32px;
  height: 32px;
  border-radius: 50%;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #666;
  cursor: pointer;
  transition: background-color 0.2s;
}

.close-preview-btn:hover {
  background-color: #f2f2f2;
  color: #333;
}

.preview-body {
  padding: 1.5rem;
  overflow-y: auto;
  flex: 1;
}

.preview-info-group {
  margin-bottom: 1.5rem;
}

.preview-info-group h4 {
  margin: 0 0 1rem;
  font-size: 1rem;
  color: #546e7a;
  font-weight: 500;
}

.preview-info-item {
  display: flex;
  margin-bottom: 0.75rem;
  font-size: 0.95rem;
}

.info-label {
  width: 120px;
  min-width: 120px;
  color: #78909c;
}

.status-label {
  padding: 0.2rem 0.6rem;
  border-radius: 3px;
  font-size: 0.85rem;
  font-weight: 500;
}

.status-label.open {
  background-color: #b3e5fc;
  color: #0277bd;
}

.status-label.assigned {
  background-color: #c8e6c9;
  color: #2e7d32;
}

.status-label.in_progress {
  background-color: #e1bee7;
  color: #6a1b9a;
}

.status-label.resolved {
  background-color: #dcedc8;
  color: #33691e;
}

.status-label.closed {
  background-color: #cfd8dc;
  color: #455a64;
}

.preview-description-group {
  margin-top: 1.5rem;
}

.preview-description-group h4 {
  margin: 0 0 1rem;
  font-size: 1rem;
  color: #546e7a;
  font-weight: 500;
}

.ticket-description {
  line-height: 1.6;
  margin: 0;
  color: #455a64;
}

.preview-tag-group {
  margin-top: 1.5rem;
}

.preview-tags {
  margin-bottom: 1rem;
  display: flex;
  flex-wrap: wrap;
  gap: 0.5rem;
}

.preview-tag {
  display: flex;
  align-items: center;
  padding: 0.35rem 0.75rem;
  border-radius: 999px;
  font-size: 0.85rem;
  font-weight: 500;
  text-shadow: 0 1px 1px rgba(0, 0, 0, 0.2);
  color: rgba(255, 255, 255, 0.95);
  border: 1px solid rgba(0, 0, 0, 0.1);
  background-image: linear-gradient(to bottom, rgba(255, 255, 255, 0.1), rgba(0, 0, 0, 0.1));
}

.tag-remove {
  margin-left: 0.5rem;
  width: 18px;
  height: 18px;
  display: inline-flex;
  align-items: center;
  justify-content: center;
  border-radius: 50%;
  background-color: rgba(255, 255, 255, 0.2);
  color: white;
  cursor: pointer;
  font-size: 0.7rem;
  transition: all 0.2s ease;
  
  &:hover {
    background-color: rgba(255, 255, 255, 0.4);
    transform: scale(1.1);
  }
}

.tag-selector {
  margin-top: 1rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  
  select {
    flex-grow: 1;
    padding: 0.65rem;
    border-radius: 8px;
    border: 1px solid var(--border-color);
    background-color: var(--input-bg);
    color: var(--text-primary);
    
    &:focus {
      outline: none;
      border-color: var(--primary-color);
    }
  }
  
  .add-tag-btn {
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 0.65rem 1rem;
    font-weight: 500;
    cursor: pointer;
    transition: transform 0.2s ease, box-shadow 0.2s ease;
    
    &:hover {
      background-color: var(--primary-color) !important;
      transform: translateY(-2px);
      box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
      color: white !important;
    }
    
    &:active, &:focus {
      background-color: var(--primary-color) !important;
      color: white !important;
      outline: none;
    }
    
    &:disabled {
      background-color: var(--disabled-bg);
      color: var(--disabled-color);
      cursor: not-allowed;
    }
  }
}

.preview-actions {
  padding: 1.25rem;
  display: flex;
  justify-content: flex-end;
  border-top: 1px solid #f0f0f0;
}

.view-details-btn {
  background-color: var(--primary-color);
  color: white;
  text-decoration: none;
  display: inline-block;
  padding: 0.75rem 1.5rem;
  border-radius: 8px;
  font-weight: 600;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  transition: transform 0.3s ease, box-shadow 0.3s ease;
  
  &:hover {
    background-color: var(--primary-color) !important;
    transform: translateY(-2px);
    box-shadow: 0 6px 10px rgba(0, 0, 0, 0.15);
    color: white !important;
  }
  
  &:active, &:focus {
    background-color: var(--primary-color) !important;
    color: white !important;
    outline: none;
  }
}

/* Responsive styles */
@media (max-width: 768px) {
  .kanban-container {
    flex-direction: column;
    align-items: stretch;
  }
  
  .kanban-column {
    min-width: auto;
    max-width: none;
    margin-bottom: 1.5rem;
  }
  
  .ticket-preview-content {
    width: 95%;
    max-width: 95%;
  }
}

/* Gestión de etiquetas */
.tag-management {
  margin-bottom: 1.5rem;
  
  .tag-manager-toggle {
    background-color: var(--primary-color);
    color: white;
    border: none;
    border-radius: 8px;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 600;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    cursor: pointer;
    transition: all 0.2s ease;
    box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
    
    &:hover {
      background-color: var(--primary-dark-color);
      transform: translateY(-2px);
      box-shadow: 0 6px 8px rgba(0, 0, 0, 0.15);
    }
    
    i {
      font-size: 1.1rem;
    }
  }
  
  .tag-manager-panel {
    background-color: var(--card-bg);
    border-radius: 16px;
    padding: 1.5rem;
    margin-top: 1rem;
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    
    @media (max-width: 768px) {
      grid-template-columns: 1fr;
    }
    
    .tag-list {
      display: flex;
      flex-direction: column;
      gap: 0.75rem;
      
      .tag-item {
        display: flex;
        align-items: center;
        padding: 0.75rem;
        background-color: var(--bg-tertiary);
        border-radius: 10px;
        transition: all 0.2s ease;
        border: 1px solid var(--border-color);
        
        &:hover {
          background-color: var(--hover-bg);
          transform: translateY(-2px);
          box-shadow: 0 4px 6px rgba(0, 0, 0, 0.05);
        }
        
        .tag-color {
          width: 16px;
          height: 16px;
          border-radius: 4px;
          margin-right: 0.75rem;
          border: 1px solid rgba(0, 0, 0, 0.1);
        }
        
        .tag-name {
          flex-grow: 1;
          font-weight: 500;
          font-size: 0.95rem;
          color: var(--text-primary);
        }
        
        .tag-actions {
          display: flex;
          gap: 0.5rem;
          
          .tag-action-btn {
            display: flex;
            align-items: center;
            justify-content: center;
            width: 32px;
            height: 32px;
            border-radius: 8px;
            border: none;
            background-color: transparent;
            color: var(--text-secondary);
            cursor: pointer;
            transition: all 0.2s ease;
            
            &:hover {
              color: var(--text-primary);
              background-color: var(--bg-hover);
            }
            
            &.edit:hover {
              background-color: rgba(var(--primary-rgb), 0.1);
              color: var(--primary-color);
            }
            
            &.delete:hover {
              background-color: rgba(220, 38, 38, 0.1);
              color: #dc2626;
            }
            
            i {
              font-size: 1rem;
            }
          }
        }
      }
      
      .empty-tags {
        padding: 1.5rem;
        text-align: center;
        color: var(--text-secondary);
        background-color: var(--bg-tertiary);
        border-radius: 10px;
        font-size: 0.95rem;
      }
    }
    
    .tag-form {
      padding: 1.25rem;
      background-color: var(--bg-tertiary);
      border-radius: 12px;
      border: 1px solid var(--border-color);
      
      .form-title {
        margin-top: 0;
        margin-bottom: 1.25rem;
        font-size: 1.1rem;
        color: var(--text-primary);
        font-weight: 600;
        padding-bottom: 0.75rem;
        border-bottom: 1px solid var(--border-color);
      }
      
      .form-group {
        margin-bottom: 1.25rem;
        
        label {
          display: block;
          margin-bottom: 0.5rem;
          font-weight: 500;
          font-size: 0.9rem;
          color: var(--text-secondary);
        }
        
        .form-control {
          width: 100%;
          padding: 0.75rem;
          border-radius: 8px;
          border: 1px solid var(--border-color);
          background-color: var(--input-bg);
          color: var(--text-primary);
          font-size: 0.95rem;
          transition: all 0.2s ease;
          
          &:focus {
            outline: none;
            border-color: var(--primary-color);
            box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
          }
        }
        
        .color-input-group {
          display: flex;
          gap: 0.5rem;
          align-items: center;
          
          .color-input {
            flex-grow: 1;
          }
          
          .color-picker {
            width: 40px;
            height: 40px;
            padding: 0;
            border: 1px solid var(--border-color);
            border-radius: 8px;
            background-color: transparent;
            cursor: pointer;
            overflow: hidden;
            
            &::-webkit-color-swatch-wrapper {
              padding: 0;
            }
            
            &::-webkit-color-swatch {
              border: none;
              border-radius: 6px;
            }
          }
        }
      }
      
      .color-options {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-top: 0.75rem;
        
        .color-option {
          width: 28px;
          height: 28px;
          border-radius: 6px;
          cursor: pointer;
          border: 2px solid transparent;
          transition: all 0.2s ease;
          position: relative;
          
          &:hover {
            transform: scale(1.1);
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            border-color: rgba(255, 255, 255, 0.5);
          }
          
          &.active {
            border-color: var(--text-primary);
            transform: scale(1.15);
            
            &:after {
              content: "✓";
              position: absolute;
              top: 50%;
              left: 50%;
              transform: translate(-50%, -50%);
              color: white;
              font-weight: bold;
              text-shadow: 0 0 2px rgba(0, 0, 0, 0.5);
            }
          }
        }
      }
    }
  }
}

.ticket-tags {
  display: flex;
  flex-wrap: wrap;
  gap: 0.3rem;
  margin: 0.75rem 0;
  
  .ticket-tag {
    font-size: 0.75rem;
    padding: 0.25rem 0.6rem;
    border-radius: 999px;
    font-weight: 500;
    white-space: nowrap;
    text-overflow: ellipsis;
    overflow: hidden;
    max-width: 100%;
    box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
    text-shadow: 0 1px 1px rgba(0, 0, 0, 0.2);
    color: rgba(255, 255, 255, 0.95);
    border: 1px solid rgba(0, 0, 0, 0.1);
    background-image: linear-gradient(to bottom, rgba(255, 255, 255, 0.1), rgba(0, 0, 0, 0.1));
  }
}

.no-tags {
  padding: 0.75rem;
  text-align: center;
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  border-radius: 8px;
  font-style: italic;
  font-size: 0.9rem;
  margin-bottom: 1rem;
}

.column-title-container {
  display: flex;
  flex: 1;
}

.column-title {
  cursor: pointer;
  transition: all 0.2s;
  padding: 5px 8px;
  margin: 0;
  border-radius: 5px;
  font-size: 1.05rem;
  font-weight: 600;
  color: #37474f;
  letter-spacing: 0.02em;
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.06);
  }
  
  &::after {
    content: "✎";
    opacity: 0;
    margin-left: 6px;
    font-size: 0.85rem;
    transition: opacity 0.2s;
    color: #78909c;
  }
  
  &:hover::after {
    opacity: 1;
  }
}

.column-title-input {
  width: 100%;
  font-size: 1.05rem;
  font-weight: 600;
  color: #37474f;
  padding: 5px 8px;
  border: 1px solid var(--primary-color);
  border-radius: 5px;
  background: white;
  box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
  outline: none;
  
  &:focus {
    border-color: var(--primary-color);
  }
}

.column-drag-handle {
  display: flex;
  align-items: center;
  padding: 0 8px;
  margin-right: 8px;
  cursor: grab;
  color: #78909c;
  border-radius: 4px;
  
  &:hover {
    background-color: rgba(0, 0, 0, 0.05);
    color: #546e7a;
  }
  
  i {
    font-size: 1rem;
  }
  
  &:active {
    cursor: grabbing;
  }
}

/* Estilos para la columna que se está arrastrando */
.kanban-column.column-dragging {
  opacity: 0.4;
  transform: rotate(1deg) scale(0.98);
  background: #f7f9fc;
  z-index: 10;
  box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
  
  * {
    pointer-events: none;
  }
  
  &::after {
    border-color: transparent;
    box-shadow: none;
  }
}

/* Columnas que pueden ser destino */
.kanban-column.column-droppable {
  transition: all 0.25s cubic-bezier(0.2, 0.9, 0.5, 1);
}

/* Efecto al pasar sobre una columna destino */
.kanban-column.column-dragover {
  background-color: #f0f7ff;
  transform: translateY(-4px);
  box-shadow: 0 8px 16px rgba(0, 0, 0, 0.1);
  
  /* Indicador visual lateral */
  &::before {
    content: "";
    position: absolute;
    top: 25%;
    left: -8px;
    width: 6px;
    height: 50%;
    background-color: var(--primary-color);
    border-radius: 3px;
    box-shadow: 0 0 8px rgba(var(--primary-rgb), 0.4);
  }
  
  .column-header {
    background-color: rgba(var(--primary-rgb), 0.08);
    transition: background-color 0.2s ease;
  }
}

/* Mejora de transiciones para el arrastre */
.kanban-column {
  transition: all 0.25s cubic-bezier(0.2, 0.9, 0.5, 1);
  position: relative;
  
  &::after {
    content: "";
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    pointer-events: none;
    transition: all 0.2s ease;
    border-radius: 10px;
  }
  
  &:hover::after {
    box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
  }
}
</style> 
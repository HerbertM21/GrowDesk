/* eslint-disable */
<template>
  <div class="ticket-detail">
    <div v-if="isLoading" class="loading">Cargando...</div>
    <div v-else-if="errorMessage" class="error">{{ errorMessage }}</div>
    <div v-else class="ticket-detail-layout">
      <!-- Notificación tipo toast -->
      <transition name="toast-fade">
        <div v-if="notification.show" :class="['notification-toast', notification.type]">
          <div class="notification-message">{{ notification.message }}</div>
          <button @click="notification.show = false" class="close-notification">&times;</button>
        </div>
      </transition>
      
      <!-- Encabezado principal del ticket -->
      <div v-if="currentTicket" class="ticket-header">
        <div class="ticket-header-main">
          <h1>{{ currentTicket.title }}</h1>
          <div class="ticket-quick-meta">
            <div class="ticket-id">
              <strong>ID:</strong> {{ currentTicket.id }}
            </div>
            <div class="ticket-status">
              <span :class="['status-badge', currentTicket.status.toLowerCase()]">{{ translateStatus(currentTicket.status) }}</span>
            </div>
            <div class="ticket-priority">
              <span :class="['priority-badge', normalizePriority(currentTicket.priority)]">{{ translatePriority(currentTicket.priority) }}</span>
            </div>
            <!-- Botones de acción -->
            <div v-if="hasAdminAccess || hasAssistantAccess" class="header-actions">
              <!-- Botón para asignar usuario -->
              <div class="admin-button header-button">
                <button class="filter-btn assign-btn" @click="toggleAssignMenu($event)"
                  :class="{ 'active': showAssignMenu }">
                  <i class="pi pi-user"></i>
                  <span>Asignar</span>
                </button>
                <div class="dropdown-content assign-dropdown" :class="{ 'show': showAssignMenu }">
                  <div v-for="user in supportUsers" :key="user.id" class="dropdown-item" @click.stop="assignToUserQuick(user.id)">
                    <i class="pi pi-user"></i> {{ user.firstName }} {{ user.lastName }}
                  </div>
                  <div class="dropdown-divider"></div>
                  <div class="dropdown-item" @click.stop="removeAssignment(currentTicket.id)">
                    <i class="pi pi-times"></i> Quitar asignación
                  </div>
                </div>
              </div>

              <!-- Botón para cambiar prioridad -->
              <div class="admin-button header-button">
                <button class="filter-btn priority-btn" @click="togglePriorityMenu($event)" 
                  :class="{ 'active': showPriorityMenu }">
                  <i class="pi pi-flag"></i>
                  <span>Prioridad</span>
                </button>
                <div class="dropdown-content priority-dropdown" :class="{ 'show': showPriorityMenu }">
                  <div class="dropdown-item" @click.stop="updatePriorityTo('LOW')">
                    <span class="priority-indicator low"></span> Baja
                  </div>
                  <div class="dropdown-item" @click.stop="updatePriorityTo('MEDIUM')">
                    <span class="priority-indicator medium"></span> Media
                  </div>
                  <div class="dropdown-item" @click.stop="updatePriorityTo('HIGH')">
                    <span class="priority-indicator high"></span> Alta
                  </div>
                  <div class="dropdown-item" @click.stop="updatePriorityTo('URGENT')">
                    <span class="priority-indicator urgent"></span> Urgente
                  </div>
                </div>
              </div>
              
              <!-- Botón para cambiar categoría -->
              <div class="admin-button header-button">
                <button class="filter-btn category-btn" @click="toggleCategoryMenu($event)" 
                  :class="{ 'active': showCategoryMenu }">
                  <i class="pi pi-tag"></i>
                  <span>Categoría</span>
                </button>
                <div class="dropdown-content category-dropdown" :class="{ 'show': showCategoryMenu }">
                  <div v-for="category in availableCategories" :key="category.id" 
                       class="dropdown-item" @click.stop="updateCategoryTo(category.name)">
                    <span class="category-indicator"></span> {{ category.name }}
                  </div>
                </div>
              </div>
              
              <!-- Botón para cerrar ticket -->
              <div class="admin-button header-button" v-if="!isTicketClosed">
                <button class="filter-btn close-btn" @click="showCloseTicketModal = true">
                  <i class="pi pi-check-circle"></i>
                  <span>Cerrar</span>
                </button>
              </div>
            </div>
          </div>
        </div>
        <button @click="toggleSidebar" class="toggle-sidebar-btn">
          <i :class="['pi', showSidebar ? 'pi-chevron-right' : 'pi-chevron-left']"></i>
        </button>
      </div>
      <div v-else class="ticket-header">
        <h1>Ticket #{{ route.params.id }}</h1>
      </div>

      <!-- Contenido principal (estructura de dos columnas) -->
      <div class="ticket-content-wrapper">
        <!-- Columna principal: Chat -->
        <div class="main-column">
          <div class="chat-section">
            <h2>
              Conversación 
              <span v-if="connectionStatus" class="connected-tag">
                <i class="pi pi-check-circle"></i> En tiempo real
              </span>
              <span v-if="isTicketClosed" class="ticket-closed-badge">
                <i class="pi pi-lock"></i> Ticket cerrado
              </span>
            </h2>
            
            <div class="messages" ref="messagesContainer">
              <template v-if="currentMessages && currentMessages.length > 0">
                <div v-for="message in currentMessages" :key="message.id" 
                    :class="['message', message.isClient ? 'client-message' : 'agent-message']">
                  <div class="message-avatar">
                    <div class="avatar-circle" :class="{ 'client-avatar': message.isClient, 'agent-avatar': !message.isClient }">
                      <template v-if="message.isClient">
                        <i class="pi pi-user"></i>
                      </template>
                      <template v-else>
                        <img v-if="getAgentAvatar(message.userId || currentTicket?.assignedTo)" :src="getAgentAvatar(message.userId || currentTicket?.assignedTo)" alt="Foto de perfil" />
                        <i v-else class="pi pi-user-edit"></i>
                      </template>
                    </div>
                  </div>
                  <div class="message-content">
                    <div class="message-header">
                      <span class="sender">{{ message.isClient ? (currentTicket?.customer?.name || 'Cliente') : getMessageSenderName(message) }}</span>
                      <span class="timestamp">{{ formatTimestamp(message.timestamp || message.createdAt) }}</span>
                    </div>
                    <p class="content">{{ message.content }}</p>
                  </div>
                </div>
              </template>
              <div v-else class="no-messages">
                No hay mensajes en este ticket
              </div>
            </div>
            <form @submit.prevent="handleSendMessage" class="message-form" v-if="!isTicketClosed">
              <textarea 
                v-model="newMessage" 
                placeholder="Escribe tu respuesta..." 
                required
                @keydown.enter.prevent="handleEnterKey"
              ></textarea>
              <button type="submit" class="btn btn-primary">
                <i class="pi pi-send"></i> Enviar
              </button>
            </form>
            <div v-else class="chat-closed-message">
              <i class="pi pi-lock"></i>
              <p>Este ticket está cerrado. No se pueden agregar más mensajes.</p>
            </div>
          </div>
        </div>
        
        <!-- Sidebar: Información del ticket -->
        <div class="sidebar" :class="{ 'sidebar-open': showSidebar }">
          <div class="sidebar-header">
            <h3>Detalles del Ticket</h3>
            <button @click="toggleSidebar" class="close-sidebar-btn">
              <i class="pi pi-times"></i>
            </button>
          </div>
          
          <div class="sidebar-content">
            <!-- Sección de datos del ticket -->
            <div class="sidebar-section">
              <h4 class="section-title">
                <i class="pi pi-info-circle"></i>
                Información General
              </h4>
              <div class="ticket-info">
                <div class="info-item">
                  <span class="info-label">Categoría:</span>
                  <span class="info-value">{{ translateCategory(currentTicket.category) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">Creado:</span>
                  <span class="info-value">{{ formatDate(currentTicket.createdAt) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">Última actualización:</span>
                  <span class="info-value">{{ formatDate(currentTicket.updatedAt) }}</span>
                </div>
                <div class="info-item">
                  <span class="info-label">Estado:</span>
                  <span class="info-value">
                    <span :class="['status-badge', currentTicket.status]">{{ translateStatus(currentTicket.status) }}</span>
                  </span>
                </div>
                <div class="info-item">
                  <span class="info-label">Prioridad:</span>
                  <span class="info-value">
                    <span :class="['priority-badge', normalizePriority(currentTicket.priority)]">
                      {{ translatePriority(currentTicket.priority) }}
                    </span>
                  </span>
                </div>
              </div>
            </div>
            
            <!-- Sección de descripción -->
            <div class="sidebar-section">
              <h4 class="section-title">
                <i class="pi pi-align-left"></i>
                Descripción
              </h4>
              <p class="ticket-description">{{ currentTicket.description }}</p>
            </div>
            
            <!-- Sección de asignación (solo información, sin botón) -->
            <div class="sidebar-section">
              <h4 class="section-title">
                <i class="pi pi-user"></i>
                Asignación
              </h4>
              <div class="assignment-info">
                <div class="info-item">
                  <span class="info-label">Asignado a:</span>
                  <span v-if="currentTicket.assignedTo" class="assigned-user">
                    {{ getAssignedUserName(currentTicket.assignedTo) }}
                  </span>
                  <span v-else class="not-assigned">No asignado</span>
                </div>
              </div>
            </div>
            
            <!-- Acciones de administración (en la barra lateral ya no se incluyen botones) -->
            <div class="sidebar-section" v-if="hasAdminAccess || hasAssistantAccess">
              <h4 class="section-title">
                <i class="pi pi-cog"></i>
                Acciones
              </h4>
              <div class="sidebar-info">
                <p class="sidebar-note">Las acciones para este ticket están disponibles en el encabezado principal.</p>
              </div>
            </div>
          </div>
        </div>
      </div>
    </div>
  </div>

  <!-- Modal para cerrar ticket -->
  <div v-if="showCloseTicketModal" class="modal-overlay" @click.self="showCloseTicketModal = false">
    <div class="modal-content">
      <div class="modal-header">
        <h3>Cerrar Ticket</h3>
        <button class="modal-close" @click="showCloseTicketModal = false">&times;</button>
      </div>
      <div class="modal-body">
        <p>Ingresa el motivo para cerrar este ticket:</p>
        <select v-model="closeReason" class="close-reason-select" required>
          <option value="" disabled>Selecciona un motivo</option>
          <option value="solved">Problema resuelto</option>
          <option value="customer_request">Por solicitud del cliente</option>
          <option value="duplicate">Ticket duplicado</option>
          <option value="irrelevant">Ya no es relevante</option>
          <option value="other">Otro motivo</option>
        </select>
        <textarea 
          v-if="closeReason === 'other'" 
          v-model="closeReasonText" 
          placeholder="Especifica el motivo..." 
          class="close-reason-textarea"
          required
        ></textarea>
        <div class="modal-footer">
          <button type="button" class="btn btn-secondary" @click="showCloseTicketModal = false">Cancelar</button>
          <button type="button" class="btn btn-primary" @click="closeTicket" :disabled="!canCloseTicket">Cerrar Ticket</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup>
import { ref, onMounted, onBeforeUnmount, watch, computed } from 'vue';
import { useRoute, useRouter } from 'vue-router';
import { useTicketStore } from '@/stores/tickets';
import { useChatStore } from '@/stores/chat';
import { useAuthStore } from '@/stores/auth';
import { useUsersStore } from '@/stores/users';
import { useCategoriesStore } from '@/stores/categories';

// Obtener stores y route
const route = useRoute();
const router = useRouter();
const ticketStore = useTicketStore();
const chatStore = useChatStore();
const authStore = useAuthStore();
const usersStore = useUsersStore();
const categoriesStore = useCategoriesStore();

// Variables reactivas 
const isLoading = ref(true);
const errorMessage = ref('');
const currentTicket = ref(null);
const currentMessages = ref([]);
const connectionStatus = ref(false);
const newMessage = ref('');
const messagesContainer = ref(null);
let connectionCheckInterval = null;
let monitorInterval = null;

// Variable para controlar la visibilidad de la sidebar
const showSidebar = ref(true);

// Función para alternar la visibilidad de la sidebar
const toggleSidebar = () => {
  showSidebar.value = !showSidebar.value;
};

// Variables para la asignación
const selectedUserId = ref('');
const assignmentRequested = ref(false);
const selectedPriority = ref('');

// Variables para controlar la visibilidad de los menús desplegables
const showPriorityMenu = ref(false);
const showAssignMenu = ref(false);
const showCategoryMenu = ref(false);

// Variables para el cierre de ticket
const showCloseTicketModal = ref(false);
const closeReason = ref('');
const closeReasonText = ref('');

// Sistema de notificación
const notification = ref({
  show: false,
  message: '',
  type: 'success', // 'success', 'error', 'warning'
  timeout: null
});

// Computed properties para roles
const isAdmin = computed(() => {
  const adminStatus = authStore.isAdmin;
  console.log('Estado isAdmin (computed property):', adminStatus);
  return adminStatus;
});

const isAssistant = computed(() => {
  const assistantStatus = authStore.isAssistant;
  console.log('Estado isAssistant (computed property):', assistantStatus);
  return assistantStatus;
});

const isEmployee = computed(() => {
  const employeeStatus = !!authStore.user && authStore.user.role === 'employee';
  console.log('Estado isEmployee (computed property):', employeeStatus);
  return employeeStatus;
});

// Propiedades para verificar acceso
const hasAdminAccess = computed(() => {
  return isAdmin.value;
});

const hasAssistantAccess = computed(() => {
  return isAssistant.value;
});

// Computed property para verificar si se ha cambiado la prioridad
const isPriorityChanged = computed(() => {
  return currentTicket.value && selectedPriority.value !== currentTicket.value.priority;
});

// Computed property para verificar si el ticket está asignado
const isTicketAssigned = computed(() => {
  return !!currentTicket.value && !!currentTicket.value.assignedTo;
});

// Computed property para obtener usuarios de soporte (admin y asistentes)
const supportUsers = computed(() => {
  return usersStore.users.filter(user => 
    user.active && (user.role === 'admin' || user.role === 'assistant')
  );
});

// Computed property para verificar si el ticket está cerrado
const isTicketClosed = computed(() => {
  return currentTicket.value && currentTicket.value.status === 'closed';
});

// Computed property para validar si se puede cerrar el ticket
const canCloseTicket = computed(() => {
  if (closeReason.value === 'other') {
    return closeReasonText.value.trim().length > 0;
  }
  return closeReason.value !== '';
});

// Computed property for categories
const availableCategories = computed(() => {
  return categoriesStore.categories;
});

// Funciones auxiliares
const formatDate = (dateString) => {
  if (!dateString) return '';
  try {
    return new Date(dateString).toLocaleString();
  } catch (e) {
    return dateString;
  }
};

const formatTimestamp = (timestamp) => {
  if (!timestamp) return '';
  try {
    return new Date(timestamp).toLocaleString();
  } catch (e) {
    return timestamp;
  }
};

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

const translateCategory = (category) => {
  const categoryMap = {
    'technical': 'Técnico',
    'billing': 'Facturación',
    'general': 'General',
    'feature': 'Solicitud de Función',
    'soporte': 'Soporte'
  };
  return categoryMap[category] || category;
};

const translateRole = (role) => {
  const roleMap = {
    'admin': 'Administrador',
    'assistant': 'Asistente',
    'employee': 'Empleado'
  };
  return roleMap[role] || role;
};

const getAssignedUserName = (userId) => {
  const user = usersStore.users.find(u => u.id === userId);
  if (user) {
    return `${user.firstName} ${user.lastName}`;
  }
  return userId || 'Desconocido';
};

const scrollToBottom = () => {
  setTimeout(() => {
    if (messagesContainer.value) {
      messagesContainer.value.scrollTop = messagesContainer.value.scrollHeight;
    }
  }, 50);
};

// Función para mostrar notificaciones
const showNotification = (message, type = 'success', duration = 5000) => {
  // Limpiar timeout anterior si existe
  if (notification.value.timeout) {
    clearTimeout(notification.value.timeout);
  }
  
  // Configurar nueva notificación
  notification.value = {
    show: true,
    message,
    type,
    timeout: setTimeout(() => {
      notification.value.show = false;
    }, duration)
  };
};

// Función para normalizar el valor de prioridad
const normalizePriority = (priority) => {
  if (!priority) return 'medium';
  
  // Convertir a minúsculas
  const lowerPriority = priority.toLowerCase();
  
  // Mapear diferentes formatos posibles al formato estándar
  if (lowerPriority === 'baja' || lowerPriority === 'low') return 'low';
  if (lowerPriority === 'media' || lowerPriority === 'medium') return 'medium';
  if (lowerPriority === 'alta' || lowerPriority === 'high') return 'high';
  if (lowerPriority === 'urgente' || lowerPriority === 'urgent') return 'urgent';
  
  // Si no coincide con ninguno, devolver medio por defecto
  return 'medium';
};

// Funciones para manejo de asignaciones
const assignTicketToUser = async () => {
  if (!selectedUserId.value || !currentTicket.value) {
    console.error('No se puede asignar: falta el ticket o el usuario seleccionado');
    showNotification('Falta información para realizar la asignación', 'error');
    return;
  }
  
  try {
    console.log('Intentando asignar ticket:', currentTicket.value.id, 'a usuario:', selectedUserId.value);
    
    // Guardar ID localmente en caso de que currentTicket.value se actualice durante la operación
    const ticketId = currentTicket.value.id;
    const userName = getAssignedUserName(selectedUserId.value);
    
    // Realizar la asignación
    const result = await ticketStore.assignTicket(ticketId, selectedUserId.value);
    
    if (!result) {
      throw new Error('La operación de asignación devolvió un resultado nulo');
    }
    
    console.log('Ticket asignado correctamente:', result);
    
    // Actualizar ticket local - usamos el ID guardado por si acaso
    await ticketStore.fetchTicket(ticketId);
    currentTicket.value = ticketStore.currentTicket;
    
    // Comprobar que el ticket se actualizó correctamente
    if (!currentTicket.value) {
      throw new Error('El ticket no se pudo actualizar después de la asignación');
    }
    
    // Añadir mensaje al chat sobre la asignación
    const systemMessage = `El ticket ha sido asignado a ${userName} por un administrador.`;
    await chatStore.sendMessage(ticketId, systemMessage);
    
    // Mostrar notificación de éxito
    showNotification(`Ticket asignado a ${userName} correctamente`);
    
    // Limpiar selección
    selectedUserId.value = '';
  } catch (error) {
    console.error('Error detallado al asignar ticket:', error);
    showNotification(`Error al asignar el ticket: ${error.message || 'Error desconocido'}`, 'error');
  }
};

const assignToSelf = async () => {
  if (!currentTicket.value || !authStore.user) {
    console.error('No se puede auto-asignar: falta el ticket o la información de usuario');
    showNotification('Falta información para realizar la auto-asignación', 'error');
    return;
  }
  
  try {
    console.log('Auto-asignando ticket:', currentTicket.value.id, 'al usuario actual:', authStore.user.id);
    
    // Guardar ID localmente en caso de que currentTicket.value se actualice durante la operación
    const ticketId = currentTicket.value.id;
    
    // Realizar la asignación
    const result = await ticketStore.assignTicket(ticketId, authStore.user.id);
    
    if (!result) {
      throw new Error('La operación de auto-asignación devolvió un resultado nulo');
    }
    
    console.log('Ticket auto-asignado correctamente:', result);
    
    // Actualizar ticket local - usamos el ID guardado por si acaso
    await ticketStore.fetchTicket(ticketId);
    currentTicket.value = ticketStore.currentTicket;
    
    // Añadir mensaje al chat sobre la auto-asignación
    const systemMessage = `${authStore.user.firstName} ${authStore.user.lastName} se ha asignado a este ticket.`;
    await chatStore.sendMessage(ticketId, systemMessage);
    
    // Mostrar notificación de éxito
    showNotification('Te has asignado a este ticket correctamente');
  } catch (error) {
    console.error('Error detallado al auto-asignar ticket:', error);
    showNotification(`Error al auto-asignar el ticket: ${error.message || 'Error desconocido'}`, 'error');
  }
};

const requestAssignment = async () => {
  if (!currentTicket.value || !authStore.user) {
    console.error('No se puede solicitar asignación: falta el ticket o la información de usuario');
    showNotification('Falta información para solicitar la asignación', 'error');
    return;
  }
  
  try {
    // Guardar ID localmente en caso de que currentTicket.value se actualice durante la operación
    const ticketId = currentTicket.value.id;
    
    // Enviar un mensaje de solicitud
    const message = `[Sistema] ${authStore.user.firstName} ${authStore.user.lastName} ha solicitado que este ticket sea asignado.`;
    await chatStore.sendMessage(ticketId, message);
    
    // Marcar la solicitud como enviada
    assignmentRequested.value = true;
    
    // Mostrar notificación de éxito
    showNotification('Solicitud de asignación enviada correctamente');
  } catch (error) {
    console.error('Error detallado al solicitar asignación:', error);
    showNotification(`Error al solicitar la asignación: ${error.message || 'Error desconocido'}`, 'error');
  }
};

const removeAssignment = async (ticketId) => {
  if (!ticketId) {
    console.error('No se puede eliminar asignación: falta el ID del ticket');
    showNotification('Falta información para eliminar la asignación', 'error');
    hideAssignMenu();
    return;
  }
  
  try {
    isLoading.value = true;
    console.log('Eliminando asignación del ticket:', ticketId);
    
    // Actualizar datos localmente primero
    const originalAssignedTo = currentTicket.value?.assignedTo;
    const originalStatus = currentTicket.value?.status;
    
    if (currentTicket.value) {
      currentTicket.value.assignedTo = null;
      currentTicket.value.status = 'open';
    }
    
    // Actualizar el ticket para quitar la asignación
    try {
      const result = await ticketStore.updateTicket(ticketId, {
        assignedTo: null,
        status: 'open' // Volver a estado abierto
      });
      
      if (!result) {
        throw new Error('La operación de eliminación de asignación devolvió un resultado nulo');
      }
      
      console.log('Asignación eliminada correctamente:', result);
    } catch (updateError) {
      console.error('Error al eliminar la asignación en el backend:', updateError);
      
      // Revertir cambios locales si falla
      if (currentTicket.value) {
        currentTicket.value.assignedTo = originalAssignedTo;
        currentTicket.value.status = originalStatus;
      }
      
      showNotification('Error al eliminar la asignación en el servidor', 'error');
      isLoading.value = false;
      hideAssignMenu();
      return;
    }
    
    // Actualizar ticket local
    try {
      await ticketStore.fetchTicket(ticketId);
      currentTicket.value = ticketStore.currentTicket;
      console.log('Ticket actualizado después de eliminar asignación:', currentTicket.value);
    } catch (fetchError) {
      console.error('Error al recargar el ticket después de eliminar asignación:', fetchError);
      // No mostrar error ya que la operación principal fue exitosa
    }
    
    // Añadir mensaje al chat sobre la eliminación de asignación
    try {
      const systemMessage = `Un administrador ha eliminado la asignación de este ticket.`;
      await chatStore.sendMessage(ticketId, systemMessage);
      console.log('Mensaje de eliminación de asignación enviado al chat');
    } catch (chatError) {
      console.error('Error al enviar mensaje al chat:', chatError);
      // Continuar aunque falle el envío del mensaje
    }
    
    // Mostrar notificación de éxito
    showNotification('Asignación eliminada correctamente');
  } catch (error) {
    console.error('Error general al eliminar asignación:', error);
    showNotification(`Error al eliminar la asignación: ${error.message || 'Error desconocido'}`, 'error');
  } finally {
    isLoading.value = false;
    hideAssignMenu();
  }
};

// Actualizar la prioridad de un ticket directamente
const updatePriorityTo = async (newPriority) => {
  if (!currentTicket.value) {
    console.error('No se puede actualizar la prioridad: ticket actual no está definido');
    hidePriorityMenu();
    return;
  }
  
  try {
    console.log('=== DEPURACIÓN CAMBIO DE PRIORIDAD ===');
    console.log(`Intentando actualizar prioridad de ${currentTicket.value.id} a ${newPriority}`);
    
    // Guardar valores actuales para posible reversión
    const ticketId = currentTicket.value.id;
    const originalPriority = currentTicket.value.priority;
    
    // Actualizar localmente primero para respuesta inmediata
    currentTicket.value.priority = newPriority;
    
    // Mostrar carga
    isLoading.value = true;
    
    // Intentar realizar la actualización en el servidor
    try {
      const response = await ticketStore.updateTicket(ticketId, { priority: newPriority });
      console.log('Respuesta del servidor:', response);
      
      if (!response) {
        throw new Error('La respuesta del servidor está vacía');
      }
      
      // Actualizar vista con datos del servidor
      if (response.priority !== newPriority) {
        console.warn(`La prioridad devuelta por el servidor (${response.priority}) no coincide con la solicitada (${newPriority})`);
      }
      
      // Asegurar que el ticket local tiene los datos correctos
      currentTicket.value.priority = response.priority || newPriority;
      
      // Mensaje de éxito
      showNotification(`La prioridad se ha actualizado a "${translatePriority(newPriority)}" correctamente`);
      
      // Añadir mensaje al chat sobre el cambio de prioridad
      const userName = authStore.user ? `${authStore.user.firstName} ${authStore.user.lastName}` : 'Un administrador';
      const systemMessage = `${userName} ha cambiado la prioridad del ticket a "${translatePriority(newPriority)}".`;
      await chatStore.sendMessage(ticketId, systemMessage);
      
    } catch (error) {
      console.error('Error al actualizar prioridad en el servidor:', error);
      
      // Revertir cambio local si falla
      currentTicket.value.priority = originalPriority;
      
      // Mostrar error
      showNotification(`Error al actualizar prioridad: ${error.message || 'Error de conexión con el servidor'}`, 'error');
    }
  } catch (error) {
    console.error('Error general en actualización de prioridad:', error);
    showNotification('Ocurrió un error al procesar la actualización de prioridad', 'error');
  } finally {
    // Ocultar menú y finalizar carga
    hidePriorityMenu();
    isLoading.value = false;
  }
};

// Asignar rápidamente el ticket a un usuario
const assignToUserQuick = async (userId) => {
  if (!currentTicket.value || !userId) {
    console.error('No se puede asignar: falta el ticket o el ID de usuario');
    showNotification('Falta información para realizar la asignación', 'error');
    hideAssignMenu();
    return;
  }
  
  try {
    console.log('=== DEPURACIÓN ASIGNACIÓN DE TICKET ===');
    console.log(`Intentando asignar ticket ${currentTicket.value.id} a usuario ${userId}`);
    
    // Guardar valores actuales para posible reversión
    const ticketId = currentTicket.value.id;
    const originalAssignedTo = currentTicket.value.assignedTo;
    const originalStatus = currentTicket.value.status;
    
    // Obtener nombre del usuario para notificaciones
    const userName = getAssignedUserName(userId);
    
    // Actualizar localmente primero para respuesta inmediata
    currentTicket.value.assignedTo = userId;
    currentTicket.value.status = 'assigned';
    
    // Mostrar carga
    isLoading.value = true;
    
    // Intentar realizar la actualización en el servidor
    try {
      const response = await ticketStore.updateTicket(ticketId, { 
        assignedTo: userId,
        status: 'assigned'
      });
      
      console.log('Respuesta del servidor:', response);
      
      if (!response) {
        throw new Error('La respuesta del servidor está vacía');
      }
      
      // Asegurar que el ticket local tiene los datos correctos
      currentTicket.value.assignedTo = response.assignedTo || userId;
      currentTicket.value.status = response.status || 'assigned';
      
      // Mensaje de éxito
      showNotification(`El ticket ha sido asignado a ${userName} correctamente`);
      
      // Añadir mensaje al chat sobre la asignación
      const adminName = authStore.user ? `${authStore.user.firstName} ${authStore.user.lastName}` : 'Un administrador';
      const systemMessage = `${adminName} ha asignado este ticket a ${userName}.`;
      await chatStore.sendMessage(ticketId, systemMessage);
      
    } catch (error) {
      console.error('Error al asignar ticket en el servidor:', error);
      
      // Revertir cambio local si falla
      currentTicket.value.assignedTo = originalAssignedTo;
      currentTicket.value.status = originalStatus;
      
      // Mostrar error
      showNotification(`Error al asignar ticket: ${error.message || 'Error de conexión con el servidor'}`, 'error');
    }
  } catch (error) {
    console.error('Error general en asignación de ticket:', error);
    showNotification('Ocurrió un error al procesar la asignación del ticket', 'error');
  } finally {
    // Ocultar menú y finalizar carga
    hideAssignMenu();
    isLoading.value = false;
  }
};

// Lógica principal
const loadTicketData = async () => {
  console.log('Cargando datos del ticket...');
  
  try {
    isLoading.value = true;
    errorMessage.value = '';
    
    // Cargar usuarios inmediatamente
    if (usersStore.users.length === 0) {
      console.log('Inicializando usuarios...');
      usersStore.initMockUsers();
    }

    // Asegurar que las categorías estén cargadas
    console.log('Cargando categorías explícitamente...');
    await categoriesStore.fetchCategories();
    console.log('Categorías disponibles:', categoriesStore.categories);
    
    // Verificar autenticación del usuario
    const userId = localStorage.getItem('userId');
    if (userId && (!authStore.user || !authStore.isAuthenticated)) {
      console.log('Configurando usuario manualmente...');
      const user = usersStore.users.find(u => u.id === userId);
      if (user) {
        // Forzar inicialización manual
        authStore.user = { ...user, role: 'admin' }; // Forzar rol admin
        localStorage.setItem('token', `force-jwt-token-${Date.now()}`);
        authStore.token = localStorage.getItem('token');
        console.log('Usuario configurado manualmente:', user);
      }
    }
    
    // Verificar si ya se tienen los datos de usuarios para asignaciones
    await usersStore.fetchUsers();
    
    // Preparar usuarios para asignación
    const adminUsers = usersStore.users.filter(user => 
      user.active && (user.role === 'admin' || user.role === 'assistant')
    );
    console.log('Usuarios disponibles para asignar:', adminUsers);
    
    // Configurar el ticket actual en el chatStore
    chatStore.setCurrentTicket(route.params.id.toString());
    
    try {
      // Intentar cargar el ticket primero
      await ticketStore.fetchTicket(route.params.id);
      // Usar el currentTicket del store
      currentTicket.value = ticketStore.currentTicket;
      console.log('Ticket obtenido del store:', currentTicket.value);
    } catch (err) {
      console.warn('No se pudo cargar el ticket, pero continuaremos con los mensajes:', err);
    }

    // Cargar mensajes del ticket
    await chatStore.fetchMessages(route.params.id.toString());

    // Conectar WebSocket
    chatStore.connectToWebSocket(route.params.id.toString());

    // Obtener mensajes y estado de conexión
    currentMessages.value = chatStore.getMessagesForTicket(route.params.id.toString());
    connectionStatus.value = chatStore.isConnected();
    
    // Configurar monitoreo de cambios
    setupMonitoring();
    
  } catch (err) {
    console.error('Error al cargar datos del ticket:', err);
    errorMessage.value = 'Error al cargar los datos del ticket';
  } finally {
    isLoading.value = false;
    scrollToBottom();
  }
};

const setupMonitoring = () => {
  // Limpiar intervalos previos
  if (monitorInterval) clearInterval(monitorInterval);
  if (connectionCheckInterval) clearInterval(connectionCheckInterval);
  
  // Monitorear cambios en mensajes
  monitorInterval = setInterval(() => {
    const ticketId = route.params.id;
    const storeMessages = chatStore.messages[ticketId];
    
    if (storeMessages && JSON.stringify(storeMessages) !== JSON.stringify(currentMessages.value)) {
      console.log('Actualizando mensajes desde store:', storeMessages);
      currentMessages.value = [...storeMessages]; // Copia para forzar reactividad
      scrollToBottom();
    }
    
    connectionStatus.value = chatStore.connected;
  }, 1000);
  
  // Verificar conexión WebSocket
  connectionCheckInterval = setInterval(() => {
    if (!chatStore.connected) {
      console.log('Intentando reconectar WebSocket...');
      chatStore.connectToWebSocket(route.params.id);
    }
  }, 5000);
};

const handleSendMessage = async () => {
  if (!newMessage.value.trim()) return;
  
  try {
    const ticketId = route.params.id;
    await chatStore.sendMessage(ticketId, newMessage.value);
    newMessage.value = '';
    
    // Actualizar mensajes locales
    currentMessages.value = chatStore.messages[ticketId] || [];
    scrollToBottom();
  } catch (err) {
    console.error('Error al enviar mensaje:', err);
  }
};

// Añadir la función para manejar el Enter después de las otras funciones
const handleEnterKey = (event) => {
  // Si shift+enter, permitir multilinea
  if (event.shiftKey) {
    return;
  }
  
  // Si no hay texto, no hacer nada
  if (!newMessage.value.trim()) {
    return;
  }
  
  // Enviar el mensaje
  handleSendMessage();
};

// Debug para ver estructura de stores
watch(() => ticketStore.ticket, (newTicket) => {
  console.log('Cambio en ticket:', newTicket);
}, { immediate: true, deep: true });

watch(() => chatStore.messages, (newMessages) => {
  console.log('Cambio en mensajes:', newMessages);
}, { immediate: true, deep: true });

// Cerrar menús al hacer clic fuera de ellos
onMounted(async () => {
  console.log("Componente montado - configurando escuchadores de eventos");
  
  // Forzar inicialización manual
  if (usersStore.users.length === 0) {
    usersStore.initMockUsers();
  }
  
  const currentUserId = localStorage.getItem('userId');
  if (currentUserId) {
    const user = usersStore.users.find(u => u.id === currentUserId);
    if (user) {
      // Forzar asignación manual
      authStore.user = { ...user, role: 'admin' }; // Forzar rol admin
      localStorage.setItem('token', `force-jwt-token-${Date.now()}`);
      authStore.token = localStorage.getItem('token');
    }
  }
  
  // Remover event listeners anteriores para evitar duplicados
  document.removeEventListener('click', handleOutsideClick);
  
  // Configurar event listener para cerrar menús al hacer clic fuera
  document.addEventListener('click', handleOutsideClick);
  
  // Cargar datos inmediatamente
  loadTicketData();
  
  // Cargar categorías
  await categoriesStore.fetchCategories();
});

// Función separada para manejar clics fuera de los menús
const handleOutsideClick = (event) => {
  // Evitamos la ejecución si los menús ya están cerrados
  if (!showPriorityMenu.value && !showAssignMenu.value && !showCategoryMenu.value) {
    return;
  }
  
  // Comprobar si el clic fue en un botón de acción o dentro de un menú desplegable
  const clickedElement = event.target;
  const isClickInsideButton = clickedElement.closest('.filter-btn');
  const isClickInsideDropdown = clickedElement.closest('.dropdown-content');
  
  console.log("Click detectado - verificando si cerrar menús:", {
    isClickInsideButton,
    isClickInsideDropdown,
    currentMenus: {
      priority: showPriorityMenu.value,
      assign: showAssignMenu.value,
      category: showCategoryMenu.value
    }
  });
  
  // Solo cerrar menús si el clic no fue en un botón ni dentro de un menú
  if (!isClickInsideButton && !isClickInsideDropdown) {
    hidePriorityMenu();
    hideAssignMenu();
    hideCategoryMenu();
  }
};

onBeforeUnmount(() => {
  // Limpiar intervalos
  if (connectionCheckInterval) clearInterval(connectionCheckInterval);
  if (monitorInterval) clearInterval(monitorInterval);
  
  // Desconectar WebSocket
  chatStore.disconnectWebSocket();
  
  // Limpiar el event listener cuando se desmonta el componente
  document.removeEventListener('click', handleOutsideClick);
  
  console.log("Componente desmontado - limpiando recursos");
});

// Funciones para mostrar/ocultar menús desplegables
const togglePriorityMenu = (event) => {
  // Evitar propagación para que no cierre inmediatamente
  if (event) {
    event.stopPropagation();
  }
  
  console.log("Toggling menú de prioridad");
  showPriorityMenu.value = !showPriorityMenu.value;
  
  // Si se abre el menú de prioridad, cerrar el menú de asignación
  if (showPriorityMenu.value) {
    showAssignMenu.value = false;
    showCategoryMenu.value = false;
  }
  
  // Actualizar visualización de los menús
  updateDropdownVisibility();
};

const hidePriorityMenu = () => {
  if (!showPriorityMenu.value) return; // Evitar operaciones innecesarias
  
  console.log("Cerrando menú de prioridad");
  showPriorityMenu.value = false;
  updateDropdownVisibility();
};

const toggleAssignMenu = (event) => {
  // Evitar propagación para que no cierre inmediatamente
  if (event) {
    event.stopPropagation();
  }
  
  console.log("Toggling menú de asignación");
  showAssignMenu.value = !showAssignMenu.value;
  
  // Si se abre el menú de asignación, cerrar el menú de prioridad
  if (showAssignMenu.value) {
    showPriorityMenu.value = false;
    showCategoryMenu.value = false;
  }
  
  // Actualizar visualización de los menús
  updateDropdownVisibility();
};

const hideAssignMenu = () => {
  if (!showAssignMenu.value) return; // Evitar operaciones innecesarias
  
  console.log("Cerrando menú de asignación");
  showAssignMenu.value = false;
  updateDropdownVisibility();
};

const toggleCategoryMenu = (event) => {
  // Evitar propagación para que no cierre inmediatamente
  if (event) {
    event.stopPropagation();
  }
  
  console.log("Toggling menú de categoría");
  showCategoryMenu.value = !showCategoryMenu.value;
  
  // Si se abre el menú de categoría, cerrar el menú de asignación
  if (showCategoryMenu.value) {
    showAssignMenu.value = false;
    showPriorityMenu.value = false;
  }
  
  // Actualizar visualización de los menús
  updateDropdownVisibility();
};

const hideCategoryMenu = () => {
  if (!showCategoryMenu.value) return; // Evitar operaciones innecesarias
  
  console.log("Cerrando menú de categoría");
  showCategoryMenu.value = false;
  updateDropdownVisibility();
};

// Función para actualizar la visualización de los menús desplegables
const updateDropdownVisibility = () => {
  // Actualizar menú de prioridad
  const priorityDropdown = document.querySelector('.priority-btn + .dropdown-content');
  if (priorityDropdown) {
    priorityDropdown.style.display = showPriorityMenu.value ? 'block' : 'none';
  }
  
  // Actualizar menú de asignación
  const assignDropdown = document.querySelector('.assign-btn + .dropdown-content');
  if (assignDropdown) {
    assignDropdown.style.display = showAssignMenu.value ? 'block' : 'none';
  }
  
  // Actualizar menú de categoría
  const categoryDropdown = document.querySelector('.category-btn + .dropdown-content');
  if (categoryDropdown) {
    categoryDropdown.style.display = showCategoryMenu.value ? 'block' : 'none';
  }
};

// Función para cerrar el ticket
const closeTicket = async () => {
  if (!currentTicket.value || !closeReason.value) {
    return;
  }
  
  try {
    isLoading.value = true;
    
    // Preparar el mensaje del motivo de cierre
    let reasonMessage = "Ticket cerrado. Motivo: ";
    
    switch(closeReason.value) {
      case 'solved':
        reasonMessage += "Problema resuelto";
        break;
      case 'customer_request':
        reasonMessage += "Por solicitud del cliente";
        break;
      case 'duplicate':
        reasonMessage += "Ticket duplicado";
        break;
      case 'irrelevant':
        reasonMessage += "Ya no es relevante";
        break;
      case 'other':
        reasonMessage += closeReasonText.value;
        break;
      default:
        reasonMessage += closeReason.value;
    }
    
    // Agregar mensaje sobre el cierre al chat
    try {
      await chatStore.sendMessage(currentTicket.value.id, reasonMessage);
      console.log('Mensaje de cierre enviado al chat correctamente');
    } catch (chatError) {
      console.error('Error al enviar mensaje de cierre al chat:', chatError);
      // Continuar con el proceso aunque falle el mensaje
    }
    
    console.log('Actualizando estado del ticket a cerrado...');
    
    // Actualizar localmente primero para mejorar la experiencia del usuario
    if (currentTicket.value) {
      currentTicket.value.status = 'closed';
    }
    
    // Actualizar el estado del ticket a cerrado usando try/catch separado
    try {
      await ticketStore.updateTicketStatus(currentTicket.value.id, 'closed');
      console.log('Ticket cerrado exitosamente en el backend');
    } catch (updateError) {
      console.error('Error al actualizar estado del ticket en el backend:', updateError);
      showNotification('El ticket se ha cerrado, pero hubo un error al sincronizar con el servidor', 'warning');
    }
    
    // Intentar recargar el ticket para obtener la versión actualizada
    try {
      await ticketStore.fetchTicket(currentTicket.value.id);
      currentTicket.value = ticketStore.currentTicket;
      console.log('Ticket recargado después del cierre:', currentTicket.value);
    } catch (fetchError) {
      console.error('Error al recargar el ticket después del cierre:', fetchError);
      // No mostrar error al usuario ya que el ticket ya está cerrado visualmente
    }
    
    // Cerrar el modal
    showCloseTicketModal.value = false;
    closeReason.value = '';
    closeReasonText.value = '';
    
    // Mostrar notificación
    showNotification('El ticket ha sido cerrado exitosamente');
    
  } catch (error) {
    console.error('Error general al cerrar el ticket:', error);
    showNotification('Error al cerrar el ticket: ' + (error.message || 'Error desconocido'), 'error');
  } finally {
    isLoading.value = false;
  }
};

// Función para actualizar la categoría de un ticket
const updateCategoryTo = async (newCategory) => {
  if (!currentTicket.value) {
    console.error('No se puede actualizar la categoría: ticket actual no está definido');
    hideCategoryMenu();
    return;
  }
  
  try {
    console.log(`Intentando actualizar categoría de ${currentTicket.value.id} a ${newCategory}`);
    
    // Guardar valores actuales para posible reversión
    const ticketId = currentTicket.value.id;
    const originalCategory = currentTicket.value.category;
    
    // Actualizar localmente primero para respuesta inmediata
    currentTicket.value.category = newCategory;
    
    // Mostrar carga
    isLoading.value = true;
    
    // Intentar realizar la actualización en el servidor
    try {
      const response = await ticketStore.updateTicket(ticketId, { category: newCategory });
      console.log('Respuesta del servidor:', response);
      
      if (!response) {
        throw new Error('La respuesta del servidor está vacía');
      }
      
      // Asegurar que el ticket local tiene los datos correctos
      currentTicket.value.category = response.category || newCategory;
      
      // Mensaje de éxito
      showNotification(`La categoría se ha actualizado a "${newCategory}" correctamente`);
      
      // Añadir mensaje al chat sobre el cambio de categoría
      const userName = authStore.user ? `${authStore.user.firstName} ${authStore.user.lastName}` : 'Un administrador';
      const systemMessage = `${userName} ha cambiado la categoría del ticket a "${newCategory}".`;
      await chatStore.sendMessage(ticketId, systemMessage);
      
    } catch (error) {
      console.error('Error al actualizar categoría en el servidor:', error);
      
      // Revertir cambio local si falla
      currentTicket.value.category = originalCategory;
      
      // Mostrar error
      showNotification(`Error al actualizar categoría: ${error.message || 'Error de conexión con el servidor'}`, 'error');
    }
  } catch (error) {
    console.error('Error general en actualización de categoría:', error);
    showNotification('Ocurrió un error al procesar la actualización de categoría', 'error');
  } finally {
    // Ocultar menú y finalizar carga
    hideCategoryMenu();
    isLoading.value = false;
  }
};

// Función para obtener el avatar del agente basado en su ID
const getAgentAvatar = (userId) => {
  if (!userId) return null;
  
  const user = usersStore.users.find(u => u.id === userId);
  if (user && user.avatarUrl) {
    return user.avatarUrl;
  }
  
  // Si no tiene avatar, retornar null para que se muestre el icono por defecto
  return null;
};

// Función para obtener el nombre del remitente del mensaje
const getMessageSenderName = (message) => {
  if (!message.userId) {
    return 'Agente';
  }
  
  const user = usersStore.users.find(u => u.id === message.userId);
  if (user) {
    return `${user.firstName} ${user.lastName}`;
  }
  
  return 'Agente';
};
</script>

<style lang="scss" scoped>
.ticket-detail {
  max-width: 1300px;
  margin: 0 auto;
  background: linear-gradient(to bottom, #f8faff, #eef2ff);
  border-radius: 10px;
  padding: 1.5rem;
  box-shadow: 0 5px 15px rgba(35, 38, 110, 0.05);

  .loading, .error {
    text-align: center;
    padding: 2rem;
  }

  .error {
    color: #7e22ce;
  }

  .ticket-detail-layout {
    display: flex;
    flex-direction: column;
    min-height: calc(100vh - 200px);
  }

  .ticket-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 1.5rem;
    padding-bottom: 1rem;
    border-bottom: 2px solid #c7d2fe;
  }

  .ticket-header .ticket-header-main {
    flex: 1;
  }

  .ticket-header .ticket-header-main h1 {
    margin: 0 0 0.5rem 0;
    color: #1e293b;
    font-weight: 600;
    font-size: 1.75rem;
  }

  .ticket-header .ticket-quick-meta {
    display: flex;
    flex-wrap: wrap;
    gap: 0.75rem;
    align-items: center;
  }

  .ticket-header .ticket-id, 
  .ticket-header .ticket-status, 
  .ticket-header .ticket-priority {
    padding: 0.35rem 0.6rem;
    border-radius: 6px;
    background-color: #eef2ff;
    font-size: 0.85rem;
    display: flex;
    align-items: center;
    box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
  }

  .ticket-header .ticket-id strong, 
  .ticket-header .ticket-status strong, 
  .ticket-header .ticket-priority strong {
    margin-right: 0.5rem;
    font-weight: 600;
    color: #444;
  }

  /* Estilos mejorados para los botones en el encabezado */
  .ticket-header .header-actions {
    display: flex;
    gap: 0.75rem;
    margin-left: auto;
    flex-wrap: wrap;
    padding-right: 2rem; /* Añadir espacio a la derecha para separar del botón de sidebar */
  }

  .ticket-header .header-button {
    position: relative;
  }

  .ticket-header .header-button .filter-btn {
    padding: 0.35rem 0.75rem;
    font-size: 0.85rem;
    background-color: #eef2ff;
    color: #4f46e5;
    border: 1px solid #ddd6fe;
    border-radius: 6px;
    display: flex;
    align-items: center;
    gap: 0.4rem;
    height: 32px;
    transition: all 0.2s ease;
  }

  .ticket-header .header-button .filter-btn:hover {
    background-color: #ddd6fe;
    transform: translateY(-1px);
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
  }

  .ticket-header .header-button .filter-btn i {
    color: #4f46e5;
    font-size: 0.85rem;
  }

  .ticket-header .header-button .filter-btn span {
    font-weight: 500;
  }

  .ticket-header .header-button .filter-btn.active {
    background-color: #4f46e5;
    color: white;
    border-color: #4f46e5;
  }

  .ticket-header .header-button .filter-btn.active i {
    color: white;
  }

  .ticket-header .header-button .filter-btn.close-btn {
    background-color: #f1f5f9;
    color: #7e22ce;
    border-color: #e2e8f0;
  }

  .ticket-header .header-button .filter-btn.close-btn:hover {
    background-color: #e2e8f0;
  }

  .ticket-header .header-button .filter-btn.close-btn i {
    color: #7e22ce;
  }

  .ticket-header .header-button .dropdown-content {
    right: 0;
    left: auto;
    min-width: 220px;
    z-index: 1000;
  }

  .ticket-header .toggle-sidebar-btn {
    display: flex;
    justify-content: center;
    align-items: center;
    width: 36px;
    height: 36px;
    border-radius: 50%;
    background: #e0e7ff;
    color: #4f46e5;
    border: none;
    cursor: pointer;
    transition: all 0.2s ease;
  }

  .ticket-header .toggle-sidebar-btn:hover {
    background: #c7d2fe;
  }

  .ticket-header .toggle-sidebar-btn i {
    font-size: 1rem;
  }

  .ticket-content-wrapper {
    display: flex;
    flex: 1;
    position: relative;
    overflow: hidden;

    .main-column {
      flex: 1;
      margin-right: 1rem;
      transition: margin-right 0.3s ease;
      
      &.full-width {
        margin-right: 0;
      }
    }

    .sidebar {
      width: 350px;
      background: white;
      border-radius: 12px;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.08);
      position: absolute;
      right: -350px;
      top: 0;
      bottom: 0;
      transition: right 0.3s ease;
      display: flex;
      flex-direction: column;
      z-index: 100;
      border: 1px solid #e2e8f0;
      
      &.sidebar-open {
        right: 0;
      }
      
      .sidebar-header {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem 1.25rem;
        border-bottom: 1px solid #e2e8f0;
        
        h3 {
          margin: 0;
          font-size: 1.1rem;
          font-weight: 600;
          color: #1e293b;
        }
        
        .close-sidebar-btn {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          background: #f1f5f9;
          border: none;
          color: #64748b;
          cursor: pointer;
          transition: all 0.2s ease;
          
          &:hover {
            background: #e2e8f0;
            color: #1e293b;
          }
          
          i {
            font-size: 0.9rem;
          }
        }
      }
      
      .sidebar-content {
        flex: 1;
        overflow-y: auto;
        padding: 1.25rem;
        
        .sidebar-section {
          margin-bottom: 1.5rem;
          
          &:last-child {
            margin-bottom: 0;
          }
          
          .section-title {
            display: flex;
            align-items: center;
            margin: 0 0 0.75rem 0;
            font-size: 1rem;
            font-weight: 600;
            color: #475569;
            
            i {
              margin-right: 0.5rem;
              color: #4f46e5;
              font-size: 1rem;
            }
          }
          
          .ticket-info {
            .info-item {
              display: flex;
              padding: 0.5rem 0;
              border-bottom: 1px solid #f1f5f9;
              
              &:last-child {
                border-bottom: none;
              }
              
              .info-label {
                width: 40%;
                color: #64748b;
                font-size: 0.9rem;
              }
              
              .info-value {
                flex: 1;
                font-weight: 500;
                color: #1e293b;
                font-size: 0.9rem;
              }
            }
          }
          
          .ticket-description {
            font-size: 0.95rem;
            line-height: 1.6;
            color: #334155;
            white-space: pre-wrap;
            margin: 0;
          }
          
          .assignment-info {
            .assigned-user {
              font-weight: 500;
              color: #4F46E5;
              background-color: #EEF2FF;
              padding: 0.35rem 0.75rem;
              border-radius: 4px;
              display: inline-block;
              font-size: 0.9rem;
            }
            
            .not-assigned {
              font-style: italic;
              color: #6B7280;
              background-color: #F3F4F6;
              padding: 0.35rem 0.75rem;
              border-radius: 4px;
              display: inline-block;
              font-size: 0.9rem;
            }
          }
          
          .admin-buttons {
            display: flex;
            flex-direction: column;
            gap: 0.75rem;
            
            .admin-button {
              position: relative;
            }
            
            .filter-btn {
              width: 100%;
              text-align: left;
              justify-content: flex-start;
              font-size: 0.9rem;
            }
          }
        }
      }
    }
  }

  .chat-section {
    background: white;
    padding: 1.5rem;
    border-radius: 12px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.05);
    border: 1px solid #e2e8f0;
    height: 100%;
    display: flex;
    flex-direction: column;

    h2 {
      margin: 0 0 1.25rem;
      display: flex;
      align-items: center;
      font-size: 1.25rem;
      
      .connected-tag {
        margin-left: 1rem;
        font-size: 0.8rem;
        color: #4caf50;
        display: flex;
        align-items: center;
        background: #e8f5e9;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        
        i {
          margin-right: 0.25rem;
        }
      }
      
      .ticket-closed-badge {
        background-color: #ef4444;
        color: white;
        padding: 0.25rem 0.5rem;
        border-radius: 4px;
        font-size: 0.8rem;
        margin-left: 0.5rem;
        display: inline-flex;
        align-items: center;
        
        i {
          margin-right: 0.25rem;
        }
      }
    }

    .messages {
      flex: 1;
      min-height: 300px;
      max-height: 500px;
      overflow-y: auto;
      margin-bottom: 1.5rem;
      padding: 1rem;
      background: #f9f9f9;
      border-radius: 8px;
      display: flex;
      flex-direction: column;
      gap: 1rem;

      .no-messages {
        text-align: center;
        color: #666;
        padding: 2rem;
        font-style: italic;
      }

      .message {
        display: flex;
        padding: 0;
        width: 90%;
        position: relative;
        
        &.client-message {
          align-self: flex-start;
          
          .message-content {
            background-color: #e0e7ff;
            border-radius: 0 12px 12px 12px;
          }
        }
        
        &.agent-message {
          align-self: flex-end;
          flex-direction: row-reverse;
          
          .message-content {
            background-color: #ddd6fe;
            border-radius: 12px 0 12px 12px;
          }
          
          .message-avatar {
            margin-left: 12px;
            margin-right: 0;
          }
        }
        
        .message-avatar {
          margin-right: 12px;
          
          .avatar-circle {
            width: 38px;
            height: 38px;
            border-radius: 50%;
            display: flex;
            align-items: center;
            justify-content: center;
            background-color: #f1f5f9;
            border: 2px solid #fff;
            box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
            overflow: hidden;
            
            &.client-avatar {
              background-color: #e0e7ff;
              color: #4f46e5;
            }
            
            &.agent-avatar {
              background-color: #ddd6fe;
              color: #7e22ce;
            }
            
            i {
              font-size: 1.1rem;
            }
            
            img {
              width: 100%;
              height: 100%;
              object-fit: cover;
            }
          }
        }
        
        .message-content {
          padding: 1rem;
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.05);
          flex: 1;

          .message-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 0.5rem;

            .sender {
              font-weight: 600;
              color: #4f46e5;
            }

            .timestamp {
              font-size: 0.8rem;
              color: #666;
            }
          }

          .content {
            margin: 0;
            white-space: pre-wrap;
            word-break: break-word;
          }
        }
      }
    }

    .message-form {
      display: flex;
      flex-direction: column;
      gap: 1rem;

      textarea {
        resize: vertical;
        min-height: 100px;
        padding: 0.75rem;
        border: 1px solid #ddd;
        border-radius: 8px;
        font-family: inherit;
        font-size: 1rem;
        
        &:focus {
          outline: none;
          border-color: #4f46e5;
          box-shadow: 0 0 0 2px rgba(79, 70, 229, 0.2);
        }
      }

      button {
        align-self: flex-end;
        padding: 0.6rem 1.25rem;
        border: none;
        border-radius: 8px;
        background: linear-gradient(to right, #4f46e5, #6366f1);
        color: white;
        font-weight: 600;
        cursor: pointer;
        transition: all 0.2s;
        display: flex;
        align-items: center;
        gap: 0.5rem;

        &:hover {
          background: linear-gradient(to right, #4338ca, #4f46e5);
          box-shadow: 0 2px 4px rgba(37, 99, 235, 0.3);
        }
        
        i {
          font-size: 1rem;
        }
      }
    }
    
    .chat-closed-message {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 1.5rem;
      background-color: #f9fafb;
      border: 1px dashed #d1d5db;
      border-radius: 0.5rem;
      margin-top: 1rem;
      color: #6b7280;
      
      i {
        font-size: 1.5rem;
        margin-bottom: 0.5rem;
        color: #ef4444;
      }
      
      p {
        margin: 0;
        font-style: italic;
      }
    }
  }
}

/* Estilos para el sistema de notificaciones */
.notification-toast {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 1000;
  min-width: 300px;
  max-width: 450px;
  padding: 15px 20px;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  display: flex;
  align-items: center;
  justify-content: space-between;
  transition: all 0.3s ease;
}

.notification-toast.success {
  background-color: #e0e7ff;
  border-left: 4px solid #4f46e5;
  color: #3730a3;
}

.notification-toast.error {
  background-color: #ede9fe;
  border-left: 4px solid #7c3aed;
  color: #5b21b6;
}

.notification-toast.warning {
  background-color: #dbeafe;
  border-left: 4px solid #3b82f6;
  color: #1e40af;
}

/* Transiciones para el toast */
.toast-fade-enter-active,
.toast-fade-leave-active {
  transition: all 0.3s ease;
}

.toast-fade-enter-from {
  transform: translateX(100%);
  opacity: 0;
}

.toast-fade-leave-to {
  transform: translateY(-20px);
  opacity: 0;
}

.notification-message {
  flex-grow: 1;
  margin-right: 15px;
  font-weight: 500;
}

.close-notification {
  background: transparent;
  border: none;
  cursor: pointer;
  font-size: 18px;
  color: inherit;
  opacity: 0.7;
  padding: 0;
  display: flex;
  align-items: center;
  justify-content: center;
  width: 24px;
  height: 24px;
}

.close-notification:hover {
  opacity: 1;
}

/* Estilos para las badges y elementos compartidos */
.status-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  text-transform: capitalize;
  display: inline-flex;
  align-items: center;
  font-size: 0.85rem;
  font-weight: 500;
  
  &::before {
    content: "";
    display: inline-block;
    width: 8px;
    height: 8px;
    border-radius: 50%;
    margin-right: 0.5rem;
  }
  
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
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  text-transform: capitalize;
  display: inline-flex;
  align-items: center;
  font-size: 0.85rem;
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

/* Estilos para botones de acciones admin */
.filter-btn {
  background-color: #f1f5f9;
  border: 1px solid #e2e8f0;
  color: #475569;
  font-size: 0.9rem;
  font-weight: 500;
  padding: 0.65rem 1.2rem;
  border-radius: 8px;
  cursor: pointer;
  transition: all 0.2s ease;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  
  &:hover {
    background-color: #e2e8f0;
    transform: translateY(-2px);
  }
  
  &.active {
    background-color: #4f46e5;
    color: white;
    border-color: transparent;
    font-weight: 600;
    box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
  }
  
  &.priority-btn i {
    color: #3b82f6;
  }
  
  &.category-btn i {
    color: #8b5cf6;
  }
  
  &.assign-btn i {
    color: #10b981;
  }
  
  &.close-btn i {
    color: #f43f5e;
  }
  
  &.active i {
    color: white;
  }
}

/* Estilos para menús desplegables */
.dropdown-content {
  position: absolute;
  z-index: 1000;
  background-color: white;
  border: 1px solid #e2e8f0;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
  min-width: 200px;
  display: none;
  margin-top: 5px;
  animation: fadeInDown 0.3s ease both;
  right: 0;
  
  &.show {
    display: block;
  }
  
  .dropdown-item {
    padding: 0.75rem 1rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    color: #475569;
    font-weight: 500;
    font-size: 0.9rem;
    transition: all 0.2s ease;
    
    &:hover {
      background-color: #f8fafc;
    }
  }
  
  .dropdown-divider {
    height: 1px;
    background-color: #e2e8f0;
    margin: 8px 0;
  }
}

/* Animación para menús desplegables */
@keyframes fadeInDown {
  from {
    opacity: 0;
    transform: translateY(-10px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

/* Indicadores de prioridad y categoría */
.priority-indicator {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 10px;
}

.priority-indicator.low {
  background-color: #0369a1;
}

.priority-indicator.medium {
  background-color: #4338ca;
}

.priority-indicator.high {
  background-color: #3730a3;
}

.priority-indicator.urgent {
  background-color: #5b21b6;
}

.category-indicator {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 10px;
  background-color: #8b5cf6;
}

/* Estilos para el modal */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 12px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
  border-bottom: 1px solid #e2e8f0;
  
  h3 {
    margin: 0;
    color: #111827;
    font-size: 1.25rem;
    font-weight: 600;
  }
  
  .modal-close {
    background: transparent;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #6b7280;
    line-height: 1;
  }
}

.modal-body {
  padding: 1.5rem;
  
  p {
    margin-top: 0;
    margin-bottom: 1rem;
    color: #4b5563;
  }
}

.close-reason-select {
  width: 100%;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 1rem;
  color: #1f2937;
}

.close-reason-textarea {
  width: 100%;
  height: 100px;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 8px;
  margin-bottom: 1rem;
  font-size: 1rem;
  color: #1f2937;
  resize: vertical;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 0.75rem;
  margin-top: 1rem;
  
  .btn {
    padding: 0.6rem 1.25rem;
    border-radius: 8px;
    font-weight: 500;
    cursor: pointer;
    font-size: 0.95rem;
    
    &.btn-secondary {
      background-color: #f3f4f6;
      color: #374151;
      border: 1px solid #d1d5db;
      
      &:hover {
        background-color: #e5e7eb;
      }
    }
    
    &.btn-primary {
      background-color: #ef4444;
      color: white;
      border: none;
      
      &:hover {
        background-color: #dc2626;
      }
      
      &:disabled {
        background-color: #fca5a5;
        cursor: not-allowed;
      }
    }
  }
}

/* Estilos responsivos */
@media (max-width: 768px) {
  .ticket-detail {
    padding: 1rem;
  }
  
  .ticket-detail-layout {
    min-height: calc(100vh - 150px);
  }
  
  .ticket-content-wrapper {
    flex-direction: column;
    
    .main-column {
      margin-right: 0;
      margin-bottom: 1rem;
    }
    
    .sidebar {
      position: fixed;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      width: 100%;
      max-width: 100%;
      z-index: 1000;
      transform: translateX(100%);
      
      &.sidebar-open {
        transform: translateX(0);
      }
    }
  }
  
  .ticket-header {
    flex-direction: column;
    align-items: flex-start;
    
    .ticket-header-main {
      width: 100%;
      margin-bottom: 1rem;
    }
    
    .toggle-sidebar-btn {
      position: absolute;
      top: 1rem;
      right: 1rem;
    }
  }
  
  .chat-section {
    .messages {
      max-height: 400px;
    }
    
    .message {
      width: 90%;
    }
  }
}

/* Estilo para el mensaje en la sección de acciones de la barra lateral */
.sidebar-info {
  padding: 1rem;
  background-color: #f1f5f9;
  border-radius: 8px;
  border: 1px dashed #cbd5e1;
}

.sidebar-note {
  margin: 0;
  font-size: 0.9rem;
  color: #64748b;
  font-style: italic;
  text-align: center;
}
</style>
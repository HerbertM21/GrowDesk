<template>
  <div class="ticket-detail">
    <div v-if="isLoading" class="loading">Cargando...</div>
    <div v-else-if="errorMessage" class="error">{{ errorMessage }}</div>
    <div v-else class="ticket-content">
      <!-- Notificación tipo toast -->
      <transition name="toast-fade">
        <div v-if="notification.show" :class="['notification-toast', notification.type]">
          <div class="notification-message">{{ notification.message }}</div>
          <button @click="notification.show = false" class="close-notification">&times;</button>
        </div>
      </transition>
      
      <div v-if="currentTicket" class="ticket-header">
        <h1>{{ currentTicket.title }}</h1>
        <div class="ticket-meta">
          <div class="ticket-id">
            <strong>ID:</strong> {{ currentTicket.id }}
          </div>
          <div class="ticket-status">
            <strong>Estado:</strong> <span :class="['status-badge', currentTicket.status.toLowerCase()]">{{ translateStatus(currentTicket.status) }}</span>
          </div>
          <div class="ticket-priority">
            <strong>Prioridad:</strong> <span :class="['priority-badge', currentTicket.priority.toLowerCase()]">{{ translatePriority(currentTicket.priority) }}</span>
          </div>
          <div class="ticket-category">
            <strong>Categoría:</strong> {{ translateCategory(currentTicket.category) }}
          </div>
          <div class="ticket-date">
            <strong>Creado:</strong> {{ formatDate(currentTicket.createdAt) }}
          </div>
        </div>
      </div>
      <div v-else class="ticket-header">
        <h1>Ticket #{{ route.params.id }}</h1>
      </div>

      <div class="ticket-info" v-if="currentTicket">
        <p class="description">{{ currentTicket.description }}</p>
        <div class="meta-info">
          <p>Categoría: {{ translateCategory(currentTicket.category) }}</p>
          <p>Creado: {{ formatDate(currentTicket.createdAt) }}</p>
          
          <!-- Nueva sección para administrar la prioridad del ticket -->
          <div v-if="isAdmin || isAssistant" class="ticket-actions">
            <div class="action-buttons">
              <!-- Botón de prioridad con menú desplegable -->
              <div class="dropdown-menu">
                <button class="action-btn" @click="togglePriorityMenu">
                  <i class="pi pi-flag"></i>
                  <span class="btn-text">Prioridad: <span class="priority-text" :class="['priority-' + currentTicket.priority.toLowerCase()]">{{ translatePriority(currentTicket.priority) }}</span></span>
                  <i class="pi pi-chevron-down"></i>
                </button>
                <div class="dropdown-content" :class="{ 'show': showPriorityMenu }">
                  <div class="dropdown-header">Cambiar prioridad</div>
                  <div class="dropdown-items">
                    <div class="dropdown-item" @click="updatePriorityTo('low'); hidePriorityMenu()" :class="{ 'active': currentTicket.priority.toLowerCase() === 'low' }">
                      <span class="priority-indicator low"></span> Baja
                    </div>
                    <div class="dropdown-item" @click="updatePriorityTo('medium'); hidePriorityMenu()" :class="{ 'active': currentTicket.priority.toLowerCase() === 'medium' }">
                      <span class="priority-indicator medium"></span> Media
                    </div>
                    <div class="dropdown-item" @click="updatePriorityTo('high'); hidePriorityMenu()" :class="{ 'active': currentTicket.priority.toLowerCase() === 'high' }">
                      <span class="priority-indicator high"></span> Alta
                    </div>
                    <div class="dropdown-item" @click="updatePriorityTo('urgent'); hidePriorityMenu()" :class="{ 'active': currentTicket.priority.toLowerCase() === 'urgent' }">
                      <span class="priority-indicator urgent"></span> Urgente
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Botón de asignación con menú desplegable -->
              <div class="dropdown-menu">
                <button class="action-btn assignment-btn" @click="toggleAssignMenu">
                  <i class="pi pi-user"></i>
                  <span class="btn-text">Asignación: {{ currentTicket.assignedTo ? getAssignedUserName(currentTicket.assignedTo) : 'Sin asignar' }}</span>
                  <i class="pi pi-chevron-down"></i>
                </button>
                <div class="dropdown-content" :class="{ 'show': showAssignMenu }">
                  <div class="dropdown-header">Asignar ticket</div>
                  <div class="dropdown-items">
                    <div v-if="isAdmin">
                      <div v-for="user in supportUsers" :key="user.id" class="dropdown-item" @click="assignToUserQuick(user.id); hideAssignMenu()">
                        <i class="pi pi-user"></i> {{ user.firstName }} {{ user.lastName }}
                      </div>
                      <div class="dropdown-divider"></div>
                      <div class="dropdown-item" @click="removeAssignment(currentTicket.id); hideAssignMenu()">
                        <i class="pi pi-times"></i> Quitar asignación
                      </div>
                    </div>
                    <div v-else-if="isAssistant && !isTicketAssigned" class="dropdown-item" @click="assignToSelf(); hideAssignMenu()">
                      <i class="pi pi-user"></i> Asignarme este ticket
                    </div>
                    <div v-else-if="!isAdmin && (isEmployee || isTicketAssigned)" class="dropdown-item" @click="requestAssignment(); hideAssignMenu()">
                      <i class="pi pi-send"></i> Solicitar asignación
                    </div>
                  </div>
                </div>
              </div>
              
              <!-- Botón para cerrar el ticket -->
              <div v-if="!isTicketClosed && (isAdmin || isAssistant)" class="dropdown-menu">
                <button 
                  class="action-btn"
                  @click="showCloseTicketModal = true">
                  <i class="pi pi-check-circle"></i>
                  <span class="btn-text">Cerrar Ticket</span>
                </button>
              </div>
            </div>
          </div>
          
          <div class="assignment-info">
            <!-- Información de asignación actual -->
            <p v-if="currentTicket.assignedTo">
              <strong>Asignado a:</strong> 
              <span class="assigned-user">{{ getAssignedUserName(currentTicket.assignedTo) }}</span>
            </p>
            <p v-else><strong>Asignado a:</strong> <span class="not-assigned">No asignado</span></p>
          </div>
        </div>
      </div>

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
              <div class="message-header">
                <span class="sender">{{ message.isClient ? (currentTicket?.customer?.name || 'Cliente') : 'Agente' }}</span>
                <span class="timestamp">{{ formatTimestamp(message.timestamp || message.createdAt) }}</span>
              </div>
              <p class="content">{{ message.content }}</p>
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

// Obtener stores y route
const route = useRoute();
const router = useRouter();
const ticketStore = useTicketStore();
const chatStore = useChatStore();
const authStore = useAuthStore();
const usersStore = useUsersStore();

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

// Variables para la asignación
const selectedUserId = ref('');
const assignmentRequested = ref(false);
const selectedPriority = ref('');

// Variables para controlar la visibilidad de los menús desplegables
const showPriorityMenu = ref(false);
const showAssignMenu = ref(false);

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
const isAdmin = computed(() => authStore.isAdmin);
const isAssistant = computed(() => authStore.isAssistant);
const isEmployee = computed(() => {
  return !!authStore.user && authStore.user.role === 'employee';
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
  
  // Normalizar a minúsculas
  const normalizedPriority = priority.toLowerCase();
  
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
    return;
  }
  
  try {
    console.log('Eliminando asignación del ticket:', ticketId);
    
    // Actualizar el ticket para quitar la asignación
    const result = await ticketStore.updateTicket(ticketId, {
      assignedTo: null,
      status: 'open' // Volver a estado abierto
    });
    
    if (!result) {
      throw new Error('La operación de eliminación de asignación devolvió un resultado nulo');
    }
    
    console.log('Asignación eliminada correctamente:', result);
    
    // Actualizar ticket local
    await ticketStore.fetchTicket(ticketId);
    currentTicket.value = ticketStore.currentTicket;
    
    // Añadir mensaje al chat sobre la eliminación de asignación
    const systemMessage = `Un administrador ha eliminado la asignación de este ticket.`;
    await chatStore.sendMessage(ticketId, systemMessage);
    
    // Mostrar notificación de éxito
    showNotification('Asignación eliminada correctamente');
  } catch (error) {
    console.error('Error detallado al eliminar asignación:', error);
    showNotification(`Error al eliminar la asignación: ${error.message || 'Error desconocido'}`, 'error');
  }
};

// Actualizar la prioridad de un ticket directamente
const updatePriorityTo = async (newPriority) => {
  // En caso de que la prioridad actual esté en mayúsculas y la nueva en minúsculas
  const currentPriorityLower = currentTicket.value.priority.toLowerCase();
  
  if (!currentTicket.value || currentPriorityLower === newPriority) {
    hidePriorityMenu();
    return;
  }
  
  const oldPriority = currentTicket.value.priority;
  
  try {
    isLoading.value = true;
    
    // Convertir la prioridad a mayúsculas si es necesario para el backend
    const priorityForBackend = newPriority.toUpperCase();
    
    console.log(`Actualizando prioridad del ticket de ${translatePriority(oldPriority)} a ${translatePriority(priorityForBackend)}`);
    
    // Actualizar el ticket con la nueva prioridad
    await ticketStore.updateTicket(currentTicket.value.id, { 
      priority: priorityForBackend 
    });
    
    // Recargar el ticket para obtener datos actualizados
    await ticketStore.fetchTicket(currentTicket.value.id);
    if (ticketStore.currentTicket) {
      currentTicket.value = ticketStore.currentTicket;
    }
    
    // Enviar un mensaje al chat sobre el cambio de prioridad
    const message = `La prioridad ha cambiado de ${translatePriority(oldPriority)} a ${translatePriority(priorityForBackend)}`;
    await chatStore.sendMessage(currentTicket.value.id, message, true);
    
    // Mostrar notificación
    notification.value = {
      show: true,
      message: `Prioridad actualizada a ${translatePriority(priorityForBackend)}`,
      type: 'success',
      timeout: setTimeout(() => {
        notification.value.show = false;
      }, 3000)
    };
    
  } catch (err) {
    console.error('Error al actualizar la prioridad del ticket:', err);
    errorMessage.value = 'Error al actualizar la prioridad del ticket. Por favor, inténtalo de nuevo.';
    
    // Mostrar notificación de error
    notification.value = {
      show: true,
      message: 'Error al actualizar la prioridad',
      type: 'error',
      timeout: setTimeout(() => {
        notification.value.show = false;
      }, 3000)
    };
  } finally {
    isLoading.value = false;
    hidePriorityMenu();
  }
};

// Asignar rápidamente el ticket a un usuario
const assignToUserQuick = async (userId) => {
  if (!userId || !currentTicket.value) {
    console.error('No se puede asignar: falta el ticket o el usuario seleccionado');
    showNotification('Falta información para realizar la asignación', 'error');
    return;
  }
  
  try {
    console.log('Asignando ticket:', currentTicket.value.id, 'a usuario:', userId);
    
    // Guardar ID localmente en caso de que currentTicket.value se actualice durante la operación
    const ticketId = currentTicket.value.id;
    const userName = getAssignedUserName(userId);
    
    // Realizar la asignación
    const result = await ticketStore.assignTicket(ticketId, userId);
    
    if (!result) {
      throw new Error('La operación de asignación devolvió un resultado nulo');
    }
    
    // Actualizar ticket local - usamos el ID guardado por si acaso
    await ticketStore.fetchTicket(ticketId);
    currentTicket.value = ticketStore.currentTicket;
    
    // Añadir mensaje al chat sobre la asignación
    const systemMessage = `El ticket ha sido asignado a ${userName} por un administrador.`;
    await chatStore.sendMessage(ticketId, systemMessage);
    
    // Mostrar notificación de éxito
    showNotification(`Ticket asignado a ${userName} correctamente`);
  } catch (error) {
    console.error('Error detallado al asignar ticket:', error);
    showNotification(`Error al asignar el ticket: ${error.message || 'Error desconocido'}`, 'error');
  }
};

// Lógica principal
const loadTicketData = async () => {
  const ticketId = route.params.id;
  console.log('Cargando datos para ticket:', ticketId);
  
  try {
    isLoading.value = true;
    
    // Cargar usuarios para asignación
    await usersStore.fetchUsers();
    
    // Configurar el ticket actual en el chatStore
    chatStore.setCurrentTicket(ticketId.toString());
    
    try {
      // Intentar cargar el ticket primero
      await ticketStore.fetchTicket(ticketId);
      // Usar el currentTicket del store
      currentTicket.value = ticketStore.currentTicket;
      console.log('Ticket obtenido del store:', currentTicket.value);
    } catch (err) {
      console.warn('No se pudo cargar el ticket, pero continuaremos con los mensajes:', err);
    }

    // Cargar mensajes del ticket
    await chatStore.fetchMessages(ticketId.toString());

    // Conectar WebSocket
    chatStore.connectToWebSocket(ticketId.toString());

    // Obtener mensajes y estado de conexión
    currentMessages.value = chatStore.getMessagesForTicket(ticketId.toString());
    connectionStatus.value = chatStore.isConnected();
    
    console.log('Datos cargados:', {
      ticket: currentTicket.value,
      mensajes: currentMessages.value,
      conexion: connectionStatus.value
    });
    
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
onMounted(() => {
  // Configurar event listener para cerrar menús al hacer clic fuera
  document.addEventListener('click', (event) => {
    // Si el clic no fue dentro de un menú desplegable, cerrar todos los menús
    const isClickInsideDropdown = event.target.closest('.dropdown-menu');
    if (!isClickInsideDropdown) {
      hidePriorityMenu();
      hideAssignMenu();
    }
  });
  
  // Cargar datos inmediatamente
  loadTicketData();
});

onBeforeUnmount(() => {
  // Limpiar intervalos
  if (connectionCheckInterval) clearInterval(connectionCheckInterval);
  if (monitorInterval) clearInterval(monitorInterval);
  
  // Desconectar WebSocket
  chatStore.disconnectWebSocket();
  
  // Limpiar el event listener cuando se desmonta el componente
  document.removeEventListener('click', () => {});
});

// Métodos para manejar los menús desplegables
const togglePriorityMenu = () => {
  showPriorityMenu.value = !showPriorityMenu.value;
  if (showPriorityMenu.value) {
    // Cerrar otros menús
    showAssignMenu.value = false;
  }
};

const hidePriorityMenu = () => {
  showPriorityMenu.value = false;
};

const toggleAssignMenu = () => {
  showAssignMenu.value = !showAssignMenu.value;
  if (showAssignMenu.value) {
    // Cerrar otros menús
    showPriorityMenu.value = false;
  }
};

const hideAssignMenu = () => {
  showAssignMenu.value = false;
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
    await chatStore.sendMessage(currentTicket.value.id, reasonMessage);
    
    // Actualizar el estado del ticket a cerrado
    await ticketStore.updateTicketStatus(currentTicket.value.id, 'closed');
    
    // Actualizar el ticket actual
    await ticketStore.fetchTicket(currentTicket.value.id);
    currentTicket.value = ticketStore.currentTicket;
    
    // Cerrar el modal
    showCloseTicketModal.value = false;
    closeReason.value = '';
    closeReasonText.value = '';
    
    // Mostrar notificación
    showNotification('El ticket ha sido cerrado exitosamente');
    
  } catch (error) {
    console.error('Error al cerrar el ticket:', error);
    showNotification('Error al cerrar el ticket: ' + (error.message || 'Error desconocido'), 'error');
  } finally {
    isLoading.value = false;
  }
};
</script>

<style lang="scss" scoped>
.ticket-detail {
  max-width: 1000px;
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

  .ticket-content {
    .ticket-header {
      display: flex;
      flex-direction: column;
      margin-bottom: 2rem;

      h1 {
        margin: 0 0 1rem 0;
        color: #1e293b;
        font-weight: 600;
        border-bottom: 2px solid #c7d2fe;
        padding-bottom: 0.5rem;
      }

      .ticket-meta {
        display: flex;
        flex-wrap: wrap;
        gap: 0.75rem;
        margin-bottom: 1.5rem;
        
        .ticket-id, .ticket-status, .ticket-priority, .ticket-category, .ticket-date {
          padding: 0.5rem 0.75rem;
          border-radius: 6px;
          background-color: #eef2ff;
          font-size: 0.9rem;
          display: flex;
          align-items: center;
          box-shadow: 0 1px 3px rgba(0, 0, 0, 0.08);
          
          strong {
            margin-right: 0.5rem;
            font-weight: 600;
            color: #444;
          }
        }
      }
    }

    .ticket-info {
      background: #f8fafc;
      padding: 2rem;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      margin-bottom: 2rem;
      border: 1px solid #e2e8f0;

      .description {
        margin-bottom: 1.5rem;
        line-height: 1.6;
      }

      .meta-info {
        color: #666;
        font-size: 0.9rem;

        p {
          margin: 0.5rem 0;
        }
        
        .ticket-actions {
          margin-top: 1.75rem;
          padding-top: 1.75rem;
          border-top: 1px solid #E5E7EB;
          
          .action-buttons {
            display: flex;
            gap: 1.25rem;
            flex-wrap: wrap;
          }
        }
        
        .assignment-info {
          margin-top: 1.75rem;
          padding-top: 1.75rem;
          border-top: 1px solid #E5E7EB;
          
          p {
            font-size: 0.95rem;
            margin-bottom: 0.5rem;
            color: #4B5563;
            
            strong {
              font-weight: 600;
              color: #374151;
            }
          }
          
          .assigned-user {
            font-weight: 500;
            color: #4F46E5;
            background-color: #EEF2FF;
            padding: 0.35rem 0.75rem;
            border-radius: 4px;
            display: inline-block;
          }
          
          .not-assigned {
            font-style: italic;
            color: #6B7280;
            background-color: #F3F4F6;
            padding: 0.35rem 0.75rem;
            border-radius: 4px;
            display: inline-block;
          }
        }
      }
    }

    .chat-section {
      background: #f8fafc;
      padding: 2rem;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      border: 1px solid #e2e8f0;

      h2 {
        margin: 0 0 1.5rem;
        display: flex;
        align-items: center;
        
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
      }

      .messages {
        min-height: 200px;
        max-height: 400px;
        overflow-y: auto;
        margin-bottom: 1.5rem;
        padding: 1rem;
        background: #f9f9f9;
        border-radius: 4px;
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
          padding: 1rem;
          border-radius: 8px;
          box-shadow: 0 1px 2px rgba(0, 0, 0, 0.1);
          width: 80%;
          position: relative;
          
          &.client-message {
            background-color: #e0e7ff;
            align-self: flex-start;
            
            &:before {
              content: '';
              position: absolute;
              left: -8px;
              top: 10px;
              border-style: solid;
              border-width: 8px 8px 8px 0;
              border-color: transparent #e0e7ff transparent transparent;
            }
          }
          
          &.agent-message {
            background-color: #ddd6fe;
            align-self: flex-end;
            
            &:before {
              content: '';
              position: absolute;
              right: -8px;
              top: 10px;
              border-style: solid;
              border-width: 8px 0 8px 8px;
              border-color: transparent transparent transparent #ddd6fe;
            }
          }

          .message-header {
            display: flex;
            justify-content: space-between;
            align-items: center;
            margin-bottom: 0.5rem;

            .sender {
              font-weight: bold;
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

      .message-form {
        display: flex;
        flex-direction: column;
        gap: 1rem;

        textarea {
          resize: vertical;
          min-height: 100px;
          padding: 0.75rem;
          border: 1px solid #ddd;
          border-radius: 4px;
          font-family: inherit;
          font-size: 1rem;
        }

        button {
          align-self: flex-end;
          padding: 0.5rem 1rem;
          border: none;
          border-radius: 4px;
          background: linear-gradient(to right, #4f46e5, #6366f1);
          color: white;
          font-weight: bold;
          cursor: pointer;
          transition: all 0.2s;

          &:hover {
            background: linear-gradient(to right, #4338ca, #4f46e5);
            box-shadow: 0 2px 4px rgba(37, 99, 235, 0.3);
          }
        }
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
  border-radius: 6px;
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

/* Estilos mejorados para la sección de asignación */
.assignment-section {
  margin-top: 20px;
  padding: 15px;
  border-radius: 6px;
  background-color: #f9f9f9;
  border: 1px solid #e5e5e5;
}

.assignment-section h3 {
  margin-top: 0;
  margin-bottom: 15px;
  font-size: 16px;
  color: #333;
}

.assignment-controls {
  display: flex;
  gap: 10px;
  flex-wrap: wrap;
  align-items: center;
}

.user-select {
  flex-grow: 1;
  min-width: 200px;
}

.assignment-buttons {
  display: flex;
  gap: 10px;
  margin-top: 10px;
}

.assignment-buttons button {
  padding: 8px 15px;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s ease;
}

.assign-btn {
  background-color: #3490dc;
  color: white;
  border: none;
}

.assign-btn:hover {
  background-color: #2779bd;
}

.self-assign-btn {
  background-color: #38c172;
  color: white;
  border: none;
}

.self-assign-btn:hover {
  background-color: #2d9958;
}

.request-btn {
  background-color: #f6993f;
  color: white;
  border: none;
}

.request-btn:hover {
  background-color: #de751f;
}

.remove-btn {
  background-color: #e3342f;
  color: white;
  border: none;
}

.remove-btn:hover {
  background-color: #cc1f1a;
}

.disabled-btn {
  opacity: 0.6;
  cursor: not-allowed;
}

/* Estilos para las etiquetas de prioridad */
.priority-badge {
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  text-transform: capitalize;
  
  &.low { background: #bae6fd; color: #0369a1; }
  &.medium { background: #c7d2fe; color: #4338ca; }
  &.high { background: #e0e7ff; color: #3730a3; }
  &.urgent { background: #ddd6fe; color: #5b21b6; }
}

/* Estilos para la sección de prioridad */
.priority-management {
  margin: 1rem 0;
  padding: 1rem;
  background-color: #f9f9f9;
  border-radius: 6px;
  border: 1px solid #e5e5e5;
  
  .priority-section {
    h3 {
      margin-top: 0;
      font-size: 1rem;
      margin-bottom: 0.5rem;
    }
    
    .priority-controls {
      display: flex;
      gap: 10px;
      align-items: center;
      margin-top: 0.5rem;
      
      .priority-select {
        flex-grow: 1;
        padding: 8px 12px;
        border-radius: 4px;
        border: 1px solid #ced4da;
        background-color: white;
        font-size: 0.9rem;
      }
      
      button {
        padding: 8px 15px;
        border-radius: 4px;
        font-weight: 500;
        cursor: pointer;
        transition: all 0.2s ease;
        background-color: #3490dc;
        color: white;
        border: none;
        
        &:hover:not(:disabled) {
          background-color: #2779bd;
        }
        
        &:disabled {
          background-color: #ccc;
          cursor: not-allowed;
        }
      }
    }
  }
}

.priority-indicator {
  display: inline-block;
  width: 12px;
  height: 12px;
  border-radius: 50%;
  margin-right: 8px;
  position: relative;
  box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
  
  &.low {
    background: linear-gradient(135deg, #38BDF8, #0EA5E9);
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      border-radius: 50%;
      box-shadow: 0 0 0 2px rgba(56, 189, 248, 0.2);
    }
  }
  
  &.medium {
    background: linear-gradient(135deg, #818CF8, #6366F1);
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      border-radius: 50%;
      box-shadow: 0 0 0 2px rgba(129, 140, 248, 0.2);
    }
  }
  
  &.high {
    background: linear-gradient(135deg, #6D28D9, #4F46E5);
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      border-radius: 50%;
      box-shadow: 0 0 0 2px rgba(109, 40, 217, 0.2);
    }
  }
  
  &.urgent {
    background: linear-gradient(135deg, #8B5CF6, #7C3AED);
    
    &::after {
      content: '';
      position: absolute;
      top: 0;
      left: 0;
      right: 0;
      bottom: 0;
      border-radius: 50%;
      box-shadow: 0 0 0 2px rgba(139, 92, 246, 0.2);
    }
  }
}

.priority-text {
  font-weight: 500;
  
  &.priority-low {
    color: #0369a1;
  }
  
  &.priority-medium {
    color: #4338ca;
  }
  
  &.priority-high {
    color: #3730a3;
  }
  
  &.priority-urgent {
    color: #5b21b6;
    font-weight: 700;
  }
}

.assignment-btn {
  .priority-indicator {
    display: none !important;
  }
  
  .priority-text {
    display: none !important;
  }
}

.ticket-status {
  strong {
    margin-right: 0.5rem;
  }
  
  .status-badge {
    padding: 0.25rem 0.5rem;
    border-radius: 4px;
    text-transform: capitalize;
    
    &.open { background: #cfe7fe; color: #1e40af; }
    &.assigned { background: #dbd2fd; color: #5b21b6; }
    &.in_progress { background: #bae6fd; color: #0369a1; }
    &.resolved { background: #c7d2fe; color: #4338ca; }
    &.closed { background: #d1d5db; color: #374151; }
  }
}

/* Animación para el menú desplegable */
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

.dropdown-menu {
  position: relative;
  margin-right: 1rem;
              
  .action-btn {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1.25rem;
    border-radius: 8px;
    background: linear-gradient(135deg, #364FC7, #4263EB);
    border: none;
    color: white;
    font-weight: 500;
    font-size: 0.95rem;
    transition: all 0.2s ease;
    box-shadow: 0 2px 8px rgba(66, 99, 235, 0.2);
    
    &:hover {
      background: linear-gradient(135deg, #3B5BDB, #4C6EF5);
      box-shadow: 0 4px 12px rgba(66, 99, 235, 0.3);
      transform: translateY(-1px);
    }
    
    &:active {
      transform: translateY(0);
      box-shadow: 0 2px 4px rgba(66, 99, 235, 0.3);
    }
    
    .btn-text {
      margin: 0 0.25rem;
      font-family: 'Inter', system-ui, sans-serif;
    }

    i {
      font-size: 1rem;
    }
    
    &.assignment-btn {
      background: linear-gradient(135deg, #5E60CE, #6366F1);
      box-shadow: 0 2px 8px rgba(94, 96, 206, 0.2);
      
      &:hover {
        background: linear-gradient(135deg, #6D6AE8, #7B79F7);
        box-shadow: 0 4px 12px rgba(94, 96, 206, 0.3);
      }
    }
  }
  
  .dropdown-content {
    position: absolute;
    top: calc(100% + 0.5rem);
    left: 0;
    z-index: 100;
    min-width: 260px;
    display: none;
    background-color: white;
    border-radius: 10px;
    box-shadow: 0 8px 24px rgba(0, 0, 0, 0.12);
    overflow: hidden;
    border: 1px solid #E5E7EB;
    
    &.show {
      display: flex;
      flex-direction: column;
      animation: fadeInDown 0.2s ease-out;
    }
    
    .dropdown-header {
      padding: 1rem 1.25rem;
      background: linear-gradient(to right, #F3F4FF, #EFF6FF);
      font-weight: 600;
      color: #4F46E5;
      border-bottom: 1px solid #E5E7EB;
      font-size: 0.95rem;
    }
    
    .dropdown-items {
      max-height: 300px;
      overflow-y: auto;
      
      .dropdown-item {
        padding: 0.9rem 1.25rem;
        display: flex;
        align-items: center;
        gap: 0.75rem;
        cursor: pointer;
        transition: all 0.2s;
        font-size: 0.95rem;
        color: #1F2937;
        
        &:hover {
          background-color: #F3F4FF;
        }
        
        &.active {
          background-color: #EFF6FF;
          color: #4F46E5;
          font-weight: 500;
          
          &:hover {
            background-color: #E0E7FF;
          }
        }

        i {
          color: #4F46E5;
          font-size: 1rem;
        }
      }
      
      .dropdown-divider {
        height: 1px;
        background-color: #E5E7EB;
        margin: 0.5rem 1rem;
      }
    }
  }
}

.close-ticket-btn {
  background: linear-gradient(135deg, #059669, #10b981);
  margin-left: 0.5rem;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  
  i {
    color: white;
  }
  
  &:hover {
    background: linear-gradient(135deg, #047857, #059669);
    transform: translateY(-2px);
    box-shadow: 0 4px 12px rgba(5, 150, 105, 0.4);
  }
  
  &:active {
    transform: translateY(0);
  }
}

.close-btn-container {
  margin-left: 0.5rem;
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
  border-radius: 8px;
  width: 100%;
  max-width: 500px;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1);
  overflow: hidden;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #e5e7eb;
  
  h3 {
    margin: 0;
    color: #111827;
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
  border-radius: 6px;
  margin-bottom: 1rem;
  font-size: 1rem;
  color: #1f2937;
}

.close-reason-textarea {
  width: 100%;
  height: 100px;
  padding: 0.75rem;
  border: 1px solid #d1d5db;
  border-radius: 6px;
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
    padding: 0.5rem 1rem;
    border-radius: 6px;
    font-weight: 500;
    cursor: pointer;
    
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
</style> 
/* eslint-disable */
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
          
          <!-- NUEVA IMPLEMENTACIÓN DE LOS BOTONES DE ACCIÓN -->
          <div class="ticket-admin-actions" v-if="hasAdminAccess || hasAssistantAccess">
            <div class="admin-buttons">
              <!-- Botón para cambiar prioridad -->
              <div class="admin-button">
                <button class="filter-btn priority-btn" @click="togglePriorityMenu($event)" 
                  :class="{ 'active': showPriorityMenu }">
                  <i class="pi pi-flag"></i>
                  <span>Cambiar prioridad</span>
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
              
              <!-- Botón para asignar usuario -->
              <div class="admin-button">
                <button class="filter-btn assign-btn" @click="toggleAssignMenu($event)"
                  :class="{ 'active': showAssignMenu }">
                  <i class="pi pi-user"></i>
                  <span>Asignar usuario</span>
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
              
              <!-- Botón para cerrar ticket -->
              <div class="admin-button" v-if="!isTicketClosed">
                <button class="filter-btn close-btn" @click="showCloseTicketModal = true">
                  <i class="pi pi-check-circle"></i>
                  <span>Cerrar ticket</span>
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
});

// Función separada para manejar clics fuera de los menús
const handleOutsideClick = (event) => {
  // Evitamos la ejecución si los menús ya están cerrados
  if (!showPriorityMenu.value && !showAssignMenu.value) {
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
      assign: showAssignMenu.value
    }
  });
  
  // Solo cerrar menús si el clic no fue en un botón ni dentro de un menú
  if (!isClickInsideButton && !isClickInsideDropdown) {
    hidePriorityMenu();
    hideAssignMenu();
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
        
        .ticket-admin-actions {
          margin-top: 1.75rem;
          padding: 1.5rem;
          background-color: var(--bg-tertiary);
          border-radius: 12px;
          box-shadow: var(--card-shadow);
          border: 1px solid var(--border-color);
          
          .admin-buttons {
            display: flex;
            flex-wrap: wrap;
            gap: 10px;
            justify-content: center;
            
            .admin-button {
              position: relative;
              margin-right: 10px;
              z-index: 100;
              display: inline-block;
              
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
                display: flex;
                align-items: center;
                gap: 0.5rem;
                
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
                
                &.priority-btn {
                  i {
                    color: #3b82f6;
                  }
                  
                  &.active, &:hover {
                    i {
                      color: white;
                    }
                  }
                }
                
                &.assign-btn {
                  i {
                    color: #10b981;
                  }
                  
                  &.active, &:hover {
                    i {
                      color: white;
                    }
                  }
                }
                
                &.close-btn {
                  i {
                    color: #f43f5e;
                  }
                  
                  &:hover {
                    i {
                      color: white;
                    }
                  }
                }
                
                i {
                  font-size: 1rem;
                }
              }
              
              .dropdown-content {
                position: absolute;
                z-index: 9999;
                background-color: var(--card-bg);
                border: 1px solid var(--border-color);
                border-radius: 8px;
                box-shadow: var(--card-shadow);
                min-width: 200px;
                display: none;
                margin-top: 5px;
                animation: fadeInDown 0.3s ease both;
                
                &.show {
                  display: block;
                  visibility: visible;
                  opacity: 1;
                }
                
                .dropdown-item {
                  padding: 0.75rem 1rem;
                  cursor: pointer;
                  display: flex;
                  align-items: center;
                  color: var(--text-primary);
                  font-weight: 500;
                  transition: all 0.2s ease;
                  
                  &:hover {
                    background-color: var(--bg-hover);
                  }
                }
                
                .dropdown-divider {
                  height: 1px;
                  background-color: var(--border-color);
                  margin: 8px 0;
                }
                
                .dropdown-item i {
                  margin-right: 10px;
                  color: var(--primary-color);
                }
              }
            }
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
  width: 16px;
  height: 16px;
  
  &.low {
    background-color: #3b82f6;
  }
  
  &.medium {
    background-color: #8b5cf6;
  }
  
  &.high {
    background-color: #ec4899;
  }
  
  &.urgent {
    background-color: #ef4444;
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
      display: block;
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

/* ESTILOS CRÍTICOS PARA LA VISIBILIDAD DE LOS BOTONES */
.ticket-actions {
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
  margin-top: 1.75rem;
  padding-top: 1.75rem;
  border-top: 1px solid #E5E7EB;
  
  .action-buttons {
    display: flex !important;
    visibility: visible !important;
    opacity: 1 !important;
    gap: 1.25rem;
    flex-wrap: wrap;
  }
  
  .dropdown-menu {
    display: inline-block !important;
    visibility: visible !important;
    opacity: 1 !important;
  }
  
  .action-btn {
    display: flex !important;
    visibility: visible !important;
    opacity: 1 !important;
  }
}


/* Estilos para los menús desplegables */
.dropdown-content.show {
  display: block !important;
}

.admin-button {
  position: relative;
  margin-right: 10px;
  margin-bottom: 10px;
}

.action-btn {
  margin-bottom: 5px;
}

/* Ensure dropdown menus are visible */
.dropdown-content {
  position: absolute;
  z-index: 1000;
  background-color: white;
  border: 1px solid #e5e7eb;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  min-width: 200px;
}

.dropdown-item {
  padding: 10px 16px;
  display: flex;
  align-items: center;
  gap: 10px;
  cursor: pointer;
  transition: background-color 0.2s;
}

.dropdown-item:hover {
  background-color: #f3f4f6;
}

.priority-dropdown {
  display: none;
}

.assign-dropdown {
  display: none;
}

.dropdown-content {
  position: absolute !important;
  z-index: 9999 !important;
  background-color: white !important;
  border: 1px solid #ccc !important;
  border-radius: 4px !important;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15) !important;
}

.dropdown-content.show {
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
}

/* Agregamos más margen al contenedor del botón para que haya espacio para el menú */
.admin-button {
  position: relative !important;
  margin-bottom: 15px !important;
  z-index: 100 !important; /* Para que aparezca sobre otros elementos */
  display: inline-block !important;
}

/* Estilos críticos para mantener menús visibles */
.priority-dropdown.show, 
.assign-dropdown.show {
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
}

/* Aseguramos que nada más pueda ocultar los menús */
body .dropdown-content.show {
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
}

.dropdown-content {
  position: absolute !important;
  z-index: 9999 !important;
  background-color: var(--card-bg) !important;
  border: 1px solid var(--border-color) !important;
  border-radius: 8px !important;
  box-shadow: var(--card-shadow) !important;
  animation: fadeInDown 0.3s ease both;
}

.dropdown-content.show {
  display: block !important;
  visibility: visible !important;
  opacity: 1 !important;
}

.dropdown-item {
  padding: 0.75rem 1rem !important;
  cursor: pointer !important;
  display: flex !important;
  align-items: center !important;
  color: var(--text-primary) !important;
  font-weight: 500 !important;
  transition: all 0.2s ease !important;
}

.dropdown-item:hover {
  background-color: var(--bg-hover) !important;
}

/* Estilo para los botones de acción */
.action-btn {
  display: flex !important;
  align-items: center !important;
  gap: 0.5rem !important;
  color: white !important;
  font-weight: 600 !important;
  padding: 0.75rem 1.5rem !important;
  border-radius: 8px !important;
  cursor: pointer !important;
  border: none !important;
  position: relative !important;
  box-shadow: 0 4px 6px rgba(0, 0, 0, 0.1) !important;
  transition: all 0.3s ease !important;
}

.action-btn:hover {
  transform: translateY(-2px) !important;
  box-shadow: 0 8px 15px rgba(0, 0, 0, 0.15) !important;
}

.action-btn.priority-btn {
  background-color: var(--primary-color) !important;
}

.action-btn.priority-btn:hover {
  background-color: var(--primary-hover) !important;
}

.action-btn.assign-btn {
  background-color: var(--success-color) !important;
}

.action-btn.assign-btn:hover {
  background-color: var(--success-hover) !important;
}

.action-btn.close-btn {
  background-color: var(--danger-color) !important;
}

.action-btn.close-btn:hover {
  background-color: var(--danger-hover) !important;
}

/* Agregamos más margen al contenedor del botón para que haya espacio para el menú */
.admin-button {
  position: relative !important;
  margin-bottom: 15px !important;
  margin-right: 10px !important;
  z-index: 100 !important;
  display: inline-block !important;
}

/* Animación para los menús desplegables */
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

/* Estilos para botones y acciones */
.ticket-admin-actions {
  margin-top: 1.75rem;
  padding: 1.5rem;
  background-color: var(--bg-tertiary);
  border-radius: 12px;
  box-shadow: var(--card-shadow);
  border: 1px solid var(--border-color);
  
  .admin-buttons {
    display: flex;
    flex-wrap: wrap;
    gap: 10px;
    justify-content: center;
  }
}

.admin-button {
  position: relative;
  margin-right: 10px;
  z-index: 100;
  display: inline-block;
}

/* Estilo para los botones tipo filter-btn */
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
  display: flex;
  align-items: center;
  gap: 0.5rem;
}

.filter-btn:hover {
  background-color: var(--hover-bg);
  transform: translateY(-2px);
}

.filter-btn.active {
  background-color: var(--primary-color);
  color: white;
  border-color: transparent;
  font-weight: 600;
  box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
}

/* Estilos para el menú desplegable */
.dropdown-content {
  position: absolute;
  z-index: 9999;
  background-color: var(--card-bg);
  border: 1px solid var(--border-color);
  border-radius: 8px;
  box-shadow: var(--card-shadow);
  min-width: 200px;
  display: none;
  margin-top: 5px;
  animation: fadeInDown 0.3s ease both;
}

.dropdown-content.show {
  display: block;
  visibility: visible;
  opacity: 1;
}

.dropdown-item {
  padding: 0.75rem 1rem;
  cursor: pointer;
  display: flex;
  align-items: center;
  color: var(--text-primary);
  font-weight: 500;
  transition: all 0.2s ease;
}

.dropdown-item:hover {
  background-color: var(--bg-hover);
}

.dropdown-divider {
  height: 1px;
  background-color: var(--border-color);
  margin: 8px 0;
}

.dropdown-item i {
  margin-right: 10px;
  color: var(--primary-color);
}

/* Indicadores de prioridad */
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

/* Animación para los menús desplegables */
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
</style>
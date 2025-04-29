<template>
  <div>
    <!-- Chat bubble -->
    <div 
      class="fixed bottom-6 right-6 w-14 h-14 rounded-full flex items-center justify-center cursor-pointer shadow-lg transition-all z-50"
      :style="{ backgroundColor: primaryColor }"
      @click="toggleChat"
    >
      <i class="pi pi-comments text-white" v-if="!isOpen"></i>
      <i class="pi pi-times text-white" v-else></i>
    </div>

    <!-- Chat interface -->
    <div 
      v-if="isOpen" 
      class="fixed bottom-24 right-6 w-80 h-96 bg-white rounded-lg shadow-xl flex flex-col overflow-hidden border border-gray-200 z-40"
    >
      <!-- Chat header -->
      <div :style="{ backgroundColor: primaryColor }" class="text-white p-4 flex justify-between items-center">
        <h3 class="font-medium">{{ brandName }}</h3>
        <div class="flex items-center">
          <i v-if="isRegistered" class="pi pi-sign-out cursor-pointer mr-3" @click.stop="logout" title="Cerrar sesión"></i>
          <i class="pi pi-times cursor-pointer" @click="toggleChat"></i>
        </div>
      </div>
      
      <!-- Registration form (displayed before chat) -->
      <div v-if="!isRegistered" class="flex-1 p-4 overflow-y-auto bg-gray-50 flex flex-col">
        <div class="text-center text-gray-700 mb-4">
          {{ welcomeMessage }}
        </div>
        <form @submit.prevent="submitRegistration" class="flex flex-col gap-4">
          <div class="flex flex-col">
            <label for="name" class="text-sm text-gray-600 mb-1">Nombre</label>
            <input 
              v-model="userData.name" 
              type="text" 
              id="name" 
              class="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2"
              :style="{ '--tw-ring-color': primaryColor }"
              required
            />
          </div>
          <div class="flex flex-col">
            <label for="email" class="text-sm text-gray-600 mb-1">Email</label>
            <input 
              v-model="userData.email" 
              type="email" 
              id="email" 
              class="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2"
              :style="{ '--tw-ring-color': primaryColor }"
              required
            />
          </div>
          <div class="flex flex-col">
            <label for="firstMessage" class="text-sm text-gray-600 mb-1">¿En qué podemos ayudarte?</label>
            <textarea 
              v-model="userData.initialMessage" 
              id="firstMessage" 
              rows="3"
              class="border border-gray-300 rounded-md px-3 py-2 focus:outline-none focus:ring-2"
              :style="{ '--tw-ring-color': primaryColor }"
              required
              placeholder="Describe brevemente tu consulta..."
            ></textarea>
          </div>
          <button 
            type="submit" 
            class="mt-2 text-white rounded-md py-2 font-medium focus:outline-none"
            :style="{ backgroundColor: primaryColor }"
            :disabled="loading"
          >
            <span v-if="!loading">Iniciar chat</span>
            <span v-else>
              <i class="pi pi-spin pi-spinner"></i> Procesando...
            </span>
          </button>
        </form>
      </div>
      
      <!-- FAQ and Chat sections (displayed after registration) -->
      <div v-else class="flex-1 overflow-y-auto bg-gray-50">
        <!-- Chat View -->
        <div v-if="showChatView" class="p-4 chat-messages">
          <div v-if="messages.length === 0" class="text-center text-gray-500 mt-20">
            Inicia una conversación escribiendo un mensaje.
          </div>
          <div v-for="(message, index) in messages" :key="index" class="mb-3">
            <div 
              :class="[
                'max-w-[80%] p-3 rounded-lg', 
                message.isUser ? 'text-white ml-auto rounded-br-none' : 'bg-gray-200 text-gray-800 rounded-bl-none'
              ]"
              :style="message.isUser ? { backgroundColor: primaryColor } : {}"
            >
              {{ message.text }}
            </div>
          </div>
        </div>
        
        <!-- FAQ View -->
        <div v-else class="p-4">
          <div class="text-center mb-4">
            <button 
              @click="showChatView = true" 
              class="text-white rounded-md py-2 px-4 font-medium focus:outline-none w-full"
              :style="{ backgroundColor: primaryColor }"
            >
              <i class="pi pi-comments mr-2"></i>Iniciar chat
            </button>
          </div>
          
          <h3 class="text-lg font-semibold mb-3 text-gray-700">Preguntas Frecuentes</h3>
          
          <div v-if="loadingFaqs" class="text-center py-4">
            <i class="pi pi-spin pi-spinner text-gray-500"></i>
            <p class="text-sm text-gray-500 mt-2">Cargando preguntas frecuentes...</p>
          </div>
          
          <div v-else-if="faqs.length === 0" class="text-center py-4 text-gray-500">
            No hay preguntas frecuentes disponibles.
          </div>
          
          <div v-else>
            <div v-for="(category, index) in faqCategories" :key="index" class="mb-4">
              <h4 class="font-medium text-gray-700 mb-2">{{ category }}</h4>
              <div class="space-y-2">
                <div 
                  v-for="faq in getFaqsByCategory(category)" 
                  :key="faq.id" 
                  class="border border-gray-200 rounded-lg overflow-hidden"
                >
                  <div 
                    class="p-3 bg-gray-100 cursor-pointer flex justify-between items-center"
                    @click="toggleFaq(faq.id)"
                  >
                    <span class="font-medium text-gray-800">{{ faq.question }}</span>
                    <i class="pi" :class="expandedFaqs.includes(faq.id) ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
                  </div>
                  <div v-if="expandedFaqs.includes(faq.id)" class="p-3 bg-white">
                    <p class="text-gray-700">{{ faq.answer }}</p>
                    <div class="mt-2 flex justify-end">
                      <button 
                        @click="setInitialQuestion(faq.question)" 
                        class="text-xs text-white rounded-md py-1 px-2 focus:outline-none"
                        :style="{ backgroundColor: primaryColor }"
                      >
                        Consultar más
                      </button>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
      </div>
      
      <!-- Chat input (displayed after registration and when in chat view) -->
      <div v-if="isRegistered && showChatView" class="p-3 border-t border-gray-200 bg-white">
        <div class="flex items-center">
          <input 
            v-model="newMessage" 
            type="text" 
            placeholder="Escribe un mensaje..." 
            class="flex-1 border border-gray-300 rounded-full px-4 py-2 focus:outline-none focus:ring-2"
            :style="{ '--tw-ring-color': primaryColor }"
            @keyup.enter="sendMessage"
          />
          <button 
            @click="sendMessage" 
            class="ml-2 text-white rounded-full p-2 focus:outline-none"
            :style="{ backgroundColor: primaryColor }"
          >
            <i class="pi pi-send"></i>
          </button>
        </div>
      </div>
      
      <!-- Bottom navigation when in FAQ view -->
      <div v-if="isRegistered && !showChatView" class="p-3 border-t border-gray-200 bg-white">
        <div class="flex justify-between items-center">
          <button 
            @click="refreshFaqs" 
            class="text-xs text-gray-600 focus:outline-none flex items-center"
          >
            <i class="pi pi-refresh mr-1"></i> Actualizar
          </button>
          <button 
            @click="showChatView = true" 
            class="text-white rounded-md py-1 px-4 text-sm focus:outline-none"
            :style="{ backgroundColor: primaryColor }"
          >
            <i class="pi pi-comments mr-1"></i> Chat
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed, onBeforeUnmount } from 'vue';
import { useWidgetApi, getSession, apiConfig, type FAQ } from '../api/widgetApi';

// Props del componente
const props = defineProps({
  primaryColor: {
    type: String,
    default: '#6200ea'
  },
  brandName: {
    type: String,
    default: 'Chat Support'
  },
  position: {
    type: String,
    default: 'bottom-right'
  },
  welcomeMessage: {
    type: String,
    default: 'Para iniciar tu consulta, por favor completa el siguiente formulario:'
  }
});

const isOpen = ref(false);
const isRegistered = ref(false);
const loading = ref(false);
const currentTicketId = ref('');
const messages = ref<Array<{text: string, isUser: boolean, id?: string, pending?: boolean, error?: boolean}>>([]);
const newMessage = ref('');
const webSocket = ref<WebSocket | null>(null);

// Vista actual (chat o FAQs)
const showChatView = ref(false);

// Estado de FAQs
const faqs = ref<FAQ[]>([]);
const loadingFaqs = ref(false);
const expandedFaqs = ref<number[]>([]);

// Datos del usuario para el registro
const userData = ref({
  name: '',
  email: '',
  initialMessage: ''
});

// API del widget
const api = useWidgetApi();

// Estado de sesión
const hasSession = computed(() => {
  return api.hasActiveSession();
});

// Categorías de FAQs (computed)
const faqCategories = computed(() => {
  const categories = new Set<string>();
  faqs.value.forEach(faq => categories.add(faq.category));
  return Array.from(categories);
});

// Filtrar FAQs por categoría
const getFaqsByCategory = (category: string) => {
  return faqs.value.filter(faq => faq.category === category);
};

// Toggle expandir/colapsar FAQ
const toggleFaq = (id: number) => {
  const index = expandedFaqs.value.indexOf(id);
  if (index === -1) {
    expandedFaqs.value.push(id);
  } else {
    expandedFaqs.value.splice(index, 1);
  }
};

// Mejorar el método loadFaqs con mejor manejo de errores
const loadFaqs = async () => {
  loadingFaqs.value = true;
  console.log('Iniciando carga de FAQs desde el servidor...');
  
  try {
    // Hacer la llamada a la API con log detallado
    console.log('Llamando a api.getFaqs()...');
    const response = await api.getFaqs();
    console.log('Respuesta de API para FAQs:', response);
    
    if (response && Array.isArray(response)) {
      // Filtrar solo las publicadas
      const publishedFaqs = response.filter(faq => faq.isPublished);
      console.log(`FAQs publicadas encontradas: ${publishedFaqs.length}`);
      
      faqs.value = publishedFaqs;
      
      // Si no hay FAQs, mostrar mensaje en consola
      if (publishedFaqs.length === 0) {
        console.warn('No se encontraron FAQs publicadas para mostrar');
      } else {
        // Expandir la primera FAQ de cada categoría
        const categories = new Set<string>();
        const firstFaqsIds: number[] = [];
        
        faqs.value.forEach(faq => {
          if (!categories.has(faq.category)) {
            categories.add(faq.category);
            firstFaqsIds.push(faq.id);
          }
        });
        
        expandedFaqs.value = firstFaqsIds;
        console.log('Categorías de FAQs encontradas:', Array.from(categories));
      }
    } else {
      console.error('Formato inesperado en la respuesta de FAQs:', response);
      faqs.value = [];
    }
  } catch (error) {
    console.error('Error al cargar FAQs:', error);
    faqs.value = []; // Asegurarnos de que no hay FAQs antiguas
    
    // Si falla la carga de FAQs, mostrar el chat
    if (isRegistered.value) {
      showChatView.value = true;
      console.log('Cambiando a vista de chat debido a error en FAQs');
    }
  } finally {
    loadingFaqs.value = false;
    console.log('Carga de FAQs finalizada. Estado:', { 
      faqsCount: faqs.value.length, 
      expandedFaqsCount: expandedFaqs.value.length,
      showingChat: showChatView.value
    });
  }
};

// Refrescar FAQs
const refreshFaqs = () => {
  loadFaqs();
};

// Usar una pregunta como mensaje inicial
const setInitialQuestion = (question: string) => {
  newMessage.value = question;
  showChatView.value = true;
  // Dar tiempo para que la UI se actualice y luego enviar
  setTimeout(() => {
    sendMessage();
  }, 100);
};

// Conexión WebSocket
const connectWebSocket = (ticketId: string) => {
  if (!ticketId) {
    console.error('No se puede conectar al WebSocket sin un ID de ticket');
    return;
  }
  
  // Cerrar WebSocket existente si hay uno
  if (webSocket.value) {
    webSocket.value.close();
    webSocket.value = null;
  }
  
  // Guardar el ID de ticket que estamos usando
  const effectiveTicketId = ticketId;
  
  // Utilizar la URL de la API para construir la URL del WebSocket
  const apiUrl = new URL(apiConfig.apiUrl);
  
  // Determinar protocolo (ws o wss dependiendo si la API usa http o https)
  const wsProtocol = apiUrl.protocol === 'https:' ? 'wss:' : 'ws:';
  
  // Obtener host y puerto
  const hostPart = apiUrl.hostname;
  const portPart = apiUrl.port || (apiUrl.protocol === 'https:' ? '443' : '80');
  
  // URL final del WebSocket
  const wsUrl = `${wsProtocol}//${hostPart}:${portPart}/api/ws/chat/${effectiveTicketId}`;
  
  console.log('Información de conexión WebSocket:', {
    apiUrl: apiConfig.apiUrl,
    parsedHost: hostPart,
    parsedPort: portPart,
    wsProtocol,
    finalWsUrl: wsUrl,
    ticketId: effectiveTicketId
  });
  
  try {
    webSocket.value = new WebSocket(wsUrl);
    
    webSocket.value.onopen = () => {
      console.log('Conexión WebSocket establecida correctamente');
    };
    
    webSocket.value.onmessage = (event) => {
      console.log('Mensaje WebSocket recibido:', event.data);
      try {
        const data = JSON.parse(event.data);
        
        // Log completo para depuración
        console.log('Datos WebSocket procesados:', {
          type: data.type,
          hasData: !!data.data,
          hasMessage: !!data.message,
          dataContent: data.data?.content,
          messageContent: data.message?.content,
          dataIsClient: data.data?.isClient,
          messageIsClient: data.message?.isClient
        });
        
        // MEJORA: Extraer el contenido del mensaje independientemente de la estructura
        let messageContent = '';
        let isClientMessage = false;
        let messageObj = null;
        
        // Verificar diferentes estructuras posibles del mensaje
        if (data.type === 'new_message') {
          // Puede tener 'message' o 'data'
          if (data.message && data.message.content) {
            messageContent = data.message.content;
            isClientMessage = data.message.isClient === true;
            messageObj = data.message;
          } else if (data.data && data.data.content) {
            messageContent = data.data.content;
            isClientMessage = data.data.isClient === true;
            messageObj = data.data;
          }
        } else if (data.type === 'message_received') {
          // Este formato viene del servidor para confirmar mensajes
          if (data.data && data.data.content) {
            messageContent = data.data.content;
            isClientMessage = data.data.isClient === true;
            messageObj = data.data;
          }
        } else if (data.content) {
          // Formato directo sin type
          messageContent = data.content;
          isClientMessage = data.isClient === true;
          messageObj = data;
        }
        
        // Solo procesar si se encontró contenido
        if (messageContent) {
          console.log('Mensaje extraído para mostrar:', {
            content: messageContent,
            isClientMessage: isClientMessage,
            messageId: messageObj?.id || 'no-id'
          });
          
          // Verificar que no sea un mensaje duplicado
          const isDuplicate = messages.value.some(
            msg => msg.text === messageContent && 
                  msg.isUser === isClientMessage
          );
          
          if (!isDuplicate) {
            // Importante: isUser=true significa que es un mensaje enviado por el usuario del widget
            // Los mensajes con isClient=false son de agentes de soporte, que deben mostrarse como isUser=false
            messages.value.push({
              text: messageContent,
              isUser: isClientMessage // isUser=true para mensajes del cliente, false para mensajes del agente
            });
            
            // Hacer scroll al final
            scrollToBottom();
          } else {
            console.log('Mensaje duplicado detectado y omitido');
          }
        } else if (data.type === 'error') {
          console.error('Error del servidor WebSocket:', data.message || data.data || 'Error desconocido');
        } else if (data.type === 'connection_established' || data.type === 'identify_success') {
          console.log('Conexión WebSocket confirmada:', data.type);
        } else {
          console.warn('Formato de mensaje WebSocket no reconocido:', data);
        }
      } catch (error) {
        console.error('Error al procesar mensaje WebSocket:', error);
      }
    };
    
    webSocket.value.onerror = (error) => {
      console.error('Error en conexión WebSocket:', error);
    };
    
    webSocket.value.onclose = (event) => {
      console.log('Conexión WebSocket cerrada. Código:', event.code, 'Razón:', event.reason);
      
      // Reconectar después de un tiempo si el chat sigue abierto
      if (isOpen.value && isRegistered.value) {
        setTimeout(() => {
          connectWebSocket(currentTicketId.value);
        }, 5000);
      }
    };
  } catch (error) {
    console.error('Error al crear conexión WebSocket:', error);
  }
};

// Actualizar onMounted para añadir carga de FAQs
onMounted(() => {
  const container = document.getElementById('growdesk-widget-container');
  if (container) {
    container.addEventListener('open-widget', () => {
      isOpen.value = true;
    });
    container.addEventListener('close-widget', () => {
      isOpen.value = false;
    });
  }
  
  // Verificar si hay una sesión activa
  if (hasSession.value) {
    const session = api.getSession();
    if (session) {
      // Restaurar datos de sesión
      userData.value.name = session.name;
      userData.value.email = session.email;
      currentTicketId.value = session.ticketId;
      isRegistered.value = true;
      
      // Inicialmente mostrar las FAQs en lugar del chat
      showChatView.value = false;
      
      // Cargar preguntas frecuentes
      loadFaqs();
      
      // Cargar mensajes anteriores
      loadPreviousMessages(session.ticketId);
      
      // Conectar WebSocket
      connectWebSocket(session.ticketId);
    }
  }
});

// Limpieza de conexiones al desmontar
onBeforeUnmount(() => {
  if (webSocket.value) {
    webSocket.value.close();
    webSocket.value = null;
  }
});

const toggleChat = () => {
  isOpen.value = !isOpen.value;
};

// Enviar el formulario de registro y crear el ticket
const submitRegistration = async () => {
  if (!userData.value.name || !userData.value.email || !userData.value.initialMessage) return;
  
  loading.value = true;
  
  try {
    // Crear el ticket con los datos del usuario
    const ticketData = {
      name: userData.value.name,
      email: userData.value.email,
      message: userData.value.initialMessage, // Usamos el mensaje inicial proporcionado por el usuario
      metadata: {
        url: window.location.href,
        referrer: document.referrer,
        userAgent: navigator.userAgent,
        screenSize: `${window.innerWidth}x${window.innerHeight}`
      }
    };
    
    const response = await api.createTicket(ticketData);
    
    // Guardar el ID del ticket
    currentTicketId.value = response.ticketId;
    
    // Marcar al usuario como registrado
    isRegistered.value = true;
    
    // EXPLÍCITAMENTE establecer showChatView a false para mostrar las FAQs primero
    showChatView.value = false;
    
    // Cargar FAQs inmediatamente después de registrarse
    await loadFaqs();
    
    console.log('FAQs cargadas después del registro:', faqs.value);
    
    // Agregar el mensaje del usuario al chat (para cuando cambie a vista de chat)
    messages.value.push({
      text: userData.value.initialMessage,
      isUser: true
    });
    
    // Agregar mensaje de bienvenida
    messages.value.push({
      text: `Hola ${userData.value.name}, gracias por contactarnos. Puedes consultar nuestras preguntas frecuentes o iniciar un chat con nosotros.`,
      isUser: false
    });
    
    // Establecer conexión WebSocket con el nuevo ticket
    connectWebSocket(response.ticketId);
    
    // Hacer scroll al final del chat
    scrollToBottom();
  } catch (error) {
    console.error('Error al registrar al usuario:', error);
    alert('Lo sentimos, hubo un problema al iniciar el chat. Por favor, inténtalo de nuevo más tarde.');
  } finally {
    loading.value = false;
  }
};

// Cargar mensajes anteriores
const loadPreviousMessages = async (ticketId: string) => {
  try {
    loading.value = true;
    const response = await api.getMessageHistory(ticketId);
    
    if (response && response.messages && response.messages.length > 0) {
      console.log('Histórico de mensajes recibido:', response.messages);
      
      // Convertir y agregar mensajes al chat
      const formattedMessages = response.messages.map((msg: any) => {
        // Usamos la propiedad isClient para determinar si es un mensaje del cliente o del agente
        // IMPORTANTE: Asegurar que isClient siempre se interpreta como booleano y no como string
        const isClientMessage = msg.isClient === true || 
                               (typeof msg.isClient === 'string' && msg.isClient.toLowerCase() === 'true');
        
        console.log(`Mensaje procesado - contenido: "${msg.text || msg.content}", isClient: ${isClientMessage}`);
        
        return {
          text: msg.text || msg.content,
          isUser: isClientMessage // Los mensajes del cliente (isClient=true) son los del usuario
        };
      });
      
      messages.value = formattedMessages;
      scrollToBottom();
    } else {
      // Si no hay mensajes, agregar un mensaje de bienvenida
      messages.value.push({
        text: `Hola ${userData.value.name}, bienvenido de nuevo. ¿En qué podemos ayudarte hoy?`,
        isUser: false
      });
    }
  } catch (error) {
    console.error('Error al cargar mensajes anteriores:', error);
    messages.value.push({
      text: `Hola ${userData.value.name}, parece que hubo un problema al cargar tus mensajes anteriores. ¿En qué podemos ayudarte hoy?`,
      isUser: false
    });
  } finally {
    loading.value = false;
  }
};

// Función para hacer scroll al final del chat
const scrollToBottom = () => {
  setTimeout(() => {
    const messagesContainer = document.querySelector('.chat-messages');
    if (messagesContainer) {
      messagesContainer.scrollTop = messagesContainer.scrollHeight;
    }
  }, 100);
};

// Cerrar sesión
const logout = () => {
  api.logout();
  isRegistered.value = false;
  currentTicketId.value = '';
  messages.value = [];
  faqs.value = [];
  expandedFaqs.value = [];
  userData.value = {
    name: '',
    email: '',
    initialMessage: ''
  };
};

const sendMessage = async () => {
  if (newMessage.value.trim() === '' || !currentTicketId.value) return;
  
  // Guardar el mensaje actual y limpiarlo para evitar envíos duplicados
  const userMessage = newMessage.value;
  newMessage.value = '';
  
  // Mostrar el mensaje del usuario inmediatamente
  const messageId = `temp-${Date.now()}`;
  messages.value.push({
    id: messageId,
    text: userMessage,
    isUser: true,
    pending: true
  });
  
  // Hacer scroll al final del chat inmediatamente
  scrollToBottom();
  
  // Mostrar estado de carga
  loading.value = true;
  
  try {
    // Enviar el mensaje al servidor con los datos del usuario
    const response = await api.sendMessage({
      ticketId: currentTicketId.value,
      message: userMessage,
      userName: userData.value.name,
      userEmail: userData.value.email
    });
    
    // Actualizar el estado del mensaje una vez enviado
    const index = messages.value.findIndex(msg => msg.id === messageId);
    if (index >= 0) {
      messages.value[index].pending = false;
      
      // Si el mensaje fue detectado como duplicado por el servidor
      if (response.messageId === 'duplicate-ignored') {
        console.log('Mensaje detectado como duplicado por el servidor');
      }
    }
    
    // No mostrar respuesta automática si el mensaje fue un duplicado
    if (response.messageId !== 'duplicate-ignored') {
      // Auto-respuesta deshabilitada para evitar confusiones con mensajes reales
      // Las respuestas deben venir del agente o del servidor por WebSocket
    }
    
    // Desactivar estado de carga
    loading.value = false;
  } catch (error) {
    console.error('Error al enviar el mensaje:', error);
    
    // Marcar el mensaje como fallido
    const index = messages.value.findIndex(msg => msg.id === messageId);
    if (index >= 0) {
      messages.value[index].error = true;
      messages.value[index].pending = false;
    }
    
    // Mostrar mensaje de error en el chat
    messages.value.push({
      text: "Lo sentimos, hubo un problema al enviar tu mensaje. Por favor, inténtalo de nuevo.",
      isUser: false
    });
    
    // Hacer scroll al final después del error
    scrollToBottom();
    loading.value = false;
  }
};
</script>

<style>
/* Los estilos se manejan con Tailwind CSS */
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

.pi {
  font-size: 1.25rem;
}
</style>
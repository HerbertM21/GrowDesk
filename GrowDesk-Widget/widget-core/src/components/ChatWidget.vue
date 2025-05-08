<template>
  <div>
    <!-- Chat bubble con animación de pulso restaurada -->
    <div 
      class="fixed bottom-6 right-6 w-20 h-20 rounded-full flex items-center justify-center cursor-pointer shadow-xl transition-all z-50 hover:scale-110 animate-pulse-slow"
      :style="{ backgroundColor: primaryColor }"
      @click="toggleChat"
    >
      <i class="pi pi-comments text-white text-3xl" v-if="!isOpen"></i>
      <i class="pi pi-times text-white text-3xl" v-else></i>
    </div>

    <!-- Chat interface con sombras restauradas -->
    <div 
      v-if="isOpen" 
      class="fixed bottom-24 right-6 w-[420px] h-[600px] bg-white rounded-2xl shadow-2xl flex flex-col overflow-hidden border border-gray-200 z-40 transition-all duration-300 ease-in-out animate-slide-up"
    >
      <!-- Chat header - Diseño plano (sin gradiente) -->
      <div 
        class="p-5 flex justify-between items-center shadow-sm"
        :style="{ backgroundColor: primaryColor }"
      >
        <div class="flex items-center">
          <h3 class="font-semibold text-white text-lg">{{ brandName }}</h3>
        </div>
        <div class="flex items-center space-x-3">
          <i v-if="isRegistered" class="pi pi-sign-out cursor-pointer text-white hover:scale-110 transition-transform" @click.stop="logout" title="Cerrar sesión"></i>
          <i class="pi pi-times cursor-pointer text-white hover:scale-110 transition-transform" @click="toggleChat"></i>
        </div>
      </div>
      
      <!-- Registration form sin efectos de sombra -->
      <div v-if="!isRegistered" class="flex-1 p-8 overflow-y-auto bg-white flex flex-col animate-fade-in">
        <div class="text-center text-gray-700 mb-8">
          <div class="text-2xl font-semibold mb-2" :style="{ color: primaryColor }">Bienvenido</div>
          <p>{{ welcomeMessage }}</p>
        </div>
        <form @submit.prevent="submitRegistration" class="flex flex-col gap-6">
          <div class="flex flex-col">
            <label for="name" class="text-sm font-medium text-gray-700 mb-2">Nombre</label>
            <input 
              v-model="userData.name" 
              type="text" 
              id="name" 
              class="border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:border-2 transition-all"
              :style="{ '--tw-border-opacity': 1, borderColor: primaryColor }"
              required
            />
          </div>
          <div class="flex flex-col">
            <label for="email" class="text-sm font-medium text-gray-700 mb-2">Email</label>
            <input 
              v-model="userData.email" 
              type="email" 
              id="email" 
              class="border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:border-2 transition-all"
              :style="{ '--tw-border-opacity': 1, borderColor: primaryColor }"
              required
            />
          </div>
          <div class="flex flex-col">
            <label for="firstMessage" class="text-sm font-medium text-gray-700 mb-2">¿En qué podemos ayudarte?</label>
            <textarea 
              v-model="userData.initialMessage" 
              id="firstMessage" 
              rows="5"
              class="border border-gray-300 rounded-lg px-4 py-3 focus:outline-none focus:border-2 transition-all"
              :style="{ '--tw-border-opacity': 1, borderColor: primaryColor }"
              required
              placeholder="Describe brevemente tu consulta..."
            ></textarea>
          </div>
          <button 
            type="submit" 
            class="mt-2 text-white rounded-lg py-4 font-medium focus:outline-none transition-all hover:shadow-lg flex items-center justify-center"
            :style="{ backgroundColor: primaryColor }"
            :disabled="loading"
          >
            <span v-if="!loading" class="flex items-center"><i class="pi pi-send mr-2"></i> Iniciar chat</span>
            <span v-else class="flex items-center">
              <i class="pi pi-spin pi-spinner mr-2"></i> Procesando...
            </span>
          </button>
        </form>
      </div>
      
      <!-- FAQ and Chat sections (displayed after registration) -->
      <div v-else class="flex-1 overflow-y-auto bg-white">
        <!-- Chat View con mensajes mejorados -->
        <div v-if="showChatView" class="p-6 chat-messages">
          <div v-if="messages.length === 0" class="flex flex-col items-center justify-center h-full">
            <div class="text-center text-gray-500">
              <i class="pi pi-comments text-5xl mb-4" :style="{ color: primaryColor }"></i>
              <p>Inicia una conversación escribiendo un mensaje.</p>
            </div>
          </div>
          <div v-for="(message, index) in messages" :key="index" class="mb-5 animate-fade-in">
            <div 
              :class="[
                'max-w-[90%] p-5 rounded-lg', 
                message.isUser ? 'text-white ml-auto' : 'bg-gray-100 text-gray-800'
              ]"
              :style="message.isUser ? { backgroundColor: primaryColor } : {}"
            >
              <div class="flex items-start">
                <!-- Avatar solo para mensajes que no son del usuario -->
                <div v-if="!message.isUser" class="w-9 h-9 rounded-full bg-white border border-gray-200 flex items-center justify-center mr-3 flex-shrink-0">
                  <i class="pi pi-user text-gray-500"></i>
                </div>
                <div class="flex-1">
                  <div v-if="!message.isUser" class="font-semibold text-sm mb-1 text-gray-600">Soporte</div>
                  <div class="text-base">{{ message.text }}</div>
                  <div class="text-right text-xs mt-2" :class="message.isUser ? 'text-white opacity-70' : 'text-gray-500'">
                    {{ message.timestamp ? formatMessageTime(message.timestamp) : '' }}
                  </div>
                </div>
              </div>
              <!-- Indicador de estado del mensaje -->
              <div v-if="message.pending" class="text-right mt-1">
                <i class="pi pi-spin pi-spinner text-white text-xs"></i>
              </div>
              <div v-if="message.error" class="text-right mt-1 text-xs text-red-500">
                Error al enviar
              </div>
            </div>
          </div>
        </div>
        
        <!-- FAQ View mejorado -->
        <div v-else class="p-6">
          <div class="text-center mb-6">
            <button 
              @click="showChatView = true" 
              class="text-white rounded-lg py-3 px-5 font-medium focus:outline-none w-full transition-all hover:opacity-90 flex items-center justify-center"
              :style="{ backgroundColor: primaryColor }"
            >
              <i class="pi pi-comments mr-2"></i>Iniciar chat
            </button>
          </div>
          
          <h3 class="text-xl font-bold mb-5 text-gray-700 flex items-center">
            <i class="pi pi-question-circle mr-2" :style="{ color: primaryColor }"></i>
            Preguntas Frecuentes
          </h3>
          
          <div v-if="loadingFaqs" class="text-center py-6">
            <i class="pi pi-spin pi-spinner text-3xl" :style="{ color: primaryColor }"></i>
            <p class="text-sm text-gray-500 mt-3">Cargando preguntas frecuentes...</p>
          </div>
          
          <div v-else-if="faqs.length === 0" class="text-center py-6 text-gray-500">
            <i class="pi pi-inbox text-3xl mb-3" :style="{ color: primaryColor }"></i>
            <p>No hay preguntas frecuentes disponibles.</p>
          </div>
          
          <div v-else class="space-y-6">
            <div v-for="(category, index) in faqCategories" :key="index" class="mb-6 animate-fade-in"
                 :style="{ animationDelay: `${index * 0.1}s` }">
              <h4 class="font-semibold text-gray-700 mb-4 pb-2 border-b border-gray-200">{{ category }}</h4>
              <div class="space-y-4">
                <div 
                  v-for="faq in getFaqsByCategory(category)" 
                  :key="faq.id" 
                  class="border border-gray-200 rounded-lg overflow-hidden hover:border-gray-300 transition-all"
                >
                  <div 
                    class="p-4 bg-white cursor-pointer flex justify-between items-center"
                    @click="toggleFaq(faq.id)"
                    :style="expandedFaqs.includes(faq.id) ? { borderLeft: `4px solid ${primaryColor}` } : {}"
                  >
                    <span class="font-medium text-gray-800 flex-1">{{ faq.question }}</span>
                    <i class="pi transition-transform" 
                       :class="expandedFaqs.includes(faq.id) ? 'pi-chevron-up' : 'pi-chevron-down'"
                       :style="{ color: primaryColor }"></i>
                  </div>
                  <div v-if="expandedFaqs.includes(faq.id)" 
                       class="p-4 bg-gray-50 animate-fade-in">
                    <p class="text-gray-700">{{ faq.answer }}</p>
                    <div class="mt-4 flex justify-end">
                      <button 
                        @click="setInitialQuestion(faq.question)" 
                        class="text-sm text-white rounded-lg py-2 px-4 focus:outline-none flex items-center transition-all hover:opacity-90"
                        :style="{ backgroundColor: primaryColor }"
                      >
                        <i class="pi pi-comments mr-2"></i>
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
      
      <!-- Chat input (displayed after registration and when in chat view) mejorado -->
      <div v-if="isRegistered && showChatView" class="p-5 border-t border-gray-200 bg-white">
        <div class="flex items-center">
          <input 
            v-model="newMessage" 
            type="text" 
            placeholder="Escribe un mensaje..." 
            class="flex-1 border border-gray-300 rounded-lg px-5 py-4 focus:outline-none focus:border-2 transition-all"
            :style="{ '--tw-border-opacity': 1, borderColor: primaryColor }"
            @keyup.enter="sendMessage"
          />
          <button 
            @click="sendMessage" 
            class="ml-3 text-white rounded-lg p-4 focus:outline-none transition-all hover:opacity-90"
            :style="{ backgroundColor: primaryColor, border: '2px solid #000' }"
          >
            <i class="pi pi-send"></i>
          </button>
        </div>
      </div>
      
      <!-- Bottom navigation when in FAQ view mejorado -->
      <div v-if="isRegistered && !showChatView" class="p-5 border-t border-gray-200 bg-white">
        <div class="flex justify-between items-center">
          <button 
            @click="refreshFaqs" 
            class="text-sm text-gray-600 focus:outline-none flex items-center hover:text-gray-800 transition-colors"
          >
            <i class="pi pi-refresh mr-1"></i> Actualizar FAQs
          </button>
          <button 
            @click="showChatView = true" 
            class="text-white rounded-lg py-2 px-4 text-sm focus:outline-none transition-all hover:opacity-90"
            :style="{ backgroundColor: primaryColor, border: '2px solid #000' }"
          >
            <i class="pi pi-comments mr-2"></i> Ir al chat
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

// Función para oscurecer un color (para gradientes)
const darkenColor = (color: string, amount: number): string => {
  // Convertir a formato rgb
  let usePound = false;
  
  if (color[0] === "#") {
    color = color.slice(1);
    usePound = true;
  }
  
  const num = parseInt(color, 16);
  
  let r = (num >> 16) - amount;
  if (r > 255) r = 255;
  else if (r < 0) r = 0;
  
  let g = ((num >> 8) & 0x00FF) - amount;
  if (g > 255) g = 255;
  else if (g < 0) g = 0;
  
  let b = (num & 0x0000FF) - amount;
  if (b > 255) b = 255;
  else if (b < 0) b = 0;
  
  return (usePound ? "#" : "") + (g | (r << 8) | (b << 16)).toString(16).padStart(6, '0');
};

// Formateador de tiempo para mensajes
const formatMessageTime = (timestamp: string): string => {
  try {
    const date = new Date(timestamp);
    return date.toLocaleTimeString([], { hour: '2-digit', minute: '2-digit' });
  } catch (error) {
    return '';
  }
};

const isOpen = ref(false);
const isRegistered = ref(false);
const loading = ref(false);
const isSubmitting = ref(false);
const error = ref('');
const currentTicketId = ref('');
const messages = ref<Array<{text: string, isUser: boolean, id?: string, pending?: boolean, error?: boolean, timestamp?: string}>>([]);
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
  faqs.value.forEach((faq: FAQ) => categories.add(faq.category));
  return Array.from(categories);
});

// Filtrar FAQs por categoría
const getFaqsByCategory = (category: string) => {
  return faqs.value.filter((faq: FAQ) => faq.category === category);
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
  } finally {
    loadingFaqs.value = false;
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
  console.log('Intentando registrar usuario con:', {
    name: userData.value.name, 
    email: userData.value.email,
    message: userData.value.initialMessage
  });

  try {
    isSubmitting.value = true;
    error.value = '';

    // Validación
    if (!userData.value.name || userData.value.name.trim() === '') {
      console.error('Error de validación: nombre vacío');
      error.value = 'Por favor ingresa tu nombre';
      isSubmitting.value = false;
      return;
    }

    if (!userData.value.email || !userData.value.email.includes('@')) {
      console.error('Error de validación: email inválido:', userData.value.email);
      error.value = 'Por favor ingresa un email válido';
      isSubmitting.value = false;
      return;
    }

    if (!userData.value.initialMessage || userData.value.initialMessage.trim() === '') {
      console.error('Error de validación: mensaje vacío');
      error.value = 'Por favor ingresa un mensaje';
      isSubmitting.value = false;
      return;
    }

    // Crear ticket en el servidor
    console.log('Iniciando creación de ticket con datos validados');
    const ticketResult = await api.createTicket({
      name: userData.value.name,
      email: userData.value.email,
      message: userData.value.initialMessage,
      subject: `Solicitud de soporte - ${userData.value.name}`,
      metadata: {
        referrer: document.referrer,
        url: window.location.href,
        userAgent: navigator.userAgent,
        screenSize: `${window.innerWidth}x${window.innerHeight}`
      }
    });

    console.log('Ticket creado exitosamente:', ticketResult);
    
    // Guardar el ID del ticket
    currentTicketId.value = ticketResult.ticketId;
    
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
    connectWebSocket(ticketResult.ticketId);
    
    // Hacer scroll al final del chat
    scrollToBottom();
  } catch (error) {
    console.error('Error al registrar al usuario:', error);
    alert('Lo sentimos, hubo un problema al iniciar el chat. Por favor, inténtalo de nuevo más tarde.');
  } finally {
    loading.value = false;
    isSubmitting.value = false;
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

// Mejorar el método logout para recargar FAQs
const logout = async () => {
  try {
    await api.logout();
    isRegistered.value = false;
    currentTicketId.value = '';
    messages.value = [];
    // Recargar FAQs después del logout
    await loadFaqs();
  } catch (error) {
    console.error('Error al cerrar sesión:', error);
  }
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
/* Los estilos base se manejan con Tailwind CSS */
@import 'tailwindcss/base';
@import 'tailwindcss/components';
@import 'tailwindcss/utilities';

/* Estilos para iconos */
.pi {
  font-size: 1.25rem;
}

/* Animaciones personalizadas */
@keyframes pulse-slow {
  0% {
    transform: scale(1);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  }
  50% {
    transform: scale(1.05);
    box-shadow: 0 8px 25px rgba(0, 0, 0, 0.25);
  }
  100% {
    transform: scale(1);
    box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
  }
}

@keyframes slide-up {
  0% {
    transform: translateY(50px);
    opacity: 0;
  }
  100% {
    transform: translateY(0);
    opacity: 1;
  }
}

@keyframes fade-in {
  0% {
    opacity: 0;
  }
  100% {
    opacity: 1;
  }
}

@keyframes bounce-slow {
  0%, 100% {
    transform: translateY(0);
  }
  50% {
    transform: translateY(-10px);
  }
}

.animate-pulse-slow {
  animation: pulse-slow 3s infinite ease-in-out;
}

.animate-slide-up {
  animation: slide-up 0.4s ease-out forwards;
}

.animate-fade-in {
  animation: fade-in 0.3s ease-out forwards;
}

.animate-bounce-slow {
  animation: bounce-slow 2s infinite ease-in-out;
}

/* Mejoras para la experiencia en dispositivos móviles */
@media (max-width: 640px) {
  .fixed.bottom-24.right-6 {
    bottom: 5rem;
    right: 1rem;
    width: calc(100vw - 2rem) !important;
    max-width: 100% !important;
    height: calc(85vh - 5rem) !important;
  }
}

/* Mejoras para scroll en mensajes y FAQs */
.chat-messages {
  scrollbar-width: thin;
  scrollbar-color: rgba(0, 0, 0, 0.2) transparent;
}

.chat-messages::-webkit-scrollbar {
  width: 6px;
}

.chat-messages::-webkit-scrollbar-track {
  background: transparent;
}

.chat-messages::-webkit-scrollbar-thumb {
  background-color: rgba(0, 0, 0, 0.2);
  border-radius: 20px;
}

/* Transiciones suaves para todos los elementos interactivos */
button, 
input, 
textarea {
  transition: all 0.2s ease-in-out;
}

/* Efecto hover para FAQs */
.faq-item:hover {
  transform: translateX(3px);
}
</style>
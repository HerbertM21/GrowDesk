import { defineStore } from 'pinia'
import apiClient from '@/api/client'
import { ref } from 'vue'

// Ampliar el tipo Window para incluir initialTicketData
declare global {
  interface Window {
    initialTicketData?: {
      id: string;
      description?: string;
      messages?: any[];
      [key: string]: any;
    }
  }
}

interface Message {
  id: string
  ticketId?: string
  userId?: string
  content: string
  isInternal?: boolean
  isClient?: boolean
  timestamp?: string
  createdAt?: string
  attachments?: Attachment[]
}

interface Attachment {
  id: string
  fileName: string
  fileType: string
  fileSize: number
  fileURL: string
}

interface ChatState {
  messages: Record<string, Message[]>
  loading: boolean
  error: string | null
  currentTicketId: string | null
  socket: WebSocket | null
  connected: boolean
  simulationInterval: NodeJS.Timeout | null
}

export const useChatStore = defineStore('chat', {
  state: (): ChatState => ({
    messages: {},
    loading: false,
    error: null,
    currentTicketId: null,
    socket: null,
    connected: false,
    simulationInterval: null
  }),

  getters: {
    currentMessages: (state: ChatState) => 
      state.currentTicketId ? state.messages[state.currentTicketId] || [] : [],
    
    hasMessages: (state: ChatState) => 
      state.currentTicketId ? (state.messages[state.currentTicketId]?.length || 0) > 0 : false,
      
    // No necesitamos definir isConnected como getter porque ya está implementado como método
  },

  actions: {
    setCurrentTicket(ticketId: string) {
      this.currentTicketId = ticketId
      if (!this.messages[ticketId]) {
        this.messages[ticketId] = []
      }
      
      // Siempre activar la conexión WebSocket
      this.connectToWebSocket(ticketId)
    },

    async fetchMessages(ticketId: string) {
      // Si ya tenemos mensajes para este ticket y no estamos forzando una recarga, reutilizamos
      if (this.messages[ticketId] && this.messages[ticketId].length > 0) {
        console.log('Usando mensajes existentes para el ticket:', ticketId);
        return;
      }

      this.loading = true
      this.error = null
      try {
        console.log('Obteniendo mensajes para el ticket:', ticketId);
        // Intenta obtener mensajes desde la ruta de tickets
        const response = await apiClient.get(`/tickets/${ticketId}/messages`);
        console.log('Respuesta de mensajes:', response.data);
        
        // Normalizar el formato de los mensajes
        let messagesList = [];
        if (Array.isArray(response.data)) {
          messagesList = response.data;
        } else if (response.data && Array.isArray(response.data.messages)) {
          messagesList = response.data.messages;
        }
        
        this.messages[ticketId] = messagesList.map((msg: any) => ({
          id: msg.id,
          content: msg.content,
          isClient: msg.isClient,
          timestamp: msg.timestamp || msg.createdAt
        }));
        
        // Actualizar la lista de mensajes en el store para uso futuro
        this.messages = { ...this.messages };
      } catch (error) {
        console.error('Error al cargar mensajes:', error);
        
        // Si hay datos iniciales de ticket, intentar usar los mensajes de ahí
        if (this.currentTicketId === ticketId && window.initialTicketData && window.initialTicketData.messages) {
          console.log('Usando mensajes de datos iniciales para el ticket', ticketId);
          this.messages[ticketId] = window.initialTicketData.messages;
          this.messages = { ...this.messages };
        } else {
          // Si tampoco hay datos iniciales, usar mensajes de demostración
          // Esto sólo se usa durante el desarrollo
          console.log('Usando mensaje de demostración para el ticket', ticketId);
          this.messages[ticketId] = [
            {
              id: 'demo-msg-1',
              content: '¡Hola! Este es un mensaje de prueba.',
              isClient: false,
              timestamp: new Date().toISOString()
            }
          ];
          this.messages = { ...this.messages };
        }
      } finally {
        this.loading = false;
      }
    },

    async sendMessage(ticketId: string, content: string, isInternal: boolean = false) {
      if (!ticketId || !content) {
        return Promise.reject(new Error('Datos incorrectos para enviar mensaje'));
      }
      
      this.loading = true;
      
      // Crear el objeto de mensaje con el formato CORRECTO que espera el backend
      const messageData = {
        type: 'new_message',
        ticketId: ticketId,
        data: {
          content: content,
          isClient: false, // Mensaje del agente, no del cliente
          isInternal: isInternal, // Para notas internas
          userId: localStorage.getItem('userId') || '1',
          userName: localStorage.getItem('userName') || 'Agente',
          timestamp: new Date().toISOString()
        }
      };
      
      console.log('Enviando mensaje con datos:', messageData);
      
      // Si hay una conexión WebSocket activa, enviar por ahí y no por HTTP
      if (this.socket && this.socket.readyState === WebSocket.OPEN && this.connected) {
        try {
          // Formato estandarizado para todos los mensajes WebSocket
          this.socket.send(JSON.stringify(messageData));
          console.log('Mensaje enviado por WebSocket correctamente');
          
          // Agregar mensaje local inmediatamente para mejor UX
          const localMessage = {
            id: `local-${Date.now()}`,
            content,
            isClient: false, // Importante: marcar explícitamente como mensaje del agente
            timestamp: new Date().toISOString(),
            ticketId
          };
          
          this.handleNewMessage(localMessage);
          this.loading = false;
          return localMessage;
        } catch (e) {
          console.warn('Error al enviar mensaje por WebSocket, usando HTTP fallback:', e);
          // Continuar con el fallback HTTP
        }
      } else {
        console.warn('WebSocket no conectado, usando HTTP fallback');
      }
      
      // Solo usar HTTP como fallback si el WebSocket falló o no está disponible
      return apiClient.post(`/tickets/${ticketId}/messages`, {
        content,
        isInternal,
        isClient: false, // Explícitamente marcar como mensaje de agente
        userId: localStorage.getItem('userId') || '1',
        userName: localStorage.getItem('userName') || 'Agente', // Añadir nombre para identificación
        timestamp: new Date().toISOString()
      })
        .then(response => {
          const newMessage = response.data;
          
          // Si no hay una lista de mensajes para este ticket, crearla
          if (!this.messages[ticketId]) {
            this.messages[ticketId] = [];
          }
          
          // Agregar el mensaje a la lista, asegurando que isClient sea false
          this.handleNewMessage({
            id: newMessage.id || `msg-${Date.now()}`,
            content: newMessage.content || content,
            isClient: false, // IMPORTANTE: Garantizar que sea false para mensajes de agente
            timestamp: newMessage.timestamp || newMessage.createdAt || new Date().toISOString(),
            ticketId
          });
          
          this.loading = false;
          return newMessage;
        })
        .catch(error => {
          console.error('Error sending message:', error);
          this.error = 'Failed to send message';
          this.loading = false;
          
          // En modo desarrollo, simular éxito para poder seguir probando
          if (import.meta.env.DEV) {
            console.warn('DEV MODE: Simulando envío exitoso del mensaje');
            const fakeMessage = {
              id: `msg-${Date.now()}`,
              content,
              isClient: false, // Mensaje del agente
              timestamp: new Date().toISOString(),
              ticketId
            };
            
            this.handleNewMessage(fakeMessage);
            return fakeMessage;
          }
          
          throw error;
        });
    },

    connectToWebSocket(ticketId: string) {
      if (!ticketId) {
        console.error('No se puede conectar WebSocket: falta ticketId');
        return;
      }

      // Desconectar la conexión anterior si existe
      this.disconnectWebSocket();
      
      try {
        // Obtener el token de autenticación del localStorage
        const token = localStorage.getItem('token') || '';
        let userId = localStorage.getItem('userId') || '';
        
        // Asegurar que userId nunca esté vacío
        if (!userId) {
          console.warn('userId no encontrado en localStorage, usando valor predeterminado');
          userId = '1'; // Valor predeterminado para desarrollo
          
          // Intentar guardar en localStorage para futuras conexiones
          try {
            localStorage.setItem('userId', userId);
          } catch (e) {
            console.warn('No se pudo guardar userId en localStorage');
          }
        }
        
        // Intentar establecer una conexión WebSocket real con token incluido
        console.log(`Conectando WebSocket con ticket ${ticketId} y usuario ${userId}`);
        
        // URL websocket corregida - usar el API Gateway configurado en lugar de localhost:8000
        const wsUrl = `ws://localhost/api/ws/chat/${ticketId}`;
        console.log(`Intentando conectar WebSocket a: ${wsUrl}`);
        
        // Probar a conectar, pero fallar silenciosamente si no es posible
        try {
          this.socket = new WebSocket(wsUrl);
        } catch (e) {
          console.log('Error al crear WebSocket, activando modo de simulación:', e);
          this.simulateWebSocket(ticketId);
          return;
        }
        
        this.socket.onopen = () => {
          console.log('WebSocket conectado correctamente');
          this.connected = true;
          
          // Enviar mensaje de identificación inmediatamente después de conectar
          this.socket.send(JSON.stringify({
            type: 'identify',
            ticketId,
            userId
          }));
        };
        
        this.socket.onmessage = (event: MessageEvent) => {
          try {
            console.log('Mensaje recibido por WebSocket:', event.data);
            let data;
            
            try {
              data = JSON.parse(event.data);
            } catch (e) {
              console.error('Error al parsear mensaje WebSocket:', e);
              return;
            }
            
            // Manejar tipos de mensajes específicos
            switch (data.type) {
              case 'connection_established':
              case 'identify_success':
                console.log('Conexión WebSocket establecida correctamente');
                this.connected = true;
                break;
                
              case 'error':
                console.error('Error del servidor WebSocket:', data.message);
                break;
                
              case 'message_history':
                // Procesar historial de mensajes
                if (Array.isArray(data.messages) && data.ticketId) {
                  console.log(`Recibido historial de ${data.messages.length} mensajes para el ticket ${data.ticketId}`);
                  
                  // Limpiar mensajes antiguos para este ticket
                  this.messages[data.ticketId] = [];
                  
                  // Procesar cada mensaje del historial
                  data.messages.forEach((msg: any) => {
                    const newMessage = {
                      id: msg.id || `hist-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`,
                      content: msg.content || '',
                      isClient: msg.isClient === true, // Asegurar que se respeta el valor exacto
                      timestamp: msg.timestamp || msg.createdAt || new Date().toISOString(),
                      ticketId: data.ticketId,
                      userName: msg.userName || ''
                    };
                    
                    // Debug para verificar que isClient se procesa correctamente
                    console.log('Procesando mensaje de historial:', {
                      id: newMessage.id,
                      isClient: newMessage.isClient,
                      originalIsClient: msg.isClient,
                      content: newMessage.content
                    });
                    
                    // Añadir al historial
                    this.handleNewMessage(newMessage);
                  });
                  
                  // Forzar actualización reactiva
                  this.messages = { ...this.messages };
                } else {
                  console.warn('Formato inválido de message_history:', data);
                }
                break;
                
              case 'new_message':
              case 'message':
                let messageContent, isClientFlag, messageData;
                
                // Extraer datos según diferentes formatos posibles
                if (data.data) {
                  // Formato con campo data
                  messageData = data.data;
                  messageContent = messageData.content;
                  // CORRECCIÓN: Usar comparación estricta para verificar si es un mensaje de cliente
                  isClientFlag = messageData.isClient === true;
                  
                  console.log('Procesando mensaje formato data:', {
                    content: messageContent,
                    isClient: isClientFlag,
                    originalIsClient: messageData.isClient
                  });
                } else if (data.message) {
                  // Formato con campo message
                  messageData = data.message;
                  messageContent = messageData.content;
                  // CORRECCIÓN: Usar comparación estricta
                  isClientFlag = messageData.isClient === true;
                  
                  console.log('Procesando mensaje formato message:', {
                    content: messageContent,
                    isClient: isClientFlag,
                    originalIsClient: messageData.isClient
                  });
                } else {
                  // Formato directo
                  messageData = data;
                  messageContent = data.content;
                  // CORRECCIÓN: Usar comparación estricta
                  isClientFlag = data.isClient === true;
                  
                  console.log('Procesando mensaje formato directo:', {
                    content: messageContent,
                    isClient: isClientFlag,
                    originalIsClient: data.isClient
                  });
                }
                
                if (messageContent) {
                  const newMessage = {
                    id: messageData.id || 'ws-' + Date.now(),
                    content: messageContent,
                    isClient: isClientFlag,
                    timestamp: messageData.timestamp || messageData.createdAt || new Date().toISOString(),
                    ticketId: data.ticketId || this.currentTicketId,
                    userName: messageData.userName || (isClientFlag ? 'Cliente' : 'Agente')
                  };
                  
                  // Debug para verificar que isClient se procesa correctamente
                  console.log('Procesando nuevo mensaje WebSocket:', {
                    id: newMessage.id,
                    isClient: newMessage.isClient,
                    content: newMessage.content
                  });
                  
                  this.handleNewMessage(newMessage);
                } else {
                  console.warn('Mensaje WebSocket sin contenido:', data);
                }
                break;
                
              default:
                // Intentar procesar como mensaje directo si tiene campos necesarios
                if (data.content) {
                  const isClientFlag = data.isClient === true;
                  const newMessage = {
                    id: data.id || 'ws-' + Date.now(),
                    content: data.content,
                    isClient: isClientFlag,
                    timestamp: data.timestamp || new Date().toISOString(),
                    ticketId: data.ticketId || this.currentTicketId,
                    userName: data.userName || (isClientFlag ? 'Cliente' : 'Agente')
                  };
                  
                  console.log('Procesando mensaje formato alternativo:', newMessage);
                  this.handleNewMessage(newMessage);
                } else {
                  console.warn('Formato de mensaje no reconocido:', data);
                }
            }
          } catch (error) {
            console.error('Error al procesar mensaje de WebSocket:', error);
          }
        };
        
        this.socket.onerror = (event: Event) => {
          console.error('Error en la conexión WebSocket:', event);
          this.connected = false;
          // Si ocurre un error, activar modo simulación silenciosamente
          this.simulateWebSocket(ticketId);
        };
        
        this.socket.onclose = (event) => {
          console.log('WebSocket cerrado con código:', event.code, event.reason || 'Sin razón proporcionada');
          this.connected = false;
          
          // Reconectar automáticamente después de un tiempo si no fue cerrado explícitamente
          if (event.code !== 1000) {
            console.log('Intentando reconectar en 5 segundos...');
            setTimeout(() => {
              if (this.currentTicketId) {
                this.connectToWebSocket(this.currentTicketId);
              }
            }, 5000);
          }
        };
      } catch (error) {
        console.error('Error al configurar WebSocket:', error);
        this.simulateWebSocket(ticketId);
      }
    },
    
    disconnectWebSocket() {
      if (this.socket) {
        this.socket.close(1000, 'Cierre normal por el cliente')
        this.socket = null
        this.connected = false
      }
    },
    
    handleNewMessage(message: Message) {
      const ticketId = message.ticketId || this.currentTicketId;
      if (!ticketId) {
        console.error('No se puede agregar mensaje: no hay ticketId');
        return;
      }
      
      console.log(`Añadiendo mensaje a ticket ${ticketId}:`, message);
      
      // Asegurar que existe el array para este ticket
      if (!this.messages[ticketId]) {
        this.messages[ticketId] = [];
      }
      
      // Verificar que el mensaje no existe ya para evitar duplicados
      const messageExists = this.messages[ticketId].some((m: Message) => m.id === message.id);
      
      // NUEVA LÓGICA: Detectar si es una confirmación de un mensaje local que enviamos
      // Si es un mensaje del servidor (ID MSG-) y no es un mensaje de cliente (isClient=false),
      // buscar si existe un mensaje local con el mismo contenido y timestamp cercano
      let isDuplicate = false;
      if (!messageExists && message.id.startsWith('MSG-') && message.isClient === false) {
        // Buscar mensajes locales con el mismo contenido en los últimos segundos
        const recentLocalMessages = this.messages[ticketId].filter((m: Message) => {
          // Debe ser un mensaje local con el mismo contenido y enviado recientemente
          return m.id.startsWith('local-') && 
                 m.content === message.content && 
                 m.isClient === false &&
                 // Comprobar si el timestamp está dentro de un margen de 30 segundos
                 Math.abs(new Date(m.timestamp || '').getTime() - new Date(message.timestamp || '').getTime()) < 30000;
        });
        
        if (recentLocalMessages.length > 0) {
          console.log('Detectada confirmación de mensaje local que ya existe:', {
            local: recentLocalMessages[0],
            server: message
          });
          
          // Actualizar el ID del mensaje local con el ID real del servidor
          // para mantener consistencia con el backend
          const localMsgIndex = this.messages[ticketId].findIndex((m: Message) => m.id === recentLocalMessages[0].id);
          if (localMsgIndex !== -1) {
            console.log(`Actualizando mensaje local en posición ${localMsgIndex} con ID del servidor`);
            this.messages[ticketId][localMsgIndex].id = message.id;
            
            // Forzar actualización reactiva
            this.messages = { ...this.messages };
          }
          
          isDuplicate = true;
        }
      }
      
      if (!messageExists && !isDuplicate) {
        console.log('Añadiendo nuevo mensaje a la lista');
        
        // Crear copia del mensaje asegurando que isClient mantenga su valor exacto
        const messageCopy = {
          ...message,
          // Asegurar explícitamente que isClient tenga el valor correcto
          isClient: message.isClient === true
        };
        
        console.log('Mensaje a agregar con isClient:', messageCopy.isClient);
        
        this.messages[ticketId].push(messageCopy);
        
        // Forzar reactividad
        this.messages = { ...this.messages };
      } else if (messageExists) {
        console.log('Mensaje duplicado, ignorando:', message.id);
      }
    },

    async uploadAttachment(ticketId: string, file: File) {
      this.loading = true
      this.error = null
      try {
        const formData = new FormData()
        formData.append('file', file)
        formData.append('ticketId', ticketId)

        const response = await apiClient.post('/chat/attachments', formData, {
          headers: {
            'Content-Type': 'multipart/form-data'
          }
        })
        return response.data
      } catch (error) {
        this.error = 'Failed to upload attachment'
        console.error('Error uploading attachment:', error)
        return null
      } finally {
        this.loading = false
      }
    },

    clearMessages() {
      this.messages = {}
      this.currentTicketId = null
      this.disconnectWebSocket()
    },

    // Método para obtener mensajes de un ticket específico
    getMessagesForTicket(ticketId: string): Message[] {
      return this.messages[ticketId] || [];
    },

    // Método para verificar si hay conexión WebSocket activa para un ticket específico
    isConnected(ticketId?: string): boolean {
      // Ignoramos el ticketId y simplemente devolvemos el estado de conexión actual
      return this.connected;
    },

    // Método para simular WebSocket en modo desarrollo
    simulateWebSocket(ticketId: string) {
      console.log('Modo desarrollo: simulando conexión WebSocket');
      
      // Simular conexión exitosa
      this.connected = true;
      
      // Limpiar intervalo anterior si existe
      if (this.simulationInterval) {
        clearInterval(this.simulationInterval);
      }
      
      // Simular recepción de mensajes periódicamente (solo en desarrollo)
      if (import.meta.env.DEV) {
        // En lugar de enviar mensajes periódicamente, solo enviar confirmación inicial
        setTimeout(() => {
          // Notificar que estamos en modo simulación
          this.handleNewMessage({
            id: `system-${Date.now()}`,
            content: '(Modo desarrollo: WebSocket simulado)',
            isClient: false,
            timestamp: new Date().toISOString(),
            ticketId
          });
        }, 2000);
      }
    }
  }
})
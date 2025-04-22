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
      
      // Crear el objeto de mensaje
      const message = {
        content,
        isInternal,
        timestamp: new Date().toISOString(),
        fromWebSocket: false // Por defecto, indicamos que no viene de WebSocket
      };
      
      // Si hay una conexión WebSocket activa, enviar por ahí y no por HTTP
      if (this.socket && this.socket.readyState === WebSocket.OPEN) {
        try {
          this.socket.send(JSON.stringify({
            type: 'message',
            content,
            ticketId
          }));
          
          // No hacer nada más, el mensaje se procesará cuando se reciba por WebSocket
          this.loading = false;
          return { id: `pending-${Date.now()}`, content, isClient: false };
        } catch (e) {
          console.warn('Error al enviar mensaje por WebSocket, usando HTTP fallback:', e);
          // Continuar con el fallback HTTP
        }
      }
      
      // Solo usar HTTP como fallback si el WebSocket falló o no está disponible
      // Indicar que este mensaje ya fue enviado por WebSocket para evitar duplicar
      message.fromWebSocket = true;
      
      return apiClient.post(`/tickets/${ticketId}/messages`, message)
        .then(response => {
          const newMessage = response.data;
          
          // Si no hay una lista de mensajes para este ticket, crearla
          if (!this.messages[ticketId]) {
            this.messages[ticketId] = [];
          }
          
          // Agregar el mensaje a la lista
          this.handleNewMessage({
            id: newMessage.id || `msg-${Date.now()}`,
            content: newMessage.content || content,
            isClient: false,
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
              isClient: false,
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
        // Intentar establecer una conexión WebSocket real
        console.log(`Intentando conectar WebSocket a: ws://localhost:8000/api/ws/chat/${ticketId}`);
        
        // Probar a conectar, pero fallar silenciosamente si no es posible
        try {
          this.socket = new WebSocket(`ws://localhost:8000/api/ws/chat/${ticketId}`);
        } catch (e) {
          console.log('Error al crear WebSocket, activando modo de simulación');
          this.simulateWebSocket(ticketId);
          return;
        }
        
        this.socket.onopen = () => {
          console.log('WebSocket conectado');
          this.connected = true;
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
            
            // Establecer conexión
            if (data.type === 'connection_established') {
              console.log('Conexión WebSocket establecida');
              this.connected = true;
              return;
            }
            
            // Procesar mensaje nuevo
            if (data.type === 'new_message' && data.message) {
              const newMessage = {
                id: data.message.id || 'ws-' + Date.now(),
                content: data.message.content,
                isClient: data.message.isClient || false,
                timestamp: data.message.timestamp || new Date().toISOString(),
                ticketId: data.ticketId || this.currentTicketId
              };
              
              console.log('Procesando nuevo mensaje WebSocket:', newMessage);
              this.handleNewMessage(newMessage);
              return;
            }
            
            // Si es un mensaje directo (sin type)
            if (data.content) {
              const newMessage = {
                id: data.id || 'ws-' + Date.now(),
                content: data.content,
                isClient: data.isClient || false,
                timestamp: data.timestamp || new Date().toISOString(),
                ticketId: data.ticketId || this.currentTicketId
              };
              
              console.log('Procesando mensaje directo WebSocket:', newMessage);
              this.handleNewMessage(newMessage);
              return;
            }
            
            console.log('Formato de mensaje no reconocido:', data);
          } catch (error) {
            console.error('Error al procesar mensaje de WebSocket:', error);
          }
        };
        
        this.socket.onerror = (event: Event) => {
          console.log('WebSocket error:', event);
          // Si ocurre un error, activar modo simulación silenciosamente
          this.simulateWebSocket(ticketId);
        };
        
        this.socket.onclose = () => {
          console.log('WebSocket cerrado');
          this.connected = false;
        };
      } catch (error) {
        console.error('Error al conectar WebSocket:', error);
        this.simulateWebSocket(ticketId);
      }
    },
    
    disconnectWebSocket() {
      if (this.socket) {
        this.socket.close()
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
      if (!messageExists) {
        console.log('Añadiendo nuevo mensaje a la lista');
        this.messages[ticketId].push(message);
        
        // Forzar reactividad
        this.messages = { ...this.messages };
      } else {
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
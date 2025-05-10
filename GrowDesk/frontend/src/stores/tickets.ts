import { defineStore } from 'pinia'
import apiClient from '@/api/client'
import { useActivityStore } from './activity'
import { useAuthStore } from './auth'

// Constante para almacenar las etiquetas en el localStorage
const TAGS_STORAGE_KEY = 'growdesk_tags'

// añade el key para los tickets
const TICKETS_STORAGE_KEY = 'growdesk_tickets'

export interface Ticket {
  id: string
  title: string
  description: string
  status: 'open' | 'assigned' | 'in_progress' | 'resolved' | 'closed'
  priority: 'LOW' | 'MEDIUM' | 'HIGH' | 'URGENT'
  category: string
  createdBy: string
  assignedTo: string | null
  createdAt: string
  updatedAt: string
  tags?: Tag[] | string[]
}

export interface Tag {
  id: string
  name: string
  color: string
  category?: string
}

interface TicketState {
  tickets: Ticket[]
  currentTicket: Ticket | null
  loading: boolean
  error: string | null
  tags: Tag[]
}

// Helper function to save tags to localStorage
const saveTagsToStorage = (tags: Tag[]) => {
  try {
    localStorage.setItem(TAGS_STORAGE_KEY, JSON.stringify(tags))
    console.log('Tags saved to localStorage:', tags.length)
  } catch (error) {
    console.error('Error saving tags to localStorage:', error)
  }
}

// Helper function to load tags from localStorage
const loadTagsFromStorage = (): Tag[] => {
  try {
    const tagsJson = localStorage.getItem(TAGS_STORAGE_KEY)
    if (tagsJson) {
      const tags = JSON.parse(tagsJson)
      console.log('Tags loaded from localStorage:', tags.length)
      return tags
    }
  } catch (error) {
    console.error('Error loading tags from localStorage:', error)
  }
  return []
}

// Helper function to save tickets to localStorage
const saveTicketsToStorage = (tickets: Ticket[]) => {
  try {
    localStorage.setItem(TICKETS_STORAGE_KEY, JSON.stringify(tickets))
    console.log('Tickets saved to localStorage:', tickets.length)
  } catch (error) {
    console.error('Error saving tickets to localStorage:', error)
  }
}

// Helper function to load tickets from localStorage
const loadTicketsFromStorage = (): Ticket[] => {
  try {
    // Importar la función de validación para limpiar los datos
    import('@/utils/validators').then(({ filterValidTickets }) => {
      try {
        const data = localStorage.getItem(TICKETS_STORAGE_KEY);
        if (data) {
          const parsed = JSON.parse(data);
          if (Array.isArray(parsed)) {
            const valid = filterValidTickets(parsed);
            if (valid.length !== parsed.length) {
              localStorage.setItem(TICKETS_STORAGE_KEY, JSON.stringify(valid));
              console.warn(`Se limpiaron ${parsed.length - valid.length} tickets inválidos`);
            }
          }
        }
      } catch (e) {
        console.error('Error al limpiar tickets del localStorage:', e);
      }
    }).catch(e => console.error('Error al importar utilidades de validación:', e));
    
    const ticketsJson = localStorage.getItem(TICKETS_STORAGE_KEY)
    if (ticketsJson) {
      const parsed = JSON.parse(ticketsJson)
      
      // Verificar que es un array válido
      if (!Array.isArray(parsed)) {
        console.error('Datos de tickets no son un array válido');
        return [];
      }
      
      // Filtrar para asegurar que solo tenemos objetos válidos
      const validTickets = parsed.filter(item => 
        typeof item === 'object' && 
        item !== null && 
        'id' in item && 
        'title' in item &&
        'status' in item
      );
      
      if (validTickets.length !== parsed.length) {
        console.warn(`Se filtraron ${parsed.length - validTickets.length} tickets inválidos`);
        localStorage.setItem(TICKETS_STORAGE_KEY, JSON.stringify(validTickets));
      }
      
      console.log('Tickets loaded from localStorage:', validTickets.length)
      return validTickets
    }
  } catch (error) {
    console.error('Error loading tickets from localStorage:', error)
  }
  return []
}

export const useTicketStore = defineStore('tickets', {
  state: (): TicketState => ({
    tickets: [],
    currentTicket: null,
    loading: false,
    error: null,
    tags: []
  }),

  getters: {
    openTickets: (state: TicketState) => state.tickets.filter((ticket: Ticket) => ticket.status === 'open'),
    inProgressTickets: (state: TicketState) => state.tickets.filter((ticket: Ticket) => ticket.status === 'in_progress'),
    urgentTickets: (state: TicketState) => state.tickets.filter((ticket: Ticket) => ticket.priority === 'URGENT'),
    assignedTickets: (state: TicketState) => state.tickets.filter((ticket: Ticket) => ticket.assignedTo),
    resolvedTickets: (state: TicketState) => state.tickets.filter((ticket: Ticket) => ticket.status === 'resolved'),
    closedTickets: (state: TicketState) => state.tickets.filter((ticket: Ticket) => ticket.status === 'closed'),
    ticketsPerDay: (state: TicketState) => {
      const ticketsByDay = new Map<string, number>()
      
      state.tickets.forEach((ticket: Ticket) => {
        const date = new Date(ticket.createdAt).toISOString().split('T')[0]
        ticketsByDay.set(date, (ticketsByDay.get(date) || 0) + 1)
      })
      
      // Convertir a array y ordenar por fecha
      return Array.from(ticketsByDay.entries())
        .sort(([dateA], [dateB]) => new Date(dateB).getTime() - new Date(dateA).getTime())
    }
  },

  actions: {
    async fetchTickets() {
      this.loading = true
      
      try {
        // First try to load tickets from localStorage
        const savedTickets = loadTicketsFromStorage()
        
        if (savedTickets && savedTickets.length > 0) {
          this.tickets = savedTickets
          this.loading = false
          return
        }
        
        const response = await apiClient.get('/tickets')
        this.tickets = response.data
        
        // Save to localStorage
        saveTicketsToStorage(this.tickets)
      } catch (error) {
        console.error('Error fetching tickets:', error)
        this.error = 'Error al cargar los tickets'
      } finally {
        this.loading = false
      }
    },

    async fetchUserTickets(userId: string) {
      this.loading = true
      this.error = null
      
      try {
        console.log(`Buscando tickets asignados al usuario ${userId}`);
        
        // Primero cargar todos los tickets si aún no están cargados
        if (this.tickets.length === 0) {
          console.log('La lista de tickets está vacía, cargando todos los tickets primero');
          await this.fetchTickets();
        }
        
        // Intentar encontrar el ticket específico que mencionó el usuario
        const specificTicket = this.tickets.find((ticket: Ticket) => ticket.id === 'TICKET-20250327041753');
        if (specificTicket) {
          console.log('Ticket específico mencionado encontrado:', specificTicket);
          console.log('Estado de asignación:', specificTicket.assignedTo === userId ? 'Asignado al usuario actual' : 'No asignado al usuario actual');
        } else {
          console.log('Ticket específico mencionado no encontrado en el sistema');
        }
        
        // Filtrar los tickets donde el usuario es asignado con log detallado
        console.log(`Total de tickets en el sistema: ${this.tickets.length}`);
        console.log('IDs de todos los tickets:', this.tickets.map((t: Ticket) => t.id).join(', '));
        console.log('Detalles de asignación de todos los tickets:');
        this.tickets.forEach((ticket: Ticket) => {
          console.log(`Ticket ${ticket.id}: assignedTo=${ticket.assignedTo}, userId=${userId}, match=${ticket.assignedTo === userId}`);
        });
        
        const assignedTickets = this.tickets.filter((ticket: Ticket) => {
          // Comparación insensible a mayúsculas/minúsculas y trim para mayor flexibilidad
          const ticketAssignedTo = String(ticket.assignedTo || '').trim();
          const currentUserId = String(userId || '').trim();
          const isMatch = ticketAssignedTo.toLowerCase() === currentUserId.toLowerCase();
          
          console.log(`Comparación ticket ${ticket.id}: "${ticketAssignedTo}" vs "${currentUserId}" = ${isMatch}`);
          return isMatch;
        });
        
        console.log(`Encontrados ${assignedTickets.length} tickets asignados al usuario ${userId}:`, assignedTickets);
        
        if (assignedTickets.length > 0) {
          return assignedTickets;
        }
        
        // Si no hay tickets asignados, mostrar mensaje claro
        console.warn('No se encontraron tickets asignados para el usuario:', userId);
        return [];
      } catch (error) {
        this.error = 'Error al cargar los tickets del usuario';
        console.error('Error fetching user tickets:', error);
        return [];
      } finally {
        this.loading = false;
      }
    },

    async fetchTicket(id: string) {
      this.loading = true
      this.error = null
      
      console.log('Intentando cargar ticket:', id)
      
      try {
        const response = await apiClient.get(`/tickets/${id}`)
        console.log('Respuesta al cargar ticket:', response.data)
        this.currentTicket = response.data
      } catch (error) {
        console.error('Error al cargar ticket:', error)
        this.error = 'Failed to fetch ticket'
        
        // Para depuración, crear un ticket ficticio
        if (import.meta.env.DEV) {
          console.warn('Usando ticket de prueba en modo desarrollo')
          this.currentTicket = {
            id: id,
            title: 'Ticket de prueba',
            description: 'Este es un ticket de prueba generado automáticamente.',
            status: 'open',
            priority: 'medium',
            category: 'technical',
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString(),
            createdBy: 'system',
            assignedTo: null,
            customer: {
              name: 'Cliente de prueba',
              email: 'cliente@example.com'
            }
          }
        } else {
          throw error
        }
      } finally {
        this.loading = false
      }
    },

    async createTicket(ticketData: {
      title: string
      description: string
      priority: Ticket['priority']
      category: string
    }) {
      // Esta función ha sido deshabilitada porque conceptualmente
      // los tickets solo deben ser creados por los clientes a través del widget.
      console.warn('La creación de tickets desde el panel de soporte está deshabilitada. Los tickets deben ser creados por los clientes a través del widget.');
      this.error = 'La creación de tickets desde el panel de soporte no está permitida';
      return null;
    },

    async updateTicket(ticketOrId: string | Ticket, ticketData?: Partial<Ticket>) {
      this.loading = true
      this.error = null
      
      let id: string;
      let updateData: Partial<Ticket>;
      let originalTicket: Ticket | undefined;
      
      // Comprobar si estamos recibiendo un objeto de ticket completo o un ID con datos parciales
      if (typeof ticketOrId === 'string') {
        // Caso tradicional: id + datos parciales
        id = ticketOrId;
        updateData = ticketData || {};
        originalTicket = this.tickets.find((t: Ticket) => t.id === id);
      } else {
        // Caso nuevo: objeto completo de ticket
        id = ticketOrId.id;
        updateData = ticketOrId;
        originalTicket = this.tickets.find((t: Ticket) => t.id === id);
      }
      
      console.log(`Iniciando actualización de ticket: ${id} con datos:`, updateData);
      
      try {
        if (!id) {
          throw new Error('ID de ticket requerido para actualizar');
        }
        
        // Normalizar prioridad si está presente en los datos de actualización
        if (updateData.priority) {
          const originalPriority = updateData.priority;
          updateData.priority = this.normalizePriority(updateData.priority) as Ticket['priority'];
          console.log(`Prioridad normalizada: de "${originalPriority}" a "${updateData.priority}"`);
        }
        
        // Asegurarse de que id no tenga espacios ni caracteres especiales
        const cleanId = id.trim();
        console.log(`Llamando a API para actualizar ticket. URL: '/tickets/${cleanId}'`);
        
        // Usar directamente la URL limpia
        const response = await apiClient.put(`/tickets/${cleanId}`, updateData)
        console.log(`Actualización exitosa, respuesta:`, response.data);
        
        if (!response.data) {
          throw new Error('La respuesta no contiene datos válidos');
        }
        
        // Actualizar en la lista de tickets
        const index = this.tickets.findIndex((t: Ticket) => t.id === id)
        if (index !== -1) {
          this.tickets[index] = response.data
          console.log(`Ticket actualizado en el índice ${index} de la colección`);
          
          // Actualizar también el ticket actual si coincide el ID
          if (this.currentTicket && this.currentTicket.id === id) {
            this.currentTicket = response.data;
            console.log('Ticket actual actualizado con nueva información');
          }
          
          // Registrar actividad de actualización
          const activityStore = useActivityStore()
          const authStore = useAuthStore()
          const currentUserId = authStore.user?.id
          
          if (currentUserId) {
            // Si se cambió el estado, registrar una actividad específica
            if (updateData.status && originalTicket && originalTicket.status !== updateData.status) {
              await activityStore.logActivity({
                userId: currentUserId,
                type: 'ticket_status_changed',
                targetId: id,
                description: `Cambió el estado del ticket #${id} de ${this.translateStatus(originalTicket.status)} a ${this.translateStatus(updateData.status)}`,
                metadata: {
                  oldStatus: originalTicket.status,
                  newStatus: updateData.status
                }
              })
            }
            
            // Si se cambió la prioridad, registrar una actividad específica
            if (updateData.priority && originalTicket && originalTicket.priority !== updateData.priority) {
              await activityStore.logActivity({
                userId: currentUserId,
                type: 'ticket_priority_changed',
                targetId: id,
                description: `Cambió la prioridad del ticket #${id} de ${originalTicket.priority} a ${updateData.priority}`,
                metadata: {
                  oldPriority: originalTicket.priority,
                  newPriority: updateData.priority
                }
              })
            }
            
            // Para cualquier otra actualización
            if (!updateData.status && !updateData.priority) {
              await activityStore.logActivity({
                userId: currentUserId,
                type: 'ticket_updated',
                targetId: id,
                description: `Actualizó la información del ticket #${id}`
              })
            }
          }
        } else {
          console.warn(`Ticket ${id} no encontrado en la colección de tickets`);
        }
        
        // Actualizar ticket actual si corresponde
        if (this.currentTicket?.id === id) {
          console.log('Actualizando ticket actual con datos nuevos');
          this.currentTicket = response.data
        }
        
        return response.data
      } catch (error: any) {
        // Mejorar el mensaje de error con información más detallada
        let errorMessage = 'Error al actualizar ticket';
        
        if (error.response) {
          // Error desde el servidor
          console.error(`Error de respuesta al actualizar ticket ${id}:`, {
            status: error.response.status,
            statusText: error.response.statusText,
            data: error.response.data
          });
          
          if (error.response.data && error.response.data.error) {
            errorMessage = `Error al actualizar ticket: ${error.response.data.error}`;
          } else if (error.response.status === 403) {
            errorMessage = 'No tienes permiso para actualizar este ticket';
          } else if (error.response.status === 404) {
            errorMessage = 'El ticket no fue encontrado';
          } else if (error.response.status >= 500) {
            errorMessage = 'Error del servidor al actualizar ticket';
          }
        } else if (error.request) {
          // La solicitud fue hecha pero no se recibió respuesta
          console.error(`No hubo respuesta del servidor al actualizar ticket ${id}:`, error.request);
          errorMessage = 'No se recibió respuesta del servidor';
        } else {
          // Otros errores
          console.error(`Error general al actualizar ticket ${id}:`, error.message);
          errorMessage = `Error al actualizar ticket: ${error.message}`;
        }
        
        this.error = errorMessage;
        
        // Para desarrollo, simular una actualización exitosa
        if (import.meta.env.DEV) {
          console.log('Modo desarrollo: Simulando actualización exitosa');
          const index = this.tickets.findIndex((t: Ticket) => t.id === id)
          if (index !== -1) {
            // Crear versión actualizada del ticket
            const updatedTicket = {
              ...this.tickets[index],
              ...updateData,
              updatedAt: new Date().toISOString()
            }
            
            // Actualizar en el array
            this.tickets[index] = updatedTicket
            
            // Si es el ticket actual, actualizar también
            if (this.currentTicket?.id === id) {
              this.currentTicket = updatedTicket
            }
            
            console.log('Actualización simulada completada', updatedTicket);
            return updatedTicket
          }
        }
        
        throw error; // Re-lanzar el error para permitir manejo externo
      } finally {
        this.loading = false
      }
    },

    async deleteTicket(id: string) {
      this.loading = true
      this.error = null
      try {
        await apiClient.delete(`/tickets/${id}`)
        this.tickets = this.tickets.filter((t: Ticket) => t.id !== id)
        if (this.currentTicket?.id === id) {
          this.currentTicket = null
        }
        
        // Registrar actividad de eliminación de ticket
        const activityStore = useActivityStore()
        const authStore = useAuthStore()
        const currentUserId = authStore.user?.id
        
        if (currentUserId) {
          await activityStore.logActivity({
            userId: currentUserId,
            type: 'ticket_closed',
            targetId: id,
            description: `Eliminó el ticket #${id}`
          })
        }
        
        return true
      } catch (error) {
        this.error = 'Failed to delete ticket'
        console.error('Error deleting ticket:', error)
        return false
      } finally {
        this.loading = false
      }
    },

    async assignTicket(id: string, agentId: string) {
      this.loading = true
      this.error = null
      console.log(`Iniciando asignación de ticket: ${id} a usuario: ${agentId}`);
      
      try {
        if (!id) {
          throw new Error('ID de ticket requerido para asignar');
        }
        
        if (!agentId) {
          throw new Error('ID de agente requerido para asignar');
        }
        
        console.log(`Llamando a API para asignar ticket ${id} a ${agentId}`);
        
        // Preparar datos de actualización
        const updateData = {
          status: 'assigned' as const,
          assignedTo: agentId
        };
        
        console.log('Datos de actualización para asignación:', updateData);
        
        try {
          const response = await this.updateTicket(id, updateData);
          console.log(`Asignación exitosa, respuesta:`, response);
          
          if (!response) {
            throw new Error('No se recibieron datos en la respuesta de asignación');
          }
          
          // Registrar actividad de asignación
          const activityStore = useActivityStore()
          const authStore = useAuthStore()
          const currentUserId = authStore.user?.id
          
          if (currentUserId) {
            await activityStore.logActivity({
              userId: currentUserId,
              type: 'ticket_assigned',
              targetId: id,
              description: `Asignó el ticket #${id} al usuario #${agentId}`,
              metadata: {
                assigneeId: agentId
              }
            })
          }
          
          return response;
        } catch (updateError: any) {
          // Manejar error específico de la actualización
          throw new Error(`Error al asignar: ${updateError.message || 'Error de actualización'}`);
        }
      } catch (error: any) {
        // Crear mensaje de error detallado
        let errorMessage = 'Error al asignar ticket';
        
        if (error.message) {
          errorMessage = `Error al asignar ticket: ${error.message}`;
        }
        
        this.error = errorMessage;
        console.error('Error detallado durante la asignación del ticket:', error);
        
        // Para desarrollo, simular una actualización exitosa
        if (import.meta.env.DEV) {
          console.log('Modo desarrollo: Simulando asignación exitosa');
          const index = this.tickets.findIndex((t: Ticket) => t.id === id)
          if (index !== -1) {
            const updatedTicket = {
              ...this.tickets[index],
              status: 'assigned',
              assignedTo: agentId,
              updatedAt: new Date().toISOString()
            }
            
            // Actualizar en el array
            this.tickets[index] = updatedTicket
            
            // Si es el ticket actual, actualizar también
            if (this.currentTicket?.id === id) {
              this.currentTicket = updatedTicket
            }
            
            console.log('Asignación simulada completada', updatedTicket);
            return updatedTicket;
          }
        }
        
        throw error; // Re-lanzar el error para permitir manejo externo
      } finally {
        this.loading = false
      }
    },

    async updateTicketStatus(id: string, status: Ticket['status']) {
      // Registrar el estado anterior
      const originalTicket = this.tickets.find((t: Ticket) => t.id === id);
      const oldStatus = originalTicket?.status;
      
      try {
        // Añade el index del ticket en el array
        const index = this.tickets.findIndex((t: Ticket) => t.id === id);
        
        if (index === -1) {
          throw new Error('Ticket not found');
        }
        
        // Crea un nuevo objeto de ticket con el estado actualizado mientras se preservan todas las demás propiedades
        const updatedTicket = {
          ...this.tickets[index],
          status,
          updatedAt: new Date().toISOString()
        };
        
        // Actualiza el ticket en el array
        this.tickets[index] = updatedTicket;
        
        // Guarda los tickets en el localStorage para persistir los cambios
        saveTicketsToStorage(this.tickets);
        
        // Intenta actualizar el ticket via API
        try {
          await this.updateTicket(id, { status });
        } catch (apiError) {
          console.error('API update failed, but local changes were preserved:', apiError);
          // Continúa la ejecución incluso si la llamada a la API falla - ya hemos actualizado localmente
        }
        
        // Si la actualización fue exitosa y hubo un cambio de estado, registra la actividad
        if (oldStatus !== status) {
          const activityStore = useActivityStore();
          const authStore = useAuthStore();
          const currentUserId = authStore.user?.id;
          
          if (currentUserId) {
            let activityType: 'ticket_status_changed' | 'ticket_closed' | 'ticket_reopened' = 'ticket_status_changed';
            
            // Determina el tipo de actividad basado en el cambio de estado
            if (status === 'closed') {
              activityType = 'ticket_closed';
            } else if (oldStatus === 'closed') {
              activityType = 'ticket_reopened';
            }
            
            await activityStore.logActivity({
              userId: currentUserId,
              type: activityType,
              targetId: id,
              description: `Cambió el estado del ticket #${id} de ${this.translateStatus(oldStatus || 'desconocido')} a ${this.translateStatus(status)}`,
              metadata: {
                oldStatus,
                newStatus: status
              }
            });
          }
        }
        
        return updatedTicket;
      } catch (error) {
        console.error('Error updating ticket status:', error);
        throw error;
      }
    },

    // Función para normalizar la prioridad (asegurarse de que esté en mayúsculas)
    normalizePriority(priority: string): string {
      const p = String(priority).toLowerCase()
      if (p === 'high' || p === 'alta') return 'HIGH'
      if (p === 'medium' || p === 'media' || p === 'normal') return 'MEDIUM'
      if (p === 'low' || p === 'baja') return 'LOW'
      if (p === 'urgent' || p === 'urgente' || p === 'crítica' || p === 'critical') return 'URGENT'
      return 'MEDIUM'
    },
    
    // Método de ayuda para traducir estados
    translateStatus(status: string): string {
      const statusMap: Record<string, string> = {
        'open': 'Abierto',
        'assigned': 'Asignado',
        'in_progress': 'En Progreso',
        'resolved': 'Resuelto',
        'closed': 'Cerrado'
      }
      
      return statusMap[status] || status
    },
    
    // Gestión de etiquetas
    async fetchTags() {
      this.loading = true
      
      try {
        // Primero intenta cargar desde localStorage
        const savedTags = loadTagsFromStorage()
        
        if (savedTags && savedTags.length > 0) {
          this.tags = savedTags
          return
        }
        
        // En un entorno real, esto haría una llamada a la API para obtener las etiquetas
        // En este caso, vamos a simular la llamada con datos de ejemplo
        console.log('No tags found in localStorage, loading example tags...')
        
        // En un entorno de desarrollo, podemos crear algunas etiquetas de ejemplo
        if (import.meta.env.DEV && this.tags.length === 0) {
          this.tags = [
            { id: '1', name: 'Bug', color: '#ff0000', category: 'technical' },
            { id: '2', name: 'Mejora', color: '#00ff00', category: 'feature' },
            { id: '3', name: 'Consulta', color: '#0000ff', category: 'general' },
            { id: '4', name: 'Documentación', color: '#ffff00', category: 'documentation' },
            { id: '5', name: 'Crítico', color: '#ff00ff', category: 'technical' }
          ]
          
          // Guarda las etiquetas de ejemplo en localStorage
          saveTagsToStorage(this.tags)
        }
      } catch (error) {
        console.error('Error al obtener etiquetas:', error)
        this.error = 'Error al cargar las etiquetas'
      } finally {
        this.loading = false
      }
    },
    
    async createTag(tagData: Omit<Tag, 'id'>) {
      try {
        // En un entorno real, esto haría una llamada a la API para crear una etiqueta
        console.log('Creando etiqueta:', tagData)
        
        // Generar un ID único para la nueva etiqueta
        const id = Date.now().toString()
        
        // Crear la nueva etiqueta
        const newTag: Tag = {
          id,
          ...tagData
        }
        
        // Añadir la etiqueta a la lista
        this.tags.push(newTag)
        
        // Guarda en localStorage
        saveTagsToStorage(this.tags)
        
        return newTag
      } catch (error) {
        console.error('Error al crear etiqueta:', error)
        this.error = 'Error al crear la etiqueta'
        throw error
      }
    },
    
    async updateTag(id: string, tagData: Partial<Tag>) {
      try {
        // En un entorno real, esto haría una llamada a la API para actualizar una etiqueta
        console.log('Actualizando etiqueta:', id, tagData)
        
        // Buscar la etiqueta a actualizar
        const index = this.tags.findIndex((tag: Tag) => tag.id === id)
        if (index === -1) {
          throw new Error('Etiqueta no encontrada')
        }
        
        // Actualizar la etiqueta
        this.tags[index] = {
          ...this.tags[index],
          ...tagData
        }
        
        // Update tag references in all tickets that use this tag
        this.tickets.forEach((ticket: Ticket) => {
          if (ticket.tags && Array.isArray(ticket.tags)) {
            ticket.tags = ticket.tags.map((tag: Tag | string) => {
              if (typeof tag === 'string') {
                return tag === id ? this.tags[index] : tag
              } else {
                // If it's a tag object, update it if the ID matches
                return tag.id === id ? this.tags[index] : tag
              }
            })
          }
        })
        
        // Guarda las etiquetas en localStorage
        saveTagsToStorage(this.tags)
        
        // Guarda los tickets en localStorage para persistir los cambios
        saveTicketsToStorage(this.tickets)
        
        return this.tags[index]
      } catch (error) {
        console.error('Error al actualizar etiqueta:', error)
        this.error = 'Error al actualizar la etiqueta'
        throw error
      }
    },
    
    async deleteTag(id: string) {
      try {
        // En un entorno real, esto haría una llamada a la API para eliminar una etiqueta
        console.log('Eliminando etiqueta:', id)
        
        // Eliminar la etiqueta de la lista
        this.tags = this.tags.filter((tag: Tag) => tag.id !== id)
        
        // También eliminar la etiqueta de todos los tickets que la tengan
        this.tickets.forEach((ticket: Ticket) => {
          if (ticket.tags) {
            // Crea un array correctamente tipado basado en lo que detectamos
            const newTags = ticket.tags.filter((tagId: Tag | string) => 
              typeof tagId === 'string' ? tagId !== id : tagId.id !== id
            );
            
            // Comprueba si el primer elemento es una cadena o un objeto Tag para determinar el tipo de array
            if (newTags.length > 0 && typeof newTags[0] === 'string') {
              // Es un array de cadenas
              ticket.tags = newTags as string[];
            } else {
              // Es un array de objetos Tag o un array vacío
              ticket.tags = newTags as Tag[];
            }
          }
        })
        
        // Guarda en localStorage
        saveTagsToStorage(this.tags)
        saveTicketsToStorage(this.tickets)
        
        return true
      } catch (error) {
        console.error('Error al eliminar etiqueta:', error)
        this.error = 'Error al eliminar la etiqueta'
        throw error
      }
    },
    
    async addTagToTicket(ticketId: string, tagId: string) {
      try {
        // En un entorno real, esto haría una llamada a la API para añadir una etiqueta a un ticket
        console.log('Añadiendo etiqueta a ticket:', ticketId, tagId)
        
        // Buscar el ticket
        const ticketIndex = this.tickets.findIndex((ticket: Ticket) => ticket.id === ticketId)
        if (ticketIndex === -1) {
          throw new Error('Ticket no encontrado')
        }
        
        // Buscar la etiqueta
        const tag = this.tags.find((tag: Tag) => tag.id === tagId)
        if (!tag) {
          throw new Error('Etiqueta no encontrada')
        }
        
        // Verificar si la etiqueta ya está en el ticket
        const ticket = this.tickets[ticketIndex]
        if (!ticket.tags) {
          ticket.tags = []
        }
        
        // Si la etiqueta ya está en el ticket, no hacer nada
        const hasTag = ticket.tags.some((t: Tag | string) => 
          typeof t === 'string' ? t === tagId : t.id === tagId
        )
        
        if (!hasTag) {
          // Añadir la etiqueta al ticket
          ticket.tags.push(tag)
          
          // Guarda las etiquetas en localStorage
          saveTagsToStorage(this.tags)
          
          // Guarda los tickets en localStorage para persistir la relación
          saveTicketsToStorage(this.tickets)
        }
        
        return ticket
      } catch (error) {
        console.error('Error al añadir etiqueta a ticket:', error)
        this.error = 'Error al añadir la etiqueta al ticket'
        throw error
      }
    },
    
    async removeTagFromTicket(ticketId: string, tagId: string) {
      try {
        // En un entorno real, esto haría una llamada a la API para eliminar una etiqueta de un ticket
        console.log('Eliminando etiqueta de ticket:', ticketId, tagId)
        
        // Buscar el ticket
        const ticketIndex = this.tickets.findIndex((ticket: Ticket) => ticket.id === ticketId)
        if (ticketIndex === -1) {
          throw new Error('Ticket no encontrado')
        }
        
        // Verificar si la etiqueta está en el ticket
        const ticket = this.tickets[ticketIndex]
        if (!ticket.tags) {
          return ticket
        }
        
        // Eliminar la etiqueta del ticket
        ticket.tags = ticket.tags.filter((tag: Tag | string) => 
          typeof tag === 'string' ? tag !== tagId : tag.id !== tagId
        )
        
        // Guarda los tickets en localStorage para persistir los cambios
        saveTicketsToStorage(this.tickets)
        
        return ticket
      } catch (error) {
        console.error('Error al eliminar etiqueta de ticket:', error)
        this.error = 'Error al eliminar la etiqueta del ticket'
        throw error
      }
    }
  }
})
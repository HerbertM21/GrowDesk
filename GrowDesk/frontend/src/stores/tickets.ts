import { defineStore } from 'pinia'
import apiClient from '@/api/client'
import { useActivityStore } from './activity'
import { useAuthStore } from './auth'

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
}

interface TicketState {
  tickets: Ticket[]
  currentTicket: Ticket | null
  loading: boolean
  error: string | null
}

export const useTicketStore = defineStore('tickets', {
  state: (): TicketState => ({
    tickets: [],
    currentTicket: null,
    loading: false,
    error: null
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
      this.error = null
      try {
        console.log('Obteniendo todos los tickets del sistema');
        const response = await apiClient.get('/tickets')
        
        // Asegurarse de que tenemos una respuesta válida
        if (response.data && Array.isArray(response.data)) {
          console.log(`API devolvió ${response.data.length} tickets`);
          
          this.tickets = response.data.map((ticket: any) => ({
            ...ticket,
            status: ticket.status || 'open',
            priority: ticket.priority || 'MEDIUM'
          }))
        } else {
          console.warn('La API no devolvió tickets o la respuesta no es un array');
          this.tickets = [];
        }
      } catch (error) {
        this.error = 'Error al cargar los tickets'
        console.error('Error fetching tickets:', error)
        this.tickets = [];
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
      
      // Actualizar el ticket
      const updatedTicket = await this.updateTicket(id, { status });
      
      // Si la actualización fue exitosa y hubo un cambio de estado, registrar la actividad
      if (updatedTicket && oldStatus !== status) {
        const activityStore = useActivityStore();
        const authStore = useAuthStore();
        const currentUserId = authStore.user?.id;
        
        if (currentUserId) {
          let activityType: 'ticket_status_changed' | 'ticket_closed' | 'ticket_reopened' = 'ticket_status_changed';
          
          // Determinar tipo de actividad según el cambio de estado
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
    },

    // Método para normalizar prioridades en diferentes formatos
    normalizePriority(priority: string): string {
      if (!priority) return 'medium';
      
      // Convertir a minúsculas para comparación
      const lowerPriority = typeof priority === 'string' ? priority.toLowerCase() : '';
      
      // Mapear diferentes formatos posibles al formato estándar
      if (lowerPriority === 'baja' || lowerPriority === 'low') return 'low';
      if (lowerPriority === 'media' || lowerPriority === 'medium') return 'medium';
      if (lowerPriority === 'alta' || lowerPriority === 'high') return 'high';
      if (lowerPriority === 'urgente' || lowerPriority === 'urgent') return 'urgent';
      
      // Si no coincide con ninguno, devolver medio por defecto
      console.log(`Store: Prioridad no reconocida "${priority}", utilizando "medium" por defecto`);
      return 'medium';
    },
    
    // Método de ayuda para traducir estados
    translateStatus(status: string): string {
      const statuses: Record<string, string> = {
        'open': 'Abierto',
        'assigned': 'Asignado',
        'in_progress': 'En Progreso',
        'resolved': 'Resuelto',
        'closed': 'Cerrado'
      };
      return statuses[status] || status;
    }
  }
})
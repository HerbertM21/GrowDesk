import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import apiClient from '@/api/client';
import type { User } from './users';

// Definición de tipos para actividades
export interface Activity {
  id: string;
  userId: string;
  type: string;
  targetId?: string; // ID del objeto relacionado (ticket, usuario, etc.)
  description: string;
  timestamp: string;
  metadata?: Record<string, any>; // Datos adicionales específicos de la actividad
}

// Tipos de actividad posibles
export type ActivityType = 
  | 'ticket_created'
  | 'ticket_updated'
  | 'ticket_closed'
  | 'ticket_reopened'
  | 'ticket_assigned'
  | 'ticket_status_changed'
  | 'ticket_priority_changed'
  | 'comment_added'
  | 'profile_updated'
  | 'user_login'
  | 'user_logout';

// Interfaz para el estado del store
interface ActivityState {
  activities: Activity[];
  loading: boolean;
  error: string | null;
}

// Definición del store
export const useActivityStore = defineStore('activity', {
  state: (): ActivityState => ({
    activities: [],
    loading: false,
    error: null
  }),

  getters: {
    // Obtener actividades de un usuario específico
    getUserActivities: (state) => (userId: string) => {
      return state.activities.filter((activity) => activity.userId === userId);
    },
    
    // Obtener actividades relacionadas con un ticket específico
    getTicketActivities: (state) => (ticketId: string) => {
      return state.activities.filter((activity) => activity.targetId === ticketId);
    },
    
    // Obtener actividades recientes (últimas 20)
    getRecentActivities: (state) => {
      return [...state.activities]
        .sort((a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime())
        .slice(0, 20);
    }
  },

  actions: {
    // Cargar todas las actividades del sistema
    async fetchActivities() {
      this.loading = true;
      this.error = null;
      
      try {
        const response = await apiClient.get('/activities');
        
        if (response.data && Array.isArray(response.data)) {
          this.activities = response.data;
        } else {
          console.warn('La respuesta de actividades no es un array o está vacía');
          this.activities = [];
        }
      } catch (error) {
        console.error('Error al cargar actividades:', error);
        this.error = 'Error al cargar actividades';
        
        // No generar actividades de ejemplo automáticamente
        this.activities = [];
      } finally {
        this.loading = false;
      }
      
      return this.activities;
    },
    
    // Cargar actividades de un usuario específico
    async fetchUserActivities(userId: string) {
      this.loading = true;
      this.error = null;
      
      try {
        // Asegurarse de que tenemos actividades cargadas
        if (this.activities.length === 0) {
          await this.fetchActivities();
        }
        
        // Filtrar las actividades del usuario
        const userActivities = this.activities.filter(
          (activity) => activity.userId === userId
        );
        
        console.log(`Encontradas ${userActivities.length} actividades para el usuario ${userId}`);
        
        // No generar actividades falsas
        return userActivities;
      } catch (error) {
        console.error('Error al cargar actividades del usuario:', error);
        this.error = 'Error al cargar actividades del usuario';
        return [];
      } finally {
        this.loading = false;
      }
    },
    
    // Registrar una nueva actividad
    async logActivity({
      userId, 
      type, 
      targetId = '', 
      description, 
      metadata = {}
    }: {
      userId: string;
      type: string;
      targetId?: string;
      description: string;
      metadata?: Record<string, any>;
    }) {
      this.loading = true;
      this.error = null;
      
      // Crear el objeto de actividad
      const activityData = {
        userId,
        type,
        targetId,
        description,
        timestamp: new Date().toISOString(),
        metadata
      };
      
      // En modo desarrollo, siempre usar la simulación local
      if (import.meta.env.DEV) {
        try {
          const localActivity = {
            id: `dev-${Date.now()}`,
            ...activityData,
            metadata: {
              ...metadata,
              _simulated: true
            }
          };
          console.log('Actividad simulada registrada (modo desarrollo):', localActivity);
          this.activities.push(localActivity);
          return localActivity;
        } catch (error) {
          console.error('Error al simular actividad:', error);
          this.error = 'Error al simular registro de actividad';
          return null;
        } finally {
          this.loading = false;
        }
      }
      
      // En modo producción, intentar registrar en la API
      try {
        const response = await apiClient.post('/activities', activityData);
        
        if (response.data) {
          // Añadir a la lista local
          this.activities.push(response.data);
          return response.data;
        } else {
          console.warn('La respuesta al crear actividad está vacía');
          // Añadir localmente con ID generado
          const localActivity = {
            ...activityData,
            id: `local-${Date.now()}`
          };
          this.activities.push(localActivity);
          return localActivity;
        }
      } catch (error) {
        console.error('Error al registrar actividad:', error);
        this.error = 'Error al registrar actividad';
        
        // Crear actividad local de respaldo
        const fallbackActivity = {
          id: `fallback-${Date.now()}`,
          ...activityData
        };
        this.activities.push(fallbackActivity);
        return fallbackActivity;
      } finally {
        this.loading = false;
      }
    },
    
    // Crear actividades de ejemplo para desarrollo
    initMockActivities() {
      const mockUsers = ['1', '2', '3']; // IDs de usuarios de ejemplo
      const mockTickets = ['101', '102', '103', '104', '105']; // IDs de tickets de ejemplo
      const activities = [];
      
      // Generar actividades para cada usuario
      for (const userId of mockUsers) {
        const userActivities = this.generateMockUserActivities(userId, mockTickets);
        activities.push(...userActivities);
      }
      
      this.activities = activities;
    },
    
    // Generar actividades de ejemplo para un usuario específico
    generateMockUserActivities(userId: string, ticketIds?: string[]) {
      const activities: Activity[] = [];
      const tickets = ticketIds || [`10${userId}1`, `10${userId}2`, `10${userId}3`];
      
      // Añadir actividad de login
      activities.push({
        id: `mock-login-${Date.now()}-${userId}`,
        userId,
        type: 'user_login',
        description: 'Inicio de sesión en el sistema',
        timestamp: new Date(Date.now() - Math.floor(Math.random() * 24) * 60 * 60 * 1000).toISOString()
      });
      
      // Añadir creación de tickets
      for (const ticketId of tickets) {
        // Crear ticket
        activities.push({
          id: `mock-create-${Date.now()}-${ticketId}`,
          userId,
          type: 'ticket_created',
          targetId: ticketId,
          description: `Creó el ticket #${ticketId}: "Problema con el sistema"`,
          timestamp: new Date(Date.now() - Math.floor(Math.random() * 20) * 24 * 60 * 60 * 1000).toISOString()
        });
        
        // Añadir comentario
        if (Math.random() > 0.3) {
          activities.push({
            id: `mock-comment-${Date.now()}-${ticketId}`,
            userId,
            type: 'comment_added',
            targetId: ticketId,
            description: `Comentó en el ticket #${ticketId}: "Añadiendo más información sobre el problema"`,
            timestamp: new Date(Date.now() - Math.floor(Math.random() * 10) * 24 * 60 * 60 * 1000).toISOString()
          });
        }
        
        // Cambiar prioridad
        if (Math.random() > 0.5) {
          activities.push({
            id: `mock-priority-${Date.now()}-${ticketId}`,
            userId,
            type: 'ticket_priority_changed',
            targetId: ticketId,
            description: `Cambió la prioridad del ticket #${ticketId} a Alta`,
            timestamp: new Date(Date.now() - Math.floor(Math.random() * 8) * 24 * 60 * 60 * 1000).toISOString(),
            metadata: {
              oldPriority: 'MEDIUM',
              newPriority: 'HIGH'
            }
          });
        }
        
        // Asignar ticket
        if (Math.random() > 0.4) {
          const assigneeId = String(Math.floor(Math.random() * 3) + 1);
          activities.push({
            id: `mock-assign-${Date.now()}-${ticketId}`,
            userId,
            type: 'ticket_assigned',
            targetId: ticketId,
            description: `Asignó el ticket #${ticketId} al usuario #${assigneeId}`,
            timestamp: new Date(Date.now() - Math.floor(Math.random() * 7) * 24 * 60 * 60 * 1000).toISOString(),
            metadata: {
              assigneeId
            }
          });
        }
        
        // Cambiar estado
        if (Math.random() > 0.3) {
          const statuses = ['in_progress', 'resolved', 'closed'];
          const newStatus = statuses[Math.floor(Math.random() * statuses.length)];
          
          activities.push({
            id: `mock-status-${Date.now()}-${ticketId}`,
            userId,
            type: 'ticket_status_changed',
            targetId: ticketId,
            description: `Cambió el estado del ticket #${ticketId} a ${this.translateStatus(newStatus)}`,
            timestamp: new Date(Date.now() - Math.floor(Math.random() * 5) * 24 * 60 * 60 * 1000).toISOString(),
            metadata: {
              oldStatus: 'open',
              newStatus
            }
          });
        }
      }
      
      // Actualización de perfil
      if (Math.random() > 0.7) {
        activities.push({
          id: `mock-profile-${Date.now()}-${userId}`,
          userId,
          type: 'profile_updated',
          description: 'Actualizó su información de perfil',
          timestamp: new Date(Date.now() - Math.floor(Math.random() * 14) * 24 * 60 * 60 * 1000).toISOString()
        });
      }
      
      // Ordenar por fecha (más reciente primero)
      return activities.sort(
        (a, b) => new Date(b.timestamp).getTime() - new Date(a.timestamp).getTime()
      );
    },
    
    // Traducciones para estados de tickets
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
}); 
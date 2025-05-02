import { defineStore } from 'pinia';
import { ref } from 'vue';

export type NotificationType = 'success' | 'error' | 'warning' | 'info';

export interface Notification {
  id: string;
  message: string;
  type: NotificationType;
  duration: number;
  show: boolean;
}

export const useNotificationsStore = defineStore('notifications', () => {
  const notifications = ref<Notification[]>([]);

  /**
   * Añade una nueva notificación
   */
  const addNotification = (
    message: string,
    type: NotificationType = 'info',
    duration: number = 5000
  ) => {
    const id = `notification-${Date.now()}`;
    
    // Crear la notificación
    const notification: Notification = {
      id,
      message,
      type,
      duration,
      show: true
    };
    
    // Añadir la notificación al array
    notifications.value.push(notification);
    
    // Programar la eliminación automática
    if (duration > 0) {
      setTimeout(() => {
        removeNotification(id);
      }, duration);
    }
    
    return id;
  };

  /**
   * Elimina una notificación por su ID
   */
  const removeNotification = (id: string) => {
    const index = notifications.value.findIndex((notif: Notification) => notif.id === id);
    if (index !== -1) {
      // Primero ocultamos la notificación para la animación
      notifications.value[index].show = false;
      
      // Después de un pequeño retraso, la eliminamos del array
      setTimeout(() => {
        notifications.value = notifications.value.filter((notif: Notification) => notif.id !== id);
      }, 300);
    }
  };

  /**
   * Elimina todas las notificaciones
   */
  const clearAllNotifications = () => {
    // Primero ocultamos todas las notificaciones
    notifications.value.forEach((notif: Notification) => {
      notif.show = false;
    });
    
    // Después de un pequeño retraso, vaciamos el array
    setTimeout(() => {
      notifications.value = [];
    }, 300);
  };

  /**
   * Helpers para diferentes tipos de notificaciones
   */
  const success = (message: string, duration: number = 5000) => {
    return addNotification(message, 'success', duration);
  };

  const error = (message: string, duration: number = 5000) => {
    return addNotification(message, 'error', duration);
  };

  const warning = (message: string, duration: number = 5000) => {
    return addNotification(message, 'warning', duration);
  };

  const info = (message: string, duration: number = 5000) => {
    return addNotification(message, 'info', duration);
  };

  return {
    notifications,
    addNotification,
    removeNotification,
    clearAllNotifications,
    success,
    error,
    warning,
    info
  };
}); 
// app.ts - Store para configuración global de la aplicación
import { defineStore } from "pinia";

export interface AppState {
  darkMode: boolean;
  sidebarCollapsed: boolean;
  notifications: Notification[];
  maintenanceMode: boolean;
}

interface Notification {
  id: string;
  message: string;
  type: 'info' | 'success' | 'warning' | 'error';
  read: boolean;
  timestamp: Date;
}

export const useAppStore = defineStore({
  id: "app",
  
  state: (): AppState => ({
    darkMode: false,
    sidebarCollapsed: false,
    notifications: [],
    maintenanceMode: false
  }),
  
  actions: {
    toggleDarkMode() {
      this.darkMode = !this.darkMode;
      localStorage.setItem("darkMode", this.darkMode ? "true" : "false");
      
      // Aplicar modo oscuro al documento
      if (this.darkMode) {
        document.documentElement.classList.add('dark-theme');
      } else {
        document.documentElement.classList.remove('dark-theme');
      }
    },
    
    toggleSidebar() {
      this.sidebarCollapsed = !this.sidebarCollapsed;
      localStorage.setItem("sidebarCollapsed", this.sidebarCollapsed ? "true" : "false");
    },
    
    addNotification(notification: Omit<Notification, 'id' | 'timestamp' | 'read'>) {
      const id = `notification-${Date.now()}-${Math.random().toString(36).substr(2, 9)}`;
      
      this.notifications.unshift({
        ...notification,
        id,
        read: false,
        timestamp: new Date()
      });
      
      // Limitar a 30 notificaciones máximo
      if (this.notifications.length > 30) {
        this.notifications.pop();
      }
    },
    
    markNotificationAsRead(id: string) {
      const notification = this.notifications.find((n: Notification) => n.id === id);
      if (notification) {
        notification.read = true;
      }
    },
    
    clearNotifications() {
      this.notifications = [];
    },
    
    // Cargar preferencias guardadas al iniciar la app
    initApp() {
      // Cargar modo oscuro
      const savedDarkMode = localStorage.getItem("darkMode");
      if (savedDarkMode) {
        this.darkMode = savedDarkMode === "true";
        if (this.darkMode) {
          document.documentElement.classList.add('dark-theme');
        }
      }
      
      // Cargar estado de la barra lateral
      const savedSidebarState = localStorage.getItem("sidebarCollapsed");
      if (savedSidebarState) {
        this.sidebarCollapsed = savedSidebarState === "true";
      }
    }
  }
}); 
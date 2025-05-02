<template>
  <div class="notifications-container">
    <transition-group name="notification">
      <div 
        v-for="notification in notificationsStore.notifications" 
        :key="notification.id"
        class="notification-toast"
        :class="[notification.type, { 'show': notification.show }]"
      >
        <div class="notification-content">
          <i :class="getIconForType(notification.type)"></i>
          <span class="notification-message">{{ notification.message }}</span>
        </div>
        <button 
          class="close-notification"
          @click="notificationsStore.removeNotification(notification.id)"
        >
          &times;
        </button>
      </div>
    </transition-group>
  </div>
</template>

<script setup lang="ts">
import { useNotificationsStore } from '@/stores/notifications';
import type { NotificationType } from '@/stores/notifications';

const notificationsStore = useNotificationsStore();

const getIconForType = (type: NotificationType): string => {
  switch (type) {
    case 'success':
      return 'pi pi-check-circle';
    case 'error':
      return 'pi pi-times-circle';
    case 'warning':
      return 'pi pi-exclamation-triangle';
    case 'info':
    default:
      return 'pi pi-info-circle';
  }
};
</script>

<style scoped lang="scss">
.notifications-container {
  position: fixed;
  top: 20px;
  right: 20px;
  z-index: 9999;
  display: flex;
  flex-direction: column;
  gap: 10px;
  max-width: 350px;
  pointer-events: none;
}

.notification-toast {
  background-color: white;
  border-radius: 8px;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  padding: 16px;
  margin-bottom: 10px;
  display: flex;
  align-items: center;
  justify-content: space-between;
  pointer-events: auto;
  transition: all 0.3s ease;
  transform: translateX(120%);
  opacity: 0;
  
  &.show {
    transform: translateX(0);
    opacity: 1;
  }
  
  &.success {
    border-left: 4px solid #10b981;
    
    i {
      color: #10b981;
    }
  }
  
  &.error {
    border-left: 4px solid #ef4444;
    
    i {
      color: #ef4444;
    }
  }
  
  &.warning {
    border-left: 4px solid #f59e0b;
    
    i {
      color: #f59e0b;
    }
  }
  
  &.info {
    border-left: 4px solid #3b82f6;
    
    i {
      color: #3b82f6;
    }
  }
}

.notification-content {
  display: flex;
  align-items: center;
  gap: 12px;
  
  i {
    font-size: 1.25rem;
  }
}

.notification-message {
  font-size: 0.9rem;
  color: #374151;
}

.close-notification {
  background: none;
  border: none;
  color: #9ca3af;
  font-size: 1.25rem;
  cursor: pointer;
  margin-left: 8px;
  padding: 0;
  
  &:hover {
    color: #374151;
  }
}

/* Transiciones */
.notification-enter-active,
.notification-leave-active {
  transition: all 0.3s ease;
}

.notification-enter-from {
  transform: translateX(120%);
  opacity: 0;
}

.notification-leave-to {
  transform: translateX(120%);
  opacity: 0;
}
</style> 
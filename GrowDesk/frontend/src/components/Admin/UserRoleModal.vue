<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Cambiar rol de usuario</h2>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>
      
      <div class="modal-body">
        <div class="user-info">
          <p><strong>Usuario:</strong> {{ user?.firstName }} {{ user?.lastName }}</p>
          <p><strong>Email:</strong> {{ user?.email }}</p>
          <p><strong>Rol actual:</strong> <span :class="['role-badge', user?.role]">{{ translateRole(user?.role || '') }}</span></p>
        </div>
        
        <div class="role-selection">
          <h3>Seleccionar nuevo rol</h3>
          <div class="role-options">
            <div 
              v-for="role in availableRoles" 
              :key="role.value"
              class="role-option"
              :class="{ 'selected': selectedRole === role.value }"
              @click="selectedRole = role.value"
            >
              <div class="role-header">
                <span :class="['role-badge', role.value]">{{ role.label }}</span>
              </div>
              <p class="role-description">{{ role.description }}</p>
            </div>
          </div>
        </div>
        
        <div v-if="error" class="error-message">
          {{ error }}
        </div>
      </div>
      
      <div class="modal-footer">
        <button class="btn btn-secondary" @click="$emit('close')">Cancelar</button>
        <button 
          class="btn btn-primary" 
          @click="changeRole" 
          :disabled="!hasRoleChanged || isSubmitting"
        >
          {{ isSubmitting ? 'Guardando...' : 'Guardar Cambios' }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed } from 'vue';
import { useUsersStore } from '@/stores/users';
import type { User } from '@/stores/users';

// Props y emits
const props = defineProps<{
  user: User | null
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'role-changed'): void;
}>();

// Store
const userStore = useUsersStore();

// Estado
const selectedRole = ref(props.user?.role || 'employee');
const isSubmitting = ref(false);
const error = ref('');

// Roles disponibles
const availableRoles = [
  {
    value: 'admin',
    label: 'Administrador',
    description: 'Acceso completo a todas las funcionalidades del sistema, incluyendo administración de usuarios y configuración.'
  },
  {
    value: 'assistant',
    label: 'Asistente',
    description: 'Puede gestionar tickets, realizar asignaciones y acceder a informes, pero con permisos limitados de administración.'
  },
  {
    value: 'employee',
    label: 'Empleado',
    description: 'Acceso básico para crear tickets, revisar su estado y comunicarse con el equipo de soporte.'
  }
];

// Comprobar si el rol ha cambiado
const hasRoleChanged = computed(() => {
  return selectedRole.value !== props.user?.role;
});

// Traducir rol
const translateRole = (role: string) => {
  const roles: Record<string, string> = {
    'admin': 'Administrador',
    'assistant': 'Asistente',
    'employee': 'Empleado'
  };
  return roles[role] || role;
};

// Cambiar rol del usuario
const changeRole = async () => {
  if (!props.user || !hasRoleChanged.value) return;
  
  isSubmitting.value = true;
  error.value = '';
  
  try {
    await userStore.changeUserRole(props.user.id, selectedRole.value as 'admin' | 'assistant' | 'employee');
    console.log(`Rol de usuario cambiado a: ${translateRole(selectedRole.value)}`);
    
    // Emitir evento para actualizar la lista
    setTimeout(() => {
      userStore.fetchUsers();
      emit('role-changed');
    }, 300);
  } catch (err) {
    console.error('Error changing user role:', err);
    error.value = 'Error al cambiar el rol del usuario';
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style lang="scss" scoped>
.modal-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background: white;
  border-radius: 8px;
  width: 100%;
  max-width: 600px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  max-height: 90vh;
  overflow-y: auto;
  display: flex;
  flex-direction: column;
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
  
  h2 {
    margin: 0;
    font-size: 1.5rem;
    color: #333;
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #666;
    
    &:hover {
      color: #333;
    }
  }
}

.modal-body {
  padding: 1.5rem;
  flex-grow: 1;
  
  .user-info {
    background-color: #f8f9fa;
    padding: 1rem;
    border-radius: 8px;
    margin-bottom: 1.5rem;
    
    p {
      margin: 0.5rem 0;
    }
    
    strong {
      font-weight: 600;
      color: #555;
    }
  }
  
  .role-selection {
    h3 {
      margin: 0 0 1rem 0;
      font-size: 1.2rem;
      color: #333;
    }
    
    .role-options {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(250px, 1fr));
      gap: 1rem;
      
      .role-option {
        border: 2px solid #e0e0e0;
        border-radius: 8px;
        padding: 1rem;
        cursor: pointer;
        transition: all 0.2s;
        
        &:hover {
          border-color: #bbdefb;
          background-color: #f5f9ff;
        }
        
        &.selected {
          border-color: #1976d2;
          background-color: #e3f2fd;
        }
        
        .role-header {
          margin-bottom: 0.75rem;
        }
        
        .role-description {
          margin: 0;
          font-size: 0.9rem;
          color: #666;
        }
      }
    }
  }
}

.modal-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
}

.role-badge {
  display: inline-block;
  padding: 0.25rem 0.5rem;
  border-radius: 4px;
  font-size: 0.875rem;
  font-weight: 500;
  
  &.admin {
    background-color: #e3f2fd;
    color: #1976d2;
  }
  
  &.assistant {
    background-color: #fff3e0;
    color: #f57c00;
  }
  
  &.employee {
    background-color: #e8f5e9;
    color: #388e3c;
  }
}

.error-message {
  margin-top: 1rem;
  color: #d32f2f;
  font-size: 0.9rem;
  background-color: #ffebee;
  padding: 0.5rem;
  border-radius: 4px;
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  
  &.btn-primary {
    background-color: #1976d2;
    color: white;
    
    &:hover {
      background-color: #1565c0;
    }
    
    &:disabled {
      background-color: #64b5f6;
      cursor: not-allowed;
    }
  }
  
  &.btn-secondary {
    background-color: #f5f5f5;
    color: #333;
    
    &:hover {
      background-color: #e0e0e0;
    }
  }
}
</style> 
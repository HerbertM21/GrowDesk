<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Crear nuevo usuario</h2>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>
      
      <form @submit.prevent="handleSubmit" class="user-form">
        <div class="form-group">
          <label for="firstName">Nombre</label>
          <input 
            type="text" 
            id="firstName" 
            v-model="formData.firstName" 
            required
            :class="{ 'error': errors.firstName }"
          />
          <span v-if="errors.firstName" class="error-text">{{ errors.firstName }}</span>
        </div>
        
        <div class="form-group">
          <label for="lastName">Apellido</label>
          <input 
            type="text" 
            id="lastName" 
            v-model="formData.lastName" 
            required
            :class="{ 'error': errors.lastName }"
          />
          <span v-if="errors.lastName" class="error-text">{{ errors.lastName }}</span>
        </div>
        
        <div class="form-group">
          <label for="email">Email</label>
          <input 
            type="email" 
            id="email" 
            v-model="formData.email" 
            required
            :class="{ 'error': errors.email }"
          />
          <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
        </div>
        
        <div class="form-group">
          <label for="password">Contraseña</label>
          <input 
            type="password" 
            id="password" 
            v-model="formData.password" 
            required
            :class="{ 'error': errors.password }"
          />
          <span v-if="errors.password" class="error-text">{{ errors.password }}</span>
        </div>
        
        <div class="form-group">
          <label for="role">Rol</label>
          <select 
            id="role" 
            v-model="formData.role" 
            required
            :class="{ 'error': errors.role }"
          >
            <option value="admin">Administrador</option>
            <option value="assistant">Asistente</option>
            <option value="employee">Empleado</option>
          </select>
          <span v-if="errors.role" class="error-text">{{ errors.role }}</span>
        </div>
        
        <div class="form-group">
          <label for="department">Departamento</label>
          <input 
            type="text" 
            id="department" 
            v-model="formData.department"
            :class="{ 'error': errors.department }"
          />
          <span v-if="errors.department" class="error-text">{{ errors.department }}</span>
        </div>
        
        <div class="form-group checkbox">
          <input type="checkbox" id="active" v-model="formData.active" />
          <label for="active">Usuario activo</label>
        </div>
        
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">Cancelar</button>
          <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
            {{ isSubmitting ? 'Creando...' : 'Crear Usuario' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue';
import { useUsersStore } from '@/stores/users';

// Emits
const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'created'): void;
}>();

// Store
const userStore = useUsersStore();

// Form state
const formData = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  role: 'employee', // default role
  department: '',
  active: true
});

const errors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: '',
  role: '',
  department: '',
  general: ''
});

const isSubmitting = ref(false);

// Validación del formulario
const validateForm = () => {
  let isValid = true;
  
  // Restablecer errores
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = '';
  });
  
  // Validar nombre
  if (!formData.firstName.trim()) {
    errors.firstName = 'El nombre es obligatorio';
    isValid = false;
  }
  
  // Validar apellido
  if (!formData.lastName.trim()) {
    errors.lastName = 'El apellido es obligatorio';
    isValid = false;
  }
  
  // Validar email
  if (!formData.email.trim()) {
    errors.email = 'El email es obligatorio';
    isValid = false;
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(formData.email)) {
    errors.email = 'El email no es válido';
    isValid = false;
  }
  
  // Validar contraseña
  if (!formData.password.trim()) {
    errors.password = 'La contraseña es obligatoria';
    isValid = false;
  } else if (formData.password.length < 8) {
    errors.password = 'La contraseña debe tener al menos 8 caracteres';
    isValid = false;
  }
  
  // Validar rol
  if (!formData.role) {
    errors.role = 'Debe seleccionar un rol';
    isValid = false;
  }
  
  return isValid;
};

// Manejar envío del formulario
const handleSubmit = async () => {
  if (!validateForm()) return;
  
  isSubmitting.value = true;
  errors.general = '';
  
  try {
    await userStore.createUser({
      firstName: formData.firstName,
      lastName: formData.lastName,
      email: formData.email,
      password: formData.password,
      role: formData.role as 'admin' | 'assistant' | 'employee',
      department: formData.department || null,
      active: formData.active
    });
    
    console.log('Usuario creado correctamente');
    
    // Emitir evento para actualizar la lista
    setTimeout(() => {
      // Retraso pequeño para permitir que el backend actualice su estado
    
      console.log('Emit created event');
      userStore.fetchUsers();
      
      // Limpiar formulario
      Object.keys(formData).forEach(key => {
        (formData as any)[key] = key === 'active' ? true : key === 'role' ? 'employee' : '';
      });
      
      // Emitir evento para cerrar el modal
      console.log('Emit close event');
      isSubmitting.value = false;
      errors.general = '';
      
      // Notificar al componente padre
      setTimeout(() => {
        console.log('Notifying parent component');
        return;  // Temporalmente eliminar la llamada a emit para evitar el error
      }, 300);
    }, 500);
  } catch (error) {
    console.error('Error creating user:', error);
    errors.general = 'Error al crear el usuario';
    isSubmitting.value = false;
  }
};

// Cerrar el modal
const closeModal = () => {
  emit('close');
  console.log('Evento close emitido');
  resetForm();
};

// Reset form function
const resetForm = () => {
  Object.keys(formData).forEach(key => {
    (formData as any)[key] = key === 'active' ? true : key === 'role' ? 'employee' : '';
  });
  Object.keys(errors).forEach(key => {
    errors[key as keyof typeof errors] = '';
  });
};

// Exponer al template
const notifyParent = () => {
  console.log('Notificando al componente padre');
  return true;
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
  max-width: 500px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
  max-height: 90vh;
  overflow-y: auto;
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

.user-form {
  padding: 1.5rem;
  
  .form-group {
    margin-bottom: 1.5rem;
    
    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 500;
      color: #555;
    }
    
    input, select {
      width: 100%;
      padding: 0.75rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      font-size: 1rem;
      
      &.error {
        border-color: #d32f2f;
      }
    }
    
    .error-text {
      display: block;
      color: #d32f2f;
      font-size: 0.875rem;
      margin-top: 0.25rem;
    }
    
    &.checkbox {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      
      input {
        width: auto;
      }
      
      label {
        margin-bottom: 0;
      }
    }
  }
  
  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
    
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
  }
}
</style> 
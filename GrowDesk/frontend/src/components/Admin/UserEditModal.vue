<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>Editar usuario</h2>
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
            {{ isSubmitting ? 'Guardando...' : 'Guardar Cambios' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue';
import { useUsersStore } from '@/stores/users';
import type { User } from '@/stores/users';

// Props y emits
const props = defineProps<{
  user: User | null
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'updated'): void;
}>();

// Store
const userStore = useUsersStore();

// Form state
const formData = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  active: true
});

const errors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  general: ''
});

const isSubmitting = ref(false);

// Cargar datos del usuario cuando se monta el componente
onMounted(() => {
  if (props.user) {
    formData.firstName = props.user.firstName;
    formData.lastName = props.user.lastName;
    formData.email = props.user.email;
    formData.department = props.user.department || '';
    formData.active = props.user.active;
  }
});

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
  
  return isValid;
};

// Manejar envío del formulario
const handleSubmit = async () => {
  if (!props.user || !validateForm()) return;
  
  isSubmitting.value = true;
  errors.general = '';
  
  try {
    // Crear un objeto con los datos actualizados
    const userData = {
      firstName: formData.firstName.trim(),
      lastName: formData.lastName.trim(),
      email: formData.email.trim(),
      department: formData.department.trim() || null,
      active: formData.active
    };
    
    console.log('---------------------------------------------');
    console.log('Datos del usuario original:', props.user);
    console.log('Datos que se enviarán al servidor:', userData);
    console.log('ID del usuario a actualizar:', props.user.id);
    
    // Intentar actualizar usando el método del store
    const updatedUser = await userStore.updateUser(props.user.id, userData);
    
    console.log('Respuesta de la actualización:', updatedUser);
    console.log('Usuario actualizado exitosamente en el store');
    
    // Emitir evento para actualizar la lista
    console.log('Emitiendo evento updated');
    emit('updated');
    
    // Cerrar el modal después de guardar
    setTimeout(() => {
      emit('close');
    }, 500);
  } catch (err: any) {
    console.error('Error updating user:', err);
    errors.general = `Error al actualizar el usuario: ${err.message}`;
  } finally {
    isSubmitting.value = false;
    console.log('---------------------------------------------');
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
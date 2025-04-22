<template>
  <div class="modal-backdrop" @click.self="$emit('close')">
    <div class="modal-content">
      <div class="modal-header">
        <h2>
          <i class="pi pi-user-edit"></i>
          Modificar Perfil de Usuario
        </h2>
        <button class="close-btn" @click="$emit('close')">&times;</button>
      </div>
      
      <div v-if="user" class="user-summary">
        <div class="avatar" :style="{ backgroundColor: getAvatarColor() }">
          {{ getInitials() }}
        </div>
        <div class="user-info">
          <h3>{{ user.firstName }} {{ user.lastName }}</h3>
          <span :class="['role-badge', user.role]">{{ translateRole(user.role) }}</span>
        </div>
      </div>
      
      <form @submit.prevent="handleSubmit" class="profile-form">
        <div class="form-section">
          <h3><i class="pi pi-user"></i> Información Básica</h3>
          
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
            <div class="input-with-icon">
              <i class="pi pi-envelope"></i>
              <input 
                type="email" 
                id="email" 
                v-model="formData.email" 
                required
                :class="{ 'error': errors.email }"
              />
            </div>
            <span v-if="errors.email" class="error-text">{{ errors.email }}</span>
          </div>
        </div>
        
        <div class="form-section">
          <h3><i class="pi pi-briefcase"></i> Información Profesional</h3>
          
          <div class="form-group">
            <label for="department">Departamento</label>
            <div class="input-with-icon">
              <i class="pi pi-building"></i>
              <input 
                type="text" 
                id="department" 
                v-model="formData.department"
                :class="{ 'error': errors.department }"
              />
            </div>
            <span v-if="errors.department" class="error-text">{{ errors.department }}</span>
          </div>
          
          <div class="form-group">
            <label for="position">Cargo</label>
            <div class="input-with-icon">
              <i class="pi pi-id-card"></i>
              <input 
                type="text" 
                id="position" 
                v-model="formData.position"
                :class="{ 'error': errors.position }"
              />
            </div>
            <span v-if="errors.position" class="error-text">{{ errors.position }}</span>
          </div>
          
          <div class="form-group">
            <label for="phone">Teléfono</label>
            <div class="input-with-icon">
              <i class="pi pi-phone"></i>
              <input 
                type="tel" 
                id="phone" 
                v-model="formData.phone"
                placeholder="+34 600 000 000"
                :class="{ 'error': errors.phone }"
              />
            </div>
            <span v-if="errors.phone" class="error-text">{{ errors.phone }}</span>
          </div>
        </div>
        
        <div class="form-section">
          <h3><i class="pi pi-cog"></i> Configuración</h3>
          
          <div class="form-group checkbox">
            <input type="checkbox" id="active" v-model="formData.active" />
            <label for="active">Usuario activo</label>
          </div>
          
          <div class="form-group">
            <label for="language">Idioma preferido</label>
            <div class="input-with-icon">
              <i class="pi pi-globe"></i>
              <select id="language" v-model="formData.language">
                <option value="es">Español</option>
                <option value="en">Inglés</option>
                <option value="fr">Francés</option>
                <option value="de">Alemán</option>
              </select>
            </div>
          </div>
        </div>
        
        <div class="form-section" v-if="isAdminUser">
          <h3><i class="pi pi-shield"></i> Solo Administradores</h3>
          <div class="admin-note">
            <p>Como administrador, puedes modificar la información de todos los usuarios.</p>
          </div>
          
          <div class="form-group" v-if="showRoleField">
            <label for="role">Rol del usuario</label>
            <div class="input-with-icon">
              <i class="pi pi-users"></i>
              <select id="role" v-model="formData.role">
                <option value="admin">Administrador</option>
                <option value="assistant">Asistente</option>
                <option value="employee">Empleado</option>
              </select>
            </div>
          </div>
        </div>
        
        <div v-if="errors.general" class="general-error">
          {{ errors.general }}
        </div>
        
        <div class="form-actions">
          <button type="button" class="btn btn-secondary" @click="$emit('close')">
            <i class="pi pi-times"></i> Cancelar
          </button>
          <button type="submit" class="btn btn-primary" :disabled="isSubmitting">
            <i class="pi" :class="isSubmitting ? 'pi-spin pi-spinner' : 'pi-save'"></i>
            {{ isSubmitting ? 'Guardando...' : 'Guardar Cambios' }}
          </button>
        </div>
      </form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue';
import { useUsersStore } from '@/stores/users';
import { useAuthStore } from '@/stores/auth';
import type { User } from '@/stores/users';

// Props y emits
const props = defineProps<{
  user: User | null
}>();

const emit = defineEmits<{
  (e: 'close'): void;
  (e: 'updated'): void;
}>();

// Stores
const userStore = useUsersStore();
const authStore = useAuthStore();

// Estado del formulario
const formData = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  position: '',
  phone: '',
  language: 'es',
  active: true,
  role: 'employee'
});

const errors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  position: '',
  phone: '',
  general: ''
});

const isSubmitting = ref(false);

// Computados
const currentUser = computed(() => authStore.user);

const isAdminUser = computed(() => {
  return currentUser.value?.role === 'admin';
});

const showRoleField = computed(() => {
  // Solo los administradores pueden cambiar roles y no pueden cambiar su propio rol
  return isAdminUser.value && currentUser.value?.id !== props.user?.id;
});

// Métodos para información de usuario
const getInitials = () => {
  if (!props.user) return '';
  return (props.user.firstName.charAt(0) + props.user.lastName.charAt(0)).toUpperCase();
};

const getAvatarColor = () => {
  if (!props.user) return 'var(--primary-color)';
  
  // Colores disponibles
  const colors = [
    'var(--primary-color, #1976d2)',
    'var(--success-color, #388e3c)',
    'var(--warning-color, #f57c00)',
    'var(--info-color, #0288d1)',
    'var(--purple-color, #7b1fa2)'
  ];
  
  // Usar el ID para seleccionar un color
  const colorIndex = parseInt(props.user.id) % colors.length;
  return colors[colorIndex];
};

const translateRole = (role: string) => {
  const roles: Record<string, string> = {
    'admin': 'Administrador',
    'assistant': 'Asistente',
    'employee': 'Empleado'
  };
  return roles[role] || role;
};

// Cargar datos del usuario cuando se monta el componente
onMounted(() => {
  if (props.user) {
    formData.firstName = props.user.firstName;
    formData.lastName = props.user.lastName;
    formData.email = props.user.email;
    formData.department = props.user.department || '';
    formData.active = props.user.active;
    formData.role = props.user.role;
    
    // Cargar campos adicionales si existen
    if ('position' in props.user) {
      formData.position = (props.user as any).position || '';
    }
    if ('phone' in props.user) {
      formData.phone = (props.user as any).phone || '';
    }
    if ('language' in props.user) {
      formData.language = (props.user as any).language || 'es';
    }
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
  
  // Validar teléfono (opcional)
  if (formData.phone && !/^(\+\d{1,3}\s?)?\d{9,10}$/.test(formData.phone.replace(/\s/g, ''))) {
    errors.phone = 'El formato del teléfono no es válido';
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
      active: formData.active,
      position: formData.position.trim() || null,
      phone: formData.phone.trim() || null,
      language: formData.language
    };
    
    // Solo incluir el rol si el usuario es admin y está modificando a otro usuario
    if (isAdminUser.value && currentUser.value?.id !== props.user.id) {
      (userData as any).role = formData.role;
    }
    
    console.log('Datos del usuario original:', props.user);
    console.log('Datos que se enviarán al servidor:', userData);
    
    // Actualizar el usuario
    await userStore.updateUser(props.user.id, userData);
    
    // Emitir evento para actualizar la lista
    emit('updated');
  } catch (error) {
    console.error('Error al actualizar perfil:', error);
    errors.general = 'Ha ocurrido un error al guardar los cambios';
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
  align-items: center;
  justify-content: center;
  z-index: 1000;
  backdrop-filter: blur(3px);
}

.modal-content {
  background-color: var(--bg-secondary, white);
  border-radius: 8px;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
  position: relative;
  border: 1px solid var(--border-color, #eee);
  animation: modalFadeIn 0.3s ease;
}

@keyframes modalFadeIn {
  from {
    opacity: 0;
    transform: translateY(-20px);
  }
  to {
    opacity: 1;
    transform: translateY(0);
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.5rem;
  border-bottom: 1px solid var(--border-color, #eee);
  background-color: var(--header-bg, white);
  
  h2 {
    margin: 0;
    color: var(--text-primary, #333);
    font-size: 1.5rem;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    
    i {
      color: var(--primary-color);
    }
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text-secondary, #777);
    transition: color 0.2s;
    
    &:hover {
      color: var(--text-primary, #333);
    }
  }
}

.user-summary {
  display: flex;
  align-items: center;
  gap: 1rem;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  background-color: var(--bg-tertiary, #f9f9f9);
  
  .avatar {
    width: 50px;
    height: 50px;
    border-radius: 50%;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-weight: bold;
    font-size: 1.2rem;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  }
  
  .user-info {
    display: flex;
    flex-direction: column;
    
    h3 {
      margin: 0 0 0.25rem 0;
      font-size: 1.1rem;
      color: var(--text-primary);
    }
    
    .role-badge {
      display: inline-block;
      padding: 0.25rem 0.5rem;
      border-radius: 4px;
      font-size: 0.75rem;
      align-self: flex-start;
      
      &.admin {
        background-color: var(--admin-bg, #e3f2fd);
        color: var(--admin-color, #1976d2);
      }
      
      &.assistant {
        background-color: var(--assistant-bg, #fff3e0);
        color: var(--assistant-color, #f57c00);
      }
      
      &.employee {
        background-color: var(--employee-bg, #e8f5e9);
        color: var(--employee-color, #388e3c);
      }
    }
  }
}

.profile-form {
  padding: 0 1.5rem 1.5rem;
  
  .form-section {
    margin-bottom: 1.5rem;
    
    h3 {
      font-size: 1.1rem;
      color: var(--text-secondary, #555);
      margin: 1.5rem 0 1rem;
      padding-bottom: 0.5rem;
      border-bottom: 1px solid var(--border-color, #eee);
      display: flex;
      align-items: center;
      gap: 0.5rem;
      
      i {
        color: var(--primary-color);
      }
    }
    
    .admin-note {
      background-color: var(--primary-color-light, #e7f5ff);
      padding: 0.75rem;
      border-radius: 4px;
      margin-bottom: 1rem;
      
      p {
        margin: 0;
        color: var(--primary-color, #1976d2);
        font-size: 0.9rem;
      }
    }
  }
  
  .form-group {
    margin-bottom: 1rem;
    
    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 500;
      color: var(--text-secondary, #555);
    }
    
    .input-with-icon {
      position: relative;
      
      i {
        position: absolute;
        left: 10px;
        top: 50%;
        transform: translateY(-50%);
        color: var(--text-secondary);
      }
      
      input, select {
        padding-left: 2.5rem !important;
      }
    }
    
    input[type="text"],
    input[type="email"],
    input[type="tel"],
    select {
      width: 100%;
      padding: 0.75rem;
      border: 1px solid var(--border-color, #ddd);
      border-radius: 4px;
      font-size: 1rem;
      background-color: var(--bg-tertiary, white);
      color: var(--text-primary, #333);
      
      &:focus {
        outline: none;
        border-color: var(--primary-color, #1976d2);
        box-shadow: 0 0 0 2px var(--primary-color-light, rgba(25, 118, 210, 0.2));
      }
      
      &.error {
        border-color: var(--error-color, #dc3545);
      }
    }
    
    .error-text {
      color: var(--error-color, #dc3545);
      font-size: 0.875rem;
      margin-top: 0.5rem;
      display: block;
    }
    
    &.checkbox {
      display: flex;
      align-items: center;
      
      input {
        margin-right: 0.5rem;
        width: 18px;
        height: 18px;
        accent-color: var(--primary-color, #1976d2);
      }
      
      label {
        margin-bottom: 0;
        cursor: pointer;
      }
    }
  }
  
  .general-error {
    background-color: var(--error-bg, #f8d7da);
    color: var(--error-color, #dc3545);
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1.5rem;
    border: 1px solid var(--error-border, #f5c6cb);
  }
  
  .form-actions {
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
    margin-top: 2rem;
    
    button {
      padding: 0.75rem 1.5rem;
      border-radius: 4px;
      font-weight: 500;
      cursor: pointer;
      transition: all 0.2s;
      display: flex;
      align-items: center;
      gap: 0.5rem;
      
      &.btn-secondary {
        background-color: var(--bg-tertiary, #f5f5f5);
        color: var(--text-primary, #333);
        border: 1px solid var(--border-color, #ddd);
        
        &:hover {
          background-color: var(--bg-hover, #e0e0e0);
        }
        
        &:active {
          transform: translateY(1px);
        }
      }
      
      &.btn-primary {
        background-color: var(--primary-color, #1976d2);
        color: white;
        border: none;
        
        &:hover:not(:disabled) {
          background-color: var(--primary-color-dark, #1565c0);
        }
        
        &:active:not(:disabled) {
          transform: translateY(1px);
        }
        
        &:disabled {
          background-color: var(--primary-color-light, #90caf9);
          cursor: not-allowed;
        }
      }
    }
  }
}

// Responsive styles
@media (max-width: 768px) {
  .modal-content {
    width: 95%;
    max-height: 95vh;
  }
  
  .modal-header h2 {
    font-size: 1.3rem;
  }
  
  .form-actions {
    flex-direction: column;
    
    button {
      width: 100%;
      
      &.btn-secondary {
        order: 2;
      }
      
      &.btn-primary {
        order: 1;
      }
    }
  }
}
</style> 
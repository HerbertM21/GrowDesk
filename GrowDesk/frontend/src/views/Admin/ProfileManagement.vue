<template>
  <div class="profile-management">
    <div class="profile-header">
      <h1 class="profile-title">Gestión de Perfiles de Usuarios</h1>
      <p class="profile-subtitle">Administre los perfiles de los usuarios del sistema</p>
    </div>
    
    <div class="profile-content">
      <!-- Vista de selección de usuario -->
      <div class="user-selection">
        <label for="user-select">Seleccionar Usuario:</label>
        <div class="select-container">
          <select id="user-select" v-model="selectedUserId" class="form-control" @change="loadUserProfile">
            <option value="" disabled>Seleccione un usuario</option>
            <option v-for="user in users" :key="user.id" :value="user.id">
              {{ user.firstName }} {{ user.lastName }} ({{ translateRole(user.role) }})
            </option>
          </select>
          <i class="pi pi-chevron-down"></i>
        </div>
      </div>
      
      <!-- Sección principal con tarjetas para diferentes aspectos del perfil del usuario seleccionado -->
      <div class="profile-grid" v-if="selectedUserId">
        <!-- Información personal -->
        <div class="profile-card">
          <div class="profile-card-header">
            <i class="pi pi-user"></i>
            <h2>Información Personal</h2>
          </div>
          <div class="profile-card-body">
            <form @submit.prevent="updatePersonalInfo">
              <div class="profile-form-group">
                <label for="firstName">Nombre</label>
                <input 
                  id="firstName" 
                  v-model="personalInfo.firstName" 
                  type="text" 
                  class="form-control"
                  required
                  :class="{ 'is-invalid': errors.firstName }"
                />
                <div v-if="errors.firstName" class="invalid-feedback">
                  {{ errors.firstName }}
                </div>
              </div>
              
              <div class="profile-form-group">
                <label for="lastName">Apellido</label>
                <input 
                  id="lastName" 
                  v-model="personalInfo.lastName" 
                  type="text" 
                  class="form-control"
                  required
                  :class="{ 'is-invalid': errors.lastName }"
                />
                <div v-if="errors.lastName" class="invalid-feedback">
                  {{ errors.lastName }}
                </div>
              </div>
              
              <div class="profile-form-group">
                <label for="email">Email</label>
                <input 
                  id="email" 
                  v-model="personalInfo.email" 
                  type="email" 
                  class="form-control"
                  required
                  :class="{ 'is-invalid': errors.email }"
                  disabled
                />
                <div v-if="errors.email" class="invalid-feedback">
                  {{ errors.email }}
                </div>
                <small class="form-text text-muted">
                  El email no se puede cambiar por motivos de seguridad.
                </small>
              </div>
              
              <div class="profile-form-group">
                <label for="department">Departamento</label>
                <select 
                  id="department" 
                  v-model="personalInfo.department" 
                  class="form-control"
                  :class="{ 'is-invalid': errors.department }"
                >
                  <option value="" disabled>Seleccione un departamento</option>
                  <option value="IT">IT</option>
                  <option value="RRHH">RRHH</option>
                  <option value="Ventas">Ventas</option>
                  <option value="Marketing">Marketing</option>
                  <option value="Soporte">Soporte</option>
                </select>
                <div v-if="errors.department" class="invalid-feedback">
                  {{ errors.department }}
                </div>
              </div>
              
              <div class="profile-card-actions">
                <button 
                  type="submit" 
                  class="btn btn-primary"
                  :disabled="isSubmitting"
                >
                  <i class="pi pi-save"></i> Guardar Cambios
                </button>
              </div>
            </form>
          </div>
        </div>
        
        <!-- Cambio de rol y estado -->
        <div class="profile-card">
          <div class="profile-card-header">
            <i class="pi pi-users"></i>
            <h2>Rol y Estado</h2>
          </div>
          <div class="profile-card-body">
            <form @submit.prevent="updateUserRole">
              <div class="profile-form-group">
                <label for="role">Rol</label>
                <select 
                  id="role" 
                  v-model="userRole.role" 
                  class="form-control"
                  :class="{ 'is-invalid': errors.role }"
                >
                  <option value="" disabled>Seleccione un rol</option>
                  <option value="admin">Administrador</option>
                  <option value="assistant">Asistente</option>
                  <option value="employee">Empleado</option>
                </select>
                <div v-if="errors.role" class="invalid-feedback">
                  {{ errors.role }}
                </div>
              </div>
              
              <div class="profile-form-group">
                <label>Estado de la cuenta</label>
                <div class="toggle-switch-container">
                  <label class="toggle-switch">
                    <input 
                      type="checkbox" 
                      v-model="userRole.active"
                    >
                    <span class="slider round"></span>
                  </label>
                  <span class="toggle-label">
                    {{ userRole.active ? 'Activo' : 'Inactivo' }}
                  </span>
                </div>
              </div>
              
              <div class="profile-card-actions">
                <button 
                  type="submit" 
                  class="btn btn-primary"
                  :disabled="isSubmitting"
                >
                  <i class="pi pi-users"></i> Actualizar Rol
                </button>
              </div>
            </form>
          </div>
        </div>
        
        <!-- Resetear contraseña -->
        <div class="profile-card">
          <div class="profile-card-header">
            <i class="pi pi-lock"></i>
            <h2>Resetear Contraseña</h2>
          </div>
          <div class="profile-card-body">
            <form @submit.prevent="resetPassword">
              <div class="profile-form-group">
                <label for="newPassword">Nueva Contraseña</label>
                <div class="password-input-container">
                  <input 
                    id="newPassword" 
                    v-model="passwordData.newPassword" 
                    :type="showNewPassword ? 'text' : 'password'" 
                    class="form-control"
                    required
                    :class="{ 'is-invalid': errors.newPassword }"
                  />
                  <button 
                    type="button" 
                    class="password-toggle" 
                    @click="showNewPassword = !showNewPassword"
                  >
                    <i :class="showNewPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
                  </button>
                </div>
                <div v-if="errors.newPassword" class="invalid-feedback">
                  {{ errors.newPassword }}
                </div>
                <div class="password-strength" v-if="passwordData.newPassword">
                  <div class="strength-meter">
                    <div 
                      class="strength-meter-fill" 
                      :style="{ width: passwordStrength.percent + '%' }"
                      :class="passwordStrength.class"
                    ></div>
                  </div>
                  <span class="strength-text" :class="passwordStrength.class">
                    {{ passwordStrength.label }}
                  </span>
                </div>
              </div>
              
              <div class="profile-form-group">
                <label for="confirmPassword">Confirmar Contraseña</label>
                <div class="password-input-container">
                  <input 
                    id="confirmPassword" 
                    v-model="passwordData.confirmPassword" 
                    :type="showConfirmPassword ? 'text' : 'password'" 
                    class="form-control"
                    required
                    :class="{ 'is-invalid': errors.confirmPassword }"
                  />
                  <button 
                    type="button" 
                    class="password-toggle" 
                    @click="showConfirmPassword = !showConfirmPassword"
                  >
                    <i :class="showConfirmPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
                  </button>
                </div>
                <div v-if="errors.confirmPassword" class="invalid-feedback">
                  {{ errors.confirmPassword }}
                </div>
              </div>
              
              <div class="profile-card-actions">
                <button 
                  type="submit" 
                  class="btn btn-primary"
                  :disabled="isSubmitting"
                >
                  <i class="pi pi-lock"></i> Resetear Contraseña
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
      
      <div v-else class="empty-state">
        <i class="pi pi-user-edit" style="font-size: 3rem"></i>
        <p>Seleccione un usuario para gestionar su perfil</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useUsersStore } from '@/stores/users';
import type { User } from '@/stores/users';

const usersStore = useUsersStore();

// Lista de usuarios para seleccionar
const users = computed(() => usersStore.users);
const selectedUserId = ref('');

// Estado de carga y errores
const isSubmitting = ref(false);
const errors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  role: '',
  newPassword: '',
  confirmPassword: ''
});

// Información personal
const personalInfo = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: ''
});

// Datos de rol
const userRole = reactive({
  role: '',
  active: true
});

// Datos de contraseña
const passwordData = reactive({
  newPassword: '',
  confirmPassword: ''
});

// Mostrar/ocultar contraseñas
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);

// Cargar datos iniciales
onMounted(async () => {
  try {
    // Cargar la lista de usuarios
    await usersStore.fetchUsers();
  } catch (error) {
    console.error('Error al cargar la lista de usuarios:', error);
  }
});

// Cargar perfil de usuario cuando se selecciona uno
const loadUserProfile = async () => {
  if (!selectedUserId.value) return;
  
  try {
    const user = await usersStore.getUser(selectedUserId.value);
    
    // Llenar información personal
    personalInfo.firstName = user.firstName || '';
    personalInfo.lastName = user.lastName || '';
    personalInfo.email = user.email || '';
    personalInfo.department = user.department || '';
    
    // Llenar datos de rol
    userRole.role = user.role || '';
    userRole.active = user.active || false;
    
    // Resetear datos de contraseña y errores
    passwordData.newPassword = '';
    passwordData.confirmPassword = '';
    resetErrors();
  } catch (error) {
    console.error('Error al cargar datos del usuario:', error);
  }
};

// Watch para resetear datos cuando cambia el usuario seleccionado
watch(selectedUserId, () => {
  loadUserProfile();
});

// Resetear errores
const resetErrors = () => {
  errors.firstName = '';
  errors.lastName = '';
  errors.email = '';
  errors.department = '';
  errors.role = '';
  errors.newPassword = '';
  errors.confirmPassword = '';
};

// Actualizar información personal
const updatePersonalInfo = async () => {
  // Reset de errores
  errors.firstName = '';
  errors.lastName = '';
  errors.email = '';
  errors.department = '';
  
  // Validación básica
  let isValid = true;
  
  if (!personalInfo.firstName.trim()) {
    errors.firstName = 'El nombre es obligatorio';
    isValid = false;
  }
  
  if (!personalInfo.lastName.trim()) {
    errors.lastName = 'El apellido es obligatorio';
    isValid = false;
  }
  
  if (!isValid) return;
  
  // Enviar datos
  isSubmitting.value = true;
  
  try {
    const userData = {
      firstName: personalInfo.firstName,
      lastName: personalInfo.lastName,
      department: personalInfo.department
    };
    
    if (selectedUserId.value) {
      await usersStore.updateUser(selectedUserId.value, userData);
      // Mostrar mensaje de éxito
      alert('Información personal actualizada correctamente');
    }
  } catch (error) {
    console.error('Error al actualizar información personal:', error);
    // Mostrar mensaje de error
    alert('Ocurrió un error al actualizar la información personal');
  } finally {
    isSubmitting.value = false;
  }
};

// Actualizar rol del usuario
const updateUserRole = async () => {
  // Reset de errores
  errors.role = '';
  
  // Validación básica
  let isValid = true;
  
  if (!userRole.role) {
    errors.role = 'El rol es obligatorio';
    isValid = false;
  }
  
  if (!isValid) return;
  
  // Enviar datos
  isSubmitting.value = true;
  
  try {
    if (selectedUserId.value) {
      // Actualizar rol
      await usersStore.updateUserRole(selectedUserId.value, userRole.role);
      
      // Actualizar estado
      await usersStore.toggleUserActive(selectedUserId.value);
      
      // Mostrar mensaje de éxito
      alert('Rol y estado del usuario actualizados correctamente');
    }
  } catch (error) {
    console.error('Error al actualizar rol del usuario:', error);
    // Mostrar mensaje de error
    alert('Ocurrió un error al actualizar el rol del usuario');
  } finally {
    isSubmitting.value = false;
  }
};

// Resetear contraseña
const resetPassword = async () => {
  // Reset de errores
  errors.newPassword = '';
  errors.confirmPassword = '';
  
  // Validación básica
  let isValid = true;
  
  if (!passwordData.newPassword) {
    errors.newPassword = 'La nueva contraseña es obligatoria';
    isValid = false;
  } else if (passwordData.newPassword.length < 8) {
    errors.newPassword = 'La contraseña debe tener al menos 8 caracteres';
    isValid = false;
  }
  
  if (passwordData.newPassword !== passwordData.confirmPassword) {
    errors.confirmPassword = 'Las contraseñas no coinciden';
    isValid = false;
  }
  
  if (!isValid) return;
  
  // Enviar datos
  isSubmitting.value = true;
  
  try {
    if (selectedUserId.value) {
      // Aquí iría la lógica para resetear la contraseña
      // Por ahora simulamos un delay
      await new Promise(resolve => setTimeout(resolve, 1000));
      
      // Mostrar mensaje de éxito
      alert('Contraseña reseteada correctamente');
      
      // Limpiar campos
      passwordData.newPassword = '';
      passwordData.confirmPassword = '';
    }
  } catch (error) {
    console.error('Error al resetear contraseña:', error);
    // Mostrar mensaje de error
    alert('Ocurrió un error al resetear la contraseña');
  } finally {
    isSubmitting.value = false;
  }
};

// Calcular fuerza de la contraseña
const passwordStrength = computed(() => {
  const password = passwordData.newPassword;
  if (!password) {
    return { percent: 0, class: '', label: '' };
  }
  
  let strength = 0;
  let tips = [];
  
  // Longitud
  if (password.length >= 8) {
    strength += 25;
  } else {
    tips.push('Debe tener al menos 8 caracteres');
  }
  
  // Mayúsculas y minúsculas
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) {
    strength += 25;
  } else {
    tips.push('Debe incluir mayúsculas y minúsculas');
  }
  
  // Números
  if (/\d/.test(password)) {
    strength += 25;
  } else {
    tips.push('Debe incluir al menos un número');
  }
  
  // Caracteres especiales
  if (/[^a-zA-Z0-9]/.test(password)) {
    strength += 25;
  } else {
    tips.push('Debe incluir al menos un carácter especial');
  }
  
  // Establecer clase y etiqueta
  let strengthClass = '';
  let strengthLabel = '';
  
  if (strength <= 25) {
    strengthClass = 'weak';
    strengthLabel = 'Débil';
  } else if (strength <= 50) {
    strengthClass = 'medium';
    strengthLabel = 'Media';
  } else if (strength <= 75) {
    strengthClass = 'good';
    strengthLabel = 'Buena';
  } else {
    strengthClass = 'strong';
    strengthLabel = 'Fuerte';
  }
  
  return {
    percent: strength,
    class: strengthClass,
    label: strengthLabel,
    tips: tips
  };
});

// Función para traducir roles
const translateRole = (role: string) => {
  const roles: Record<string, string> = {
    admin: 'Administrador',
    assistant: 'Asistente',
    employee: 'Empleado'
  };
  
  return roles[role] || role;
};
</script>

<style scoped lang="scss">
.profile-management {
  padding: 1.5rem;
  max-width: 1200px;
  margin: 0 auto;
}

.profile-header {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
  background-color: var(--bg-secondary);
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  border: 1px solid var(--border-color);
  
  .profile-title {
    margin: 0 0 0.5rem 0;
    font-size: 1.75rem;
    font-weight: 600;
    color: var(--text-primary);
  }
  
  .profile-subtitle {
    margin: 0;
    color: var(--text-secondary);
  }
}

.profile-content {
  display: flex;
  flex-direction: column;
  gap: 1.5rem;
}

.user-selection {
  margin-bottom: 1.5rem;
  padding: 1.5rem;
  background-color: var(--bg-secondary);
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  border: 1px solid var(--border-color);
  
  label {
    display: block;
    margin-bottom: 0.75rem;
    font-weight: 500;
    color: var(--text-secondary);
    font-size: 1rem;
    font-family: inherit;
  }
  
  .select-container {
    position: relative;
    
    select {
      width: 100%;
      padding: 0.75rem;
      padding-right: 2.5rem;
      background-color: var(--bg-tertiary);
      border: 1px solid var(--border-color);
      border-radius: var(--border-radius);
      color: var(--text-primary);
      appearance: none;
      font-size: 0.95rem;
      font-family: inherit;
      
      &:focus {
        outline: none;
        border-color: var(--primary-color);
        box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
      }
    }
    
    i {
      position: absolute;
      right: 0.75rem;
      top: 50%;
      transform: translateY(-50%);
      color: var(--text-secondary);
      pointer-events: none;
    }
  }
}

.profile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 1.5rem;
  margin-bottom: 1.5rem;
}

.profile-card {
  background-color: var(--bg-secondary);
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  border: 1px solid var(--border-color);
  overflow: hidden;
  
  .profile-card-header {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 1.25rem 1.5rem;
    background-color: var(--bg-tertiary);
    border-bottom: 1px solid var(--border-color);
    
    i {
      font-size: 1.2rem;
      color: var(--primary-color);
    }
    
    h2 {
      margin: 0;
      font-size: 1.2rem;
      color: var(--text-primary);
    }
  }
  
  .profile-card-body {
    padding: 1.5rem;
  }
  
  .profile-card-actions {
    margin-top: 1.5rem;
    display: flex;
    justify-content: flex-end;
    
    .btn {
      padding: 0.75rem 1.5rem;
      border-radius: var(--border-radius);
      font-weight: 500;
      cursor: pointer;
      border: none;
      display: inline-flex;
      align-items: center;
      gap: 0.5rem;
      
      &.btn-primary {
        background-color: var(--primary-color);
        color: white;
        
        &:hover {
          background-color: var(--primary-hover);
        }
        
        &:disabled {
          opacity: 0.6;
          cursor: not-allowed;
        }
      }
    }
  }
}

.profile-form-group {
  margin-bottom: 1.25rem;
  
  &:last-child {
    margin-bottom: 0;
  }
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-secondary);
    font-size: 0.95rem;
    font-family: inherit;
  }
  
  input, select, textarea {
    width: 100%;
    padding: 0.75rem;
    background-color: var(--bg-tertiary);
    border: 1px solid var(--border-color);
    border-radius: var(--border-radius);
    color: var(--text-primary);
    font-size: 0.95rem;
    font-family: inherit;
    
    &:focus {
      outline: none;
      border-color: var(--primary-color);
      box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
    }
    
    &.is-invalid {
      border-color: var(--danger-color);
    }

    &::placeholder {
      color: var(--text-tertiary);
      opacity: 0.7;
    }

    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }
  
  .invalid-feedback {
    display: block;
    margin-top: 0.25rem;
    color: var(--danger-color);
    font-size: 0.85rem;
    font-family: inherit;
  }
  
  .password-input-container {
    position: relative;
    
    input {
      padding-right: 3rem;
    }
    
    .password-toggle {
      position: absolute;
      right: 0.75rem;
      top: 50%;
      transform: translateY(-50%);
      background: none;
      border: none;
      cursor: pointer;
      color: var(--text-secondary);
      
      &:hover {
        color: var(--text-primary);
      }
    }
  }
  
  .toggle-switch-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    
    .toggle-switch {
      position: relative;
      display: inline-block;
      width: 3.5rem;
      height: 1.75rem;
      
      input {
        opacity: 0;
        width: 0;
        height: 0;
        
        &:checked + .slider {
          background-color: var(--primary-color);
        }
        
        &:checked + .slider:before {
          transform: translateX(1.75rem);
        }
      }
      
      .slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: var(--text-muted);
        transition: .4s;
        
        &:before {
          position: absolute;
          content: "";
          height: 1.25rem;
          width: 1.25rem;
          left: 0.25rem;
          bottom: 0.25rem;
          background-color: white;
          transition: .4s;
        }
        
        &.round {
          border-radius: 34px;
          
          &:before {
            border-radius: 50%;
          }
        }
      }
    }
    
    .toggle-label {
      font-weight: 500;
      color: var(--text-primary);
      font-size: 0.95rem;
      font-family: inherit;
    }
  }
}

.toast-notification {
  position: fixed;
  bottom: 2rem;
  right: 2rem;
  padding: 1rem 1.5rem;
  display: flex;
  align-items: center;
  gap: 0.75rem;
  border-radius: var(--border-radius);
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  transform: translateY(150%);
  transition: transform 0.3s ease-in-out;
  z-index: 1000;
  
  &.show {
    transform: translateY(0);
  }
  
  &.success {
    background-color: rgba(46, 213, 115, 0.95);
    color: white;
  }
  
  &.error {
    background-color: rgba(255, 71, 87, 0.95);
    color: white;
  }
  
  i {
    font-size: 1.2rem;
  }
}

@media (max-width: 768px) {
  .profile-grid {
    grid-template-columns: 1fr;
  }
}
</style> 
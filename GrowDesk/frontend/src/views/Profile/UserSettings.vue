<template>
  <div class="profile-management">
    <div class="profile-header">
      <h1 class="profile-title">Configuración de mi cuenta</h1>
      <p class="profile-subtitle">Actualiza tu información personal y preferencias</p>
    </div>
    
    <div class="profile-content">
      <!-- Sección principal con tarjetas para diferentes aspectos del perfil -->
      <div class="profile-grid">
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
              
              <div class="profile-form-group">
                <label for="position">Posición</label>
                <input 
                  id="position" 
                  v-model="personalInfo.position" 
                  type="text" 
                  class="form-control"
                  :class="{ 'is-invalid': errors.position }"
                />
                <div v-if="errors.position" class="invalid-feedback">
                  {{ errors.position }}
                </div>
              </div>
              
              <div class="profile-form-group">
                <label for="phone">Teléfono</label>
                <input 
                  id="phone" 
                  v-model="personalInfo.phone" 
                  type="tel" 
                  class="form-control"
                  :class="{ 'is-invalid': errors.phone }"
                />
                <div v-if="errors.phone" class="invalid-feedback">
                  {{ errors.phone }}
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
        
        <!-- Cambio de contraseña -->
        <div class="profile-card">
          <div class="profile-card-header">
            <i class="pi pi-lock"></i>
            <h2>Cambiar Contraseña</h2>
          </div>
          <div class="profile-card-body">
            <form @submit.prevent="updatePassword">
              <div class="profile-form-group">
                <label for="currentPassword">Contraseña Actual</label>
                <div class="password-input-container">
                  <input 
                    id="currentPassword" 
                    v-model="passwordData.currentPassword" 
                    :type="showCurrentPassword ? 'text' : 'password'" 
                    class="form-control"
                    required
                    :class="{ 'is-invalid': errors.currentPassword }"
                  />
                  <button 
                    type="button" 
                    class="password-toggle" 
                    @click="showCurrentPassword = !showCurrentPassword"
                  >
                    <i :class="showCurrentPassword ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
                  </button>
                </div>
                <div v-if="errors.currentPassword" class="invalid-feedback">
                  {{ errors.currentPassword }}
                </div>
              </div>
              
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
                  <i class="pi pi-lock"></i> Actualizar Contraseña
                </button>
              </div>
            </form>
          </div>
        </div>
        
        <!-- Preferencias -->
        <div class="profile-card">
          <div class="profile-card-header">
            <i class="pi pi-cog"></i>
            <h2>Preferencias</h2>
          </div>
          <div class="profile-card-body">
            <form @submit.prevent="updatePreferences">
              <div class="profile-form-group">
                <label>Notificaciones por Email</label>
                <div class="toggle-switch-container">
                  <label class="toggle-switch">
                    <input 
                      type="checkbox" 
                      v-model="preferences.emailNotifications"
                    >
                    <span class="slider round"></span>
                  </label>
                  <span class="toggle-label">
                    {{ preferences.emailNotifications ? 'Activadas' : 'Desactivadas' }}
                  </span>
                </div>
              </div>
              
              <div class="profile-form-group">
                <label>Notificaciones del Sistema</label>
                <div class="toggle-switch-container">
                  <label class="toggle-switch">
                    <input 
                      type="checkbox" 
                      v-model="preferences.systemNotifications"
                    >
                    <span class="slider round"></span>
                  </label>
                  <span class="toggle-label">
                    {{ preferences.systemNotifications ? 'Activadas' : 'Desactivadas' }}
                  </span>
                </div>
              </div>
              
              <div class="profile-form-group">
                <label for="language">Idioma</label>
                <select 
                  id="language" 
                  v-model="preferences.language" 
                  class="form-control"
                >
                  <option value="es">Español</option>
                  <option value="en">English</option>
                </select>
              </div>
              
              <div class="profile-card-actions">
                <button 
                  type="submit" 
                  class="btn btn-primary"
                  :disabled="isSubmitting"
                >
                  <i class="pi pi-save"></i> Guardar Preferencias
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>
    
    <!-- Toast notifications -->
    <div 
      class="toast-notification"
      :class="{ 
        'show': notification.show,
        'success': notification.type === 'success',
        'error': notification.type === 'error'
      }"
      v-if="notification.show"
    >
      <i :class="notification.icon"></i>
      <span>{{ notification.message }}</span>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue';
import { useAuthStore } from '@/stores/auth';
import { useUserStore } from '@/stores/users';
import { storeToRefs } from 'pinia';

// Stores
const authStore = useAuthStore();
const userStore = useUserStore();

// Estado local para formularios
const personalInfo = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  position: '',
  phone: ''
});

const passwordData = reactive({
  currentPassword: '',
  newPassword: '',
  confirmPassword: ''
});

const preferences = reactive({
  emailNotifications: true,
  systemNotifications: true,
  language: 'es'
});

// Estados para mostrar/ocultar contraseñas
const showCurrentPassword = ref(false);
const showNewPassword = ref(false);
const showConfirmPassword = ref(false);

// Estado de carga y envío
const isSubmitting = ref(false);
const isLoading = ref(true);

// Estado para notificaciones
const notification = reactive({
  show: false,
  message: '',
  type: 'success',
  icon: 'pi pi-check'
});

// Validación de errores
const errors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  department: '',
  position: '',
  phone: '',
  currentPassword: '',
  newPassword: '',
  confirmPassword: '',
  general: ''
});

// Cargar datos del usuario
onMounted(async () => {
  isLoading.value = true;
  
  try {
    // Verificar si hay un usuario autenticado
    if (!authStore.user) {
      await authStore.checkAuth();
    }
    
    // Cargar datos del usuario
    if (authStore.user) {
      personalInfo.firstName = authStore.user.firstName;
      personalInfo.lastName = authStore.user.lastName;
      personalInfo.email = authStore.user.email;
      personalInfo.department = authStore.user.department || '';
      personalInfo.position = authStore.user.position || '';
      personalInfo.phone = authStore.user.phone || '';
      
      // Cargar preferencias (en un escenario real, vendrían de la API)
      // Aquí usamos valores por defecto para la demo
    }
  } catch (error) {
    console.error('Error al cargar datos del usuario:', error);
  } finally {
    isLoading.value = false;
  }
});

// Calcular la fortaleza de la contraseña
const passwordStrength = computed(() => {
  const password = passwordData.newPassword;
  
  if (!password) {
    return { class: '', percent: 0, label: '' };
  }
  
  let strength = 0;
  let feedback = '';
  
  // Longitud mínima
  if (password.length >= 8) strength += 25;
  
  // Contiene números
  if (/\d/.test(password)) strength += 25;
  
  // Contiene letras minúsculas y mayúsculas
  if (/[a-z]/.test(password) && /[A-Z]/.test(password)) strength += 25;
  
  // Contiene caracteres especiales
  if (/[^a-zA-Z0-9]/.test(password)) strength += 25;
  
  // Determinar nivel de fortaleza
  let strengthClass = '';
  
  if (strength <= 25) {
    strengthClass = 'weak';
    feedback = 'Débil';
  } else if (strength <= 50) {
    strengthClass = 'medium';
    feedback = 'Media';
  } else if (strength <= 75) {
    strengthClass = 'good';
    feedback = 'Buena';
  } else {
    strengthClass = 'strong';
    feedback = 'Fuerte';
  }
  
  return {
    class: strengthClass,
    percent: strength,
    label: feedback
  };
});

// Validación del formulario de información personal
const validatePersonalInfo = () => {
  let isValid = true;
  
  // Restablecer errores
  errors.firstName = '';
  errors.lastName = '';
  errors.email = '';
  errors.department = '';
  errors.position = '';
  errors.phone = '';
  
  // Validar nombre
  if (!personalInfo.firstName.trim()) {
    errors.firstName = 'El nombre es obligatorio';
    isValid = false;
  }
  
  // Validar apellido
  if (!personalInfo.lastName.trim()) {
    errors.lastName = 'El apellido es obligatorio';
    isValid = false;
  }
  
  // Validar teléfono (opcional)
  if (personalInfo.phone && !/^[+]?[(]?[0-9]{3}[)]?[-\s.]?[0-9]{3}[-\s.]?[0-9]{4,6}$/.test(personalInfo.phone)) {
    errors.phone = 'Formato de teléfono inválido';
    isValid = false;
  }
  
  return isValid;
};

// Validación del formulario de contraseña
const validatePasswordForm = () => {
  let isValid = true;
  
  // Restablecer errores
  errors.currentPassword = '';
  errors.newPassword = '';
  errors.confirmPassword = '';
  
  // Validar contraseña actual
  if (!passwordData.currentPassword) {
    errors.currentPassword = 'La contraseña actual es obligatoria';
    isValid = false;
  }
  
  // Validar nueva contraseña
  if (!passwordData.newPassword) {
    errors.newPassword = 'La nueva contraseña es obligatoria';
    isValid = false;
  } else if (passwordData.newPassword.length < 8) {
    errors.newPassword = 'La contraseña debe tener al menos 8 caracteres';
    isValid = false;
  }
  
  // Validar confirmación de contraseña
  if (!passwordData.confirmPassword) {
    errors.confirmPassword = 'Debe confirmar la nueva contraseña';
    isValid = false;
  } else if (passwordData.newPassword !== passwordData.confirmPassword) {
    errors.confirmPassword = 'Las contraseñas no coinciden';
    isValid = false;
  }
  
  return isValid;
};

// Mostrar notificación
const showNotification = (message: string, type: 'success' | 'error') => {
  notification.message = message;
  notification.type = type;
  notification.icon = type === 'success' ? 'pi pi-check' : 'pi pi-times';
  notification.show = true;
  
  // Ocultar después de 3 segundos
  setTimeout(() => {
    notification.show = false;
  }, 3000);
};

// Actualizar información personal
const updatePersonalInfo = async () => {
  if (!validatePersonalInfo()) return;
  
  isSubmitting.value = true;
  
  try {
    // Preparar datos para la actualización
    const userData = {
      firstName: personalInfo.firstName.trim(),
      lastName: personalInfo.lastName.trim(),
      department: personalInfo.department || null,
      position: personalInfo.position || null,
      phone: personalInfo.phone || null
    };
    
    // Llamar al método de actualización de perfil
    if (authStore.user && authStore.user.id) {
      await userStore.updateUser(authStore.user.id, userData);
      
      // Actualizar el usuario en el authStore
      await authStore.updateProfile(userData);
      
      showNotification('Información personal actualizada con éxito', 'success');
    }
  } catch (error) {
    console.error('Error al actualizar información personal:', error);
    showNotification('Error al actualizar la información personal', 'error');
  } finally {
    isSubmitting.value = false;
  }
};

// Actualizar contraseña
const updatePassword = async () => {
  if (!validatePasswordForm()) return;
  
  isSubmitting.value = true;
  
  try {
    // En un entorno real, aquí llamaríamos a la API para cambiar la contraseña
    // Por ahora, simulamos una respuesta exitosa
    
    // Simular demora
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    showNotification('Contraseña actualizada con éxito', 'success');
    
    // Limpiar formulario
    passwordData.currentPassword = '';
    passwordData.newPassword = '';
    passwordData.confirmPassword = '';
  } catch (error) {
    console.error('Error al actualizar contraseña:', error);
    showNotification('Error al actualizar la contraseña', 'error');
  } finally {
    isSubmitting.value = false;
  }
};

// Actualizar preferencias
const updatePreferences = async () => {
  isSubmitting.value = true;
  
  try {
    // En un entorno real, aquí llamaríamos a la API para guardar las preferencias
    // Por ahora, simulamos una respuesta exitosa
    
    // Simular demora
    await new Promise(resolve => setTimeout(resolve, 1000));
    
    showNotification('Preferencias guardadas con éxito', 'success');
  } catch (error) {
    console.error('Error al actualizar preferencias:', error);
    showNotification('Error al guardar las preferencias', 'error');
  } finally {
    isSubmitting.value = false;
  }
};
</script>

<style scoped lang="scss">
.profile-management {
  max-width: 1200px;
  margin: 0 auto;
  padding: 1.5rem;
}

.profile-header {
  margin-bottom: 2rem;
  
  .profile-title {
    font-size: 1.75rem;
    margin-bottom: 0.5rem;
    color: var(--text-primary);
  }
  
  .profile-subtitle {
    font-size: 1rem;
    color: var(--text-secondary);
  }
}

.profile-grid {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(400px, 1fr));
  gap: 1.5rem;
  
  @media (max-width: 768px) {
    grid-template-columns: 1fr;
  }
}

.profile-card {
  background-color: var(--card-bg);
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  overflow: hidden;
  border: 1px solid var(--border-color);
  margin-bottom: 1.5rem;
  
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
  }
}

.profile-form-group {
  margin-bottom: 1.25rem;
  
  label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-secondary);
  }
  
  .form-control {
    width: 100%;
    padding: 0.75rem;
    border: 1px solid var(--input-border);
    border-radius: var(--border-radius);
    font-size: 1rem;
    background-color: var(--input-bg);
    color: var(--text-primary);
    
    &:focus {
      outline: none;
      border-color: var(--primary-color);
      box-shadow: 0 0 0 2px rgba(79, 70, 229, 0.2);
    }
    
    &.is-invalid {
      border-color: var(--danger-color);
    }
    
    &:disabled {
      background-color: var(--bg-tertiary);
      cursor: not-allowed;
    }
  }
  
  .invalid-feedback {
    display: block;
    color: var(--danger-color);
    font-size: 0.875rem;
    margin-top: 0.25rem;
  }
  
  .form-text {
    font-size: 0.875rem;
    color: var(--text-muted);
    margin-top: 0.25rem;
  }
}

.password-input-container {
  position: relative;
  
  .password-toggle {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    background: none;
    border: none;
    padding: 5px;
    color: var(--text-secondary);
    cursor: pointer;
    
    &:hover {
      color: var(--text-primary);
    }
  }
}

.password-strength {
  margin-top: 0.5rem;
  
  .strength-meter {
    height: 5px;
    background-color: var(--bg-tertiary);
    border-radius: 3px;
    margin-bottom: 0.25rem;
    overflow: hidden;
    
    .strength-meter-fill {
      height: 100%;
      border-radius: 3px;
      transition: width 0.3s;
      
      &.weak {
        background-color: var(--danger-color);
      }
      
      &.medium {
        background-color: var(--warning-color);
      }
      
      &.good {
        background-color: var(--success-color);
      }
      
      &.strong {
        background-color: var(--success-color);
      }
    }
  }
  
  .strength-text {
    font-size: 0.75rem;
    font-weight: 500;
    
    &.weak {
      color: var(--danger-color);
    }
    
    &.medium {
      color: var(--warning-color);
    }
    
    &.good {
      color: var(--success-color);
    }
    
    &.strong {
      color: var(--success-color);
    }
  }
}

.toggle-switch-container {
  display: flex;
  align-items: center;
  gap: 1rem;
  
  .toggle-label {
    font-weight: 500;
  }
}

.toggle-switch {
  position: relative;
  display: inline-block;
  width: 50px;
  height: 24px;
  
  input {
    opacity: 0;
    width: 0;
    height: 0;
    
    &:checked + .slider {
      background-color: var(--primary-color);
    }
    
    &:checked + .slider:before {
      transform: translateX(26px);
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
      height: 16px;
      width: 16px;
      left: 4px;
      bottom: 4px;
      background-color: white;
      transition: .4s;
    }
    
    &.round {
      border-radius: 24px;
      
      &:before {
        border-radius: 50%;
      }
    }
  }
}

.btn {
  padding: 0.75rem 1.5rem;
  border-radius: var(--border-radius);
  font-weight: 500;
  cursor: pointer;
  border: none;
  display: flex;
  align-items: center;
  gap: 0.5rem;
  
  &.btn-primary {
    background-color: var(--primary-color);
    color: white;
    
    &:hover {
      background-color: var(--primary-hover);
    }
    
    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }
}

.toast-notification {
  position: fixed;
  bottom: 20px;
  right: 20px;
  padding: 1rem 1.5rem;
  background-color: var(--card-bg);
  color: var(--text-primary);
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  display: flex;
  align-items: center;
  gap: 0.75rem;
  z-index: 1000;
  transform: translateY(100px);
  opacity: 0;
  transition: transform 0.3s, opacity 0.3s;
  
  i {
    font-size: 1.25rem;
  }
  
  &.show {
    transform: translateY(0);
    opacity: 1;
  }
  
  &.success {
    border-left: 4px solid var(--success-color);
    
    i {
      color: var(--success-color);
    }
  }
  
  &.error {
    border-left: 4px solid var(--danger-color);
    
    i {
      color: var(--danger-color);
    }
  }
}
</style> 
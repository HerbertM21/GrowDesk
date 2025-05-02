<template>
  <div class="admin-section">
    <!-- Permiso verificado: Solo admins y asistentes -->
    <div v-if="!isAdminOrAssistant" class="permission-error">
      <div class="error-container">
        <i class="pi pi-exclamation-triangle"></i>
        <h2>Acceso restringido</h2>
        <p>No tienes permiso para acceder a esta sección. Esta página solo está disponible para administradores y asistentes.</p>
        <button class="btn btn-primary" @click="goToDashboard">Volver al Dashboard</button>
      </div>
    </div>

    <div v-else>
      <div class="admin-section">
        <div class="profile-management">
          <!-- Sección del encabezado con fondo de color primario y forma ondulada -->
          <div class="hero-section">
            <div class="hero-content">
              <h1 class="hero-title">Gestión de Perfiles de Usuarios</h1>
              <p class="hero-subtitle">Administre los perfiles de los usuarios del sistema</p>
            </div>
            <div class="wave-shape"></div>
          </div>
          
          <div class="content-wrapper">
            <!-- Vista de selección de usuario -->
            <div class="filters-section">
              <h2 class="section-title">
                <span class="title-icon"><i class="pi pi-user"></i></span>
                Selección de Usuario
              </h2>
              
              <div class="filters-container">
                <div class="filter-group">
                  <label for="user-select">Seleccione un usuario para gestionar su perfil:</label>
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
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted, watch } from 'vue';
import { useUsersStore } from '@/stores/users';
import type { User } from '@/stores/users';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useNotificationsStore } from '@/stores/notifications';

const usersStore = useUsersStore();
const router = useRouter();
const authStore = useAuthStore();
const notificationsStore = useNotificationsStore();

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

// Verificación de permisos
const isAdminOrAssistant = computed(() => {
  return authStore.isAdmin || authStore.isAssistant;
});

// Función para redireccionar al dashboard
function goToDashboard() {
  router.push('/dashboard');
}

// Mostrar notificación si no tiene permiso
onMounted(() => {
  if (!isAdminOrAssistant.value) {
    notificationsStore.addNotification({
      type: 'error',
      message: 'No tienes permiso para acceder a la gestión de perfiles.',
      timeout: 5000
    });
  }
});

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
  --primary-gradient: linear-gradient(135deg, var(--primary-color) 0%, #3b82f6 100%);
  --secondary-gradient: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  --border-radius-lg: 1.25rem;
  --transition-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
  
  background-color: var(--bg-secondary);
  position: relative;
  overflow-x: hidden;
  
  // Sección hero con fondo primario
  .hero-section {
    position: relative;
    padding: 2.5rem 2rem 4.5rem;
    background-color: var(--primary-color);
    color: white;
    text-align: center;
    overflow: hidden;
    
    .hero-content {
      position: relative;
      z-index: 2;
      max-width: 800px;
      margin: 0 auto;
    }
    
    .hero-title {
      font-size: 2.25rem;
      font-weight: 700;
      margin-bottom: 0.75rem;
      text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      color: white;
    }
    
    .hero-subtitle {
      font-size: 1.1rem;
      margin-bottom: 1.75rem;
      opacity: 0.9;
      color: white;
    }
    
    // Forma ondulada en la parte inferior
    .wave-shape {
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 100%;
      height: 4rem;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' fill='%23f8fafc' opacity='.25'%3E%3C/path%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' fill='%23f8fafc' opacity='.5'%3E%3C/path%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%23f8fafc'%3E%3C/path%3E%3C/svg%3E");
      background-size: cover;
      background-position: center;
    }
  }
  
  .content-wrapper {
    max-width: 1300px;
    margin: 0 auto;
    padding: 3rem 1.5rem;
  }
  
  // Títulos de sección con iconos
  .section-title {
    display: flex;
    align-items: center;
    justify-content: flex-start;
    margin-bottom: 1.5rem;
    color: var(--text-primary);
    font-size: 1.5rem;
    font-weight: 600;
    text-align: left;
    background-color: var(--bg-tertiary);
    border-radius: 12px;
    padding: 0.75rem 1.5rem;
    box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
    border-left: 4px solid var(--primary-color);
    
    .title-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 38px;
      height: 38px;
      background-color: var(--primary-color);
      border-radius: 10px;
      margin-right: 1rem;
      color: white;
      box-shadow: 0 4px 10px rgba(var(--primary-rgb), 0.25);
      
      i {
        font-size: 1.2rem;
      }
    }
  }
  
  // Sección de filtros (adaptada del estilo de TicketList)
  .filters-section {
    margin-bottom: 3rem;
    
    .filters-container {
      background-color: var(--card-bg);
      border-radius: var(--border-radius-lg);
      box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
      padding: 1.75rem;
      border: 1px solid var(--border-color);
      
      .filter-group {
        margin-bottom: 1.5rem;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        label {
          display: block;
          margin-bottom: 0.75rem;
          font-weight: 600;
          color: var(--text-primary);
          font-size: 1rem;
        }
        
        .select-container {
          position: relative;
          
          select {
            width: 100%;
            padding: 0.85rem;
            padding-right: 2.5rem;
            background-color: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
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
            right: 0.85rem;
            top: 50%;
            transform: translateY(-50%);
            color: var(--text-secondary);
            pointer-events: none;
          }
        }
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
    background-color: var(--card-bg);
    border-radius: 24px; /* Aumentando el radio de bordes como en TicketList */
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    overflow: hidden;
    transition: all 0.3s var(--transition-bounce);
    
    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                  0 10px 10px -5px rgba(0, 0, 0, 0.04);
    }
    
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
  }
  
  // Estado vacío cuando no hay usuario seleccionado
  .empty-state {
    text-align: center;
    padding: 3rem;
    background-color: var(--card-bg);
    border-radius: var(--border-radius-lg);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    border: 1px solid var(--border-color);
    
    i {
      color: var(--text-secondary);
      margin-bottom: 1rem;
    }
    
    p {
      font-size: 1.1rem;
      color: var(--text-secondary);
    }
  }

  // Botones de acción mejorados con márgenes apropiados
  .profile-card-actions {
    display: flex;
    justify-content: flex-end;
    margin-top: 1.5rem;
    padding-top: 1rem;
    border-top: 1px solid var(--border-color);
    
    button {
      min-width: 160px;
      margin-left: 0.75rem;
      
      &:first-child {
        margin-left: 0;
      }
      
      i {
        margin-right: 0.5rem;
      }
    }
  }
  
  // Estilos mejorados para botones
  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    padding: 0.65rem 1.25rem;
    border-radius: 0.5rem;
    font-weight: 600;
    font-size: 0.9rem;
    transition: all 0.2s ease;
    cursor: pointer;
    border: none;
    
    &.btn-primary {
      background: var(--primary-color); // Color sólido en lugar de gradiente
      color: white;
      
      &:hover, &:focus {
        box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.3);
        transform: translateY(-2px);
      }
      
      &:disabled {
        opacity: 0.7;
        cursor: not-allowed;
        transform: none !important;
      }
    }
  }
  
  // Grupos de formulario con espaciado apropiado
  .profile-form-group {
    margin-bottom: 1.5rem;
    
    label {
      display: block;
      margin-bottom: 0.5rem;
      font-weight: 600;
      color: var(--text-primary);
    }
    
    .form-control {
      width: 100%;
      padding: 0.75rem 1rem;
      border: 1px solid var(--border-color);
      border-radius: 0.5rem;
      font-size: 0.95rem;
      color: var(--text-primary);
      background-color: var(--bg-tertiary);
      transition: all 0.2s ease;
      
      &:focus {
        border-color: var(--primary-color);
        box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.15);
        outline: none;
      }
      
      &:disabled {
        opacity: 0.7;
        cursor: not-allowed;
      }
      
      &.is-invalid {
        border-color: var(--danger-color);
        
        &:focus {
          box-shadow: 0 0 0 3px rgba(220, 38, 38, 0.15);
        }
      }
    }
    
    .invalid-feedback {
      display: block;
      margin-top: 0.5rem;
      color: var(--danger-color);
      font-size: 0.85rem;
    }
    
    .form-text {
      display: block;
      margin-top: 0.5rem;
      font-size: 0.85rem;
      color: var(--text-secondary);
    }
  }
  
  // Contenedor de entrada de contraseña mejorado
  .password-input-container {
    position: relative;
    display: flex;
    align-items: center;
    
    .form-control {
      padding-right: 3rem;
      border-top-right-radius: 0;
      border-bottom-right-radius: 0;
      flex: 1;
    }
    
    // Botón de visibilidad de contraseña mejorado
    .password-toggle {
      position: absolute;
      right: 0;
      top: 0;
      height: 90%;
      width: 42px;
      display: flex;
      align-items: center;
      justify-content: center;
      background-color: var(--bg-tertiary);
      border: 1px solid var(--border-color);
      border-left: none;
      border-top-right-radius: 0.5rem;
      border-bottom-right-radius: 0.5rem;
      cursor: pointer;
      color: var(--text-secondary);
      transition: all 0.2s ease;
      padding: 0;
      
      &:hover, &:focus {
        background-color: rgba(var(--primary-rgb), 0.1);
        color: var(--primary-color);
      }
      
      i {
        font-size: 1rem;
      }
    }
  }
  
  // Medidor de fortaleza de contraseña
  .password-strength {
    margin-top: 0.75rem;
    
    .strength-meter {
      height: 6px;
      background-color: var(--border-color);
      border-radius: 3px;
      overflow: hidden;
      
      .strength-meter-fill {
        height: 100%;
        border-radius: 3px;
        transition: width 0.3s ease;
        
        &.weak {
          background-color: #ef4444;
        }
        
        &.medium {
          background-color: #f59e0b;
        }
        
        &.good {
          background-color: #10b981;
        }
        
        &.strong {
          background-color: #0ea5e9;
        }
      }
    }
    
    .strength-text {
      display: block;
      margin-top: 0.35rem;
      font-size: 0.85rem;
      font-weight: 500;
      
      &.weak {
        color: #ef4444;
      }
      
      &.medium {
        color: #f59e0b;
      }
      
      &.good {
        color: #10b981;
      }
      
      &.strong {
        color: #0ea5e9;
      }
    }
  }
  
  // Interruptor de toggle mejorado para estado activo/inactivo
  .toggle-switch-container {
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-top: 0.5rem;
    
    .toggle-label {
      font-weight: 500;
    }
    
    .toggle-switch {
      position: relative;
      display: inline-block;
      width: 50px;
      height: 26px;
      
      input {
        opacity: 0;
        width: 0;
        height: 0;
        
        &:checked + .slider {
          background-color: var(--primary-color);
          
          &:before {
            transform: translateX(24px);
          }
        }
        
        &:focus + .slider {
          box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.25);
        }
      }
      
      .slider {
        position: absolute;
        cursor: pointer;
        top: 0;
        left: 0;
        right: 0;
        bottom: 0;
        background-color: var(--text-tertiary);
        transition: .4s;
        
        &:before {
          position: absolute;
          content: "";
          height: 18px;
          width: 18px;
          left: 4px;
          bottom: 4px;
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
  }

  // Responsive adjustments
  @media (max-width: 768px) {
    .hero-section {
      padding: 2rem 1rem 4rem;
      
      .hero-title {
        font-size: 2rem;
      }
      
      .hero-subtitle {
        font-size: 1rem;
      }
    }
    
    .section-title {
      font-size: 1.5rem;
      flex-direction: column;
      
      .title-icon {
        margin-right: 0;
        margin-bottom: 0.75rem;
      }
    }
    
    .profile-grid {
      grid-template-columns: 1fr;
    }
    
    // Ajustes adicionales para dispositivos móviles
    .profile-card-actions {
      flex-direction: column;
      align-items: stretch;
      
      button {
        width: 100%;
        margin-left: 0;
        margin-top: 0.75rem;
        
        &:first-child {
          margin-top: 0;
        }
      }
    }
  }
}

/* Estilos para el mensaje de error de permisos */
.permission-error {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 120px);
  padding: 2rem;
}

.error-container {
  max-width: 500px;
  text-align: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.error-container i {
  font-size: 3rem;
  color: #f59e0b;
  margin-bottom: 1rem;
}

.error-container h2 {
  margin-bottom: 1rem;
  color: #334155;
}

.error-container p {
  margin-bottom: 2rem;
  color: #64748b;
}
</style>
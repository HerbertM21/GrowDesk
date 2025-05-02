/* eslint-disable */
<template>
  <div class="admin-section">
    <div v-if="!isAdmin" class="authorization-error">
      <div class="alert alert-danger">
        <i class="pi pi-exclamation-triangle"></i>
        <span>No tienes permisos suficientes para acceder a esta página</span>
      </div>
      <button @click="redirectToDashboard" class="btn btn-primary">
        <i class="pi pi-home"></i> Volver al Dashboard
      </button>
    </div>
    
    <template v-else>
      <div class="users-list">
        <!-- Sección del encabezado con fondo de color primario y forma ondulada -->
        <div class="hero-section">
          <div class="hero-content">
            <h1 class="hero-title">Administración de Usuarios</h1>
            <p class="hero-subtitle">Gestione todos los usuarios del sistema, asigne roles y controle su estado.</p>
            
            <button class="create-user-btn" @click="showCreateModal = true">
              <i class="pi pi-plus"></i>
              Nuevo Usuario
            </button>
          </div>
          <div class="wave-shape"></div>
        </div>
        
        <!-- Contenido principal -->
        <div class="content-wrapper">
          <!-- Filtros con nuevo diseño visual -->
          <div class="filters-section">
            <h2 class="section-title">
              <span class="title-icon"><i class="pi pi-filter"></i></span>
              Filtrar Usuarios
            </h2>
            
            <div class="filters-container">
              <div class="filter-group">
                <label>Rol:</label>
                <div class="filter-options">
                  <button 
                    @click="roleFilter = 'all'" 
                    :class="['filter-btn', roleFilter === 'all' ? 'active' : '']"
                  >
                    Todos
                  </button>
                  <button 
                    @click="roleFilter = 'admin'" 
                    :class="['filter-btn', roleFilter === 'admin' ? 'active' : '']"
                  >
                    Administradores
                  </button>
                  <button 
                    @click="roleFilter = 'assistant'" 
                    :class="['filter-btn', roleFilter === 'assistant' ? 'active' : '']"
                  >
                    Asistentes
                  </button>
                  <button 
                    @click="roleFilter = 'employee'" 
                    :class="['filter-btn', roleFilter === 'employee' ? 'active' : '']"
                  >
                    Empleados
                  </button>
                </div>
              </div>
              
              <div class="filter-group">
                <label>Estado:</label>
                <div class="filter-options">
                  <button 
                    @click="statusFilter = 'all'" 
                    :class="['filter-btn', statusFilter === 'all' ? 'active' : '']"
                  >
                    Todos
                  </button>
                  <button 
                    @click="statusFilter = 'active'" 
                    :class="['filter-btn', statusFilter === 'active' ? 'active' : '']"
                  >
                    Activos
                  </button>
                  <button 
                    @click="statusFilter = 'inactive'" 
                    :class="['filter-btn', statusFilter === 'inactive' ? 'active' : '']"
                  >
                    Inactivos
                  </button>
                </div>
              </div>
              
              <div class="filter-group">
                <label>Buscar:</label>
                <div class="search-container">
                  <input 
                    v-model="searchQuery" 
                    type="text" 
                    placeholder="Buscar por nombre o correo..." 
                    class="search-input"
                  />
                  <i class="pi pi-search search-icon"></i>
                </div>
              </div>
            </div>
          </div>

          <!-- Estado de carga, error o vacío -->
          <div v-if="loading" class="status-message loading">
            <i class="pi pi-spin pi-spinner"></i>
            <p>Cargando usuarios...</p>
          </div>
          
          <div v-else-if="error" class="status-message error">
            <i class="pi pi-exclamation-triangle"></i>
            <p>{{ error }}</p>
          </div>
          
          <div v-else-if="filteredUsers.length === 0" class="status-message empty">
            <i class="pi pi-inbox"></i>
            <p>No se encontraron usuarios con los filtros seleccionados</p>
          </div>
          
          <!-- Lista de usuarios con nuevo diseño -->
          <div v-else class="users-grid">
            <div 
              v-for="user in filteredUsers" 
              :key="user.id" 
              class="user-card"
            >
              <div class="user-header">
                <div class="user-badges">
                  <span :class="['role-badge', user.role]">{{ translateRole(user.role) }}</span>
                  <span :class="['status-badge', user.active ? 'active' : 'inactive']">
                    {{ user.active ? 'Activo' : 'Inactivo' }}
                  </span>
                </div>
                <h3 class="user-title">{{ user.firstName }} {{ user.lastName }}</h3>
              </div>
              
              <div class="user-body">
                <p class="user-email">{{ user.email }}</p>
                <div class="user-detail" v-if="user.department">
                  <span class="detail-label">Departamento:</span>
                  <span>{{ user.department }}</span>
                </div>
              </div>
              
              <div class="user-meta">
                <div class="meta-item" v-if="user.createdAt">
                  <i class="pi pi-calendar"></i>
                  <span>{{ formatDate(user.createdAt) }}</span>
                </div>
                
                <div class="meta-item">
                  <i class="pi pi-id-card"></i>
                  <span>{{ user.id }}</span>
                </div>
              </div>
              
              <div class="user-actions">
                <button class="action-btn" @click="editUser(user)" title="Editar usuario">
                  <i class="pi pi-pencil"></i>
                </button>
                <button class="action-btn" @click="changeRole(user)" title="Cambiar rol">
                  <i class="pi pi-users"></i>
                </button>
                <button 
                  class="action-btn" 
                  :class="user.active ? 'warning' : 'success'"
                  @click="toggleActive(user)" 
                  :title="user.active ? 'Desactivar usuario' : 'Activar usuario'"
                >
                  <i :class="user.active ? 'pi pi-eye-slash' : 'pi pi-eye'"></i>
                </button>
                <button 
                  class="action-btn danger" 
                  @click="confirmDelete(user)" 
                  v-if="currentUser?.id !== user.id"
                  title="Eliminar usuario"
                >
                  <i class="pi pi-trash"></i>
                </button>
              </div>
            </div>
          </div>
        </div>
        
        <!-- Modales -->
        <UserCreateModal v-if="showCreateModal" @close="showCreateModal = false" @created="handleUserCreated" />
        <UserEditModal v-if="showEditModal" :user="selectedUser" @close="showEditModal = false" @updated="handleUserUpdated" />
        <UserRoleModal v-if="showRoleModal" :user="selectedUser" @close="showRoleModal = false" @role-changed="handleRoleChanged" />
        <ConfirmDialog 
          v-if="showDeleteConfirm" 
          title="Confirmar eliminación" 
          message="¿Estás seguro de que quieres eliminar este usuario? Esta acción no se puede deshacer."
          @confirm="deleteUser"
          @cancel="showDeleteConfirm = false" 
        />
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useUsersStore } from '@/stores/users';
import { useAuthStore } from '@/stores/auth';
import { storeToRefs } from 'pinia';
import type { User } from '@/stores/users';
import UserCreateModal from '@/components/Admin/UserCreateModal.vue';
import UserEditModal from '@/components/Admin/UserEditModal.vue';
import UserRoleModal from '@/components/Admin/UserRoleModal.vue';
import ConfirmDialog from '@/components/common/ConfirmDialog.vue';
import { useRouter } from 'vue-router';
import { useNotificationsStore } from '@/stores/notifications';

// Stores
const userStore = useUsersStore();
const authStore = useAuthStore();
const notificationsStore = useNotificationsStore();
const router = useRouter();
const { users, loading, error } = storeToRefs(userStore);
const currentUser = computed(() => authStore.user);

// Filtros y búsqueda
const roleFilter = ref('all');
const statusFilter = ref('all');
const searchQuery = ref('');

// Estado de modales
const showCreateModal = ref(false);
const showEditModal = ref(false);
const showRoleModal = ref(false);
const showDeleteConfirm = ref(false);
const selectedUser = ref<User | null>(null);

// Verificar permisos
const isAdmin = computed(() => {
  return authStore.isAdmin;
});

// Filtrar usuarios según criterios
const filteredUsers = computed(() => {
  return users.value.filter(user => {
    // Filtro por rol
    if (roleFilter.value !== 'all' && user.role !== roleFilter.value) return false;
    
    // Filtro por estado
    if (statusFilter.value === 'active' && !user.active) return false;
    if (statusFilter.value === 'inactive' && user.active) return false;
    
    // Búsqueda por nombre o email
    if (searchQuery.value) {
      const query = searchQuery.value.toLowerCase();
      const fullName = `${user.firstName} ${user.lastName}`.toLowerCase();
      return fullName.includes(query) || user.email.toLowerCase().includes(query);
    }
    
    return true;
  });
});

// Traducir roles
const translateRole = (role: string) => {
  const roles: Record<string, string> = {
    'admin': 'Administrador',
    'assistant': 'Asistente',
    'employee': 'Empleado'
  };
  return roles[role] || role;
};

// Obtener iniciales del usuario
const getInitials = (user: User) => {
  return (user.firstName.charAt(0) + user.lastName.charAt(0)).toUpperCase();
};

// Generar color para avatar
const getAvatarColor = (user: User) => {
  const colors = [
    'var(--primary-color, #1976d2)',
    'var(--success-color, #388e3c)',
    'var(--warning-color, #f57c00)',
    'var(--info-color, #0288d1)',
    'var(--purple-color, #7b1fa2)'
  ];
  
  const colorIndex = parseInt(user.id) % colors.length;
  return colors[colorIndex];
};

// Formatear fecha
const formatDate = (dateString: string) => {
  if (!dateString) return '';
  const date = new Date(dateString);
  return date.toLocaleDateString();
};

// Funciones para acciones
const editUser = (user: User) => {
  selectedUser.value = user;
  showEditModal.value = true;
};

const changeRole = (user: User) => {
  selectedUser.value = user;
  showRoleModal.value = true;
};

const toggleActive = async (user: User) => {
  if (currentUser.value?.id === user.id) {
    alert('No puedes desactivar tu propia cuenta');
    return;
  }
  
  try {
    await userStore.toggleUserActive(user.id);
    console.log(`Usuario ${!user.active ? 'activado' : 'desactivado'} correctamente`);
  } catch (error) {
    console.error('Error al cambiar estado de usuario:', error);
  }
};

const confirmDelete = (user: User) => {
  selectedUser.value = user;
  showDeleteConfirm.value = true;
};

const deleteUser = async () => {
  if (selectedUser.value) {
    const success = await userStore.deleteUser(selectedUser.value.id);
    if (success) {
      console.log('Usuario eliminado correctamente');
    }
    showDeleteConfirm.value = false;
  }
};

// Handlers para eventos de modales
const handleUserCreated = () => {
  showCreateModal.value = false;
  userStore.fetchUsers(); // Recargar la lista
  console.log('Usuario creado: lista recargada');
};

const handleUserUpdated = () => {
  showEditModal.value = false;
  console.log('Evento de actualización recibido en UsersList');
  console.log('Recargando lista de usuarios...');
  userStore.fetchUsers(); // Recargar la lista
};

const handleRoleChanged = () => {
  showRoleModal.value = false;
  userStore.fetchUsers(); // Recargar la lista
  console.log('Rol actualizado: lista recargada');
};

// Función para redirigir al dashboard
const redirectToDashboard = () => {
  notificationsStore.warning('Has sido redirigido al Dashboard debido a permisos insuficientes');
  router.push('/dashboard');
};

// Hook de ciclo de vida
onMounted(() => {
  // Solo cargamos datos si el usuario es administrador
  if (isAdmin.value) {
    userStore.fetchUsers();
  } else {
    // Si no es admin, redirigimos al dashboard
    redirectToDashboard();
  }
  
  // Para modo de desarrollo con el mock server
  if (process.env.NODE_ENV === 'development' && users.value.length === 0) {
    userStore.initMockUsers();
  }
});
</script>

<style scoped lang="scss">
.users-list {
  --primary-gradient: linear-gradient(135deg, var(--primary-color) 0%, #3b82f6 100%);
  --secondary-gradient: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  --border-radius-lg: 1.25rem;
  --transition-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
  
  background-color: var(--bg-secondary);
  position: relative;
  overflow-x: hidden;
  
  // Sección hero con fondo más sutil
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
    
    .create-user-btn {
      display: inline-flex;
      align-items: center;
      gap: 0.75rem;
      background-color: white;
      color: var(--primary-color);
      border: none;
      padding: 0.85rem 1.75rem;
      border-radius: 8px;
      font-size: 1rem;
      font-weight: 600;
      text-decoration: none;
      transition: all 0.3s ease;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      cursor: pointer;
      
      &:hover {
        transform: translateY(-3px);
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
      }
      
      i {
        font-size: 1rem;
      }
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
  
  // Sección de filtros
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
        
        .filter-options {
          display: flex;
          flex-wrap: wrap;
          gap: 0.75rem;
          
          .filter-btn {
            background-color: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            color: var(--text-primary);
            font-size: 0.9rem;
            font-weight: 500;
            padding: 0.65rem 1.2rem;
            border-radius: 8px;
            cursor: pointer;
            transition: all 0.2s ease;
            
            &:hover {
              background-color: var(--hover-bg);
              transform: translateY(-2px);
            }
            
            &.active {
              background-color: var(--primary-color);
              color: white;
              border-color: transparent;
              font-weight: 600;
              box-shadow: 0 4px 8px rgba(0, 0, 0, 0.15);
            }
          }
        }
        
        .search-container {
          position: relative;
          
          .search-input {
            width: 100%;
            padding: 0.85rem;
            padding-right: 2.5rem;
            background-color: var(--bg-tertiary);
            border: 1px solid var(--border-color);
            border-radius: 8px;
            color: var(--text-primary);
            font-size: 0.95rem;
            
            &:focus {
              outline: none;
              border-color: var(--primary-color);
              box-shadow: 0 0 0 2px rgba(var(--primary-rgb), 0.2);
            }
            
            &::placeholder {
              color: var(--text-tertiary);
            }
          }
          
          .search-icon {
            position: absolute;
            right: 1rem;
            top: 50%;
            transform: translateY(-50%);
            color: var(--text-secondary);
          }
        }
      }
    }
  }

  // Mensajes de estado (loading, error, empty)
  .status-message {
    text-align: center;
    padding: 3rem;
    background-color: var(--card-bg);
    border-radius: var(--border-radius-lg);
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    margin-bottom: 2rem;
    
    i {
      font-size: 3rem;
      margin-bottom: 1rem;
      display: block;
    }
    
    p {
      font-size: 1.1rem;
      color: var(--text-secondary);
    }
    
    &.loading i {
      color: var(--primary-color);
    }
    
    &.error {
      i, p {
        color: #ef4444;
      }
    }
    
    &.empty i {
      color: #6b7280;
    }
  }

  // Grid de usuarios
  .users-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(340px, 1fr));
    gap: 1.5rem;
  }

  // Tarjeta de usuario con nuevo diseño
  .user-card {
    background-color: var(--card-bg);
    border-radius: 24px;
    box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
    overflow: hidden;
    transition: all 0.3s var(--transition-bounce);
    position: relative;
    display: flex;
    flex-direction: column;
    border: 1px solid var(--border-color);
    
    &:hover {
      transform: translateY(-8px);
      box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                  0 10px 10px -5px rgba(0, 0, 0, 0.04);
    }
    
    .user-header {
      padding: 1.5rem 1.5rem 0.5rem;
      display: flex;
      flex-direction: column;
      justify-content: center;
      
      .user-badges {
        display: flex;
        flex-wrap: wrap;
        gap: 0.5rem;
        margin-bottom: 1rem;
        
        .role-badge,
        .status-badge {
          display: inline-flex;
          align-items: center;
          padding: 0.35rem 0.75rem;
          border-radius: 99px;
          font-size: 0.8rem;
          font-weight: 500;
          
          &::before {
            content: "";
            display: inline-block;
            width: 8px;
            height: 8px;
            border-radius: 50%;
            margin-right: 0.5rem;
          }
        }
        
        .role-badge {
          &.admin { 
            background: rgba(79, 70, 229, 0.1); 
            color: #4f46e5;
            &::before { background-color: #4f46e5; }
          }
          &.assistant { 
            background: rgba(249, 115, 22, 0.1); 
            color: #f97316; 
            &::before { background-color: #f97316; }
          }
          &.employee { 
            background: rgba(16, 185, 129, 0.1); 
            color: #10b981; 
            &::before { background-color: #10b981; }
          }
        }
        
        .status-badge {
          &.active { 
            background: rgba(16, 185, 129, 0.1); 
            color: #10b981; 
            &::before { background-color: #10b981; }
          }
          &.inactive { 
            background: rgba(75, 85, 99, 0.1); 
            color: #4b5563; 
            &::before { background-color: #4b5563; }
          }
        }
      }
      
      .user-title {
        font-size: 1.25rem;
        font-weight: 600;
        color: var(--text-primary);
        margin: 0 0 1rem 0;
        line-height: 1.4;
      }
    }
    
    .user-body {
      padding: 0 1.5rem 1.5rem;
      flex-grow: 1;
      
      .user-email {
        color: var(--text-secondary);
        margin-bottom: 1rem;
        font-size: 0.95rem;
      }
      
      .user-detail {
        display: flex;
        justify-content: space-between;
        margin-bottom: 0.5rem;
        font-size: 0.9rem;
        
        .detail-label {
          font-weight: 500;
          color: var(--text-secondary);
        }
      }
    }
    
    .user-meta {
      padding: 1rem 1.5rem;
      background-color: var(--bg-tertiary);
      border-top: 1px solid var(--border-color);
      
      .meta-item {
        display: flex;
        align-items: center;
        margin-bottom: 0.5rem;
        
        &:last-child {
          margin-bottom: 0;
        }
        
        i {
          width: 24px;
          height: 24px;
          background-color: rgba(99, 102, 241, 0.1);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          margin-right: 0.75rem;
          color: var(--primary-color);
          font-size: 0.9rem;
        }
        
        span {
          font-size: 0.9rem;
          color: var(--text-secondary);
          white-space: nowrap;
          overflow: hidden;
          text-overflow: ellipsis;
        }
      }
    }
    
    .user-actions {
      padding: 1rem 1.5rem;
      border-top: 1px solid var(--border-color);
      display: flex;
      justify-content: flex-end;
      gap: 0.5rem;
      
      .action-btn {
        display: inline-flex;
        align-items: center;
        justify-content: center;
        width: 36px;
        height: 36px;
        border-radius: 50%;
        border: none;
        background-color: var(--bg-tertiary);
        color: var(--text-primary);
        cursor: pointer;
        transition: all 0.2s ease;
        
        &:hover {
          background-color: var(--hover-bg);
          transform: translateY(-2px);
        }
        
        &.warning {
          color: var(--warning-color);
          
          &:hover {
            background-color: rgba(var(--warning-rgb, 255, 152, 0), 0.15);
          }
        }
        
        &.success {
          color: var(--success-color);
          
          &:hover {
            background-color: rgba(var(--success-rgb, 76, 175, 80), 0.15);
          }
        }
        
        &.danger {
          color: var(--danger-color);
          
          &:hover {
            background-color: rgba(var(--danger-rgb, 244, 67, 54), 0.15);
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
    
    .filters-container .filter-group .filter-options {
      flex-direction: column;
      
      .filter-btn {
        width: 100%;
        text-align: center;
      }
    }
    
    .users-grid {
      grid-template-columns: 1fr;
    }
  }
}
</style>
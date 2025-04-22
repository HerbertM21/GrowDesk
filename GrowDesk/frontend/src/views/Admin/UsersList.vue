<template>
  <div class="admin-section">
    <div class="admin-section-header">
      <div>
        <h1 class="admin-section-title">Administración de Usuarios</h1>
        <p class="admin-section-description">
          Gestione todos los usuarios del sistema, asigne roles y controle su estado.
        </p>
      </div>
      
      <div class="admin-section-actions">
        <button class="btn btn-primary" @click="showCreateModal = true">
          <i class="pi pi-plus"></i> Nuevo Usuario
        </button>
      </div>
    </div>
    
    <div v-if="loading" class="loading-container">
      <div class="loading">
        <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
        <p>Cargando usuarios...</p>
      </div>
    </div>
    
    <div v-else-if="error" class="alert alert-danger">
      <i class="pi pi-exclamation-circle"></i> {{ error }}
    </div>
    
    <div v-else class="users-content">
      <div class="filter-controls">
        <div class="filter-header">
          <i class="pi pi-filter"></i>
          <h3>Filtros y búsqueda</h3>
        </div>
        
        <div class="filter-body">
          <div class="filter-group">
            <label class="form-label">Filtrar por rol:</label>
            <div class="select-container">
              <select v-model="roleFilter" class="form-control">
                <option value="all">Todos</option>
                <option value="admin">Administrador</option>
                <option value="assistant">Asistente</option>
                <option value="employee">Empleado</option>
              </select>
              <i class="pi pi-chevron-down"></i>
            </div>
          </div>
          
          <div class="filter-group">
            <label class="form-label">Mostrar:</label>
            <div class="select-container">
              <select v-model="statusFilter" class="form-control">
                <option value="all">Todos</option>
                <option value="active">Activos</option>
                <option value="inactive">Inactivos</option>
              </select>
              <i class="pi pi-chevron-down"></i>
            </div>
          </div>
          
          <div class="filter-group search-box">
            <label class="form-label">Buscar:</label>
            <div class="search-input-container">
              <input v-model="searchQuery" class="form-control" placeholder="Buscar por nombre o email..." />
              <i class="pi pi-search search-icon"></i>
            </div>
          </div>
        </div>
      </div>
      
      <div v-if="filteredUsers.length === 0" class="empty-state">
        <i class="pi pi-users" style="font-size: 3rem; color: var(--text-muted)"></i>
        <p v-if="searchQuery || roleFilter !== 'all' || statusFilter !== 'all'">
          No se encontraron usuarios con los filtros seleccionados.
        </p>
        <p v-else>No hay usuarios registrados en el sistema.</p>
      </div>
      
      <div v-else class="users-grid">
        <div 
          v-for="user in filteredUsers" 
          :key="user.id" 
          class="user-card"
        >
          <div class="user-card-header">
            <div class="avatar" :style="{ backgroundColor: getAvatarColor(user) }">
              {{ getInitials(user) }}
            </div>
            <div class="user-info">
              <h3>{{ user.firstName }} {{ user.lastName }}</h3>
              <p class="user-email">{{ user.email }}</p>
            </div>
          </div>
          
          <div class="user-card-body">
            <div class="user-details">
              <div class="detail-row">
                <span class="detail-label">Rol:</span>
                <span :class="['role-badge', user.role]">
                  {{ translateRole(user.role) }}
                </span>
              </div>
              <div class="detail-row">
                <span class="detail-label">Estado:</span>
                <span :class="['status-badge', user.active ? 'active' : 'inactive']">
                  {{ user.active ? 'Activo' : 'Inactivo' }}
                </span>
              </div>
              <div class="detail-row" v-if="user.department">
                <span class="detail-label">Departamento:</span>
                <span>{{ user.department }}</span>
              </div>
              <div class="detail-row" v-if="user.createdAt">
                <span class="detail-label">Creado:</span>
                <span>{{ formatDate(user.createdAt) }}</span>
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

// Stores
const userStore = useUsersStore();
const authStore = useAuthStore();
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

// Cargar usuarios al montar el componente
onMounted(() => {
  userStore.fetchUsers();
  
  // Para modo de desarrollo con el mock server
  if (process.env.NODE_ENV === 'development' && users.value.length === 0) {
    userStore.initMockUsers();
  }
});
</script>

<style scoped lang="scss">
.admin-section {
  padding: 1rem;
  
  .admin-section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 2rem;
    
    .admin-section-title {
      font-size: 1.75rem;
      font-weight: 600;
      color: var(--text-primary);
      margin: 0;
    }
    
    .admin-section-description {
      margin: 0.5rem 0 0;
      color: var(--text-secondary);
      font-size: 1rem;
    }
    
    .btn-primary {
      display: flex;
      align-items: center;
      gap: 0.5rem;
      padding: 0.75rem 1.25rem;
      background-color: var(--primary-color);
      color: white;
      border: none;
      border-radius: 4px;
      font-weight: 500;
      cursor: pointer;
      transition: background-color 0.2s;
      
      &:hover {
        background-color: var(--primary-color-dark, #8a15c0);
      }
      
      i {
        font-size: 1rem;
      }
    }
  }
}

.loading-container {
  display: flex;
  justify-content: center;
  padding: 3rem 0;
  
  .loading {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    color: var(--text-secondary);
    
    p {
      margin-top: 1rem;
    }
  }
}

.empty-state {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  text-align: center;
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  border-radius: var(--border-radius);
  
  p {
    margin-top: 1rem;
    font-size: 1.1rem;
  }
}

.users-content {
  display: grid;
  grid-template-columns: 1fr;
  gap: 1.5rem;
}

.filter-controls {
  background-color: var(--bg-secondary, white);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  
  .filter-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 1rem 1.5rem;
    background-color: var(--primary-color-light, #e3f2fd);
    border-bottom: 1px solid var(--border-color, #eee);
    
    i {
      color: var(--primary-color);
      font-size: 1.25rem;
    }
    
    h3 {
      margin: 0;
      font-size: 1.1rem;
      font-weight: 500;
      color: var(--text-primary);
    }
  }
  
  .filter-body {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
    gap: 1.5rem;
    padding: 1.5rem;
  }
}

.filter-group {
  .form-label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-secondary);
  }
  
  .select-container {
    position: relative;
    
    select {
      width: 100%;
      padding: 0.75rem;
      padding-right: 2.5rem;
      border: 1px solid var(--border-color);
      border-radius: 4px;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      appearance: none;
      cursor: pointer;
      
      &:focus {
        outline: none;
        border-color: var(--primary-color);
        box-shadow: 0 0 0 2px var(--primary-color-light, rgba(25, 118, 210, 0.2));
      }
    }
    
    i {
      position: absolute;
      right: 10px;
      top: 50%;
      transform: translateY(-50%);
      color: var(--text-secondary);
      pointer-events: none;
    }
  }
  
  .search-input-container {
    position: relative;
    
    input {
      width: 100%;
      padding: 0.75rem;
      padding-right: 2.5rem;
      border: 1px solid var(--border-color);
      border-radius: 4px;
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
      
      &:focus {
        outline: none;
        border-color: var(--primary-color);
        box-shadow: 0 0 0 2px var(--primary-color-light, rgba(25, 118, 210, 0.2));
      }
    }
    
    .search-icon {
      position: absolute;
      right: 10px;
      top: 50%;
      transform: translateY(-50%);
      color: var(--text-secondary);
    }
  }
}

.users-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 1.5rem;
}

.user-card {
  background-color: var(--bg-secondary, white);
  border-radius: 8px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.1);
  overflow: hidden;
  transition: transform 0.2s, box-shadow 0.2s;
  
  &:hover {
    transform: translateY(-3px);
    box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
  }
  
  .user-card-header {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.25rem;
    background-color: var(--primary-color-light, #e3f2fd);
    
    .avatar {
      width: 48px;
      height: 48px;
      border-radius: 50%;
      background-color: var(--primary-color);
      color: white;
      display: flex;
      align-items: center;
      justify-content: center;
      font-weight: 600;
      font-size: 1.1rem;
    }
    
    .user-info {
      h3 {
        margin: 0;
        font-size: 1.1rem;
        color: var(--text-primary);
        font-weight: 600;
      }
      
      .user-email {
        margin: 0.25rem 0 0;
        font-size: 0.9rem;
        color: var(--text-secondary);
      }
    }
  }
  
  .user-card-body {
    padding: 1.25rem;
  }
  
  .user-details {
    .detail-row {
      display: flex;
      justify-content: space-between;
      margin-bottom: 0.75rem;
      
      &:last-child {
        margin-bottom: 0;
      }
      
      .detail-label {
        font-weight: 500;
        color: var(--text-secondary);
        font-size: 0.9rem;
      }
    }
  }
  
  .user-actions {
    display: flex;
    justify-content: flex-end;
    gap: 0.5rem;
    margin-top: 1.25rem;
    padding-top: 1.25rem;
    border-top: 1px solid var(--border-color, #eee);
  }
  
  .action-btn {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: none;
    background-color: var(--bg-tertiary);
    color: var(--text-primary);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    transition: background-color 0.2s, color 0.2s;
    
    &:hover {
      background-color: var(--bg-hover);
    }
    
    &.warning {
      color: var(--warning-color);
      
      &:hover {
        background-color: var(--warning-color-light);
      }
    }
    
    &.success {
      color: var(--success-color);
      
      &:hover {
        background-color: var(--success-color-light);
      }
    }
    
    &.danger {
      color: var(--error-color);
      
      &:hover {
        background-color: var(--error-color-light);
      }
    }
  }
}

.role-badge,
.status-badge {
  display: inline-block;
  padding: 0.35em 0.65em;
  font-size: 0.75em;
  font-weight: 700;
  line-height: 1;
  text-align: center;
  white-space: nowrap;
  vertical-align: baseline;
  border-radius: 0.25rem;
}

.role-badge {
  &.admin {
    background-color: var(--primary-color);
    color: white;
  }
  
  &.assistant {
    background-color: var(--secondary-color);
    color: white;
  }
  
  &.employee {
    background-color: var(--text-secondary);
    color: white;
  }
}

.status-badge {
  &.active {
    background-color: var(--success-color);
    color: white;
  }
  
  &.inactive {
    background-color: var(--text-muted);
    color: white;
  }
}

// Responsive adjustments
@media (max-width: 768px) {
  .admin-section-header {
    flex-direction: column;
    align-items: flex-start;
    
    .admin-section-actions {
      margin-top: 1rem;
      width: 100%;
      
      .btn-primary {
        width: 100%;
        justify-content: center;
      }
    }
  }
  
  .users-grid {
    grid-template-columns: 1fr;
  }
}
</style> 
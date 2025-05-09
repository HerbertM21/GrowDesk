import { defineStore } from 'pinia';
import apiClient from '@/api/client';
import { ref } from 'vue';

// Nombre clave para localStorage
const STORAGE_KEY = 'growdesk-users';

// Interfaces
export interface User {
  id: string;
  email: string;
  firstName: string;
  lastName: string;
  role: 'admin' | 'assistant' | 'employee';
  department: string | null;
  active: boolean;
  createdAt: string;
  updatedAt: string;
  // Campos adicionales de perfil
  position?: string | null;
  phone?: string | null;
  language?: string;
}

interface UsersState {
  users: User[];
  currentProfile: User | null;
  loading: boolean;
  error: string | null;
}

// Datos mock para desarrollo
const mockUsers: User[] = [
  {
    id: '1',
    email: 'admin@example.com',
    firstName: 'Herbert',
    lastName: 'Usuario',
    role: 'admin',
    department: 'Tecnología',
    active: true,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    position: 'Gerente de TI',
    phone: '+569 1234 5678',
    language: 'es'
  },
  {
    id: '2',
    email: 'asistente@example.com',
    firstName: 'Asistente',
    lastName: 'Soporte',
    role: 'assistant',
    department: 'Soporte',
    active: true,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    position: 'Coordinador de Soporte',
    phone: '+34 600 234 567',
    language: 'es'
  },
  {
    id: '3',
    email: 'empleado@example.com',
    firstName: 'Empleado',
    lastName: 'Regular',
    role: 'employee',
    department: 'Ventas',
    active: true,
    createdAt: new Date().toISOString(),
    updatedAt: new Date().toISOString(),
    position: 'Representante de Ventas',
    phone: '+34 600 345 678',
    language: 'en'
  }
];

export const useUsersStore = defineStore('users', () => {
  // Estado
  const users = ref<User[]>([]);
  const currentProfile = ref<User | null>(null);
  const loading = ref<boolean>(false);
  const error = ref<string | null>(null);

  // Guardar usuarios en localStorage
  function saveUsersToLocalStorage() {
    console.log('Guardando usuarios en localStorage:', users.value);
    localStorage.setItem(STORAGE_KEY, JSON.stringify(users.value));
  }

  // Cargar usuarios desde localStorage
  function loadUsersFromLocalStorage(): User[] {
    try {
      const storedData = localStorage.getItem(STORAGE_KEY);
      if (storedData) {
        console.log('Usuarios cargados desde localStorage');
        return JSON.parse(storedData);
      }
    } catch (err) {
      console.error('Error al cargar usuarios desde localStorage:', err);
    }
    console.log('No se encontraron usuarios en localStorage, usando datos iniciales');
    return [];
  }

  // Getters
  const getUsersByRole = (roleFilter: string) => {
    return users.value.filter((user: User) => user.role === roleFilter);
  };

  const getActiveUsers = () => {
    return users.value.filter((user: User) => user.active);
  };

  const getInactiveUsers = () => {
    return users.value.filter((user: User) => !user.active);
  };

  // Acciones
  const fetchUsers = async () => {
    if (loading.value) return;
    
    loading.value = true;
    error.value = null;
    
    try {
      console.log('Obteniendo usuarios...');
      
      // Intentar cargar desde localStorage primero
      const storedUsers = loadUsersFromLocalStorage();
      
      if (storedUsers.length > 0) {
        users.value = storedUsers;
        console.log('Usuarios cargados desde localStorage:', users.value);
      } else {
        // Intentar cargar desde la API
        try {
          const response = await apiClient.get('/users');
          users.value = response.data;
          // Guardar en localStorage
          saveUsersToLocalStorage();
        } catch (apiErr) {
          console.log('Error al cargar desde API, cargando datos mock');
          initMockUsers();
        }
      }
    } catch (err: unknown) {
      error.value = 'Error al cargar usuarios';
      console.error('Error al cargar usuarios:', err);
      
      if (import.meta.env.DEV) {
        console.log('Cargando datos mock de usuarios');
        initMockUsers();
      }
    } finally {
      loading.value = false;
    }
  };

  const fetchUserProfile = async (userId: string) => {
    if (loading.value) return null;
    
    loading.value = true;
    error.value = null;
    
    try {
      console.log(`Obteniendo perfil de usuario para ID: ${userId}`);
      
      // Primero, buscar en el array de usuarios
      const userFromArray = users.value.find((u: User) => u.id === userId);
      if (userFromArray) {
        currentProfile.value = userFromArray;
        return currentProfile.value;
      }
      
      // Si no se encuentra, intentar obtenerlo de la API
      const response = await apiClient.get(`/users/${userId}`);
      currentProfile.value = response.data;
      return currentProfile.value;
    } catch (err: unknown) {
      error.value = 'Error al cargar perfil de usuario';
      console.error('Error al cargar perfil de usuario:', err);
      return null;
    } finally {
      loading.value = false;
    }
  };

  const createUser = async (userData: Partial<User>) => {
    loading.value = true;
    error.value = null;

    try {
      const response = await apiClient.post('/users', userData);
      // Añadir a la lista local
      users.value.push(response.data);
      // Guardar en localStorage
      saveUsersToLocalStorage();
      return response.data;
    } catch (err: unknown) {
      console.error('Error al crear usuario:', err);
      error.value = 'Error al crear el usuario';
      
      // Para desarrollo, simular creación
      const newUser: User = {
        id: `${users.value.length + 1}`,
        firstName: userData.firstName || '',
        lastName: userData.lastName || '',
        email: userData.email || '',
        role: userData.role || 'employee',
        department: userData.department || null,
        active: userData.active !== undefined ? userData.active : true,
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString()
      };
      
      // Añadir a la lista local
      users.value.push(newUser);
      // Guardar en localStorage
      saveUsersToLocalStorage();
      return newUser;
    } finally {
      loading.value = false;
    }
  };

  const updateUser = async (userId: string, userData: Partial<User>) => {
    loading.value = true;
    error.value = null;

    try {
      console.log('Actualizando usuario con ID:', userId);
      console.log('Datos enviados:', userData);
      
      const response = await apiClient.put(`/users/${userId}`, userData);
      
      // Actualizar en el array local
      const index = users.value.findIndex((u: User) => u.id === userId);
      if (index !== -1) {
        users.value[index] = { ...users.value[index], ...response.data };
      }
      
      // Si es el perfil actual, actualizar también
      if (currentProfile.value?.id === userId) {
        currentProfile.value = { ...currentProfile.value, ...response.data };
      }
      
      // Guardar en localStorage
      saveUsersToLocalStorage();
      
      return users.value[index];
    } catch (err: unknown) {
      console.error('Error al actualizar usuario:', err);
      error.value = 'Error al actualizar el usuario';
      
      // Para desarrollo, simular actualización
      if (import.meta.env.DEV) {
        const index = users.value.findIndex((u: User) => u.id === userId);
        if (index !== -1) {
          users.value[index] = { 
            ...users.value[index],
            ...userData,
            updatedAt: new Date().toISOString()
          };
          
          // Si es el perfil actual, actualizar también
          if (currentProfile.value?.id === userId) {
            currentProfile.value = { 
              ...currentProfile.value, 
              ...userData,
              updatedAt: new Date().toISOString()
            };
          }
          
          // Guardar en localStorage
          saveUsersToLocalStorage();
          
          console.log('Usuario actualizado localmente:', users.value[index]);
          return users.value[index];
        }
      }
      
      return null;
    } finally {
      loading.value = false;
    }
  };

  const changeUserRole = async (userId: string, newRole: 'admin' | 'assistant' | 'employee') => {
    return updateUser(userId, { role: newRole });
  };
  
  const toggleUserActive = async (userId: string) => {
    const userToToggle = users.value.find((u: User) => u.id === userId);
    if (userToToggle) {
      return updateUser(userId, { active: !userToToggle.active });
    }
    return null;
  };

  const deleteUser = async (userId: string) => {
    loading.value = true;
    error.value = null;

    try {
      await apiClient.delete(`/users/${userId}`);
      
      // Eliminar del array local
      users.value = users.value.filter((u: User) => u.id !== userId);
      
      // Si es el perfil actual, limpiar
      if (currentProfile.value?.id === userId) {
        currentProfile.value = null;
      }
      
      // Guardar en localStorage
      saveUsersToLocalStorage();
      
      return true;
    } catch (err: unknown) {
      console.error('Error al eliminar usuario:', err);
      error.value = 'Error al eliminar el usuario';
      
      // Para desarrollo, simular eliminación
      if (import.meta.env.DEV) {
        // Eliminar del array local
        users.value = users.value.filter((u: User) => u.id !== userId);
        
        // Si es el perfil actual, limpiar
        if (currentProfile.value?.id === userId) {
          currentProfile.value = null;
        }
        
        // Guardar en localStorage
        saveUsersToLocalStorage();
        
        console.log('Usuario eliminado localmente');
        return true;
      }
      
      return false;
    } finally {
      loading.value = false;
    }
  };

  // Para desarrollo rápido - inicializar con usuarios de ejemplo
  const initMockUsers = () => {
    // Cargar desde localStorage primero
    const storedUsers = loadUsersFromLocalStorage();
    
    if (storedUsers.length > 0) {
      users.value = storedUsers;
      console.log('Usuarios cargados desde localStorage:', users.value);
      return;
    }
    
    // Si no hay datos en localStorage, usar datos mock
    if (users.value.length === 0) {
      users.value = [...mockUsers];
      // Guardar en localStorage
      saveUsersToLocalStorage();
      console.log('Usuarios mock inicializados y guardados en localStorage:', users.value);
    }
  };

  // Obtener perfil del usuario actual (desde auth)
  const fetchCurrentUserProfile = async () => {
    if (loading.value) return null;
    
    loading.value = true;
    error.value = null;
    
    try {
      console.log('Obteniendo perfil del usuario actual');
      const response = await apiClient.get('/auth/me');
      currentProfile.value = response.data;
      return currentProfile.value;
    } catch (err: unknown) {
      error.value = 'Error al cargar perfil del usuario actual';
      console.error('Error al cargar perfil del usuario actual:', err);
      
      // En modo desarrollo, asegurarse de tener usuarios mock cargados
      if (import.meta.env.DEV) {
        console.log('Usando perfil mock en desarrollo');
        
        // Inicializar usuarios mock si no hay ninguno cargado
        if (users.value.length === 0) {
          console.log('Cargando usuarios mock primero');
          initMockUsers();
        }
        
        if (users.value.length > 0) {
          currentProfile.value = users.value[0];
          console.log('Perfil mock cargado:', currentProfile.value);
          return currentProfile.value;
        }
      }
      
      return null;
    } finally {
      loading.value = false;
    }
  };
  
  // Alias para hacer más consistente la API con el auth store
  const getCurrentUser = async (): Promise<User | null> => {
    if (currentProfile.value) {
      return currentProfile.value;
    }
    return fetchCurrentUserProfile();
  };

  // Obtener un usuario por su ID
  const fetchUser = async (userId: string): Promise<User | null> => {
    if (loading.value) return null;
    
    loading.value = true;
    error.value = null;
    
    try {
      console.log(`Obteniendo usuario con ID: ${userId}`);
      
      // Si hay usuarios cargados, intentar encontrarlo en la caché primero
      if (users.value.length > 0) {
        const cachedUser = users.value.find((u: User) => u.id === userId);
        if (cachedUser) {
          console.log('Usuario encontrado en caché:', cachedUser);
          return cachedUser;
        }
      }
      
      // Si no está en caché o no hay usuarios cargados, hacer petición al API
      const response = await apiClient.get(`/users/${userId}`);
      const userData = response.data;
      
      // Actualizar en la caché local si ya existe
      const index = users.value.findIndex((u: User) => u.id === userId);
      if (index !== -1) {
        users.value[index] = userData;
      } else {
        // Si no existe, añadirlo a la lista de usuarios
        users.value.push(userData);
      }
      
      // Guardar en localStorage
      saveUsersToLocalStorage();
      
      return userData;
    } catch (err) {
      console.error(`Error al obtener usuario con ID ${userId}:`, err);
      error.value = 'Error al obtener usuario';
      
      // En desarrollo, buscar en los usuarios mock
      if (import.meta.env.DEV) {
        if (users.value.length === 0) {
          initMockUsers();
        }
        
        const mockUser = users.value.find((u: User) => u.id === userId);
        if (mockUser) {
          console.log('Usando usuario mock:', mockUser);
          return mockUser;
        }
        
        // Si estamos buscando el usuario actual en localStorage
        if (userId === localStorage.getItem('userId')) {
          console.log('Creando usuario actual mock');
          return {
            id: userId,
            firstName: 'Usuario',
            lastName: 'Actual',
            email: 'current@example.com',
            role: 'employee',
            department: 'General',
            active: true,
            createdAt: new Date().toISOString(),
            updatedAt: new Date().toISOString()
          };
        }
      }
      
      return null;
    } finally {
      loading.value = false;
    }
  };

  return {
    // Estado
    users,
    currentProfile,
    loading,
    error,

    // Getters
    getUsersByRole,
    getActiveUsers,
    getInactiveUsers,

    // Acciones
    fetchUsers,
    fetchUser,
    fetchUserProfile,
    fetchCurrentUserProfile,
    getCurrentUser,
    createUser,
    updateUser,
    changeUserRole,
    toggleUserActive,
    deleteUser,
    initMockUsers,
    saveUsersToLocalStorage,
    loadUsersFromLocalStorage
  };
});

// Alias para mantener compatibilidad con código existente
export const useUserStore = useUsersStore; 
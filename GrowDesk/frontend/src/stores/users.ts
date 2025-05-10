import { defineStore } from 'pinia';
import apiClient from '@/api/client';
import { ref } from 'vue';

// Nombre clave para localStorage
const STORAGE_KEY = 'growdesk_users';

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
    
    // Sincronizar con el backend
    syncUsersWithBackend();
  }

  // Función para sincronizar los usuarios con el backend
  async function syncUsersWithBackend() {
    try {
      // Paso adicional: Limpiar localStorage antes de sincronizar
      const storedData = localStorage.getItem(STORAGE_KEY);
      if (storedData) {
        let parsedData;
        try {
          parsedData = JSON.parse(storedData);
          
          // Validar que es un array
          if (!Array.isArray(parsedData)) {
            console.error('Datos en localStorage no son un array, reiniciando');
            localStorage.setItem(STORAGE_KEY, JSON.stringify([]));
            users.value = [];
            return; // No sincronizar datos inválidos
          }
          
          // Filtrar elementos no válidos
          const validItems = parsedData.filter(item => 
            typeof item === 'object' && 
            item !== null &&
            !(typeof item === 'string' && item.includes('page not found'))
          );
          
          if (validItems.length !== parsedData.length) {
            console.warn(`Limpiando ${parsedData.length - validItems.length} elementos inválidos antes de sincronizar`);
            localStorage.setItem(STORAGE_KEY, JSON.stringify(validItems));
            parsedData = validItems;
          }
        } catch (e) {
          console.error('Error en datos almacenados, reiniciando localStorage:', e);
          localStorage.setItem(STORAGE_KEY, JSON.stringify([]));
          users.value = [];
          return; // No sincronizar datos inválidos
        }
      }

      // Asegurar que los datos a enviar son válidos
      const validUsers = users.value.filter(user => 
        typeof user === 'object' && 
        user !== null &&
        'id' in user && 
        'email' in user && 
        'firstName' in user && 
        'lastName' in user && 
        'role' in user
      );
      
      // Si hay usuarios inválidos, actualizamos el array local
      if (validUsers.length !== users.value.length) {
        console.warn(`Se filtraron ${users.value.length - validUsers.length} usuarios inválidos antes de sincronizar`);
        users.value = validUsers;
        // Guardar array filtrado en localStorage
        localStorage.setItem(STORAGE_KEY, JSON.stringify(validUsers));
      }
      
      // URL del servidor de sincronización (asegúrate de que apunta al contenedor)
      const syncApiUrl = import.meta.env.VITE_SYNC_API_URL || 'http://localhost:8000/api/sync/users';
      
      console.log('Sincronizando usuarios con el backend en:', syncApiUrl);
      console.log('Datos a enviar:', validUsers);
      
      // Usar fetch directamente para tener más control sobre la solicitud
      const response = await fetch(syncApiUrl, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
          'Accept': 'application/json'
        },
        body: JSON.stringify(validUsers),
        // Agregar timeout para evitar bloqueos largos
        signal: AbortSignal.timeout(5000)
      });
      
      if (!response.ok) {
        throw new Error(`Error de sincronización: ${response.status} ${response.statusText}`);
      }
      
      const data = await response.json();
      console.log('Usuarios sincronizados con el backend:', data);
    } catch (error) {
      console.error('Error al sincronizar usuarios con el backend:', error);
      
      // No es crítico, ya que los datos ya están en localStorage
      // En producción, podríamos implementar un mecanismo de reintentos
      if (import.meta.env.DEV) {
        console.log('Funcionando en modo local sin sincronización con el backend');
      }
    }
  }

  // Función para cargar usuarios desde localStorage
  function loadUsersFromLocalStorage(): User[] {
    try {
      const storedData = localStorage.getItem(STORAGE_KEY);
      if (storedData) {
        let parsedData;
        try {
          parsedData = JSON.parse(storedData);

          // Validar que es un array
          if (!Array.isArray(parsedData)) {
            console.error('Datos en localStorage no son un array, reiniciando');
            users.value = [];
            return [];
          }

          // Filtrar elementos no válidos (404 page not found u otros strings/nulls)
          const validUsers = parsedData.filter(item => 
            typeof item === 'object' && 
            item !== null &&
            'id' in item &&
            'email' in item &&
            'firstName' in item &&
            'lastName' in item &&
            'role' in item &&
            !(typeof item === 'string' && item.includes('page not found'))
          );

          if (validUsers.length !== parsedData.length) {
            console.warn(`Se filtraron ${parsedData.length - validUsers.length} elementos inválidos del localStorage`);
            localStorage.setItem(STORAGE_KEY, JSON.stringify(validUsers));
          }

          // Asignar únicamente usuarios válidos
          users.value = validUsers;
          console.log('Usuarios cargados desde localStorage:', validUsers.length);
          return validUsers;
        } catch (e) {
          console.error('Error al analizar localStorage, reiniciando:', e);
          users.value = [];
          localStorage.setItem(STORAGE_KEY, JSON.stringify([]));
          return [];
        }
      } else {
        console.log('No hay usuarios en localStorage');
        users.value = [];
        return [];
      }
    } catch (error) {
      console.error('Error al cargar usuarios:', error);
      // En caso de error, inicializar con array vacío
      users.value = [];
      return [];
    }
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
      
      // Buscar en el array de usuarios local
      const userFromArray = users.value.find((u: User) => u.id === userId);
      if (userFromArray) {
        console.log('Usuario encontrado en caché local:', userFromArray);
        currentProfile.value = userFromArray;
        return currentProfile.value;
      }
      
      // Si no hay usuarios cargados, cargar desde localStorage
      if (users.value.length === 0) {
        loadUsersFromLocalStorage();
        
        // Intentar buscar de nuevo
        const userFromStorage = users.value.find((u: User) => u.id === userId);
        if (userFromStorage) {
          console.log('Usuario encontrado en localStorage:', userFromStorage);
          currentProfile.value = userFromStorage;
          return currentProfile.value;
        }
      }
      
      console.log(`Usuario con ID ${userId} no encontrado localmente`);
      error.value = 'Usuario no encontrado';
      return null;
    } catch (err: unknown) {
      error.value = 'Error al cargar perfil de usuario';
      console.error('Error al cargar perfil de usuario:', err);
      return null;
    } finally {
      loading.value = false;
    }
  };

  // Crea un nuevo usuario
  async function createUser(userData: Partial<User>) {
    if (!userData.firstName || !userData.lastName || !userData.email) {
      console.error('Datos incompletos para crear usuario');
      return null;
    }

    try {
      // SIEMPRE crear localmente sin importar el modo
      console.log('Creando usuario localmente:', userData);
      
      const newUser: User = {
        id: 'user_' + Date.now().toString(),
        firstName: userData.firstName,
        lastName: userData.lastName,
        email: userData.email,
        role: userData.role || 'employee',
        active: userData.active !== undefined ? userData.active : true,
        department: userData.department || null,
        position: userData.position || null,
        phone: userData.phone || null,
        language: userData.language || 'es',
        createdAt: new Date().toISOString(),
        updatedAt: new Date().toISOString(),
      };
      
      // Guardar en memoria
      users.value.push(newUser);
      
      // Guardar en localStorage
      saveUsersToLocalStorage();
      
      // Notificar éxito
      console.log(`✅ Usuario creado: ${newUser.firstName} ${newUser.lastName}`);
      
      // Intentar sincronizar con el backend (no es crítico si falla)
      try {
        await syncUsersWithBackend();
      } catch (syncError) {
        console.warn('No se pudo sincronizar con el backend, pero el usuario se guardó localmente', syncError);
      }
      
      return newUser;
    } catch (error) {
      console.error('Error al crear usuario:', error);
      return null;
    }
  }

  const updateUser = async (userId: string, userData: Partial<User>) => {
    loading.value = true;
    error.value = null;

    try {
      // Buscar el usuario en el array local
      const index = users.value.findIndex((u: User) => u.id === userId);
      
      if (index === -1) {
        console.error(`Usuario con ID ${userId} no encontrado`);
        return null;
      }
      
      console.log('Actualizando usuario localmente:', userId);
      console.log('Datos a actualizar:', userData);
      
      // Actualizar en el array local
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
      
      // Notificar éxito
      console.log(`✅ Usuario actualizado: ${users.value[index].firstName} ${users.value[index].lastName}`);
      
      // Intentar sincronizar con el backend (no es crítico si falla)
      try {
        await syncUsersWithBackend();
      } catch (syncError) {
        console.warn('No se pudo sincronizar con el backend, pero el usuario se actualizó localmente', syncError);
      }
      
      return users.value[index];
    } catch (err: unknown) {
      console.error('Error al actualizar usuario:', err);
      error.value = 'Error al actualizar el usuario';
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
      console.log('Eliminando usuario localmente:', userId);
      
      // Eliminar del array local
      const userToDelete = users.value.find((u: User) => u.id === userId);
      if (!userToDelete) {
        console.error(`Usuario con ID ${userId} no encontrado`);
        return false;
      }
      
      // Eliminar del array local
      users.value = users.value.filter((u: User) => u.id !== userId);
      
      // Si es el perfil actual, limpiar
      if (currentProfile.value?.id === userId) {
        currentProfile.value = null;
      }
      
      // Guardar en localStorage
      saveUsersToLocalStorage();
      
      // Notificar éxito
      console.log(`✅ Usuario eliminado: ${userToDelete.firstName} ${userToDelete.lastName}`);
      
      // Intentar sincronizar con el backend (no es crítico si falla)
      try {
        await syncUsersWithBackend();
      } catch (syncError) {
        console.warn('No se pudo sincronizar con el backend, pero el usuario se eliminó localmente', syncError);
      }
      
      return true;
    } catch (err: unknown) {
      console.error('Error al eliminar usuario:', err);
      error.value = 'Error al eliminar el usuario';
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
      
      // Intentar obtener userId del localStorage
      const userId = localStorage.getItem('userId');
      
      if (userId) {
        console.log('ID de usuario actual encontrado en localStorage:', userId);
        
        // Buscar en el array de usuarios
        if (users.value.length > 0) {
          const userFromArray = users.value.find((u: User) => u.id === userId);
          if (userFromArray) {
            console.log('Usuario actual encontrado en caché:', userFromArray);
            currentProfile.value = userFromArray;
            return currentProfile.value;
          }
        }
        
        // Si no hay usuarios cargados, cargar desde localStorage
        if (users.value.length === 0) {
          loadUsersFromLocalStorage();
          
          // Intentar buscar de nuevo
          const userFromStorage = users.value.find((u: User) => u.id === userId);
          if (userFromStorage) {
            console.log('Usuario actual encontrado en localStorage:', userFromStorage);
            currentProfile.value = userFromStorage;
            return currentProfile.value;
          }
        }
        
        // Si no se encuentra, pero tenemos el ID, crear un perfil mock
        console.log('Creando perfil mock para el usuario actual');
        const mockUser: User = {
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
        
        // Guardar en memoria y localStorage
        users.value.push(mockUser);
        saveUsersToLocalStorage();
        
        currentProfile.value = mockUser;
        return currentProfile.value;
      }
      
      // Si no hay userId en localStorage pero tenemos usuarios
      if (users.value.length === 0) {
        loadUsersFromLocalStorage();
      }
      
      // Si hay usuarios, usar el primero como perfil actual
      if (users.value.length > 0) {
        console.log('Usando el primer usuario como perfil actual');
        currentProfile.value = users.value[0];
        
        // Guardar en localStorage para referencia futura
        localStorage.setItem('userId', currentProfile.value.id);
        
        return currentProfile.value;
      }
      
      // Si aún no hay usuarios, inicializar con mock
      console.log('No hay usuarios, inicializando con datos mock');
      initMockUsers();
      
      if (users.value.length > 0) {
        currentProfile.value = users.value[0];
        localStorage.setItem('userId', currentProfile.value.id);
        return currentProfile.value;
      }
      
      console.error('No se pudo obtener o crear un perfil de usuario');
      error.value = 'No se pudo obtener el perfil de usuario';
      return null;
    } catch (err: unknown) {
      error.value = 'Error al cargar perfil del usuario actual';
      console.error('Error al cargar perfil del usuario actual:', err);
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
      
      // Si no hay usuarios cargados, intentar cargar desde localStorage
      if (users.value.length === 0) {
        loadUsersFromLocalStorage();
        
        // Buscar nuevamente después de cargar
        const userFromStorage = users.value.find((u: User) => u.id === userId);
        if (userFromStorage) {
          console.log('Usuario encontrado en localStorage:', userFromStorage);
          return userFromStorage;
        }
      }
      
      console.log(`Usuario con ID ${userId} no encontrado localmente`);
      
      // Si se trata del usuario actual según localStorage
      if (userId === localStorage.getItem('userId')) {
        console.log('Creando usuario actual mock');
        const mockCurrentUser: User = {
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
        
        // Guardar en memoria para futuras consultas
        users.value.push(mockCurrentUser);
        saveUsersToLocalStorage();
        
        return mockCurrentUser;
      }
      
      error.value = 'Usuario no encontrado';
      return null;
    } catch (err) {
      console.error(`Error al obtener usuario con ID ${userId}:`, err);
      error.value = 'Error al obtener usuario';
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
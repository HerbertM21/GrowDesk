import apiClient from '../api/client';
import type { User } from '@/stores/users';

interface LoginCredentials {
  email: string;
  password: string;
}

interface AuthResponse {
  token: string;
  user: User;
}

export const authService = {
  async login(credentials: LoginCredentials): Promise<AuthResponse> {
    try {
      // En un entorno real, esto se conectaría con el backend
      // Por ahora, simulamos la respuesta
      const response = await apiClient.post<AuthResponse>('/auth/login', credentials);
      
      // Guardar datos importantes en localStorage
      if (response.data && response.data.token) {
        localStorage.setItem('token', response.data.token);
        
        // Guardar ID de usuario para uso en WebSocket y otras partes
        if (response.data.user && response.data.user.id) {
          localStorage.setItem('userId', String(response.data.user.id));
        }
        
        // Guardar rol de usuario para controles de acceso
        if (response.data.user && response.data.user.role) {
          localStorage.setItem('userRole', String(response.data.user.role));
        }
        
        // Guardar objeto de usuario completo para uso local
        try {
          localStorage.setItem('user', JSON.stringify(response.data.user));
        } catch (e) {
          console.warn('No se pudo guardar el objeto de usuario completo');
        }
      }
      
      return response.data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout');
      this.clearUserData();
    } catch (error) {
      console.error('Logout error:', error);
      // Aun si falla la petición, limpiamos el storage
      this.clearUserData();
    }
  },
  
  clearUserData(): void {
    // Eliminar todos los datos del usuario al cerrar sesión
    localStorage.removeItem('token');
    localStorage.removeItem('user');
    localStorage.removeItem('userId');
    localStorage.removeItem('userRole');
  },

  async checkAuth(): Promise<User | null> {
    try {
      const token = localStorage.getItem('token');
      if (!token) {
        return null;
      }
      
      const response = await apiClient.get<User>('/auth/me');
      
      // Actualizar datos del usuario si es necesario
      if (response.data && response.data.id) {
        localStorage.setItem('userId', String(response.data.id));
      }
      
      return response.data;
    } catch (error) {
      // En caso de error, no eliminar el token automáticamente
      // para permitir reintentos cuando hay problemas temporales
      console.error('Error verificando autenticación:', error);
      return null;
    }
  },

  async updateProfile(userData: Partial<User>): Promise<User> {
    try {
      const response = await apiClient.put<User>('/auth/profile', userData);
      return response.data;
    } catch (error) {
      console.error('Profile update error:', error);
      throw error;
    }
  },

  async changePassword(currentPassword: string, newPassword: string): Promise<void> {
    try {
      await apiClient.put('/auth/password', {
        currentPassword,
        newPassword
      });
    } catch (error) {
      console.error('Password change error:', error);
      throw error;
    }
  }
};

export default authService;
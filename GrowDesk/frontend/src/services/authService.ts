import apiClient from './apiClient';
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
      return response.data;
    } catch (error) {
      console.error('Login error:', error);
      throw error;
    }
  },

  async logout(): Promise<void> {
    try {
      await apiClient.post('/auth/logout');
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    } catch (error) {
      console.error('Logout error:', error);
      // Aun si falla la petición, limpiamos el storage
      localStorage.removeItem('token');
      localStorage.removeItem('user');
    }
  },

  async checkAuth(): Promise<User | null> {
    try {
      const response = await apiClient.get<User>('/auth/me');
      return response.data;
    } catch (error) {
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
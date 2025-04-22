import apiClient from '../api/client';
import type { User } from '../stores/users';

export interface UserCreateData {
  firstName: string;
  lastName: string;
  email: string;
  password: string;
  role: string;
  department?: string;
  isActive?: boolean;
}

export interface UserUpdateData {
  firstName?: string;
  lastName?: string;
  email?: string;
  role?: string;
  department?: string;
  isActive?: boolean;
}

const userService = {

  async getAllUsers(): Promise<User[]> {
    try {
      const response = await apiClient.get('/users');
      return response.data;
    } catch (error) {
      console.error('Error fetching users:', error);
      throw error;
    }
  },


  async getUser(id: string): Promise<User> {
    try {
      const response = await apiClient.get(`/users/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error fetching user ${id}:`, error);
      throw error;
    }
  },


  async createUser(userData: UserCreateData): Promise<User> {
    try {
      const response = await apiClient.post('/users', userData);
      return response.data;
    } catch (error) {
      console.error('Error creating user:', error);
      throw error;
    }
  },


  async updateUser(id: string, userData: UserUpdateData): Promise<User> {
    try {
      const response = await apiClient.put(`/users/${id}`, userData);
      return response.data;
    } catch (error) {
      console.error(`Error updating user ${id}:`, error);
      throw error;
    }
  },


  async deleteUser(id: string): Promise<void> {
    try {
      await apiClient.delete(`/users/${id}`);
    } catch (error) {
      console.error(`Error deleting user ${id}:`, error);
      throw error;
    }
  },

 
  async updateUserRole(id: string, role: string): Promise<User> {
    try {
      const response = await apiClient.patch(`/users/${id}/role`, { role });
      return response.data;
    } catch (error) {
      console.error(`Error updating role for user ${id}:`, error);
      throw error;
    }
  },

  
  async toggleUserActive(id: string, isActive: boolean): Promise<User> {
    try {
      const response = await apiClient.patch(`/users/${id}/status`, { isActive });
      return response.data;
    } catch (error) {
      console.error(`Error toggling active status for user ${id}:`, error);
      throw error;
    }
  },


  async getCurrentUserProfile(): Promise<User> {
    try {
      const response = await apiClient.get('/users/profile');
      return response.data;
    } catch (error) {
      console.error('Error fetching current user profile:', error);
      throw error;
    }
  }
};

export default userService; 
import apiClient from '../api/client';
import type { FAQ, FaqCreateData, FaqUpdateData } from '../types/faq.types';

// Servicio para manejar las preguntas frecuentes
const faqService = {
  /**
   * Obtener todas las preguntas frecuentes
   */
  async getAllFaqs(): Promise<FAQ[]> {
    try {
      const response = await apiClient.get('/faqs');
      console.log('FAQs cargadas:', response.data);
      return response.data;
    } catch (error) {
      console.error('Error al obtener preguntas frecuentes:', error);
      throw error;
    }
  },

  /**
   * Obtener una pregunta frecuente por su ID
   */
  async getFaq(id: number): Promise<FAQ> {
    try {
      const response = await apiClient.get(`/faqs/${id}`);
      return response.data;
    } catch (error) {
      console.error(`Error al obtener la pregunta frecuente ${id}:`, error);
      throw error;
    }
  },

  /**
   * Crear una nueva pregunta frecuente
   */
  async createFaq(faqData: FaqCreateData): Promise<FAQ> {
    try {
      const response = await apiClient.post('/faqs', faqData);
      return response.data;
    } catch (error) {
      console.error('Error al crear pregunta frecuente:', error);
      throw error;
    }
  },

  /**
   * Actualizar una pregunta frecuente existente
   */
  async updateFaq(faqData: FaqUpdateData): Promise<FAQ> {
    try {
      const response = await apiClient.put(`/faqs/${faqData.id}`, faqData);
      return response.data;
    } catch (error) {
      console.error(`Error al actualizar la pregunta frecuente ${faqData.id}:`, error);
      throw error;
    }
  },

  /**
   * Eliminar una pregunta frecuente
   */
  async deleteFaq(id: number): Promise<void> {
    try {
      await apiClient.delete(`/faqs/${id}`);
    } catch (error) {
      console.error(`Error al eliminar la pregunta frecuente ${id}:`, error);
      throw error;
    }
  },

  /**
   * Cambiar el estado de publicación de una pregunta frecuente
   */
  async togglePublishStatus(id: number): Promise<FAQ> {
    try {
      const response = await apiClient.patch(`/faqs/${id}/toggle-publish`);
      return response.data;
    } catch (error) {
      console.error(`Error al cambiar el estado de publicación de la pregunta frecuente ${id}:`, error);
      throw error;
    }
  }
};

export default faqService;
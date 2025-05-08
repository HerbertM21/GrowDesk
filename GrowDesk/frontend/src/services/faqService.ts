import apiClient from '@/api/client';
import type { FAQ } from '@/types/faq.types';

export const getAllFaqs = async (): Promise<FAQ[]> => {
  try {
    console.log('Obteniendo todas las FAQs...');
    const response = await apiClient.get('/faqs');
    console.log('Respuesta de FAQs:', response.data);
    
    if (!response.data) {
      console.error('Respuesta vacía del servidor');
      throw new Error('No se recibieron datos del servidor');
    }
    
    if (!Array.isArray(response.data)) {
      console.error('Respuesta inválida del servidor:', response.data);
      throw new Error('La respuesta del servidor no es un array de FAQs');
    }
    
    return response.data;
  } catch (error: any) {
    console.error('Error al obtener preguntas frecuentes:', {
      message: error.message,
      code: error.code,
      response: error.response?.data
    });
    throw error;
  }
};

export const getFaq = async (id: number): Promise<FAQ> => {
  try {
    console.log('Obteniendo FAQ:', id);
    const response = await apiClient.get(`/faqs/${id}`);
    console.log('FAQ obtenida:', response.data);
    
    if (!response.data || !response.data.id) {
      throw new Error('La respuesta del servidor no contiene una FAQ válida');
    }
    
    return response.data;
  } catch (error: any) {
    console.error('Error al obtener FAQ:', error);
    throw error;
  }
};

export const createFaq = async (faq: Omit<FAQ, 'id' | 'createdAt' | 'updatedAt'>): Promise<FAQ> => {
  try {
    console.log('Creando nueva FAQ:', faq);
    const response = await apiClient.post('/faqs', faq);
    console.log('Respuesta de creación:', response.data);
    
    if (!response.data) {
      console.error('Respuesta vacía del servidor');
      throw new Error('No se recibió respuesta del servidor');
    }
    
    if (!response.data.id) {
      console.error('Respuesta inválida del servidor:', response.data);
      throw new Error('La respuesta del servidor no contiene un ID válido');
    }
    
    return response.data;
  } catch (error: any) {
    console.error('Error al crear FAQ:', {
      message: error.message,
      code: error.code,
      response: error.response?.data,
      faq
    });
    throw error;
  }
};

export const updateFaq = async (id: number, faq: Partial<FAQ>): Promise<FAQ> => {
  try {
    console.log('Actualizando FAQ:', { id, faq });
    const response = await apiClient.put(`/faqs/${id}`, faq);
    console.log('Respuesta de actualización:', response.data);
    
    if (!response.data) {
      throw new Error('No se recibió respuesta del servidor');
    }
    
    if (!response.data.id) {
      throw new Error('La respuesta del servidor no contiene un ID válido');
    }
    
    return response.data;
  } catch (error: any) {
    console.error('Error al actualizar FAQ:', error);
    throw error;
  }
};

export const deleteFaq = async (id: number): Promise<void> => {
  try {
    console.log('Eliminando FAQ:', id);
    await apiClient.delete(`/faqs/${id}`);
    console.log('FAQ eliminada exitosamente');
  } catch (error: any) {
    console.error('Error al eliminar FAQ:', error);
    throw error;
  }
};

export const togglePublishStatus = async (id: number): Promise<FAQ> => {
  try {
    console.log('Cambiando estado de publicación de FAQ:', id);
    const response = await apiClient.patch(`/faqs/${id}/toggle-publish`);
    console.log('Respuesta de cambio de estado:', response.data);
    
    if (!response.data) {
      throw new Error('No se recibió respuesta del servidor');
    }
    
    if (!response.data.id) {
      throw new Error('La respuesta del servidor no contiene un ID válido');
    }
    
    return response.data;
  } catch (error: any) {
    console.error('Error al cambiar estado de publicación:', error);
    throw error;
  }
};

const faqService = {
  getAllFaqs,
  getFaq,
  createFaq,
  updateFaq,
  deleteFaq,
  togglePublishStatus
};

export default faqService;
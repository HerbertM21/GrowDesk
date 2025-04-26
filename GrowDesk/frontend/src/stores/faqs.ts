import { defineStore } from 'pinia';
import { ref } from 'vue';
import faqService from '../services/faqService';
import type { FAQ, FaqCreateData, FaqUpdateData } from '../types/faq.types';

export const useFaqsStore = defineStore('faqs', () => {
  const faqs = ref<FAQ[]>([]);
  const isLoading = ref(false);
  const error = ref<string | null>(null);

  /**
   * Carga todas las FAQs (alias para compatibilidad con FaqManagement)
   */
  const fetchFaqs = async () => {
    return loadFaqs();
  };

  /**
   * Carga todas las FAQs
   */
  const loadFaqs = async () => {
    isLoading.value = true;
    error.value = null;
    
    try {
      faqs.value = await faqService.getAllFaqs();
    } catch (err) {
      error.value = 'Error al cargar preguntas frecuentes';
      console.error(error.value, err);
    } finally {
      isLoading.value = false;
    }
  };

  /**
   * Añade una nueva FAQ
   */
  const addFaq = async (faq: FaqCreateData) => {
    isLoading.value = true;
    error.value = null;
    
    try {
      const newFaq = await faqService.createFaq(faq);
      faqs.value.push(newFaq);
      return newFaq;
    } catch (err) {
      error.value = 'Error al crear la pregunta frecuente';
      console.error(error.value, err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  /**
   * Actualiza una FAQ existente
   */
  const updateFaq = async (faq: FaqUpdateData) => {
    isLoading.value = true;
    error.value = null;
    
    try {
      const updatedFaq = await faqService.updateFaq(faq);
      const index = faqs.value.findIndex(f => f.id === faq.id);
      if (index !== -1) {
        faqs.value[index] = updatedFaq;
      }
      return updatedFaq;
    } catch (err) {
      error.value = 'Error al actualizar la pregunta frecuente';
      console.error(error.value, err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  /**
   * Alias de updateFaq para mantener compatibilidad con FaqManagement
   */
  const editFaq = async (faq: FaqUpdateData) => {
    return updateFaq(faq);
  };

  /**
   * Elimina una FAQ
   */
  const deleteFaq = async (id: number) => {
    isLoading.value = true;
    error.value = null;
    
    try {
      await faqService.deleteFaq(id);
      faqs.value = faqs.value.filter(f => f.id !== id);
    } catch (err) {
      error.value = 'Error al eliminar la pregunta frecuente';
      console.error(error.value, err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  /**
   * Cambia el estado de publicación de una FAQ
   */
  const togglePublish = async (id: number) => {
    isLoading.value = true;
    error.value = null;
    
    try {
      const updatedFaq = await faqService.togglePublishStatus(id);
      const index = faqs.value.findIndex(f => f.id === id);
      if (index !== -1) {
        faqs.value[index] = updatedFaq;
      }
      return updatedFaq;
    } catch (err) {
      error.value = 'Error al cambiar el estado de publicación';
      console.error(error.value, err);
      throw err;
    } finally {
      isLoading.value = false;
    }
  };

  /**
   * Alias para togglePublish para mantener compatibilidad con FaqManagement
   */
  const togglePublishStatus = async (id: number) => {
    return togglePublish(id);
  };

  /**
   * Obtiene todas las categorías únicas de las FAQs
   */
  const getCategories = () => {
    const categories = new Set<string>();
    faqs.value.forEach(faq => {
      if (faq.category) {
        categories.add(faq.category);
      }
    });
    return Array.from(categories);
  };

  return { 
    faqs, 
    isLoading, 
    error, 
    loadFaqs,
    fetchFaqs,
    addFaq, 
    editFaq,
    updateFaq,
    removeFaq: deleteFaq,
    deleteFaq,
    togglePublish,
    togglePublishStatus,
    getCategories
  };
});
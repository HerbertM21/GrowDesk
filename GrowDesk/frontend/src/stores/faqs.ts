import { defineStore } from 'pinia';
import { ref, computed } from 'vue';
import type { FAQ } from '@/types/faq.types';
import faqService from '@/services/faqService';

export const useFaqStore = defineStore('faqs', () => {
  const faqs = ref<FAQ[]>([]);
  const loading = ref(false);
  const error = ref<string | null>(null);

  const DEFAULT_CATEGORIES = [
    'General',
    'Cuenta',
    'Pagos',
    'Soporte',
    'Técnico',
    'Otros'
  ];

  const fetchFaqs = async () => {
    try {
      loading.value = true;
      error.value = null;
      console.log('Iniciando fetch de FAQs...');
      const data = await faqService.getAllFaqs();
      console.log('FAQs obtenidas:', data);
      faqs.value = data;
    } catch (err: any) {
      console.error('Error en fetchFaqs:', err);
      error.value = err.message || 'Error al cargar las preguntas frecuentes';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const addFaq = async (faq: Omit<FAQ, 'id' | 'createdAt' | 'updatedAt'>) => {
    try {
      loading.value = true;
      error.value = null;
      console.log('Agregando nueva FAQ:', faq);
      const newFaq = await faqService.createFaq(faq);
      console.log('FAQ creada:', newFaq);
      faqs.value = [...faqs.value, newFaq];
      return newFaq;
    } catch (err: any) {
      console.error('Error en addFaq:', err);
      error.value = err.message || 'Error al crear la pregunta frecuente';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const updateFaq = async (id: number, faq: Partial<FAQ>) => {
    try {
      loading.value = true;
      error.value = null;
      console.log('Actualizando FAQ:', { id, faq });
      const updatedFaq = await faqService.updateFaq(id, faq);
      console.log('FAQ actualizada:', updatedFaq);
      faqs.value = faqs.value.map(f => f.id === id ? updatedFaq : f);
      return updatedFaq;
    } catch (err: any) {
      console.error('Error en updateFaq:', err);
      error.value = err.message || 'Error al actualizar la pregunta frecuente';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const deleteFaq = async (id: number) => {
    try {
      loading.value = true;
      error.value = null;
      console.log('Eliminando FAQ:', id);
      await faqService.deleteFaq(id);
      faqs.value = faqs.value.filter(f => f.id !== id);
    } catch (err: any) {
      console.error('Error en deleteFaq:', err);
      error.value = err.message || 'Error al eliminar la pregunta frecuente';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const togglePublish = async (id: number) => {
    try {
      loading.value = true;
      error.value = null;
      console.log('Cambiando estado de publicación de FAQ:', id);
      const updatedFaq = await faqService.togglePublishStatus(id);
      console.log('Estado de publicación actualizado:', updatedFaq);
      faqs.value = faqs.value.map(f => f.id === id ? updatedFaq : f);
      return updatedFaq;
    } catch (err: any) {
      console.error('Error en togglePublish:', err);
      error.value = err.message || 'Error al cambiar el estado de publicación';
      throw err;
    } finally {
      loading.value = false;
    }
  };

  const getCategories = () => {
    if (!faqs.value || faqs.value.length === 0) return DEFAULT_CATEGORIES;
    const categories = new Set<string>();
    faqs.value.forEach(faq => {
      if (faq.category) categories.add(faq.category);
    });
    return Array.from(categories);
  };

  const getFaqById = (id: number) => {
    return faqs.value.find(faq => faq.id === id);
  };

  return {
    faqs,
    loading,
    error,
    fetchFaqs,
    addFaq,
    updateFaq,
    deleteFaq,
    togglePublish,
    getCategories,
    getFaqById
  };
});
<template>
  <div class="faq-manager">
    <div class="d-flex justify-content-between align-items-center mb-4">
      <h1>Gestión de Preguntas Frecuentes</h1>
      <button class="btn btn-primary" @click="openAddModal">
        <i class="bi bi-plus"></i> Crear Nueva FAQ
      </button>
    </div>
    
    <div v-if="faqsStore.error" class="alert alert-danger" role="alert">
      {{ faqsStore.error }}
    </div>
    
    <div class="card mb-4">
      <div class="card-body">
        <div v-if="faqsStore.isLoading" class="text-center py-3">
          <div class="spinner-border" role="status">
            <span class="visually-hidden">Cargando...</span>
          </div>
        </div>
        <div v-else-if="faqsStore.faqs.length === 0" class="text-center py-3">
          <p>No hay preguntas frecuentes disponibles.</p>
        </div>
        <div v-else>
          <table class="table table-hover">
            <thead>
              <tr>
                <th>ID</th>
                <th>Pregunta</th>
                <th>Categoría</th>
                <th>Estado</th>
                <th>Acciones</th>
              </tr>
            </thead>
            <tbody>
              <tr v-for="faq in faqsStore.faqs" :key="faq.id">
                <td>{{ faq.id }}</td>
                <td>{{ truncateText(faq.question, 70) }}</td>
                <td>{{ faq.category }}</td>
                <td>
                  <span 
                    :class="['badge', faq.isPublished ? 'bg-success' : 'bg-secondary']"
                  >
                    {{ faq.isPublished ? 'Publicado' : 'Borrador' }}
                  </span>
                </td>
                <td>
                  <div class="btn-group">
                    <button 
                      class="btn btn-sm btn-outline-secondary" 
                      @click="togglePublishFaq(faq)"
                      :title="faq.isPublished ? 'Pasar a borrador' : 'Publicar'"
                    >
                      <i :class="['bi', faq.isPublished ? 'bi-eye-slash' : 'bi-eye']"></i>
                    </button>
                    <button 
                      class="btn btn-sm btn-outline-primary" 
                      @click="openEditModal(faq)"
                      title="Editar"
                    >
                      <i class="bi bi-pencil"></i>
                    </button>
                    <button 
                      class="btn btn-sm btn-outline-danger" 
                      @click="confirmDelete(faq)"
                      title="Eliminar"
                    >
                      <i class="bi bi-trash"></i>
                    </button>
                  </div>
                </td>
              </tr>
            </tbody>
          </table>
        </div>
      </div>
    </div>

    <!-- Modal para crear/editar FAQ -->
    <div class="modal fade" id="faqModal" tabindex="-1" aria-labelledby="faqModalLabel" aria-hidden="true">
      <div class="modal-dialog modal-lg">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="faqModalLabel">{{ isEditing ? 'Editar' : 'Crear' }} Pregunta Frecuente</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <form @submit.prevent="saveFaq">
              <div class="mb-3">
                <label for="question" class="form-label">Pregunta</label>
                <input 
                  type="text" 
                  class="form-control" 
                  id="question" 
                  v-model="currentFaq.question" 
                  required
                  placeholder="¿Cómo puedo...?"
                >
              </div>
              
              <div class="mb-3">
                <label for="answer" class="form-label">Respuesta</label>
                <textarea 
                  class="form-control" 
                  id="answer" 
                  v-model="currentFaq.answer" 
                  rows="5" 
                  required
                  placeholder="Para hacer esto, usted debe..."
                ></textarea>
              </div>
              
              <div class="mb-3">
                <label for="category" class="form-label">Categoría</label>
                <select class="form-select" id="category" v-model="currentFaq.category" required>
                  <option value="" disabled>Seleccione una categoría</option>
                  <option value="General">General</option>
                  <option value="Cuenta">Cuenta</option>
                  <option value="Facturación">Facturación</option>
                  <option value="Soporte Técnico">Soporte Técnico</option>
                  <option value="Configuración">Configuración</option>
                </select>
              </div>
              
              <div class="form-check mb-3">
                <input class="form-check-input" type="checkbox" id="isPublished" v-model="currentFaq.isPublished">
                <label class="form-check-label" for="isPublished">
                  Publicar inmediatamente (visible en el widget de chat)
                </label>
              </div>
              
              <div class="d-flex justify-content-end gap-2">
                <button 
                  type="button" 
                  class="btn btn-secondary" 
                  data-bs-dismiss="modal"
                >
                  Cancelar
                </button>
                <button 
                  type="submit" 
                  class="btn btn-primary" 
                  :disabled="faqsStore.isLoading"
                >
                  <span v-if="faqsStore.isLoading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
                  {{ isEditing ? 'Guardar Cambios' : 'Crear Pregunta' }}
                </button>
              </div>
            </form>
          </div>
        </div>
      </div>
    </div>

    <!-- Modal de confirmación de eliminación -->
    <div class="modal fade" id="deleteModal" tabindex="-1" aria-labelledby="deleteModalLabel" aria-hidden="true">
      <div class="modal-dialog">
        <div class="modal-content">
          <div class="modal-header">
            <h5 class="modal-title" id="deleteModalLabel">Confirmar eliminación</h5>
            <button type="button" class="btn-close" data-bs-dismiss="modal" aria-label="Close"></button>
          </div>
          <div class="modal-body">
            <p>¿Está seguro que desea eliminar esta pregunta frecuente?</p>
            <p class="fw-bold">{{ currentFaq.question }}</p>
            <p class="text-danger">Esta acción no se puede deshacer.</p>
          </div>
          <div class="modal-footer">
            <button type="button" class="btn btn-secondary" data-bs-dismiss="modal">Cancelar</button>
            <button 
              type="button" 
              class="btn btn-danger" 
              @click="deleteFaq" 
              :disabled="faqsStore.isLoading"
            >
              <span v-if="faqsStore.isLoading" class="spinner-border spinner-border-sm me-1" role="status" aria-hidden="true"></span>
              Eliminar
            </button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue';
import { useFaqsStore } from '@/stores/faqs';
import { FAQ } from '@/types/faq.types';
import { Modal } from 'bootstrap';

const faqsStore = useFaqsStore();
const isEditing = ref(false);
const currentFaq = ref<FAQ>({
  question: '',
  answer: '',
  category: '',
  isPublished: false
});

// Referencias a los modales de Bootstrap
let faqModal: Modal | null = null;
let deleteModal: Modal | null = null;

onMounted(async () => {
  try {
    // Inicializar los modales de Bootstrap
    faqModal = new Modal(document.getElementById('faqModal')!);
    deleteModal = new Modal(document.getElementById('deleteModal')!);
    
    // Cargar las FAQs al montar el componente
    await faqsStore.loadFaqs();
  } catch (error) {
    console.error('Error al inicializar la vista de FAQs:', error);
  }
});

/**
 * Abre el modal para crear una nueva FAQ
 */
const openAddModal = () => {
  isEditing.value = false;
  currentFaq.value = {
    question: '',
    answer: '',
    category: '',
    isPublished: false
  };
  faqModal?.show();
};

/**
 * Abre el modal para editar una FAQ existente
 */
const openEditModal = (faq: FAQ) => {
  isEditing.value = true;
  currentFaq.value = { ...faq };
  faqModal?.show();
};

/**
 * Guarda una FAQ (nueva o existente)
 */
const saveFaq = async () => {
  try {
    if (isEditing.value) {
      await faqsStore.editFaq(currentFaq.value);
    } else {
      await faqsStore.addFaq(currentFaq.value);
    }
    faqModal?.hide();
  } catch (error) {
    console.error('Error al guardar FAQ:', error);
  }
};

/**
 * Cambia el estado de publicación de una FAQ
 */
const togglePublishFaq = async (faq: FAQ) => {
  try {
    if (!faq.id) return;
    await faqsStore.togglePublish(faq.id);
  } catch (error) {
    console.error('Error al cambiar estado de publicación:', error);
  }
};

/**
 * Abre el modal de confirmación para eliminar una FAQ
 */
const confirmDelete = (faq: FAQ) => {
  currentFaq.value = { ...faq };
  deleteModal?.show();
};

/**
 * Elimina una FAQ
 */
const deleteFaq = async () => {
  try {
    if (!currentFaq.value.id) return;
    await faqsStore.removeFaq(currentFaq.value.id);
    deleteModal?.hide();
  } catch (error) {
    console.error('Error al eliminar FAQ:', error);
  }
};

/**
 * Trunca un texto largo para mostrarlo en la tabla
 */
const truncateText = (text: string, maxLength: number): string => {
  if (text.length <= maxLength) return text;
  return text.substring(0, maxLength) + '...';
};
</script>

<style scoped>
.faq-manager {
  padding: 1rem;
}
</style>
<template>
  <div class="admin-section">
    <div class="admin-section-header">
      <div>
        <h1 class="admin-section-title">Gestión de Preguntas Frecuentes</h1>
        <p class="admin-section-description">
          En esta sección el equipo de soporte técnico puede crear y administrar las preguntas frecuentes que se mostrarán a los clientes.
        </p>
      </div>
      
      <div class="admin-section-actions">
        <button class="btn btn-primary" @click="openCreateModal">
          <i class="pi pi-plus"></i> Nueva Pregunta
        </button>
      </div>
    </div>
    
    <div v-if="loading" class="loading">
      <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
      <p>Cargando preguntas frecuentes...</p>
    </div>
    <div v-else-if="error" class="alert alert-danger">
      <i class="pi pi-exclamation-circle"></i> {{ error }}
    </div>
    
    <div v-else>
      <div class="filter-controls mb-4">
        <div class="search-box">
          <div class="form-group">
            <label class="form-label">Buscar:</label>
            <div class="search-input-container">
              <input type="text" v-model="searchQuery" class="form-control" placeholder="Buscar preguntas..." />
              <i class="pi pi-search search-icon"></i>
            </div>
          </div>
        </div>
        
        <div class="filter-group">
          <label class="form-label">Filtrar por categoría:</label>
          <select v-model="categoryFilter" class="form-control">
            <option value="all">Todas las categorías</option>
            <option v-for="category in categories" :key="category" :value="category">
              {{ category }}
            </option>
          </select>
        </div>
        
        <div class="filter-group">
          <label class="form-label">Estado:</label>
          <select v-model="publishFilter" class="form-control">
            <option value="all">Todos</option>
            <option value="published">Publicados</option>
            <option value="unpublished">No publicados</option>
          </select>
        </div>
      </div>
      
      <div class="faq-list-container">
        <table class="admin-table">
          <thead>
            <tr>
              <th>Pregunta</th>
              <th>Categoría</th>
              <th>Estado</th>
              <th>Última actualización</th>
              <th>Acciones</th>
            </tr>
          </thead>
          <tbody>
            <tr v-if="filteredFaqs.length === 0">
              <td colspan="5" class="no-results">
                No se encontraron preguntas frecuentes que coincidan con los filtros.
              </td>
            </tr>
            <tr v-for="faq in filteredFaqs" :key="faq.id">
              <td class="question-cell">{{ faq.question }}</td>
              <td><span class="category-tag">{{ faq.category }}</span></td>
              <td>
                <span :class="['status-badge', faq.isPublished ? 'published' : 'unpublished']">
                  {{ faq.isPublished ? 'Publicado' : 'No publicado' }}
                </span>
              </td>
              <td>{{ formatDate(faq.updatedAt) }}</td>
              <td class="actions-cell">
                <button @click="openViewModal(faq)" class="btn btn-sm btn-outline-primary" title="Ver detalles">
                  <i class="pi pi-eye"></i>
                </button>
                <button @click="openEditModal(faq)" class="btn btn-sm btn-outline-primary" title="Editar">
                  <i class="pi pi-pencil"></i>
                </button>
                <button @click="confirmDelete(faq)" class="btn btn-sm btn-outline-danger" title="Eliminar">
                  <i class="pi pi-trash"></i>
                </button>
                <button 
                  @click="togglePublish(faq.id)" 
                  class="btn btn-sm" 
                  :class="faq.isPublished ? 'btn-outline-warning' : 'btn-outline-success'"
                  :title="faq.isPublished ? 'Despublicar' : 'Publicar'"
                >
                  <i :class="['pi', faq.isPublished ? 'pi-eye-slash' : 'pi-check-circle']"></i>
                </button>
              </td>
            </tr>
          </tbody>
        </table>
      </div>
    </div>
    
    <!-- Modal para ver detalles de una pregunta -->
    <div v-if="showViewModal" class="modal-overlay" @click.self="showViewModal = false">
      <div class="modal-content">
        <div class="modal-header">
          <h2>Detalles de la Pregunta</h2>
          <button class="close-btn" @click="showViewModal = false">&times;</button>
        </div>
        <div class="modal-body">
          <div class="faq-details">
            <div class="detail-section">
              <h3>Pregunta</h3>
              <p>{{ currentFaq.question }}</p>
            </div>
            
            <div class="detail-section">
              <h3>Respuesta</h3>
              <div class="answer-content">{{ currentFaq.answer }}</div>
            </div>
            
            <div class="detail-group">
              <div class="detail-item">
                <h4>Categoría</h4>
                <p class="category-tag">{{ currentFaq.category }}</p>
              </div>
              
              <div class="detail-item">
                <h4>Estado</h4>
                <span :class="['status-badge', currentFaq.isPublished ? 'published' : 'unpublished']">
                  {{ currentFaq.isPublished ? 'Publicado' : 'No publicado' }}
                </span>
              </div>
            </div>
            
            <div class="detail-group">
              <div class="detail-item">
                <h4>Creado</h4>
                <p>{{ formatDate(currentFaq.createdAt) }}</p>
              </div>
              
              <div class="detail-item">
                <h4>Última actualización</h4>
                <p>{{ formatDate(currentFaq.updatedAt) }}</p>
              </div>
            </div>
          </div>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline-secondary" @click="showViewModal = false">Cerrar</button>
          <button class="btn btn-primary" @click="openEditModal(currentFaq)">Editar</button>
        </div>
      </div>
    </div>
    
    <!-- Modal para crear/editar una pregunta -->
    <div v-if="showEditModal" class="modal-overlay" @click.self="cancelEdit">
      <div class="modal-content">
        <div class="modal-header">
          <h2>{{ isEditing ? 'Editar Pregunta' : 'Nueva Pregunta' }}</h2>
          <button class="close-btn" @click="cancelEdit">&times;</button>
        </div>
        <div class="modal-body">
          <form @submit.prevent="saveFaq" class="faq-form">
            <div class="form-group">
              <label for="question" class="form-label">Pregunta <span class="required">*</span></label>
              <input 
                type="text" 
                id="question" 
                v-model="formData.question" 
                required
                class="form-control"
                :class="{ 'is-invalid': formErrors.question }"
              />
              <span v-if="formErrors.question" class="error-text">{{ formErrors.question }}</span>
            </div>
            
            <div class="form-group">
              <label for="answer" class="form-label">Respuesta <span class="required">*</span></label>
              <textarea 
                id="answer" 
                v-model="formData.answer" 
                rows="6" 
                required
                class="form-control"
                :class="{ 'is-invalid': formErrors.answer }"
              ></textarea>
              <span v-if="formErrors.answer" class="error-text">{{ formErrors.answer }}</span>
            </div>
            
            <div class="form-group">
              <label for="category" class="form-label">Categoría <span class="required">*</span></label>
              <div class="category-input">
                <select 
                  id="category" 
                  v-model="formData.category" 
                  required
                  class="form-control"
                  :class="{ 'is-invalid': formErrors.category }"
                >
                  <option value="" disabled>Seleccione una categoría</option>
                  <option v-for="category in categories" :key="category" :value="category">
                    {{ category }}
                  </option>
                </select>
                <span v-if="formErrors.category" class="error-text">{{ formErrors.category }}</span>
              </div>
            </div>
            
            <div class="form-group">
              <div class="checkbox-container">
                <input type="checkbox" id="isPublished" v-model="formData.isPublished" />
                <label for="isPublished" class="checkbox-label">Publicar esta pregunta</label>
              </div>
              <p class="help-text">
                Las preguntas publicadas serán visibles para los usuarios en la sección de ayuda.
              </p>
            </div>
          </form>
        </div>
        <div class="modal-footer">
          <button type="button" class="btn btn-outline-secondary" @click="cancelEdit">Cancelar</button>
          <button type="button" class="btn btn-primary" @click="saveFaq">
            {{ isEditing ? 'Actualizar' : 'Guardar' }}
          </button>
        </div>
      </div>
    </div>
    
    <!-- Modal de confirmación para eliminar -->
    <div v-if="showDeleteConfirm" class="modal-overlay" @click.self="showDeleteConfirm = false">
      <div class="modal-content confirm-modal">
        <div class="modal-header">
          <h2>Confirmar eliminación</h2>
          <button class="close-btn" @click="showDeleteConfirm = false">&times;</button>
        </div>
        <div class="modal-body">
          <p>¿Estás seguro de que deseas eliminar esta pregunta frecuente?</p>
          <p class="warning-text">Esta acción no se puede deshacer.</p>
        </div>
        <div class="modal-footer">
          <button class="btn btn-outline-secondary" @click="showDeleteConfirm = false">Cancelar</button>
          <button class="btn btn-danger" @click="deleteFaq">Eliminar</button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useFaqsStore, type FAQ } from '@/stores/faqs';
import { storeToRefs } from 'pinia';

// Store para FAQs
const faqStore = useFaqsStore();
const { faqs, loading, error } = storeToRefs(faqStore);

// Estado para filtros
const searchQuery = ref('');
const categoryFilter = ref('all');
const publishFilter = ref('all');

// Estado para modales
const showViewModal = ref(false);
const showEditModal = ref(false);
const showDeleteConfirm = ref(false);
const showNewCategoryInput = ref(false);

// Estado para formularios
const isEditing = ref(false);
const isSaving = ref(false);
const isDeleting = ref(false);
const currentFaq = ref<FAQ>({
  id: 0,
  question: '',
  answer: '',
  category: '',
  isPublished: false,
  createdAt: '',
  updatedAt: ''
});
const formData = ref({
  question: '',
  answer: '',
  category: '',
  isPublished: false
});
const formErrors = ref({
  question: '',
  answer: '',
  category: ''
});
const newCategory = ref('');
const faqToDelete = ref<number | null>(null);

// Cargar FAQs al montar el componente
onMounted(() => {
  faqStore.fetchFaqs();
});

// Computed properties
const categories = computed(() => {
  return faqStore.getCategories();
});

const filteredFaqs = computed(() => {
  let filtered = [...faqs.value];
  
  // Filtrar por búsqueda
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(faq => 
      faq.question.toLowerCase().includes(query) || 
      faq.answer.toLowerCase().includes(query)
    );
  }
  
  // Filtrar por categoría
  if (categoryFilter.value !== 'all') {
    filtered = filtered.filter(faq => faq.category === categoryFilter.value);
  }
  
  // Filtrar por estado de publicación
  if (publishFilter.value !== 'all') {
    const isPublished = publishFilter.value === 'published';
    filtered = filtered.filter(faq => faq.isPublished === isPublished);
  }
  
  // Ordenar por fecha de actualización más reciente
  return filtered.sort((a, b) => {
    return new Date(b.updatedAt).getTime() - new Date(a.updatedAt).getTime();
  });
});

// Métodos
const formatDate = (dateString: string) => {
  const options: Intl.DateTimeFormatOptions = { 
    year: 'numeric', 
    month: 'short', 
    day: 'numeric',
    hour: '2-digit',
    minute: '2-digit'
  };
  return new Date(dateString).toLocaleDateString('es-ES', options);
};

const openViewModal = (faq: FAQ) => {
  currentFaq.value = { ...faq };
  showViewModal.value = true;
};

const openCreateModal = () => {
  formData.value = {
    question: '',
    answer: '',
    category: '',
    isPublished: false
  };
  formErrors.value = {
    question: '',
    answer: '',
    category: ''
  };
  isEditing.value = false;
  showEditModal.value = true;
};

const openEditModal = (faq: FAQ) => {
  formData.value = {
    question: faq.question,
    answer: faq.answer,
    category: faq.category,
    isPublished: faq.isPublished
  };
  currentFaq.value = { ...faq };
  formErrors.value = {
    question: '',
    answer: '',
    category: ''
  };
  isEditing.value = true;
  showViewModal.value = false;
  showEditModal.value = true;
};

const cancelEdit = () => {
  showEditModal.value = false;
  showNewCategoryInput.value = false;
};

const validateForm = () => {
  let isValid = true;
  formErrors.value = {
    question: '',
    answer: '',
    category: ''
  };
  
  if (!formData.value.question.trim()) {
    formErrors.value.question = 'La pregunta es obligatoria';
    isValid = false;
  }
  
  if (!formData.value.answer.trim()) {
    formErrors.value.answer = 'La respuesta es obligatoria';
    isValid = false;
  }
  
  if (!formData.value.category) {
    formErrors.value.category = 'La categoría es obligatoria';
    isValid = false;
  }
  
  return isValid;
};

const saveFaq = async () => {
  if (!validateForm()) return;
  
  isSaving.value = true;
  
  try {
    if (isEditing.value) {
      await faqStore.updateFaq({
        id: currentFaq.value.id,
        ...formData.value
      });
    } else {
      await faqStore.addFaq(formData.value);
    }
    
    showEditModal.value = false;
    showNewCategoryInput.value = false;
  } catch (error) {
    console.error('Error al guardar la pregunta frecuente:', error);
  } finally {
    isSaving.value = false;
  }
};

const confirmDelete = (faq: FAQ) => {
  faqToDelete.value = faq.id;
  showDeleteConfirm.value = true;
};

const deleteFaq = async () => {
  if (!faqToDelete.value) return;
  
  isDeleting.value = true;
  
  try {
    await faqStore.deleteFaq(faqToDelete.value);
    showDeleteConfirm.value = false;
  } catch (error) {
    console.error('Error al eliminar la pregunta frecuente:', error);
  } finally {
    isDeleting.value = false;
  }
};

const togglePublish = async (id: number) => {
  try {
    await faqStore.togglePublishStatus(id);
  } catch (error) {
    console.error('Error al cambiar el estado de publicación:', error);
  }
};

const addNewCategory = () => {
  if (newCategory.value.trim()) {
    formData.value.category = newCategory.value.trim();
    showNewCategoryInput.value = false;
    newCategory.value = '';
  }
};
</script>

<style scoped lang="scss">
.admin-section-description {
  color: var(--text-secondary);
  margin-top: 0.5rem;
}

.loading {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  padding: 3rem;
  color: var(--text-secondary);
  
  p {
    margin-top: 1rem;
  }
}

.filter-controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1rem;
  margin-bottom: 2rem;
}

.search-input-container {
  position: relative;
  
  .search-icon {
    position: absolute;
    right: 10px;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-muted);
  }
}

.no-results {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

.question-cell {
  font-weight: 500;
  max-width: 40%;
}

.category-tag {
  display: inline-block;
  background-color: var(--bg-tertiary);
  color: var(--text-secondary);
  padding: 0.3rem 0.6rem;
  border-radius: var(--border-radius);
  font-size: 0.85rem;
}

.status-badge {
  display: inline-flex;
  align-items: center;
  padding: 0.35em 0.65em;
  font-size: 0.75em;
  font-weight: 500;
  line-height: 1;
  text-align: center;
  white-space: nowrap;
  vertical-align: baseline;
  border-radius: 0.25rem;
  
  &.published {
    background-color: var(--success-color);
    color: white;
  }
  
  &.unpublished {
    background-color: var(--text-muted);
    color: white;
  }
}

.actions-cell {
  display: flex;
  gap: 0.5rem;
}

/* Estilos para modales */
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--card-bg);
  border-radius: var(--border-radius);
  width: 90%;
  max-width: 650px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 4px 20px rgba(0, 0, 0, 0.15);
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid var(--border-color);
  
  h2 {
    margin: 0;
    font-size: 1.5rem;
    color: var(--text-primary);
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: var(--text-secondary);
    
    &:hover {
      color: var(--text-primary);
    }
  }
}

.modal-body {
  padding: 1.5rem;
}

.modal-footer {
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  padding: 1rem 1.5rem;
  border-top: 1px solid var(--border-color);
}

/* Estilos para detalles */
.faq-details {
  .detail-section {
    margin-bottom: 2rem;
    
    h3 {
      font-size: 1.1rem;
      font-weight: 600;
      margin-bottom: 0.5rem;
      color: var(--text-primary);
    }
    
    p {
      margin: 0;
      line-height: 1.6;
    }
    
    .answer-content {
      line-height: 1.6;
      white-space: pre-line;
      padding: 1rem;
      background-color: var(--bg-tertiary);
      border-radius: var(--border-radius);
    }
  }
  
  .detail-group {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-bottom: 1.5rem;
    
    .detail-item {
      h4 {
        font-size: 0.9rem;
        font-weight: 600;
        margin-bottom: 0.25rem;
        color: var(--text-secondary);
      }
      
      p {
        margin: 0;
      }
    }
  }
}

/* Estilos para formulario */
.faq-form {
  .form-group {
    margin-bottom: 1.5rem;
  }
  
  .required {
    color: var(--danger-color);
  }
  
  .help-text {
    font-size: 0.85rem;
    color: var(--text-secondary);
    margin-top: 0.5rem;
  }
  
  .error-text {
    color: var(--danger-color);
    font-size: 0.85rem;
    margin-top: 0.25rem;
    display: block;
  }
  
  .is-invalid {
    border-color: var(--danger-color);
  }
  
  .checkbox-container {
    display: flex;
    align-items: center;
    
    input[type="checkbox"] {
      margin-right: 0.5rem;
    }
    
    .checkbox-label {
      font-weight: 500;
    }
  }
}

/* Modal de confirmación */
.confirm-modal {
  max-width: 450px;
  
  .warning-text {
    color: var(--danger-color);
    margin-top: 0.5rem;
  }
}
</style> 
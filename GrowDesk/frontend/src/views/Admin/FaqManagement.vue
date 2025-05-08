/* eslint-disable */
<template>
  <div class="admin-section">
    <!-- Permiso verificado: Solo admins y asistentes -->
    <div v-if="!isAdminOrAssistant" class="permission-error">
      <div class="error-container">
        <i class="pi pi-exclamation-triangle"></i>
        <h2>Acceso restringido</h2>
        <p>No tienes permiso para acceder a esta sección. Esta página solo está disponible para administradores y asistentes.</p>
        <button class="btn btn-primary" @click="goToDashboard">Volver al Dashboard</button>
      </div>
    </div>

    <div v-else>
      <!-- Sección del encabezado con fondo de gradiente y forma ondulada -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">Gestión de Preguntas Frecuentes</h1>
          <p class="hero-subtitle">
            Crear y administrar las preguntas frecuentes que se mostrarán a los clientes.
          </p>
          
          <button class="btn btn-accent" @click="openCreateModal">
            <i class="pi pi-plus"></i> Nueva Pregunta
          </button>
        </div>
        <div class="wave-shape"></div>
      </div>
      
      <div class="content-wrapper">
        <div v-if="loading" class="loading">
          <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
          <p>Cargando preguntas frecuentes...</p>
        </div>
        <div v-else-if="error" class="alert alert-danger">
          <i class="pi pi-exclamation-circle"></i> {{ error }}
        </div>
        
        <div v-else>
          <div class="admin-card">
            <div class="filter-controls">
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
                      <button @click="openViewModal(faq)" class="btn btn-icon" title="Ver detalles">
                        <i class="pi pi-eye"></i>
                      </button>
                      <button @click="openEditModal(faq)" class="btn btn-icon" title="Editar">
                        <i class="pi pi-pencil"></i>
                      </button>
                      <button @click="confirmDelete(faq)" class="btn btn-icon btn-danger" title="Eliminar">
                        <i class="pi pi-trash"></i>
                      </button>
                      <button 
                        @click="togglePublish(faq)" 
                        class="btn btn-icon" 
                        :class="faq.isPublished ? 'btn-warning' : 'btn-success'"
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
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useFaqStore } from '@/stores/faqs';
import { useNotificationsStore } from '@/stores/notifications';
import { useAuthStore } from '@/stores/auth';
import { storeToRefs } from 'pinia';
import type { FAQ } from '@/types/faq.types';

const router = useRouter();
const faqStore = useFaqStore();
const notificationsStore = useNotificationsStore();
const authStore = useAuthStore();
const { faqs, loading, error } = storeToRefs(faqStore);

// Estado
const currentFaq = ref<FAQ | null>(null);
const isEditing = ref(false);
const showEditModal = ref(false);
const showDeleteConfirm = ref(false);
const faqToDelete = ref<number | null>(null);
const searchQuery = ref('');
const categoryFilter = ref('');
const publishFilter = ref('all');
const newCategory = ref('');

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

// Computed properties
const isAdminOrAssistant = computed(() => {
  return authStore.isAdmin || authStore.isAssistant;
});

const categories = computed(() => {
  const storeCategories = faqStore.getCategories();
  return storeCategories;
});

const filteredFaqs = computed(() => {
  let filtered = [...faqs.value];
  
  if (searchQuery.value) {
    const query = searchQuery.value.toLowerCase();
    filtered = filtered.filter(faq => 
      faq.question.toLowerCase().includes(query) ||
      faq.answer.toLowerCase().includes(query)
    );
  }
  
  if (categoryFilter.value) {
    filtered = filtered.filter(faq => faq.category === categoryFilter.value);
  }
  
  if (publishFilter.value !== 'all') {
    const isPublished = publishFilter.value === 'published';
    filtered = filtered.filter(faq => faq.isPublished === isPublished);
  }
  
  return filtered;
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
  showEditModal.value = true;
};

const openCreateModal = () => {
  formData.value = {
    question: '',
    answer: '',
    category: '',
    isPublished: false
  };
  isEditing.value = false;
  showEditModal.value = true;
};

const openEditModal = (faq: FAQ) => {
  formData.value = { ...faq };
  isEditing.value = true;
  showEditModal.value = true;
};

const closeModal = () => {
  showEditModal.value = false;
  formData.value = {
    question: '',
    answer: '',
    category: '',
    isPublished: false
  };
  isEditing.value = false;
};

const cancelEdit = () => {
  showEditModal.value = false;
  formData.value = {
    question: '',
    answer: '',
    category: '',
    isPublished: false
  };
  isEditing.value = false;
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

  try {
    if (isEditing.value && currentFaq.value?.id) {
      await faqStore.updateFaq(currentFaq.value.id, formData.value);
      showSuccess('Pregunta actualizada exitosamente');
    } else {
      await faqStore.addFaq(formData.value);
      showSuccess('Pregunta creada exitosamente');
    }
    
    // Cerrar el modal y limpiar el formulario
    closeModal();
    
    // Recargar la lista de FAQs
    await faqStore.fetchFaqs();
  } catch (error: any) {
    showError(error.message || 'Error al guardar la pregunta');
  }
};

const confirmDelete = (faq: FAQ) => {
  if (faq.id) {
    faqToDelete.value = faq.id;
    showDeleteConfirm.value = true;
  }
};

const deleteFaq = async () => {
  if (!faqToDelete.value) return;
  
  try {
    await faqStore.deleteFaq(faqToDelete.value);
    showDeleteConfirm.value = false;
    faqToDelete.value = null;
    await faqStore.fetchFaqs();
  } catch (error) {
    console.error('Error al eliminar la pregunta frecuente:', error);
  }
};

const togglePublish = async (faq: FAQ) => {
  try {
    await faqStore.togglePublish(faq.id);
    await faqStore.fetchFaqs();
  } catch (error) {
    console.error('Error al cambiar el estado de publicación:', error);
  }
};

const addNewCategory = () => {
  if (newCategory.value && !categories.value.includes(newCategory.value)) {
    categoryFilter.value = newCategory.value;
    newCategory.value = '';
  }
};

const showSuccess = (message: string) => {
  notificationsStore.addNotification({
    type: 'success',
    message: message,
    timeout: 5000
  });
};

const showError = (message: string) => {
  notificationsStore.addNotification({
    type: 'error',
    message: message,
    timeout: 5000
  });
};

// Función para redireccionar al dashboard
function goToDashboard() {
  router.push('/dashboard');
}

// Mostrar notificación si no tiene permiso
onMounted(async () => {
  if (!isAdminOrAssistant.value) {
    notificationsStore.addNotification({
      type: 'error',
      message: 'No tienes permiso para acceder a la gestión de FAQs.',
      timeout: 5000
    });
  } else {
    try {
      await faqStore.fetchFaqs();
    } catch (err: any) {
      console.error('Error al cargar FAQs:', err);
    }
  }
});
</script>

<style scoped lang="scss">
// Estilos comunes para secciones administrativas
.admin-section {
  display: flex;
  flex-direction: column;
  min-height: 100%;
  background-color: var(--bg-secondary);
}

.content-wrapper {
  padding: 2rem;
}

// Hero Section con ondas
.hero-section {
  position: relative;
  background-color: var(--primary-color);
  color: #fff;
  padding: 2.5rem 2rem 6rem;
  text-align: center;
  margin-bottom: 1rem;
  overflow: hidden;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.hero-content {
  position: relative;
  z-index: 2;
  max-width: 1200px;
  margin: 0 auto;
}

.hero-title {
  font-size: 2.25rem;
  font-weight: 700;
  margin-bottom: 0.75rem;
  text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  color:#fff;
}

.hero-subtitle {
  font-size: 1.1rem;
  margin: 0.5rem auto 1.5rem;
  opacity: 0.9;
  max-width: 700px;
}

.wave-shape {
  position: absolute;
  bottom: -2px;
  left: 0;
  width: 100%;
  height: 4rem;
  background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' fill='%23f8fafc' opacity='.25'%3E%3C/path%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' fill='%23f8fafc' opacity='.5'%3E%3C/path%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%23f8fafc'%3E%3C/path%3E%3C/svg%3E");
  background-size: cover;
  background-position: center;
}

// Card para el contenedor principal
.admin-card {
  background-color: var(--card-bg);
  border-radius: var(--border-radius-lg);
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
}

// Loading state
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

// Filtros
.filter-controls {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
  gap: 1.5rem;
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

.form-group {
  margin-bottom: 1.25rem;
  
  .form-label {
    display: block;
    margin-bottom: 0.5rem;
    font-weight: 500;
    color: var(--text-primary);
  }
  
  .form-control {
    width: 100%;
    padding: 0.65rem 1rem;
    background-color: var(--input-bg);
    border: 1px solid var(--border-color);
    border-radius: var(--border-radius);
    font-size: 1rem;
    transition: border-color 0.2s, box-shadow 0.2s;
    color: var(--text-primary);
    
    &:focus {
      border-color: var(--primary-color);
      box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.1);
      outline: none;
    }
    
    &::placeholder {
      color: var(--text-muted);
    }
    
    &.is-invalid {
      border-color: var(--danger-color);
    }
  }
  
  .error-text {
    color: var(--danger-color);
    font-size: 0.85rem;
    margin-top: 0.4rem;
    display: block;
  }
}

// Tabla de FAQs
.admin-table {
  width: 100%;
  border-collapse: separate;
  border-spacing: 0;
  
  th, td {
    padding: 1rem;
    text-align: left;
    border-bottom: 1px solid var(--border-color);
  }
  
  th {
    font-weight: 600;
    color: var(--text-primary);
    background-color: var(--bg-tertiary);
    position: sticky;
    top: 0;
    z-index: 10;
    
    &:first-child {
      border-top-left-radius: var(--border-radius);
    }
    
    &:last-child {
      border-top-right-radius: var(--border-radius);
    }
  }
  
  tr:last-child td {
    border-bottom: none;
  }
  
  tbody tr {
    transition: background-color 0.2s;
    
    &:hover {
      background-color: var(--bg-hover);
    }
  }
}

.no-results {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
}

.question-cell {
  font-weight: 500;
  max-width: 350px;
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
  gap: 0.4rem;
}

// Botones
.btn {
  display: inline-flex;
  align-items: center;
  justify-content: center;
  gap: 0.5rem;
  border-radius: var(--border-radius);
  padding: 0.65rem 1.5rem;
  font-size: 1rem;
  font-weight: 500;
  cursor: pointer;
  transition: all 0.2s;
  border: none;
  
  &.btn-primary {
    background-color: var(--primary-color);
    color: white;
    
    &:hover, &:focus {
      background-color: var(--primary-color-dark);
      box-shadow: 0 4px 8px rgba(var(--primary-rgb), 0.25);
    }
  }
  
  &.btn-accent {
      display: inline-flex;
      align-items: center;
      gap: 0.75rem;
      background-color: white;
      color: var(--primary-color);
      border: none;
      padding: 0.85rem 1.75rem;
      border-radius: 8px;
      font-size: 1rem;
      font-weight: 600;
      text-decoration: none;
      transition: all 0.3s ease;
      box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
      cursor: pointer;
      
      &:hover {
        transform: translateY(-3px);
        box-shadow: 0 8px 16px rgba(0, 0, 0, 0.2);
      }
      
      i {
        font-size: 1rem;
      }
    }
  
  &.btn-danger {
    background-color: var(--danger-color);
    color: white;
    
    &:hover, &:focus {
      background-color: var(--danger-color-dark);
    }
  }
  
  &.btn-warning {
    background-color: var(--warning-color);
    color: white;
    
    &:hover, &:focus {
      background-color: var(--warning-color-dark);
    }
  }
  
  &.btn-success {
    background-color: var(--success-color);
    color: white;
    
    &:hover, &:focus {
      background-color: var(--success-color-dark);
    }
  }
  
  &.btn-outline-secondary {
    background-color: transparent;
    color: var(--text-secondary);
    border: 1px solid var(--border-color);
    
    &:hover, &:focus {
      background-color: var(--bg-tertiary);
      border-color: var(--text-secondary);
    }
  }
  
  &.btn-icon {
    padding: 0.5rem;
    border-radius: 50%;
    background-color: transparent;
    color: var(--text-secondary);
    
    &:hover {
      background-color: var(--bg-tertiary);
      color: var(--text-primary);
    }
    
    &.btn-danger {
      color: var(--danger-color);
      
      &:hover {
        background-color: rgba(var(--danger-rgb), 0.1);
      }
    }
    
    &.btn-warning {
      color: var(--warning-color);
      
      &:hover {
        background-color: rgba(var(--warning-rgb), 0.1);
      }
    }
    
    &.btn-success {
      color: var(--success-color);
      
      &:hover {
        background-color: rgba(var(--success-rgb), 0.1);
      }
    }
  }
}

// Modales
.modal-overlay {
  position: fixed;
  top: 0;
  left: 0;
  right: 0;
  bottom: 0;
  background-color: rgba(0, 0, 0, 0.5);
  backdrop-filter: blur(4px);
  display: flex;
  align-items: center;
  justify-content: center;
  z-index: 1000;
}

.modal-content {
  background-color: var(--card-bg);
  border-radius: var(--border-radius-lg);
  width: 90%;
  max-width: 650px;
  max-height: 90vh;
  overflow-y: auto;
  box-shadow: 0 8px 30px rgba(0, 0, 0, 0.15);
  
  &.confirm-modal {
    max-width: 450px;
  }
}

.modal-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1.25rem 1.5rem;
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
  padding: 1.25rem 1.5rem;
  border-top: 1px solid var(--border-color);
}

// Formulario en modal
.faq-form {
  .checkbox-container {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 0.5rem;
    
    input[type="checkbox"] {
      width: 18px;
      height: 18px;
    }
    
    .checkbox-label {
      cursor: pointer;
      user-select: none;
    }
  }
  
  .help-text {
    margin: 0.5rem 0 0;
    font-size: 0.9rem;
    color: var(--text-secondary);
  }
  
  .required {
    color: var(--danger-color);
    margin-left: 2px;
  }
}

// Detalles de FAQ
.faq-details {
  .detail-section {
    margin-bottom: 1.5rem;
    
    h3 {
      font-size: 1.1rem;
      font-weight: 600;
      margin-bottom: 0.75rem;
      color: var(--text-primary);
      padding-bottom: 0.5rem;
      border-bottom: 1px solid var(--border-color);
    }
    
    p {
      margin: 0;
      line-height: 1.6;
    }
    
    .answer-content {
      white-space: pre-wrap;
      line-height: 1.6;
    }
  }
  
  .detail-group {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(200px, 1fr));
    gap: 1.5rem;
    margin-bottom: 1.5rem;
    
    .detail-item {
      h4 {
        font-size: 0.9rem;
        font-weight: 600;
        margin: 0 0 0.5rem;
        color: var(--text-secondary);
      }
      
      p {
        margin: 0;
      }
    }
  }
}

.warning-text {
  color: var(--danger-color);
  font-weight: 500;
}

/* Estilos para el mensaje de error de permisos */
.permission-error {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 120px);
  padding: 2rem;
}

.error-container {
  max-width: 500px;
  text-align: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.error-container i {
  font-size: 3rem;
  color: #f59e0b;
  margin-bottom: 1rem;
}

.error-container h2 {
  margin-bottom: 1rem;
  color: #334155;
}

.error-container p {
  margin-bottom: 2rem;
  color: #64748b;
}
</style>
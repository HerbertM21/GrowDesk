/* eslint-disable */
<template>
  <div class="admin-section">
    <div v-if="!isAdmin" class="authorization-error">
      <div class="alert alert-danger">
        <i class="pi pi-exclamation-triangle"></i>
        <span>No tienes permisos suficientes para acceder a esta página</span>
      </div>
      <button @click="redirectToDashboard" class="btn btn-primary">
        <i class="pi pi-home"></i> Volver al Dashboard
      </button>
    </div>
    
    <template v-else>
      <!-- Sección del encabezado con fondo de gradiente y forma ondulada -->
      <div class="hero-section">
        <div class="hero-content">
          <h1 class="hero-title">Gestión de Categorías</h1>
          <p class="hero-subtitle">Administre las categorías para la clasificación de tickets</p>
        </div>
        <div class="wave-shape"></div>
      </div>
      
      <div class="content-wrapper">
        <div v-if="loading" class="loading">
          <i class="pi pi-spin pi-spinner" style="font-size: 2rem"></i>
          <p>Cargando categorías...</p>
        </div>
        <div v-else-if="error" class="alert alert-danger">
          <i class="pi pi-exclamation-circle"></i> {{ error }}
        </div>
        
        <div v-else class="admin-form-row">
          <div class="admin-form-col">
            <div class="admin-card">
              <div class="section-title">
                <div class="title-icon">
                  <i class="pi pi-tags"></i>
                </div>
                <h2>Categorías Existentes</h2>
              </div>
              <div v-if="categories.length === 0" class="empty-list">
                <p>No hay categorías disponibles</p>
              </div>
              <ul v-else class="category-list">
                <li v-for="category in categories" :key="category.id" class="category-item">
                  <div class="category-info">
                    <span class="category-name">{{ category.name }}</span>
                    <div class="category-actions">
                      <button @click="editCategory(category)" class="btn btn-sm btn-outline-primary">
                        <i class="pi pi-pencil"></i>
                      </button>
                      <button @click="deleteCategory(category.id)" class="btn btn-sm btn-outline-danger">
                        <i class="pi pi-trash"></i>
                      </button>
                    </div>
                  </div>
                  <p v-if="category.description" class="category-description">
                    {{ category.description }}
                  </p>
                </li>
              </ul>
            </div>
          </div>
          
          <div class="admin-form-col">
            <div class="admin-card">
              <div class="section-title">
                <div class="title-icon">
                  <i class="pi pi-plus-circle"></i>
                </div>
                <h2>{{ isEditing ? 'Editar Categoría' : 'Añadir Nueva Categoría' }}</h2>
              </div>
              <form @submit.prevent="saveCategory" class="admin-form">
                <div class="form-group">
                  <label for="categoryName" class="form-label">Nombre</label>
                  <input 
                    type="text" 
                    id="categoryName" 
                    v-model="currentCategory.name" 
                    required
                    class="form-control"
                    placeholder="Nombre de la categoría"
                  />
                </div>
                
                <div class="form-group">
                  <label for="categoryDescription" class="form-label">Descripción</label>
                  <textarea 
                    id="categoryDescription" 
                    v-model="currentCategory.description" 
                    class="form-control"
                    rows="4"
                    placeholder="Descripción de la categoría"
                  ></textarea>
                </div>
                
                <div class="form-actions">
                  <button type="submit" class="btn btn-primary">
                    <i class="pi" :class="isEditing ? 'pi-check' : 'pi-plus'"></i>
                    {{ isEditing ? 'Actualizar' : 'Guardar' }}
                  </button>
                  <button v-if="isEditing" @click="cancelEdit" type="button" class="btn btn-outline-secondary">
                    <i class="pi pi-times"></i> Cancelar
                  </button>
                </div>
              </form>
            </div>
          </div>
        </div>
      </div>
    </template>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useCategoriesStore } from '@/stores/categories';
import { storeToRefs } from 'pinia';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useNotificationsStore } from '@/stores/notifications';

// Store
const categoryStore = useCategoriesStore();
const { categories, loading, error } = storeToRefs(categoryStore);
const authStore = useAuthStore();
const notificationsStore = useNotificationsStore();
const router = useRouter();

// Verificar permisos
const isAdmin = computed(() => {
  return authStore.isAdmin;
});

// Función para redirigir al dashboard
const redirectToDashboard = () => {
  notificationsStore.warning('Has sido redirigido al Dashboard debido a permisos insuficientes');
  router.push('/dashboard');
};

// Estado local
const isEditing = ref(false);
const currentCategory = ref({
  id: null,
  name: '',
  description: ''
});

// Cargar categorías al montar el componente
onMounted(() => {
  console.log('CategoriesManagement: componente montado')
  // Solo cargamos datos si el usuario es administrador
  if (isAdmin.value) {
    console.log('CategoriesManagement: usuario es admin, cargando categorías')
    categoryStore.fetchCategories().then(() => {
      console.log('CategoriesManagement: categorías cargadas:', categoryStore.categories)
    })
  } else {
    // Si no es admin, redirigimos al dashboard
    console.log('CategoriesManagement: usuario no es admin, redirigiendo al dashboard')
    redirectToDashboard();
  }
});

// Métodos
const editCategory = (category: { id: number, name: string, description: string }) => {
  console.log('CategoriesManagement: editando categoría:', category)
  currentCategory.value = { ...category };
  isEditing.value = true;
};

const cancelEdit = () => {
  console.log('CategoriesManagement: cancelando edición')
  currentCategory.value = {
    id: null,
    name: '',
    description: ''
  };
  isEditing.value = false;
};

const saveCategory = async () => {
  console.log('CategoriesManagement: guardando categoría:', currentCategory.value)
  try {
    if (isEditing.value) {
      console.log('CategoriesManagement: actualizando categoría existente')
      await categoryStore.updateCategory(currentCategory.value);
      console.log('CategoriesManagement: categoría actualizada')
    } else {
      console.log('CategoriesManagement: añadiendo nueva categoría')
      await categoryStore.addCategory({
        name: currentCategory.value.name,
        description: currentCategory.value.description
      });
      console.log('CategoriesManagement: categoría añadida')
    }
    // Resetear formulario
    cancelEdit();
    // Mostrar notificación de éxito
    notificationsStore.success(`Categoría ${isEditing.value ? 'actualizada' : 'creada'} correctamente`);
  } catch (error) {
    console.error('Error al guardar categoría:', error);
    notificationsStore.error(`Error al ${isEditing.value ? 'actualizar' : 'crear'} la categoría`);
  }
};

const deleteCategory = async (id: number) => {
  console.log('CategoriesManagement: confirmando eliminación de categoría:', id)
  if (confirm('¿Está seguro de que desea eliminar esta categoría?')) {
    try {
      console.log('CategoriesManagement: eliminando categoría:', id)
      await categoryStore.deleteCategory(id);
      console.log('CategoriesManagement: categoría eliminada')
      notificationsStore.success('Categoría eliminada correctamente');
    } catch (error) {
      console.error('Error al eliminar categoría:', error);
      notificationsStore.error('Error al eliminar la categoría');
    }
  }
};
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
  color: white;
}

.hero-subtitle {
  font-size: 1.1rem;
  margin-top: 0.5rem;
  opacity: 0.9;
  max-width: 600px;
  margin-left: auto;
  margin-right: auto;
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

// Layout para el contenido principal
.admin-form-row {
  display: grid;
  grid-template-columns: repeat(auto-fit, minmax(350px, 1fr));
  gap: 2rem;
  margin-bottom: 2rem;
}

.admin-form-col {
  display: flex;
  flex-direction: column;
}

// Card para las secciones
.admin-card {
  background-color: var(--card-bg);
  border-radius: var(--border-radius-lg);
  padding: 1.5rem;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.08);
  height: 100%;
}

.admin-card-title {
  font-size: 1.25rem;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
  font-weight: 600;
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

// Lista de categorías
.category-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.category-item {
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
  transition: all 0.2s ease;
  
  &:hover {
    background-color: rgba(var(--primary-rgb), 0.05);
  }
  
  &:last-child {
    border-bottom: none;
  }
}

.category-info {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.category-name {
  font-weight: 600;
  font-size: 1.1rem;
  color: var(--text-primary);
}

.category-description {
  margin-top: 0.5rem;
  font-size: 0.9rem;
  color: var(--text-secondary);
  line-height: 1.5;
}

.category-actions {
  display: flex;
  gap: 0.5rem;
}

// Formulario
.admin-form {
  .form-group {
    margin-bottom: 1.5rem;
  }
  
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
  }
  
  textarea.form-control {
    resize: vertical;
  }
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}

// Empty state
.empty-list {
  text-align: center;
  padding: 3rem;
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  border-radius: var(--border-radius);
  
  p {
    font-size: 1.1rem;
    margin: 0;
  }
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
  
  &.btn-primary {
    background-color: var(--primary-color);
    color: white;
    border: none;
    
    &:hover, &:focus {
      background-color: var(--primary-color-dark);
      box-shadow: 0 4px 8px rgba(var(--primary-rgb), 0.25);
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
  
  &.btn-outline-danger {
    background-color: transparent;
    color: var(--danger-color);
    border: 1px solid var(--border-color);
    
    &:hover, &:focus {
      background-color: rgba(220, 53, 69, 0.1);
      border-color: var(--danger-color);
    }
  }
  
  &.btn-sm {
    padding: 0.35rem 0.65rem;
    font-size: 0.85rem;
  }
  
  i {
    font-size: 0.9em;
  }
}

.section-title {
  display: flex;
  align-items: center;
  justify-content: flex-start;
  margin-bottom: 1rem;
  color: var(--text-primary);
  font-size: 1.25rem;
  font-weight: 600;
  text-align: left;
  background-color: var(--bg-tertiary);
  border-radius: 10px;
  padding: 0.5rem 1rem;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  border-left: 4px solid var(--primary-color);
  
  .title-icon {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    width: 32px;
    height: 32px;
    background-color: var(--primary-color);
    border-radius: 8px;
    margin-right: 0.75rem;
    color: white;
    box-shadow: 0 4px 10px rgba(var(--primary-rgb), 0.25);
    
    i {
      font-size: 1rem;
    }
  }
  
  h2 {
    margin: 0;
    font-size: 1.1rem;
  }
}
</style>
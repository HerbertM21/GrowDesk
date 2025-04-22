<template>
  <div class="admin-section">
    <div class="admin-section-header">
      <h1 class="admin-section-title">Gestión de Categorías</h1>
    </div>
    
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
          <h2 class="admin-card-title">Categorías Existentes</h2>
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
          <h2 class="admin-card-title">{{ isEditing ? 'Editar Categoría' : 'Añadir Nueva Categoría' }}</h2>
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

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useCategoriesStore } from '@/stores/categories';
import { storeToRefs } from 'pinia';

// Store
const categoryStore = useCategoriesStore();
const { categories, loading, error } = storeToRefs(categoryStore);

// Estado local
const isEditing = ref(false);
const currentCategory = ref({
  id: null,
  name: '',
  description: ''
});

// Cargar categorías al montar el componente
onMounted(() => {
  categoryStore.fetchCategories();
});

// Métodos
const editCategory = (category: { id: number, name: string, description: string }) => {
  currentCategory.value = { ...category };
  isEditing.value = true;
};

const cancelEdit = () => {
  currentCategory.value = {
    id: null,
    name: '',
    description: ''
  };
  isEditing.value = false;
};

const saveCategory = async () => {
  try {
    if (isEditing.value) {
      await categoryStore.updateCategory(currentCategory.value);
    } else {
      await categoryStore.addCategory({
        name: currentCategory.value.name,
        description: currentCategory.value.description
      });
    }
    // Resetear formulario
    cancelEdit();
  } catch (error) {
    console.error('Error al guardar categoría:', error);
  }
};

const deleteCategory = async (id: number) => {
  if (confirm('¿Está seguro de que desea eliminar esta categoría?')) {
    try {
      await categoryStore.deleteCategory(id);
    } catch (error) {
      console.error('Error al eliminar categoría:', error);
    }
  }
};
</script>

<style scoped lang="scss">
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

.admin-card-title {
  font-size: 1.25rem;
  margin-bottom: 1.5rem;
  padding-bottom: 0.75rem;
  border-bottom: 1px solid var(--border-color);
  color: var(--text-primary);
}

.empty-list {
  text-align: center;
  padding: 2rem;
  color: var(--text-secondary);
  background-color: var(--bg-tertiary);
  border-radius: var(--border-radius);
}

.category-list {
  list-style: none;
  padding: 0;
  margin: 0;
}

.category-item {
  padding: 1rem;
  border-bottom: 1px solid var(--border-color);
  transition: background-color 0.2s;
  
  &:hover {
    background-color: var(--hover-bg);
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
}

.category-actions {
  display: flex;
  gap: 0.5rem;
}

.form-actions {
  display: flex;
  gap: 1rem;
  margin-top: 2rem;
}
</style> 
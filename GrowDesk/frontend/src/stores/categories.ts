import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface Category {
  id: number
  name: string
  description: string
}

// Nombre clave para localStorage
const STORAGE_KEY = 'growdesk-categories'

export const useCategoriesStore = defineStore('categories', () => {
  const categories = ref<Category[]>([])
  const loading = ref(false)
  const error = ref('')

  // Categorías de ejemplo (mock)
  const mockCategories: Category[] = [
    { id: 1, name: 'Soporte Técnico', description: 'Problemas técnicos y asistencia' },
    { id: 2, name: 'Ventas', description: 'Consultas sobre productos y servicios' },
    { id: 3, name: 'Facturación', description: 'Problemas con pagos y facturas' },
    { id: 4, name: 'General', description: 'Consultas generales' }
  ]

  // Guardar categorías en localStorage
  function saveCategoriesToLocalStorage() {
    console.log('Guardando categorías en localStorage:', categories.value)
    localStorage.setItem(STORAGE_KEY, JSON.stringify(categories.value))
  }

  // Cargar categorías desde localStorage
  function loadCategoriesFromLocalStorage(): Category[] {
    try {
      // Limpiar datos potencialmente inválidos
      import('@/utils/validators').then(({ filterValidCategories }) => {
        try {
          const data = localStorage.getItem(STORAGE_KEY);
          if (data) {
            const parsed = JSON.parse(data);
            if (Array.isArray(parsed)) {
              const valid = filterValidCategories(parsed);
              if (valid.length !== parsed.length) {
                localStorage.setItem(STORAGE_KEY, JSON.stringify(valid));
                console.warn(`Se limpiaron ${parsed.length - valid.length} categorías inválidas`);
              }
            }
          }
        } catch (e) {
          console.error('Error al limpiar categorías:', e);
        }
      }).catch(e => console.error('Error al importar utilidades:', e));
      
      const storedData = localStorage.getItem(STORAGE_KEY)
      if (storedData) {
        console.log('Categorías cargadas desde localStorage')
        const parsed = JSON.parse(storedData);
        // Verificar que es un array válido
        if (!Array.isArray(parsed)) {
          console.error('Datos de categorías no son un array válido');
          return [];
        }
        // Filtrar para asegurar que solo tenemos objetos válidos
        const validCategories = parsed.filter(item => 
          typeof item === 'object' && 
          item !== null && 
          'id' in item && 
          'name' in item
        );
        if (validCategories.length !== parsed.length) {
          console.warn(`Se filtraron ${parsed.length - validCategories.length} categorías inválidas`);
          localStorage.setItem(STORAGE_KEY, JSON.stringify(validCategories));
        }
        return validCategories;
      }
    } catch (err) {
      console.error('Error al cargar categorías desde localStorage:', err)
    }
    console.log('No se encontraron categorías en localStorage, usando datos iniciales')
    return []
  }

  // Inicializar el store
  function initializeStore() {
    // Intentar cargar desde localStorage primero
    const storedCategories = loadCategoriesFromLocalStorage()
    
    // Si no hay categorías almacenadas, usar las categorías de ejemplo
    if (storedCategories.length === 0) {
      categories.value = [...mockCategories]
      // Guardar las categorías iniciales en localStorage
      saveCategoriesToLocalStorage()
    } else {
      categories.value = storedCategories
    }
  }

  // Cargar categorías
  async function fetchCategories() {
    loading.value = true
    error.value = ''
    
    try {
      // Intentar cargar desde localStorage primero
      const storedCategories = loadCategoriesFromLocalStorage()
      
      // Si no hay categorías almacenadas, usar las categorías de ejemplo
      if (storedCategories.length === 0) {
        categories.value = [...mockCategories]
        // Guardar las categorías iniciales en localStorage
        saveCategoriesToLocalStorage()
      } else {
        categories.value = storedCategories
      }
    } catch (err) {
      error.value = 'Error al cargar las categorías'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Añadir una categoría
  async function addCategory(category: Omit<Category, 'id'>) {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      
      // Generar un ID único (en una aplicación real, el backend lo generaría)
      const newId = categories.value.length > 0 
        ? Math.max(...categories.value.map((c: Category) => c.id)) + 1 
        : 1
        
      const newCategory: Category = {
        id: newId,
        ...category
      }
      
      // Añadir al array local
      categories.value.push(newCategory)
      
      // Guardar en localStorage
      saveCategoriesToLocalStorage()
      
      return newCategory
    } catch (err) {
      error.value = 'Error al añadir la categoría'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Actualizar una categoría
  async function updateCategory(category: Category) {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      
      const index = categories.value.findIndex((c: Category) => c.id === category.id)
      if (index !== -1) {
        categories.value[index] = { ...category }
        
        // Guardar en localStorage
        saveCategoriesToLocalStorage()
        
        return category
      } else {
        throw new Error('Categoría no encontrada')
      }
    } catch (err) {
      error.value = 'Error al actualizar la categoría'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Eliminar una categoría
  async function deleteCategory(id: number) {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      
      // Filtrar categorías y actualizar el estado
      categories.value = categories.value.filter((c: Category) => c.id !== id)
      
      // Guardar en localStorage
      saveCategoriesToLocalStorage()
    } catch (err) {
      error.value = 'Error al eliminar la categoría'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Inicializar automáticamente el store al importarlo
  initializeStore()

  return {
    categories,
    loading,
    error,
    fetchCategories,
    addCategory,
    updateCategory,
    deleteCategory,
    initializeStore
  }
}) 
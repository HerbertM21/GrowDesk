import { defineStore } from 'pinia'
import { ref } from 'vue'

interface Category {
  id: number
  name: string
  description: string
}

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

  // Inicializar el store
  function initializeStore() {
    //  aquí cargaríamos las categorías desde el backend
    categories.value = [...mockCategories]
  }

  // Cargar categorías
  async function fetchCategories() {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API con una promesa
      await new Promise(resolve => setTimeout(resolve, 500))
      categories.value = [...mockCategories]
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
      
      categories.value.push(newCategory)
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
      
      categories.value = categories.value.filter((c: Category) => c.id !== id)
    } catch (err) {
      error.value = 'Error al eliminar la categoría'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

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
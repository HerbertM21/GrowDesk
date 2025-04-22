import { defineStore } from 'pinia'
import { ref } from 'vue'

export interface FAQ {
  id: number
  question: string
  answer: string
  category: string
  isPublished: boolean
  createdAt: string
  updatedAt: string
}

export const useFaqsStore = defineStore('faqs', () => {
  const faqs = ref<FAQ[]>([])
  const loading = ref(false)
  const error = ref('')

  // FAQs de ejemplo (mock)
  const mockFaqs: FAQ[] = [
    {
      id: 1,
      question: '¿Cómo puedo restablecer mi contraseña?',
      answer: 'Para restablecer su contraseña, haga clic en el enlace "¿Olvidó su contraseña?" en la pantalla de inicio de sesión y siga las instrucciones enviadas a su correo electrónico.',
      category: 'Cuenta',
      isPublished: true,
      createdAt: '2023-01-15T10:30:00Z',
      updatedAt: '2023-01-15T10:30:00Z'
    },
    {
      id: 2,
      question: '¿Cómo puedo actualizar mi información de contacto?',
      answer: 'Inicie sesión en su cuenta, vaya a "Mi Perfil" y haga clic en "Editar Información". Allí podrá actualizar su dirección de correo electrónico, número de teléfono y dirección postal.',
      category: 'Cuenta',
      isPublished: true,
      createdAt: '2023-01-20T14:15:00Z',
      updatedAt: '2023-02-05T09:45:00Z'
    },
    {
      id: 3,
      question: '¿Cómo puedo reportar un problema técnico?',
      answer: 'Para reportar un problema técnico, vaya a la sección "Soporte", haga clic en "Crear Ticket" y complete el formulario con los detalles del problema. Un técnico se pondrá en contacto con usted lo antes posible.',
      category: 'Soporte Técnico',
      isPublished: true,
      createdAt: '2023-02-10T11:20:00Z',
      updatedAt: '2023-02-10T11:20:00Z'
    },
    {
      id: 4,
      question: '¿Cuáles son los horarios de atención al cliente?',
      answer: 'Nuestro equipo de atención al cliente está disponible de lunes a viernes, de 9:00 a.m. a 6:00 p.m. (hora local). Para asistencia de emergencia fuera de este horario, por favor envíe un correo electrónico a soporte@ejemplo.com.',
      category: 'General',
      isPublished: false,
      createdAt: '2023-03-01T16:40:00Z',
      updatedAt: '2023-03-10T08:30:00Z'
    }
  ]

  // Obtener todas las categorías disponibles
  const getCategories = () => {
    const uniqueCategories = new Set<string>()
    faqs.value.forEach(faq => uniqueCategories.add(faq.category))
    return Array.from(uniqueCategories).sort()
  }

  // Inicializar el store
  function initializeStore() {
    faqs.value = [...mockFaqs]
  }

  // Cargar FAQs
  async function fetchFaqs() {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      faqs.value = [...mockFaqs]
    } catch (err) {
      error.value = 'Error al cargar las preguntas frecuentes'
      console.error(err)
    } finally {
      loading.value = false
    }
  }

  // Añadir una FAQ
  async function addFaq(faq: Omit<FAQ, 'id' | 'createdAt' | 'updatedAt'>) {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      
      const now = new Date().toISOString()
      const newId = faqs.value.length > 0 
        ? Math.max(...faqs.value.map((f: FAQ) => f.id)) + 1 
        : 1
        
      const newFaq: FAQ = {
        id: newId,
        createdAt: now,
        updatedAt: now,
        ...faq
      }
      
      faqs.value.push(newFaq)
      return newFaq
    } catch (err) {
      error.value = 'Error al añadir la pregunta frecuente'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Actualizar una FAQ
  async function updateFaq(faq: Omit<FAQ, 'createdAt' | 'updatedAt'>) {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      
      const now = new Date().toISOString()
      const index = faqs.value.findIndex((f: FAQ) => f.id === faq.id)
      
      if (index !== -1) {
        const updatedFaq = {
          ...faq,
          updatedAt: now,
          createdAt: faqs.value[index].createdAt
        }
        
        faqs.value[index] = updatedFaq
        return updatedFaq
      } else {
        throw new Error('Pregunta frecuente no encontrada')
      }
    } catch (err) {
      error.value = 'Error al actualizar la pregunta frecuente'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Eliminar una FAQ
  async function deleteFaq(id: number) {
    loading.value = true
    error.value = ''
    
    try {
      // Simulamos la llamada a la API
      await new Promise(resolve => setTimeout(resolve, 500))
      
      faqs.value = faqs.value.filter((f: FAQ) => f.id !== id)
    } catch (err) {
      error.value = 'Error al eliminar la pregunta frecuente'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  // Cambiar el estado de publicación
  async function togglePublishStatus(id: number) {
    loading.value = true
    error.value = ''
    
    try {
      const index = faqs.value.findIndex((f: FAQ) => f.id === id)
      
      if (index !== -1) {
        const faq = { ...faqs.value[index] }
        faq.isPublished = !faq.isPublished
        faq.updatedAt = new Date().toISOString()
        
        // Simulamos la llamada a la API
        await new Promise(resolve => setTimeout(resolve, 300))
        
        faqs.value[index] = faq
        return faq
      } else {
        throw new Error('Pregunta frecuente no encontrada')
      }
    } catch (err) {
      error.value = 'Error al cambiar el estado de publicación'
      console.error(err)
      throw err
    } finally {
      loading.value = false
    }
  }

  return {
    faqs,
    loading,
    error,
    fetchFaqs,
    addFaq,
    updateFaq,
    deleteFaq,
    togglePublishStatus,
    getCategories,
    initializeStore
  }
}) 
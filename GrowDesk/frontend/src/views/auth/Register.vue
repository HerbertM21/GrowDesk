<template>
  <div class="register-container">
    <div class="register-card">
      <div class="logo-container">
        <img src="@/assets/logo.png" alt="GrowDesk Logo" class="login-logo">
      </div>
      <h2>Registro</h2>
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      <form @submit.prevent="handleRegister">
        <div class="form-group">
          <label for="firstName">Nombre</label>
          <input type="text" id="firstName" v-model="firstName" :class="{'input-error': validationErrors.firstName}" required>
          <small v-if="validationErrors.firstName" class="error-text">{{ validationErrors.firstName }}</small>
        </div>
        <div class="form-group">
          <label for="lastName">Apellido</label>
          <input type="text" id="lastName" v-model="lastName" :class="{'input-error': validationErrors.lastName}" required>
          <small v-if="validationErrors.lastName" class="error-text">{{ validationErrors.lastName }}</small>
        </div>
        <div class="form-group">
          <label for="email">Correo Electrónico</label>
          <input type="email" id="email" v-model="email" :class="{'input-error': validationErrors.email}" required>
          <small v-if="validationErrors.email" class="error-text">{{ validationErrors.email }}</small>
        </div>
        <div class="form-group">
          <label for="password">Contraseña</label>
          <input type="password" id="password" v-model="password" :class="{'input-error': validationErrors.password}" required minlength="6">
          <small v-if="validationErrors.password" class="error-text">{{ validationErrors.password }}</small>
          <small v-else>La contraseña debe tener al menos 6 caracteres</small>
        </div>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Registrando...' : 'Registrarse' }}
        </button>
      </form>
      <p class="login-link">
        ¿Ya tienes una cuenta? <router-link to="/login">Iniciar Sesión</router-link>
      </p>
      <div v-if="debugInfo.length > 0" class="debug-info">
        <h3>Información de depuración:</h3>
        <ul>
          <li v-for="(item, index) in debugInfo" :key="index">{{ item }}</li>
        </ul>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'
import apiClient from '@/api/client'

const router = useRouter()
const authStore = useAuthStore()
const firstName = ref('')
const lastName = ref('')
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = computed(() => authStore.loading)
const debugInfo = ref<string[]>([])

const validationErrors = reactive({
  firstName: '',
  lastName: '',
  email: '',
  password: ''
})

// Función para añadir información de depuración
const addDebugInfo = (info: string) => {
  console.log('DEBUG:', info)
  debugInfo.value.push(`${new Date().toISOString()}: ${info}`)
}

// Función para validar el formulario
const validateForm = () => {
  let isValid = true
  
  // Resetear errores
  validationErrors.firstName = ''
  validationErrors.lastName = ''
  validationErrors.email = ''
  validationErrors.password = ''
  
  // Validar nombre
  if (!firstName.value.trim()) {
    validationErrors.firstName = 'El nombre es obligatorio'
    isValid = false
  }
  
  // Validar apellido
  if (!lastName.value.trim()) {
    validationErrors.lastName = 'El apellido es obligatorio'
    isValid = false
  }
  
  // Validar email
  if (!email.value.trim()) {
    validationErrors.email = 'El correo electrónico es obligatorio'
    isValid = false
  } else if (!/^[^\s@]+@[^\s@]+\.[^\s@]+$/.test(email.value)) {
    validationErrors.email = 'Formato de correo electrónico inválido'
    isValid = false
  }
  
  // Validar contraseña
  if (!password.value) {
    validationErrors.password = 'La contraseña es obligatoria'
    isValid = false
  } else if (password.value.length < 6) {
    validationErrors.password = 'La contraseña debe tener al menos 6 caracteres'
    isValid = false
  }
  
  return isValid
}

const handleRegister = async () => {
  try {
    // Limpiar la información de depuración para esta nueva solicitud
    debugInfo.value = []
    
    // Resetear mensaje de error global
    errorMessage.value = ''
    
    // Validar formulario
    if (!validateForm()) {
      addDebugInfo('Validación del formulario fallida')
      return
    }
    
    addDebugInfo('Validación del formulario exitosa')
    addDebugInfo(`API URL: ${import.meta.env.VITE_API_URL || '/api'}`)
    
    const userData = {
      firstName: firstName.value,
      lastName: lastName.value,
      email: email.value,
      password: password.value
    }
    
    addDebugInfo(`Enviando datos de registro para: ${email.value}`)
    
    try {
      addDebugInfo('Llamando a authStore.register()')
      const success = await authStore.register(userData)
      
      if (success) {
        addDebugInfo('Registro exitoso, redirigiendo al dashboard')
        router.push('/dashboard')
      } else {
        addDebugInfo('Registro fallido: Respuesta success=false')
        errorMessage.value = 'No se pudo completar el registro. Por favor intenta de nuevo.'
      }
    } catch (error: any) {
      addDebugInfo(`Error en el registro: ${error.message || 'Error desconocido'}`)
      console.error('Error detallado durante el registro:', error)
      
      if (error.response) {
        // El servidor respondió con un código de estado diferente de 2xx
        const statusCode = error.response.status
        const responseData = error.response.data
        
        addDebugInfo(`Error del servidor (${statusCode}): ${JSON.stringify(responseData)}`)
        console.error(`Error de servidor (${statusCode}):`, responseData)
        
        if (statusCode === 409) {
          errorMessage.value = 'Este correo electrónico ya está registrado. Por favor utiliza otro o inicia sesión.'
        } else if (responseData && responseData.error) {
          errorMessage.value = responseData.error
        } else {
          errorMessage.value = `Error en el registro (${statusCode}). Por favor intenta de nuevo.`
        }
      } else if (error.request) {
        // La solicitud se realizó pero no se recibió respuesta
        addDebugInfo('Sin respuesta del servidor')
        console.error('Sin respuesta del servidor:', error.request)
        errorMessage.value = 'No se pudo conectar con el servidor. Por favor verifica tu conexión a internet.'
      } else if (error.friendlyMessage) {
        // Error personalizado del interceptor de axios
        addDebugInfo(`Error amigable: ${error.friendlyMessage}`)
        errorMessage.value = error.friendlyMessage
      } else {
        // Error durante la configuración de la solicitud
        addDebugInfo(`Error de configuración: ${error.message}`)
        console.error('Error durante la configuración de la solicitud:', error.message)
        errorMessage.value = 'Error en el registro. Por favor intenta de nuevo.'
      }
    }
  } finally {
    // La propiedad loading ya se maneja en el store
  }
}
</script>

<style lang="scss" scoped>
.register-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
  padding: 2rem;
}

.register-card {
  background: white;
  padding: 2rem;
  border-radius: 8px;
  box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
  width: 100%;
  max-width: 400px;

  .logo-container {
    display: flex;
    justify-content: center;
    margin-bottom: 1.5rem;
    
    .login-logo {
      height: 100px;
      width: auto;
    }
  }

  h2 {
    margin: 0 0 2rem;
    text-align: center;
    color: #333;
    font-weight: 600;
  }

  .error-message {
    background-color: #ffebee;
    color: #d32f2f;
    padding: 0.75rem;
    border-radius: 4px;
    margin-bottom: 1rem;
    text-align: center;
  }

  .form-group {
    margin-bottom: 1.5rem;

    label {
      display: block;
      margin-bottom: 0.5rem;
      color: #666;
    }

    input {
      width: 100%;
      padding: 0.5rem;
      border: 1px solid #ddd;
      border-radius: 4px;
      
      &.input-error {
        border-color: #d32f2f;
        background-color: #fff8f8;
      }
    }
    
    .error-text {
      display: block;
      color: #d32f2f;
      margin-top: 0.25rem;
      font-size: 0.8rem;
    }

    small {
      display: block;
      color: #666;
      margin-top: 0.25rem;
      font-size: 0.8rem;
    }
  }

  .btn {
    width: 100%;
    margin-bottom: 1rem;
    
    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }

  .login-link {
    text-align: center;
    margin: 0;
    color: #666;

    a {
      color: #1976d2;
      text-decoration: none;

      &:hover {
        text-decoration: underline;
      }
    }
  }
  
  .debug-info {
    margin-top: 2rem;
    padding: 1rem;
    border: 1px dashed #ccc;
    background-color: #f8f8f8;
    border-radius: 4px;
    
    h3 {
      margin-top: 0;
      font-size: 1rem;
      color: #666;
    }
    
    ul {
      margin: 0;
      padding-left: 1.5rem;
      
      li {
        font-size: 0.8rem;
        color: #666;
        margin-bottom: 0.25rem;
        word-break: break-all;
      }
    }
  }
}
</style> 
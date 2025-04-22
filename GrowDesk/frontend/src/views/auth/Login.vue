<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-container">
        <img src="@/assets/logo.png" alt="GrowDesk Logo" class="login-logo">
      </div>
      <h2>Iniciar Sesión</h2>
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      <form @submit.prevent="handleLogin">
        <div class="form-group">
          <label for="email">Correo Electrónico</label>
          <input 
            type="email" 
            id="email" 
            v-model="email" 
            :class="{'input-error': validationErrors.email}"
            required
          >
          <small v-if="validationErrors.email" class="error-text">{{ validationErrors.email }}</small>
        </div>
        <div class="form-group">
          <label for="password">Contraseña</label>
          <input 
            type="password" 
            id="password" 
            v-model="password" 
            :class="{'input-error': validationErrors.password}"
            required
          >
          <small v-if="validationErrors.password" class="error-text">{{ validationErrors.password }}</small>
        </div>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          {{ loading ? 'Iniciando sesión...' : 'Iniciar Sesión' }}
        </button>
      </form>
      <p class="register-link">
        ¿No tienes una cuenta? <router-link to="/register">Regístrate</router-link>
      </p>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useAuthStore } from '@/stores/auth'

const router = useRouter()
const authStore = useAuthStore()
const email = ref('')
const password = ref('')
const errorMessage = ref('')
const loading = computed(() => authStore.loading)

const validationErrors = reactive({
  email: '',
  password: ''
})

// Función para validar el formulario
const validateForm = () => {
  let isValid = true
  
  // Resetear errores
  validationErrors.email = ''
  validationErrors.password = ''
  
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
  }
  
  return isValid
}

const handleLogin = async () => {
  if (!validateForm()) return

  errorMessage.value = ''

  try {
    console.log('Iniciando sesión con:', { email: email.value, password: '******' })
    console.log('URL de la API:', import.meta.env.VITE_API_URL)
    
    const success = await authStore.login(email.value, password.value)
    console.log('Respuesta de login:', success)
    
    if (success) {
      router.push({ name: 'dashboard' })
    } else {
      errorMessage.value = 'Credenciales inválidas'
    }
  } catch (e: any) {
    console.error('Error durante el inicio de sesión:', e)
    
    if (e.response) {
      console.error('Respuesta del servidor:', e.response.status, e.response.data)
      errorMessage.value = e.friendlyMessage || 'No se pudo iniciar sesión'
    } else if (e.request) {
      console.error('No hubo respuesta del servidor:', e.request)
      errorMessage.value = 'No se pudo conectar con el servidor. Por favor verifica tu conexión a internet.'
    } else {
      errorMessage.value = e.message || 'Ocurrió un error durante el inicio de sesión'
    }
  }
}
</script>

<style lang="scss" scoped>
.login-container {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 200px);
  padding: 2rem;
}

.login-card {
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
  }

  .btn {
    width: 100%;
    margin-bottom: 1rem;
    
    &:disabled {
      opacity: 0.7;
      cursor: not-allowed;
    }
  }

  .register-link {
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
}
</style> 
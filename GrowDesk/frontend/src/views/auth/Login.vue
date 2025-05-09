<template>
  <div class="login-container">
    <div class="login-card">
      <div class="logo-container">
        <img src="@/assets/logo.png" alt="GrowDesk Logo" class="login-logo">
        <div class="logo-circle"></div>
      </div>
      <h2>Iniciar Sesión</h2>
      <div v-if="errorMessage" class="error-message">
        {{ errorMessage }}
      </div>
      <form @submit.prevent="handleLogin" class="login-form">
        <div class="form-group">
          <label for="email">Correo Electrónico</label>
          <div class="input-wrapper">
            <i class="pi pi-envelope"></i>
            <input 
              type="email" 
              id="email" 
              v-model="email" 
              :class="{'input-error': validationErrors.email}"
              required
              placeholder="Ingresa tu correo electrónico"
            >
          </div>
          <small v-if="validationErrors.email" class="error-text">{{ validationErrors.email }}</small>
        </div>
        <div class="form-group">
          <label for="password">Contraseña</label>
          <div class="input-wrapper">
            <i class="pi pi-lock"></i>
            <input 
              type="password" 
              id="password" 
              v-model="password" 
              :class="{'input-error': validationErrors.password}"
              required
              placeholder="Ingresa tu contraseña"
            >
          </div>
          <small v-if="validationErrors.password" class="error-text">{{ validationErrors.password }}</small>
        </div>
        <button type="submit" class="btn btn-primary" :disabled="loading">
          <span class="btn-text">{{ loading ? 'Iniciando sesión...' : 'Iniciar Sesión' }}</span>
          <i class="pi pi-sign-in"></i>
        </button>
      </form>
      
      <p class="register-link">
        ¿No tienes una cuenta? <router-link to="/register">Regístrate</router-link>
      </p>
    </div>
    
    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
      <div class="circle circle-4"></div>
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
  min-height: 100vh;
  padding: 2rem;
  position: relative;
  overflow: hidden;
  background-color: var(--bg-primary);
}

.background-decoration {
  position: absolute;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  z-index: 0;
  overflow: hidden;
}

.circle {
  position: absolute;
  border-radius: 50%;
  opacity: 0.15;
  
  &.circle-1 {
    width: 400px;
    height: 400px;
    background: linear-gradient(45deg, var(--primary-color), #3b82f6);
    top: -100px;
    left: -100px;
    animation: float 15s infinite ease-in-out;
  }
  
  &.circle-2 {
    width: 300px;
    height: 300px;
    background: linear-gradient(135deg, var(--primary-color), #4f46e5);
    bottom: -50px;
    right: -50px;
    animation: float 20s infinite ease-in-out reverse;
  }
  
  &.circle-3 {
    width: 200px;
    height: 200px;
    background: linear-gradient(225deg, var(--primary-color), #6366f1);
    top: 60%;
    left: 15%;
    animation: float 18s infinite ease-in-out 2s;
  }
  
  &.circle-4 {
    width: 250px;
    height: 250px;
    background: linear-gradient(315deg, var(--primary-color), #818cf8);
    top: 10%;
    right: 15%;
    animation: float 22s infinite ease-in-out 1s;
  }
}

@keyframes float {
  0% {
    transform: translate(0, 0) rotate(0deg);
  }
  25% {
    transform: translate(10px, 10px) rotate(2deg);
  }
  50% {
    transform: translate(0, 20px) rotate(0deg);
  }
  75% {
    transform: translate(-10px, 10px) rotate(-2deg);
  }
  100% {
    transform: translate(0, 0) rotate(0deg);
  }
}

.login-card {
  background: white;
  padding: 2.5rem;
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  width: 100%;
  max-width: 450px;
  position: relative;
  z-index: 10;
  animation: fadeIn 0.8s ease-out;
  
  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }

  .logo-container {
    position: relative;
    display: flex;
    justify-content: center;
    margin-bottom: 2rem;
    
    .login-logo {
      height: 100px;
      width: auto;
      position: relative;
      z-index: 2;
    }
    
    .logo-circle {
      position: absolute;
      width: 80px;
      height: 80px;
      background-color: var(--light-bg);
      border-radius: 50%;
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      z-index: 1;
    }
  }

  h2 {
    text-align: center;
    color: var(--text-primary);
    margin-bottom: 1.5rem;
    font-weight: 600;
  }

  .login-form {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
    margin-bottom: 1.5rem;
  }

  .form-group {
    display: flex;
    flex-direction: column;
    
    label {
      font-size: 0.9rem;
      font-weight: 500;
      margin-bottom: 0.5rem;
      color: var(--text-secondary);
    }
    
    .input-wrapper {
      position: relative;
      
      i {
        position: absolute;
        left: 1rem;
        top: 50%;
        transform: translateY(-50%);
        color: var(--text-tertiary);
      }
      
      input {
        width: 100%;
        padding: 1rem 1rem 1rem 2.5rem;
        border: 1px solid var(--border-color);
        border-radius: var(--border-radius-sm);
        font-size: 1rem;
        transition: all 0.3s;
        
        &:focus {
          border-color: var(--primary-color);
          box-shadow: 0 0 0 2px rgba(var(--primary-color-rgb), 0.1);
          outline: none;
        }
        
        &.input-error {
          border-color: var(--error-color);
        }
      }
    }
    
    .error-text {
      color: var(--error-color);
      font-size: 0.8rem;
      margin-top: 0.5rem;
    }
  }

  .btn {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.75rem 1.5rem;
    font-size: 1rem;
    font-weight: 500;
    border-radius: var(--border-radius-sm);
    border: none;
    cursor: pointer;
    transition: all 0.3s;
    
    &-primary {
      background-color: var(--primary-color);
      color: white;
      
      &:hover:not(:disabled) {
        background-color: var(--primary-dark);
      }
      
      &:disabled {
        opacity: 0.7;
        cursor: not-allowed;
      }
    }
    
    &-secondary {
      background-color: var(--secondary-color);
      color: white;
      
      &:hover {
        background-color: var(--secondary-dark);
      }
    }
    
    i {
      font-size: 1.1rem;
    }
  }

  .error-message {
    background-color: rgba(var(--error-color-rgb), 0.1);
    color: var(--error-color);
    padding: 0.75rem;
    border-radius: var(--border-radius-sm);
    margin-bottom: 1.5rem;
    font-size: 0.9rem;
    border-left: 3px solid var(--error-color);
  }

  .register-link {
    text-align: center;
    margin-top: 1.5rem;
    font-size: 0.9rem;
    color: var(--text-tertiary);
    
    a {
      color: var(--primary-color);
      text-decoration: none;
      font-weight: 500;
      
      &:hover {
        text-decoration: underline;
      }
    }
  }
}

@media (max-width: 768px) {
  .login-card {
    padding: 2rem;
    
    .login-logo {
      height: 80px;
    }
  }
}

@media (max-width: 480px) {
  .login-container {
    padding: 1rem;
  }
  
  .login-card {
    padding: 1.5rem;
    
    .login-logo {
      height: 70px;
    }
    
    .form-group .input-wrapper input {
      padding: 0.9rem 1rem 0.9rem 2.5rem;
    }
  }
}
</style> 
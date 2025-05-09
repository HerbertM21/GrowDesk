<template>
  <div class="not-found-container">
    <div class="not-found-content">
      <div class="logo-container">
        <img src="@/assets/logo.png" alt="GrowDesk Logo" class="not-found-logo">
        <div class="logo-circle"></div>
      </div>

      <div class="error-content">
        <div class="error-code">404</div>
        <h1>Página no encontrada</h1>
        <p>Lo sentimos, la página que estás buscando no existe o no está disponible.</p>
      </div>

      <div class="actions">
        <router-link to="/" class="btn btn-primary">
          <i class="pi pi-home"></i>
          Volver al inicio
        </router-link>
        <router-link to="/tickets" class="btn btn-outline-primary">
          <i class="pi pi-ticket"></i>
          Ver tickets
        </router-link>
      </div>
    </div>

    <div class="background-decoration">
      <div class="circle circle-1"></div>
      <div class="circle circle-2"></div>
      <div class="circle circle-3"></div>
      <div class="circle circle-4"></div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent } from 'vue'

export default defineComponent({
  name: 'NotFound'
})
</script>

<style lang="scss" scoped>
.not-found-container {
  min-height: 100vh;
  display: flex;
  align-items: center;
  justify-content: center;
  background-color: var(--bg-primary);
  position: relative;
  overflow: hidden;
  padding: 1rem;
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

.not-found-content {
  background: white;
  padding: 3.5rem;
  border-radius: var(--border-radius);
  box-shadow: var(--card-shadow);
  text-align: center;
  max-width: 550px;
  width: 100%;
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
}

.logo-container {
  position: relative;
  margin-bottom: 2.5rem;
  display: flex;
  justify-content: center;
  
  .not-found-logo {
    height: 90px;
    width: auto;
    position: relative;
    z-index: 2;
    animation: pulse 3s infinite ease-in-out;
    
    @keyframes pulse {
      0%, 100% {
        transform: scale(1);
      }
      50% {
        transform: scale(1.05);
      }
    }
  }
  
  .logo-circle {
    position: absolute;
    width: 120px;
    height: 120px;
    border-radius: 50%;
    background: linear-gradient(135deg, rgba(var(--primary-rgb), 0.1), rgba(99, 102, 241, 0.1));
    top: 50%;
    left: 50%;
    transform: translate(-50%, -50%);
    z-index: 1;
    animation: pulse-circle 3s infinite ease-in-out;
    
    @keyframes pulse-circle {
      0%, 100% {
        transform: translate(-50%, -50%) scale(1);
        opacity: 0.3;
      }
      50% {
        transform: translate(-50%, -50%) scale(1.2);
        opacity: 0.5;
      }
    }
    
    &::after {
      content: '';
      position: absolute;
      width: 150px;
      height: 150px;
      border-radius: 50%;
      border: 2px dashed rgba(var(--primary-rgb), 0.2);
      top: 50%;
      left: 50%;
      transform: translate(-50%, -50%);
      animation: rotate 30s infinite linear;
      
      @keyframes rotate {
        from {
          transform: translate(-50%, -50%) rotate(0deg);
        }
        to {
          transform: translate(-50%, -50%) rotate(360deg);
        }
      }
    }
  }
}

.error-content {
  margin-bottom: 2.5rem;
  animation: slideUp 0.8s ease-out 0.2s both;
  
  @keyframes slideUp {
    from {
      opacity: 0;
      transform: translateY(20px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
}

.error-code {
  font-size: 7rem;
  font-weight: 800;
  background: linear-gradient(135deg, var(--primary-color), #6366f1);
  -webkit-background-clip: text;
  -webkit-text-fill-color: transparent;
  background-clip: text;
  color: transparent;
  line-height: 1;
  margin-bottom: 1rem;
  animation: shimmer 3s infinite;
  
  @keyframes shimmer {
    0% {
      background-position: -200px;
    }
    100% {
      background-position: 200px;
    }
  }
}

h1 {
  font-size: 2.25rem;
  color: var(--text-primary);
  margin-bottom: 1rem;
  font-weight: 700;
}

p {
  color: var(--text-secondary);
  margin-bottom: 1rem;
  font-size: 1.1rem;
  line-height: 1.6;
  max-width: 400px;
  margin-left: auto;
  margin-right: auto;
}

.actions {
  display: flex;
  gap: 1rem;
  justify-content: center;
  flex-wrap: wrap;
  animation: fadeIn 0.8s ease-out 0.4s both;

  .btn {
    display: inline-flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 0.85rem 1.8rem;
    border-radius: var(--border-radius);
    font-size: 1rem;
    font-weight: 600;
    text-decoration: none;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    cursor: pointer;
    position: relative;
    overflow: hidden;
    
    &::before {
      content: '';
      position: absolute;
      top: 50%;
      left: 50%;
      width: 0;
      height: 0;
      background: rgba(255, 255, 255, 0.2);
      border-radius: 50%;
      transform: translate(-50%, -50%);
      transition: width 0.6s, height 0.6s;
    }
    
    &:hover::before {
      width: 300px;
      height: 300px;
    }

    i {
      font-size: 1.2rem;
      transition: transform 0.3s ease;
    }

    &:hover i {
      transform: translateX(-3px);
    }

    &.btn-primary {
      background-color: var(--primary-color);
      color: white;
      border: none;
      box-shadow: 0 4px 6px rgba(var(--primary-rgb), 0.25);

      &:hover {
        background-color: var(--primary-hover);
        transform: translateY(-3px);
        box-shadow: 0 6px 12px rgba(var(--primary-rgb), 0.3);
      }
      
      &:active {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(var(--primary-rgb), 0.3);
      }
    }

    &.btn-outline-primary {
      background-color: transparent;
      color: var(--primary-color);
      border: 2px solid var(--primary-color);

      &:hover {
        background-color: rgba(var(--primary-rgb), 0.05);
        transform: translateY(-3px);
        box-shadow: 0 6px 12px rgba(var(--primary-rgb), 0.15);
      }
      
      &:active {
        transform: translateY(-1px);
        box-shadow: 0 2px 4px rgba(var(--primary-rgb), 0.15);
      }
    }
  }
}

// Responsive adjustments
@media (max-width: 650px) {
  .not-found-content {
    padding: 2.5rem 1.5rem;
  }
  
  .error-code {
    font-size: 5rem;
  }

  h1 {
    font-size: 1.75rem;
  }

  p {
    font-size: 1rem;
  }
}

@media (max-width: 450px) {
  .not-found-content {
    padding: 2rem 1.25rem;
  }
  
  .logo-container {
    margin-bottom: 2rem;
    
    .not-found-logo {
      height: 70px;
    }
    
    .logo-circle {
      width: 100px;
      height: 100px;
    }
  }
  
  .error-code {
    font-size: 4.5rem;
  }

  h1 {
    font-size: 1.5rem;
  }

  .actions {
    flex-direction: column;
    
    .btn {
      width: 100%;
      padding: 0.75rem 1.5rem;
    }
  }
}
</style> 
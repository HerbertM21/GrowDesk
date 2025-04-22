<template>
  <div class="confirm-dialog-backdrop" @click.self="$emit('cancel')">
    <div class="confirm-dialog">
      <div class="confirm-dialog-header">
        <h3>{{ title }}</h3>
        <button class="close-btn" @click="$emit('cancel')">&times;</button>
      </div>
      
      <div class="confirm-dialog-body">
        <p>{{ message }}</p>
      </div>
      
      <div class="confirm-dialog-footer">
        <button class="btn btn-secondary" @click="$emit('cancel')">
          {{ cancelText }}
        </button>
        <button class="btn btn-danger" @click="$emit('confirm')">
          {{ confirmText }}
        </button>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Definir props con valores por defecto
const props = defineProps({
  title: {
    type: String,
    default: 'Confirmar acción'
  },
  message: {
    type: String,
    default: '¿Estás seguro de que quieres realizar esta acción?'
  },
  confirmText: {
    type: String,
    default: 'Confirmar'
  },
  cancelText: {
    type: String,
    default: 'Cancelar'
  }
});

// Emits
defineEmits<{
  (e: 'confirm'): void;
  (e: 'cancel'): void;
}>();
</script>

<style lang="scss" scoped>
.confirm-dialog-backdrop {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1100; // Mayor que otros modales
}

.confirm-dialog {
  background: white;
  border-radius: 8px;
  width: 100%;
  max-width: 400px;
  box-shadow: 0 5px 15px rgba(0, 0, 0, 0.3);
}

.confirm-dialog-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding: 1rem 1.5rem;
  border-bottom: 1px solid #eee;
  
  h3 {
    margin: 0;
    font-size: 1.25rem;
    color: #333;
  }
  
  .close-btn {
    background: none;
    border: none;
    font-size: 1.5rem;
    cursor: pointer;
    color: #666;
    
    &:hover {
      color: #333;
    }
  }
}

.confirm-dialog-body {
  padding: 1.5rem;
  
  p {
    margin: 0;
    color: #555;
  }
}

.confirm-dialog-footer {
  padding: 1rem 1.5rem;
  border-top: 1px solid #eee;
  display: flex;
  justify-content: flex-end;
  gap: 1rem;
  
  .btn {
    padding: 0.5rem 1rem;
    border-radius: 4px;
    font-weight: 500;
    cursor: pointer;
    border: none;
    
    &.btn-secondary {
      background-color: #f5f5f5;
      color: #333;
      
      &:hover {
        background-color: #e0e0e0;
      }
    }
    
    &.btn-danger {
      background-color: #f44336;
      color: white;
      
      &:hover {
        background-color: #d32f2f;
      }
    }
  }
}
</style> 
<template>
  <div class="create-ticket">
    <h1>Crear Nuevo Ticket</h1>
    <form @submit.prevent="handleSubmit" class="ticket-form">
      <div class="form-group">
        <label for="title">Título</label>
        <input type="text" id="title" v-model="form.title" required>
      </div>

      <div class="form-group">
        <label for="description">Descripción</label>
        <textarea id="description" v-model="form.description" required></textarea>
      </div>

      <div class="form-group">
        <label for="priority">Prioridad</label>
        <select id="priority" v-model="form.priority" required>
          <option value="low">Baja</option>
          <option value="medium">Media</option>
          <option value="high">Alta</option>
          <option value="urgent">Urgente</option>
        </select>
      </div>

      <div class="form-group">
        <label for="category">Categoría</label>
        <select id="category" v-model="form.category" required>
          <option value="technical">Técnico</option>
          <option value="billing">Facturación</option>
          <option value="general">General</option>
          <option value="feature">Solicitud de Función</option>
        </select>
      </div>

      <div class="form-actions">
        <button type="submit" class="btn btn-primary">Crear Ticket</button>
        <router-link to="/tickets" class="btn btn-secondary">Cancelar</router-link>
      </div>
    </form>
  </div>
</template>

<script setup lang="ts">
import { ref } from 'vue'
import { useRouter } from 'vue-router'
import { useTicketStore } from '@/stores/tickets'

const router = useRouter()
const ticketStore = useTicketStore()

const form = ref({
  title: '',
  description: '',
  priority: 'medium',
  category: 'general'
})

const handleSubmit = async () => {
  try {
    await ticketStore.createTicket(form.value)
    router.push('/tickets')
  } catch (error) {
    console.error('Error al crear el ticket:', error)
  }
}
</script>

<style lang="scss" scoped>
.create-ticket {
  max-width: 800px;
  margin: 0 auto;
  padding: 2rem;

  h1 {
    margin-bottom: 2rem;
  }

  .ticket-form {
    background: white;
    padding: 2rem;
    border-radius: 8px;
    box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);

    .form-group {
      margin-bottom: 1.5rem;

      label {
        display: block;
        margin-bottom: 0.5rem;
        font-weight: 500;
      }

      input, select, textarea {
        width: 100%;
        padding: 0.5rem;
        border: 1px solid #ddd;
        border-radius: 4px;
      }

      textarea {
        min-height: 150px;
        resize: vertical;
      }
    }

    .form-actions {
      display: flex;
      gap: 1rem;
      margin-top: 2rem;
    }
  }
}
</style> 
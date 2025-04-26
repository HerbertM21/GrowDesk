<template>
  <div class="faq-section">
    <h4 class="mb-3">Preguntas frecuentes</h4>
    
    <div v-if="isLoading" class="text-center py-3">
      <div class="spinner-border spinner-border-sm" role="status">
        <span class="visually-hidden">Cargando...</span>
      </div>
      <p class="mb-0 mt-2">Cargando preguntas frecuentes...</p>
    </div>
    
    <div v-else-if="error" class="alert alert-danger py-2" role="alert">
      <i class="bi bi-exclamation-triangle-fill me-2"></i>
      Error al cargar preguntas frecuentes
    </div>
    
    <div v-else-if="faqs.length === 0" class="text-center py-3">
      <p class="text-muted">No hay preguntas frecuentes disponibles</p>
    </div>
    
    <div v-else class="accordion accordion-flush mb-4" id="faqAccordion">
      <div v-for="(faq, index) in faqs" :key="faq.id" class="accordion-item">
        <h2 class="accordion-header" :id="'heading-' + faq.id">
          <button 
            class="accordion-button collapsed" 
            type="button" 
            data-bs-toggle="collapse" 
            :data-bs-target="'#collapse-' + faq.id" 
            aria-expanded="false" 
            :aria-controls="'collapse-' + faq.id"
          >
            {{ faq.question }}
          </button>
        </h2>
        <div 
          :id="'collapse-' + faq.id" 
          class="accordion-collapse collapse" 
          :aria-labelledby="'heading-' + faq.id" 
          data-bs-parent="#faqAccordion"
        >
          <div class="accordion-body">
            {{ faq.answer }}
          </div>
        </div>
      </div>
    </div>
    
    <div class="d-grid gap-2">
      <button class="btn btn-primary" @click="startChat">
        <i class="bi bi-chat-dots-fill me-2"></i> Iniciar Chat
      </button>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { getFaqs, type FAQ } from '@/services/faqs';

const faqs = ref<FAQ[]>([]);
const isLoading = ref(true);
const error = ref<string | null>(null);

onMounted(async () => {
  try {
    console.log('Cargando FAQs para el widget...');
    const response = await getFaqs();
    // Solo mostrar las FAQs publicadas
    faqs.value = response.filter(faq => faq.isPublished);
    console.log(`FAQs cargadas: ${faqs.value.length} preguntas publicadas`);
  } catch (err) {
    console.error('Error al cargar FAQs:', err);
    error.value = 'No se pudieron cargar las preguntas frecuentes';
  } finally {
    isLoading.value = false;
  }
});

const emit = defineEmits(['start-chat']);

const startChat = () => {
  emit('start-chat');
};
</script>

<style scoped>
.faq-section {
  padding: 1rem;
}

.accordion-button {
  font-size: 0.9rem;
  font-weight: 500;
  padding: 0.75rem 1rem;
}

.accordion-body {
  font-size: 0.85rem;
  padding: 1rem;
  background-color: #f8f9fa;
  white-space: pre-line;
}
</style>
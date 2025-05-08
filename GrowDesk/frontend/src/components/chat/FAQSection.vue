/* eslint-disable */
<template>
  <div class="faq-section">
    <div v-if="loading" class="loading-state">
      <i class="pi pi-spin pi-spinner"></i>
      <span>Cargando preguntas frecuentes...</span>
    </div>
    
    <div v-else-if="error" class="error-state">
      <i class="pi pi-exclamation-triangle"></i>
      <span>{{ error }}</span>
    </div>
    
    <div v-else-if="publishedFaqs.length === 0" class="empty-state">
      <i class="pi pi-info-circle"></i>
      <span>No hay preguntas frecuentes disponibles.</span>
    </div>
    
    <div v-else class="faq-list">
      <div v-for="faq in publishedFaqs" :key="faq.id" class="faq-item">
        <div class="faq-question" @click="toggleFaq(faq.id || 0)">
          <span>{{ faq.question }}</span>
          <i class="pi" :class="expandedFaqs.includes(faq.id || 0) ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
        </div>
        <div v-if="expandedFaqs.includes(faq.id || 0)" class="faq-answer">
          <p>{{ faq.answer }}</p>
          <button class="start-chat-btn" @click="emit('start-chat', faq)">
            <i class="pi pi-comments"></i>
            Iniciar chat sobre este tema
          </button>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue';
import { useFaqStore } from '@/stores/faqs';
import { storeToRefs } from 'pinia';
import type { FAQ } from '@/types/faq.types';

// Define props and emits
const emit = defineEmits<{
  (e: 'start-chat', faq: FAQ): void;
}>();

// Initialize store and extract reactive refs
const faqStore = useFaqStore();
const { faqs, loading, error } = storeToRefs(faqStore);

// Local state
const expandedFaqs = ref<number[]>([]);

// Computed property for published FAQs
const publishedFaqs = computed<FAQ[]>(() => {
  return faqs.value.filter((faq: FAQ) => faq.isPublished);
});

// Methods
const toggleFaq = (id: number) => {
  const index = expandedFaqs.value.indexOf(id);
  if (index === -1) {
    expandedFaqs.value.push(id);
  } else {
    expandedFaqs.value.splice(index, 1);
  }
};

// Lifecycle hooks
onMounted(async () => {
  try {
    await faqStore.fetchFaqs();
  } catch (err) {
    console.error('Error fetching FAQs:', err);
  }
});
</script>

<style scoped lang="scss">
.faq-section {
  padding: 1rem;
  
  .loading-state,
  .error-state,
  .empty-state {
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    padding: 1rem;
    text-align: center;
    color: var(--text-color-secondary);
    
    i {
      font-size: 1.2rem;
    }
  }
  
  .error-state {
    color: var(--red-500);
  }
  
  .faq-list {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    
    .faq-item {
      border: 1px solid var(--surface-border);
      border-radius: var(--border-radius);
      overflow: hidden;
      
      .faq-question {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 1rem;
        background-color: var(--surface-ground);
        cursor: pointer;
        transition: background-color 0.2s;
        
        &:hover {
          background-color: var(--surface-hover);
        }
        
        span {
          flex: 1;
          margin-right: 1rem;
        }
        
        i {
          font-size: 1rem;
          color: var(--text-color-secondary);
        }
      }
      
      .faq-answer {
        padding: 1rem;
        background-color: var(--surface-card);
        border-top: 1px solid var(--surface-border);
        
        p {
          margin: 0 0 1rem 0;
          line-height: 1.5;
        }
        
        .start-chat-btn {
          display: inline-flex;
          align-items: center;
          gap: 0.5rem;
          padding: 0.5rem 1rem;
          border: none;
          border-radius: var(--border-radius);
          background-color: var(--primary-color);
          color: var(--primary-color-text);
          cursor: pointer;
          transition: background-color 0.2s;
          
          &:hover {
            background-color: var(--primary-600);
          }
          
          i {
            font-size: 1rem;
          }
        }
      }
    }
  }
}
</style>
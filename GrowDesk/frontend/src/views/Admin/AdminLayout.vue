<template>
  <div class="admin-layout">
    <div class="admin-content">
      <div class="action-header" v-if="$slots.actions">
        <slot name="actions"></slot>
      </div>
      
      <div class="admin-main-content">
        <slot></slot>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
// Sin estados adicionales, solo un contenedor para el contenido
</script>

<style lang="scss" scoped>
.admin-layout {
  width: 100%;
  
  .admin-content {
    width: 100%;
    
    .action-header {
      display: flex;
      justify-content: flex-end;
      margin-bottom: 1.5rem;
      padding: 0.5rem 0;
    }
    
    .admin-main-content {
      width: 100%;
    }
  }
}
</style> 
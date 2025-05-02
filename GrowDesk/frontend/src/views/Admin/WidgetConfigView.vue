<template>
  <div class="widget-config">
    <!-- Permiso verificado: Solo admins -->
    <div v-if="!authStore.isAdmin" class="permission-error">
      <div class="error-container">
        <i class="pi pi-exclamation-triangle"></i>
        <h2>Acceso restringido</h2>
        <p>No tienes permiso para acceder a esta sección. Esta página solo está disponible para administradores.</p>
        <button class="btn btn-primary" @click="goToDashboard">Volver al Dashboard</button>
      </div>
    </div>

    <div v-else>
    <h1 class="page-title">Configuración del Widget</h1>
    
    <div v-if="loading" class="loading-indicator">
      <div class="spinner"></div>
      <p>Cargando configuración...</p>
    </div>
    
    <div v-else class="config-container">
      <div class="config-card active-config" v-if="currentConfig">
        <h2>Widget Activo</h2>
        
        <div class="config-details">
          <div class="config-field">
            <label>Nombre:</label>
            <p>{{ currentConfig.name }}</p>
          </div>
          
          <div class="config-field">
            <label>API Key:</label>
            <div class="api-key-container">
              <p class="api-key">{{ currentConfig.apiKey }}</p>
              <button @click="regenerateApiKey" class="btn btn-secondary">
                <i class="fas fa-sync-alt"></i> Regenerar
              </button>
            </div>
          </div>
          
          <div class="config-field">
            <label>Marca:</label>
            <p>{{ currentConfig.brandName }}</p>
          </div>
          
          <div class="config-field">
            <label>Mensaje de bienvenida:</label>
            <p>{{ currentConfig.welcomeMessage }}</p>
          </div>
          
          <div class="config-field">
            <label>Color principal:</label>
            <div class="color-preview" :style="{ backgroundColor: currentConfig.primaryColor }"></div>
            <p>{{ currentConfig.primaryColor }}</p>
          </div>
          
          <div class="config-field">
            <label>Posición:</label>
            <p>{{ currentConfig.position }}</p>
          </div>
          
          <div class="config-field">
            <label>Dominios permitidos:</label>
            <ul class="domains-list">
              <li v-for="(domain, index) in currentConfig.allowedDomains" :key="index">
                {{ domain }}
              </li>
            </ul>
          </div>
        </div>
        
        <div class="embed-code">
          <h3>Código para incrustar:</h3>
          <pre>{{ currentConfig.embedCode }}</pre>
          <button @click="copyEmbedCode" class="btn btn-primary">
            <i class="fas fa-copy"></i> Copiar código
          </button>
        </div>
        
        <div class="action-buttons">
          <button @click="editConfig" class="btn btn-primary">
            <i class="fas fa-edit"></i> Editar
          </button>
        </div>
      </div>
      
      <div v-else class="no-config">
        <p>No hay configuración de widget. Crea una nueva para comenzar.</p>
        <button @click="showCreateForm = true" class="btn btn-primary">
          <i class="fas fa-plus"></i> Crear Widget
        </button>
      </div>
      
      <!-- Información de la arquitectura y puertos -->
      <div class="architecture-info">
        <h3>Información de la arquitectura</h3>
        <p>El widget de chat se comunica con diferentes componentes del sistema GrowDesk:</p>
        
        <div class="port-info">
          <h4>Puertos utilizados:</h4>
          <ul>
            <li>
              <strong>8082</strong>: API del Widget - <em>widget-api</em>
              <p>Intermediario entre el widget y el sistema principal. Recibe los mensajes y crea tickets.</p>
            </li>
            <li>
              <strong>8000</strong>: Backend de GrowDesk 
              <p>Sistema principal donde se procesan los tickets y se maneja la lógica de negocio.</p>
            </li>
            <li>
              <strong>3000</strong>: Frontend de GrowDesk (Panel de administración)
              <p>Interfaz donde los agentes de soporte pueden ver y responder tickets.</p>
            </li>
          </ul>
        </div>
        
        <div class="flow-info">
          <h4>Flujo de comunicación:</h4>
          <ol>
            <li>El usuario envía un mensaje a través del widget en el sitio web.</li>
            <li>El widget envía el mensaje a la <strong>API del Widget</strong> (puerto 8082).</li>
            <li>La API del Widget crea un ticket o añade el mensaje a un ticket existente en el <strong>Backend de GrowDesk</strong> (puerto 8000).</li>
            <li>Los agentes ven el nuevo ticket en el <strong>Frontend de GrowDesk</strong> (puerto 3000).</li>
            <li>Cuando responden, la respuesta se envía al usuario a través de la <strong>API del Widget</strong> al widget en el sitio web.</li>
          </ol>
        </div>
      </div>
      
      <!-- Formulario de edición -->
      <div class="modal" v-if="showEditForm">
        <div class="modal-content">
          <h2>Editar Widget</h2>
          
          <form @submit.prevent="updateConfig">
            <div class="form-group">
              <label for="name">Nombre:</label>
              <input type="text" id="name" v-model="editForm.name" required>
            </div>
            
            <div class="form-group">
              <label for="brandName">Marca:</label>
              <input type="text" id="brandName" v-model="editForm.brandName" required>
            </div>
            
            <div class="form-group">
              <label for="welcomeMessage">Mensaje de bienvenida:</label>
              <input type="text" id="welcomeMessage" v-model="editForm.welcomeMessage">
            </div>
            
            <div class="form-group">
              <label for="primaryColor">Color principal:</label>
              <input type="color" id="primaryColor" v-model="editForm.primaryColor">
            </div>
            
            <div class="form-group">
              <label for="position">Posición:</label>
              <select id="position" v-model="editForm.position">
                <option value="bottom-right">Abajo a la derecha</option>
                <option value="bottom-left">Abajo a la izquierda</option>
                <option value="top-right">Arriba a la derecha</option>
                <option value="top-left">Arriba a la izquierda</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>Dominios permitidos:</label>
              <div v-for="(domain, index) in editForm.allowedDomains" :key="index" class="domain-input">
                <input type="text" v-model="editForm.allowedDomains[index]">
                <button type="button" @click="removeDomain(index)" class="btn btn-danger">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <button type="button" @click="addDomain" class="btn btn-secondary">
                <i class="fas fa-plus"></i> Añadir dominio
              </button>
            </div>
            
            <div class="form-actions">
              <button type="button" @click="showEditForm = false" class="btn btn-secondary">Cancelar</button>
              <button type="submit" class="btn btn-primary">Guardar</button>
            </div>
          </form>
        </div>
      </div>
      
      <!-- Formulario de creación -->
      <div class="modal" v-if="showCreateForm">
        <div class="modal-content">
          <h2>Crear Nuevo Widget</h2>
          
          <form @submit.prevent="createConfig">
            <div class="form-group">
              <label for="create-name">Nombre:</label>
              <input type="text" id="create-name" v-model="createForm.name" required>
            </div>
            
            <div class="form-group">
              <label for="create-brandName">Marca:</label>
              <input type="text" id="create-brandName" v-model="createForm.brandName" required>
            </div>
            
            <div class="form-group">
              <label for="create-welcomeMessage">Mensaje de bienvenida:</label>
              <input type="text" id="create-welcomeMessage" v-model="createForm.welcomeMessage" placeholder="¿En qué podemos ayudarte hoy?">
            </div>
            
            <div class="form-group">
              <label for="create-primaryColor">Color principal:</label>
              <input type="color" id="create-primaryColor" v-model="createForm.primaryColor">
            </div>
            
            <div class="form-group">
              <label for="create-position">Posición:</label>
              <select id="create-position" v-model="createForm.position">
                <option value="bottom-right">Abajo a la derecha</option>
                <option value="bottom-left">Abajo a la izquierda</option>
                <option value="top-right">Arriba a la derecha</option>
                <option value="top-left">Arriba a la izquierda</option>
              </select>
            </div>
            
            <div class="form-group">
              <label>Dominios permitidos:</label>
              <div v-for="(domain, index) in createForm.allowedDomains" :key="index" class="domain-input">
                <input type="text" v-model="createForm.allowedDomains[index]">
                <button type="button" @click="removeCreateDomain(index)" class="btn btn-danger">
                  <i class="fas fa-times"></i>
                </button>
              </div>
              <button type="button" @click="addCreateDomain" class="btn btn-secondary">
                <i class="fas fa-plus"></i> Añadir dominio
              </button>
            </div>
            
            <div class="form-actions">
              <button type="button" @click="showCreateForm = false" class="btn btn-secondary">Cancelar</button>
              <button type="submit" class="btn btn-primary">Crear</button>
            </div>
          </form>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue';
import { useRouter } from 'vue-router';
import { useAuthStore } from '@/stores/auth';
import { useNotificationsStore } from '@/stores/notifications';

// Stores para la verificación de permisos
const authStore = useAuthStore();
const router = useRouter();
const notificationsStore = useNotificationsStore();

// Función para redireccionar al dashboard
function goToDashboard() {
  router.push('/dashboard');
}

// Mostrar notificación si no tiene permiso
onMounted(() => {
  if (!authStore.isAdmin) {
    notificationsStore.addNotification({
      type: 'error',
      message: 'No tienes permiso para acceder a la configuración del widget.',
      timeout: 5000
    });
  }
  
  // Cargar la configuración actual si tiene permiso
  if (authStore.isAdmin) {
    loadConfig();
  }
});

// ... existing code ...

// Función para cargar la configuración actual
function loadConfig() {
  // Código de carga existente
}

// ... existing code ...
</script>

<style scoped>
.widget-config {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.page-title {
  margin-bottom: 30px;
  font-size: 24px;
  color: #333;
}

.loading-indicator {
  display: flex;
  flex-direction: column;
  align-items: center;
  margin-top: 50px;
}

.spinner {
  border: 4px solid #f3f3f3;
  border-top: 4px solid #4caf50;
  border-radius: 50%;
  width: 40px;
  height: 40px;
  animation: spin 1s linear infinite;
  margin-bottom: 20px;
}

@keyframes spin {
  0% { transform: rotate(0deg); }
  100% { transform: rotate(360deg); }
}

.config-container {
  background-color: #fff;
  border-radius: 8px;
  box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
}

.config-card {
  padding: 20px;
}

.config-card h2 {
  margin-top: 0;
  padding-bottom: 15px;
  border-bottom: 1px solid #eee;
}

.config-details {
  margin-bottom: 30px;
}

.config-field {
  margin-bottom: 15px;
}

.config-field label {
  font-weight: bold;
  display: block;
  margin-bottom: 5px;
  color: #555;
}

.api-key-container {
  display: flex;
  align-items: center;
}

.api-key {
  background-color: #f5f5f5;
  padding: 8px 12px;
  border-radius: 4px;
  font-family: monospace;
  margin-right: 15px;
  margin-top: 0;
  margin-bottom: 0;
}

.color-preview {
  width: 20px;
  height: 20px;
  border-radius: 4px;
  display: inline-block;
  margin-right: 10px;
  border: 1px solid #ddd;
}

.domains-list {
  margin: 0;
  padding-left: 20px;
}

.domains-list li {
  margin-bottom: 3px;
}

.embed-code {
  margin-top: 30px;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 6px;
}

.embed-code h3 {
  margin-top: 0;
  margin-bottom: 15px;
}

.embed-code pre {
  background-color: #fff;
  padding: 15px;
  border-radius: 4px;
  border: 1px solid #ddd;
  white-space: pre-wrap;
  word-break: break-all;
  font-family: monospace;
  margin-bottom: 15px;
}

.action-buttons {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}

.btn {
  padding: 8px 16px;
  border-radius: 4px;
  font-weight: 500;
  cursor: pointer;
  border: none;
  display: inline-flex;
  align-items: center;
  gap: 8px;
}

.btn-primary {
  background-color: #4caf50;
  color: white;
}

.btn-secondary {
  background-color: #f5f5f5;
  color: #333;
}

.btn-danger {
  background-color: #f44336;
  color: white;
}

.no-config {
  padding: 40px 20px;
  text-align: center;
}

.no-config p {
  margin-bottom: 20px;
  color: #666;
}

.modal {
  position: fixed;
  top: 0;
  left: 0;
  width: 100%;
  height: 100%;
  background-color: rgba(0, 0, 0, 0.5);
  display: flex;
  justify-content: center;
  align-items: center;
  z-index: 1000;
}

.modal-content {
  background-color: white;
  border-radius: 8px;
  padding: 30px;
  width: 90%;
  max-width: 600px;
  max-height: 90vh;
  overflow-y: auto;
}

.modal-content h2 {
  margin-top: 0;
  margin-bottom: 20px;
}

.form-group {
  margin-bottom: 20px;
}

.form-group label {
  display: block;
  margin-bottom: 8px;
  font-weight: 500;
  color: #555;
}

.form-group input,
.form-group select {
  width: 100%;
  padding: 10px;
  border: 1px solid #ddd;
  border-radius: 4px;
  font-size: 14px;
}

.form-group input[type="color"] {
  height: 40px;
}

.domain-input {
  display: flex;
  margin-bottom: 10px;
}

.domain-input input {
  flex: 1;
  margin-right: 10px;
}

.form-actions {
  display: flex;
  justify-content: flex-end;
  gap: 10px;
  margin-top: 30px;
}

.architecture-info {
  margin-top: 30px;
  padding: 20px;
  background-color: #f8f9fa;
  border-radius: 6px;
}

.architecture-info h3 {
  margin-top: 0;
  margin-bottom: 15px;
}

.port-info {
  margin-bottom: 20px;
}

.port-info h4 {
  margin-bottom: 10px;
}

.port-info ul {
  margin: 0;
  padding-left: 20px;
}

.port-info ul li {
  margin-bottom: 5px;
}

.flow-info {
  margin-top: 20px;
}

.flow-info h4 {
  margin-bottom: 10px;
}

.flow-info ol {
  margin: 0;
  padding-left: 20px;
}

.flow-info ol li {
  margin-bottom: 5px;
}

/* Estilos para el mensaje de error de permisos */
.permission-error {
  display: flex;
  justify-content: center;
  align-items: center;
  min-height: calc(100vh - 120px);
  padding: 2rem;
}

.error-container {
  max-width: 500px;
  text-align: center;
  background-color: #fff;
  border-radius: 8px;
  padding: 2rem;
  box-shadow: 0 4px 12px rgba(0, 0, 0, 0.1);
}

.error-container i {
  font-size: 3rem;
  color: #f59e0b;
  margin-bottom: 1rem;
}

.error-container h2 {
  margin-bottom: 1rem;
  color: #334155;
}

.error-container p {
  margin-bottom: 2rem;
  color: #64748b;
}
</style>
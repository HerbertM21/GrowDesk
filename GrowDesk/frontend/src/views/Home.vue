<template>
  <div class="employee-support">
    <!-- Sección del encabezado con fondo de gradiente y forma ondulada -->
    <div class="hero-section">
      <div class="hero-content">
        <h1 class="hero-title">¡Bienvenido al Portal de Soporte!</h1>
        <p class="hero-subtitle">Tu espacio para encontrar ayuda y recursos</p>
        
        <!-- Buscador con diseño mejorado -->
        <div class="search-container">
          <div class="search-box">
            <input 
              v-model="searchQuery" 
              type="text" 
              placeholder="¿En qué podemos ayudarte hoy?" 
            />
            <button @click="search" class="search-button">
              <i class="pi pi-search"></i>
            </button>
          </div>
        </div>
      </div>
      <div class="wave-shape"></div>
    </div>

    <div class="content-wrapper">
      <!-- Enlaces rápidos con iconos más grandes y animaciones -->
      <div class="quick-links-section">
        <h2 class="section-title">
          <span class="title-icon"><i class="pi pi-bolt"></i></span>
          Accesos Rápidos
        </h2>
        <div class="quick-links-grid">
          <div 
            class="quick-link-card" 
            v-for="(link, index) in quickLinks" 
            :key="index"
            :class="`card-${index + 1}`"
          >
            <div class="card-icon">
              <i :class="link.icon"></i>
            </div>
            <h3>{{ link.title }}</h3>
            <p>{{ link.description }}</p>
          </div>
        </div>
      </div>

      <!-- Categorías de soporte con nuevo diseño más visual -->
      <div class="categories-section">
        <h2 class="section-title">
          <span class="title-icon"><i class="pi pi-th-large"></i></span>
          Recursos de Ayuda
        </h2>
        <div class="categories-grid">
          <div class="category-card" v-for="(category, index) in categories" :key="index">
            <div class="category-icon">
              <i :class="category.icon"></i>
            </div>
            <h3>{{ category.title }}</h3>
            <ul class="category-list">
              <li v-for="(item, i) in category.items" :key="i">
                <a href="#">
                  <i class="pi pi-chevron-right"></i>
                  <span>{{ item }}</span>
                </a>
              </li>
            </ul>
          </div>
        </div>
      </div>

      <!-- Preguntas frecuentes con diseño moderno -->
      <div class="faq-section">
        <h2 class="section-title">
          <span class="title-icon"><i class="pi pi-question-circle"></i></span>
          Preguntas Frecuentes
        </h2>
        <div class="faq-container">
          <div 
            v-for="(faq, index) in faqs" 
            :key="index" 
            class="faq-item"
            :class="{ 'active': openFaq === index }"
          >
            <div class="faq-question" @click="toggleFaq(index)">
              <span>{{ faq.question }}</span>
              <div class="faq-icon">
                <i class="pi" :class="openFaq === index ? 'pi-minus' : 'pi-plus'"></i>
              </div>
            </div>
            <div class="faq-answer" v-show="openFaq === index">
              <p>{{ faq.answer }}</p>
            </div>
          </div>
        </div>
      </div>

      <!-- Contacto con diseño tipo tarjeta flotante -->
      <div class="contact-section">
        <h2 class="section-title">
          <span class="title-icon"><i class="pi pi-comments"></i></span>
          ¿Necesitas más ayuda?
        </h2>
        <div class="contact-card">
          <div class="contact-graphics">
            <div class="contact-image">
              <i class="pi pi-users"></i>
            </div>
            <div class="contact-decoration"></div>
          </div>
          <div class="contact-content">
            <h3>Equipo de Soporte</h3>
            <p>Estamos aquí para ayudarte con cualquier problema que tengas.</p>
            <div class="contact-info">
              <div class="contact-item">
                <i class="pi pi-clock"></i>
                <span>Lunes a Viernes, 8am-5pm</span>
              </div>
              <div class="contact-item">
                <i class="pi pi-phone"></i>
                <span>(+569) 99313589 </span>
              </div>
              <div class="contact-item">
                <i class="pi pi-envelope"></i>
                <span>soporte@growdesk.com</span>
              </div>
            </div>
            <router-link to="/tickets" class="create-ticket-btn">
              <i class="pi pi-list"></i>
              Listado de Tickets
            </router-link>
          </div>
        </div>
      </div>
    </div>

    <!-- Footer con mensaje motivacional -->
    <div class="support-footer">
      <div class="footer-wave"></div>
      <div class="footer-content">
        <p class="footer-tagline">"Avanzamos juntos hacia el éxito"</p>
        <p class="footer-info">Portal de Soporte Growdesk &copy; 2025</p>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import { defineComponent, ref } from 'vue';

// Definición de interfaces
interface QuickLink {
  title: string;
  description: string;
  icon: string;
}

interface Category {
  title: string;
  icon: string;
  items: string[];
}

interface Faq {
  question: string;
  answer: string;
}

export default defineComponent({
  name: 'EmployeeSupport',
  setup() {
    // Search functionality
    const searchQuery = ref('');
    const search = () => {
      // Implementar funcionalidad de búsqueda aquí
      console.log('Buscando:', searchQuery.value);
      // Normalmente harías una llamada a la API o filtrarías resultados aquí
    };

    // FAQ toggle functionality
    const openFaq = ref<number | null>(null);
    const toggleFaq = (index: number) => {
      openFaq.value = openFaq.value === index ? null : index;
    };

    // Quick links data
    const quickLinks = ref<QuickLink[]>([
      {
        title: 'Crear Ticket',
        description: 'Reportar un problema o solicitar ayuda',
        icon: 'pi pi-ticket'
      },
      {
        title: 'Guías de Equipos',
        description: 'Cómo usar herramientas de corte de vidrio',
        icon: 'pi pi-book'
      },
      {
        title: 'Procedimientos de Seguridad',
        description: 'Protocolos de seguridad en el trabajo',
        icon: 'pi pi-shield'
      },
      {
        title: 'Directorio de Empleados',
        description: 'Encuentra información de contacto',
        icon: 'pi pi-users'
      }
    ]);

    // Support categories data
    const categories = ref<Category[]>([
      {
        title: 'Software y Aplicaciones',
        icon: 'pi pi-desktop',
        items: [
          'Sistema de Punto de Venta',
          'Gestión de Inventario',
          'Base de Datos de Clientes',
          'Correo y Calendario',
          'Software de Diseño'
        ]
      },
      {
        title: 'Hardware y Equipos',
        icon: 'pi pi-cog',
        items: [
          'Máquinas de Corte de Vidrio',
          'Impresoras y Escáneres',
          'Estaciones de Trabajo',
          'Dispositivos Móviles',
          'Equipo de Seguridad'
        ]
      },
      {
        title: 'Políticas y Procedimientos',
        icon: 'pi pi-file',
        items: [
          'Manual del Empleado',
          'Directrices de Seguridad',
          'Protocolo de Servicio al Cliente',
          'Política de Devoluciones',
          'Estándares de Control de Calidad'
        ]
      },
      {
        title: 'Recursos de Capacitación',
        icon: 'pi pi-briefcase',
        items: [
          'Incorporación de Nuevos Empleados',
          'Técnicas de Corte de Vidrio',
          'Capacitación en Servicio al Cliente',
          'Certificación de Seguridad',
          'Tutoriales de Software'
        ]
      }
    ]);

    // FAQ data
    const faqs = ref<Faq[]>([
      {
        question: '¿Cómo restablezco mi contraseña?',
        answer: 'Puedes restablecer tu contraseña haciendo clic en el enlace "Olvidé mi Contraseña" en la página de inicio de sesión, o contactando al soporte de TI en la extensión 4321.'
      },
      {
        question: '¿Dónde puedo encontrar equipo de seguridad?',
        answer: 'El equipo de seguridad se almacena en los gabinetes designados en cada área de trabajo. Si los suministros están agotándose, notifica a tu supervisor o envía una solicitud a través del portal.'
      },
      {
        question: '¿Cómo reporto un problema de software?',
        answer: 'Haz clic en el botón "Crear Ticket de Soporte" en la parte inferior de esta página, selecciona "Problema de Software" del menú desplegable y proporciona detalles sobre el problema que estás experimentando.'
      },
      {
        question: '¿Cuándo son los períodos de mantenimiento del sistema?',
        answer: 'El mantenimiento regular del sistema está programado para todos los domingos de 10pm a 2am. Durante este tiempo, algunos sistemas pueden no estar disponibles o funcionar con funcionalidad limitada.'
      },
      {
        question: '¿Cómo accedo a los materiales de capacitación?',
        answer: 'Los materiales de capacitación se pueden encontrar en la sección "Recursos de Capacitación" del portal de soporte. También puedes solicitar capacitación específica enviando un ticket al departamento de RRHH.'
      }
    ]);

    return {
      searchQuery,
      search,
      openFaq,
      toggleFaq,
      quickLinks,
      categories,
      faqs
    };
  }
});
</script>

<style lang="scss" scoped>
.employee-support {
  --primary-gradient: linear-gradient(135deg, var(--primary-color) 0%, #3b82f6 100%);
  --secondary-gradient: linear-gradient(135deg, #6366f1 0%, #a855f7 100%);
  --border-radius-lg: 1.25rem;
  --transition-bounce: cubic-bezier(0.34, 1.56, 0.64, 1);
  
  background-color: var(--bg-secondary);
  position: relative;
  overflow-x: hidden;
  
  // Sección hero con fondo gradiente
  .hero-section {
    position: relative;
    padding: 4rem 2rem 6rem;
    background: var(--primary-gradient);
    color: white;
    text-align: center;
    overflow: hidden;
    
    .hero-content {
      position: relative;
      z-index: 2;
      max-width: 800px;
      margin: 0 auto;
    }
    
    .hero-title {
      font-size: 3rem;
      font-weight: 700;
      margin-bottom: 1rem;
      text-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      color: white;
    }
    
    .hero-subtitle {
      font-size: 1.5rem;
      margin-bottom: 2.5rem;
      opacity: 0.9;
    }
    
    .search-container {
      max-width: 600px;
      margin: 0 auto;
      
      .search-box {
        display: flex;
        position: relative;
        
        input {
          width: 100%;
          padding: 1.25rem 1.5rem;
          border-radius: 99px;
          border: none;
          font-size: 1.1rem;
          box-shadow: 0 5px 15px rgba(0, 0, 0, 0.2);
          transition: transform 0.3s ease;
          
          &:focus {
            outline: none;
            transform: scale(1.02);
          }
        }
        
        .search-button {
          position: absolute;
          right: 8px;
          top: 8px;
          background: var(--primary-color);
          border: none;
          color: white;
          width: 45px;
          height: 45px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          cursor: pointer;
          transition: all 0.3s ease;
          
          &:hover {
            transform: rotate(15deg) scale(1.1);
          }
          
          i {
            font-size: 1.2rem;
          }
        }
      }
    }
    
    // Forma ondulada en la parte inferior
    .wave-shape {
      position: absolute;
      bottom: -2px;
      left: 0;
      width: 100%;
      height: 4rem;
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' fill='%23f8fafc' opacity='.25'%3E%3C/path%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' fill='%23f8fafc' opacity='.5'%3E%3C/path%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%23f8fafc'%3E%3C/path%3E%3C/svg%3E");
      background-size: cover;
      background-position: center;
    }
  }
  
  .content-wrapper {
    max-width: 1300px;
    margin: 0 auto;
    padding: 3rem 1.5rem;
  }
  
  // Títulos de sección con iconos
  .section-title {
    display: flex;
    align-items: center;
    justify-content: center;
    margin-bottom: 3rem;
    color: var(--text-primary);
    font-size: 2rem;
    font-weight: 600;
    text-align: center;
    
    .title-icon {
      display: inline-flex;
      align-items: center;
      justify-content: center;
      width: 50px;
      height: 50px;
      background: var(--primary-gradient);
      border-radius: 50%;
      margin-right: 1rem;
      color: white;
      font-size: 1.4rem;
      
      i {
        font-size: 1.25rem;
      }
    }
  }
  
  // Sección de enlaces rápidos
  .quick-links-section {
    margin-bottom: 5rem;
    
    .quick-links-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
      gap: 1.5rem;
      margin: 0 auto;
    }
    
    .quick-link-card {
      background-color: var(--card-bg);
      border-radius: 24px; /* Aumentando el radio de las esquinas para que sean más redondeadas */
      padding: 2rem 1.5rem;
      text-align: center;
      box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1);
      transition: all 0.3s var(--transition-bounce);
      position: relative;
      overflow: hidden;
      z-index: 1;
      display: flex; /* Agregando para centrar verticalmente el contenido */
      flex-direction: column;
      justify-content: center; /* Centrado vertical */
      
      &::before {
        content: "";
        position: absolute;
        top: 0;
        left: 0;
        width: 100%;
        height: 5px;
        background: var(--secondary-gradient);
        transition: height 0.3s ease;
        z-index: -1;
      }
      
      &:hover {
        transform: translateY(-10px);
        box-shadow: 0 20px 25px -5px rgba(0, 0, 0, 0.1), 
                    0 10px 10px -5px rgba(0, 0, 0, 0.04);
                    
        &::before {
          height: 100%;
          opacity: 0.1;
        }
        
        .card-icon {
          transform: scale(1.2);
          
          i {
            color: var(--primary-color);
          }
        }
      }
      
      &.card-1 .card-icon {
        background-color: rgba(239, 68, 68, 0.1);
      }
      
      &.card-2 .card-icon {
        background-color: rgba(59, 130, 246, 0.1);
      }
      
      &.card-3 .card-icon {
        background-color: rgba(16, 185, 129, 0.1);
      }
      
      &.card-4 .card-icon {
        background-color: rgba(245, 158, 11, 0.1);
      }
      
      .card-icon {
        width: 80px;
        height: 80px;
        margin: 0 auto 1.5rem;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        transition: transform 0.3s var(--transition-bounce);
        
        i {
          font-size: 2.5rem;
          color: var(--text-primary);
          transition: color 0.3s ease;
        }
      }
      
      h3 {
        font-size: 1.25rem;
        color: var(--text-primary);
        font-weight: 600;
        margin-bottom: 0.75rem;
      }
      
      p {
        color: var(--text-secondary);
        font-size: 0.95rem;
      }
    }
  }
  
  // Sección de categorías
  .categories-section {
    margin-bottom: 5rem;
    
    .categories-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
      gap: 2rem;
    }
    
    .category-card {
      border-radius: var(--border-radius-lg);
      background-color: var(--card-bg);
      box-shadow: 0 4px 6px -1px rgba(0, 0, 0, 0.1), 
                  0 2px 4px -1px rgba(0, 0, 0, 0.06);
      padding: 2rem;
      transition: transform 0.3s ease, box-shadow 0.3s ease;
      display: flex;
      flex-direction: column;
      
      &:hover {
        transform: translateY(-5px);
        box-shadow: 0 10px 15px -3px rgba(0, 0, 0, 0.1),
                    0 4px 6px -2px rgba(0, 0, 0, 0.05);
                    
        .category-icon {
          transform: rotate(15deg);
        }
      }
      
      .category-icon {
        width: 60px;
        height: 60px;
        background: var(--primary-gradient);
        border-radius: var(--border-radius);
        display: flex;
        align-items: center;
        justify-content: center;
        color: white;
        margin-bottom: 1.5rem;
        transition: transform 0.3s var(--transition-bounce);
        
        i {
          font-size: 1.75rem;
        }
      }
      
      h3 {
        font-size: 1.35rem;
        color: var(--text-primary);
        margin-bottom: 1.25rem;
        font-weight: 600;
      }
      
      .category-list {
        list-style: none;
        padding: 0;
        margin: 0;
        
        li {
          margin-bottom: 0.5rem;
          
          &:last-child {
            margin-bottom: 0;
          }
          
          a {
            display: flex;
            align-items: center;
            color: var(--text-secondary);
            text-decoration: none;
            padding: 0.5rem 0;
            transition: color 0.2s ease;
            
            i {
              font-size: 0.75rem;
              margin-right: 0.75rem;
              color: var(--primary-color);
            }
            
            span {
              flex: 1;
            }
            
            &:hover {
              color: var(--primary-color);
            }
          }
        }
      }
    }
  }
  
  // Sección de FAQ
  .faq-section {
    margin-bottom: 5rem;
    
    .faq-container {
      max-width: 800px;
      margin: 0 auto;
    }
    
    .faq-item {
      margin-bottom: 1rem;
      border-radius: var(--border-radius-lg);
      overflow: hidden;
      box-shadow: 0 1px 3px rgba(0, 0, 0, 0.1);
      
      &.active {
        .faq-question {
          background-color: var(--primary-color);
          color: white;
          font-weight: 600;
        }
        
        .faq-icon {
          background-color: white;
          
          i {
            color: var(--primary-color);
          }
        }
      }
      
      .faq-question {
        padding: 1.25rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        cursor: pointer;
        font-weight: 500;
        color: var(--text-primary);
        background-color: var(--card-bg);
        transition: all 0.3s ease;
        
        &:hover {
          background-color: var(--bg-tertiary);
        }
        
        .faq-icon {
          width: 24px;
          height: 24px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          background-color: var(--bg-tertiary);
          transition: all 0.3s ease;
          
          i {
            font-size: 0.85rem;
            color: var(--primary-color);
          }
        }
      }
      
      .faq-answer {
        padding: 1.5rem;
        color: var(--text-secondary);
        line-height: 1.6;
        background-color: var(--card-bg);
        border-top: 1px solid var(--border-color);
      }
    }
  }
  
  // Sección de contacto
  .contact-section {
    margin-bottom: 5rem;
    
    .contact-card {
      max-width: 900px;
      margin: 0 auto;
      display: flex;
      flex-direction: column;
      background-color: var(--card-bg);
      border-radius: var(--border-radius-lg);
      box-shadow: 0 15px 30px rgba(0, 0, 0, 0.1);
      overflow: hidden;
      
      @media (min-width: 768px) {
        flex-direction: row;
      }
      
      .contact-graphics {
        position: relative;
        padding: 2rem;
        display: flex;
        align-items: center;
        justify-content: center;
        background: var(--secondary-gradient);
        flex: none;
        
        @media (min-width: 768px) {
          width: 30%;
        }
        
        .contact-image {
          width: 100px;
          height: 100px;
          background-color: rgba(255, 255, 255, 0.2);
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          z-index: 2;
          
          i {
            font-size: 3rem;
            color: white;
          }
        }
        
        .contact-decoration {
          position: absolute;
          width: 200px;
          height: 200px;
          border-radius: 50%;
          background-color: rgba(255, 255, 255, 0.1);
          top: 50%;
          left: 50%;
          transform: translate(-50%, -50%);
          z-index: 1;
        }
      }
      
      .contact-content {
        padding: 2.5rem;
        flex: 1;
        
        h3 {
          font-size: 1.75rem;
          color: var(--text-primary);
          margin-bottom: 0.75rem;
        }
        
        p {
          color: var(--text-secondary);
          margin-bottom: 2rem;
        }
        
        .contact-info {
          margin-bottom: 2rem;
          
          .contact-item {
            display: flex;
            align-items: center;
            margin-bottom: 1rem;
            
            i {
              width: 40px;
              height: 40px;
              border-radius: 50%;
              background-color: var(--bg-tertiary);
              display: flex;
              align-items: center;
              justify-content: center;
              margin-right: 1rem;
              color: var(--primary-color);
              font-size: 1.2rem;
            }
            
            span {
              color: var(--text-secondary);
            }
          }
        }
        
        .create-ticket-btn {
          background: var(--primary-gradient);
          color: white;
          border: none;
          padding: 0.85rem 1.75rem;
          border-radius: 99px;
          display: inline-flex;
          align-items: center;
          font-weight: 600;
          cursor: pointer;
          transition: transform 0.3s ease, box-shadow 0.3s ease;
          
          i {
            margin-right: 0.75rem;
          }
          
          &:hover {
            transform: translateY(-3px);
            box-shadow: 0 5px 15px rgba(0, 0, 0, 0.1);
          }
        }
      }
    }
  }
  
  // Footer con estilo
  .support-footer {
    position: relative;
    text-align: center;
    color: white;
    padding: 5rem 1rem 2rem;
    background: var(--primary-gradient);
    
    .footer-wave {
      position: absolute;
      top: -2px;
      left: 0;
      width: 100%;
      height: 3rem;
      transform: rotate(180deg);
      background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' viewBox='0 0 1200 120' preserveAspectRatio='none'%3E%3Cpath d='M0,0V46.29c47.79,22.2,103.59,32.17,158,28,70.36-5.37,136.33-33.31,206.8-37.5C438.64,32.43,512.34,53.67,583,72.05c69.27,18,138.3,24.88,209.4,13.08,36.15-6,69.85-17.84,104.45-29.34C989.49,25,1113-14.29,1200,52.47V0Z' fill='%23f8fafc' opacity='.25'%3E%3C/path%3E%3Cpath d='M0,0V15.81C13,36.92,27.64,56.86,47.69,72.05,99.41,111.27,165,111,224.58,91.58c31.15-10.15,60.09-26.07,89.67-39.8,40.92-19,84.73-46,130.83-49.67,36.26-2.85,70.9,9.42,98.6,31.56,31.77,25.39,62.32,62,103.63,73,40.44,10.79,81.35-6.69,119.13-24.28s75.16-39,116.92-43.05c59.73-5.85,113.28,22.88,168.9,38.84,30.2,8.66,59,6.17,87.09-7.5,22.43-10.89,48-26.93,60.65-49.24V0Z' fill='%23f8fafc' opacity='.5'%3E%3C/path%3E%3Cpath d='M0,0V5.63C149.93,59,314.09,71.32,475.83,42.57c43-7.64,84.23-20.12,127.61-26.46,59-8.63,112.48,12.24,165.56,35.4C827.93,77.22,886,95.24,951.2,90c86.53-7,172.46-45.71,248.8-84.81V0Z' fill='%23f8fafc'%3E%3C/path%3E%3C/svg%3E");
      background-size: cover;
      background-position: center;
    }
    
    .footer-content {
      position: relative;
      z-index: 2;
      
      .footer-tagline {
        font-size: 1.5rem;
        font-weight: 600;
        margin-bottom: 1rem;
      }
      
      .footer-info {
        font-size: 0.9rem;
        opacity: 0.8;
      }
    }
  }
  
  // Responsive adjustments
  @media (max-width: 768px) {
    .hero-section {
      padding: 3rem 1.5rem 5rem;
      
      .hero-title {
        font-size: 2.5rem;
      }
      
      .hero-subtitle {
        font-size: 1.25rem;
      }
    }
    
    .section-title {
      font-size: 1.75rem;
      flex-direction: column;
      
      .title-icon {
        margin-right: 0;
        margin-bottom: 0.75rem;
      }
    }
    
    .contact-section .contact-card {
      .contact-content {
        padding: 1.5rem;
      }
    }
  }
}
</style>
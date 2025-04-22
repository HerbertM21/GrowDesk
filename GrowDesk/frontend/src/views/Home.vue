<template>
  <div class="employee-support">
    <section class="hero">
      <h1>Portal de Soporte para Empleados</h1>
      <p class="subtitle">Recursos y asistencia exclusiva para el personal de GlassWorks</p>
      <div class="search-container">
        <div class="search-box">
          <input 
            v-model="searchQuery" 
            type="text" 
            placeholder="Buscar recursos de ayuda..." 
          />
          <button @click="search" class="search-button">
            <i class="pi pi-search"></i>
          </button>
        </div>
      </div>
    </section>

    <section class="quick-links">
      <h2>Enlaces Rápidos</h2>
      <div class="feature-grid">
        <div class="feature-card" v-for="(link, index) in quickLinks" :key="index">
          <i :class="link.icon"></i>
          <h3>{{ link.title }}</h3>
          <p>{{ link.description }}</p>
        </div>
      </div>
    </section>

    <section class="support-categories">
      <h2>Categorías de Soporte</h2>
      <div class="category-grid">
        <div class="category-card" v-for="(category, index) in categories" :key="index">
          <h3>
            <i :class="category.icon"></i>
            {{ category.title }}
          </h3>
          <ul class="category-items">
            <li v-for="(item, i) in category.items" :key="i">
              <a href="#">
                <i class="pi pi-chevron-right"></i>
                {{ item }}
              </a>
            </li>
          </ul>
        </div>
      </div>
    </section>

    <section class="faq-section">
      <h2>Preguntas Frecuentes</h2>
      <div class="faq-container">
        <div 
          v-for="(faq, index) in faqs" 
          :key="index" 
          class="faq-item"
          :class="{ 'active': openFaq === index }"
        >
          <div class="faq-question" @click="toggleFaq(index)">
            <span>{{ faq.question }}</span>
            <i class="pi" :class="openFaq === index ? 'pi-chevron-up' : 'pi-chevron-down'"></i>
          </div>
          <div class="faq-answer" v-show="openFaq === index">
            <p>{{ faq.answer }}</p>
          </div>
        </div>
      </div>
    </section>

    <section class="contact-support">
      <h2>¿Necesitas Más Ayuda?</h2>
      <div class="support-contact-card">
        <div class="support-icon">
          <i class="pi pi-phone"></i>
        </div>
        <div class="support-details">
          <h3>Equipo de Soporte TI</h3>
          <p>Disponible de Lunes a Viernes, 8am-5pm</p>
          <div class="contact-methods">
            <p><i class="pi pi-phone"></i> (555) 123-4567</p>
            <p><i class="pi pi-envelope"></i> soporte@glassworks.ejemplo</p>
            <p><i class="pi pi-phone-fill"></i> Extensión: 4321</p>
          </div>
          <button class="btn btn-primary">Crear Ticket de Soporte</button>
        </div>
      </div>
    </section>
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
  .hero {
    text-align: center;
    padding: 4rem 2rem;
    background: linear-gradient(135deg, #1976d2 0%, #1565c0 100%);
    color: white;
    border-radius: 8px;
    margin-bottom: 4rem;

    h1 {
      font-size: 3rem;
      margin-bottom: 1rem;
    }

    .subtitle {
      font-size: 1.5rem;
      margin-bottom: 2rem;
      opacity: 0.9;
    }

    .search-container {
      max-width: 600px;
      margin: 0 auto;
    }

    .search-box {
      display: flex;
      position: relative;
      
      input {
        width: 100%;
        padding: 1rem 1.5rem;
        border-radius: 50px;
        border: none;
        font-size: 1rem;
        box-shadow: 0 2px 10px rgba(0, 0, 0, 0.1);
      }

      .search-button {
        position: absolute;
        right: 5px;
        top: 5px;
        background: #1565c0;
        border: none;
        color: white;
        width: 40px;
        height: 40px;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
        
        &:hover {
          background: #0d47a1;
        }
      }
    }
  }

  .quick-links, .support-categories, .faq-section, .contact-support {
    padding: 4rem 2rem;
    
    h2 {
      text-align: center;
      margin-bottom: 3rem;
      color: #2c3e50;
    }
  }

  .quick-links {
    .feature-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(250px, 1fr));
      gap: 2rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .feature-card {
      text-align: center;
      padding: 2rem;
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      transition: transform 0.2s;

      &:hover {
        transform: translateY(-5px);
      }

      i {
        font-size: 2.5rem;
        color: #1976d2;
        margin-bottom: 1rem;
      }

      h3 {
        margin-bottom: 1rem;
        color: #2c3e50;
      }

      p {
        color: #666;
        line-height: 1.6;
      }
    }
  }

  .support-categories {
    background-color: #f5f5f5;

    .category-grid {
      display: grid;
      grid-template-columns: repeat(auto-fit, minmax(300px, 1fr));
      gap: 2rem;
      max-width: 1200px;
      margin: 0 auto;
    }

    .category-card {
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      padding: 2rem;

      h3 {
        display: flex;
        align-items: center;
        color: #2c3e50;
        margin-bottom: 1.5rem;
        font-size: 1.25rem;

        i {
          color: #1976d2;
          margin-right: 0.75rem;
          font-size: 1.5rem;
        }
      }

      .category-items {
        list-style: none;
        padding: 0;
        margin: 0;

        li {
          margin-bottom: 0.75rem;

          a {
            display: flex;
            align-items: center;
            color: #1976d2;
            text-decoration: none;
            transition: color 0.2s;

            &:hover {
              color: #0d47a1;
              text-decoration: underline;
            }

            i {
              font-size: 0.75rem;
              margin-right: 0.5rem;
            }
          }
        }
      }
    }
  }

  .faq-section {
    .faq-container {
      max-width: 800px;
      margin: 0 auto;
    }

    .faq-item {
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      margin-bottom: 1rem;
      overflow: hidden;

      &.active {
        .faq-question {
          background-color: #f0f7ff;
        }
      }

      .faq-question {
        padding: 1.5rem;
        display: flex;
        justify-content: space-between;
        align-items: center;
        cursor: pointer;
        font-weight: 500;
        color: #2c3e50;
        transition: background-color 0.2s;

        &:hover {
          background-color: #f0f7ff;
        }

        i {
          color: #1976d2;
        }
      }

      .faq-answer {
        padding: 0 1.5rem 1.5rem;
        color: #666;
        line-height: 1.6;
      }
    }
  }

  .contact-support {
    background-color: #f5f5f5;

    .support-contact-card {
      max-width: 800px;
      margin: 0 auto;
      background: white;
      border-radius: 8px;
      box-shadow: 0 2px 4px rgba(0, 0, 0, 0.1);
      padding: 2rem;
      display: flex;
      flex-direction: column;
      align-items: center;
      text-align: center;

      @media (min-width: 768px) {
        flex-direction: row;
        text-align: left;
      }

      .support-icon {
        width: 80px;
        height: 80px;
        background-color: #e3f2fd;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        margin-bottom: 1.5rem;
        margin-right: 0;

        @media (min-width: 768px) {
          margin-right: 2rem;
          margin-bottom: 0;
        }

        i {
          font-size: 2.5rem;
          color: #1976d2;
        }
      }

      .support-details {
        flex: 1;

        h3 {
          color: #2c3e50;
          margin-bottom: 0.5rem;
          font-size: 1.5rem;
        }

        p {
          color: #666;
          margin-bottom: 1.5rem;
        }

        .contact-methods {
          margin-bottom: 1.5rem;

          p {
            display: flex;
            align-items: center;
            margin-bottom: 0.5rem;

            i {
              color: #1976d2;
              margin-right: 0.75rem;
            }
          }
        }

        .btn-primary {
          background-color: #1976d2;
          color: white;
          border: none;
          padding: 0.75rem 1.5rem;
          border-radius: 4px;
          font-weight: 500;
          cursor: pointer;
          transition: background-color 0.2s;

          &:hover {
            background-color: #1565c0;
          }
        }
      }
    }
  }
}
</style>
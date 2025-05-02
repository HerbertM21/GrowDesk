# GrowDesk Widget

Widget de chat embebible para soporte al cliente, parte del sistema GrowDesk.

## Estructura del Proyecto

El proyecto está organizado de la siguiente manera:

- **widget-core**: Contiene el código fuente del widget como una librería JavaScript embebible.
- **widget-api**: API que sirve como intermediario entre el widget y el sistema principal de GrowDesk.
- **examples**: Ejemplos de uso del widget:
  - **simple-demo**: Una página HTML simple que muestra cómo integrar el widget.

## Requisitos

- Docker y Docker Compose
- Node.js 18 o superior (para desarrollo local)

## Inicio Rápido

1. Clonar el repositorio:
   ```bash
   git clone https://github.com/tu-usuario/growdesk-widget.git
   cd growdesk-widget
   ```

2. Iniciar los servicios con Docker Compose:
   ```bash
   docker compose up -d
   ```

3. Acceder a la demo:
   - Demo site: http://localhost:8090

## Desarrollo

### Widget Core

El widget está construido como una librería utilizando Vue 3 y Vite. Para desarrollo:

```bash
cd widget-core
npm install
npm run dev
```

### Widget API

La API está construida con Node.js. Para desarrollo:

```bash
cd widget-api
npm install
npm run dev
```

## Uso del Widget en un Sitio Web

Para incluir el widget en un sitio web, añade el siguiente código antes del cierre del `</body>`:

```html
<!-- Configuración global para el widget -->
<script>
  window.GrowDeskConfig = {
    brandName: "Tu Empresa",
    primaryColor: "#6200ea",
    position: "bottom-right",
    welcomeMessage: "¿En qué podemos ayudarte hoy?"
  };
</script>

<!-- Script del Widget GrowDesk -->
<script 
  src="https://cdn.tudominio.com/growdesk-widget.umd.js" 
  id="growdesk-widget"
  data-widget-id="tu-widget-id" 
  data-widget-token="tu-token" 
  data-api-url="https://api.tudominio.com"
></script>
```

## Licencia

MIT 
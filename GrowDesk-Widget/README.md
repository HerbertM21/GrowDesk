# GrowDesk-Widget

GrowDesk-Widget es un componente de chat embebible para integrar soporte al cliente en cualquier sitio web. Este widget se conecta con el sistema de gestión de tickets GrowDesk para ofrecer una experiencia de soporte completa.

## Estructura del Proyecto

El proyecto está organizado en tres componentes principales:

- **widget-src**: Código fuente del widget desarrollado con Vue.js, listo para ser integrado en cualquier sitio web.
- **widget-api**: API dedicada para manejar las comunicaciones entre el widget y el sistema principal de GrowDesk.
- **demo-site**: Sitio de demostración para probar el widget en un entorno real.

## Características

- Widget de chat flotante y personalizable
- Formulario de contacto para visitantes anónimos
- Chat en tiempo real con agentes de soporte
- Integración simple con cualquier sitio web (solo añadiendo un script)
- Personalización de apariencia (colores, mensajes, posición)
- Opciones de seguridad para dominios permitidos
- **Comunicación con GrowDesk**: El widget envía tickets y mensajes al sistema principal de GrowDesk

## Arquitectura y Puertos del Sistema

El sistema GrowDesk-Widget utiliza varios componentes que se comunican entre sí. A continuación, se describe la configuración de puertos y servicios:

### Componentes y Puertos

| Componente | Puerto | Descripción |
|------------|--------|-------------|
| **API del Widget** | 8082 | API que funciona como intermediario entre el widget y el sistema GrowDesk. Recibe mensajes del widget, crea tickets, etc. |
| **Backend de GrowDesk** | 8000 | Sistema principal de GrowDesk donde se procesan los tickets y se maneja la lógica de negocio. |
| **Frontend de GrowDesk** | 3000 | Panel de administración donde los agentes pueden ver y responder tickets. |
| **Widget (Desarrollo)** | 3030 | Servidor de desarrollo para el widget. |
| **Demo Site** | 8090 | Página de demostración que implementa el widget. |

### Flujo de Comunicación

1. El usuario envía un mensaje a través del widget en un sitio web.
2. El widget envía el mensaje a la **API del Widget** (puerto 8082).
3. La API del Widget crea un ticket o añade el mensaje a un ticket existente en el **Backend de GrowDesk** (puerto 8000).
4. Los agentes de soporte ven el nuevo ticket en el **Frontend de GrowDesk** (puerto 3000).
5. Cuando responden, la respuesta se envía al usuario a través de la API del Widget al widget en el sitio web.

### Implementación

Para implementar el widget en un sitio web, se debe incluir el siguiente código:

```html
<script 
    src="URL_DEL_WIDGET/widget.js" 
    id="growdesk-widget"
    data-widget-id="ID_DEL_WIDGET" 
    data-widget-token="TOKEN_DEL_WIDGET" 
    data-api-url="URL_DE_LA_API"
>
</script>
```

Este código puede ser generado automáticamente desde el panel de administración de GrowDesk en la sección de configuración del widget.

## Integración

Para integrar el widget en tu sitio web, simplemente añade el siguiente código a tu HTML:

```html
<script src="https://cdn.growdesk.com/widget.js" id="growdesk-widget"
  data-brand-name="Mi Empresa"
  data-welcome-message="¿En qué podemos ayudarte hoy?"
  data-primary-color="#4caf50"
  data-position="bottom-right"
  data-api-key="tu-api-key">
</script>
```

## Comunicación con GrowDesk

El widget se comunica con el sistema principal de GrowDesk a través de su API dedicada. Cuando un usuario envía un mensaje o crea un ticket desde el widget:

1. El widget envía los datos a la API del widget (`widget-api`)
2. La API del widget procesa la información y la reenvía al sistema de GrowDesk
3. GrowDesk recibe y gestiona los tickets/mensajes en su panel de administración
4. Los agentes de soporte pueden responder desde GrowDesk y las respuestas se muestran en el widget

### Configuración de la comunicación

Para configurar la comunicación entre el Widget y GrowDesk:

1. En el archivo `.env` de widget-api, configura:
   ```
   GROWDESK_API_URL=http://tu-servidor-growdesk/api
   GROWDESK_API_KEY=tu-api-key-de-growdesk
   ```

2. En GrowDesk, asegúrate de tener habilitada la API para widgets:
   - Configura `WIDGET_API_KEY` en las variables de entorno de GrowDesk
   - Habilita los endpoints necesarios en la configuración

## Tecnologías

- **Frontend**: Vue.js 3, TypeScript, SCSS
- **API**: Golang con Gin
- **Comunicación en tiempo real**: WebSockets

## Desarrollo

### Requisitos previos

- Node.js 18 o superior
- Go 1.20 o superior
- Docker (opcional para despliegue)

### Configuración del entorno de desarrollo

1. Clona este repositorio
2. Configura el widget-src:
   ```bash
   cd widget-src
   npm install
   npm run dev
   ```

3. Configura el widget-api:
   ```bash
   cd widget-api
   go mod tidy
   go run main.go
   ```

4. Lanza el sitio de demostración:
   ```bash
   cd demo-site
   npm install
   npm run dev
   ```

## Licencia

MIT

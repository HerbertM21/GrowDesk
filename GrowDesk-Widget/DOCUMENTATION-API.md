# Documentación de la API de GrowDesk Widget

## Descripción General

GrowDesk Widget es un componente embebible de chat para soporte al cliente que se integra con el sistema de gestión de tickets GrowDesk. Este widget permite a los visitantes de cualquier sitio web crear tickets de soporte que son gestionados por el equipo de atención al cliente a través de GrowDesk.

## Arquitectura

El sistema consta de dos componentes principales:

1. **GrowDesk**: Sistema principal de gestión de tickets
   - Backend (Go/Gin)
   - Frontend (Vue.js)
   - Base de datos (PostgreSQL)

2. **GrowDesk Widget**: Widget embebible para sitios web
   - widget-src: Componente frontend (Vue.js)
   - widget-api: API para comunicación (Go/Gin)
   - demo-site: Sitio web de ejemplo

## Configuración y Ejecución

### Requisitos Previos

- Go 1.20+
- Node.js 18+
- PostgreSQL 15+
- Docker (opcional)

### Iniciar la API del Widget

```bash
# Acceder al directorio de la API
cd /home/hmdev/GrowProyect/GrowDesk-Widget/widget-api

# Instalar dependencias (si es necesario)
go mod tidy

# Configurar variables de entorno
cp .env.example .env
# Editar .env según sea necesario, especialmente:
# GROWDESK_API_URL=http://localhost:8000/api
# GROWDESK_API_KEY=tu_api_key_aquí

# Ejecutar la API (por defecto en puerto 8080)
go run main.go

# O especificar un puerto diferente
PORT=8081 go run main.go
```

### Iniciar el Frontend del Widget

```bash
# Acceder al directorio del widget
cd /home/hmdev/GrowProyect/GrowDesk-Widget/widget-src

# Instalar dependencias
npm install

# Ejecutar en modo desarrollo
npm run dev

# Construir para producción
npm run build
```

### Iniciar el Sistema GrowDesk

```bash
# Acceder al directorio de GrowDesk
cd /home/hmdev/GrowProyect/GrowDesk

# Iniciar backend
cd backend
go mod tidy
go run main.go

# Iniciar frontend (en otra terminal)
cd frontend
npm install
npm run dev

# Alternativamente, usar Docker Compose
cd /home/hmdev/GrowProyect/GrowDesk
docker-compose up -d
```

## Endpoints de la API del Widget

| Método | Endpoint | Descripción | Autenticación |
|--------|----------|-------------|---------------|
| GET | `/` | Información sobre la API | No |
| GET | `/widget/status` | Comprobar estado de la API | Sí |
| POST | `/widget/tickets` | Crear un nuevo ticket | Sí |
| POST | `/widget/messages` | Enviar mensaje a un ticket existente | Sí |
| GET | `/widget/tickets/:ticketId/messages` | Obtener mensajes de un ticket | Sí |

### Autenticación

Todas las peticiones al widget (excepto la ruta raíz) requieren dos cabeceras de autenticación:

```
X-Widget-ID: tu-widget-id
X-Widget-Token: tu-widget-token
```

Para el entorno de prueba, puedes usar:
```
X-Widget-ID: demo-widget
X-Widget-Token: demo-token
```

## Integración en un Sitio Web

### Método Básico

Añade el siguiente script en tu HTML (preferiblemente antes del cierre de `</body>`):

```html
<script src="https://cdn.growdesk.com/widget.js" id="growdesk-widget"
  data-widget-id="tu-widget-id"
  data-widget-token="tu-widget-token"
  data-brand-name="Nombre de tu Empresa"
  data-welcome-message="¿En qué podemos ayudarte hoy?"
  data-primary-color="#4caf50"
  data-position="bottom-right">
</script>
```

### Para Desarrollo Local

```html
<script src="http://localhost:3000/widget.js" id="growdesk-widget"
  data-widget-id="demo-widget"
  data-widget-token="demo-token"
  data-brand-name="Mi Empresa"
  data-welcome-message="¿En qué podemos ayudarte hoy?"
  data-primary-color="#4caf50"
  data-position="bottom-right">
</script>
```

## Comunicación entre Widget y GrowDesk

El flujo de comunicación funciona así:

1. **Usuario → Widget**: El usuario interactúa con el widget en el sitio web
2. **Widget → widget-api**: El widget envía los datos a la API del widget
3. **widget-api → GrowDesk**: La API del widget reenvía los datos al sistema GrowDesk
4. **GrowDesk → Agentes**: Los agentes de soporte ven y responden a los tickets
5. **GrowDesk → widget-api → Widget → Usuario**: Las respuestas llegan al widget

Para configurar esta comunicación:

1. En el archivo `.env` de widget-api:
   ```
   GROWDESK_API_URL=http://localhost:8000/api
   GROWDESK_API_KEY=tu-api-key-de-growdesk
   ```

2. En GrowDesk:
   ```
   WIDGET_API_KEY=el-mismo-api-key
   ```

## Pruebas

### Probar la API del Widget

```bash
# Comprobar estado
curl http://localhost:8081/widget/status -H "X-Widget-ID: demo-widget" -H "X-Widget-Token: demo-token"

# Crear un ticket
curl -X POST http://localhost:8081/widget/tickets \
  -H "X-Widget-ID: demo-widget" \
  -H "X-Widget-Token: demo-token" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Usuario Prueba",
    "email": "test@example.com",
    "message": "Necesito ayuda con mi cuenta",
    "metadata": {"url": "https://ejemplo.com/pagina"}
  }'

# Enviar mensaje a un ticket existente
curl -X POST http://localhost:8081/widget/messages \
  -H "X-Widget-ID: demo-widget" \
  -H "X-Widget-Token: demo-token" \
  -H "Content-Type: application/json" \
  -d '{
    "ticketId": "TICKET-XXXXXXXX-XXXXXX",
    "message": "¿Alguna novedad sobre mi consulta?"
  }'

# Obtener mensajes de un ticket
curl http://localhost:8081/widget/tickets/TICKET-XXXXXXXX-XXXXXX/messages \
  -H "X-Widget-ID: demo-widget" \
  -H "X-Widget-Token: demo-token"
```

### Probar el Sitio de Demostración

1. Accede al directorio demo-site
   ```bash
   cd /home/hmdev/GrowProyect/GrowDesk-Widget/demo-site
   ```

2. Sirve el sitio con un servidor web simple
   ```bash
   python -m http.server 8000
   ```

3. Abre un navegador y visita `http://localhost:8000`
   - Deberías ver el sitio demo con el widget de chat en la esquina inferior derecha

## Solución de Problemas

### 404 Page Not Found

- Si obtienes un 404 al acceder a la raíz, es normal si no has configurado un endpoint para `/`
- Asegúrate de estar utilizando las rutas correctas: `/widget/status`, `/widget/tickets`, etc.

### Error de Conexión con GrowDesk

Si ves mensajes como:
```
Error al enviar datos a GrowDesk: Post "http://growdesk-backend:8000/api/tickets": dial tcp: lookup growdesk-backend: no such host
```

- Verifica la URL en `GROWDESK_API_URL`
- Para desarrollo local, usa `http://localhost:8000/api` en lugar de `http://growdesk-backend:8000/api`
- Asegúrate de que GrowDesk esté en ejecución

### Puerto en Uso

Si ves:
```
[GIN-debug] [ERROR] listen tcp :8080: bind: address already in use
```

- Cambia el puerto con la variable de entorno `PORT`
- Asegúrate de no tener instancias duplicadas en ejecución
- Usa comandos como `lsof -i :8080` para ver qué proceso está usando el puerto

## Mejores Prácticas

1. **Seguridad**:
   - En producción, limita los dominios permitidos en CORS
   - Utiliza HTTPS para todas las comunicaciones
   - Nunca expongas claves API en el código del cliente

2. **Rendimiento**:
   - Implementa caché para respuestas frecuentes
   - Optimiza el tamaño del widget para carga rápida
   - Considera utilizar un CDN para servir el script del widget

3. **Experiencia de Usuario**:
   - Personaliza los colores del widget para que coincidan con tu marca
   - Configura mensajes de bienvenida amigables
   - Implementa notificaciones para nuevos mensajes

## Desarrollo Futuro

- Implementación de WebSockets para comunicación en tiempo real
- Soporte para archivos adjuntos en los mensajes
- Chatbot con IA para respuestas automáticas a preguntas comunes
- Ampliación de las opciones de personalización del widget
- Añadir soporte para múltiples idiomas

## Ejemplo de Logs de Ejecución

```
[GIN-debug] [WARNING] Creating an Engine instance with the Logger and Recovery middleware already attached.
[GIN-debug] [WARNING] Running in "debug" mode. Switch to "release" mode in production.
 - using env:   export GIN_MODE=release
 - using code:  gin.SetMode(gin.ReleaseMode)
[GIN-debug] GET    /                         --> main.main.func1 (3 handlers)
[GIN-debug] GET    /widget/status            --> main.main.func3 (5 handlers)
[GIN-debug] POST   /widget/tickets           --> main.main.func4 (5 handlers)
[GIN-debug] POST   /widget/messages          --> main.main.func5 (5 handlers)
[GIN-debug] GET    /widget/tickets/:ticketId/messages --> main.main.func6 (5 handlers)
2025/03/21 20:51:55 Servidor iniciado en el puerto 8081
[GIN-debug] [WARNING] You trusted all proxies, this is NOT safe. We recommend you to set a value.
Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.
[GIN-debug] Listening and serving HTTP on :8081
[GIN] 2025/03/21 - 20:50:26 | 404 |       1.201µs |             ::1 | GET      "/"
[GIN] 2025/03/21 - 20:50:56 | 200 |     466.271µs |             ::1 | GET      "/widget/status"
[GIN] 2025/03/21 - 20:51:06 | 201 |     211.868µs |             ::1 | POST     "/widget/tickets"
2025/03/21 20:51:09 Error al enviar datos a GrowDesk: Post "http://growdesk-backend:8000/api/tickets": dial tcp: lookup growdesk-backend: no such host
[GIN] 2025/03/21 - 20:51:19 | 200 |  3.686858762s |             ::1 | GET      "/widget/tickets/TICKET-20250321-205106/messages"
```

---

Este documento proporciona una visión general de cómo funciona la API de GrowDesk Widget. Para preguntas específicas o asistencia adicional, contacta al equipo de desarrollo. 
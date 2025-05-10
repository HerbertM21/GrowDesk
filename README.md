# GrowDesk V2

Sistema de mesa de ayuda y soporte al cliente con widget embebible para sitios web.

## Descripción General

GrowDesk es una plataforma completa para la gestión de tickets de soporte que incluye:

- Panel de administración para agentes y administradores
- API RESTful para gestión de tickets, usuarios y categorías
- Widget embebible para integración en sitios web
- Comunicación en tiempo real mediante WebSockets
- Sistema de chat integrado

## Estructura del Proyecto

El proyecto está organizado en tres componentes principales:

- **GrowDesk**: Core del sistema (backend y frontend principal)
- **GrowDesk-Widget**: Widget embebible y su API
- **api-gateway**: Puerta de enlace API con Traefik

## Componentes y Puertos

| Componente      | Puerto  | Descripción                                  | URL Directa              | URL Gateway            |
|-----------------|---------|----------------------------------------------|--------------------------|------------------------|
| Backend API     | 8080    | API principal para el sistema de tickets     | http://localhost:8080    | http://localhost/api   |
| Frontend Admin  | 3001    | Panel de administración                      | http://localhost:3001    | http://localhost/admin |
| Widget API      | 3002    | API para el widget                           | http://localhost:3002    | -                      |
| Widget Core     | 3030    | Widget embebible JavaScript                  | http://localhost:3030    | http://localhost/widget|
| Demo Site       | 8091    | Sitio de demostración para el widget         | http://localhost:8091    | http://localhost/demo  |
| API Gateway     | 80      | Puerta de enlace Traefik                     | http://localhost         | -                      |
| Traefik Dashboard | 8082  | Panel de control de Traefik                  | http://localhost:8082    | http://localhost/dashboard |
| PostgreSQL      | 5432    | Base de datos                                | localhost:5432           | -                      |

## Requisitos Previos

- Docker v20.10.0 o superior
- Docker Compose v2.0.0 o superior
- Git

## Instalación y Ejecución

### 1. Clonar el repositorio

```bash
git clone https://github.com/HerbertM21/GrowDesk.git
cd GrowDesk
```

### 2. Iniciar el sistema completo

El script `start-system.sh` automatiza el proceso de inicio de todos los componentes:

```bash
chmod +x start-system.sh
./start-system.sh
```

Alternativamente, puedes iniciar cada componente manualmente:

```bash
# Iniciar todos los servicios
docker-compose up -d

# O iniciar servicios específicos
docker-compose up -d api-gateway postgres backend frontend widget-api widget-core demo-site
```

### 3. Verificar que los contenedores estén funcionando

```bash
docker ps
```

Deberías ver todos los contenedores en estado "Up".

### 4. Acceder al panel de administración

Abre en tu navegador: http://localhost/admin o http://localhost:3001

**Credenciales por defecto:**
- Usuario: admin@example.com
- Contraseña: password

## Configuración

### Variables de Entorno

El sistema utiliza archivos `.env` para configuración en cada componente. Los principales son:

**Proyecto principal (.env):**
```
COMPOSE_PROJECT_NAME=growdeskv2
```

**Backend (.env):**
```
# PostgreSQL
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=growdesk

# Aplicación
PORT=8080
DATA_DIR=/app/data
MOCK_AUTH=true
JWT_SECRET=your_jwt_secret
LOG_LEVEL=debug
```

**Frontend (.env):**
```
VITE_API_URL=http://localhost/api
VITE_WS_URL=ws://localhost/api/ws
VITE_APP_NAME=GrowDesk
VITE_APP_VERSION=1.0.0
```

**Widget API (.env):**
```
PORT=3000
DATA_DIR=/app/data
GROWDESK_API_URL=http://growdesk-backend:8080
ALLOWED_ORIGINS=http://localhost:3000,http://localhost:5173,http://localhost:3030,http://localhost:8091
```

### Configuración del Widget para tu Sitio Web

Para integrar el widget en tu sitio, añade este código antes del cierre de `</body>`:

```html
<script id="growdesk-widget" 
        src="http://localhost/widget/growdesk-widget.umd.js" 
        data-api-url="http://localhost:3002/widget" 
        data-widget-id="testwidget" 
        data-widget-token="testtoken" 
        async>
</script>
```

## Base de Datos

GrowDesk utiliza PostgreSQL como base de datos principal:

- Los datos se persisten en un volumen Docker
- La primera vez que se inicia, se migran automáticamente datos de ejemplo
- Puedes conectarte directamente con un cliente SQL:
  ```
  Host: localhost
  Puerto: 5432
  Usuario: postgres
  Contraseña: postgres
  Base de datos: growdesk
  ```

Para acceder desde la línea de comandos:

```bash
docker exec -it growdesk-db psql -U postgres -d growdesk
```

## Solución de Problemas Comunes

### Widget API no Accesible a través del Gateway

El Widget API (puerto 3002) debe accederse directamente debido a limitaciones en la configuración del gateway. Para resolver esto:

1. Asegúrate de que el widget esté configurado para conectarse a http://localhost:3002/widget
2. Verifica que no haya un firewall bloqueando el puerto 3002

### Errores de CORS

Si encuentras errores de CORS en el frontend o widget:

1. Verifica la configuración CORS en `api-gateway/dynamic_conf/api_routes.yml`
2. Asegúrate de que tu origen esté permitido en las cabeceras `Access-Control-Allow-Origin`

### Tickets No Aparecen en el Panel de Administración

Si los tickets creados desde el widget no aparecen en el panel:

1. Verifica que el ticket existe en la base de datos:
   ```bash
   docker exec -it growdesk-db psql -U postgres -d growdesk -c "SELECT id, title FROM tickets ORDER BY created_at DESC LIMIT 5"
   ```
2. Comprueba que el frontend esté utilizando la URL de API correcta (http://localhost/api)
3. Reinicia el frontend para que refresque su caché:
   ```bash
   docker restart growdesk-frontend
   ```

### Problemas de WebSocket

Si el chat en tiempo real no funciona:

1. Asegúrate de que los puertos para WebSockets estén abiertos
2. Verifica la configuración de WebSocket en el frontend y backend
3. Comprueba que el API Gateway esté configurado para permitir conexiones WebSocket

## Desarrollo y Extensión

### Desarrollo en Modo Local

Para desarrollo, puedes ejecutar componentes individuales en modo local:

```bash
# Backend
cd GrowDesk/backend
make run

# Frontend
cd GrowDesk/frontend
npm install
npm run dev

# Widget Core
cd GrowDesk-Widget/widget-core
npm install
npm run dev
```

### Reconstruir Contenedores Específicos

```bash
docker-compose build backend frontend
docker-compose up -d backend frontend
```

### Ver Logs

```bash
# Ver logs de un servicio específico
docker logs -f growdesk-backend

# Ver logs de múltiples servicios
docker-compose logs -f backend frontend
```

## Documentación Técnica

Consulta estos documentos para una comprensión más profunda de la arquitectura:

- [Arquitectura del Backend](docs/backend_architecture.md)
- [Arquitectura del Frontend](docs/frontend_architecture.md)
- [API Documentation](docs/API.md)

## Estructura de Directorios Completa

```
GrowDeskV2/
├── docs/                   # Documentación del proyecto
├── GrowDesk/               # Sistema principal
│   ├── backend/            # API backend en Go
│   │   ├── api/            # Definiciones API
│   │   ├── cmd/            # Puntos de entrada
│   │   ├── internal/       # Código interno
│   │   └── ...
│   ├── frontend/           # Panel de admin Vue.js
│   │   ├── public/
│   │   ├── src/
│   │   └── ...
│   └── docker-compose.yml  # Configuración Docker
├── GrowDesk-Widget/        # Widget de chat
│   ├── widget-api/         # API del widget (Go)
│   ├── widget-core/        # Core del widget (Vue.js)
│   ├── examples/           # Ejemplos de uso
│   └── docker-compose.yml  # Configuración Docker
├── api-gateway/            # Configuración Traefik
│   ├── dynamic_conf/       # Configuraciones dinámicas
│   ├── traefik.yml         # Configuración principal
│   └── ...
└── docker-compose.yml      # Configuración principal
```
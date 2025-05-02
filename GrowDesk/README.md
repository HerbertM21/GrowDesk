# GrowDesk - Sistema de Tickets y Soporte

GrowDesk es una plataforma moderna de help desk y gestión de tickets de soporte, diseñada para simplificar la comunicación con clientes y la gestión eficiente de solicitudes de soporte.

## Requisitos previos

- [Docker](https://docs.docker.com/get-docker/)
- [Docker Compose](https://docs.docker.com/compose/install/)

## Cómo ejecutar el programa

Sigue estos pasos para iniciar la aplicación GrowDesk en tu entorno local:

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/GrowDeskV2.git
cd GrowDeskV2/GrowDesk
```

### 2. Iniciar los servicios con Docker Compose

```bash
docker compose up -d
```

Este comando iniciará todos los servicios necesarios:
- **Frontend**: Interfaz de usuario Vue.js (accesible en http://localhost:3000)
- **Backend**: API y servidor de mock para desarrollo (accesible en http://localhost:8080)
- **Base de datos**: PostgreSQL (puerto 5433)
- **Redis**: Para caché y gestión de sesiones (puerto 6379)

### 3. Verificar que los contenedores estén funcionando

```bash
docker ps
```

Deberías ver todos los contenedores ejecutándose: `growdesk-frontend`, `growdesk-backend`, `growdesk-db`, y `growdesk-redis`.

### 4. Acceder a la aplicación

Abre tu navegador web y navega a:
- **Aplicación principal**: http://localhost:3000
- **API Backend (para pruebas)**: http://localhost:8080/api/health

### 5. Credenciales de inicio de sesión

La aplicación viene preconfigurada con estos usuarios de prueba:

| Email | Contraseña | Rol |
|-------|------------|-----|
| admin@example.com | password | Administrador |
| asistente@example.com | password | Asistente |
| empleado@example.com | password | Empleado |

### 6. Detener los servicios

Cuando hayas terminado de usar la aplicación, puedes detener los servicios con:

```bash
docker compose down
```

Si quieres eliminar todos los volúmenes y datos persistentes:

```bash
docker compose down -v
```

## Solución de problemas

### El frontend no carga o muestra errores

Si tienes problemas para acceder al frontend:

```bash
# Ver los logs del frontend
docker logs growdesk-frontend

# Reiniciar solo el frontend
docker compose restart frontend
```

### El backend no responde

Si la API no está respondiendo:

```bash
# Ver los logs del backend
docker logs growdesk-backend

# Reiniciar solo el backend
docker compose restart backend
```

### Limpiar completamente y reconstruir los contenedores

Si necesitas reiniciar desde cero:

```bash
# Detener y eliminar todos los contenedores
docker compose down

# Eliminar volúmenes
docker compose down -v

# Reconstruir las imágenes
docker compose up -d --build
```

## Características principales

- Sistema de tickets multi-usuario
- Chat en tiempo real con WebSockets
- Gestión de usuarios y permisos
- Panel de administración
- Base de conocimientos y FAQs
- Widget personalizable para sitios web

## Estructura del proyecto

- `frontend/` - Aplicación Vue.js (TypeScript)
- `backend/` - Servidor Node.js/Express (modo desarrollo)
- `docker-compose.yml` - Configuración de Docker Compose 
# GrowDeskV2

## Sistema Integrado de Gestión de Tickets y Soporte al Cliente

GrowDeskV2 es una plataforma moderna de help desk y gestión de tickets de soporte, diseñada para simplificar la comunicación con clientes y la gestión eficiente de solicitudes de soporte.

## Requisitos Previos

- **Docker**: versión 20.10.0 o superior
- **Docker Compose**: versión 2.0.0 o superior
- **Git**: para clonar el repositorio

## Estructura del Proyecto

El proyecto está organizado en dos componentes principales:

- **GrowDesk**: Sistema principal de gestión de tickets
- **GrowDesk-Widget**: Widget embebible para sitios web de clientes

## Inicio Rápido con Docker

### 1. Clonar el repositorio

```bash
git clone https://github.com/HerbertM21/GrowDesk.git
cd GrowDeskV2
```

### 2. Iniciar el Sistema Principal (GrowDesk)

```bash
# Navegar al directorio GrowDesk
cd GrowDesk

# Iniciar todos los servicios con Docker Compose
docker compose up -d
```

Este comando iniciará todos los servicios necesarios:
- **Frontend**: Interfaz de usuario Vue.js (accesible en http://localhost:3000)
- **Backend**: API y servidor mock para desarrollo (accesible en http://localhost:8080)
- **Base de datos**: PostgreSQL (puerto 5433)
- **Redis**: Para caché y gestión de sesiones (puerto 6379)

### 3. Iniciar el Widget (opcional)

```bash
# Navegar al directorio GrowDesk-Widget
cd ../GrowDesk-Widget

# Iniciar todos los servicios para el widget
docker compose up -d
```

Este comando iniciará todos los servicios del widget:
- **Widget Core**: Núcleo del widget (accesible en http://localhost:5174)
- **Widget API**: API para la comunicación del widget (accesible en http://localhost:8082)
- **Demo Site**: Sitio web de demostración con el widget integrado (accesible en http://localhost:8090)
- **Base de datos**: PostgreSQL para el widget (puerto 5434)

## Acceso a las Aplicaciones

Una vez que todos los contenedores estén en funcionamiento, puedes acceder a:

- **Panel de Administración**: http://localhost:3000
- **API Backend (para pruebas)**: http://localhost:8080/api/health
- **Demo del Widget**: http://localhost:8090

## Credenciales de Prueba

Para acceder al sistema de administración, usa las siguientes credenciales:

- **Administrador**:
  - Email: admin@example.com
  - Contraseña: password
- **Asistente**:
  - Email: asistente@example.com
  - Contraseña: password
- **Empleado**:
  - Email: empleado@example.com
  - Contraseña: password

## Gestión de los Contenedores

### Detener los servicios

```bash
# Para el sistema principal
cd GrowDesk
docker compose down

# Para el widget
cd ../GrowDesk-Widget
docker compose down
```

### Ver logs

```bash
# Logs del frontend
docker logs growdesk-frontend

# Logs del backend
docker logs growdesk-backend
```

### Reiniciar servicios específicos

```bash
# Reiniciar solo el frontend
docker compose restart frontend

# Reiniciar solo el backend
docker compose restart backend
```

### Limpieza completa

Para eliminar todos los contenedores, redes y volúmenes:

```bash
# Sistema principal
cd GrowDesk
docker compose down -v

# Widget
cd ../GrowDesk-Widget
docker compose down -v
```

## Configuración de Puertos

| Componente            | Puerto | Descripción                          |
| --------------------- | ------ | ------------------------------------- |
| Frontend              | 3000   | Panel de administración de GrowDesk  |
| Backend               | 8080   | Servidor mock para datos de prueba    |
| PostgreSQL (GrowDesk) | 5433   | Base de datos para el sistema principal |
| Redis                 | 6379   | Caché y gestión de sesiones          |
| Widget Core           | 5174   | Interfaz del widget                  |
| Widget API            | 8082   | API para la comunicación del widget  |
| Demo Site             | 8090   | Sitio web de demostración con widget |
| PostgreSQL (Widget)   | 5434   | Base de datos para el widget         |

# GrowDeskV2

Sistema completo de mesa de ayuda con panel de administración y widget de chat para integración en sitios web.

## Componentes

El sistema está compuesto por dos componentes principales:

1. **GrowDesk**: Panel de administración y backend principal
   - Frontend: Panel de administración
   - Backend: API REST y servidor de WebSockets
   - Sync-Server: Servidor de sincronización entre instancias

2. **GrowDesk-Widget**: Widget de chat para sitios web
   - Widget-Core: Componente frontend del widget
   - Widget-API: Backend para el widget
   - Demo-Site: Sitio de demostración para probar el widget

## Requisitos

- Docker
- Docker Compose

## Configuración

El sistema incluye archivos de configuración predeterminados, pero puedes personalizarlos modificando los archivos `.env` en cada directorio.

### Variables de entorno principales

- **PostgreSQL**:
  - `DB_HOST`: Host de la base de datos (por defecto: postgres)
  - `DB_PORT`: Puerto de la base de datos (por defecto: 5432)
  - `DB_USER`: Usuario de la base de datos (por defecto: postgres)
  - `DB_PASSWORD`: Contraseña de la base de datos (por defecto: postgres)
  - `DB_NAME`: Nombre de la base de datos (por defecto: growdesk)

- **Backend**:
  - `PORT`: Puerto del backend (por defecto: 8080)
  - `DATA_DIR`: Directorio de datos (por defecto: /app/data)
  - `MOCK_AUTH`: Usar autenticación simulada para desarrollo (por defecto: true)

## Ejecución

Para ejecutar todo el sistema:

```bash
docker compose up
```

Para ejecutar en segundo plano:

```bash
docker compose up -d
```

## Acceso a los servicios

Una vez que el sistema esté en ejecución, podrás acceder a los siguientes servicios:

- **Panel de administración**: http://localhost:3001
- **API Backend**: http://localhost:8080
- **Widget Demo**: http://localhost:8090
- **Widget API**: http://localhost:3000
- **Widget Core**: http://localhost:3030

## Base de datos PostgreSQL

El sistema utiliza PostgreSQL para almacenar todos los datos. La base de datos está configurada para persistir los datos en un volumen Docker, por lo que no se perderán al reiniciar los contenedores.

### Datos iniciales

La primera vez que se inicie el sistema, se migrarán automáticamente los datos de los archivos JSON a PostgreSQL si existen.

### Acceso directo a PostgreSQL

Puedes acceder directamente a la base de datos mediante:

```bash
docker exec -it growdesk-db psql -U postgres -d growdesk
```

## Desarrollo

Para desarrollar el sistema, puedes utilizar los siguientes comandos:

### Reconstruir servicios específicos

```bash
docker compose build frontend backend
docker compose up -d frontend backend
```

### Ver logs

```bash
docker compose logs -f backend
```

## Estructura de directorios

```
GrowDeskV2/
├── GrowDesk/                # Panel de administración y backend
│   ├── backend/             # Backend API
│   ├── frontend/            # Panel de administración
│   └── docker-compose.yml   # Configuración de Docker para GrowDesk
├── GrowDesk-Widget/         # Widget de chat
│   ├── widget-api/          # API del widget
│   ├── widget-core/         # Componente frontend del widget
│   ├── examples/            # Ejemplos de integración
│   └── docker-compose.yml   # Configuración de Docker para Widget
└── docker-compose.yml       # Configuración de Docker principal
```

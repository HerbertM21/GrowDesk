# GrowDesk V2

Sistema de mesa de ayuda y soporte al cliente con widget embebible.

## Estructura del Proyecto

El proyecto está organizado en dos carpetas principales:

- **GrowDesk**: Contiene el backend y frontend principal
- **GrowDesk-Widget**: Contiene el widget embebible y su API

## Componentes

- **Backend API**: API principal para el sistema de tickets (puerto 8080)
- **Frontend Admin**: Panel de administración (puerto 3001)
- **Widget API**: API para el widget (puerto 3002)
- **Widget Core**: Widget embebible JavaScript (puerto 3030)
- **Demo Site**: Sitio de demostración para el widget (puerto 8091)
- **API Gateway**: Traefik como puerta de enlace (puerto 80)

## Configuración de Acceso

| Componente      | URL Directa              | URL a través del Gateway |
|-----------------|--------------------------|--------------------------|
| Backend API     | http://localhost:8080    | http://localhost/api     |
| Frontend Admin  | http://localhost:3001    | http://localhost/admin   |
| Widget API      | http://localhost:3002    | http://localhost/widget-api* |
| Widget Core     | http://localhost:3030    | http://localhost/widget  |
| Demo Site       | http://localhost:8091    | http://localhost/demo    |
| Traefik Dashboard | http://localhost:8080/dashboard | http://localhost/dashboard |

> *Nota: El Widget API debe accederse directamente a través del puerto 3002 debido a limitaciones de configuración en el gateway. Este comportamiento está configurado mediante una redirección.

## Configuración del API Gateway

El API Gateway está configurado con Traefik y proporciona:

- Enrutamiento a los diferentes servicios
- Redirecciones para servicios front-end
- Soporte para CORS
- Dashboard para monitoreo

### Detalles de Redirección

- Los endpoints `/admin` y `/demo` están configurados como redirecciones 307 para preservar la funcionalidad completa de las aplicaciones frontend
- El endpoint `/widget-api` redirige al puerto 3002 donde se encuentra la API del widget
- El endpoint `/widget` utiliza un stripPrefix para acceder al widget core

## Solución de Problemas

### Widget API

El Widget API (puerto 3002) debe accederse directamente debido a incompatibilidades con el gateway. Esto se debe a:

1. Configuración de red entre contenedores Docker
2. Requisitos específicos del servicio de API del widget
3. Redirección configurada en el gateway

### Aplicaciones Frontend

Las aplicaciones frontend (admin, demo) utilizan redirecciones en lugar de proxy inverso porque:

1. Evita problemas de recursos estáticos (CSS, JS) que comúnmente ocurren con stripPrefix
2. Mantiene la navegación y rutas internas intactas
3. Preserve la funcionalidad de desarrollador como Hot Module Replacement

## Inicio del Sistema

```bash
# En la carpeta GrowDesk
docker-compose up -d

# En la carpeta GrowDesk-Widget
docker-compose up -d
```

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

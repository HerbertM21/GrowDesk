# Arquitectura del Backend de GrowDesk

## Visión General

El backend de GrowDesk está desarrollado en Go (Golang), siguiendo una arquitectura modular y limpia que separa claramente las responsabilidades. Está diseñado para ofrecer una API RESTful que gestiona tickets de soporte, usuarios, categorías y FAQs.

## Estructura de Directorios

```
/backend
├── api/               # Definiciones de la API y documentación
├── cmd/               # Puntos de entrada de la aplicación
│   ├── server/        # Servidor principal de la aplicación
│   └── sync-server/   # Servidor de sincronización (componente auxiliar)
├── data/              # Datos persistentes y archivos de almacenamiento
├── internal/          # Código interno de la aplicación 
│   ├── data/          # Capa de acceso a datos y almacenamiento
│   ├── db/            # Configuración y conexiones de base de datos
│   ├── handlers/      # Manejadores HTTP para las rutas de la API
│   ├── middleware/    # Middleware HTTP (auth, logging, etc.)
│   ├── models/        # Definición de modelos y entidades
│   ├── service/       # Lógica de negocio
│   ├── utils/         # Utilidades y funciones auxiliares
│   └── websocket/     # Gestión de conexiones WebSocket
├── temp/              # Archivos temporales
├── Dockerfile         # Configuración para construir la imagen Docker
├── Dockerfile.sync    # Configuración para el servidor de sincronización
├── go.mod             # Dependencias de Go
└── Makefile           # Comandos para compilar y gestionar el proyecto
```

## Componentes Principales

### 1. Servidor Principal (cmd/server)

Es el punto de entrada de la aplicación que configura e inicia:
- El router HTTP
- La conexión a la base de datos
- Los servicios de autenticación
- Los manejadores de API

### 2. Manejadores (internal/handlers)

Procesan las solicitudes HTTP y gestionan las respuestas:
- `auth.go`: Gestiona autenticación y usuarios
- `tickets.go`: Operaciones CRUD para tickets
- `categories.go`: Gestión de categorías
- `faqs.go`: Gestión de FAQs (preguntas frecuentes)

Cada manejador sigue un patrón similar:
1. Recibe una solicitud HTTP
2. Valida los datos de entrada
3. Interactúa con la capa de acceso a datos
4. Devuelve la respuesta apropiada

### 3. Modelos (internal/models)

Define las estructuras de datos principales:
- Ticket: Representa un ticket de soporte con su estado, mensajes, etc.
- User: Información de usuarios del sistema
- Category: Categorías para clasificar tickets
- FAQ: Preguntas frecuentes para ayuda

### 4. Capa de Datos (internal/data)

Proporciona una abstracción para acceder a los datos, implementando el patrón Repository:
- Admite almacenamiento en PostgreSQL
- Ofrece operaciones CRUD para todas las entidades
- Maneja transacciones y lógica de persistencia

### 5. WebSockets (internal/websocket)

Gestiona conexiones en tiempo real para:
- Actualizaciones de tickets
- Notificaciones
- Chat en tiempo real con agentes

## Flujos de Trabajo Principales

### 1. Creación de Tickets

1. Cliente envía solicitud POST a `/api/tickets`
2. `ticketHandler.CreateTicket` valida los datos
3. Genera un ID único para el ticket
4. Almacena el ticket en la base de datos
5. Devuelve el ticket creado con su ID

### 2. Widget Integration

La integración con el widget externo funciona así:
1. Widget envía solicitudes al endpoint `/widget/tickets`
2. `ticketHandler.CreateWidgetTicket` procesa estas solicitudes
3. Crea un ticket en el sistema principal
4. Establece una conexión WebSocket para chat en tiempo real

### 3. Autenticación

1. Usuario envía credenciales a `/api/auth/login`
2. `authHandler` valida las credenciales contra la base de datos
3. Genera un token JWT con la información del usuario y sus permisos
4. El cliente usa este token en el header `Authorization` para solicitudes futuras

## Servicios Auxiliares

### Servidor de Sincronización (cmd/sync-server)

Un componente independiente que:
- Sincroniza datos entre sistemas externos
- Importa usuarios o tickets desde otras plataformas
- Ejecuta procesos periódicos de mantenimiento

## Tecnologías Utilizadas

- **Go**: Lenguaje principal de desarrollo
- **PostgreSQL**: Base de datos principal
- **Gorilla/Mux**: Router HTTP
- **JWT**: Para autenticación
- **WebSockets**: Para comunicación en tiempo real
- **Docker**: Para contenerización y despliegue

## Patrones de Diseño

1. **Repository Pattern**: Separa la lógica de negocio del acceso a datos
2. **Middleware Pattern**: Para funcionalidades transversales como autenticación y logging
3. **Dependency Injection**: Para facilitar pruebas y desacoplamiento
4. **MVC-like**: Separación entre modelos, manejadores (controladores) y vistas (respuestas JSON) 
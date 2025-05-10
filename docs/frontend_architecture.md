# Arquitectura del Frontend de GrowDesk

## Visión General

El frontend de GrowDesk está desarrollado con Vue.js y TypeScript, utilizando un enfoque modular y reactivo. La aplicación sigue una arquitectura basada en componentes con gestión centralizada del estado mediante Pinia, facilitando una experiencia de usuario fluida y reactiva.

## Estructura de Directorios

```
/frontend
├── public/            # Archivos estáticos y favicon
├── src/               # Código fuente de la aplicación
│   ├── api/           # Clientes y configuración de API
│   ├── assets/        # Recursos estáticos (imágenes, fuentes, CSS)
│   ├── components/    # Componentes reutilizables
│   ├── router/        # Configuración de rutas y navegación
│   ├── services/      # Servicios y lógica compartida
│   ├── stores/        # Gestión de estado con Pinia
│   ├── types/         # Definiciones de tipos TypeScript
│   ├── utils/         # Utilidades y funciones auxiliares
│   ├── views/         # Componentes de página/vista
│   │   ├── Admin/     # Vistas de administración
│   │   ├── Profile/   # Vistas de perfil de usuario
│   │   ├── auth/      # Vistas de autenticación
│   │   ├── dashboard/ # Vistas de dashboard
│   │   └── tickets/   # Vistas relacionadas con tickets
│   ├── App.vue        # Componente raíz
│   └── main.ts        # Punto de entrada de la aplicación
├── .env               # Variables de entorno para desarrollo
├── Dockerfile         # Configuración para construir la imagen Docker
├── package.json       # Dependencias y scripts
└── vite.config.ts     # Configuración de Vite 
```

## Componentes Principales

### 1. Sistema de Rutas (src/router)

La navegación se gestiona mediante Vue Router:
- Rutas definidas con rutas anidadas y carga diferida (lazy loading)
- Protección de rutas basada en autenticación y roles
- Transiciones de página fluidas

El archivo principal, `index.ts`, establece todas las rutas y los middleware de protección:
- Comprueba si el usuario está autenticado
- Verifica permisos y roles para acceder a ciertas rutas
- Redirige según el estado de autenticación

### 2. Gestión de Estado (src/stores)

La aplicación utiliza Pinia para gestión de estado centralizada:
- `auth.ts`: Estado de autenticación y usuario actual
- `tickets.ts`: Estado y operaciones relacionadas con tickets
- `users.ts`: Gestión de usuarios y permisos
- `categories.ts`: Categorías de tickets
- `notifications.ts`: Notificaciones del sistema
- `faqs.ts`: Preguntas frecuentes
- `activity.ts`: Registro de actividades
- `chat.ts`: Estado de las conversaciones de chat

Cada store sigue un patrón similar:
1. Define un estado inicial
2. Proporciona getters para consultar el estado
3. Implementa acciones para modificar el estado
4. Se conecta con la API para operaciones asíncronas

### 3. Comunicación con API (src/api)

El cliente API está configurado en `api/client.ts` usando Axios:
- Interceptores para manejo de tokens
- Gestión global de errores
- Integración con el estado de autenticación
- Retry automático para operaciones críticas

### 4. Vistas Principales (src/views)

#### Dashboard (views/dashboard)
- Muestra métricas y KPIs principales
- Gráficos de tickets por categoría/estado
- Actividad reciente

#### Tickets (views/tickets)
- `TicketList.vue`: Lista filtrable de tickets
- `TicketDetail.vue`: Detalles y conversación de un ticket
- `KanbanBoard.vue`: Vista de tickets en formato Kanban

#### Admin (views/Admin)
- Gestión de usuarios y permisos
- Configuración del sistema
- Categorías y FAQs

### 5. Componentes Reutilizables (src/components)

Widgets y elementos de UI reutilizables:
- Barras de navegación
- Formularios
- Tarjetas de tickets
- Selectores
- Componentes de chat

## Flujos de Trabajo Principales

### 1. Autenticación

1. Usuario introduce credenciales en `Login.vue`
2. `auth.ts` store ejecuta la acción `login`
3. Se realiza petición al backend y se recibe token JWT
4. El token se almacena y se actualiza el estado
5. Router redirige al dashboard o admin según el rol

### 2. Gestión de Tickets

1. Cliente navega a `TicketList.vue`
2. `tickets.ts` store carga tickets desde la API
3. Usuario puede filtrar/ordenar/buscar tickets
4. Al seleccionar un ticket, se navega a `TicketDetail.vue`
5. Se pueden añadir mensajes o actualizar el estado del ticket

### 3. Integración con Widget

- La interfaz de administración permite configurar el widget
- Se proporciona código de embedding para sitios web
- Los tickets creados desde el widget aparecen en el sistema

## Tecnologías Utilizadas

- **Vue.js 3**: Framework principal con Composition API
- **TypeScript**: Para tipado estático
- **Pinia**: Gestión de estado
- **Vue Router**: Navegación
- **Axios**: Cliente HTTP
- **PrimeVue**: Componentes UI
- **Chart.js**: Visualización de datos
- **Vite**: Bundler y servidor de desarrollo

## Patrones de Diseño

1. **Composition API**: Para lógica reutilizable y tipado fuerte
2. **Store Pattern**: Centralización de estado y lógica
3. **Component-Based Architecture**: Modularidad y reutilización
4. **Adapter Pattern**: En la capa de API para aislar cambios
5. **Responsive Design**: Adaptación a diferentes dispositivos 
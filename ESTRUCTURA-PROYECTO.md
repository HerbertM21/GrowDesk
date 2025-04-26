# Estructura del Proyecto GrowProyect

Este documento describe la estructura y organización del proyecto GrowProyect, detallando la funcionalidad de cada carpeta tanto en el backend como en el frontend.

## Visión General

GrowProyect está organizado en dos aplicaciones principales:

1. **GrowDesk**: Sistema principal de gestión de tickets de soporte
2. **GrowDesk-Widget**: Widget integrable para que los clientes puedan crear y gestionar tickets

## Backend (GrowDesk/backend)

### Estructura de Carpetas del Backend

```
backend/
├── cmd/                # Puntos de entrada de la aplicación
│   └── api/            # Código específico para la API
├── internal/           # Código privado de la aplicación
│   ├── auth/           # Autenticación y autorización
│   ├── config/         # Configuración de la aplicación
│   ├── controllers/    # Controladores de la lógica de negocio
│   ├── handlers/       # Manejadores HTTP
│   ├── middleware/     # Middleware para peticiones HTTP
│   ├── models/         # Definiciones de modelos de datos
│   ├── repository/     # Capa de acceso a datos
│   ├── routes/         # Definición de rutas de la API
│   ├── server/         # Configuración del servidor
│   └── service/        # Servicios de negocio
├── models/             # Modelos de dominio
│   ├── attachment.go   # Modelo para adjuntos
│   ├── message.go      # Modelo para mensajes
│   ├── ticket.go       # Modelo para tickets
│   ├── user.go         # Modelo para usuarios
│   ├── utils.go        # Utilidades para los modelos
│   └── widget.go       # Modelo para configuración del widget
├── pkg/                # Código que puede ser utilizado por aplicaciones externas
│   ├── auth/           # Utilidades de autenticación
│   ├── database/       # Utilidades para conexión a base de datos
│   ├── logger/         # Sistema de registro de eventos
│   ├── utils/          # Utilidades generales
│   └── websocket/      # Implementación de WebSockets
├── main.go             # Punto de entrada principal
├── go.mod              # Dependencias de Go
├── go.sum              # Checksums de dependencias de Go
├── mock-auth-server.js # Servidor de autenticación simulado para desarrollo
├── simple-mock-server.js # Servidor simulado simple para desarrollo
├── categories.json     # Datos de categorías para desarrollo
├── faqs.json           # Datos de preguntas frecuentes para desarrollo
├── tickets.json        # Datos de tickets para desarrollo
└── users.json          # Datos de usuarios para desarrollo
```

### Funcionalidad del Backend

- **cmd/api**: Contiene el punto de entrada para iniciar el servidor API.
  
- **internal**: Contiene la implementación principal del backend:
  - **auth**: Gestiona la autenticación y autorización de usuarios.
  - **config**: Maneja la configuración de la aplicación.
  - **controllers**: Implementa la lógica de negocio.
  - **handlers**: Gestiona las peticiones HTTP y respuestas.
  - **middleware**: Procesa peticiones antes de llegar a los handlers.
  - **models**: Define estructuras de datos internas.
  - **repository**: Implementa el acceso a datos.
  - **routes**: Define las rutas de la API.
  - **server**: Configura y gestiona el servidor HTTP.
  - **service**: Implementa servicios de negocio.

- **models**: Define las estructuras de datos principales para la aplicación:
  - **ticket.go**: Modelo para tickets de soporte.
  - **user.go**: Modelo para usuarios del sistema.
  - **message.go**: Modelo para mensajes en los tickets.
  - **attachment.go**: Modelo para archivos adjuntos.
  - **widget.go**: Modelo para la configuración del widget.

- **pkg**: Contiene paquetes que podrían ser utilizados por aplicaciones externas:
  - **auth**: Utilidades para autenticación.
  - **database**: Funciones para interactuar con la base de datos.
  - **logger**: Sistema de registro para la aplicación.
  - **utils**: Funciones de utilidad general.
  - **websocket**: Implementación para comunicación en tiempo real.

- **Archivos JSON**: Contienen datos de muestra para desarrollo.

- **Servidores Mock**: Proporcionan endpoints simulados para desarrollo y pruebas.

## Frontend (GrowDesk/frontend)

### Estructura de Carpetas del Frontend

```
frontend/
├── public/             # Archivos estáticos públicos
├── src/                # Código fuente
│   ├── api/            # Conexión con la API del backend
│   ├── assets/         # Recursos estáticos (imágenes, fuentes, etc.)
│   ├── components/     # Componentes reutilizables de Vue
│   ├── contexts/       # Contextos para gestión de estado
│   ├── hooks/          # Hooks personalizados
│   ├── pages/          # Componentes a nivel de página
│   ├── router/         # Configuración de rutas de la aplicación
│   ├── services/       # Servicios para lógica de negocio
│   ├── stores/         # Almacenamiento de estado (Pinia/Vuex)
│   ├── types/          # Definiciones de tipos TypeScript
│   ├── utils/          # Utilidades generales
│   ├── views/          # Vistas de la aplicación
│   ├── App.vue         # Componente raíz
│   ├── main.ts         # Punto de entrada de la aplicación
│   └── shims-vue.d.ts  # Declaraciones de tipos para Vue
├── index.html          # Plantilla HTML principal
├── tsconfig.json       # Configuración de TypeScript
├── vite.config.ts      # Configuración de Vite
├── Dockerfile          # Configuración para Docker
└── nginx.conf          # Configuración de Nginx para producción
```

### Funcionalidad del Frontend

- **src**: Contiene todo el código fuente de la aplicación:
  - **api**: Módulos para comunicarse con la API del backend.
  - **assets**: Recursos como imágenes, íconos y hojas de estilo.
  - **components**: Componentes Vue reutilizables.
  - **contexts**: Contextos para gestionar el estado en ciertos ámbitos.
  - **hooks**: Funcionalidades reutilizables para componentes.
  - **pages**: Componentes que representan páginas completas.
  - **router**: Configuración del enrutador de Vue, define las rutas de la aplicación.
  - **services**: Servicios para manejar lógica de negocio compleja.
  - **stores**: Almacenamiento global usando Pinia/Vuex para gestión de estado.
  - **types**: Definiciones de tipos TypeScript para mejorar el desarrollo.
  - **utils**: Funciones de utilidad general.
  - **views**: Componentes de vista que conforman la interfaz de usuario.

- **public**: Contiene archivos que serán servidos sin procesamiento:
  - **favicon.ico**: Ícono del sitio.

- **Archivos de configuración**:
  - **tsconfig.json**: Configuración de TypeScript.
  - **vite.config.ts**: Configuración del bundler Vite.
  - **nginx.conf**: Configuración de Nginx para despliegue.
  - **Dockerfile**: Instrucciones para crear la imagen Docker.

## Widget (GrowDesk-Widget)

### Estructura de Carpetas del Widget

```
GrowDesk-Widget/
├── demo-site/          # Sitio de demostración para el widget
├── widget-api/         # API específica para el widget
│   ├── data/           # Datos almacenados por el widget
│   └── ...             # Archivos de configuración y código
└── widget-src/         # Código fuente del widget
    ├── api/            # Cliente API para comunicación con widget-api
    ├── public/         # Archivos públicos
    ├── src/            # Código fuente principal del widget
    └── ...             # Archivos de configuración
```

### Funcionalidad del Widget

- **demo-site**: Sitio web de demostración para mostrar la integración del widget.
  
- **widget-api**: API dedicada que sirve como backend para el widget:
  - **data**: Almacena datos relacionados con los tickets creados desde el widget.

- **widget-src**: Código fuente del widget que se puede integrar en sitios web:
  - **api**: Cliente para comunicarse con widget-api.
  - **src**: Implementación del widget utilizando JavaScript/TypeScript.

## Arquitectura Global

GrowProyect implementa una arquitectura de tres partes:

1. **GrowDesk Backend**: API REST en Go que maneja la lógica de negocio y almacenamiento de datos.
   
2. **GrowDesk Frontend**: Aplicación Vue.js que proporciona la interfaz administrativa para gestionar tickets.
   
3. **GrowDesk Widget**: Componente embebible que permite a los usuarios finales crear y ver tickets desde cualquier sitio web.

Las tres partes funcionan en conjunto para proporcionar un sistema completo de gestión de soporte al cliente.
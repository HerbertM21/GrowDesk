# GrowDesk

## Guía de Instalación y Ejecución

Este documento proporciona las instrucciones detalladas para instalar y ejecutar el proyecto GrowDesk, un sistema completo de gestión de tickets de soporte con un widget embebible para sitios web.

## Requisitos Previos

Antes de comenzar, asegúrate de tener instalado lo siguiente:

- **Node.js**: versión 16.0.0 o superior
- **npm**: versión 8.0.0 o superior
- **Go**: versión 1.18 o superior
- **Git**: para clonar el repositorio

## Estructura del Proyecto

El proyecto está organizado en dos componentes principales:

- **GrowDesk**: Sistema principal de gestión de tickets
- **GrowDesk-Widget**: Widget embebible para sitios web de clientes

Para más detalles sobre la estructura del proyecto, consulta el archivo [ESTRUCTURA-PROYECTO.md](ESTRUCTURA-PROYECTO.md).

## Pasos para la Instalación

### 1. Configuración del Backend

```bash
# Navegar al directorio del backend
cd GrowDesk/backend

# Instalar dependencias de Node.js para los servidores mock
npm install

# Instalar dependencias de Go
go mod tidy
```

### 2. Configuración del Frontend

```bash
# Navegar al directorio del frontend
cd GrowDesk/frontend

# Instalar dependencias
npm install
```

### 3. Configuración del Widget

```bash
# Configuración de la API del Widget
cd GrowDesk-Widget/widget-api
go mod tidy

# Configuración del código fuente del Widget
cd ../widget-src
npm install

# Configuración del sitio de demostración
cd ../demo-site
npm install
```

## Ejecución del Proyecto

### Backend

```bash
# Desde el directorio GrowDesk/backend
# Para iniciar el servidor mock (datos de prueba)
node simple-mock-server.js

# En otra terminal, para iniciar el servidor de autenticación mock
node mock-auth-server.js
```

### Frontend

```bash
# Desde el directorio GrowDesk/frontend
npm run dev
```

El panel de administración estará disponible en: http://localhost:3000

### Widget API

```bash
# Desde el directorio GrowDesk-Widget/widget-api
go run main.go
```

La API del widget estará disponible en: http://localhost:8082

### Widget Demo Site

```bash
# Desde el directorio GrowDesk-Widget/demo-site
npm run dev
```

El sitio de demostración estará disponible en: http://localhost:8090

## Credenciales de Prueba

Para acceder al sistema de administración, usa las siguientes credenciales:

- **Administrador**:

  - Email: admin@example.com
  - Contraseña: password

## Configuración de Puertos

| Componente            | Puerto | Descripción                          |
| --------------------- | ------ | ------------------------------------- |
| Backend (Mock Server) | 8000   | Servidor mock para datos de prueba    |
| Backend (Auth Server) | 8001   | Servidor mock de autenticación       |
| Frontend              | 3000   | Panel de administración de GrowDesk  |
| Widget API            | 8082   | API para la comunicación del widget  |
| Demo Site             | 8090   | Sitio web de demostración con widget |

## Desarrollo del Widget

Si necesitas modificar y compilar el widget:

```bash
# Desde el directorio GrowDesk-Widget/widget-src
# Iniciar en modo desarrollo
npm run dev

# Compilar para producción
npm run build
```

Los archivos compilados se generarán en el directorio `dist` y deberán copiarse a la carpeta `demo-site/assets/` para probar los cambios.

## Licencia

Este proyecto está bajo la licencia MIT. Para más detalles, consulta el archivo LICENSE.

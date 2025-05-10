# API Gateway para GrowDesk V2

Este componente implementa un API Gateway utilizando Traefik para centralizar y gestionar todas las comunicaciones entre los diferentes servicios de GrowDesk V2.

## ¿Qué es Traefik?

Traefik es un proxy inverso y balanceador de carga moderno diseñado específicamente para microservicios. A diferencia de soluciones tradicionales como Nginx o HAProxy, Traefik ofrece descubrimiento dinámico de servicios, configuración automática de rutas y un potente sistema de middlewares.

## Arquitectura de GrowDesk con Traefik

```
   Usuario Final
        |
        v
        |
+-------+--------+
|                |
| Traefik API    |
| Gateway (80)   |
|                |
+----------------+
     |    |    |
     v    v    v
+-----+ +-----+ +------+
|     | |     | |      |
|Back-| |Front| |Widget|
|end  | |end  | |API   |
|(8080)| |(3000)| |(3000)|
|     | |     | |      |
+-----+ +-----+ +------+
   |               |
   v               v
+------+       +------+
|      |       |Widget|
|Postgre|      |Core  |
|SQL    |      |(3030)|
|      |       |      |
+------+       +------+
   |    \
   v     \
+------+  \    +-------+
|      |   --->|       |
|Redis |       |Sync   |
|      |       |Server |
+------+       |       |
               +-------+
```

## Flujo de Solicitudes

1. El cliente envía una solicitud a `http://localhost/`
2. Traefik recibe la solicitud y analiza la URL
3. Basándose en las reglas de enrutamiento, Traefik determina a qué servicio interno debe dirigir la solicitud
4. Se aplican middlewares según sea necesario (CORS, rate limiting, etc.)
5. La solicitud se reenvía al servicio correspondiente
6. El servicio procesa la solicitud y devuelve una respuesta
7. Traefik devuelve la respuesta al cliente

## Componentes Configurados

### Entrypoints

Los entrypoints son los puntos de entrada donde Traefik escucha las conexiones entrantes:

| Entrypoint | Puerto | Descripción |
|------------|--------|-------------|
| web        | 80     | Punto de entrada principal para todas las solicitudes HTTP |
| dashboard  | 8082   | Puerto dedicado para el dashboard de Traefik |
| traefik    | 8080   | Puerto interno para la API de administración de Traefik |

### Routers

Los routers definen las reglas para dirigir el tráfico a los servicios:

| Router | Regla | Servicio | Prioridad |
|--------|-------|----------|-----------|
| backend-api | `Host('localhost') && PathPrefix('/api') && !PathPrefix('/dashboard') && !PathPrefix('/api/version') && !PathPrefix('/api/entrypoints') && !PathPrefix('/api/overview') && !PathPrefix('/api/http')` | backend-service | 10 |
| widget-api | `PathPrefix('/widget-api')` | widget-api-service | 20 |
| admin-ui | `PathPrefix('/admin')` | frontend-service | 30 |
| widget-core | `PathPrefix('/widget')` | widget-core-service | 40 |
| demo-site | `PathPrefix('/demo')` | demo-site-service | 50 |
| traefik-api | `PathPrefix('/api/version') \|\| PathPrefix('/api/entrypoints') \|\| PathPrefix('/api/overview') \|\| PathPrefix('/dashboard') \|\| PathPrefix('/api/http')` | api@internal | 200 |

### Servicios

Los servicios definen cómo acceder a los backends:

| Servicio | URL |
|----------|-----|
| backend-service | http://backend:8080 |
| widget-api-service | http://widget-api:3000 |
| frontend-service | http://frontend:3000 |
| widget-core-service | http://widget-core:3030 |
| demo-site-service | http://demo-site:8090 |

### Middlewares

Los middlewares proporcionan funcionalidades adicionales:

| Middleware | Tipo | Descripción |
|------------|------|-------------|
| widget-api-stripprefix | stripPrefix | Elimina el prefijo `/widget-api` de las URLs |
| admin-ui-stripprefix | stripPrefix | Elimina el prefijo `/admin` de las URLs |
| widget-core-stripprefix | stripPrefix | Elimina el prefijo `/widget` de las URLs |
| demo-site-stripprefix | stripPrefix | Elimina el prefijo `/demo` de las URLs |
| cors-headers | cors | Configura los encabezados CORS para permitir solicitudes cross-origin |
| rate-limit-api | rateLimit | Limita el número de solicitudes a la API (100 req/s) |
| rate-limit-widget | rateLimit | Limita el número de solicitudes al widget (200 req/s) |

## Beneficios de la Arquitectura

1. **Simplificación de URLs**: Los clientes acceden a todos los servicios a través de un único dominio (`localhost`).
2. **Seguridad Mejorada**: Los servicios internos no están expuestos directamente.
3. **Escalabilidad**: Fácil adición de nuevos servicios sin cambiar la configuración del cliente.
4. **Observabilidad**: Dashboard centralizado para monitorear todos los servicios.
5. **Control de Tráfico**: Limitación de tasas y otras políticas aplicadas de manera centralizada.

## Accediendo a los Servicios

Con esta configuración, puedes acceder a todos los servicios a través de Traefik:

- Backend API: `http://localhost/api/*`
- Admin UI: `http://localhost/admin/*`
- Widget API: `http://localhost/widget-api/*`
- Widget Cliente: `http://localhost/widget/*`
- Demo Site: `http://localhost/demo/*`
- Traefik Dashboard: `http://localhost/dashboard/`

## Monitoreo y Administración

El dashboard de Traefik proporciona una interfaz visual para:

1. **Monitorear** el estado de los servicios y routers
2. **Visualizar** la configuración actual
3. **Depurar** problemas de enrutamiento
4. **Inspeccionar** el tráfico y las métricas

Accede al dashboard en `http://localhost/dashboard/`

## Descripción

El API Gateway funciona como punto de entrada único para todo el sistema GrowDesk V2, proporcionando:

- **Enrutamiento unificado**: Todas las APIs y UIs disponibles a través de una sola URL.
- **Seguridad**: Limita el acceso directo a los servicios internos.
- **Control de tráfico**: Implementa rate limiting para prevenir abusos.
- **Observabilidad**: Métricas y logs centralizados para monitoreo.
- **Gestión de CORS**: Configuración centralizada para encabezados CORS.

## Estructura de rutas

| Ruta               | Servicio              | Descripción                        |
|--------------------|----------------------|-----------------------------------|
| `/api/*`           | Backend API          | API principal para el panel admin  |
| `/widget-api/*`    | Widget API           | API para comunicación del widget   |
| `/admin/*`         | Frontend Admin       | Panel de administración            |
| `/widget/*`        | Widget Core          | Widget para sitios de clientes     |
| `/demo/*`          | Demo Site            | Sitio de demostración              |

## Configuración

La configuración del API Gateway se divide en varios archivos:

- **traefik.yml**: Configuración principal de Traefik.
- **dynamic_conf/api_routes.yml**: Definición de rutas y servicios.
- **dynamic_conf/cors.yml**: Configuración de CORS.
- **dynamic_conf/rate_limit.yml**: Límites de tasa para prevención de abusos.
- **dynamic_conf/metrics.yml**: Configuración de métricas y observabilidad.

## Uso

El API Gateway se inicia automáticamente con Docker Compose:

```bash
docker-compose up -d
```

Una vez iniciado, puedes acceder a los siguientes endpoints:

- **API Principal**: http://localhost/api
- **Panel de Admin**: http://localhost/admin
- **Widget**: http://localhost/widget
- **Demo Site**: http://localhost/demo
- **Dashboard de Traefik**: http://localhost/dashboard/

## Seguridad

Por defecto, el API Gateway está configurado para desarrollo local. En un entorno de producción, se recomienda:

1. Habilitar HTTPS/TLS
2. Configurar autenticación para el dashboard de Traefik
3. Revisar y ajustar los límites de tasa
4. Implementar un firewall de aplicación web (WAF)

## Extensión

Para añadir nuevos servicios al API Gateway:

1. Agregar el servicio al archivo `docker-compose.yml`
2. Definir las rutas en `dynamic_conf/api_routes.yml`
3. Actualizar la configuración CORS si es necesario
4. Aplicar rate limits apropiados

# Dockerización e Instalación de GrowDesk

Este documento contiene las instrucciones para configurar y ejecutar GrowDesk usando Docker.

## Requisitos previos

- Docker y Docker Compose instalados en tu máquina
- Git para clonar el repositorio

## Estructura del proyecto

El proyecto GrowDesk está dividido en dos partes principales:

1. **GrowDesk**: La aplicación principal con frontend y backend
2. **GrowDesk-Widget**: El widget para integración en otros sitios web

## Pasos para instalar y ejecutar GrowDesk

### 1. Clonar el repositorio

```bash
git clone https://github.com/tu-usuario/GrowDeskV2.git
cd GrowDeskV2
```

### 2. Configuración inicial

#### 2.1 Variables de entorno

Verifica los archivos de variables de entorno:

- `GrowDesk/frontend/.env`
- `GrowDesk/backend/.env`

#### 2.2 Permisos (en sistemas Unix/Linux)

```bash
chmod +x GrowDesk/frontend/entrypoint.sh
```

### 3. Iniciar la aplicación

#### Modo desarrollo

```bash
cd GrowDesk
docker-compose up -d
```

Este comando:
- Construye las imágenes necesarias con el entorno configurado para desarrollo
- Inicia el frontend, backend, PostgreSQL y Redis
- El frontend estará disponible en http://localhost:3000
- El backend estará disponible en http://localhost:8080

#### Modo producción

```bash
cd GrowDesk
docker-compose -f docker-compose.prod.yml up -d
```

En modo producción:
- Utiliza el mismo Dockerfile pero construido con el argumento NODE_ENV=production
- El frontend estará disponible en http://localhost:3000
- El backend seguirá disponible en http://localhost:8080

### 4. Verificar que los servicios están funcionando

```bash
docker-compose ps
```

### 5. Acceder a la aplicación

- Abre un navegador y visita http://localhost:3000
- Inicia sesión con las credenciales de prueba (en modo desarrollo)
  - Email: admin@example.com
  - Contraseña: password

## Compilación personalizada

Si necesitas personalizar la compilación, puedes pasar argumentos adicionales:

```bash
# Para desarrollo con host específico
docker-compose build --build-arg VITE_API_URL=http://miservidor:8080/api frontend

# Para producción con configuración personalizada
docker-compose -f docker-compose.prod.yml build --build-arg NODE_ENV=production frontend
```

## Iniciar GrowDesk-Widget (opcional)

Si necesitas ejecutar el widget de integración:

```bash
cd GrowDesk-Widget
docker-compose up -d
```

Esto iniciará:
- El API del widget en http://localhost:8082
- El sitio de demostración en http://localhost:8090
- El código fuente del widget en http://localhost:5174

## Detener los servicios

### Para detener GrowDesk

```bash
cd GrowDesk
docker-compose down
```

### Para detener GrowDesk y eliminar volúmenes (elimina datos persistentes)

```bash
cd GrowDesk
docker-compose down -v
```

## Solución de problemas comunes

### Error de conexión a la base de datos

Verifica que PostgreSQL esté corriendo:

```bash
docker logs growdesk-db
```

### El frontend no puede conectarse al backend

- Asegúrate de que las variables de entorno `VITE_API_URL` estén correctamente configuradas
- Verifica los logs del backend: `docker logs growdesk-backend`

### Problema de autenticación

Si tienes problemas al iniciar sesión, verifica los logs del backend y asegúrate de que las credenciales sean correctas.

### Error con el archivo app.ts

Si ves errores relacionados con el import de "./app" en auth.ts:

```
Failed to resolve import "./app" from "src/stores/auth.ts"
```

Necesitas crear el archivo `src/stores/app.ts` o eliminar la importación de ese archivo en `auth.ts`. 
# Solución: Servidor de Sincronización en Contenedor Docker

## Problema resuelto

Se ha implementado una solución para ejecutar el servidor de sincronización de GrowDesk como un contenedor Docker integrado con el resto de la aplicación. Esto resuelve los siguientes problemas:

1. Sincronización automática de datos entre el frontend y los archivos JSON en `/backend/data/`
2. Gestión correcta de los nuevos usuarios creados en la aplicación
3. Persistencia de datos entre reinicios del sistema

## Implementación

### 1. Dockerfile específico

Se ha creado un Dockerfile dedicado para el servidor de sincronización en:
```
/GrowDesk/backend/cmd/sync-server/Dockerfile
```

Este Dockerfile:
- Utiliza una compilación en dos etapas para minimizar el tamaño de la imagen
- Compila el código Go del servidor de sincronización
- Crea un contenedor liviano con Alpine Linux
- Configura el directorio de datos en `/app/data`

### 2. Configuración en docker-compose.yml

Se ha añadido la configuración del servicio sync-server al archivo docker-compose.yml principal:

```yaml
sync-server:
  build:
    context: ./GrowDesk/backend
    dockerfile: cmd/sync-server/Dockerfile
  container_name: growdesk-sync-server
  ports:
    - "8000:8000"
  volumes:
    - backend_data:/app/data
  environment:
    - PORT=8000
    - DATA_DIR=/app/data
  networks:
    - grow-network
  restart: unless-stopped
  depends_on:
    - backend
```

### 3. Volumen compartido

El volumen compartido `backend_data` permite que:
- El servidor de sincronización escriba en `/app/data`
- El backend acceda a los mismos archivos
- Los datos persistan entre reinicios

## Cómo usarlo

Para iniciar toda la aplicación con el servidor de sincronización:

```bash
cd /home/hmdev/Repositorios/GrowDeskV2
docker compose up -d
```

Para iniciar solo el servidor de sincronización:

```bash
cd /home/hmdev/Repositorios/GrowDeskV2
docker compose up -d sync-server
```

## Verificación

El servidor está funcionando correctamente si ves en los logs:

```
Servidor de sincronización iniciado en http://0.0.0.0:8000
Endpoint de sincronización: http://localhost:8000/api/sync/users
```

Ahora las actualizaciones desde el frontend se sincronizarán automáticamente con los archivos JSON en el directorio `data` del backend. 
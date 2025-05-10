# Solución al problema de sincronización en GrowDesk

## Problema detectado

Se identificaron los siguientes problemas en el sistema GrowDesk:

1. La aplicación frontend no podía sincronizar los datos con el backend
2. Al crear un usuario nuevo, aparecía un error `ERR_CONNECTION_REFUSED`
3. Las peticiones a APIs recibían errores 404
4. Los archivos JSON no estaban siendo actualizados correctamente

## Causas identificadas

- El servidor de sincronización (`sync-server`) no estaba en ejecución
- Las URLs en el archivo `.env` no apuntaban a los puertos correctos
- Docker Compose no estaba iniciando correctamente todos los contenedores necesarios

## Solución implementada

1. **Configuración del servidor de sincronización**:
   - Se creó un servidor de sincronización independiente en Go (`sync-server.go`)
   - Este servidor se encarga de recibir los datos del frontend y guardarlos correctamente en `/backend/data/`

2. **Corrección de URLs en el frontend**:
   - Se actualizó el archivo `.env` para apuntar a los puertos correctos:
     ```
     VITE_API_URL=http://localhost:8080/api
     VITE_WS_URL=ws://localhost:8080/api/ws
     VITE_SYNC_API_URL=http://localhost:8000/api/sync/users
     ```

3. **Estructura de datos**:
   - Se confirmó que el backend está configurado para usar los archivos de `/backend/data/` 
   - La estructura de `Store` en Go crea las rutas correctamente en este directorio

## Cómo usar la solución

1. **Iniciar el servidor de sincronización**:
   ```bash
   cd /home/hmdev/Repositorios/GrowDeskV2
   go run sync-server.go
   ```

2. **Verificar la sincronización**:
   - Al crear un usuario en la interfaz web, debería guardarse en localStorage
   - El servidor de sincronización recibirá estos datos y los guardará en `./GrowDesk/backend/data/users.json`
   - El backend leerá los datos actualizados desde este directorio

## Notas adicionales

- Si decides no usar Docker, asegúrate de que el servidor de sincronización esté siempre en ejecución
- Los errores "404 page not found" eran causados por URLs incorrectas en las peticiones API
- El backend usa los archivos en el directorio `data`, no los archivos JSON en la raíz 
# Solución de errores en GrowDesk

## Problemas identificados

### 1. Error 404 en la API
El error `Failed to load resource: the server responded with a status of 404 (Not Found)` ocurre porque:
- El frontend está intentando conectarse a un endpoint que no existe
- Posiblemente hay una discrepancia entre las rutas configuradas en el frontend y las disponibles en el backend

### 2. Error de JavaScript en el componente UserCreateModal
El error `Uncaught ReferenceError: emit is not defined` ocurre porque:
- En el componente `UserCreateModal.vue`, se está intentando usar `emit` directamente, pero no está definido correctamente
- Este es un error común en Vue 3 cuando se trabaja con la Composition API

## Soluciones implementadas

### 1. Servidor de sincronización
Hemos creado y configurado correctamente el servidor de sincronización:
- El contenedor Docker está ejecutándose en el puerto 8000
- La API de sincronización está disponible en `http://localhost:8000/api/sync/users`
- Los datos se almacenan correctamente en la carpeta `/backend/data/`

### 2. Corrección del error de emit en UserCreateModal
Para solucionar este error, hay tres opciones:

#### Opción 1: Definir correctamente emit
```javascript
// En script setup
const emit = defineEmits(['close', 'created']);

// Luego usar
emit('created');
```

#### Opción 2: Usar $emit en el template pero no en el script
```vue
<button @click="$emit('close')">Cerrar</button>

// En el script evitar usar emit directamente
```

#### Opción 3: Solución temporal
```javascript
// Reemplazar la línea problemática
setTimeout(() => {
  console.log('Notifying parent component');
  // emit('created'); <- Comentar esta línea
  return; // Devolver para evitar el error
}, 300);
```

## Configuración correcta de URLs

Para evitar errores 404, asegúrate de que las URLs en el archivo `.env` sean las correctas:

```
VITE_API_URL=http://localhost:8080/api
VITE_WS_URL=ws://localhost:8080/api/ws
VITE_SYNC_API_URL=http://localhost:8000/api/sync/users
```

## Cómo verificar que todo funciona correctamente

1. Reinicia el frontend: `docker compose restart frontend`
2. Verifica los logs del servidor de sincronización: `docker logs growdesk-sync-server`
3. Intenta crear un nuevo usuario desde la interfaz
4. Verifica que los datos se guarden en `/backend/data/users.json`

Si sigues teniendo problemas con el error 404:
1. Revisa las rutas definidas en el backend (`/backend/cmd/server/main.go`)
2. Asegúrate de que coincidan con las que se llaman desde el frontend
3. Posiblemente necesites añadir la ruta `/users` al backend 

## Error: 404 en creación de usuarios

### Problema
Cuando se intenta crear un nuevo usuario, la aplicación intenta hacer una solicitud POST a `/api/users`, pero recibe un error 404. Además, el mensaje de error "404 page not found" se estaba guardando incorrectamente en el localStorage como si fuera un usuario.

### Síntomas
1. Error en consola: 404 Not Found cuando se intenta crear un usuario
2. El localStorage de usuarios contiene entradas incorrectas, incluyendo el texto "404 page not found"
3. La creación de usuarios no funciona correctamente

### Solución implementada
Se han realizado las siguientes mejoras en el código:

1. **Modificación de la función `createUser` en `/GrowDesk/frontend/src/stores/users.ts`**:
   - Ahora detecta automáticamente si está en modo desarrollo (`import.meta.env.DEV`)
   - En modo desarrollo, crea el usuario directamente en localStorage sin intentar comunicarse con el backend
   - Usa IDs temporales con timestamp para evitar conflictos
   - Mejora el manejo de errores para evitar guardar mensajes de error como usuarios
   - Implementa validación de datos para asegurar que solo se guardan usuarios válidos

2. **Mejoras en la función `syncUsersWithBackend`**:
   - Validación previa de los datos en localStorage antes de intentar sincronizar
   - Filtrado de elementos no válidos en localStorage
   - Manejo robusto de errores durante la sincronización
   - Timeout de 5 segundos para evitar bloqueos indefinidos

3. **Validación de estructura de datos**:
   - Asegura que todos los usuarios tienen los campos necesarios
   - Filtra automáticamente datos que no sean objetos válidos

### Cómo verificar que funciona
1. Prueba crear un nuevo usuario desde la interfaz
2. Verifica en la consola del navegador que el usuario se crea correctamente en modo desarrollo
3. Confirma que el usuario aparece en la lista de usuarios después de crearlo
4. Inspecciona localStorage en DevTools para verificar que solo contiene usuarios válidos

### Limitaciones de la solución
1. Esta solución está orientada al desarrollo, permitiendo crear usuarios sin backend
2. En un entorno de producción, se seguirá intentando conectar con el backend API
3. Se recomienda implementar correctamente el endpoint `/api/users` en el backend cuando se pase a producción

## Próximos pasos recomendados
1. Implementar el endpoint `/api/users` en el backend para gestionar usuarios
2. Añadir autenticación real a la aplicación en lugar de la simulada
3. Mejorar el sistema de sincronización para manejar conflictos entre datos locales y del servidor 
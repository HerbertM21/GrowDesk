# Instrucciones de Acceso GrowDesk

## Acceso de administrador

### Credenciales por defecto:
- **Email:** admin@example.com
- **Password:** password

## Almacenamiento y sincronización de datos

GrowDesk almacena los datos en varias ubicaciones para asegurar una experiencia fluida:

### 1. Frontend (localStorage)
- Los datos iniciales se almacenan en el navegador usando localStorage
- La información de usuarios, tickets y categorías se guarda localmente
- Incluye limpieza automática de datos inválidos para evitar corrupción

### 2. Backend (archivos JSON)
- La información sincronizada se guarda en `/backend/data/`
- Los archivos importantes son:
  - `users.json` - Información de usuarios
  - `tickets.json` - Tickets y mensajes
  - `categories.json` - Categorías de soporte
  - `faqs.json` - Preguntas frecuentes

### Flujo de sincronización
1. Cuando creas un usuario nuevo en la interfaz de GrowDesk, primero se guarda en localStorage
2. El servidor de sincronización (`sync-server`) copia los datos a `/backend/data/`
3. Si usas el comando `start-sync-server.sh`, la sincronización está habilitada automáticamente

### Posibles problemas y soluciones
- Si los cambios no aparecen en la aplicación, verifica que el servidor de sincronización esté activo
- Para forzar la limpieza de localStorage, abre el Inspector del navegador -> Application -> Storage -> Clear Site Data
- El backend siempre usa los archivos en `/backend/data/`, no los de la raíz

## Ejecutar la aplicación con sincronización
```bash
# Desde la raíz del proyecto
cd GrowDesk
./start-sync-server.sh
```

Este comando inicia todos los servicios necesarios en contenedores Docker:
- Frontend (http://localhost:3000)
- API Backend (http://localhost:8080)
- Servidor de sincronización (http://localhost:8000)

## Estructura de directorio
```
GrowDesk/
├── frontend/        # Aplicación Vue.js
├── backend/         # API en Go
│   ├── data/        # Archivos de datos sincronizados (IMPORTANTE)
│   └── ...
└── start-sync-server.sh  # Script para iniciar todos los servicios
``` 
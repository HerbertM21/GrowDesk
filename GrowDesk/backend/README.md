# GrowDesk Backend - Servidor de Sincronización

Este servidor proporciona una API simple para sincronizar datos entre el frontend (almacenados en localStorage) y archivos JSON en el servidor.

## Funcionalidades

- Sincronización de usuarios entre localStorage y users.json

## Requisitos

- Go 1.18 o superior

## Instalación

1. Clona este repositorio
2. Navega al directorio del backend
3. Ejecuta `go mod tidy` para instalar dependencias

## Uso

### Iniciar el servidor de sincronización

```bash
cd cmd/sync-server
go run main.go
```

El servidor se ejecutará en http://localhost:8000 por defecto.

### Endpoints disponibles

- `POST /api/sync/users` - Sincroniza los usuarios desde localStorage al archivo users.json

## Cómo funciona la sincronización

1. El frontend almacena datos de usuarios en localStorage
2. Cuando se realiza una operación CRUD en usuarios (crear, actualizar, eliminar), los datos se envían al servidor
3. El servidor guarda estos datos en el archivo users.json
4. Al reiniciar, el backend lee los datos del archivo para mantener la persistencia

## Configuración

Puedes configurar el puerto usando la variable de entorno `PORT`:

```bash
PORT=9000 go run main.go
```

## Solución de problemas

Si encuentras problemas con la sincronización:

1. Verifica que el servidor esté ejecutándose y accesible
2. Revisa la consola del navegador para errores
3. Comprueba los logs del servidor 
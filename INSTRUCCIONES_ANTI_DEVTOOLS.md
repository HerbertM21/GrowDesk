# Protección contra Inspección de Código en GrowDesk

Este documento describe las medidas implementadas para proteger el código fuente de GrowDesk contra la inspección no autorizada mediante las herramientas de desarrollo del navegador.

## Medidas Implementadas

### 1. Ofuscación del Código (Production Build)

Se ha configurado Vite para ofuscar el código JavaScript durante el proceso de compilación:

```typescript
// En vite.config.ts
build: {
  minify: 'terser',
  terserOptions: {
    compress: {
      drop_console: true,
      drop_debugger: true,
      pure_funcs: ['console.log', 'console.info', 'console.debug']
    },
    format: {
      comments: false
    },
    mangle: {
      toplevel: true,
      safari10: true
    },
    keep_classnames: false,
    keep_fnames: false
  },
  rollupOptions: {
    output: {
      manualChunks: {
        'vendor': ['vue', 'vue-router', 'pinia'],
        'primevue': ['primevue', 'primeicons', 'primeflex'],
        'charts': ['chart.js', 'vue-chartjs']
      },
      entryFileNames: 'assets/[name]-[hash].js',
      chunkFileNames: 'assets/[name]-[hash].js',
      assetFileNames: 'assets/[name]-[hash].[ext]'
    }
  },
  sourcemap: false,
}
```

Esta configuración:
- Elimina todos los `console.log` y comentarios
- Ofusca nombres de variables y funciones
- Divide el código en chunks con nombres aleatorios
- Elimina los sourcemaps para que no se pueda reconstruir el código original

### 2. Bloqueo de DevTools - Script Independiente

Se ha creado el archivo `/public/anti-devtools.js` que:
- Detecta cuando las herramientas de desarrollo están abiertas mediante múltiples técnicas
- Bloquea la interfaz de usuario mostrando un mensaje de advertencia
- Permite restaurar la aplicación cuando se cierran las herramientas
- **Bypass para administradores**: Los usuarios con rol "admin" pueden utilizar las herramientas de desarrollo sin restricciones

El diseño de la pantalla de bloqueo incluye:
- Fondo con degradado en tonos morados (de `#5e35b1` a `#7b1fa2`)
- Contenedor con efecto de vidrio esmerilado (glassmorphism)
- Icono SVG de advertencia
- Botón redondeado con efectos de hover
- Texto elegante con sombras suaves
- Mensaje informativo sobre el acceso para administradores

### 3. Prevención de Accesos Directos

En el archivo `index.html` se han añadido scripts para:
- Bloquear teclas de acceso directo a DevTools (F12, Ctrl+Shift+I, etc.)
- Deshabilitar el menú contextual (clic derecho)

Estas restricciones no se aplican a los usuarios con rol de administrador.

## Cómo Funciona

El sistema de protección utiliza varias técnicas para detectar las herramientas de desarrollo:

1. **Diferencias de tamaño de ventana**: Detecta la diferencia entre `window.outerWidth/Height` y `window.innerWidth/Height`
2. **Punto de interrupción del depurador**: Utiliza la instrucción `debugger` para detectar si está activo un depurador
3. **Comprobación de objetos específicos**: Detecta objetos específicos de las herramientas de desarrollo
4. **Medición de tiempos de ejecución**: Comprueba el tiempo de ejecución de operaciones de consola

### Bypass para Administradores

El sistema ahora incluye un mecanismo de verificación de permisos:

1. Cuando se detectan herramientas de desarrollo, el script verifica si el usuario actual tiene rol de administrador:
   - Consulta el localStorage para obtener el ID del usuario actual
   - Busca al usuario en el array de usuarios almacenado en localStorage
   - Verifica si el rol del usuario es 'admin'

2. Si el usuario es administrador:
   - No se bloquea el acceso a las herramientas de desarrollo
   - Se permite el uso de teclas de acceso directo (F12, Ctrl+Shift+I)
   - Se permite el uso del menú contextual (clic derecho)
   - Se muestra un mensaje en la consola indicando que el modo desarrollador está activo

3. Si el usuario no es administrador:
   - Se muestra el bloqueo normal con un mensaje adicional
   - Se indica que los administradores pueden acceder a las herramientas

## Limitaciones

Es importante tener en cuenta que estas medidas dificultan, pero no hacen imposible, la inspección del código. Un desarrollador experimentado con suficiente determinación siempre podrá encontrar formas de eludir estas protecciones. Estas medidas están diseñadas para:

1. Disuadir a usuarios casuales de inspeccionar el código
2. Proteger la propiedad intelectual de usuarios no técnicos
3. Ocultar detalles de implementación sensibles

## Cómo Aplicar estos Cambios en Producción

Para aplicar estos cambios en el entorno de contenedores:

1. Los archivos modificados ya están en su lugar en el repositorio
2. Al construir la imagen Docker, estos archivos se incluirán automáticamente
3. Para reconstruir la aplicación con estas protecciones, ejecute:

```bash
docker-compose exec frontend npm run build
```

O reconstruya todo el contenedor:

```bash
docker-compose build frontend
docker-compose up -d
```

## Pruebas

Para verificar que la protección funciona correctamente:

1. Abra la aplicación en el navegador como usuario normal
   - Intente abrir las herramientas de desarrollo (F12 o clic derecho + Inspeccionar)
   - Debería ver un mensaje de bloqueo con fondo morado y un botón "Continuar"
   - La aplicación se detendrá hasta que cierre las herramientas de desarrollo

2. Inicie sesión como administrador (usuario con rol 'admin')
   - Intente abrir las herramientas de desarrollo (F12 o clic derecho + Inspeccionar)
   - Las herramientas deberían abrirse normalmente sin bloqueo
   - Debería ver un mensaje en la consola indicando que el modo desarrollador está activado

## Consideraciones Finales

Estas medidas de protección deben utilizarse junto con otras buenas prácticas de seguridad:

1. Implementar autenticación y autorización robustas en el backend
2. No confiar nunca en la seguridad del lado del cliente para datos sensibles
3. Realizar validaciones de datos tanto en el cliente como en el servidor
4. Mantener actualizado el software y sus dependencias 
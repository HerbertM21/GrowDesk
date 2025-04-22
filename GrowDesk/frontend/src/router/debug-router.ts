// Este archivo es solo para depuración
// eslint-disable-next-line @typescript-eslint/no-var-requires
const router = require('./index');

console.log('Router importado correctamente:', router);

// No hace nada más, solo verifica que la importación funcione
export default function debugRouter() {
  return 'Router importado correctamente';
} 
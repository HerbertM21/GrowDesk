#!/bin/sh

echo "Configurando variables de entorno..."

# Verificar si node_modules existe y tiene contenido
if [ ! -d "node_modules" ] || [ -z "$(ls -A node_modules 2>/dev/null)" ]; then
  echo "Directorio node_modules vacío o inexistente, instalando dependencias..."
  npm install
  if [ $? -ne 0 ]; then
    echo "Error al instalar las dependencias. Abortando."
    exit 1
  fi
fi

# Verificar específicamente si vite está instalado
if [ ! -f "node_modules/.bin/vite" ]; then
  echo "Vite no encontrado, instalando vite específicamente..."
  npm install --save-dev vite
  if [ $? -ne 0 ]; then
    echo "Error al instalar vite. Abortando."
    exit 1
  fi
fi

# Verificar el entorno
if [ "$NODE_ENV" = "production" ]; then
  echo "Iniciando en modo PRODUCCIÓN"
  
  # Instalamos Nginx si no está instalado
  apk add --no-cache nginx
  
  # Copiamos la configuración de Nginx
  if [ -f "nginx.conf" ]; then
    cp nginx.conf /etc/nginx/http.d/default.conf
  else
    echo "ADVERTENCIA: No se encontró nginx.conf, usando configuración por defecto"
    cat > /etc/nginx/http.d/default.conf << 'EOF'
server {
    listen 80;
    server_name _;
    root /app/dist;
    index index.html;
    
    location / {
        try_files $uri $uri/ /index.html;
    }
    
    location /api {
        proxy_pass ${VITE_API_URL};
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
EOF
  fi
  
  # Iniciamos Nginx
  echo "Iniciando Nginx..."
  exec nginx -g 'daemon off;'
else
  # En desarrollo, ejecutamos el servidor de Vite
  echo "Iniciando en modo DESARROLLO"
  echo "Iniciando servidor de desarrollo Vite..."
  # Usar la ruta absoluta a vite para asegurarse de que lo encuentre
  exec ./node_modules/.bin/vite --host 0.0.0.0 --port 3000
fi 
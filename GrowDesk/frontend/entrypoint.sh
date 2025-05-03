#!/bin/sh

echo "Configurando variables de entorno..."

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
  exec npm run dev -- --host 0.0.0.0 --port 3000
fi 
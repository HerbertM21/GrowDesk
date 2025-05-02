#!/bin/bash

cd /home/hmdev/Repositorios/GrowDeskV2/GrowDesk-Widget

# Copiar el archivo modificado al contenedor
echo "Copiando ChatWidget.vue actualizado al contenedor..."
docker cp widget-core/src/components/ChatWidget.vue growdesk-widget-widget-core-1:/app/src/components/ChatWidget.vue

# compilación 
echo "Compilando dentro del contenedor..."
docker exec growdesk-widget-widget-core-1 npm run build

echo "Esperando 3 segundos..."
sleep 3

echo "Verificando archivos generados..."
docker exec growdesk-widget-widget-core-1 ls -la /app/dist

# Copiar los archivos al demo-site
echo "Copiando archivos al demo-site..."
docker cp growdesk-widget-widget-core-1:/app/dist/style.css /tmp/style.css 
docker cp /tmp/style.css growdesk-widget-demo-site-1:/app/style.css 
docker cp growdesk-widget-widget-core-1:/app/dist/growdesk-widget.umd.js /tmp/growdesk-widget.umd.js 
docker cp /tmp/growdesk-widget.umd.js growdesk-widget-demo-site-1:/app/growdesk-widget.umd.js 
docker cp growdesk-widget-widget-core-1:/app/dist/growdesk-widget.umd.js.map /tmp/growdesk-widget.umd.js.map 
docker cp /tmp/growdesk-widget.umd.js.map growdesk-widget-demo-site-1:/app/growdesk-widget.umd.js.map 

# Reiniciar demo-site
echo "Reiniciando demo-site..."
docker restart growdesk-widget-demo-site-1

echo "¡Proceso completado! Ahora puedes acceder a http://localhost:8090" 
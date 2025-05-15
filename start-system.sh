#!/bin/bash

# Script para iniciar el sistema GrowDesk completo de forma controlada

echo "=== Iniciando GrowDesk V2 ==="

# Detener y eliminar los contenedores existentes si están en ejecución
echo "Limpiando contenedores existentes..."
docker stop growdesk-db growdesk-redis growdesk-backend growdesk-sync-server growdesk-frontend growdesk-widget-api growdesk-widget-core growdesk-demo-site 2>/dev/null || true
docker rm growdesk-db growdesk-redis growdesk-backend growdesk-sync-server growdesk-frontend growdesk-widget-api growdesk-widget-core growdesk-demo-site 2>/dev/null || true

# Limpiar volúmenes 
# echo "Limpiando volúmenes..."
# docker volume rm growdeskv2_postgres_data growdeskv2_backend_data growdeskv2_backend_uploads 2>/dev/null || true

# Iniciar servicios de GrowDesk
echo "Iniciando servicios de GrowDesk..."
cd GrowDesk
docker compose up -d postgres redis
echo "Esperando a que PostgreSQL esté listo..."
sleep 10

# Iniciar el backend
echo "Iniciando el backend de GrowDesk..."
docker compose up -d backend
echo "Esperando a que el backend esté listo..."
sleep 10

# Iniciar los servicios restantes de GrowDesk
echo "Iniciando el resto de servicios de GrowDesk..."
docker compose up -d
cd ..

# Iniciar servicios de GrowDesk-Widget
echo "Iniciando servicios de GrowDesk-Widget..."
cd GrowDesk-Widget
docker compose up -d
cd ..

echo "=== Sistema GrowDesk V2 iniciado ==="
echo "Panel de administración: http://localhost:3001"
echo "API Backend: http://localhost:8080"
echo "Widget Demo: http://localhost:8090"
echo "Widget API: http://localhost:3000"
echo "Widget Core: http://localhost:3030"

# Mostrar logs del backend para diagnóstico
echo ""
echo "=== Logs del backend (Ctrl+C para salir) ==="
docker logs -f growdesk-backend 
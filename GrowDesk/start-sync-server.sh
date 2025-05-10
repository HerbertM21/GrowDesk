#!/bin/bash

# Colores para mensajes
GREEN='\033[0;32m'
BLUE='\033[0;34m'
RED='\033[0;31m'
YELLOW='\033[0;33m'
NC='\033[0m' # Sin color

echo -e "${BLUE}=== Iniciando GrowDesk con sincronización automática ===${NC}"
echo -e "${YELLOW}Este script ejecutará la aplicación con todos los servicios incluyendo la sincronización${NC}"

# Comprobar si Docker está instalado
if ! command -v docker &> /dev/null; then
    echo -e "${RED}Error: Docker no está instalado.${NC}"
    echo -e "Por favor, instala Docker antes de continuar."
    exit 1
fi

# Comprobar si docker-compose está instalado
if ! command -v docker compose &> /dev/null; then
    echo -e "${RED}Error: docker-compose no está instalado.${NC}"
    echo -e "Por favor, instala docker-compose antes de continuar."
    exit 1
fi

# Directorio actual
CURRENT_DIR="$(pwd)"
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

# Navegar al directorio del proyecto
cd "$SCRIPT_DIR"

echo -e "${BLUE}Construyendo las imágenes Docker...${NC}"
docker compose build

echo -e "${GREEN}Iniciando todos los servicios de GrowDesk...${NC}"
docker compose up -d

echo -e "${GREEN}=== GrowDesk iniciado correctamente con sincronización ===${NC}"
echo -e "Aplicación Frontend: ${BLUE}http://localhost:3000${NC}"
echo -e "API Backend: ${BLUE}http://localhost:8080/api${NC}"
echo -e "Servidor de sincronización: ${BLUE}http://localhost:8000${NC}"
echo -e "${YELLOW}Para detener todos los servicios, ejecute: docker-compose down${NC}"
echo -e "${YELLOW}Para ver los logs en tiempo real: docker compose logs -f${NC}"

# Volver al directorio original
cd "$CURRENT_DIR" 
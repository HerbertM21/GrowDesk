#!/bin/bash

# Script para iniciar el servidor de sincronización de GrowDesk

echo "Iniciando servidor de sincronización GrowDesk..."

# Verificar que Go está instalado
if ! command -v go &> /dev/null; then
    echo "Error: Go no está instalado. Por favor instala Go antes de continuar."
    exit 1
fi

go run main.go

echo "Servidor detenido." 
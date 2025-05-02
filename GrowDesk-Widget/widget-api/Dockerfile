FROM node:18-alpine

# Instalar Go
RUN apk add --no-cache go git

WORKDIR /app

# Copiar archivos del proyecto
COPY . .

# Instalar dependencias de Node.js
RUN npm install

# Construir el proyecto Go en segundo plano (por si es necesario)
RUN go version || true

# Exponer puerto
EXPOSE 3000

# Comando para ejecutar la aplicación - inicialmente ejecutaremos main.go con go run
CMD ["go", "run", "main.go"] 
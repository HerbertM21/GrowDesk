FROM node:18-alpine

WORKDIR /app

# Copiar archivos del proyecto
COPY package*.json ./

# Instalar dependencias
RUN npm install

# Copiar el resto del código
COPY . .

# Construir la aplicación
RUN npm run build

#  servir archivos estáticos
RUN npm install -g serve

# Exponer el puerto
EXPOSE 3000

# Iniciar el servidor estático sirviendo los archivos de dist
CMD ["serve", "-s", "dist", "-l", "3000"] 
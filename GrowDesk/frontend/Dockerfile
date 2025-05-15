# This is a placeholder file
# Use Dockerfile.dev for development environment
# Use Dockerfile.prod for production environment

# Stage 1: Base
FROM node:18-alpine AS base

# Set working directory
WORKDIR /app

# Install dependencies
RUN apk add --no-cache bash curl

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source files
COPY . .

# Make entrypoint script executable
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# For production, build the app
ARG NODE_ENV=development
ENV NODE_ENV=${NODE_ENV}

RUN if [ "$NODE_ENV" = "production" ] ; then npm run build ; fi

# Configurar para que Vite escuche en todas las interfaces
ENV HOST=0.0.0.0
ENV PORT=3000
ENV VITE_API_URL=http://localhost:8080/api

# Expose port
EXPOSE 3000

# Set entrypoint script
ENTRYPOINT ["/entrypoint.sh"] 
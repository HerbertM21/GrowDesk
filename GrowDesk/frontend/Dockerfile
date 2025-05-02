# This is a placeholder file
# Use Dockerfile.dev for development environment
# Use Dockerfile.prod for production environment

FROM node:18-alpine

# Set working directory
WORKDIR /app

# Build arguments - default to development
ARG NODE_ENV=development
ENV NODE_ENV=${NODE_ENV}

# Copy package files
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source files
COPY . .

# Copy entrypoint script
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

# For production, build the app
RUN if [ "$NODE_ENV" = "production" ] ; then npm run build ; fi

# Configurar para que Vite escuche en todas las interfaces
ENV HOST=0.0.0.0
ENV PORT=3000
ENV VITE_API_URL=http://localhost:8080/api

# Expose port
EXPOSE 3000

# Start development server
CMD ["npm", "run", "dev", "--", "--host", "0.0.0.0", "--port", "3000"] 
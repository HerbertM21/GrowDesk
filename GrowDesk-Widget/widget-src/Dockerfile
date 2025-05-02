FROM node:18-alpine

WORKDIR /app

# Copy package.json and package-lock.json
COPY package*.json ./

# Install dependencies
RUN npm install

# Copy source code
COPY . .

# Build the widget (if needed)
RUN npm run build

# Expose port
EXPOSE 8000

# Start development server or serve built files
CMD ["npm", "start"] 
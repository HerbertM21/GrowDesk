# Node.js stage for mock server
FROM node:18-alpine

WORKDIR /app

# Copy package.json and install dependencies
COPY package.json ./
RUN npm install

# Copy mock server files and data files
COPY *.json ./
COPY *.js ./

# Expose port
EXPOSE 8080

# Set environment variables
ENV PORT=8080

# Command to run
CMD ["node", "simple-mock-server.js"] 
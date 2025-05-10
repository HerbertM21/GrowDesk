FROM golang:1.20-alpine AS builder

WORKDIR /app

# Instalar git para ir a buscar dependencias
RUN apk add --no-cache git

# Copiar archivos de código fuente
COPY . /app/

# Compilar la aplicación
WORKDIR /app/cmd/sync-server
RUN go mod tidy && \
    CGO_ENABLED=0 GOOS=linux go build -o sync-server .

# Segunda etapa para la imagen final
FROM alpine:latest

# Agregar certificados y zona horaria
RUN apk --no-cache add tzdata ca-certificates

# Directorio de trabajo
WORKDIR /app

# Copiar el binario compilado
COPY --from=builder /app/cmd/sync-server/sync-server /app/sync-server

# Crear el directorio de datos
RUN mkdir -p /app/data

# Exponer el puerto
EXPOSE 8000

# Comando para ejecutar
CMD ["/app/sync-server"] 
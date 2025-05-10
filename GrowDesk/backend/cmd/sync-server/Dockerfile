FROM golang:1.20-alpine AS builder

WORKDIR /app

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

WORKDIR /app

COPY --from=builder /app/cmd/sync-server/sync-server /app/sync-server

RUN mkdir -p /app/data

EXPOSE 8000

# Comando para ejecutar
CMD ["/app/sync-server"] 
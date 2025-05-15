FROM golang:1.20-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add tzdata ca-certificates wget

WORKDIR /app

COPY --from=builder /app/server /app/server

COPY --from=builder /app/internal/db/schema.sql /app/internal/db/schema.sql

RUN mkdir -p /app/data

EXPOSE 8080

ENV PORT=8080
ENV DATA_DIR=/app/data

CMD ["/app/server"] 
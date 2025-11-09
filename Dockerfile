# Build stage
FROM golang:1.21-alpine AS builder

WORKDIR /app

# Copiar go.mod y go.sum
COPY go.mod go.sum ./

# Descargar dependencias
RUN go mod download

# Copiar código fuente
COPY . .

# Compilar aplicación
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o whatsapp-api ./cmd/server

# Runtime stage
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

WORKDIR /root/

# Copiar binario compilado
COPY --from=builder /app/whatsapp-api .

# Puerto
EXPOSE 8080

# Variables de entorno por defecto
ENV API_PORT=8080
ENV LOG_LEVEL=INFO

# Ejecutar aplicación
CMD ["./whatsapp-api"]


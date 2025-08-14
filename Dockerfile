# --- Etapa 1: Compilación (Builder) ---
# --- CORRECCIÓN AQUÍ ---
# Actualizamos la versión de Go para que coincida con tu go.mod (>= 1.24.6)
FROM golang:1.24-alpine AS builder

# Establecemos el directorio de trabajo dentro del contenedor.
WORKDIR /app

# Copiamos los archivos de dependencias primero.
COPY go.mod go.sum ./
# Descargamos todas las dependencias.
RUN go mod download

# Copiamos el resto del código fuente de la aplicación al contenedor.
COPY . .

# Compilamos la aplicación.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# --- Etapa 2: Final ---
# Usamos una imagen base mínima para la producción.
FROM alpine:latest

# Establecemos el directorio de trabajo.
WORKDIR /app

# Copiamos únicamente el binario compilado desde la etapa de 'builder'.
COPY --from=builder /app/main .

# Copiamos el archivo .env si existe.
# El comando no fallará si el archivo .env no está en el repositorio.
COPY .env ./.env

# Creamos la carpeta 'uploads' directamente en la imagen.
# El volumen de docker-compose se montará sobre esta carpeta sin problemas.
RUN mkdir -p /app/uploads

# Exponemos el puerto en el que la aplicación de Go se ejecuta.
EXPOSE 8080

# El comando que se ejecutará cuando el contenedor se inicie.
CMD ["/app/main"]

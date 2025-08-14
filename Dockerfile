# --- Etapa 1: Compilación (Builder) ---
# Usamos una imagen oficial de Go con Alpine Linux como base para compilar la aplicación.
# Alpine es una distribución ligera, lo que acelera la etapa de construcción.
FROM golang:1.22-alpine AS builder

# Establecemos el directorio de trabajo dentro del contenedor.
WORKDIR /app

# Copiamos los archivos de dependencias primero.
# Docker guardará esta capa en caché y solo la volverá a ejecutar si go.mod o go.sum cambian.
COPY go.mod go.sum ./
# Descargamos todas las dependencias.
RUN go mod download

# Copiamos el resto del código fuente de la aplicación al contenedor.
COPY . .

# Compilamos la aplicación.
# CGO_ENABLED=0 crea un binario estático que no depende de librerías C del sistema.
# GOOS=linux asegura que el binario sea compatible con nuestro sistema operativo final (Alpine).
# -o /app/main especifica que el ejecutable se llamará 'main' y se guardará en /app.
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# --- Etapa 2: Final ---
# Usamos una imagen base mínima para la producción.
# Alpine es una excelente opción por su pequeño tamaño.
FROM alpine:latest

# Establecemos el directorio de trabajo.
WORKDIR /app

# Copiamos únicamente el binario compilado desde la etapa de 'builder'.
# Esto hace que nuestra imagen final sea muy pequeña y segura, ya que no contiene
# el código fuente ni las herramientas de compilación.
COPY --from=builder /app/main .

# Copiamos el directorio 'uploads' y el archivo .env.
# Aunque 'uploads' será un volumen, es una buena práctica que el directorio exista en la imagen.
COPY uploads ./uploads
COPY .env ./.env

# Exponemos el puerto en el que la aplicación de Go se ejecuta.
# Cambia el 8080 si tu aplicación usa un puerto diferente.
EXPOSE 8080

# El comando que se ejecutará cuando el contenedor se inicie.
# Simplemente ejecuta el binario de nuestra aplicación.
CMD ["/app/main"]

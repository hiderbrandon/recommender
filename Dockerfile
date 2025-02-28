# Usa una imagen base ligera de Go
FROM golang:1.24 as builder

# Configurar el directorio de trabajo
WORKDIR /app

# Copiar los archivos del proyecto al contenedor
COPY go.mod go.sum ./
RUN go mod download

# Copiar el resto del código fuente
COPY . .

# Compilar el binario de la aplicación
RUN go build -o main cmd/main.go

# Crear una imagen final más ligera
FROM gcr.io/distroless/base-debian12

# Configurar el directorio de trabajo
WORKDIR /app

# Copiar el binario compilado desde la fase anterior
COPY --from=builder /app/main .

# Exponer el puerto en el que corre la API
EXPOSE 8080

# Ejecutar la aplicación
CMD ["/app/main"]

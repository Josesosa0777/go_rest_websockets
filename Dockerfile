# Variable para definir la versión de go a utilizar
ARG GO_VERSION=1.19.1

# Para hacer pull de la imagen, alpine es para generar un binario para compilar el archivo
FROM golang:${GO_VERSION}-alpine AS builder

# Configurar la variable de entorno, no hay proxy, por eso es direct
RUN go env -w GOPROXY=direct
# Necesitamos comandos git para las instalaciones
RUN apk add --no-cache git
# Se requieren certificados de seguridad para la instalación de la aplicacion
RUN apk --no-cache add ca-certificates && update-ca-certificates

# Directorio en el que se ejecutarán los comandos
WORKDIR /src

# Copiar al directorio local los archivos go.mod y go.sum
COPY ./go.mod ./go.sum ./
# Instalación de las dependecias de go.mod
RUN go mod download

# Copiar los directorios en el contenedor
COPY ./ ./

# Ejecutar el comando que construirá la aplicación, a veces go quiere usar compilador de c++, entonces para que no lo haga, CGO_ENABLED=0
# Es necesario especificar los flags que se requieren para que el ejecutable funcione en el siguiente contenedor que se encargará de ejecutarlo (-installsuffix 'static')
# El ejecutable de la imagen será /platzi-rest-ws
RUN CGO_ENABLED=0 go build \
    -installsuffix 'static' \
    -o /platzi-rest-ws .

# Otra imagen (scratch) encargada de ejecutar la aplicación con el servidor:
FROM scratch AS runner

# Especificar de donde se copia (builder), y se copian los certificados previamente descargados (/etc/ssl/certs/ca-certificates.crt) en /ect/ssl/certs/
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /ect/ssl/certs/

# Copiar archivo de entorno a la ruta pricipal
COPY .env ./
# Copiar el ejecutable que se generó anteriormente (/platzi-rest-ws) de la imagen del builder a la imagen del runner (en /platzi-rest-ws):
COPY --from=builder /platzi-rest-ws /platzi-rest-ws

# Exponer el Puerto 5050
EXPOSE 5050

# Definir el comando que se quiere ejecutar
ENTRYPOINT ["/platzi-rest-ws"]
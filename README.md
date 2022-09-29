Inicialmente, si se crea un nuevo poryecto, se requiere crear carpeta del proyecto Rest-websockets e inicializar
`go mod init platzi.com/go/rest-ws`
Sin embargo, si este proyecto se está clonando, ya debería traer la info.

# Instalar módulo para crear web tokens, ejecutar:
`go get github.com/golang-jwt/jwt`
# paquete para tener ruteador y websocket:
`go get github.com/gorilla/mux`
`go get github.com/gorilla/websocket`
# Para variables de entorno:
`go get github.com/joho/godotenv`
# instalar libreria para usar postgres 
`go get github.com/lib/pq`
# Instalar dependencia que permitirá retornar el id como texto:
`go get github.com/segmentio/ksuid`

# Ejecución
En carpeta database, ejecutar:
`docker build . -t platzi-ws-rest-db`
`docker run -p 54321:5432 platzi-ws-rest-db`

En carpeta principal, ejecutar:
`go run main.go`
Ir a:
http://localhost:5050/signup
Hacer post, por ejemplo en postman:
{
    "email": "josephsosa@gmail.com",
    "password": "mypassword"
}

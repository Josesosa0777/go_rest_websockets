Inicialmente, si se crea un nuevo poryecto, se requiere crear carpeta del proyecto Rest-websockets e inicializar

`go mod init platzi.com/go/rest-ws`

Sin embargo, si este proyecto se está clonando, ya debería traer la info.

# Para instalar módulo para crear web tokens, ejecutar:
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
# Instalar bcrypt:
`go get golang.org/x/crypto/bcrypt`

# Ejecución
En carpeta database, ejecutar:

`docker build . -t platzi-ws-rest-db`

`docker run -p 54321:5432 platzi-ws-rest-db`

En carpeta principal, ejecutar:

`go run main.go`

Otro caso con otro tag, en carpeta database, ejecutar:

`docker build . -t platzi-rs-ws-db`

`docker run -p 54321:5432 platzi-rs-ws-db`

En carpeta principal, ejecutar:

`go run main.g`

ejemplo de info en el archivo .env:
```
PORT=:5050
JWT_SECRET=secret
DATABASE_URL=postgres://postgres:postgres@localhost:54321/postgres?sslmode=disable
```

# -------------

# Para ejecutar con Docker, 
En el archivo .env se debe pasar algo como:

```
PORT=:5050
JWT_SECRET=secret
```
se puede hacer:
`docker build . -t platzi-rest-ws-app`

`docker run -p 5050:5050 platzi-rest-ws-app`

pero sólo funciona el GET en http://localhost:5050

# Por tanto, para que se conecte con la db usando Docker, ejecutar directamente:
`docker-compose up -d`

Ir a:

http://localhost:5050

**Con postman:**

http://localhost:5050/signup

Hacer post, por ejemplo:
```
{
    "email": "josephsosa@gmail.com",
    "password": "mypassword"
}
```

http://localhost:5050/login

Hacer post, por ejemplo en postman:
```
{
    "email": "josephsosa@gmail.com",
    "password": "mypassword"
}
```

Devuelve el token, ejemplo:
```
{
    "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRml2TmxBYmVOaW1WUk1MNFVVVTVnUzk5UTIiLCJleHAiOjE2NjUxNjA0ODR9.oMT9_vwseGUR3h_NbthFwFK02r1aJGfQ8kyFQZk81sM"
}
```

http://localhost:5050/api/v1/me

Hacer get pasando en el campo de Headers como Key Authorization, y como Value el Token:

Key: Authorization

Value: eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VySWQiOiIyRml2TmxBYmVOaW1WUk1MNFVVVTVnUzk5UTIiLCJleHAiOjE2NjUxNjA0ODR9.oMT9_vwseGUR3h_NbthFwFK02r1aJGfQ8kyFQZk81sM

devolverá algo como:
```
{
    "id": "2FivNlAbeNimVRML4UUU5gS99Q2",
    "email": "josephsosa@gmail.com",
    "password": ""
}
```

En postman crear new WebSocket Request, y hacer conexión pasando en el campo de Headers como Key Authorization, y como Value el Token:

http://localhost:5050/ws

Hacer nuevo post:

http://localhost:5050/api/v1/posts

Ejemplo, en Body pasar como raw tipo Json:
```
{
    "post_content": "New socket post"
}
```

Devolverá algo como:
```
{
    "id": "2Fivq80aJ0z7OM8RdpJ9TdHznpi",
    "post_content": "New socket post"
}
```

Si se va a donde se hizo la onexión del WebSocket, en http://localhost:5050/ws se verá la comunicación del nuevo post realizado.

Se puede visualizar los posts realizados al ir a:

http://localhost:5050/posts?page=0

o bien sin paginar en http://localhost:5050/posts

Igual se puede actualizar un post pasando en el campo de Headers como Key Authorization, y como Value el Token:
http://localhost:5050/api/v1/posts/2Fivq80aJ0z7OM8RdpJ9TdHznpi

y pasando:
```
{
    "post_content": "Nuevo post websocket editado"
}
```

o eliminar el post igual pasando en el campo de Headers como Key Authorization, y como Value el Token:

http://localhost:5050/api/v1/posts/2Fivq80aJ0z7OM8RdpJ9TdHznpi


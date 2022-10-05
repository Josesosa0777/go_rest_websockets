// Definir un paquete al cual pertenece este archivo:
package server

import (
	"context"
	"errors"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
	database "platzi.com/go/rest-ws/database"
	repository "platzi.com/go/rest-ws/repository"
	websocket "platzi.com/go/rest-ws/websocket"
)

// Definir un struct para la configuración que el servidor requiere para poder ejecutarse,
// definir el puerto, la llave secreta y la conexión a la db:
type Config struct {
	Port        string
	JWTSecret   string
	DatabaseUrl string
}

// Interfaz para el tipo Server, se requiere un método Config() que lo que hará será retornar algo de tipo Config
// Lo que indica que para que algo sea llamado un servidor, deberá tener algo de Config que retorne una Configuración tal como se definió arriba (Con port, llave y conexión a db)
type Server interface {
	Config() *Config
	Hub() *websocket.Hub // con websocket, pasar el HUB en el server
}

// Definir el broker que será quien maneje los servidores que tendrá un archivo de configuración (config) con las propiedades definidas previamente (Con port, llave y conexión a db)
// También tendrá un router que definirá las rutas que el API tendrá, se usa la dependencia de tipo mux (*mux.Router)
type Broker struct {
	config *Config
	router *mux.Router
	hub    *websocket.Hub
}

// Se requiere que el broker satisfaga la interface, se crea un receiver function llamado Config() que retornará una configuración (*Config)
// y lo que se hace es devolver la configuración
func (b *Broker) Config() *Config {
	return b.config
}

// crear nueva función llamada Hub perteneciente al broker, que devolverá un websocket.Hub:
func (b *Broker) Hub() *websocket.Hub {
	return b.hub
}

// Definir el constructor para el struct, que recibe 2 parámetros, primero un contexto que se usará para encontrar posibles problemas en el código,
// y el segundo parámetro será la configuración, de tipo previamente definido (*Config), luego hay que retornar el Broker y si hay un error, devolver ese error
func NewServer(ctx context.Context, config *Config) (*Broker, error) {
	// A los campos de configuración, por defecto go les da strings vacíos a cada uno, por tanto hay que asegurarse de que no sean campos vacíos,
	// Si algún campo está vacío, que retorne el error sobre lo que sucede:
	if config.Port == "" {
		return nil, errors.New("port is required")
	}
	if config.JWTSecret == "" {
		return nil, errors.New("jwt secret is required")
	}
	if config.DatabaseUrl == "" {
		return nil, errors.New("database url is required")
	}
	// Si no hay errores, entonces retornar el broker con su configuración y el router nuevo:
	broker := &Broker{
		config: config,
		router: mux.NewRouter(), // Define una nueva instancia del broker
		hub:    websocket.NewHub(),
	}
	return broker, nil
}

// Agregar un método al broker que le permita ejecutarse, en ese caso se llama Start() que recibe una función como parámetro (binder),
// La función binder recibe como parámetro un servidor de tipo Server y un routeador:
func (b *Broker) Start(binder func(s Server, r *mux.Router)) {
	// Trae un router nuevo, la manera como la librería lo hace es con NewRouter():
	b.router = mux.NewRouter()
	// Se pasa el binder que lleva los parámetros b y b.router (b de tipo servidor y b.router del router)
	binder(b, b.router)
	// agregar un handler para manejar las conexiones:
	handler := cors.Default().Handler(b.router)
	// crear repository con la configuracion de la db:
	repo, err := database.NewPostgresRepository(b.config.DatabaseUrl)
	if err != nil {
		log.Fatal(err)
	}
	// Para inicializar el hub:
	go b.hub.Run()
	// se envía el repo en el repository
	repository.SetRepository(repo)
	// Imprimir mensaje con el puerto en el que se ejecuta:
	log.Println("starting server on port", b.config.Port)
	// Ejecutar el servidor:
	// if err := http.ListenAndServe(b.config.Port, b.router); err != nil {
	// cambiar el segundo parámetro que era definido en el router (b.router), por el handler que da la función de Default().Handler
	if err := http.ListenAndServe(b.config.Port, handler); err != nil {
		log.Println("error starting server:", err)
	} else {
		log.Fatalf("server stopped")
	}
}

package websocket

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

// el upgrader es necesario para que la conexion http pueda utilizar un websocket
var upgrader = websocket.Upgrader{
	// se requiere la función CheckOrigin que retorna un true, que quiere decir que todos los clientes? que quieran serán parte de la conexión
	CheckOrigin: func(r *http.Request) bool { return true },
}

// el hub tendrá un slide de clientes, un canal especifico para todos los que se deseen registrar y para todos los clientes que se esten registrando de nuestro hub
// también tendremos un mutex para evitar condiciones de carrera en el programa
type Hub struct {
	clients    []*Client
	register   chan *Client
	unregister chan *Client
	mutex      *sync.Mutex
}

// crear constructor para el hub que devolverá el Hub, :
func NewHub() *Hub {
	return &Hub{
		// el hub tendrá un nuevo slide para clients de longitud 0, y para register y unregister crear el canal y para el mutex se usa lo de la librería de &sync.Mutex{}
		clients:    make([]*Client, 0),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		mutex:      &sync.Mutex{},
	}
}

// Definir la ruta que será usada para los diferentes websockets
func (hub *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	// usar un upgrade para la conexion que permita usar los websockets
	socket, err := upgrader.Upgrade(w, r, nil) // en este caso no se requere un response header, se pone como nil
	if err != nil {
		log.Println(err)
		// si hay un error, no se habrá podido abrir la conexion
		http.Error(w, "Error upgrading connection", http.StatusInternalServerError)
		return
	}
	// crear nuevo client pasando el hub y el socket
	client := NewClient(hub, socket)
	// al hub se le va a registrar el cliente:
	hub.register <- client

	go client.Write() // go routina que se encarga de estar escribiendo
}

// crear receiver function para el hub que le permitirá ejecutarse
func (hub *Hub) Run() {
	for {
		select {
		// crear multiplexacion de los chanels, caso donde un cliente se acaba de registar en la aplicacion:
		case client := <-hub.register:
			hub.onConnect(client)
		// caso de cliente que se está desregistrando:
		case client := <-hub.unregister:
			hub.onDisconnect(client)
		}
	}
}

// onConnect pasar el parametro tipo client
func (hub *Hub) onConnect(client *Client) {
	// imprime que un cliente se está conectando, y la dirección que usa para coinectarse, con client.socket.RemoteAddr()
	log.Println("Client connected", client.socket.RemoteAddr())

	// se bloquea el programa para evitar condiciones de carrera, porque se hará una modificación a los clientes que están conectados:
	hub.mutex.Lock()
	// al final se debe desbloquear el programa:
	defer hub.mutex.Unlock()
	// signar un id al client (que será la conexión que está usando para conectarse):
	client.id = client.socket.RemoteAddr().String()
	// agregar el nuevo cliente al hub de clientes que ya se tiene en existencia
	hub.clients = append(hub.clients, client)
}

// la funcipon inDisconnect recibe también un client como parámetro
func (hub *Hub) onDisconnect(client *Client) {
	log.Println("Client disconnected", client.socket.RemoteAddr())

	// Se cierra la conexión de ese cliente
	client.Close()
	// hay que remover el cliente del NewHub, por lo tanto primero se bloquea:
	hub.mutex.Lock()
	// al finalizar desbloquearlo:
	defer hub.mutex.Unlock()

	// iterar a través de los clientes para buscar el que se ha desconectado
	i := -1
	for j, c := range hub.clients {
		if c.id == client.id {
			i = j
			break
		}
	}
	// remover del hub el client encontrado:
	copy(hub.clients[i:], hub.clients[i+1:])
	hub.clients[len(hub.clients)-1] = nil
	hub.clients = hub.clients[:len(hub.clients)-1]

}

// función Broadcast que recibe como parámetro un mensaje que será una interface, es decir podemos transmitir cualquier tipo de data,
// y un cliente que puede enviar un mensaje a otro y el no quere recibir al mismo tiempo ese mensaje, por eso se pone usa el parámetro de ignore
func (hub *Hub) Broadcast(message interface{}, ignore *Client) {
	// serializar la data usando json.Marshal:
	data, _ := json.Marshal(message)
	// enviar a través del hub de client,
	for _, client := range hub.clients {
		// se hace la evaluación de que el cliente sea diferente, si es el caso se usa el outbpund para enviar la data que se desea:
		if client != ignore {
			client.outbound <- data
		}
	}
}

package websocket

import "github.com/gorilla/websocket"

// definir un struct para el client que manejará las conexiones de diferentes clientes
type Client struct {
	// recibe un hub, un id para identificar al cliente, un socket de tipo conexion del websocket, y un canal de go llamaado outbound que servirá para enviar mensajes como si fueran byte
	hub      *Hub
	id       string
	socket   *websocket.Conn
	outbound chan []byte
}

// crear new client que recibe un hub y un socket y devuelve un client
func NewClient(hub *Hub, socket *websocket.Conn) *Client {
	return &Client{
		hub:      hub,
		socket:   socket,
		outbound: make(chan []byte), // canal nuevo
	}
}

// para enviar los mensajes que yo quiera creo una función llamad Write que pertenece al cliente, lo queremos en este caso para transmitir posts en tiempo real:
func (c *Client) Write() {
	for {
		// va a escuchar la multiplexacion de los diferentes chanels
		select {
		// caso de un mensaje que vendrá de un canal outbound
		case message, ok := <-c.outbound:
			if !ok {
				// si falla, cerrar la conexión de ese mensaje fallido, escribir un mensaje diciendo que la conexion se ha cerrado, envía un arreglo de bytes vacío, porque no requiero enviar data en este caso
				c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}
			// En caso de que se recibió el mensaje, escribirlo para transmitirlo
			c.socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

func (c Client) Close() {
	c.socket.Close()
	close(c.outbound)
}

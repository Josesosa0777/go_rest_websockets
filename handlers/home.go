// Definir el paquete al cual pertenece este archivo:
package handlers

import (
	"encoding/json"
	"net/http"

	"platzi.com/go/rest-ws/server" // busca la dependencia dentro del paquete server
)

// Crear HomeResponse que es lo que se devolverá al cliente:
type HomeResponse struct {
	Message string `json:"message"` // nota que en go usa Message, pero al serializar a json, será message
	Status  bool   `json:"status"`
}

// Definir primer handler que se llama HomeHandler y recibirá como parámetro un servidor y se devolverá como respuesta del llamado un http.HandlerFunc
func HomeHandler(s server.Server) http.HandlerFunc {
	// Se retorna una función que recibe un escritor (http.ResponseWriter) encargado de responderle al cliente, y también tendrá un request, que es la data que envía el cliente
	return func(w http.ResponseWriter, r *http.Request) {
		// Le ponemos un header a la petición http, con el content-Type de tipo application/json, y luego otro header con el status de la llamada:
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK) // status 200
		// Crear la respuesta que se dará, en donde se pasa el escritor (w) al codificador de tipo json, y se hace encode de la respuesta de tipo HomeResponse que tiene las propiedades de Message y Status:
		json.NewEncoder(w).Encode(HomeResponse{
			Message: "Welcome to Platzi Go",
			Status:  true,
		})
	}
}

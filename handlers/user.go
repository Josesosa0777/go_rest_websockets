package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/segmentio/ksuid"
	"platzi.com/go/rest-ws/models"
	"platzi.com/go/rest-ws/repository"
	"platzi.com/go/rest-ws/server"
)

// Lo que se espera devolver
type SignUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// Tipo de request que se espera para que el usuario se pueda registrar
type SignUpRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Función encargada de manejar el handler, recibe un servidor como parametro y devuelve como respuesta un HandlerFunction
func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Creación del request de tipo SignUpRequest
		var request = SignUpRequest{}
		// Se pasa el cuerpo de la petición para ser decodificado
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			// Si hay error, responder un 400 con http.StatusBadRequest, error del lado del cliente
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// retornar id random de librería instalada:
		id, err := ksuid.NewRandom()
		if err != nil {
			// Si hay error, responder un 500 con http.StatusInternalServerError
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// si no hay errores, crear variable del usuario:
		var user = models.User{
			Email:    request.Email,
			Password: request.Password,
			Id:       id.String(),
		}
		// Insertar el usuario a la db usando el repository
		err = repository.InsertUser(r.Context(), &user)
		if err != nil {
			// Si hay error, responder un 500 con http.StatusInternalServerError
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// pasar el header de tipo application/json
		w.Header().Set("Content-Type", "application/json")
		// enviar la respuesta codificada:
		json.NewEncoder(w).Encode(SignUpResponse{
			ID:    user.Id,
			Email: user.Email,
		})

	}
}

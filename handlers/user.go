package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/segmentio/ksuid"
	"golang.org/x/crypto/bcrypt"
	"platzi.com/go/rest-ws/models"
	"platzi.com/go/rest-ws/repository"
	"platzi.com/go/rest-ws/server"
)

const (
	HASH_COST = 8
)

// Lo que se espera devolver
type SignUpResponse struct {
	ID    string `json:"id"`
	Email string `json:"email"`
}

// se envia token de tipo string que cuando se pase a json se llamará token
type LoginResponse struct {
	Token string `json:"token"`
}

// Tipo de request que se espera para que el usuario se pueda registrar
type SignUpLoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Función encargada de manejar el handler, recibe un servidor como parametro y devuelve como respuesta un HandlerFunction
func SignUpHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Creación del request de tipo SignUpRequest
		var request = SignUpLoginRequest{}
		// Se pasa el cuerpo de la petición para ser decodificado
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			// Si hay error, responder un 400 con http.StatusBadRequest, error del lado del cliente
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		// función para generar password hashed:
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(request.Password), HASH_COST) // HASH_COST es cuantas veces va a estar pasando el algoritmo, si no pusieramos, por defecto se debería poner bcrypt.DefaultCost
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
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
			Password: string(hashedPassword),
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

// Función LoginHandler que recibe un server y retorna un htt..HandlerFunc
func LoginHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var request = SignUpLoginRequest{} // para ver los valores que envía el cliente:
		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		user, err := repository.GetUserByEmail(r.Context(), request.Email)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		if user == nil {
			http.Error(w, "Invalid Credentials", http.StatusUnauthorized)
			return
		}
		// Comparar lo que está almacenado en la db con lo que se está pasando por el usuario:
		log.Println(user.Password)
		if err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(request.Password)); err != nil {
			log.Println(err)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// Crear claim en el que se va a pasar el userId, y de StandardClaims pasar el parámetro de cuanto tiempo durará el token:
		claims := models.AppClaims{
			UserId: user.Id,
			StandardClaims: jwt.StandardClaims{
				ExpiresAt: time.Now().Add(2 * time.Hour * 24).Unix(), // token que dure 2 días
			},
		}
		// creación de token con parámetro algoritmo de firmado HS256
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		// creación de string a partir del token usando el secret que se pasa en la configuración
		tokenString, err := token.SignedString([]byte(s.Config().JWTSecret))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		// pasar en la structura de LoginResponse el token
		json.NewEncoder(w).Encode(LoginResponse{
			Token: tokenString,
		})

	}
}

// handler para recibir un token, decodificarlo, validarlo y devolver la data del usuario registrado con ese token
func MeHandler(s server.Server) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Obtener el token:
		tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
		// validar el token
		token, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			return []byte(s.Config().JWTSecret), nil
		})
		// Si hay un error:
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}
		// crear variable llamada claim para devolver un modelo como el definido en models/claims.go
		if claims, ok := token.Claims.(*models.AppClaims); ok && token.Valid {
			// Devolver el usuario:
			user, err := repository.GetUserByID(r.Context(), claims.UserId)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			// Si no hay error:
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(user)
		} else {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	}
}

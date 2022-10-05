package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"
	"platzi.com/go/rest-ws/models"
	"platzi.com/go/rest-ws/server"
)

// crear variable con rutas de login y signup, para que si el middleware las encuentra no las revise, porque no es necesario que lo haga el middleware,
// porque el login devolverá el token, y si lo pidiera y aun no se tiene, no se podria ir al login, o nadie podria registrarse
var (
	NO_AUTH_NEEDED = []string{
		"login",
		"signup",
	}
)

// helper que recibe como parámetro la ruta y devolverá un booleano, revisa si las ruta que se pasa está en las que que no se necesitan
func shoulCheckToken(route string) bool {
	for _, p := range NO_AUTH_NEEDED {
		if strings.Contains(route, p) {
			return false // si la ruta existe, devolver false porque no queremos revisar los tokens de esas rutas ya que son rutas no protegidas
		}
	}
	return true
}

// función que recibe como parámetro un server y devuelve una función que recibe un handler de tipo http.Handler y esto a su vez devolverá un http.Handler:
// se recibe una funcion y se devuelve otra de tipo http.Handler porque el middleware debe recibir la función para hacer un salto al siguiente handler
func CheckAuthMiddleware(s server.Server) func(h http.Handler) http.Handler {
	// retorna funcion next de http.Handler y a su vez devuelve un http.Handler:
	return func(next http.Handler) http.Handler {
		// retornar el valor de http.Handler:
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// revisar si la ruta se debe autenticar:
			if !shoulCheckToken(r.URL.Path) {
				next.ServeHTTP(w, r) // si la ruta no está protegida, entonces puede seguir sin el token, por eso se llama al next (al siguiente handler)
				return
			}
			// Si la ruta está protegida validar el token
			tokenString := strings.TrimSpace(r.Header.Get("Authorization"))
			// En ParseWithClaims pasar el token, los claims (tipo de dato para descompilar el token), y recibe una funcion que tiene como parametro un token que devolverá una interface vacía o un error, y debe devolver el secret:
			_, err := jwt.ParseWithClaims(tokenString, &models.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
				return []byte(s.Config().JWTSecret), nil
			})
			// si existe un error (token vencido, token inválido, etc) devolver el error:
			if err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}
			// si todo va bien, entonces se enviará al siguiente handler:
			next.ServeHTTP(w, r)
		})

	}
}

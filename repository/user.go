// Definir el paquete al cual pertenece este archivo:
package repository

import (
	"context"

	"platzi.com/go/rest-ws/models"
)

// Definir interfaz UserRepository que tendrá un insertar usuario y tendrá un contexto como parámetro y un user de nuestros modelos, y devuelve un error si lo hay
// También traerá un usuario by Id, se pasa el contexto y el id que sea tipo string. Retornará un usuario y un error si lo hay
type UserRepository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	Close() error // Se agrega close, para cerrar conexiones a la db cuando la app no esté corriendo, en este caso, también agregamos que devuelva un error si existe
}

// Crear variable implementation que será de tipo UserRepository:
var implementation UserRepository

// Crear funcion SetRepository que recibe un repository de tipo UserRepository y lo asigna, no importa si usa Postgress, Mongo, etc,
// Se inyecta el SetRepository en la interface UserRepository para que la implementación haga lo que tiene que hacer
func SetRepository(repository UserRepository) {
	implementation = repository
}

// Se crea la funcion InsertUser y se le pasa el contexto, y el user, y lo que hace, relaciona el user con lo de la implementation de manera abstracta:
func InsertUser(ctx context.Context, user *models.User) error {
	return implementation.InsertUser(ctx, user)
}

// Se crea la funcion GetUserByID y se le pasa el contexto, el id, y devuelve el User:
func GetUserByID(ctx context.Context, id string) (*models.User, error) {
	return implementation.GetUserByID(ctx, id)
}

// Se crea la funcion Close, que devolverá lo que la implementación esté haciendo:
func Close() error {
	return implementation.Close()
}

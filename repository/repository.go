// Definir el paquete al cual pertenece este archivo:
package repository

import (
	"context"

	"platzi.com/go/rest-ws/models"
)

// Definir interfaz Repository que tendrá un insertar usuario y tendrá un contexto como parámetro y un user de nuestros modelos, y devuelve un error si lo hay
// También traerá un usuario by Id, se pasa el contexto y el id que sea tipo string. Retornará un usuario y un error si lo hay
type Repository interface {
	InsertUser(ctx context.Context, user *models.User) error
	GetUserByID(ctx context.Context, id string) (*models.User, error)
	GetUserByEmail(ctx context.Context, email string) (*models.User, error) // para autenticar a un usuario
	InsertPost(ctx context.Context, post *models.Post) error
	GetPostByID(ctx context.Context, id string) (*models.Post, error)
	DeletePost(ctx context.Context, id string, userId string) error
	UpdatePost(ctx context.Context, post *models.Post, userId string) error
	ListPost(ctx context.Context, page uint64) ([]*models.Post, error)
	Close() error // Se agrega close, para cerrar conexiones a la db cuando la app no esté corriendo, en este caso, también agregamos que devuelva un error si existe
}

// Crear variable implementation que será de tipo Repository:
var implementation Repository

// Crear funcion SetRepository que recibe un repository de tipo Repository y lo asigna, no importa si usa Postgress, Mongo, etc,
// Se inyecta el SetRepository en la interface Repository para que la implementación haga lo que tiene que hacer
func SetRepository(repository Repository) {
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

func GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	return implementation.GetUserByEmail(ctx, email)
}

func InsertPost(ctx context.Context, post *models.Post) error {
	return implementation.InsertPost(ctx, post)
}

func GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	return implementation.GetPostByID(ctx, id)
}

func DeletePost(ctx context.Context, id string, userId string) error {
	return implementation.DeletePost(ctx, id, userId)
}

func UpdatePost(ctx context.Context, post *models.Post, userId string) error {
	return implementation.UpdatePost(ctx, post, userId)
}

func ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	return implementation.ListPost(ctx, page)
}

// Se crea la funcion Close, que devolverá lo que la implementación esté haciendo:
func Close() error {
	return implementation.Close()
}

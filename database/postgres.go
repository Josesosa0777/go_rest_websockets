package database

import (
	"context"
	"database/sql"
	"log"

	_ "github.com/lib/pq"
	"platzi.com/go/rest-ws/models"
)

// crear la representación de la conexión con la db en PostgresRepository:
type PostgresRepository struct {
	db *sql.DB
}

// Crear el constructor que recibe como parametro la URL que indica a donde se debe hacer la conexcion de la db, se retorna el repositorio PostgresRepository
func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

// Crear la funcion de tipo PostgresRepository, para insertar el User a la db, se crea como un receiver function, a la función se le pasa el context y el user que viene de los modelos de Usuario, y devolver un error si existe:
func (repo *PostgresRepository) InsertUser(ctx context.Context, user *models.User) error {
	// EJecución de sql para insertar el usuario, el ExecContext devuelve el resultado de sql y el error, si no requiero el resultado de sql, le pongo _
	_, err := repo.db.ExecContext(ctx, "INSERT INTO users (id, email, password) VALUES ($1, $2, $3)", user.Id, user.Email, user.Password)
	return err
}

// Crear la funcion de tipo PostgresRepository, para insertar el Post a la db, se crea como un receiver function, a la función se le pasa el context y el post que viene de los modelos de Post, y devolver un error si existe:
func (repo *PostgresRepository) InsertPost(ctx context.Context, post *models.Post) error {
	// EJecución de sql para insertar el post, el ExecContext devuelve el resultado de sql y el error, si no requiero el resultado de sql, le pongo _
	_, err := repo.db.ExecContext(ctx, "INSERT INTO posts (id, post_content, user_id) VALUES ($1, $2, $3)", post.Id, post.PostContent, post.UserId)
	return err
}

// Crear funcion de tipo PostgresRepository, que se llama GetUserByID, se le pasa el contexto y el id de tipo string, devolverá un usuario o un error
func (repo *PostgresRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	// Se hace la query a la db, en la que se pasa el contexto y la query, y lo que devuelve sería las filas de la query y si hay algun error:
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email FROM users WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	// Ya que el row devuelve lecturas, hay que cerrar la conexión cuando se termine de ejecutar:
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// crear la variable que se va a devolver:
	var user = models.User{}
	// Crear la función de parseo que pase los rows al user:
	for rows.Next() {
		// Checar si hay un error al hacer un Scan, Scan permite copiar las columnas que se leen dentro de un la interfaz que se definió (en user)
		if err = rows.Scan(&user.Id, &user.Email); err == nil {
			return &user, nil
		}
	}
	// Si hay error, retornar usuario nulo, y el error:
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// Si no hay problema, retornar el usuario y error nulo:
	return &user, nil
}

// Crear funcion de tipo PostgresRepository, que se llama GetUserByEmail, se le pasa el contexto y el email de tipo string, devolverá un usuario o un error
func (repo *PostgresRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Se hace la query a la db, en la que se pasa el contexto y la query, y lo que devuelve sería las filas de la query y si hay algun error:
	rows, err := repo.db.QueryContext(ctx, "SELECT id, email, password FROM users WHERE email = $1", email)
	if err != nil {
		return nil, err
	}
	// Ya que el row devuelve lecturas, hay que cerrar la conexión cuando se termine de ejecutar:
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	// crear la variable que se va a devolver:
	var user = models.User{}
	// Crear la función de parseo que pase los rows al user:
	for rows.Next() {
		// Checar si hay un error al hacer un Scan, Scan permite copiar las columnas que se leen dentro de un la interfaz que se definió (en user)
		if err = rows.Scan(&user.Id, &user.Email, &user.Password); err == nil {
			return &user, nil
		}
	}
	// Si hay error, retornar usuario nulo, y el error:
	if err = rows.Err(); err != nil {
		return nil, err
	}
	// Si no hay problema, retornar el usuario y error nulo:
	return &user, nil
}

// Crear función que se encarga de cerrar la conexión de la db cuando ya no se requiera
func (repo *PostgresRepository) Close() error {
	return repo.db.Close()
}

func (repo *PostgresRepository) GetPostByID(ctx context.Context, id string) (*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts WHERE id = $1", id)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var post = models.Post{}
	for rows.Next() {
		// van en orden id, post_content, user_id, created_at
		if err = rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err == nil {
			return &post, nil
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return &post, nil
}

func (repo *PostgresRepository) DeletePost(ctx context.Context, id string, userId string) error {
	_, err := repo.db.ExecContext(ctx, "DELETE FROM posts WHERE id = $1 and user_id = $2", id, userId)
	return err
}

func (repo *PostgresRepository) UpdatePost(ctx context.Context, post *models.Post, userId string) error {
	_, err := repo.db.ExecContext(ctx, "UPDATE posts SET post_content = $1 WHERE id = $2 and user_id = $3", post.PostContent, post.Id, userId)
	return err
}

func (repo *PostgresRepository) ListPost(ctx context.Context, page uint64) ([]*models.Post, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, post_content, user_id, created_at FROM posts LIMIT $1 OFFSET $2", 5, page*5)
	if err != nil {
		return nil, err
	}
	defer func() {
		err = rows.Close()
		if err != nil {
			log.Fatal(err)
		}
	}()
	var posts []*models.Post
	for rows.Next() {
		var post = models.Post{}
		if err = rows.Scan(&post.Id, &post.PostContent, &post.UserId, &post.CreatedAt); err == nil {
			posts = append(posts, &post)
		}
	}
	if err = rows.Err(); err != nil {
		return nil, err
	}
	return posts, nil
}

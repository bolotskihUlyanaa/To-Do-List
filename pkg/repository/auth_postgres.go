package repository

import (
	"fmt"

	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/jmoiron/sqlx"
)

// структура имплементирует интерфейс репозитория и работает с базой postgres
type AuthPostgres struct {
	db *sqlx.DB
}

func NewAuthPostgres(db *sqlx.DB) *AuthPostgres {
	return &AuthPostgres{db: db}
}

func (r *AuthPostgres) CreateUser(user todolist.User) (int, error) {
	var id int
	//числа с $ - плейсхолдер, в которые будут подставлены значения,
	//которые мы передадим в качестве аргументов к функции для выполнения запроса к бд
	//возвращает id новой записи после операции insert
	query := fmt.Sprintf("INSERT INTO %s (name, username, password_hash) values ($1, $2, $3) RETURNING id", usersTable)
	//возвращает объект row, те хранит в себе информацию о возвращаемой строке из базы,
	//в нашел случае одна строка с полем id
	row := r.db.QueryRow(query, user.Name, user.Username, user.Password) //user.Name, user.Username, user.Password будут подставлены на места плейсхолдеров
	if err := row.Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (r *AuthPostgres) GetUser(username, password string) (todolist.User, error) {
	var user todolist.User
	query := fmt.Sprintf("SELECT id FROM %s WHERE username=$1 AND password_hash=$2", usersTable)
	err := r.db.Get(&user, query, username, password)
	return user, err
}

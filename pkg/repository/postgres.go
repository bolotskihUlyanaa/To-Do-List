// реализация логики подключения бд
package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

// названия таблиц
const (
	usersTable      = "users"
	todoListTable   = "todo_lists"
	usersListsTable = "users_lists"
	todoItemsTable  = "todo_items"
	listsItemsTable = "lists_items"
)

// структура параметров, необходимых для подключения к бд
type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	SSLMode  string
}

func NewPostgresDB(cfg Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("postgres", fmt.Sprintf("host=%s port=%s user=%s dbname=%s password=%s sslmode=%s",
		cfg.Host, cfg.Port, cfg.Username, cfg.DBName, cfg.Password, cfg.SSLMode))
	if err != nil {
		return nil, err
	}
	//проверим подключение
	err = db.Ping()
	if err != nil {
		return nil, err
	}
	return db, nil
}

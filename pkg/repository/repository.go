package repository

import (
	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/jmoiron/sqlx"
)

//заготовки интерфейсов для сущностей
//интерфейсы называем исходя из их доменной зоны
//те участки бизнес логики приложения за которые они ответственны

type Authorization interface {
	CreateUser(user todolist.User) (int, error)
	GetUser(username, password string) (todolist.User, error)
}

type TodoList interface {
	Create(userId int, list todolist.ToDoList) (int, error)
	GetAll(userId int) ([]todolist.ToDoList, error)
	GetById(userId, listId int) (todolist.ToDoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todolist.UpdateListInput) error
}

type TodoItem interface {
	Create(listId int, item todolist.ToDoItem) (int, error)
	GetAll(userId, listId int) ([]todolist.ToDoItem, error)
	GetById(userId, itemId int) (todolist.ToDoItem, error)
	Update(userId, itemId int, input todolist.UpdateItemInput) error
	Delete(userId, itemId int) error
}

// структура которая будет собирать все сервисы в одном месте
type Repository struct {
	Authorization
	TodoList
	TodoItem
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Authorization: NewAuthPostgres(db),
		TodoList:      NewTodoListPostgres(db),
		TodoItem:      NewTodoItemPostgres(db),
	}
}

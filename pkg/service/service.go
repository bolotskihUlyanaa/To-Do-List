package service

import (
	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/repository"
)

//заготовки интерфейсов для сущностей
//интерфейсы называем исходя из их доменной зоны
//те участки бизнес логики приложения за которые они ответственны

type Authorization interface {
	CreateUser(user todolist.User) (int, error)              //возвращает id созданного пользователя и ошибку
	GenerateToken(username, password string) (string, error) //возвращает сгенерированный токен и ошибку
	ParseToken(token string) (int, error)                    //возвращает id клиента при успешном парсинге
}

type TodoList interface {
	Create(userId int, list todolist.ToDoList) (int, error)
	GetAll(userId int) ([]todolist.ToDoList, error)
	GetById(userId, listId int) (todolist.ToDoList, error)
	Delete(userId, listId int) error
	Update(userId, listId int, input todolist.UpdateListInput) error
}

type TodoItem interface {
	Create(userId, listId int, item todolist.ToDoItem) (int, error)
	GetAll(userId, listId int) ([]todolist.ToDoItem, error)
	GetById(userId, itemId int) (todolist.ToDoItem, error)
	Update(userId, itemId int, item todolist.UpdateItemInput) error
	Delete(userId, itemId int) error
}

// структура которая будет собирать все сервисы в одном месте
type Service struct {
	Authorization
	TodoList
	TodoItem
}

// сервисы будут обращаться к базе данных поэтому указатель на структуру репозиторий(внедрение зависимостей)
func NewService(repos *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repos.Authorization),
		TodoList:      NewToDoListService(repos.TodoList),
		TodoItem:      NewToDoItemService(repos.TodoItem, repos.TodoList),
	}
}

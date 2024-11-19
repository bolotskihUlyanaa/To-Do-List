// сервис для работы со списком
package service

import (
	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/repository"
)

type ToDoListService struct {
	repo repository.TodoList
}

func NewToDoListService(repo repository.TodoList) *ToDoListService {
	return &ToDoListService{repo: repo}
}

// передаем данные на след уровень в репозиторий
func (s *ToDoListService) Create(userId int, list todolist.ToDoList) (int, error) {
	return s.repo.Create(userId, list)
}

func (s *ToDoListService) GetAll(userId int) ([]todolist.ToDoList, error) {
	return s.repo.GetAll(userId)
}

func (s *ToDoListService) GetById(userId, listId int) (todolist.ToDoList, error) {
	return s.repo.GetById(userId, listId)
}

func (s *ToDoListService) Delete(userId, listId int) error {
	return s.repo.Delete(userId, listId)
}

func (s *ToDoListService) Update(userId, listId int, input todolist.UpdateListInput) error {
	if err := input.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, listId, input)
}

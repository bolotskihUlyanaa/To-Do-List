package service

import (
	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/repository"
)

type ToDoItemService struct {
	repo     repository.TodoItem
	listRepo repository.TodoList
}

func NewToDoItemService(repo repository.TodoItem, listRepo repository.TodoList) *ToDoItemService {
	return &ToDoItemService{
		repo:     repo,
		listRepo: listRepo,
	}
}

func (s *ToDoItemService) Create(userId, listId int, item todolist.ToDoItem) (int, error) {
	//проверка на наличие списка и принадлежность его определенному пользователю
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return 0, err
	}
	return s.repo.Create(listId, item)
}

func (s *ToDoItemService) GetAll(userId, listId int) ([]todolist.ToDoItem, error) {
	_, err := s.listRepo.GetById(userId, listId)
	if err != nil {
		return nil, err
	}
	return s.repo.GetAll(userId, listId)
}

func (s *ToDoItemService) GetById(userId, itemId int) (todolist.ToDoItem, error) {
	return s.repo.GetById(userId, itemId)
}

func (s *ToDoItemService) Update(userId, itemId int, item todolist.UpdateItemInput) error {
	if err := item.Validate(); err != nil {
		return err
	}
	return s.repo.Update(userId, itemId, item)
}

func (s *ToDoItemService) Delete(userId, itemId int) error {
	return s.repo.Delete(userId, itemId)
}

package repository

import (
	"fmt"
	"strings"

	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/jmoiron/sqlx"
)

type TodoItemPostgres struct {
	db *sqlx.DB
}

func NewTodoItemPostgres(db *sqlx.DB) *TodoItemPostgres {
	return &TodoItemPostgres{db: db}
}

func (r *TodoItemPostgres) Create(listId int, item todolist.ToDoItem) (int, error) {
	tx, err := r.db.Begin()
	if err != nil {
		return 0, err
	}
	var id int
	createItemQuery := fmt.Sprintf("INSERT INTO %s (title, description, done) VALUES($1, $2, $3) RETURNING ID", todoItemsTable)
	row := tx.QueryRow(createItemQuery, item.Title, item.Description, item.Done)
	if err = row.Scan(&id); err != nil {
		tx.Rollback()
		return 0, err
	}
	createTodoLists := fmt.Sprintf("INSERT INTO %s (item_id, list_id) VALUES($1, $2) RETURNING id", listsItemsTable)
	_, err = tx.Exec(createTodoLists, id, listId)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	tx.Commit()
	return id, nil
}

func (r *TodoItemPostgres) GetAll(userId, listId int) ([]todolist.ToDoItem, error) {
	var items []todolist.ToDoItem
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, tl.done FROM %s tl INNER JOIN %s ul ON tl.id = ul.item_id INNER JOIN %s nl ON ul.list_id = nl.list_id WHERE ul.list_id = $1 AND nl.user_id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	err := r.db.Select(&items, query, listId, userId)
	if err != nil {
		return nil, err
	}
	return items, nil

}

func (r *TodoItemPostgres) GetById(userId, itemId int) (todolist.ToDoItem, error) {
	var item todolist.ToDoItem
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description, tl.done FROM %s tl INNER JOIN %s ul ON tl.id = ul.item_id INNER JOIN %s nl ON ul.list_id = nl.list_id WHERE nl.user_id = $1 AND tl.id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	if err := r.db.Get(&item, query, userId, itemId); err != nil {

		return item, err
	}
	return item, nil
}

func (r *TodoItemPostgres) Update(userId, itemId int, input todolist.UpdateItemInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	if input.Title != nil {
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}
	if input.Done != nil {
		setValues = append(setValues, fmt.Sprintf("done=$%d", argId))
		args = append(args, *input.Done)
		argId++
	}
	args = append(args, userId, itemId)
	querySet := strings.Join(setValues, ",")
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul INNER JOIN %s nl ON ul.list_id = nl.list_id WHERE nl.user_id = $%d AND tl.id = $%d", todoItemsTable, querySet, listsItemsTable, usersListsTable, argId, argId+1)
	_, err := r.db.Exec(query, args...)
	return err
}

func (r *TodoItemPostgres) Delete(userId, itemId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul, %s nl WHERE tl.id = ul.item_id AND ul.list_id=nl.list_id AND tl.id = $1 AND nl.user_id = $2", todoItemsTable, listsItemsTable, usersListsTable)
	_, err := r.db.Exec(query, itemId, userId)
	return err
}

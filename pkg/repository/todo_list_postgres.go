package repository

import (
	"fmt"
	"strings"

	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
)

type TodoListPostgres struct {
	db *sqlx.DB
}

func NewTodoListPostgres(db *sqlx.DB) *TodoListPostgres {
	return &TodoListPostgres{db: db}
}

func (r *TodoListPostgres) Create(userId int, list todolist.ToDoList) (int, error) {
	tx, err := r.db.Begin() //создание транзакции
	if err != nil {
		return 0, err
	}
	var id int
	//запрос для вставки в таблицу todolists
	createListQuery := fmt.Sprintf("INSERT INTO %s (title, description) VALUES ($1, $2) RETURNING id", todoListTable)
	row := tx.QueryRow(createListQuery, list.Title, list.Description)
	if err := row.Scan(&id); err != nil {
		//откатывает все изменения базы данных до начала выполнения транзакции
		tx.Rollback()
		return 0, err
	}
	//вставка в таблицу userslists, в которой свяжем id пользователя и id нового списка
	createUsersLists := fmt.Sprintf("INSERT INTO %s (user_id, list_id) VALUES ($1, $2)", usersListsTable)
	//exec для простого выполнения запроса без чтения возвращаемого значения
	_, err = tx.Exec(createUsersLists, userId, id)
	if err != nil {
		tx.Rollback()
		return 0, err
	}
	return id, tx.Commit() //commit - применит изменения к бд и закончит транзакцию
}

func (r *TodoListPostgres) GetAll(userId int) ([]todolist.ToDoList, error) {
	var lists []todolist.ToDoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id = ul.list_id WHERE ul.user_id = $1", todoListTable, usersListsTable)
	err := r.db.Select(&lists, query, userId) //работает как get но применяется при выборки больше одного элемента и записи в слайс
	if err != nil {
		return nil, err
	}
	return lists, nil
}

func (r *TodoListPostgres) GetById(userId, listId int) (todolist.ToDoList, error) {
	var list todolist.ToDoList
	query := fmt.Sprintf("SELECT tl.id, tl.title, tl.description FROM %s tl INNER JOIN %s ul on tl.id=ul.list_id WHERE ul.user_id=$1 AND ul.list_id=$2", todoListTable, usersListsTable)
	err := r.db.Get(&list, query, userId, listId)
	return list, err
}

func (r *TodoListPostgres) Delete(userId, listId int) error {
	query := fmt.Sprintf("DELETE FROM %s tl USING %s ul WHERE tl.id = $1 AND ul.list_id=$2", todoListTable, usersListsTable)
	_, err := r.db.Exec(query, userId, listId)
	return err
}

func (r *TodoListPostgres) Update(userId, listId int, input todolist.UpdateListInput) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1
	//проверка полей
	if input.Title != nil {
		//добавляем в слайсы элементы для формирования запросов в бд для их обновления
		setValues = append(setValues, fmt.Sprintf("title=$%d", argId))
		args = append(args, *input.Title)
		argId++
	}
	if input.Description != nil {
		setValues = append(setValues, fmt.Sprintf("Description=$%d", argId))
		args = append(args, *input.Description)
		argId++
	}

	//title =$1
	//description =$2
	//title =$1, description =$2
	setQuery := strings.Join(setValues, ", ") //соединим элементы строк в одну строку через запятую
	query := fmt.Sprintf("UPDATE %s tl SET %s FROM %s ul WHERE tl.id = ul.list_id AND ul.list_id=$%d AND ul.user_id=$%d",
		todoListTable, setQuery, usersListsTable, argId, argId+1)
	args = append(args, listId, userId)
	//залогируем запрос и аргументы в консоль
	logrus.Debugf("updateQuery: %s", query)
	logrus.Debugf("args: %s", args)
	_, err := r.db.Exec(query, args...)
	return err
}

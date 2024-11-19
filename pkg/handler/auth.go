package handler

import (
	"net/http"

	todolist "github.com/bolotskihUlyanaa/To-Do-List"
	"github.com/gin-gonic/gin"
)

// не уверен что это надо
//чтобы убрать ошибку Running in "debug" mode. Switch to "release" mode in production.
//func init() {
//gin.SetMode(gin.ReleaseMode)
//}

// регистрация, при регистрации от пользователя мы получаем имя, username и password
func (h *Handler) signUp(c *gin.Context) {
	//структура в которую будем записывать данные из json от пользователя
	var input todolist.User

	//BindJSON принимает ссылку на объект в который мы хотим распарсить тело json
	//тк прописали теги для структуры: значения полей из тела json будут присвоены в поля с аналогичными тегами
	if err := c.BindJSON(&input); err != nil {
		//своя функция для создания ответа с ошибкой
		//http.StatusBadRequest = 400 пользователь предоставил некоректные данные в запросе
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//передаем данные на слоц ниже - сервис
	id, err := h.services.Authorization.CreateUser(input)
	if err != nil {
		//StatusInternalServerError = 500 - внутреняя ошибка на сервере
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//StatusOK = 200
	c.JSON(http.StatusOK, map[string]interface{}{
		"id": id,
	})
}

// на endpoint аутентификация необходимо получать логин и пароль от пользователя
type signInInput struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// аутентификация
func (h *Handler) signIn(c *gin.Context) {
	//структура в которую будем записывать данные из json от пользователя
	var input signInInput

	//BindJSON принимает ссылку на объект в который мы хотим распарсить тело json
	//тк прописали теги для структуры: значения полей из тела json будут присвоены в поля с аналогичными тегами
	if err := c.BindJSON(&input); err != nil {
		//своя функция для создания ответа с ошибкой
		//http.StatusBadRequest = 400 пользователь предоставил некоректные данные в запросе
		newErrorResponse(c, http.StatusBadRequest, err.Error())
		return
	}
	//передаем данные на слоц ниже - сервис
	token, err := h.services.Authorization.GenerateToken(input.Username, input.Password)
	if err != nil {
		//StatusInternalServerError = 500 - внутреняя ошибка на сервере
		newErrorResponse(c, http.StatusInternalServerError, err.Error())
		return
	}
	//StatusOK = 200
	c.JSON(http.StatusOK, map[string]interface{}{
		"token": token,
	})

}

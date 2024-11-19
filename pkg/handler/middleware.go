// middleware - прослойка которая парсит токены из запроса и предоставляет доступ к endpoint "/api"
package handler

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	userCtx             = "userId"
)

func (h *Handler) userIdentity(c *gin.Context) {
	//получаем значения из hander авторизации и валидириуем его(не пустая строка и состоит из 2х частей)
	header := c.GetHeader(authorizationHeader)
	if header == "" { //валидируем что он не пустой
		//статус код 401 - пользователь не авторизирован
		newErrorResponse(c, http.StatusUnauthorized, "empty auth header")
		return
	}
	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 {
		newErrorResponse(c, http.StatusUnauthorized, "invalid auth header")
		return
	}
	//парсим токен и записываем пользователя в контекст
	userId, err := h.services.Authorization.ParseToken(headerParts[1])
	if err != nil {
		newErrorResponse(c, http.StatusUnauthorized, err.Error())
	}
	//запишем значение id  в контекст, чтобы иметь доступ к id пользователя, который делает запрос
	//в последующих обработчиках, которые вызываются после данной прослойки
	c.Set(userCtx, userId)
}

// вынесено в отдельную функцию тк Get возвращает интерфейс и всегда нужно приводить тип
func getUserId(c *gin.Context) (int, error) {
	id, ok := c.Get(userCtx)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	idInt, ok := id.(int)
	if !ok {
		newErrorResponse(c, http.StatusInternalServerError, "user id not found")
		return 0, errors.New("user id not found")
	}
	return idInt, nil
}

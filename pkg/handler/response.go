package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

// структура ошибки
type errorResponse struct {
	Message string `json:"message"`
}

type statusResponce struct {
	Status string `json:"status"`
}

// для стандартной обработки ошибок
func newErrorResponse(c *gin.Context, statusCode int, message string) {
	logrus.Errorf(message) //выводим сообщение об ошибке в консоль

	//принимает статус код и тело ответа,
	//блокирует выполнение последующих обработчиков (когда используется несколько подряд обработчиков)
	//а также записывает в ответ статус код и тело сообщения в формате json
	c.AbortWithStatusJSON(statusCode, errorResponse{message})
}

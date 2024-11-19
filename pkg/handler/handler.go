package handler

import (
	"github.com/bolotskihUlyanaa/To-Do-List/pkg/service"
	"github.com/gin-gonic/gin"
)

type Handler struct {
	services *service.Service
}

func NewHandler(services *service.Service) *Handler {
	return &Handler{services: services}
}

// инициализирует endpointы - это конечная точка API, конкретный URL с помощью которого мы можем взаимодействовать с системой
func (h *Handler) InitRoutes() *gin.Engine { //gin.Engine реализует http.handler
	router := gin.New() //инициализируем роутер

	//объявим методы сгруппировав их по маршрутам
	auth := router.Group("/auth")
	{
		auth.POST("/sign-up", h.signUp) //endpoint для регистрации
		auth.POST("/sign-in", h.signIn) //endpoint для авторизации
	}

	//для работы со списками и задачами
	api := router.Group("/api", h.userIdentity)
	{
		lists := api.Group("/lists") //группа для работы со списками
		{
			lists.POST("/", h.createList)
			lists.GET("/", h.getAllLists)
			lists.GET("/:id", h.getListById)
			lists.PUT("/:id", h.updateList)
			lists.DELETE("/:id", h.deleteList)

			items := lists.Group(":id/items") //группа для задач списков
			{
				items.POST("/", h.createItem)
				items.GET("/", h.getAllItems)
			}
		}
		items := api.Group("items")
		{
			items.GET("/:item_id", h.getItemById)
			items.PUT("/:item_id", h.updateItem)
			items.DELETE("/:item_id", h.deleteItem)
		}
	}
	return router
}

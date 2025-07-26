package routes

import (
	"goProject/postgreapi/controllers"

	"github.com/gin-gonic/gin"
)

func TodoRoutes(router *gin.Engine) {
	userGroup := router.Group("/todos")
	{
		userGroup.POST("/", controllers.TodoCreate)
		userGroup.GET("/", controllers.TodoIndex)
		userGroup.GET("/:id", controllers.TodoShow)
		userGroup.PUT("/:id", controllers.TodoUpdate)
		userGroup.DELETE("/:id", controllers.TodoDelete)
	}
}

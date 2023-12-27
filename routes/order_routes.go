package routes

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-service-go/controllers"
)

func SetupOrderRoutes(router *gin.Engine) {
	orderGroup := router.Group("/order")
	{
		orderGroup.POST("/", controllers.CreateOrder)
		orderGroup.GET("/", controllers.GetAllOrders)
		orderGroup.GET("/:id", controllers.GetOrderByID)
		orderGroup.PUT("/:id", controllers.UpdateOrder)
		orderGroup.DELETE("/:id", controllers.DeleteOrder)
	}
}

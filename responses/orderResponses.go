package responses

import (
	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-service-go/models"
)

func CreateErrorResponse(errors []string) gin.H {
	return gin.H{
		"status":  "error",
		"message": "Request failed",
		"data": gin.H{
			"errors": errors,
		},
	}
}

func CreateSuccessResponseForMultipleOrders(orders []models.Order) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Orders retrieved successfully",
		"data": gin.H{
			"orders": orders,
		},
	}
}

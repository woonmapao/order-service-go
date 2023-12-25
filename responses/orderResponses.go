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

func CreateSuccessResponse(data interface{}) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Operation successful",
		"data":    data,
	}
}

func CreateSuccessResponseForSingleOrder(order models.Order) gin.H {
	return gin.H{
		"status":  "success",
		"message": "Order retrieved successfully",
		"data": gin.H{
			"order": order,
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

func CreateSuccessResponseForMultipleOrderDetails(orderDetails []models.OrderDetail) gin.H {
	// Convert order details to a format suitable for the response
	var formattedOrderDetails []gin.H
	for _, orderDetail := range orderDetails {
		formattedOrderDetails = append(formattedOrderDetails, gin.H{
			"id":        orderDetail.ID,
			"orderID":   orderDetail.OrderID,
			"productID": orderDetail.ProductID,
			"quantity":  orderDetail.Quantity,
			"subtotal":  orderDetail.Subtotal,
		})
	}

	return gin.H{
		"status":  "success",
		"message": "Order details retrieved successfully",
		"data": gin.H{
			"orderDetails": formattedOrderDetails,
		},
	}
}

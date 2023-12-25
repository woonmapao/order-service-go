package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-service-go/initializer"
	"github.com/woonmapao/order-service-go/models"
	"github.com/woonmapao/order-service-go/responses"
)

func GetAllOrders(c *gin.Context) {
	// Fetch a list of all orders from the database

	// Retrieve orders from the database
	var orders []models.Order
	err := initializer.DB.Find(&orders).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch orders",
			}))
		return
	}

	// Check if no orders were found
	if len(orders) == 0 {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"No orders found",
			}))
		return
	}

	// Return a JSON response with the list of orders
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForMultipleOrders(orders),
	)
}

func GetOrderByID(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Convert order ID to integer (validations)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order ID",
			}))
		return
	}

	// Query the database for the order with the specified ID
	var order models.Order
	err = initializer.DB.First(&order, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order",
			}))
		return
	}

	// Check if the order was not found
	if order == (models.Order{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Order not found",
			}))
		return
	}

	// Return a JSON response with the order
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForSingleOrder(order),
	)
}

func CreateOrder(c *gin.Context) {
	// Extract order data from the request body
	var orderData models.Order
	err := c.ShouldBindJSON(&orderData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate the input data
	err = validators.ValidateOrderData(orderData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Create a new order in the database
	err = initializer.DB.Create(&orderData).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create order",
		})
		return
	}

	// Return a JSON response with the newly created order
	c.JSON(http.StatusCreated, gin.H{
		"order": orderData,
	})
}

// UpdateOrder handles the update of an existing order
func UpdateOrder(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Extract updated order data from the request body
	var updatedOrderData models.Order
	err := c.ShouldBindJSON(&updatedOrderData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Validate the input data
	err = validators.ValidateOrderData(updatedOrderData)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	// Get the existing order from the database
	var existingOrder models.Order
	err = initializer.DB.First(&existingOrder, orderID).Error
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error": "Order not found",
		})
		return
	}

	// Update the order in the database
	initializer.DB.Model(&existingOrder).Updates(updatedOrderData)

	// Return a JSON response with the updated order
	c.JSON(http.StatusOK, gin.H{
		"updated_order": existingOrder,
	})
}

// DeleteOrder deletes an order based on its ID
func DeleteOrder(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Delete the order from the database
	err := initializer.DB.Delete(&models.Order{}, orderID).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to delete order",
		})
		return
	}

	// Return a JSON response indicating success
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Order deleted successfully",
	})
}

// GetOrderDetails fetches all details (products) associated with a specific order
func GetOrderDetails(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Query the database for details (products) associated with the order
	var orderDetails []models.OrderDetail
	err := initializer.DB.Where("order_id = ?", orderID).Find(&orderDetails).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to fetch order details",
		})
		return
	}

	// Return a JSON response with the order details
	c.JSON(http.StatusOK, gin.H{
		"order_details": orderDetails,
	})
}

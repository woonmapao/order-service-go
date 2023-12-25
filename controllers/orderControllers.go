package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-service-go/initializer"
	"github.com/woonmapao/order-service-go/models"
	"github.com/woonmapao/order-service-go/responses"
	"github.com/woonmapao/order-service-go/validations"
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
	var body models.Order
	err := c.ShouldBindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
			}))
		return
	}

	// Validate the order data
	err = validations.ValidateOrderData(body)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Create a new order in the database
	order := models.Order{
		UserID:      body.UserID,
		OrderDate:   body.OrderDate,
		TotalAmount: body.TotalAmount,
		Status:      body.Status,
	}
	err = initializer.DB.Create(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create order",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForSingleOrder(order),
	)
}

// UpdateOrder handles the update of an existing order
func UpdateOrder(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Convert order ID to integer (validations)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{"Invalid order ID"}))
		return
	}

	// Extract updated order data from the request body
	var updateData models.Order
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
			}))
		return
	}

	// Validate order data
	err = validations.ValidateOrderData(updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				err.Error(),
			}))
		return
	}

	// Check if the order with the given ID exists
	var order models.Order
	err = initializer.DB.First(&order, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order",
			}))
		return
	}
	if order == (models.Order{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Order not found",
			}))
		return
	}

	// Update order fields
	order.UserID = updateData.UserID
	order.OrderDate = updateData.OrderDate
	order.TotalAmount = updateData.TotalAmount
	order.Status = updateData.Status

	// Save the updated order to the database
	err = initializer.DB.Save(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update order",
			}))
		return
	}

	/// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForSingleOrder(order),
	)
}

// DeleteOrder deletes an order based on its ID
func DeleteOrder(c *gin.Context) {
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

	// Check if the order with the given ID exists
	var order models.Order
	err = initializer.DB.First(&order, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order",
			}))
		return
	}
	if order == (models.Order{}) {
		c.JSON(http.StatusNotFound,
			responses.CreateErrorResponse([]string{
				"Order not found",
			}))
		return
	}

	// Delete the order
	err = initializer.DB.Delete(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete order",
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(nil))
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

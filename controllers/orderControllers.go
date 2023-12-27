package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-service-go/initializer"
	"github.com/woonmapao/order-service-go/models"
	"github.com/woonmapao/order-service-go/responses"
	"github.com/woonmapao/order-service-go/services"
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
				err.Error(),
			}))
		return
	}
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
				err.Error(),
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
				err.Error(),
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
				err.Error(),
			}))
		return
	}

	// Check for empty values
	if body.UserID == 0 || body.OrderDate.IsZero() || body.TotalAmount == 0.0 || body.Status == "" {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"UserID, OrderDate, TotalAmount, and Status are required fields",
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

	// Start a transaction
	tx := initializer.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to begin transaction",
				tx.Error.Error(),
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
	err = tx.Create(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to create order",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to commit transaction",
				err.Error(),
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
			responses.CreateErrorResponse([]string{
				"Invalid order ID",
				err.Error(),
			}))
		return
	}

	// Extract updated order data from the request body
	var updateData models.Order
	err = c.ShouldBindJSON(&updateData)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid request format",
				err.Error(),
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

	// Start a transaction
	tx := initializer.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to begin transaction",
				tx.Error.Error(),
			}))
		return
	}

	// Check if the order with the given ID exists
	var order models.Order
	err = tx.First(&order, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order",
				err.Error(),
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

	// Only update the fields that are present in the request
	if updateData.UserID != 0 {
		order.UserID = updateData.UserID
	}
	if !updateData.OrderDate.IsZero() {
		order.OrderDate = updateData.OrderDate
	}
	if updateData.TotalAmount != 0 {
		order.TotalAmount = updateData.TotalAmount
	}
	if updateData.Status != "" {
		order.Status = updateData.Status
	}

	// Save the updated order to the database
	err = tx.Save(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to update order",
				err.Error(),
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to commit transaction",
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

	// Start a transaction
	tx := initializer.DB.Begin()
	if tx.Error != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to begin transaction",
				tx.Error.Error(),
			}))
		return
	}

	// Check if the order with the given ID exists
	var order models.Order
	err = tx.First(&order, id).Error
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
	err = tx.Delete(&order).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to delete order",
			}))
		return
	}

	// Commit the transaction and check for commit errors
	err = tx.Commit().Error
	if err != nil {
		tx.Rollback()
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to commit transaction",
				err.Error(),
			}))
		return
	}

	// Return success response
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponse(&order))
}

// GetOrderDetails from OrderID
func GetOrderDetails(c *gin.Context) {
	// Extract order ID from the request parameters
	orderID := c.Param("id")

	// Convert order ID to integer (validations)
	id, err := strconv.Atoi(orderID)
	if err != nil {
		c.JSON(http.StatusBadRequest,
			responses.CreateErrorResponse([]string{
				"Invalid order ID",
				err.Error(),
			}))
		return
	}

	// Get the order from the database
	var order models.Order
	err = initializer.DB.First(&order, id).Error
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order",
				err.Error(),
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

	// Fetch order details for the order from order detail service
	orderDetails, err := services.FetchOrderDetailsFromOrderDetailService(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError,
			responses.CreateErrorResponse([]string{
				"Failed to fetch order details",
				err.Error(),
			}))
		return
	}

	// Return success response with order details
	c.JSON(http.StatusOK,
		responses.CreateSuccessResponseForMultipleOrderDetails(orderDetails))
}

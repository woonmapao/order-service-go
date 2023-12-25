package services

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/woonmapao/order-service-go/models"
)

func FetchOrderDetailsFromOrderDetailService(orderID int) ([]models.OrderDetail, error) {
	// Build the URL for the order-detail-service API endpoint
	url := fmt.Sprintf("http://order-detail-service/api/order-details?orderID=%d", orderID)

	// Make an HTTP GET request to the order-detail-service
	response, err := http.Get(url)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch order details: %v", err)
	}
	defer response.Body.Close()

	// Check the HTTP status code
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch order details. Status code: %d", response.StatusCode)
	}

	// Decode the JSON response into a slice of OrderDetail structs
	var orderDetails []models.OrderDetail
	if err := json.NewDecoder(response.Body).Decode(&orderDetails); err != nil {
		return nil, fmt.Errorf("failed to decode order details: %v", err)
	}

	return orderDetails, nil
}

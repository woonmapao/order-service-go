package validations

import (
	"fmt"
	"log"
	"net/http"

	"github.com/woonmapao/order-service-go/models"
)

func ValidateOrderData(data models.Order) error {

	if !userExists(data.UserID) {
		return fmt.Errorf("user with ID %d does not exist", data.UserID)
	}

	return nil
}

// userExists checks if a user with the given ID exists in the database.
func userExists(userID int) bool {
	url := fmt.Sprintf("http://user-service/api/users/%d", userID)

	// Make HTTP GET request to check user existence
	resp, err := http.Get(url)
	if err != nil {
		log.Printf("Error checking user existence")
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

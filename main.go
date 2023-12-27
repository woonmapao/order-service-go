package main

import (
	"os"

	"github.com/gin-gonic/gin"
	"github.com/woonmapao/order-service-go/initializer"
	"github.com/woonmapao/order-service-go/routes"
)

func init() {
	initializer.LoadEnvVariables()
	initializer.DBInitializer()
}

func main() {

	r := gin.Default()

	routes.SetupOrderRoutes(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	r.Run(":" + port)

}

package main

import (
	"goProject/postgreapi/initializers"
	"goProject/postgreapi/routes"
	"log"

	"github.com/gin-gonic/gin"
)

func init() {
	// Initialize the database connection
	initializers.LoadEnvVariables()
	initializers.ConnectDB()

}
func main() {
	r := gin.Default()
	routes.TodoRoutes(r)
	if err := r.Run(":9000"); err != nil {
		log.Fatalf("Failed to run server: %v", err)
	}
}

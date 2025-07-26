package main

import (
	"goProject/postgreapi/initializers"
	"goProject/postgreapi/models"
)

func init() {
	// Initialize the database connection
	initializers.LoadEnvVariables()
	initializers.ConnectDB()

}
func main() {
	initializers.DB.AutoMigrate(&models.Todo{})
}

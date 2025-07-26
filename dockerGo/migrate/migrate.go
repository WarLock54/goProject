package main

import (
	"goProject/dockerGo/initializers"

	"goProject/dockerGo/models"
)

func init() {
	// Initialize the database connection
	initializers.LoadEnvVariables()
	initializers.ConnectDB()

}
func main() {
	initializers.DB.AutoMigrate(&models.Book{},
		&models.Credentials{},
		&models.Claims{})
}

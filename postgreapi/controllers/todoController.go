package controllers

import (
	"goProject/postgreapi/initializers"
	"goProject/postgreapi/models"

	"github.com/gin-gonic/gin"
)

func TodoCreate1(c *gin.Context) {
	var todo models.Todo
	if err := c.ShouldBindJSON(&todo); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}
	if err := initializers.DB.Create(&todo).Error; err != nil {
		c.JSON(500, gin.H{"error": "Failed to create todo"})
		return
	}
	c.JSON(201, todo)
}
func TodoCreate(c *gin.Context) {
	// Get data from req body
	var body struct {
		Content string `json:"Content"`
		Status  bool   `json:"Status"`
	}
	c.Bind(&body)
	todo := models.Todo{Content: body.Content, Status: body.Status}
	result := initializers.DB.Create(&todo)

	if result.Error != nil {
		c.Status(400)
		return
	}
	c.JSON(200, gin.H{"todo": todo})
}
func TodoIndex(c *gin.Context) {
	var todos []models.Todo
	initializers.DB.Find(&todos)
	if len(todos) == 0 {
		c.JSON(404, gin.H{"message": "No todos found"})
		return
	}
	// Return todos in response
	c.JSON(200, gin.H{"todos": todos})
}
func TodoShow(c *gin.Context) {
	var todo models.Todo
	id := c.Param("id")
	initializers.DB.First(&todo, id)

	if todo.ID == 0 {
		c.JSON(404, gin.H{"message": "Todo not found"})
		return
	}
	c.JSON(200, gin.H{"todo": todo})
}
func TodoUpdate(c *gin.Context) {

	id := c.Param("id")
	var body struct {
		Content string
		Status  bool
	}
	c.Bind(&body)

	var todo models.Todo
	initializers.DB.First(&todo, id)
	initializers.DB.Model(&todo).Updates(models.Todo{Content: body.Content, Status: body.Status})
	if initializers.DB.First(&todo, id).RowsAffected == 0 {
		c.JSON(404, gin.H{"message": "Todo not found"})
		return
	}
	c.JSON(200, gin.H{"todo": todo})
}
func TodoDelete(c *gin.Context) {
	id := c.Param("id")
	var todo models.Todo
	initializers.DB.First(&todo, id)

	if todo.ID == 0 {
		c.JSON(404, gin.H{"message": "Todo not found"})
		return
	}
	initializers.DB.Delete(&todo)
	c.JSON(200, gin.H{"message": "Todo deleted successfully"})
}

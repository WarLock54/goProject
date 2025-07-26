package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type todo struct {
	ID        string `json:"id"`
	Item      string `json:"item"`
	Completed bool   `json:"completed"`
}

var todos = []todo{
	{ID: "1", Item: "Learn Go", Completed: false},
	{ID: "2", Item: "Build an API", Completed: false},
	{ID: "3", Item: "Test the API", Completed: false},
}

// basic CRUD operations for a todo list API using Gin framework
func main() {
	r := gin.Default()
	r.GET("/todos/:id", getTodo)
	r.GET("/todos", getTodos)
	r.POST("/todos", addTodo)
	r.PUT("/todos/:id", updateTodo)
	r.DELETE("/todos/:id", deleteTodo)
	if err := r.Run(":9000"); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
func createTodo(todoItem todo) (todo, error) {
	if todoItem.Item == "" {
		return todo{}, errors.New("item cannot be empty")
	}
	todoItem.ID = fmt.Sprintf("%d", len(todos)+1)
	todos = append(todos, todoItem)
	return todoItem, nil
}
func addTodo(context *gin.Context) {
	var newTodo = []todo{}
	err := context.BindJSON(&newTodo)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	todos = append(todos, newTodo...)
	context.IndentedJSON(http.StatusCreated, newTodo)
}
func getTodos(context *gin.Context) {
	context.IndentedJSON(http.StatusOK, todos)
}
func getTodo(context *gin.Context) {
	id := context.Param("id")
	todo, err := getTodoIndex(id)
	if err != nil {
		context.IndentedJSON(http.StatusNotFound, gin.H{"error": err.Error()})

	}
	context.IndentedJSON(http.StatusOK, todos[todo])
}
func getTodoIndex(id string) (int, error) {

	for i, t := range todos {
		if t.ID == id {
			return i, nil
		}
	}
	return -1, errors.New("todo not found")
}
func getTodos2(c *gin.Context) {
	c.JSON(http.StatusOK, todos)
}

func updateTodo(context *gin.Context) {
	id := context.Param("id")

	todo, err := getTodoIndex(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if err := context.BindJSON(&todos[todo]); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	context.IndentedJSON(http.StatusOK, todos[todo])
}
func deleteTodo(context *gin.Context) {
	id := context.Param("id")

	index, err := getTodoIndex(id)
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	todos = append(todos[:index], todos[index+1:]...)

	context.IndentedJSON(http.StatusOK, gin.H{"message": "Todo deleted successfully"})
}

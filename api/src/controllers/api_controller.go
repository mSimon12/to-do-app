package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func createTask(c *gin.Context) {
	c.String(http.StatusOK, "Creating Task!")
}

func getTask(c *gin.Context) {
	c.String(http.StatusOK, "Getting Task!")
}

func updateTask(c *gin.Context) {
	c.String(http.StatusOK, "Updating Task!")
}

func deleteTask(c *gin.Context) {
	c.String(http.StatusOK, "Deleting Task!")
}

func StartAPI() {
	r := gin.Default()
	r.POST("api/tasks", createTask)
	r.GET("api/tasks", getTask)
	r.PUT("api/tasks", updateTask)
	r.DELETE("api/tasks", deleteTask)

	r.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

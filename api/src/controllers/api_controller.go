package controllers

import (
	"github.com/gin-gonic/gin"
)

func StartAPI() {
	router := gin.Default()

	// General endpoints
	router.POST("api/tasks", createTask)
	router.GET("api/tasks", getTasksList)

	// Task based endpoints
	router.GET("api/tasks/:taskId", getTask)
	router.PUT("api/tasks/:taskId", updateTask)
	router.DELETE("api/tasks/:taskId", deleteTask)

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

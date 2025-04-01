package controllers

import (
	_ "to-do-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	router.GET("/api/documentation/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.Run() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}

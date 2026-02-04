package controllers

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func StartAPI() {
	router := gin.Default()

	// Configure CORS
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"http://localhost:4200", "http://localhost:3000"}
	config.AllowMethods = []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}
	config.AllowHeaders = []string{"Content-Type", "Authorization"}
	config.AllowCredentials = true
	router.Use(cors.New(config))

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

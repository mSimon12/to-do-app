package controllers

import (
	"fmt"
	"net/http"
	"to-do-api/service"

	"github.com/gin-gonic/gin"
)

func createTask(c *gin.Context) {
	var requestBody service.TaskRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO:  Create new task and process possible errors
	taskId := service.CreateNewTask(requestBody)

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "taskId": taskId})
}

func getTasksList(c *gin.Context) {
	tasks := map[uint]string{1: "task 1", 2: "task 2"}

	// TODO: Get tasks list and process possible errors

	c.JSON(http.StatusOK, gin.H{"message": "Tasks queried successfully", "tasks": tasks})
}

func getTask(c *gin.Context) {
	taskId := c.Param("taskId")

	// TODO: Get task and process possible errors

	c.JSON(http.StatusOK, gin.H{"message": "Task queried successfully", "task": taskId})
}

func updateTask(c *gin.Context) {
	taskId := c.Param("taskId")

	fmt.Println(taskId)
	// TODO: Update task and process possible errors

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func deleteTask(c *gin.Context) {
	taskId := c.Param("taskId")

	fmt.Println(taskId)
	// TODO: Delete task and process possible errors

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})

}

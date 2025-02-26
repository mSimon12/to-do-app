package controllers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

type taskRequestBody struct {
	Title       string `json:"title" binding:"required"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

func createTask(c *gin.Context) {
	var requestBody taskRequestBody
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		// Respond with error
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// TODO:  Create new task and process possible errors

	// Respond with success
	c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "task": requestBody})
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

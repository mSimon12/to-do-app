package controllers

import (
	"fmt"
	"net/http"
	"to-do-api/service"

	"github.com/gin-gonic/gin"
)

func createTask(c *gin.Context) {
	var requestBody service.TaskRequestBody
	var err error
	if err = c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var taskId uint
	if taskId, err = service.CreateNewTask(requestBody); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

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
	taskIdString := c.Param("taskId")

	taskId, err := service.ValidateTaskId(taskIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = service.DeleteTask(uint(taskId)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})

}

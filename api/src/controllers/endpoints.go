package controllers

import (
	"errors"
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

	if err = service.ValidateNewTaskInput(requestBody); err != nil {
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

func getTask(c *gin.Context) {
	taskIdString := c.Param("taskId")
	taskId, err := service.ValidateTaskIdInput(taskIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, err := service.GetTaskById(taskId)
	if err != nil {
		if errors.Is(err, service.ErrRowNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrDatabaseGeneral) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task queried successfully", "task": task})
}

func updateTask(c *gin.Context) {
	var err error
	taskIdString := c.Param("taskId")
	taskId, err := service.ValidateTaskIdInput(taskIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var requestBody service.TaskRequestBody
	if err = c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = service.ValidateUpdateTaskInput(requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = service.UpdateTask(uint(taskId), requestBody); err != nil {
		if errors.Is(err, service.ErrRowNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrDatabaseGeneral) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	// c.JSON(http.StatusCreated, gin.H{"message": "Task created successfully", "taskId": taskId})

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

func deleteTask(c *gin.Context) {
	taskIdString := c.Param("taskId")
	taskId, err := service.ValidateTaskIdInput(taskIdString)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err = service.DeleteTask(uint(taskId)); err != nil {
		if errors.Is(err, service.ErrRowNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrDatabaseGeneral) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task deleted successfully"})

}

// Endpoint for querying tasks list
func getTasksList(c *gin.Context) {

	// Filtering
	titleFilter := c.Query("title_contains")
	descriptionFilter := c.Query("description_contains")
	statusFilter := c.Query("status")
	priorityFilter := c.Query("priority")

	filtersConfig, err := service.CreateFilterConfig(titleFilter, descriptionFilter, statusFilter, priorityFilter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Pagination
	offset := c.Query("offset")
	limit := c.Query("limit")
	sortBy := c.Query("sort_by")
	sortOrder := c.Query("sort_order")

	pageConfig, err := service.CreatePageConfig(offset, limit, sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tasks, err := service.GetTasksList(filtersConfig, pageConfig)

	if err != nil {
		if errors.Is(err, service.ErrRowNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if errors.Is(err, service.ErrDatabaseGeneral) {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}

		return
	}

	pagination, sorting := service.GetReturnInfo(pageConfig)
	c.JSON(http.StatusOK, gin.H{"message": "Tasks queried successfully", "data": tasks, "pagination": pagination, "sorting": sorting})
}

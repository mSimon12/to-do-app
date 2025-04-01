package controllers

import (
	"errors"
	"net/http"
	"to-do-api/service"

	"github.com/gin-gonic/gin"
)

// CreateTask Creates a new task
//
//	@Summary		Create a new task
//	@Description	Adds a new task to the To-Do List
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			task	body		service.TaskRequestBody	true	"Task data"
//	@Success		201		{object}	map[string]interface{}	"Task created successfully"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/api/tasks [post]
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

// GetTask Retrieves a single task by ID
//
//	@Summary		Get a task by ID
//	@Description	Fetch a specific task from the To-Do List
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			taskId	path		int						true	"Task ID"
//	@Success		200		{object}	map[string]interface{}	"Task retrieved successfully"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		404		{object}	map[string]interface{}	"Task not found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/api/tasks/{taskId} [get]
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

	c.JSON(http.StatusOK, gin.H{"message": "Task retrieved successfully", "task": task})
}

// UpdateTask Updates an existing task
//
//	@Summary		Update a task
//	@Description	Modifies an existing task in the To-Do List
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			taskId	path		int						true	"Task ID"
//	@Param			task	body		service.TaskRequestBody	true	"Updated task data"
//	@Success		200		{object}	map[string]interface{}	"Task updated successfully"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		404		{object}	map[string]interface{}	"Task not found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/api/tasks/{taskId} [put]
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

	c.JSON(http.StatusOK, gin.H{"message": "Task updated successfully"})
}

// DeleteTask Deletes a task by ID
//
//	@Summary		Delete a task
//	@Description	Removes a task from the To-Do List
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			taskId	path		int						true	"Task ID"
//	@Success		200		{object}	map[string]interface{}	"Task deleted successfully"
//	@Failure		400		{object}	map[string]interface{}	"Bad request"
//	@Failure		404		{object}	map[string]interface{}	"Task not found"
//	@Failure		500		{object}	map[string]interface{}	"Internal server error"
//	@Router			/api/tasks/{taskId} [delete]
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

// ListTasks Getting tasks in the To-Do List
//
//	@Summary		Getting tasks in the To-Do List
//	@Description	Get all tasks from To-Do List that match filter conditions
//	@Tags			Tasks
//	@Accept			json
//	@Produce		json
//	@Param			offset					query		int						false	"Pagination offset (default: 0)"
//	@Param			limit					query		int						false	"Pagination limit (default: 10)"
//	@Param			sort_by					query		string					false	"Sort by field (e.g., 'title', 'description')"
//	@Param			sort_order				query		string					false	"Sort order (ASC or DESC)"
//	@Param			title_contains			query		string					false	"Filter by title (substring match)"
//	@Param			description_contains	query		string					false	"Filter by description (substring match)"
//	@Param			status					query		string					false	"Filter by task status"
//	@Param			priority				query		string					false	"Filter by task priority"
//	@Success		200						{object}	map[string]interface{}	"Successful response"
//	@Failure		400						{object}	map[string]interface{}	"Bad request"
//	@Failure		404						{object}	map[string]interface{}	"Tasks not found"
//	@Failure		500						{object}	map[string]interface{}	"Internal server error"
//	@Router			/api/tasks [get]
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

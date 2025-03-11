package controllers

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"to-do-api/service"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestUpdateTask(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: "new title"}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "1"}} // Proper parameter setup

	// Mock internal functions
	monkey.Patch(service.UpdateTask, func(taskId uint, task service.TaskRequestBody) error {
		return nil
	})
	defer monkey.UnpatchAll()

	// Call the handler
	updateTask(context)

	// Validate response
	expectedResponse := "{\"message\":\"Task updated successfully\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusOK, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestUpdateTaskInvalidId(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: "new title"}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "a"}} // Proper parameter setup

	// Call the handler
	updateTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"invalid task id\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestUpdateTaskInexistentId(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: "new title"}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "10"}} // Proper parameter setup

	// Mock internal functions
	monkey.Patch(service.UpdateTask, func(taskId uint, task service.TaskRequestBody) error {
		return service.ErrRowNotFound
	})
	defer monkey.UnpatchAll()

	// Call the handler
	updateTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"requested resource not found on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusNotFound, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestUpdateTaskServerError(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: "new title"}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "10"}} // Proper parameter setup

	// Mock internal functions
	monkey.Patch(service.UpdateTask, func(taskId uint, task service.TaskRequestBody) error {
		return service.ErrDatabaseGeneral
	})
	defer monkey.UnpatchAll()

	// Call the handler
	updateTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"fail processing request on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestUpdateTaskNoInfo(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "10"}} // Proper parameter setup
	context.Request = nil

	// Call the handler
	updateTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"invalid request\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
	"to-do-api/models"
	"to-do-api/service"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func TestGetTask(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "1"}} // Proper parameter setup

	// Mocking the GetTaskById function
	createdAt := time.Now()
	dueDate := createdAt.AddDate(0, 0, 7)
	expectedTask := models.Task{
		Id:          1,
		Title:       "test",
		Description: "nothing",
		Status:      "done",
		Priority:    5,
		CreatedAt:   createdAt,
		DueDate:     dueDate,
	}

	monkey.Patch(service.GetTaskById, func(taskId uint) (models.Task, error) {
		return expectedTask, nil
	})
	defer monkey.UnpatchAll()

	// Call the handler function
	getTask(context)

	// Read actual response body
	responseBody, _ := io.ReadAll(recorder.Body)

	// Construct expected JSON response
	expectedResponseMap := map[string]interface{}{
		"message": "Task queried successfully",
		"task": map[string]interface{}{
			"Id":          expectedTask.Id,
			"Title":       expectedTask.Title,
			"Description": expectedTask.Description,
			"Status":      expectedTask.Status,
			"Priority":    expectedTask.Priority,
			"CreatedAt":   expectedTask.CreatedAt.Format("2006-01-02T15:04:05.999999999Z"), // Ensure consistent format
			"DueDate":     expectedTask.DueDate.Format("2006-01-02T15:04:05.999999999Z"),
		},
	}

	// Convert expected response to JSON
	expectedResponseBytes, _ := json.Marshal(expectedResponseMap)
	expectedResponseString := string(expectedResponseBytes)

	// Validate response
	assert.Equal(t, http.StatusOK, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponseString, string(responseBody), "Invalid response JSON")
}

func TestGetTaskInvalidId(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "a"}} // Proper parameter setup

	// Call the handler function
	getTask(context)

	// Read actual response body
	responseBody, _ := io.ReadAll(recorder.Body)

	expectedResponse := "{\"error\":\"invalid task id\"}"

	// Validate response
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

func TestGetTaskInexistentId(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "10"}} // Proper parameter setup

	// Mocking the GetTaskById function
	monkey.Patch(service.GetTaskById, func(taskId uint) (models.Task, error) {
		return models.Task{}, service.ErrRowNotFound
	})
	defer monkey.UnpatchAll()

	// Call the handler function
	getTask(context)

	// Read actual response body
	responseBody, _ := io.ReadAll(recorder.Body)

	expectedResponse := "{\"error\":\"requested resource not found on database\"}"

	// Validate response
	assert.Equal(t, http.StatusNotFound, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

func TestGetTaskServerError(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "10"}} // Proper parameter setup

	// Mocking the GetTaskById function
	monkey.Patch(service.GetTaskById, func(taskId uint) (models.Task, error) {
		return models.Task{}, service.ErrDatabaseGeneral
	})
	defer monkey.UnpatchAll()

	// Call the handler function
	getTask(context)

	// Read actual response body
	responseBody, _ := io.ReadAll(recorder.Body)

	expectedResponse := "{\"error\":\"fail processing request on database\"}"

	// Validate response
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

package controllers

import (
	"fmt"
	"io"
	"net/http"
	"testing"
	"to-do-api/service"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestCreateTask(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: "new task"}
	context, recorder := getTestGinContextAndRecorder(requestBody)

	// Mock internal functions
	monkey.Patch(service.CreateNewTask, func(task service.TaskRequestBody) (uint, error) {
		return 1, nil
	})
	defer monkey.UnpatchAll()

	// Call the handler
	createTask(context)

	// Validate response
	expectedResponse := "{\"message\":\"Task created successfully\",\"taskId\":1}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusCreated, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestCreateTaskNoTitle(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: ""}
	context, recorder := getTestGinContextAndRecorder(requestBody)

	// Call the handler
	createTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"title must not be empty\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestCreateTaskNoInfo(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Request = nil

	// Call the handler
	createTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"invalid request\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestCreateTaskServerError(t *testing.T) {
	requestBody := TestTaskRequestBody{Title: "new task"}
	context, recorder := getTestGinContextAndRecorder(requestBody)

	// Mock internal functions
	monkey.Patch(service.CreateNewTask, func(task service.TaskRequestBody) (uint, error) {
		return 1, service.ErrDatabaseGeneral
	})
	defer monkey.UnpatchAll()

	// Call the handler
	createTask(context)

	// Validate response
	expectedResponse := "{\"error\":\"fail processing request on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.Equal(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

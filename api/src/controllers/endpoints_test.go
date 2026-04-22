package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"to-do-api/models"
	"to-do-api/service"

	"bou.ke/monkey"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

type TestTaskRequestBody struct {
	Title       string
	Priority    uint
	Description string
	Status      string
	DueDate     string
}

func getTestGinContextAndRecorder(reqBody TestTaskRequestBody) (*gin.Context, *httptest.ResponseRecorder) {
	gin.SetMode(gin.TestMode)

	w := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(w)

	// Convert requestBody to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Failed to marshal JSON: %v\n", err)
	}

	// Initialize request properly
	ctx.Request = httptest.NewRequest(http.MethodGet, "/tasks", nil)
	ctx.Request.Header = make(http.Header)
	ctx.Request.Body = io.NopCloser(bytes.NewReader(jsonBody)) // Properly set Body

	return ctx, w
}

func TestDeleteTask(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "1"}}

	monkey.Patch(service.DeleteTask, func(taskId uint) error {
		return nil
	})
	defer monkey.UnpatchAll()

	deleteTask(context)

	expectedResponse := "{\"message\":\"Task deleted successfully\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusOK, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestDeleteTaskInvalidId(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "abc"}}

	deleteTask(context)

	expectedResponse := "{\"error\":\"invalid task id\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestDeleteTaskInexistentId(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "10"}}

	monkey.Patch(service.DeleteTask, func(taskId uint) error {
		return service.ErrRowNotFound
	})
	defer monkey.UnpatchAll()

	deleteTask(context)

	expectedResponse := "{\"error\":\"requested resource not found on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusNotFound, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestDeleteTaskServerError(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Params = []gin.Param{{Key: "taskId", Value: "1"}}

	monkey.Patch(service.DeleteTask, func(taskId uint) error {
		return service.ErrDatabaseGeneral
	})
	defer monkey.UnpatchAll()

	deleteTask(context)

	expectedResponse := "{\"error\":\"fail processing request on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response pattern")
}

func TestGetTasksList(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)

	monkey.Patch(service.GetTasksList, func(filterConfig []models.TasksFilterQuery, pageConfig models.TasksPaginationQuery) ([]service.TaskInfo, error) {
		return []service.TaskInfo{}, nil
	})
	monkey.Patch(service.GetReturnInfo, func(pageConfig models.TasksPaginationQuery) (map[string]uint, map[string]string) {
		return map[string]uint{"offset": 0, "limit": 10, "total_tasks": 0},
			map[string]string{"by": "id", "order": "ASC"}
	})
	defer monkey.UnpatchAll()

	getTasksList(context)

	expectedResponse := `{"data":[],"message":"Tasks queried successfully","pagination":{"limit":10,"offset":0,"total_tasks":0},"sorting":{"by":"id","order":"ASC"}}`
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusOK, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

func TestGetTasksListInvalidFilter(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Request = httptest.NewRequest(http.MethodGet, "/tasks?title_contains=@test", nil)

	getTasksList(context)

	expectedResponse := "{\"error\":\"invalid title filter: must be alphanumeric\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

func TestGetTasksListInvalidPageConfig(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)
	context.Request = httptest.NewRequest(http.MethodGet, "/tasks?offset=abc", nil)

	getTasksList(context)

	expectedResponse := "{\"error\":\"invalid 'offset' value, must be int > 0\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusBadRequest, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

func TestGetTasksListNotFound(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)

	monkey.Patch(service.GetTasksList, func(filterConfig []models.TasksFilterQuery, pageConfig models.TasksPaginationQuery) ([]service.TaskInfo, error) {
		return []service.TaskInfo{}, service.ErrRowNotFound
	})
	defer monkey.UnpatchAll()

	getTasksList(context)

	expectedResponse := "{\"error\":\"requested resource not found on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusNotFound, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

func TestGetTasksListServerError(t *testing.T) {
	requestBody := TestTaskRequestBody{}
	context, recorder := getTestGinContextAndRecorder(requestBody)

	monkey.Patch(service.GetTasksList, func(filterConfig []models.TasksFilterQuery, pageConfig models.TasksPaginationQuery) ([]service.TaskInfo, error) {
		return []service.TaskInfo{}, service.ErrDatabaseGeneral
	})
	defer monkey.UnpatchAll()

	getTasksList(context)

	expectedResponse := "{\"error\":\"fail processing request on database\"}"
	responseBody, _ := io.ReadAll(recorder.Body)
	assert.Equal(t, http.StatusInternalServerError, recorder.Code, fmt.Sprintf("Unexpected status code: %d", recorder.Code))
	assert.JSONEq(t, expectedResponse, string(responseBody), "Invalid response JSON")
}

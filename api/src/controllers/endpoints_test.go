package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
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

}

func TestGetTasksList(t *testing.T) {

}

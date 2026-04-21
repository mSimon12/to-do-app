package service

import (
	"database/sql"
	"errors"
	"testing"
	"time"
	"to-do-api/models"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

const testCreatedAt = int64(1770843800)
const testDueDate = int64(1770844785)

// Create New Task test /////////////////////////////////////////////////
func TestCreateNewTask(t *testing.T) {
	mockTaskID := uint(1)

	// Mock models.AddTask function
	monkey.Patch(models.AddTask, func(task models.Task) (uint, error) {
		return mockTaskID, nil
	})
	defer monkey.UnpatchAll()

	title := "Test Task"
	priority := uint(1)
	description := "This is a test task"
	status := "pending"
	dueDate := testDueDate

	taskRequest := TaskRequestBody{
		Title:       &title,
		Priority:    &priority,
		Description: &description,
		Status:      &status,
		DueDate:     &dueDate,
	}

	// Run function
	taskID, err := CreateNewTask(taskRequest)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, mockTaskID, taskID)
}

func TestCreateNewTaskDBError(t *testing.T) {
	// Mock models.AddTask to return an error
	monkey.Patch(models.AddTask, func(task models.Task) (uint, error) {
		return 0, errors.New("database error")
	})
	defer monkey.UnpatchAll()

	title := "Test Task"
	priority := uint(1)
	description := "This is a test task"
	status := "pending"
	dueDate := testDueDate

	taskRequest := TaskRequestBody{
		Title:       &title,
		Priority:    &priority,
		Description: &description,
		Status:      &status,
		DueDate:     &dueDate,
	}

	// Run function
	_, err := CreateNewTask(taskRequest)

	// Assertions
	assert.Equal(t, errors.New("fail processing request on database"), err)
}

// Query Task by Id test /////////////////////////////////////////////////
func TestGetTaskById(t *testing.T) {
	createdAt := time.Unix(testCreatedAt, 0)
	dueDate := time.Unix(testDueDate, 0)

	mockTask := models.Task{
		Id:          1,
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    uint16(5),
		Status:      "done",
		CreatedAt:   createdAt,
		DueDate:     dueDate,
	}

	expectedTask := TaskResponseBody{
		Id:          1,
		Title:       "Test Task",
		Description: "This is a test task",
		Status:      "done",
		Priority:    uint16(5),
		CreatedAt:   testCreatedAt,
		DueDate:     testDueDate,
	}

	// Mock models.CheckExistence function
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return true, nil
	})

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, nil
	})
	defer monkey.UnpatchAll()

	// Run function
	task, err := GetTaskById(1)

	// Assertions
	assert.Nil(t, err)
	assert.Equal(t, expectedTask, task)
}

func TestGetTaskByIdInvalidId(t *testing.T) {

	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return false, sql.ErrNoRows
	})
	defer monkey.UnpatchAll()

	// Run function
	_, err := GetTaskById(1)

	// Assertions
	assert.Equal(t, errors.New("requested resource not found on database"), err)
}

func TestGetTaskByIdServerError(t *testing.T) {

	var mockTask models.Task

	// Mock models.CheckExistence function
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return true, nil
	})

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, sql.ErrTxDone
	})
	defer monkey.UnpatchAll()

	// Run function
	_, err := GetTaskById(1)

	// Assertions
	assert.Equal(t, errors.New("fail processing request on database"), err)
}

// Update Task test /////////////////////////////////////////////////
func TestUpdateTask(t *testing.T) {
	mockCreatedAt := time.Unix(testCreatedAt, 0)
	mockDueDate := time.Unix(testDueDate, 0)

	mockTask := models.Task{
		Id:          1,
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    uint16(5),
		Status:      "done",
		CreatedAt:   mockCreatedAt,
		DueDate:     mockDueDate,
	}

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, nil
	})

	// Mock models.UpdateTask function
	monkey.Patch(models.UpdateTask, func(models.Task) error {
		return nil
	})
	defer monkey.UnpatchAll()

	title := "Update Test Task Title"
	priority := uint(2)
	description := "New  test description"
	status := "done"
	dueDate := testDueDate

	taskRequest := TaskRequestBody{
		Title:       &title,
		Priority:    &priority,
		Description: &description,
		Status:      &status,
		DueDate:     &dueDate,
	}

	// Run function
	err := UpdateTask(1, taskRequest)

	// Assertions
	assert.Nil(t, err)
}

func TestUpdateTaskInvalidId(t *testing.T) {
	var mockTask models.Task

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, sql.ErrNoRows
	})

	defer monkey.UnpatchAll()

	title := "Update Test Task Title"

	taskRequest := TaskRequestBody{
		Title: &title,
	}

	// Run function
	err := UpdateTask(1, taskRequest)

	// Assertions
	assert.Equal(t, errors.New("requested resource not found on database"), err)
}

func TestUpdateTaskQueryServerError(t *testing.T) {
	var mockTask models.Task

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, sql.ErrConnDone
	})

	defer monkey.UnpatchAll()

	var taskRequest TaskRequestBody

	// Run function
	err := UpdateTask(1, taskRequest)

	// Assertions
	assert.Equal(t, errors.New("fail processing request on database"), err)
}

func TestUpdateTaskExecServerError(t *testing.T) {
	var mockTask models.Task

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, nil
	})

	// Mock models.UpdateTask function
	monkey.Patch(models.UpdateTask, func(models.Task) error {
		return sql.ErrTxDone
	})
	defer monkey.UnpatchAll()

	var taskRequest TaskRequestBody

	// Run function
	err := UpdateTask(1, taskRequest)

	// Assertions
	assert.Equal(t, errors.New("fail processing request on database"), err)
}

// Delete Task test /////////////////////////////////////////////////
func TestDeleteTask(t *testing.T) {

	// Mock models.CheckExistence function
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return true, nil
	})

	// Mock models.DeleteTask function
	monkey.Patch(models.DeleteTask, func(taskId uint) error {
		return nil
	})
	defer monkey.UnpatchAll()

	// Run function
	err := DeleteTask(1)

	// Assertions
	assert.Nil(t, err)
}

func TestDeleteTaskInvalidId(t *testing.T) {

	// Mock models.CheckExistence function
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return false, nil
	})

	defer monkey.UnpatchAll()

	// Run function
	err := DeleteTask(1)

	// Assertions
	assert.Equal(t, errors.New("requested resource not found on database"), err)
}

func TestDeleteTaskServerError(t *testing.T) {

	// Mock models.CheckExistence function
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return true, nil
	})

	// Mock models.DeleteTask function
	monkey.Patch(models.DeleteTask, func(taskId uint) error {
		return sql.ErrTxDone
	})
	defer monkey.UnpatchAll()

	// Run function
	err := DeleteTask(1)

	// Assertions
	assert.Equal(t, errors.New("fail processing request on database"), err)
}

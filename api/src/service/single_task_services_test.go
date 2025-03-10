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

// Data conversion test /////////////////////////////////////////////////
func TestDateStrToTime(t *testing.T) {

	stringDate := "2024-10-26"

	// Run function
	timeDate, _ := dateStrToTime(stringDate)

	// Assertions
	expectedDate, _ := time.Parse(time.DateOnly, stringDate)
	assert.Equal(t, expectedDate, timeDate, "Invalid date conversion")
}

func TestDateStrToTimeInvalidMonth(t *testing.T) {

	stringDate := "2024-15-26"

	// Run function
	_, err := dateStrToTime(stringDate)

	// Assertions
	expectedErr := time.ParseError{
		Layout:     time.DateOnly,
		Value:      stringDate,
		LayoutElem: "01",
		ValueElem:  "-26",
		Message:    ": month out of range",
	}

	assert.Equal(t, &expectedErr, err, "Allowed Invalid Month")
}

func TestDateStrToTimeInvalidDay(t *testing.T) {

	stringDate := "2024-01-38"

	// Run function
	_, err := dateStrToTime(stringDate)

	// Assertions
	expectedErr := time.ParseError{
		Layout:     time.DateOnly,
		Value:      stringDate,
		LayoutElem: "",
		ValueElem:  "",
		Message:    ": day out of range",
	}

	assert.Equal(t, &expectedErr, err, "Allowed Invalid Month")
}

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
	dueDate := "2025-03-10"

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
	dueDate := "2025-03-10"

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

func TestCreateNewTaskInvalidDueDate(t *testing.T) {

	title := "Test Task"
	priority := uint(1)
	description := "This is a test task"
	status := "pending"
	dueDate := "2025-13-10"

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
	assert.Equal(t, errors.New("invalid due_date format, expects: 'yyyy-mm-dd'"), err)
}

// Query Task by Id test /////////////////////////////////////////////////
func TestGetTaskById(t *testing.T) {
	createdAt, _ := dateStrToTime("2025-03-01")
	dueDate, _ := dateStrToTime("2025-03-10")

	mockTask := models.Task{
		Id:          1,
		Title:       "Test Task",
		Description: "This is a test task",
		Priority:    uint16(5),
		Status:      "done",
		CreatedAt:   createdAt,
		DueDate:     dueDate,
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
	assert.Equal(t, mockTask, task)
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
	mockCreatedAt, _ := dateStrToTime("2025-03-01")
	mockDueDate, _ := dateStrToTime("2025-03-10")

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
	dueDate := "2025-03-11"

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

func TestUpdateTaskInvalidDueDate(t *testing.T) {
	var mockTask models.Task

	// Mock models.QueryTask function
	monkey.Patch(models.QueryTask, func(taskId uint) (models.Task, error) {
		return mockTask, nil
	})

	defer monkey.UnpatchAll()

	dueDate := "invalid date"

	taskRequest := TaskRequestBody{
		DueDate: &dueDate,
	}

	// Run function
	err := UpdateTask(1, taskRequest)

	// Assertions
	assert.Equal(t, errors.New("invalid due_date format, expects: 'yyyy-mm-dd'"), err)
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

package service

import (
	"database/sql"
	"errors"
	"testing"
	"to-do-api/models"

	"bou.ke/monkey"
	"github.com/stretchr/testify/assert"
)

func TestValidateNewTaskInput(t *testing.T) {
	var err error

	// Check valid Info
	title := "New Task"
	desc := "New description"
	prio := uint(8)
	status := "done"
	dueDate := "2025-10-03"

	err = ValidateNewTaskInput(TaskRequestBody{Title: &title, Description: &desc, Priority: &prio, Status: &status, DueDate: &dueDate})
	assert.Nil(t, err, "Should return no error for input containing all info")

	err = ValidateNewTaskInput(TaskRequestBody{Title: &title})
	assert.Nil(t, err, "Should return no error for input containing only Title")

	// Check invalid Info
	err = ValidateNewTaskInput(TaskRequestBody{Description: &desc, Priority: &prio, Status: &status, DueDate: &dueDate})
	assert.Equal(t, errors.New("missing required field: 'title'"), err, "Should return Error for missing title")

	invalidTitle := ""
	err = ValidateNewTaskInput(TaskRequestBody{Title: &invalidTitle})
	assert.Equal(t, errors.New("title must not be empty"), err, "Should return Error for missing title")
}

func TestValidateUpdateTaskInput(t *testing.T) {
	var err error

	// Check valid Info
	title := "New Task"
	desc := "New description"
	prio := uint(8)
	status := "done"
	dueDate := "2025-10-03"

	err = ValidateUpdateTaskInput(TaskRequestBody{Title: &title})
	assert.Nil(t, err, "Should return no error for input containing Title")

	err = ValidateUpdateTaskInput(TaskRequestBody{Description: &desc})
	assert.Nil(t, err, "Should return no error for input containing Description")

	err = ValidateUpdateTaskInput(TaskRequestBody{Priority: &prio})
	assert.Nil(t, err, "Should return no error for input containing Priority")

	err = ValidateUpdateTaskInput(TaskRequestBody{Status: &status})
	assert.Nil(t, err, "Should return no error for input containing Status")

	err = ValidateUpdateTaskInput(TaskRequestBody{DueDate: &dueDate})
	assert.Nil(t, err, "Should return no error for input containing DueDate")

	err = ValidateUpdateTaskInput(TaskRequestBody{Title: &title, Description: &desc, Priority: &prio, Status: &status, DueDate: &dueDate})
	assert.Nil(t, err, "Should return no error for input containing all info")

	// Check invalid Info
	err = ValidateUpdateTaskInput(TaskRequestBody{})
	assert.Equal(t, errors.New("at least one field must be present: 'title', 'description', 'priority', 'status, 'due_date'"), err, "Should return Error for invalid update input")

}

func TestValidateTaskIdInput(t *testing.T) {
	var err error

	// Check valid Info
	config, err := ValidateTaskIdInput("5")
	assert.Nil(t, err, "Should return no error for input '5'")
	assert.Equal(t, uint(5), config, "Should return config as uint(5)")

	// Check invalid Info
	_, err = ValidateTaskIdInput("alpha")
	assert.Equal(t, errors.New("invalid task id"), err, "Should return Error for input 'alpha'")

	_, err = ValidateTaskIdInput("?/_")
	assert.Equal(t, errors.New("invalid task id"), err, "Should return Error for input symbols")

	_, err = ValidateTaskIdInput("a1l2p3h4a")
	assert.Equal(t, errors.New("invalid task id"), err, "Should return Error for alphanumeric input")

	_, err = ValidateTaskIdInput("-10")
	assert.Equal(t, errors.New("invalid task id"), err, "Should return Error for negative input")
}

func TestCheckIdExistValidId(t *testing.T) {
	// Mock models.AddTask to return an error
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return true, nil
	})
	defer monkey.UnpatchAll()

	idExist, _ := checkIdExist(uint(2))
	assert.True(t, idExist, "Should return True for valid Id")
}

func TestCheckIdExistInvalidId(t *testing.T) {
	// Mock models.AddTask to return an error
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return false, nil
	})
	defer monkey.UnpatchAll()

	idExist, _ := checkIdExist(uint(2))
	assert.False(t, idExist, "Should return False for invalid Id")
}

func TestCheckIdExistDBError(t *testing.T) {
	// Mock models.AddTask to return an error
	monkey.Patch(models.CheckExistence, func(taskId uint) (bool, error) {
		return false, sql.ErrConnDone
	})
	defer monkey.UnpatchAll()

	_, err := checkIdExist(uint(2))
	assert.Equal(t, ErrDatabaseGeneral, err, "Should return Error for problems in DB")
}

func TestIsValidPageConfig(t *testing.T) {
	var isValid bool

	// Check valid Info
	config, isValid := isValidPageConfig("5")
	assert.True(t, isValid, "Should return True for input '5'")
	assert.Equal(t, uint(5), config, "Should return config as uint(5)")

	// Check invalid Info
	_, isValid = isValidPageConfig("alpha")
	assert.False(t, isValid, "Should return False for input 'alpha'")

	_, isValid = isValidPageConfig("?/_")
	assert.False(t, isValid, "Should return False for input symbols")

	_, isValid = isValidPageConfig("a1l2p3h4a")
	assert.False(t, isValid, "Should return False for alphanumeric input")

	_, isValid = isValidPageConfig("-10")
	assert.False(t, isValid, "Should return False for negative input")
}

func TestIsValidSortCriteria(t *testing.T) {
	var isValid bool

	// Check valid Info
	isValid = isValidSortCriteria("id")
	assert.True(t, isValid, "Should return True for input 'id'")

	isValid = isValidSortCriteria("Title")
	assert.True(t, isValid, "Should return True for input 'Title'")

	isValid = isValidSortCriteria("STATUS")
	assert.True(t, isValid, "Should return True for input 'STATUS'")

	isValid = isValidSortCriteria("priority")
	assert.True(t, isValid, "Should return True for input 'priority'")

	isValid = isValidSortCriteria("created_at")
	assert.True(t, isValid, "Should return True for input 'created_at'")

	isValid = isValidSortCriteria("due_date")
	assert.True(t, isValid, "Should return True for input 'due_date'")

	// Check invalid Info
	isValid = isValidSortOrder("name")
	assert.False(t, isValid, "Should return False for input 'name'")

	isValid = isValidSortOrder("PRIO")
	assert.False(t, isValid, "Should return False for input 'PRIO'")

	isValid = isValidSortOrder("CreatedAt")
	assert.False(t, isValid, "Should return False for input 'CreatedAt'")

	isValid = isValidSortOrder("DueDate")
	assert.False(t, isValid, "Should return False for input 'DueDate'")

	isValid = isValidSortCriteria("")
	assert.False(t, isValid, "Should return False for empty input ")
}

func TestIsValidSortOrder(t *testing.T) {
	var isValid bool

	// Check valid Info
	isValid = isValidSortOrder("asc")
	assert.True(t, isValid, "Should return True for input 'asc'")

	isValid = isValidSortOrder("Asc")
	assert.True(t, isValid, "Should return True for input 'Asc'")

	isValid = isValidSortOrder("DeSc")
	assert.True(t, isValid, "Should return True for input 'DeSc'")

	isValid = isValidSortOrder("DESC")
	assert.True(t, isValid, "Should return True for input 'DESC'")

	// Check invalid Info
	isValid = isValidSortOrder("ascending")
	assert.False(t, isValid, "Should return False for input 'ascending'")

	isValid = isValidSortOrder("Descending")
	assert.False(t, isValid, "Should return False for input 'Descending'")

	isValid = isValidSortOrder("growing")
	assert.False(t, isValid, "Should return False for input 'growing'")

	isValid = isValidSortOrder("")
	assert.False(t, isValid, "Should return False for empty input ")

}

func TestIsValidTextFilter(t *testing.T) {
	var isValid bool

	// Check valid Info
	isValid = isValidTextFilter("545111")
	assert.True(t, isValid, "Should return True for input with only numbers")

	isValid = isValidTextFilter("test Text")
	assert.True(t, isValid, "Should return True for input with only letters")

	isValid = isValidTextFilter("this is 1 text for with 37 characters")
	assert.True(t, isValid, "Should return True for alphanumeric input")

	isValid = isValidTextFilter("Text with comma, points .?! and -, as well as _")
	assert.True(t, isValid, "Should return True for the provided symbols")

	// Check invalid Info
	isValid = isValidTextFilter("Invalid element: § $ % & @ () {}")
	assert.False(t, isValid, "Should return False for inputs containing § $ % & @ () {}")

	isValid = isValidTextFilter("")
	assert.False(t, isValid, "Should return False for empty inputs")
}

func TestIsValidPriorityFilter(t *testing.T) {
	var isValid bool

	// Check valid Info
	isValid = isValidPriorityFilter("5")
	assert.True(t, isValid, "Should return True for input '5'")

	// Check invalid Info
	isValid = isValidPriorityFilter("alpha")
	assert.False(t, isValid, "Should return False for input 'alpha'")

	isValid = isValidPriorityFilter("?/_")
	assert.False(t, isValid, "Should return False for input symbols")

	isValid = isValidPriorityFilter("a1l2p3h4a")
	assert.False(t, isValid, "Should return False for alphanumeric input")

	isValid = isValidPriorityFilter("-10")
	assert.False(t, isValid, "Should return False for negative input")

}

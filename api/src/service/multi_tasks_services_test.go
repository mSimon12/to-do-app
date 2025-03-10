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

func getTestTasksList() []models.Task {
	layout := "2006-01-02"
	testCreatedAt, _ := time.Parse(layout, "2025-02-03")
	testDueDate, _ := time.Parse(layout, "2025-02-10")
	laterDueDate, _ := time.Parse(layout, "2025-03-10")

	testTasks := []models.Task{
		{
			Id:          1,
			Title:       "First Task",
			Description: "Example of task",
			Status:      "pending",
			Priority:    uint16(1),
			CreatedAt:   testCreatedAt,
			DueDate:     testDueDate,
		},
		{
			Id:          2,
			Title:       "Second Task",
			Description: "Optional description",
			Status:      "done",
			Priority:    uint16(5),
			CreatedAt:   testCreatedAt,
			DueDate:     laterDueDate,
		},
		{
			Id:          3,
			Title:       "Third Task",
			Description: "Another description with 'task'",
			Status:      "pending",
			Priority:    uint16(2),
			CreatedAt:   testCreatedAt,
			DueDate:     testDueDate,
		},
	}
	return testTasks
}

// Create New Task test /////////////////////////////////////////////////
func TestGetTasksList(t *testing.T) {

	testList := getTestTasksList()

	// Mock models.QueryTasks function
	monkey.Patch(models.QueryTasks, func(filterConfig []models.TasksFilterQuery, pageConfig models.TasksPaginationQuery) ([]models.Task, error) {
		return testList, nil
	})
	defer monkey.UnpatchAll()

	pageConfig := models.TasksPaginationQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []models.TasksFilterQuery{
		{Query: "priority= $1 ", Value: "1"},
	}

	// Run function
	taskList, err := GetTasksList(filterConfig, pageConfig)

	// Assertions
	expectedOutput := map[int]taskInfo{
		0: {taskList[0].Id, taskList[0].Title},
		1: {taskList[1].Id, taskList[1].Title},
		2: {taskList[2].Id, taskList[2].Title},
	}

	assert.Nil(t, err)
	assert.Equal(t, expectedOutput, taskList)
}

func TestGetTasksListNoRows(t *testing.T) {

	testList := getTestTasksList()

	// Mock models.QueryTasks function
	monkey.Patch(models.QueryTasks, func(filterConfig []models.TasksFilterQuery, pageConfig models.TasksPaginationQuery) ([]models.Task, error) {
		return testList, sql.ErrNoRows
	})
	defer monkey.UnpatchAll()

	pageConfig := models.TasksPaginationQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []models.TasksFilterQuery{
		{Query: "priority= $1 ", Value: "1"},
	}

	// Run function
	_, err := GetTasksList(filterConfig, pageConfig)

	assert.Equal(t, errors.New("requested resource not found on database"), err)
}

func TestGetTasksListServerError(t *testing.T) {

	testList := getTestTasksList()

	// Mock models.QueryTasks function
	monkey.Patch(models.QueryTasks, func(filterConfig []models.TasksFilterQuery, pageConfig models.TasksPaginationQuery) ([]models.Task, error) {
		return testList, sql.ErrConnDone
	})
	defer monkey.UnpatchAll()

	pageConfig := models.TasksPaginationQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []models.TasksFilterQuery{
		{Query: "priority= $1 ", Value: "1"},
	}

	// Run function
	_, err := GetTasksList(filterConfig, pageConfig)

	assert.Equal(t, errors.New("fail processing request on database"), err)
}

// Append Filter test /////////////////////////////////////////////////
func TestAppendFilters(t *testing.T) {
	filterConfig := []models.TasksFilterQuery{}

	assert.Len(t, filterConfig, 0)

	expectedTitleFilter := models.TasksFilterQuery{Query: "title LIKE $1", Value: "%title_value%"}
	expectedDescFilter := models.TasksFilterQuery{Query: "description LIKE $2", Value: "%description_value%"}
	expectedStatusFilter := models.TasksFilterQuery{Query: "status = $3", Value: "status_value"}
	expectedPriorityFilter := models.TasksFilterQuery{Query: "priority = $4", Value: "priority_value"}

	filterConfig = appendFilter(filterConfig, "title", "title_value")
	assert.Len(t, filterConfig, 1, "Wrong length for filter config")
	assert.Equal(t, expectedTitleFilter, filterConfig[0], "Fail appending title filter")

	filterConfig = appendFilter(filterConfig, "description", "description_value")
	assert.Len(t, filterConfig, 2, "Wrong length for filter config")
	assert.Equal(t, expectedDescFilter, filterConfig[1], "Fail appending title filter")

	filterConfig = appendFilter(filterConfig, "status", "status_value")
	assert.Len(t, filterConfig, 3, "Wrong length for filter config")
	assert.Equal(t, expectedStatusFilter, filterConfig[2], "Fail appending title filter")

	filterConfig = appendFilter(filterConfig, "priority", "priority_value")
	assert.Len(t, filterConfig, 4, "Wrong length for filter config")
	assert.Equal(t, expectedPriorityFilter, filterConfig[3], "Fail appending title filter")

}

func TestAppendFiltersInvalidType(t *testing.T) {
	filterConfig := []models.TasksFilterQuery{}

	assert.Len(t, filterConfig, 0, "Wrong length for filter config")

	filterConfig = appendFilter(filterConfig, "dueDate", "DueDate_value")
	assert.Len(t, filterConfig, 0, "Invalid filter appended")

}

// Create Filter Config test /////////////////////////////////////////////////
func TestCreateFilters(t *testing.T) {
	expectedFilterConfig := []models.TasksFilterQuery{
		{Query: "title LIKE $1", Value: "%title_value%"},
		{Query: "description LIKE $2", Value: "%description_value%"},
		{Query: "status = $3", Value: "status_value"},
		{Query: "priority = $4", Value: "1"},
	}

	filterConfig, err := CreateFilterConfig("title_value", "description_value", "status_value", "1")

	assert.Nil(t, err, "Create Filter returned error")
	assert.Len(t, filterConfig, 4, "Wrong length for filter config")
	assert.Equal(t, expectedFilterConfig, filterConfig, "Fail creating filter")

}

func TestCreateFiltersInvalidTitle(t *testing.T) {
	_, err := CreateFilterConfig("{}", "description_value", "status_value", "1")

	assert.Equal(t, errors.New("invalid title filter: must be alphanumeric"), err, "Did not raise error with invalid Title format")
}

func TestCreateFiltersInvalidDescription(t *testing.T) {
	_, err := CreateFilterConfig("", "()", "status_value", "1")

	assert.Equal(t, errors.New("invalid description filter: must be alphanumeric"), err, "Did not raise error with invalid Description format")
}

func TestCreateFiltersInvalidStatus(t *testing.T) {
	_, err := CreateFilterConfig("title_value", "description_value", ">", "1")

	assert.Equal(t, errors.New("invalid status filter: must be alphanumeric"), err, "Did not raise error with invalid Status format")
}

func TestCreateFiltersInvalidPriority(t *testing.T) {
	_, err := CreateFilterConfig("title_value", "description_value", "status_value", "alpha")

	assert.Equal(t, errors.New("invalid priority filter: must be positive integer"), err, "Did not raise error with invalid Priority format")
}

// Create Page Config test /////////////////////////////////////////////////
func TestCreatePageConfig(t *testing.T) {
	expectedPageConfig := models.TasksPaginationQuery{
		Offset:    0,
		SortBy:    "priority",
		SortOrder: "ASC",
		Limit:     10,
	}

	pageConfig, err := CreatePageConfig("0", "10", "priority", "ASC")

	assert.Nil(t, err, "Create Page Config returned error")
	assert.Equal(t, expectedPageConfig, pageConfig, "Fail creating page config")
}

func TestCreatePageConfigInvalidOffset(t *testing.T) {
	_, err := CreatePageConfig("alpha", "10", "priority", "ASC")

	assert.Equal(t, errors.New("invalid 'offset' value, must be int > 0"), err, "Did not raise error with invalid Offset format")
}

func TestCreatePageConfigInvalidLimit(t *testing.T) {
	_, err := CreatePageConfig("0", "alpha", "priority", "ASC")

	assert.Equal(t, errors.New("invalid 'limit' value, must be int > 0"), err, "Did not raise error with invalid Limit format")
}

func TestCreatePageConfigInvalidSortBy(t *testing.T) {
	_, err := CreatePageConfig("0", "10", "author", "ASC")

	assert.Equal(t, errors.New("invalid 'sortBy' value. Valid values: [id title status priority created_at due_date]"), err, "Did not raise error with invalid SortBy format")
}

func TestCreatePageConfigInvalidOrder(t *testing.T) {
	_, err := CreatePageConfig("0", "10", "priority", "UP")

	assert.Equal(t, errors.New("invalid 'sortOrder' value. Valid values: ['asc', 'desc']"), err, "Did not raise error with invalid SortOrder format")
}

// Get Return Info test /////////////////////////////////////////////////
func TestGetReturnInfo(t *testing.T) {
	expectedPaginationInfo := map[string]uint{
		"offset":      0,
		"limit":       10,
		"total_tasks": 5,
	}

	expectedSortingInfo := map[string]string{
		"by":    "priority",
		"order": "ASC",
	}

	// Mock models.AddTask function
	monkey.Patch(models.GetAmountOfTasks, func() (uint, error) {
		return 5, nil
	})
	defer monkey.UnpatchAll()

	pageConfig, err := CreatePageConfig("0", "10", "priority", "ASC")
	assert.Nil(t, err, "Create Page Config returned error")

	paginationInfo, sortingInfo := GetReturnInfo(pageConfig)

	assert.Equal(t, expectedPaginationInfo, paginationInfo, "Returned wrong pagination info")
	assert.Equal(t, expectedSortingInfo, sortingInfo, "Returned wrong sorting info")
}

func TestGetReturnInfoErrorGettingTotal(t *testing.T) {
	expectedPaginationInfo := map[string]uint{
		"offset": 0,
		"limit":  10,
	}

	expectedSortingInfo := map[string]string{
		"by":    "priority",
		"order": "ASC",
	}

	// Mock models.AddTask function
	monkey.Patch(models.GetAmountOfTasks, func() (uint, error) {
		return 5, sql.ErrNoRows
	})
	defer monkey.UnpatchAll()

	pageConfig, err := CreatePageConfig("0", "10", "priority", "ASC")
	assert.Nil(t, err, "Create Page Config returned error")

	paginationInfo, sortingInfo := GetReturnInfo(pageConfig)

	assert.Equal(t, expectedPaginationInfo, paginationInfo, "Returned wrong pagination info")
	assert.Equal(t, expectedSortingInfo, sortingInfo, "Returned wrong sorting info")
}

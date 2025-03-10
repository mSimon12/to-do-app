package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// Multi Tasks Tests ///////////////////////////////////

func getTestTasksList() []Task {
	layout := "2006-01-02"
	testCreatedAt, _ := time.Parse(layout, "2025-02-03")
	testDueDate, _ := time.Parse(layout, "2025-02-10")
	laterDueDate, _ := time.Parse(layout, "2025-03-10")

	testTasks := []Task{
		Task{
			Id:          1,
			Title:       "First Task",
			Description: "Example of task",
			Status:      "pending",
			Priority:    uint16(1),
			CreatedAt:   testCreatedAt,
			DueDate:     testDueDate,
		},
		Task{
			Id:          2,
			Title:       "Second Task",
			Description: "Optional description",
			Status:      "done",
			Priority:    uint16(5),
			CreatedAt:   testCreatedAt,
			DueDate:     laterDueDate,
		},
		Task{
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

func TestQueryTasksListPrioFilterIdSorted(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []TasksFilterQuery{
		{Query: "priority= $1 ", Value: "1"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE priority= \\$1 ORDER BY id ASC LIMIT \\$2 OFFSET \\$3;"
	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[0].Id, testTasks[0].Title, testTasks[0].Description, testTasks[0].Status, testTasks[0].Priority, testTasks[0].CreatedAt, testTasks[0].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 1, len(queriedTask), "Returned list should have only 1 element")
	assert.Equal(t, testTasks[0].Id, queriedTask[0].Id, "Returned Id should be 1")
	assert.Equal(t, testTasks[0].Title, queriedTask[0].Title, "Returned Title should match Task 1")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestQueryTasksListStatusFilterPrioSorted(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 0, SortBy: "priority", SortOrder: "ASC", Limit: 10}
	filterConfig := []TasksFilterQuery{
		{Query: "status= $1 ", Value: "pending"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE status= \\$1 ORDER BY priority ASC LIMIT \\$2 OFFSET \\$3;"

	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[2].Id, testTasks[2].Title, testTasks[2].Description, testTasks[2].Status, testTasks[2].Priority, testTasks[2].CreatedAt, testTasks[2].DueDate)
	expectedReturn.AddRow(testTasks[0].Id, testTasks[0].Title, testTasks[0].Description, testTasks[0].Status, testTasks[0].Priority, testTasks[0].CreatedAt, testTasks[0].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 2, len(queriedTask), "Returned list should have 2 elements")
	assert.Equal(t, testTasks[2].Id, queriedTask[0].Id, "First Id should be from Task 3")
	assert.Equal(t, testTasks[2].Title, queriedTask[0].Title, "First Title should be from Task 3")
	assert.Equal(t, testTasks[0].Id, queriedTask[1].Id, "Second Id should be from Task 1")
	assert.Equal(t, testTasks[0].Title, queriedTask[1].Title, "Second Title should be from Task 1")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestQueryTasksListTitleFilterDueDateSortedDESC(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 0, SortBy: "due_date", SortOrder: "DESC", Limit: 10}
	filterConfig := []TasksFilterQuery{
		{Query: "title= $1 ", Value: "Task"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE title= \\$1 ORDER BY due_date DESC LIMIT \\$2 OFFSET \\$3;"

	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[1].Id, testTasks[1].Title, testTasks[1].Description, testTasks[1].Status, testTasks[1].Priority, testTasks[1].CreatedAt, testTasks[1].DueDate)
	expectedReturn.AddRow(testTasks[0].Id, testTasks[0].Title, testTasks[0].Description, testTasks[0].Status, testTasks[0].Priority, testTasks[0].CreatedAt, testTasks[0].DueDate)
	expectedReturn.AddRow(testTasks[2].Id, testTasks[2].Title, testTasks[2].Description, testTasks[2].Status, testTasks[2].Priority, testTasks[2].CreatedAt, testTasks[2].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 3, len(queriedTask), "Returned list should have 2 elements")
	assert.Equal(t, testTasks[1].Id, queriedTask[0].Id, "First Id should be from Task 2")
	assert.Equal(t, testTasks[1].Title, queriedTask[0].Title, "First Title should be from Task 2")
	assert.Equal(t, testTasks[0].Id, queriedTask[1].Id, "Second Id should be from Task 1")
	assert.Equal(t, testTasks[0].Title, queriedTask[1].Title, "Second Title should be from Task 1")
	assert.Equal(t, testTasks[2].Id, queriedTask[2].Id, "Third Id should be from Task 3")
	assert.Equal(t, testTasks[2].Title, queriedTask[2].Title, "Third Title should be from Task 3")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestQueryTasksListTitleFilterDueDateSortedLimited(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 0, SortBy: "due_date", SortOrder: "ASC", Limit: 2}
	filterConfig := []TasksFilterQuery{
		{Query: "title= $1 ", Value: "Task"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE title= \\$1 ORDER BY due_date ASC LIMIT \\$2 OFFSET \\$3;"

	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[1].Id, testTasks[1].Title, testTasks[1].Description, testTasks[1].Status, testTasks[1].Priority, testTasks[1].CreatedAt, testTasks[1].DueDate)
	expectedReturn.AddRow(testTasks[0].Id, testTasks[0].Title, testTasks[0].Description, testTasks[0].Status, testTasks[0].Priority, testTasks[0].CreatedAt, testTasks[0].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 2, len(queriedTask), "Returned list should have 2 elements")
	assert.Equal(t, testTasks[1].Id, queriedTask[0].Id, "First Id should be from Task 2")
	assert.Equal(t, testTasks[1].Title, queriedTask[0].Title, "First Title should be from Task 2")
	assert.Equal(t, testTasks[0].Id, queriedTask[1].Id, "Second Id should be from Task 1")
	assert.Equal(t, testTasks[0].Title, queriedTask[1].Title, "Second Title should be from Task 1")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestQueryTasksListDescFilterIdSorted(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []TasksFilterQuery{
		{Query: "description= $1 ", Value: "task"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE description= \\$1 ORDER BY id ASC LIMIT \\$2 OFFSET \\$3;"

	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[0].Id, testTasks[0].Title, testTasks[0].Description, testTasks[0].Status, testTasks[0].Priority, testTasks[0].CreatedAt, testTasks[0].DueDate)
	expectedReturn.AddRow(testTasks[2].Id, testTasks[2].Title, testTasks[2].Description, testTasks[2].Status, testTasks[2].Priority, testTasks[2].CreatedAt, testTasks[2].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 2, len(queriedTask), "Returned list should have 2 elements")
	assert.Equal(t, testTasks[0].Id, queriedTask[0].Id, "First Id should be from Task 1")
	assert.Equal(t, testTasks[0].Title, queriedTask[0].Title, "First Title should be from Task 1")
	assert.Equal(t, testTasks[2].Id, queriedTask[1].Id, "Second Id should be from Task 3")
	assert.Equal(t, testTasks[2].Title, queriedTask[1].Title, "Second Title should be from Task 3")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestQueryTasksListDescFilterIdSortedOffset(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 1, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []TasksFilterQuery{
		{Query: "description= $1 ", Value: "task"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE description= \\$1 ORDER BY id ASC LIMIT \\$2 OFFSET \\$3;"

	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[2].Id, testTasks[2].Title, testTasks[2].Description, testTasks[2].Status, testTasks[2].Priority, testTasks[2].CreatedAt, testTasks[2].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 1, len(queriedTask), "Returned list should have 2 elements")
	assert.Equal(t, testTasks[2].Id, queriedTask[0].Id, "Second Id should be from Task 3")
	assert.Equal(t, testTasks[2].Title, queriedTask[0].Title, "Second Title should be from Task 3")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestQueryTasksListTitleFilterStatusFilter(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	testTasks := getTestTasksList()

	pagConfig := TasksPaginationQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}
	filterConfig := []TasksFilterQuery{
		{Query: "title= $1 ", Value: "Task"},
		{Query: "status= $2 ", Value: "done"},
	}

	// Set SQL mock expectation
	expectedQuery := "SELECT \\* FROM tasks WHERE title= \\$1 AND status= \\$2 ORDER BY id ASC LIMIT \\$3 OFFSET \\$4;"

	expectedReturn := pgxmock.NewRows([]string{"id", "title", "description", "status", "priority", "created_at", "due_date"})
	expectedReturn.AddRow(testTasks[1].Id, testTasks[1].Title, testTasks[1].Description, testTasks[1].Status, testTasks[1].Priority, testTasks[1].CreatedAt, testTasks[1].DueDate)

	mockConn.ExpectQuery(expectedQuery).
		WithArgs(filterConfig[0].Value, filterConfig[1].Value, pagConfig.Limit, pagConfig.Offset).
		WillReturnRows(expectedReturn)

	// Run function
	queriedTask, err := QueryTasks(filterConfig, pagConfig)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, 1, len(queriedTask), "Returned list should have 2 elements")
	assert.Equal(t, testTasks[1].Id, queriedTask[0].Id, "Second Id should be from Task 2")
	assert.Equal(t, testTasks[1].Title, queriedTask[0].Title, "Second Title should be from Task 2")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestGetAmountOfTasks(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	// Set SQL mock expectation
	expectedQuery := "SELECT COUNT\\(\\*\\) FROM tasks;"

	mockConn.ExpectQuery(expectedQuery).
		WillReturnRows(pgxmock.NewRows([]string{"count"}).AddRow(uint(3)))

	// Run function
	tasksAmount, err := GetAmountOfTasks()

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, uint(3), tasksAmount, "Returned wrong tasks amount")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")

}

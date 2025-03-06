package models

import (
	"fmt"
	"testing"
	"time"

	"github.com/pashagolub/pgxmock/v4"
	"github.com/stretchr/testify/assert"
)

// Inject mock database connection
func setMockConnection() pgxmock.PgxPoolIface {
	mockConn, _ := pgxmock.NewPool()

	dbConn = mockConn
	return mockConn
}

func TestAddTaskPattern(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	// Define test task
	newTask := Task{
		Title:       "Unit Test Task",
		Description: "Mocked DB test",
		Status:      "pending",
		Priority:    1,
		CreatedAt:   time.Now(),
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	expectedQuery := "INSERT INTO tasks \\(title, description, status, priority, created_at, due_date\\) VALUES \\(\\$1, \\$2, \\$3, \\$4, \\$5, \\$6\\) RETURNING id; "

	// Expect an INSERT query with correct values
	mockConn.ExpectQuery(expectedQuery).
		WithArgs(newTask.Title, newTask.Description, newTask.Status, newTask.Priority, newTask.CreatedAt, newTask.DueDate).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint(1)))

	// Run function
	taskID, err := AddTask(newTask)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, uint(1), taskID, "Returned value should be 1")

	// Ensure all expected queries were executed
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

// func TestDeleteTaskPattern(t *testing.T) {
// 	mockConn := setMockConnection()
// 	defer mockConn.Close()

// 	expectedQuery := "DELETE FROM tasks WHERE id=\\$1; "

// 	// Expect an INSERT query with correct values
// 	mockConn.ExpectQuery(expectedQuery).
// 		WithArgs(uint(1)).
// 		WillReturnResult(psxmock.NewResult(1, 1))

// 	// Run function
// 	taskID, err := AddTask(newTask)

// 	// Assertions
// 	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
// 	assert.Equal(t, uint(1), taskID, "Returned value should be 1")

// 	// Ensure all expected queries were executed
// 	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
// }

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

	// Set SQL mock expectation
	mockConn.ExpectQuery(expectedQuery).
		WithArgs(newTask.Title, newTask.Description, newTask.Status, newTask.Priority, newTask.CreatedAt, newTask.DueDate).
		WillReturnRows(pgxmock.NewRows([]string{"id"}).AddRow(uint(1)))

	// Run function
	taskID, err := AddTask(newTask)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.Equal(t, uint(1), taskID, "Returned value should be 1")
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestUpdateValidTask(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	// Define test task
	updatedTask := Task{
		Id:          1,
		Title:       "Unit Test Task",
		Description: "Mocked DB test",
		Status:      "pending",
		Priority:    1,
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	expectedQuery := "UPDATE tasks SET title = \\$1, description= \\$2, status= \\$3, priority= \\$4, due_date= \\$5 WHERE id = \\$6;"

	// Set SQL mock expectation
	mockConn.ExpectExec(expectedQuery).
		WithArgs(updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Priority, updatedTask.DueDate, updatedTask.Id).
		WillReturnResult(pgxmock.NewResult("UPDATE", 1))

	// Run function
	err := UpdateTask(updatedTask)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestUpdateInvalidTask(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	// Define test task
	updatedTask := Task{
		Id:          1,
		Title:       "Unit Test Task",
		Description: "Mocked DB test",
		Status:      "pending",
		Priority:    1,
		DueDate:     time.Now().Add(24 * time.Hour),
	}

	expectedQuery := "UPDATE tasks SET title = \\$1, description= \\$2, status= \\$3, priority= \\$4, due_date= \\$5 WHERE id = \\$6;"

	// Set SQL mock expectation
	mockConn.ExpectExec(expectedQuery).
		WithArgs(updatedTask.Title, updatedTask.Description, updatedTask.Status, updatedTask.Priority, updatedTask.DueDate, updatedTask.Id).
		WillReturnResult(pgxmock.NewResult("UPDATE", 0))

	// Run function
	err := UpdateTask(updatedTask)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestDeleteValidTask(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	taskId := uint(1)

	expectedQuery := "DELETE FROM tasks WHERE id=\\$1; "

	// Set SQL mock expectation
	mockConn.ExpectExec(expectedQuery).
		WithArgs(taskId).
		WillReturnResult(pgxmock.NewResult("DELETE", 1))

	// Run function
	err := DeleteTask(taskId)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

func TestDeleteInvalidTask(t *testing.T) {
	mockConn := setMockConnection()
	defer mockConn.Close()

	taskId := uint(1)

	expectedQuery := "DELETE FROM tasks WHERE id=\\$1; "

	// Set SQL mock expectation
	mockConn.ExpectExec(expectedQuery).
		WithArgs(taskId).
		WillReturnResult(pgxmock.NewResult("DELETE", 0))

	// Run function
	err := DeleteTask(taskId)

	// Assertions
	assert.NoError(t, err, fmt.Sprintf("Unexpected error from function - err: %v", err))
	assert.NoError(t, mockConn.ExpectationsWereMet(), "Expectations from SQL not met!")
}

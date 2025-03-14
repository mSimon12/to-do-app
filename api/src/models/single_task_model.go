package models

import (
	"context"
	"time"
)

type Task struct {
	Id          uint
	Title       string
	Description string
	Status      string
	Priority    uint16
	CreatedAt   time.Time
	DueDate     time.Time
}

func AddTask(newTask Task) (uint, error) {
	conn := getDatabaseConnection()
	defer conn.Close()

	newTaskQuery := `
		INSERT INTO tasks (title, description, status, priority, created_at, due_date) 
		VALUES ($1, $2, $3, $4, $5, $6) 
		RETURNING id;
	`

	var taskId uint
	err := conn.QueryRow(context.Background(), newTaskQuery,
		newTask.Title,
		newTask.Description,
		newTask.Status,
		newTask.Priority,
		newTask.CreatedAt,
		newTask.DueDate,
	).Scan(&taskId)

	return taskId, err
}

func QueryTask(taskId uint) (Task, error) {
	conn := getDatabaseConnection()
	defer conn.Close()

	var queriedTask Task
	err := conn.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id=$1;", taskId).Scan(
		&queriedTask.Id,
		&queriedTask.Title,
		&queriedTask.Description,
		&queriedTask.Status,
		&queriedTask.Priority,
		&queriedTask.CreatedAt,
		&queriedTask.DueDate)

	return queriedTask, err
}

func UpdateTask(updatedTask Task) error {
	conn := getDatabaseConnection()
	defer conn.Close()

	// Update task from DB
	newTaskQuery := "UPDATE tasks SET title = $1, description= $2, status= $3, priority= $4, due_date= $5 WHERE id = $6;"
	_, err := conn.Exec(context.Background(), newTaskQuery,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.Status,
		updatedTask.Priority,
		updatedTask.DueDate,
		updatedTask.Id)

	return err
}

func DeleteTask(taskId uint) error {
	conn := getDatabaseConnection()
	defer conn.Close()

	// Delete task from DB
	_, err := conn.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1;", taskId)

	return err
}

func CheckExistence(taskId uint) (bool, error) {
	conn := getDatabaseConnection()
	defer conn.Close()

	// Delete task from DB

	var idExist bool
	err := conn.QueryRow(context.Background(), "SELECT EXISTS(SELECT * from tasks WHERE id=$1);", taskId).Scan(&idExist)

	return idExist, err
}

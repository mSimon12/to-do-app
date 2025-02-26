package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
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

func CreateTask(newTask Task) uint16 {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	// Add new task to DB
	newTaskQuery := "INSERT INTO tasks (title, description, status, priority, created_at, due_date) VALUES ($1, $2, $3, $4, $5, $6);"
	_, err := conn.Exec(context.Background(), newTaskQuery,
		newTask.Title,
		newTask.Description,
		newTask.Status,
		newTask.Priority,
		newTask.CreatedAt,
		newTask.DueDate,
	)

	if err != nil {
		panic(err)
	}

	return 1
}

func QueryTask(taskId uint) {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	var queriedTask Task
	err := conn.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id=$1;", taskId).Scan(&queriedTask.Id, &queriedTask.Title,
		&queriedTask.Description,
		&queriedTask.Status,
		&queriedTask.Priority,
		&queriedTask.CreatedAt,
		&queriedTask.DueDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	fmt.Println(queriedTask)
}

func UpdateTask(updatedTask Task) {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	// update task from DB
	newTaskQuery := "UPDATE tasks SET title = $1, description= $2, status= $3, priority= $4, due_date= $5 WHERE id = $6;"
	_, err := conn.Exec(context.Background(), newTaskQuery,
		updatedTask.Title,
		updatedTask.Description,
		updatedTask.Status,
		updatedTask.Priority,
		updatedTask.DueDate,
		updatedTask.Id)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}
}

func DeleteTask(taskId uint) {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	// Delete task from DB
	_, err := conn.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1;", taskId)

	if err != nil {
		panic(err)
	}

}

func getDatabaseConnection() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), getDatabaseUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	return conn
}

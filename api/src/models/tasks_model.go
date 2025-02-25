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
	createdAt   time.Time
	DueDate     time.Time
}

func CreateTask(newTask Task) {
	createdAt := time.Now()

	newTaskQuery := "INSERT INTO tasks (title, description, status, priority, created_at, due_date) VALUES ($1, $2, $3, $4, $5, $6);"

	conn, err := pgx.Connect(context.Background(), getDatabaseUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Add new task to DB
	_, err = conn.Exec(context.Background(), newTaskQuery,
		newTask.Title,
		newTask.Description,
		newTask.Status,
		newTask.Priority,
		createdAt,
		newTask.DueDate,
	)

	if err != nil {
		panic(err)
	}

}

func QueryTask(taskId uint) {
	conn, err := pgx.Connect(context.Background(), getDatabaseUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	var queriedTask Task
	err = conn.QueryRow(context.Background(), "SELECT * FROM tasks WHERE id=$1;", taskId).Scan(&queriedTask.Id, &queriedTask.Title,
		&queriedTask.Description,
		&queriedTask.Status,
		&queriedTask.Priority,
		&queriedTask.createdAt,
		&queriedTask.DueDate)
	if err != nil {
		fmt.Fprintf(os.Stderr, "QueryRow failed: %v\n", err)
	}

	fmt.Println(queriedTask)
}

func DeleteTask(taskId uint) {
	conn, err := pgx.Connect(context.Background(), getDatabaseUrl())
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(context.Background())

	// Delete task from DB
	_, err = conn.Exec(context.Background(), "DELETE FROM tasks WHERE id=$1;", taskId)

	if err != nil {
		panic(err)
	}

}

package models

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

type Task struct {
	Title       string
	Description string
	Status      string
	Priority    uint16
	start_date  time.Time
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

	// Add new task to db
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

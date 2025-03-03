package models

import (
	"context"
	"fmt"
	"log"
)

func QueryTasks(orderBy string) ([]Task, error) {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	taskQuery := fmt.Sprintf("SELECT * FROM tasks ORDER BY %s ASC;", orderBy)

	rows, err := conn.Query(context.Background(), taskQuery)
	tasks := []Task{}

	if err != nil {
		return tasks, err
	}

	for rows.Next() {
		task := Task{}
		err := rows.Scan(&task.Id,
			&task.Title,
			&task.Description,
			&task.Status,
			&task.Priority,
			&task.CreatedAt,
			&task.DueDate)

		if err != nil {
			log.Fatalln(err)
		}
		tasks = append(tasks, task)
	}

	return tasks, err
}

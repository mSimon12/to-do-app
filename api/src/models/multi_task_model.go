package models

import (
	"context"
	"fmt"
	"log"
)

type TasksListQuery struct {
	Offset    uint
	SortBy    string
	SortOrder string
	Limit     uint
}

func QueryTasks(queryConfig TasksListQuery) ([]Task, error) {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	// taskQuery := fmt.Sprintf("SELECT * FROM tasks ORDER BY %s ASC;", orderBy)
	taskQuery := fmt.Sprintf("SELECT * FROM tasks ORDER BY %s %s LIMIT %d OFFSET %d;",
		queryConfig.SortBy,
		queryConfig.SortOrder,
		queryConfig.Limit,
		queryConfig.Offset)

	fmt.Println(taskQuery)

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

func GetAmountOfTasks() (uint, error) {
	conn := getDatabaseConnection()
	defer conn.Close(context.Background())

	var tableSize uint
	err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks;").Scan(&tableSize)

	return tableSize, err
}

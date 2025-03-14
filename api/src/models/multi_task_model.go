package models

import (
	"context"
	"fmt"
	"log"
	"strings"
)

type TasksPaginationQuery struct {
	Offset    uint
	SortBy    string
	SortOrder string
	Limit     uint
}

type TasksFilterQuery struct {
	Query string // columns
	Value string // searched value
}

func QueryTasks(filterConfig []TasksFilterQuery, pageConfig TasksPaginationQuery) ([]Task, error) {
	conn := getDatabaseConnection()
	defer conn.Close()

	var queryBuilder strings.Builder
	var queryParams []interface{} // Slice to store query values

	queryBuilder.WriteString("SELECT * FROM tasks")

	// Add filter queries
	filterElements := len(filterConfig)
	if filterElements > 0 {
		queryBuilder.WriteString(" WHERE ")
	}

	for filter_idx, filter := range filterConfig {
		if filter_idx > 0 {
			queryBuilder.WriteString(" AND ")
		}
		queryBuilder.WriteString(filter.Query)
		queryParams = append(queryParams, filter.Value)
	}

	// Add pagination query
	paginationQuery := fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d;", pageConfig.SortBy, pageConfig.SortOrder, filterElements+1, filterElements+2)
	queryBuilder.WriteString(paginationQuery)
	queryParams = append(queryParams, pageConfig.Limit, pageConfig.Offset) // Add limit and offset

	taskQuery := queryBuilder.String()

	rows, err := conn.Query(context.Background(), taskQuery, queryParams...)
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
	defer conn.Close()

	var tableSize uint
	err := conn.QueryRow(context.Background(), "SELECT COUNT(*) FROM tasks;").Scan(&tableSize)

	return tableSize, err
}

package service

import (
	"database/sql"
	"errors"
	"fmt"
	"to-do-api/models"
)

type taskInfo struct {
	Id    uint
	Title string
}

func GetTasksList(filter string) (map[int]taskInfo, error) {
	orderedTasks := map[int]taskInfo{}

	queriedTasks, err := models.QueryTasks(filter)

	for taskIdx, task := range queriedTasks {
		newTask := taskInfo{Id: task.Id, Title: task.Title}
		orderedTasks[taskIdx] = newTask
	}

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return orderedTasks, ErrRowNotFound
		} else {
			fmt.Printf("Query Task failed: %v\n", err)
			return orderedTasks, ErrDatabaseGeneral
		}
	}

	return orderedTasks, nil

}

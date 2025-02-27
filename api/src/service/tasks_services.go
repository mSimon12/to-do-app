package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"to-do-api/models"
)

type TaskRequestBody struct {
	Title       string `json:"title" binding:"required"`
	Priority    int    `json:"priority"`
	Description string `json:"description"`
	Status      string `json:"status"`
	DueDate     string `json:"due_date"`
}

var ErrDatabaseGeneral = errors.New("fail processing request on database")
var ErrRowNotFound = errors.New("requested resource not found on database")

func dateStrToTime(date string) time.Time {
	layout := "2006-01-22"
	convertedDate, err := time.Parse(layout, date)

	if err != nil {
		fmt.Println(err)
		return time.Now().AddDate(0, 0, 1)
	}

	return convertedDate
}

func CreateNewTask(task TaskRequestBody) (uint, error) {
	newTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    uint16(task.Priority),
		CreatedAt:   time.Now(),
		DueDate:     dateStrToTime(task.DueDate),
	}

	newTaskId, err := models.AddTask(newTask)

	if err != nil {
		fmt.Printf("Create Task failed: %v\n", err)
		return 0, ErrDatabaseGeneral
	}

	return newTaskId, nil
}

func GetTaskById(taskId uint) (models.Task, error) {
	var task models.Task

	if !checkIdExist(taskId) {
		return task, ErrRowNotFound
	}

	task, err := models.QueryTask(taskId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return task, ErrRowNotFound
		} else {
			fmt.Printf("Get Task failed: %v\n", err)
			return task, ErrDatabaseGeneral
		}
	}

	return task, nil
}

func DeleteTask(taskId uint) error {

	if !checkIdExist(taskId) {
		return ErrRowNotFound
	}

	err := models.DeleteTask(taskId)

	if err != nil {
		fmt.Printf("Delete Task failed: %v\n", err)
		return ErrDatabaseGeneral
	}

	return nil
}

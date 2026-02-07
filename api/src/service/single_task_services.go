package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"to-do-api/models"
)

var ErrDatabaseGeneral = errors.New("fail processing request on database")
var ErrRowNotFound = errors.New("requested resource not found on database")
var ErrInvalidInput = errors.New("invalid input")

type TaskRequestBody struct {
	Title       *string `json:"title"`
	Priority    *uint   `json:"priority"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *int64  `json:"due_date"`
}

type TaskResponseBody struct {
	Id          uint
	Title       string
	Description string
	Status      string
	Priority    uint16
	CreatedAt   int64
	DueDate     int64
}

func CreateNewTask(task TaskRequestBody) (uint, error) {
	dueDate := time.Unix(*task.DueDate, 0)

	newTask := models.Task{
		Title:       *task.Title,
		Description: *task.Description,
		Status:      *task.Status,
		Priority:    uint16(*task.Priority),
		CreatedAt:   time.Now(),
		DueDate:     dueDate,
	}

	newTaskId, err := models.AddTask(newTask)

	if err != nil {
		fmt.Printf("Create Task failed: %v\n", err)
		return 0, ErrDatabaseGeneral
	}

	return newTaskId, nil
}

func GetTaskById(taskId uint) (TaskResponseBody, error) {

	idExist, err := checkIdExist(taskId)
	if !idExist {
		return TaskResponseBody{}, ErrRowNotFound
	} else if err != nil {
		return TaskResponseBody{}, ErrDatabaseGeneral
	}

	task, err := models.QueryTask(taskId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return TaskResponseBody{}, ErrRowNotFound
		} else {
			fmt.Printf("Query Task failed: %v\n", err)
			return TaskResponseBody{}, ErrDatabaseGeneral
		}
	}

	return TaskResponseBody{
		Id:          task.Id,
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    task.Priority,
		CreatedAt:   task.CreatedAt.Unix(),
		DueDate:     task.DueDate.Unix(),
	}, nil

}

func UpdateTask(taskId uint, task TaskRequestBody) error {
	currentTask, err := models.QueryTask(taskId)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrRowNotFound
		} else {
			fmt.Printf("Query Task failed: %v\n", err)
			return ErrDatabaseGeneral
		}
	}

	if task.Title != nil {
		currentTask.Title = *task.Title
	}
	if task.Description != nil {
		currentTask.Description = *task.Description
	}
	if task.Priority != nil {
		currentTask.Priority = uint16(*task.Priority)
	}
	if task.Status != nil {
		currentTask.Status = *task.Status
	}
	if task.DueDate != nil {
		// Convert Unix timestamp (seconds) to time.Time
		currentTask.DueDate = time.Unix(*task.DueDate, 0)
	}

	err = models.UpdateTask(currentTask)
	if err != nil {
		fmt.Printf("Update Task failed: %v\n", err)
		return ErrDatabaseGeneral
	}

	return nil
}

func DeleteTask(taskId uint) error {

	idExist, err := checkIdExist(taskId)
	if !idExist {
		return ErrRowNotFound
	} else if err != nil {
		return ErrDatabaseGeneral
	}

	err = models.DeleteTask(taskId)
	if err != nil {
		fmt.Printf("Delete Task failed: %v\n", err)
		return ErrDatabaseGeneral
	}

	return nil
}

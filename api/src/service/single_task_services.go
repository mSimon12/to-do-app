package service

import (
	"database/sql"
	"errors"
	"fmt"
	"time"
	"to-do-api/models"
)

type TaskRequestBody struct {
	Title       *string `json:"title"`
	Priority    *uint   `json:"priority"`
	Description *string `json:"description"`
	Status      *string `json:"status"`
	DueDate     *string `json:"due_date"`
}

var ErrDatabaseGeneral = errors.New("fail processing request on database")
var ErrRowNotFound = errors.New("requested resource not found on database")

func dateStrToTime(date string) (time.Time, error) {
	layout := "2006-01-02"
	convertedDate, err := time.Parse(layout, date)

	return convertedDate, err
}

func CreateNewTask(task TaskRequestBody) (uint, error) {
	dueDate, err := dateStrToTime(*task.DueDate)
	if err != nil {
		return 0, errors.New("invalid due_date format, expects: 'yyyy-mm-dd'")
	}

	if task.DueDate == nil {
		dueDate = time.Now().AddDate(0, 0, 7)
	}

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
			fmt.Printf("Query Task failed: %v\n", err)
			return task, ErrDatabaseGeneral
		}
	}

	return task, nil
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
		dueDate, err := dateStrToTime(*task.DueDate)
		if err != nil {
			return errors.New("invalid due_date format, expects: 'yyyy-mm-dd'")
		}
		currentTask.DueDate = dueDate
	}

	err = models.UpdateTask(currentTask)
	if err != nil {
		fmt.Printf("Update Task failed: %v\n", err)
		return ErrDatabaseGeneral
	}

	return nil
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

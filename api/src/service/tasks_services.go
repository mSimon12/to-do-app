package service

import (
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

func CreateNewTask(task TaskRequestBody) uint16 {
	newTask := models.Task{
		Title:       task.Title,
		Description: task.Description,
		Status:      task.Status,
		Priority:    uint16(task.Priority),
		CreatedAt:   time.Now(),
		DueDate:     dateStrToTime(task.DueDate),
	}

	newTaskId := models.CreateTask(newTask)
	return newTaskId
}

func dateStrToTime(date string) time.Time {
	layout := "2006-01-22"
	convertedDate, err := time.Parse(layout, date)

	if err != nil {
		fmt.Println(err)
		return time.Now().AddDate(0, 0, 1)
	}

	return convertedDate
}

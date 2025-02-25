package main

import (
	"time"
	"to-do-api/models"
)

func main() {
	//controllers.StartAPI()
	// models.InitDatabase()
	newTask := models.Task{Title: "task number 1", Description: "my description", Status: "backlog", Priority: 1, DueDate: time.Now().AddDate(0, 0, 7)}
	models.CreateTask(newTask)
}

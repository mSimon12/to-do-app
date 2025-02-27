package service

import (
	"errors"
	"fmt"
	"strconv"
	"to-do-api/models"
)

func ValidateNewTaskInput(requestInput TaskRequestBody) error {
	if requestInput.Title == nil {
		return errors.New("missing required field: 'title'")
	}
	return nil
}

func ValidateUpdateTaskInput(requestInput TaskRequestBody) error {

	if (TaskRequestBody{}) == requestInput {
		return errors.New("at least one field must be present: 'title', 'description', 'priority', 'status, 'due_date'")
	}
	return nil
}

func ValidateTaskIdInput(taskIdString string) (uint, error) {
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		fmt.Printf("failed in Id type conversion: %v\n", err)
		return 0, errors.New("invalid task id")
	}

	return uint(taskId), nil

}

func checkIdExist(taskId uint) bool {
	validId, _ := models.CheckExistence(taskId)

	return validId
}

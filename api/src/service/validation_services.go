package service

import (
	"errors"
	"fmt"
	"slices"
	"strconv"
	"strings"
	"to-do-api/models"
)

var validSortCriteria []string = []string{"id", "title", "status", "priority", "created_at", "due_date"}

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

func isValidPageConfig(strConfig string) (uint, bool) {
	config_int, err := strconv.Atoi(strConfig)
	if (err != nil) || (config_int < 0) {
		return uint(config_int), false
	}

	return uint(config_int), true
}

func isValidSortCriteria(criteria string) bool {
	return slices.Contains(validSortCriteria, strings.ToLower(criteria))
}

func isValidSortOrder(order string) bool {
	validOrder := []string{"asc", "desc"}
	return slices.Contains(validOrder, strings.ToLower(order))
}

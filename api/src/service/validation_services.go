package service

import (
	"errors"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"
	"to-do-api/models"
)

var validSortCriteria []string = []string{"id", "title", "status", "priority", "created_at", "due_date"}

func ValidateNewTaskInput(requestInput TaskRequestBody) error {
	if requestInput.Title == nil {
		return errors.New("missing required field: 'title'")
	} else if *requestInput.Title == "" {
		return errors.New("title must not be empty")
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
	if err != nil || taskId < 0 {
		fmt.Printf("failed in Id type conversion: %v\n", err)
		return 0, errors.New("invalid task id")
	}

	return uint(taskId), nil

}

func checkIdExist(taskId uint) (bool, error) {
	validId, err := models.CheckExistence(taskId)
	if err != nil {
		return false, ErrDatabaseGeneral
	}

	return validId, err
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

// Ensures only letters, numbers, spaces, and basic punctuation
func isValidTextFilter(input string) bool {
	re := regexp.MustCompile(`^[a-zA-Z0-9\s.,!?_-]+$`)
	return re.MatchString(input)
}

func isValidPriorityFilter(input string) bool {
	value, err := strconv.Atoi(input)
	if err != nil {
		return false
	}
	return value >= 0
}

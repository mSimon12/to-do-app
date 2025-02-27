package service

import (
	"errors"
	"fmt"
	"strconv"
)

func ValidateTaskId(taskIdString string) (uint, error) {
	taskId, err := strconv.Atoi(taskIdString)
	if err != nil {
		fmt.Printf("failed in Id type conversion: %v\n", err)
		return 0, errors.New("invalid task id")
	}

	return uint(taskId), nil

}

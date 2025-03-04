package service

import (
	"database/sql"
	"errors"
	"fmt"
	"to-do-api/models"
)

var defaultPageConfig models.TasksListQuery = models.TasksListQuery{Offset: 0, SortBy: "id", SortOrder: "ASC", Limit: 10}

type taskInfo struct {
	Id    uint
	Title string
}

func CreatePageConfig(offset string, limit string, sortBy string, sortOrder string) (models.TasksListQuery, error) {
	pageConfig := defaultPageConfig

	if offset != "" {
		offset_int, valid := isValidPageConfig(offset)
		if !valid {
			return pageConfig, errors.New("invalid 'offset' value, must be int > 0")
		}
		pageConfig.Offset = offset_int
	}

	if limit != "" {
		limit_int, valid := isValidPageConfig(limit)
		if !valid {
			return pageConfig, errors.New("invalid 'limit' value, must be int > 0")
		}
		pageConfig.Limit = limit_int
	}

	if sortBy != "" {
		if !isValidSortCriteria(sortBy) {
			err_string := fmt.Sprintf("invalid 'sortBy' value. Valid values: %v", validSortCriteria)
			return pageConfig, errors.New(err_string)
		}
		pageConfig.SortBy = sortBy
	}

	if sortOrder != "" {
		if !isValidSortOrder(sortOrder) {
			return pageConfig, errors.New("invalid 'sortOrder' value. Valid values: ['asc', 'desc']")
		}
		pageConfig.SortOrder = sortOrder
	}

	return pageConfig, nil
}

func GetTasksList(pageConfig models.TasksListQuery) (map[int]taskInfo, error) {
	orderedTasks := map[int]taskInfo{}

	queriedTasks, err := models.QueryTasks(pageConfig)

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

func GetReturnInfo(pageConfig models.TasksListQuery) (map[string]uint, map[string]string) {
	paginationInfo := map[string]uint{"offset": pageConfig.Offset, "limit": pageConfig.Limit}
	sortingInfo := map[string]string{"by": pageConfig.SortBy, "order": pageConfig.SortOrder}

	totalTasks, err := models.GetAmountOfTasks()
	if err != nil {
		fmt.Printf("Error getting total tasks. e: %v\n", err)
	} else {
		paginationInfo["total_tasks"] = totalTasks
	}

	return paginationInfo, sortingInfo
}

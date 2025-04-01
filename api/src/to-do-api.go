package main

import (
	"to-do-api/controllers"
	"to-do-api/models"
)

// @title			To-Do API
// @version		1.0
// @description	To-Do List API for managing tasks to be done. It provides the basic functions of Creating, Updating, Deleting and Querying single or multiple tasks.
// @termsOfService	http://swagger.io/terms/
func main() {
	models.InitDatabase()

	controllers.StartAPI()
}

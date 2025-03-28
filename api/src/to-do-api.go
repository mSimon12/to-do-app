package main

import (
	"to-do-api/controllers"
	"to-do-api/models"
)

func main() {
	models.InitDatabase()

	controllers.StartAPI()
}
